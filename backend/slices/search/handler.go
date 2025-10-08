package search

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Search handles search requests
// GET /search?q=query&type=users|videos|hashtags&cursor=&limit=20
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	searchType := SearchType(r.URL.Query().Get("type"))
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	// Parse limit
	limit := 20 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create request
	req := &SearchRequest{
		Query:  query,
		Type:   searchType,
		Cursor: cursor,
		Limit:  limit,
	}

	// Perform search
	response, err := h.service.Search(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// GetTrendingHashtags handles trending hashtags requests
// GET /trending/hashtags?limit=20
func (h *Handler) GetTrendingHashtags(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")

	// Parse limit
	limit := 20 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get trending hashtags
	response, err := h.service.GetTrendingHashtags(r.Context(), limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// GetVideosByHashtag handles hashtag videos requests
// GET /hashtags/:tag/videos?cursor=&limit=20
func (h *Handler) GetVideosByHashtag(w http.ResponseWriter, r *http.Request) {
	// Parse URL parameters
	tag := chi.URLParam(r, "tag")

	// Parse query parameters
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	// Parse limit
	limit := 20 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create request
	req := &HashtagVideosRequest{
		Tag:    tag,
		Cursor: cursor,
		Limit:  limit,
	}

	// Get videos by hashtag
	response, err := h.service.GetVideosByHashtag(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, response)
}

// Helper functions for HTTP responses

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
