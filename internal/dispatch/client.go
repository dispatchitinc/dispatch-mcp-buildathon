package dispatch

import (
	"dispatch-mcp-server/internal/auth"
	"dispatch-mcp-server/internal/config"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client     *resty.Client
	config     *config.Config
	authClient *auth.Client
	mockClient *MockClient
}

func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Check if we should use mock mode (when no auth token is provided)
	useMockMode := cfg.AuthToken == "" && !cfg.UseIDP

	if useMockMode {
		// Return a mock client for demo purposes
		mockClient, err := NewMockClient()
		if err != nil {
			return nil, err
		}

		// Create a wrapper that implements the same interface
		return &Client{
			client:     nil, // Not used in mock mode
			config:     cfg,
			authClient: nil,
			mockClient: mockClient,
		}, nil
	}

	client := resty.New()
	client.SetBaseURL(cfg.GraphQLEndpoint)
	client.SetHeader("Content-Type", "application/json")

	// Set up authentication
	var authClient *auth.Client
	if cfg.UseIDP {
		authConfig, err := auth.LoadConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load auth config: %v", err)
		}
		authClient = auth.NewClient(authConfig)
	} else {
		client.SetHeader("Authorization", "Bearer "+cfg.AuthToken)
	}

	return &Client{
		client:     client,
		config:     cfg,
		authClient: authClient,
		mockClient: nil,
	}, nil
}

func (c *Client) getAuthToken() (string, error) {
	if c.config.UseIDP && c.authClient != nil {
		return c.authClient.GetValidToken()
	}
	return c.config.AuthToken, nil
}

func (c *Client) CreateEstimate(input CreateEstimateInput) (*CreateEstimateResponse, error) {
	// Use mock client if available
	if c.mockClient != nil {
		return c.mockClient.CreateEstimate(input)
	}
	query := `
		mutation CreateEstimate($input: CreateEstimateInput!) {
			createEstimate(input: $input) {
				estimate {
					availableOrderOptions {
						serviceType
						estimatedDeliveryTimeUtc
						estimatedOrderCost
						vehicleType
						pickupLocationInfo {
							googlePlaceId
							lat
							lng
						}
						dropOffLocationsInfo {
							googlePlaceId
							lat
							lng
						}
						estimateInfo {
							serviceType
							vehicleType
							tollAmount
							estimatedOrderCost
							dedicatedVehicleRequested
							dedicatedVehicleFee
						}
						addOns
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Get auth token
	authToken, err := c.getAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %v", err)
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+authToken).
		SetBody(requestBody).
		Post("")

	if err != nil {
		return nil, fmt.Errorf("failed to create estimate: %v", err)
	}

	var response CreateEstimateResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &response, nil
}

func (c *Client) CreateOrder(input CreateOrderInput) (*CreateOrderResponse, error) {
	// Use mock client if available
	if c.mockClient != nil {
		return c.mockClient.CreateOrder(input)
	}
	query := `
		mutation CreateOrder($input: CreateOrderInput!) {
			createOrder(input: $input) {
				order {
					id
					status
					scheduledAt
					totalCost
					trackingNumber
					estimatedArrival
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Get auth token
	authToken, err := c.getAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %v", err)
	}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+authToken).
		SetBody(requestBody).
		Post("")

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	var response CreateOrderResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &response, nil
}
