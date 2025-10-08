package following

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
	followCollection *mongo.Collection
	userCollection   *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		followCollection: db.Collection("follows"),
		userCollection:   db.Collection("users"),
	}
}

// FollowUser creates a follow relationship between two users
func (r *Repository) FollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error {
	// Check if already following
	exists, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following this user")
	}

	// Create follow relationship
	follow := &Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	_, err = r.followCollection.InsertOne(ctx, follow)
	if err != nil {
		return err
	}

	// Update counts
	return r.updateFollowCounts(ctx, followerID, followingID, 1)
}

// UnfollowUser removes a follow relationship between two users
func (r *Repository) UnfollowUser(ctx context.Context, followerID, followingID primitive.ObjectID) error {
	// Check if following
	exists, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("not following this user")
	}

	// Delete follow relationship
	filter := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	}

	_, err = r.followCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Update counts
	return r.updateFollowCounts(ctx, followerID, followingID, -1)
}

// IsFollowing checks if followerID is following followingID
func (r *Repository) IsFollowing(ctx context.Context, followerID, followingID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	}

	count, err := r.followCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetFollowers returns a paginated list of users who follow the specified user
func (r *Repository) GetFollowers(ctx context.Context, userID primitive.ObjectID, cursor string, limit int) ([]UserProfile, string, error) {
	// Build the filter
	filter := bson.M{"following_id": userID}

	// If cursor is provided, add it to filter for pagination
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", err
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	// Find follow documents
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit + 1)) // Fetch one extra to check if there are more

	cursor2, err := r.followCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursor2.Close(ctx)

	var follows []Follow
	if err = cursor2.All(ctx, &follows); err != nil {
		return nil, "", err
	}

	// Check if there are more results
	hasMore := len(follows) > limit
	if hasMore {
		follows = follows[:limit]
	}

	// Extract follower IDs
	followerIDs := make([]primitive.ObjectID, len(follows))
	for i, follow := range follows {
		followerIDs[i] = follow.FollowerID
	}

	// Fetch user profiles
	profiles, err := r.getUserProfiles(ctx, followerIDs)
	if err != nil {
		return nil, "", err
	}

	// Calculate next cursor
	var nextCursor string
	if hasMore && len(follows) > 0 {
		nextCursor = follows[len(follows)-1].ID.Hex()
	}

	return profiles, nextCursor, nil
}

// GetFollowing returns a paginated list of users that the specified user follows
func (r *Repository) GetFollowing(ctx context.Context, userID primitive.ObjectID, cursor string, limit int) ([]UserProfile, string, error) {
	// Build the filter
	filter := bson.M{"follower_id": userID}

	// If cursor is provided, add it to filter for pagination
	if cursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", err
		}
		filter["_id"] = bson.M{"$lt": cursorID}
	}

	// Find follow documents
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit + 1)) // Fetch one extra to check if there are more

	cursor2, err := r.followCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursor2.Close(ctx)

	var follows []Follow
	if err = cursor2.All(ctx, &follows); err != nil {
		return nil, "", err
	}

	// Check if there are more results
	hasMore := len(follows) > limit
	if hasMore {
		follows = follows[:limit]
	}

	// Extract following IDs
	followingIDs := make([]primitive.ObjectID, len(follows))
	for i, follow := range follows {
		followingIDs[i] = follow.FollowingID
	}

	// Fetch user profiles
	profiles, err := r.getUserProfiles(ctx, followingIDs)
	if err != nil {
		return nil, "", err
	}

	// Calculate next cursor
	var nextCursor string
	if hasMore && len(follows) > 0 {
		nextCursor = follows[len(follows)-1].ID.Hex()
	}

	return profiles, nextCursor, nil
}

// GetFollowerCount returns the count of followers for a user
func (r *Repository) GetFollowerCount(ctx context.Context, userID primitive.ObjectID) (int, error) {
	count, err := r.followCollection.CountDocuments(ctx, bson.M{"following_id": userID})
	return int(count), err
}

// GetFollowingCount returns the count of users that a user follows
func (r *Repository) GetFollowingCount(ctx context.Context, userID primitive.ObjectID) (int, error) {
	count, err := r.followCollection.CountDocuments(ctx, bson.M{"follower_id": userID})
	return int(count), err
}

// updateFollowCounts updates the follower and following counts for both users
func (r *Repository) updateFollowCounts(ctx context.Context, followerID, followingID primitive.ObjectID, delta int) error {
	// Update follower's following_count
	_, err := r.userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followerID},
		bson.M{"$inc": bson.M{"following_count": delta}},
	)
	if err != nil {
		return err
	}

	// Update following's follower_count
	_, err = r.userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followingID},
		bson.M{"$inc": bson.M{"follower_count": delta}},
	)
	return err
}

// getUserProfiles fetches user profiles by their IDs
func (r *Repository) getUserProfiles(ctx context.Context, userIDs []primitive.ObjectID) ([]UserProfile, error) {
	if len(userIDs) == 0 {
		return []UserProfile{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	projection := bson.M{
		"_id":             1,
		"username":        1,
		"display_name":    1,
		"avatar_url":      1,
		"follower_count":  1,
		"following_count": 1,
	}

	opts := options.Find().SetProjection(projection)
	cursor, err := r.userCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []UserProfile
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	// Maintain the order of the input userIDs
	profileMap := make(map[primitive.ObjectID]UserProfile)
	for _, profile := range profiles {
		profileMap[profile.ID] = profile
	}

	orderedProfiles := make([]UserProfile, 0, len(userIDs))
	for _, id := range userIDs {
		if profile, ok := profileMap[id]; ok {
			orderedProfiles = append(orderedProfiles, profile)
		}
	}

	return orderedProfiles, nil
}

// GetUserByID fetches a user by ID and returns their follow counts
func (r *Repository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*UserProfile, error) {
	filter := bson.M{"_id": userID}
	projection := bson.M{
		"_id":             1,
		"username":        1,
		"display_name":    1,
		"avatar_url":      1,
		"follower_count":  1,
		"following_count": 1,
	}

	opts := options.FindOne().SetProjection(projection)
	var profile UserProfile
	err := r.userCollection.FindOne(ctx, filter, opts).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &profile, nil
}
