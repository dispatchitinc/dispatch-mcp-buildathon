# Makefile for Dispatch MCP Server Testing

.PHONY: all test test-unit test-integration test-responses test-coverage build clean help

# Default target
all: build test

# Build the project
build:
	@echo "🔨 Building project..."
	go mod tidy
	go build -o bin/dispatch-cli cmd/cli/main.go
	go build -o bin/dispatch-mcp-server cmd/server/main.go
	@echo "✅ Build complete"

# Run all tests
test: test-unit test-integration test-responses
	@echo "🎉 All tests completed!"

# Run unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	go test ./test -v
	@echo "✅ Unit tests completed"

# Run integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	./test_chat.sh
	@echo "✅ Integration tests completed"

# Run response validation tests
test-responses:
	@echo "🧪 Running response validation tests..."
	./test_cli_responses.sh
	@echo "✅ Response validation tests completed"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test ./test -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run specific test
test-specific:
	@echo "🧪 Running specific test..."
	go test ./test -v -run $(TEST)
	@echo "✅ Specific test completed"

# Run conversation tests only
test-conversation:
	@echo "🧪 Running conversation tests..."
	go test ./test -v -run TestConversation
	@echo "✅ Conversation tests completed"

# Run context tests only
test-context:
	@echo "🧪 Running context tests..."
	go test ./test -v -run TestContext
	@echo "✅ Context tests completed"

# Run pricing tests only
test-pricing:
	@echo "🧪 Running pricing tests..."
	go test ./test -v -run TestPricing
	@echo "✅ Pricing tests completed"

# Run CLI chat demo
demo-chat:
	@echo "🗣️  Running CLI chat demo..."
	@echo "Type 'quit' to exit"
	./bin/dispatch-cli chat

# Run pricing comparison demo
demo-pricing:
	@echo "💰 Running pricing comparison demo..."
	./bin/dispatch-cli pricing

# Run estimate demo
demo-estimate:
	@echo "📊 Running estimate demo..."
	./bin/dispatch-cli estimate

# Run order demo
demo-order:
	@echo "📦 Running order demo..."
	./bin/dispatch-cli order

# Run all demos
demo: demo-chat demo-pricing demo-estimate demo-order
	@echo "🎉 All demos completed!"

# Clean up build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -f bin/dispatch-cli bin/dispatch-mcp-server
	rm -f coverage.out coverage.html
	rm -f test_*.txt test_*.log
	@echo "✅ Cleanup complete"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Code linted"

# Security scan
security:
	@echo "🔒 Running security scan..."
	gosec ./...
	@echo "✅ Security scan completed"

# Performance test
perf:
	@echo "⚡ Running performance tests..."
	go test ./test -v -run TestPerformance -bench=.
	@echo "✅ Performance tests completed"

# Memory test
memory:
	@echo "🧠 Running memory tests..."
	go test ./test -v -run TestMemory -memprofile=mem.prof
	go tool pprof mem.prof
	@echo "✅ Memory tests completed"

# Stress test
stress:
	@echo "💪 Running stress tests..."
	go test ./test -v -run TestStress -count=100
	@echo "✅ Stress tests completed"

# Generate test report
report:
	@echo "📊 Generating test report..."
	go test ./test -v -json > test-results.json
	@echo "✅ Test report generated: test-results.json"

# Run tests in parallel
test-parallel:
	@echo "🧪 Running tests in parallel..."
	go test ./test -v -parallel 4
	@echo "✅ Parallel tests completed"

# Run tests with race detection
test-race:
	@echo "🧪 Running tests with race detection..."
	go test ./test -v -race
	@echo "✅ Race detection tests completed"

# Run tests with timeout
test-timeout:
	@echo "🧪 Running tests with timeout..."
	go test ./test -v -timeout 30s
	@echo "✅ Timeout tests completed"

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests with verbose output..."
	go test ./test -v -v
	@echo "✅ Verbose tests completed"

# Run tests with short flag
test-short:
	@echo "🧪 Running short tests..."
	go test ./test -v -short
	@echo "✅ Short tests completed"

# Run tests with failfast flag
test-failfast:
	@echo "🧪 Running tests with failfast..."
	go test ./test -v -failfast
	@echo "✅ Failfast tests completed"

# Help
help:
	@echo "🚀 Dispatch MCP Server Testing Makefile"
	@echo "======================================"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build and run all tests"
	@echo "  build            - Build the project"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-responses   - Run response validation tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-specific    - Run specific test (use TEST=TestName)"
	@echo "  test-conversation- Run conversation tests only"
	@echo "  test-context     - Run context tests only"
	@echo "  test-pricing     - Run pricing tests only"
	@echo "  demo-chat        - Run CLI chat demo"
	@echo "  demo-pricing     - Run pricing comparison demo"
	@echo "  demo-estimate    - Run estimate demo"
	@echo "  demo-order       - Run order demo"
	@echo "  demo             - Run all demos"
	@echo "  clean            - Clean up build artifacts"
	@echo "  deps             - Install dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo "  security         - Run security scan"
	@echo "  perf             - Run performance tests"
	@echo "  memory           - Run memory tests"
	@echo "  stress           - Run stress tests"
	@echo "  report           - Generate test report"
	@echo "  test-parallel    - Run tests in parallel"
	@echo "  test-race        - Run tests with race detection"
	@echo "  test-timeout     - Run tests with timeout"
	@echo "  test-verbose     - Run tests with verbose output"
	@echo "  test-short       - Run short tests"
	@echo "  test-failfast    - Run tests with failfast"
	@echo "  help             - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make test                    # Run all tests"
	@echo "  make test-unit               # Run unit tests only"
	@echo "  make test-specific TEST=TestConversationEngine  # Run specific test"
	@echo "  make demo-chat              # Run CLI chat demo"
	@echo "  make test-coverage          # Run tests with coverage"
	@echo "  make clean                  # Clean up build artifacts"
