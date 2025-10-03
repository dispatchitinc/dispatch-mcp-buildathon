# GraphQL Integration Plan for Order Creation

## Overview
Instead of integrating directly with ms-monolith services, we'll use GraphQL endpoints to handle order creation. This provides a clean API interface and follows proper microservice patterns.

## GraphQL Endpoints Analysis

### **1. Order Creation GraphQL Mutations**

Based on the ms-monolith analysis, we need these GraphQL mutations:

```graphql
# Create a new order
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    order {
      id
      druid
      status
      pricing {
        totalPrice
        basePrice
        discounts {
          type
          amount
        }
      }
    }
    errors {
      field
      message
    }
  }
}

# Validate order before creation
mutation ValidateOrder($input: ValidateOrderInput!) {
  validateOrder(input: $input) {
    valid
    errors {
      field
      message
    }
    warnings {
      field
      message
    }
  }
}

# Get pricing estimate
query GetOrderPricing($input: PricingInput!) {
  getOrderPricing(input: $input) {
    basePrice
    totalPrice
    discounts {
      type
      amount
      description
    }
    availableVehicleTypes {
      id
      name
      price
    }
  }
}
```

### **2. Input Types for GraphQL**

```graphql
input CreateOrderInput {
  organizationId: ID!
  jobName: String!
  pickupInfo: PickupInfoInput!
  dropOffs: [DropOffInfoInput!]!
  vehicleTypeId: ID!
  capabilities: [String!]
  scheduling: SchedulingInput!
  specialRequirements: [String!]
}

input PickupInfoInput {
  businessName: String!
  contactName: String!
  contactPhone: String!
  address: AddressInput!
  notes: String
}

input DropOffInfoInput {
  businessName: String!
  contactName: String!
  contactPhone: String!
  address: AddressInput!
  notes: String
  packageCount: Int
  packageWeight: Float
}

input AddressInput {
  street: String!
  city: String!
  state: String!
  zipCode: String!
  country: String!
  googlePlaceId: String
  coordinates: CoordinatesInput
}

input CoordinatesInput {
  latitude: Float!
  longitude: Float!
}

input SchedulingInput {
  pickupTime: String!
  deliveryTime: String!
  pickupDate: String!
  deliveryDate: String!
  timeZone: String
}
```

## Implementation Plan

### **Phase 1: GraphQL Client Setup**

#### **1.1 Create GraphQL Client**
```go
// internal/graphql/client.go
package graphql

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type GraphQLClient struct {
    endpoint string
    client   *http.Client
    headers  map[string]string
}

func NewGraphQLClient(endpoint string) *GraphQLClient {
    return &GraphQLClient{
        endpoint: endpoint,
        client:   &http.Client{},
        headers:  make(map[string]string),
    }
}

func (c *GraphQLClient) SetHeader(key, value string) {
    c.headers[key] = value
}

func (c *GraphQLClient) Execute(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
    payload := GraphQLRequest{
        Query:     query,
        Variables: variables,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result GraphQLResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

#### **1.2 GraphQL Response Types**
```go
type GraphQLRequest struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
    Data   interface{} `json:"data"`
    Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLError struct {
    Message   string                 `json:"message"`
    Locations []GraphQLLocation     `json:"locations,omitempty"`
    Path      []interface{}         `json:"path,omitempty"`
}

type GraphQLLocation struct {
    Line   int `json:"line"`
    Column int `json:"column"`
}
```

### **Phase 2: Order Creation Service**

#### **2.1 Order Creation Service**
```go
// internal/order/creator.go
package order

import (
    "dispatch-mcp-server/internal/graphql"
    "dispatch-mcp-server/internal/conversation"
    "fmt"
)

type OrderCreator struct {
    graphqlClient *graphql.GraphQLClient
}

func NewOrderCreator(graphqlEndpoint string) *OrderCreator {
    return &OrderCreator{
        graphqlClient: graphql.NewGraphQLClient(graphqlEndpoint),
    }
}

func (oc *OrderCreator) CreateOrder(context *conversation.ConversationContext) (*OrderResult, error) {
    // Convert conversation context to GraphQL input
    input := oc.convertContextToGraphQLInput(context)
    
    // Validate order first
    if err := oc.validateOrder(input); err != nil {
        return nil, err
    }
    
    // Get pricing estimate
    pricing, err := oc.getPricing(input)
    if err != nil {
        return nil, err
    }
    
    // Create the order
    result, err := oc.executeCreateOrder(input)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}

