# Dispatch MCP Buildathon Project

A comprehensive Model Context Protocol (MCP) server implementation for Dispatch's order creation and estimate APIs, developed during the Fall 2025 Buildathon. This project combines conversational AI capabilities with Dispatch's logistics platform to create an intelligent interface for order management and pricing optimization.

## 🚀 Project Overview

This repository contains both the **implementation** (MCP server) and **planning documentation** from the Fall 2025 Buildathon project. The goal was to create a conversational AI interface for Dispatch orders and pricing estimates using the Model Context Protocol framework.

### Key Features

- **🤖 Claude AI Integration**: Natural language pricing conversations powered by Anthropic's Claude via AI Hub
- **💬 Conversational Pricing Advisor**: Interactive chat interface for pricing recommendations
- **🌐 Web Chat Interface**: Modern web-based chat interface with real-time order tracking
- **📦 Order Management**: Create estimates and orders through conversational interface
- **💰 Advanced Pricing Models**: Compare different pricing strategies (multi-delivery, volume discounts, loyalty programs)
- **📊 Real-time Progress Tracking**: Visual order creation progress and pricing recommendations
- **✅ Input Validation**: Comprehensive validation for all user inputs with detailed error messages
- **🔐 IDP Authentication**: Support for Identity Provider token management with auto-refresh
- **🛡️ Production Ready**: Secure authentication with OAuth 2.0 and comprehensive error handling

## 📁 Repository Structure

```
dispatch-mcp-buildathon/
├── docs/                               # All documentation
│   ├── README.md                       # Main documentation
│   ├── API_REFERENCE.md                # API documentation
│   ├── PRICING_GUIDE.md                # Pricing guide
│   ├── QUICK_START.md                  # Quick start guide
│   ├── TESTING_GUIDE.md                # Testing documentation
│   ├── TROUBLESHOOTING.md              # Troubleshooting guide
│   ├── AI_HUB_INTEGRATION.md           # AI Hub integration
│   └── buildathon-planning/            # Buildathon planning documents
├── scripts/                            # Shell scripts
│   ├── build.sh                        # Build script
│   ├── demo.sh                         # Demo script
│   ├── demo-cli.sh                      # CLI demo
│   ├── chat_with_ai_hub.sh             # AI Hub chat script
│   └── test_*.sh                        # Test scripts
├── cmd/                                # Application entry points
│   ├── cli/main.go                     # Command-line interface
│   ├── server/main.go                  # MCP server
│   └── web/main.go                     # Web chat server
├── internal/                           # Core application logic
│   ├── auth/                          # Authentication handling
│   ├── claude/                        # Claude AI integration
│   ├── config/                        # Configuration management
│   ├── conversation/                  # Conversational AI engine
│   ├── dispatch/                      # Dispatch API client
│   ├── mcp/                          # MCP server implementation
│   ├── pricing/                      # Pricing model engine
│   └── validation/                   # Input validation
├── static/                            # Web interface assets
│   └── index.html                     # Web chat interface
├── test/                              # Test suites
├── samples/                           # Sample data files
├── bin/                               # Built binaries
│   ├── dispatch-cli                   # CLI tool
│   ├── dispatch-mcp-server            # MCP server
│   └── dispatch-web                   # Web server
└── [configuration files]              # go.mod, Makefile, etc.
```

## 🏗️ Buildathon Timeline

This project was developed over 3 buildathon days:

- **Day 1** (Oct 3, 2025): Foundation and Setup
- **Day 2** (Oct 10, 2025): Integration and Testing  
- **Day 3** (Oct 17, 2025): Polish and Demo

See `docs/buildathon-planning/` for detailed planning documents and timeline.

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- Dispatch API credentials
- (Optional) Anthropic API key for Claude AI features

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd dispatch-mcp-buildathon
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configure authentication** (choose one method):

   **Option 1: IDP Authentication (Recommended)**
   ```bash
   export USE_IDP_AUTH=true
   export IDP_ENDPOINT=https://id.dispatchfog.io
   export IDP_CLIENT_ID=your_client_id
   export IDP_CLIENT_SECRET=your_client_secret
   export IDP_SCOPE=dispatch:api
   export IDP_TOKEN_ENDPOINT=https://id.dispatchfog.io/oauth/token
   export DISPATCH_ORGANIZATION_ID=your_org_id_here
   ```

   **Option 2: Static Token (Development)**
   ```bash
   export USE_IDP_AUTH=false
   export DISPATCH_AUTH_TOKEN=your_static_token_here
   export DISPATCH_ORGANIZATION_ID=your_org_id_here
   ```

4. **Configure Claude AI via AI Hub (Recommended)**:
   ```bash
   export USE_AI_HUB=true
   export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1
   export ANTHROPIC_API_KEY=your_ai_hub_api_key_here
   export AI_HUB_MODEL=claude-sonnet  # or claude-haiku for faster responses
   ```

   **Alternative: Direct Claude API**:
   ```bash
   export USE_AI_HUB=false
   export ANTHROPIC_API_KEY=your_anthropic_api_key_here
   export USE_CLAUDE_AI=true
   ```

5. **Build the server**:
   ```bash
   chmod +x scripts/build.sh
   ./scripts/build.sh
   ```

6. **Run the server**:
   ```bash
   ./bin/dispatch-mcp-server
   ```

## 🛠️ Usage

### Web Chat Interface (Recommended)

The project includes a modern web-based chat interface:

