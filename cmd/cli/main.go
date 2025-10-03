package main

import (
	"bufio"
	"dispatch-mcp-server/internal/config"
	"dispatch-mcp-server/internal/conversation"
	"dispatch-mcp-server/internal/dispatch"
	"dispatch-mcp-server/internal/pricing"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "estimate":
		handleEstimate()
	case "order":
		handleOrder()
	case "pricing":
		handlePricingComparison()
	case "chat":
		handleConversationalPricing()
	case "interactive":
		handleInteractive()
	case "login":
		handleLogin()
	case "logout":
		handleLogout()
	case "subenv":
		handleSubenv()
	case "status":
		showStatus()
	case "help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("🚀 Dispatch MCP CLI Demo Tool")
	fmt.Println("=============================")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  ./dispatch-cli estimate     - Create a cost estimate")
	fmt.Println("  ./dispatch-cli order        - Create a delivery order")
	fmt.Println("  ./dispatch-cli pricing      - Compare different pricing models")
	fmt.Println("  ./dispatch-cli chat         - Conversational pricing advisor")
	fmt.Println("  ./dispatch-cli interactive  - Interactive mode")
	fmt.Println("  ./dispatch-cli login        - Authenticate with Dispatch API")
	fmt.Println("  ./dispatch-cli logout       - Clear authentication")
	fmt.Println("  ./dispatch-cli subenv       - Set subenv (monkey, staging, prod)")
	fmt.Println("  ./dispatch-cli status       - Show connection status")
	fmt.Println("  ./dispatch-cli help         - Show this help")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  USE_IDP_AUTH=true/false     - Use IDP authentication")
	fmt.Println("  DISPATCH_ORGANIZATION_ID    - Your organization ID")
	fmt.Println("  (See env.example for full list)")
}

func handleEstimate() {
	fmt.Println("📊 Creating Cost Estimate...")
	fmt.Println("============================")

	// Create sample pickup location
	pickupInfo := dispatch.PickupInfoInput{
		BusinessName: "Demo Business",
		Location: dispatch.LocationInput{
			Address: &dispatch.AddressInput{
				Street:  "123 Market St",
				City:    "San Francisco",
				State:   "CA",
				ZipCode: "94105",
				Country: "US",
			},
		},
	}

	// Create sample drop-off locations
	dropOffs := []dispatch.DropOffInfoInput{
		{
			BusinessName: "Customer Location",
			Location: dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  "456 Oak Ave",
					City:    "Oakland",
					State:   "CA",
					ZipCode: "94610",
					Country: "US",
				},
			},
		},
	}

	// Create estimate input
	input := dispatch.CreateEstimateInput{
		PickupInfo:  pickupInfo,
		DropOffs:    dropOffs,
		VehicleType: "cargo_van",
	}

	// Create client and make API call
	client, err := dispatch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("🔄 Calling Dispatch API...")
	response, err := client.CreateEstimate(input)
	if err != nil {
		log.Fatalf("Failed to create estimate: %v", err)
	}

	// Display results
	fmt.Println("✅ Estimate created successfully!")
	fmt.Println("")

	if len(response.Data.CreateEstimate.Estimate.AvailableOrderOptions) > 0 {
		option := response.Data.CreateEstimate.Estimate.AvailableOrderOptions[0]
		fmt.Printf("💰 Estimated Cost: $%.2f\n", option.EstimatedOrderCost)
		fmt.Printf("🚚 Vehicle Type: %s\n", option.VehicleType)
		fmt.Printf("⏰ Estimated Delivery: %s\n", option.EstimatedDeliveryTimeUTC)
		fmt.Printf("🏢 Service Type: %s\n", option.ServiceType)
	} else {
		fmt.Println("⚠️  No delivery options available")
	}

	// Show full response in JSON
	fmt.Println("\n📋 Full Response:")
	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonData))
}

