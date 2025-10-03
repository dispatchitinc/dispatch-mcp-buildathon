package validation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult contains the result of validation
type ValidationResult struct {
	Valid   bool              `json:"valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Message string            `json:"message,omitempty"`
}

// Validator provides validation functions for MCP tool inputs
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateVehicleType validates the vehicle type parameter
func (v *Validator) ValidateVehicleType(vehicleType string) *ValidationResult {
	validTypes := []string{"pickup_truck", "cargo_van", "sprinter_van", "box_truck"}

	if vehicleType == "" {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "vehicle_type",
					Message: "vehicle_type is required",
				},
			},
			Message: "Vehicle type is required",
		}
	}

	for _, validType := range validTypes {
		if strings.EqualFold(vehicleType, validType) {
			return &ValidationResult{Valid: true}
		}
	}

	return &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{
				Field:   "vehicle_type",
				Message: fmt.Sprintf("Invalid vehicle type '%s'. Must be one of: %s", vehicleType, strings.Join(validTypes, ", ")),
			},
		},
		Message: "Invalid vehicle type",
	}
}

// ValidateDeliveryScenario validates the delivery scenario parameter
func (v *Validator) ValidateDeliveryScenario(scenario string) *ValidationResult {
	validScenarios := []string{"fastest", "asap", "urgent", "cheapest", "economy", "sometime_today"}

	if scenario == "" {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "delivery_scenario",
					Message: "delivery_scenario is required",
				},
			},
			Message: "Delivery scenario is required",
		}
	}

	for _, validScenario := range validScenarios {
		if strings.EqualFold(scenario, validScenario) {
			return &ValidationResult{Valid: true}
		}
	}

	return &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{
				Field:   "delivery_scenario",
				Message: fmt.Sprintf("Invalid delivery scenario '%s'. Must be one of: %s", scenario, strings.Join(validScenarios, ", ")),
			},
		},
		Message: "Invalid delivery scenario",
	}
}

// ValidateAddress validates address information
func (v *Validator) ValidateAddress(address map[string]interface{}) *ValidationResult {
	var errors []ValidationError

	// Check required fields
	requiredFields := []string{"street", "city", "state", "zip_code", "country"}
	for _, field := range requiredFields {
		if value, exists := address[field]; !exists || value == "" {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s is required", field),
			})
		}
	}

	// Validate zip code format (basic US zip code validation)
	if zipCode, exists := address["zip_code"]; exists && zipCode != "" {
		zipStr := fmt.Sprintf("%v", zipCode)
		zipRegex := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
		if !zipRegex.MatchString(zipStr) {
			errors = append(errors, ValidationError{
				Field:   "zip_code",
				Message: "Invalid zip code format. Must be 5 digits or 5+4 format (e.g., 12345 or 12345-6789)",
			})
		}
	}

	// Validate state format (basic US state validation)
	if state, exists := address["state"]; exists && state != "" {
		stateStr := fmt.Sprintf("%v", state)
		stateRegex := regexp.MustCompile(`^[A-Z]{2}$`)
		if !stateRegex.MatchString(stateStr) {
			errors = append(errors, ValidationError{
				Field:   "state",
				Message: "Invalid state format. Must be 2-letter state code (e.g., CA, NY, TX)",
			})
		}
	}

	if len(errors) > 0 {
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Address validation failed",
		}
	}

	return &ValidationResult{Valid: true}
}

// ValidatePickupInfo validates pickup information
func (v *Validator) ValidatePickupInfo(pickupInfo map[string]interface{}) *ValidationResult {
	var errors []ValidationError

	// Check if location exists
	location, exists := pickupInfo["location"]
	if !exists {
		errors = append(errors, ValidationError{
			Field:   "location",
			Message: "location is required",
		})
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Pickup location is required",
		}
	}

	// Validate location structure
	locationMap, ok := location.(map[string]interface{})
	if !ok {
		errors = append(errors, ValidationError{
			Field:   "location",
			Message: "location must be an object",
		})
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Invalid location format",
		}
	}

	// Check if address exists in location
	address, exists := locationMap["address"]
	if !exists {
		errors = append(errors, ValidationError{
			Field:   "location.address",
			Message: "address is required in location",
		})
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Address is required in location",
		}
	}

	// Validate address
	addressMap, ok := address.(map[string]interface{})
	if !ok {
		errors = append(errors, ValidationError{
			Field:   "location.address",
			Message: "address must be an object",
		})
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Invalid address format",
		}
	}

	// Validate the address
	addressResult := v.ValidateAddress(addressMap)
	if !addressResult.Valid {
		// Add location prefix to field names
		for i, err := range addressResult.Errors {
			addressResult.Errors[i].Field = "location.address." + err.Field
		}
		return addressResult
	}

	return &ValidationResult{Valid: true}
}

// ValidateDropOffs validates drop-off locations
func (v *Validator) ValidateDropOffs(dropOffs []map[string]interface{}) *ValidationResult {
	var errors []ValidationError

	if len(dropOffs) == 0 {
		errors = append(errors, ValidationError{
			Field:   "drop_offs",
			Message: "At least one drop-off location is required",
		})
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Drop-off locations are required",
		}
	}

	// Validate each drop-off location
	for i, dropOff := range dropOffs {
		dropOffResult := v.ValidatePickupInfo(dropOff) // Reuse pickup validation logic
		if !dropOffResult.Valid {
			// Add index prefix to field names
			for _, err := range dropOffResult.Errors {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("drop_offs[%d].%s", i, err.Field),
					Message: err.Message,
				})
			}
		}
	}

	if len(errors) > 0 {
		return &ValidationResult{
			Valid:   false,
			Errors:  errors,
			Message: "Drop-off validation failed",
		}
	}

	return &ValidationResult{Valid: true}
}

// ValidateJSONString validates that a string contains valid JSON
func (v *Validator) ValidateJSONString(jsonStr, fieldName string) *ValidationResult {
	if jsonStr == "" {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   fieldName,
					Message: fmt.Sprintf("%s is required", fieldName),
				},
			},
			Message: fmt.Sprintf("%s is required", fieldName),
		}
	}

	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   fieldName,
					Message: fmt.Sprintf("Invalid JSON format: %v", err),
				},
			},
			Message: "Invalid JSON format",
		}
	}

	return &ValidationResult{Valid: true}
}

// ValidateCustomerTier validates customer tier parameter
func (v *Validator) ValidateCustomerTier(tier string) *ValidationResult {
	validTiers := []string{"bronze", "silver", "gold"}

	if tier == "" {
		return &ValidationResult{Valid: true} // Optional field
	}

	for _, validTier := range validTiers {
		if strings.EqualFold(tier, validTier) {
			return &ValidationResult{Valid: true}
		}
	}

	return &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{
				Field:   "customer_tier",
				Message: fmt.Sprintf("Invalid customer tier '%s'. Must be one of: %s", tier, strings.Join(validTiers, ", ")),
			},
		},
		Message: "Invalid customer tier",
	}
}

// ValidateNumericString validates that a string contains a valid number
func (v *Validator) ValidateNumericString(value, fieldName string, min, max int) *ValidationResult {
	if value == "" {
		return &ValidationResult{Valid: true} // Optional field
	}

	// Basic numeric validation
	numericRegex := regexp.MustCompile(`^\d+$`)
	if !numericRegex.MatchString(value) {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   fieldName,
					Message: fmt.Sprintf("%s must be a positive integer", fieldName),
				},
			},
			Message: "Invalid numeric format",
		}
	}

	return &ValidationResult{Valid: true}
}

// ValidateBooleanString validates that a string contains a valid boolean
func (v *Validator) ValidateBooleanString(value, fieldName string) *ValidationResult {
	if value == "" {
		return &ValidationResult{Valid: true} // Optional field
	}

	if value != "true" && value != "false" {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   fieldName,
					Message: fmt.Sprintf("%s must be 'true' or 'false'", fieldName),
				},
			},
			Message: "Invalid boolean format",
		}
	}

	return &ValidationResult{Valid: true}
}
