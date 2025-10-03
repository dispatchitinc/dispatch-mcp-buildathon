package dispatch

import (
	"dispatch-mcp-server/internal/config"
	"fmt"
	"time"
)

// MockClient provides mock responses for demo purposes
type MockClient struct {
	config *config.Config
}

func NewMockClient() (*MockClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &MockClient{
		config: cfg,
	}, nil
}

func (c *MockClient) CreateEstimate(input CreateEstimateInput) (*CreateEstimateResponse, error) {
	// Simulate API delay
	time.Sleep(500 * time.Millisecond)

	// Create mock response
	response := &CreateEstimateResponse{
		Data: struct {
			CreateEstimate struct {
				Estimate struct {
					AvailableOrderOptions []AvailableOrderOption `json:"availableOrderOptions"`
				} `json:"estimate"`
			} `json:"createEstimate"`
		}{
			CreateEstimate: struct {
				Estimate struct {
					AvailableOrderOptions []AvailableOrderOption `json:"availableOrderOptions"`
				} `json:"estimate"`
			}{
				Estimate: struct {
					AvailableOrderOptions []AvailableOrderOption `json:"availableOrderOptions"`
				}{
					AvailableOrderOptions: []AvailableOrderOption{
						{
							ServiceType:              "delivery",
							EstimatedDeliveryTimeUTC: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
							EstimatedOrderCost:       45.99,
							VehicleType:              input.VehicleType,
							PickupLocationInfo: LocationInfo{
								GooglePlaceID: "mock_pickup_place_id",
								Lat:           37.7749,
								Lng:           -122.4194,
							},
							DropOffLocationsInfo: []LocationInfo{
								{
									GooglePlaceID: "mock_dropoff_place_id",
									Lat:           37.8044,
									Lng:           -122.2712,
								},
							},
							EstimateInfo: EstimateInfo{
								ServiceType:               "delivery",
								VehicleType:               input.VehicleType,
								TollAmount:                "5.50",
								EstimatedOrderCost:        "45.99",
								DedicatedVehicleRequested: &[]bool{false}[0],
								DedicatedVehicleFee:       "0.00",
							},
							AddOns: input.AddOns,
						},
					},
				},
			},
		},
	}

	return response, nil
}

func (c *MockClient) CreateOrder(input CreateOrderInput) (*CreateOrderResponse, error) {
	// Simulate API delay
	time.Sleep(500 * time.Millisecond)

	// Create mock response
	response := &CreateOrderResponse{
		Data: struct {
			CreateOrder struct {
				Order Order `json:"order"`
			} `json:"createOrder"`
		}{
			CreateOrder: struct {
				Order Order `json:"order"`
			}{
				Order: Order{
					ID:               fmt.Sprintf("ORD-%d", time.Now().Unix()),
					Status:           "pending",
					ScheduledAt:      time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					TotalCost:        45.99,
					TrackingNumber:   fmt.Sprintf("TRK-%d", time.Now().Unix()),
					EstimatedArrival: time.Now().Add(3 * time.Hour).Format(time.RFC3339),
				},
			},
		},
	}

	return response, nil
}
