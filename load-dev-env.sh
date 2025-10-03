#!/bin/bash
# Development Environment Loader
# This script loads development environment variables

echo "🔧 Loading development environment..."

# AI Hub Configuration (Primary AI Service)
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1
export ANTHROPIC_API_KEY=sk-hB4t-5i9G701-YI_gVc2Hw

# Dispatch API Configuration (Monkey Subenv)
export DISPATCH_GRAPHQL_ENDPOINT=https://graphql-gateway.monkey.dispatchfog.org/graphql
export DISPATCH_ORGANIZATION_ID=your_org_id_here

# Authentication Method (choose one)
export USE_IDP_AUTH=false

# Option 1: Static Token Authentication
export DISPATCH_AUTH_TOKEN=your_static_auth_token_here

echo "✅ Development environment loaded!"
echo "🤖 AI Hub: Enabled"
echo "🔑 API Key: Set"
echo "🌐 Endpoint: $AI_HUB_ENDPOINT"
echo ""
echo "You can now run:"
echo "  ./bin/dispatch-web"
echo "  ./bin/dispatch-cli chat"