func handleOrder() {
	fmt.Println("📦 Creating Delivery Order...")
	fmt.Println("==============================")

	// Create sample delivery info
	deliveryInfo := dispatch.DeliveryInfoInput{
		ServiceType: "delivery",
	}

	// Create sample pickup info
	pickupInfo := dispatch.CreateOrderPickupInfoInput{
		BusinessName:       stringPtr("Demo Business"),
		ContactName:        stringPtr("John Doe"),
		ContactPhoneNumber: stringPtr("555-123-4567"),
		Location: &dispatch.LocationInput{
			Address: &dispatch.AddressInput{
				Street:  "123 Market St",
				City:    "San Francisco",
				State:   "CA",
				ZipCode: "94105",
				Country: "US",
			},
		},
	}

	// Create sample drop-off info
	dropOffs := []dispatch.CreateOrderDropOffInfoInput{
		{
			BusinessName:       stringPtr("Customer Location"),
			ContactName:        stringPtr("Jane Smith"),
			ContactPhoneNumber: stringPtr("555-987-6543"),
			Location: &dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  "456 Oak Ave",
					City:    "Oakland",
					State:   "CA",
					ZipCode: "94610",
					Country: "US",
				},
			},
		},
	}

	// Create order input
	input := dispatch.CreateOrderInput{
		DeliveryInfo: deliveryInfo,
		PickupInfo:   pickupInfo,
		DropOffs:     dropOffs,
	}

	// Create client and make API call
	client, err := dispatch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("🔄 Calling Dispatch API...")
	response, err := client.CreateOrder(input)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	// Display results
	fmt.Println("✅ Order created successfully!")
	fmt.Println("")

	order := response.Data.CreateOrder.Order
	fmt.Printf("🆔 Order ID: %s\n", order.ID)
	fmt.Printf("📊 Status: %s\n", order.Status)
	fmt.Printf("💰 Total Cost: $%.2f\n", order.TotalCost)
	fmt.Printf("📦 Tracking Number: %s\n", order.TrackingNumber)
	fmt.Printf("⏰ Scheduled At: %s\n", order.ScheduledAt)

	// Show full response in JSON
	fmt.Println("\n📋 Full Response:")
	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonData))
}

func handleInteractive() {
	fmt.Println("🎮 Interactive Dispatch CLI")
	fmt.Println("============================")
	fmt.Println("Type 'help' for commands, 'quit' to exit")
	fmt.Println("")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("dispatch> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "help":
			showInteractiveHelp()
		case "estimate":
			handleEstimate()
		case "order":
			handleOrder()
		case "login":
			handleLogin()
		case "logout":
			handleLogout()
		case "subenv":
			handleSubenv()
		case "status":
			showStatus()
		case "quit", "exit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}

		fmt.Println("")
	}
}

func showInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  estimate  - Create a cost estimate")
	fmt.Println("  order     - Create a delivery order")
	fmt.Println("  login     - Authenticate with Dispatch API")
	fmt.Println("  logout    - Clear authentication")
	fmt.Println("  subenv    - Set subenv (monkey, staging, prod)")
	fmt.Println("  status    - Show connection status")
	fmt.Println("  help      - Show this help")
	fmt.Println("  quit      - Exit the program")
}

func showStatus() {
	fmt.Println("🔍 Connection Status")
	fmt.Println("===================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config Error: %v\n", err)
		return
	}

	fmt.Printf("🔐 Authentication: %s\n", map[bool]string{true: "IDP", false: "Static Token"}[cfg.UseIDP])
	fmt.Printf("🏢 Organization ID: %s\n", cfg.OrganizationID)
	fmt.Printf("📡 GraphQL Endpoint: %s\n", cfg.GraphQLEndpoint)

	// Test client creation
	_, err = dispatch.NewClient()
	if err != nil {
		fmt.Printf("❌ Client Error: %v\n", err)
	} else {
		fmt.Println("✅ Client created successfully")
	}
}

func handleLogin() {
	fmt.Println("🔐 Dispatch API Login")
	fmt.Println("===================")
	fmt.Println("")

	// Check if already authenticated
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config Error: %v\n", err)
		return
	}

	if cfg.AuthToken != "" || cfg.UseIDP {
		fmt.Println("✅ Already authenticated!")
		fmt.Printf("🔐 Method: %s\n", map[bool]string{true: "IDP", false: "Static Token"}[cfg.UseIDP])
		fmt.Printf("🏢 Organization: %s\n", cfg.OrganizationID)
		return
	}

	fmt.Println("Choose authentication method:")
	fmt.Println("1. Static Token (API Key)")
	fmt.Println("2. IDP Authentication")
	fmt.Println("3. Cancel")
	fmt.Print("Enter choice (1-3): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		handleStaticTokenLogin()
	case "2":
		handleIDPLogin()
	case "3":
		fmt.Println("❌ Login cancelled")
	default:
		fmt.Println("❌ Invalid choice")
	}
}

