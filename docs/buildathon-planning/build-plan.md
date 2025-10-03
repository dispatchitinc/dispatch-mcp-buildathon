# Buildathon MCP Service - Build Plan

## Project Overview
Build a Go-based MCP (Model Context Protocol) service that can call Dispatch's order creation and estimate mutations.

## Target Mutations Analysis
**Need to analyze**:
- Create Estimate: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_estimate.rb
- Create Order: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_order.rb

## MCP Service Design

### Core Structure
Based on Chris's Jira MCP server and mcp-go framework:

```
dispatch-mcp-service/
├── main.go
├── internal/
│   ├── mcp/
│   │   ├── server.go
│   │   ├── handlers.go
│   │   └── types.go
│   ├── dispatch/
│   │   ├── client.go
│   │   ├── mutations.go
│   │   └── types.go
│   └── config/
│       └── config.go
├── go.mod
└── README.md
```

### MCP Endpoints to Implement
1. **create_estimate** - Call Dispatch's create estimate mutation
2. **create_order** - Call Dispatch's create order mutation
3. **health_check** - Service health endpoint
4. **list_services** - List available MCP services

### Authentication Strategy
- **MVP**: Environment variables for auth tokens
- **Future**: Proper Dispatch auth integration

## Implementation Plan

### Phase 1: Setup & Foundation (Day 1)
1. **Project Setup**
   - Initialize Go module
   - Set up project structure
   - Add mcp-go dependency
   - Create basic server skeleton

2. **Research & Analysis**
   - Study target GraphQL mutations
   - Understand mutation parameters and responses
   - Review Chris's Jira MCP implementation
   - Understand Dispatch API structure

### Phase 2: Core MCP Implementation (Day 2)
1. **MCP Server Setup**
   - Implement basic MCP server using mcp-go
   - Create server configuration
   - Add health check endpoint

2. **Dispatch Client**
   - Create GraphQL client for Dispatch API
   - Implement authentication (env vars for MVP)
   - Add basic error handling

### Phase 3: Mutation Integration (Day 3)
1. **Create Estimate Handler**
   - Implement create_estimate MCP endpoint
   - Map MCP parameters to GraphQL mutation
   - Add response formatting

2. **Create Order Handler**
   - Implement create_order MCP endpoint
   - Map MCP parameters to GraphQL mutation
   - Add response formatting

3. **Testing & Validation**
   - Test both endpoints
   - Validate responses
   - Add error handling

## Technical Requirements

### Dependencies
```go
// go.mod
module dispatch-mcp-service

go 1.21

require (
    github.com/mark3labs/mcp-go v0.1.0
    // Add GraphQL client library
    // Add HTTP client library
)
```

### Environment Variables
```bash
DISPATCH_API_URL=https://api.dispatch.com/graphql
DISPATCH_AUTH_TOKEN=your_token_here
MCP_SERVER_PORT=8080
```

### MCP Service Configuration
```go
type Config struct {
    DispatchAPIURL    string
    DispatchAuthToken string
    ServerPort        string
}
```

## Development Steps

### Step 1: Analyze Target Mutations
- [ ] Review create_estimate.rb mutation
- [ ] Review create_order.rb mutation
- [ ] Understand required parameters
- [ ] Understand response structure
- [ ] Document mutation schemas

### Step 2: Set Up Project
- [ ] Create project directory
- [ ] Initialize go.mod
- [ ] Add dependencies
- [ ] Create basic structure
- [ ] Add configuration

### Step 3: Implement MCP Server
- [ ] Create MCP server using mcp-go
- [ ] Add health check endpoint
- [ ] Add service listing endpoint
- [ ] Test basic MCP functionality

### Step 4: Implement Dispatch Client
- [ ] Create GraphQL client
- [ ] Add authentication
- [ ] Implement mutation calls
- [ ] Add error handling

### Step 5: Implement MCP Endpoints
- [ ] Create estimate endpoint
- [ ] Create order endpoint
- [ ] Add parameter mapping
- [ ] Add response formatting

### Step 6: Testing & Validation
- [ ] Test create estimate
- [ ] Test create order
- [ ] Validate responses
- [ ] Add comprehensive error handling

## Success Criteria
- [ ] MCP service can be started and responds to health checks
- [ ] Service can call Dispatch create_estimate mutation
- [ ] Service can call Dispatch create_order mutation
- [ ] Proper error handling and logging
- [ ] Clear documentation for usage

## Next Steps
1. Analyze the target GraphQL mutations
2. Set up the project structure
3. Implement basic MCP server
4. Add Dispatch client integration
5. Implement the two main endpoints
6. Test and validate functionality

---
*This plan will be updated as we learn more about the target mutations and requirements.*
