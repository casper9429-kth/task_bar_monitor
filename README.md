# Ubuntu System Monitor

<img src="assets/icons/app-icon.svg" align="right" width="128">

Ubuntu System Monitor is a lightweight application developed in Go that displays real-time system metrics including CPU usage, memory usage, network usage, and disk usage directly in the Ubuntu taskbar. The application is designed to be user-friendly and customizable.

## Features

- Real-time monitoring of:
  - CPU Usage
  - Memory Usage
  - Network Usage
  - Disk Usage
- Customizable settings:
  - Choose which metrics to display
  - Configure what appears in the taskbar
  - Set refresh interval
- System tray integration for easy access and visibility

## Requirements

- Go 1.16 or later
- Ubuntu 22.04 or later (or another Linux distribution with system tray support)
- Required system dependencies:
  - libayatana-appindicator3-dev
  - pkg-config
  - build-essential

## Installation

1. Install required system dependencies:
   ```bash
   # For Ubuntu/Debian
   sudo apt update
   sudo apt install -y libayatana-appindicator3-dev pkg-config build-essential

   # For Fedora
   sudo dnf install -y libayatana-appindicator3-devel pkg-config gcc
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ubuntu-system-monitor.git
   cd ubuntu-system-monitor
   ```

3. Build the application:
   ```bash
   go build -o ubuntu-system-monitor ./cmd/main.go
   ```

4. Run the application:
   ```bash
   ./ubuntu-system-monitor
   ```

## Usage

- After launching the application, you will see an icon in the system tray.
- Click on the icon to view the current system metrics.
- Right-click the icon and select "Settings" to customize your preferences.
- To have the application start automatically at login, use the provided autostart script:
  ```bash
  ./utils/setup_autostart.sh
  ```
- Remove the autostart entry with:
  ```bash
  ./utils/setup_autostart.sh --remove
  ```

### Command Line Options

The application supports the following command line options:

- `-debug`: Enable debug logging

Example:
```bash
./ubuntu-system-monitor -debug
```

## Troubleshooting

If you encounter any issues, try running the application with debugging enabled:

```bash
./ubuntu-system-monitor -debug
```

For more detailed troubleshooting, check the debugging.md file.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.