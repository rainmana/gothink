package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainmana/gothink/internal/intelligence"
	"github.com/rainmana/gothink/internal/models"
)

// IntelligenceHandler handles intelligence-related MCP requests
type IntelligenceHandler struct {
	intelligenceService *intelligence.IntelligenceService
}

// NewIntelligenceHandler creates a new intelligence handler
func NewIntelligenceHandler(apiKey string) *IntelligenceHandler {
	return &IntelligenceHandler{
		intelligenceService: intelligence.NewIntelligenceService(apiKey),
	}
}

// AddIntelligenceTools adds intelligence tools to the MCP server
func (h *IntelligenceHandler) AddIntelligenceTools(s *server.MCPServer) {
	// Query NVD CVE data
	s.AddTool(
		mcp.NewTool("query_nvd",
			mcp.WithDescription("Query NVD CVE data for security vulnerabilities"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query for CVEs")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of results to return")),
			mcp.WithNumber("offset", mcp.Description("Number of results to skip")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, _ := req.RequireString("query")
			limit := req.GetInt("limit", 10)
			offset := req.GetInt("offset", 0)

			// Create intelligence query
			intelQuery := models.IntelligenceQuery{
				Query:     query,
				Limit:     limit,
				Offset:    offset,
				SortBy:    "published",
				SortOrder: "desc",
			}

			// Query NVD data
			response, err := h.intelligenceService.QueryNVDData(ctx, intelQuery)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to query NVD data: %v", err)), nil
			}

			// Create response
			result := map[string]interface{}{
				"status":    "success",
				"source":    "NVD",
				"query":     query,
				"total":     response.Total,
				"limit":     response.Limit,
				"offset":    response.Offset,
				"results":   response.Results,
				"timestamp": response.Timestamp.Format(time.RFC3339),
			}

			resultJSON, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	)

	// Query MITRE ATT&CK data
	s.AddTool(
		mcp.NewTool("query_attack",
			mcp.WithDescription("Query MITRE ATT&CK techniques and tactics"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query for ATT&CK techniques")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of results to return")),
			mcp.WithNumber("offset", mcp.Description("Number of results to skip")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, _ := req.RequireString("query")
			limit := req.GetInt("limit", 10)
			offset := req.GetInt("offset", 0)

			// Create intelligence query
			intelQuery := models.IntelligenceQuery{
				Query:     query,
				Limit:     limit,
				Offset:    offset,
				SortBy:    "name",
				SortOrder: "asc",
			}

			// Query MITRE data
			response, err := h.intelligenceService.QueryMITREData(ctx, intelQuery)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to query MITRE data: %v", err)), nil
			}

			// Create response
			result := map[string]interface{}{
				"status":    "success",
				"source":    "MITRE ATT&CK",
				"query":     query,
				"total":     response.Total,
				"limit":     response.Limit,
				"offset":    response.Offset,
				"results":   response.Results,
				"timestamp": response.Timestamp.Format(time.RFC3339),
			}

			resultJSON, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	)

	// Query OWASP data
	s.AddTool(
		mcp.NewTool("query_owasp",
			mcp.WithDescription("Query OWASP testing procedures and guidelines"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query for OWASP procedures")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of results to return")),
			mcp.WithNumber("offset", mcp.Description("Number of results to skip")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, _ := req.RequireString("query")
			limit := req.GetInt("limit", 10)
			offset := req.GetInt("offset", 0)

			// Create intelligence query
			intelQuery := models.IntelligenceQuery{
				Query:     query,
				Limit:     limit,
				Offset:    offset,
				SortBy:    "title",
				SortOrder: "asc",
			}

			// Query OWASP data
			response, err := h.intelligenceService.QueryOWASPData(ctx, intelQuery)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to query OWASP data: %v", err)), nil
			}

			// Create response
			result := map[string]interface{}{
				"status":    "success",
				"source":    "OWASP",
				"query":     query,
				"total":     response.Total,
				"limit":     response.Limit,
				"offset":    response.Offset,
				"results":   response.Results,
				"timestamp": response.Timestamp.Format(time.RFC3339),
			}

			resultJSON, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	)

	// Refresh intelligence data
	s.AddTool(
		mcp.NewTool("refresh_intelligence",
			mcp.WithDescription("Refresh all intelligence data from external sources"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Refresh intelligence data
			if err := h.intelligenceService.RefreshIntelligenceData(ctx); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to refresh intelligence data: %v", err)), nil
			}

			// Get updated stats
			stats := h.intelligenceService.GetIntelligenceStats(ctx)

			// Create response
			result := map[string]interface{}{
				"status":    "success",
				"message":   "Intelligence data refreshed successfully",
				"stats":     stats,
				"timestamp": time.Now().Format(time.RFC3339),
			}

			resultJSON, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	)

	// Get intelligence stats
	s.AddTool(
		mcp.NewTool("intelligence_stats",
			mcp.WithDescription("Get statistics about available intelligence data"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get intelligence stats
			stats := h.intelligenceService.GetIntelligenceStats(ctx)

			// Create response
			result := map[string]interface{}{
				"status":    "success",
				"stats":     stats,
				"timestamp": time.Now().Format(time.RFC3339),
			}

			resultJSON, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	)
}
