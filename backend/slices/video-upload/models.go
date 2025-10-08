package videoupload

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)

type Video struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title            string             `bson:"title" json:"title"`
	Description      string             `bson:"description" json:"description"`
	VideoURL         string             `bson:"video_url" json:"video_url"`
	ThumbnailURL     string             `bson:"thumbnail_url" json:"thumbnail_url"`
	Duration         int                `bson:"duration" json:"duration"` // in seconds
	Hashtags         []string           `bson:"hashtags" json:"hashtags"`
	ViewCount        int                `bson:"view_count" json:"view_count"`
	LikeCount        int                `bson:"like_count" json:"like_count"`
	CommentCount     int                `bson:"comment_count" json:"comment_count"`
	ShareCount       int                `bson:"share_count" json:"share_count"`
	ProcessingStatus ProcessingStatus   `bson:"processing_status" json:"processing_status"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type UploadRequest struct {
	Title       string   `form:"title" binding:"required"`
	Description string   `form:"description"`
	Hashtags    []string `form:"hashtags"`
}

type UploadResponse struct {
	VideoID string           `json:"video_id"`
	Status  ProcessingStatus `json:"status"`
}

type StatusResponse struct {
	VideoID string           `json:"video_id"`
	Status  ProcessingStatus `json:"status"`
	VideoURL string          `json:"video_url,omitempty"`
}
