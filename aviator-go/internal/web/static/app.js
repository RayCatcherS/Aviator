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

async function fetchInfo() {
    try {
        const response = await fetch(`${API_BASE}/api/info`);
        const data = await response.json();
        hostnameSpan.innerText = data.hostname;
    } catch (e) {
        hostnameSpan.innerText = "Offline?";
        hostnameSpan.classList.add("text-red-500");
    }
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
        card.onclick = () => launchApp(app.id, app.name);

        // Generate a random gradient based on name char code for "icon"
        const hue = app.name.charCodeAt(0) * 10 % 360;

        card.innerHTML = `
            <div class="app-icon w-16 h-16 rounded-2xl bg-gradient-to-br from-[hsl(${hue},70%,50%)] to-[hsl(${hue + 40},70%,30%)] shadow-lg mb-4 flex items-center justify-center text-2xl font-bold text-white group-hover:scale-110 transition-transform duration-300">
                ${app.name.substring(0, 2).toUpperCase()}
            </div>
            <div class="flex flex-col">
                <h3 class="text-lg font-semibold text-slate-100 group-hover:text-cyan-400 transition-colors">${app.name}</h3>
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

// Setup WebSocket for real-time updates
const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const wsUrl = `${protocol}//${window.location.host}/ws`;

function connectWebSocket() {
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log("Connected to Real-time updates");
        showToast("Connected to Server", 2000);
    };

    ws.onmessage = (event) => {
        console.log("Update received:", event.data);
        // Simple update trigger
        fetchApps();
        showToast("App list updated", 2000);
    };

    ws.onclose = () => {
        console.log("Disconnected. Reconnecting in 3s...");
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = (err) => {
        console.error("Socket error", err);
        ws.close();
    };
}

connectWebSocket();
