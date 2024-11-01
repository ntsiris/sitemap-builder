# Project variables
APP_NAME := sitemap
VERSION := 0.0.1
BIN_DIR := bin
SRC_DIR := cmd
PKG_DIR := ./...

# Docker variables
DOCKER_IMAGE = $(APP_NAME):$(VERSION)
DOCKER_FILE = ./Dockerfile

# Compiler flags
GO := go
GOFLAGS := -v

# Test configuration
TEST_DIR := ./...

# Flag variable
URL ?= ""
DEPTH ?= ""
FILE ?= ""

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(SRC_DIR)/*

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -cover $(GOFLAGS) $(TEST_DIR)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)/*
	$(GO) clean -testcache

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt $(PKG_DIR)

# Lint code (requires golangci-lint to be installed)
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run $(PKG_DIR)

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Run application
.PHONY: run
run: build
	@echo "Running the application with URL=$(URL)"
	$(BIN_DIR)/$(APP_NAME) -url=$(URL) -depth=$(DEPTH) -out=$(FILE)


## Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f $(DOCKER_FILE) .