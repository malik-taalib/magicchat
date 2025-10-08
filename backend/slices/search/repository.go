package search

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
	usersCollection    *mongo.Collection
	videosCollection   *mongo.Collection
	hashtagsCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		usersCollection:    db.Collection("users"),
		videosCollection:   db.Collection("videos"),
		hashtagsCollection: db.Collection("hashtags"),
	}
}

// SearchUsers performs a text search on username and display_name fields
// Uses cursor-based pagination for efficient scrolling
func (r *Repository) SearchUsers(ctx context.Context, query string, cursor string, limit int) ([]*UserSearchResult, error) {
	// Build filter for text search
	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": query, "$options": "i"}},
			{"display_name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// If cursor is provided, only return users with IDs less than cursor
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	// Sort by follower count (descending) for relevance, then by _id for consistent pagination
	opts := options.Find().
		SetSort(bson.D{{"follower_count", -1}, {"_id", -1}}).
		SetLimit(int64(limit))

	cursor_db, err := r.usersCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor_db.Close(ctx)

	var users []*UserSearchResult
	for cursor_db.Next(ctx) {
		var doc bson.M
		if err := cursor_db.Decode(&doc); err != nil {
			return nil, err
		}

		user := &UserSearchResult{
			ID:            doc["_id"].(primitive.ObjectID),
			Username:      getStringOrEmpty(doc, "username"),
			DisplayName:   getStringOrEmpty(doc, "display_name"),
			Bio:           getStringOrEmpty(doc, "bio"),
			AvatarURL:     getStringOrEmpty(doc, "avatar_url"),
			FollowerCount: getIntOrZero(doc, "follower_count"),
			VideoCount:    getIntOrZero(doc, "video_count"),
			IsVerified:    getBoolOrFalse(doc, "is_verified"),
		}
		users = append(users, user)
	}

	if err := cursor_db.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// SearchVideos performs a text search on title, description, and hashtags fields
// Uses cursor-based pagination for efficient scrolling
func (r *Repository) SearchVideos(ctx context.Context, query string, cursor string, limit int) ([]*VideoSearchResult, error) {
	// Build filter for cursor-based pagination
	matchFilter := bson.M{
		"processing_status": "completed", // Only show completed videos
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"hashtags": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// If cursor is provided, only return videos with IDs less than cursor
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		matchFilter["_id"] = bson.M{"$lt": cursorID}
	}

	// Use aggregation pipeline to join with users collection
	pipeline := mongo.Pipeline{
		// Match videos with search criteria
		{{"$match", matchFilter}},
		// Sort by relevance score (view_count + like_count), then by created_at
		{{"$addFields", bson.M{
			"relevance_score": bson.M{
				"$add": bson.A{
					"$view_count",
					bson.M{"$multiply": bson.A{"$like_count", 2}},
				},
			},
		}}},
		{{"$sort", bson.D{{"relevance_score", -1}, {"created_at", -1}}}},
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
			"_id":           1,
			"user_id":       1,
			"username":      "$user.username",
			"display_name":  "$user.display_name",
			"avatar_url":    "$user.avatar_url",
			"title":         1,
			"description":   1,
			"video_url":     1,
			"thumbnail_url": 1,
			"duration":      1,
			"hashtags":      1,
			"view_count":    1,
			"like_count":    1,
			"comment_count": 1,
			"created_at":    1,
		}}},
	}

	cursor_db, err := r.videosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor_db.Close(ctx)

	var videos []*VideoSearchResult
	if err := cursor_db.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

