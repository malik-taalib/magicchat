package videoupload

import (
	"encoding/json"
	"net/http"

	"magicchat/slices/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse multipart form (max 100MB)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("video")
	if err != nil {
		respondError(w, http.StatusBadRequest, "video file is required")
		return
	}
	defer file.Close()

	// Parse request
	req := &UploadRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Hashtags:    []string{}, // TODO: Parse hashtags from form
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	// Upload video
	video, err := h.service.UploadVideo(r.Context(), userID, req, fileHeader)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, UploadResponse{
		VideoID: video.ID.Hex(),
		Status:  video.ProcessingStatus,
	})
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "video ID is required")
		return
	}

	video, err := h.service.GetVideoStatus(r.Context(), videoID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, StatusResponse{
		VideoID:  video.ID.Hex(),
		Status:   video.ProcessingStatus,
		VideoURL: video.VideoURL,
	})
}

func (h *Handler) ProcessWebhook(w http.ResponseWriter, r *http.Request) {
	// This endpoint would be called by video processing workers
	var req struct {
		VideoID string `json:"video_id"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.ProcessVideo(r.Context(), req.VideoID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{"message": "video processing completed"})
}

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
