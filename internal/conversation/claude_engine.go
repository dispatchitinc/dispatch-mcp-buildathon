package conversation

import (
	"dispatch-mcp-server/internal/claude"
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/pricing"
	"fmt"
	"strings"
	"time"
)

// ClaudeConversationEngine handles natural language pricing conversations using Claude AI
type ClaudeConversationEngine struct {
	claudeClient   *claude.Client
	pricingEngine  *pricing.PricingEngine
	contextManager *ContextManager
	useClaude      bool
}

// NewClaudeConversationEngine creates a new Claude-powered conversation engine
func NewClaudeConversationEngine() (*ClaudeConversationEngine, error) {
	claudeClient, err := claude.NewClient()
	if err != nil {
		// If Claude is not available, fall back to rule-based engine
		return &ClaudeConversationEngine{
			claudeClient:   nil,
			pricingEngine:  pricing.NewPricingEngine(),
			contextManager: NewContextManager(),
			useClaude:      false,
		}, nil
	}

	return &ClaudeConversationEngine{
		claudeClient:   claudeClient,
		pricingEngine:  pricing.NewPricingEngine(),
		contextManager: NewContextManager(),
		useClaude:      true,
	}, nil
}

// ProcessMessage processes a natural language message using Claude AI
func (ce *ClaudeConversationEngine) ProcessMessage(message string, context *ConversationContext) (*ConversationResponse, error) {
	// If Claude is not available, fall back to rule-based processing
	if !ce.useClaude || ce.claudeClient == nil {
		return ce.processWithRules(message, context)
	}

	// Convert our context to Claude's pricing context
	pricingContext := ce.convertToPricingContext(context)

	// Get Claude's response
	claudeResponse, err := ce.claudeClient.CreatePricingAdvisorMessage(message, pricingContext)
	if err != nil {
		// If Claude fails, fall back to rule-based processing
		return ce.processWithRules(message, context)
	}

	// Extract Claude's response text
	responseText := ""
	if len(claudeResponse.Content) > 0 {
		responseText = claudeResponse.Content[0].Text
	}

	// Generate pricing recommendations using our pricing engine
	recommendations := ce.generateRecommendations(context)

	// Generate next questions based on context
	nextQuestions := ce.generateNextQuestions(context)

	// Update context with new information from the message
	updatedContext := ce.updateContextFromMessage(message, context)

	return &ConversationResponse{
		Message:         responseText,
		Recommendations: recommendations,
		NextQuestions:   nextQuestions,
		UpdatedContext:  updatedContext,
	}, nil
}

// processWithRules processes the message using our rule-based system
func (ce *ClaudeConversationEngine) processWithRules(message string, context *ConversationContext) (*ConversationResponse, error) {
	// Use the original rule-based engine
	originalEngine := NewConversationEngine()
	return originalEngine.ProcessMessage(message, context)
}

// convertToPricingContext converts our conversation context to Claude's pricing context
func (ce *ClaudeConversationEngine) convertToPricingContext(context *ConversationContext) *claude.PricingContext {
	if context == nil {
		return &claude.PricingContext{
			DeliveryCount:   1,
			CustomerTier:     "bronze",
			OrderFrequency:   1,
			TotalOrderValue:  50.0,
			IsBulkOrder:      false,
		}
	}

	return &claude.PricingContext{
		DeliveryCount:   context.CustomerProfile.OrderFrequency,
		CustomerTier:    context.CustomerProfile.Tier,
		OrderFrequency:   context.CustomerProfile.OrderFrequency,
		TotalOrderValue:  context.CustomerProfile.AverageOrderValue,
		IsBulkOrder:     false, // Can be enhanced based on conversation
	}
}

