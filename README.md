# Task Bar Monitor

<img src="assets/icons/app-icon.svg" align="right" width="128">

Task Bar Monitor is a lightweight application developed in Go that displays real-time system metrics including CPU usage, memory usage, network usage, and disk usage directly in the Ubuntu taskbar. The application is designed to be user-friendly and customizable.

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

### Using Debian Package (Recommended)
1. Download the latest Debian package or build it yourself:
   ```bash
   ./utils/build_deb_package.sh
   ```

2. Install the package:
   ```bash
   sudo dpkg -i task-bar-monitor_1.0.0_amd64.deb
   ```

3. If there are dependency issues, resolve them with:
   ```bash
   sudo apt-get install -f
   ```

The application will be automatically installed and configured to start when you log in.

### Uninstalling
To remove the application:
```bash
sudo apt remove task-bar-monitor
```

To completely remove the application including any configuration files:
```bash
sudo apt purge task-bar-monitor
```

### Building from Source

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
   git clone https://github.com/casper9429-kth/task_bar_monitor.git
   cd task_bar_monitor
   ```

3. Build the application:
   You can build the application in two ways:

   a. Using the install script (recommended):
   ```bash
   ./utils/install.sh
   ```
   
   b. Manual compilation:
   ```bash
   go build -o task_bar_monitor ./cmd/main.go
   ```

4. Run the application:
   ```bash
   ./task_bar_monitor
   ```

5. Optional: Set up autostart to run at login:
   ```bash
   ./utils/setup_autostart.sh
   ```

   To remove the autostart entry:
   ```bash
    ./utils/setup_autostart.sh --remove
    ```

6. For developers: If you want to build a Debian package from source:
   ```bash
   ./utils/build_deb_package.sh
   ```
   This will create a `.deb` file that can be installed with `sudo dpkg -i task-bar-monitor_1.0.0_amd64.deb`

## Usage

- After launching the application, you will see an icon `...` in the system tray.
- Click on the icon to view the current system metrics.
- Right-click the icon and select "Settings" to customize your preferences.

### Command Line Options

The application supports the following command line options:

- `-debug`: Enable debug logging

Example:
```bash
./task_bar_monitor -debug
```

## Troubleshooting

If you encounter any issues, try running the application with debugging enabled:

```bash
./task_bar_monitor -debug
```

For more detailed troubleshooting, check the debugging.md file.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License.