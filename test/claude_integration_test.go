package test

import (
	"dispatch-mcp-server/internal/conversation"
	"os"
	"testing"
)

func TestClaudeConversationEngine(t *testing.T) {
	// Test creating the Claude conversation engine
	engine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		t.Logf("Claude not available (expected if no API key): %v", err)
		return
	}

	// Test engine info
	info := engine.GetEngineInfo()
	if info["engine_type"] != "claude_hybrid" {
		t.Errorf("Expected engine_type to be 'claude_hybrid', got %v", info["engine_type"])
	}

	// Test Claude availability
	claudeAvailable := engine.IsClaudeAvailable()
	t.Logf("Claude available: %v", claudeAvailable)

	// Test processing a simple message
	context := &conversation.ConversationContext{
		SessionID: "test_session",
		CustomerProfile: conversation.CustomerProfile{
			Tier:               "gold",
			OrderFrequency:     5,
			AverageOrderValue:  100.0,
		},
	}

	response, err := engine.ProcessMessage("I need 3 deliveries", context)
	if err != nil {
		t.Errorf("Failed to process message: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Response message: %s", response.Message)
	t.Logf("Recommendations count: %d", len(response.Recommendations))
	t.Logf("Next questions count: %d", len(response.NextQuestions))

	// Test that we get some response
	if response.Message == "" {
		t.Error("Expected non-empty response message")
	}
}

func TestClaudeFallback(t *testing.T) {
	// Test that the engine falls back to rule-based when Claude is not available
	// This simulates the case where ANTHROPIC_API_KEY is not set
	
	// Temporarily unset the API key
	originalKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Setenv("ANTHROPIC_API_KEY", originalKey)

	engine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		t.Logf("Expected error when Claude is not available: %v", err)
		return
	}

	// Should fall back to rule-based processing
	if engine.IsClaudeAvailable() {
		t.Error("Expected Claude to be unavailable")
	}

	// Test that it still works with rule-based processing
	context := &conversation.ConversationContext{
		SessionID: "test_session",
		CustomerProfile: conversation.CustomerProfile{
			Tier:               "bronze",
			OrderFrequency:     2,
			AverageOrderValue:  50.0,
		},
	}

	response, err := engine.ProcessMessage("What pricing options do I have?", context)
	if err != nil {
		t.Errorf("Failed to process message with fallback: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Fallback response: %s", response.Message)
}

func TestClaudeContextConversion(t *testing.T) {
	engine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		t.Logf("Claude not available: %v", err)
		return
	}

	// Test context conversion
	context := &conversation.ConversationContext{
		SessionID: "test_session",
		CustomerProfile: conversation.CustomerProfile{
			Tier:               "gold",
			OrderFrequency:     5,
			AverageOrderValue:  150.0,
		},
	}

	// Test that the engine can handle the context
	response, err := engine.ProcessMessage("I'm a gold tier customer with 5 orders per month", context)
	if err != nil {
		t.Errorf("Failed to process message: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	// Check that the context was updated
	if response.UpdatedContext == nil {
		t.Error("Expected updated context")
		return
	}

	t.Logf("Updated context tier: %s", response.UpdatedContext.CustomerProfile.Tier)
	t.Logf("Updated context frequency: %d", response.UpdatedContext.CustomerProfile.OrderFrequency)
}

func TestClaudeRecommendations(t *testing.T) {
	engine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		t.Logf("Claude not available: %v", err)
		return
	}

	// Test with a context that should generate recommendations
	context := &conversation.ConversationContext{
		SessionID: "test_session",
		CustomerProfile: conversation.CustomerProfile{
			Tier:               "gold",
			OrderFrequency:     5,
			AverageOrderValue:  200.0,
		},
	}

	response, err := engine.ProcessMessage("What's the best pricing for me?", context)
	if err != nil {
		t.Errorf("Failed to process message: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	// Check that we got recommendations
	if len(response.Recommendations) == 0 {
		t.Error("Expected pricing recommendations")
		return
	}

	t.Logf("Got %d recommendations:", len(response.Recommendations))
	for i, rec := range response.Recommendations {
		t.Logf("  %d. %s: %s (%.1f%% savings)", i+1, rec.Name, rec.Reason, rec.SavingsPercent)
	}
}
