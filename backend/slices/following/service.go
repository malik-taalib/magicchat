package following

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// FollowUser allows a user to follow another user
func (s *Service) FollowUser(ctx context.Context, followerID, followingID string) (*FollowResponse, error) {
	// Validate IDs
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return nil, errors.New("invalid follower ID")
	}

	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return nil, errors.New("invalid following ID")
	}

	// Check if user is trying to follow themselves
	if followerID == followingID {
		return nil, errors.New("cannot follow yourself")
	}

	// Check if the user to be followed exists
	_, err = s.repo.GetUserByID(ctx, followingObjID)
	if err != nil {
		return nil, errors.New("user to follow not found")
	}

	// Create follow relationship
	err = s.repo.FollowUser(ctx, followerObjID, followingObjID)
	if err != nil {
		return nil, err
	}

	// Get updated counts
	followerCount, err := s.repo.GetFollowerCount(ctx, followingObjID)
	if err != nil {
		return nil, err
	}

	followingCount, err := s.repo.GetFollowingCount(ctx, followerObjID)
	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		Success:        true,
		IsFollowing:    true,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		Message:        "Successfully followed user",
	}, nil
}

// UnfollowUser allows a user to unfollow another user
func (s *Service) UnfollowUser(ctx context.Context, followerID, followingID string) (*FollowResponse, error) {
	// Validate IDs
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return nil, errors.New("invalid follower ID")
	}

	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return nil, errors.New("invalid following ID")
	}

	// Check if user is trying to unfollow themselves
	if followerID == followingID {
		return nil, errors.New("cannot unfollow yourself")
	}

	// Remove follow relationship
	err = s.repo.UnfollowUser(ctx, followerObjID, followingObjID)
	if err != nil {
		return nil, err
	}

	// Get updated counts
	followerCount, err := s.repo.GetFollowerCount(ctx, followingObjID)
	if err != nil {
		return nil, err
	}

	followingCount, err := s.repo.GetFollowingCount(ctx, followerObjID)
	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		Success:        true,
		IsFollowing:    false,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		Message:        "Successfully unfollowed user",
	}, nil
}

// IsFollowing checks if one user is following another
func (s *Service) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	// Validate IDs
	followerObjID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return false, errors.New("invalid follower ID")
	}

	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		return false, errors.New("invalid following ID")
	}

	return s.repo.IsFollowing(ctx, followerObjID, followingObjID)
}

// GetFollowers returns a paginated list of users who follow the specified user
func (s *Service) GetFollowers(ctx context.Context, userID, cursor string, limit int) (*FollowListResponse, error) {
	// Validate user ID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	// Check if user exists
	_, err = s.repo.GetUserByID(ctx, userObjID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get followers
	profiles, nextCursor, err := s.repo.GetFollowers(ctx, userObjID, cursor, limit)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.GetFollowerCount(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	return &FollowListResponse{
		Users:      profiles,
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}, nil
}

// GetFollowing returns a paginated list of users that the specified user follows
func (s *Service) GetFollowing(ctx context.Context, userID, cursor string, limit int) (*FollowListResponse, error) {
	// Validate user ID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	// Check if user exists
	_, err = s.repo.GetUserByID(ctx, userObjID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get following
	profiles, nextCursor, err := s.repo.GetFollowing(ctx, userObjID, cursor, limit)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.GetFollowingCount(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	return &FollowListResponse{
		Users:      profiles,
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}, nil
}

// GetUserProfile returns a user's profile with follow counts
func (s *Service) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	// Validate user ID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetUserByID(ctx, userObjID)
}
