# API Endpoints Analysis for Buildathon Project

## Overview

This document analyzes the available API endpoints for creating orders and estimates in the Dispatch system, comparing GraphQL and REST approaches for the MCP server implementation.

## GraphQL Endpoints

### 1. Create Estimate Mutation

**Endpoint**: GraphQL Mutation `createEstimate`
**File**: `/app/graphql/mutations/order/create_estimate.rb`

**Input Parameters**:
```graphql
input CreateEstimateInput {
  add_ons: [String]                    # Optional add-ons for delivery
  dedicated_vehicle: Boolean           # Whether dedicated vehicle is requested
  drop_offs: [DropOffInfoInput!]!     # Required: Drop-off location information
  drop_off_date_time_utc: DateTime    # Optional: Drop-off date/time in UTC
  organization_druid: ID              # Optional: Organization ID
  pickup_info: PickupInfoInput!       # Required: Pickup location information
  vehicle_type: String!               # Required: Type of vehicle needed
}
```

**Response Structure**:
```graphql
type Estimate {
  available_order_options: [AvailableOrderOption]
}

type AvailableOrderOption {
  service_type: String
  estimated_delivery_time_utc: DateTime
  estimated_order_cost: Float
  vehicle_type: String
  pickup_location_info: LocationInfo
  drop_off_locations_info: [LocationInfo]
  estimate_info: EstimateInfo
  add_ons: [String]
}
```

**Authentication**: Requires authenticated user
**Analytics**: Tracks "GQL Order Estimate" events

### 2. Create Order Mutation

**Endpoint**: GraphQL Mutation `createOrder`
**File**: `/app/graphql/mutations/order/create_order.rb`

**Input Parameters**:
```graphql
input CreateOrderInput {
  add_ons: [String]                           # Optional add-ons
  delivery_info: DeliveryInfoInput!            # Required: Delivery information
  drop_offs: [CreateOrderDropOffInfoInput!]!  # Required: Drop-off locations
  pickup_info: CreateOrderPickupInfoInput!   # Required: Pickup information
  tags: [TagsInput]                           # Optional: Order tags
}
```

**Response Structure**:
```graphql
type CreateOrderResponse {
  order: Order!
}
```

**Authentication**: Requires authenticated user
**Business Logic**: 
- Validates organization access
- Handles service type auto-changes
- Creates unassigned multi-stop orders
- Processes credit transactions for discounts

## REST API Endpoints

### 1. Order Estimates Controller

**Endpoint**: `POST /order_estimates/estimate`
**File**: `/app/controllers/order_estimates_controller.rb`

**Input Parameters**:
```ruby
{
  pickup_lat: Float,
  pickup_lng: Float,
  dropoff_lat: Float,
  dropoff_lng: Float,
  drop_off_at_time: String,
  needed_vehicle_type: String,
  scheduled_at_picker: String,
  scheduled_at_date: String,
  scheduled_at_time: String,
  pickup_at_right_away: Boolean,
  pickup_state: String,
  pickup_postal_code: String,
  dropoff_postal_code: String
}
```

**Features**:
- No authentication required (public endpoint)
- Tracks "Quick Estimate" analytics events
- Returns JSON response with pricing information

### 2. External Orders API (v1)

**Endpoint**: `GET /api/external/v1/orders`
**File**: `/app/controllers/api/external/v1/orders_controller.rb`

**Capabilities**:
- List orders with pagination (default 50 per page)
- Filter by status and updated_after date
- Expand order details
- Get order events/history

**Authentication**: Requires API authentication
**Features**:
- Feature flag controlled (`api_external_list_orders_endpoint`)
- Policy-based authorization
- Includes deliveries and capabilities

## Comparison: GraphQL vs REST

### GraphQL Advantages

1. **Type Safety**: Strong typing with input/output schemas
2. **Single Endpoint**: One endpoint for all operations
3. **Flexible Queries**: Request only needed fields
4. **Built-in Validation**: Schema validation built-in
5. **Rich Tooling**: GraphQL playground, introspection
6. **Consistent Authentication**: Unified auth model

### GraphQL Disadvantages

1. **Complexity**: More complex for simple operations
2. **Learning Curve**: Team needs GraphQL knowledge
3. **Caching**: More complex caching strategies
4. **Over-fetching Risk**: Can request too much data

### REST Advantages

1. **Simplicity**: Simple HTTP requests
2. **Familiar**: Standard HTTP methods and status codes
3. **Caching**: Standard HTTP caching works well
4. **Tooling**: Standard HTTP tools and libraries
5. **Public Endpoints**: Some endpoints don't require auth

### REST Disadvantages

1. **Multiple Endpoints**: Need to know specific endpoints
2. **Over-fetching**: Always get full response
3. **Versioning**: Need to version endpoints
4. **Inconsistent**: Different patterns across endpoints

## Recommendation for MCP Server

### Primary Choice: GraphQL

**Rationale**:
1. **Type Safety**: Better for MCP tool parameter validation
2. **Single Endpoint**: Simpler MCP server configuration
3. **Rich Schema**: Better documentation and introspection
4. **Consistent Auth**: Unified authentication model
5. **Future-Proof**: More extensible for additional operations

### Implementation Strategy

1. **Start with GraphQL**: Use `createEstimate` and `createOrder` mutations
2. **Fallback to REST**: Keep REST endpoints as backup for simple operations
3. **Hybrid Approach**: Use GraphQL for complex operations, REST for simple ones

### MCP Server Design

```go
// MCP Tools for GraphQL
type CreateEstimateTool struct {
    client *GraphQLClient
}

type CreateOrderTool struct {
    client *GraphQLClient
}

// MCP Tools for REST (backup)
type QuickEstimateTool struct {
    client *HTTPClient
}
```

## Authentication Considerations

### GraphQL Authentication
- Requires authenticated user session
- Uses `current_user` context
- Organization-based authorization
- Policy-based access control

### REST Authentication
- Some endpoints are public (estimates)
- API key authentication for external endpoints
- Feature flag controlled access

### MCP Server Auth Strategy
1. **Environment Variables**: Store auth tokens in env vars (MVP approach)
2. **Session Management**: Handle user sessions if needed
3. **Organization Context**: Pass organization context for multi-tenant access

## Next Steps

1. **Test GraphQL Endpoints**: Use GraphQL playground to test mutations
2. **Understand Input Types**: Review the input type definitions
3. **Authentication Research**: Determine auth token requirements
4. **Error Handling**: Understand error response formats
5. **Rate Limits**: Check for any rate limiting

## Files to Review Further

1. **Input Type Definitions**:
   - `/app/graphql/types/create_estimate/request/`
   - `/app/graphql/types/create_order/request/`

2. **Response Type Definitions**:
   - `/app/graphql/types/create_estimate/response/`
   - `/app/graphql/types/order_type.rb`

3. **Business Logic**:
   - `/app/models/order_estimate_entry.rb`
   - `/app/models/order_entry.rb`
   - `/app/interactors/create_unassigned_multi_stop_order.rb`

4. **Authentication**:
   - `/app/controllers/application_controller.rb`
   - `/app/policies/`

---
*This analysis will be updated as we gain more information about the actual API usage patterns and requirements.*
