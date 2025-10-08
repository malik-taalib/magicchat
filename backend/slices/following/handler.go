package following

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

// FollowUser handles POST /:id/follow
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from context
	followerID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user ID to follow from URL parameter
	followingID := chi.URLParam(r, "id")
	if followingID == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Call service to follow user
	response, err := h.service.FollowUser(r.Context(), followerID, followingID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// UnfollowUser handles DELETE /:id/follow
func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from context
	followerID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user ID to unfollow from URL parameter
	followingID := chi.URLParam(r, "id")
	if followingID == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Call service to unfollow user
	response, err := h.service.UnfollowUser(r.Context(), followerID, followingID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// GetFollowers handles GET /:id/followers
func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from URL parameter
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Get pagination parameters from query string
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			respondError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	// Call service to get followers
	response, err := h.service.GetFollowers(r.Context(), userID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// GetFollowing handles GET /:id/following
func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from URL parameter
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Get pagination parameters from query string
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			respondError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	// Call service to get following
	response, err := h.service.GetFollowing(r.Context(), userID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// IsFollowing handles GET /:id/following/check (optional endpoint to check follow status)
func (h *Handler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from context
	followerID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the user ID to check from URL parameter
	followingID := chi.URLParam(r, "id")
	if followingID == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Call service to check if following
	isFollowing, err := h.service.IsFollowing(r.Context(), followerID, followingID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]interface{}{
		"is_following": isFollowing,
		"user_id":      followingID,
	})
}

// Helper functions
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
