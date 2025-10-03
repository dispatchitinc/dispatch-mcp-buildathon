package order

import (
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/graphql"
	"fmt"
	"os"
)

// OrderCreator handles order creation via GraphQL
type OrderCreator struct {
	graphqlClient *graphql.GraphQLClient
}

// OrderResult represents the result of order creation
type OrderResult struct {
	Order struct {
		ID      string `json:"id"`
		Druid   string `json:"druid"`
		Status  string `json:"status"`
		Pricing struct {
			TotalPrice float64 `json:"totalPrice"`
			BasePrice  float64 `json:"basePrice"`
			Discounts  []struct {
				Type   string  `json:"type"`
				Amount float64 `json:"amount"`
			} `json:"discounts"`
		} `json:"pricing"`
	} `json:"order"`
	Errors []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors"`
}

// NewOrderCreator creates a new order creator
func NewOrderCreator(graphqlEndpoint string) *OrderCreator {
	client := graphql.NewGraphQLClient(graphqlEndpoint)

	// Set authentication headers if available
	if apiKey := os.Getenv("GRAPHQL_API_KEY"); apiKey != "" {
		client.SetHeader("X-API-Key", apiKey)
	}
	if authToken := os.Getenv("GRAPHQL_AUTH_TOKEN"); authToken != "" {
		client.SetHeader("Authorization", "Bearer "+authToken)
	}

	return &OrderCreator{
		graphqlClient: client,
	}
}

// CreateOrder creates an order from order creation input
func (oc *OrderCreator) CreateOrder(input *OrderCreationInput) (*OrderResult, error) {
	// Validate order first
	if err := oc.validateOrder(input); err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Create the order
	result, err := oc.executeCreateOrder(input)
	if err != nil {
		return nil, fmt.Errorf("order creation failed: %w", err)
	}

	return result, nil
}

// ConvertFromDispatchTypes converts dispatch types to order types
func ConvertFromDispatchTypes(pickupInfo *dispatch.CreateOrderPickupInfoInput, dropOffs []dispatch.CreateOrderDropOffInfoInput, vehicleType *VehicleTypeInfo, capabilities []string, scheduling *SchedulingInput) *OrderCreationInput {
	input := &OrderCreationInput{
		OrganizationID: "default-org-id",       // TODO: Get from context or user
		JobName:        "Conversational Order", // TODO: Get from context
		VehicleTypeID:  "cargo_van",            // Default
		Capabilities:   capabilities,
		Scheduling:     scheduling,
	}

	// Convert pickup info
	if pickupInfo != nil {
		notes := ""
		if pickupInfo.PickupNotes != nil {
			notes = *pickupInfo.PickupNotes
		}
		input.PickupInfo = &PickupInfoInput{
			BusinessName: *pickupInfo.BusinessName,
			ContactName:  *pickupInfo.ContactName,
			ContactPhone: *pickupInfo.ContactPhoneNumber,
			Address:      convertAddress(pickupInfo.Location.Address),
			Notes:        notes,
		}
	}

	// Convert drop-offs
	input.DropOffs = make([]DropOffInfoInput, len(dropOffs))
	for i, dropOff := range dropOffs {
		notes := ""
		if dropOff.DropOffNotes != nil {
			notes = *dropOff.DropOffNotes
		}
		input.DropOffs[i] = DropOffInfoInput{
			BusinessName: *dropOff.BusinessName,
			ContactName:  *dropOff.ContactName,
			ContactPhone: *dropOff.ContactPhoneNumber,
			Address:      convertAddress(dropOff.Location.Address),
			Notes:        notes,
		}
	}

	// Set vehicle type
	if vehicleType != nil {
		input.VehicleTypeID = vehicleType.VehicleTypeID
	}

	return input
}

// convertAddress converts dispatch address to order address
func convertAddress(address *dispatch.AddressInput) *AddressInput {
	if address == nil {
		return nil
	}

	return &AddressInput{
		Street:  address.Street,
		City:    address.City,
		State:   address.State,
		ZipCode: address.ZipCode,
		Country: address.Country,
	}
}

// validateOrder validates the order before creation
func (oc *OrderCreator) validateOrder(input *OrderCreationInput) error {
	// Check required fields
	if input.PickupInfo == nil {
		return fmt.Errorf("pickup information is required")
	}
	if len(input.DropOffs) == 0 {
		return fmt.Errorf("at least one delivery location is required")
	}
	if input.VehicleTypeID == "" {
		return fmt.Errorf("vehicle type is required")
	}
	return nil
}

// executeCreateOrder executes the GraphQL create order mutation
func (oc *OrderCreator) executeCreateOrder(input *OrderCreationInput) (*OrderResult, error) {
	// For now, return a mock result since we don't have a real GraphQL endpoint
	// In production, this would call the actual GraphQL endpoint

	mockResult := &OrderResult{
		Order: struct {
			ID      string `json:"id"`
			Druid   string `json:"druid"`
			Status  string `json:"status"`
			Pricing struct {
				TotalPrice float64 `json:"totalPrice"`
				BasePrice  float64 `json:"basePrice"`
				Discounts  []struct {
					Type   string  `json:"type"`
					Amount float64 `json:"amount"`
				} `json:"discounts"`
			} `json:"pricing"`
		}{
			ID:     "ORD-12345",
			Druid:  "ABC123",
			Status: "pending",
			Pricing: struct {
				TotalPrice float64 `json:"totalPrice"`
				BasePrice  float64 `json:"basePrice"`
				Discounts  []struct {
					Type   string  `json:"type"`
					Amount float64 `json:"amount"`
				} `json:"discounts"`
			}{
				TotalPrice: 45.00,
				BasePrice:  50.00,
				Discounts: []struct {
					Type   string  `json:"type"`
					Amount float64 `json:"amount"`
				}{
					{Type: "multi_delivery", Amount: 5.00},
				},
			},
		},
	}

	return mockResult, nil
}

// TODO: Implement actual GraphQL calls when endpoint is available
/*
func (oc *OrderCreator) executeCreateOrder(input map[string]interface{}) (*OrderResult, error) {
	response, err := oc.graphqlClient.Execute(CreateOrderMutation, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return nil, err
	}

	// Parse the response
	var result OrderResult
	if err := json.Unmarshal(response.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse order result: %w", err)
	}

	return &result, nil
}
*/
