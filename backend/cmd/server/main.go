package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"magicchat/pkg/cache"
	"magicchat/pkg/config"
	"magicchat/pkg/database"
	"magicchat/pkg/storage"
	"magicchat/slices/auth"
	"magicchat/slices/engagement"
	"magicchat/slices/following"
	"magicchat/slices/notifications"
	"magicchat/slices/search"
	videofeed "magicchat/slices/video-feed"
	videoupload "magicchat/slices/video-upload"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting MagicChat server in %s mode on port %s", cfg.Server.Env, cfg.Server.Port)

	// Connect to MongoDB
	db, err := database.ConnectMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectMongoDB()
	log.Println("✓ MongoDB connected")

	// Connect to Redis
	_, err = cache.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.DisconnectRedis()
	log.Println("✓ Redis connected")

	// Initialize storage client
	storageClient, err := storage.InitStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Println("✓ Storage client initialized")

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"magicchat"}`))
	})

	// Mount API routes
	r.Route("/api", func(r chi.Router) {
		// Authentication routes
		r.Mount("/auth", auth.Routes(db))

		// Video upload routes (POST /videos/upload, GET /videos/:id/status)
		r.Mount("/videos", videoupload.Routes(db, storageClient))

		// Video feed routes (GET /feed/for-you, GET /feed/following)
		r.Mount("/feed", videofeed.Routes(db))

		// Engagement routes (POST /engage/:id/like, POST /engage/:id/comments, etc)
		// Changed from /videos to /engage to avoid conflict
		r.Mount("/engage", engagement.Routes(db))

		// Following routes
		r.Mount("/users", following.Routes(db))

		// Search & discovery routes
		r.Mount("/search", search.Routes(db))
		r.Mount("/trending", search.Routes(db))
		r.Mount("/hashtags", search.Routes(db))

		// Notifications routes (includes WebSocket)
		r.Mount("/notifications", notifications.Routes(db))
	})

	// Print registered routes
	log.Println("\n📋 Registered API routes:")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("  %s %s", method, route)
		return nil
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("\n🚀 Server starting on http://localhost:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server exited gracefully")
}
