package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rainmana/gothink/internal/models"
)

// SecurityRepository handles database operations for security intelligence data
type SecurityRepository struct {
	// In a real implementation, this would be a database connection
	// For now, we'll use in-memory storage
	cves       map[string]models.CVE
	techniques map[string]models.AttackTechnique
	procedures map[string]models.OWASPProcedure
}

// NewSecurityRepository creates a new security repository
func NewSecurityRepository() *SecurityRepository {
	return &SecurityRepository{
		cves:       make(map[string]models.CVE),
		techniques: make(map[string]models.AttackTechnique),
		procedures: make(map[string]models.OWASPProcedure),
	}
}

// CVE Operations

// StoreCVE stores a CVE in the repository
func (r *SecurityRepository) StoreCVE(ctx context.Context, cve models.CVE) error {
	r.cves[cve.ID] = cve
	return nil
}

// StoreCVEs stores multiple CVEs in the repository
func (r *SecurityRepository) StoreCVEs(ctx context.Context, cves []models.CVE) error {
	for _, cve := range cves {
		if err := r.StoreCVE(ctx, cve); err != nil {
			return fmt.Errorf("failed to store CVE %s: %w", cve.ID, err)
		}
	}
	return nil
}

// GetCVE retrieves a CVE by ID
func (r *SecurityRepository) GetCVE(ctx context.Context, id string) (*models.CVE, error) {
	cve, exists := r.cves[id]
	if !exists {
		return nil, fmt.Errorf("CVE %s not found", id)
	}
	return &cve, nil
}

// QueryCVEs searches for CVEs based on query parameters
func (r *SecurityRepository) QueryCVEs(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	var results []interface{}

	for _, cve := range r.cves {
		// Simple text search in description
		if query.Query == "" || contains(cve.Description, query.Query) || contains(cve.ID, query.Query) {
			results = append(results, cve)
		}
	}

	// Apply pagination
	total := len(results)
	start := query.Offset
	end := start + query.Limit
	if end > len(results) {
		end = len(results)
	}
	if start > len(results) {
		start = len(results)
	}

	paginatedResults := results[start:end]

	return &models.IntelligenceResponse{
		Results:   paginatedResults,
		Total:     total,
		Limit:     query.Limit,
		Offset:    query.Offset,
		Query:     query.Query,
		Source:    "NVD",
		Timestamp: time.Now(),
	}, nil
}

// Attack Technique Operations

// StoreTechnique stores an attack technique in the repository
func (r *SecurityRepository) StoreTechnique(ctx context.Context, technique models.AttackTechnique) error {
	r.techniques[technique.ID] = technique
	return nil
}

// StoreTechniques stores multiple attack techniques in the repository
func (r *SecurityRepository) StoreTechniques(ctx context.Context, techniques []models.AttackTechnique) error {
	for _, technique := range techniques {
		if err := r.StoreTechnique(ctx, technique); err != nil {
			return fmt.Errorf("failed to store technique %s: %w", technique.ID, err)
		}
	}
	return nil
}

// GetTechnique retrieves an attack technique by ID
func (r *SecurityRepository) GetTechnique(ctx context.Context, id string) (*models.AttackTechnique, error) {
	technique, exists := r.techniques[id]
	if !exists {
		return nil, fmt.Errorf("technique %s not found", id)
	}
	return &technique, nil
}

// QueryTechniques searches for attack techniques based on query parameters
func (r *SecurityRepository) QueryTechniques(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	var results []interface{}

	for _, technique := range r.techniques {
		// Simple text search in name, description, and tactics
		if query.Query == "" ||
			contains(technique.Name, query.Query) ||
			contains(technique.Description, query.Query) ||
			contains(technique.ID, query.Query) {
			results = append(results, technique)
		}
	}

	// Apply pagination
	total := len(results)
	start := query.Offset
	end := start + query.Limit
	if end > len(results) {
		end = len(results)
	}
	if start > len(results) {
		start = len(results)
	}

	paginatedResults := results[start:end]

	return &models.IntelligenceResponse{
		Results:   paginatedResults,
		Total:     total,
		Limit:     query.Limit,
		Offset:    query.Offset,
		Query:     query.Query,
		Source:    "MITRE ATT&CK",
		Timestamp: time.Now(),
	}, nil
}

// OWASP Procedure Operations

// StoreProcedure stores an OWASP procedure in the repository
func (r *SecurityRepository) StoreProcedure(ctx context.Context, procedure models.OWASPProcedure) error {
	r.procedures[procedure.ID] = procedure
	return nil
}

// StoreProcedures stores multiple OWASP procedures in the repository
func (r *SecurityRepository) StoreProcedures(ctx context.Context, procedures []models.OWASPProcedure) error {
	for _, procedure := range procedures {
		if err := r.StoreProcedure(ctx, procedure); err != nil {
			return fmt.Errorf("failed to store procedure %s: %w", procedure.ID, err)
		}
	}
	return nil
}

// GetProcedure retrieves an OWASP procedure by ID
func (r *SecurityRepository) GetProcedure(ctx context.Context, id string) (*models.OWASPProcedure, error) {
	procedure, exists := r.procedures[id]
	if !exists {
		return nil, fmt.Errorf("procedure %s not found", id)
	}
	return &procedure, nil
}

// QueryProcedures searches for OWASP procedures based on query parameters
func (r *SecurityRepository) QueryProcedures(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	var results []interface{}

	for _, procedure := range r.procedures {
		// Simple text search in title, description, and category
		if query.Query == "" ||
			contains(procedure.Title, query.Query) ||
			contains(procedure.Description, query.Query) ||
			contains(procedure.Category, query.Query) ||
			contains(procedure.ID, query.Query) {
			results = append(results, procedure)
		}
	}

	// Apply pagination
	total := len(results)
	start := query.Offset
	end := start + query.Limit
	if end > len(results) {
		end = len(results)
	}
	if start > len(results) {
		start = len(results)
	}

	paginatedResults := results[start:end]

	return &models.IntelligenceResponse{
		Results:   paginatedResults,
		Total:     total,
		Limit:     query.Limit,
		Offset:    query.Offset,
		Query:     query.Query,
		Source:    "OWASP",
		Timestamp: time.Now(),
	}, nil
}

// Utility Functions

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetStats returns statistics about the repository
func (r *SecurityRepository) GetStats(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"cves":       len(r.cves),
		"techniques": len(r.techniques),
		"procedures": len(r.procedures),
		"total":      len(r.cves) + len(r.techniques) + len(r.procedures),
	}
}
