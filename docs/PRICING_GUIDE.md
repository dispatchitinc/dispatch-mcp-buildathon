# 💰 Pricing Model Comparison - Quick Reference Guide

## 🚀 Quick Start

### 1. Run the Demo
```bash
./bin/dispatch-cli pricing
```

### 2. Use in AI Agents
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

## 📊 Pricing Models Overview

| Model | Discount | Requirements | Best For |
|-------|----------|--------------|----------|
| **Standard** | 0% | None | New customers |
| **Multi-Delivery** | 15% | 2+ deliveries | Multiple stops |
| **Volume** | 20% | 5+ deliveries + 3+ orders/month | Regular customers |
| **Loyalty** | 10% | Gold tier | VIP customers |
| **Bulk Order** | 25% | 10+ deliveries + bulk flag | Large orders |

## 🎯 Customer Scenarios

### Scenario 1: New Customer
- **Profile**: 1 delivery, bronze tier, 1 order/month
- **Best Model**: Standard Pricing
- **Savings**: $0.00 (0%)

### Scenario 2: Multi-Delivery Customer  
- **Profile**: 2 deliveries, silver tier, 3 orders/month
- **Best Model**: Multi-Delivery Discount
- **Savings**: $6.90 (15%)

### Scenario 3: High-Volume Customer
- **Profile**: 5 deliveries, gold tier, 8 orders/month
- **Best Model**: Multi-Delivery Discount
- **Savings**: $11.20 (24.4%)

### Scenario 4: Bulk Order Customer
- **Profile**: 10 deliveries, gold tier, 15 orders/month, bulk order
- **Best Model**: Multi-Delivery Discount
- **Savings**: $15.11 (32.8%)

## 🔧 MCP Tool Parameters

### Required Parameters
- `original_estimate`: JSON string of estimate data

### Optional Parameters
- `delivery_count`: Number of deliveries (default: 1)
- `customer_tier`: bronze, silver, gold (default: bronze)
- `order_frequency`: Orders per month (default: 1)
- `total_order_value`: Total order value (default: original cost)
- `is_bulk_order`: true/false (default: false)

## 📈 Business Impact

### Revenue Optimization
- **Customer Segmentation**: Different pricing for different customer types
- **Upselling**: Show customers how to unlock better pricing
- **Retention**: Demonstrate value of loyalty programs

### Sales Enablement
- **Pricing Transparency**: Show potential savings upfront
- **Value Demonstration**: Prove ROI of volume commitments
- **Competitive Advantage**: Flexible pricing models

## 🛠 Technical Implementation

### Go Code Example
```go
engine := pricing.NewPricingEngine()
context := pricing.PricingContext{
    DeliveryCount:    3,
    CustomerTier:     "gold",
    OrderFrequency:   5,
    TotalOrderValue:  150.00,
    IsBulkOrder:      false,
}
comparison := engine.ComparePricingModels(estimate, context)
```

### Response Structure
```json
{
  "original_estimate": { /* estimate data */ },
  "pricing_models": [
    {
      "model": "multi_delivery",
      "name": "Multi-Delivery Discount",
      "original_cost": 45.99,
      "adjusted_cost": 39.09,
      "discount": 6.90,
      "discount_percent": 15.0,
      "savings": 6.90,
      "eligible": true
    }
  ],
  "best_option": { /* best pricing model */ },
  "savings": 6.90,
  "savings_percentage": 15.0
}
```

## 🎨 Customization

### Adding New Pricing Models
1. Add model type to `PricingModel` enum
2. Add rule in `initializeDefaultRules()`
3. Update eligibility logic in `isEligibleForModel()`
4. Add custom discount logic in `calculateAdditionalDiscount()`

### Modifying Existing Models
- **Discount Rates**: Change `BaseMultiplier` in pricing rules
- **Eligibility**: Update thresholds in `isEligibleForModel()`
- **Additional Discounts**: Modify `calculateAdditionalDiscount()`

## 🔍 Troubleshooting

### Common Issues

**Q: No eligible pricing models found**
- Check customer tier and delivery count
- Verify order frequency meets requirements
- Ensure bulk order flag is set correctly

**Q: Pricing seems too low**
- Check minimum price protection (50% of original)
- Verify discount calculations
- Review additional discount logic

**Q: MCP tool not working**
- Ensure original_estimate is valid JSON
- Check parameter types (strings for numbers)
- Verify authentication is working

### Debug Mode
```bash
# Enable debug logging
export DEBUG=true
./bin/dispatch-cli pricing
```

## 📞 Support

- **Documentation**: See main README.md
- **Examples**: Check `samples/` directory
- **Issues**: Report on GitHub issues
- **Questions**: Contact development team

---

*Last updated: $(date)*
