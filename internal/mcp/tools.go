package mcp

import (
	"context"
	"dispatch-mcp-server/internal/conversation"
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/pricing"
	"dispatch-mcp-server/internal/validation"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) createEstimateTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Cast arguments to the correct type
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Parse and validate pickup_info
	pickupInfoRaw, ok := arguments["pickup_info"].(string)
	if !ok {
		return mcp.NewToolResultError("pickup_info is required and must be a string"), nil
	}

	// Validate JSON format
	if result := validator.ValidateJSONString(pickupInfoRaw, "pickup_info"); !result.Valid {
		errorMsg := fmt.Sprintf("pickup_info validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	var pickupInfo dispatch.PickupInfoInput
	if err := json.Unmarshal([]byte(pickupInfoRaw), &pickupInfo); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse pickup_info: %v", err)), nil
	}

	// Parse and validate drop_offs
	dropOffsRaw, ok := arguments["drop_offs"].(string)
	if !ok {
		return mcp.NewToolResultError("drop_offs is required and must be a string"), nil
	}

	// Validate JSON format
	if result := validator.ValidateJSONString(dropOffsRaw, "drop_offs"); !result.Valid {
		errorMsg := fmt.Sprintf("drop_offs validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	var dropOffs []dispatch.DropOffInfoInput
	if err := json.Unmarshal([]byte(dropOffsRaw), &dropOffs); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse drop_offs: %v", err)), nil
	}

	// Validate vehicle_type
	vehicleType := getStringArg(arguments, "vehicle_type")
	if result := validator.ValidateVehicleType(vehicleType); !result.Valid {
		errorMsg := fmt.Sprintf("vehicle_type validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	// Build input
	input := dispatch.CreateEstimateInput{
		PickupInfo:  pickupInfo,
		DropOffs:    dropOffs,
		VehicleType: vehicleType,
	}

	// Optional fields
	if addOns, ok := arguments["add_ons"].(string); ok && addOns != "" {
		var addOnsList []string
		if err := json.Unmarshal([]byte(addOns), &addOnsList); err == nil {
			input.AddOns = addOnsList
		}
	}

	if dedicatedVehicle, ok := arguments["dedicated_vehicle"].(string); ok && dedicatedVehicle != "" {
		if dedicatedVehicle == "true" {
			input.DedicatedVehicle = &[]bool{true}[0]
		} else if dedicatedVehicle == "false" {
			input.DedicatedVehicle = &[]bool{false}[0]
		}
	}

	if orgDruid, ok := arguments["organization_druid"].(string); ok && orgDruid != "" {
		input.OrganizationDruid = &orgDruid
	}

	// Call API
	response, err := s.dispatchClient.CreateEstimate(input)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create estimate: %v", err)), nil
	}

	// Format response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (s *MCPServer) createOrderTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Cast arguments to the correct type
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	// Parse delivery_info
	deliveryInfoRaw, ok := arguments["delivery_info"].(string)
	if !ok {
		return mcp.NewToolResultError("delivery_info is required and must be a string"), nil
	}

	var deliveryInfo dispatch.DeliveryInfoInput
	if err := json.Unmarshal([]byte(deliveryInfoRaw), &deliveryInfo); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse delivery_info: %v", err)), nil
	}

	// Parse pickup_info
	pickupInfoRaw, ok := arguments["pickup_info"].(string)
	if !ok {
		return mcp.NewToolResultError("pickup_info is required and must be a string"), nil
	}

	var pickupInfo dispatch.CreateOrderPickupInfoInput
	if err := json.Unmarshal([]byte(pickupInfoRaw), &pickupInfo); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse pickup_info: %v", err)), nil
	}

	// Parse drop_offs
	dropOffsRaw, ok := arguments["drop_offs"].(string)
	if !ok {
		return mcp.NewToolResultError("drop_offs is required and must be a string"), nil
	}

	var dropOffs []dispatch.CreateOrderDropOffInfoInput
	if err := json.Unmarshal([]byte(dropOffsRaw), &dropOffs); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse drop_offs: %v", err)), nil
	}

	// Build input
	input := dispatch.CreateOrderInput{
		DeliveryInfo: deliveryInfo,
		PickupInfo:   pickupInfo,
		DropOffs:     dropOffs,
	}

	// Optional fields
	if tags, ok := arguments["tags"].(string); ok && tags != "" {
		var tagsList []dispatch.TagInput
		if err := json.Unmarshal([]byte(tags), &tagsList); err == nil {
			input.Tags = tagsList
		}
	}

	// Call API
	response, err := s.dispatchClient.CreateOrder(input)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create order: %v", err)), nil
	}

	// Format response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (s *MCPServer) comparePricingModelsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Cast arguments to the correct type
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Parse and validate original_estimate
	originalEstimateRaw, ok := arguments["original_estimate"].(string)
	if !ok {
		return mcp.NewToolResultError("original_estimate is required and must be a string"), nil
	}

	// Validate JSON format
	if result := validator.ValidateJSONString(originalEstimateRaw, "original_estimate"); !result.Valid {
		errorMsg := fmt.Sprintf("original_estimate validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	var originalEstimate dispatch.AvailableOrderOption
	if err := json.Unmarshal([]byte(originalEstimateRaw), &originalEstimate); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse original_estimate: %v", err)), nil
	}

	// Parse context parameters with defaults
	context := pricing.PricingContext{
		DeliveryCount:   1,
		CustomerTier:    "bronze",
		OrderFrequency:  1,
		TotalOrderValue: originalEstimate.EstimatedOrderCost,
		IsBulkOrder:     false,
	}

	// Parse and validate delivery_count
	if deliveryCountStr, ok := arguments["delivery_count"].(string); ok && deliveryCountStr != "" {
		if result := validator.ValidateNumericString(deliveryCountStr, "delivery_count", 1, 100); !result.Valid {
			errorMsg := fmt.Sprintf("delivery_count validation failed: %s", result.Message)
			if len(result.Errors) > 0 {
				errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
			}
			return mcp.NewToolResultError(errorMsg), nil
		}
		if count, err := strconv.Atoi(deliveryCountStr); err == nil {
			context.DeliveryCount = count
		}
	}

	// Parse and validate customer_tier
	if customerTier, ok := arguments["customer_tier"].(string); ok && customerTier != "" {
		if result := validator.ValidateCustomerTier(customerTier); !result.Valid {
			errorMsg := fmt.Sprintf("customer_tier validation failed: %s", result.Message)
			if len(result.Errors) > 0 {
				errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
			}
			return mcp.NewToolResultError(errorMsg), nil
		}
		context.CustomerTier = customerTier
	}

	// Parse and validate order_frequency
	if orderFreqStr, ok := arguments["order_frequency"].(string); ok && orderFreqStr != "" {
		if result := validator.ValidateNumericString(orderFreqStr, "order_frequency", 1, 100); !result.Valid {
			errorMsg := fmt.Sprintf("order_frequency validation failed: %s", result.Message)
			if len(result.Errors) > 0 {
				errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
			}
			return mcp.NewToolResultError(errorMsg), nil
		}
		if freq, err := strconv.Atoi(orderFreqStr); err == nil {
			context.OrderFrequency = freq
		}
	}

	// Parse and validate total_order_value
	if totalValueStr, ok := arguments["total_order_value"].(string); ok && totalValueStr != "" {
		if value, err := strconv.ParseFloat(totalValueStr, 64); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("total_order_value must be a valid number: %v", err)), nil
		} else {
			context.TotalOrderValue = value
		}
	}

	// Parse and validate is_bulk_order
	if isBulkStr, ok := arguments["is_bulk_order"].(string); ok && isBulkStr != "" {
		if result := validator.ValidateBooleanString(isBulkStr, "is_bulk_order"); !result.Valid {
			errorMsg := fmt.Sprintf("is_bulk_order validation failed: %s", result.Message)
			if len(result.Errors) > 0 {
				errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
			}
			return mcp.NewToolResultError(errorMsg), nil
		}
		context.IsBulkOrder = (isBulkStr == "true")
	}

	// Create pricing engine and compare models
	engine := pricing.NewPricingEngine()
	comparison := engine.ComparePricingModels(&originalEstimate, context)

	// Format response
	responseJSON, _ := json.MarshalIndent(comparison, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (s *MCPServer) selectDeliveryOptionTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Cast arguments to the correct type
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Parse and validate estimate response
	estimateResponseRaw, ok := arguments["estimate_response"].(string)
	if !ok {
		return mcp.NewToolResultError("estimate_response is required and must be a string"), nil
	}

	// Validate JSON format
	if result := validator.ValidateJSONString(estimateResponseRaw, "estimate_response"); !result.Valid {
		errorMsg := fmt.Sprintf("estimate_response validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	var estimateResponse dispatch.CreateEstimateResponse
	if err := json.Unmarshal([]byte(estimateResponseRaw), &estimateResponse); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse estimate_response: %v", err)), nil
	}

	// Parse and validate delivery scenario
	scenario, ok := arguments["delivery_scenario"].(string)
	if !ok {
		return mcp.NewToolResultError("delivery_scenario is required and must be a string"), nil
	}

	// Validate delivery scenario
	if result := validator.ValidateDeliveryScenario(scenario); !result.Valid {
		errorMsg := fmt.Sprintf("delivery_scenario validation failed: %s", result.Message)
		if len(result.Errors) > 0 {
			errorMsg += fmt.Sprintf(" - %s", result.Errors[0].Message)
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	// Get available options
	options := estimateResponse.Data.CreateEstimate.Estimate.AvailableOrderOptions
	if len(options) == 0 {
		return mcp.NewToolResultError("no delivery options available"), nil
	}

	var selectedOption dispatch.AvailableOrderOption
	var scenarioDescription string

	switch scenario {
	case "fastest", "asap", "urgent":
		// First option = fastest delivery, most expensive
		selectedOption = options[0]
		scenarioDescription = "Fastest delivery (most expensive)"
	case "cheapest", "economy", "sometime_today":
		// Last option = slowest delivery, cheapest
		selectedOption = options[len(options)-1]
		scenarioDescription = "Cheapest delivery (slowest)"
	default:
		return mcp.NewToolResultError("delivery_scenario must be 'fastest' or 'cheapest'"), nil
	}

	// Create response with selected option and context
	response := map[string]interface{}{
		"selected_option": selectedOption,
		"scenario":        scenario,
		"description":     scenarioDescription,
		"total_options":   len(options),
		"all_options":     options,
	}

	// Format response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (s *MCPServer) conversationalPricingAdvisorTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Cast arguments to the correct type
	arguments, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	// Parse user message
	userMessage, ok := arguments["user_message"].(string)
	if !ok {
		return mcp.NewToolResultError("user_message is required and must be a string"), nil
	}

	// Parse conversation context (optional)
	var conversationContext *conversation.ConversationContext
	if contextRaw, ok := arguments["conversation_context"].(string); ok && contextRaw != "" {
		if err := json.Unmarshal([]byte(contextRaw), &conversationContext); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse conversation_context: %v", err)), nil
		}
	}

	// Parse customer profile (optional)
	if profileRaw, ok := arguments["customer_profile"].(string); ok && profileRaw != "" {
		if conversationContext == nil {
			conversationContext = &conversation.ConversationContext{}
		}
		if err := json.Unmarshal([]byte(profileRaw), &conversationContext.CustomerProfile); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse customer_profile: %v", err)), nil
		}
	}

	// Process the message through the conversation engine
	response, err := s.conversationEngine.ProcessMessage(userMessage, conversationContext)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to process message: %v", err)), nil
	}

	// Format response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func getStringArg(arguments map[string]interface{}, key string) string {
	if value, ok := arguments[key].(string); ok {
		return value
	}
	return ""
}
