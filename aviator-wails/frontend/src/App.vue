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
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
                Start Server
              </button>
              <button 
                @click="stopServer"
                :disabled="!serverInfo.running"
                class="glass-button flex-none flex items-center justify-center gap-2 border-red-500/40 text-red-100 hover:bg-red-500/20"
                :class="{ 'opacity-50 cursor-not-allowed': !serverInfo.running }"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect></svg>
                Stop Server
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
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-slate-500/50">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
            <div class="text-xs font-semibold uppercase tracking-widest text-slate-500 mt-2">Server Offline</div>
          </div>
        </div>
      </div>

      <!-- Applications Section -->
      <div v-tilt class="glass-card flex-1 flex flex-col overflow-hidden min-h-0">
        <div class="p-6 pb-4 flex justify-between items-center border-b border-white/5">
          <h2 class="text-xl font-semibold text-slate-200">Applications</h2>
          <button @click="openAddDialog" class="glass-button primary flex items-center gap-2 text-sm">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            Add App
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
                      <button @click="editApp(app)" class="p-1.5 rounded hover:bg-white/10 text-slate-400 hover:text-white transition-colors" title="Edit">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                      </button>
                      <button @click="removeApp(app.id)" class="p-1.5 rounded hover:bg-red-500/20 text-slate-400 hover:text-red-400 transition-colors" title="Remove">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
                      </button>
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
          <div v-else class="h-full flex flex-col items-center justify-center text-slate-500/50 py-12">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="opacity-20 mb-4">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="3" y1="9" x2="21" y2="9"></line>
              <line x1="9" y1="21" x2="9" y2="9"></line>
            </svg>
            <h3 class="text-xl font-medium mb-1">No applications yet</h3>
            <p class="text-sm">Click "Add App" to get started</p>
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

          <!-- Web Access Security -->
          <div class="space-y-3 p-4 glass-card bg-white/5 rounded-xl border border-white/5">
            <div class="flex justify-between items-start">
              <div>
                <div class="font-semibold text-slate-200">Web Access Security</div>
                <div class="text-xs text-slate-400 mb-3">Remote deck protection</div>
              </div>
              <div v-if="settings.auth_enabled" class="px-2 py-0.5 rounded bg-green-500/10 text-green-400 text-[10px] font-bold border border-green-500/20">
                ENABLED
              </div>
              <div v-else class="px-2 py-0.5 rounded bg-red-500/10 text-red-400 text-[10px] font-bold border border-red-500/20">
                DISABLED
              </div>
            </div>
            
            <div class="flex items-center gap-3 p-3 bg-black/30 rounded-lg border border-white/5">
                <div class="flex-1 font-mono tracking-widest text-sm text-center text-slate-400">
                   {{ settings.auth_enabled ? '••••••' : 'NO PIN SET' }}
                </div>
            </div>

            <button @click="openPinDialog" class="glass-button primary w-full flex items-center justify-center gap-2 py-2">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3L15.5 7.5z"></path></svg>
              Modify PIN
            </button>
          </div>

          <!-- Version Display -->
          <div class="pt-2 text-center">
            <span class="text-[10px] font-mono text-slate-600 tracking-widest uppercase">Aviator {{ appVersion }}</span>
          </div>
        </div>

        <div class="flex gap-4 mt-8">
          <button @click="closeSettings" class="glass-button w-full font-bold">Done</button>
        </div>
      </div>
    </div>

    <!-- Modify PIN Dialog -->
    <div v-if="showPinDialog" class="dialog-overlay">
      <div class="glass-card p-8 rounded-2xl w-full max-w-sm shadow-2xl m-4 animate-fade-in-up">
        <h2 class="text-2xl font-bold mb-2 text-white">Modify Security PIN</h2>
        <p class="text-xs text-slate-400 mb-6 text-center italic">Type 4-6 digits. Leave empty to disable security.</p>
        
        <div class="space-y-4">
          <div class="relative">
            <input 
              :type="showNewPinValue ? 'text' : 'password'" 
              v-model="newPinData" 
              maxlength="6"
              placeholder="Enter new PIN"
              class="glass-input text-center font-mono tracking-widest text-base h-14 pr-12 focus:border-cyan-500"
              ref="pinInputRef"
            />
            <button 
                @click="showNewPinValue = !showNewPinValue" 
                class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-500 hover:text-white transition-colors p-2"
                title="Toggle visibility"
            >
                <svg v-if="showNewPinValue" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
            </button>
          </div>
        </div>

        <div class="flex gap-4 mt-8">
          <button @click="closePinDialog" class="glass-button flex-1 bg-white/5 hover:bg-white/10">Cancel</button>
          <button @click="updatePin" class="glass-button primary flex-1 font-bold">Update</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { GetApps, AddApp, UpdateApp, RemoveApp, GetServerInfo, SelectFile, StartServer, StopServer, GetProcessStatuses, GetSettings, UpdateSettings, SetWebPIN, GetVersion } from '../wailsjs/go/main/App';
import { BrowserOpenURL, EventsOn, WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime';
import QRCode from 'qrcode';

const apps = ref([]);
const appVersion = ref('');
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
const settings = ref({ auto_start: false, auth_enabled: false });
const webPin = ref(''); // Internal state for the input in settings (now mostly for display)

const showPinDialog = ref(false);
const newPinData = ref('');
const showNewPinValue = ref(false);
const pinInputRef = ref(null);

const qrCanvas = ref(null);
let statusPollInterval = null;

onMounted(async () => {
  await loadApps();
  await loadServerInfo();
  await loadProcessStatuses();
  await loadSettings();
  appVersion.value = await GetVersion();
  
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

function openPinDialog() {
  newPinData.value = '';
  showNewPinValue.value = false;
  showPinDialog.value = true;
}

function closePinDialog() {
  showPinDialog.value = false;
}

async function updatePin() {
  try {
    await SetWebPIN(newPinData.value);
    await loadSettings();
    closePinDialog();
    alert(settings.value.auth_enabled ? '✅ Security PIN updated successfully!' : '🔓 Security PIN removed.');
  } catch (err) {
    alert('Failed to update PIN: ' + err);
  }
}

async function saveWebPin() {
  // Keeping this for compatibility or simpler use, but now we use updatePin
  updatePin();
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
