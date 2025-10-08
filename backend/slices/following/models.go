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
	AvatarURL      string             `bson:"avatar_url" json:"avatar_url"`
	FollowerCount  int                `bson:"follower_count" json:"follower_count"`
	FollowingCount int                `bson:"following_count" json:"following_count"`
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