// SearchHashtags performs a search on hashtag tags
func (r *Repository) SearchHashtags(ctx context.Context, query string, limit int) ([]*HashtagSearchResult, error) {
	// Build filter for text search
	filter := bson.M{
		"tag": bson.M{"$regex": query, "$options": "i"},
	}

	// Sort by trending score (descending) and video count (descending)
	opts := options.Find().
		SetSort(bson.D{{"trending_score", -1}, {"video_count", -1}}).
		SetLimit(int64(limit))

	cursor, err := r.hashtagsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hashtags []*HashtagSearchResult
	for cursor.Next(ctx) {
		var hashtag Hashtag
		if err := cursor.Decode(&hashtag); err != nil {
			return nil, err
		}

		result := &HashtagSearchResult{
			Tag:           hashtag.Tag,
			VideoCount:    hashtag.VideoCount,
			TrendingScore: hashtag.TrendingScore,
			LastUsed:      hashtag.LastUsed,
		}
		hashtags = append(hashtags, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return hashtags, nil
}

// GetTrendingHashtags returns the top trending hashtags by video count and trending score
func (r *Repository) GetTrendingHashtags(ctx context.Context, limit int) ([]*HashtagSearchResult, error) {
	// Calculate trending score: (video_count * recency_factor)
	// Recency factor: higher for recently used hashtags
	pipeline := mongo.Pipeline{
		// Calculate trending score with recency factor
		{{"$addFields", bson.M{
			"days_since_last_used": bson.M{
				"$divide": bson.A{
					bson.M{"$subtract": bson.A{time.Now(), "$last_used"}},
					86400000, // Convert milliseconds to days
				},
			},
		}}},
		{{"$addFields", bson.M{
			"recency_factor": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$lte": bson.A{"$days_since_last_used", 1}},
					"then": 1.0,
					"else": bson.M{
						"$divide": bson.A{
							1,
							bson.M{"$add": bson.A{
								1,
								bson.M{"$multiply": bson.A{"$days_since_last_used", 0.1}},
							}},
						},
					},
				},
			},
		}}},
		{{"$addFields", bson.M{
			"trending_score": bson.M{
				"$multiply": bson.A{"$video_count", "$recency_factor"},
			},
		}}},
		// Sort by trending score (descending)
		{{"$sort", bson.D{{"trending_score", -1}}}},
		// Limit results
		{{"$limit", limit}},
		// Project final shape
		{{"$project", bson.M{
			"tag":            1,
			"video_count":    1,
			"trending_score": 1,
			"last_used":      1,
		}}},
	}

	cursor, err := r.hashtagsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hashtags []*HashtagSearchResult
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		hashtag := &HashtagSearchResult{
			Tag:           getStringOrEmpty(doc, "tag"),
			VideoCount:    getIntOrZero(doc, "video_count"),
			TrendingScore: getFloat64OrZero(doc, "trending_score"),
			LastUsed:      getTimeOrNow(doc, "last_used"),
		}
		hashtags = append(hashtags, hashtag)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return hashtags, nil
}

// GetVideosByHashtag returns videos that contain a specific hashtag
func (r *Repository) GetVideosByHashtag(ctx context.Context, tag string, cursor string, limit int) ([]*VideoSearchResult, error) {
	// Build filter for hashtag search
	matchFilter := bson.M{
		"processing_status": "completed",
		"hashtags":          tag, // Exact match on hashtag
	}

	// If cursor is provided, only return videos with IDs less than cursor
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		matchFilter["_id"] = bson.M{"$lt": cursorID}
	}

	// Use aggregation pipeline to join with users collection
	pipeline := mongo.Pipeline{
		// Match videos with hashtag
		{{"$match", matchFilter}},
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
			"_id":           1,
			"user_id":       1,
			"username":      "$user.username",
			"display_name":  "$user.display_name",
			"avatar_url":    "$user.avatar_url",
			"title":         1,
			"description":   1,
			"video_url":     1,
			"thumbnail_url": 1,
			"duration":      1,
			"hashtags":      1,
			"view_count":    1,
			"like_count":    1,
			"comment_count": 1,
			"created_at":    1,
		}}},
	}

	cursor_db, err := r.videosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor_db.Close(ctx)

	var videos []*VideoSearchResult
	if err := cursor_db.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

// Helper functions to safely extract values from bson.M
func getStringOrEmpty(doc bson.M, key string) string {
	if val, ok := doc[key].(string); ok {
		return val
	}
	return ""
}

func getIntOrZero(doc bson.M, key string) int {
	if val, ok := doc[key].(int); ok {
		return val
	}
	if val, ok := doc[key].(int32); ok {
		return int(val)
	}
	if val, ok := doc[key].(int64); ok {
		return int(val)
	}
	return 0
}

func getBoolOrFalse(doc bson.M, key string) bool {
	if val, ok := doc[key].(bool); ok {
		return val
	}
	return false
}

func getFloat64OrZero(doc bson.M, key string) float64 {
	if val, ok := doc[key].(float64); ok {
		return val
	}
	if val, ok := doc[key].(float32); ok {
		return float64(val)
	}
	if val, ok := doc[key].(int); ok {
		return float64(val)
	}
	if val, ok := doc[key].(int32); ok {
		return float64(val)
	}
	if val, ok := doc[key].(int64); ok {
		return float64(val)
	}
	return 0
}

func getTimeOrNow(doc bson.M, key string) time.Time {
	if val, ok := doc[key].(time.Time); ok {
		return val
	}
	if val, ok := doc[key].(primitive.DateTime); ok {
		return val.Time()
	}
	return time.Now()
}
