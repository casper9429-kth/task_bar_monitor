#!/bin/bash
set -e

echo "Building Task Bar Monitor Debian Package"
echo "---------------------------------------"

# Ensure the directory exists
mkdir -p debian_package/usr/bin

# Build the application
echo "Building the Go application..."
go build -o debian_package/usr/bin/task_bar_monitor ./cmd/main.go

# Make sure the scripts are executable
chmod +x debian_package/DEBIAN/prerm
chmod +x debian_package/DEBIAN/postinst
chmod +x debian_package/DEBIAN/postrm

# Build the package
echo "Building Debian package..."
dpkg-deb --build debian_package task-bar-monitor_1.0.0_amd64.deb

echo "Debian package built successfully: task-bar-monitor_1.0.0_amd64.deb"
echo "You can install it with: sudo dpkg -i task-bar-monitor_1.0.0_amd64.deb"
echo "After installation, you may need to run: sudo apt-get install -f"
echo "To uninstall: sudo apt remove task-bar-monitor"