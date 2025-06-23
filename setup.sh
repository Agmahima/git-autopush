#!/bin/bash

# Git Auto-Push Installation and Usage Script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        echo "Visit: https://golang.org/doc/install"
        exit 1
    fi
    print_success "Go is installed: $(go version)"
}

# Build the Go application
build_app() {
    print_info "Building Git Auto-Push tool..."
    
    if go build -o git-autopush main.go; then
        print_success "Build completed successfully!"
    else
        print_error "Build failed!"
        exit 1
    fi
}

# Install the binary to system PATH
install_binary() {
    print_info "Installing binary to system PATH..."
    
    # Create ~/bin if it doesn't exist
    mkdir -p ~/bin
    
    # Copy binary to ~/bin
    cp git-autopush ~/bin/
    
    # Make it executable
    chmod +x ~/bin/git-autopush
    
    # Add ~/bin to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
        echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
        echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
        print_warning "Added ~/bin to PATH. Please restart your terminal or run:"
        echo "export PATH=\"\$HOME/bin:\$PATH\""
    fi
    
    print_success "Installation completed!"
}

# Create alias for easier usage
create_alias() {
    print_info "Creating alias 'gp' for git-autopush..."
    
    # Add alias to shell config files
    echo "alias gp='git-autopush'" >> ~/.bashrc
    echo "alias gp='git-autopush'" >> ~/.zshrc
    
    print_success "Alias 'gp' created! Use 'gp' to run git auto-push."
    print_warning "Please restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"
}

# Show usage instructions
show_usage() {
    echo ""
    print_info "Usage Instructions:"
    echo "1. Navigate to your project directory"
    echo "2. Run: git-autopush (or 'gp' if using alias)"
    echo "3. Follow the prompts"
    echo ""
    echo "Features:"
    echo "• First-time setup: Initializes git, adds remote, commits, and pushes"
    echo "• Subsequent pushes: Asks for confirmation, commits, and pushes automatically"
    echo "• Smart branch detection and creation"
    echo "• Interactive prompts for user control"
    echo ""
}

# Main function
main() {
    echo "========================================="
    echo "    Git Auto-Push Tool Setup Script"
    echo "========================================="
    echo ""
    
    case "${1:-install}" in
        "build")
            check_go
            build_app
            ;;
        "install")
            check_go
            build_app
            install_binary
            create_alias
            show_usage
            ;;
        "usage")
            show_usage
            ;;
        *)
            print_error "Unknown command: $1"
            echo "Usage: $0 [build|install|usage]"
            exit 1
            ;;
    esac
}

main "$@"

