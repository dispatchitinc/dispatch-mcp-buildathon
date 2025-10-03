package test

import (
	"dispatch-mcp-server/internal/conversation"
	"fmt"
	"strings"
	"testing"
)

func TestConversationEngine(t *testing.T) {
	engine := conversation.NewConversationEngine()

	t.Run("basic_message_processing", func(t *testing.T) {
		response, err := engine.ProcessMessage("I need 3 deliveries", nil)
		if err != nil {
			t.Fatalf("ProcessMessage failed: %v", err)
		}

		if response.Message == "" {
			t.Error("Expected non-empty response message")
		}

		if response.UpdatedContext == nil {
			t.Error("Expected updated context")
		}
	})

	t.Run("context_building", func(t *testing.T) {
		// Start with no context
		response1, err := engine.ProcessMessage("I need 3 deliveries", nil)
		if err != nil {
			t.Fatalf("First message failed: %v", err)
		}

		// Continue with context
		response2, err := engine.ProcessMessage("I'm a gold tier customer", response1.UpdatedContext)
		if err != nil {
			t.Fatalf("Second message failed: %v", err)
		}

		// Check that context is maintained
		if response2.UpdatedContext.CustomerProfile.Tier != "gold" {
			t.Errorf("Expected gold tier, got %s", response2.UpdatedContext.CustomerProfile.Tier)
		}
	})

	t.Run("pricing_recommendations", func(t *testing.T) {
		context := &conversation.ConversationContext{
			CustomerProfile: conversation.CustomerProfile{
				Tier:           "gold",
				OrderFrequency: 5,
			},
		}

		response, err := engine.ProcessMessage("show me pricing options", context)
		if err != nil {
			t.Fatalf("ProcessMessage failed: %v", err)
		}

		if len(response.Recommendations) == 0 {
			t.Error("Expected pricing recommendations")
		}

		// Check that at least one recommendation is eligible
		eligibleCount := 0
		for _, rec := range response.Recommendations {
			if rec.Eligible {
				eligibleCount++
			}
		}

		if eligibleCount == 0 {
			t.Error("Expected at least one eligible recommendation for gold tier customer")
		}
	})

	t.Run("intent_recognition", func(t *testing.T) {
		testCases := []struct {
			message          string
			expectedKeywords []string
		}{
			{
				"I need 3 deliveries",
				[]string{"delivery", "Multi-Delivery"},
			},
			{
				"I'm a gold tier customer",
				[]string{"gold", "Loyalty"},
			},
			{
				"show me pricing options",
				[]string{"pricing", "options"},
			},
		}

		for _, tc := range testCases {
			response, err := engine.ProcessMessage(tc.message, nil)
			if err != nil {
				t.Fatalf("ProcessMessage failed for '%s': %v", tc.message, err)
			}

			for _, keyword := range tc.expectedKeywords {
				if !contains(response.Message, keyword) {
					t.Errorf("Expected response to contain '%s' for message '%s', got: %s",
						keyword, tc.message, response.Message)
				}
			}
		}
	})

	t.Run("next_questions", func(t *testing.T) {
		response, err := engine.ProcessMessage("hello", nil)
		if err != nil {
			t.Fatalf("ProcessMessage failed: %v", err)
		}

		if len(response.NextQuestions) == 0 {
			t.Error("Expected next questions for new conversation")
		}

		for _, question := range response.NextQuestions {
			if question == "" {
				t.Error("Expected non-empty next question")
			}
		}
	})

	t.Run("error_handling", func(t *testing.T) {
		// Test with empty message
		response, err := engine.ProcessMessage("", nil)
		if err != nil {
			t.Fatalf("Empty message should not cause error: %v", err)
		}

		if response.Message == "" {
			t.Error("Expected response even for empty message")
		}
	})
}

func TestContextManager(t *testing.T) {
	manager := conversation.NewContextManager()

	t.Run("session_management", func(t *testing.T) {
		// Create a context
		context := &conversation.ConversationContext{
			SessionID: "test_session",
			CustomerProfile: conversation.CustomerProfile{
				Tier: "gold",
			},
		}

		// Save session
		manager.SaveSession(context)

		// Retrieve session
		retrieved := manager.GetSession("test_session")
		if retrieved == nil {
			t.Error("Expected to retrieve saved session")
		}

		if retrieved.CustomerProfile.Tier != "gold" {
			t.Errorf("Expected gold tier, got %s", retrieved.CustomerProfile.Tier)
		}
	})

	t.Run("session_stats", func(t *testing.T) {
		stats := manager.GetSessionStats()

		if stats["total_sessions"] == nil {
			t.Error("Expected session stats")
		}
	})
}

func TestConversationFlow(t *testing.T) {
	engine := conversation.NewConversationEngine()

	// Simulate a complete conversation
	conversationSteps := []string{
		"Hello",
		"I need 3 deliveries to different locations",
		"I'm a gold tier customer",
		"show me pricing options",
		"what's the best option for me",
	}

	var context *conversation.ConversationContext

	for i, message := range conversationSteps {
		t.Run(fmt.Sprintf("step_%d", i+1), func(t *testing.T) {
			response, err := engine.ProcessMessage(message, context)
			if err != nil {
				t.Fatalf("Step %d failed: %v", i+1, err)
			}

			if response.Message == "" {
				t.Errorf("Step %d: Expected non-empty response", i+1)
			}

			context = response.UpdatedContext
		})
	}

	// Check final context
	if context == nil {
		t.Error("Expected final context")
	}

	if context.CustomerProfile.Tier != "gold" {
		t.Errorf("Expected gold tier, got %s", context.CustomerProfile.Tier)
	}
}

func TestPricingScenarios(t *testing.T) {
	engine := conversation.NewConversationEngine()

	scenarios := []struct {
		name             string
		context          *conversation.ConversationContext
		expectedEligible int
	}{
		{
			name: "bronze_tier_single_delivery",
			context: &conversation.ConversationContext{
				CustomerProfile: conversation.CustomerProfile{
					Tier:           "bronze",
					OrderFrequency: 1,
				},
			},
			expectedEligible: 1, // Only standard pricing
		},
		{
			name: "gold_tier_multiple_deliveries",
			context: &conversation.ConversationContext{
				CustomerProfile: conversation.CustomerProfile{
					Tier:           "gold",
					OrderFrequency: 5,
				},
			},
			expectedEligible: 3, // Standard, Multi-Delivery, Loyalty
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			response, err := engine.ProcessMessage("show me pricing options", scenario.context)
			if err != nil {
				t.Fatalf("ProcessMessage failed: %v", err)
			}

			// Count eligible recommendations
			eligibleCount := 0
			for _, rec := range response.Recommendations {
				if rec.Eligible {
					eligibleCount++
				}
			}

			if eligibleCount < scenario.expectedEligible {
				t.Errorf("Expected at least %d eligible recommendations, got %d",
					scenario.expectedEligible, eligibleCount)
			}
		})
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
