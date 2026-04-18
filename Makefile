.PHONY: all build test clean example tidy

# Default target
all: tidy test build

# Tidy module dependencies
tidy:
	@echo "Tidying module dependencies..."
	go mod tidy

# Build the project
build:
	@echo "Building..."
	go build -v ./...

# Run tests with race detector
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run the provided example
example:
	@echo "Running example..."
	go run examples/basic_scan/main.go

# Clean build cache
clean:
	@echo "Cleaning..."
	go clean
