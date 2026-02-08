# Aviator ✈️

> **Modern desktop application launcher with glassmorphism UI and network control**

Launch your Windows applications remotely from any device on your local network with a beautiful, modern interface powered by Wails.

![Version](https://img.shields.io/badge/version-2.0.0-blue)
![Wails](https://img.shields.io/badge/Wails-v2.11.0-00ACD7)
![Vue](https://img.shields.io/badge/Vue-3-42b883)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8)

## ✨ Features

- 🎨 **Glassmorphism UI** - Modern, premium interface with blur effects and smooth animations
- 📱 **Mobile Access** - Control your PC applications from your phone or tablet
- 🔍 **Auto-Discovery** - Zero-configuration setup with mDNS (scan QR code and go)
- ⚡ **Fast & Lightweight** - Native performance with small footprint (~12 MB)
- 🌐 **Cross-Platform Ready** - Built with Wails (Windows, macOS, Linux support)
- 🎯 **Smart Launcher** - Launch apps with custom arguments and working directories
- 📡 **Real-time Monitoring** - View running/stopped status of apps on both Desktop and Web Dashboard

## 📸 Screenshots

*Coming soon - Glassmorphism interface with gradient backgrounds, blur effects, and smooth animations*

## 🚀 Quick Start

### Prerequisites

- Windows 10/11 (WebView2 included)
- For development: Go 1.23+ and Node.js 18+

### Installation

**Option 1: Download Release** (Recommended)
1. Download `aviator-wails.exe` from [Releases](https://github.com/RayCatcherS/Aviator/releases)
2. Run the executable
3. Add your favorite applications
4. Scan the QR code from your phone!

**Option 2: Build from Source**
```bash
# Clone the repository
git clone https://github.com/RayCatcherS/Aviator.git
cd Aviator/aviator-wails

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build the application
wails build

# Run
./build/bin/aviator-wails.exe
```

## 💡 Usage

### Desktop Interface

1. **Launch Aviator** - Double-click `aviator-wails.exe`
2. **Add Applications** - Click "Add App" and select your .exe files
3. **Configure** - Optionally add launch arguments
4. **Monitor** - Check the status LED next to each app (Green = Running, Grey = Stopped)
5. **Access** - Use the displayed URLs or scan the QR code

### Mobile/Web Access

1. **Same Network** - Ensure your phone/device is on the same WiFi
2. **Scan QR Code** - Use your phone's camera to scan the QR code on the desktop app
3. **Launch & Monitor** - Tap any app card to launch it and view its real-time status
4. **Auto-Discovery** - The app announces itself via mDNS

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│         Frontend (Vue 3 + CSS)          │
│     Glassmorphism UI Components         │
└──────────────┬──────────────────────────┘
               │ Wails Bridge
┌──────────────┴──────────────────────────┐
│           Backend (Go)                   │
├──────────────────────────────────────────┤
│  • Config Manager (JSON persistence)    │
│  • HTTP Server (REST API + Static)      │
│  • Process Launcher (Detached exec)     │
│  • mDNS Discovery (Network announce)    │
└──────────────────────────────────────────┘
```

### Technology Stack

- **Frontend**: Vue 3, Glassmorphism CSS, QRCode.js
- **Backend**: Go 1.23, Wails v2
- **Networking**: HTTP server, mDNS/Zeroconf
- **Storage**: JSON configuration in `%LOCALAPPDATA%\Aviator\`

## 📚 Documentation

- **[Overview](docs/index.html)** - General introduction and features
- **[User Guide](docs/user_guide.html)** - How to use Aviator
- **[Architecture](docs/architecture.html)** - Technical overview
- **[API Reference](docs/technical.html)** - REST API documentation

> **Note:** To view these docs as a website, enable **GitHub Pages** in your repository settings (Source: `/docs` folder). Otherwise, you will see the HTML source code.

## 🛠️ Development

### Prerequisites

```bash
# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify installation
wails doctor
```

### Development Mode

```bash
cd aviator-wails
wails dev
```

This starts:
- Backend Go server with hot reload
- Frontend Vite dev server with HMR
- Application window with developer tools

### Project Structure

```
aviator-wails/
├── app.go              # Wails bindings (Go ↔ JS)
├── main.go             # Entry point
├── internal/           # Backend modules
│   ├── config/         # Configuration management
│   ├── server/         # HTTP server + API
│   ├── launcher/       # Process execution
│   ├── discovery/      # mDNS service
│   └── web/            # Embedded static files
├── frontend/           # Vue 3 application
│   ├── src/
│   │   ├── App.vue     # Main component
│   │   └── style.css   # Glassmorphism design system
│   └── wailsjs/        # Auto-generated bindings
└── build/              # Build output
    └── bin/
        └── aviator-wails.exe
```

### Building

```bash
# Production build
wails build

# Clean build
wails build -clean

# With custom flags
wails build -ldflags "-w -s"
```

⚠️ **Important**: Never use `go build` directly. Always use `wails build` to ensure proper build tags.

## 🔧 Configuration

Configuration is stored in: `%LOCALAPPDATA%\Aviator\config.json`

```json
[
  {
    "id": "unique-uuid",
    "name": "My Application",
    "path": "C:\\path\\to\\app.exe",
    "args": "--flag value"
  }
]
```

## 🌐 API Endpoints

When the server is running on port 8000:

- `GET /` - Web interface (glassmorphism UI)
- `GET /api/apps` - List configured applications
- `POST /api/launch/{id}` - Launch an application
- `GET /api/monitoring/status` - Get real-time process statuses
- `GET /api/info` - Server information

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Wails](https://wails.io) - Amazing Go + Web framework
- [Vue.js](https://vuejs.org) - Progressive JavaScript framework
- [Glassmorphism](https://uxdesign.cc/glassmorphism-in-user-interfaces-1f39bb1308c9) - Design trend inspiration
- [Zeroconf](https://github.com/grandcat/zeroconf) - mDNS library

## 🐛 Known Issues

- **WebView2 Required**: Windows 10/11 include WebView2 by default. Older systems may need manual installation.
- **Firewall**: Windows Firewall may prompt for network access on first run.
- **Arguments Parsing**: Complex arguments with quotes may need escaping.

## 🗺️ Roadmap

- [ ] System tray integration
- [ ] Dark/Light mode toggle
- [ ] Custom app icons
- [ ] Application usage statistics
- [ ] Auto-start on Windows boot
- [ ] macOS and Linux builds
- [ ] Mobile companion app

## 📧 Contact

- **Author**: RayCatcherS
- **Project Link**: https://github.com/RayCatcherS/Aviator

---

<p align="center">Made with ❤️ and Glassmorphism</p>
