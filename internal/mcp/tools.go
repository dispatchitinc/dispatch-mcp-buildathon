package mcp

import (
	"context"
	"dispatch-mcp-server/internal/conversation"
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/pricing"
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

	// Parse pickup_info
	pickupInfoRaw, ok := arguments["pickup_info"].(string)
	if !ok {
		return mcp.NewToolResultError("pickup_info is required and must be a string"), nil
	}

	var pickupInfo dispatch.PickupInfoInput
	if err := json.Unmarshal([]byte(pickupInfoRaw), &pickupInfo); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse pickup_info: %v", err)), nil
	}

	// Parse drop_offs
	dropOffsRaw, ok := arguments["drop_offs"].(string)
	if !ok {
		return mcp.NewToolResultError("drop_offs is required and must be a string"), nil
	}

	var dropOffs []dispatch.DropOffInfoInput
	if err := json.Unmarshal([]byte(dropOffsRaw), &dropOffs); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse drop_offs: %v", err)), nil
	}

	// Build input
	input := dispatch.CreateEstimateInput{
		PickupInfo:  pickupInfo,
		DropOffs:    dropOffs,
		VehicleType: getStringArg(arguments, "vehicle_type"),
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

	// Parse original_estimate
	originalEstimateRaw, ok := arguments["original_estimate"].(string)
	if !ok {
		return mcp.NewToolResultError("original_estimate is required and must be a string"), nil
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

	// Parse delivery_count
	if deliveryCountStr, ok := arguments["delivery_count"].(string); ok && deliveryCountStr != "" {
		if count, err := strconv.Atoi(deliveryCountStr); err == nil {
			context.DeliveryCount = count
		}
	}

	// Parse customer_tier
	if customerTier, ok := arguments["customer_tier"].(string); ok && customerTier != "" {
		context.CustomerTier = customerTier
	}

	// Parse order_frequency
	if orderFreqStr, ok := arguments["order_frequency"].(string); ok && orderFreqStr != "" {
		if freq, err := strconv.Atoi(orderFreqStr); err == nil {
			context.OrderFrequency = freq
		}
	}

	// Parse total_order_value
	if totalValueStr, ok := arguments["total_order_value"].(string); ok && totalValueStr != "" {
		if value, err := strconv.ParseFloat(totalValueStr, 64); err == nil {
			context.TotalOrderValue = value
		}
	}

	// Parse is_bulk_order
	if isBulkStr, ok := arguments["is_bulk_order"].(string); ok && isBulkStr != "" {
		context.IsBulkOrder = (isBulkStr == "true")
	}

	// Create pricing engine and compare models
	engine := pricing.NewPricingEngine()
	comparison := engine.ComparePricingModels(&originalEstimate, context)

	// Format response
	responseJSON, _ := json.MarshalIndent(comparison, "", "  ")
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
