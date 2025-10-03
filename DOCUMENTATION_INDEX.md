# 📚 Documentation Index - Dispatch MCP Server

## 🎯 Getting Started

### For New Users
- **[Quick Start Guide](QUICK_START.md)** - Get up and running in 5 minutes
- **[README.md](README.md)** - Main documentation with setup instructions
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions

## 🔧 Technical Documentation

### API Reference
- **[API_REFERENCE.md](API_REFERENCE.md)** - Complete API documentation
  - Tool parameters and responses
  - Data types and structures
  - Error handling
  - Usage patterns

### Pricing System
- **[PRICING_GUIDE.md](PRICING_GUIDE.md)** - Quick reference for pricing functionality
  - Available pricing models
  - Customer scenarios
  - MCP tool usage
  - Business use cases

### AI Integration
- **[CLAUDE_INTEGRATION_GUIDE.md](CLAUDE_INTEGRATION_GUIDE.md)** - Claude AI integration guide
  - Setup and configuration
  - Natural language processing
  - Conversational capabilities
  - Fallback mechanisms

## 🚀 Usage Guides

### CLI Commands
```bash
./bin/dispatch-cli help          # Show all commands
./bin/dispatch-cli estimate      # Create cost estimate
./bin/dispatch-cli order         # Create delivery order
./bin/dispatch-cli pricing       # Compare pricing models
./bin/dispatch-cli interactive   # Interactive mode
./bin/dispatch-cli status        # Check connection
```

### MCP Tools
- `create_estimate` - Create delivery cost estimates
- `create_order` - Create delivery orders  
- `compare_pricing_models` - Compare pricing options

## 💰 Pricing Models Reference

| Model | Discount | Requirements | Best For |
|-------|----------|--------------|----------|
| **Standard** | 0% | None | New customers |
| **Multi-Delivery** | 15% | 2+ deliveries | Multiple stops |
| **Volume** | 20% | 5+ deliveries + 3+ orders/month | Regular customers |
| **Loyalty** | 10% | Gold tier | VIP customers |
| **Bulk Order** | 25% | 10+ deliveries + bulk flag | Large orders |

## 🔧 Development

### Project Structure
```
mcp-dry-run/
├── cmd/
│   ├── cli/           # CLI application
│   └── server/       # MCP server
├── internal/
│   ├── auth/         # Authentication
│   ├── config/       # Configuration
│   ├── dispatch/     # Dispatch API client
│   ├── mcp/          # MCP server implementation
│   └── pricing/      # Pricing engine
├── samples/          # Sample data
├── test/            # Tests
└── docs/            # Documentation
```

### Key Files
- `internal/pricing/models.go` - Pricing model definitions
- `internal/mcp/tools.go` - MCP tool implementations
- `cmd/cli/main.go` - CLI application
- `internal/dispatch/client.go` - API client

### Adding New Features
1. **New Pricing Models**: Edit `internal/pricing/models.go`
2. **New MCP Tools**: Add to `internal/mcp/server.go` and `tools.go`
3. **New CLI Commands**: Add to `cmd/cli/main.go`

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/pricing
go test ./internal/dispatch
```

### Integration Tests
```bash
go test ./test
```

### Manual Testing
```bash
./bin/dispatch-cli pricing     # Test pricing comparison
./bin/dispatch-cli estimate    # Test estimate creation
./bin/dispatch-cli order       # Test order creation
```

## 🔒 Security

### Authentication Methods
1. **IDP Authentication** (Production)
   - OAuth 2.0 with Identity Provider
   - Automatic token refresh
   - Organization-based access

2. **Static Token** (Development)
   - Simple token-based auth
   - Good for testing
   - No automatic refresh

### Environment Variables
```bash
# IDP Authentication
export USE_IDP_AUTH=true
export IDP_CLIENT_ID=your_client_id
export IDP_CLIENT_SECRET=your_client_secret
export DISPATCH_ORGANIZATION_ID=your_org_id

# Static Token
export USE_IDP_AUTH=false
export DISPATCH_AUTH_TOKEN=your_token
```

## 📊 Business Use Cases

### Customer Onboarding
- Show potential savings with different pricing models
- Demonstrate value of loyalty programs
- Guide customers to optimal pricing

### Sales Optimization
- Prove ROI of volume commitments
- Demonstrate competitive advantages
- Enable data-driven pricing decisions

### Revenue Optimization
- Find optimal pricing for each customer segment
- Test different pricing strategies
- Maximize customer lifetime value

## 🎨 Customization

### Adding New Pricing Models
1. Define model type in `PricingModel` enum
2. Add rule in `initializeDefaultRules()`
3. Update eligibility logic in `isEligibleForModel()`
4. Add custom discount logic in `calculateAdditionalDiscount()`

### Modifying Existing Models
- **Discount Rates**: Change `BaseMultiplier` values
- **Eligibility**: Update thresholds in `isEligibleForModel()`
- **Additional Discounts**: Modify `calculateAdditionalDiscount()`

## 🔍 Troubleshooting

### Common Issues
- **Authentication failures**: Check credentials and environment variables
- **No eligible pricing models**: Verify customer context parameters
- **Pricing too low**: Check minimum price protection (50% of original)
- **MCP tool errors**: Validate JSON format and parameter types

### Debug Mode
```bash
export DEBUG=true
./bin/dispatch-cli pricing
```

### Mock Mode
```bash
./bin/dispatch-cli logout  # Switches to mock mode
./bin/dispatch-cli estimate
```

## 📞 Support

### Documentation
- **README.md**: Main documentation
- **API_REFERENCE.md**: Complete API reference
- **PRICING_GUIDE.md**: Pricing functionality guide
- **TROUBLESHOOTING.md**: Common issues and solutions

### Getting Help
- **GitHub Issues**: Report bugs and feature requests
- **Development Team**: Contact for urgent issues
- **Community**: Ask questions in discussions

### Log Collection
When reporting issues, include:
1. Error messages
2. Configuration (without secrets)
3. Steps to reproduce
4. System information
5. Relevant logs

## 🚀 Quick Commands

### Build and Run
```bash
go mod tidy
chmod +x build.sh
./build.sh
./bin/dispatch-cli pricing
```

### Development
```bash
go build -o bin/dispatch-cli cmd/cli/main.go
go build -o bin/dispatch-mcp-server cmd/server/main.go
```

### Testing
```bash
go test ./...
./bin/dispatch-cli pricing
./bin/dispatch-cli estimate
```

---

*This index provides a comprehensive overview of all documentation and resources available for the Dispatch MCP Server project.*
