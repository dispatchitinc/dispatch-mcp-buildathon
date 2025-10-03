# 🔧 Troubleshooting Guide - Dispatch MCP Server

## 🚨 Common Issues

### Authentication Problems

#### Issue: "Failed to get auth token"
**Symptoms:**
```
Error: failed to get auth token: authentication failed
```

**Solutions:**
1. **Check IDP Configuration:**
   ```bash
   export USE_IDP_AUTH=true
   export IDP_ENDPOINT=https://id.dispatchfog.io
   export IDP_CLIENT_ID=your_client_id
   export IDP_CLIENT_SECRET=your_client_secret
   export IDP_SCOPE=dispatch:api
   export IDP_TOKEN_ENDPOINT=https://id.dispatchfog.io/oauth/token
   ```

2. **Verify Organization ID:**
   ```bash
   export DISPATCH_ORGANIZATION_ID=your_org_id_here
   ```

3. **Test with Static Token:**
   ```bash
   export USE_IDP_AUTH=false
   export DISPATCH_AUTH_TOKEN=your_static_token_here
   ```

#### Issue: "Invalid credentials"
**Solutions:**
- Verify client ID and secret are correct
- Check if credentials have expired
- Ensure proper permissions for the organization

### Pricing Model Issues

#### Issue: "No eligible pricing models found"
**Symptoms:**
```json
{
  "pricing_models": [
    {
      "eligible": false,
      "reason": "Requires 2+ deliveries, you have 1"
    }
  ]
}
```

**Solutions:**
1. **Check Delivery Count:**
   ```json
   {
     "delivery_count": "2"  // Minimum for multi-delivery discount
   }
   ```

2. **Verify Customer Tier:**
   ```json
   {
     "customer_tier": "gold"  // Required for loyalty discount
   }
   ```

3. **Check Order Frequency:**
   ```json
   {
     "order_frequency": "5"  // Required for volume discount
   }
   ```

#### Issue: "Pricing seems too low"
**Symptoms:**
- Adjusted cost is significantly lower than expected
- Discount percentage seems too high

**Solutions:**
1. **Check Minimum Price Protection:**
   - System enforces 50% minimum of original cost
   - This prevents unrealistic pricing

2. **Verify Discount Calculations:**
   ```go
   // Base multiplier: 0.85 (15% discount)
   // Additional discounts: up to 25% more
   // Total maximum: 40% discount
   ```

3. **Review Context Parameters:**
   - Ensure `total_order_value` is realistic
   - Check if `is_bulk_order` is set correctly

### MCP Tool Issues

#### Issue: "Invalid arguments format"
**Symptoms:**
```
Error: invalid arguments format
```

**Solutions:**
1. **Check JSON Format:**
   ```json
   {
     "original_estimate": "{\"estimatedOrderCost\":45.99}"  // Valid JSON string
   }
   ```

2. **Verify Parameter Types:**
   ```json
   {
     "delivery_count": "3",      // String, not number
     "customer_tier": "gold",    // String
     "is_bulk_order": "true"    // String "true"/"false"
   }
   ```

#### Issue: "Failed to parse original_estimate"
**Symptoms:**
```
Error: failed to parse original_estimate: invalid character 'x' looking for beginning of value
```

**Solutions:**
1. **Validate JSON:**
   ```bash
   echo '{"estimatedOrderCost":45.99}' | jq .
   ```

2. **Check Escaping:**
   ```json
   {
     "original_estimate": "{\"estimatedOrderCost\":45.99,\"serviceType\":\"delivery\"}"
   }
   ```

### API Connection Issues

#### Issue: "Failed to create estimate"
**Symptoms:**
```
Error: failed to create estimate: connection refused
```

**Solutions:**
1. **Check Network Connectivity:**
   ```bash
   curl -I https://graphql-gateway.monkey.dispatchfog.org/graphql
   ```

2. **Verify Endpoint:**
   ```bash
   export DISPATCH_GRAPHQL_ENDPOINT=https://graphql-gateway.monkey.dispatchfog.org/graphql
   ```

3. **Test with Mock Mode:**
   ```bash
   ./bin/dispatch-cli logout  # Switches to mock mode
   ./bin/dispatch-cli estimate
   ```

#### Issue: "No delivery options available"
**Symptoms:**
```
⚠️  No delivery options available
```

**Solutions:**
1. **Check Location Validity:**
   - Ensure addresses are complete and valid
   - Verify zip codes are correct
   - Check if locations are within service area

