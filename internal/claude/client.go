package claude

import (
	"bytes"
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

	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
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
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Content      []struct {
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
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

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
	systemPrompt := `You are a helpful pricing advisor for a delivery service. You help customers find the best pricing options based on their needs.

Available pricing models:
- Standard Pricing: 0% discount (baseline)
- Multi-Delivery Discount: 15% off for 2+ deliveries
- Volume Discount: 20% off for 5+ deliveries + 3+ orders/month
- Loyalty Discount: 10% off for gold tier customers
- Bulk Order Discount: 25% off for 10+ deliveries + bulk order flag

Customer context:
- Delivery Count: ` + fmt.Sprintf("%d", context.DeliveryCount) + `
- Customer Tier: ` + context.CustomerTier + `
- Order Frequency: ` + fmt.Sprintf("%d", context.OrderFrequency) + ` orders/month
- Total Order Value: $` + fmt.Sprintf("%.2f", context.TotalOrderValue) + `
- Is Bulk Order: ` + fmt.Sprintf("%t", context.IsBulkOrder) + `

Provide helpful, conversational advice about pricing options. Be specific about which discounts they qualify for and how much they could save.`

	request := MessageRequest{
		Model:     "claude-3-sonnet-20240229",
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
	DeliveryCount    int     `json:"delivery_count"`
	CustomerTier     string  `json:"customer_tier"`
	OrderFrequency   int     `json:"order_frequency"`
	TotalOrderValue  float64 `json:"total_order_value"`
	IsBulkOrder      bool    `json:"is_bulk_order"`
}
