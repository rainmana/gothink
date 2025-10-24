package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rainmana/gothink/internal/handlers"
	"github.com/rainmana/gothink/internal/intelligence"
	"github.com/rainmana/gothink/internal/models"
)

func main() {
	// Create intelligence service
	intelligenceService := intelligence.NewIntelligenceService("")

	// Download initial data (skip NVD for now due to rate limits)
	fmt.Println("Downloading intelligence data...")
	ctx := context.Background()

	// Download MITRE and OWASP data only
	fmt.Println("Downloading MITRE data...")
	if err := intelligenceService.DownloadAndStoreMITREData(ctx); err != nil {
		log.Printf("Warning: Failed to download MITRE data: %v", err)
	} else {
		fmt.Println("MITRE data downloaded successfully!")
	}

	fmt.Println("Downloading OWASP data...")
	if err := intelligenceService.DownloadAndStoreOWASPData(ctx); err != nil {
		log.Printf("Warning: Failed to download OWASP data: %v", err)
	} else {
		fmt.Println("OWASP data downloaded successfully!")
	}

	// Create intelligence handler with the same service instance
	intelligenceHandler := handlers.NewIntelligenceHandler("")
	// Share the same intelligence service instance
	intelligenceHandler.SetIntelligenceService(intelligenceService)

	// Create HTTP server for testing
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Handle MCP request
		response := handleMCPRequest(req, intelligenceHandler)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// Start server
	port := "8090"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Starting test MCP server on port %s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("MCP endpoint: http://localhost:%s/mcp\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleMCPRequest(req map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Extract method and params
	method, ok := req["method"].(string)
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": map[string]interface{}{
				"code":    -32600,
				"message": "Invalid Request",
			},
		}
	}

	// Handle different MCP methods
	switch method {
	case "tools/call":
		return handleToolCall(req, handler)
	case "tools/list":
		return handleToolList()
	default:
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

func handleToolCall(req map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	params, ok := req["params"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Tool name required",
			},
		}
	}

	// Handle intelligence tools
	switch toolName {
	case "query_nvd":
		return handleQueryNVD(params, handler)
	case "query_attack":
		return handleQueryAttack(params, handler)
	case "query_owasp":
		return handleQueryOWASP(params, handler)
	case "refresh_intelligence":
		return handleRefreshIntelligence(params, handler)
	case "intelligence_stats":
		return handleIntelligenceStats(params, handler)
	default:
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "Tool not found",
			},
		}
	}
}

func handleQueryNVD(params map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Extract query parameters
	query := "test"
	if q, ok := params["arguments"].(map[string]interface{})["query"].(string); ok {
		query = q
	}
	limit := 10
	if l, ok := params["arguments"].(map[string]interface{})["limit"].(float64); ok {
		limit = int(l)
	}
	offset := 0
	if o, ok := params["arguments"].(map[string]interface{})["offset"].(float64); ok {
		offset = int(o)
	}

	// Create intelligence query
	intelQuery := models.IntelligenceQuery{
		Query:     query,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "published",
		SortOrder: "desc",
	}

	// Query NVD data using the intelligence service
	ctx := context.Background()
	response, err := handler.QueryNVDData(ctx, intelQuery)
	if err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf(`{"status":"error","message":"%s","timestamp":"%s"}`, err.Error(), time.Now().Format(time.RFC3339)),
					},
				},
			},
		}
	}

	// Add status field to response
	response.Status = "success"
	response.Timestamp = time.Now()

	// Convert response to JSON
	responseJSON, _ := json.Marshal(response)

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": string(responseJSON),
				},
			},
		},
	}
}

func handleQueryAttack(params map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Extract query parameters
	query := "test"
	if q, ok := params["arguments"].(map[string]interface{})["query"].(string); ok {
		query = q
	}
	limit := 10
	if l, ok := params["arguments"].(map[string]interface{})["limit"].(float64); ok {
		limit = int(l)
	}
	offset := 0
	if o, ok := params["arguments"].(map[string]interface{})["offset"].(float64); ok {
		offset = int(o)
	}

	// Create intelligence query
	intelQuery := models.IntelligenceQuery{
		Query:     query,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "name",
		SortOrder: "asc",
	}

	// Query MITRE data using the intelligence service
	ctx := context.Background()
	response, err := handler.QueryMITREData(ctx, intelQuery)
	if err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf(`{"status":"error","message":"%s","timestamp":"%s"}`, err.Error(), time.Now().Format(time.RFC3339)),
					},
				},
			},
		}
	}

	// Add status field to response
	response.Status = "success"
	response.Timestamp = time.Now()

	// Convert response to JSON
	responseJSON, _ := json.Marshal(response)

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": string(responseJSON),
				},
			},
		},
	}
}

