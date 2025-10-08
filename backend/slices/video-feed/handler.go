package videofeed

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"magicchat/slices/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetForYouFeed handles requests for the "For You" feed
// GET /for-you?cursor=<cursor>&limit=<limit>
func (h *Handler) GetForYouFeed(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters
	req := &FeedRequest{
		Cursor: r.URL.Query().Get("cursor"),
		Limit:  parseLimit(r.URL.Query().Get("limit")),
	}

	// Get feed from service
	feed, err := h.service.GetForYouFeed(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, feed)
}

// GetFollowingFeed handles requests for the "Following" feed
// GET /following?cursor=<cursor>&limit=<limit>
func (h *Handler) GetFollowingFeed(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters
	req := &FeedRequest{
		Cursor: r.URL.Query().Get("cursor"),
		Limit:  parseLimit(r.URL.Query().Get("limit")),
	}

	// Get feed from service
	feed, err := h.service.GetFollowingFeed(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, feed)
}

// GetVideo handles requests for a single video
// GET /:id
func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	// Get video ID from URL parameter
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Get video from service
	video, err := h.service.GetVideoByID(r.Context(), videoID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, VideoResponse{Video: video})
}

// parseLimit parses the limit query parameter
func parseLimit(limitStr string) int {
	if limitStr == "" {
		return 0 // Will use default in service
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0 // Will use default in service
	}
	return limit
}

// respondSuccess writes a successful JSON response
func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// respondError writes an error JSON response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
