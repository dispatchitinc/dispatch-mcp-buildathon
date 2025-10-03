package conversation

import (
	"dispatch-mcp-server/internal/claude"
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/order"
	"dispatch-mcp-server/internal/pricing"
	"dispatch-mcp-server/internal/validation"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ConversationMessage represents a message in the conversation history
type ConversationMessage = claude.ConversationMessage

// ClaudeConversationEngine handles natural language pricing conversations using Claude AI
type ClaudeConversationEngine struct {
	claudeClient   *claude.Client
	pricingEngine  *pricing.PricingEngine
	contextManager *ContextManager
	useClaude      bool
}

// NewClaudeConversationEngine creates a new AI Hub-powered conversation engine with rule-based fallback
func NewClaudeConversationEngine() (*ClaudeConversationEngine, error) {
	// Try AI Hub first
	useAIHub := os.Getenv("USE_AI_HUB")
	if useAIHub == "true" {
		claudeClient, err := claude.NewClient()
		if err == nil {
			// AI Hub is available
			return &ClaudeConversationEngine{
				claudeClient:   claudeClient,
				pricingEngine:  pricing.NewPricingEngine(),
				contextManager: NewContextManager(),
				useClaude:      true,
			}, nil
		}
		// AI Hub failed, fall back to rule-based engine
		fmt.Printf("⚠️  AI Hub not available, using rule-based engine: %v\n", err)
	} else {
		// AI Hub not configured, use rule-based engine
		fmt.Printf("⚠️  AI Hub not configured (USE_AI_HUB=false), using rule-based engine\n")
	}

	// Fall back to rule-based engine
	return &ClaudeConversationEngine{
		claudeClient:   nil,
		pricingEngine:  pricing.NewPricingEngine(),
		contextManager: NewContextManager(),
		useClaude:      false,
	}, nil
}

// ProcessMessage processes a natural language message using Claude AI
func (ce *ClaudeConversationEngine) ProcessMessage(message string, context *ConversationContext) (*ConversationResponse, error) {
	return ce.ProcessMessageWithHistory(message, context, nil)
}

// ProcessMessageWithHistory processes a message with conversation history
func (ce *ClaudeConversationEngine) ProcessMessageWithHistory(message string, context *ConversationContext, history []ConversationMessage) (*ConversationResponse, error) {
	// If Claude is not available, fall back to rule-based processing
	if !ce.useClaude || ce.claudeClient == nil {
		return ce.processWithRules(message, context)
	}

	// Update context with new information from the message FIRST
	updatedContext := ce.updateContextFromMessage(message, context)

	// Check for validation errors and handle them
	validationErrorMsg := ce.handleValidationErrors(updatedContext)
	if validationErrorMsg != "" {
		// Return validation error response
		return &ConversationResponse{
			Message:         validationErrorMsg,
			Recommendations: []PricingRecommendation{},
			NextQuestions:   []string{},
			UpdatedContext:  updatedContext,
		}, nil
	}

	// Convert updated context to Claude's pricing context
	pricingContext := ce.convertToPricingContext(updatedContext)

	// Get Claude's response with updated context and conversation history
	claudeResponse, err := ce.claudeClient.CreatePricingAdvisorMessageWithHistory(message, pricingContext, history)
	if err != nil {
		// If Claude fails, fall back to rule-based processing
		return ce.processWithRules(message, context)
	}

	// Extract Claude's response text
	responseText := ""
	if len(claudeResponse.Content) > 0 {
		responseText = claudeResponse.Content[0].Text
	}

	// Generate pricing recommendations using our pricing engine with updated context
	recommendations := ce.generateRecommendations(updatedContext)

	// Context is now properly updated and passed to Claude

	return &ConversationResponse{
		Message:         responseText,
		Recommendations: recommendations,
		NextQuestions:   []string{}, // No generic next steps - AI handles conversation flow
		UpdatedContext:  updatedContext,
	}, nil
}

// processWithRules processes the message using our rule-based system
func (ce *ClaudeConversationEngine) processWithRules(message string, context *ConversationContext) (*ConversationResponse, error) {
	// Check for validation errors first
	validationErrorMsg := ce.handleValidationErrors(context)
	if validationErrorMsg != "" {
		// Return validation error response
		return &ConversationResponse{
			Message:         validationErrorMsg,
			Recommendations: []PricingRecommendation{},
			NextQuestions:   []string{},
			UpdatedContext:  context,
		}, nil
	}

	// Use the original rule-based engine
	originalEngine := NewConversationEngine()
	response, err := originalEngine.ProcessMessage(message, context)
	if err != nil {
		return response, err
	}

	// Remove generic next steps - let the AI handle conversation flow
	response.NextQuestions = []string{}
	return response, nil
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
					// Check previous word for number (e.g., "3 packages")
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

	// Also try to parse pickup and delivery information if not in formal order creation mode
	ce.parseOrderInformation(message, context)

	return context
}

// parseOrderInformation parses pickup and delivery information from messages
func (ce *ClaudeConversationEngine) parseOrderInformation(message string, context *ConversationContext) {
	// Initialize order creation if not started
	if !context.OrderCreation.InProgress {
		// Check if this looks like address information
		if ce.looksLikeAddressInfo(message) {
			context.OrderCreation.InProgress = true
			context.OrderCreation.Step = "pickup"
			context.OrderCreation.PickupInfo = &dispatch.CreateOrderPickupInfoInput{}
			context.OrderCreation.DropOffs = []dispatch.CreateOrderDropOffInfoInput{}
			context.OrderCreation.CompletedFields = []string{}
			context.OrderCreation.MissingFields = []string{}
		}
	}

	// If order creation is in progress, try to parse the information
	if context.OrderCreation.InProgress {
		// Check if this looks like address information
		if ce.looksLikeAddressInfo(message) {
			// If we don't have pickup info yet, this is pickup
			if context.OrderCreation.PickupInfo == nil || context.OrderCreation.PickupInfo.BusinessName == nil {
				ce.parsePickupInfo(message, context)
			} else {
				// We already have pickup info, so this must be delivery
				ce.parseDeliveryInfo(message, context)

				// After adding delivery info, validate service area
				if len(context.OrderCreation.DropOffs) > 0 {
					if err := ce.validateServiceArea(context); err != nil {
						context.OrderCreation.ValidationErrors = append(context.OrderCreation.ValidationErrors,
							fmt.Sprintf("Service area validation failed: %s", err.Error()))
					}
				}
			}
		}
	}
}

// looksLikeAddressInfo checks if the message contains address information
func (ce *ClaudeConversationEngine) looksLikeAddressInfo(message string) bool {
	// Look for patterns like "Business Name, Contact Name, Address, Phone"
	parts := strings.Split(message, ",")
	if len(parts) < 4 {
		return false
	}

	// Check for address indicators
	addressIndicators := []string{"drive", "street", "avenue", "road", "lane", "way", "blvd", "boulevard", "st", "ave", "rd", "ln", "pkwy", "parkway"}
	messageLower := strings.ToLower(message)

	for _, indicator := range addressIndicators {
		if strings.Contains(messageLower, indicator) {
			return true
		}
	}

	// Check for zip code pattern (5 digits)
	zipRegex := regexp.MustCompile(`\b\d{5}\b`)
	return zipRegex.MatchString(message)
}

// looksLikePickupInfo checks if the message contains pickup information
func (ce *ClaudeConversationEngine) looksLikePickupInfo(message string) bool {
	// Look for patterns like "Business Name, Contact Name, Address, Phone"
	// This is a simple heuristic - could be enhanced with better NLP
	parts := strings.Split(message, ",")
	return len(parts) >= 4 && (strings.Contains(strings.ToLower(message), "drive") || strings.Contains(strings.ToLower(message), "street") || strings.Contains(strings.ToLower(message), "avenue"))
}

// looksLikeDeliveryInfo checks if the message contains delivery information
func (ce *ClaudeConversationEngine) looksLikeDeliveryInfo(message string) bool {
	// Look for patterns like "Business Name, Contact Name, Address, Phone"
	parts := strings.Split(message, ",")
	return len(parts) >= 4 && (strings.Contains(strings.ToLower(message), "drive") || strings.Contains(strings.ToLower(message), "street") || strings.Contains(strings.ToLower(message), "avenue"))
}

// parsePickupInfo parses pickup information from the message
func (ce *ClaudeConversationEngine) parsePickupInfo(message string, context *ConversationContext) {
	parts := strings.Split(message, ",")
	if len(parts) >= 4 {
		// Parse: "Business Name, Contact Name, Address, City, State, Zip, Phone"
		businessName := strings.TrimSpace(parts[0])
		contactName := strings.TrimSpace(parts[1])
		address := strings.TrimSpace(parts[2])
		city := strings.TrimSpace(parts[3])

		// Handle optional fields with bounds checking
		state := ""
		if len(parts) > 4 {
			state = strings.TrimSpace(parts[4])
		}

		zip := ""
		if len(parts) > 5 {
			zip = strings.TrimSpace(parts[5])
		}

		phone := ""
		if len(parts) > 6 {
			phone = strings.TrimSpace(parts[6])
		}

		// Normalize state code (handle full state names)
		state = ce.normalizeStateCode(state)

		// Validate address using existing validator
		validator := validation.NewValidator()
		addressMap := map[string]interface{}{
			"street":   address,
			"city":     city,
			"state":    state,
			"zip_code": zip,
			"country":  "US",
		}

		validationResult := validator.ValidateAddress(addressMap)
		if !validationResult.Valid {
			// Store validation errors in context for AI to handle
			context.OrderCreation.ValidationErrors = append(context.OrderCreation.ValidationErrors,
				fmt.Sprintf("Pickup address validation failed: %s", validationResult.Message))
			return
		}

		context.OrderCreation.PickupInfo = &dispatch.CreateOrderPickupInfoInput{
			BusinessName:       &businessName,
			ContactName:        &contactName,
			ContactPhoneNumber: &phone,
			Location: &dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  address,
					City:    city,
					State:   state,
					ZipCode: zip,
					Country: "US",
				},
			},
		}

		context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, "pickup_business", "pickup_contact", "pickup_address", "pickup_phone")
	}
}

