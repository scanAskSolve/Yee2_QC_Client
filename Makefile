# 顏色定義
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
PURPLE = \033[0;35m
CYAN = \033[0;36m
WHITE = \033[0;37m
NC = \033[0m # No Color


SRC_DEVICE=./device/client.go
SRC_SERVICE=./service

# Variables
APP_NAME = yee2_qc
BUILD_DIR = build

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS = -ldflags "\
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.GitBranch=$(GIT_BRANCH)' \
	-s -w"

# Default targets
.PHONY: all clean linux windows darwin version info help

all: linux-arm64 linux windows darwin

# Show version information
version:
	@echo -e "$(CYAN)Version Information:$(NC)"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Git Branch: $(GIT_BRANCH)"

# Show build information
info:
	@echo -e "$(BLUE)Build Information:$(NC)"
	@echo "  App Name: $(APP_NAME)"
	@echo "  Build Dir: $(BUILD_DIR)"
	@echo "  LDFLAGS: $(LDFLAGS)"

# Show help
help:
	@echo -e "$(PURPLE)Available Commands:$(NC)"
	@echo ""
	@echo "  make all          - Build for all platforms"
	@echo "  make linux        - Build for Linux"
	@echo "  make windows      - Build for Windows"
	@echo "  make darwin       - Build for macOS"
	@echo "  make linux-arm64  - Build for Linux ARM64"
	@echo "  make clean        - Clean build files"
	@echo "  make deps         - Install dependencies"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make build        - Complete build process"
	@echo "  make quick        - Quick build (Linux only)"
	@echo "  make version      - Show version info"
	@echo "  make info         - Show build info"
	@echo "  make help         - Show this help"
	@echo ""

# Linux build
linux:
	@echo -e "$(YELLOW)Building Linux version ($(VERSION))...$(NC)"
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(APP_NAME)_device $(SRC_DEVICE)
	#GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(APP_NAME)_service $(SRC_SERVICE)
	@echo -e "$(GREEN)Linux build completed: $(BUILD_DIR)/linux/$(APP_NAME)$(NC)"

# Windows build
windows:
	@echo -e "$(YELLOW)Building Windows version ($(VERSION))...$(NC)"
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(APP_NAME)_device.exe $(SRC_DEVICE)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(APP_NAME)_service.exe $(SRC_SERVICE)
	@echo -e "$(GREEN)Windows build completed: $(BUILD_DIR)/windows/$(APP_NAME).exe$(NC)"

# macOS build
darwin:
	@echo -e "$(YELLOW)Building macOS version ($(VERSION))...$(NC)"
	@mkdir -p $(BUILD_DIR)/darwin
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(APP_NAME)_device $(SRC_DEVICE)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(APP_NAME)_service $(SRC_SERVICE)
	@echo -e "$(GREEN)macOS build completed: $(BUILD_DIR)/darwin/$(APP_NAME)$(NC)"

# ARM build
linux-arm64:
	@echo -e "$(YELLOW)Building Linux ARM64 version ($(VERSION))...$(NC)"
	@mkdir -p $(BUILD_DIR)/linux-arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-arm64/$(APP_NAME)_device $(SRC_DEVICE)
	#GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/linux-arm64/$(APP_NAME)_service $(SRC_SERVICE)
	@echo -e "$(GREEN)Linux ARM64 build completed: $(BUILD_DIR)/linux-arm64/$(APP_NAME)$(NC)"

# Clean
clean:
	@echo -e "$(RED)Cleaning build files...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo -e "$(GREEN)Clean completed$(NC)"

# Install dependencies
deps:
	@echo -e "$(BLUE)Installing dependencies...$(NC)"
	go mod tidy
	go mod download
	@echo -e "$(GREEN)Dependencies installed$(NC)"

# Run tests
test:
	@echo -e "$(BLUE)Running tests...$(NC)"
	go test -v ./...
	@echo -e "$(GREEN)Tests completed$(NC)"

# Format code
fmt:
	@echo -e "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@echo -e "$(GREEN)Code formatted$(NC)"

# Lint code
lint:
	@echo -e "$(BLUE)Linting code...$(NC)"
	go vet ./...
	@echo -e "$(GREEN)Code linted$(NC)"

# Complete build process
build: clean deps fmt lint test all
	@echo -e "$(GREEN)Complete build process finished!$(NC)"

# Quick build (Linux only)
quick: linux
	@echo -e "$(GREEN)Quick build completed!$(NC)"
