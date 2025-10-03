package pricing

import (
	"dispatch-mcp-server/internal/dispatch"
	"fmt"
	"math"
)

// PricingModel represents different pricing strategies
type PricingModel string

const (
	StandardPricing        PricingModel = "standard"
	MultiDeliveryPricing   PricingModel = "multi_delivery"
	VolumeDiscountPricing  PricingModel = "volume_discount"
	LoyaltyDiscountPricing PricingModel = "loyalty_discount"
	BulkOrderPricing       PricingModel = "bulk_order"
)

// PricingRule defines how a pricing model should be applied
type PricingRule struct {
	Model           PricingModel `json:"model"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	BaseMultiplier  float64      `json:"base_multiplier"`  // Base price multiplier (1.0 = no change)
	MinDiscount     float64      `json:"min_discount"`     // Minimum discount percentage
	MaxDiscount     float64      `json:"max_discount"`     // Maximum discount percentage
	VolumeThreshold int          `json:"volume_threshold"` // Minimum deliveries for discount
	LoyaltyTier     string       `json:"loyalty_tier"`     // Required loyalty tier
}

// PricingComparison represents the result of comparing pricing models
type PricingComparison struct {
	OriginalEstimate  *dispatch.AvailableOrderOption `json:"original_estimate"`
	PricingModels     []PricingResult                `json:"pricing_models"`
	BestOption        *PricingResult                 `json:"best_option"`
	Savings           float64                        `json:"savings"`
	SavingsPercentage float64                        `json:"savings_percentage"`
}

// PricingResult represents the result of applying a specific pricing model
type PricingResult struct {
	Model           PricingModel `json:"model"`
	Name            string       `json:"name"`
	OriginalCost    float64      `json:"original_cost"`
	AdjustedCost    float64      `json:"adjusted_cost"`
	Discount        float64      `json:"discount"`
	DiscountPercent float64      `json:"discount_percent"`
	Savings         float64      `json:"savings"`
	Eligible        bool         `json:"eligible"`
	Reason          string       `json:"reason,omitempty"`
}

// PricingEngine handles pricing model calculations
type PricingEngine struct {
	rules map[PricingModel]PricingRule
}

// NewPricingEngine creates a new pricing engine with default rules
func NewPricingEngine() *PricingEngine {
	engine := &PricingEngine{
		rules: make(map[PricingModel]PricingRule),
	}

	// Initialize default pricing rules
	engine.initializeDefaultRules()

	return engine
}

// initializeDefaultRules sets up default pricing models
func (pe *PricingEngine) initializeDefaultRules() {
	pe.rules[StandardPricing] = PricingRule{
		Model:          StandardPricing,
		Name:           "Standard Pricing",
		Description:    "Standard pricing with no discounts",
		BaseMultiplier: 1.0,
		MinDiscount:    0.0,
		MaxDiscount:    0.0,
	}

	pe.rules[MultiDeliveryPricing] = PricingRule{
		Model:           MultiDeliveryPricing,
		Name:            "Multi-Delivery Discount",
		Description:     "Discount for multiple deliveries in the same order",
		BaseMultiplier:  0.85, // 15% discount
		MinDiscount:     5.0,  // 5% minimum
		MaxDiscount:     25.0, // 25% maximum
		VolumeThreshold: 2,    // 2+ deliveries
	}

	pe.rules[VolumeDiscountPricing] = PricingRule{
		Model:           VolumeDiscountPricing,
		Name:            "Volume Discount",
		Description:     "Discount based on order volume and frequency",
		BaseMultiplier:  0.80, // 20% discount
		MinDiscount:     10.0, // 10% minimum
		MaxDiscount:     30.0, // 30% maximum
		VolumeThreshold: 5,    // 5+ deliveries
	}

	pe.rules[LoyaltyDiscountPricing] = PricingRule{
		Model:          LoyaltyDiscountPricing,
		Name:           "Loyalty Discount",
		Description:    "Discount for loyal customers",
		BaseMultiplier: 0.90, // 10% discount
		MinDiscount:    5.0,  // 5% minimum
		MaxDiscount:    15.0, // 15% maximum
		LoyaltyTier:    "gold",
	}

	pe.rules[BulkOrderPricing] = PricingRule{
		Model:           BulkOrderPricing,
		Name:            "Bulk Order Discount",
		Description:     "Discount for large bulk orders",
		BaseMultiplier:  0.75, // 25% discount
		MinDiscount:     15.0, // 15% minimum
		MaxDiscount:     40.0, // 40% maximum
		VolumeThreshold: 10,   // 10+ deliveries
	}
}

// ComparePricingModels compares different pricing models against an original estimate
func (pe *PricingEngine) ComparePricingModels(originalEstimate *dispatch.AvailableOrderOption, context PricingContext) *PricingComparison {
	comparison := &PricingComparison{
		OriginalEstimate: originalEstimate,
		PricingModels:    []PricingResult{},
	}

	originalCost := originalEstimate.EstimatedOrderCost

	// Apply each pricing model
	for _, rule := range pe.rules {
		result := pe.applyPricingModel(originalCost, rule, context)
		comparison.PricingModels = append(comparison.PricingModels, result)
	}

	// Find the best option (lowest cost)
	bestOption := pe.findBestOption(comparison.PricingModels)
	comparison.BestOption = bestOption

	if bestOption != nil {
		comparison.Savings = originalCost - bestOption.AdjustedCost
		comparison.SavingsPercentage = (comparison.Savings / originalCost) * 100
	}

	return comparison
}

// PricingContext provides context for pricing calculations
type PricingContext struct {
	DeliveryCount     int     `json:"delivery_count"`
	CustomerTier      string  `json:"customer_tier"`
	OrderFrequency    int     `json:"order_frequency"` // orders per month
	TotalOrderValue   float64 `json:"total_order_value"`
	IsBulkOrder       bool    `json:"is_bulk_order"`
	OrganizationDruid string  `json:"organization_druid"`
}

// applyPricingModel applies a specific pricing model to calculate adjusted cost
func (pe *PricingEngine) applyPricingModel(originalCost float64, rule PricingRule, context PricingContext) PricingResult {
	result := PricingResult{
		Model:        rule.Model,
		Name:         rule.Name,
		OriginalCost: originalCost,
		Eligible:     true,
	}

	// Check eligibility based on context
	if !pe.isEligibleForModel(rule, context) {
		result.Eligible = false
		result.AdjustedCost = originalCost
		result.Discount = 0.0
		result.DiscountPercent = 0.0
		result.Savings = 0.0
		result.Reason = pe.getIneligibilityReason(rule, context)
		return result
	}

	// Calculate adjusted cost
	adjustedCost := originalCost * rule.BaseMultiplier

	// Apply additional discounts based on context
	additionalDiscount := pe.calculateAdditionalDiscount(rule, context)
	if additionalDiscount > 0 {
		discountAmount := adjustedCost * (additionalDiscount / 100)
		adjustedCost -= discountAmount
	}

	// Ensure we don't go below minimum cost (e.g., 50% of original)
	minCost := originalCost * 0.5
	if adjustedCost < minCost {
		adjustedCost = minCost
	}

	result.AdjustedCost = adjustedCost
	result.Discount = originalCost - adjustedCost
	result.DiscountPercent = (result.Discount / originalCost) * 100
	result.Savings = result.Discount

	return result
}

// isEligibleForModel checks if the context makes the customer eligible for a pricing model
func (pe *PricingEngine) isEligibleForModel(rule PricingRule, context PricingContext) bool {
	switch rule.Model {
	case StandardPricing:
		return true
	case MultiDeliveryPricing:
		return context.DeliveryCount >= rule.VolumeThreshold
	case VolumeDiscountPricing:
		return context.DeliveryCount >= rule.VolumeThreshold && context.OrderFrequency >= 3
	case LoyaltyDiscountPricing:
		return context.CustomerTier == rule.LoyaltyTier
	case BulkOrderPricing:
		return context.IsBulkOrder && context.DeliveryCount >= rule.VolumeThreshold
	default:
		return false
	}
}

// getIneligibilityReason returns a human-readable reason why a model is not eligible
func (pe *PricingEngine) getIneligibilityReason(rule PricingRule, context PricingContext) string {
	switch rule.Model {
	case MultiDeliveryPricing:
		return fmt.Sprintf("Requires %d+ deliveries, you have %d", rule.VolumeThreshold, context.DeliveryCount)
	case VolumeDiscountPricing:
		return fmt.Sprintf("Requires %d+ deliveries and 3+ orders/month, you have %d deliveries and %d orders/month",
			rule.VolumeThreshold, context.DeliveryCount, context.OrderFrequency)
	case LoyaltyDiscountPricing:
		return fmt.Sprintf("Requires %s tier, you are %s", rule.LoyaltyTier, context.CustomerTier)
	case BulkOrderPricing:
		return fmt.Sprintf("Requires bulk order with %d+ deliveries, you have %d", rule.VolumeThreshold, context.DeliveryCount)
	default:
		return "Not eligible for this pricing model"
	}
}

// calculateAdditionalDiscount calculates additional discounts based on context
func (pe *PricingEngine) calculateAdditionalDiscount(rule PricingRule, context PricingContext) float64 {
	additionalDiscount := 0.0

	// Volume-based additional discount
	if context.DeliveryCount > rule.VolumeThreshold {
		extraDeliveries := context.DeliveryCount - rule.VolumeThreshold
		additionalDiscount += float64(extraDeliveries) * 2.0 // 2% per extra delivery
	}

	// Frequency-based additional discount
	if context.OrderFrequency > 5 {
		additionalDiscount += 5.0 // 5% for high frequency customers
	}

	// Value-based additional discount
	if context.TotalOrderValue > 1000 {
		additionalDiscount += 3.0 // 3% for high value orders
	}

	// Cap at maximum discount
	if additionalDiscount > rule.MaxDiscount {
		additionalDiscount = rule.MaxDiscount
	}

	return additionalDiscount
}

// findBestOption finds the pricing model with the lowest cost
func (pe *PricingEngine) findBestOption(results []PricingResult) *PricingResult {
	var best *PricingResult
	lowestCost := math.Inf(1)

	for i := range results {
		if results[i].Eligible && results[i].AdjustedCost < lowestCost {
			lowestCost = results[i].AdjustedCost
			best = &results[i]
		}
	}

	return best
}

// GetAvailableModels returns all available pricing models
func (pe *PricingEngine) GetAvailableModels() []PricingRule {
	var models []PricingRule
	for _, rule := range pe.rules {
		models = append(models, rule)
	}
	return models
}
