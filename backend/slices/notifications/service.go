package notifications

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo      *Repository
	usersColl *mongo.Collection // To fetch actor details
	wsManager *WebSocketManager
}

func NewService(repo *Repository, db *mongo.Database, wsManager *WebSocketManager) *Service {
	return &Service{
		repo:      repo,
		usersColl: db.Collection("users"),
		wsManager: wsManager,
	}
}

// CreateNotification creates a notification and broadcasts it via WebSocket
func (s *Service) CreateNotification(ctx context.Context, notification *Notification) error {
	// Check for duplicate notifications (optional, prevents spam)
	duplicate, err := s.repo.CheckDuplicateNotification(
		ctx,
		notification.UserID,
		notification.ActorID,
		notification.Type,
		notification.VideoID,
	)
	if err != nil {
		log.Printf("Error checking duplicate notification: %v", err)
	}
	if duplicate {
		return nil // Silent skip for duplicate notifications
	}

	// Create the notification
	err = s.repo.CreateNotification(ctx, notification)
	if err != nil {
		return err
	}

	// Fetch actor details for the response
	response, err := s.buildNotificationResponse(ctx, notification)
	if err != nil {
		log.Printf("Error building notification response: %v", err)
		// Don't fail the entire operation if we can't build the response
		return nil
	}

	// Broadcast via WebSocket to the recipient
	s.wsManager.BroadcastToUser(notification.UserID.Hex(), response)

	return nil
}

// GetNotifications retrieves notifications for a user
func (s *Service) GetNotifications(ctx context.Context, userID string, cursor string, limit int) (*NotificationListResponse, error) {
	// Get notifications
	notifications, err := s.repo.GetNotifications(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}

	// Build responses with actor details
	responses := make([]NotificationResponse, 0, len(notifications))
	for _, notif := range notifications {
		response, err := s.buildNotificationResponse(ctx, notif)
		if err != nil {
			log.Printf("Error building notification response: %v", err)
			continue
		}
		responses = append(responses, *response)
	}

	// Get unread count
	unreadCount, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		log.Printf("Error getting unread count: %v", err)
		unreadCount = 0
	}

	// Determine if there are more notifications
	hasMore := len(notifications) == limit
	var nextCursor string
	if hasMore && len(notifications) > 0 {
		nextCursor = notifications[len(notifications)-1].ID.Hex()
	}

	return &NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unreadCount,
		HasMore:       hasMore,
		NextCursor:    nextCursor,
	}, nil
}

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// GetUnreadCount returns the count of unread notifications
func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

// buildNotificationResponse builds a NotificationResponse with actor details
func (s *Service) buildNotificationResponse(ctx context.Context, notification *Notification) (*NotificationResponse, error) {
	response := &NotificationResponse{
		ID:        notification.ID.Hex(),
		UserID:    notification.UserID.Hex(),
		Type:      notification.Type,
		ActorID:   notification.ActorID.Hex(),
		Text:      notification.Text,
		Read:      notification.Read,
		CreatedAt: notification.CreatedAt,
	}

	// Add optional fields
	if notification.VideoID != nil {
		response.VideoID = notification.VideoID.Hex()
	}
	if notification.CommentID != nil {
		response.CommentID = notification.CommentID.Hex()
	}

	// Fetch actor details
	actor, err := s.getActorInfo(ctx, notification.ActorID)
	if err != nil {
		log.Printf("Error fetching actor info: %v", err)
		// Don't fail, just leave actor as nil
	} else {
		response.Actor = actor
	}

	return response, nil
}

// getActorInfo fetches user information for the actor
func (s *Service) getActorInfo(ctx context.Context, actorID primitive.ObjectID) (*ActorInfo, error) {
	var user struct {
		ID          primitive.ObjectID `bson:"_id"`
		Username    string             `bson:"username"`
		DisplayName string             `bson:"display_name"`
		AvatarURL   string             `bson:"avatar_url"`
	}

	err := s.usersColl.FindOne(ctx, primitive.M{"_id": actorID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &ActorInfo{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
	}, nil
}

// Helper functions to create notifications easily

// NotifyLike creates a notification for a video like
func (s *Service) NotifyLike(ctx context.Context, videoOwnerID, actorID, videoID string) error {
	videoOwnerObjID, err := primitive.ObjectIDFromHex(videoOwnerID)
	if err != nil {
		return err
	}
	actorObjID, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return err
	}
	videoObjID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return err
	}

	// Don't notify if user liked their own video
	if videoOwnerID == actorID {
		return nil
	}

	notification := &Notification{
		UserID:  videoOwnerObjID,
		Type:    NotificationTypeLike,
		ActorID: actorObjID,
		VideoID: &videoObjID,
		Text:    "liked your video",
	}

	return s.CreateNotification(ctx, notification)
}

// NotifyComment creates a notification for a video comment
func (s *Service) NotifyComment(ctx context.Context, videoOwnerID, actorID, videoID, commentID string, commentText string) error {
	videoOwnerObjID, err := primitive.ObjectIDFromHex(videoOwnerID)
	if err != nil {
		return err
	}
	actorObjID, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return err
	}
	videoObjID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	// Don't notify if user commented on their own video
	if videoOwnerID == actorID {
		return nil
	}

	// Truncate comment text if too long
	displayText := commentText
	if len(displayText) > 50 {
		displayText = displayText[:50] + "..."
	}

	notification := &Notification{
		UserID:    videoOwnerObjID,
		Type:      NotificationTypeComment,
		ActorID:   actorObjID,
		VideoID:   &videoObjID,
		CommentID: &commentObjID,
		Text:      fmt.Sprintf("commented: %s", displayText),
	}

	return s.CreateNotification(ctx, notification)
}

// NotifyFollow creates a notification for a new follower
func (s *Service) NotifyFollow(ctx context.Context, followedUserID, followerID string) error {
	followedObjID, err := primitive.ObjectIDFromHex(followedUserID)
	if err != nil {
		return err
	}
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return err
	}

	notification := &Notification{
		UserID:  followedObjID,
		Type:    NotificationTypeFollow,
		ActorID: followerObjID,
		Text:    "started following you",
	}

	return s.CreateNotification(ctx, notification)
}

// NotifyMention creates a notification for a mention in a comment
func (s *Service) NotifyMention(ctx context.Context, mentionedUserID, actorID, videoID, commentID string) error {
	mentionedObjID, err := primitive.ObjectIDFromHex(mentionedUserID)
	if err != nil {
		return err
	}
	actorObjID, err := primitive.ObjectIDFromHex(actorID)
	if err != nil {
		return err
	}
	videoObjID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return err
	}
	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	notification := &Notification{
		UserID:    mentionedObjID,
		Type:      NotificationTypeMention,
		ActorID:   actorObjID,
		VideoID:   &videoObjID,
		CommentID: &commentObjID,
		Text:      "mentioned you in a comment",
	}

	return s.CreateNotification(ctx, notification)
}
