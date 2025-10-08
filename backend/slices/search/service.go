package search

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Search performs a search based on the search type
func (s *Service) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// Validate request
	if err := s.validateSearchRequest(req); err != nil {
		return nil, err
	}

	// Normalize query (trim spaces)
	req.Query = strings.TrimSpace(req.Query)

	response := &SearchResponse{}

	switch req.Type {
	case SearchTypeUsers:
		users, err := s.repo.SearchUsers(ctx, req.Query, req.Cursor, req.Limit)
		if err != nil {
			return nil, err
		}
		response.Users = users
		response.HasMore = len(users) == req.Limit
		if response.HasMore && len(users) > 0 {
			response.NextCursor = users[len(users)-1].ID.Hex()
		}

	case SearchTypeVideos:
		videos, err := s.repo.SearchVideos(ctx, req.Query, req.Cursor, req.Limit)
		if err != nil {
			return nil, err
		}
		response.Videos = videos
		response.HasMore = len(videos) == req.Limit
		if response.HasMore && len(videos) > 0 {
			response.NextCursor = videos[len(videos)-1].ID.Hex()
		}

	case SearchTypeHashtags:
		hashtags, err := s.repo.SearchHashtags(ctx, req.Query, req.Limit)
		if err != nil {
			return nil, err
		}
		response.Hashtags = hashtags
		response.HasMore = false // Hashtag search doesn't use cursor pagination

	default:
		return nil, errors.New("invalid search type")
	}

	return response, nil
}

// GetTrendingHashtags returns the top trending hashtags
func (s *Service) GetTrendingHashtags(ctx context.Context, limit int) (*TrendingHashtagsResponse, error) {
	// Validate limit
	if limit <= 0 {
		limit = 20 // Default
	}
	if limit > 50 {
		limit = 50 // Max
	}

	hashtags, err := s.repo.GetTrendingHashtags(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &TrendingHashtagsResponse{
		Hashtags: hashtags,
	}, nil
}

// GetVideosByHashtag returns videos for a specific hashtag
func (s *Service) GetVideosByHashtag(ctx context.Context, req *HashtagVideosRequest) (*HashtagVideosResponse, error) {
	// Validate request
	if err := s.validateHashtagVideosRequest(req); err != nil {
		return nil, err
	}

	// Normalize tag (trim spaces, remove # if present, lowercase)
	req.Tag = strings.TrimSpace(req.Tag)
	req.Tag = strings.TrimPrefix(req.Tag, "#")
	req.Tag = strings.ToLower(req.Tag)

	videos, err := s.repo.GetVideosByHashtag(ctx, req.Tag, req.Cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	response := &HashtagVideosResponse{
		Tag:     req.Tag,
		Videos:  videos,
		HasMore: len(videos) == req.Limit,
	}

	if response.HasMore && len(videos) > 0 {
		response.NextCursor = videos[len(videos)-1].ID.Hex()
	}

	return response, nil
}

// Validation helpers

func (s *Service) validateSearchRequest(req *SearchRequest) error {
	// Validate query
	if strings.TrimSpace(req.Query) == "" {
		return errors.New("query cannot be empty")
	}

	if len(req.Query) < 2 {
		return errors.New("query must be at least 2 characters")
	}

	if len(req.Query) > 100 {
		return errors.New("query must be at most 100 characters")
	}

	// Validate search type
	if req.Type == "" {
		return errors.New("search type is required")
	}

	validTypes := map[SearchType]bool{
		SearchTypeUsers:    true,
		SearchTypeVideos:   true,
		SearchTypeHashtags: true,
	}

	if !validTypes[req.Type] {
		return errors.New("invalid search type: must be users, videos, or hashtags")
	}

	// Validate limit
	if req.Limit <= 0 {
		req.Limit = 20 // Default
	}
	if req.Limit > 50 {
		req.Limit = 50 // Max
	}

	return nil
}

func (s *Service) validateHashtagVideosRequest(req *HashtagVideosRequest) error {
	// Validate tag
	if strings.TrimSpace(req.Tag) == "" {
		return errors.New("hashtag tag cannot be empty")
	}

	tag := strings.TrimPrefix(strings.TrimSpace(req.Tag), "#")
	if len(tag) < 1 {
		return errors.New("hashtag tag must be at least 1 character")
	}

	if len(tag) > 50 {
		return errors.New("hashtag tag must be at most 50 characters")
	}

	// Validate limit
	if req.Limit <= 0 {
		req.Limit = 20 // Default
	}
	if req.Limit > 50 {
		req.Limit = 50 // Max
	}

	return nil
}
