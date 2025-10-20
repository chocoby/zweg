.PHONY: build test clean install lint fmt vet coverage help

# Build variables
BINARY_NAME=zweg
BUILD_DIR=./bin
CMD_DIR=./cmd/zweg
GO=go

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -timeout 30s ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(CMD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean

# Run the application (requires input and output arguments)
run: build
	@echo "Running $(BINARY_NAME)..."
	@if [ -z "$(INPUT)" ] || [ -z "$(OUTPUT)" ]; then \
		echo "Usage: make run INPUT=<input.json> OUTPUT=<output.gpx> [TRACK=<track_name>]"; \
		exit 1; \
	fi
	$(BUILD_DIR)/$(BINARY_NAME) $(INPUT) $(OUTPUT) $(TRACK)

# Check code quality (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  test       - Run all tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  bench      - Run benchmarks"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run golangci-lint"
	@echo "  install    - Install the binary"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Run the application (INPUT=<file> OUTPUT=<file> [TRACK=<name>])"
	@echo "  check      - Run fmt, vet, and test"
	@echo "  help       - Show this help message"
