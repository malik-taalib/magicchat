package auth

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(db *mongo.Database) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Post("/logout", handler.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/me", handler.Me)
		r.Put("/profile", handler.UpdateProfile)
	})

	return r
}
