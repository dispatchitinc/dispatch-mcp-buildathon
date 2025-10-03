package config

import (
	"os"
)

type Config struct {
	GraphQLEndpoint string
	AuthToken       string
	OrganizationID  string
	UseIDP          bool
	IDPEndpoint     string
	ClientID        string
	ClientSecret    string
	Scope           string
}

func Load() (*Config, error) {
	useIDP := getEnv("USE_IDP_AUTH", "false") == "true"

	config := &Config{
		GraphQLEndpoint: getEnv("DISPATCH_GRAPHQL_ENDPOINT", "https://monkey.graph.qa.dispatchfog.io/graphql"),
		OrganizationID:  getEnv("DISPATCH_ORGANIZATION_ID", ""),
		UseIDP:          useIDP,
	}

	if useIDP {
		config.IDPEndpoint = getEnv("IDP_ENDPOINT", "")
		config.ClientID = getEnv("IDP_CLIENT_ID", "")
		config.ClientSecret = getEnv("IDP_CLIENT_SECRET", "")
		config.Scope = getEnv("IDP_SCOPE", "dispatch:api")
	} else {
		config.AuthToken = getEnv("DISPATCH_AUTH_TOKEN", "")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
