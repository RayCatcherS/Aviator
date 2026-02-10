const API_BASE = "";
// console.log("Aviator Web App: v2.4 Loaded");

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

    if (btnGrid) btnGrid.blur();
    if (btnList) btnList.blur();

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
        const response = await fetch(`${API_BASE}/api/info?t=${new Date().getTime()}`);
        if (!response.ok) throw new Error(`Server returned ${response.status}`);

        const data = await response.json();
        currentHostname = data.hostname;

        // Update version display if element exists
        const versionEl = document.getElementById('web-version');
        if (versionEl && data.version) {
            versionEl.innerText = data.version;
        }

        if (data.auth_required && !data.is_authorized) {
            showAuthModal();
            document.getElementById('header-controls').classList.add('hidden');
            document.getElementById('view-section').classList.add('hidden');
        } else {
            hideAuthModal();
            document.getElementById('header-controls').classList.remove('hidden');
            document.getElementById('view-section').classList.remove('hidden');
        }

        updateServerStatus(true);
        return data.is_authorized || !data.auth_required;
    } catch (e) {
        console.error('Failed to fetch server info:', e);
        updateServerStatus(false);
        return false;
    }
}

function showAuthModal() {
    console.log("Displaying PIN Modal...");
    const modal = document.getElementById('login-overlay');
    if (modal) {
        modal.classList.remove('hidden');
        // Usiamo un timeout per il focus su mobile
        setTimeout(() => document.getElementById('pin-input').focus(), 100);
    } else {
        console.error("CRITICAL: login-overlay element NOT FOUND in DOM!");
    }
}

function hideAuthModal() {
    document.getElementById('login-overlay').classList.add('hidden');
}

async function handleLogin() {
    const pin = document.getElementById('pin-input').value;
    const btn = document.getElementById('btn-login');
    const err = document.getElementById('login-error');

    btn.disabled = true;
    err.classList.add('opacity-0');

    try {
        const response = await fetch(`${API_BASE}/api/auth`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ pin })
        });

        if (response.ok) {
            hideAuthModal();
            await fetchInfo(); // Update header controls
            await fetchApps();
            startPolling();
            showToast('✅ Authorization successful!', 2000);
        } else {
            err.classList.remove('opacity-0');
            document.getElementById('pin-input').value = '';
            document.getElementById('pin-input').focus();
        }
    } catch (e) {
        showToast('❌ Auth server error', 3000);
    } finally {
        btn.disabled = false;
    }
}

async function handleLogout() {
    try {
        await fetch(`${API_BASE}/api/logout`, { method: 'POST' });
        showToast('👋 Logged out');
        location.reload(); // Refresh to trigger auth check
    } catch (e) {
        showToast('❌ Logout error');
    }
}

function toggleSettingsMenu(event) {
    event.stopPropagation();
    const menu = document.getElementById('settings-menu');
    const btn = document.getElementById('btn-settings');
    const isHidden = menu.classList.toggle('hidden');

    // Toggle visual state on button
    if (!isHidden) {
        btn.classList.add('active-cyan');
    } else {
        btn.classList.remove('active-cyan');
        btn.blur();
    }
}

function closeSettingsMenu() {
    const menu = document.getElementById('settings-menu');
    const btn = document.getElementById('btn-settings');
    if (menu && !menu.classList.contains('hidden')) {
        menu.classList.add('hidden');
        if (btn) {
            btn.classList.remove('active-cyan');
            btn.blur();
        }
    }
}

