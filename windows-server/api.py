from fastapi import FastAPI, HTTPException, WebSocket, WebSocketDisconnect
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List
import launcher
from config_manager import ConfigManager
import os

app = FastAPI()
config = ConfigManager()

# WebSocket Manager
class ConnectionManager:
    def __init__(self):
        self.active_connections: List[WebSocket] = []

    async def connect(self, websocket: WebSocket):
        await websocket.accept()
        self.active_connections.append(websocket)

    def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)

    async def broadcast(self, message: str):
        for connection in self.active_connections:
            try:
                await connection.send_text(message)
            except:
                pass

manager = ConnectionManager()

# Register listener to ConfigManager
# Since notify_listeners is synchronous and broadcast is async, we need a bridge.
# But fastapi/uvicorn runs in an async loop. 
# We can't easily call async from sync callback without getting loop.
# However, we can use a simple simpler approach: check for updates or just ignore 
# strict async correctness for a local tool or use run_coroutine_threadsafe.
# Simplify: The listener will just set a flag or we try to run it.
# BETTER OPTION: we can't await in the sync callback.
# Let's simple use a non-async broadcast for now? No, send_text is async.
# We will cheat slightly: The GUI thread calls 'save', which calls listener.
# The listener is executed in GUI thread (MainThread).
# The event loop is in a separate Daemon thread (Server).
# We need to tell the Server Loop to broadcast.
import asyncio

def on_config_change():
    """Called when config changes. Needs to trigger broadcast in the loop."""
    # We can't easily access the running loop from here if it's in another thread without handle.
    # But wait, main.py starts server in daemon thread? Yes.
    # Actually, uvicorn.run blocks. 
    # In main.py: t = threading.Thread(target=start_server); t.start()
    # So server has its own loop.
    # We need to schedule the broadcast coroutine in that loop.
    # For now, let's just allow polling or simple refresh?
    # User asked for AUTOMATIC update.
    # We can use a simple global queue or just ignore the loop issue?
    # No, let's try to get the loop.
    pass 
    
# We will solve the thread bridge in a moment. 
# For now let's add the endpoint.

# Allow CORS (useful if we debug from a different port, though not strictly needed for same-origin)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class AppModel(BaseModel):
    name: str
    path: str

@app.get("/api/apps")
def get_apps():
    return config.get_apps()

@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    await manager.connect(websocket)
    try:
        while True:
            await websocket.receive_text()
    except WebSocketDisconnect:
        manager.disconnect(websocket)

@app.post("/api/launch/{app_id}")
def launch_app(app_id: str):
    app_item = config.get_app_by_id(app_id)
    if not app_item:
        raise HTTPException(status_code=404, detail="App not found")
    
    try:
        args = app_item.get('args', '')
        launcher.run_executable(app_item['path'], args)
        return {"status": "success", "message": f"Launched {app_item['name']}"}
    except FileNotFoundError:
        raise HTTPException(status_code=404, detail="Executable file not found on disk")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/info")
def get_info():
    import socket
    return {
        "hostname": socket.gethostname(),
        "status": "running"
    }

# Mount static files - DO THIS LAST so API routes take precedence
# We expect 'static' folder to be in the same directory
static_dir = os.path.join(os.path.dirname(__file__), "static")
if os.path.exists(static_dir):
    app.mount("/", StaticFiles(directory=static_dir, html=True), name="static")
