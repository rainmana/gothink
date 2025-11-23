package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rainmana/gothink/internal/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// MentalModelConfig represents the YAML configuration for custom mental models
type MentalModelConfig struct {
	Models map[string]MentalModel `yaml:"models"`
}

// MentalModel represents a mental model with priority support
type MentalModel struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Steps       []string `yaml:"steps" json:"steps"`
	Category    string   `yaml:"category" json:"category"`
	Priority    int      `yaml:"priority,omitempty" json:"priority,omitempty"`
}

// MentalModelWithKey represents a mental model with its key for sorting
type MentalModelWithKey struct {
	Key   string
	Model MentalModel
}

// Loader handles loading and managing mental models
type Loader struct {
	logger *logrus.Logger
}

// NewLoader creates a new mental models loader
func NewLoader(logger *logrus.Logger) *Loader {
	return &Loader{
		logger: logger,
	}
}

// LoadMentalModels loads mental models from core types and optional custom YAML file
func (l *Loader) LoadMentalModels(configPath string) (map[string]MentalModel, error) {
	// Start with core models (always available as fallback)
	models := make(map[string]MentalModel)

	// Convert core models to our format
	for key, coreModel := range types.MentalModels {
		models[key] = MentalModel{
			Name:        coreModel.Name,
			Description: coreModel.Description,
			Steps:       coreModel.Steps,
			Category:    coreModel.Category,
			Priority:    0, // Core models have default priority
		}
	}

	l.logger.Infof("Loaded %d core mental models", len(models))

	// Load custom models if file exists
	if configPath != "" {
		customModels, err := l.loadCustomModels(configPath)
		if err != nil {
			l.logger.Warnf("Failed to load custom mental models from %s: %v", configPath, err)
			// Continue with core models only
		} else {
			// Merge custom models (they can override core models)
			for key, model := range customModels {
				models[key] = model
				l.logger.Infof("Loaded custom mental model: %s (priority: %d)", key, model.Priority)
			}
		}
	}

	return models, nil
}

// loadCustomModels loads mental models from a YAML file or directory
func (l *Loader) loadCustomModels(path string) (map[string]MentalModel, error) {
	// Check if path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("mental models path does not exist: %s", path)
	}

	models := make(map[string]MentalModel)

	if info.IsDir() {
		// Walk directory
		err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
				fileModels, err := l.loadModelsFromFile(path)
				if err != nil {
					l.logger.Warnf("Failed to load models from %s: %v", path, err)
					return nil // Continue loading other files
				}
				for k, v := range fileModels {
					models[k] = v
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		// Load single file
		models, err = l.loadModelsFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return models, nil
}

// loadModelsFromFile loads mental models from a single YAML file
func (l *Loader) loadModelsFromFile(filePath string) (map[string]MentalModel, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mental models file: %w", err)
	}

	// Parse YAML
	var config MentalModelConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse mental models YAML: %w", err)
	}

	// Validate models
	if err := l.validateModels(config.Models); err != nil {
		return nil, fmt.Errorf("invalid mental models configuration: %w", err)
	}

	return config.Models, nil
}

// validateModels validates the mental models configuration
func (l *Loader) validateModels(models map[string]MentalModel) error {
	for key, model := range models {
		// Check required fields
		if strings.TrimSpace(model.Name) == "" {
			return fmt.Errorf("model '%s' has empty name", key)
		}
		if strings.TrimSpace(model.Description) == "" {
			return fmt.Errorf("model '%s' has empty description", key)
		}
		if len(model.Steps) == 0 {
			return fmt.Errorf("model '%s' has no steps", key)
		}
		if strings.TrimSpace(model.Category) == "" {
			return fmt.Errorf("model '%s' has empty category", key)
		}

		// Validate steps
		for i, step := range model.Steps {
			if strings.TrimSpace(step) == "" {
				return fmt.Errorf("model '%s' has empty step at index %d", key, i)
			}
		}

		// Set default priority if not specified
		if model.Priority == 0 {
			models[key] = MentalModel{
				Name:        model.Name,
				Description: model.Description,
				Steps:       model.Steps,
				Category:    model.Category,
				Priority:    1, // Custom models get priority 1 by default
			}
		}
	}

	return nil
}

// GetModelsByPriority returns models sorted by priority (highest first)
func (l *Loader) GetModelsByPriority(models map[string]MentalModel) []MentalModelWithKey {
	var modelsWithKeys []MentalModelWithKey

	for key, model := range models {
		modelsWithKeys = append(modelsWithKeys, MentalModelWithKey{
			Key:   key,
			Model: model,
		})
	}

	// Sort by priority (highest first), then by name
	sort.Slice(modelsWithKeys, func(i, j int) bool {
		if modelsWithKeys[i].Model.Priority != modelsWithKeys[j].Model.Priority {
			return modelsWithKeys[i].Model.Priority > modelsWithKeys[j].Model.Priority
		}
		return modelsWithKeys[i].Model.Name < modelsWithKeys[j].Model.Name
	})

	return modelsWithKeys
}

// GetModelsByCategory returns models grouped by category
func (l *Loader) GetModelsByCategory(models map[string]MentalModel) map[string][]MentalModelWithKey {
	categories := make(map[string][]MentalModelWithKey)

	for key, model := range models {
		categories[model.Category] = append(categories[model.Category], MentalModelWithKey{
			Key:   key,
			Model: model,
		})
	}

	// Sort each category by priority
	for category, models := range categories {
		sort.Slice(models, func(i, j int) bool {
			if models[i].Model.Priority != models[j].Model.Priority {
				return models[i].Model.Priority > models[j].Model.Priority
			}
			return models[i].Model.Name < models[j].Model.Name
		})
		categories[category] = models
	}

	return categories
}

// GetAvailableModels returns a list of available model keys and names
func (l *Loader) GetAvailableModels(models map[string]MentalModel) []string {
	var available []string
	for key, model := range models {
		available = append(available, fmt.Sprintf("%s: %s", key, model.Name))
	}
	return available
}
