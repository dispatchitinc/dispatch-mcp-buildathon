# Complete Order Collection Plan

## Overview
This document outlines the enhanced conversation flow needed to collect all required information for complete order creation, based on the ms-monolith database schema analysis.

## Current vs. Required Information Collection

### **Phase 1: Basic Information (✅ Currently Implemented)**
- Multi-stop selection
- Pickup location and contact
- Delivery location(s) and contact
- Vehicle type selection
- Basic scheduling

### **Phase 2: Package Details (❌ Missing - Critical)**
```
🤖 "Now I need details about your packages. How many packages are you shipping?"

User: "3 packages"

🤖 "Great! What's the total weight of all packages, or the weight of each individual package?"

User: "Total weight is about 50 pounds"

🤖 "Do you have reference names for these packages? This helps with tracking."

User: "Package A, Package B, Package C"

🤖 "Are there any special handling requirements? (fragile, temperature sensitive, etc.)"

User: "One package is fragile"
```

### **Phase 3: Organization Context (❌ Missing - Critical)**
```
🤖 "Which organization is this order for?"

User: "Acme Corporation"

🤖 "Which branch or division? (if applicable)"

User: "Main Office"

🤖 "What's your name and email for order notifications?"

User: "John Smith, john@acme.com"
```

### **Phase 4: Advanced Scheduling (❌ Missing - Important)**
```
🤖 "When do you need this picked up and delivered?"

User: "Pickup today by 2 PM, deliver by 6 PM"

🤖 "Is this an ASAP delivery or do you have flexibility?"

User: "ASAP"

🤖 "What time zone are we working with?"

User: "Pacific Time"
```

### **Phase 5: Service Level Options (❌ Missing - Critical)**
```
🤖 "What service level do you need?"

Options:
- Standard (next business day)
- Express (same day)
- Rush (within 4 hours)
- Scheduled (specific time)

User: "Express"

🤖 "Do you need any special services?"
- Signature required
- Photo proof of delivery
- White glove service
- Temperature control
- Unloading assistance

User: "Signature required and unloading assistance"
```

### **Phase 6: Special Requirements (❌ Missing - Important)**
```
🤖 "Any special delivery instructions?"

User: "Call when 15 minutes away, use side entrance"

🤖 "Do you need a dedicated vehicle (not shared with other deliveries)?"

User: "No, shared is fine"

🤖 "Any access instructions for the delivery location?"

User: "Gate code 1234, building B"
```

### **Phase 7: Billing Information (❌ Missing - Essential)**
```
🤖 "How would you like to pay for this delivery?"

Options:
- Bill to account
- Credit card
- Invoice
- COD (Cash on Delivery)

User: "Bill to account"

🤖 "What's the billing address? (if different from pickup)"

User: "Same as pickup"
```

## Enhanced Data Structures

### **Updated OrderCreationState**
```go
type OrderCreationState struct {
    // Basic info (existing)
    InProgress           bool
    Step                 string
    CurrentQuestion      string
    MultiStop            bool
    PickupInfo           *dispatch.CreateOrderPickupInfoInput
    DropOffs             []dispatch.CreateOrderDropOffInfoInput
    VehicleType          *VehicleTypeInfo
    Capabilities         []string
    DeliveryInfo         *dispatch.DeliveryInfoInput
    SchedulingInfo       *SchedulingInfo
    
    // NEW: Package Details
    PackageDetails       *PackageDetailsInfo
    
    // NEW: Organization Context
    OrganizationInfo     *OrganizationInfo
    
    // NEW: Service Level
    ServiceLevel         *ServiceLevelInfo
    
    // NEW: Billing Information
    BillingInfo          *BillingInfo
    
    // NEW: Special Requirements
    SpecialRequirements  *SpecialRequirementsInfo
    
    // State management
    MissingFields        []string
    CompletedFields      []string
    CurrentDeliveryIndex int
    ValidationErrors     []string
}

type PackageDetailsInfo struct {
    PackageCount         int      `json:"package_count"`
    TotalWeight          float64  `json:"total_weight"`
    IndividualWeights    []float64 `json:"individual_weights,omitempty"`
    ReferenceNames       []string `json:"reference_names,omitempty"`
    SpecialHandling      []string `json:"special_handling,omitempty"`
    Dimensions           *PackageDimensions `json:"dimensions,omitempty"`
}

type PackageDimensions struct {
    Length float64 `json:"length"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
    Unit   string  `json:"unit"` // "inches", "cm"
}

type OrganizationInfo struct {
    OrganizationID       string `json:"organization_id"`
    OrganizationName     string `json:"organization_name"`
    BranchID            string `json:"branch_id,omitempty"`
    BranchName          string `json:"branch_name,omitempty"`
    MarketID            string `json:"market_id,omitempty"`
    MarketName          string `json:"market_name,omitempty"`
    CreatedByUserID     string `json:"created_by_user_id"`
    CreatedByUserName   string `json:"created_by_user_name"`
    CreatedByUserEmail  string `json:"created_by_user_email"`
}

type ServiceLevelInfo struct {
    ServiceLevel         string   `json:"service_level"` // "standard", "express", "rush", "scheduled"
    DeliveryTimeWindow   string   `json:"delivery_time_window,omitempty"`
    UrgencyLevel         string   `json:"urgency_level"` // "asap", "scheduled", "flexible"
    SpecialServices      []string `json:"special_services,omitempty"`
}

type BillingInfo struct {
    BillingMethod        string `json:"billing_method"` // "account", "credit_card", "invoice", "cod"
    BillingAddress       *AddressInput `json:"billing_address,omitempty"`
    PaymentTerms         string `json:"payment_terms,omitempty"`
    ContactEmail         string `json:"contact_email"`
    NotificationPhone    string `json:"notification_phone,omitempty"`
}

type SpecialRequirementsInfo struct {
    UnloadingAssistance  bool     `json:"unloading_assistance"`
    DedicatedVehicle    bool     `json:"dedicated_vehicle"`
    DeliveryInstructions string  `json:"delivery_instructions,omitempty"`
    AccessInstructions   string  `json:"access_instructions,omitempty"`
    SpecialNotes        string  `json:"special_notes,omitempty"`
}
```

## Implementation Steps

### **Step 1: Update Conversation Engine**
- Add new parsing functions for each information type
- Update step management to include new phases
- Add validation for all required fields

### **Step 2: Enhance Data Collection**
- Implement package details collection
- Add organization context gathering
- Include service level selection
- Add billing information collection

### **Step 3: Update GraphQL Integration**
- Extend OrderCreationInput to include all new fields
- Update validation rules
- Enhance error handling

### **Step 4: Improve User Experience**
- Add progress indicators
- Implement smart defaults
- Add field validation with helpful error messages
- Create summary review before submission

## Benefits of Complete Information Collection

1. **Complete Orders**: All required fields populated
2. **Better Pricing**: Accurate pricing based on complete information
3. **Improved Service**: Better delivery experience with all details
4. **Reduced Errors**: Fewer order failures due to missing information
5. **Enhanced Tracking**: Better package tracking with reference names
6. **Proper Billing**: Correct billing and payment processing

## Next Steps

1. Implement package details collection
2. Add organization context gathering
3. Include service level selection
4. Add billing information collection
5. Update validation and error handling
6. Test complete order creation flow
