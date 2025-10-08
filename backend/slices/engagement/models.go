package engagement

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like represents a video like
type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	VideoID   primitive.ObjectID `bson:"video_id" json:"video_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Comment represents a video comment with support for nested replies
type Comment struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID   `bson:"user_id" json:"user_id"`
	VideoID   primitive.ObjectID   `bson:"video_id" json:"video_id"`
	Text      string               `bson:"text" json:"text"`
	Replies   []primitive.ObjectID `bson:"replies" json:"replies"` // Array of comment IDs
	ParentID  *primitive.ObjectID  `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

// Share represents a video share action
type Share struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	VideoID   primitive.ObjectID `bson:"video_id" json:"video_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Request/Response models

type LikeResponse struct {
	VideoID   string `json:"video_id"`
	Liked     bool   `json:"liked"`
	LikeCount int    `json:"like_count"`
}

type CreateCommentRequest struct {
	Text     string `json:"text" binding:"required"`
	ParentID string `json:"parent_id,omitempty"` // Optional, for replies
}

type CommentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	VideoID   string    `json:"video_id"`
	Text      string    `json:"text"`
	ParentID  string    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ShareResponse struct {
	VideoID    string `json:"video_id"`
	ShareCount int    `json:"share_count"`
}
