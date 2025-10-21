package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/rainmana/gothink/internal/storage"
	"github.com/rainmana/gothink/internal/types"
)

// DecisionHandler handles decision framework operations
type DecisionHandler struct {
	storage *storage.Storage
	logger  *logrus.Logger
}

// NewDecisionHandler creates a new decision handler
func NewDecisionHandler(storage *storage.Storage, logger *logrus.Logger) *DecisionHandler {
	return &DecisionHandler{
		storage: storage,
		logger:  logger,
	}
}

// DecisionFramework handles decision framework requests
func (h *DecisionHandler) DecisionFramework(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SessionID         string                    `json:"session_id"`
		DecisionStatement string                    `json:"decision_statement"`
		Options           []types.DecisionOption    `json:"options"`
		Criteria          []types.DecisionCriterion `json:"criteria,omitempty"`
		Stakeholders      []string                  `json:"stakeholders,omitempty"`
		Constraints       []string                  `json:"constraints,omitempty"`
		TimeHorizon       string                    `json:"time_horizon,omitempty"`
		RiskTolerance     string                    `json:"risk_tolerance,omitempty"`
		AnalysisType      string                    `json:"analysis_type"`
		Stage             string                    `json:"stage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create decision data
	decision := &types.DecisionData{
		ID:                "",
		DecisionStatement: request.DecisionStatement,
		Options:           request.Options,
		Criteria:          request.Criteria,
		Stakeholders:      request.Stakeholders,
		Constraints:       request.Constraints,
		TimeHorizon:       request.TimeHorizon,
		RiskTolerance:     request.RiskTolerance,
		AnalysisType:      request.AnalysisType,
		Stage:             request.Stage,
		Iteration:         1,
		NextStageNeeded:   true,
		CreatedAt:         time.Now(),
	}

	// Add to storage
	if err := h.storage.AddDecision(request.SessionID, decision); err != nil {
		h.logger.WithError(err).Error("Failed to add decision")
		h.respondWithError(w, "Failed to add decision", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"decision_id":   decision.ID,
		"status":        "success",
		"has_options":   len(request.Options) > 0,
		"has_criteria":  len(request.Criteria) > 0,
		"analysis_type": request.AnalysisType,
		"stage":         request.Stage,
	}

	h.respondWithJSON(w, response)
}

// ExpectedUtility handles expected utility analysis requests
func (h *DecisionHandler) ExpectedUtility(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Expected utility analysis not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// MultiCriteria handles multi-criteria analysis requests
func (h *DecisionHandler) MultiCriteria(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Multi-criteria analysis not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// RiskAnalysis handles risk analysis requests
func (h *DecisionHandler) RiskAnalysis(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	response := map[string]interface{}{
		"message": "Risk analysis not yet implemented",
		"status":  "coming_soon",
	}
	h.respondWithJSON(w, response)
}

// Helper methods

func (h *DecisionHandler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *DecisionHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
