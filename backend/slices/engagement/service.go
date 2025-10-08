package engagement

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Like operations

func (s *Service) LikeVideo(ctx context.Context, userID, videoID string) (*LikeResponse, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	// Like the video
	err = s.repo.LikeVideo(ctx, userObjectID, videoObjectID)
	if err != nil {
		return nil, err
	}

	// Get updated like count
	stats, err := s.repo.GetVideoStats(ctx, videoObjectID)
	if err != nil {
		return nil, err
	}

	return &LikeResponse{
		VideoID:   videoID,
		Liked:     true,
		LikeCount: stats["like_count"],
	}, nil
}

func (s *Service) UnlikeVideo(ctx context.Context, userID, videoID string) (*LikeResponse, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	// Unlike the video
	err = s.repo.UnlikeVideo(ctx, userObjectID, videoObjectID)
	if err != nil {
		return nil, err
	}

	// Get updated like count
	stats, err := s.repo.GetVideoStats(ctx, videoObjectID)
	if err != nil {
		return nil, err
	}

	return &LikeResponse{
		VideoID:   videoID,
		Liked:     false,
		LikeCount: stats["like_count"],
	}, nil
}

func (s *Service) IsVideoLiked(ctx context.Context, userID, videoID string) (bool, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID")
	}

	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return false, errors.New("invalid video ID")
	}

	return s.repo.IsVideoLiked(ctx, userObjectID, videoObjectID)
}

// Comment operations

func (s *Service) CreateComment(ctx context.Context, userID, videoID string, req *CreateCommentRequest) (*CommentResponse, error) {
	// Validate input
	if err := s.validateComment(req); err != nil {
		return nil, err
	}

	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	comment := &Comment{
		UserID:  userObjectID,
		VideoID: videoObjectID,
		Text:    strings.TrimSpace(req.Text),
	}

	// Handle parent ID for replies
	if req.ParentID != "" {
		parentObjectID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent comment ID")
		}

		// Verify parent comment exists and belongs to the same video
		parentComment, err := s.repo.GetCommentByID(ctx, parentObjectID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}

		if parentComment.VideoID != videoObjectID {
			return nil, errors.New("parent comment does not belong to this video")
		}

		comment.ParentID = &parentObjectID
	}

	// Create comment
	err = s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	response := &CommentResponse{
		ID:        comment.ID.Hex(),
		UserID:    comment.UserID.Hex(),
		VideoID:   comment.VideoID.Hex(),
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}

	if comment.ParentID != nil {
		response.ParentID = comment.ParentID.Hex()
	}

	return response, nil
}

func (s *Service) GetComments(ctx context.Context, videoID string, limit, offset int64) ([]*CommentResponse, error) {
	// Convert video ID to ObjectID
	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	// Set default pagination values
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Get comments
	comments, err := s.repo.GetComments(ctx, videoObjectID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responses := make([]*CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = &CommentResponse{
			ID:        comment.ID.Hex(),
			UserID:    comment.UserID.Hex(),
			VideoID:   comment.VideoID.Hex(),
			Text:      comment.Text,
			CreatedAt: comment.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetCommentReplies(ctx context.Context, commentID string) ([]*CommentResponse, error) {
	// Convert comment ID to ObjectID
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, errors.New("invalid comment ID")
	}

	// Verify parent comment exists
	_, err = s.repo.GetCommentByID(ctx, commentObjectID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Get replies
	replies, err := s.repo.GetCommentReplies(ctx, commentObjectID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responses := make([]*CommentResponse, len(replies))
	for i, reply := range replies {
		responses[i] = &CommentResponse{
			ID:        reply.ID.Hex(),
			UserID:    reply.UserID.Hex(),
			VideoID:   reply.VideoID.Hex(),
			Text:      reply.Text,
			ParentID:  reply.ParentID.Hex(),
			CreatedAt: reply.CreatedAt,
		}
	}

	return responses, nil
}

// Share operations

func (s *Service) RecordShare(ctx context.Context, userID, videoID string) (*ShareResponse, error) {
	// Convert IDs to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	// Record share
	err = s.repo.RecordShare(ctx, userObjectID, videoObjectID)
	if err != nil {
		return nil, err
	}

	// Get updated share count
	stats, err := s.repo.GetVideoStats(ctx, videoObjectID)
	if err != nil {
		return nil, err
	}

	return &ShareResponse{
		VideoID:    videoID,
		ShareCount: stats["share_count"],
	}, nil
}

func (s *Service) GetShareCount(ctx context.Context, videoID string) (int, error) {
	// Convert video ID to ObjectID
	videoObjectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return 0, errors.New("invalid video ID")
	}

	return s.repo.GetShareCount(ctx, videoObjectID)
}

// Validation helpers

func (s *Service) validateComment(req *CreateCommentRequest) error {
	text := strings.TrimSpace(req.Text)

	if text == "" {
		return errors.New("comment text is required")
	}

	if len(text) > 500 {
		return errors.New("comment text must be less than 500 characters")
	}

	return nil
}
