# Monolith-Aligned Conversational Order Creation Flow

## Overview
Our chat flow should work exactly like the ms-monolith order creation steps, but make them conversational and interactive. This document outlines how to align our current implementation with the monolith's proven workflow.

## Current Monolith Flow vs Our Chat Flow

### **MS Monolith Steps:**
1. **Multi-Stop Selection** (`/orders/begin`) - Choose single vs multi-stop
2. **Pickup Step** (`/orders/begin/pickup`) - Pickup location and contact
3. **Drop-off Step** (`/orders/begin/drop_off`) - Delivery locations and packages
4. **Vehicle Step** (`/orders/begin/vehicle`) - Vehicle type selection
5. **Add-ons Step** (`/orders/begin/add_ons`) - Capabilities and special requirements
6. **Delivery Step** (`/orders/begin/delivery`) - Scheduling and logistics
7. **Review Step** (`/orders/begin/review`) - Final review and submission

### **Our Current Chat Flow:**
- ✅ **Pickup Information** - Working well
- ✅ **Delivery Information** - Working well  
- ❌ **Missing**: Vehicle selection step
- ❌ **Missing**: Add-ons/capabilities step
- ❌ **Missing**: Delivery scheduling step
- ❌ **Missing**: Review and submission step

## Implementation Plan: Align Chat with Monolith Steps

### **Step 1: Enhanced Order Creation State Management**

```go
// Update OrderCreationState to match monolith steps exactly
type OrderCreationState struct {
    InProgress           bool                                   `json:"in_progress"`
    Step                 string                                 `json:"step"`             // "multi_stop", "pickup", "drop_off", "vehicle", "add_ons", "delivery", "review"
    CurrentQuestion      string                                 `json:"current_question"`
    PickupInfo           *dispatch.CreateOrderPickupInfoInput   `json:"pickup_info,omitempty"`
    DropOffs             []dispatch.CreateOrderDropOffInfoInput `json:"drop_offs,omitempty"`
    VehicleType          *VehicleTypeInfo                       `json:"vehicle_type,omitempty"`
    Capabilities         []string                               `json:"capabilities,omitempty"`
    DeliveryInfo         *dispatch.DeliveryInfoInput            `json:"delivery_info,omitempty"`
    SchedulingInfo       *SchedulingInfo                        `json:"scheduling_info,omitempty"`
    MissingFields        []string                               `json:"missing_fields"`
    CompletedFields      []string                               `json:"completed_fields"`
    CurrentDeliveryIndex int                                    `json:"current_delivery_index"`
    ValidationErrors     []string                               `json:"validation_errors,omitempty"`
}

// New supporting types
type VehicleTypeInfo struct {
    VehicleTypeID   string `json:"vehicle_type_id"`
    VehicleTypeName string `json:"vehicle_type_name"`
    CustomTypes     []string `json:"custom_types,omitempty"`
}

type SchedulingInfo struct {
    PickupTime      string `json:"pickup_time"`
    DeliveryTime    string `json:"delivery_time"`
    PickupDate      string `json:"pickup_date"`
    DeliveryDate    string `json:"delivery_date"`
    TimeZone        string `json:"time_zone"`
    SpecialTiming   []string `json:"special_timing,omitempty"`
}
```

### **Step 2: Conversational Step Management**

```go
// Add step management to conversation engine
func (ce *ClaudeConversationEngine) determineNextStep(context *ConversationContext) string {
    if !context.OrderCreation.InProgress {
        return "multi_stop" // Start with multi-stop selection
    }
    
    // Check what information we have and what's missing
    if context.OrderCreation.PickupInfo == nil {
        return "pickup"
    }
    
    if len(context.OrderCreation.DropOffs) == 0 {
        return "drop_off"
    }
    
    if context.OrderCreation.VehicleType == nil {
        return "vehicle"
    }
    
    if len(context.OrderCreation.Capabilities) == 0 {
        return "add_ons"
    }
    
    if context.OrderCreation.SchedulingInfo == nil {
        return "delivery"
    }
    
    return "review"
}

func (ce *ClaudeConversationEngine) getStepPrompt(step string, context *ConversationContext) string {
    switch step {
    case "multi_stop":
        return "I'll help you create a delivery order. First, let me know: do you need to deliver to one location or multiple locations?"
    
    case "pickup":
        return "Great! Now I need your pickup information. What's the business name and address where we'll be picking up the package?"
    
    case "drop_off":
        return "Perfect! Now I need the delivery information. Where should we deliver this package?"
    
    case "vehicle":
        return "What type of vehicle do you need for this delivery? (cargo van, pickup truck, etc.)"
    
    case "add_ons":
        return "Do you need any special services or capabilities for this delivery? (temperature control, white glove service, etc.)"
    
    case "delivery":
        return "When would you like this delivered? I need pickup time and delivery time."
    
    case "review":
        return "Let me review your order details before we create it..."
    
    default:
        return "I'm ready to help you create a delivery order. What do you need?"
    }
}
```

