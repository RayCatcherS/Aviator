const API_BASE = ""; // Relative path since we are served by the same server

// Logic involved in sticky header
window.addEventListener('scroll', () => {
    const scrollPos = window.scrollY;
    // Add 'scrolled' class to body if scrolled down more than 100px
    if (scrollPos > 100) {
        document.body.classList.add('scrolled');
    } else {
        document.body.classList.remove('scrolled');
    }
});

// Logic for View Switching
let currentView = 'grid';

function setView(mode) {
    currentView = mode;
    const gridEl = document.getElementById('app-grid');
    const btnGrid = document.getElementById('btn-grid');
    const btnList = document.getElementById('btn-list');

    if (mode === 'list') {
        gridEl.classList.add('list-view');
        btnList.classList.replace('text-slate-500', 'text-cyan-400');
        btnGrid.classList.replace('text-cyan-400', 'text-slate-500');
    } else {
        gridEl.classList.remove('list-view');
        btnGrid.classList.replace('text-slate-500', 'text-cyan-400');
        btnList.classList.replace('text-cyan-400', 'text-slate-500');
    }
}

const grid = document.getElementById('app-grid');
const hostnameSpan = document.getElementById('hostname');
const toastEl = document.getElementById('toast');
const toastMsg = document.getElementById('toast-message');

// Track process statuses
let processStatuses = {};
let serverOnline = true;
let processStatusInterval = null;
let serverStatusInterval = null;

let currentHostname = '...';

async function fetchInfo() {
    try {
        // Add timestamp to prevent caching
        const response = await fetch(`${API_BASE}/api/info?t=${new Date().getTime()}`);
        if (!response.ok) throw new Error(`Server returned ${response.status}`);

        const data = await response.json();
        currentHostname = data.hostname;
        updateServerStatus(true);
        return true;
    } catch (e) {
        console.error('Failed to fetch server info:', e);
        updateServerStatus(false);
        return false;
    }
}

function updateServerStatus(online) {
    const wasOffline = !serverOnline;
    serverOnline = online;
    const badge = document.getElementById('connection-badge');
    const dot = document.getElementById('status-dot');
    const statusText = document.getElementById('status-text');
    const offlineBanner = document.getElementById('offline-banner');
    const viewControls = document.querySelector('.view-controls');

    if (online) {
        // Server online - green badge
        badge.className = 'inline-flex items-center px-3 py-1 rounded-full bg-green-500/10 border border-green-500/20 text-green-400 text-sm';
        dot.className = 'w-2 h-2 rounded-full bg-green-500 mr-2 animate-pulse';
        statusText.innerHTML = `Connected to <span id="hostname" class="font-bold ml-1">${currentHostname}</span>`;
        if (offlineBanner) offlineBanner.classList.add('hidden');
        if (viewControls) viewControls.style.display = ''; // Restore default display

        // If was offline, restart polling and reload apps
        if (wasOffline) {
            showToast('✅ Server reconnected!', 3000);
            startPolling();
            fetchApps();
            connectWebSocket();
        }
    } else {
        // Server offline - red badge
        badge.className = 'inline-flex items-center px-3 py-1 rounded-full bg-red-500/10 border border-red-500/20 text-red-400 text-sm';
        dot.className = 'w-2 h-2 rounded-full bg-red-500 mr-2';
        statusText.innerHTML = 'Server Not Available';
        if (offlineBanner) offlineBanner.classList.remove('hidden');
        if (viewControls) viewControls.style.display = 'none';

        // Clear apps grid when offline (completely remove content)
        grid.innerHTML = '';

        // Stop polling to save resources
        stopPolling();

        // Show toast notification only on first disconnect
        if (wasOffline === false) {
            showToast('⚠️ Server is offline. Please start the server on the host PC.', 5000);
        }
    }
}

function startPolling() {
    // Clear any existing intervals
    stopPolling();

    // Start process status polling every 2 seconds
    processStatusInterval = setInterval(fetchProcessStatuses, 2000);

    // Start server status polling every 5 seconds
    serverStatusInterval = setInterval(fetchInfo, 5000);
}

function stopPolling() {
    if (processStatusInterval) {
        clearInterval(processStatusInterval);
        processStatusInterval = null;
    }
    if (serverStatusInterval) {
        clearInterval(serverStatusInterval);
        serverStatusInterval = null;
    }
}

async function retryConnection() {
    const btn = document.getElementById('btn-retry');
    const icon = document.getElementById('retry-icon');
    const text = document.getElementById('retry-text');

    // UI Loading state
    if (btn) btn.disabled = true;
    if (icon) icon.classList.add('rotate-180');
    if (text) text.innerText = 'Connecting...';

    showToast('🔄 Connecting to server...', 2000);

    // Small delay to make the interaction feel responsive before network call
    await new Promise(r => setTimeout(r, 500));

    try {
        // Force check server info and capture result
        const serverIsBack = await fetchInfo();

        if (serverIsBack && serverOnline) {
            // Success!
            if (icon) icon.classList.remove('rotate-180');
            if (text) text.innerText = 'Success!';

            // Reload everything
            console.log('Retry successful, reloading apps...');
            await fetchApps();
            await fetchProcessStatuses();
            startPolling();

            showToast('✅ Connected successfully!', 3000);

            // Reset button after short delay
            setTimeout(() => {
                if (btn) btn.disabled = false;
                if (text) text.innerText = 'Retry Connection';
            }, 1000);
        } else {
            throw new Error('Still offline or fetchInfo failed');
        }
    } catch (e) {
        console.error('Retry failed:', e);
        // Failed
        if (btn) btn.disabled = false;
        if (icon) icon.classList.remove('rotate-180');
        if (text) text.innerText = 'Retry Connection';
        showToast('❌ Server still unreachable', 3000);
    }
}


