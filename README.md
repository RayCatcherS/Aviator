# WebApp Aviator

**WebApp Aviator** is a premium local network solution designed to bridge the gap between your mobile devices and your Windows desktop applications. Trigger games, efficiency tools, or scripts on your PC directly from your phone's browser with a beautiful, glassmorphism-inspired interface.

![Aviator Banner](docs/style.css) 
<!-- Note: Ideally we would have a screenshot here. Using a placeholder or omitting for now. -->

## ✨ Features

- **🚀 Launch Remotely**: Start applications on your Windows PC from any device on your local network.
- **🔌 Zero Config**: Automatic discovery via mDNS/Zeroconf. No need to type IP addresses—just scan a QR code.
- **🎨 Premium UI**: A modern, responsive web interface that looks native on iOS and Android.
- **🔧 Custom Arguments**: Define specific launch flags (e.g., `-bigpicture`, `--incognito`) for your applications.
- **⚡ Lightweight**: Pure Python backend with a Vanilla JS/HTML frontend. No Node.js required.

## 🚀 Quick Start

### Prerequisites
- Windows 10/11
- Python 3.8+ (Automatically checked by startup script)
- Mobile device on the **same Wi-Fi network** as your PC.

### Installation & Run
1.  **Clone/Download** this repository.
2.  Double-click the `start_aviator.bat` script in the root directory.
    -   *On first run, this will automatically create a virtual environment and install dependencies.*
3.  The **Aviator Server** window will open.

## 📖 User Guide

### 1. Configure Applications
-   In the Server window, click **Add Application...**.
-   Select the `.exe` file you want to be able to launch (e.g., Steam, Notepad, Custom Scripts).
-   (Optional) Select the app in the list and click **Edit Args** to add command-line arguments.

### 2. Connect Your Mobile Device
-   Look at the top-right corner of the Server window for a **QR Code**.
-   Open your phone's camera and scan the code.
-   Tap the link to open the **Aviator Web App**.
-   *Alternatively, open the URL displayed (e.g., `http://192.168.1.X:8000`) in your phone's browser.*

### 3. Launch!
-   Tap any icon on your phone to instantly launch the corresponding app on your PC.

## 🏗 Architecture

Aviator operates on a strictly local loop. No data leaves your network.

-   **Backend**: Python (FastAPI + Uvicorn) handles REST requests and process management.
-   **Frontend**: Pure HTML/JS with Tailwind CSS, served statically by FastAPI.
-   **Discovery**: `zeroconf` broadcasts `_aviator._tcp` services for easy finding.
-   **GUI**: `tkinter` provides a lightweight control panel on the host machine.

## 🛠 Technical Details

### Code Structure
-   `windows-server/main.py`: Entry point. Orchestrates the UI thread and the Server daemon thread.
-   `windows-server/api.py`: FastAPI application defining REST endpoints (`/api/apps`, `/api/launch/{id}`).
-   `windows-server/launcher.py`: Handles subprocess execution with `subprocess.Popen`.
-   `windows-server/discovery.py`: Manages mDNS broadcasting.
-   `windows-server/config_manager.py`: Persists app list to `config.json`.

### Development
To run from source without the batch script:
```bash
cd windows-server
python -m venv venv
.\venv\Scripts\activate
pip install -r requirements.txt
python main.py
```

## 📄 License
This project is open source.
