# 📚 API Reference - Dispatch MCP Server

## Overview

The Dispatch MCP Server provides four main tools for delivery and pricing operations:

1. **`create_estimate`** - Create cost estimates for delivery orders
2. **`create_order`** - Create new delivery orders
3. **`compare_pricing_models`** - Compare different pricing models against estimates
4. **`select_delivery_option`** - Select appropriate delivery option based on customer scenario

## 🚚 Delivery Scenarios

The Dispatch MCP service supports two primary delivery scenarios that align with common customer needs:

### Scenario 1: "I need this delivered as soon as possible"
- **Use**: First option in the `availableOrderOptions` array
- **Characteristics**: Fastest delivery time, highest cost
- **Best for**: Urgent deliveries, time-sensitive items
- **Example**: Same-day delivery for urgent business documents

### Scenario 2: "I need this delivered sometime today"  
- **Use**: Last option in the `availableOrderOptions` array
- **Characteristics**: Slowest delivery time, lowest cost
- **Best for**: Non-urgent deliveries, cost-conscious customers
- **Example**: End-of-day delivery for regular packages

### Understanding the Response Format
The API returns multiple delivery options sorted by speed and cost:
- **First option**: Fastest delivery (most expensive)
- **Middle options**: Various time/cost trade-offs (can be ignored for most use cases)
- **Last option**: Slowest delivery (cheapest)

This design allows customers to choose between speed and cost based on their specific needs.

## 🔧 Tool Reference

### create_estimate

