#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

print_status "visuche installer starting..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_warning "Go is not installed."
    print_status "Installing Go..."
    
    # Try Homebrew first (macOS/Linux)
    if command -v brew &> /dev/null; then
        print_status "Using Homebrew to install Go..."
        brew install go
        print_success "Go installed via Homebrew!"
    # Try package managers on Linux
    elif command -v apt-get &> /dev/null; then
        print_status "Using apt to install Go..."
        sudo apt update
        sudo apt install golang-go -y
        print_success "Go installed via apt!"
    elif command -v yum &> /dev/null; then
        print_status "Using yum to install Go..."
        sudo dnf install golang -y
        print_success "Go installed via yum!"
    elif command -v pacman &> /dev/null; then
        print_status "Using pacman to install Go..."
        sudo pacman -S go --noconfirm
        print_success "Go installed via pacman!"
    # Fallback to manual binary download
    else
        print_warning "No package manager found. Downloading Go binary..."
        
        # Detect OS and architecture
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)
        case $ARCH in
            x86_64) ARCH="amd64" ;;
            aarch64|arm64) ARCH="arm64" ;;
            armv7l) ARCH="armv6l" ;;
            i386|i686) ARCH="386" ;;
            *) print_error "Unsupported architecture: $ARCH" ;;
        esac
        
        # Get latest Go version
        GO_VERSION=$(curl -s https://go.dev/VERSION?m=text | head -n1)
        DOWNLOAD_URL="https://go.dev/dl/${GO_VERSION}.${OS}-${ARCH}.tar.gz"
        
        print_status "Downloading Go ${GO_VERSION} from: $DOWNLOAD_URL"
        curl -L "$DOWNLOAD_URL" -o /tmp/go.tar.gz
        
        # Install to /usr/local (requires sudo) or ~/go
        if [ -w "/usr/local" ] || sudo -n true 2>/dev/null; then
            print_status "Installing Go to /usr/local..."
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf /tmp/go.tar.gz
            
            # Add to PATH
            if ! grep -q "/usr/local/go/bin" ~/.bashrc ~/.zshrc 2>/dev/null; then
                if [[ "$SHELL" == *"zsh"* ]]; then
                    echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.zshrc
                    print_success "Added Go to ~/.zshrc"
                else
                    echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc
                    print_success "Added Go to ~/.bashrc"
                fi
                export PATH="/usr/local/go/bin:$PATH"
            fi
        else
            print_status "Installing Go to ~/go..."
            mkdir -p ~/go
            tar -C ~/go -xzf /tmp/go.tar.gz --strip-components=1
            
            # Add to PATH
            if [[ "$SHELL" == *"zsh"* ]]; then
                echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
                print_success "Added Go to ~/.zshrc"
            else
                echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
                print_success "Added Go to ~/.bashrc"
            fi
            export PATH="$HOME/go/bin:$PATH"
        fi
        
        # Clean up
        rm -f /tmp/go.tar.gz
        
        print_success "Go installed successfully!"
        print_status "You may need to restart your shell or run: source ~/.bashrc (or ~/.zshrc)"
    fi
fi

print_success "Go found: $(go version)"

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    print_warning "GitHub CLI (gh) is not installed."
    print_status "Installing GitHub CLI..."
    
    # Try Homebrew first (macOS/Linux)
    if command -v brew &> /dev/null; then
        print_status "Using Homebrew to install GitHub CLI..."
        brew install gh
        print_success "GitHub CLI installed via Homebrew!"
    # Try package managers on Linux
    elif command -v apt-get &> /dev/null; then
        print_status "Using apt to install GitHub CLI..."
        curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
        sudo apt update
        sudo apt install gh -y
        print_success "GitHub CLI installed via apt!"
    elif command -v yum &> /dev/null; then
        print_status "Using yum to install GitHub CLI..."
        sudo dnf install 'dnf-command(config-manager)' -y
        sudo dnf config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo
        sudo dnf install gh -y
        print_success "GitHub CLI installed via yum!"
    elif command -v pacman &> /dev/null; then
        print_status "Using pacman to install GitHub CLI..."
        sudo pacman -S github-cli --noconfirm
        print_success "GitHub CLI installed via pacman!"
    # Fallback to manual binary download
    else
        print_warning "No package manager found. Downloading GitHub CLI binary..."
        
        # Detect OS and architecture
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)
        case $ARCH in
            x86_64) ARCH="amd64" ;;
            aarch64|arm64) ARCH="arm64" ;;
            armv7l) ARCH="armv6" ;;
            *) print_error "Unsupported architecture: $ARCH" ;;
        esac
        
        # Download and install
        GH_VERSION=$(curl -s https://api.github.com/repos/cli/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4 | sed 's/v//')
        DOWNLOAD_URL="https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_${OS}_${ARCH}.tar.gz"
        
        print_status "Downloading from: $DOWNLOAD_URL"
        curl -L "$DOWNLOAD_URL" -o /tmp/gh.tar.gz
        tar -xzf /tmp/gh.tar.gz -C /tmp
        
        # Install to ~/bin
        mkdir -p ~/bin
        cp "/tmp/gh_${GH_VERSION}_${OS}_${ARCH}/bin/gh" ~/bin/
        chmod +x ~/bin/gh
        
        # Clean up
        rm -rf /tmp/gh.tar.gz "/tmp/gh_${GH_VERSION}_${OS}_${ARCH}"
        
        print_success "GitHub CLI installed to ~/bin/gh!"
        print_status "Make sure ~/bin is in your PATH"
    fi
fi

print_success "GitHub CLI found: $(gh --version | head -n1)"

# Build the binary
print_status "Building visuche..."
go build -ldflags="-s -w" -o visuche

if [ ! -f "visuche" ]; then
    print_error "Build failed. Binary not found."
fi

print_success "Build completed successfully!"

# Create ~/bin directory if it doesn't exist
mkdir -p ~/bin

# Install the binary
print_status "Installing visuche to ~/bin..."
cp visuche ~/bin/
chmod +x ~/bin/visuche

print_success "visuche installed successfully!"

# Check if ~/bin is in PATH
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    print_warning "~/bin is not in your PATH"
    print_status "Adding ~/bin to your PATH..."
    
    # Detect shell and add to appropriate config file
    if [[ "$SHELL" == *"zsh"* ]]; then
        echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
        print_success "Added to ~/.zshrc"
        print_status "Run: source ~/.zshrc"
    elif [[ "$SHELL" == *"bash"* ]]; then
        echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
        print_success "Added to ~/.bashrc"
        print_status "Run: source ~/.bashrc"
    else
        print_warning "Unknown shell. Please manually add ~/bin to your PATH"
    fi
fi

# Test installation
print_status "Testing installation..."
if ~/bin/visuche --help > /dev/null 2>&1; then
    print_success "Installation test passed!"
else
    print_error "Installation test failed!"
fi

echo ""
print_success "🎉 visuche installation completed!"
echo ""
print_status "Usage:"
echo "  visuche                    # Interactive mode"
echo "  visuche --repo owner/repo  # Direct analysis"
echo "  visuche --help            # Show help"
echo ""
print_status "Don't forget to authenticate with GitHub:"
echo "  gh auth login"