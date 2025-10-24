package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rainmana/gothink/internal/models"
)

// NVDDownloader handles downloading CVE data from the National Vulnerability Database
type NVDDownloader struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewNVDDownloader creates a new NVD downloader
func NewNVDDownloader(apiKey string) *NVDDownloader {
	return &NVDDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://services.nvd.nist.gov/rest/json/cves/2.0",
		apiKey:  apiKey,
	}
}

// NVDResponse represents the response structure from NVD API
type NVDResponse struct {
	ResultsPerPage  int `json:"resultsPerPage"`
	StartIndex      int `json:"startIndex"`
	TotalResults    int `json:"totalResults"`
	Vulnerabilities []struct {
		CVE struct {
			ID               string `json:"id"`
			SourceIdentifier string `json:"sourceIdentifier"`
			Published        string `json:"published"`
			LastModified     string `json:"lastModified"`
			VulnStatus       string `json:"vulnStatus"`
			Descriptions     []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"descriptions"`
			References []struct {
				URL    string   `json:"url"`
				Source string   `json:"source"`
				Tags   []string `json:"tags"`
			} `json:"references"`
			Metrics struct {
				CvssMetricV31 []struct {
					Source   string `json:"source"`
					Type     string `json:"type"`
					CvssData struct {
						Version               string  `json:"version"`
						VectorString          string  `json:"vectorString"`
						AttackVector          string  `json:"attackVector"`
						AttackComplexity      string  `json:"attackComplexity"`
						PrivilegesRequired    string  `json:"privilegesRequired"`
						UserInteraction       string  `json:"userInteraction"`
						Scope                 string  `json:"scope"`
						ConfidentialityImpact string  `json:"confidentialityImpact"`
						IntegrityImpact       string  `json:"integrityImpact"`
						AvailabilityImpact    string  `json:"availabilityImpact"`
						BaseScore             float64 `json:"baseScore"`
						BaseSeverity          string  `json:"baseSeverity"`
					} `json:"cvssData"`
				} `json:"cvssMetricV31"`
			} `json:"metrics"`
			Weaknesses []struct {
				Source      string `json:"source"`
				Type        string `json:"type"`
				Description []struct {
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"description"`
			} `json:"weaknesses"`
			Configurations []struct {
				Nodes []struct {
					Operator string `json:"operator"`
					Negate   bool   `json:"negate"`
					CpeMatch []struct {
						Vulnerable            bool   `json:"vulnerable"`
						Cpe23Uri              string `json:"cpe23Uri"`
						VersionStartIncluding string `json:"versionStartIncluding"`
						VersionEndIncluding   string `json:"versionEndIncluding"`
					} `json:"cpeMatch"`
				} `json:"nodes"`
			} `json:"configurations"`
		} `json:"cve"`
	} `json:"vulnerabilities"`
}

// DownloadCVEs downloads CVE data from NVD
func (n *NVDDownloader) DownloadCVEs(ctx context.Context, startIndex int, resultsPerPage int) ([]models.CVE, error) {
	url := fmt.Sprintf("%s?startIndex=%d&resultsPerPage=%d", n.baseURL, startIndex, resultsPerPage)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if available
	if n.apiKey != "" {
		req.Header.Set("apiKey", n.apiKey)
	}

	req.Header.Set("User-Agent", "GoThink-Security-Intelligence/1.0")

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("NVD API rate limit exceeded (429) - too many requests")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NVD API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var nvdResp NVDResponse
	if err := json.Unmarshal(body, &nvdResp); err != nil {
		return nil, fmt.Errorf("failed to parse NVD response: %w", err)
	}

	// Convert NVD response to our CVE models
	var cves []models.CVE
	for _, vuln := range nvdResp.Vulnerabilities {
		cve := models.CVE{
			ID:        vuln.CVE.ID,
			Published: parseTime(vuln.CVE.Published),
			Modified:  parseTime(vuln.CVE.LastModified),
		}

		// Extract description
		for _, desc := range vuln.CVE.Descriptions {
			if desc.Lang == "en" {
				cve.Description = desc.Value
				break
			}
		}

		// Extract CVSS score and severity
		if len(vuln.CVE.Metrics.CvssMetricV31) > 0 {
			cvss := vuln.CVE.Metrics.CvssMetricV31[0]
			cve.CVSSScore = cvss.CvssData.BaseScore
			cve.CVSSVector = cvss.CvssData.VectorString
			cve.Severity = cvss.CvssData.BaseSeverity
		}

		// Extract references
		for _, ref := range vuln.CVE.References {
			cve.References = append(cve.References, ref.URL)
		}

		// Extract products and vendors from configurations
		products := make(map[string]bool)
		vendors := make(map[string]bool)
		for _, config := range vuln.CVE.Configurations {
			for _, node := range config.Nodes {
				for _, cpe := range node.CpeMatch {
					if cpe.Vulnerable {
						// Parse CPE URI to extract vendor and product
						// CPE format: cpe:2.3:a:vendor:product:version:update:edition:language:sw_edition:target_sw:target_hw:other
						parts := splitCPE(cpe.Cpe23Uri)
						if len(parts) >= 4 {
							vendors[parts[3]] = true
							if len(parts) >= 5 {
								products[parts[4]] = true
							}
						}
					}
				}
			}
		}

		// Convert maps to slices
		for vendor := range vendors {
			cve.Vendors = append(cve.Vendors, vendor)
		}
		for product := range products {
			cve.Products = append(cve.Products, product)
		}

		cves = append(cves, cve)
	}

	return cves, nil
}

// DownloadAllCVEs downloads all CVE data from NVD (with pagination)
func (n *NVDDownloader) DownloadAllCVEs(ctx context.Context) ([]models.CVE, error) {
	var allCVEs []models.CVE
	startIndex := 0
	resultsPerPage := 2000 // NVD API max

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		cves, err := n.DownloadCVEs(ctx, startIndex, resultsPerPage)
		if err != nil {
			return nil, fmt.Errorf("failed to download CVEs at index %d: %w", startIndex, err)
		}

		if len(cves) == 0 {
			break
		}

		allCVEs = append(allCVEs, cves...)
		startIndex += len(cves)

		// Rate limiting - NVD API allows 5 requests per 30 seconds without API key
		// Use 7 seconds to be safe
		time.Sleep(7 * time.Second)
	}

	return allCVEs, nil
}

// parseTime parses a time string from NVD API
func parseTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000", timeStr)
	if err != nil {
		// Try alternative format
		t, err = time.Parse("2006-01-02T15:04:05", timeStr)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}

// splitCPE splits a CPE URI into its components
func splitCPE(cpeURI string) []string {
	// Simple CPE parsing - in production, use a proper CPE library
	parts := make([]string, 0)
	current := ""

	for _, char := range cpeURI {
		if char == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
