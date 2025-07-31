# Build and Installation Instructions

## Prerequisites

- Go 1.19 or later
- GitHub CLI (`gh`) - for API access
- Git (for repository detection)

## Quick Install (Recommended)

### Option 1: Automated Install Script
```bash
chmod +x install.sh
./install.sh
```

This script will:
- Check for Go and GitHub CLI
- Install GitHub CLI via Homebrew if missing
- Build the binary with optimizations
- Install to `~/bin/`
- Configure your PATH automatically

### Option 2: Manual Build with Makefile
```bash
# Build only
make build

# Build and install
make install

# Show help
make help
```

## Manual Installation

If you prefer manual installation:

```bash
# 1. Build the binary
go build -ldflags="-s -w" -o visuche

# 2. Move to your PATH
mkdir -p ~/bin
cp visuche ~/bin/
chmod +x ~/bin/visuche

# 3. Add ~/bin to PATH (if not already)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc  # or ~/.bashrc
```

## Usage

After installation:

```bash
# Interactive mode (recommended for first use)
visuche

# Direct analysis
visuche --repo owner/repo

# With time period
visuche --repo owner/repo --since 2024-01-01 --until 2024-01-31

# Export to CSV
visuche --repo owner/repo --csv

# Show help
visuche --help
```

## GitHub Authentication

Before first use, authenticate with GitHub:

```bash
gh auth login
```

Choose "GitHub.com" and follow the prompts to authenticate via web browser.

## Uninstall

```bash
# Using Makefile
make uninstall

# Manual removal
rm ~/bin/visuche
```

## Development

For development builds:

```bash
# Development build with verbose output
make dev-build

# Clean build artifacts
make clean
```

## Troubleshooting

**Command not found**: Ensure `~/bin` is in your PATH
**Authentication errors**: Run `gh auth login` to authenticate with GitHub
**Build errors**: Ensure Go 1.19+ is installed and GOPATH is configured