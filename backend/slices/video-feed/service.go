package videofeed

import (
	"context"
	"errors"
)

const (
	DefaultFeedLimit = 10
	MaxFeedLimit     = 50
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetForYouFeed retrieves the personalized "For You" feed for a user
// This feed shows videos from all users, ranked by engagement metrics
func (s *Service) GetForYouFeed(ctx context.Context, userID string, req *FeedRequest) (*FeedResponse, error) {
	// Validate and set defaults for pagination
	limit := s.validateLimit(req.Limit)

	// Fetch one extra video to determine if there are more results
	videos, err := s.repo.GetForYouFeed(ctx, userID, req.Cursor, limit+1)
	if err != nil {
		return nil, err
	}

	// Build response with pagination metadata
	return s.buildFeedResponse(videos, limit), nil
}

// GetFollowingFeed retrieves videos from users that the current user follows
// This feed shows only videos from followed users, sorted by recency
func (s *Service) GetFollowingFeed(ctx context.Context, userID string, req *FeedRequest) (*FeedResponse, error) {
	// Validate and set defaults for pagination
	limit := s.validateLimit(req.Limit)

	// Fetch one extra video to determine if there are more results
	videos, err := s.repo.GetFollowingFeed(ctx, userID, req.Cursor, limit+1)
	if err != nil {
		return nil, err
	}

	// Build response with pagination metadata
	return s.buildFeedResponse(videos, limit), nil
}

// GetVideoByID retrieves a single video by its ID
// This also increments the view count when called
func (s *Service) GetVideoByID(ctx context.Context, videoID string) (*FeedVideo, error) {
	if videoID == "" {
		return nil, errors.New("video ID is required")
	}

	// Get video with user information
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously (fire and forget)
	// In production, this might be done via a message queue to avoid blocking
	go func() {
		// Create a new context for the background operation
		bgCtx := context.Background()
		_ = s.repo.IncrementViewCount(bgCtx, videoID)
	}()

	return video, nil
}

// IncrementViewCount explicitly increments the view count for a video
func (s *Service) IncrementViewCount(ctx context.Context, videoID string) error {
	if videoID == "" {
		return errors.New("video ID is required")
	}

	return s.repo.IncrementViewCount(ctx, videoID)
}

// validateLimit ensures the limit is within acceptable bounds
func (s *Service) validateLimit(limit int) int {
	if limit <= 0 {
		return DefaultFeedLimit
	}
	if limit > MaxFeedLimit {
		return MaxFeedLimit
	}
	return limit
}

// buildFeedResponse constructs a FeedResponse with pagination metadata
func (s *Service) buildFeedResponse(videos []*FeedVideo, limit int) *FeedResponse {
	hasMore := len(videos) > limit

	// If we have more results than requested, trim to limit
	if hasMore {
		videos = videos[:limit]
	}

	// Build response
	response := &FeedResponse{
		Videos:  videos,
		HasMore: hasMore,
	}

	// Set next cursor if there are more results
	if hasMore && len(videos) > 0 {
		lastVideo := videos[len(videos)-1]
		response.NextCursor = lastVideo.ID.Hex()
	}

	return response
}
