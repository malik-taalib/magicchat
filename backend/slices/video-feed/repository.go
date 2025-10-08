package videofeed

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
	videosCollection  *mongo.Collection
	usersCollection   *mongo.Collection
	followsCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		videosCollection:  db.Collection("videos"),
		usersCollection:   db.Collection("users"),
		followsCollection: db.Collection("follows"),
	}
}

// GetForYouFeed returns videos for the "For You" feed sorted by engagement metrics
// Uses cursor-based pagination for efficient scrolling
func (r *Repository) GetForYouFeed(ctx context.Context, userID string, cursor string, limit int) ([]*FeedVideo, error) {
	// Build filter for cursor-based pagination
	filter := bson.M{
		"processing_status": "completed", // Only show completed videos
	}

	// If cursor is provided, only return videos older than cursor
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	// Calculate engagement score and sort by it
	// Formula: (like_count * 3 + view_count * 0.5 + comment_count * 5 + share_count * 10) / age_in_hours
	pipeline := mongo.Pipeline{
		// Match completed videos with cursor filter
		{{"$match", filter}},
		// Add engagement score calculation
		{{"$addFields", bson.M{
			"engagement_score": bson.M{
				"$divide": bson.A{
					bson.M{"$add": bson.A{
						bson.M{"$multiply": bson.A{"$like_count", 3}},
						bson.M{"$multiply": bson.A{"$view_count", 0.5}},
						bson.M{"$multiply": bson.A{"$comment_count", 5}},
						bson.M{"$multiply": bson.A{"$share_count", 10}},
					}},
					bson.M{"$max": bson.A{
						bson.M{"$divide": bson.A{
							bson.M{"$subtract": bson.A{time.Now(), "$created_at"}},
							3600000, // Convert milliseconds to hours
						}},
						1, // Minimum 1 hour to avoid division by zero
					}},
				},
			},
		}}},
		// Sort by engagement score (descending), then by created_at (descending)
		{{"$sort", bson.D{{"engagement_score", -1}, {"created_at", -1}}}},
		// Limit results
		{{"$limit", limit}},
		// Lookup user information
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		// Unwind user array
		{{"$unwind", "$user"}},
		// Project final shape
		{{"$project", bson.M{
			"_id":               1,
			"user_id":           1,
			"username":          "$user.username",
			"display_name":      "$user.display_name",
			"avatar_url":        "$user.avatar_url",
			"title":             1,
			"description":       1,
			"video_url":         1,
			"thumbnail_url":     1,
			"duration":          1,
			"hashtags":          1,
			"view_count":        1,
			"like_count":        1,
			"comment_count":     1,
			"share_count":       1,
			"processing_status": 1,
			"created_at":        1,
			"updated_at":        1,
		}}},
	}

	cursor_db, err := r.videosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor_db.Close(ctx)

	var videos []*FeedVideo
	if err := cursor_db.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

// GetFollowingFeed returns videos from users that the given user follows
// Uses cursor-based pagination for efficient scrolling
func (r *Repository) GetFollowingFeed(ctx context.Context, userID string, cursor string, limit int) ([]*FeedVideo, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Build filter for cursor-based pagination
	filter := bson.M{
		"processing_status": "completed",
	}

	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	// Aggregate pipeline to get videos from followed users
	pipeline := mongo.Pipeline{
		// Get all users that current user follows
		{{"$match", bson.M{"follower_id": userObjectID}}},
		// Lookup videos from followed users
		{{"$lookup", bson.M{
			"from": "videos",
			"let":  bson.M{"followed_user_id": "$following_id"},
			"pipeline": mongo.Pipeline{
				{{"$match", bson.M{"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$user_id", "$$followed_user_id"}},
						bson.M{"$eq": bson.A{"$processing_status", "completed"}},
					},
				}}}},
			},
			"as": "videos",
		}}},
		// Unwind videos
		{{"$unwind", "$videos"}},
		// Replace root with video document
		{{"$replaceRoot", bson.M{"newRoot": "$videos"}}},
		// Apply cursor filter if needed
		{{"$match", filter}},
		// Sort by created_at (descending) - most recent first
		{{"$sort", bson.D{{"created_at", -1}}}},
		// Limit results
		{{"$limit", limit}},
		// Lookup user information
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		// Unwind user array
		{{"$unwind", "$user"}},
		// Project final shape
		{{"$project", bson.M{
			"_id":               1,
			"user_id":           1,
			"username":          "$user.username",
			"display_name":      "$user.display_name",
			"avatar_url":        "$user.avatar_url",
			"title":             1,
			"description":       1,
			"video_url":         1,
			"thumbnail_url":     1,
			"duration":          1,
			"hashtags":          1,
			"view_count":        1,
			"like_count":        1,
			"comment_count":     1,
			"share_count":       1,
			"processing_status": 1,
			"created_at":        1,
			"updated_at":        1,
		}}},
	}

	cursor_db, err := r.followsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor_db.Close(ctx)

	var videos []*FeedVideo
	if err := cursor_db.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

// GetVideoByID returns a single video with user information
func (r *Repository) GetVideoByID(ctx context.Context, videoID string) (*FeedVideo, error) {
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return nil, errors.New("invalid video ID")
	}

	pipeline := mongo.Pipeline{
		// Match video by ID
		{{"$match", bson.M{"_id": objectID}}},
		// Lookup user information
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		// Unwind user array
		{{"$unwind", "$user"}},
		// Project final shape
		{{"$project", bson.M{
			"_id":               1,
			"user_id":           1,
			"username":          "$user.username",
			"display_name":      "$user.display_name",
			"avatar_url":        "$user.avatar_url",
			"title":             1,
			"description":       1,
			"video_url":         1,
			"thumbnail_url":     1,
			"duration":          1,
			"hashtags":          1,
			"view_count":        1,
			"like_count":        1,
			"comment_count":     1,
			"share_count":       1,
			"processing_status": 1,
			"created_at":        1,
			"updated_at":        1,
		}}},
	}

	cursor, err := r.videosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []*FeedVideo
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		return nil, errors.New("video not found")
	}

	return videos[0], nil
}

// IncrementViewCount atomically increments the view count for a video
func (r *Repository) IncrementViewCount(ctx context.Context, videoID string) error {
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return errors.New("invalid video ID")
	}

	update := bson.M{
		"$inc": bson.M{"view_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	opts := options.Update().SetUpsert(false)
	result, err := r.videosCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("video not found")
	}

	return nil
}
