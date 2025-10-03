# AI Agent Implementation Guide: Dispatch MCP Server

## Project Overview
Build a Model Context Protocol (MCP) server in Go that provides conversational AI access to Dispatch's order creation and estimate APIs. This server will expose two main tools: `create_estimate` and `create_order`.

## Architecture Pattern
Based on Chris's Jira MCP implementation, use the `mcp-go` framework with the following structure:

```go
// Main server structure
type MCPServer struct {
    dispatchClient *resty.Client
    config        *Config
}

// Configuration
type Config struct {
    GraphQLEndpoint string
    AuthToken       string
    OrganizationID   string
}
```

## Step 1: Project Setup

### 1.1 Initialize Go Module
```bash
mkdir dispatch-mcp-server
cd dispatch-mcp-server
go mod init dispatch-mcp-server
```

### 1.2 Add Dependencies
```go
// go.mod
module dispatch-mcp-server

go 1.21

require (
    github.com/go-resty/resty/v2 v2.10.0
    github.com/mark3labs/mcp-go v0.1.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### 1.3 Project Structure
```
dispatch-mcp-server/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── mcp/
│   │   ├── server.go
│   │   └── tools.go
│   ├── dispatch/
│   │   ├── client.go
│   │   └── types.go
│   └── config/
│       └── config.go
├── go.mod
├── go.sum
└── README.md
```

## Step 2: Core MCP Server Implementation

### 2.1 Main Server File (`cmd/server/main.go`)
```go
package main

import (
    "fmt"
    "log"
    "os"
    "dispatch-mcp-server/internal/mcp"
)

func main() {
    server, err := mcp.NewMCPServer()
    if err != nil {
        log.Fatalf("Failed to create MCP server: %v", err)
    }

    fmt.Fprintf(os.Stderr, "Starting Dispatch MCP server...\n")
    if err := server.Run(); err != nil {
        log.Fatalf("MCP server error: %v", err)
    }
}
```

### 2.2 MCP Server Implementation (`internal/mcp/server.go`)
```go
package mcp

import (
    "fmt"
    "log"
    "os"
    "dispatch-mcp-server/internal/dispatch"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
    dispatchClient *dispatch.Client
}

func NewMCPServer() (*MCPServer, error) {
    dispatchClient, err := dispatch.NewClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create dispatch client: %v", err)
    }

    return &MCPServer{
        dispatchClient: dispatchClient,
    }, nil
}

func (s *MCPServer) Run() error {
    fmt.Fprintf(os.Stderr, "Starting Dispatch MCP server...\n")

    srv := server.NewMCPServer(
        "dispatch-mcp-server",
        "1.0.0",
    )

    // Register create_estimate tool
    estimateTool := mcp.NewTool("create_estimate",
        mcp.WithDescription("Create a cost estimate for a delivery or service order"),
        mcp.WithString("pickup_info", mcp.Required(), mcp.Description("Pickup location information")),
        mcp.WithString("drop_offs", mcp.Required(), mcp.Description("Drop-off locations array")),
        mcp.WithString("vehicle_type", mcp.Required(), mcp.Description("Type of vehicle required")),
        mcp.WithString("add_ons", mcp.Description("Optional add-ons for delivery")),
        mcp.WithString("dedicated_vehicle", mcp.Description("Whether dedicated vehicle is requested")),
        mcp.WithString("organization_druid", mcp.Description("Organization ID")),
    )

    srv.AddTool(estimateTool, s.createEstimateTool)

    // Register create_order tool
    orderTool := mcp.NewTool("create_order",
        mcp.WithDescription("Create a new order for delivery or service"),
        mcp.WithString("delivery_info", mcp.Required(), mcp.Description("Delivery information")),
        mcp.WithString("pickup_info", mcp.Required(), mcp.Description("Pickup information")),
        mcp.WithString("drop_offs", mcp.Required(), mcp.Description("Drop-off locations array")),
        mcp.WithString("tags", mcp.Description("Optional order tags")),
    )

    srv.AddTool(orderTool, s.createOrderTool)

    fmt.Fprintf(os.Stderr, "MCP server initialized and listening...\n")

    if err := server.ServeStdio(srv); err != nil {
        log.Fatalf("Server error: %v", err)
    }
    return nil
}
```

## Step 3: Dispatch API Client Implementation

### 3.1 Dispatch Client (`internal/dispatch/client.go`)
```go
package dispatch

