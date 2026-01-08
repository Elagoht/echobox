.PHONY: build run test clean install fmt vet lint

BINARY_NAME=echobox
BUILD_DIR=bin
CMD_DIR=cmd/echobox

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

run:
	@echo "Running $(BINARY_NAME)..."
	@go run ./$(CMD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(CMD_DIR)

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Vetting code..."
	@go vet ./...

lint:
	@echo "Linting code..."
	@golangci-lint run ./... || echo "golangci-lint not installed, skipping..."

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

mod-verify:
	@echo "Verifying dependencies..."
	@go mod verify

help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install the binary to GOPATH/bin"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  mod-verify     - Verify dependencies"
	@echo "  help           - Show this help message"