### **Step 3: Enhanced Information Parsing**

```go
// Add parsing for each monolith step
func (ce *ClaudeConversationEngine) parseMultiStopInfo(message string, context *ConversationContext) {
    // Parse whether user wants single or multi-stop delivery
    if strings.Contains(strings.ToLower(message), "multiple") || 
       strings.Contains(strings.ToLower(message), "several") ||
       strings.Contains(strings.ToLower(message), "many") {
        context.OrderCreation.MultiStop = true
    } else {
        context.OrderCreation.MultiStop = false
    }
}

func (ce *ClaudeConversationEngine) parseVehicleInfo(message string, context *ConversationContext) {
    // Parse vehicle type from message
    vehicleTypes := map[string]string{
        "cargo van": "cargo_van",
        "pickup truck": "pickup_truck", 
        "box truck": "box_truck",
        "van": "cargo_van",
        "truck": "pickup_truck",
    }
    
    messageLower := strings.ToLower(message)
    for keyword, vehicleType := range vehicleTypes {
        if strings.Contains(messageLower, keyword) {
            context.OrderCreation.VehicleType = &VehicleTypeInfo{
                VehicleTypeID: vehicleType,
                VehicleTypeName: keyword,
            }
            break
        }
    }
}

func (ce *ClaudeConversationEngine) parseCapabilitiesInfo(message string, context *ConversationContext) {
    // Parse special requirements and capabilities
    capabilities := []string{}
    
    if strings.Contains(strings.ToLower(message), "temperature") {
        capabilities = append(capabilities, "temperature_control")
    }
    if strings.Contains(strings.ToLower(message), "white glove") {
        capabilities = append(capabilities, "white_glove_service")
    }
    if strings.Contains(strings.ToLower(message), "fragile") {
        capabilities = append(capabilities, "fragile_handling")
    }
    if strings.Contains(strings.ToLower(message), "signature") {
        capabilities = append(capabilities, "signature_required")
    }
    
    context.OrderCreation.Capabilities = capabilities
}

func (ce *ClaudeConversationEngine) parseSchedulingInfo(message string, context *ConversationContext) {
    // Parse scheduling information
    // This would use a more sophisticated time parsing library
    // For now, simple keyword matching
    
    scheduling := &SchedulingInfo{}
    
    // Look for time patterns
    timePattern := regexp.MustCompile(`(\d{1,2}):?(\d{2})?\s*(am|pm|AM|PM)?`)
    matches := timePattern.FindAllString(message, -1)
    
    if len(matches) > 0 {
        scheduling.PickupTime = matches[0]
        if len(matches) > 1 {
            scheduling.DeliveryTime = matches[1]
        }
    }
    
    context.OrderCreation.SchedulingInfo = scheduling
}
```

### **Step 4: Conversational Flow Implementation**

```go
// Update the main conversation processing to follow monolith steps
func (ce *ClaudeConversationEngine) ProcessMessageWithHistory(message string, context *ConversationContext, history []ConversationMessage) (*ConversationResponse, error) {
    // Update context from message
    updatedContext := ce.updateContextFromMessage(message, context)
    
    // Determine current step
    currentStep := ce.determineNextStep(updatedContext)
    updatedContext.OrderCreation.Step = currentStep
    
    // Handle validation errors first
    validationErrorMsg := ce.handleValidationErrors(updatedContext)
    if validationErrorMsg != "" {
        return &ConversationResponse{
            Message:         validationErrorMsg,
            Recommendations: []PricingRecommendation{},
            NextQuestions:   []string{},
            UpdatedContext:  updatedContext,
        }, nil
    }
    
    // Generate response based on current step
    response := ce.generateStepResponse(currentStep, updatedContext, message)
    
    return &ConversationResponse{
        Message:         response,
        Recommendations: []PricingRecommendation{},
        NextQuestions:   []string{},
        UpdatedContext:  updatedContext,
    }, nil
}

func (ce *ClaudeConversationEngine) generateStepResponse(step string, context *ConversationContext, message string) string {
    switch step {
    case "multi_stop":
        return ce.handleMultiStopStep(message, context)
    case "pickup":
        return ce.handlePickupStep(message, context)
    case "drop_off":
        return ce.handleDropOffStep(message, context)
    case "vehicle":
        return ce.handleVehicleStep(message, context)
    case "add_ons":
        return ce.handleAddOnsStep(message, context)
    case "delivery":
        return ce.handleDeliveryStep(message, context)
    case "review":
        return ce.handleReviewStep(message, context)
    default:
        return "I'm ready to help you create a delivery order. What do you need?"
    }
}
```