2. **Try Different Vehicle Types:**
   ```json
   {
     "vehicle_type": "sprinter_van"  // Instead of cargo_van
   }
   ```

3. **Verify Service Area:**
   - Check if pickup/drop-off locations are in supported areas
   - Try locations in major metropolitan areas

### CLI Issues

#### Issue: "Command not found"
**Symptoms:**
```
./bin/dispatch-cli: command not found
```

**Solutions:**
1. **Build the CLI:**
   ```bash
   go build -o bin/dispatch-cli cmd/cli/main.go
   ```

2. **Make Executable:**
   ```bash
   chmod +x bin/dispatch-cli
   ```

3. **Check Path:**
   ```bash
   ls -la bin/dispatch-cli
   ```

#### Issue: "Permission denied"
**Symptoms:**
```
./bin/dispatch-cli: Permission denied
```

**Solutions:**
1. **Fix Permissions:**
   ```bash
   chmod +x bin/dispatch-cli
   ```

2. **Check Ownership:**
   ```bash
   ls -la bin/dispatch-cli
   ```

### Performance Issues

#### Issue: "Slow response times"
**Symptoms:**
- Estimate creation takes >5 seconds
- Pricing comparison takes >1 second

**Solutions:**
1. **Check Network Latency:**
   ```bash
   ping graphql-gateway.monkey.dispatchfog.org
   ```

2. **Use Mock Mode for Testing:**
   ```bash
   ./bin/dispatch-cli logout
   ```

3. **Optimize Requests:**
   - Use minimal required parameters
   - Avoid unnecessary add-ons
   - Cache results when possible

#### Issue: "Memory usage high"
**Symptoms:**
- System becomes slow
- Out of memory errors

**Solutions:**
1. **Check Go Version:**
   ```bash
   go version  # Should be 1.19+
   ```

2. **Monitor Memory:**
   ```bash
   top -p $(pgrep dispatch-cli)
   ```

3. **Restart if Needed:**
   ```bash
   pkill dispatch-cli
   ./bin/dispatch-cli pricing
   ```

## 🔍 Debug Mode

### Enable Debug Logging
```bash
export DEBUG=true
./bin/dispatch-cli pricing
```

### Verbose Output
```bash
./bin/dispatch-cli pricing --verbose
```

### Check Logs
```bash
# Check system logs
journalctl -u dispatch-mcp-server

# Check application logs
tail -f /var/log/dispatch-mcp-server.log
```

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
# Test authentication
./bin/dispatch-cli status

# Test estimate creation
./bin/dispatch-cli estimate

# Test pricing comparison
./bin/dispatch-cli pricing

# Test order creation
./bin/dispatch-cli order
```

## 📊 Monitoring

### Health Checks
```bash
# Check server status
./bin/dispatch-cli status

# Test API connectivity
curl -f https://graphql-gateway.monkey.dispatchfog.org/graphql
```

### Metrics
- **Response Time**: Monitor API response times
- **Error Rate**: Track failed requests
- **Success Rate**: Monitor successful operations
- **Memory Usage**: Track memory consumption

## 🆘 Getting Help

### Documentation
- **README.md**: Main documentation
- **API_REFERENCE.md**: Complete API documentation
- **PRICING_GUIDE.md**: Pricing model guide

### Support Channels
- **GitHub Issues**: Report bugs and feature requests
- **Development Team**: Contact for urgent issues
- **Community**: Ask questions in discussions

### Log Collection
When reporting issues, include:
1. **Error Messages**: Full error text
2. **Configuration**: Environment variables (without secrets)
3. **Steps to Reproduce**: Exact commands run
4. **System Info**: OS, Go version, etc.
5. **Logs**: Relevant log entries

### Example Bug Report
```
**Issue**: Pricing comparison returns no eligible models

**Steps to Reproduce**:
1. Run: ./bin/dispatch-cli pricing
2. See: All models show "Not Eligible"

**Expected**: At least Standard Pricing should be eligible

**Environment**:
- OS: macOS 14.6.0
- Go: 1.21.0
- CLI Version: 1.0.0

**Configuration**:
- USE_IDP_AUTH=false
- DISPATCH_AUTH_TOKEN=*** (redacted)

**Logs**:
[ERROR] No eligible pricing models found
[DEBUG] Context: {DeliveryCount:1, CustomerTier:bronze, ...}
```

---

*Last updated: $(date)*