// parseDeliveryInfo parses delivery information from the message
func (ce *ClaudeConversationEngine) parseDeliveryInfo(message string, context *ConversationContext) {
	parts := strings.Split(message, ",")
	if len(parts) >= 4 {
		// Parse: "Business Name, Contact Name, Address, City, State, Zip, Phone"
		businessName := strings.TrimSpace(parts[0])
		contactName := strings.TrimSpace(parts[1])
		address := strings.TrimSpace(parts[2])
		city := strings.TrimSpace(parts[3])

		// Handle optional fields with bounds checking
		state := ""
		if len(parts) > 4 {
			state = strings.TrimSpace(parts[4])
		}

		zip := ""
		if len(parts) > 5 {
			zip = strings.TrimSpace(parts[5])
		}

		phone := ""
		if len(parts) > 6 {
			phone = strings.TrimSpace(parts[6])
		}

		// Normalize state code (handle full state names)
		state = ce.normalizeStateCode(state)

		// Validate address using existing validator
		validator := validation.NewValidator()
		addressMap := map[string]interface{}{
			"street":   address,
			"city":     city,
			"state":    state,
			"zip_code": zip,
			"country":  "US",
		}

		validationResult := validator.ValidateAddress(addressMap)
		if !validationResult.Valid {
			// Store validation errors in context for AI to handle
			context.OrderCreation.ValidationErrors = append(context.OrderCreation.ValidationErrors,
				fmt.Sprintf("Delivery address validation failed: %s", validationResult.Message))
			return
		}

		dropOff := dispatch.CreateOrderDropOffInfoInput{
			BusinessName:       &businessName,
			ContactName:        &contactName,
			ContactPhoneNumber: &phone,
			Location: &dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  address,
					City:    city,
					State:   state,
					ZipCode: zip,
					Country: "US",
				},
			},
		}

		context.OrderCreation.DropOffs = append(context.OrderCreation.DropOffs, dropOff)
		context.OrderCreation.CompletedFields = append(context.OrderCreation.CompletedFields, fmt.Sprintf("delivery_%d", len(context.OrderCreation.DropOffs)))
	}
}

