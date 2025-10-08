package videoupload

import (
	"context"
	"io"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"magicchat/slices/auth"
)

type StorageProvider interface {
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
}

func Routes(db *mongo.Database, storage StorageClient) chi.Router {
	repo := NewRepository(db)
	service := NewService(repo, storage)
	handler := NewHandler(service)

	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/upload", handler.Upload)
		r.Get("/{id}/status", handler.GetStatus)
	})

	// Webhook for video processing (should be protected with API key in production)
	r.Post("/process", handler.ProcessWebhook)

	return r
}