func handleStaticTokenLogin() {
	fmt.Println("")
	fmt.Println("🔑 Static Token Authentication")
	fmt.Println("==============================")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your Dispatch API token: ")
	scanner.Scan()
	token := strings.TrimSpace(scanner.Text())

	if token == "" {
		fmt.Println("❌ Token cannot be empty")
		return
	}

	fmt.Print("Enter your Organization ID: ")
	scanner.Scan()
	orgID := strings.TrimSpace(scanner.Text())

	if orgID == "" {
		fmt.Println("❌ Organization ID cannot be empty")
		return
	}

	// Set environment variables for current session
	os.Setenv("DISPATCH_AUTH_TOKEN", token)
	os.Setenv("DISPATCH_ORGANIZATION_ID", orgID)
	os.Setenv("USE_IDP_AUTH", "false")

	fmt.Println("")
	fmt.Println("✅ Authentication configured!")
	fmt.Println("🔐 Method: Static Token")
	fmt.Printf("🏢 Organization: %s\n", orgID)
	fmt.Println("")
	fmt.Println("💡 Note: These settings are for this session only.")
	fmt.Println("   To persist, set environment variables in your shell.")
}

func handleIDPLogin() {
	fmt.Println("")
	fmt.Println("🔐 IDP Authentication")
	fmt.Println("===================")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter IDP Endpoint (default: https://id.dispatchfog.io): ")
	scanner.Scan()
	idpEndpoint := strings.TrimSpace(scanner.Text())
	if idpEndpoint == "" {
		idpEndpoint = "https://id.dispatchfog.io"
	}

	fmt.Print("Enter Client ID: ")
	scanner.Scan()
	clientID := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Client Secret: ")
	scanner.Scan()
	clientSecret := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Organization ID: ")
	scanner.Scan()
	orgID := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Scope (default: dispatch:api): ")
	scanner.Scan()
	scope := strings.TrimSpace(scanner.Text())
	if scope == "" {
		scope = "dispatch:api"
	}

	// Set environment variables for current session
	os.Setenv("USE_IDP_AUTH", "true")
	os.Setenv("IDP_ENDPOINT", idpEndpoint)
	os.Setenv("IDP_CLIENT_ID", clientID)
	os.Setenv("IDP_CLIENT_SECRET", clientSecret)
	os.Setenv("IDP_SCOPE", scope)
	os.Setenv("IDP_TOKEN_ENDPOINT", idpEndpoint+"/oauth/token")
	os.Setenv("DISPATCH_ORGANIZATION_ID", orgID)

	fmt.Println("")
	fmt.Println("✅ IDP Authentication configured!")
	fmt.Println("🔐 Method: IDP")
	fmt.Printf("🏢 Organization: %s\n", orgID)
	fmt.Printf("🔗 IDP Endpoint: %s\n", idpEndpoint)
	fmt.Println("")
	fmt.Println("💡 Note: These settings are for this session only.")
	fmt.Println("   To persist, set environment variables in your shell.")
}

func handleLogout() {
	fmt.Println("🚪 Logging out...")
	fmt.Println("=================")

	// Clear environment variables
	os.Unsetenv("DISPATCH_AUTH_TOKEN")
	os.Unsetenv("DISPATCH_ORGANIZATION_ID")
	os.Unsetenv("USE_IDP_AUTH")
	os.Unsetenv("IDP_ENDPOINT")
	os.Unsetenv("IDP_CLIENT_ID")
	os.Unsetenv("IDP_CLIENT_SECRET")
	os.Unsetenv("IDP_SCOPE")
	os.Unsetenv("IDP_TOKEN_ENDPOINT")

	fmt.Println("✅ Logged out successfully!")
	fmt.Println("🔄 Switched to mock mode")
	fmt.Println("")
	fmt.Println("💡 To login again, use: ./dispatch-cli login")
}

func handleSubenv() {
	fmt.Println("🌍 Dispatch Subenv Configuration")
	fmt.Println("===============================")
	fmt.Println("")

	// Show current subenv
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config Error: %v\n", err)
		return
	}

	fmt.Printf("Current endpoint: %s\n", cfg.GraphQLEndpoint)
	fmt.Println("")

	// Define available subenvs
	subenvs := map[string]string{
		"1": "monkey",
		"2": "staging",
		"3": "production",
		"4": "development",
		"5": "custom",
	}

	endpoints := map[string]string{
		"monkey":      "https://monkey.graph.qa.dispatchfog.io/graphql",
		"staging":     "https://qa.graph.dispatchfog.io/graphql",
		"production":  "https://graph.dispatchfog.io/graphql",
		"development": "https://eng.graph.dispatchfog.io/graphql",
	}

	fmt.Println("Available subenvs:")
	fmt.Println("1. monkey      - Monkey subenv (testing)")
	fmt.Println("2. staging     - Staging environment")
	fmt.Println("3. production  - Production environment")
	fmt.Println("4. development - Development environment")
	fmt.Println("5. custom      - Enter custom endpoint")
	fmt.Println("6. cancel      - Cancel")
	fmt.Print("Enter choice (1-6): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	if choice == "6" {
		fmt.Println("❌ Subenv change cancelled")
		return
	}

	subenv, exists := subenvs[choice]
	if !exists {
		fmt.Println("❌ Invalid choice")
		return
	}

	var endpoint string

	if subenv == "custom" {
		fmt.Print("Enter custom GraphQL endpoint: ")
		scanner.Scan()
		endpoint = strings.TrimSpace(scanner.Text())

		if endpoint == "" {
			fmt.Println("❌ Endpoint cannot be empty")
			return
		}

		// Validate URL format
		if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
			endpoint = "https://" + endpoint
		}
	} else {
		endpoint = endpoints[subenv]
	}

	// Set the endpoint
	os.Setenv("DISPATCH_GRAPHQL_ENDPOINT", endpoint)

	fmt.Println("")
	fmt.Println("✅ Subenv configured successfully!")
	fmt.Printf("🌍 Environment: %s\n", subenv)
	fmt.Printf("📡 Endpoint: %s\n", endpoint)
	fmt.Println("")
	fmt.Println("💡 Note: This setting is for this session only.")
	fmt.Println("   To persist, set DISPATCH_GRAPHQL_ENDPOINT in your shell.")
}

