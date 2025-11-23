package types

import "time"

// ============================================================================
// Core Thinking Types
// ============================================================================

// ThoughtData represents a single thought in a sequential thinking process
type ThoughtData struct {
	ID                string    `json:"id"`
	Thought           string    `json:"thought"`
	ThoughtNumber     int       `json:"thought_number"`
	TotalThoughts     int       `json:"total_thoughts"`
	IsRevision        bool      `json:"is_revision,omitempty"`
	RevisesThought    *int      `json:"revises_thought,omitempty"`
	BranchFromThought *int      `json:"branch_from_thought,omitempty"`
	BranchID          string    `json:"branch_id,omitempty"`
	NeedsMoreThoughts bool      `json:"needs_more_thoughts,omitempty"`
	NextThoughtNeeded bool      `json:"next_thought_needed"`
	CreatedAt         time.Time `json:"created_at"`
}

// MentalModelData represents the application of a mental model to a problem
type MentalModelData struct {
	ID         string    `json:"id"`
	ModelName  string    `json:"model_name"`
	Problem    string    `json:"problem"`
	Steps      []string  `json:"steps"`
	Reasoning  string    `json:"reasoning"`
	Conclusion string    `json:"conclusion"`
	Confidence float64   `json:"confidence,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// ============================================================================
// Session Management Types
// ============================================================================

// SessionExport represents exported session data
type SessionExport struct {
	Version     string                 `json:"version"`
	Timestamp   time.Time              `json:"timestamp"`
	SessionID   string                 `json:"session_id"`
	SessionType string                 `json:"session_type"`
	Data        interface{}            `json:"data"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessResult represents the result of processing a thinking operation
type ProcessResult struct {
	Success bool `json:"success"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error,omitempty"`
}

// SessionStatistics represents comprehensive session statistics
type SessionStatistics struct {
	SessionID         string                 `json:"session_id"`
	CreatedAt         time.Time              `json:"created_at"`
	LastAccessedAt    time.Time              `json:"last_accessed_at"`
	ThoughtCount      int                    `json:"thought_count"`
	ToolsUsed         []string               `json:"tools_used"`
	TotalOperations   int                    `json:"total_operations"`
	IsActive          bool                   `json:"is_active"`
	RemainingThoughts int                    `json:"remaining_thoughts"`
	Stores            map[string]interface{} `json:"stores"`
}

// ============================================================================
// Tool Request/Response Types
// ============================================================================

// ToolRequest represents a request to execute a tool
type ToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResponse represents a response from a tool execution
type ToolResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"is_error,omitempty"`
}

// ============================================================================
// Mental Model Types
// ============================================================================

// MentalModel represents a specific mental model
type MentalModel struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Examples    []string `json:"examples"`
	Category    string   `json:"category"`
}

// Available mental models
var MentalModels = map[string]MentalModel{
	"first_principles": {
		Name:        "First Principles Thinking",
		Description: "Break down complex problems into fundamental components",
		Steps: []string{
			"Identify the problem clearly",
			"Break it down into basic components",
			"Question assumptions",
			"Build up from the basics",
		},
		Category: "analytical",
	},
	"opportunity_cost": {
		Name:        "Opportunity Cost Analysis",
		Description: "Consider what you give up when making a choice",
		Steps: []string{
			"Identify all available options",
			"List the benefits of each option",
			"Identify what you give up with each choice",
			"Compare opportunity costs",
		},
		Category: "decision-making",
	},
	"bayesian_thinking": {
		Name:        "Bayesian Thinking",
		Description: "Update beliefs based on new evidence",
		Steps: []string{
			"Start with prior beliefs",
			"Gather new evidence",
			"Update beliefs using Bayes' theorem",
			"Consider alternative explanations",
		},
		Category: "probabilistic",
	},
	"systems_thinking": {
		Name:        "Systems Thinking",
		Description: "Understand how parts of a system interact",
		Steps: []string{
			"Identify system boundaries",
			"Map system components",
			"Identify relationships and feedback loops",
			"Consider emergent properties",
		},
		Category: "holistic",
	},
}