Creates a cost estimate for a delivery or service order.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pickup_info` | string | ✅ | JSON string containing pickup location information |
| `drop_offs` | string | ✅ | JSON string containing array of drop-off locations |
| `vehicle_type` | string | ✅ | Type of vehicle required (see Vehicle Types below) |
| `add_ons` | string | ❌ | JSON string containing optional add-ons array |
| `dedicated_vehicle` | string | ❌ | Whether dedicated vehicle is requested ("true"/"false") |
| `organization_druid` | string | ❌ | Organization ID for the request |

#### Vehicle Types

The `vehicle_type` parameter is **required** and must be one of the following values:

| Vehicle Type | Description | Best For |
|--------------|-------------|----------|
| `cargo_van` | Standard cargo van | Medium to large packages, furniture |
| `sprinter_van` | Large sprinter van | Large items, multiple packages |
| `box_truck` | Box truck | Very large items, bulk deliveries |
| `pickup_truck` | Pickup truck | Small to medium items, quick deliveries |

**Important**: Always ask the user what type of vehicle they need before calling the `create_estimate` tool. The vehicle type affects both pricing and delivery options.

#### Example Request

```json
{
  "tool": "create_estimate",
  "arguments": {
    "pickup_info": "{\"business_name\":\"Test Business\",\"location\":{\"address\":{\"street\":\"123 Main St\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94105\",\"country\":\"US\"}}}",
    "drop_offs": "[{\"business_name\":\"Drop Off Business\",\"location\":{\"address\":{\"street\":\"456 Oak Ave\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94110\",\"country\":\"US\"}}}]",
    "vehicle_type": "cargo_van",
    "add_ons": "[\"white_glove\", \"signature_required\"]",
    "dedicated_vehicle": "false"
  }
}
```

#### Response Format

The `/orders/estimates` response provides a list of delivery options where:

- **First option** = fastest delivery, most expensive
- **Last option** = slowest delivery, cheapest  
- **Everything in between** = various time/cost trade-offs

The system supports two key delivery scenarios:
1. **"I need this delivered as soon as possible"** - Use the first option (fastest)
2. **"I need this delivered sometime today"** - Use the last option (cheapest)

```json
{
  "data": {
    "createEstimate": {
      "estimate": {
        "availableOrderOptions": [
          {
            "serviceType": "delivery",
            "estimatedDeliveryTimeUtc": "2024-01-15T14:30:00Z",
            "estimatedOrderCost": 89.99,
            "vehicleType": "cargo_van",
            "pickupLocationInfo": {
              "googlePlaceId": "ChIJ...",
              "lat": 37.7749,
              "lng": -122.4194
            },
            "dropOffLocationsInfo": [
              {
                "googlePlaceId": "ChIJ...",
                "lat": 37.8044,
                "lng": -122.2712
              }
            ],
            "estimateInfo": {
              "serviceType": "delivery",
              "vehicleType": "cargo_van",
              "tollAmount": "5.50",
              "estimatedOrderCost": "89.99",
              "dedicatedVehicleRequested": false,
              "dedicatedVehicleFee": "0.00"
            },
            "addOns": ["white_glove", "signature_required"]
          },
          {
            "serviceType": "delivery",
            "estimatedDeliveryTimeUtc": "2024-01-15T16:30:00Z",
            "estimatedOrderCost": 65.99,
            "vehicleType": "cargo_van",
            "pickupLocationInfo": {
              "googlePlaceId": "ChIJ...",
              "lat": 37.7749,
              "lng": -122.4194
            },
            "dropOffLocationsInfo": [
              {
                "googlePlaceId": "ChIJ...",
                "lat": 37.8044,
                "lng": -122.2712
              }
            ],
            "estimateInfo": {
              "serviceType": "delivery",
              "vehicleType": "cargo_van",
              "tollAmount": "3.50",
              "estimatedOrderCost": "65.99",
              "dedicatedVehicleRequested": false,
              "dedicatedVehicleFee": "0.00"
            },
            "addOns": ["white_glove", "signature_required"]
          },
          {
            "serviceType": "delivery",
            "estimatedDeliveryTimeUtc": "2024-01-15T18:30:00Z",
            "estimatedOrderCost": 45.99,
            "vehicleType": "cargo_van",
            "pickupLocationInfo": {
              "googlePlaceId": "ChIJ...",
              "lat": 37.7749,
              "lng": -122.4194
            },
            "dropOffLocationsInfo": [
              {
                "googlePlaceId": "ChIJ...",
                "lat": 37.8044,
                "lng": -122.2712
              }
            ],
            "estimateInfo": {
              "serviceType": "delivery",
              "vehicleType": "cargo_van",
              "tollAmount": "2.50",
              "estimatedOrderCost": "45.99",
              "dedicatedVehicleRequested": false,
              "dedicatedVehicleFee": "0.00"
            },
            "addOns": ["white_glove", "signature_required"]
          }
        ]
      }
    }
  }
}
```

### create_order

Creates a new order for delivery or service.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `delivery_info` | string | ✅ | JSON string containing delivery information |
| `pickup_info` | string | ✅ | JSON string containing pickup information |
| `drop_offs` | string | ✅ | JSON string containing array of drop-off locations |
| `tags` | string | ❌ | JSON string containing optional order tags |

#### Example Request

```json
{
  "tool": "create_order",
  "arguments": {
    "delivery_info": "{\"service_type\":\"delivery\",\"organization_druid\":\"org_123\"}",
    "pickup_info": "{\"business_name\":\"Test Business\",\"contact_name\":\"John Doe\",\"contact_phone_number\":\"555-123-4567\",\"location\":{\"address\":{\"street\":\"123 Main St\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94105\",\"country\":\"US\"}}}",
    "drop_offs": "[{\"business_name\":\"Drop Off Business\",\"contact_name\":\"Jane Smith\",\"contact_phone_number\":\"555-987-6543\",\"location\":{\"address\":{\"street\":\"456 Oak Ave\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94110\",\"country\":\"US\"}}}]",
    "tags": "[{\"name\":\"priority\",\"value\":\"high\"}]"
  }
}
```

#### Response Format

```json
{
  "data": {
    "createOrder": {
      "order": {
        "id": "ORD-1234567890",
        "status": "pending",
        "scheduledAt": "2024-01-15T14:30:00Z",
        "totalCost": 45.99,
        "trackingNumber": "TRK-1234567890",
        "estimatedArrival": "2024-01-15T16:30:00Z"
      }
    }
  }
}
```

### select_delivery_option

Selects the appropriate delivery option from an estimate response based on the customer's delivery scenario.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `estimate_response` | string | ✅ | Full estimate response from create_estimate tool |
| `delivery_scenario` | string | ✅ | Delivery scenario: "fastest" for urgent delivery, "cheapest" for economy delivery |

#### Example Request

```json
{
  "tool": "select_delivery_option",
  "arguments": {
    "estimate_response": "{\"data\":{\"createEstimate\":{\"estimate\":{\"availableOrderOptions\":[{\"serviceType\":\"delivery\",\"estimatedDeliveryTimeUtc\":\"2024-01-15T14:30:00Z\",\"estimatedOrderCost\":89.99,\"vehicleType\":\"cargo_van\"},{\"serviceType\":\"delivery\",\"estimatedDeliveryTimeUtc\":\"2024-01-15T18:30:00Z\",\"estimatedOrderCost\":45.99,\"vehicleType\":\"cargo_van\"}]}}}}",
    "delivery_scenario": "fastest"
  }
}
```

#### Response Format

```json
{
  "selected_option": {
    "serviceType": "delivery",
    "estimatedDeliveryTimeUtc": "2024-01-15T14:30:00Z",
    "estimatedOrderCost": 89.99,
    "vehicleType": "cargo_van",
    "pickupLocationInfo": {
      "googlePlaceId": "ChIJ...",
      "lat": 37.7749,
      "lng": -122.4194
    },
    "dropOffLocationsInfo": [
      {
        "googlePlaceId": "ChIJ...",
        "lat": 37.8044,
        "lng": -122.2712
      }
    ],
    "estimateInfo": {
      "serviceType": "delivery",
      "vehicleType": "cargo_van",
      "tollAmount": "5.50",
      "estimatedOrderCost": "89.99",
      "dedicatedVehicleRequested": false,
      "dedicatedVehicleFee": "0.00"
    },
    "addOns": ["white_glove", "signature_required"]
  },
  "scenario": "fastest",
  "description": "Fastest delivery (most expensive)",
  "total_options": 2,
  "all_options": [
    {
      "serviceType": "delivery",
      "estimatedDeliveryTimeUtc": "2024-01-15T14:30:00Z",
      "estimatedOrderCost": 89.99,
      "vehicleType": "cargo_van"
    },
    {
      "serviceType": "delivery",
      "estimatedDeliveryTimeUtc": "2024-01-15T18:30:00Z",
      "estimatedOrderCost": 45.99,
      "vehicleType": "cargo_van"
    }
  ]
}
```

### compare_pricing_models

Compares different pricing models against an existing estimate to find the best pricing option.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `original_estimate` | string | ✅ | JSON string containing original estimate data |
| `delivery_count` | string | ❌ | Number of deliveries in the order (default: "1") |
| `customer_tier` | string | ❌ | Customer loyalty tier: "bronze", "silver", "gold" (default: "bronze") |
| `order_frequency` | string | ❌ | Number of orders per month (default: "1") |
| `total_order_value` | string | ❌ | Total value of the order (default: original cost) |
| `is_bulk_order` | string | ❌ | Whether this is a bulk order: "true"/"false" (default: "false") |

#### Example Request

```json
{
  "tool": "compare_pricing_models",
  "arguments": {
    "original_estimate": "{\"serviceType\":\"delivery\",\"estimatedOrderCost\":45.99,\"vehicleType\":\"cargo_van\",\"estimatedDeliveryTimeUtc\":\"2024-01-15T14:30:00Z\"}",
    "delivery_count": "3",
    "customer_tier": "gold",
    "order_frequency": "5",
    "total_order_value": "150.00",
    "is_bulk_order": "false"
  }
}
```

#### Response Format

```json
{
  "original_estimate": {
    "serviceType": "delivery",
    "estimatedOrderCost": 45.99,
    "vehicleType": "cargo_van",
    "estimatedDeliveryTimeUtc": "2024-01-15T14:30:00Z"
  },
  "pricing_models": [
    {
      "model": "standard",
      "name": "Standard Pricing",
      "original_cost": 45.99,
      "adjusted_cost": 45.99,
      "discount": 0.0,
      "discount_percent": 0.0,
      "savings": 0.0,
      "eligible": true
    },
    {
      "model": "multi_delivery",
      "name": "Multi-Delivery Discount",
      "original_cost": 45.99,
      "adjusted_cost": 39.09,
      "discount": 6.90,
      "discount_percent": 15.0,
      "savings": 6.90,
      "eligible": true
    },
    {
      "model": "volume_discount",
      "name": "Volume Discount",
      "original_cost": 45.99,
      "adjusted_cost": 36.79,
      "discount": 9.20,
      "discount_percent": 20.0,
      "savings": 9.20,
      "eligible": true
    },
    {
      "model": "loyalty_discount",
      "name": "Loyalty Discount",
      "original_cost": 45.99,
      "adjusted_cost": 41.39,
      "discount": 4.60,
      "discount_percent": 10.0,
      "savings": 4.60,
      "eligible": true
    },
    {
      "model": "bulk_order",
      "name": "Bulk Order Discount",
      "original_cost": 45.99,
      "adjusted_cost": 45.99,
      "discount": 0.0,
      "discount_percent": 0.0,
      "savings": 0.0,
      "eligible": false,
      "reason": "Requires bulk order with 10+ deliveries, you have 3"
    }
  ],
  "best_option": {
    "model": "volume_discount",
    "name": "Volume Discount",
    "original_cost": 45.99,
    "adjusted_cost": 36.79,
    "discount": 9.20,
    "discount_percent": 20.0,
    "savings": 9.20,
    "eligible": true
  },
  "savings": 9.20,
  "savings_percentage": 20.0
}
```

## 📊 Data Types

### PricingModel Enum

```go
type PricingModel string

