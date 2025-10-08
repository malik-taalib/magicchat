package following

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"magicchat/slices/auth"
)

// Routes sets up the following slice routes
func Routes(db *mongo.Database) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()

	// All routes are protected with authentication
	r.Use(auth.AuthMiddleware)

	// Follow/Unfollow endpoints
	r.Post("/{id}/follow", handler.FollowUser)
	r.Delete("/{id}/follow", handler.UnfollowUser)

	// Check if following (optional)
	r.Get("/{id}/following/check", handler.IsFollowing)

	// Get followers and following lists
	r.Get("/{id}/followers", handler.GetFollowers)
	r.Get("/{id}/following", handler.GetFollowing)

	return r
}
