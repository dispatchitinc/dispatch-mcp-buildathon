package conversation

import (
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/pricing"
	"fmt"
	"regexp"
	"strings"
)

// ConversationEngine handles natural language pricing conversations
type ConversationEngine struct {
	pricingEngine  *pricing.PricingEngine
	contextManager *ContextManager
}

// NewConversationEngine creates a new conversation engine
func NewConversationEngine() *ConversationEngine {
	return &ConversationEngine{
		pricingEngine:  pricing.NewPricingEngine(),
		contextManager: NewContextManager(),
	}
}

// ProcessMessage processes a natural language message and returns a conversational response
func (ce *ConversationEngine) ProcessMessage(message string, context *ConversationContext) (*ConversationResponse, error) {
	// Extract intent and entities from the message
	intent := ce.extractIntent(message)

	// Update context with new information
	updatedContext := ce.contextManager.Update(context, intent)

	// Generate pricing recommendations based on context
	recommendations := ce.generateRecommendations(updatedContext)

	// Generate natural language response
	response := ce.generateResponse(intent, recommendations, updatedContext)

	return &ConversationResponse{
		Message:         response,
		Recommendations: recommendations,
		NextQuestions:   ce.generateNextQuestions(updatedContext),
		UpdatedContext:  updatedContext,
	}, nil
}

// ConversationContext holds the state of a conversation
type ConversationContext struct {
	SessionID       string                `json:"session_id"`
	CustomerProfile CustomerProfile       `json:"customer_profile"`
	DeliveryHistory []DeliveryRequirement `json:"delivery_history"`
	PricingHistory  []PricingComparison   `json:"pricing_history"`
	CurrentGoal     string                `json:"current_goal"`
	Preferences     CustomerPreferences   `json:"preferences"`
	OrderCreation   OrderCreationState    `json:"order_creation"`
}

// OrderCreationState tracks the progress of order creation
type OrderCreationState struct {
	InProgress           bool                                   `json:"in_progress"`
	Step                 string                                 `json:"step"`                 // "multi_stop", "pickup", "drop_off", "packages", "organization", "service_level", "vehicle", "add_ons", "delivery", "billing", "special_requirements", "review"
	CurrentQuestion      string                                 `json:"current_question"`     // "pickup_business", "pickup_address", "pickup_contact", "pickup_phone", etc.
	MultiStop            bool                                   `json:"multi_stop,omitempty"` // Whether this is a multi-stop delivery
	PickupInfo           *dispatch.CreateOrderPickupInfoInput   `json:"pickup_info,omitempty"`
	DropOffs             []dispatch.CreateOrderDropOffInfoInput `json:"drop_offs,omitempty"`
	VehicleType          *VehicleTypeInfo                       `json:"vehicle_type,omitempty"`
	Capabilities         []string                               `json:"capabilities,omitempty"`
	DeliveryInfo         *dispatch.DeliveryInfoInput            `json:"delivery_info,omitempty"`
	SchedulingInfo       *SchedulingInfo                        `json:"scheduling_info,omitempty"`
	
	// NEW: Package Details
	PackageDetails       *PackageDetailsInfo                    `json:"package_details,omitempty"`
	
	// NEW: Organization Context
	OrganizationInfo     *OrganizationInfo                      `json:"organization_info,omitempty"`
	
	// NEW: Service Level
	ServiceLevel         *ServiceLevelInfo                      `json:"service_level,omitempty"`
	
	// NEW: Billing Information
	BillingInfo          *BillingInfo                           `json:"billing_info,omitempty"`
	
	// NEW: Special Requirements
	SpecialRequirements  *SpecialRequirementsInfo              `json:"special_requirements,omitempty"`
	
	MissingFields        []string                               `json:"missing_fields"`
	CompletedFields      []string                               `json:"completed_fields"`
	CurrentDeliveryIndex int                                    `json:"current_delivery_index"`      // Which delivery we're collecting info for
	ValidationErrors     []string                               `json:"validation_errors,omitempty"` // Address validation errors
}

