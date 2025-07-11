.PHONY: build run test clean fmt lint help

# Binary name
BINARY_NAME=api-server

# Build the application
build:
	go build -o bin/$(BINARY_NAME) ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Show help
help:
	@echo "Available commands:"
	@echo "  build  - Build the application"
	@echo "  run    - Run the application"
	@echo "  test   - Run tests"
	@echo "  clean  - Clean build artifacts"
	@echo "  fmt    - Format code"
	@echo "  lint   - Lint code"
	@echo "  deps   - Install dependencies"
	@echo "  help   - Show this help" 