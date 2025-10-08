package engagement

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

	// All engagement routes are protected with auth middleware
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		// Like endpoints
		r.Post("/{id}/like", handler.LikeVideo)
		r.Delete("/{id}/like", handler.UnlikeVideo)

		// Comment endpoints
		r.Post("/{id}/comments", handler.CreateComment)
		r.Get("/{id}/comments", handler.GetComments)

		// Get replies for a specific comment
		r.Get("/comments/{id}/replies", handler.GetCommentReplies)

		// Share endpoint
		r.Post("/{id}/share", handler.ShareVideo)
	})

	return r
}
