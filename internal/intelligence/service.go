package intelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/rainmana/gothink/internal/models"
	"github.com/rainmana/gothink/internal/repository"
)

// IntelligenceService orchestrates intelligence data downloads and storage
type IntelligenceService struct {
	nvdDownloader   *NVDDownloader
	mitreDownloader *MITREDownloader
	owaspDownloader *OWASPDownloader
	securityRepo    *repository.SecurityRepository
}

// NewIntelligenceService creates a new intelligence service
func NewIntelligenceService(apiKey string) *IntelligenceService {
	return &IntelligenceService{
		nvdDownloader:   NewNVDDownloader(apiKey),
		mitreDownloader: NewMITREDownloader(),
		owaspDownloader: NewOWASPDownloader(),
		securityRepo:    repository.NewSecurityRepository(),
	}
}

// DownloadAndStoreAllIntelligence downloads and stores all intelligence data
func (s *IntelligenceService) DownloadAndStoreAllIntelligence(ctx context.Context) error {
	// Download NVD data
	if err := s.DownloadAndStoreNVDData(ctx); err != nil {
		return fmt.Errorf("failed to download NVD data: %w", err)
	}

	// Download MITRE ATT&CK data
	if err := s.DownloadAndStoreMITREData(ctx); err != nil {
		return fmt.Errorf("failed to download MITRE data: %w", err)
	}

	// Download OWASP data
	if err := s.DownloadAndStoreOWASPData(ctx); err != nil {
		return fmt.Errorf("failed to download OWASP data: %w", err)
	}

	return nil
}

// DownloadAndStoreNVDData downloads and stores NVD CVE data
func (s *IntelligenceService) DownloadAndStoreNVDData(ctx context.Context) error {
	// Download CVEs from NVD with retry logic
	var cves []models.CVE
	err := Retry(ctx, func() error {
		var err error
		cves, err = s.nvdDownloader.DownloadAllCVEs(ctx)
		if err != nil && IsRetryableError(err) {
			return err
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to download CVEs: %w", err)
	}

	// Store CVEs in repository
	if err := s.securityRepo.StoreCVEs(ctx, cves); err != nil {
		return fmt.Errorf("failed to store CVEs: %w", err)
	}

	return nil
}

// DownloadAndStoreMITREData downloads and stores MITRE ATT&CK data
func (s *IntelligenceService) DownloadAndStoreMITREData(ctx context.Context) error {
	// Download techniques from MITRE with retry logic
	var techniques []models.AttackTechnique
	err := Retry(ctx, func() error {
		var err error
		techniques, err = s.mitreDownloader.DownloadTechniques(ctx)
		if err != nil && IsRetryableError(err) {
			return err
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to download techniques: %w", err)
	}

	// Store techniques in repository
	if err := s.securityRepo.StoreTechniques(ctx, techniques); err != nil {
		return fmt.Errorf("failed to store techniques: %w", err)
	}

	return nil
}

// DownloadAndStoreOWASPData downloads and stores OWASP data
func (s *IntelligenceService) DownloadAndStoreOWASPData(ctx context.Context) error {
	// Download procedures from OWASP with retry logic
	var procedures []models.OWASPProcedure
	err := Retry(ctx, func() error {
		var err error
		procedures, err = s.owaspDownloader.DownloadProcedures(ctx)
		if err != nil && IsRetryableError(err) {
			return err
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to download procedures: %w", err)
	}

	// Store procedures in repository
	if err := s.securityRepo.StoreProcedures(ctx, procedures); err != nil {
		return fmt.Errorf("failed to store procedures: %w", err)
	}

	return nil
}

// QueryNVDData queries NVD CVE data
func (s *IntelligenceService) QueryNVDData(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	return s.securityRepo.QueryCVEs(ctx, query)
}

// QueryMITREData queries MITRE ATT&CK data
func (s *IntelligenceService) QueryMITREData(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	return s.securityRepo.QueryTechniques(ctx, query)
}

// QueryOWASPData queries OWASP data
func (s *IntelligenceService) QueryOWASPData(ctx context.Context, query models.IntelligenceQuery) (*models.IntelligenceResponse, error) {
	return s.securityRepo.QueryProcedures(ctx, query)
}

// GetIntelligenceStats returns statistics about the intelligence data
func (s *IntelligenceService) GetIntelligenceStats(ctx context.Context) map[string]interface{} {
	return s.securityRepo.GetStats(ctx)
}

// RefreshIntelligenceData refreshes all intelligence data
func (s *IntelligenceService) RefreshIntelligenceData(ctx context.Context) error {
	// Set a timeout for the refresh operation
	refreshCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Download and store all intelligence data
	if err := s.DownloadAndStoreAllIntelligence(refreshCtx); err != nil {
		return fmt.Errorf("failed to refresh intelligence data: %w", err)
	}

	return nil
}
