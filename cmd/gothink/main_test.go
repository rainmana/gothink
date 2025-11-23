package main

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rainmana/gothink/internal/config"
	"github.com/rainmana/gothink/internal/models"
	"github.com/rainmana/gothink/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerSetup(t *testing.T) {
	// Setup dependencies
	cfg := config.DefaultConfig()
	store, err := storage.New(cfg)
	require.NoError(t, err)

	logger := logrus.New()
	modelsLoader := models.NewLoader(logger)

	// Create server
	s := server.NewMCPServer(
		"GoThink Test Server",
		"1.0.0",
	)

	// Add tools
	addThinkingTools(s, store, modelsLoader, cfg)
	addSessionTools(s, store)

	// Verify tools are registered
	// Note: mcp-go doesn't expose a way to list tools directly from the server struct easily without using the protocol,
	// but we can verify that adding them didn't panic.
	// In a real integration test, we would start the server and connect a client.
	// For now, we'll just verify the setup functions run without error.
}

func TestAddThinkingTools(t *testing.T) {
	cfg := config.DefaultConfig()
	store, _ := storage.New(cfg)
	logger := logrus.New()
	modelsLoader := models.NewLoader(logger)
	s := server.NewMCPServer("Test", "1.0.0")

	addThinkingTools(s, store, modelsLoader, cfg)

	// We can't easily inspect s.tools without private access or running the server,
	// but successful execution implies tools were added.
}

func TestAddSessionTools(t *testing.T) {
	cfg := config.DefaultConfig()
	store, _ := storage.New(cfg)
	s := server.NewMCPServer("Test", "1.0.0")

	addSessionTools(s, store)
}

func TestHandleSequentialThinking(t *testing.T) {
	cfg := config.DefaultConfig()
	store, err := storage.New(cfg)
	require.NoError(t, err)

	sessionID := "test-session"
	thought := "This is a test thought"
	thoughtNumber := 1
	totalThoughts := 5
	nextThoughtNeeded := true

	// Initialize session
	_, err = store.CreateSession(sessionID)
	require.NoError(t, err)

	// Call handler
	result, err := handleSequentialThinking(store, sessionID, thought, thoughtNumber, totalThoughts, nextThoughtNeeded)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify thought was stored
	thoughts, err := store.GetThoughts(sessionID)
	require.NoError(t, err)
	require.Len(t, thoughts, 1)
	assert.Equal(t, thought, thoughts[0].Thought)
	assert.Equal(t, thoughtNumber, thoughts[0].ThoughtNumber)

	// Verify response contains expected fields
	assert.Contains(t, result, "success")
	assert.Contains(t, result, thoughts[0].ID)
}
