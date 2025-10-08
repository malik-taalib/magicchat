package videoupload

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("videos"),
	}
}

func (r *Repository) CreateVideo(ctx context.Context, video *Video) error {
	video.ID = primitive.NewObjectID()
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.ViewCount = 0
	video.LikeCount = 0
	video.CommentCount = 0
	video.ShareCount = 0
	video.ProcessingStatus = StatusPending

	_, err := r.collection.InsertOne(ctx, video)
	return err
}

func (r *Repository) GetVideoByID(ctx context.Context, id string) (*Video, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var video Video
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("video not found")
		}
		return nil, err
	}
	return &video, nil
}

func (r *Repository) UpdateVideoStatus(ctx context.Context, id string, status ProcessingStatus, videoURL string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"processing_status": status,
			"video_url":         videoURL,
			"updated_at":        time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *Repository) UpdateVideoMetadata(ctx context.Context, id string, duration int, thumbnailURL string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"duration":      duration,
			"thumbnail_url": thumbnailURL,
			"updated_at":    time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
