package videoupload

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"magicchat/pkg/config"
)

type StorageClient interface {
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
}

type Service struct {
	repo    *Repository
	storage StorageClient
}

func NewService(repo *Repository, storage StorageClient) *Service {
	return &Service{
		repo:    repo,
		storage: storage,
	}
}

func (s *Service) UploadVideo(ctx context.Context, userID string, req *UploadRequest, file *multipart.FileHeader) (*Video, error) {
	// Validate file
	if err := s.validateVideo(file); err != nil {
		return nil, err
	}

	// Convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create video record
	video := &Video{
		UserID:      userObjectID,
		Title:       req.Title,
		Description: req.Description,
		Hashtags:    req.Hashtags,
	}

	err = s.repo.CreateVideo(ctx, video)
	if err != nil {
		return nil, err
	}

	// Upload file to storage
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	videoURL, err := s.storage.UploadFile(ctx, src, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	// Update video with URL and set to processing
	err = s.repo.UpdateVideoStatus(ctx, video.ID.Hex(), StatusProcessing, videoURL)
	if err != nil {
		return nil, err
	}

	video.VideoURL = videoURL
	video.ProcessingStatus = StatusProcessing

	// TODO: Trigger async video processing job (transcoding, thumbnail generation)
	// This would typically be done via Redis queue + worker

	return video, nil
}

func (s *Service) GetVideoStatus(ctx context.Context, videoID string) (*Video, error) {
	return s.repo.GetVideoByID(ctx, videoID)
}

func (s *Service) ProcessVideo(ctx context.Context, videoID string) error {
	// This would be called by a background worker
	// TODO: Implement actual video processing with FFmpeg
	// - Generate thumbnails
	// - Transcode to multiple qualities
	// - Extract duration
	// - Generate preview clips

	// For now, just mark as completed
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Simulate processing completion
	err = s.repo.UpdateVideoStatus(ctx, videoID, StatusCompleted, video.VideoURL)
	if err != nil {
		return err
	}

	// TODO: Update duration and thumbnail URL
	// err = s.repo.UpdateVideoMetadata(ctx, videoID, duration, thumbnailURL)

	return nil
}

func (s *Service) validateVideo(file *multipart.FileHeader) error {
	cfg := config.Load()

	// Check file size
	maxSize := int64(cfg.Video.MaxSizeMB) * 1024 * 1024
	if file.Size > maxSize {
		return errors.New("video file too large")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")

	allowed := false
	for _, format := range cfg.Video.AllowedFormats {
		if ext == format {
			allowed = true
			break
		}
	}

	if !allowed {
		return errors.New("invalid video format")
	}

	return nil
}
