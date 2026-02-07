<template>
  <div id="app">
    <!-- Status Header Panel -->
    <div class="glass-panel fade-in" style="flex-shrink: 0;">
      <div class="flex justify-between items-center">
        <div class="flex flex-col gap-sm">
          <h1 style="font-size: 28px; font-weight: 700; margin: 0;">Aviator</h1>
          <div class="status-badge" :class="serverInfo.status">
            {{ serverInfo.status.toUpperCase() }}
          </div>
          
          <!-- Server Control Buttons -->
          <div class="flex gap-sm" style="margin-top: 12px;">
            <button 
              @click="startServer" 
              :disabled="serverInfo.running"
              class="glass-button success"
              :class="{ 'opacity-50 cursor-not-allowed': serverInfo.running }"
            >
              ▶️ Start Server
            </button>
            <button 
              @click="stopServer"
              :disabled="!serverInfo.running"
              class="glass-button"
              :class="{ 'opacity-50 cursor-not-allowed': !serverInfo.running }"
             style="background: rgba(239, 68, 68, 0.2); border-color: rgba(239, 68, 68, 0.3);"
            >
              ⏹️ Stop Server
            </button>
          </div>
          
          <div v-if="serverInfo.running" class="flex flex-col gap-sm" style="margin-top: 12px;">
            <div style="font-size: 13px; opacity: 0.9;">
              <strong>Local:</strong> <a :href="serverInfo.localURL" @click.prevent="openURL(serverInfo.localURL)">{{ serverInfo.localURL }}</a>
            </div>
            <div style="font-size: 13px; opacity: 0.9;">
              <strong>Network:</strong> <a :href="serverInfo.networkURL" @click.prevent="openURL(serverInfo.networkURL)">{{ serverInfo.networkURL }}</a>
            </div>
          </div>
          <div v-else style="margin-top: 12px; font-size: 13px; opacity: 0.7;">
            Server stopped. Click "Start Server" to enable web access.
          </div>
        </div>
        
        <!-- QR Code -->
        <div v-if="serverInfo.running" class="qr-container glass-card" style="padding: 16px; cursor: default;">
          <canvas ref="qrCanvas" width="140" height="140"></canvas>
        </div>
        <div v-else class="qr-container glass-card" style="padding: 16px; width: 172px; height: 172px; display: flex; align-items: center; justify-content: center; opacity: 0.3;">
          <div style="text-align: center;">
            <div style="font-size: 48px;">🔒</div>
            <div style="font-size: 11px; margin-top: 8px;">Server Offline</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Applications Section -->
    <div class="glass-panel" style="flex: 1; display: flex; flex-direction: column; overflow: hidden;">
      <div class="flex justify-between items-center" style="margin-bottom: 16px;">
        <h2 style="font-size: 20px; font-weight: 600; margin: 0;">Applications</h2>
        <button @click="openAddDialog" class="glass-button primary">
          + Add App
        </button>
      </div>

      <!-- Apps Grid -->
      <div class="app-grid" style="flex: 1; overflow-y: auto;">
        <div v-for="app in apps" :key="app.id" class="app-card glass-card">
          <div class="flex justify-between items-center" style="margin-bottom: 12px;">
            <h3 style="font-size: 16px; font-weight: 600; margin: 0;">{{ app.name }}</h3>
            <div class="flex gap-sm">
              <button @click="editApp(app)" class="icon-button glass-button">✏️</button>
              <button @click="removeApp(app.id)" class="icon-button glass-button">🗑️</button>
            </div>
          </div>
          <div style="font-size: 12px; opacity: 0.8; margin-bottom: 8px; word-break: break-all;">
            {{ app.path }}
          </div>
          <div v-if="app.args" style="font-size: 11px; opacity: 0.6;">
            Args: {{ app.args }}
          </div>
        </div>

        <!-- Empty state -->
        <div v-if="apps.length === 0" style="grid-column: 1 / -1; text-align: center; padding: 60px 20px; opacity: 0.7;">
          <div style="font-size: 48px; margin-bottom: 16px;">📱</div>
          <h3 style="font-size: 18px; margin-bottom: 8px;">No applications yet</h3>
          <p style="opacity: 0.8;">Click "Add App" to get started</p>
        </div>
      </div>
    </div>

    <!-- Add/Edit Dialog -->
    <div v-if="showDialog" class="dialog-overlay" @click.self="closeDialog">
      <div class="dialog glass-panel">
        <h2 style="margin-bottom: 20px;">{{ editingApp ? 'Edit' : 'Add' }} Application</h2>
        
        <div class="form-group">
          <label>Application Name *</label>
          <input v-model="dialogData.name" class="glass-input" placeholder="My App" />
        </div>

        <div class="form-group">
          <label>Executable Path *</label>
          <div class="flex gap-sm">
            <input v-model="dialogData.path" class="glass-input" placeholder="C:\path\to\app.exe" />
            <button @click="selectFile" class="glass-button">Browse</button>
          </div>
        </div>

        <div class="form-group">
          <label>Arguments (optional)</label>
          <input v-model="dialogData.args" class="glass-input" placeholder="--flag value" />
        </div>

        <div class="flex gap-md" style="margin-top: 24px;">
          <button @click="closeDialog" class="glass-button" style="flex: 1;">Cancel</button>
          <button @click="saveApp" class="glass-button primary" style="flex: 1;">Save</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { GetApps, AddApp, UpdateApp, RemoveApp, GetServerInfo, SelectFile, StartServer, StopServer } from '../wailsjs/go/main/App';
