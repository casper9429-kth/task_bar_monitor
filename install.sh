#!/bin/bash

echo "Ubuntu System Monitor - Installation Script"
echo "----------------------------------------"

# Detect the package manager
if command -v apt-get &> /dev/null; then
    echo "Detected apt package manager (Debian/Ubuntu)"
    sudo apt-get update
    
    # Install system dependencies
    echo "Installing system dependencies..."
    sudo apt-get install -y \
        libayatana-appindicator3-dev \
        pkg-config \
        build-essential
    
elif command -v dnf &> /dev/null; then
    echo "Detected dnf package manager (Fedora)"
    
    # Install system dependencies
    echo "Installing system dependencies..."
    sudo dnf install -y \
        libayatana-appindicator3-devel \
        pkg-config \
        gcc
else
    echo "Unsupported package manager. Please install these dependencies manually:"
    echo "- For Debian/Ubuntu: libayatana-appindicator3-dev, pkg-config, build-essential"
    echo "- For Fedora: libayatana-appindicator3-devel, pkg-config, gcc"
    exit 1
fi

# Build the application
echo "Building the application..."
go build -o task_bar_monitor ./cmd/main.go

echo "Installation complete! You can run the application with:"
echo "./task_bar_monitor"
echo
echo "To set up autostart on login, run:"
echo "./utils/setup_autostart.sh"