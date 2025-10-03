# MCP Implementation Research

## 1. Chris's Jira MCP Implementation Pattern

Based on research of MCP implementations and the general MCP pattern, here's what we understand about how MCP servers typically work:

### MCP Server Architecture

**Core Components:**
```
MCP Server
├── Server (Orchestrator)
│   ├── Handles AI client requests
│   ├── Routes to appropriate tools/capabilities
│   └── Manages responses
├── Tools/Capabilities
│   ├── Specific action handlers (create, update, query)
│   ├── API integration logic
│   └── Data transformation
└── Configuration
    ├── Environment variables
    ├── API credentials
    └── Endpoint configuration
```

### Key Patterns from MCP Implementations:

1. **Tool Definition**:
   - Each tool is a discrete function that AI can call
   - Tools have defined input parameters and output structures
   - Tools handle specific business operations

2. **API Integration**:
   - HTTP/GraphQL clients for external services
   - Authentication via environment variables
   - Request/response transformation

3. **Error Handling**:
   - Comprehensive validation
   - Meaningful error messages
   - Graceful failure handling

### Typical Go MCP Structure:

```go
// Server structure
type MCPServer struct {
    tools    map[string]Tool
    client   *http.Client
    config   *Config
}

// Tool interface
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// Configuration
type Config struct {
    APIEndpoint   string
    AuthToken     string
    ServerPort    string
}
```

## 2. Target GraphQL Mutations Analysis

### Expected Structure for `create_estimate.rb` and `create_order.rb`

Based on typical GraphQL mutation patterns for order/estimate systems:

### Create Estimate Mutation

**Likely Input Parameters:**
```ruby
# create_estimate.rb
module Mutations
  module Order
    class CreateEstimate < BaseMutation
      argument :customer_id, ID, required: true
      argument :service_type, String, required: true
      argument :location, Types::LocationInput, required: true
      argument :pickup_location, Types::LocationInput, required: false
      argument :dropoff_location, Types::LocationInput, required: false
      argument :items, [Types::ItemInput], required: false
      argument :requested_date, GraphQL::Types::ISO8601DateTime, required: false
      argument :notes, String, required: false
      
      field :estimate, Types::EstimateType, null: true
      field :errors, [String], null: true
      
      def resolve(...)
        # Business logic
        # Validation
        # Calculation
        # Return estimate
      end
    end
  end
end
```

**Expected Fields:**
- `customer_id`: Identifier for the customer
- `service_type`: Type of service (e.g., "delivery", "moving")
- `location`: Service location details (address, coordinates)
- `pickup_location`: Where to pick up items (for delivery)
- `dropoff_location`: Where to drop off items (for delivery)
- `items`: List of items to be transported
- `requested_date`: When the service is needed
- `notes`: Additional information

**Expected Response:**
```json
{
  "data": {
    "createEstimate": {
      "estimate": {
        "id": "est_123",
        "totalCost": 150.00,
        "estimatedDuration": 120,
        "validUntil": "2025-10-15T12:00:00Z",
        "breakdown": {
          "baseFare": 100.00,
          "distanceCost": 30.00,
          "serviceFee": 20.00
        }
      },
      "errors": null
    }
  }
}
```

### Create Order Mutation

**Likely Input Parameters:**
```ruby
# create_order.rb
module Mutations
  module Order
    class CreateOrder < BaseMutation
      argument :estimate_id, ID, required: false
      argument :customer_id, ID, required: true
      argument :service_type, String, required: true
      argument :location, Types::LocationInput, required: true
      argument :pickup_location, Types::LocationInput, required: false
      argument :dropoff_location, Types::LocationInput, required: false
      argument :items, [Types::ItemInput], required: false
      argument :scheduled_at, GraphQL::Types::ISO8601DateTime, required: true
      argument :payment_method_id, ID, required: true
      argument :special_instructions, String, required: false
      argument :contact_phone, String, required: true
      
      field :order, Types::OrderType, null: true
      field :errors, [String], null: true
      
      def resolve(...)
        # Business logic
        # Validation
        # Order creation
        # Payment processing initiation
        # Return order
      end
    end
  end
end
```

**Expected Fields:**
- `estimate_id`: Optional reference to a previously created estimate
- `customer_id`: Identifier for the customer
- `service_type`: Type of service being ordered
- `location`: Service location details
- `pickup_location`: Pickup location (for delivery services)
- `dropoff_location`: Dropoff location (for delivery services)
- `items`: List of items to transport
- `scheduled_at`: When the service should occur
- `payment_method_id`: How the customer will pay
- `special_instructions`: Additional delivery/service instructions
- `contact_phone`: Contact number for the service

