package order

// OrderCreationInput represents the input for order creation
type OrderCreationInput struct {
	OrganizationID string             `json:"organizationId"`
	JobName        string             `json:"jobName"`
	PickupInfo     *PickupInfoInput   `json:"pickupInfo"`
	DropOffs       []DropOffInfoInput `json:"dropOffs"`
	VehicleTypeID  string             `json:"vehicleTypeId"`
	Capabilities   []string           `json:"capabilities"`
	Scheduling     *SchedulingInput   `json:"scheduling"`
}

// PickupInfoInput represents pickup information
type PickupInfoInput struct {
	BusinessName string        `json:"businessName"`
	ContactName  string        `json:"contactName"`
	ContactPhone string        `json:"contactPhone"`
	Address      *AddressInput `json:"address"`
	Notes        string        `json:"notes"`
}

// DropOffInfoInput represents delivery information
type DropOffInfoInput struct {
	BusinessName string        `json:"businessName"`
	ContactName  string        `json:"contactName"`
	ContactPhone string        `json:"contactPhone"`
	Address      *AddressInput `json:"address"`
	Notes        string        `json:"notes"`
}

// AddressInput represents address information
type AddressInput struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

// SchedulingInput represents scheduling information
type SchedulingInput struct {
	PickupTime   string `json:"pickupTime"`
	DeliveryTime string `json:"deliveryTime"`
	PickupDate   string `json:"pickupDate"`
	DeliveryDate string `json:"deliveryDate"`
	TimeZone     string `json:"timeZone"`
}

// VehicleTypeInfo represents vehicle type information
type VehicleTypeInfo struct {
	VehicleTypeID   string   `json:"vehicleTypeId"`
	VehicleTypeName string   `json:"vehicleTypeName"`
	CustomTypes     []string `json:"customTypes,omitempty"`
}
