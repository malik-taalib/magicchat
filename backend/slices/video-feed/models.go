package videofeed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FeedVideo combines video information with user details for feed display
type FeedVideo struct {
	ID               primitive.ObjectID `json:"id"`
	UserID           primitive.ObjectID `json:"user_id"`
	Username         string             `json:"username"`
	DisplayName      string             `json:"display_name"`
	AvatarURL        string             `json:"avatar_url"`
	Title            string             `json:"title"`
	Description      string             `json:"description"`
	VideoURL         string             `json:"video_url"`
	ThumbnailURL     string             `json:"thumbnail_url"`
	Duration         int                `json:"duration"`
	Hashtags         []string           `json:"hashtags"`
	ViewCount        int                `json:"view_count"`
	LikeCount        int                `json:"like_count"`
	CommentCount     int                `json:"comment_count"`
	ShareCount       int                `json:"share_count"`
	ProcessingStatus string             `json:"processing_status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
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
