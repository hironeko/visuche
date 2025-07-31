.PHONY: build install uninstall clean help

# Default target
help:
	@echo "visuche - GitHub Repository Analytics Tool"
	@echo ""
	@echo "Available commands:"
	@echo "  make build     - Build the binary"
	@echo "  make install   - Build and install to ~/bin"
	@echo "  make uninstall - Remove from ~/bin"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make help      - Show this help"

# Build the binary
build:
	@echo "🔨 Building visuche..."
	go build -ldflags="-s -w" -o visuche

# Install to ~/bin
install: build
	@echo "📦 Installing visuche to ~/bin..."
	@mkdir -p ~/bin
	@cp visuche ~/bin/
	@echo "✅ visuche installed successfully!"
	@echo "💡 Make sure ~/bin is in your PATH:"
	@echo "   export PATH=\"\$$HOME/bin:\$$PATH\""
	@echo ""
	@echo "🎯 You can now run: visuche"

# Uninstall from ~/bin
uninstall:
	@echo "🗑️  Removing visuche from ~/bin..."
	@rm -f ~/bin/visuche
	@echo "✅ visuche uninstalled successfully!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f visuche
	@rm -f *.csv
	@echo "✅ Clean complete!"

# Development build with verbose output
dev-build:
	@echo "🔨 Building visuche (development mode)..."
	go build -v -o visuche