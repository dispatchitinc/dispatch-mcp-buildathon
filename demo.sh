#!/bin/bash

# Demo script for Dispatch MCP Server with IDP authentication
echo "🚀 Dispatch MCP Server Demo"
echo "=========================="

# Check if IDP authentication is enabled
if [ "$USE_IDP_AUTH" = "true" ]; then
    echo "🔐 Using IDP Authentication"
    echo "IDP Endpoint: $IDP_ENDPOINT"
    echo "Client ID: $IDP_CLIENT_ID"
    echo "Scope: $IDP_SCOPE"
else
    echo "🔑 Using Static Token Authentication"
    echo "Auth Token: ${DISPATCH_AUTH_TOKEN:0:20}..."
fi

echo ""
echo "📡 Dispatch GraphQL Endpoint: $DISPATCH_GRAPHQL_ENDPOINT"
echo "🏢 Organization ID: $DISPATCH_ORGANIZATION_ID"
echo ""

# Build the server
echo "🔨 Building MCP Server..."
./build.sh

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"
echo ""

# Test the server startup
echo "🧪 Testing server startup..."
timeout 3s ./bin/dispatch-mcp-server &
SERVER_PID=$!
sleep 2
kill $SERVER_PID 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Server starts successfully!"
else
    echo "❌ Server startup test failed!"
    exit 1
fi

echo ""
echo "🎯 Demo completed successfully!"
echo ""
echo "To run the server:"
echo "  ./bin/dispatch-mcp-server"
echo ""
echo "To test with MCP client:"
echo "  # Set up your MCP client to connect to this server"
echo "  # Use tools: create_estimate, create_order"
