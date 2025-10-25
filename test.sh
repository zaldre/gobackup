#!/bin/bash

# Test runner script for gobackup
set -e

echo "Running tests for gobackup..."

# Run all tests
echo "Running unit tests..."
go test -v ./cmd/...

echo "Running integration tests..."
go test -v ./test/...

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -coverprofile=coverage.out ./cmd/... ./test/...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Run linting
echo "Running linting..."
go vet ./cmd/... ./test/...

# Check formatting
echo "Checking formatting..."
go fmt ./cmd/... ./test/...

echo "Tests completed successfully!"
echo "Coverage report generated: coverage.html"
