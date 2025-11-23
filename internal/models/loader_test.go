package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	assert.NotNil(t, loader)
	assert.Equal(t, logger, loader.logger)
}

func TestLoadMentalModels_NoCustomFile(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Test loading with no custom file (should load core models only)
	models, err := loader.LoadMentalModels("")

	require.NoError(t, err)
	assert.NotEmpty(t, models)

	// Check that core models are loaded
	assert.Contains(t, models, "first_principles")
	assert.Contains(t, models, "opportunity_cost")
	assert.Contains(t, models, "bayesian_thinking")
	assert.Contains(t, models, "systems_thinking")

	// Check core model properties
	firstPrinciples := models["first_principles"]
	assert.Equal(t, "First Principles Thinking", firstPrinciples.Name)
	assert.Equal(t, "analytical", firstPrinciples.Category)
	assert.Equal(t, 0, firstPrinciples.Priority) // Core models have priority 0
	assert.NotEmpty(t, firstPrinciples.Steps)
}

func TestLoadMentalModels_WithCustomFile(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary YAML file
	yamlContent := `
models:
  custom_model_1:
    name: "My Custom Model"
    description: "A custom mental model for testing"
    steps:
      - "Step 1: Define the problem"
      - "Step 2: Gather information"
      - "Step 3: Analyze patterns"
      - "Step 4: Generate solutions"
    category: "custom"
    priority: 10

  custom_model_2:
    name: "Another Custom Model"
    description: "Another custom approach"
    steps:
      - "Step 1: Start here"
      - "Step 2: Continue here"
    category: "analytical"
    priority: 5
`

	// Create temporary file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "mental_models.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Load models
	models, err := loader.LoadMentalModels(configPath)

	require.NoError(t, err)
	assert.NotEmpty(t, models)

	// Check that core models are still there
	assert.Contains(t, models, "first_principles")
	assert.Contains(t, models, "opportunity_cost")

	// Check that custom models are loaded
	assert.Contains(t, models, "custom_model_1")
	assert.Contains(t, models, "custom_model_2")

	// Check custom model properties
	customModel1 := models["custom_model_1"]
	assert.Equal(t, "My Custom Model", customModel1.Name)
	assert.Equal(t, "A custom mental model for testing", customModel1.Description)
	assert.Equal(t, "custom", customModel1.Category)
	assert.Equal(t, 10, customModel1.Priority)
	assert.Len(t, customModel1.Steps, 4)
	assert.Equal(t, "Step 1: Define the problem", customModel1.Steps[0])

	customModel2 := models["custom_model_2"]
	assert.Equal(t, "Another Custom Model", customModel2.Name)
	assert.Equal(t, "analytical", customModel2.Category)
	assert.Equal(t, 5, customModel2.Priority)
}

func TestLoadMentalModels_InvalidFile(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Test with non-existent file
	models, err := loader.LoadMentalModels("/nonexistent/file.yaml")

	require.NoError(t, err)    // Should not error, just log warning
	assert.NotEmpty(t, models) // Should still have core models

	// Check that core models are loaded
	assert.Contains(t, models, "first_principles")
}

func TestLoadMentalModels_InvalidYAML(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary file with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	require.NoError(t, err)

	// Load models
	models, err := loader.LoadMentalModels(configPath)

	require.NoError(t, err)    // Should not error, just log warning
	assert.NotEmpty(t, models) // Should still have core models

	// Check that core models are loaded
	assert.Contains(t, models, "first_principles")
}

