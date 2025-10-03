# Location Services Implementation Plan

## Immediate Next Steps for MCP Project

### 1. Add Location Validation to MCP Server

**File**: `internal/mcp/tools.go`
**Action**: Add new location validation tools

```go
// Add these new tools to the MCP server
func (s *MCPServer) validateAddress(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Implementation for address validation
}

func (s *MCPServer) geocodeAddress(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Implementation for geocoding
}

func (s *MCPServer) checkServiceArea(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    // Implementation for service area checking
}
```

### 2. Create Location Service Client

**File**: `internal/location/client.go`
**Purpose**: Handle all location-related API calls

```go
package location

type Client struct {
    googleAPIKey string
    dispatchAPI   string
}

func (c *Client) ValidateAddress(address string) (*AddressValidation, error)
func (c *Client) GeocodeAddress(address string) (*Coordinates, error)
func (c *Client) CheckServiceArea(coords *Coordinates) (*ServiceAreaStatus, error)
```

### 3. Integrate with Conversation Engine

**File**: `internal/conversation/claude_engine.go`
**Action**: Add location validation to address parsing

```go
func (ce *ClaudeConversationEngine) parseOrderInformation(message string, context *ConversationContext) {
    // Existing parsing logic...
    
    // Add location validation
    if ce.looksLikeAddressInfo(message) {
        // Validate address before storing
        validation, err := ce.locationClient.ValidateAddress(message)
        if err != nil {
            // Handle validation error
            return
        }
        
        // Store validated address with coordinates
        // Update context with validated location data
    }
}
```

### 4. Environment Configuration

**File**: `env.example`
**Add**:
```bash
# Location Services Configuration
GOOGLE_PLACES_API_KEY=your_google_places_api_key
GOOGLE_GEOCODING_API_KEY=your_geocoding_api_key
ENABLE_LOCATION_VALIDATION=true
```

## Implementation Phases

### Phase 1: Basic Address Validation (Week 1)
- [ ] Add Google Places API integration
- [ ] Implement basic address validation
- [ ] Add validation to conversation engine
- [ ] Test with sample addresses

### Phase 2: Service Area Integration (Week 2)
- [ ] Add service area checking
- [ ] Integrate with Dispatch API
- [ ] Add delivery time estimates
- [ ] Handle service area errors

### Phase 3: Advanced Features (Week 3)
- [ ] Add address autocomplete
- [ ] Implement caching
- [ ] Add international support
- [ ] Performance optimization

## Quick Start Implementation

### 1. Add Google Places API Key
```bash
# Add to load-dev-env.sh
export GOOGLE_PLACES_API_KEY="your_api_key_here"
```

### 2. Create Location Service
```go
// internal/location/service.go
type LocationService struct {
    googleClient *googleplaces.Client
}

func (ls *LocationService) ValidateAddress(address string) (*AddressResult, error) {
    // Call Google Places API
    // Return validated address with coordinates
}
```

### 3. Update Conversation Engine
```go
// In parseOrderInformation function
if ce.looksLikeAddressInfo(message) {
    // Validate address
    result, err := ce.locationService.ValidateAddress(message)
    if err != nil {
        // Return error to user
        return
    }
    
    // Store validated address
    // Update context with coordinates
}
```

## Testing Strategy

### 1. Unit Tests
- Test address validation with various formats
- Test geocoding accuracy
- Test service area boundaries

### 2. Integration Tests
- Test with real Google Places API
- Test with Dispatch service areas
- Test error handling

### 3. User Testing
- Test with real user addresses
- Test autocomplete functionality
- Test error messages

## Cost Management

### 1. Caching Strategy
- Cache validated addresses for 24 hours
- Cache service area results
- Implement smart cache invalidation

### 2. Rate Limiting
- Limit API calls per user session
- Implement request queuing
- Monitor usage patterns

### 3. Fallback Handling
- Graceful degradation when APIs fail
- Use cached results when possible
- Clear error messages for users

## Success Metrics

### 1. Technical Metrics
- Address validation accuracy: >95%
- Geocoding success rate: >90%
- API response time: <2 seconds
- Cache hit rate: >60%

### 2. User Experience Metrics
- Reduced address entry errors
- Faster order creation
- Better delivery time estimates
- Improved user satisfaction

## Risk Mitigation

### 1. API Dependencies
- Implement fallback validation
- Cache critical data
- Monitor API health

### 2. Cost Control
- Set usage limits
- Implement cost alerts
- Optimize API calls

### 3. Data Privacy
- No storage of personal addresses
- Secure API key management
- GDPR compliance

## Next Actions

1. **Get Google Places API key** from Google Cloud Console
2. **Add location service** to the MCP project
3. **Integrate with conversation engine** for real-time validation
4. **Test with sample addresses** to verify accuracy
5. **Deploy and monitor** performance

This implementation will significantly improve the MCP project's location handling capabilities and provide a better user experience for address entry and validation.
