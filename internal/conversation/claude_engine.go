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

	// Update context with new information from the message first
	updatedContext := ce.updateContextFromMessage(message, context)

	// Generate pricing recommendations using our pricing engine with updated context
	recommendations := ce.generateRecommendations(updatedContext)

	// Generate next questions based on updated context
	nextQuestions := ce.generateNextQuestions(updatedContext)

	// Debug: Print context information
	fmt.Printf("DEBUG: Updated context - DeliveryCount: %d, Tier: %s, OrderFrequency: %d\n",
		updatedContext.CustomerProfile.CurrentDeliveryCount,
		updatedContext.CustomerProfile.Tier,
		updatedContext.CustomerProfile.OrderFrequency)

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
			CustomerTier:    "bronze",
			OrderFrequency:  1,
			TotalOrderValue: 50.0,
			IsBulkOrder:     false,
			OrderCreation:   claude.OrderCreationState{},
		}
	}

	return &claude.PricingContext{
		DeliveryCount:   context.CustomerProfile.CurrentDeliveryCount,
		CustomerTier:    context.CustomerProfile.Tier,
		OrderFrequency:  context.CustomerProfile.OrderFrequency,
		TotalOrderValue: context.CustomerProfile.AverageOrderValue,
		IsBulkOrder:     false, // Can be enhanced based on conversation
		OrderCreation: claude.OrderCreationState{
			InProgress:           context.OrderCreation.InProgress,
			Step:                 context.OrderCreation.Step,
			CurrentQuestion:      context.OrderCreation.CurrentQuestion,
			PickupInfo:           context.OrderCreation.PickupInfo,
			DropOffs:             context.OrderCreation.DropOffs,
			DeliveryInfo:         context.OrderCreation.DeliveryInfo,
			MissingFields:        context.OrderCreation.MissingFields,
			CompletedFields:      context.OrderCreation.CompletedFields,
			CurrentDeliveryIndex: context.OrderCreation.CurrentDeliveryIndex,
		},
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

	// Extract delivery count with more sophisticated parsing
	// Look for delivery-related terms and numbers
	deliveryTerms := []string{"deliver", "delivery", "deliveries", "package", "packages", "shipment", "shipments", "drop", "drops"}

	for _, term := range deliveryTerms {
		if strings.Contains(message, term) {
			// Look for numbers in the message
			words := strings.Fields(message)
			for i, word := range words {
				// Check if current word is a delivery term
				if word == term || word == term+"s" {
					// Check previous word for number
					if i > 0 {
						switch words[i-1] {
						case "one", "1":
							context.CustomerProfile.CurrentDeliveryCount = 1
						case "two", "2":
							context.CustomerProfile.CurrentDeliveryCount = 2
						case "three", "3":
							context.CustomerProfile.CurrentDeliveryCount = 3
						case "four", "4":
							context.CustomerProfile.CurrentDeliveryCount = 4
						case "five", "5":
							context.CustomerProfile.CurrentDeliveryCount = 5
						case "six", "6":
							context.CustomerProfile.CurrentDeliveryCount = 6
						case "seven", "7":
							context.CustomerProfile.CurrentDeliveryCount = 7
						case "eight", "8":
							context.CustomerProfile.CurrentDeliveryCount = 8
						case "nine", "9":
							context.CustomerProfile.CurrentDeliveryCount = 9
						case "ten", "10":
							context.CustomerProfile.CurrentDeliveryCount = 10
						}
					}
					// Also check next word for number (e.g., "deliver 3 packages")
					if i < len(words)-1 {
						switch words[i+1] {
						case "one", "1":
							context.CustomerProfile.CurrentDeliveryCount = 1
						case "two", "2":
							context.CustomerProfile.CurrentDeliveryCount = 2
						case "three", "3":
							context.CustomerProfile.CurrentDeliveryCount = 3
						case "four", "4":
							context.CustomerProfile.CurrentDeliveryCount = 4
						case "five", "5":
							context.CustomerProfile.CurrentDeliveryCount = 5
						case "six", "6":
							context.CustomerProfile.CurrentDeliveryCount = 6
						case "seven", "7":
							context.CustomerProfile.CurrentDeliveryCount = 7
						case "eight", "8":
							context.CustomerProfile.CurrentDeliveryCount = 8
						case "nine", "9":
							context.CustomerProfile.CurrentDeliveryCount = 9
						case "ten", "10":
							context.CustomerProfile.CurrentDeliveryCount = 10
						}
					}
				}
			}
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

	// Handle order creation progress
	ce.updateOrderCreationProgress(message, context)

	return context
}

// updateOrderCreationProgress handles step-by-step order creation
func (ce *ClaudeConversationEngine) updateOrderCreationProgress(message string, context *ConversationContext) {
	// Initialize order creation if not started
	if !context.OrderCreation.InProgress {
		// Check if user wants to create an order
		if strings.Contains(message, "create") && strings.Contains(message, "order") {
			context.OrderCreation.InProgress = true
			context.OrderCreation.Step = "pickup"
			context.OrderCreation.CurrentQuestion = "pickup_business"
			context.OrderCreation.PickupInfo = &dispatch.CreateOrderPickupInfoInput{}
			context.OrderCreation.DropOffs = []dispatch.CreateOrderDropOffInfoInput{}
			context.OrderCreation.CompletedFields = []string{}
			context.OrderCreation.MissingFields = []string{"pickup_business", "pickup_address", "pickup_contact", "pickup_phone"}
		}
		return
	}

	// Handle current question based on step and question
	switch context.OrderCreation.Step {
	case "pickup":
		ce.handlePickupStep(message, context)
	case "deliveries":
		ce.handleDeliveriesStep(message, context)
	case "review":
		ce.handleReviewStep(message, context)
	}
}

// handlePickupStep processes pickup information step by step
func (ce *ClaudeConversationEngine) handlePickupStep(message string, context *ConversationContext) {
	switch context.OrderCreation.CurrentQuestion {
	case "pickup_business":
		// Extract business name from message
		if len(strings.TrimSpace(message)) > 0 {
			context.OrderCreation.PickupInfo.BusinessName = &message
			context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, "pickup_business")
			context.OrderCreation.CurrentQuestion = "pickup_address"
		}
	case "pickup_address":
		// Extract address from message
		if len(strings.TrimSpace(message)) > 0 {
			// Parse address into components (simplified)
			context.OrderCreation.PickupInfo.Location = &dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  message,         // Simplified - would need better parsing
					City:    "San Francisco", // Default - would extract from message
					State:   "CA",
					ZipCode: "94105",
					Country: "US",
				},
			}
			context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, "pickup_address")
			context.OrderCreation.CurrentQuestion = "pickup_contact"
		}
	case "pickup_contact":
		// Extract contact name from message
		if len(strings.TrimSpace(message)) > 0 {
			context.OrderCreation.PickupInfo.ContactName = &message
			context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, "pickup_contact")
			context.OrderCreation.CurrentQuestion = "pickup_phone"
		}
	case "pickup_phone":
		// Extract phone number from message
		if len(strings.TrimSpace(message)) > 0 {
			context.OrderCreation.PickupInfo.ContactPhoneNumber = &message
			context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, "pickup_phone")
			// Move to deliveries step
			context.OrderCreation.Step = "deliveries"
			context.OrderCreation.CurrentQuestion = "delivery_count"
		}
	}
}

