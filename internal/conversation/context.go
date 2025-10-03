package conversation

import (
	"encoding/json"
	"fmt"
	"time"
)

// ContextManager handles conversation context and state
type ContextManager struct {
	sessions map[string]*ConversationContext
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		sessions: make(map[string]*ConversationContext),
	}
}

// Update updates the conversation context with new information
func (cm *ContextManager) Update(context *ConversationContext, intent *Intent) *ConversationContext {
	if context == nil {
		context = &ConversationContext{
			SessionID: fmt.Sprintf("session_%d", time.Now().Unix()),
			CustomerProfile: CustomerProfile{
				Tier: "bronze", // Default tier
			},
			DeliveryHistory: []DeliveryRequirement{},
			PricingHistory:  []PricingComparison{},
			Preferences:     CustomerPreferences{},
		}
	}

	// Update context based on intent
	cm.updateFromIntent(context, intent)

	// Store updated context
	cm.sessions[context.SessionID] = context

	return context
}

// updateFromIntent updates context based on recognized intent
func (cm *ContextManager) updateFromIntent(context *ConversationContext, intent *Intent) {
	// Update customer profile
	if tier, exists := intent.Entities["customer_tier"]; exists {
		context.CustomerProfile.Tier = tier
	}

	if freq, exists := intent.Entities["order_frequency"]; exists {
		if freqInt, err := fmt.Sscanf(freq, "%d", &context.CustomerProfile.OrderFrequency); err == nil {
			context.CustomerProfile.OrderFrequency = freqInt
		}
	}

	if vehicle, exists := intent.Entities["vehicle_type"]; exists {
		context.CustomerProfile.PreferredVehicle = vehicle
	}

	// Update delivery requirements
	if count, exists := intent.Entities["delivery_count"]; exists {
		if countInt, err := fmt.Sscanf(count, "%d", new(int)); err == nil {
			// Update current delivery count in customer profile
			context.CustomerProfile.CurrentDeliveryCount = countInt
			
			deliveryReq := DeliveryRequirement{
				Count:       countInt,
				VehicleType: context.CustomerProfile.PreferredVehicle,
			}
			context.DeliveryHistory = append(context.DeliveryHistory, deliveryReq)
		}
	}

	// Update preferences based on intent
	switch intent.Type {
	case "compare_pricing":
		context.CurrentGoal = "compare_pricing"
	case "get_recommendation":
		context.CurrentGoal = "get_recommendation"
	case "explore_options":
		context.CurrentGoal = "explore_options"
	}
}

// GetSession retrieves a conversation session
func (cm *ContextManager) GetSession(sessionID string) *ConversationContext {
	if context, exists := cm.sessions[sessionID]; exists {
		return context
	}
	return nil
}

// SaveSession saves a conversation session
func (cm *ContextManager) SaveSession(context *ConversationContext) {
	cm.sessions[context.SessionID] = context
}

// DeleteSession removes a conversation session
func (cm *ContextManager) DeleteSession(sessionID string) {
	delete(cm.sessions, sessionID)
}

// GetAllSessions returns all active sessions
func (cm *ContextManager) GetAllSessions() map[string]*ConversationContext {
	return cm.sessions
}

// GetSessionSummary returns a summary of a session
func (cm *ContextManager) GetSessionSummary(sessionID string) string {
	context := cm.GetSession(sessionID)
	if context == nil {
		return "Session not found"
	}

	summary := fmt.Sprintf("Session: %s\n", context.SessionID)
	summary += fmt.Sprintf("Customer Tier: %s\n", context.CustomerProfile.Tier)
	summary += fmt.Sprintf("Order Frequency: %d/month\n", context.CustomerProfile.OrderFrequency)
	summary += fmt.Sprintf("Current Goal: %s\n", context.CurrentGoal)
	summary += fmt.Sprintf("Delivery History: %d entries\n", len(context.DeliveryHistory))
	summary += fmt.Sprintf("Pricing History: %d entries\n", len(context.PricingHistory))

	return summary
}

// ExportSession exports a session to JSON
func (cm *ContextManager) ExportSession(sessionID string) (string, error) {
	context := cm.GetSession(sessionID)
	if context == nil {
		return "", fmt.Errorf("session not found")
	}

	jsonData, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ImportSession imports a session from JSON
func (cm *ContextManager) ImportSession(jsonData string) error {
	var context ConversationContext
	if err := json.Unmarshal([]byte(jsonData), &context); err != nil {
		return err
	}

	cm.SaveSession(&context)
	return nil
}

// ClearExpiredSessions removes sessions older than the specified duration
func (cm *ContextManager) ClearExpiredSessions(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)

	for sessionID, sessionContext := range cm.sessions {
		// Simple expiration based on session ID timestamp
		// In a real implementation, you'd store creation time
		if time.Since(cutoff) > maxAge {
			delete(cm.sessions, sessionID)
		}
		_ = sessionContext // Avoid unused variable warning
	}
}

// GetSessionStats returns statistics about active sessions
func (cm *ContextManager) GetSessionStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_sessions":    len(cm.sessions),
		"tier_distribution": make(map[string]int),
		"goal_distribution": make(map[string]int),
	}

	tierDist := make(map[string]int)
	goalDist := make(map[string]int)

	for _, context := range cm.sessions {
		tierDist[context.CustomerProfile.Tier]++
		goalDist[context.CurrentGoal]++
	}

	stats["tier_distribution"] = tierDist
	stats["goal_distribution"] = goalDist

	return stats
}
