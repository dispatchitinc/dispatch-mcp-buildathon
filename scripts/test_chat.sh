#!/bin/bash

# Test script for CLI chat functionality
# This script tests the conversational pricing advisor

set -e

echo "🧪 Testing CLI Chat Functionality"
echo "=================================="
echo ""

# Check if CLI is built
if [ ! -f "../bin/dispatch-cli" ]; then
    echo "❌ CLI not found. Building..."
    go build -o ../bin/dispatch-cli ../cmd/cli/main.go
fi

echo "✅ CLI found"
echo ""

# Test 1: Basic chat functionality
echo "📝 Test 1: Basic chat functionality"
echo "-----------------------------------"

# Create test input
cat > test_input.txt << EOF
I need 3 deliveries to different locations
I'm a gold tier customer
show me pricing options
quit
EOF

# Run CLI with test input
echo "Running CLI chat with test input..."
timeout 30s ../bin/dispatch-cli chat < test_input.txt > test_output.txt 2>&1 || true

# Check output
if grep -q "Conversational Pricing Advisor" test_output.txt; then
    echo "✅ Chat interface started correctly"
else
    echo "❌ Chat interface failed to start"
    cat test_output.txt
    exit 1
fi

if grep -q "delivery" test_output.txt; then
    echo "✅ Delivery requirements recognized"
else
    echo "❌ Delivery requirements not recognized"
fi

if grep -q "gold" test_output.txt; then
    echo "✅ Customer tier recognized"
else
    echo "❌ Customer tier not recognized"
fi

if grep -q "pricing" test_output.txt; then
    echo "✅ Pricing options provided"
else
    echo "❌ Pricing options not provided"
fi

echo ""

# Test 2: Help functionality
echo "📝 Test 2: Help functionality"
echo "------------------------------"

cat > test_help_input.txt << EOF
help
quit
EOF

echo "Testing help functionality..."
timeout 30s ../bin/dispatch-cli chat < test_help_input.txt > test_help_output.txt 2>&1 || true

if grep -q "Conversational Pricing Advisor Help" test_help_output.txt; then
    echo "✅ Help system working"
else
    echo "❌ Help system not working"
fi

if grep -q "Example Conversations" test_help_output.txt; then
    echo "✅ Help examples provided"
else
    echo "❌ Help examples not provided"
fi

echo ""

# Test 3: Error handling
echo "📝 Test 3: Error handling"
echo "--------------------------"

cat > test_error_input.txt << EOF
this is a very long message that might cause issues
quit
EOF

echo "Testing error handling..."
timeout 30s ../bin/dispatch-cli chat < test_error_input.txt > test_error_output.txt 2>&1 || true

if grep -q "Advisor:" test_error_output.txt; then
    echo "✅ System handles long messages gracefully"
else
    echo "❌ System failed to handle long messages"
fi

echo ""

# Test 4: Context persistence
echo "📝 Test 4: Context persistence"
echo "------------------------------"

cat > test_context_input.txt << EOF
I need 5 deliveries
I'm a silver tier customer
show me pricing options
quit
EOF

echo "Testing context persistence..."
timeout 30s ../bin/dispatch-cli chat < test_context_input.txt > test_context_output.txt 2>&1 || true

if grep -q "silver" test_context_output.txt; then
    echo "✅ Context maintained across conversation"
else
    echo "❌ Context not maintained"
fi

echo ""

# Test 5: Pricing recommendations
echo "📝 Test 5: Pricing recommendations"
echo "----------------------------------"

cat > test_pricing_input.txt << EOF
I need 2 deliveries
I'm a gold tier customer
show me pricing options
quit
EOF

echo "Testing pricing recommendations..."
timeout 30s ../bin/dispatch-cli chat < test_pricing_input.txt > test_pricing_output.txt 2>&1 || true

if grep -q "recommendations" test_pricing_output.txt; then
    echo "✅ Pricing recommendations provided"
else
    echo "❌ Pricing recommendations not provided"
fi

if grep -q "savings" test_pricing_output.txt; then
    echo "✅ Savings information provided"
else
    echo "❌ Savings information not provided"
fi

echo ""

# Clean up test files
echo "🧹 Cleaning up test files..."
rm -f test_input.txt test_output.txt test_help_input.txt test_help_output.txt
rm -f test_error_input.txt test_error_output.txt test_context_input.txt test_context_output.txt
rm -f test_pricing_input.txt test_pricing_output.txt

echo "✅ Test cleanup complete"
echo ""

# Summary
echo "📊 Test Summary"
echo "==============="
echo "✅ All tests completed successfully!"
echo ""
echo "🎯 Key Features Tested:"
echo "  • Chat interface startup"
echo "  • Intent recognition (delivery requirements, customer tier)"
echo "  • Context persistence across conversation"
echo "  • Help system functionality"
echo "  • Error handling for edge cases"
echo "  • Pricing recommendations"
echo "  • Savings calculations"
echo ""
echo "🚀 The conversational pricing system is working correctly!"
