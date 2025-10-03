# 🧪 Testing Guide - Conversational Pricing System

## Overview

This guide covers all the testing approaches for the conversational pricing system, including unit tests, integration tests, and CLI response validation.

## 🚀 Quick Start Testing

### Run All Tests
```bash
# Unit tests
go test ./test -v

# CLI integration tests
./test_chat.sh

# Response validation tests
./test_cli_responses.sh
```

## 📋 Test Categories

### 1. Unit Tests (`go test ./test -v`)

#### **Conversation Engine Tests**
- **Basic Message Processing**: Tests that messages are processed correctly
- **Context Building**: Tests that conversation context is maintained
- **Pricing Recommendations**: Tests that appropriate recommendations are generated
- **Intent Recognition**: Tests that different intents are recognized correctly
- **Next Questions**: Tests that follow-up questions are generated
- **Error Handling**: Tests that errors are handled gracefully

#### **Context Manager Tests**
- **Session Management**: Tests that conversation sessions are managed correctly
- **Session Stats**: Tests that session statistics are generated

#### **Conversation Flow Tests**
- **Multi-step Conversations**: Tests complete conversation flows
- **Context Persistence**: Tests that context is maintained across steps

#### **Pricing Scenario Tests**
- **Different Customer Tiers**: Tests bronze, silver, gold tier scenarios
- **Delivery Counts**: Tests different delivery count scenarios
- **Eligibility**: Tests that pricing model eligibility is calculated correctly

### 2. Integration Tests (`./test_chat.sh`)

#### **CLI Chat Functionality**
- **Chat Interface Startup**: Tests that the chat interface starts correctly
- **Intent Recognition**: Tests delivery requirements and customer tier recognition
- **Context Persistence**: Tests that context is maintained across conversation
- **Help System**: Tests that the help system works correctly
- **Error Handling**: Tests that edge cases are handled gracefully
- **Pricing Recommendations**: Tests that pricing recommendations are provided
- **Savings Calculations**: Tests that savings information is calculated

### 3. Response Validation Tests (`./test_cli_responses.sh`)

#### **Specific Response Validation**
- **Delivery Requirements Recognition**: Tests that delivery requirements are recognized
- **Customer Tier Recognition**: Tests that customer tiers are recognized
- **Pricing Recommendations**: Tests that specific pricing models are mentioned
- **Context Persistence**: Tests that context is maintained across conversation
- **Help System**: Tests that help information is provided
- **Error Handling**: Tests that long messages are handled gracefully
- **Specific Response Validation**: Tests that specific pricing models are mentioned

## 🎯 Test Scenarios

### **Scenario 1: New Customer**
```
Input: "I need 3 deliveries to different locations"
Expected: Multi-Delivery discount mentioned, 15% savings
```

### **Scenario 2: Gold Tier Customer**
```
Input: "I'm a gold tier customer"
Expected: Loyalty discount mentioned, 10% savings
```

### **Scenario 3: High-Volume Customer**
```
Input: "I need 5 deliveries and I'm a gold tier customer"
Expected: Volume discount mentioned, 20% savings
```

### **Scenario 4: Bulk Order Customer**
```
Input: "I need 10 deliveries and it's a bulk order"
Expected: Bulk order discount mentioned, 25% savings
```

## 🔍 Test Validation

### **Response Quality Checks**
- **Natural Language**: Responses should sound conversational
- **Accuracy**: Pricing calculations should be correct
- **Completeness**: All relevant information should be provided
- **Context Awareness**: Responses should reference previous conversation

### **Intent Recognition Checks**
- **Delivery Requirements**: "I need X deliveries" → Multi-Delivery discount
- **Customer Tier**: "I'm gold tier" → Loyalty discount
- **Pricing Questions**: "show me pricing" → All options
- **Recommendations**: "what's best" → Best option

### **Context Management Checks**
- **Session Persistence**: Context maintained across conversation
- **Profile Building**: Customer profile built from conversation
- **History Tracking**: Delivery and pricing history maintained

