# 🔗 AI Hub Integration Guide

## Overview

This guide explains how to configure the Dispatch MCP Buildathon project to work with AI Hub for cost tracking and API management.

## 🎯 **Why AI Hub?**

AI Hub provides:
- **Cost Tracking**: Monitor API usage and costs across all AI services
- **Rate Limiting**: Control API usage and prevent overages
- **Authentication**: Centralized API key management
- **Analytics**: Detailed usage reports and insights
- **Proxy Services**: Route requests through your infrastructure

## 🔧 **Configuration**

### Environment Variables

Set these environment variables to enable AI Hub integration:

```bash
# Enable AI Hub proxy
export USE_AI_HUB=true

# Your AI Hub endpoint
export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1

# Your AI Hub API key (this might be different from direct Anthropic key)
export ANTHROPIC_API_KEY=your_ai_hub_api_key_here

# Optional: Specify model (defaults to claude-sonnet for conversational pricing)
export AI_HUB_MODEL=claude-sonnet  # or claude-haiku for faster responses
```

### Example Configuration

```bash
# For development
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1
export ANTHROPIC_API_KEY=sk-hB4t-5i9G701-YI_gVc2Hw

# For production
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1
export ANTHROPIC_API_KEY=your_production_ai_hub_key
```

## 🤖 **Model Selection**

Your AI Hub supports these models:

| Model | Purpose | Best For |
|-------|---------|----------|
| `claude-haiku` | Fast responses, quick thinking | Simple queries, high-volume scenarios |
| `claude-sonnet` | Deeper logic chains, better human text | Complex reasoning, conversational AI |

**For Conversational Pricing**: Use `claude-sonnet` (default) for better reasoning and more natural responses.

**For High-Volume Scenarios**: Use `claude-haiku` for faster, more cost-effective responses.

## 🚀 **Usage**

### 1. **Direct Anthropic API** (Default)
```bash
# No AI Hub - direct API calls
export USE_AI_HUB=false
export ANTHROPIC_API_KEY=your_direct_anthropic_key
./bin/dispatch-cli chat
```

### 2. **AI Hub Proxy** (Recommended)
```bash
# Use AI Hub for cost tracking
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://your-ai-hub.com/v1
export ANTHROPIC_API_KEY=your_ai_hub_key
./bin/dispatch-cli chat
```

## 🔍 **How It Works**

### Request Flow

```
User Input → Dispatch MCP → Claude Client → AI Hub → Anthropic API
                ↓
            Response ← AI Hub ← Anthropic API
```

### Authentication

**Direct API:**
- Uses `x-api-key` header
- Requires `anthropic-version` header
- Direct connection to Anthropic

**AI Hub:**
- Uses `Authorization: Bearer` header
- No version header required
- Routes through your AI Hub infrastructure

## 🧪 **Testing**

### Test AI Hub Integration

```bash
# Set up AI Hub configuration
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://your-ai-hub.com/v1
export ANTHROPIC_API_KEY=your_ai_hub_key

# Test the integration
go run debug_entity_extraction.go
```

### Test Direct API

```bash
# Set up direct API configuration
export USE_AI_HUB=false
export ANTHROPIC_API_KEY=your_direct_anthropic_key

# Test the integration
go run debug_entity_extraction.go
```

## 📊 **Monitoring**

### AI Hub Benefits

1. **Cost Tracking**: See exactly how much each conversation costs
2. **Usage Analytics**: Track API calls, tokens, and response times
3. **Rate Limiting**: Prevent API overages and unexpected costs
4. **Audit Trail**: Complete log of all AI interactions

### Metrics Available

- **API Calls**: Number of requests per session
- **Token Usage**: Input and output tokens consumed
- **Response Times**: Latency measurements
- **Cost Per Request**: Detailed cost breakdown
- **Error Rates**: Failed request tracking

## 🔧 **Troubleshooting**

### Common Issues

**1. Authentication Errors**
```
Error: API request failed with status 401
```
**Solution**: Check your AI Hub API key and endpoint

**2. Endpoint Not Found**
```
Error: connection refused
```
**Solution**: Verify your AI Hub endpoint URL

**3. Rate Limiting**
```
Error: API request failed with status 429
```
**Solution**: Check your AI Hub rate limits and quotas

### Debug Mode

Enable debug logging to see request details:

```bash
export DEBUG_CLAUDE=true
./bin/dispatch-cli chat
```

## 🎯 **Best Practices**

### 1. **Environment Management**
- Use different AI Hub endpoints for dev/staging/prod
- Rotate API keys regularly
- Use environment-specific configurations

### 2. **Cost Optimization**
- Monitor token usage in AI Hub dashboard
- Set up alerts for high usage
- Use caching when possible

### 3. **Security**
- Never commit API keys to version control
- Use environment variables for configuration
- Rotate keys regularly

## 📈 **Advanced Configuration**

### Custom Headers

If your AI Hub requires custom headers, modify `internal/claude/client.go`:

```go
if useAIHub == "true" {
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("X-Client-ID", "dispatch-mcp-server")
    req.Header.Set("X-Environment", "production")
}
```

### Request Transformation

If AI Hub requires request transformation, add logic in `CreateMessage`:

```go
// Transform request for AI Hub if needed
if useAIHub == "true" {
    // Add any request transformation logic here
}
```

## 🚀 **Production Deployment**

### 1. **Set Environment Variables**
```bash
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://ai-hub.your-company.com/v1
export ANTHROPIC_API_KEY=your_production_key
```

### 2. **Test Integration**
```bash
./bin/dispatch-cli chat
```

### 3. **Monitor Usage**
- Check AI Hub dashboard for usage metrics
- Set up alerts for cost thresholds
- Monitor response times and error rates

## 📚 **Additional Resources**

- **AI Hub Documentation**: Your internal AI Hub docs
- **Cost Tracking**: AI Hub usage dashboard
- **API Reference**: AI Hub API documentation
- **Support**: Contact your AI Hub team

---

**Note**: Replace `your-ai-hub-domain.com` and other placeholders with your actual AI Hub configuration.
