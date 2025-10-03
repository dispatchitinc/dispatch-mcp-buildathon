package mcp

import (
	"dispatch-mcp-server/internal/conversation"
	"dispatch-mcp-server/internal/dispatch"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	dispatchClient     *dispatch.Client
	conversationEngine *conversation.ClaudeConversationEngine
}

func NewMCPServer() (*MCPServer, error) {
	dispatchClient, err := dispatch.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create dispatch client: %v", err)
	}

	conversationEngine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation engine: %v", err)
	}

	return &MCPServer{
		dispatchClient:     dispatchClient,
		conversationEngine: conversationEngine,
	}, nil
}

func (s *MCPServer) Run() error {
	fmt.Fprintf(os.Stderr, "Starting Dispatch MCP server...\n")

	srv := server.NewMCPServer(
		"dispatch-mcp-server",
		"1.0.0",
	)

	// Register create_estimate tool
	estimateTool := mcp.NewTool("create_estimate",
		mcp.WithDescription("Create a cost estimate for a delivery or service order"),
		mcp.WithString("pickup_info", mcp.Required(), mcp.Description("Pickup location information")),
		mcp.WithString("drop_offs", mcp.Required(), mcp.Description("Drop-off locations array")),
		mcp.WithString("vehicle_type", mcp.Required(), mcp.Description("Type of vehicle required")),
		mcp.WithString("add_ons", mcp.Description("Optional add-ons for delivery")),
		mcp.WithString("dedicated_vehicle", mcp.Description("Whether dedicated vehicle is requested")),
		mcp.WithString("organization_druid", mcp.Description("Organization ID")),
	)

	srv.AddTool(estimateTool, s.createEstimateTool)

	// Register create_order tool
	orderTool := mcp.NewTool("create_order",
		mcp.WithDescription("Create a new order for delivery or service"),
		mcp.WithString("delivery_info", mcp.Required(), mcp.Description("Delivery information")),
		mcp.WithString("pickup_info", mcp.Required(), mcp.Description("Pickup information")),
		mcp.WithString("drop_offs", mcp.Required(), mcp.Description("Drop-off locations array")),
		mcp.WithString("tags", mcp.Description("Optional order tags")),
	)

	srv.AddTool(orderTool, s.createOrderTool)

	// Register compare_pricing_models tool
	pricingTool := mcp.NewTool("compare_pricing_models",
		mcp.WithDescription("Compare different pricing models (multi-delivery, volume discounts, etc.) against an existing estimate"),
		mcp.WithString("original_estimate", mcp.Required(), mcp.Description("Original estimate data in JSON format")),
		mcp.WithString("delivery_count", mcp.Description("Number of deliveries in the order (default: 1)")),
		mcp.WithString("customer_tier", mcp.Description("Customer loyalty tier (bronze, silver, gold)")),
		mcp.WithString("order_frequency", mcp.Description("Number of orders per month (default: 1)")),
		mcp.WithString("total_order_value", mcp.Description("Total value of the order")),
		mcp.WithString("is_bulk_order", mcp.Description("Whether this is a bulk order (true/false)")),
	)

	srv.AddTool(pricingTool, s.comparePricingModelsTool)

	// Register conversational pricing advisor tool
	advisorTool := mcp.NewTool("conversational_pricing_advisor",
		mcp.WithDescription("Get personalized pricing advice through natural conversation"),
		mcp.WithString("user_message", mcp.Required(), mcp.Description("User's natural language message")),
		mcp.WithString("conversation_context", mcp.Description("Previous conversation context in JSON format")),
		mcp.WithString("customer_profile", mcp.Description("Customer information and preferences in JSON format")),
	)

	srv.AddTool(advisorTool, s.conversationalPricingAdvisorTool)

	fmt.Fprintf(os.Stderr, "MCP server initialized and listening...\n")

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	return nil
}