const (
    StandardPricing     PricingModel = "standard"
    MultiDeliveryPricing PricingModel = "multi_delivery"
    VolumeDiscountPricing PricingModel = "volume_discount"
    LoyaltyDiscountPricing PricingModel = "loyalty_discount"
    BulkOrderPricing    PricingModel = "bulk_order"
)
```

### PricingContext

```go
type PricingContext struct {
    DeliveryCount    int     `json:"delivery_count"`
    CustomerTier     string  `json:"customer_tier"`
    OrderFrequency   int     `json:"order_frequency"`
    TotalOrderValue  float64 `json:"total_order_value"`
    IsBulkOrder      bool    `json:"is_bulk_order"`
    OrganizationDruid string `json:"organization_druid"`
}
```

### PricingResult

```go
type PricingResult struct {
    Model           PricingModel `json:"model"`
    Name            string       `json:"name"`
    OriginalCost    float64      `json:"original_cost"`
    AdjustedCost    float64      `json:"adjusted_cost"`
    Discount        float64      `json:"discount"`
    DiscountPercent float64      `json:"discount_percent"`
    Savings         float64      `json:"savings"`
    Eligible        bool         `json:"eligible"`
    Reason          string       `json:"reason,omitempty"`
}
```

## ✅ Input Validation

The MCP server includes comprehensive input validation to ensure data quality and provide clear error messages:

### Validation Features

- **Vehicle Type Validation**: Ensures only valid vehicle types are accepted
- **Delivery Scenario Validation**: Validates delivery scenario parameters
- **Address Validation**: Validates address format and required fields
- **JSON Format Validation**: Ensures all JSON parameters are properly formatted
- **Numeric Validation**: Validates numeric parameters with range checking
- **Boolean Validation**: Ensures boolean parameters are properly formatted
- **Customer Tier Validation**: Validates customer loyalty tiers

### Validation Error Format

When validation fails, the server returns detailed error messages:

```json
{
  "error": "vehicle_type validation failed: Invalid vehicle type 'invalid'. Must be one of: pickup_truck, cargo_van, sprinter_van, box_truck"
}
```

### Supported Validation Rules

| Parameter | Validation Rules |
|-----------|------------------|
| `vehicle_type` | Must be one of: pickup_truck, cargo_van, sprinter_van, box_truck |
| `delivery_scenario` | Must be one of: fastest, asap, urgent, cheapest, economy, sometime_today |
| `customer_tier` | Must be one of: bronze, silver, gold (optional) |
| `delivery_count` | Must be a positive integer between 1-100 (optional) |
| `order_frequency` | Must be a positive integer between 1-100 (optional) |
| `is_bulk_order` | Must be 'true' or 'false' (optional) |
| `pickup_info` | Must be valid JSON with required address fields |
| `drop_offs` | Must be valid JSON array with at least one location |
| `original_estimate` | Must be valid JSON with estimate data |

## 🔍 Error Handling

### Common Error Responses

#### Invalid Arguments Format
```json
{
  "error": "invalid arguments format"
}
```

#### Missing Required Parameter
```json
{
  "error": "original_estimate is required and must be a string"
}
```

#### JSON Parsing Error
```json
{
  "error": "failed to parse original_estimate: invalid character 'x' looking for beginning of value"
}
```

#### API Call Failure
```json
{
  "error": "failed to create estimate: authentication failed"
}
```

#### Validation Errors
```json
{
  "error": "vehicle_type validation failed: Invalid vehicle type 'invalid'. Must be one of: pickup_truck, cargo_van, sprinter_van, box_truck"
}
```

```json
{
  "error": "delivery_scenario validation failed: Invalid delivery scenario 'slow'. Must be one of: fastest, asap, urgent, cheapest, economy, sometime_today"
}
```

```json
{
  "error": "pickup_info validation failed: Invalid JSON format - invalid character 'x' looking for beginning of value"
}
```

## 🚀 Usage Patterns

### 1. Getting Vehicle Type Information

Before creating an estimate, you must collect the vehicle type from the user. Here's a suggested approach:

**Ask the user**: "What type of vehicle do you need for this delivery?"
- For small to medium packages: `pickup_truck`
- For medium to large packages: `cargo_van` 
- For large items or multiple packages: `sprinter_van`
- For very large items or bulk deliveries: `box_truck`

**Example conversation flow**:
```
User: "I need to deliver a couch from my store to a customer"
AI: "What type of vehicle do you need for this delivery? 
     - pickup_truck (small to medium items)
     - cargo_van (medium to large packages) 
     - sprinter_van (large items, multiple packages)
     - box_truck (very large items, bulk deliveries)"
