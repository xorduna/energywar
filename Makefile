# Energy War Game Makefile

# Variables
BINARY_NAME=energywar
GO=go
SWAG=$(shell go env GOPATH)/bin/swag
MAIN_FILE=cmd/server/main.go

# Default target
.PHONY: all
all: build

# Install dependencies
.PHONY: deps
deps:
	$(GO) get -u github.com/swaggo/swag/cmd/swag
	$(GO) get -u github.com/labstack/echo/v4
	$(GO) get -u github.com/swaggo/echo-swagger

# Generate Swagger documentation
.PHONY: docs
docs:
	$(SWAG) init -g $(MAIN_FILE)

# Build the binary
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf docs

# Test the application
.PHONY: test
test:
	$(GO) test ./...

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all      - Build the application (default)"
	@echo "  deps     - Install dependencies"
	@echo "  docs  - Generate Swagger documentation"
	@echo "  build    - Build the binary"
	@echo "  run      - Run the application"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help message"
