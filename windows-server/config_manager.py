import json
import os
from typing import List, Dict, Optional, Callable
import uuid

CONFIG_FILE = "config.json"

class ConfigManager:
    def __init__(self):
        self.apps = []
        self.listeners: List[Callable] = []
        self.load()

    def add_listener(self, callback: Callable):
        self.listeners.append(callback)

    def _notify_listeners(self):
        for callback in self.listeners:
            try:
                callback()
            except Exception as e:
                print(f"Error in listener: {e}")

    def load(self):
        if os.path.exists(CONFIG_FILE):
            try:
                with open(CONFIG_FILE, 'r') as f:
                    self.apps = json.load(f)
            except Exception as e:
                print(f"Error loading config: {e}")
                self.apps = []
        else:
            self.apps = []

    def save(self):
        try:
            with open(CONFIG_FILE, 'w') as f:
                json.dump(self.apps, f, indent=4)
            self._notify_listeners()
        except Exception as e:
            print(f"Error saving config: {e}")

    def add_app(self, name: str, path: str, args: str = ""):
        app = {
            "id": str(uuid.uuid4()),
            "name": name,
            "path": path,
            "args": args
        }
        self.apps.append(app)
        self.save()
        return app

    def update_app_args(self, app_id: str, new_args: str):
        for app in self.apps:
            if app['id'] == app_id:
                app['args'] = new_args
                self.save()
                return True
        return False

    def remove_app(self, app_id: str):
        self.apps = [app for app in self.apps if app['id'] != app_id]
        self.save()

    def get_apps(self) -> List[Dict]:
        return self.apps

    def get_app_by_id(self, app_id: str) -> Optional[Dict]:
        for app in self.apps:
            if app['id'] == app_id:
                return app
        return None
