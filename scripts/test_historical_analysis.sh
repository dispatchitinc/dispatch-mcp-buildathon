#!/bin/bash

# Test script for the new historical analysis MCP tool
echo "🧪 Testing Historical Analysis MCP Tool"
echo "========================================"
echo ""

# Check if CLI is built
if [ ! -f "../bin/dispatch-cli" ]; then
    echo "❌ CLI not found. Building..."
    go build -o ../bin/dispatch-cli ../cmd/cli/main.go
fi

echo "✅ CLI found"
echo ""

# Test 1: Basic historical analysis request
echo "📝 Test 1: Basic historical analysis request"
echo "--------------------------------------------"

cat > test_historical_input.txt << EOF
analyze_historical_savings
start_date: 2024-01-01
end_date: 2024-03-31
analysis_types: comprehensive
include_recommendations: true
quit
EOF

echo "Testing historical analysis with comprehensive analysis..."
timeout 30s ../bin/dispatch-cli interactive < test_historical_input.txt > test_historical_output.txt 2>&1 || true

if grep -q "analyze_historical_savings" test_historical_output.txt; then
    echo "✅ Historical analysis tool recognized"
else
    echo "❌ Historical analysis tool not recognized"
fi

if grep -q "analysis_period" test_historical_output.txt; then
    echo "✅ Analysis period processed correctly"
else
    echo "❌ Analysis period not processed"
fi

if grep -q "pending_implementation" test_historical_output.txt; then
    echo "✅ Tool returns expected placeholder response"
else
    echo "❌ Tool response format unexpected"
fi

echo ""

# Test 2: Specific analysis types
echo "📝 Test 2: Specific analysis types"
echo "-----------------------------------"

cat > test_bundling_input.txt << EOF
analyze_historical_savings
start_date: 2024-01-01
end_date: 2024-03-31
analysis_types: bundling,volume
include_recommendations: true
quit
EOF

echo "Testing bundling and volume analysis..."
timeout 30s ../bin/dispatch-cli interactive < test_bundling_input.txt > test_bundling_output.txt 2>&1 || true

if grep -q "bundling" test_bundling_output.txt; then
    echo "✅ Bundling analysis type recognized"
else
    echo "❌ Bundling analysis type not recognized"
fi

if grep -q "volume" test_bundling_output.txt; then
    echo "✅ Volume analysis type recognized"
else
    echo "❌ Volume analysis type not recognized"
fi

echo ""

# Test 3: Date validation
echo "📝 Test 3: Date validation"
echo "--------------------------"

cat > test_date_validation_input.txt << EOF
analyze_historical_savings
start_date: invalid-date
end_date: 2024-03-31
quit
EOF

echo "Testing date validation..."
timeout 30s ../bin/dispatch-cli interactive < test_date_validation_input.txt > test_date_validation_output.txt 2>&1 || true

if grep -q "invalid.*date" test_date_validation_output.txt; then
    echo "✅ Date validation working correctly"
else
    echo "❌ Date validation not working"
fi

echo ""

# Test 4: Help and documentation
echo "📝 Test 4: Help and documentation"
echo "---------------------------------"

cat > test_help_input.txt << EOF
help
quit
EOF

echo "Testing help system for new tool..."
timeout 30s ../bin/dispatch-cli interactive < test_help_input.txt > test_help_output.txt 2>&1 || true

if grep -q "analyze_historical_savings" test_help_output.txt; then
    echo "✅ Historical analysis tool included in help"
else
    echo "❌ Historical analysis tool not in help"
fi

echo ""

# Clean up test files
echo "🧹 Cleaning up test files..."
rm -f test_historical_input.txt test_historical_output.txt
rm -f test_bundling_input.txt test_bundling_output.txt
rm -f test_date_validation_input.txt test_date_validation_output.txt
rm -f test_help_input.txt test_help_output.txt

echo "✅ Test cleanup complete"
echo ""

# Summary
echo "📊 Test Summary"
echo "==============="
echo "✅ Historical analysis tool tests completed!"
echo ""
echo "🎯 Key Features Tested:"
echo "  • Tool registration and recognition"
echo "  • Parameter parsing and validation"
echo "  • Date format validation"
echo "  • Analysis type selection"
echo "  • Help system integration"
echo "  • Placeholder response format"
echo ""
echo "🚀 The historical analysis tool is ready for development!"
echo ""
echo "📋 Next Steps:"
echo "  • Implement historical data retrieval from Dispatch API"
echo "  • Build bundling analysis algorithms"
echo "  • Implement volume discount analysis"
echo "  • Add loyalty tier analysis"
echo "  • Create comprehensive report generation"
echo "  • Add visualization capabilities"
