package notifications

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeMention NotificationType = "mention"
)

// Notification represents a user notification
type Notification struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"user_id"`                   // Recipient of the notification
	Type      NotificationType    `bson:"type" json:"type"`                         // Type of notification
	ActorID   primitive.ObjectID  `bson:"actor_id" json:"actor_id"`                 // User who triggered the notification
	VideoID   *primitive.ObjectID `bson:"video_id,omitempty" json:"video_id,omitempty"` // Related video (if applicable)
	CommentID *primitive.ObjectID `bson:"comment_id,omitempty" json:"comment_id,omitempty"` // Related comment (if applicable)
	Text      string              `bson:"text" json:"text"`                         // Notification text/message
	Read      bool                `bson:"read" json:"read"`                         // Whether the notification has been read
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
}

// NotificationResponse represents the API response for a notification
type NotificationResponse struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	ActorID   string           `json:"actor_id"`
	VideoID   string           `json:"video_id,omitempty"`
	CommentID string           `json:"comment_id,omitempty"`
	Text      string           `json:"text"`
	Read      bool             `json:"read"`
	CreatedAt time.Time        `json:"created_at"`
	Actor     *ActorInfo       `json:"actor,omitempty"` // Actor details for UI
}

// ActorInfo contains information about the user who triggered the notification
type ActorInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// NotificationListResponse represents paginated notification list
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int                    `json:"unread_count"`
	HasMore       bool                   `json:"has_more"`
	NextCursor    string                 `json:"next_cursor,omitempty"`
}

// MarkAsReadRequest represents the request to mark notifications as read
type MarkAsReadRequest struct {
	NotificationID string `json:"notification_id"`
}
