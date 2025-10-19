package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/waldzellai/gothink/internal/storage"
	"github.com/waldzellai/gothink/internal/types"
)

// VisualHandler handles visualization operations
type VisualHandler struct {
	storage *storage.Storage
	logger  *logrus.Logger
}

// NewVisualHandler creates a new visual handler
func NewVisualHandler(storage *storage.Storage, logger *logrus.Logger) *VisualHandler {
	return &VisualHandler{
		storage: storage,
		logger:  logger,
	}
}

// ConceptMap handles concept map requests
func (h *VisualHandler) ConceptMap(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID           string                `json:"session_id"`
		DiagramID           string                `json:"diagram_id"`
		Operation           string                `json:"operation"`
		Elements            []types.VisualElement `json:"elements,omitempty"`
		Iteration           int                   `json:"iteration"`
		Observation         string                `json:"observation,omitempty"`
		Insight             string                `json:"insight,omitempty"`
		Hypothesis          string                `json:"hypothesis,omitempty"`
		NextOperationNeeded bool                  `json:"next_operation_needed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create visual data
	visual := &types.VisualData{
		ID:                  "",
		Operation:           request.Operation,
		Elements:            request.Elements,
		DiagramID:           request.DiagramID,
		DiagramType:         "concept-map",
		Iteration:           request.Iteration,
		Observation:         request.Observation,
		Insight:             request.Insight,
		Hypothesis:          request.Hypothesis,
		NextOperationNeeded: request.NextOperationNeeded,
		CreatedAt:           time.Now(),
	}

	// Add to storage
	if err := h.storage.AddVisualData(request.SessionID, visual); err != nil {
		h.logger.WithError(err).Error("Failed to add visual data")
		h.respondWithError(w, "Failed to add visual data", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"visual_id":    visual.ID,
		"status":       "success",
		"diagram_type": "concept-map",
		"operation":    request.Operation,
		"elements":     len(request.Elements),
	}

	h.respondWithJSON(w, response)
}

// MindMap handles mind map requests
func (h *VisualHandler) MindMap(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Mind map not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Flowchart handles flowchart requests
func (h *VisualHandler) Flowchart(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Flowchart not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// DecisionTree handles decision tree requests
func (h *VisualHandler) DecisionTree(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Decision tree not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// ProbabilityTree handles probability tree requests
func (h *VisualHandler) ProbabilityTree(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Probability tree not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// BayesianNetwork handles Bayesian network requests
func (h *VisualHandler) BayesianNetwork(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Bayesian network not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Helper methods

func (h *VisualHandler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *VisualHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