// handleValidationErrors processes validation errors and returns user-friendly messages
func (ce *ClaudeConversationEngine) handleValidationErrors(context *ConversationContext) string {
	if len(context.OrderCreation.ValidationErrors) == 0 {
		return ""
	}

	var errorMessages []string
	for _, error := range context.OrderCreation.ValidationErrors {
		if strings.Contains(error, "zip code") {
			errorMessages = append(errorMessages, "Please provide a valid zip code (5 digits or 5+4 format like 12345 or 12345-6789)")
		} else if strings.Contains(error, "state") {
			errorMessages = append(errorMessages, "Please provide a valid 2-letter state code (e.g., CA, NY, TX)")
		} else if strings.Contains(error, "required") {
			errorMessages = append(errorMessages, "Please provide a complete address with street, city, state, and zip code")
		} else if strings.Contains(error, "Service area validation failed") {
			errorMessages = append(errorMessages, "Sorry, we don't currently deliver to this location. Please try a different address or contact support for service area information.")
		} else if strings.Contains(error, "no delivery options available") {
			errorMessages = append(errorMessages, "No delivery options are available for this location. Please try a different address.")
		} else {
			errorMessages = append(errorMessages, "Please check your address format and try again")
		}
	}

	// Clear validation errors after processing
	context.OrderCreation.ValidationErrors = []string{}

	return strings.Join(errorMessages, "\n")
}