import (
    "bytes"
    "encoding/json"
    "fmt"
    "dispatch-mcp-server/internal/config"
    "github.com/go-resty/resty/v2"
)

type Client struct {
    client *resty.Client
    config *config.Config
}

func NewClient() (*Client, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    client := resty.New()
    client.SetBaseURL(cfg.GraphQLEndpoint)
    client.SetHeader("Authorization", "Bearer "+cfg.AuthToken)
    client.SetHeader("Content-Type", "application/json")

    return &Client{
        client: client,
        config: cfg,
    }, nil
}

func (c *Client) CreateEstimate(input CreateEstimateInput) (*CreateEstimateResponse, error) {
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

    resp, err := c.client.R().
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

    resp, err := c.client.R().
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
```

### 3.2 Type Definitions (`internal/dispatch/types.go`)
```go
package dispatch

// CreateEstimateInput represents the input for creating an estimate
type CreateEstimateInput struct {
    AddOns              []string                `json:"add_ons,omitempty"`
    DedicatedVehicle    *bool                  `json:"dedicated_vehicle,omitempty"`
    DropOffs            []DropOffInfoInput     `json:"drop_offs"`
    DropOffDateTimeUTC  *string                `json:"drop_off_date_time_utc,omitempty"`
    OrganizationDruid   *string                `json:"organization_druid,omitempty"`
    PickupInfo          PickupInfoInput        `json:"pickup_info"`
    VehicleType         string                 `json:"vehicle_type"`
}

type PickupInfoInput struct {
    BusinessName        string        `json:"business_name"`
    Location           LocationInput `json:"location"`
    PickupDateTimeUTC  *string       `json:"pickup_date_time_utc,omitempty"`
}

type DropOffInfoInput struct {
    BusinessName    string        `json:"business_name"`
    EstimatedWeight *int          `json:"estimated_weight,omitempty"`
    Location        LocationInput `json:"location"`
}

type LocationInput struct {
    Address         *AddressInput      `json:"address,omitempty"`
    GeoCoordinates  *GeoCoordinatesInput `json:"geo_coordinates,omitempty"`
}

type AddressInput struct {
    Street    string `json:"street"`
    City      string `json:"city"`
    State     string `json:"state"`
    ZipCode   string `json:"zip_code"`
    Country   string `json:"country"`
}

type GeoCoordinatesInput struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
    AddOns        []string                `json:"add_ons,omitempty"`
    DeliveryInfo  DeliveryInfoInput       `json:"delivery_info"`
    DropOffs       []CreateOrderDropOffInfoInput `json:"drop_offs"`
    PickupInfo     CreateOrderPickupInfoInput    `json:"pickup_info"`
    Tags          []TagInput              `json:"tags,omitempty"`
}

type DeliveryInfoInput struct {
    ServiceType        string  `json:"service_type"`
    OrganizationDruid  *string `json:"organization_druid,omitempty"`
}

type CreateOrderPickupInfoInput struct {
    BusinessName       *string       `json:"business_name,omitempty"`
    ContactName        *string       `json:"contact_name,omitempty"`
    ContactPhoneNumber *string       `json:"contact_phone_number,omitempty"`
    Location           *LocationInput `json:"location,omitempty"`
    PickupNotes        *string       `json:"pickup_notes,omitempty"`
}

