package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the server configuration
type Config struct {
	// Server settings
	Port         string        `json:"port" yaml:"port"`
	Host         string        `json:"host" yaml:"host"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`

	// Session settings
	SessionTimeout        time.Duration `json:"session_timeout" yaml:"session_timeout"`
	MaxThoughtsPerSession int           `json:"max_thoughts_per_session" yaml:"max_thoughts_per_session"`

	// Feature flags
	EnableStochasticAlgorithms bool `json:"enable_stochastic_algorithms" yaml:"enable_stochastic_algorithms"`
	EnableSystematicThinking   bool `json:"enable_systematic_thinking" yaml:"enable_systematic_thinking"`
	EnableVisualization        bool `json:"enable_visualization" yaml:"enable_visualization"`
	EnableHybridThinking       bool `json:"enable_hybrid_thinking" yaml:"enable_hybrid_thinking"`

	// Algorithm settings
	MaxStochasticIterations    int     `json:"max_stochastic_iterations" yaml:"max_stochastic_iterations"`
	DefaultConfidenceThreshold float64 `json:"default_confidence_threshold" yaml:"default_confidence_threshold"`

	// Persistence settings
	EnablePersistence bool   `json:"enable_persistence" yaml:"enable_persistence"`
	PersistencePath   string `json:"persistence_path" yaml:"persistence_path"`

	// Logging settings
	EnableDetailedLogging bool   `json:"enable_detailed_logging" yaml:"enable_detailed_logging"`
	LogLevel              string `json:"log_level" yaml:"log_level"`

	// Algorithm defaults
	AlgorithmDefaults map[string]interface{} `json:"algorithm_defaults" yaml:"algorithm_defaults"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:                       "8080",
		Host:                       "localhost",
		ReadTimeout:                30 * time.Second,
		WriteTimeout:               30 * time.Second,
		SessionTimeout:             30 * time.Minute,
		MaxThoughtsPerSession:      100,
		EnableStochasticAlgorithms: true,
		EnableSystematicThinking:   true,
		EnableVisualization:        true,
		EnableHybridThinking:       true,
		MaxStochasticIterations:    1000,
		DefaultConfidenceThreshold: 0.8,
		EnablePersistence:          false,
		EnableDetailedLogging:      false,
		LogLevel:                   "info",
		AlgorithmDefaults:          make(map[string]interface{}),
	}
}

// Load loads configuration from file or environment variables
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from config file
	if configFile := os.Getenv("GOTHINK_CONFIG"); configFile != "" {
		if err := loadFromFile(cfg, configFile); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(cfg *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	if port := os.Getenv("GOTHINK_PORT"); port != "" {
		cfg.Port = port
	}
	if host := os.Getenv("GOTHINK_HOST"); host != "" {
		cfg.Host = host
	}
	if enableStochastic := os.Getenv("GOTHINK_ENABLE_STOCHASTIC"); enableStochastic == "false" {
		cfg.EnableStochasticAlgorithms = false
	}
	if enableSystematic := os.Getenv("GOTHINK_ENABLE_SYSTEMATIC"); enableSystematic == "false" {
		cfg.EnableSystematicThinking = false
	}
	if enableVisualization := os.Getenv("GOTHINK_ENABLE_VISUALIZATION"); enableVisualization == "false" {
		cfg.EnableVisualization = false
	}
	if enableHybrid := os.Getenv("GOTHINK_ENABLE_HYBRID"); enableHybrid == "false" {
		cfg.EnableHybridThinking = false
	}
	if logLevel := os.Getenv("GOTHINK_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}
