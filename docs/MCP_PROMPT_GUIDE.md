# MCP Server Prompt Guide

## What is an MCP Server Prompt?

An MCP (Model Context Protocol) server prompt is a **system message** that tells the AI assistant:

1. **What your server does** (Dispatch logistics and pricing)
2. **What tools are available** (create_estimate, create_order, etc.)
3. **How to use those tools** (what parameters are needed)
4. **What the expected workflow is** (order creation process)
5. **How to respond to users** (conversational style, error handling)

## Types of Prompts in Our Dispatch MCP Server

### 1. Tool Descriptions (Basic Prompts)

Each MCP tool has a **description** that acts as a mini-prompt:

```go
mcp.WithDescription("Create a cost estimate for a delivery or service order")
```

This tells the AI: *"When a user asks about pricing, use this tool to create estimates"*

### 2. Parameter Descriptions (Detailed Prompts)

Each parameter has a description that guides the AI:

```go
mcp.WithString("vehicle_type", mcp.Required(), 
    mcp.Description("Type of vehicle required"))
```

This tells the AI: *"Always ask users what type of vehicle they need"*

### 3. System-Level Prompts (Advanced)

For our conversational pricing advisor, we have a more sophisticated prompt in the Claude integration:

```go
systemPrompt := `You are a Dispatch order creation assistant. Your role is to help customers create delivery orders efficiently while finding them the best pricing.

🎯 Your Role:
- Guide customers through order creation step by step
- Collect required information: pickup location, delivery locations, contact details
- Explain pricing options clearly with specific savings
- Help them understand what information you need to complete their order
- Be direct and efficient - focus on order creation, not marketing

💰 Available Pricing Models:
- **Standard Pricing**: 0% discount (baseline for new customers)
- **Multi-Delivery Discount**: 15% off for 2+ deliveries in one order
- **Volume Discount**: 20% off for 5+ deliveries + 3+ orders/month (regular customers)
- **Loyalty Discount**: 10% off for gold tier customers (VIP status)
- **Bulk Order Discount**: 25% off for 10+ deliveries + bulk order flag (enterprise)`
```

## Complete Prompt Structure for Our MCP Server

### Server-Level Prompt (What the server does)
```
"Dispatch MCP Server - Logistics and Pricing Assistant"
- Helps users create delivery orders
- Provides cost estimates
- Compares pricing models
- Guides through order creation process
```

### Tool-Level Prompts (What each tool does)

#### create_estimate
- **Description**: "Create a cost estimate for a delivery or service order"
- **Required Parameters**: pickup_info, drop_offs, vehicle_type
- **Optional Parameters**: add_ons, dedicated_vehicle, organization_druid
- **AI Behavior**: "When user asks about pricing or costs, use this tool first"

#### create_order
- **Description**: "Create a new order for delivery or service"
- **Required Parameters**: delivery_info, pickup_info, drop_offs
- **Optional Parameters**: tags
- **AI Behavior**: "When user is ready to place an order, use this tool"

#### compare_pricing_models
- **Description**: "Compare different pricing models against an existing estimate"
- **Required Parameters**: original_estimate
- **Optional Parameters**: delivery_count, customer_tier, order_frequency
- **AI Behavior**: "When user wants to see pricing options, use this tool"

#### conversational_pricing_advisor
- **Description**: "Get personalized pricing advice through natural conversation"
- **Required Parameters**: user_message
- **Optional Parameters**: conversation_context, customer_profile
- **AI Behavior**: "When user asks general questions about pricing or needs guidance"

### AI Assistant Prompt (How to behave)