// VehicleTypeInfo represents vehicle type selection
type VehicleTypeInfo struct {
	VehicleTypeID   string   `json:"vehicle_type_id"`
	VehicleTypeName string   `json:"vehicle_type_name"`
	CustomTypes     []string `json:"custom_types,omitempty"`
}

// SchedulingInfo represents delivery scheduling information
type SchedulingInfo struct {
	PickupTime    string   `json:"pickup_time"`
	DeliveryTime  string   `json:"delivery_time"`
	PickupDate    string   `json:"pickup_date"`
	DeliveryDate  string   `json:"delivery_date"`
	TimeZone      string   `json:"time_zone"`
	SpecialTiming []string `json:"special_timing,omitempty"`
}

// PackageDetailsInfo represents package information
type PackageDetailsInfo struct {
	PackageCount      int       `json:"package_count"`
	TotalWeight       float64   `json:"total_weight"`
	IndividualWeights []float64 `json:"individual_weights,omitempty"`
	ReferenceNames    []string  `json:"reference_names,omitempty"`
	SpecialHandling   []string  `json:"special_handling,omitempty"`
	Dimensions        *PackageDimensions `json:"dimensions,omitempty"`
}

// PackageDimensions represents package dimensions
type PackageDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // "inches", "cm"
}

// OrganizationInfo represents organization context
type OrganizationInfo struct {
	OrganizationID     string `json:"organization_id"`
	OrganizationName   string `json:"organization_name"`
	BranchID          string `json:"branch_id,omitempty"`
	BranchName        string `json:"branch_name,omitempty"`
	MarketID          string `json:"market_id,omitempty"`
	MarketName        string `json:"market_name,omitempty"`
	CreatedByUserID   string `json:"created_by_user_id"`
	CreatedByUserName string `json:"created_by_user_name"`
	CreatedByUserEmail string `json:"created_by_user_email"`
}

// ServiceLevelInfo represents service level selection
type ServiceLevelInfo struct {
	ServiceLevel       string   `json:"service_level"` // "standard", "express", "rush", "scheduled"
	DeliveryTimeWindow string   `json:"delivery_time_window,omitempty"`
	UrgencyLevel       string   `json:"urgency_level"` // "asap", "scheduled", "flexible"
	SpecialServices    []string `json:"special_services,omitempty"`
}

// BillingInfo represents billing information
type BillingInfo struct {
	BillingMethod     string        `json:"billing_method"` // "account", "credit_card", "invoice", "cod"
	BillingAddress    *AddressInput `json:"billing_address,omitempty"`
	PaymentTerms      string        `json:"payment_terms,omitempty"`
	ContactEmail      string        `json:"contact_email"`
	NotificationPhone string        `json:"notification_phone,omitempty"`
}

// SpecialRequirementsInfo represents special delivery requirements
type SpecialRequirementsInfo struct {
	UnloadingAssistance  bool   `json:"unloading_assistance"`
	DedicatedVehicle     bool   `json:"dedicated_vehicle"`
	DeliveryInstructions string `json:"delivery_instructions,omitempty"`
	AccessInstructions   string `json:"access_instructions,omitempty"`
	SpecialNotes         string `json:"special_notes,omitempty"`
}

// CustomerProfile represents customer information
type CustomerProfile struct {
	Tier                 string   `json:"tier"`
	OrderFrequency       int      `json:"order_frequency"`
	CurrentDeliveryCount int      `json:"current_delivery_count"` // Current order delivery count
	AverageOrderValue    float64  `json:"average_order_value"`
	PreferredVehicle     string   `json:"preferred_vehicle"`
	SpecialNeeds         []string `json:"special_needs"`
}

// DeliveryRequirement represents a delivery need
type DeliveryRequirement struct {
	Count               int      `json:"count"`
	Locations           []string `json:"locations"`
	VehicleType         string   `json:"vehicle_type"`
	SpecialRequirements []string `json:"special_requirements"`
}

