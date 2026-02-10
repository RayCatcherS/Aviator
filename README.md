# Aviator

**Local network application launcher and process monitor.**

Aviator is a desktop application that allows you to launch and monitor Windows applications remotely from any device on your local network. It exposes a secure internal HTTP server that serves a mobile-optimized dashboard.

![Version](https://img.shields.io/badge/version-2.8.2-blue)
![Stack](https://img.shields.io/badge/stack-Go_|_Vue_|_Wails-00ADD8)

## Documentation

-   [**User Manual**](docs/manual.html) - Installation and configuration guide.
-   [**Architecture**](docs/architecture.html) - System design and data flow.
-   [**API Reference**](docs/api.html) - REST API specification.

## Quick Start

### Download
Download the latest binary from [Releases](https://github.com/RayCatcherS/Aviator/releases).

### Build from Source
Requirements: Go 1.23+, Node.js 18+.

```bash
# Clone
git clone https://github.com/RayCatcherS/Aviator.git
cd Aviator/aviator-wails

# Install Tools
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build
wails build
```

## Features

-   **Remote Execution**: Launch executables via HTTP requests.
-   **Process Monitoring**: Real-time status tracking of configured apps.
-   **Security**: PIN-based authentication with HttpOnly cookies.
-   **Cross-Device**: Responsive web dashboard for mobile control.
-   **System Tray**: Background operation and quick server control.

## Tech Stack

-   **Backend**: Go (Wails bindings, `net/http` server).
-   **Frontend (Desktop)**: Vue 3.
-   **Frontend (Web)**: Vanilla JS + HTML (embedded).
-   **Storage**: JSON-based configuration (`%LOCALAPPDATA%\Aviator`).

## License
MIT License.
