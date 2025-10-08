package notifications

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("notifications"),
	}
}

// CreateNotification creates a new notification
func (r *Repository) CreateNotification(ctx context.Context, notification *Notification) error {
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	notification.Read = false

	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

// GetNotifications retrieves notifications for a user with pagination
func (r *Repository) GetNotifications(ctx context.Context, userID string, cursor string, limit int) ([]*Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Build filter
	filter := bson.M{"user_id": objectID}

	// Add cursor-based pagination
	if cursor != "" {
		cursorObjectID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		filter["_id"] = bson.M{"$lt": cursorObjectID}
	}

	// Query options - sort by created_at descending (newest first)
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit))

	cursor2, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)

	var notifications []*Notification
	if err = cursor2.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetNotificationByID retrieves a single notification by ID
func (r *Repository) GetNotificationByID(ctx context.Context, notificationID string) (*Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return nil, err
	}

	var notification Notification
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}

	return &notification, nil
}

// MarkAsRead marks a single notification as read
func (r *Repository) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	notifObjectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Update only if the notification belongs to the user
	filter := bson.M{
		"_id":     notifObjectID,
		"user_id": userObjectID,
	}
	update := bson.M{
		"$set": bson.M{
			"read": true,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("notification not found or unauthorized")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *Repository) MarkAllAsRead(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"user_id": objectID,
		"read":    false,
	}
	update := bson.M{
		"$set": bson.M{
			"read": true,
		},
	}

	_, err = r.collection.UpdateMany(ctx, filter, update)
	return err
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *Repository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{
		"user_id": objectID,
		"read":    false,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// DeleteNotification deletes a notification (optional utility method)
func (r *Repository) DeleteNotification(ctx context.Context, notificationID string, userID string) error {
	notifObjectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     notifObjectID,
		"user_id": userObjectID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("notification not found or unauthorized")
	}

	return nil
}

// CheckDuplicateNotification checks if a similar notification already exists
// This prevents duplicate notifications for the same action
func (r *Repository) CheckDuplicateNotification(ctx context.Context, userID, actorID primitive.ObjectID, notifType NotificationType, videoID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"user_id":  userID,
		"actor_id": actorID,
		"type":     notifType,
	}

	if videoID != nil {
		filter["video_id"] = videoID
	}

	// Check if notification was created in the last 24 hours
	filter["created_at"] = bson.M{
		"$gte": time.Now().Add(-24 * time.Hour),
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
