#!/bin/bash

# 1. Define installation directory
DEFAULT_DIR="$HOME/.fast-wireguard"

echo "Default installation directory is: $DEFAULT_DIR"
read -p "Press Enter to continue, or enter a new path: " USER_DIR
if [ -z "$USER_DIR" ]; then
    INSTALL_DIR="$DEFAULT_DIR"
else
    INSTALL_DIR="$USER_DIR"
fi

# 2. Create installation directory and copy files
# Note: $PWD refers to the temporary directory after extraction
echo "Installing to $INSTALL_DIR ..."
mkdir -p "$INSTALL_DIR"
cp -r * "$INSTALL_DIR/"

# 3. Configure environment variables (this is the most distinctive step of Miniconda)
echo "Configuring environment variables..."
SHELL_CONFIG="$HOME/.bashrc"
# Detect if using zsh
if [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
fi

# Write PATH configuration
EXPORT_CMD="export PATH=\"$INSTALL_DIR/bin:\$PATH\""

# Check if it already exists to avoid duplicate entries
if grep -q "$INSTALL_DIR/bin" "$SHELL_CONFIG"; then
    echo "Environment variable already exists, skipping."
else
    echo "" >> "$SHELL_CONFIG"
    echo "# Fast-WireGuard configuration" >> "$SHELL_CONFIG"
    echo "$EXPORT_CMD" >> "$SHELL_CONFIG"
    echo "Environment variable added to $SHELL_CONFIG"
fi

echo "=== Installation complete! ==="
echo "Please reopen your terminal, or run 'source $SHELL_CONFIG' to apply the changes."
