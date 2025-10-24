package models

import "time"

// CVE represents a single CVE entry from the NVD
type CVE struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	CVSSScore   float64   `json:"cvss_score"`
	CVSSVector  string    `json:"cvss_vector"`
	Published   time.Time `json:"published"`
	Modified    time.Time `json:"modified"`
	References  []string  `json:"references"`
	Products    []string  `json:"products"`
	Vendors     []string  `json:"vendors"`
}

// AttackTechnique represents a MITRE ATT&CK technique
type AttackTechnique struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tactics     []string  `json:"tactics"`
	Platforms   []string  `json:"platforms"`
	KillChain   string    `json:"kill_chain"`
	References  []string  `json:"references"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// OWASPProcedure represents an OWASP testing procedure
type OWASPProcedure struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tools       []string  `json:"tools"`
	Steps       []string  `json:"steps"`
	References  []string  `json:"references"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// IntelligenceQuery represents a query for intelligence data
type IntelligenceQuery struct {
	Query     string `json:"query"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// IntelligenceResponse represents the response from an intelligence query
type IntelligenceResponse struct {
	Status    string        `json:"status"`
	Results   []interface{} `json:"results"`
	Total     int           `json:"total"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
	Query     string        `json:"query"`
	Source    string        `json:"source"`
	Timestamp time.Time     `json:"timestamp"`
}