// PricingComparison represents a pricing analysis
type PricingComparison struct {
	OriginalCost    float64                 `json:"original_cost"`
	BestOption      string                  `json:"best_option"`
	BestPrice       float64                 `json:"best_price"`
	Savings         float64                 `json:"savings"`
	SavingsPercent  float64                 `json:"savings_percent"`
	Recommendations []PricingRecommendation `json:"recommendations"`
}

// PricingRecommendation represents a pricing suggestion
type PricingRecommendation struct {
	Model          string  `json:"model"`
	Name           string  `json:"name"`
	Savings        float64 `json:"savings"`
	SavingsPercent float64 `json:"savings_percent"`
	Eligible       bool    `json:"eligible"`
	Reason         string  `json:"reason,omitempty"`
}

// ConversationResponse represents a conversational response
type ConversationResponse struct {
	Message         string                  `json:"message"`
	Recommendations []PricingRecommendation `json:"recommendations"`
	NextQuestions   []string                `json:"next_questions"`
	UpdatedContext  *ConversationContext    `json:"updated_context"`
}

// CustomerPreferences represents customer preferences
type CustomerPreferences struct {
	Priority        string   `json:"priority"` // "cost", "speed", "reliability"
	Budget          float64  `json:"budget"`
	TimeConstraints []string `json:"time_constraints"`
	SpecialNeeds    []string `json:"special_needs"`
}

// Intent represents a recognized intent from user input
type Intent struct {
	Type       string            `json:"type"`
	Entities   map[string]string `json:"entities"`
	Confidence float64           `json:"confidence"`
}

// extractIntent extracts intent and entities from natural language
func (ce *ConversationEngine) extractIntent(message string) *Intent {
	message = strings.ToLower(message)

	// Simple pattern matching for common intents
	patterns := map[string][]string{
		"compare_pricing": {
			"compare.*pricing",
			"what.*pricing.*options",
			"show.*me.*pricing",
			"pricing.*models",
		},
		"get_recommendation": {
			"what.*best.*pricing",
			"recommend.*pricing",
			"which.*pricing.*best",
			"best.*option",
		},
		"explore_options": {
			"explore.*pricing",
			"what.*options.*available",
			"show.*me.*options",
			"pricing.*choices",
		},
		"delivery_requirements": {
			"need.*deliver",
			"deliver.*to",
			"pickup.*from",
			"delivery.*count",
		},
		"customer_tier": {
			"gold.*tier",
			"silver.*tier",
			"bronze.*tier",
			"loyalty.*tier",
		},
		"volume_questions": {
			"how.*many.*deliver",
			"delivery.*count",
			"multiple.*deliver",
			"bulk.*order",
		},
	}

	// Find matching intent
	for intentType, intentPatterns := range patterns {
		for _, pattern := range intentPatterns {
			if matched, _ := regexp.MatchString(pattern, message); matched {
				entities := ce.extractEntities(message)
				return &Intent{
					Type:       intentType,
					Entities:   entities,
					Confidence: 0.8,
				}
			}
		}
	}

	// Default intent
	return &Intent{
		Type:       "general_inquiry",
		Entities:   make(map[string]string),
		Confidence: 0.5,
	}
}

