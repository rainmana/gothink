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

// MITREDownloader handles downloading ATT&CK data from MITRE
type MITREDownloader struct {
	client  *http.Client
	baseURL string
}

// NewMITREDownloader creates a new MITRE downloader
func NewMITREDownloader() *MITREDownloader {
	return &MITREDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json",
	}
}

// MITREResponse represents the response structure from MITRE ATT&CK
type MITREResponse struct {
	Type        string `json:"type"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		Type            string   `json:"type"`
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		XMitrePlatforms []string `json:"x_mitre_platforms"`
		KillChainPhases []struct {
			KillChainName string `json:"kill_chain_name"`
			PhaseName     string `json:"phase_name"`
		} `json:"kill_chain_phases"`
		ExternalReferences []struct {
			SourceName string `json:"source_name"`
			URL        string `json:"url"`
			ExternalID string `json:"external_id"`
		} `json:"external_references"`
		XMitreDataSources         []string `json:"x_mitre_data_sources"`
		XMitreDefenseBypassed     []string `json:"x_mitre_defense_bypassed"`
		XMitrePermissionsRequired []string `json:"x_mitre_permissions_required"`
		XMitreSystemRequirements  []string `json:"x_mitre_system_requirements"`
		XMitreNetworkRequirements bool     `json:"x_mitre_network_requirements"`
		XMitreRemoteSupport       bool     `json:"x_mitre_remote_support"`
		XMitreContributors        []string `json:"x_mitre_contributors"`
		XMitreVersion             string   `json:"x_mitre_version"`
		Created                   string   `json:"created"`
		Modified                  string   `json:"modified"`
		Revoked                   bool     `json:"revoked"`
		XMitreDeprecated          bool     `json:"x_mitre_deprecated"`
	} `json:"objects"`
}

// DownloadTechniques downloads ATT&CK techniques from MITRE
func (m *MITREDownloader) DownloadTechniques(ctx context.Context) ([]models.AttackTechnique, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "GoThink-Security-Intelligence/1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MITRE API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var mitreResp MITREResponse
	if err := json.Unmarshal(body, &mitreResp); err != nil {
		return nil, fmt.Errorf("failed to parse MITRE response: %w", err)
	}

	// Convert MITRE response to our AttackTechnique models
	var techniques []models.AttackTechnique
	fmt.Printf("Processing %d objects from MITRE...\n", len(mitreResp.Objects))
	attackPatternCount := 0
	for _, obj := range mitreResp.Objects {
		// Only process attack-pattern objects (techniques)
		if obj.Type == "attack-pattern" {
			attackPatternCount++
			technique := models.AttackTechnique{
				ID:          obj.ID,
				Name:        obj.Name,
				Description: obj.Description,
				Platforms:   obj.XMitrePlatforms,
				Created:     parseMITRETime(obj.Created),
				Modified:    parseMITRETime(obj.Modified),
			}

			// Extract tactics from kill chain phases
			for _, phase := range obj.KillChainPhases {
				if phase.KillChainName == "mitre-attack" {
					technique.Tactics = append(technique.Tactics, phase.PhaseName)
				}
			}

			// Extract references
			for _, ref := range obj.ExternalReferences {
				technique.References = append(technique.References, ref.URL)
			}

			// Set kill chain
			technique.KillChain = "mitre-attack"

			techniques = append(techniques, technique)
		}
	}
	
	fmt.Printf("Found %d attack-pattern objects, created %d techniques\n", attackPatternCount, len(techniques))
	return techniques, nil
}

// DownloadTactics downloads ATT&CK tactics from MITRE
func (m *MITREDownloader) DownloadTactics(ctx context.Context) ([]models.AttackTechnique, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "GoThink-Security-Intelligence/1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MITRE API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var mitreResp MITREResponse
	if err := json.Unmarshal(body, &mitreResp); err != nil {
		return nil, fmt.Errorf("failed to parse MITRE response: %w", err)
	}

	// Convert MITRE response to our AttackTechnique models
	var tactics []models.AttackTechnique
	for _, obj := range mitreResp.Objects {
		// Only process x-mitre-tactic objects (tactics)
		if obj.Type == "x-mitre-tactic" {
			tactic := models.AttackTechnique{
				ID:          obj.ID,
				Name:        obj.Name,
				Description: obj.Description,
				Platforms:   obj.XMitrePlatforms,
				Created:     parseMITRETime(obj.Created),
				Modified:    parseMITRETime(obj.Modified),
			}

			// Extract references
			for _, ref := range obj.ExternalReferences {
				tactic.References = append(tactic.References, ref.URL)
			}

			// Set kill chain
			tactic.KillChain = "mitre-attack"

			tactics = append(tactics, tactic)
		}
	}

	return tactics, nil
}

// parseMITRETime parses a time string from MITRE ATT&CK
func parseMITRETime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", timeStr)
	if err != nil {
		// Try alternative format
		t, err = time.Parse("2006-01-02T15:04:05Z", timeStr)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}
