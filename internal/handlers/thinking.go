package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/rainmana/gothink/internal/storage"
	"github.com/rainmana/gothink/internal/types"
)

// ThinkingHandler handles systematic thinking operations
type ThinkingHandler struct {
	storage *storage.Storage
	logger  *logrus.Logger
}

// NewThinkingHandler creates a new thinking handler
func NewThinkingHandler(storage *storage.Storage, logger *logrus.Logger) *ThinkingHandler {
	return &ThinkingHandler{
		storage: storage,
		logger:  logger,
	}
}

// SequentialThinking handles sequential thinking requests
func (h *ThinkingHandler) SequentialThinking(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID         string `json:"session_id"`
		Thought           string `json:"thought"`
		ThoughtNumber     int    `json:"thought_number"`
		TotalThoughts     int    `json:"total_thoughts"`
		NextThoughtNeeded bool   `json:"next_thought_needed"`
		IsRevision        bool   `json:"is_revision,omitempty"`
		RevisesThought    *int   `json:"revises_thought,omitempty"`
		BranchFromThought *int   `json:"branch_from_thought,omitempty"`
		BranchID          string `json:"branch_id,omitempty"`
		NeedsMoreThoughts bool   `json:"needs_more_thoughts,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create thought data
	thought := &types.ThoughtData{
		ID:                "",
		Thought:           request.Thought,
		ThoughtNumber:     request.ThoughtNumber,
		TotalThoughts:     request.TotalThoughts,
		IsRevision:        request.IsRevision,
		RevisesThought:    request.RevisesThought,
		BranchFromThought: request.BranchFromThought,
		BranchID:          request.BranchID,
		NeedsMoreThoughts: request.NeedsMoreThoughts,
		NextThoughtNeeded: request.NextThoughtNeeded,
		CreatedAt:         time.Now(),
	}

	// Add to storage
	if err := h.storage.AddThought(request.SessionID, thought); err != nil {
		h.logger.WithError(err).Error("Failed to add thought")
		h.respondWithError(w, "Failed to add thought", http.StatusInternalServerError)
		return
	}

	// Get session context
	stats, err := h.storage.GetSessionStats(request.SessionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session stats")
	}

	// Prepare response
	response := map[string]interface{}{
		"thought_id": thought.ID,
		"status":     "success",
		"session_context": map[string]interface{}{
			"session_id":         request.SessionID,
			"total_thoughts":     stats.ThoughtCount,
			"remaining_thoughts": stats.RemainingThoughts,
		},
	}

	h.respondWithJSON(w, response)
}

// MentalModel handles mental model application requests
func (h *ThinkingHandler) MentalModel(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID  string   `json:"session_id"`
		ModelName  string   `json:"model_name"`
		Problem    string   `json:"problem"`
		Steps      []string `json:"steps"`
		Reasoning  string   `json:"reasoning"`
		Conclusion string   `json:"conclusion"`
		Confidence float64  `json:"confidence,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate model name
	if _, exists := types.MentalModels[request.ModelName]; !exists {
		h.respondWithError(w, "Invalid mental model", http.StatusBadRequest)
		return
	}

	// Create mental model data
	model := &types.MentalModelData{
		ID:         "",
		ModelName:  request.ModelName,
		Problem:    request.Problem,
		Steps:      request.Steps,
		Reasoning:  request.Reasoning,
		Conclusion: request.Conclusion,
		Confidence: request.Confidence,
		CreatedAt:  time.Now(),
	}

	// Add to storage
	if err := h.storage.AddMentalModel(request.SessionID, model); err != nil {
		h.logger.WithError(err).Error("Failed to add mental model")
		h.respondWithError(w, "Failed to add mental model", http.StatusInternalServerError)
		return
	}

	// Get session context
	stats, err := h.storage.GetSessionStats(request.SessionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session stats")
	}

	// Prepare response
	response := map[string]interface{}{
		"model_id":       model.ID,
		"status":         "success",
		"has_steps":      len(request.Steps) > 0,
		"has_conclusion": request.Conclusion != "",
		"session_context": map[string]interface{}{
			"session_id":          request.SessionID,
			"total_mental_models": stats.Stores["mental_models"].(map[string]int)["count"],
		},
	}

	h.respondWithJSON(w, response)
}

// DebuggingApproach handles debugging approach requests
func (h *ThinkingHandler) DebuggingApproach(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID    string   `json:"session_id"`
		ApproachName string   `json:"approach_name"`
		Issue        string   `json:"issue"`
		Steps        []string `json:"steps"`
		Findings     string   `json:"findings"`
		Resolution   string   `json:"resolution"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For now, we'll store this as a mental model with a special type
	model := &types.MentalModelData{
		ID:         "",
		ModelName:  "debugging_" + request.ApproachName,
		Problem:    request.Issue,
		Steps:      request.Steps,
		Reasoning:  request.Findings,
		Conclusion: request.Resolution,
		CreatedAt:  time.Now(),
	}

	// Add to storage
	if err := h.storage.AddMentalModel(request.SessionID, model); err != nil {
		h.logger.WithError(err).Error("Failed to add debugging approach")
		h.respondWithError(w, "Failed to add debugging approach", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"approach_id":    model.ID,
		"status":         "success",
		"has_findings":   request.Findings != "",
		"has_resolution": request.Resolution != "",
	}

	h.respondWithJSON(w, response)
}

// CollaborativeReasoning handles collaborative reasoning requests
func (h *ThinkingHandler) CollaborativeReasoning(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Collaborative reasoning not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// SocraticMethod handles Socratic method requests
func (h *ThinkingHandler) SocraticMethod(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Socratic method not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// CreativeThinking handles creative thinking requests
func (h *ThinkingHandler) CreativeThinking(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Creative thinking not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// SystemsThinking handles systems thinking requests
func (h *ThinkingHandler) SystemsThinking(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Systems thinking not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// ScientificMethod handles scientific method requests
func (h *ThinkingHandler) ScientificMethod(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Scientific method not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Helper methods

func (h *ThinkingHandler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *ThinkingHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
