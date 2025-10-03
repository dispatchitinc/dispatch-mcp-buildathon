package main

import (
	"dispatch-mcp-server/internal/mcp"
	"fmt"
	"log"
	"os"
)

func main() {
	server, err := mcp.NewMCPServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Starting Dispatch MCP server...\n")
	if err := server.Run(); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
