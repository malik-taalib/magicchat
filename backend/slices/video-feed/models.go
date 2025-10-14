package videofeed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FeedVideo combines video information with user details for feed display
type FeedVideo struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username         string             `bson:"username" json:"username"`
	DisplayName      string             `bson:"display_name" json:"display_name"`
	AvatarURL        string             `bson:"avatar_url" json:"avatar_url"`
	Title            string             `bson:"title" json:"title"`
	Description      string             `bson:"description" json:"description"`
	VideoURL         string             `bson:"video_url" json:"video_url"`
	ThumbnailURL     string             `bson:"thumbnail_url" json:"thumbnail_url"`
	Duration         int                `bson:"duration" json:"duration"`
	Hashtags         []string           `bson:"hashtags" json:"hashtags"`
	ViewCount        int                `bson:"view_count" json:"view_count"`
	LikeCount        int                `bson:"like_count" json:"like_count"`
	CommentCount     int                `bson:"comment_count" json:"comment_count"`
	ShareCount       int                `bson:"share_count" json:"share_count"`
	ProcessingStatus string             `bson:"processing_status" json:"processing_status"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// FeedRequest represents pagination parameters for feed requests
type FeedRequest struct {
	Cursor string `json:"cursor" form:"cursor"` // ObjectID hex string for cursor-based pagination
	Limit  int    `json:"limit" form:"limit"`   // Number of videos to return (default: 10, max: 50)
}

// FeedResponse represents the paginated feed response
type FeedResponse struct {
	Videos     []*FeedVideo `json:"videos"`
	NextCursor string       `json:"next_cursor,omitempty"` // Empty if no more videos
	HasMore    bool         `json:"has_more"`
}

// VideoResponse represents a single video response
type VideoResponse struct {
	Video *FeedVideo `json:"video"`
}
