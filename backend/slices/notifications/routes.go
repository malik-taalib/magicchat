package notifications

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"magicchat/slices/auth"
)

// Routes sets up both HTTP and WebSocket routes for notifications
func Routes(db *mongo.Database) chi.Router {
	// Initialize WebSocket manager
	wsManager := NewWebSocketManager()

	// Start WebSocket manager in a goroutine
	go wsManager.Run()

	// Initialize repository and service
	repo := NewRepository(db)
	service := NewService(repo, db, wsManager)
	handler := NewHandler(service, wsManager)

	r := chi.NewRouter()

	// All routes require authentication
	r.Use(auth.AuthMiddleware)

	// HTTP routes
	r.Get("/", handler.GetNotifications)                // GET /notifications
	r.Put("/{id}/read", handler.MarkAsRead)            // PUT /notifications/:id/read
	r.Put("/read-all", handler.MarkAllAsRead)          // PUT /notifications/read-all
	r.Get("/unread-count", handler.GetUnreadCount)     // GET /notifications/unread-count

	// WebSocket route
	r.Get("/stream", handler.HandleWebSocket)          // WS /notifications/stream

	return r
}

// GetWebSocketManager returns the WebSocket manager instance
// This function can be used by other slices to send notifications
func GetWebSocketManager(db *mongo.Database) *WebSocketManager {
	// Create a singleton pattern or dependency injection
	// For now, we'll create a new instance
	wsManager := NewWebSocketManager()
	go wsManager.Run()
	return wsManager
}

// GetService returns a notifications service instance
// This can be used by other slices to create notifications
func GetService(db *mongo.Database, wsManager *WebSocketManager) *Service {
	repo := NewRepository(db)
	return NewService(repo, db, wsManager)
}