func handlePricingComparison() {
	fmt.Println("💰 Pricing Model Comparison")
	fmt.Println("===========================")
	fmt.Println("")

	// First, create an estimate to compare pricing models against
	fmt.Println("🔄 Creating base estimate...")

	// Create sample pickup location
	pickupInfo := dispatch.PickupInfoInput{
		BusinessName: "Demo Business",
		Location: dispatch.LocationInput{
			Address: &dispatch.AddressInput{
				Street:  "123 Market St",
				City:    "San Francisco",
				State:   "CA",
				ZipCode: "94105",
				Country: "US",
			},
		},
	}

	// Create sample drop-off locations
	dropOffs := []dispatch.DropOffInfoInput{
		{
			BusinessName: "Customer Location 1",
			Location: dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  "456 Oak Ave",
					City:    "Oakland",
					State:   "CA",
					ZipCode: "94610",
					Country: "US",
				},
			},
		},
		{
			BusinessName: "Customer Location 2",
			Location: dispatch.LocationInput{
				Address: &dispatch.AddressInput{
					Street:  "789 Pine St",
					City:    "Berkeley",
					State:   "CA",
					ZipCode: "94710",
					Country: "US",
				},
			},
		},
	}

	// Create estimate input
	input := dispatch.CreateEstimateInput{
		PickupInfo:  pickupInfo,
		DropOffs:    dropOffs,
		VehicleType: "cargo_van",
	}

	// Create client and make API call
	client, err := dispatch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.CreateEstimate(input)
	if err != nil {
		log.Fatalf("Failed to create estimate: %v", err)
	}

	if len(response.Data.CreateEstimate.Estimate.AvailableOrderOptions) == 0 {
		fmt.Println("⚠️  No delivery options available for comparison")
		return
	}

	originalEstimate := response.Data.CreateEstimate.Estimate.AvailableOrderOptions[0]

	fmt.Printf("✅ Base estimate created: $%.2f\n", originalEstimate.EstimatedOrderCost)
	fmt.Println("")

	// Now compare different pricing models
	fmt.Println("🔍 Comparing Pricing Models...")
	fmt.Println("===============================")

	// Create pricing engine
	engine := pricing.NewPricingEngine()

	// Test different scenarios
	scenarios := []struct {
		name    string
		context pricing.PricingContext
	}{
		{
			name: "Standard Customer (1 delivery, bronze tier)",
			context: pricing.PricingContext{
				DeliveryCount:   1,
				CustomerTier:    "bronze",
				OrderFrequency:  1,
				TotalOrderValue: originalEstimate.EstimatedOrderCost,
				IsBulkOrder:     false,
			},
		},
		{
			name: "Multi-Delivery Customer (2 deliveries, silver tier)",
			context: pricing.PricingContext{
				DeliveryCount:   2,
				CustomerTier:    "silver",
				OrderFrequency:  3,
				TotalOrderValue: originalEstimate.EstimatedOrderCost * 2,
				IsBulkOrder:     false,
			},
		},
		{
			name: "High-Volume Customer (5 deliveries, gold tier)",
			context: pricing.PricingContext{
				DeliveryCount:   5,
				CustomerTier:    "gold",
				OrderFrequency:  8,
				TotalOrderValue: originalEstimate.EstimatedOrderCost * 5,
				IsBulkOrder:     false,
			},
		},
		{
			name: "Bulk Order Customer (10 deliveries, gold tier)",
			context: pricing.PricingContext{
				DeliveryCount:   10,
				CustomerTier:    "gold",
				OrderFrequency:  15,
				TotalOrderValue: originalEstimate.EstimatedOrderCost * 10,
				IsBulkOrder:     true,
			},
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("\n📊 Scenario %d: %s\n", i+1, scenario.name)
		fmt.Println(strings.Repeat("-", 50))

		comparison := engine.ComparePricingModels(&originalEstimate, scenario.context)

		// Display results
		fmt.Printf("Original Cost: $%.2f\n", originalEstimate.EstimatedOrderCost)
		fmt.Println("")

		for _, result := range comparison.PricingModels {
			status := "❌ Not Eligible"
			if result.Eligible {
				status = "✅ Eligible"
			}

			fmt.Printf("🏷️  %s: %s\n", result.Name, status)
			if result.Eligible {
				fmt.Printf("   💰 Adjusted Cost: $%.2f\n", result.AdjustedCost)
				fmt.Printf("   💸 Savings: $%.2f (%.1f%%)\n", result.Savings, result.DiscountPercent)
			} else {
				fmt.Printf("   📝 Reason: %s\n", result.Reason)
			}
			fmt.Println("")
		}

		if comparison.BestOption != nil {
			fmt.Printf("🏆 Best Option: %s\n", comparison.BestOption.Name)
			fmt.Printf("💰 Best Price: $%.2f\n", comparison.BestOption.AdjustedCost)
			fmt.Printf("💸 Total Savings: $%.2f (%.1f%%)\n", comparison.Savings, comparison.SavingsPercentage)
		}
	}

	fmt.Println("\n🎯 Summary:")
	fmt.Println("===========")
	fmt.Println("• Standard Pricing: No discounts")
	fmt.Println("• Multi-Delivery: 15% discount for 2+ deliveries")
	fmt.Println("• Volume Discount: 20% discount for 5+ deliveries + 3+ orders/month")
	fmt.Println("• Loyalty Discount: 10% discount for gold tier customers")
	fmt.Println("• Bulk Order: 25% discount for 10+ deliveries in bulk orders")
	fmt.Println("")
	fmt.Println("💡 Tip: Combine multiple discounts for maximum savings!")
}

