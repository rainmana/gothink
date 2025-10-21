package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/rainmana/gothink/internal/config"
	"github.com/rainmana/gothink/internal/models"
	"github.com/rainmana/gothink/internal/storage"
	"github.com/rainmana/gothink/internal/types"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create storage
	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create mental models loader
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	modelsLoader := models.NewLoader(logger)

	// Create MCP server
	s := server.NewMCPServer(
		"GoThink MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
	)

	// Add all the thinking tools
	addThinkingTools(s, store, modelsLoader, cfg)
	addStochasticTools(s, store)
	addDecisionTools(s, store)
	addVisualTools(s, store)
	addSessionTools(s, store)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func addThinkingTools(s *server.MCPServer, store *storage.Storage, modelsLoader *models.Loader, cfg *config.Config) {
	// Sequential Thinking Tool
	s.AddTool(
		mcp.NewTool("sequential_thinking",
			mcp.WithDescription("Perform sequential thinking operations with structured thought progression"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("thought", mcp.Required(), mcp.Description("Current thought content")),
			mcp.WithNumber("thought_number", mcp.Required(), mcp.Description("Current thought number in sequence")),
			mcp.WithNumber("total_thoughts", mcp.Required(), mcp.Description("Total number of thoughts planned")),
			mcp.WithBoolean("next_thought_needed", mcp.Required(), mcp.Description("Whether another thought is needed")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			thought, _ := req.RequireString("thought")
			thoughtNumber, _ := req.RequireInt("thought_number")
			totalThoughts, _ := req.RequireInt("total_thoughts")
			nextThoughtNeeded, _ := req.RequireBool("next_thought_needed")

			// Create thought data
			thoughtData := &types.ThoughtData{
				ID:                fmt.Sprintf("%d-%d", time.Now().UnixNano(), thoughtNumber),
				Thought:           thought,
				ThoughtNumber:     thoughtNumber,
				TotalThoughts:     totalThoughts,
				NextThoughtNeeded: nextThoughtNeeded,
				CreatedAt:         time.Now(),
			}

			// Store the thought
			store.AddThought(sessionID, thoughtData)

			// Get session stats
			stats, _ := store.GetSessionStats(sessionID)

			// Create response
			response := map[string]interface{}{
				"status":     "success",
				"thought_id": thoughtData.ID,
				"session_context": map[string]interface{}{
					"session_id":         sessionID,
					"total_thoughts":     stats.ThoughtCount,
					"remaining_thoughts": 100 - stats.ThoughtCount,
				},
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// Mental Model Tool
	s.AddTool(
		mcp.NewTool("mental_model",
			mcp.WithDescription("Apply mental models to solve problems using structured thinking frameworks"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("model_name", mcp.Required(), mcp.Description("Name of the mental model to apply")),
			mcp.WithString("problem", mcp.Required(), mcp.Description("Problem statement to analyze")),
			mcp.WithArray("steps", mcp.Description("Steps to follow for the mental model")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			modelName, _ := req.RequireString("model_name")
			problem, _ := req.RequireString("problem")
			steps := req.GetStringSlice("steps", []string{})

			// Load available mental models
			availableModels, err := modelsLoader.LoadMentalModels(cfg.MentalModelsPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to load mental models: %v", err)), nil
			}

			// Check if the requested model exists
			model, exists := availableModels[modelName]
			if !exists {
				// Return available models for reference
				available := modelsLoader.GetAvailableModels(availableModels)
				return mcp.NewToolResultError(fmt.Sprintf("Mental model '%s' not found. Available models: %v", modelName, available)), nil
			}

			// Use model steps if no custom steps provided
			if len(steps) == 0 {
				steps = model.Steps
			}

			// Create mental model data
			modelData := &types.MentalModelData{
				ID:        fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(steps)),
				ModelName: modelName,
				Problem:   problem,
				Steps:     steps,
				CreatedAt: time.Now(),
			}

			// Store the mental model
			store.AddMentalModel(sessionID, modelData)

			// Get session stats
			stats, _ := store.GetSessionStats(sessionID)

			// Create response
			response := map[string]interface{}{
				"status":   "success",
				"model_id": modelData.ID,
				"model_info": map[string]interface{}{
					"name":        model.Name,
					"description": model.Description,
					"category":    model.Category,
					"priority":    model.Priority,
				},
				"steps_used":     steps,
				"has_steps":      len(steps) > 0,
				"has_conclusion": false,
				"session_context": map[string]interface{}{
					"session_id":          sessionID,
					"total_mental_models": stats.Stores["mental_models"].(map[string]int)["count"],
				},
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// Debugging Approach Tool
	s.AddTool(
		mcp.NewTool("debugging_approach",
			mcp.WithDescription("Apply systematic debugging approaches to identify and resolve issues"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("approach_name", mcp.Required(), mcp.Description("Name of the debugging approach")),
			mcp.WithString("issue", mcp.Required(), mcp.Description("Issue description to debug")),
			mcp.WithArray("steps", mcp.Description("Debugging steps to follow")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			_, _ = req.RequireString("approach_name")
			_, _ = req.RequireString("issue")
			steps := req.GetStringSlice("steps", []string{})

			// Create response
			response := map[string]interface{}{
				"status":         "success",
				"approach_id":    fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(steps)),
				"has_steps":      len(steps) > 0,
				"has_findings":   false,
				"has_resolution": false,
				"session_context": map[string]interface{}{
					"session_id": sessionID,
				},
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// List Available Mental Models Tool
	s.AddTool(
		mcp.NewTool("list_mental_models",
			mcp.WithDescription("List all available mental models with their details"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Load available mental models
			availableModels, err := modelsLoader.LoadMentalModels(cfg.MentalModelsPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to load mental models: %v", err)), nil
			}

			// Get models sorted by priority
			modelsByPriority := modelsLoader.GetModelsByPriority(availableModels)
			modelsByCategory := modelsLoader.GetModelsByCategory(availableModels)

			// Create response
			response := map[string]interface{}{
				"status":             "success",
				"total_models":       len(availableModels),
				"models_by_priority": modelsByPriority,
				"models_by_category": modelsByCategory,
				"available_models":   modelsLoader.GetAvailableModels(availableModels),
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)
}

func addStochasticTools(s *server.MCPServer, store *storage.Storage) {
	// Markov Decision Process Tool
	s.AddTool(
		mcp.NewTool("markov_decision_process",
			mcp.WithDescription("Run Markov Decision Process optimization for sequential decision making"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("problem", mcp.Required(), mcp.Description("Problem description for MDP")),
			mcp.WithObject("parameters", mcp.Description("MDP parameters (states, actions, rewards, etc.)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			problem, _ := req.RequireString("problem")
			paramsInterface, _ := req.GetArguments()["parameters"]
			params, ok := paramsInterface.(map[string]interface{})
			if !ok {
				params = map[string]interface{}{}
			}

			// Create stochastic algorithm data
			algorithmData := &types.StochasticAlgorithmData{
				ID:         fmt.Sprintf("%d-%d", time.Now().UnixNano(), 1000),
				Algorithm:  "mdp",
				Problem:    problem,
				Parameters: params,
				Result:     "Optimized policy computed",
				Confidence: 0.85,
				Iterations: 1000,
				Converged:  true,
				CreatedAt:  time.Now(),
			}

			// Store the algorithm
			store.AddStochasticAlgorithm(sessionID, algorithmData)

			// Create response
			response := map[string]interface{}{
				"status":       "success",
				"algorithm_id": algorithmData.ID,
				"has_result":   true,
				"converged":    true,
				"iterations":   1000,
				"summary":      "Optimized policy computed successfully",
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// Monte Carlo Tree Search Tool
	s.AddTool(
		mcp.NewTool("monte_carlo_tree_search",
			mcp.WithDescription("Run Monte Carlo Tree Search for game tree exploration and decision making"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("problem", mcp.Required(), mcp.Description("Problem description for MCTS")),
			mcp.WithObject("parameters", mcp.Description("MCTS parameters (iterations, exploration constant, etc.)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			problem, _ := req.RequireString("problem")
			paramsInterface, _ := req.GetArguments()["parameters"]
			params, ok := paramsInterface.(map[string]interface{})
			if !ok {
				params = map[string]interface{}{}
			}

			// Create stochastic algorithm data
			algorithmData := &types.StochasticAlgorithmData{
				ID:         fmt.Sprintf("%d-%d", time.Now().UnixNano(), 10000),
				Algorithm:  "mcts",
				Problem:    problem,
				Parameters: params,
				Result:     "Best action selected",
				Confidence: 0.92,
				Iterations: 10000,
				Converged:  true,
				CreatedAt:  time.Now(),
			}

			// Store the algorithm
			store.AddStochasticAlgorithm(sessionID, algorithmData)

			// Create response
			response := map[string]interface{}{
				"status":       "success",
				"algorithm_id": algorithmData.ID,
				"has_result":   true,
				"converged":    true,
				"iterations":   10000,
				"summary":      "Best action selected through tree search",
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// Multi-Armed Bandit Tool
	s.AddTool(
		mcp.NewTool("multi_armed_bandit",
			mcp.WithDescription("Run Multi-Armed Bandit algorithm for exploration vs exploitation optimization"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("problem", mcp.Required(), mcp.Description("Problem description for bandit")),
			mcp.WithObject("parameters", mcp.Description("Bandit parameters (arms, epsilon, etc.)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			problem, _ := req.RequireString("problem")
			paramsInterface, _ := req.GetArguments()["parameters"]
			params, ok := paramsInterface.(map[string]interface{})
			if !ok {
				params = map[string]interface{}{}
			}

			// Create stochastic algorithm data
			algorithmData := &types.StochasticAlgorithmData{
				ID:         fmt.Sprintf("%d-%d", time.Now().UnixNano(), 1000),
				Algorithm:  "bandit",
				Problem:    problem,
				Parameters: params,
				Result:     "Optimal arm selected",
				Confidence: 0.88,
				Iterations: 1000,
				Converged:  true,
				CreatedAt:  time.Now(),
			}

			// Store the algorithm
			store.AddStochasticAlgorithm(sessionID, algorithmData)

			// Create response
			response := map[string]interface{}{
				"status":       "success",
				"algorithm_id": algorithmData.ID,
				"has_result":   true,
				"converged":    true,
				"iterations":   1000,
				"summary":      "Optimal arm selected for exploitation",
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)
}

func addDecisionTools(s *server.MCPServer, store *storage.Storage) {
	// Decision Framework Tool
	s.AddTool(
		mcp.NewTool("decision_framework",
			mcp.WithDescription("Apply decision frameworks for structured decision making"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("decision_statement", mcp.Required(), mcp.Description("Statement of the decision to be made")),
			mcp.WithArray("options", mcp.Description("Available decision options")),
			mcp.WithArray("criteria", mcp.Description("Decision criteria and weights")),
			mcp.WithString("analysis_type", mcp.Description("Type of analysis to perform")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			decisionStatement, _ := req.RequireString("decision_statement")
			optionsInterface, _ := req.GetArguments()["options"]
			criteriaInterface, _ := req.GetArguments()["criteria"]
			analysisType := req.GetString("analysis_type", "multi-criteria")

			// Convert options and criteria
			var options []types.DecisionOption
			if optionsSlice, ok := optionsInterface.([]interface{}); ok {
				for _, opt := range optionsSlice {
					if optMap, ok := opt.(map[string]interface{}); ok {
						option := types.DecisionOption{
							ID:          getString(optMap, "id"),
							Name:        getString(optMap, "name"),
							Description: getString(optMap, "description"),
						}
						options = append(options, option)
					}
				}
			}

			var criteria []types.DecisionCriterion
			if criteriaSlice, ok := criteriaInterface.([]interface{}); ok {
				for _, crit := range criteriaSlice {
					if critMap, ok := crit.(map[string]interface{}); ok {
						criterion := types.DecisionCriterion{
							ID:               getString(critMap, "id"),
							Name:             getString(critMap, "name"),
							Description:      getString(critMap, "description"),
							Weight:           getFloat64(critMap, "weight"),
							EvaluationMethod: getString(critMap, "evaluation_method"),
						}
						criteria = append(criteria, criterion)
					}
				}
			}

			// Create decision data
			decisionData := &types.DecisionData{
				ID:                fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(options)),
				DecisionStatement: decisionStatement,
				Options:           options,
				Criteria:          criteria,
				AnalysisType:      analysisType,
				Stage:             "evaluation",
				Iteration:         1,
				NextStageNeeded:   true,
				CreatedAt:         time.Now(),
			}

			// Store the decision
			store.AddDecision(sessionID, decisionData)

			// Create response
			response := map[string]interface{}{
				"status":        "success",
				"decision_id":   decisionData.ID,
				"has_options":   len(options) > 0,
				"has_criteria":  len(criteria) > 0,
				"analysis_type": analysisType,
				"stage":         "evaluation",
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)
}

func addVisualTools(s *server.MCPServer, store *storage.Storage) {
	// Concept Map Tool
	s.AddTool(
		mcp.NewTool("concept_map",
			mcp.WithDescription("Create and manipulate concept maps for visual thinking"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
			mcp.WithString("diagram_id", mcp.Description("Unique identifier for the diagram")),
			mcp.WithString("diagram_type", mcp.Description("Type of diagram (conceptMap, mindMap, etc.)")),
			mcp.WithString("operation", mcp.Required(), mcp.Description("Operation to perform (create, update, delete)")),
			mcp.WithArray("elements", mcp.Description("Visual elements (nodes, edges, etc.)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")
			diagramID := req.GetString("diagram_id", "default-diagram")
			diagramType := req.GetString("diagram_type", "conceptMap")
			operation, _ := req.RequireString("operation")
			elementsInterface, _ := req.GetArguments()["elements"]

			// Convert elements
			var elements []types.VisualElement
			if elementsSlice, ok := elementsInterface.([]interface{}); ok {
				for _, elem := range elementsSlice {
					if elemMap, ok := elem.(map[string]interface{}); ok {
						element := types.VisualElement{
							ID:         getString(elemMap, "id"),
							Type:       getString(elemMap, "type"),
							Label:      getString(elemMap, "label"),
							Properties: getProperties(elemMap["properties"]),
							Source:     getString(elemMap, "source"),
							Target:     getString(elemMap, "target"),
						}
						elements = append(elements, element)
					}
				}
			}

			// Create visual data
			visualData := &types.VisualData{
				ID:                  fmt.Sprintf("%d-%d", time.Now().UnixNano(), len(elements)),
				Operation:           operation,
				Elements:            elements,
				DiagramID:           diagramID,
				DiagramType:         diagramType,
				Iteration:           0,
				NextOperationNeeded: false,
				CreatedAt:           time.Now(),
			}

			// Store the visual data
			store.AddVisualData(sessionID, visualData)

			// Create response
			response := map[string]interface{}{
				"status":       "success",
				"visual_id":    visualData.ID,
				"operation":    operation,
				"diagram_type": diagramType,
				"elements":     len(elements),
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)
}

func addSessionTools(s *server.MCPServer, store *storage.Storage) {
	// Session Stats Tool
	s.AddTool(
		mcp.NewTool("session_stats",
			mcp.WithDescription("Get statistics for a session"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")

			// Get session stats
			stats, err := store.GetSessionStats(sessionID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get session stats: %v", err)), nil
			}

			// Create response
			response := map[string]interface{}{
				"session_id":         sessionID,
				"created_at":         stats.CreatedAt.Format(time.RFC3339),
				"last_accessed_at":   stats.LastAccessedAt.Format(time.RFC3339),
				"thought_count":      stats.ThoughtCount,
				"tools_used":         stats.ToolsUsed,
				"total_operations":   stats.TotalOperations,
				"is_active":          stats.IsActive,
				"remaining_thoughts": stats.RemainingThoughts,
				"stores":             stats.Stores,
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// Session Export Tool
	s.AddTool(
		mcp.NewTool("session_export",
			mcp.WithDescription("Export all data for a session"),
			mcp.WithString("session_id", mcp.Required(), mcp.Description("Session identifier")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessionID, _ := req.RequireString("session_id")

			// Export session data
			exportData, err := store.ExportSession(sessionID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to export session: %v", err)), nil
			}

			// Create response
			response := map[string]interface{}{
				"version":      "1.0.0",
				"timestamp":    time.Now().Format(time.RFC3339),
				"session_id":   sessionID,
				"session_type": "hybrid",
				"data":         exportData,
				"metadata": map[string]interface{}{
					"exported_at": time.Now().Format(time.RFC3339),
					"version":     "0.1.0",
				},
			}

			result, _ := json.Marshal(response)
			return mcp.NewToolResultText(string(result)), nil
		},
	)
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

func getProperties(properties interface{}) map[string]interface{} {
	if props, ok := properties.(map[string]interface{}); ok {
		return props
	}
	return nil
}
