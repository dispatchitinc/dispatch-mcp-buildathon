#!/bin/bash

echo "🚀 Starting Dispatch Chat with AI Hub Integration"
echo "================================================="
echo ""

# Set up AI Hub configuration
export USE_AI_HUB=true
export AI_HUB_ENDPOINT=https://aihub.dispatchit.com/v1
export ANTHROPIC_API_KEY=sk-hB4t-5i9G701-YI_gVc2Hw
export AI_HUB_MODEL=claude-sonnet

echo "🔧 Configuration:"
echo "- AI Hub: ✅ Enabled"
echo "- Endpoint: https://aihub.dispatchit.com/v1"
echo "- Model: claude-sonnet"
echo "- Cost Tracking: ✅ Active"
echo ""

echo "💬 Starting conversational pricing advisor..."
echo "Type 'quit' to exit, 'help' for examples."
echo ""

# Start the chat interface
../bin/dispatch-cli chat
