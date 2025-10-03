package claude

import (
	"bytes"
	"dispatch-mcp-server/internal/dispatch"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client represents a Claude API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Claude client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	// Check if we should use AI Hub proxy
	useAIHub := os.Getenv("USE_AI_HUB")
	baseURL := "https://api.anthropic.com/v1"

	if useAIHub == "true" {
		// Use AI Hub proxy endpoint
		aiHubEndpoint := os.Getenv("AI_HUB_ENDPOINT")
		if aiHubEndpoint == "" {
			// Default AI Hub endpoint for Dispatch
			aiHubEndpoint = "https://aihub.dispatchit.com/v1"
		}
		baseURL = aiHubEndpoint
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// MessageRequest represents a request to Claude
type MessageRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessageResponse represents a response from Claude
type MessageResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// CreateMessage sends a message to Claude and returns the response
func (c *Client) CreateMessage(request MessageRequest) (*MessageResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Check if using AI Hub for different authentication
	useAIHub := os.Getenv("USE_AI_HUB")
	if useAIHub == "true" {
		// AI Hub might use different authentication
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		// AI Hub might not need anthropic-version header
	} else {
		// Direct Anthropic API
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

// CreatePricingAdvisorMessage creates a message for the pricing advisor
func (c *Client) CreatePricingAdvisorMessage(userMessage string, context *PricingContext) (*MessageResponse, error) {
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
- **Bulk Order Discount**: 25% off for 10+ deliveries + bulk order flag (enterprise)

📊 Current Customer Context:
- Delivery Count: ` + fmt.Sprintf("%d", context.DeliveryCount) + `
- Customer Tier: ` + context.CustomerTier + ` (bronze/silver/gold)
- Order Frequency: ` + fmt.Sprintf("%d", context.OrderFrequency) + ` orders/month
- Total Order Value: $` + fmt.Sprintf("%.2f", context.TotalOrderValue) + `
- Is Bulk Order: ` + fmt.Sprintf("%t", context.IsBulkOrder) + `

🎯 **Current Order Creation Progress:**
- In Progress: ` + fmt.Sprintf("%t", context.OrderCreation.InProgress) + `
- Current Step: ` + context.OrderCreation.Step + `
- Current Question: ` + context.OrderCreation.CurrentQuestion + `
- Completed Fields: ` + fmt.Sprintf("%v", context.OrderCreation.CompletedFields) + `
- Missing Fields: ` + fmt.Sprintf("%v", context.OrderCreation.MissingFields) + `

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

Remember: Your goal is to efficiently collect all information needed to create their delivery order while helping them get the best pricing.`

	// Check if using AI Hub for different model names
	useAIHub := os.Getenv("USE_AI_HUB")
	modelName := "claude-3-sonnet-20240229"

	if useAIHub == "true" {
		// Use AI Hub model names
		aiHubModel := os.Getenv("AI_HUB_MODEL")
		if aiHubModel != "" {
			modelName = aiHubModel
		} else {
			// Default to claude-sonnet for conversational pricing (better for complex reasoning)
			modelName = "claude-sonnet"
		}
	}

	request := MessageRequest{
		Model:     modelName,
		MaxTokens: 1000,
		Messages: []Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		System: systemPrompt,
	}

	return c.CreateMessage(request)
}

// PricingContext represents the context for pricing conversations
type PricingContext struct {
	DeliveryCount   int                `json:"delivery_count"`
	CustomerTier    string             `json:"customer_tier"`
	OrderFrequency  int                `json:"order_frequency"`
	TotalOrderValue float64            `json:"total_order_value"`
	IsBulkOrder     bool               `json:"is_bulk_order"`
	OrderCreation   OrderCreationState `json:"order_creation"`
}

// OrderCreationState tracks the progress of order creation
type OrderCreationState struct {
	InProgress           bool                                   `json:"in_progress"`
	Step                 string                                 `json:"step"`             // "pickup", "deliveries", "contact", "review"
	CurrentQuestion      string                                 `json:"current_question"` // "pickup_business", "pickup_address", "pickup_contact", "pickup_phone", etc.
	PickupInfo           *dispatch.CreateOrderPickupInfoInput   `json:"pickup_info,omitempty"`
	DropOffs             []dispatch.CreateOrderDropOffInfoInput `json:"drop_offs,omitempty"`
	DeliveryInfo         *dispatch.DeliveryInfoInput            `json:"delivery_info,omitempty"`
	MissingFields        []string                               `json:"missing_fields"`
	CompletedFields      []string                               `json:"completed_fields"`
	CurrentDeliveryIndex int                                    `json:"current_delivery_index"` // Which delivery we're collecting info for
}