func TestValidateModels(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	tests := []struct {
		name    string
		models  map[string]MentalModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid models",
			models: map[string]MentalModel{
				"valid_model": {
					Name:        "Valid Model",
					Description: "A valid model",
					Steps:       []string{"Step 1", "Step 2"},
					Category:    "test",
					Priority:    5,
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			models: map[string]MentalModel{
				"invalid_model": {
					Name:        "",
					Description: "A model with empty name",
					Steps:       []string{"Step 1"},
					Category:    "test",
				},
			},
			wantErr: true,
			errMsg:  "empty name",
		},
		{
			name: "empty description",
			models: map[string]MentalModel{
				"invalid_model": {
					Name:        "Valid Name",
					Description: "",
					Steps:       []string{"Step 1"},
					Category:    "test",
				},
			},
			wantErr: true,
			errMsg:  "empty description",
		},
		{
			name: "no steps",
			models: map[string]MentalModel{
				"invalid_model": {
					Name:        "Valid Name",
					Description: "Valid description",
					Steps:       []string{},
					Category:    "test",
				},
			},
			wantErr: true,
			errMsg:  "no steps",
		},
		{
			name: "empty category",
			models: map[string]MentalModel{
				"invalid_model": {
					Name:        "Valid Name",
					Description: "Valid description",
					Steps:       []string{"Step 1"},
					Category:    "",
				},
			},
			wantErr: true,
			errMsg:  "empty category",
		},
		{
			name: "empty step",
			models: map[string]MentalModel{
				"invalid_model": {
					Name:        "Valid Name",
					Description: "Valid description",
					Steps:       []string{"Step 1", "", "Step 3"},
					Category:    "test",
				},
			},
			wantErr: true,
			errMsg:  "empty step at index 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.validateModels(tt.models)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetModelsByPriority(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	models := map[string]MentalModel{
		"low_priority": {
			Name:     "Low Priority Model",
			Priority: 1,
		},
		"high_priority": {
			Name:     "High Priority Model",
			Priority: 10,
		},
		"medium_priority": {
			Name:     "Medium Priority Model",
			Priority: 5,
		},
		"same_priority_1": {
			Name:     "Same Priority A",
			Priority: 3,
		},
		"same_priority_2": {
			Name:     "Same Priority B",
			Priority: 3,
		},
	}

	sorted := loader.GetModelsByPriority(models)

	require.Len(t, sorted, 5)

	// Check that they're sorted by priority (highest first)
	assert.Equal(t, "high_priority", sorted[0].Key)
	assert.Equal(t, "medium_priority", sorted[1].Key)
	assert.Equal(t, "same_priority_1", sorted[2].Key) // Should be sorted by name when priority is equal
	assert.Equal(t, "same_priority_2", sorted[3].Key)
	assert.Equal(t, "low_priority", sorted[4].Key)
}

func TestGetModelsByCategory(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	models := map[string]MentalModel{
		"analytical_1": {
			Name:     "Analytical Model 1",
			Category: "analytical",
			Priority: 5,
		},
		"analytical_2": {
			Name:     "Analytical Model 2",
			Category: "analytical",
			Priority: 10,
		},
		"decision_1": {
			Name:     "Decision Model 1",
			Category: "decision-making",
			Priority: 3,
		},
		"custom_1": {
			Name:     "Custom Model 1",
			Category: "custom",
			Priority: 1,
		},
	}

	categorized := loader.GetModelsByCategory(models)

	require.Len(t, categorized, 3)

	// Check analytical category
	analytical, exists := categorized["analytical"]
	require.True(t, exists)
	require.Len(t, analytical, 2)
	assert.Equal(t, "analytical_2", analytical[0].Key) // Higher priority first
	assert.Equal(t, "analytical_1", analytical[1].Key)

	// Check decision-making category
	decision, exists := categorized["decision-making"]
	require.True(t, exists)
	require.Len(t, decision, 1)
	assert.Equal(t, "decision_1", decision[0].Key)

	// Check custom category
	custom, exists := categorized["custom"]
	require.True(t, exists)
	require.Len(t, custom, 1)
	assert.Equal(t, "custom_1", custom[0].Key)
}

func TestGetAvailableModels(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	models := map[string]MentalModel{
		"model_1": {
			Name: "First Model",
		},
		"model_2": {
			Name: "Second Model",
		},
	}

	available := loader.GetAvailableModels(models)

	require.Len(t, available, 2)
	assert.Contains(t, available, "model_1: First Model")
	assert.Contains(t, available, "model_2: Second Model")
}

func TestLoadCustomModels_FileNotExists(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	models, err := loader.loadCustomModels("/nonexistent/file.yaml")

	require.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestLoadCustomModels_InvalidYAML(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary file with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	require.NoError(t, err)

	models, err := loader.loadCustomModels(configPath)

	require.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestLoadCustomModels_ValidYAML(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary YAML file
	yamlContent := `
models:
  test_model:
    name: "Test Model"
    description: "A test model"
    steps:
      - "Step 1"
      - "Step 2"
    category: "test"
    priority: 5
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	models, err := loader.loadCustomModels(configPath)

	require.NoError(t, err)
	require.Len(t, models, 1)

	model, exists := models["test_model"]
	require.True(t, exists)
	assert.Equal(t, "Test Model", model.Name)
	assert.Equal(t, "A test model", model.Description)
	assert.Equal(t, "test", model.Category)
	assert.Equal(t, 5, model.Priority)
	assert.Len(t, model.Steps, 2)
	assert.Equal(t, "Step 1", model.Steps[0])
	assert.Equal(t, "Step 2", model.Steps[1])
}

func TestLoadCustomModels_DefaultPriority(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary YAML file without priority
	yamlContent := `
models:
  test_model:
    name: "Test Model"
    description: "A test model"
    steps:
      - "Step 1"
    category: "test"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	models, err := loader.loadCustomModels(configPath)

	require.NoError(t, err)
	require.Len(t, models, 1)

	model, exists := models["test_model"]
	require.True(t, exists)
	assert.Equal(t, 1, model.Priority) // Should get default priority of 1
}

func TestLoadMentalModels_WithDirectory(t *testing.T) {
	logger := logrus.New()
	loader := NewLoader(logger)

	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create first YAML file
	yamlContent1 := `
models:
  model_1:
    name: "Model 1"
    description: "Description 1"
    steps: ["Step 1"]
    category: "cat1"
`
	err := os.WriteFile(filepath.Join(tmpDir, "models1.yaml"), []byte(yamlContent1), 0644)
	require.NoError(t, err)

	// Create second YAML file
	yamlContent2 := `
models:
  model_2:
    name: "Model 2"
    description: "Description 2"
    steps: ["Step 1"]
    category: "cat2"
`
	err = os.WriteFile(filepath.Join(tmpDir, "models2.yaml"), []byte(yamlContent2), 0644)
	require.NoError(t, err)

	// Create a non-YAML file (should be ignored)
	err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignore me"), 0644)
	require.NoError(t, err)

	// Load models from directory
	models, err := loader.LoadMentalModels(tmpDir)

	require.NoError(t, err)
	assert.NotEmpty(t, models)

	// Check that models from both files are loaded
	assert.Contains(t, models, "model_1")
	assert.Contains(t, models, "model_2")
	assert.Contains(t, models, "first_principles") // Core models should still be there
}
