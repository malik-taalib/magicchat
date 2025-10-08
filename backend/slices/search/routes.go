package search

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"magicchat/slices/auth"
)

func Routes(db *mongo.Database) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()

	// Public routes - anyone can search
	r.Get("/search", handler.Search)
	r.Get("/trending/hashtags", handler.GetTrendingHashtags)
	r.Get("/hashtags/{tag}/videos", handler.GetVideosByHashtag)

	// Optional: Protected routes for personalized search (if needed in future)
	// r.Group(func(r chi.Router) {
	// 	r.Use(auth.AuthMiddleware)
	// 	// Add protected search endpoints here if needed
	// })

	return r
}

// ProtectedRoutes returns routes that require authentication
// This is useful if you want to separate public and protected search functionality
func ProtectedRoutes(db *mongo.Database) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware)

	// Protected search endpoints (for logged-in users only)
	// Could include personalized search results, search history, etc.
	r.Get("/search", handler.Search)
	r.Get("/trending/hashtags", handler.GetTrendingHashtags)
	r.Get("/hashtags/{tag}/videos", handler.GetVideosByHashtag)

	return r
}
