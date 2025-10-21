package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/rainmana/gothink/internal/config"
	"github.com/rainmana/gothink/internal/handlers"
	"github.com/rainmana/gothink/internal/middleware"
	"github.com/rainmana/gothink/internal/storage"
)

// Server represents the GoThink MCP server
type Server struct {
	config  *config.Config
	router  *mux.Router
	storage *storage.Storage
	logger  *logrus.Logger
	server  *http.Server
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Setup logger
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create storage
	storage, err := storage.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	// Create router
	router := mux.NewRouter()

	// Create server
	srv := &Server{
		config:  cfg,
		router:  router,
		storage: storage,
		logger:  logger,
	}

	// Setup routes
	srv.setupRoutes()

	// Create HTTP server
	srv.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return srv, nil
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.WithFields(logrus.Fields{
		"host": s.config.Host,
		"port": s.config.Port,
	}).Info("Starting GoThink server")

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down GoThink server")
	return s.server.Shutdown(ctx)
}

// setupRoutes configures all the server routes
func (s *Server) setupRoutes() {
	// Create handlers
	thinkingHandler := handlers.NewThinkingHandler(s.storage, s.logger)
	stochasticHandler := handlers.NewStochasticHandler(s.storage, s.logger)
	decisionHandler := handlers.NewDecisionHandler(s.storage, s.logger)
	visualHandler := handlers.NewVisualHandler(s.storage, s.logger)
	sessionHandler := handlers.NewSessionHandler(s.storage, s.logger)

	// Middleware
	s.router.Use(middleware.Logging(s.logger))
	s.router.Use(middleware.CORS())
	s.router.Use(middleware.JSON())

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Thinking tools
	if s.config.EnableSystematicThinking {
		api.HandleFunc("/thinking/sequential", thinkingHandler.SequentialThinking).Methods("POST")
		api.HandleFunc("/thinking/mental-model", thinkingHandler.MentalModel).Methods("POST")
		api.HandleFunc("/thinking/debugging", thinkingHandler.DebuggingApproach).Methods("POST")
		api.HandleFunc("/thinking/collaborative", thinkingHandler.CollaborativeReasoning).Methods("POST")
		api.HandleFunc("/thinking/socratic", thinkingHandler.SocraticMethod).Methods("POST")
		api.HandleFunc("/thinking/creative", thinkingHandler.CreativeThinking).Methods("POST")
		api.HandleFunc("/thinking/systems", thinkingHandler.SystemsThinking).Methods("POST")
		api.HandleFunc("/thinking/scientific", thinkingHandler.ScientificMethod).Methods("POST")
	}

	// Stochastic algorithms
	if s.config.EnableStochasticAlgorithms {
		api.HandleFunc("/stochastic/mdp", stochasticHandler.MarkovDecisionProcess).Methods("POST")
		api.HandleFunc("/stochastic/mcts", stochasticHandler.MonteCarloTreeSearch).Methods("POST")
		api.HandleFunc("/stochastic/bandit", stochasticHandler.MultiArmedBandit).Methods("POST")
		api.HandleFunc("/stochastic/bayesian", stochasticHandler.BayesianOptimization).Methods("POST")
		api.HandleFunc("/stochastic/hmm", stochasticHandler.HiddenMarkovModel).Methods("POST")
		api.HandleFunc("/stochastic/reinforcement", stochasticHandler.ReinforcementLearning).Methods("POST")
	}

	// Decision frameworks
	api.HandleFunc("/decision/framework", decisionHandler.DecisionFramework).Methods("POST")
	api.HandleFunc("/decision/expected-utility", decisionHandler.ExpectedUtility).Methods("POST")
	api.HandleFunc("/decision/multi-criteria", decisionHandler.MultiCriteria).Methods("POST")
	api.HandleFunc("/decision/risk-analysis", decisionHandler.RiskAnalysis).Methods("POST")

	// Visualization tools
	if s.config.EnableVisualization {
		api.HandleFunc("/visual/concept-map", visualHandler.ConceptMap).Methods("POST")
		api.HandleFunc("/visual/mind-map", visualHandler.MindMap).Methods("POST")
		api.HandleFunc("/visual/flowchart", visualHandler.Flowchart).Methods("POST")
		api.HandleFunc("/visual/decision-tree", visualHandler.DecisionTree).Methods("POST")
		api.HandleFunc("/visual/probability-tree", visualHandler.ProbabilityTree).Methods("POST")
		api.HandleFunc("/visual/bayesian-network", visualHandler.BayesianNetwork).Methods("POST")
	}

	// Session management
	api.HandleFunc("/session/{id}/stats", sessionHandler.GetStats).Methods("GET")
	api.HandleFunc("/session/{id}/export", sessionHandler.Export).Methods("GET")
	api.HandleFunc("/session/{id}/import", sessionHandler.Import).Methods("POST")
	api.HandleFunc("/session/{id}/clear", sessionHandler.Clear).Methods("DELETE")

	// Hybrid thinking (combines systematic and stochastic)
	if s.config.EnableHybridThinking {
		api.HandleFunc("/hybrid/adaptive-reasoning", s.hybridAdaptiveReasoning).Methods("POST")
		api.HandleFunc("/hybrid/probabilistic-decision", s.hybridProbabilisticDecision).Methods("POST")
		api.HandleFunc("/hybrid/uncertainty-analysis", s.hybridUncertaintyAnalysis).Methods("POST")
	}
}

// healthCheck handles health check requests
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "0.1.0",
		"features": map[string]bool{
			"stochastic_algorithms": s.config.EnableStochasticAlgorithms,
			"systematic_thinking":   s.config.EnableSystematicThinking,
			"visualization":         s.config.EnableVisualization,
			"hybrid_thinking":       s.config.EnableHybridThinking,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// hybridAdaptiveReasoning combines systematic and stochastic reasoning
func (s *Server) hybridAdaptiveReasoning(w http.ResponseWriter, r *http.Request) {
	// This would implement adaptive reasoning that switches between
	// systematic and stochastic approaches based on problem characteristics
	response := map[string]interface{}{
		"message": "Hybrid adaptive reasoning not yet implemented",
		"status":  "coming_soon",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// hybridProbabilisticDecision combines decision frameworks with stochastic analysis
func (s *Server) hybridProbabilisticDecision(w http.ResponseWriter, r *http.Request) {
	// This would implement probabilistic decision making that combines
	// traditional decision frameworks with stochastic algorithms
	response := map[string]interface{}{
		"message": "Hybrid probabilistic decision making not yet implemented",
		"status":  "coming_soon",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// hybridUncertaintyAnalysis provides comprehensive uncertainty analysis
func (s *Server) hybridUncertaintyAnalysis(w http.ResponseWriter, r *http.Request) {
	// This would implement uncertainty analysis that combines
	// multiple approaches to quantify and manage uncertainty
	response := map[string]interface{}{
		"message": "Hybrid uncertainty analysis not yet implemented",
		"status":  "coming_soon",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getLogLevel converts string to logrus level
func getLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
