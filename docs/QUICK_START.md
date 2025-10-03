# 🚀 Quick Start Guide - Dispatch MCP Server

## ⚡ 5-Minute Setup

### 1. Build the Project
```bash
cd mcp-dry-run
go mod tidy
chmod +x build.sh
./build.sh
```

### 2. Test the CLI
```bash
# Check if everything works
./bin/dispatch-cli help

# Run the pricing demo (works without authentication)
./bin/dispatch-cli pricing
```

### 3. Set Up Authentication (Optional)
```bash
# For production use
export USE_IDP_AUTH=true
export IDP_CLIENT_ID=your_client_id
export IDP_CLIENT_SECRET=your_client_secret
export DISPATCH_ORGANIZATION_ID=your_org_id

# Or use static token for testing
export USE_IDP_AUTH=false
export DISPATCH_AUTH_TOKEN=your_token
```

### 4. Run the Full Demo
```bash
./demo-cli.sh
```

## 🎯 What You Can Do

### 1. Create Estimates
```bash
./bin/dispatch-cli estimate
```
- Creates a sample delivery estimate
- Shows cost, vehicle type, and delivery time
- Works with mock data (no auth required)

### 2. Compare Pricing Models
```bash
./bin/dispatch-cli pricing
```
- Shows 4 different customer scenarios
- Demonstrates various pricing models
- Calculates potential savings

### 3. Create Orders
```bash
./bin/dispatch-cli order
```
- Creates a sample delivery order
- Shows order ID and tracking number
- Requires authentication for real orders

## 🔧 MCP Server Usage

### Start the Server
```bash
./bin/dispatch-mcp-server
```

### Use with AI Agents
The server provides these tools:
- `create_estimate`: Get delivery cost estimates
- `create_order`: Create delivery orders
- `compare_pricing_models`: Compare pricing options

### Example AI Integration
```json
{
  "tool": "compare_pricing_models",
  "arguments": {
    "original_estimate": "{\"estimatedOrderCost\":45.99}",
    "delivery_count": "3",
    "customer_tier": "gold"
  }
}
```

## 💰 Pricing Models Explained

| Model | Discount | When to Use |
|-------|----------|-------------|
| **Standard** | 0% | New customers, single deliveries |
| **Multi-Delivery** | 15% | 2+ deliveries in one order |
| **Volume** | 20% | 5+ deliveries + 3+ orders/month |
| **Loyalty** | 10% | Gold tier customers |
| **Bulk Order** | 25% | 10+ deliveries + bulk flag |

## 🎨 Customization

### Add Your Own Pricing Models
1. Edit `internal/pricing/models.go`
2. Add new model type
3. Define pricing rules
4. Update eligibility logic

### Modify Existing Models
- Change discount rates
- Adjust eligibility requirements
- Add custom discount logic

## 🔍 Troubleshooting

### Common Issues
- **"Command not found"**: Run `go build -o bin/dispatch-cli cmd/cli/main.go`
- **"Permission denied"**: Run `chmod +x bin/dispatch-cli`
- **"No delivery options"**: Check if locations are in service area
- **"Authentication failed"**: Verify credentials or use mock mode

### Debug Mode
```bash
export DEBUG=true
./bin/dispatch-cli pricing
```

### Mock Mode (No Auth Required)
```bash
./bin/dispatch-cli logout  # Switches to mock mode
./bin/dispatch-cli estimate
```

## 📚 Next Steps

1. **Read the Full Documentation**: [README.md](README.md)
2. **Explore the API**: [API_REFERENCE.md](API_REFERENCE.md)
3. **Learn Pricing Models**: [PRICING_GUIDE.md](PRICING_GUIDE.md)
4. **Get Help**: [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

## 🎉 Success!

You now have a working Dispatch MCP Server with pricing model comparison! 

**Try these commands:**
```bash
./bin/dispatch-cli pricing     # See pricing comparison demo
./bin/dispatch-cli estimate    # Create a sample estimate
./bin/dispatch-cli interactive # Interactive mode
```

**Need help?** Check the [Troubleshooting Guide](TROUBLESHOOTING.md) or create an issue on GitHub.

---

*Happy coding! 🚀*
