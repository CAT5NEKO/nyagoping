APP_NAME := nyagoping
BUILD_DIR := bin
MAIN_PATH := ./cmd/nyagoping

PLATFORMS := windows/amd64 windows/386 linux/amd64 linux/386 darwin/amd64 darwin/arm64
WINDOWS_EXT := .exe

.PHONY: all build clean test coverage cross-build help

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(if $(filter windows,$(shell go env GOOS)),$(WINDOWS_EXT),) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)$(if $(filter windows,$(shell go env GOOS)),$(WINDOWS_EXT),)"

cross-build: clean
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		GOOS=$(word 1,$(subst /, ,$(platform))) GOARCH=$(word 2,$(subst /, ,$(platform))) \
		go build -o $(BUILD_DIR)/$(APP_NAME)-$(subst /,-,$(platform))$(if $(filter windows,$(word 1,$(subst /, ,$(platform)))),$(WINDOWS_EXT),) $(MAIN_PATH) && \
		echo "Built $(BUILD_DIR)/$(APP_NAME)-$(subst /,-,$(platform))$(if $(filter windows,$(word 1,$(subst /, ,$(platform)))),$(WINDOWS_EXT),)" || exit 1; \
	)
	@echo "Cross-build complete!"

build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-386.exe $(MAIN_PATH)
	@echo "Windows builds complete: $(BUILD_DIR)/$(APP_NAME)-windows-{amd64,386}.exe"

build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-386 $(MAIN_PATH)
	@echo "Linux builds complete: $(BUILD_DIR)/$(APP_NAME)-linux-{amd64,386}"

build-mac:
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "macOS builds complete: $(BUILD_DIR)/$(APP_NAME)-darwin-{amd64,arm64}"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)/
	@rm -rf coverage/
	@go clean
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report: coverage/coverage.html"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

run:
	@go run $(MAIN_PATH)

help:
	@echo "Available commands:"
	@echo "  build         - Build for current platform"
	@echo "  cross-build   - Cross-compile for all platforms"
	@echo "  build-windows - Build for Windows"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-mac     - Build for macOS (Intel + Apple Silicon)"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage"
	@echo "  fmt           - Format code"
	@echo "  deps          - Install dependencies"
	@echo "  run           - Run application"
	@echo "  help          - Show this help"