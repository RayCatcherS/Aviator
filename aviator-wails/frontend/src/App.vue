<template>
  <div id="app" class="h-screen w-screen flex flex-col bg-slate-900 text-slate-100 overflow-hidden relative selection:bg-cyan-500/30 font-sans">
    
    <!-- Animated Blobs Background -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none z-0">
        <div class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-purple-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
        <div class="absolute top-[-10%] right-[-10%] w-96 h-96 bg-cyan-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
        <div class="absolute bottom-[-10%] left-[20%] w-96 h-96 bg-pink-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
    </div>

    <!-- Custom Title Bar (Frameless) -->
    <div class="title-bar z-50 flex items-center justify-between h-8 bg-black/20 backdrop-blur-md border-b border-white/5 px-3 select-none flex-shrink-0">
      <div class="title-bar-left flex items-center gap-2">
        <span class="text-xs font-semibold text-white/90 tracking-wide">Aviator</span>
      </div>
      <div class="title-bar-controls flex gap-[1px]">
        <button class="title-btn w-12 h-8 flex items-center justify-center bg-transparent hover:bg-white/10 text-white/80 hover:text-white transition-colors" @click="openSettings" title="Settings">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
        </button>
        <button class="title-btn w-12 h-8 flex items-center justify-center bg-transparent hover:bg-white/10 text-white/80 hover:text-white transition-colors" @click="minimizeWindow" title="Minimize">
          <svg width="10" height="10" viewBox="0 0 12 12"><line x1="0" y1="6" x2="12" y2="6" stroke="currentColor" stroke-width="1"/></svg>
        </button>
        <button class="title-btn w-12 h-8 flex items-center justify-center bg-transparent hover:bg-white/10 text-white/80 hover:text-white transition-colors" @click="toggleMaximize" title="Maximize">
          <svg width="10" height="10" viewBox="0 0 12 12"><rect x="1" y="1" width="10" height="10" fill="none" stroke="currentColor" stroke-width="1"/></svg>
        </button>
        <button class="title-btn w-12 h-8 flex items-center justify-center bg-transparent hover:bg-red-600 text-white/80 hover:text-white transition-colors" @click="closeWindow" title="Close">
          <svg width="10" height="10" viewBox="0 0 12 12">
            <line x1="1" y1="1" x2="11" y2="11" stroke="currentColor" stroke-width="1"/>
            <line x1="11" y1="1" x2="1" y2="11" stroke="currentColor" stroke-width="1"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col p-6 gap-6 overflow-hidden z-10 relative">
      
      <!-- Status Header Panel -->
      <div v-tilt class="glass-card p-6 flex-shrink-0">
        <div class="flex flex-row justify-between items-center gap-6">
          <div class="flex flex-col gap-3 flex-1 min-w-0">
            <div class="flex items-center gap-4">
              <h1 class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-cyan-400 to-purple-500 drop-shadow-lg">Aviator</h1>
              <div class="px-3 py-1 rounded-full text-xs font-bold border" 
                   :class="serverInfo.running ? 'bg-green-500/10 border-green-500/20 text-green-400' : 'bg-slate-500/10 border-slate-500/20 text-slate-400'">
                {{ serverInfo.status.toUpperCase() }}
              </div>
            </div>
            
            <!-- Server Control Buttons -->
            <div class="flex gap-3 mt-2">
              <button 
                @click="startServer" 
                :disabled="serverInfo.running"
                class="glass-button success flex-none flex items-center justify-center gap-2"
                :class="{ 'opacity-50 cursor-not-allowed': serverInfo.running }"
              >
                <span class="text-lg">▶️</span> Start Server
              </button>
              <button 
                @click="stopServer"
                :disabled="!serverInfo.running"
                class="glass-button flex-none flex items-center justify-center gap-2 border-red-500/40 text-red-100 hover:bg-red-500/20"
                :class="{ 'opacity-50 cursor-not-allowed': !serverInfo.running }"
              >
                <span class="text-lg">⏹️</span> Stop Server
              </button>
            </div>
            
            <div v-if="serverInfo.running" class="flex flex-col gap-1 mt-2 text-sm text-slate-300 min-w-0">
              <div class="flex items-center gap-2 min-w-0">
                <span class="font-semibold text-slate-400 w-16 flex-shrink-0">Local:</span> 
                <a :href="serverInfo.localURL" @click.prevent="openURL(serverInfo.localURL)" class="hover:text-cyan-400 transition-colors truncate block flex-1" :title="serverInfo.localURL">{{ serverInfo.localURL }}</a>
              </div>
              <div class="flex items-center gap-2 min-w-0">
                <span class="font-semibold text-slate-400 w-16 flex-shrink-0">Network:</span> 
                <a :href="serverInfo.networkURL" @click.prevent="openURL(serverInfo.networkURL)" class="hover:text-cyan-400 transition-colors truncate block flex-1" :title="serverInfo.networkURL">{{ serverInfo.networkURL }}</a>
              </div>
            </div>
            <div v-else class="mt-2 text-sm text-slate-500 italic">
              Server stopped. Click "Start Server" to enable web access.
            </div>
          </div>
          
          <!-- QR Code -->
          <div v-if="serverInfo.running" class="glass-card p-4 rounded-xl flex-shrink-0 bg-white/5">
            <canvas ref="qrCanvas" class="block w-[140px] h-[140px]"></canvas>
          </div>
          <div v-else class="glass-card w-[172px] h-[172px] flex flex-col items-center justify-center gap-2 opacity-30 flex-shrink-0">
            <div class="text-5xl">🔒</div>
            <div class="text-xs font-medium">Server Offline</div>
          </div>
        </div>
      </div>

      <!-- Applications Section -->
      <div v-tilt class="glass-card flex-1 flex flex-col overflow-hidden min-h-0">
        <div class="p-6 pb-4 flex justify-between items-center border-b border-white/5">
          <h2 class="text-xl font-semibold text-slate-200">Applications</h2>
          <button @click="openAddDialog" class="glass-button primary flex items-center gap-2 text-sm">
            <span class="text-lg">+</span> Add App
          </button>
        </div>

        <!-- Apps Grid -->
        <div class="flex-1 overflow-y-auto p-6 app-grid-container custom-scrollbar">
          <div v-if="apps.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="app in apps" :key="app.id" class="glass-card glass-interactive p-4 hover:border-white/20 transition-all group">
              <div class="flex gap-4 items-start">
                <!-- App Icon -->
                <div v-if="app.icon" class="w-16 h-16 rounded-xl overflow-hidden flex-shrink-0 bg-black/20 p-2">
                  <img :src="'data:image/png;base64,' + app.icon" :alt="app.name" class="w-full h-full object-contain" />
                </div>
                <div v-else class="w-16 h-16 rounded-xl flex-shrink-0 flex items-center justify-center bg-gradient-to-br from-cyan-500/20 to-purple-500/20 border border-white/10">
                  <span class="text-2xl font-bold text-white/80">{{ app.name.substring(0, 2).toUpperCase() }}</span>
                </div>
                
                <!-- App Info -->
                <div class="flex-1 min-w-0">
                  <div class="flex justify-between items-start mb-1">
                    <div class="flex items-center gap-2">
                      <!-- LED Status Indicator -->
                      <div class="status-led" :class="{ 'led-running': processStatuses[app.id] }" :title="processStatuses[app.id] ? 'Running' : 'Stopped'"></div>
                      <h3 class="font-semibold text-lg text-slate-100 truncate group-hover:text-cyan-400 transition-colors">{{ app.name }}</h3>
                    </div>
                    <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button @click="editApp(app)" class="p-1.5 rounded hover:bg-white/10 text-slate-400 hover:text-white transition-colors" title="Edit">✏️</button>
                      <button @click="removeApp(app.id)" class="p-1.5 rounded hover:bg-red-500/20 text-slate-400 hover:text-red-400 transition-colors" title="Remove">🗑️</button>
                    </div>
                  </div>
                  <div class="text-xs text-slate-500 mb-2 font-mono truncate" :title="app.path">
                    {{ app.path }}
                  </div>
                  <div v-if="app.args" class="text-[10px] text-slate-600 font-mono truncate">
                    Args: {{ app.args }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty state -->
          <div v-else class="h-full flex flex-col items-center justify-center text-slate-500 opacity-60">
            <div class="text-6xl mb-4">📱</div>
            <h3 class="text-xl font-medium mb-2">No applications yet</h3>
            <p>Click "Add App" to get started</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Dialog -->
    <div v-if="showDialog" class="dialog-overlay">
      <div class="glass-card p-8 rounded-2xl w-full max-w-md shadow-2xl m-4 animate-fade-in-up">
        <h2 class="text-2xl font-bold mb-6 text-white">{{ editingApp ? 'Edit' : 'Add' }} Application</h2>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-semibold text-slate-400 mb-2">Application Name *</label>
            <input v-model="dialogData.name" class="glass-input" placeholder="My App" />
          </div>

          <div>
            <label class="block text-sm font-semibold text-slate-400 mb-2">Executable Path *</label>
            <div class="flex gap-2">
              <input v-model="dialogData.path" class="glass-input" placeholder="C:\path\to\app.exe" />
              <button @click="selectFile" class="glass-button whitespace-nowrap">Browse</button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-semibold text-slate-400 mb-2">Arguments (optional)</label>
            <input v-model="dialogData.args" class="glass-input" placeholder="--flag value" />
          </div>
        </div>

        <div class="flex gap-4 mt-8">
          <button @click="closeDialog" class="glass-button flex-1 bg-white/5 hover:bg-white/10">Cancel</button>
          <button @click="saveApp" class="glass-button primary flex-1 font-bold">Save</button>
        </div>
      </div>
    </div>

    <!-- Settings Dialog -->
    <div v-if="showSettings" class="dialog-overlay">
      <div class="glass-card p-8 rounded-2xl w-full max-w-sm shadow-2xl m-4 animate-fade-in-up">
        <h2 class="text-2xl font-bold mb-6 text-white flex items-center gap-3">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
          Settings
        </h2>
        
        <div class="space-y-6">
          <div class="flex items-center justify-between p-4 glass-card bg-white/5 rounded-xl">
            <div>
              <div class="font-semibold text-slate-200">Run at Startup</div>
              <div class="text-xs text-slate-400">Launch Aviator when Windows starts</div>
            </div>
            
            <button 
              @click="toggleAutoStart"
              class="w-12 h-6 rounded-full relative transition-colors duration-200 ease-in-out focus:outline-none"
              :class="settings.auto_start ? 'bg-cyan-500' : 'bg-slate-600'"
            >
              <div 
                class="absolute top-1 left-1 bg-white w-4 h-4 rounded-full transition-transform duration-200 ease-in-out shadow"
                :class="settings.auto_start ? 'translate-x-6' : 'translate-x-0'"
              ></div>
            </button>
          </div>
        </div>

        <div class="flex gap-4 mt-8">
          <button @click="closeSettings" class="glass-button w-full font-bold">Done</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { GetApps, AddApp, UpdateApp, RemoveApp, GetServerInfo, SelectFile, StartServer, StopServer, GetProcessStatuses, GetSettings, UpdateSettings } from '../wailsjs/go/main/App';
import { BrowserOpenURL, EventsOn, WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime';
import QRCode from 'qrcode';

const apps = ref([]);
const processStatuses = ref({});
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

const showSettings = ref(false);
const settings = ref({ auto_start: false });

const qrCanvas = ref(null);
let statusPollInterval = null;

onMounted(async () => {
  await loadApps();
  await loadServerInfo();
  await loadProcessStatuses();
  await loadSettings();
  
  // Poll process statuses every 2 seconds
  statusPollInterval = setInterval(loadProcessStatuses, 2000);
  
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

onUnmounted(() => {
  if (statusPollInterval) {
    clearInterval(statusPollInterval);
  }
});

async function loadApps() {
  apps.value = await GetApps();
}

async function loadProcessStatuses() {
  try {
    processStatuses.value = await GetProcessStatuses();
  } catch (err) {
    console.error('Failed to load process statuses:', err);
  }
}

async function loadSettings() {
  try {
    settings.value = await GetSettings();
  } catch (err) {
    console.error('Failed to load settings:', err);
  }
}

async function toggleAutoStart() {
  settings.value.auto_start = !settings.value.auto_start;
  try {
    await UpdateSettings(settings.value);
  } catch (err) {
    console.error('Failed to save settings:', err);
    // Revert on error
    settings.value.auto_start = !settings.value.auto_start;
    alert('Failed to update system settings: ' + err);
  }
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

function openSettings() {
  loadSettings(); // Refresh
  showSettings.value = true;
}

function closeSettings() {
  showSettings.value = false;
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

// Window controls
function minimizeWindow() {
  WindowMinimise();
}

function toggleMaximize() {
  WindowToggleMaximise();
}

function closeWindow() {
  Quit();
}

// 3D Tilt Effect Directive
const vTilt = {
  mounted(el) {
    el.style.transition = 'transform 0.1s ease-out';
    el.style.willChange = 'transform';
    el.style.transformStyle = 'preserve-3d';
    
    // Configurazione Effetto 3D (Tilt)
    // Modifica 'maxTilt' per cambiare l'intensità della rotazione
    const maxTilt = 3.0; // Gradi di rotazione massima (es. 1.5 leggero, 5.0 forte)
    const perspective = 1000; // Prospettiva (più basso = più deformato)
    
    const handleMove = (e) => {
      const rect = el.getBoundingClientRect();
      const width = rect.width;
      const height = rect.height;
      
      const mouseX = e.clientX - rect.left;
      const mouseY = e.clientY - rect.top;
      
      const xPct = mouseX / width - 0.5;
      const yPct = mouseY / height - 0.5;
      
      const xRot = yPct * maxTilt * -1; // Invert X for natural tilt
      const yRot = xPct * maxTilt;
      
      el.style.transform = `perspective(${perspective}px) rotateX(${xRot}deg) rotateY(${yRot}deg) scale(1.005)`;
    };
    
    const handleLeave = () => {
      el.style.transition = 'transform 0.5s ease-out'; // Slower return
      el.style.transform = `perspective(${perspective}px) rotateX(0) rotateY(0) scale(1)`;
      setTimeout(() => {
          el.style.transition = 'transform 0.1s ease-out'; // Reset for next move
      }, 500);
    };

    el.addEventListener('mousemove', handleMove);
    el.addEventListener('mouseleave', handleLeave);
  }
};
</script>

<style scoped>
/* Wails Window Drag Regions */
.title-bar {
  --wails-draggable: drag;
}

.title-bar-controls {
  --wails-draggable: no-drag;
}

.title-btn {
  -webkit-app-region: no-drag;
}

/* Custom Scrollbar for the app grid */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