User: "A cargo van should work"
AI: [Proceeds with create_estimate using vehicle_type: "cargo_van"]
```

### 2. Complete Information Gathering Workflow

Before calling `create_estimate`, you must collect all required information from the user:

1. **Pickup Information**: Business name, address, contact details
2. **Drop-off Information**: Destination address(es), contact details  
3. **Vehicle Type**: Ask "What type of vehicle do you need?" and provide options
4. **Delivery Scenario**: Ask "Do you need this delivered as soon as possible or is economy delivery okay?"

**Example conversation**:
```
User: "I need to deliver a dining table from my store to a customer"
AI: "I'll help you create a delivery estimate. Let me gather some information:

1. What's your store address and business name?
2. What's the customer's delivery address?
3. What type of vehicle do you need for this delivery?
   - pickup_truck (small to medium items)
   - cargo_van (medium to large packages) 
   - sprinter_van (large items, multiple packages)
   - box_truck (very large items, bulk deliveries)
4. Do you need this delivered as soon as possible or is economy delivery okay?"

User: "My store is at 123 Main St, San Francisco. Customer is at 456 Oak Ave, Oakland. I need a cargo van, and I need it delivered as soon as possible."

AI: [Calls create_estimate with all the information, then uses select_delivery_option with scenario: "fastest"]
```

### 3. Basic Estimate Creation
```json
{
  "tool": "create_estimate",
  "arguments": {
    "pickup_info": "{\"business_name\":\"My Business\",\"location\":{\"address\":{\"street\":\"123 Main St\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94105\",\"country\":\"US\"}}}",
    "drop_offs": "[{\"business_name\":\"Customer\",\"location\":{\"address\":{\"street\":\"456 Oak Ave\",\"city\":\"Oakland\",\"state\":\"CA\",\"zip_code\":\"94610\",\"country\":\"US\"}}}]",
    "vehicle_type": "cargo_van"
  }
}
```

### 4. Pricing Comparison Workflow
```json
// Step 1: Create estimate
{
  "tool": "create_estimate",
  "arguments": { /* estimate parameters */ }
}

