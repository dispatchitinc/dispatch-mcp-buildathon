package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	config     *Config
	token      *Token
}

type Config struct {
	IDPEndpoint   string
	ClientID      string
	ClientSecret  string
	Scope         string
	TokenEndpoint string
}

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"-"`
	Scope       string    `json:"scope"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
}

func NewClient(config *Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

func (c *Client) GetValidToken() (string, error) {
	// Check if we have a valid token
	if c.token != nil && time.Now().Before(c.token.ExpiresAt) {
		return c.token.AccessToken, nil
	}

	// Token expired or doesn't exist, get a new one
	return c.refreshToken()
}

func (c *Client) refreshToken() (string, error) {
	tokenRequest := TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        c.config.Scope,
	}

	jsonData, err := json.Marshal(tokenRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %v", err)
	}

	req, err := http.NewRequest("POST", c.config.TokenEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	// Calculate expiration time
	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	c.token = &token

	return token.AccessToken, nil
}

func (c *Client) IsTokenValid() bool {
	return c.token != nil && time.Now().Before(c.token.ExpiresAt)
}
