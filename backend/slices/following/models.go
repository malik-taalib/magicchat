package following

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Follow represents a follow relationship between two users
type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID  primitive.ObjectID `bson:"follower_id" json:"follower_id"`   // User who is following
	FollowingID primitive.ObjectID `bson:"following_id" json:"following_id"` // User being followed
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// UserProfile represents a simplified user profile for following lists
type UserProfile struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Username       string             `bson:"username" json:"username"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	Bio            string             `bson:"bio" json:"bio"`
	AvatarURL      string             `bson:"avatar_url" json:"avatar_url"`
	FollowerCount  int                `bson:"follower_count" json:"follower_count"`
	FollowingCount int                `bson:"following_count" json:"following_count"`
	VideoCount     int                `bson:"video_count" json:"video_count"`
	TotalLikes     int                `bson:"total_likes" json:"total_likes"`
}

// FollowResponse represents the response after follow/unfollow action
type FollowResponse struct {
	Success        bool   `json:"success"`
	IsFollowing    bool   `json:"is_following"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
	Message        string `json:"message,omitempty"`
}

// FollowListResponse represents a paginated list of user profiles
type FollowListResponse struct {
	Users      []UserProfile `json:"users"`
	NextCursor string        `json:"next_cursor,omitempty"`
	HasMore    bool          `json:"has_more"`
	Total      int           `json:"total"`
}

// Like represents a user liking a video
type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	VideoID   primitive.ObjectID `bson:"video_id" json:"video_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// FeedVideo represents a video with all its metadata for the feed
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

// LikedVideosResponse represents the response for liked videos
type LikedVideosResponse struct {
	Videos []*FeedVideo `json:"videos"`
	Limit  int64        `json:"limit"`
	Offset int64        `json:"offset"`
	Total  int          `json:"total"`
}