func (oc *OrderCreator) convertContextToGraphQLInput(context *conversation.ConversationContext) map[string]interface{} {
    return map[string]interface{}{
        "organizationId": "default-org-id", // TODO: Get from context
        "jobName": context.OrderCreation.JobName,
        "pickupInfo": oc.convertPickupInfo(context.OrderCreation.PickupInfo),
        "dropOffs": oc.convertDropOffs(context.OrderCreation.DropOffs),
        "vehicleTypeId": context.OrderCreation.VehicleType.VehicleTypeID,
        "capabilities": context.OrderCreation.Capabilities,
        "scheduling": oc.convertScheduling(context.OrderCreation.SchedulingInfo),
    }
}

func (oc *OrderCreator) convertPickupInfo(pickupInfo *dispatch.CreateOrderPickupInfoInput) map[string]interface{} {
    if pickupInfo == nil {
        return nil
    }
    
    return map[string]interface{}{
        "businessName": *pickupInfo.BusinessName,
        "contactName":  *pickupInfo.ContactName,
        "contactPhone": *pickupInfo.ContactPhoneNumber,
        "address": oc.convertAddress(pickupInfo.Location.Address),
        "notes": pickupInfo.Notes,
    }
}

func (oc *OrderCreator) convertDropOffs(dropOffs []dispatch.CreateOrderDropOffInfoInput) []map[string]interface{} {
    result := make([]map[string]interface{}, len(dropOffs))
    
    for i, dropOff := range dropOffs {
        result[i] = map[string]interface{}{
            "businessName": *dropOff.BusinessName,
            "contactName":  *dropOff.ContactName,
            "contactPhone": *dropOff.ContactPhoneNumber,
            "address": oc.convertAddress(dropOff.Location.Address),
            "notes": dropOff.DropOffNotes,
        }
    }
    
    return result
}

func (oc *OrderCreator) convertAddress(address *dispatch.AddressInput) map[string]interface{} {
    if address == nil {
        return nil
    }
    
    return map[string]interface{}{
        "street": address.Street,
        "city":   address.City,
        "state":  address.State,
        "zipCode": address.ZipCode,
        "country": address.Country,
    }
}

func (oc *OrderCreator) convertScheduling(scheduling *conversation.SchedulingInfo) map[string]interface{} {
    if scheduling == nil {
        return nil
    }
    
    return map[string]interface{}{
        "pickupTime":   scheduling.PickupTime,
        "deliveryTime": scheduling.DeliveryTime,
        "pickupDate":   scheduling.PickupDate,
        "deliveryDate": scheduling.DeliveryDate,
        "timeZone":     scheduling.TimeZone,
    }
}
```

#### **2.2 GraphQL Queries and Mutations**
```go
// internal/order/queries.go
package order

const CreateOrderMutation = `
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    order {
      id
      druid
      status
      pricing {
        totalPrice
        basePrice
        discounts {
          type
          amount
        }
      }
    }
    errors {
      field
      message
    }
  }
}
`

const ValidateOrderMutation = `
mutation ValidateOrder($input: ValidateOrderInput!) {
  validateOrder(input: $input) {
    valid
    errors {
      field
      message
    }
    warnings {
      field
      message
    }
  }
}
`

const GetPricingQuery = `
query GetOrderPricing($input: PricingInput!) {
  getOrderPricing(input: $input) {
    basePrice
    totalPrice
    discounts {
      type
      amount
      description
    }
    availableVehicleTypes {
      id
      name
      price
    }
  }
}
`
```

### **Phase 3: Integration with Conversation Engine**

#### **3.1 Update Conversation Engine**
```go
// internal/conversation/claude_engine.go