// validateServiceArea checks if locations are in Dispatch service area using CreateEstimate
func (ce *ClaudeConversationEngine) validateServiceArea(context *ConversationContext) error {
	// Only validate if we have both pickup and delivery info
	if context.OrderCreation.PickupInfo == nil || len(context.OrderCreation.DropOffs) == 0 {
		return nil // Not ready for validation
	}

	// Convert to CreateEstimateInput format
	pickupInfo := dispatch.PickupInfoInput{
		BusinessName: *context.OrderCreation.PickupInfo.BusinessName,
		Location:     *context.OrderCreation.PickupInfo.Location,
	}

	dropOffs := make([]dispatch.DropOffInfoInput, len(context.OrderCreation.DropOffs))
	for i, dropOff := range context.OrderCreation.DropOffs {
		dropOffs[i] = dispatch.DropOffInfoInput{
			BusinessName: *dropOff.BusinessName,
			Location:     *dropOff.Location,
		}
	}

	// Create test estimate to check service area
	input := dispatch.CreateEstimateInput{
		PickupInfo:  pickupInfo,
		DropOffs:    dropOffs,
		VehicleType: "cargo_van", // Default vehicle type for validation
	}

	// Create Dispatch client
	client, err := dispatch.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Dispatch client: %v", err)
	}

	// Try to create estimate - if this fails, location might be out of service area
	response, err := client.CreateEstimate(input)
	if err != nil {
		return fmt.Errorf("service area validation failed: %v", err)
	}

	// Check if we got valid delivery options
	if len(response.Data.CreateEstimate.Estimate.AvailableOrderOptions) == 0 {
		return fmt.Errorf("no delivery options available for this location")
	}

	return nil
}

