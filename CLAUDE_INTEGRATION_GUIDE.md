# Claude AI Integration Guide

This guide explains how to set up and use Claude AI integration with the Dispatch MCP server for enhanced conversational pricing capabilities.

## 🤖 What is Claude Integration?

Claude integration adds natural language understanding and conversational capabilities to the pricing system, making it more intelligent and user-friendly.

### Key Benefits:
- **Natural Language Processing**: Understands complex user queries
- **Contextual Responses**: Provides relevant, personalized advice
- **Conversational Flow**: Maintains context across multiple interactions
- **Intelligent Recommendations**: Suggests optimal pricing strategies
- **Fallback Support**: Gracefully degrades to rule-based system if Claude is unavailable

## 🚀 Setup Instructions

### 1. Get Anthropic API Key

1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Sign up or log in to your account
3. Navigate to API Keys section
4. Create a new API key
5. Copy the key for configuration

### 2. Configure Environment

Add Claude configuration to your `.env` file:

```bash
# Claude AI Configuration
ANTHROPIC_API_KEY=your_anthropic_api_key_here
USE_CLAUDE_AI=true
```

### 3. Test the Integration

Test Claude integration with the CLI:

```bash
# Test with Claude enabled
./bin/dispatch-cli chat

# Test fallback mode (without API key)
unset ANTHROPIC_API_KEY
./bin/dispatch-cli chat
```

## 💬 Using the Conversational Interface

### Starting a Conversation

```bash
./bin/dispatch-cli chat
```

### Example Conversations

**Basic Pricing Query:**
```
You: I need 3 deliveries to different locations
Claude: Great! With 3 deliveries, you qualify for our Multi-Delivery Discount (15% off). 
       This could save you $7.50 on a $50 order. Would you like to see all available 
       pricing options?
```

**Customer Tier Inquiry:**
```
You: I'm a gold tier customer, what's the best pricing for me?
Claude: As a gold tier customer, you have access to our Loyalty Discount (10% off) 
       plus other volume-based discounts. With your status, you could save up to 20% 
       on qualifying orders. How many deliveries do you need?
```

**Complex Scenario:**
```
You: I need 5 deliveries this month and I'm planning to order 10 times next month
Claude: Excellent! You're setting up for great savings:
       - Current order: Multi-Delivery Discount (15% off)
       - Next month: Volume Discount (20% off) with 10+ orders
       - Plus your gold tier loyalty discount (10% off)
       
       This combination could save you 20-25% on future orders!
```

## 🔧 Technical Details

### Architecture

```
User Input → Claude API → Pricing Engine → Response Generation
     ↓
Fallback to Rule-Based Engine (if Claude unavailable)
```

### Hybrid Approach

The system uses a hybrid approach:

1. **Primary**: Claude AI for natural language understanding
2. **Fallback**: Rule-based engine for reliability
3. **Integration**: Pricing engine for accurate calculations

### API Usage

- **Model**: `claude-3-sonnet-20240229`
- **Max Tokens**: 1000 per request
- **Timeout**: 30 seconds
- **Rate Limits**: Follows Anthropic's API limits

## 🧪 Testing

### Unit Tests

```bash
go test ./test/claude_integration_test.go -v
```

### Integration Tests

```bash
# Test with Claude
ANTHROPIC_API_KEY=your_key ./bin/dispatch-cli chat

# Test fallback
unset ANTHROPIC_API_KEY
./bin/dispatch-cli chat
```

### Test Scenarios

1. **Claude Available**: Tests natural language processing
2. **Claude Unavailable**: Tests fallback to rule-based system
3. **Context Management**: Tests conversation memory
4. **Pricing Integration**: Tests accurate pricing calculations

## 🐛 Troubleshooting

### Common Issues

**1. Claude Not Responding**
```
Error: Claude not available, using rule-based engine
```
**Solution**: Check your `ANTHROPIC_API_KEY` environment variable

**2. API Rate Limits**
```
Error: API request failed with status 429
```
**Solution**: Wait a moment and try again, or check your Anthropic account limits

**3. Invalid API Key**
```
Error: API request failed with status 401
```
**Solution**: Verify your API key is correct and active

### Debug Mode

Enable debug logging to see Claude interactions:

```bash
export DEBUG_CLAUDE=true
./bin/dispatch-cli chat
```

## 📊 Performance Considerations

### Response Times
- **Claude API**: 1-3 seconds per request
- **Fallback Engine**: <100ms per request
- **Hybrid Mode**: Best of both worlds

### Cost Management
- **Token Usage**: Monitor via Anthropic console
- **Caching**: Responses cached for similar queries
- **Fallback**: Reduces API calls when possible

## 🔒 Security

### API Key Security
- Store API keys in environment variables
- Never commit keys to version control
- Use `.env` files for local development
- Rotate keys regularly

### Data Privacy
- Claude processes conversation context
- No sensitive customer data sent to Claude
- Pricing calculations remain local

## 🚀 Advanced Usage

### Custom Prompts

Modify the system prompt in `internal/claude/client.go`:

```go
systemPrompt := `You are a helpful pricing advisor for a delivery service...`
```

### Context Enhancement

Add more context to Claude requests:

```go
pricingContext := &claude.PricingContext{
    DeliveryCount:   5,
    CustomerTier:    "gold",
    OrderFrequency:  10,
    TotalOrderValue: 200.0,
    IsBulkOrder:     true,
}
```

### Response Customization

Customize response formatting in `internal/conversation/claude_engine.go`:

```go
response := &ConversationResponse{
    Message:         claudeResponse.Content[0].Text,
    Recommendations: recommendations,
    NextQuestions:   nextQuestions,
    UpdatedContext:  updatedContext,
}
```

## 📈 Monitoring and Analytics

### Usage Tracking

Monitor Claude usage:
- API calls per session
- Response quality metrics
- Fallback frequency
- User satisfaction

### Performance Metrics

Track system performance:
- Response times
- Error rates
- Success rates
- User engagement

## 🔄 Updates and Maintenance

### Keeping Claude Updated

1. Monitor Anthropic API changes
2. Update client code as needed
3. Test new Claude models
4. Optimize prompts for better results

### Regular Maintenance

- Review API usage and costs
- Update system prompts
- Test fallback scenarios
- Monitor user feedback

## 📚 Additional Resources

- [Anthropic API Documentation](https://docs.anthropic.com/)
- [Claude Model Information](https://www.anthropic.com/claude)
- [Pricing Guide](./PRICING_GUIDE.md)
- [API Reference](./API_REFERENCE.md)
- [Testing Guide](./TESTING_GUIDE.md)

## 🤝 Support

For issues with Claude integration:

1. Check the troubleshooting section above
2. Review Anthropic's API documentation
3. Test with fallback mode to isolate issues
4. Check logs for detailed error messages

---

**Note**: Claude integration is optional. The system works perfectly with the rule-based conversation engine as a fallback.