func handleQueryOWASP(params map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Extract query parameters
	query := "test"
	if q, ok := params["arguments"].(map[string]interface{})["query"].(string); ok {
		query = q
	}
	limit := 10
	if l, ok := params["arguments"].(map[string]interface{})["limit"].(float64); ok {
		limit = int(l)
	}
	offset := 0
	if o, ok := params["arguments"].(map[string]interface{})["offset"].(float64); ok {
		offset = int(o)
	}

	// Create intelligence query
	intelQuery := models.IntelligenceQuery{
		Query:     query,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "title",
		SortOrder: "asc",
	}

	// Query OWASP data using the intelligence service
	ctx := context.Background()
	response, err := handler.QueryOWASPData(ctx, intelQuery)
	if err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf(`{"status":"error","message":"%s","timestamp":"%s"}`, err.Error(), time.Now().Format(time.RFC3339)),
					},
				},
			},
		}
	}

	// Add status field to response
	response.Status = "success"
	response.Timestamp = time.Now()

	// Convert response to JSON
	responseJSON, _ := json.Marshal(response)

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": string(responseJSON),
				},
			},
		},
	}
}

func handleRefreshIntelligence(params map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Refresh intelligence data
	ctx := context.Background()
	if err := handler.RefreshIntelligenceData(ctx); err != nil {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf(`{"status":"error","message":"%s","timestamp":"%s"}`, err.Error(), time.Now().Format(time.RFC3339)),
					},
				},
			},
		}
	}

	// Get updated stats
	stats := handler.GetIntelligenceStats(ctx)
	statsJSON, _ := json.Marshal(stats)

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf(`{"status":"success","message":"Intelligence data refreshed successfully","stats":%s,"timestamp":"%s"}`, string(statsJSON), time.Now().Format(time.RFC3339)),
				},
			},
		},
	}
}

func handleIntelligenceStats(params map[string]interface{}, handler *handlers.IntelligenceHandler) map[string]interface{} {
	// Get intelligence stats
	ctx := context.Background()
	stats := handler.GetIntelligenceStats(ctx)
	statsJSON, _ := json.Marshal(stats)

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf(`{"status":"success","stats":%s,"timestamp":"%s"}`, string(statsJSON), time.Now().Format(time.RFC3339)),
				},
			},
		},
	}
}

func handleToolList() map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"tools": []map[string]interface{}{
				{
					"name":        "query_nvd",
					"description": "Query NVD CVE data for security vulnerabilities",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"query": map[string]interface{}{
								"type":        "string",
								"description": "Search query for CVEs",
							},
							"limit": map[string]interface{}{
								"type":        "number",
								"description": "Maximum number of results to return",
							},
							"offset": map[string]interface{}{
								"type":        "number",
								"description": "Number of results to skip",
							},
						},
						"required": []string{"query"},
					},
				},
				{
					"name":        "query_attack",
					"description": "Query MITRE ATT&CK techniques and tactics",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"query": map[string]interface{}{
								"type":        "string",
								"description": "Search query for ATT&CK techniques",
							},
							"limit": map[string]interface{}{
								"type":        "number",
								"description": "Maximum number of results to return",
							},
							"offset": map[string]interface{}{
								"type":        "number",
								"description": "Number of results to skip",
							},
						},
						"required": []string{"query"},
					},
				},
				{
					"name":        "query_owasp",
					"description": "Query OWASP testing procedures and guidelines",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"query": map[string]interface{}{
								"type":        "string",
								"description": "Search query for OWASP procedures",
							},
							"limit": map[string]interface{}{
								"type":        "number",
								"description": "Maximum number of results to return",
							},
							"offset": map[string]interface{}{
								"type":        "number",
								"description": "Number of results to skip",
							},
						},
						"required": []string{"query"},
					},
				},
				{
					"name":        "refresh_intelligence",
					"description": "Refresh all intelligence data from external sources",
					"inputSchema": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				{
					"name":        "intelligence_stats",
					"description": "Get statistics about available intelligence data",
					"inputSchema": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
			},
		},
	}
}