func handleConversationalPricing() {
	fmt.Println("🗣️  Conversational Pricing Advisor")
	fmt.Println("==================================")
	fmt.Println("")
	fmt.Println("Chat with our AI pricing advisor to find the best pricing for your needs!")
	fmt.Println("Type 'quit' to exit, 'help' for examples.")
	fmt.Println("")

	// Create AI Hub-powered conversation engine with rule-based fallback
	engine, err := conversation.NewClaudeConversationEngine()
	if err != nil {
		fmt.Printf("⚠️  Conversation engine error: %v\n", err)
	}
	var context *conversation.ConversationContext

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("💬 You: ")
	for scanner.Scan() {
		userInput := strings.TrimSpace(scanner.Text())

		if userInput == "quit" || userInput == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		if userInput == "help" {
			showConversationalHelp()
			fmt.Print("💬 You: ")
			continue
		}

		if userInput == "" {
			fmt.Print("💬 You: ")
			continue
		}

		// Show thinking indicator while processing
		thinkingDone := make(chan bool)
		go showThinkingIndicator(thinkingDone)

		// Process the message
		response, err := engine.ProcessMessage(userInput, context)

		// Stop thinking indicator
		thinkingDone <- true

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Print("💬 You: ")
			continue
		}

		// Update context
		context = response.UpdatedContext

		// Display response with markdown formatting removed
		cleanMessage := removeMarkdownFormatting(response.Message)
		fmt.Printf("🤖 Advisor: %s\n", cleanMessage)

		// Show recommendations if available
		if len(response.Recommendations) > 0 {
			fmt.Println("\n📊 Pricing Recommendations:")
			for _, rec := range response.Recommendations {
				if rec.Eligible {
					fmt.Printf("  ✅ %s: $%.2f savings (%.1f%%)\n", rec.Name, rec.Savings, rec.SavingsPercent)
				} else {
					fmt.Printf("  ❌ %s: %s\n", rec.Name, rec.Reason)
				}
			}
		}

		// Show next questions if available
		if len(response.NextQuestions) > 0 {
			fmt.Println("\n💡 Next Steps:")
			for i, question := range response.NextQuestions {
				fmt.Printf("  %d. %s\n", i+1, question)
			}
		}

		fmt.Println("")
		fmt.Print("💬 You: ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
	}
}