// Step 2: Compare pricing models
{
  "tool": "compare_pricing_models",
  "arguments": {
    "original_estimate": "/* result from step 1 */",
    "delivery_count": "3",
    "customer_tier": "gold"
  }
}
```

### 5. Delivery Scenario Selection
```json
// Step 1: Create estimate
{
  "tool": "create_estimate",
  "arguments": { /* estimate parameters */ }
}

// Step 2: Select delivery option based on scenario
{
  "tool": "select_delivery_option",
  "arguments": {
    "estimate_response": "/* result from step 1 */",
    "delivery_scenario": "fastest"  // or "cheapest"
  }
}
```

### 6. Order Creation with Best Pricing
```json
{
  "tool": "create_order",
  "arguments": {
    "delivery_info": "{\"service_type\":\"delivery\"}",
    "pickup_info": "/* pickup info */",
    "drop_offs": "/* drop-off info */"
  }
}
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `USE_IDP_AUTH` | Use IDP authentication | `false` |
| `DISPATCH_AUTH_TOKEN` | Static auth token | - |
| `DISPATCH_ORGANIZATION_ID` | Organization ID | - |
| `DISPATCH_GRAPHQL_ENDPOINT` | GraphQL endpoint | `https://graphql-gateway.monkey.dispatchfog.org/graphql` |

### IDP Authentication Variables

| Variable | Description |
|----------|-------------|
| `IDP_ENDPOINT` | IDP endpoint URL |
| `IDP_CLIENT_ID` | Client ID for OAuth |
| `IDP_CLIENT_SECRET` | Client secret for OAuth |
| `IDP_SCOPE` | OAuth scope |
| `IDP_TOKEN_ENDPOINT` | Token endpoint URL |

## 📈 Performance Considerations

- **Estimate Creation**: ~500ms average response time
- **Pricing Comparison**: ~50ms average response time
- **Order Creation**: ~1s average response time
- **Concurrent Requests**: Supports up to 100 concurrent requests
- **Rate Limiting**: 1000 requests per hour per organization

## 🔒 Security

- **Authentication**: OAuth 2.0 with IDP or static tokens
- **Authorization**: Organization-based access control
- **Data Encryption**: All data encrypted in transit
- **Token Management**: Automatic token refresh
- **Audit Logging**: All operations logged for compliance

---

*Last updated: $(date)*
