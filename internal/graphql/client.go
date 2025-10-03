package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphQLClient handles GraphQL requests
type GraphQLClient struct {
	endpoint string
	client   *http.Client
	headers  map[string]string
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message   string            `json:"message"`
	Locations []GraphQLLocation `json:"locations,omitempty"`
	Path      []interface{}     `json:"path,omitempty"`
}

// GraphQLLocation represents a location in GraphQL source
type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(endpoint string) *GraphQLClient {
	return &GraphQLClient{
		endpoint: endpoint,
		client:   &http.Client{},
		headers:  make(map[string]string),
	}
}

// SetHeader sets a header for all requests
func (c *GraphQLClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// Execute executes a GraphQL query or mutation
func (c *GraphQLClient) Execute(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	payload := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL request failed with status %d", resp.StatusCode)
	}

	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", result.Errors)
	}

	return &result, nil
}