func showConversationalHelp() {
	fmt.Println("\n📚 Conversational Pricing Advisor Help")
	fmt.Println("=====================================")
	fmt.Println("")
	fmt.Println("💬 Example Conversations:")
	fmt.Println("")
	fmt.Println("  'I need 3 deliveries to different locations'")
	fmt.Println("  'What's the best pricing for a gold customer?'")
	fmt.Println("  'How can I save money on my deliveries?'")
	fmt.Println("  'Show me bulk order discounts'")
	fmt.Println("  'I'm a bronze tier customer with 2 orders per month'")
	fmt.Println("")
	fmt.Println("🎯 Available Pricing Models:")
	fmt.Println("  • Standard Pricing: No discounts")
	fmt.Println("  • Multi-Delivery: 15% off for 2+ deliveries")
	fmt.Println("  • Volume Discount: 20% off for 5+ deliveries + 3+ orders/month")
	fmt.Println("  • Loyalty Discount: 10% off for gold tier customers")
	fmt.Println("  • Bulk Order: 25% off for 10+ deliveries + bulk flag")
	fmt.Println("")
	fmt.Println("💡 Tips:")
	fmt.Println("  • Be specific about your delivery count")
	fmt.Println("  • Mention your customer tier (bronze, silver, gold)")
	fmt.Println("  • Ask about bulk ordering for maximum savings")
	fmt.Println("  • Inquire about loyalty program benefits")
	fmt.Println("")
}

func stringPtr(s string) *string {
	return &s
}

// removeMarkdownFormatting removes markdown formatting from text for CLI display
func removeMarkdownFormatting(text string) string {
	// Remove bold formatting (**text** -> text)
	text = strings.ReplaceAll(text, "**", "")

	// Remove italic formatting (*text* -> text)
	text = strings.ReplaceAll(text, "*", "")

	// Remove code formatting (`text` -> text)
	text = strings.ReplaceAll(text, "`", "")

	// Remove headers (# ## ### -> empty)
	lines := strings.Split(text, "\n")
	var cleanLines []string
	for _, line := range lines {
		// Remove header markers but keep the text
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimPrefix(line, "##")
		line = strings.TrimPrefix(line, "###")
		line = strings.TrimPrefix(line, "####")
		line = strings.TrimSpace(line)
		cleanLines = append(cleanLines, line)
	}

	return strings.Join(cleanLines, "\n")
}

// showThinkingIndicator displays a thinking animation while waiting for Claude's response
func showThinkingIndicator(done chan bool) {
	thinkingChars := []string{"🤔", "💭", "🧠", "⚡"}
	i := 0

	for {
		select {
		case <-done:
			// Clear the thinking indicator
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r🤖 Thinking %s", thinkingChars[i%len(thinkingChars)])
			time.Sleep(300 * time.Millisecond)
			i++
		}
	}
}