### **Step 5: Step-Specific Handlers**

```go
func (ce *ClaudeConversationEngine) handleMultiStopStep(message string, context *ConversationContext) string {
    ce.parseMultiStopInfo(message, context)
    
    if context.OrderCreation.MultiStop {
        return "Perfect! I'll help you set up a multi-stop delivery. Let's start with your pickup location - what's the business name and address?"
    } else {
        return "Great! I'll help you set up a single delivery. Let's start with your pickup location - what's the business name and address?"
    }
}

func (ce *ClaudeConversationEngine) handlePickupStep(message string, context *ConversationContext) string {
    if ce.looksLikeAddressInfo(message) {
        ce.parsePickupInfo(message, context)
        return "Perfect! I have your pickup information. Now I need your delivery location - where should we deliver this package?"
    }
    return "I need your pickup location details. Please provide the business name, address, contact name, and phone number."
}

func (ce *ClaudeConversationEngine) handleDropOffStep(message string, context *ConversationContext) string {
    if ce.looksLikeAddressInfo(message) {
        ce.parseDeliveryInfo(message, context)
        
        if context.OrderCreation.MultiStop {
            return "Got it! Do you have another delivery location, or are we ready to move on to vehicle selection?"
        } else {
            return "Perfect! Now I need to know what type of vehicle you need for this delivery."
        }
    }
    return "I need your delivery location details. Please provide the business name, address, contact name, and phone number."
}

func (ce *ClaudeConversationEngine) handleVehicleStep(message string, context *ConversationContext) string {
    ce.parseVehicleInfo(message, context)
    
    if context.OrderCreation.VehicleType != nil {
        return "Great! I've noted you need a " + context.OrderCreation.VehicleType.VehicleTypeName + ". Do you need any special services or capabilities for this delivery?"
    }
    return "What type of vehicle do you need? (cargo van, pickup truck, box truck, etc.)"
}

func (ce *ClaudeConversationEngine) handleAddOnsStep(message string, context *ConversationContext) string {
    ce.parseCapabilitiesInfo(message, context)
    
    return "Perfect! Now I need to know your delivery schedule. When would you like pickup and delivery?"
}

func (ce *ClaudeConversationEngine) handleDeliveryStep(message string, context *ConversationContext) string {
    ce.parseSchedulingInfo(message, context)
    
    if context.OrderCreation.SchedulingInfo != nil {
        return "Excellent! Let me review your order details before we create the order..."
    }
    return "When would you like this delivered? I need pickup time and delivery time."
}

func (ce *ClaudeConversationEngine) handleReviewStep(message string, context *ConversationContext) string {
    // Generate order summary
    summary := ce.generateOrderSummary(context)
    
    if strings.Contains(strings.ToLower(message), "yes") || 
       strings.Contains(strings.ToLower(message), "create") ||
       strings.Contains(strings.ToLower(message), "confirm") {
        // Create the order
        return "Perfect! I'm creating your order now. You'll receive a confirmation shortly."
    }
    
    return summary + "\n\nShould I create this order for you?"
}
```

## Implementation Priority

### **Phase 1: Core Step Alignment** (Immediate)
1. ✅ **Pickup Step** - Already working well
2. ✅ **Drop-off Step** - Already working well  
3. 🔄 **Vehicle Step** - Add vehicle type selection
4. 🔄 **Review Step** - Add order summary and confirmation

### **Phase 2: Advanced Steps** (Next)
1. **Multi-Stop Step** - Handle single vs multi-stop selection
2. **Add-ons Step** - Handle capabilities and special requirements  
3. **Delivery Step** - Handle scheduling and logistics

### **Phase 3: Integration** (Final)
1. **Order Creation** - Integrate with ms-monolith services
2. **Validation** - Use monolith validation logic
3. **Pricing** - Use monolith pricing engine

## Benefits of Monolith-Aligned Flow

### **1. Familiar User Experience**
- Users already know the monolith flow
- Chat provides the same information collection
- More intuitive and predictable

### **2. Proven Workflow**
- Monolith flow is battle-tested
- Handles edge cases and complex scenarios
- Well-integrated with existing systems

### **3. Easier Integration**
- Direct mapping to monolith services
- Reuse existing validation and business logic
- Consistent data structures

### **4. Better Error Handling**
- Same validation rules as monolith
- Consistent error messages
- Proven error recovery patterns

## Next Steps

1. **Implement Vehicle Step** - Add vehicle type selection to conversation
2. **Add Review Step** - Generate order summary and confirmation
3. **Enhance Step Management** - Better step transitions
4. **Integrate with Monolith** - Use existing services for order creation

This approach ensures our conversational interface works exactly like the proven monolith flow, just more interactive and user-friendly!
