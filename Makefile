# Project variables
BINARY_NAME=gauth-extractor
MAIN_PACKAGE=./cmd/extractor
PROTO_DIR=./internal/proto

# Go related variables
GOBASE=$(shell pwd)
BUILD_DIR=$(GOBASE)/build

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test help proto docker docker-run

all: test build

## Build the application
build:
	@echo "Building..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Install the application
install:
	@echo "Installing..."
	go install $(MAIN_PACKAGE)
	@echo "Install complete"

## Run the application
run:
	@echo "Running..."
	go run $(MAIN_PACKAGE) -i

## Clean build files
clean:
	@echo "Cleaning..."
	go clean
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## Generate protobuf code
proto:
	@echo "Generating Protocol Buffer code..."
	protoc --go_out=. --go_opt=paths=source_relative $(PROTO_DIR)/google_auth.proto
	@echo "Protocol Buffer code generation complete"

## Run unit tests
test:
	@echo "Running tests..."
	go test -v ./...

## Run unit tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t google-auth-extractor:latest .
	@echo "Docker image built: google-auth-extractor:latest"

## Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -it --rm google-auth-extractor:latest

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)