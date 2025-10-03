package main

import (
	"dispatch-mcp-server/internal/conversation"
	"dispatch-mcp-server/internal/dispatch"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// ChatMessage represents a message in the chat
type ChatMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatSession represents a chat session
type ChatSession struct {
	ID          string                            `json:"id"`
	Messages    []ChatMessage                     `json:"messages"`
	Context     *conversation.ConversationContext `json:"context"`
	OrderInfo   *OrderInfo                        `json:"order_info,omitempty"`
	PricingInfo *PricingInfo                      `json:"pricing_info,omitempty"`
}

// OrderInfo represents order creation progress
type OrderInfo struct {
	InProgress      bool           `json:"in_progress"`
	Step            string         `json:"step"`
	CurrentQuestion string         `json:"current_question"`
	PickupInfo      *PickupInfo    `json:"pickup_info,omitempty"`
	Deliveries      []DeliveryInfo `json:"deliveries,omitempty"`
	CompletedFields []string       `json:"completed_fields"`
	MissingFields   []string       `json:"missing_fields"`
}

// PickupInfo represents pickup location information
type PickupInfo struct {
	BusinessName string `json:"business_name"`
	Address      string `json:"address"`
	ContactName  string `json:"contact_name"`
	PhoneNumber  string `json:"phone_number"`
}

// DeliveryInfo represents delivery location information
type DeliveryInfo struct {
	BusinessName string `json:"business_name"`
	Address      string `json:"address"`
	ContactName  string `json:"contact_name"`
	PhoneNumber  string `json:"phone_number"`
}

// PricingInfo represents pricing recommendations
type PricingInfo struct {
	Recommendations []PricingRecommendation `json:"recommendations"`
	BestOption      *PricingRecommendation  `json:"best_option,omitempty"`
	TotalSavings    float64                 `json:"total_savings"`
}

// PricingRecommendation represents a pricing option
type PricingRecommendation struct {
	Name           string  `json:"name"`
	Savings        float64 `json:"savings"`
	SavingsPercent float64 `json:"savings_percent"`
	Eligible       bool    `json:"eligible"`
	Reason         string  `json:"reason,omitempty"`
}

var (
	engine   *conversation.ClaudeConversationEngine
	sessions = make(map[string]*ChatSession)
)

func main() {
	// Initialize conversation engine
	var err error
	engine, err = conversation.NewClaudeConversationEngine()
	if err != nil {
		log.Printf("⚠️  Claude not available, using rule-based engine: %v", err)
	}

	// Set up routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/chat", handleChat)
	http.HandleFunc("/api/session", handleSession)
	http.HandleFunc("/api/health", handleHealth)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Dispatch Web Chat Server starting on port %s\n", port)
	fmt.Printf("🌐 Open http://localhost:%s in your browser\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create new session
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())
	session := &ChatSession{
		ID:       sessionID,
		Messages: []ChatMessage{},
		Context:  nil,
	}
	sessions[sessionID] = session

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		SessionID string `json:"session_id"`
		Message   string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get or create session
	session, exists := sessions[request.SessionID]
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Add user message
	userMessage := ChatMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Type:      "user",
		Content:   request.Message,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, userMessage)

	// Process message with conversation engine
	response, err := engine.ProcessMessage(request.Message, session.Context)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing message: %v", err), http.StatusInternalServerError)
		return
	}

	// Update session context
	session.Context = response.UpdatedContext

	// Add assistant response
	assistantMessage := ChatMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Type:      "assistant",
		Content:   response.Message,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, assistantMessage)

	// Update order and pricing info
	updateSessionInfo(session, response)

	// Return updated session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func updateSessionInfo(session *ChatSession, response *conversation.ConversationResponse) {
	// Update order info
	if session.Context != nil && session.Context.OrderCreation.InProgress {
		session.OrderInfo = &OrderInfo{
			InProgress:      session.Context.OrderCreation.InProgress,
			Step:            session.Context.OrderCreation.Step,
			CurrentQuestion: session.Context.OrderCreation.CurrentQuestion,
			CompletedFields: session.Context.OrderCreation.CompletedFields,
			MissingFields:   session.Context.OrderCreation.MissingFields,
		}

		// Add pickup info if available
		if session.Context.OrderCreation.PickupInfo != nil {
			pickup := session.Context.OrderCreation.PickupInfo
			session.OrderInfo.PickupInfo = &PickupInfo{
				BusinessName: getStringValue(pickup.BusinessName),
				Address:      getAddressString(pickup.Location),
				ContactName:  getStringValue(pickup.ContactName),
				PhoneNumber:  getStringValue(pickup.ContactPhoneNumber),
			}
		}

		// Add delivery info if available
		if len(session.Context.OrderCreation.DropOffs) > 0 {
			for _, dropoff := range session.Context.OrderCreation.DropOffs {
				delivery := DeliveryInfo{
					BusinessName: getStringValue(dropoff.BusinessName),
					Address:      getAddressString(dropoff.Location),
					ContactName:  getStringValue(dropoff.ContactName),
					PhoneNumber:  getStringValue(dropoff.ContactPhoneNumber),
				}
				session.OrderInfo.Deliveries = append(session.OrderInfo.Deliveries, delivery)
			}
		}
	}

	// Update pricing info
	if len(response.Recommendations) > 0 {
		session.PricingInfo = &PricingInfo{
			Recommendations: convertPricingRecommendations(response.Recommendations),
		}

		// Find best option
		bestSavings := 0.0
		for _, rec := range response.Recommendations {
			if rec.Eligible && rec.Savings > bestSavings {
				bestSavings = rec.Savings
				session.PricingInfo.BestOption = &PricingRecommendation{
					Name:           rec.Name,
					Savings:        rec.Savings,
					SavingsPercent: rec.SavingsPercent,
					Eligible:       rec.Eligible,
					Reason:         rec.Reason,
				}
			}
		}
		session.PricingInfo.TotalSavings = bestSavings
	}
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getAddressString(location *dispatch.LocationInput) string {
	if location == nil || location.Address == nil {
		return ""
	}
	addr := location.Address
	return fmt.Sprintf("%s, %s, %s %s", addr.Street, addr.City, addr.State, addr.ZipCode)
}

func convertPricingRecommendations(recs []conversation.PricingRecommendation) []PricingRecommendation {
	var result []PricingRecommendation
	for _, rec := range recs {
		result = append(result, PricingRecommendation{
			Name:           rec.Name,
			Savings:        rec.Savings,
			SavingsPercent: rec.SavingsPercent,
			Eligible:       rec.Eligible,
			Reason:         rec.Reason,
		})
	}
	return result
}
