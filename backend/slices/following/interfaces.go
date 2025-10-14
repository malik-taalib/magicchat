package following

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RepositoryInterface defines the contract for the repository layer
// This allows for easier testing with mock implementations
type RepositoryInterface interface {
	FollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error
	UnfollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error
	IsFollowing(ctx context.Context, followerID, followingID primitive.ObjectID) (bool, error)
	GetFollowers(ctx context.Context, userID primitive.ObjectID, cursor string, limit int) ([]UserProfile, string, error)
	GetFollowing(ctx context.Context, userID primitive.ObjectID, cursor string, limit int) ([]UserProfile, string, error)
	GetFollowerCount(ctx context.Context, userID primitive.ObjectID) (int, error)
	GetFollowingCount(ctx context.Context, userID primitive.ObjectID) (int, error)
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*UserProfile, error)
	GetLikedVideos(ctx context.Context, userID primitive.ObjectID, limit, offset int64) ([]*FeedVideo, error)
}
