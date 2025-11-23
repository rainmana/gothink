package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/rainmana/gothink/internal/config"
	"github.com/rainmana/gothink/internal/types"
	"github.com/sirupsen/logrus"
)

// Storage manages all data storage for the GoThink server
type Storage struct {
	config *config.Config
	logger *logrus.Logger

	// In-memory stores (in production, these would be backed by a database)
	thoughts     map[string]*types.ThoughtData
	mentalModels map[string]*types.MentalModelData
	sessions     map[string]*SessionData

	// Mutexes for thread safety
	thoughtsMutex     sync.RWMutex
	mentalModelsMutex sync.RWMutex
	sessionsMutex     sync.RWMutex
}

// SessionData represents session-specific data
type SessionData struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	LastAccessedAt    time.Time `json:"last_accessed_at"`
	ThoughtCount      int       `json:"thought_count"`
	ToolsUsed         []string  `json:"tools_used"`
	TotalOperations   int       `json:"total_operations"`
	IsActive          bool      `json:"is_active"`
	RemainingThoughts int       `json:"remaining_thoughts"`
}

// New creates a new storage instance
func New(cfg *config.Config) (*Storage, error) {

	return &Storage{
		config:       cfg,
		logger:       logrus.New(),
		thoughts:     make(map[string]*types.ThoughtData),
		mentalModels: make(map[string]*types.MentalModelData),
		sessions:     make(map[string]*SessionData),
	}, nil
}

// ============================================================================
// Thought Management
// ============================================================================

// AddThought adds a new thought to storage
func (s *Storage) AddThought(sessionID string, thought *types.ThoughtData) error {
	s.thoughtsMutex.Lock()
	defer s.thoughtsMutex.Unlock()

	// Check thought limit
	session := s.getSession(sessionID)
	if session.ThoughtCount >= s.config.MaxThoughtsPerSession {
		return fmt.Errorf("thought limit reached for session %s", sessionID)
	}

	// Generate ID if not provided
	if thought.ID == "" {
		thought.ID = generateID()
	}
	thought.CreatedAt = time.Now()

	s.thoughts[thought.ID] = thought

	// Update session
	session.ThoughtCount++
	session.LastAccessedAt = time.Now()
	s.sessions[sessionID] = session

	s.logger.WithFields(logrus.Fields{
		"session_id":     sessionID,
		"thought_id":     thought.ID,
		"thought_number": thought.ThoughtNumber,
	}).Debug("Added thought to storage")

	return nil
}

// GetThoughts retrieves all thoughts for a session
func (s *Storage) GetThoughts(sessionID string) ([]*types.ThoughtData, error) {
	s.thoughtsMutex.RLock()
	defer s.thoughtsMutex.RUnlock()

	var sessionThoughts []*types.ThoughtData
	for _, thought := range s.thoughts {
		// In a real implementation, you'd filter by session ID
		sessionThoughts = append(sessionThoughts, thought)
	}

	return sessionThoughts, nil
}

// ============================================================================
// Mental Model Management
// ============================================================================

// AddMentalModel adds a mental model application to storage
func (s *Storage) AddMentalModel(sessionID string, model *types.MentalModelData) error {
	s.mentalModelsMutex.Lock()
	defer s.mentalModelsMutex.Unlock()

	if model.ID == "" {
		model.ID = generateID()
	}
	model.CreatedAt = time.Now()

	s.mentalModels[model.ID] = model

	// Update session
	session := s.getSession(sessionID)
	session.LastAccessedAt = time.Now()
	s.sessions[sessionID] = session

	s.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"model_id":   model.ID,
		"model_name": model.ModelName,
	}).Debug("Added mental model to storage")

	return nil
}

// GetMentalModels retrieves all mental models for a session
func (s *Storage) GetMentalModels(sessionID string) ([]*types.MentalModelData, error) {
	s.mentalModelsMutex.RLock()
	defer s.mentalModelsMutex.RUnlock()

	var sessionModels []*types.MentalModelData
	for _, model := range s.mentalModels {
		sessionModels = append(sessionModels, model)
	}

	return sessionModels, nil
}

// ============================================================================
// Session Management
// ============================================================================

// GetSession retrieves session data
func (s *Storage) GetSession(sessionID string) (*SessionData, error) {
	s.sessionsMutex.RLock()
	defer s.sessionsMutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	return session, nil
}

// CreateSession creates a new session
func (s *Storage) CreateSession(sessionID string) (*SessionData, error) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	session := &SessionData{
		ID:                sessionID,
		CreatedAt:         time.Now(),
		LastAccessedAt:    time.Now(),
		ThoughtCount:      0,
		ToolsUsed:         []string{},
		TotalOperations:   0,
		IsActive:          true,
		RemainingThoughts: s.config.MaxThoughtsPerSession,
	}

	s.sessions[sessionID] = session

	s.logger.WithField("session_id", sessionID).Debug("Created new session")

	return session, nil
}

// getSession gets or creates a session
func (s *Storage) getSession(sessionID string) *SessionData {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		session = &SessionData{
			ID:                sessionID,
			CreatedAt:         time.Now(),
			LastAccessedAt:    time.Now(),
			ThoughtCount:      0,
			ToolsUsed:         []string{},
			TotalOperations:   0,
			IsActive:          true,
			RemainingThoughts: s.config.MaxThoughtsPerSession,
		}
		s.sessions[sessionID] = session
	}

	return session
}

// GetSessionStats retrieves comprehensive session statistics
func (s *Storage) GetSessionStats(sessionID string) (*types.SessionStatistics, error) {
	session := s.getSession(sessionID)

	thoughts, _ := s.GetThoughts(sessionID)
	mentalModels, _ := s.GetMentalModels(sessionID)

	// Collect tools used
	toolsUsed := make(map[string]bool)
	if len(thoughts) > 0 {
		toolsUsed["sequential-thinking"] = true
	}
	if len(mentalModels) > 0 {
		toolsUsed["mental-model"] = true
	}

	var toolsList []string
	for tool := range toolsUsed {
		toolsList = append(toolsList, tool)
	}

	stats := &types.SessionStatistics{
		SessionID:         sessionID,
		CreatedAt:         session.CreatedAt,
		LastAccessedAt:    session.LastAccessedAt,
		ThoughtCount:      len(thoughts),
		ToolsUsed:         toolsList,
		TotalOperations:   len(thoughts) + len(mentalModels),
		IsActive:          session.IsActive,
		RemainingThoughts: s.config.MaxThoughtsPerSession - len(thoughts),
		Stores: map[string]interface{}{
			"thoughts":      map[string]int{"count": len(thoughts)},
			"mental_models": map[string]int{"count": len(mentalModels)},
		},
	}

	return stats, nil
}

// ============================================================================
// Export/Import
// ============================================================================

// ExportSession exports session data
func (s *Storage) ExportSession(sessionID string) (*types.SessionExport, error) {
	thoughts, _ := s.GetThoughts(sessionID)
	mentalModels, _ := s.GetMentalModels(sessionID)

	export := &types.SessionExport{
		Version:     "1.0.0",
		Timestamp:   time.Now(),
		SessionID:   sessionID,
		SessionType: "hybrid",
		Data: map[string]interface{}{
			"thoughts":      thoughts,
			"mental_models": mentalModels,
		},
		Metadata: map[string]interface{}{
			"exported_at": time.Now(),
			"version":     "0.1.0",
		},
	}

	return export, nil
}

// ============================================================================
// Utility Functions
// ============================================================================

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}
