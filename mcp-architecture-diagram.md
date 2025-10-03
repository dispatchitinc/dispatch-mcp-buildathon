# MCP Tool Architecture Diagram

```mermaid
graph TB
    %% External Clients
    AI[AI Client<br/>Claude/Other] --> MCP[MCP Server<br/>dispatch-mcp-server]
    
    %% MCP Server Core
    MCP --> Server[MCPServer<br/>Orchestrator]
    Server --> Tools[Tool Registry]
    
    %% MCP Tools
    Tools --> T1[create_estimate<br/>Tool]
    Tools --> T2[create_order<br/>Tool]
    Tools --> T3[compare_pricing_models<br/>Tool]
    Tools --> T4[conversational_pricing_advisor<br/>Tool]
    
    %% Tool Implementations
    T1 --> DC1[Dispatch Client<br/>CreateEstimate]
    T2 --> DC2[Dispatch Client<br/>CreateOrder]
    T3 --> PE[Pricing Engine<br/>Compare Models]
    T4 --> CE[Claude Conversation Engine<br/>Natural Language Processing]
    
    %% External Services
    DC1 --> API[Dispatch GraphQL API<br/>External Service]
    DC2 --> API
    
    %% Pricing System
    PE --> PM1[Standard Pricing<br/>0% discount]
    PE --> PM2[Multi-Delivery<br/>15% discount]
    PE --> PM3[Volume Discount<br/>20% discount]
    PE --> PM4[Loyalty Discount<br/>10% discount]
    PE --> PM5[Bulk Order<br/>25% discount]
    
    %% Conversation Engine
    CE --> Claude[Claude AI<br/>Anthropic API]
    CE --> Fallback[Rule-based Engine<br/>Fallback Mode]
    
    %% Data Flow
    subgraph "Data Flow"
        direction TB
        Input[User Input<br/>Natural Language] --> CE
        CE --> Response[Structured Response<br/>JSON]
        Response --> AI
    end
    
    %% Configuration
    subgraph "Configuration"
        direction TB
        Config[Environment Config<br/>Auth Tokens, Endpoints]
        Config --> Server
        Config --> DC1
        Config --> DC2
        Config --> Claude
    end
    
    %% Mock Mode
    subgraph "Mock Mode"
        direction TB
        Mock[Mock Client<br/>Demo Data]
        DC1 -.-> Mock
        DC2 -.-> Mock
    end
    
    %% Styling
    classDef toolClass fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef engineClass fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef apiClass fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px
    classDef configClass fill:#fff3e0,stroke:#e65100,stroke-width:2px
    
    class T1,T2,T3,T4 toolClass
    class PE,CE engineClass
    class API,Claude apiClass
    class Config,Mock configClass
```

## Architecture Overview

### Core Components

1. **MCP Server** - Central orchestrator that receives requests from AI clients
2. **Tool Registry** - Manages available tools and their implementations
3. **Dispatch Client** - Handles communication with external Dispatch GraphQL API
4. **Pricing Engine** - Calculates and compares different pricing models
5. **Claude Conversation Engine** - Provides natural language processing capabilities

### Available Tools

1. **create_estimate** - Creates cost estimates for delivery orders
2. **create_order** - Creates actual delivery orders
3. **compare_pricing_models** - Compares different pricing strategies
4. **conversational_pricing_advisor** - Provides AI-powered pricing advice

### Pricing Models

- **Standard**: 0% discount (baseline)
- **Multi-Delivery**: 15% discount for 2+ deliveries
- **Volume**: 20% discount for 5+ deliveries + 3+ orders/month
- **Loyalty**: 10% discount for gold tier customers
- **Bulk Order**: 25% discount for 10+ deliveries

### Key Features

- **Fallback Mode**: Rule-based processing when Claude is unavailable
- **Mock Mode**: Demo functionality without external API dependencies
- **Context Management**: Maintains conversation state and customer profiles
- **Flexible Authentication**: Supports both token-based and IDP authentication