// normalizeStateCode converts full state names to 2-letter codes
func (ce *ClaudeConversationEngine) normalizeStateCode(state string) string {
	stateMap := map[string]string{
		"california": "CA", "calif": "CA", "ca": "CA",
		"new york": "NY", "ny": "NY",
		"texas": "TX", "tx": "TX",
		"florida": "FL", "fl": "FL",
		"illinois": "IL", "il": "IL",
		"pennsylvania": "PA", "pa": "PA",
		"ohio": "OH", "oh": "OH",
		"georgia": "GA", "ga": "GA",
		"north carolina": "NC", "nc": "NC",
		"michigan": "MI", "mi": "MI",
		"new jersey": "NJ", "nj": "NJ",
		"virginia": "VA", "va": "VA",
		"washington": "WA", "wa": "WA",
		"arizona": "AZ", "az": "AZ",
		"massachusetts": "MA", "ma": "MA",
		"tennessee": "TN", "tn": "TN",
		"indiana": "IN", "in": "IN",
		"missouri": "MO", "mo": "MO",
		"maryland": "MD", "md": "MD",
		"wisconsin": "WI", "wi": "WI",
		"colorado": "CO", "co": "CO",
		"minnesota": "MN", "mn": "MN",
		"alabama": "AL", "al": "AL",
		"louisiana": "LA", "la": "LA",
		"kentucky": "KY", "ky": "KY",
		"oregon": "OR", "or": "OR",
		"oklahoma": "OK", "ok": "OK",
		"connecticut": "CT", "ct": "CT",
		"utah": "UT", "ut": "UT",
		"iowa": "IA", "ia": "IA",
		"nevada": "NV", "nv": "NV",
		"arkansas": "AR", "ar": "AR",
		"mississippi": "MS", "ms": "MS",
		"kansas": "KS", "ks": "KS",
		"new mexico": "NM", "nm": "NM",
		"nebraska": "NE", "ne": "NE",
		"west virginia": "WV", "wv": "WV",
		"idaho": "ID", "id": "ID",
		"hawaii": "HI", "hi": "HI",
		"new hampshire": "NH", "nh": "NH",
		"maine": "ME", "me": "ME",
		"montana": "MT", "mt": "MT",
		"rhode island": "RI", "ri": "RI",
		"delaware": "DE", "de": "DE",
		"south dakota": "SD", "sd": "SD",
		"north dakota": "ND", "nd": "ND",
		"alaska": "AK", "ak": "AK",
		"vermont": "VT", "vt": "VT",
		"wyoming": "WY", "wy": "WY",
		"district of columbia": "DC", "washington dc": "DC", "dc": "DC",
	}

	stateLower := strings.ToLower(strings.TrimSpace(state))
	if code, exists := stateMap[stateLower]; exists {
		return code
	}

	// If already a 2-letter code, return as-is
	if len(state) == 2 {
		return strings.ToUpper(state)
	}

	// Return original if no match found
	return state
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

// determineNextStep determines the next step in the order creation process
func (ce *ClaudeConversationEngine) determineNextStep(context *ConversationContext) string {
	if !context.OrderCreation.InProgress {
		return "multi_stop" // Start with multi-stop selection
	}

	// Check what information we have and what's missing
	if context.OrderCreation.Step == "multi_stop" {
		return "pickup"
	}

	if context.OrderCreation.PickupInfo == nil {
		return "pickup"
	}

	if len(context.OrderCreation.DropOffs) == 0 {
		return "drop_off"
	}

	if context.OrderCreation.VehicleType == nil {
		return "vehicle"
	}

	if len(context.OrderCreation.Capabilities) == 0 {
		return "add_ons"
	}

	if context.OrderCreation.SchedulingInfo == nil {
		return "delivery"
	}

	return "review"
}

// getStepPrompt returns the appropriate prompt for the current step
func (ce *ClaudeConversationEngine) getStepPrompt(step string, context *ConversationContext) string {
	switch step {
	case "multi_stop":
		return "I'll help you create a delivery order. First, let me know: do you need to deliver to one location or multiple locations?"

	case "pickup":
		return "Great! Now I need your pickup information. What's the business name and address where we'll be picking up the package?"

	case "drop_off":
		if context.OrderCreation.MultiStop {
			return "Perfect! Now I need your delivery information. Where should we deliver this package? (You can add multiple delivery locations)"
		}
		return "Perfect! Now I need your delivery information. Where should we deliver this package?"

	case "vehicle":
		return "What type of vehicle do you need for this delivery? (cargo van, pickup truck, etc.)"

	case "add_ons":
		return "Do you need any special services or capabilities for this delivery? (temperature control, white glove service, etc.)"

	case "delivery":
		return "When would you like this delivered? I need pickup time and delivery time."

	case "review":
		return "Let me review your order details before we create it..."

	default:
		return "I'm ready to help you create a delivery order. What do you need?"
	}
}

// parseMultiStopInfo parses whether user wants single or multi-stop delivery
func (ce *ClaudeConversationEngine) parseMultiStopInfo(message string, context *ConversationContext) {
	messageLower := strings.ToLower(message)

	if strings.Contains(messageLower, "multiple") ||
		strings.Contains(messageLower, "several") ||
		strings.Contains(messageLower, "many") ||
		strings.Contains(messageLower, "more than one") {
		context.OrderCreation.MultiStop = true
	} else {
		context.OrderCreation.MultiStop = false
	}
}

// parseVehicleInfo parses vehicle type from message
func (ce *ClaudeConversationEngine) parseVehicleInfo(message string, context *ConversationContext) {
	messageLower := strings.ToLower(message)

	vehicleTypes := map[string]string{
		"cargo van":    "cargo_van",
		"pickup truck": "pickup_truck",
		"box truck":    "box_truck",
		"van":          "cargo_van",
		"truck":        "pickup_truck",
		"cargo":        "cargo_van",
		"pickup":       "pickup_truck",
	}

	for keyword, vehicleType := range vehicleTypes {
		if strings.Contains(messageLower, keyword) {
			context.OrderCreation.VehicleType = &VehicleTypeInfo{
				VehicleTypeID:   vehicleType,
				VehicleTypeName: keyword,
			}
			break
		}
	}
}

// parseCapabilitiesInfo parses special requirements and capabilities
func (ce *ClaudeConversationEngine) parseCapabilitiesInfo(message string, context *ConversationContext) {
	messageLower := strings.ToLower(message)
	capabilities := []string{}

	if strings.Contains(messageLower, "temperature") {
		capabilities = append(capabilities, "temperature_control")
	}
	if strings.Contains(messageLower, "white glove") {
		capabilities = append(capabilities, "white_glove_service")
	}
	if strings.Contains(messageLower, "fragile") {
		capabilities = append(capabilities, "fragile_handling")
	}
	if strings.Contains(messageLower, "signature") {
		capabilities = append(capabilities, "signature_required")
	}
	if strings.Contains(messageLower, "assistance") {
		capabilities = append(capabilities, "unloading_assistance")
	}

	context.OrderCreation.Capabilities = capabilities
}

// parseSchedulingInfo parses scheduling information
func (ce *ClaudeConversationEngine) parseSchedulingInfo(message string, context *ConversationContext) {
	scheduling := &SchedulingInfo{}

	// Look for time patterns (simple regex for now)
	timePattern := regexp.MustCompile(`(\d{1,2}):?(\d{2})?\s*(am|pm|AM|PM)?`)
	matches := timePattern.FindAllString(message, -1)

	if len(matches) > 0 {
		scheduling.PickupTime = matches[0]
		if len(matches) > 1 {
			scheduling.DeliveryTime = matches[1]
		}
	}

	// Look for date patterns
	datePattern := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4})`)
	dateMatches := datePattern.FindAllString(message, -1)

	if len(dateMatches) > 0 {
		scheduling.PickupDate = dateMatches[0]
		if len(dateMatches) > 1 {
			scheduling.DeliveryDate = dateMatches[1]
		}
	}

	context.OrderCreation.SchedulingInfo = scheduling
}

// parseStepInformation parses information based on the current step
func (ce *ClaudeConversationEngine) parseStepInformation(message string, context *ConversationContext) {
	if !context.OrderCreation.InProgress {
		return
	}

	switch context.OrderCreation.Step {
	case "multi_stop":
		ce.parseMultiStopInfo(message, context)
	case "pickup":
		ce.parseOrderInformation(message, context) // Existing pickup/delivery parsing
	case "drop_off":
		ce.parseOrderInformation(message, context) // Existing pickup/delivery parsing
	case "vehicle":
		ce.parseVehicleInfo(message, context)
	case "add_ons":
		ce.parseCapabilitiesInfo(message, context)
	case "delivery":
		ce.parseSchedulingInfo(message, context)
	case "review":
		// Review step - no parsing needed, just confirmation
	}
}

// handleReviewStep handles the review step and order creation
func (ce *ClaudeConversationEngine) handleReviewStep(message string, context *ConversationContext) string {
	// Generate order summary
	summary := ce.generateOrderSummary(context)

	if strings.Contains(strings.ToLower(message), "yes") ||
		strings.Contains(strings.ToLower(message), "create") ||
		strings.Contains(strings.ToLower(message), "confirm") {

		// Convert context to order creation input
		// Convert VehicleTypeInfo to order.VehicleTypeInfo
		var orderVehicleType *order.VehicleTypeInfo
		if context.OrderCreation.VehicleType != nil {
			orderVehicleType = &order.VehicleTypeInfo{
				VehicleTypeID:   context.OrderCreation.VehicleType.VehicleTypeID,
				VehicleTypeName: context.OrderCreation.VehicleType.VehicleTypeName,
				CustomTypes:     context.OrderCreation.VehicleType.CustomTypes,
			}
		}

		// Convert SchedulingInfo to order.SchedulingInput
		var orderScheduling *order.SchedulingInput
		if context.OrderCreation.SchedulingInfo != nil {
			orderScheduling = &order.SchedulingInput{
				PickupTime:   context.OrderCreation.SchedulingInfo.PickupTime,
				DeliveryTime: context.OrderCreation.SchedulingInfo.DeliveryTime,
				PickupDate:   context.OrderCreation.SchedulingInfo.PickupDate,
				DeliveryDate: context.OrderCreation.SchedulingInfo.DeliveryDate,
				TimeZone:     context.OrderCreation.SchedulingInfo.TimeZone,
			}
		}

		orderInput := order.ConvertFromDispatchTypes(
			context.OrderCreation.PickupInfo,
			context.OrderCreation.DropOffs,
			orderVehicleType,
			context.OrderCreation.Capabilities,
			orderScheduling,
		)

		// Create the order using GraphQL
		orderCreator := order.NewOrderCreator(os.Getenv("GRAPHQL_ENDPOINT"))
		result, err := orderCreator.CreateOrder(orderInput)

		if err != nil {
			return fmt.Sprintf("I encountered an error creating your order: %v. Please try again or contact support.", err)
		}

		return fmt.Sprintf("🎉 Order created successfully!\n\nOrder ID: %s\nDruid: %s\nTotal Price: $%.2f\n\nYou'll receive a confirmation email shortly.",
			result.Order.ID, result.Order.Druid, result.Order.Pricing.TotalPrice)
	}

	return summary + "\n\nShould I create this order for you?"
}

// generateOrderSummary generates a summary of the order for review
func (ce *ClaudeConversationEngine) generateOrderSummary(context *ConversationContext) string {
	var summary strings.Builder

	summary.WriteString("📋 **Order Summary**\n\n")

	// Pickup information
	if context.OrderCreation.PickupInfo != nil {
		summary.WriteString("**Pickup Location:**\n")
		summary.WriteString(fmt.Sprintf("- Business: %s\n", *context.OrderCreation.PickupInfo.BusinessName))
		summary.WriteString(fmt.Sprintf("- Contact: %s (%s)\n", *context.OrderCreation.PickupInfo.ContactName, *context.OrderCreation.PickupInfo.ContactPhoneNumber))
		summary.WriteString(fmt.Sprintf("- Address: %s, %s, %s %s\n",
			context.OrderCreation.PickupInfo.Location.Address.Street,
			context.OrderCreation.PickupInfo.Location.Address.City,
			context.OrderCreation.PickupInfo.Location.Address.State,
			context.OrderCreation.PickupInfo.Location.Address.ZipCode))
	}

	// Delivery information
	if len(context.OrderCreation.DropOffs) > 0 {
		summary.WriteString("\n**Delivery Locations:**\n")
		for i, dropOff := range context.OrderCreation.DropOffs {
			summary.WriteString(fmt.Sprintf("%d. %s\n", i+1, *dropOff.BusinessName))
			summary.WriteString(fmt.Sprintf("   Contact: %s (%s)\n", *dropOff.ContactName, *dropOff.ContactPhoneNumber))
			summary.WriteString(fmt.Sprintf("   Address: %s, %s, %s %s\n",
				dropOff.Location.Address.Street,
				dropOff.Location.Address.City,
				dropOff.Location.Address.State,
				dropOff.Location.Address.ZipCode))
		}
	}

	// Vehicle type
	if context.OrderCreation.VehicleType != nil {
		summary.WriteString(fmt.Sprintf("\n**Vehicle Type:** %s\n", context.OrderCreation.VehicleType.VehicleTypeName))
	}

	// Capabilities
	if len(context.OrderCreation.Capabilities) > 0 {
		summary.WriteString(fmt.Sprintf("\n**Special Services:** %s\n", strings.Join(context.OrderCreation.Capabilities, ", ")))
	}

	// Scheduling
	if context.OrderCreation.SchedulingInfo != nil {
		summary.WriteString(fmt.Sprintf("\n**Schedule:**\n"))
		summary.WriteString(fmt.Sprintf("- Pickup: %s %s\n", context.OrderCreation.SchedulingInfo.PickupDate, context.OrderCreation.SchedulingInfo.PickupTime))
		summary.WriteString(fmt.Sprintf("- Delivery: %s %s\n", context.OrderCreation.SchedulingInfo.DeliveryDate, context.OrderCreation.SchedulingInfo.DeliveryTime))
	}

	return summary.String()
}
