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
// Stochastic Algorithm Types
// ============================================================================

// StochasticAlgorithmData represents the application of a stochastic algorithm
type StochasticAlgorithmData struct {
	ID         string                 `json:"id"`
	Algorithm  string                 `json:"algorithm"`
	Problem    string                 `json:"problem"`
	Parameters map[string]interface{} `json:"parameters"`
	Result     string                 `json:"result,omitempty"`
	Confidence float64                `json:"confidence,omitempty"`
	Iterations int                    `json:"iterations,omitempty"`
	Converged  bool                   `json:"converged,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// MDPData represents Markov Decision Process specific data
type MDPData struct {
	StochasticAlgorithmData
	Policy        map[string]string             `json:"policy,omitempty"`
	ValueFunction map[string]float64            `json:"value_function,omitempty"`
	QValues       map[string]map[string]float64 `json:"q_values,omitempty"`
}

// MCTSData represents Monte Carlo Tree Search specific data
type MCTSData struct {
	StochasticAlgorithmData
	BestAction string                 `json:"best_action,omitempty"`
	TreeStats  map[string]interface{} `json:"tree_stats,omitempty"`
}

// BanditData represents Multi-Armed Bandit specific data
type BanditData struct {
	StochasticAlgorithmData
	ArmStats    []ArmStatistics `json:"arm_stats,omitempty"`
	SelectedArm int             `json:"selected_arm,omitempty"`
}

// ArmStatistics represents statistics for a bandit arm
type ArmStatistics struct {
	Arm           int     `json:"arm"`
	Pulls         int     `json:"pulls"`
	Rewards       float64 `json:"rewards"`
	AverageReward float64 `json:"average_reward"`
}

// BayesianOptimizationData represents Bayesian Optimization specific data
type BayesianOptimizationData struct {
	StochasticAlgorithmData
	OptimizationHistory []OptimizationStep `json:"optimization_history,omitempty"`
	BestParameters      map[string]float64 `json:"best_parameters,omitempty"`
	BestValue           float64            `json:"best_value,omitempty"`
}

// OptimizationStep represents a step in Bayesian optimization
type OptimizationStep struct {
	Iteration  int                `json:"iteration"`
	Parameters map[string]float64 `json:"parameters"`
	Value      float64            `json:"value"`
}

// HMMData represents Hidden Markov Model specific data
type HMMData struct {
	StochasticAlgorithmData
	StateSequence           []int       `json:"state_sequence,omitempty"`
	TransitionProbabilities [][]float64 `json:"transition_probabilities,omitempty"`
	EmissionProbabilities   [][]float64 `json:"emission_probabilities,omitempty"`
	InitialProbabilities    []float64   `json:"initial_probabilities,omitempty"`
}

// ============================================================================
// Decision Framework Types
// ============================================================================

// DecisionOption represents an option in a decision
type DecisionOption struct {
	ID                   string  `json:"id,omitempty"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	ExpectedValue        float64 `json:"expected_value,omitempty"`
	RiskLevel            string  `json:"risk_level,omitempty"`
	ProbabilityOfSuccess float64 `json:"probability_of_success,omitempty"`
}

// DecisionCriterion represents a criterion for evaluating options
type DecisionCriterion struct {
	ID               string  `json:"id,omitempty"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Weight           float64 `json:"weight"`
	EvaluationMethod string  `json:"evaluation_method"`
}

// DecisionData represents a complete decision framework
type DecisionData struct {
	ID                string              `json:"id"`
	DecisionStatement string              `json:"decision_statement"`
	Options           []DecisionOption    `json:"options"`
	Criteria          []DecisionCriterion `json:"criteria,omitempty"`
	Stakeholders      []string            `json:"stakeholders,omitempty"`
	Constraints       []string            `json:"constraints,omitempty"`
	TimeHorizon       string              `json:"time_horizon,omitempty"`
	RiskTolerance     string              `json:"risk_tolerance,omitempty"`
	AnalysisType      string              `json:"analysis_type"`
	Stage             string              `json:"stage"`
	Recommendation    string              `json:"recommendation,omitempty"`
	Iteration         int                 `json:"iteration"`
	NextStageNeeded   bool                `json:"next_stage_needed"`
	CreatedAt         time.Time           `json:"created_at"`
}

// ============================================================================
// Visualization Types
// ============================================================================

// VisualElement represents an element in a diagram
type VisualElement struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Label       string                 `json:"label,omitempty"`
	Properties  map[string]interface{} `json:"properties"`
	Source      string                 `json:"source,omitempty"`
	Target      string                 `json:"target,omitempty"`
	Contains    []string               `json:"contains,omitempty"`
	Probability float64                `json:"probability,omitempty"`
}

// VisualData represents a visual reasoning operation
type VisualData struct {
	ID                  string          `json:"id"`
	Operation           string          `json:"operation"`
	Elements            []VisualElement `json:"elements,omitempty"`
	TransformationType  string          `json:"transformation_type,omitempty"`
	DiagramID           string          `json:"diagram_id"`
	DiagramType         string          `json:"diagram_type"`
	Iteration           int             `json:"iteration"`
	Observation         string          `json:"observation,omitempty"`
	Insight             string          `json:"insight,omitempty"`
	Hypothesis          string          `json:"hypothesis,omitempty"`
	NextOperationNeeded bool            `json:"next_operation_needed"`
	CreatedAt           time.Time       `json:"created_at"`
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