type CreateOrderDropOffInfoInput struct {
    BusinessName       *string       `json:"business_name,omitempty"`
    ContactName        *string       `json:"contact_name,omitempty"`
    ContactPhoneNumber *string       `json:"contact_phone_number,omitempty"`
    Location           *LocationInput `json:"location,omitempty"`
    DropOffNotes       *string       `json:"drop_off_notes,omitempty"`
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
    ServiceType              string                 `json:"serviceType"`
    EstimatedDeliveryTimeUTC string                 `json:"estimatedDeliveryTimeUtc"`
    EstimatedOrderCost       float64                `json:"estimatedOrderCost"`
    VehicleType              string                 `json:"vehicleType"`
    PickupLocationInfo       LocationInfo           `json:"pickupLocationInfo"`
    DropOffLocationsInfo     []LocationInfo         `json:"dropOffLocationsInfo"`
    EstimateInfo             EstimateInfo           `json:"estimateInfo"`
    AddOns                   []string               `json:"addOns"`
}

type LocationInfo struct {
    GooglePlaceID string  `json:"googlePlaceId"`
    Lat           float64 `json:"lat"`
    Lng           float64 `json:"lng"`
}

type EstimateInfo struct {
    ServiceType                string  `json:"serviceType"`
    VehicleType                string  `json:"vehicleType"`
    TollAmount                 string  `json:"tollAmount"`
    EstimatedOrderCost         string  `json:"estimatedOrderCost"`
    DedicatedVehicleRequested *bool   `json:"dedicatedVehicleRequested"`
    DedicatedVehicleFee       string  `json:"dedicatedVehicleFee"`
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
```

### 3.3 Configuration (`internal/config/config.go`)
```go
package config

import (
    "os"
)

type Config struct {
    GraphQLEndpoint string
    AuthToken       string
    OrganizationID   string
}

func Load() (*Config, error) {
    return &Config{
        GraphQLEndpoint: getEnv("DISPATCH_GRAPHQL_ENDPOINT", "https://api.dispatchit.com/graphql"),
        AuthToken:       getEnv("DISPATCH_AUTH_TOKEN", ""),
        OrganizationID:   getEnv("DISPATCH_ORGANIZATION_ID", ""),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

## Step 4: MCP Tool Implementations

### 4.1 Create Estimate Tool (`internal/mcp/tools.go`)
```go
package mcp

import (
    "encoding/json"
    "fmt"
    "dispatch-mcp-server/internal/dispatch"
    "github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) createEstimateTool(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
    // Parse pickup_info
    pickupInfoRaw, ok := arguments["pickup_info"].(string)
    if !ok {
        return mcp.NewToolResultError("pickup_info is required and must be a string"), nil
    }
    
    var pickupInfo dispatch.PickupInfoInput
    if err := json.Unmarshal([]byte(pickupInfoRaw), &pickupInfo); err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to parse pickup_info: %v", err)), nil
    }

    // Parse drop_offs
    dropOffsRaw, ok := arguments["drop_offs"].(string)
    if !ok {
        return mcp.NewToolResultError("drop_offs is required and must be a string"), nil
    }
    
    var dropOffs []dispatch.DropOffInfoInput
    if err := json.Unmarshal([]byte(dropOffsRaw), &dropOffs); err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to parse drop_offs: %v", err)), nil
    }

    // Build input
    input := dispatch.CreateEstimateInput{
        PickupInfo: pickupInfo,
        DropOffs:   dropOffs,
        VehicleType: getStringArg(arguments, "vehicle_type"),
    }

    // Optional fields
    if addOns, ok := arguments["add_ons"].(string); ok && addOns != "" {
        var addOnsList []string
        if err := json.Unmarshal([]byte(addOns), &addOnsList); err == nil {
            input.AddOns = addOnsList
        }
    }

    if dedicatedVehicle, ok := arguments["dedicated_vehicle"].(string); ok && dedicatedVehicle != "" {
        if dedicatedVehicle == "true" {
            input.DedicatedVehicle = &[]bool{true}[0]
        } else if dedicatedVehicle == "false" {
            input.DedicatedVehicle = &[]bool{false}[0]
        }
    }

    if orgDruid, ok := arguments["organization_druid"].(string); ok && orgDruid != "" {
        input.OrganizationDruid = &orgDruid
    }

    // Call API
    response, err := s.dispatchClient.CreateEstimate(input)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to create estimate: %v", err)), nil
    }

    // Format response
    responseJSON, _ := json.MarshalIndent(response, "", "  ")
    return mcp.NewToolResultText(string(responseJSON)), nil
}

func (s *MCPServer) createOrderTool(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
    // Parse delivery_info
    deliveryInfoRaw, ok := arguments["delivery_info"].(string)
    if !ok {
        return mcp.NewToolResultError("delivery_info is required and must be a string"), nil
    }
    
    var deliveryInfo dispatch.DeliveryInfoInput
    if err := json.Unmarshal([]byte(deliveryInfoRaw), &deliveryInfo); err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to parse delivery_info: %v", err)), nil
    }

    // Parse pickup_info
    pickupInfoRaw, ok := arguments["pickup_info"].(string)
    if !ok {
        return mcp.NewToolResultError("pickup_info is required and must be a string"), nil
    }
    
    var pickupInfo dispatch.CreateOrderPickupInfoInput
    if err := json.Unmarshal([]byte(pickupInfoRaw), &pickupInfo); err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to parse pickup_info: %v", err)), nil
    }

    // Parse drop_offs
    dropOffsRaw, ok := arguments["drop_offs"].(string)
    if !ok {
        return mcp.NewToolResultError("drop_offs is required and must be a string"), nil
    }
    
    var dropOffs []dispatch.CreateOrderDropOffInfoInput
    if err := json.Unmarshal([]byte(dropOffsRaw), &dropOffs); err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to parse drop_offs: %v", err)), nil
    }

    // Build input
    input := dispatch.CreateOrderInput{
        DeliveryInfo: deliveryInfo,
        PickupInfo:   pickupInfo,
        DropOffs:     dropOffs,
    }

    // Optional fields
    if tags, ok := arguments["tags"].(string); ok && tags != "" {
        var tagsList []dispatch.TagInput
        if err := json.Unmarshal([]byte(tags), &tagsList); err == nil {
            input.Tags = tagsList
        }
    }

    // Call API
    response, err := s.dispatchClient.CreateOrder(input)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("failed to create order: %v", err)), nil
    }

    // Format response
    responseJSON, _ := json.MarshalIndent(response, "", "  ")
    return mcp.NewToolResultText(string(responseJSON)), nil
}

func getStringArg(arguments map[string]interface{}, key string) string {
    if value, ok := arguments[key].(string); ok {
        return value
    }
    return ""
}
```

## Step 5: Environment Configuration

### 5.1 Environment Variables
Create a `.env` file:
```bash
DISPATCH_GRAPHQL_ENDPOINT=https://api.dispatchit.com/graphql
DISPATCH_AUTH_TOKEN=your_auth_token_here
DISPATCH_ORGANIZATION_ID=your_org_id_here
```

### 5.2 Build Script
Create `build.sh`:
```bash
#!/bin/bash
go build -o bin/dispatch-mcp-server ./cmd/server
```

## Step 6: Testing

### 6.1 Test File (`test/mcp_test.go`)
```go
package test

import (
    "encoding/json"
    "os"
    "os/exec"
    "testing"
    "time"
)

func testMCPTool(t *testing.T, tool string, args map[string]interface{}) {
    request := map[string]interface{}{
        "jsonrpc": "2.0",
        "id":      1,
        "method":  "tools/call",
        "params": map[string]interface{}{
            "name":      tool,
            "arguments": args,
        },
    }

    reqBytes, _ := json.Marshal(request)
    reqStr := string(reqBytes) + "\n"

    cmd := exec.Command("../bin/dispatch-mcp-server")
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    
    if err := cmd.Start(); err != nil {
        t.Fatalf("Failed to start MCP server: %v", err)
    }
    defer cmd.Process.Kill()

    time.Sleep(100 * time.Millisecond)

    stdin.Write([]byte(reqStr))
    stdin.Close()

    responseChan := make(chan string, 1)
    go func() {
        output := make([]byte, 8192)
        if n, err := stdout.Read(output); err == nil {
            responseChan <- string(output[:n])
        }
    }()

    select {
    case response := <-responseChan:
        t.Logf("MCP Response: %s", response)
    case <-time.After(10 * time.Second):
        t.Log("MCP server timed out")
    }
}

func TestCreateEstimate(t *testing.T) {
    if os.Getenv("DISPATCH_AUTH_TOKEN") == "" {
        t.Skip("Skipping - Dispatch credentials not available")
    }

    pickupInfo := map[string]interface{}{
        "business_name": "Test Business",
        "location": map[string]interface{}{
            "address": map[string]interface{}{
                "street":  "123 Main St",
                "city":    "San Francisco",
                "state":   "CA",
                "zip_code": "94105",
                "country":  "US",
            },
        },
    }

    dropOffs := []map[string]interface{}{
        {
            "business_name": "Drop Off Business",
            "location": map[string]interface{}{
                "address": map[string]interface{}{
                    "street":  "456 Oak Ave",
                    "city":    "San Francisco",
                    "state":   "CA",
                    "zip_code": "94110",
                    "country":  "US",
                },
            },
        },
    }

    args := map[string]interface{}{
        "pickup_info":  jsonString(pickupInfo),
        "drop_offs":    jsonString(dropOffs),
        "vehicle_type": "cargo_van",
    }

    testMCPTool(t, "create_estimate", args)
}

func jsonString(v interface{}) string {
    b, _ := json.Marshal(v)
    return string(b)
}
```

## Questions for AI Agent

When implementing this, the AI agent should ask these questions:

1. **Authentication**: "What type of authentication does Dispatch use? (API key, OAuth, session token?)"

2. **GraphQL Endpoint**: "What is the exact GraphQL endpoint URL for Dispatch?"

3. **Organization Context**: "How do we determine which organization to use for orders?"

4. **Error Handling**: "What are the common error responses from the Dispatch API?"

5. **Rate Limits**: "Are there any rate limits or special requirements for the API?"

6. **Testing**: "Do you have test credentials or a sandbox environment for testing?"

## Implementation Checklist

- [ ] Set up Go module and dependencies
- [ ] Implement MCP server structure
- [ ] Create Dispatch API client
- [ ] Implement create_estimate tool
- [ ] Implement create_order tool
- [ ] Add error handling and validation
- [ ] Create configuration management
- [ ] Write tests
- [ ] Build and test the server
- [ ] Document usage examples

## Usage Examples

### Create Estimate Example
```json
{
  "tool": "create_estimate",
  "arguments": {
    "pickup_info": "{\"business_name\":\"Test Business\",\"location\":{\"address\":{\"street\":\"123 Main St\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94105\",\"country\":\"US\"}}}",
    "drop_offs": "[{\"business_name\":\"Drop Off Business\",\"location\":{\"address\":{\"street\":\"456 Oak Ave\",\"city\":\"San Francisco\",\"state\":\"CA\",\"zip_code\":\"94110\",\"country\":\"US\"}}}]",
    "vehicle_type": "cargo_van"
  }
}
```

### Create Order Example
```json
{
  "tool": "create_order",
  "arguments": {
    "delivery_info": "{\"service_type\":\"delivery\"}",
    "pickup_info": "{\"business_name\":\"Test Business\",\"contact_name\":\"John Doe\",\"contact_phone_number\":\"555-123-4567\"}",
    "drop_offs": "[{\"business_name\":\"Drop Off Business\",\"contact_name\":\"Jane Smith\",\"contact_phone_number\":\"555-987-6543\"}]"
  }
}
```

This guide provides everything needed to build the MCP server. The AI agent can follow this step-by-step and ask clarifying questions as needed.