// extractEntities extracts entities from the message
func (ce *ConversationEngine) extractEntities(message string) map[string]string {
	entities := make(map[string]string)

	// Extract delivery count - check multiple keywords
	deliveryKeywords := []string{"deliver", "delivery", "deliveries", "package", "packages", "shipment", "shipments"}
	for _, keyword := range deliveryKeywords {
		if count := ce.extractNumber(message, keyword); count > 0 {
			entities["delivery_count"] = fmt.Sprintf("%d", count)
			break
		}
	}

	// Extract customer tier
	if strings.Contains(message, "gold") {
		entities["customer_tier"] = "gold"
	} else if strings.Contains(message, "silver") {
		entities["customer_tier"] = "silver"
	} else if strings.Contains(message, "bronze") {
		entities["customer_tier"] = "bronze"
	}

	// Extract order frequency
	if freq := ce.extractNumber(message, "order"); freq > 0 {
		entities["order_frequency"] = fmt.Sprintf("%d", freq)
	}

	// Extract vehicle type
	if strings.Contains(message, "cargo") {
		entities["vehicle_type"] = "cargo_van"
	} else if strings.Contains(message, "sprinter") {
		entities["vehicle_type"] = "sprinter_van"
	}

	// Extract bulk order
	if strings.Contains(message, "bulk") {
		entities["is_bulk_order"] = "true"
	}

	return entities
}

// extractNumber extracts a number from text
func (ce *ConversationEngine) extractNumber(text, keyword string) int {
	// Simple number extraction - can be enhanced with NLP
	re := regexp.MustCompile(`(\d+).*` + keyword)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		if num, err := fmt.Sscanf(matches[1], "%d", new(int)); err == nil {
			return num
		}
	}
	return 0
}

// generateRecommendations creates pricing recommendations based on context
func (ce *ConversationEngine) generateRecommendations(context *ConversationContext) []PricingRecommendation {
	// Create pricing context from conversation context
	pricingContext := pricing.PricingContext{
		DeliveryCount:   context.CustomerProfile.OrderFrequency,
		CustomerTier:    context.CustomerProfile.Tier,
		OrderFrequency:  context.CustomerProfile.OrderFrequency,
		TotalOrderValue: context.CustomerProfile.AverageOrderValue,
		IsBulkOrder:     false, // Can be enhanced based on conversation
	}

	// Get available pricing models
	models := ce.pricingEngine.GetAvailableModels()
	var recommendations []PricingRecommendation

	for _, model := range models {
		// Create a sample estimate for comparison
		sampleEstimate := &dispatch.AvailableOrderOption{
			EstimatedOrderCost: 50.0, // Sample cost
		}

		// Compare pricing models
		comparison := ce.pricingEngine.ComparePricingModels(sampleEstimate, pricingContext)

		// Find the result for this model
		for _, result := range comparison.PricingModels {
			if result.Model == model.Model {
				recommendations = append(recommendations, PricingRecommendation{
					Model:          string(result.Model),
					Name:           result.Name,
					Savings:        result.Savings,
					SavingsPercent: result.DiscountPercent,
					Eligible:       result.Eligible,
					Reason:         result.Reason,
				})
				break
			}
		}
	}

	return recommendations
}

// generateResponse creates a natural language response
func (ce *ConversationEngine) generateResponse(intent *Intent, recommendations []PricingRecommendation, context *ConversationContext) string {
	switch intent.Type {
	case "compare_pricing":
		return ce.generatePricingComparisonResponse(recommendations)
	case "get_recommendation":
		return ce.generateRecommendationResponse(recommendations)
	case "explore_options":
		return ce.generateExplorationResponse(recommendations)
	case "delivery_requirements":
		return ce.generateDeliveryRequirementsResponse(intent, context)
	case "customer_tier":
		return ce.generateCustomerTierResponse(intent, context)
	default:
		return ce.generateDefaultResponse(recommendations)
	}
}

// generatePricingComparisonResponse creates a response for pricing comparison
func (ce *ConversationEngine) generatePricingComparisonResponse(recommendations []PricingRecommendation) string {
	if len(recommendations) == 0 {
		return "I'd be happy to help you compare pricing options! Could you tell me about your delivery needs?"
	}

	response := "Here are your pricing options:\n\n"

	for _, rec := range recommendations {
		if rec.Eligible {
			response += fmt.Sprintf("✅ **%s**: Save $%.2f (%.1f%%)\n", rec.Name, rec.Savings, rec.SavingsPercent)
		} else {
			response += fmt.Sprintf("❌ **%s**: %s\n", rec.Name, rec.Reason)
		}
	}

	// Find best option
	bestSavings := 0.0
	bestOption := ""
	for _, rec := range recommendations {
		if rec.Eligible && rec.Savings > bestSavings {
			bestSavings = rec.Savings
			bestOption = rec.Name
		}
	}

	if bestOption != "" {
		response += fmt.Sprintf("\n🏆 **Best Option**: %s with $%.2f savings!", bestOption, bestSavings)
	}

	return response
}

