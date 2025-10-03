#!/bin/bash

# Demo script for Dispatch CLI
echo "🚀 Dispatch CLI Demo"
echo "===================="

# Check if CLI exists
if [ ! -f "../bin/dispatch-cli" ]; then
    echo "❌ CLI not found. Building..."
    ./build.sh
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
fi

echo ""
echo "🔍 Checking configuration..."
../bin/dispatch-cli status

echo ""
echo "📊 Demo 1: Creating Cost Estimate"
echo "================================="
../bin/dispatch-cli estimate

echo ""
echo "📦 Demo 2: Creating Delivery Order"
echo "==================================="
../bin/dispatch-cli order

echo ""
echo "🎮 Demo 3: Interactive Mode"
echo "=========================="
echo "Starting interactive mode. Type 'help' for commands, 'quit' to exit."
echo ""

# Start interactive mode
../bin/dispatch-cli interactive
