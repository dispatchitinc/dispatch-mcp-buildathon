package auth

import (
	"os"
)

func LoadConfig() (*Config, error) {
	return &Config{
		IDPEndpoint:   getEnv("IDP_ENDPOINT", ""),
		ClientID:      getEnv("IDP_CLIENT_ID", ""),
		ClientSecret:  getEnv("IDP_CLIENT_SECRET", ""),
		Scope:         getEnv("IDP_SCOPE", "dispatch:api"),
		TokenEndpoint: getEnv("IDP_TOKEN_ENDPOINT", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
