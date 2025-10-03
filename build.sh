#!/bin/bash
echo "🔨 Building Dispatch MCP Server..."
go build -o bin/dispatch-mcp-server ./cmd/server

echo "🔨 Building Dispatch CLI..."
go build -o bin/dispatch-cli ./cmd/cli

echo "✅ Build complete!"
echo "📦 Binaries created:"
echo "  - bin/dispatch-mcp-server (MCP Server)"
echo "  - bin/dispatch-cli (CLI Demo Tool)"
