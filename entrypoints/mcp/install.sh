#!/bin/bash
set -e

echo "Installing Classeviva MCP Server..."

# Build the binary
echo "Building..."
go build -o classeviva-mcp ./entrypoints/mcp/main.go

# Move binary to a common location
if [ -w "/usr/local/bin" ]; then
    echo "Installing to /usr/local/bin/classeviva-mcp..."
    mv classeviva-mcp /usr/local/bin/
    BINARY_PATH="/usr/local/bin/classeviva-mcp"
else
    echo "Installing to $HOME/.local/bin/classeviva-mcp..."
    mkdir -p "$HOME/.local/bin"
    mv classeviva-mcp "$HOME/.local/bin/"
    BINARY_PATH="$HOME/.local/bin/classeviva-mcp"
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo "Adding $HOME/.local/bin to PATH..."
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
        echo "Note: Restart your shell or run: source ~/.bashrc"
    fi
fi

# Create config directory if it doesn't exist
mkdir -p "$HOME/.classeviva"

# Copy example config if it doesn't exist
if [ ! -f "$HOME/.classeviva/students.json" ]; then
    echo "Creating example configuration..."
    cp entrypoints/mcp/students.json.example "$HOME/.classeviva/students.json"
    chmod 600 "$HOME/.classeviva/students.json"
    echo ""
    echo "Configuration file created at: $HOME/.classeviva/students.json"
    echo "Please edit this file with your Classeviva credentials."
else
    echo "Configuration file already exists at: $HOME/.classeviva/students.json"
fi

echo ""
echo "Installation complete!"
echo "Binary installed at: $BINARY_PATH"
echo ""
echo "To use with Claude Desktop, add this to your config:"
echo ""
echo '{
  "mcpServers": {
    "classeviva": {
      "command": "'$BINARY_PATH'"
    }
  }
}'
echo ""
echo "Configuration file location:"
echo "  macOS:   ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "  Linux:   ~/.config/Claude/claude_desktop_config.json"
echo "  Windows: %APPDATA%\\Claude\\claude_desktop_config.json"
