package notifications

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"magicchat/slices/auth"
)

type Handler struct {
	service   *Service
	wsManager *WebSocketManager
}

func NewHandler(service *Service, wsManager *WebSocketManager) *Handler {
	return &Handler{
		service:   service,
		wsManager: wsManager,
	}
}

// GetNotifications retrieves the notification list for the authenticated user
// GET /notifications?cursor=<id>&limit=20
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	// Get notifications
	result, err := h.service.GetNotifications(r.Context(), userID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to retrieve notifications")
		return
	}

	respondSuccess(w, http.StatusOK, result)
}

// MarkAsRead marks a single notification as read
// PUT /notifications/:id/read
func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notificationID := chi.URLParam(r, "id")
	if notificationID == "" {
		respondError(w, http.StatusBadRequest, "notification ID is required")
		return
	}

	err := h.service.MarkAsRead(r.Context(), notificationID, userID)
	if err != nil {
		if err.Error() == "notification not found or unauthorized" {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "notification marked as read",
		"id":      notificationID,
	})
}

// MarkAllAsRead marks all notifications as read for the authenticated user
// PUT /notifications/read-all
func (h *Handler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.service.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark all notifications as read")
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "all notifications marked as read",
	})
}

// GetUnreadCount returns the count of unread notifications
// GET /notifications/unread-count
func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	count, err := h.service.GetUnreadCount(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get unread count")
		return
	}

	respondSuccess(w, http.StatusOK, map[string]int{
		"unread_count": count,
	})
}

// HandleWebSocket handles WebSocket connections for real-time notifications
// WS /notifications/stream
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Upgrade connection to WebSocket
	h.wsManager.ServeWS(w, r, userID)
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
