package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/rainmana/gothink/internal/storage"
)

// SessionHandler handles session management operations
type SessionHandler struct {
	storage *storage.Storage
	logger  *logrus.Logger
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(storage *storage.Storage, logger *logrus.Logger) *SessionHandler {
	return &SessionHandler{
		storage: storage,
		logger:  logger,
	}
}

// GetStats handles session statistics requests
func (h *SessionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		h.respondWithError(w, "Session ID required", http.StatusBadRequest)
		return
	}

	stats, err := h.storage.GetSessionStats(sessionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session stats")
		h.respondWithError(w, "Failed to get session stats", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, stats)
}

// Export handles session export requests
func (h *SessionHandler) Export(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		h.respondWithError(w, "Session ID required", http.StatusBadRequest)
		return
	}

	export, err := h.storage.ExportSession(sessionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to export session")
		h.respondWithError(w, "Failed to export session", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, export)
}

// Import handles session import requests
func (h *SessionHandler) Import(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Session import not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Clear handles session clear requests
func (h *SessionHandler) Clear(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Session clear not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Helper methods

func (h *SessionHandler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *SessionHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
