package videofeed

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"magicchat/slices/auth"
)

// Routes creates and returns a chi router with all video feed routes
func Routes(db *mongo.Database) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()

	// All feed routes require authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		// Feed endpoints
		r.Get("/for-you", handler.GetForYouFeed)
		r.Get("/following", handler.GetFollowingFeed)

		// Single video endpoint
		r.Get("/{id}", handler.GetVideo)
	})

	return r
}
