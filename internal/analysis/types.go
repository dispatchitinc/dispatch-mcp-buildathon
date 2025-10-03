package analysis

import (
	"time"
)

// HistoricalOrder represents a historical order for analysis
type HistoricalOrder struct {
	ID                string    `json:"id"`
	OrderDate         time.Time `json:"order_date"`
	DeliveryCount     int       `json:"delivery_count"`
	TotalCost         float64   `json:"total_cost"`
	PickupLocation    string    `json:"pickup_location"`
	DeliveryLocations []string  `json:"delivery_locations"`
	CustomerTier      string    `json:"customer_tier"`
	OrderFrequency    int       `json:"order_frequency"` // orders per month
	IsBulkOrder       bool      `json:"is_bulk_order"`
	PricingModel      string    `json:"pricing_model"` // standard, multi_delivery, volume, loyalty, bulk
}

// AnalysisPeriod represents the time period for analysis
type AnalysisPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// BundlingAnalysis represents the analysis of bundling orders
type BundlingAnalysis struct {
	CurrentOrders     int      `json:"current_orders"`
	OptimizedOrders   int      `json:"optimized_orders"`
	CurrentCost       float64  `json:"current_cost"`
	OptimizedCost     float64  `json:"optimized_cost"`
	PotentialSavings  float64  `json:"potential_savings"`
	SavingsPercentage float64  `json:"savings_percentage"`
	Recommendations   []string `json:"recommendations"`
}

// VolumeAnalysis represents the analysis of volume discounts
type VolumeAnalysis struct {
	CurrentFrequency   int      `json:"current_frequency"` // orders per month
	TargetFrequency    int      `json:"target_frequency"`  // orders per month for volume discount
	CurrentCost        float64  `json:"current_cost"`
	VolumeDiscountCost float64  `json:"volume_discount_cost"`
	PotentialSavings   float64  `json:"potential_savings"`
	SavingsPercentage  float64  `json:"savings_percentage"`
	Recommendations    []string `json:"recommendations"`
}

// LoyaltyAnalysis represents the analysis of loyalty tier benefits
type LoyaltyAnalysis struct {
	CurrentTier         string   `json:"current_tier"`
	TargetTier          string   `json:"target_tier"`
	CurrentCost         float64  `json:"current_cost"`
	LoyaltyDiscountCost float64  `json:"loyalty_discount_cost"`
	PotentialSavings    float64  `json:"potential_savings"`
	SavingsPercentage   float64  `json:"savings_percentage"`
	Recommendations     []string `json:"recommendations"`
}

// ComprehensiveAnalysis represents the complete analysis combining all strategies
type ComprehensiveAnalysis struct {
	AnalysisPeriod   AnalysisPeriod `json:"analysis_period"`
	TotalOrders      int            `json:"total_orders"`
	TotalDeliveries  int            `json:"total_deliveries"`
	CurrentTotalCost float64        `json:"current_total_cost"`

	BundlingAnalysis *BundlingAnalysis `json:"bundling_analysis,omitempty"`
	VolumeAnalysis   *VolumeAnalysis   `json:"volume_analysis,omitempty"`
	LoyaltyAnalysis  *LoyaltyAnalysis  `json:"loyalty_analysis,omitempty"`

	CombinedSavings           float64  `json:"combined_savings"`
	CombinedSavingsPercentage float64  `json:"combined_savings_percentage"`
	ImplementationTimeline    string   `json:"implementation_timeline"`
	ROI                       float64  `json:"roi"`
	Recommendations           []string `json:"recommendations"`
}

// AnalysisRequest represents the request for historical analysis
type AnalysisRequest struct {
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	CustomerID             string    `json:"customer_id,omitempty"`
	AnalysisTypes          []string  `json:"analysis_types"` // bundling, volume, loyalty, comprehensive
	IncludeRecommendations bool      `json:"include_recommendations"`
}

// AnalysisResponse represents the response from historical analysis
type AnalysisResponse struct {
	Status                string                 `json:"status"`
	Message               string                 `json:"message"`
	AnalysisRequest       AnalysisRequest        `json:"analysis_request"`
	ComprehensiveAnalysis *ComprehensiveAnalysis `json:"comprehensive_analysis,omitempty"`
	Error                 string                 `json:"error,omitempty"`
}