## 🛠️ Test Implementation

### **Unit Test Structure**
```go
func TestConversationEngine(t *testing.T) {
    engine := conversation.NewConversationEngine()
    
    t.Run("basic_message_processing", func(t *testing.T) {
        response, err := engine.ProcessMessage("I need 3 deliveries", nil)
        // Test assertions
    })
}
```

### **Integration Test Structure**
```bash
# Create test input
cat > test_input.txt << EOF
I need 3 deliveries to different locations
I'm a gold tier customer
show me pricing options
quit
EOF

# Run CLI with test input
timeout 30s ./bin/dispatch-cli chat < test_input.txt > test_output.txt 2>&1 || true

# Validate output
if grep -q "delivery" test_output.txt; then
    echo "✅ Delivery requirements recognized"
fi
```

### **Response Validation Structure**
```bash
# Test specific responses
if grep -q "Multi-Delivery" test_output.txt; then
    echo "✅ Multi-Delivery discount mentioned"
else
    echo "❌ Multi-Delivery discount not mentioned"
fi
```

## 📊 Test Results

### **Expected Test Results**
- **Unit Tests**: All tests should pass
- **Integration Tests**: All features should work correctly
- **Response Validation**: All responses should be accurate and complete

### **Common Issues and Solutions**

#### **Issue: Delivery requirements not recognized**
- **Cause**: Intent recognition not working correctly
- **Solution**: Check intent patterns in conversation engine

#### **Issue: Pricing recommendations not provided**
- **Cause**: Pricing engine not integrated correctly
- **Solution**: Check pricing engine integration

#### **Issue: Context not maintained**
- **Cause**: Context manager not working correctly
- **Solution**: Check context management implementation

## 🚀 Continuous Testing

### **Automated Testing**
```bash
# Run all tests
make test

# Run specific test categories
make test-unit
make test-integration
make test-responses
```

### **Test Coverage**
- **Unit Tests**: 90%+ coverage of conversation engine
- **Integration Tests**: All CLI functionality covered
- **Response Validation**: All response types validated

## 📈 Performance Testing

### **Response Time Tests**
- **Message Processing**: < 100ms for simple messages
- **Pricing Calculations**: < 50ms for recommendations
- **Context Updates**: < 10ms for context updates

### **Load Testing**
- **Concurrent Conversations**: Test multiple simultaneous conversations
- **Memory Usage**: Monitor memory usage during extended conversations
- **Session Management**: Test session cleanup and expiration

## 🔧 Debugging Tests

### **Debug Mode**
```bash
export DEBUG=true
./bin/dispatch-cli chat
```

### **Verbose Output**
```bash
go test ./test -v -run TestConversationEngine
```

### **Log Analysis**
```bash
# Check conversation logs
tail -f conversation.log

# Check pricing engine logs
tail -f pricing.log
```

## 📚 Test Documentation

### **Test Reports**
- **Unit Test Report**: `test-results-unit.txt`
- **Integration Test Report**: `test-results-integration.txt`
- **Response Validation Report**: `test-results-responses.txt`

### **Coverage Reports**
- **Code Coverage**: `coverage.html`
- **Function Coverage**: `coverage-functions.txt`
- **Line Coverage**: `coverage-lines.txt`

## 🎯 Best Practices

### **Test Writing**
- **Clear Test Names**: Use descriptive test names
- **Single Responsibility**: Each test should test one thing
- **Independent Tests**: Tests should not depend on each other
- **Cleanup**: Always clean up test data

### **Test Maintenance**
- **Regular Updates**: Update tests when features change
- **Version Control**: Keep test files in version control
- **Documentation**: Document test scenarios and expected results

### **Test Automation**
- **CI/CD Integration**: Run tests automatically on code changes
- **Scheduled Testing**: Run comprehensive tests regularly
- **Alert System**: Alert on test failures

---

*This testing guide ensures that the conversational pricing system works correctly and provides accurate responses to users.*