import { BrowserOpenURL } from '../wailsjs/runtime/runtime';
import { EventsOn } from '../wailsjs/runtime/runtime';
import QRCode from 'qrcode';

const apps = ref([]);
const serverInfo = ref({
  localURL: 'http://localhost:8000',
  networkURL: 'http://localhost:8000',
  status: 'stopped',
  running: false
});

const showDialog = ref(false);
const editingApp = ref(null);
const dialogData = ref({
  name: '',
  path: '',
  args: ''
});

const qrCanvas = ref(null);

onMounted(async () => {
  await loadApps();
  await loadServerInfo();
  
  // Listen for server events
  EventsOn('server:started', () => {
    loadServerInfo();
  });
  
  EventsOn('server:stopped', () => {
    loadServerInfo();
  });
  
  if (serverInfo.value.running) {
    generateQR();
  }
});

async function loadApps() {
  apps.value = await GetApps();
}

async function loadServerInfo() {
  serverInfo.value = await GetServerInfo();
  if (serverInfo.value.running && qrCanvas.value) {
    generateQR();
  }
}

function generateQR() {
  if (qrCanvas.value && serverInfo.value.running) {
    QRCode.toCanvas(qrCanvas.value, serverInfo.value.networkURL, {
      width: 140,
      margin: 1,
      color: {
        dark: '#FFFFFF',
        light: '#00000000'
      }
    });
  }
}

async function startServer() {
  try {
    await StartServer();
    await loadServerInfo();
    generateQR();
  } catch (err) {
    alert('Failed to start server: ' + err);
  }
}

async function stopServer() {
  if (confirm('Are you sure you want to stop the server? Mobile access will be disabled.')) {
    try {
      await StopServer();
      await loadServerInfo();
    } catch (err) {
      alert('Failed to stop server: ' + err);
    }
  }
}

function openAddDialog() {
  editingApp.value = null;
  dialogData.value = { name: '', path: '', args: '' };
  showDialog.value = true;
}

function editApp(app) {
  editingApp.value = app;
  dialogData.value = { ...app };
  showDialog.value = true;
}

async function saveApp() {
  if (!dialogData.value.name || !dialogData.value.path) {
    alert('Please fill in all required fields');
    return;
  }

  if (editingApp.value) {
    await UpdateApp(editingApp.value.id, dialogData.value.name, dialogData.value.path, dialogData.value.args);
  } else {
    await AddApp(dialogData.value.name, dialogData.value.path, dialogData.value.args);
  }

  await loadApps();
  closeDialog();
}

async function removeApp(id) {
  if (confirm('Are you sure you want to remove this application?')) {
    await RemoveApp(id);
    await loadApps();
  }
}

async function launchApp(id) {
  if (!serverInfo.value.running) {
    alert('Please start the server first');
    return;
  }
  
  const app = apps.value.find(a => a.id === id);
  if (app) {
    try {
      const response = await fetch(`${serverInfo.value.localURL}/api/launch/${id}`, { method: 'POST' });
      if (!response.ok) {
        alert('Failed to launch application');
      }
    } catch (err) {
      alert('Error launching application. Make sure the server is running.');
      console.error('Launch error:', err);
    }
  }
}

function closeDialog() {
  showDialog.value = false;
  editingApp.value = null;
}

async function selectFile() {
  const path = await SelectFile();
  if (path) {
    dialogData.value.path = path;
    if (!dialogData.value.name) {
      const filename = path.split('\\').pop().replace('.exe', '');
      dialogData.value.name = filename;
    }
  }
}

function openURL(url) {
  BrowserOpenURL(url);
}
</script>

<style scoped>
.qr-container canvas {
  display: block;
}

.icon-button {
  padding: 6px 10px;
  font-size: 14px;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  width: 500px;
  max-width: 90%;
  animation: fadeIn 0.3s ease-out;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  opacity: 0.9;
}
</style>
