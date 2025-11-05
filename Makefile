.PHONY: test test-unit test-integration test-coverage build clean lint vet fmt

# Default target
all: test build

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	go test -v ./cmd/...

# Run integration tests  
test-integration:
	go test -v ./test/...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./cmd/... ./test/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the binary
build:
	go build -o backup ./cmd/

# Clean build artifacts
clean:
	rm -f backup coverage.out coverage.html

# Run linting
lint:
	@which golint > /dev/null || go install golang.org/x/lint/golint@latest
	golint ./cmd/... ./test/...

# Run go vet
vet:
	go vet ./cmd/... ./test/...

# Check formatting
fmt:
	go fmt ./cmd/... ./test/...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run all checks (lint, vet, fmt, test)
check: lint vet fmt test

# Development setup
dev-setup: deps
	@echo "Development environment ready"
	@echo "Run 'make test' to run tests"
	@echo "Run 'make build' to build the binary"
	@echo "Run 'make check' to run all checks"
