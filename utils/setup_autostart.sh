#!/bin/bash
# Setup autostart entry for Ubuntu System Monitor

# Get the current directory (parent of utils)
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/.."
EXEC_PATH="$DIR/task_bar_monitor"
ICON_PATH="$DIR/assets/icons/app-icon.svg"
DESKTOP_FILE="$HOME/.config/autostart/task_bar_monitor.desktop"

# Parse command line arguments
ACTION="install"
if [ "$1" = "--remove" ] || [ "$1" = "-r" ]; then
    ACTION="remove"
fi

# Function to install autostart entry
install_autostart() {
    # Check if the executable exists
    if [ ! -f "$EXEC_PATH" ]; then
        echo "Error: Executable not found at $EXEC_PATH"
        echo "Please run 'go build -o task_bar_monitor ./cmd/main.go' first."
        exit 1
    fi

    # Check if the icon exists
    if [ ! -f "$ICON_PATH" ]; then
        echo "Warning: Custom icon not found at $ICON_PATH"
        echo "Using system icon instead."
        ICON="utilities-system-monitor"
    else
        ICON="$ICON_PATH"
    fi

    # Create autostart directory if it doesn't exist
    mkdir -p "$HOME/.config/autostart"

    # Create desktop entry file
    cat > "$DESKTOP_FILE" << EOL
[Desktop Entry]
Type=Application
Name=Ubuntu System Monitor
Exec=$EXEC_PATH
Icon=$ICON
Comment=System metrics in the taskbar
Categories=Utility;
Terminal=false
X-GNOME-Autostart-enabled=true
EOL

    echo "Auto-start entry created at $DESKTOP_FILE"
    echo "Ubuntu System Monitor will now start automatically when you log in."
}

# Function to remove autostart entry
remove_autostart() {
    if [ -f "$DESKTOP_FILE" ]; then
        rm "$DESKTOP_FILE"
        echo "Autostart entry removed."
    else
        echo "Autostart entry does not exist. Nothing to remove."
    fi
}

# Run the appropriate action
if [ "$ACTION" = "install" ]; then
    install_autostart
else
    remove_autostart
fi
