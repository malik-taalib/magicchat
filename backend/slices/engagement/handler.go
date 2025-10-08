package engagement

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

// LikeVideo handles POST /:id/like
func (h *Handler) LikeVideo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get video ID from URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Like video
	response, err := h.service.LikeVideo(r.Context(), userID, videoID)
	if err != nil {
		if err.Error() == "video already liked" {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// UnlikeVideo handles DELETE /:id/like
func (h *Handler) UnlikeVideo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get video ID from URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Unlike video
	response, err := h.service.UnlikeVideo(r.Context(), userID, videoID)
	if err != nil {
		if err.Error() == "video not liked" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// CreateComment handles POST /:id/comments
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get video ID from URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Parse request body
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create comment
	response, err := h.service.CreateComment(r.Context(), userID, videoID, &req)
	if err != nil {
		if err.Error() == "comment text is required" || err.Error() == "comment text must be less than 500 characters" {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "parent comment not found" || err.Error() == "parent comment does not belong to this video" {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, response)
}

// GetComments handles GET /:id/comments
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	// Get video ID from URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int64(20) // Default limit
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	offset := int64(0) // Default offset
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			offset = o
		}
	}

	// Get comments
	comments, err := h.service.GetComments(r.Context(), videoID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]interface{}{
		"comments": comments,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetCommentReplies handles GET /comments/:id/replies
func (h *Handler) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	// Get comment ID from URL
	commentID := chi.URLParam(r, "id")
	if commentID == "" {
		respondError(w, http.StatusBadRequest, "comment ID is required")
		return
	}

	// Get replies
	replies, err := h.service.GetCommentReplies(r.Context(), commentID)
	if err != nil {
		if err.Error() == "comment not found" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]interface{}{
		"replies": replies,
	})
}

// ShareVideo handles POST /:id/share
func (h *Handler) ShareVideo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get video ID from URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	// Record share
	response, err := h.service.RecordShare(r.Context(), userID, videoID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// Helper functions for consistent response formatting

func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