async function fetchProcessStatuses() {
    if (!serverOnline) return; // Skip if server is offline

    try {
        const response = await fetch(`${API_BASE}/api/process-statuses`);
        if (!response.ok) {
            updateServerStatus(false);
            return;
        }
        const data = await response.json();
        processStatuses = data;
        updateStatusIndicators();
        updateServerStatus(true);
    } catch (e) {
        console.error('Failed to fetch process statuses:', e);
        updateServerStatus(false);
    }
}

function updateStatusIndicators() {
    // Update LED indicators for all apps
    document.querySelectorAll('[data-app-id]').forEach(card => {
        const appId = card.getAttribute('data-app-id');
        const led = card.querySelector('.status-led');
        if (led) {
            if (processStatuses[appId]) {
                led.classList.add('led-running');
            } else {
                led.classList.remove('led-running');
            }
        }
    });
}

async function fetchApps() {
    try {
        const response = await fetch(`${API_BASE}/api/apps`);
        const apps = await response.json();
        renderGrid(apps);
    } catch (e) {
        console.error("Failed to fetch apps", e);
        grid.innerHTML = `<div class="col-span-full text-center text-red-400">Error connecting to server.</div>`;
    }
}

function renderGrid(apps) {
    grid.innerHTML = '';

    if (apps.length === 0) {
        grid.innerHTML = `
            <div class="col-span-full text-center py-20 text-slate-500 glass-card rounded-2xl p-8">
                <p class="text-xl mb-2">No apps configured.</p>
                <p class="text-sm">Go to the Host PC and add some apps!</p>
            </div>
        `;
        return;
    }

    apps.forEach(app => {
        const card = document.createElement('button');
        card.className = "glass-card rounded-2xl p-6 flex flex-col items-center justify-center group cursor-pointer text-left w-full h-full";
        card.setAttribute('data-app-id', app.id);
        card.onclick = () => launchApp(app.id, app.name);

        // Use app icon if available, otherwise use gradient fallback
        let iconHTML;
        if (app.icon) {
            iconHTML = `<img src="data:image/png;base64,${app.icon}" class="w-16 h-16 rounded-2xl shadow-lg mb-4 group-hover:scale-110 transition-transform duration-300 object-contain" alt="${app.name} icon">`;
        } else {
            // Fallback: Generate a random gradient based on name char code
            const hue = app.name.charCodeAt(0) * 10 % 360;
            iconHTML = `<div class="app-icon w-16 h-16 rounded-2xl bg-gradient-to-br from-[hsl(${hue},70%,50%)] to-[hsl(${hue + 40},70%,30%)] shadow-lg mb-4 flex items-center justify-center text-2xl font-bold text-white group-hover:scale-110 transition-transform duration-300">
                ${app.name.substring(0, 2).toUpperCase()}
            </div>`;
        }

        card.innerHTML = `
            ${iconHTML}
            <div class="flex flex-col items-center">
                <div class="flex items-center gap-2 mb-1">
                    <div class="status-led ${processStatuses[app.id] ? 'led-running' : ''}" title="${processStatuses[app.id] ? 'Running' : 'Stopped'}"></div>
                    <h3 class="text-lg font-semibold text-slate-100 group-hover:text-cyan-400 transition-colors">${app.name}</h3>
                </div>
                <span class="launch-text text-xs text-slate-500 mt-1 truncate max-w-full opacity-60">Click to Launch</span>
            </div>
        `;
        grid.appendChild(card);
    });
}

async function launchApp(id, name) {
    showToast(`Launching ${name}...`);
    try {
        const response = await fetch(`${API_BASE}/api/launch/${id}`, { method: 'POST' });
        if (response.ok) {
            showToast(`${name} launched successfully!`, 3000);
        } else {
            const err = await response.json();
            showToast(`Error: ${err.detail}`, 4000);
        }
    } catch (e) {
        showToast(`Network Error`, 3000);
    }
}

function showToast(msg, duration = 2000) {
    toastMsg.innerText = msg;
    toastEl.classList.remove('opacity-0');
    toastEl.classList.add('opacity-100', 'translate-y-[-10px]');

    setTimeout(() => {
        toastEl.classList.add('opacity-0');
        toastEl.classList.remove('opacity-100', 'translate-y-[-10px]');
    }, duration);
}

// Initial Load
fetchInfo();
fetchApps();

// Start polling
startPolling();

// Setup WebSocket for real-time updates
const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
// If API_BASE is set (remote), use that host, otherwise use window.location
const wsHost = API_BASE ? API_BASE.replace('http', 'ws') : `${protocol}//${window.location.host}`;
const wsUrl = `${wsHost}/ws`;
let wsClient = null;

function connectWebSocket() {
    if (!serverOnline) return;
    if (wsClient && wsClient.readyState === WebSocket.OPEN) return;

    try {
        wsClient = new WebSocket(wsUrl);

        wsClient.onopen = () => {
            console.log("Connected to Real-time updates");
        };

        wsClient.onmessage = (event) => {
            console.log("Update received:", event.data);
            fetchApps();
            // showToast("App list updated", 2000);
        };

        wsClient.onclose = () => {
            console.log("Socket disconnected");
            wsClient = null;
        };

        wsClient.onerror = (err) => {
            console.error("Socket error", err);
            if (wsClient) wsClient.close();
            wsClient = null;
        };
    } catch (e) {
        console.error("Failed to connect websocket", e);
    }
}

// Initial connection attempt
if (serverOnline) connectWebSocket();