// Add order creation to the review step
func (ce *ClaudeConversationEngine) handleReviewStep(message string, context *ConversationContext) string {
    // Generate order summary
    summary := ce.generateOrderSummary(context)
    
    if strings.Contains(strings.ToLower(message), "yes") || 
       strings.Contains(strings.ToLower(message), "create") ||
       strings.Contains(strings.ToLower(message), "confirm") {
        
        // Create the order using GraphQL
        orderCreator := order.NewOrderCreator(os.Getenv("GRAPHQL_ENDPOINT"))
        result, err := orderCreator.CreateOrder(context)
        
        if err != nil {
            return fmt.Sprintf("I encountered an error creating your order: %v. Please try again or contact support.", err)
        }
        
        return fmt.Sprintf("🎉 Order created successfully!\n\nOrder ID: %s\nDruid: %s\nTotal Price: $%.2f\n\nYou'll receive a confirmation email shortly.", 
            result.Order.ID, result.Order.Druid, result.Order.Pricing.TotalPrice)
    }
    
    return summary + "\n\nShould I create this order for you?"
}
```

#### **3.2 Order Summary Generation**
```go
func (ce *ClaudeConversationEngine) generateOrderSummary(context *ConversationContext) string {
    var summary strings.Builder
    
    summary.WriteString("📋 **Order Summary**\n\n")
    
    // Pickup information
    if context.OrderCreation.PickupInfo != nil {
        summary.WriteString("**Pickup Location:**\n")
        summary.WriteString(fmt.Sprintf("- Business: %s\n", *context.OrderCreation.PickupInfo.BusinessName))
        summary.WriteString(fmt.Sprintf("- Contact: %s (%s)\n", *context.OrderCreation.PickupInfo.ContactName, *context.OrderCreation.PickupInfo.ContactPhoneNumber))
        summary.WriteString(fmt.Sprintf("- Address: %s, %s, %s %s\n", 
            context.OrderCreation.PickupInfo.Location.Address.Street,
            context.OrderCreation.PickupInfo.Location.Address.City,
            context.OrderCreation.PickupInfo.Location.Address.State,
            context.OrderCreation.PickupInfo.Location.Address.ZipCode))
    }
    
    // Delivery information
    if len(context.OrderCreation.DropOffs) > 0 {
        summary.WriteString("\n**Delivery Locations:**\n")
        for i, dropOff := range context.OrderCreation.DropOffs {
            summary.WriteString(fmt.Sprintf("%d. %s\n", i+1, *dropOff.BusinessName))
            summary.WriteString(fmt.Sprintf("   Contact: %s (%s)\n", *dropOff.ContactName, *dropOff.ContactPhoneNumber))
            summary.WriteString(fmt.Sprintf("   Address: %s, %s, %s %s\n", 
                dropOff.Location.Address.Street,
                dropOff.Location.Address.City,
                dropOff.Location.Address.State,
                dropOff.Location.Address.ZipCode))
        }
    }
    
    // Vehicle type
    if context.OrderCreation.VehicleType != nil {
        summary.WriteString(fmt.Sprintf("\n**Vehicle Type:** %s\n", context.OrderCreation.VehicleType.VehicleTypeName))
    }
    
    // Capabilities
    if len(context.OrderCreation.Capabilities) > 0 {
        summary.WriteString(fmt.Sprintf("\n**Special Services:** %s\n", strings.Join(context.OrderCreation.Capabilities, ", ")))
    }
    
    // Scheduling
    if context.OrderCreation.SchedulingInfo != nil {
        summary.WriteString(fmt.Sprintf("\n**Schedule:**\n"))
        summary.WriteString(fmt.Sprintf("- Pickup: %s %s\n", context.OrderCreation.SchedulingInfo.PickupDate, context.OrderCreation.SchedulingInfo.PickupTime))
        summary.WriteString(fmt.Sprintf("- Delivery: %s %s\n", context.OrderCreation.SchedulingInfo.DeliveryDate, context.OrderCreation.SchedulingInfo.DeliveryTime))
    }
    
    return summary.String()
}
```

## Configuration

### **Environment Variables**
```bash
# GraphQL endpoint for order creation
GRAPHQL_ENDPOINT="https://api.dispatchit.com/graphql"

# Authentication
GRAPHQL_API_KEY="your-api-key"
GRAPHQL_AUTH_TOKEN="your-auth-token"
```

### **GraphQL Client Setup**
```go
// In main.go or initialization
func setupGraphQLClient() *order.OrderCreator {
    client := order.NewOrderCreator(os.Getenv("GRAPHQL_ENDPOINT"))
    
    // Set authentication headers
    client.SetHeader("Authorization", "Bearer "+os.Getenv("GRAPHQL_AUTH_TOKEN"))
    client.SetHeader("X-API-Key", os.Getenv("GRAPHQL_API_KEY"))
    
    return client
}
```

## Benefits of GraphQL Approach

### **1. Clean API Interface**
- **Single endpoint** for all order operations
- **Type-safe** queries and mutations
- **Flexible** data fetching

### **2. Better Error Handling**
- **Structured errors** from GraphQL
- **Field-level validation** errors
- **Clear error messages** for users

### **3. Efficient Data Transfer**
- **Request only needed data**
- **Single request** for complex operations
- **Reduced network overhead**

### **4. Future-Proof**
- **Easy to extend** with new fields
- **Backward compatible** with schema evolution
- **Standard GraphQL tooling** support

## Next Steps

1. **Implement GraphQL client** with proper authentication
2. **Add order creation service** with context conversion
3. **Integrate with conversation engine** for order creation
4. **Add error handling** and user feedback
5. **Test with real GraphQL endpoints**

This approach provides a clean, maintainable way to integrate with the ms-monolith order creation system while keeping our conversational interface separate and focused on user experience.
