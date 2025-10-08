package engagement

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	likesCollection    *mongo.Collection
	commentsCollection *mongo.Collection
	sharesCollection   *mongo.Collection
	videosCollection   *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		likesCollection:    db.Collection("likes"),
		commentsCollection: db.Collection("comments"),
		sharesCollection:   db.Collection("shares"),
		videosCollection:   db.Collection("videos"),
	}
}

// Like operations

func (r *Repository) LikeVideo(ctx context.Context, userID, videoID primitive.ObjectID) error {
	// Check if already liked
	exists, err := r.IsVideoLiked(ctx, userID, videoID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("video already liked")
	}

	// Create like
	like := &Like{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		VideoID:   videoID,
		CreatedAt: time.Now(),
	}

	session, err := r.likesCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Use transaction to ensure atomicity
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert like
		_, err := r.likesCollection.InsertOne(sessCtx, like)
		if err != nil {
			return nil, err
		}

		// Increment like count on video
		update := bson.M{
			"$inc": bson.M{"like_count": 1},
		}
		_, err = r.videosCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": videoID},
			update,
		)
		return nil, err
	})

	return err
}

func (r *Repository) UnlikeVideo(ctx context.Context, userID, videoID primitive.ObjectID) error {
	// Check if liked
	exists, err := r.IsVideoLiked(ctx, userID, videoID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("video not liked")
	}

	session, err := r.likesCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Use transaction to ensure atomicity
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Delete like
		_, err := r.likesCollection.DeleteOne(
			sessCtx,
			bson.M{
				"user_id":  userID,
				"video_id": videoID,
			},
		)
		if err != nil {
			return nil, err
		}

		// Decrement like count on video
		update := bson.M{
			"$inc": bson.M{"like_count": -1},
		}
		_, err = r.videosCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": videoID},
			update,
		)
		return nil, err
	})

	return err
}

func (r *Repository) IsVideoLiked(ctx context.Context, userID, videoID primitive.ObjectID) (bool, error) {
	count, err := r.likesCollection.CountDocuments(ctx, bson.M{
		"user_id":  userID,
		"video_id": videoID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) GetLikeCount(ctx context.Context, videoID primitive.ObjectID) (int, error) {
	count, err := r.likesCollection.CountDocuments(ctx, bson.M{
		"video_id": videoID,
	})
	return int(count), err
}

// Comment operations

func (r *Repository) CreateComment(ctx context.Context, comment *Comment) error {
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if comment.Replies == nil {
		comment.Replies = []primitive.ObjectID{}
	}

	session, err := r.commentsCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Use transaction to ensure atomicity
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert comment
		_, err := r.commentsCollection.InsertOne(sessCtx, comment)
		if err != nil {
			return nil, err
		}

		// If this is a reply, add to parent's replies array
		if comment.ParentID != nil {
			_, err = r.commentsCollection.UpdateOne(
				sessCtx,
				bson.M{"_id": *comment.ParentID},
				bson.M{
					"$push": bson.M{"replies": comment.ID},
					"$set":  bson.M{"updated_at": time.Now()},
				},
			)
			if err != nil {
				return nil, err
			}
		}

		// Increment comment count on video
		update := bson.M{
			"$inc": bson.M{"comment_count": 1},
		}
		_, err = r.videosCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": comment.VideoID},
			update,
		)
		return nil, err
	})

	return err
}

func (r *Repository) GetComments(ctx context.Context, videoID primitive.ObjectID, limit, offset int64) ([]*Comment, error) {
	// Only get top-level comments (no parent_id)
	filter := bson.M{
		"video_id":  videoID,
		"parent_id": bson.M{"$exists": false},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // Newest first
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.commentsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) GetCommentReplies(ctx context.Context, parentID primitive.ObjectID) ([]*Comment, error) {
	filter := bson.M{
		"parent_id": parentID,
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}) // Oldest first for replies

	cursor, err := r.commentsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) GetCommentByID(ctx context.Context, commentID primitive.ObjectID) (*Comment, error) {
	var comment Comment
	err := r.commentsCollection.FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return &comment, nil
}

// Share operations

func (r *Repository) RecordShare(ctx context.Context, userID, videoID primitive.ObjectID) error {
	share := &Share{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		VideoID:   videoID,
		CreatedAt: time.Now(),
	}

	session, err := r.sharesCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Use transaction to ensure atomicity
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert share
		_, err := r.sharesCollection.InsertOne(sessCtx, share)
		if err != nil {
			return nil, err
		}

		// Increment share count on video
		update := bson.M{
			"$inc": bson.M{"share_count": 1},
		}
		_, err = r.videosCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": videoID},
			update,
		)
		return nil, err
	})

	return err
}

func (r *Repository) GetShareCount(ctx context.Context, videoID primitive.ObjectID) (int, error) {
	count, err := r.sharesCollection.CountDocuments(ctx, bson.M{
		"video_id": videoID,
	})
	return int(count), err
}

// Video helper to get current video stats
func (r *Repository) GetVideoStats(ctx context.Context, videoID primitive.ObjectID) (map[string]int, error) {
	var result struct {
		LikeCount    int `bson:"like_count"`
		CommentCount int `bson:"comment_count"`
		ShareCount   int `bson:"share_count"`
	}

	err := r.videosCollection.FindOne(
		ctx,
		bson.M{"_id": videoID},
		options.FindOne().SetProjection(bson.M{
			"like_count":    1,
			"comment_count": 1,
			"share_count":   1,
		}),
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("video not found")
		}
		return nil, err
	}

	return map[string]int{
		"like_count":    result.LikeCount,
		"comment_count": result.CommentCount,
		"share_count":   result.ShareCount,
	}, nil
}
