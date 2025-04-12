# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMODTIDY=$(GOCMD) mod tidy
BINARY_NAME=backend-dragonhak
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
LDFLAGS=-ldflags "-s -w"

# Directories
SRC_DIR=.
TEST_DIR=./handlers
BIN_DIR=./bin
GOLANGCI_LINT=$(BIN_DIR)/golangci-lint

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build test run lint clean

# Default target - runs the most common tasks in sequence
default: deps build lint test

build:
	@echo "$(BLUE)Building application...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(SRC_DIR)
	@echo "$(GREEN)Build completed!$(NC)"

test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v -cover ./handlers/...
	@echo "$(GREEN)Tests completed!$(NC)"

test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./handlers/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated!$(NC)"

run:
	@echo "$(BLUE)Running application...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)
	./$(BINARY_NAME)

clean:
	@echo "$(BLUE)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html
	rm -rf $(BIN_DIR)
	@echo "$(GREEN)Clean completed!$(NC)"

deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMODTIDY)
	$(GOGET) -v ./...
	@echo "$(BLUE)Installing golangci-lint...$(NC)"
	@mkdir -p $(BIN_DIR)
	@if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.55.2; \
	fi
	@echo "$(GREEN)Dependencies installed!$(NC)"

lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "$(RED)golangci-lint not found. Please run 'make deps' first$(NC)"; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run
	@echo "$(GREEN)Lint completed!$(NC)"
	@echo "$(BLUE)Cleaning up...$(NC)"
	rm -f $(BINARY_NAME)
	rm -rf $(BIN_DIR)
	@echo "$(GREEN)Cleanup completed!$(NC)"

# Cross compilation
build-linux:
	@echo "$(BLUE)Building for Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) $(SRC_DIR)
	@echo "$(GREEN)Linux build completed!$(NC)"