import threading
import uvicorn
import time
import os
import sys
import tkinter as tk
from tkinter import messagebox
from api import app
from discovery import DiscoveryService
from gui import start_gui
from api import app, config as shared_config

PORT = 8000
discovery_service = None

def start_server():
    """Function to run the uvicorn server"""
    try:
        # Disable duplicate logs
        uvicorn.run(app, host="0.0.0.0", port=PORT, log_level="info")
    except Exception as e:
        print(f"[SERVER ERROR] {e}")

def main():
    global discovery_service
    
    print("Initializing components...")
    
    # 1. Start Discovery Service
    try:
        discovery_service = DiscoveryService(PORT)
        discovery_service.register()
    except Exception as e:
        print(f"[WARNING] Failed to start discovery service (mDNS): {e}")

    # 2. Start FastAPI Server in a separate thread
    print(f"Starting Web Server on port {PORT}...")
    server_thread = threading.Thread(target=start_server, daemon=True)
    server_thread.start()

    # 3. Start GUI (Must be on main thread for tkinter)
    print("Starting GUI...")
    try:
        # We use the shared config from API so updates are broadcasted
        config = shared_config
        
        # Get IP for console display
        import socket
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        try:
            s.connect(('10.255.255.255', 1))
            IP = s.getsockname()[0]
        except:
            IP = '127.0.0.1'
        finally:
            s.close()

        print("\n" + "="*50)
        print(" SERVER READY! Access via:")
        print(f" -> Local PC:    http://localhost:8000")
        print(f" -> Network/App: http://{IP}:8000")
        print("="*50 + "\n")

        start_gui(config, stop_callback=stop)
    except Exception as e:
        print(f"[CRITICAL ERROR] GUI failed to start: {e}")
        input("Press Enter to exit...")

def stop():
    """Callback to cleanup when GUI closes"""
    print("Stopping application...")
    if discovery_service:
        try:
            discovery_service.unregister()
        except:
            pass
    # Sys exit will kill the daemon server thread
    print("Goodbye.")
    os._exit(0)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        stop()
    except Exception as e:
        print(f"[CRITICAL UNHANDLED ERROR] {e}")
        import traceback
        traceback.print_exc()
        input("Press Enter to exit...")