// updateContextFromMessage updates the conversation context based on the message
func (ce *ClaudeConversationEngine) updateContextFromMessage(message string, context *ConversationContext) *ConversationContext {
	if context == nil {
		context = &ConversationContext{
			SessionID:       fmt.Sprintf("session_%d", time.Now().Unix()),
			CustomerProfile: CustomerProfile{Tier: "bronze"},
			DeliveryHistory: []DeliveryRequirement{},
			PricingHistory:  []PricingComparison{},
			Preferences:     CustomerPreferences{},
		}
	}

	// Simple entity extraction for context updates
	message = strings.ToLower(message)

	// Extract delivery count
	if strings.Contains(message, "deliver") {
		// Simple extraction - can be enhanced
		if strings.Contains(message, "3") {
			context.CustomerProfile.OrderFrequency = 3
		} else if strings.Contains(message, "5") {
			context.CustomerProfile.OrderFrequency = 5
		} else if strings.Contains(message, "10") {
			context.CustomerProfile.OrderFrequency = 10
		}
	}

	// Extract customer tier
	if strings.Contains(message, "gold") {
		context.CustomerProfile.Tier = "gold"
	} else if strings.Contains(message, "silver") {
		context.CustomerProfile.Tier = "silver"
	} else if strings.Contains(message, "bronze") {
		context.CustomerProfile.Tier = "bronze"
	}

	// Extract order frequency
	if strings.Contains(message, "month") {
		if strings.Contains(message, "5") {
			context.CustomerProfile.OrderFrequency = 5
		} else if strings.Contains(message, "10") {
			context.CustomerProfile.OrderFrequency = 10
		}
	}

	return context
}

// generateRecommendations generates pricing recommendations using our pricing engine
func (ce *ClaudeConversationEngine) generateRecommendations(context *ConversationContext) []PricingRecommendation {
	if context == nil {
		return []PricingRecommendation{}
	}

	// Create pricing context for our engine
	pricingContext := pricing.PricingContext{
		DeliveryCount:   context.CustomerProfile.OrderFrequency,
		CustomerTier:    context.CustomerProfile.Tier,
		OrderFrequency:   context.CustomerProfile.OrderFrequency,
		TotalOrderValue:  context.CustomerProfile.AverageOrderValue,
		IsBulkOrder:     false,
	}

	// Create a sample estimate for comparison
	sampleEstimate := &dispatch.AvailableOrderOption{
		EstimatedOrderCost: 50.0, // Sample cost
	}

	// Get pricing comparison
	comparison := ce.pricingEngine.ComparePricingModels(sampleEstimate, pricingContext)

	// Convert to recommendations
	var recommendations []PricingRecommendation
	for _, result := range comparison.PricingModels {
		recommendations = append(recommendations, PricingRecommendation{
			Model:          string(result.Model),
			Name:           result.Name,
			Savings:        result.Savings,
			SavingsPercent: result.DiscountPercent,
			Eligible:       result.Eligible,
			Reason:         result.Reason,
		})
	}

	return recommendations
}

// generateNextQuestions generates follow-up questions based on context
func (ce *ClaudeConversationEngine) generateNextQuestions(context *ConversationContext) []string {
	questions := []string{}

	if context == nil {
		questions = append(questions, "What's your customer tier? (bronze, silver, gold)")
		questions = append(questions, "How many deliveries do you need?")
		questions = append(questions, "How many orders do you place per month?")
		return questions
	}

	// Ask about missing information
	if context.CustomerProfile.Tier == "" {
		questions = append(questions, "What's your customer tier? (bronze, silver, gold)")
	}

	if context.CustomerProfile.OrderFrequency == 0 {
		questions = append(questions, "How many orders do you place per month?")
	}

	if len(context.DeliveryHistory) == 0 {
		questions = append(questions, "How many deliveries do you need for this order?")
	}

	// Suggest optimizations
	if context.CustomerProfile.Tier == "bronze" {
		questions = append(questions, "Would you like to learn about our loyalty program?")
	}

	if context.CustomerProfile.OrderFrequency < 3 {
		questions = append(questions, "Are you interested in increasing your order frequency for better pricing?")
	}

	return questions
}

// IsClaudeAvailable returns whether Claude is available
func (ce *ClaudeConversationEngine) IsClaudeAvailable() bool {
	return ce.useClaude && ce.claudeClient != nil
}

// GetEngineInfo returns information about the engine
func (ce *ClaudeConversationEngine) GetEngineInfo() map[string]interface{} {
	info := map[string]interface{}{
		"engine_type": "claude_hybrid",
		"claude_available": ce.IsClaudeAvailable(),
		"fallback_mode": !ce.IsClaudeAvailable(),
	}

	if ce.IsClaudeAvailable() {
		info["claude_model"] = "claude-3-sonnet-20240229"
		info["features"] = []string{
			"natural_language_understanding",
			"contextual_responses",
			"intelligent_recommendations",
			"conversational_flow",
		}
	} else {
		info["features"] = []string{
			"rule_based_processing",
			"pattern_matching",
			"basic_intent_recognition",
		}
	}

	return info
}
