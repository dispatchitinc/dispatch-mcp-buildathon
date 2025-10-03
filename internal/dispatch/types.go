package dispatch

// CreateEstimateInput represents the input for creating an estimate
type CreateEstimateInput struct {
	AddOns             []string           `json:"add_ons,omitempty"`
	DedicatedVehicle   *bool              `json:"dedicated_vehicle,omitempty"`
	DropOffs           []DropOffInfoInput `json:"drop_offs"`
	DropOffDateTimeUTC *string            `json:"drop_off_date_time_utc,omitempty"`
	OrganizationDruid  *string            `json:"organization_druid,omitempty"`
	PickupInfo         PickupInfoInput    `json:"pickup_info"`
	VehicleType        string             `json:"vehicle_type"`
}

type PickupInfoInput struct {
	BusinessName      string        `json:"business_name"`
	Location          LocationInput `json:"location"`
	PickupDateTimeUTC *string       `json:"pickup_date_time_utc,omitempty"`
}

type DropOffInfoInput struct {
	BusinessName    string        `json:"business_name"`
	EstimatedWeight *int          `json:"estimated_weight,omitempty"`
	Location        LocationInput `json:"location"`
}

type LocationInput struct {
	Address        *AddressInput        `json:"address,omitempty"`
	GeoCoordinates *GeoCoordinatesInput `json:"geo_coordinates,omitempty"`
}

type AddressInput struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type GeoCoordinatesInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	AddOns       []string                      `json:"add_ons,omitempty"`
	DeliveryInfo DeliveryInfoInput             `json:"delivery_info"`
	DropOffs     []CreateOrderDropOffInfoInput `json:"drop_offs"`
	PickupInfo   CreateOrderPickupInfoInput    `json:"pickup_info"`
	Tags         []TagInput                    `json:"tags,omitempty"`
}

type DeliveryInfoInput struct {
	ServiceType       string  `json:"service_type"`
	OrganizationDruid *string `json:"organization_druid,omitempty"`
}

type CreateOrderPickupInfoInput struct {
	BusinessName       *string        `json:"business_name,omitempty"`
	ContactName        *string        `json:"contact_name,omitempty"`
	ContactPhoneNumber *string        `json:"contact_phone_number,omitempty"`
	Location           *LocationInput `json:"location,omitempty"`
	PickupNotes        *string        `json:"pickup_notes,omitempty"`
}

type CreateOrderDropOffInfoInput struct {
	BusinessName       *string        `json:"business_name,omitempty"`
	ContactName        *string        `json:"contact_name,omitempty"`
	ContactPhoneNumber *string        `json:"contact_phone_number,omitempty"`
	Location           *LocationInput `json:"location,omitempty"`
	DropOffNotes       *string        `json:"drop_off_notes,omitempty"`
}

type TagInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response types
type CreateEstimateResponse struct {
	Data struct {
		CreateEstimate struct {
			Estimate struct {
				AvailableOrderOptions []AvailableOrderOption `json:"availableOrderOptions"`
			} `json:"estimate"`
		} `json:"createEstimate"`
	} `json:"data"`
}

type AvailableOrderOption struct {
	ServiceType              string         `json:"serviceType"`
	EstimatedDeliveryTimeUTC string         `json:"estimatedDeliveryTimeUtc"`
	EstimatedOrderCost       float64        `json:"estimatedOrderCost"`
	VehicleType              string         `json:"vehicleType"`
	PickupLocationInfo       LocationInfo   `json:"pickupLocationInfo"`
	DropOffLocationsInfo     []LocationInfo `json:"dropOffLocationsInfo"`
	EstimateInfo             EstimateInfo   `json:"estimateInfo"`
	AddOns                   []string       `json:"addOns"`
}

type LocationInfo struct {
	GooglePlaceID string  `json:"googlePlaceId"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
}

type EstimateInfo struct {
	ServiceType               string `json:"serviceType"`
	VehicleType               string `json:"vehicleType"`
	TollAmount                string `json:"tollAmount"`
	EstimatedOrderCost        string `json:"estimatedOrderCost"`
	DedicatedVehicleRequested *bool  `json:"dedicatedVehicleRequested"`
	DedicatedVehicleFee       string `json:"dedicatedVehicleFee"`
}

type CreateOrderResponse struct {
	Data struct {
		CreateOrder struct {
			Order Order `json:"order"`
		} `json:"createOrder"`
	} `json:"data"`
}

type Order struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	ScheduledAt      string  `json:"scheduledAt"`
	TotalCost        float64 `json:"totalCost"`
	TrackingNumber   string  `json:"trackingNumber"`
	EstimatedArrival string  `json:"estimatedArrival"`
}