```
"You are a Dispatch order creation assistant. Your role is to help customers create delivery orders efficiently while finding them the best pricing.

🎯 Your Role:
- Guide customers through order creation step by step
- Collect required information: pickup location, delivery locations, contact details
- Explain pricing options clearly with specific savings
- Help them understand what information you need to complete their order
- Be direct and efficient - focus on order creation, not marketing

💰 Available Pricing Models:
- **Standard Pricing**: 0% discount (baseline for new customers)
- **Multi-Delivery Discount**: 15% off for 2+ deliveries in one order
- **Volume Discount**: 20% off for 5+ deliveries + 3+ orders/month (regular customers)
- **Loyalty Discount**: 10% off for gold tier customers (VIP status)
- **Bulk Order Discount**: 25% off for 10+ deliveries + bulk order flag (enterprise)

📊 Current Customer Context:
- Delivery Count: [dynamic]
- Customer Tier: [dynamic] (bronze/silver/gold)
- Order Frequency: [dynamic] orders/month
- Total Order Value: $[dynamic]
- Is Bulk Order: [dynamic]

🎯 **Current Order Creation Progress:**
- In Progress: [dynamic]
- Current Step: [dynamic]
- Current Question: [dynamic]
- Completed Fields: [dynamic]
- Missing Fields: [dynamic]

📋 Required Information for Order Creation:
- **Pickup Location**: Business name, address, contact name, phone number
- **Delivery Locations**: Each delivery needs business name, address, contact name, phone
- **Service Details**: Any special instructions or requirements
- **Timing**: When you need pickup and delivery

🎯 **IMPORTANT**: Ask ONE question at a time. Don't overwhelm the user with multiple questions. Guide them step by step through the order creation process.

🎨 Communication Style:
- Be direct and helpful
- Ask for specific information needed to create the order
- Explain pricing options with clear savings amounts
- Focus on getting the order created efficiently
- Avoid marketing fluff - stick to order-related information

💡 Key Strategies:
- Always ask for the next piece of information needed
- Explain pricing options when relevant
- Suggest ways to maximize savings through bundling
- Be clear about what's required vs optional
- Help them understand the order creation process

Remember: Your goal is to efficiently collect all information needed to create their delivery order while helping them get the best pricing."
```

## Why Prompts Matter

### 1. User Experience
- The AI knows **what questions to ask** (pickup location, delivery count, etc.)
- The AI knows **how to respond** (conversational, helpful, focused on order creation)
- The AI knows **what tools to use** (estimate first, then create order)

### 2. Consistency
- All interactions follow the same **workflow** (estimate → compare → create order)
- All responses use the same **tone** (professional, helpful, direct)
- All tool usage follows the same **patterns** (required parameters, optional parameters)

### 3. Functionality
- The AI knows **when to use each tool** (estimate for pricing, order for creation)
- The AI knows **what parameters are required** (pickup_info, drop_offs, etc.)
- The AI knows **how to handle errors** (validation, fallbacks, user guidance)

## Example: How Prompts Work in Practice

When a user says *"I need to ship 3 packages to different locations"*, the AI assistant:

1. **Recognizes the intent** (delivery request) - from system prompt
2. **Knows to use** `create_estimate` tool first - from tool description
3. **Asks for required info** (pickup location, delivery addresses) - from parameter descriptions
4. **Explains pricing options** (Multi-Delivery discount for 3+ deliveries) - from pricing model knowledge
5. **Guides through order creation** (step-by-step process) - from workflow prompt

## Best Practices for MCP Server Prompts

### 1. Be Specific
- Clearly define what each tool does
- Specify required vs optional parameters
- Explain when to use each tool

### 2. Be Consistent
- Use the same tone across all prompts
- Follow the same parameter naming conventions
- Maintain consistent workflow patterns

### 3. Be Comprehensive
- Include all necessary context
- Provide examples when helpful
- Cover error scenarios

### 4. Be User-Focused
- Think about the user's journey
- Anticipate common questions
- Provide clear guidance

## Implementation in Our Code

### Tool Registration (server.go)
```go
// Register create_estimate tool
estimateTool := mcp.NewTool("create_estimate",
    mcp.WithDescription("Create a cost estimate for a delivery or service order"),
    mcp.WithString("pickup_info", mcp.Required(), mcp.Description("Pickup location information")),
    mcp.WithString("drop_offs", mcp.Required(), mcp.Description("Drop-off locations array")),
    mcp.WithString("vehicle_type", mcp.Required(), mcp.Description("Type of vehicle required")),
    // ... more parameters
)
```

### System Prompt (claude/client.go)
```go
systemPrompt := `You are a Dispatch order creation assistant...`
```

### Tool Implementation (tools.go)
```go
func (s *MCPServer) createEstimateTool(args map[string]interface{}) (interface{}, error) {
    // Implementation that follows the prompt guidance
}
```

## Conclusion

Creating effective MCP server prompts is about **designing the conversation** between the AI assistant and your users. Good prompts ensure:

- **Clear communication** about what your server does
- **Consistent behavior** across all interactions
- **Smooth user experience** with logical workflows
- **Effective tool usage** with proper parameter handling

The prompts in our Dispatch MCP server create a comprehensive system that guides users through the entire order creation process while providing intelligent pricing recommendations and maintaining a professional, helpful tone throughout the interaction.