// Close menu on click outside
window.addEventListener('click', () => {
    closeSettingsMenu();
});

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
        document.getElementById('view-section').classList.remove('hidden');

        // If was offline, restart polling and reload apps
        if (wasOffline) {
            showToast('✅ Server reconnected!', 3000);
            startPolling();
            fetchApps();
        }
    } else {
        // Server offline - red badge
        badge.className = 'inline-flex items-center px-3 py-1 rounded-full bg-red-500/10 border border-red-500/20 text-red-400 text-sm';
        dot.className = 'w-2 h-2 rounded-full bg-red-500 mr-2';
        statusText.innerHTML = 'Server Not Available';
        if (offlineBanner) offlineBanner.classList.remove('hidden');
        document.getElementById('view-section').classList.add('hidden');

        // Close any open modals on disconnect
        closeAppDetails();

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
    if (!serverOnline) return;

    try {
        const response = await fetch(`${API_BASE}/api/process-statuses`);
        if (!response.ok) {
            if (response.status !== 401) updateServerStatus(false);
            return;
        }
        const data = await response.json();
        processStatuses = data;
        updateStatusIndicators();
        updateServerStatus(true);
    } catch (e) {
        console.error('Failed to fetch process statuses:', e);
        if (e.message !== 'Unauthorized') updateServerStatus(false);
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
    updateModalStatus();
}

async function fetchApps() {
    try {
        const response = await fetch(`${API_BASE}/api/apps`);
        if (!response.ok) return;
        const apps = await response.json();
        renderGrid(apps);
    } catch (e) {
        console.error("Failed to fetch apps", e);
        if (e.message !== 'Unauthorized') {
            grid.innerHTML = `<div class="col-span-full text-center text-red-400">Error connecting to server.</div>`;
        }
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
        card.className = "glass-card glass-interactive rounded-2xl p-4 flex flex-col items-center justify-center group cursor-pointer text-left w-full h-full border border-white/5";
        card.setAttribute('data-app-id', app.id);
        card.onclick = (e) => {
            e.currentTarget.blur();
            openAppDetails(app);
        };

        // Use app icon if available, otherwise use gradient fallback
        let iconHTML;
        if (app.icon) {
            iconHTML = `<img src="data:image/png;base64,${app.icon}" class="app-icon w-12 h-12 rounded-xl shadow-lg mb-3 group-hover:scale-110 transition-transform duration-300 object-contain" alt="${app.name} icon">`;
        } else {
            const hue = app.name.charCodeAt(0) * 10 % 360;
            iconHTML = `<div class="app-icon w-12 h-12 rounded-xl bg-gradient-to-br from-[hsl(${hue},70%,50%)] to-[hsl(${hue + 40},70%,30%)] shadow-lg mb-3 flex items-center justify-center text-xl font-bold text-white group-hover:scale-110 transition-transform duration-300">
                ${app.name.substring(0, 2).toUpperCase()}
            </div>`;
        }

        card.innerHTML = `
            ${iconHTML}
            <div class="app-info flex flex-col items-center gap-1">
                <div class="flex items-center gap-3">
                    <div class="status-led ${processStatuses[app.id] ? 'led-running' : ''}" title="${processStatuses[app.id] ? 'Running' : 'Stopped'}"></div>
                    <h3 class="text-base font-semibold text-slate-100 group-hover:text-cyan-400 transition-colors truncate max-w-[120px]">${app.name}</h3>
                </div>
                <span class="launch-text text-[10px] text-slate-500 font-medium tracking-wider uppercase opacity-40 group-hover:opacity-100 transition-opacity">Launch App</span>
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
        } else if (response.status !== 401) {
            const err = await response.json();
            showToast(`Error: ${err.detail || 'Internal error'}`, 4000);
        }
    } catch (e) {
        if (e.message !== 'Unauthorized') showToast(`Network Error`, 3000);
    }
}

let currentlySelectedApp = null;

function openAppDetails(app) {
    currentlySelectedApp = app;
    const modal = document.getElementById('app-details-overlay');
    const iconContainer = document.getElementById('modal-app-icon');
    const nameContainer = document.getElementById('modal-app-name');
    const launchBtn = document.getElementById('modal-launch-btn');

    // Populate Data
    nameContainer.innerText = app.name;

    // Icon
    if (app.icon) {
        iconContainer.innerHTML = `<img src="data:image/png;base64,${app.icon}" class="w-full h-full object-contain">`;
    } else {
        const hue = app.name.charCodeAt(0) * 10 % 360;
        iconContainer.innerHTML = `<div class="w-full h-full bg-gradient-to-br from-[hsl(${hue},70%,50%)] to-[hsl(${hue + 40},70%,30%)] flex items-center justify-center text-3xl font-bold text-white uppercase">${app.name.substring(0, 2)}</div>`;
    }

    // Set Launch Handlers
    launchBtn.onclick = () => {
        launchApp(app.id, app.name);
        // Optional: closeAppDetails();
    };

    updateModalStatus();
    modal.classList.remove('hidden');
}

function closeAppDetails() {
    const modal = document.getElementById('app-details-overlay');
    modal.classList.add('hidden');
    currentlySelectedApp = null;
}

function updateModalStatus() {
    if (!currentlySelectedApp) return;

    const led = document.getElementById('modal-status-led');
    const text = document.getElementById('modal-status-text');
    const badge = document.getElementById('modal-app-status-badge');
    const isRunning = processStatuses[currentlySelectedApp.id];

    if (isRunning) {
        led.className = 'w-2 h-2 rounded-full bg-green-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.6)]';
        text.innerText = 'Running';
        text.className = 'text-xs font-semibold tracking-widest uppercase text-green-400';
        badge.className = 'inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-green-500/10 border border-green-500/20 mb-8';
    } else {
        led.className = 'w-2 h-2 rounded-full bg-slate-500';
        text.innerText = 'Stopped';
        text.className = 'text-xs font-semibold tracking-widest uppercase text-slate-400';
        badge.className = 'inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-white/5 border border-white/5 mb-8';
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

// Initial Load sequence
async function init() {
    grid.innerHTML = '<div class="col-span-full text-center py-20 text-cyan-400 animate-pulse">Connecting...</div>';

    try {
        const isAuthorized = await fetchInfo();

        if (isAuthorized) {
            await fetchApps();
            startPolling();
        } else {
            grid.innerHTML = '<div class="col-span-full text-center py-20 text-slate-500">Authorization required. Please enter PIN.</div>';
        }
    } catch (err) {
        console.error("Aviator Initialization Failed:", err);
        grid.innerHTML = `<div class="col-span-full text-center py-20 text-red-400">Connection Failed: ${err.message}</div>`;
    }
}

init();


