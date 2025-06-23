# Git Auto-Push Tool Makefile

# Variables
BINARY_NAME=git-autopush
MAIN_FILE=main.go
INSTALL_DIR=$(HOME)/bin

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed successfully!"

# Install the application
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation completed!"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

# Create alias
.PHONY: alias
alias:
	@echo "Creating alias 'gp' for $(BINARY_NAME)..."
	@echo "alias gp='$(BINARY_NAME)'" >> $(HOME)/.bashrc || true
	@echo "alias gp='$(BINARY_NAME)'" >> $(HOME)/.zshrc || true
	@echo "Alias created! Restart your terminal or source your shell config."

# Run the application
.PHONY: run
run: build
	@./$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@echo "Clean completed!"

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  install  - Install the application to ~/bin"
	@echo "  alias    - Create 'gp' alias for easy usage"
	@echo "  run      - Build and run the application"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  help     - Show this help message"

# Initialize Go module (run once)
.PHONY: init
init:
	@echo "Initializing Go module..."
	@go mod init git-autopush
	@echo "Go module initialized!"