// generateRecommendationResponse creates a response for recommendations
func (ce *ConversationEngine) generateRecommendationResponse(recommendations []PricingRecommendation) string {
	if len(recommendations) == 0 {
		return "I'd love to recommend the best pricing for you! Could you tell me about your delivery needs and customer tier?"
	}

	// Find best recommendation
	bestSavings := 0.0
	bestOption := ""
	for _, rec := range recommendations {
		if rec.Eligible && rec.Savings > bestSavings {
			bestSavings = rec.Savings
			bestOption = rec.Name
		}
	}

	if bestOption != "" {
		return fmt.Sprintf("Based on your profile, I recommend **%s** for $%.2f savings (%.1f%% off)! This gives you the best value for your delivery needs.",
			bestOption, bestSavings, bestSavings/50.0*100) // Assuming $50 base cost
	}

	return "I need a bit more information to give you the best recommendation. Could you tell me about your delivery count and customer tier?"
}

// generateExplorationResponse creates a response for exploration
func (ce *ConversationEngine) generateExplorationResponse(recommendations []PricingRecommendation) string {
	response := "Let's explore your pricing options! Here's what's available:\n\n"

	// Group by eligibility
	eligible := []PricingRecommendation{}
	ineligible := []PricingRecommendation{}

	for _, rec := range recommendations {
		if rec.Eligible {
			eligible = append(eligible, rec)
		} else {
			ineligible = append(ineligible, rec)
		}
	}

	if len(eligible) > 0 {
		response += "**Available Options:**\n"
		for _, rec := range eligible {
			response += fmt.Sprintf("• %s: $%.2f savings\n", rec.Name, rec.Savings)
		}
	}

	if len(ineligible) > 0 {
		response += "\n**Potential Options:**\n"
		for _, rec := range ineligible {
			response += fmt.Sprintf("• %s: %s\n", rec.Name, rec.Reason)
		}
	}

	return response
}

// generateDeliveryRequirementsResponse creates a response for delivery requirements
func (ce *ConversationEngine) generateDeliveryRequirementsResponse(intent *Intent, context *ConversationContext) string {
	deliveryCount := intent.Entities["delivery_count"]
	if deliveryCount != "" {
		return fmt.Sprintf("Great! %s deliveries gives you access to our Multi-Delivery Discount (15%% off). Would you like to see all your pricing options?", deliveryCount)
	}

	return "I'd love to help you with your delivery needs! How many deliveries are you planning?"
}

// generateCustomerTierResponse creates a response for customer tier
func (ce *ConversationEngine) generateCustomerTierResponse(intent *Intent, context *ConversationContext) string {
	tier := intent.Entities["customer_tier"]
	if tier != "" {
		return fmt.Sprintf("Excellent! Your %s tier status gives you access to our Loyalty Discount (10%% off). Let me show you all available pricing options.", tier)
	}

	return "What's your customer tier? This helps me find the best pricing options for you."
}

// generateDefaultResponse creates a default response
func (ce *ConversationEngine) generateDefaultResponse(recommendations []PricingRecommendation) string {
	return "I'd be happy to help you with pricing! Could you tell me about your delivery needs? For example, how many deliveries do you need and what's your customer tier?"
}

// generateNextQuestions creates follow-up questions
func (ce *ConversationEngine) generateNextQuestions(context *ConversationContext) []string {
	questions := []string{}

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
