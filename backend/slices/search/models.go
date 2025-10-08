package search

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchType represents the type of search to perform
type SearchType string

const (
	SearchTypeUsers    SearchType = "users"
	SearchTypeVideos   SearchType = "videos"
	SearchTypeHashtags SearchType = "hashtags"
)

// SearchRequest represents a search query
type SearchRequest struct {
	Query  string     `json:"query" form:"q"`
	Type   SearchType `json:"type" form:"type"`
	Cursor string     `json:"cursor" form:"cursor"` // ObjectID hex string for cursor-based pagination
	Limit  int        `json:"limit" form:"limit"`   // Number of results to return (default: 20, max: 50)
}

// SearchResponse represents the paginated search results
type SearchResponse struct {
	Users      []*UserSearchResult    `json:"users,omitempty"`
	Videos     []*VideoSearchResult   `json:"videos,omitempty"`
	Hashtags   []*HashtagSearchResult `json:"hashtags,omitempty"`
	NextCursor string                 `json:"next_cursor,omitempty"` // Empty if no more results
	HasMore    bool                   `json:"has_more"`
}

// UserSearchResult represents a user in search results
type UserSearchResult struct {
	ID             primitive.ObjectID `json:"id"`
	Username       string             `json:"username"`
	DisplayName    string             `json:"display_name"`
	Bio            string             `json:"bio"`
	AvatarURL      string             `json:"avatar_url"`
	FollowerCount  int                `json:"follower_count"`
	VideoCount     int                `json:"video_count"`
	IsVerified     bool               `json:"is_verified"`
}

// VideoSearchResult represents a video in search results
type VideoSearchResult struct {
	ID           primitive.ObjectID `json:"id"`
	UserID       primitive.ObjectID `json:"user_id"`
	Username     string             `json:"username"`
	DisplayName  string             `json:"display_name"`
	AvatarURL    string             `json:"avatar_url"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	VideoURL     string             `json:"video_url"`
	ThumbnailURL string             `json:"thumbnail_url"`
	Duration     int                `json:"duration"`
	Hashtags     []string           `json:"hashtags"`
	ViewCount    int                `json:"view_count"`
	LikeCount    int                `json:"like_count"`
	CommentCount int                `json:"comment_count"`
	CreatedAt    time.Time          `json:"created_at"`
}

// HashtagSearchResult represents a hashtag in search results
type HashtagSearchResult struct {
	Tag           string    `json:"tag"`
	VideoCount    int       `json:"video_count"`
	TrendingScore float64   `json:"trending_score"`
	LastUsed      time.Time `json:"last_used"`
}

// Hashtag represents a hashtag document in the database
type Hashtag struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Tag           string             `bson:"tag" json:"tag"`
	VideoCount    int                `bson:"video_count" json:"video_count"`
	TrendingScore float64            `bson:"trending_score" json:"trending_score"`
	LastUsed      time.Time          `bson:"last_used" json:"last_used"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// TrendingHashtagsResponse represents the response for trending hashtags
type TrendingHashtagsResponse struct {
	Hashtags []*HashtagSearchResult `json:"hashtags"`
}

// HashtagVideosRequest represents pagination parameters for hashtag videos
type HashtagVideosRequest struct {
	Tag    string `json:"tag" form:"tag"`
	Cursor string `json:"cursor" form:"cursor"` // ObjectID hex string for cursor-based pagination
	Limit  int    `json:"limit" form:"limit"`   // Number of videos to return (default: 20, max: 50)
}

// HashtagVideosResponse represents videos for a specific hashtag
type HashtagVideosResponse struct {
	Tag        string               `json:"tag"`
	Videos     []*VideoSearchResult `json:"videos"`
	NextCursor string               `json:"next_cursor,omitempty"` // Empty if no more videos
	HasMore    bool                 `json:"has_more"`
}