**Expected Response:**
```json
{
  "data": {
    "createOrder": {
      "order": {
        "id": "ord_456",
        "status": "pending",
        "scheduledAt": "2025-10-10T14:00:00Z",
        "totalCost": 150.00,
        "trackingNumber": "TRK123456",
        "estimatedArrival": "2025-10-10T16:00:00Z"
      },
      "errors": null
    }
  }
}
```

## Go MCP Service Design

### Mapping to Go Structures

```go
// EstimateInput represents the input for creating an estimate
type EstimateInput struct {
    CustomerID      string              `json:"customer_id"`
    ServiceType     string              `json:"service_type"`
    Location        LocationInput       `json:"location"`
    PickupLocation  *LocationInput      `json:"pickup_location,omitempty"`
    DropoffLocation *LocationInput      `json:"dropoff_location,omitempty"`
    Items           []ItemInput         `json:"items,omitempty"`
    RequestedDate   *time.Time          `json:"requested_date,omitempty"`
    Notes           string              `json:"notes,omitempty"`
}

// OrderInput represents the input for creating an order
type OrderInput struct {
    EstimateID          *string            `json:"estimate_id,omitempty"`
    CustomerID          string             `json:"customer_id"`
    ServiceType         string             `json:"service_type"`
    Location            LocationInput      `json:"location"`
    PickupLocation      *LocationInput     `json:"pickup_location,omitempty"`
    DropoffLocation     *LocationInput     `json:"dropoff_location,omitempty"`
    Items               []ItemInput        `json:"items,omitempty"`
    ScheduledAt         time.Time          `json:"scheduled_at"`
    PaymentMethodID     string             `json:"payment_method_id"`
    SpecialInstructions string             `json:"special_instructions,omitempty"`
    ContactPhone        string             `json:"contact_phone"`
}

// LocationInput represents location details
type LocationInput struct {
    Address     string  `json:"address"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    City        string  `json:"city"`
    State       string  `json:"state"`
    ZipCode     string  `json:"zip_code"`
}

// ItemInput represents an item to be transported
type ItemInput struct {
    Name        string  `json:"name"`
    Description string  `json:"description,omitempty"`
    Quantity    int     `json:"quantity"`
    Weight      float64 `json:"weight,omitempty"`
    Dimensions  string  `json:"dimensions,omitempty"`
}
```

### MCP Tool Implementation

```go
// CreateEstimateTool implements the create estimate MCP tool
type CreateEstimateTool struct {
    client *DispatchClient
}

func (t *CreateEstimateTool) Name() string {
    return "create_estimate"
}

func (t *CreateEstimateTool) Description() string {
    return "Create a cost estimate for a delivery or service order"
}

func (t *CreateEstimateTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Parse params into EstimateInput
    var input EstimateInput
    if err := mapToStruct(params, &input); err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }
    
    // Call Dispatch API
    estimate, err := t.client.CreateEstimate(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to create estimate: %w", err)
    }
    
    return estimate, nil
}

// CreateOrderTool implements the create order MCP tool
type CreateOrderTool struct {
    client *DispatchClient
}

func (t *CreateOrderTool) Name() string {
    return "create_order"
}

func (t *CreateOrderTool) Description() string {
    return "Create a new order for delivery or service"
}

func (t *CreateOrderTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Parse params into OrderInput
    var input OrderInput
    if err := mapToStruct(params, &input); err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }
    
    // Call Dispatch API
    order, err := t.client.CreateOrder(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    return order, nil
}
```

## Next Steps

### To Confirm and Refine:

1. **Access Actual Mutation Code**: Review the actual `create_estimate.rb` and `create_order.rb` files to:
   - Confirm exact parameter names and types
   - Understand validation rules
   - Identify any special business logic

2. **Test Mutations Manually**: Use GraphQL playground or similar tool to:
   - Test the mutations directly
   - Understand the exact request/response format
   - Identify required vs optional fields

3. **Reference Chris's Implementation**: Review the actual Jira MCP server.go to:
   - Understand the exact MCP pattern used
   - See how tools are registered
   - Learn authentication approach

### Questions to Answer:

- What is the exact GraphQL endpoint URL for Dispatch?
- How does authentication work? (API key, OAuth, session token?)
- Are there any rate limits or special requirements?
- What are the exact field names and types for the mutations?
- Are there any validation rules we need to know about?

---
*This document will be updated as we get access to the actual code and API documentation.*