// handleDeliveriesStep processes delivery information step by step
func (ce *ClaudeConversationEngine) handleDeliveriesStep(message string, context *ConversationContext) {
	// This would be implemented to handle delivery location collection
	// For now, just move to review step
	context.OrderCreation.Step = "review"
	context.OrderCreation.CurrentQuestion = "confirm_order"
}

// handleReviewStep processes order review and confirmation
func (ce *ClaudeConversationEngine) handleReviewStep(message string, context *ConversationContext) {
	// Handle order confirmation
	if strings.Contains(message, "confirm") || strings.Contains(message, "yes") {
		// Order confirmed - would create the actual order here
		context.OrderCreation.InProgress = false
		context.OrderCreation.Step = "completed"
	}
}

// generateRecommendations generates pricing recommendations using our pricing engine
func (ce *ClaudeConversationEngine) generateRecommendations(context *ConversationContext) []PricingRecommendation {
	if context == nil {
		return []PricingRecommendation{}
	}

	// Create pricing context for our engine
	pricingContext := pricing.PricingContext{
		DeliveryCount:   context.CustomerProfile.CurrentDeliveryCount,
		CustomerTier:    context.CustomerProfile.Tier,
		OrderFrequency:  context.CustomerProfile.OrderFrequency,
		TotalOrderValue: context.CustomerProfile.AverageOrderValue,
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
		"engine_type":      "claude_hybrid",
		"claude_available": ce.IsClaudeAvailable(),
		"fallback_mode":    !ce.IsClaudeAvailable(),
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