```bash
# Start the web server
PORT=8081 ./bin/dispatch-web

# Open in browser
open http://localhost:8081
```

**Features:**
- 💬 Real-time chat with AI assistant
- 📊 Live order creation progress tracking
- 💰 Real-time pricing recommendations
- 📱 Mobile-responsive design
- 🤔 Animated thinking indicators

### CLI Interface

The project includes a comprehensive CLI for testing and interaction:

```bash
# Check connection status
./bin/dispatch-cli status

# Set subenvironment (monkey, staging, prod)
./bin/dispatch-cli subenv

# Authenticate with Dispatch API
./bin/dispatch-cli login

# Create cost estimate
./bin/dispatch-cli estimate

# Create delivery order
./bin/dispatch-cli order

# Compare pricing models
./bin/dispatch-cli pricing

# Interactive mode
./bin/dispatch-cli interactive
```

### MCP Tools

The server provides several MCP tools for AI agent integration:

#### Create Estimate
**Note**: The `vehicle_type` parameter is required. Always ask the user what type of vehicle they need before calling this tool.

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

**Available Vehicle Types**:
- `pickup_truck` - Small to medium items, quick deliveries
- `cargo_van` - Medium to large packages, furniture  
- `sprinter_van` - Large items, multiple packages
- `box_truck` - Very large items, bulk deliveries

#### Create Order
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

#### Select Delivery Option
```json
{
  "tool": "select_delivery_option",
  "arguments": {
    "estimate_response": "{\"data\":{\"createEstimate\":{\"estimate\":{\"availableOrderOptions\":[...]}}}}",
    "delivery_scenario": "fastest"  // or "cheapest"
  }
}
```

#### Compare Pricing Models
```json
{
  "tool": "compare_pricing_models",
  "arguments": {
    "original_estimate": "{\"serviceType\":\"delivery\",\"estimatedOrderCost\":45.99,\"vehicleType\":\"cargo_van\"}",
    "delivery_count": "3",
    "customer_tier": "gold",
    "order_frequency": "5",
    "total_order_value": "150.00",
    "is_bulk_order": "false"
  }
}
```

## 🚚 Delivery Scenarios

The system supports two key delivery scenarios that align with common customer needs:

### Scenario 1: "I need this delivered as soon as possible"
- **Use**: First option in the estimate response
- **Characteristics**: Fastest delivery time, highest cost
- **Best for**: Urgent deliveries, time-sensitive items

### Scenario 2: "I need this delivered sometime today"  
- **Use**: Last option in the estimate response
- **Characteristics**: Slowest delivery time, lowest cost
- **Best for**: Non-urgent deliveries, cost-conscious customers

The API returns multiple delivery options sorted by speed and cost, allowing customers to choose between speed and cost based on their specific needs.

## ✅ Input Validation

The system includes comprehensive input validation to ensure data quality:

- **Vehicle Type Validation**: Ensures only valid vehicle types (pickup_truck, cargo_van, sprinter_van, box_truck)
- **Delivery Scenario Validation**: Validates delivery scenarios (fastest, cheapest, etc.)
- **Address Validation**: Validates address format and required fields
- **JSON Format Validation**: Ensures all JSON parameters are properly formatted
- **Numeric Validation**: Validates numeric parameters with range checking
- **Boolean Validation**: Ensures boolean parameters are properly formatted

All validation errors include detailed messages to help users correct their inputs.

## 💰 Pricing Model Comparison

The system includes advanced pricing comparison capabilities:

| Model | Description | Discount | Requirements |
|-------|-------------|----------|--------------|
| **Standard Pricing** | No discounts | 0% | None |
| **Multi-Delivery Discount** | Discount for multiple deliveries | 15% | 2+ deliveries |
| **Volume Discount** | Discount for high-volume customers | 20% | 5+ deliveries + 3+ orders/month |
| **Loyalty Discount** | Discount for loyal customers | 10% | Gold tier customer |
| **Bulk Order Discount** | Discount for large bulk orders | 25% | 10+ deliveries + bulk order flag |

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./test

# Run specific test categories
go test ./test -run TestClaudeIntegration
go test ./test -run TestConversation
go test ./test -run TestMCP
```

**Note**: Tests require valid Dispatch credentials to run.

## 📚 Documentation

- **[Quick Start Guide](QUICK_START.md)**: Get up and running in 5 minutes
- **[API Reference](API_REFERENCE.md)**: Complete API documentation
- **[Pricing Guide](PRICING_GUIDE.md)**: Quick reference for pricing functionality
- **[Testing Guide](TESTING_GUIDE.md)**: Comprehensive testing documentation
- **[Troubleshooting](TROUBLESHOOTING.md)**: Common issues and solutions
- **[MCP Prompt Guide](MCP_PROMPT_GUIDE.md)**: Understanding MCP server prompts and AI behavior
- **[Buildathon Planning](docs/buildathon-planning/)**: Original project planning documents

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏆 Buildathon Team

- **Camron**: Go/MCP Development
- **Julia**: UI/Prompts Development  
- **Tyler**: Domain Expertise
- **Chris**: Project Oversight

## 🔗 Related Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/)
- [Dispatch API Documentation](https://docs.dispatchfog.io/)
- [Anthropic Claude API](https://docs.anthropic.com/)

---

*This project was developed during the Fall 2025 Buildathon as a proof of concept for conversational AI integration with Dispatch's logistics platform.*