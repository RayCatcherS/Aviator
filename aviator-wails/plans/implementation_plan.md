# Implementation Plan: Authentication & Connectivity Fix

## Context
Aviator web app was failing to load on mobile and incognito browsers due to server-side initialization hangs and network timeouts.

## Root Causes
1. **Mutex Deadlocks**: Concurrent access to `ConfigManager` and session buckets was causing the Go server to hang.
2. **Registry Latency**: Synchronous Windows Registry calls during the HTTP request cycle were too slow.
3. **Improper Routing**: API endpoints were behind conditional logic that could be bypassed or blocked.

## Implementation Details
1. **Go Server Structure**:
   - Refactored `ServeHTTP` to handle `/ping` and `/api/info` with maximum priority.
   - Removed blocking log statements in hot paths.
2. **Configuration Management**:
   - Fixed `ConfigManager` to release locks before calling `SaveSettings`.
   - Optimized `LoadSettings` to handle registry calls more efficiently.
3. **Frontend Connectivity**:
   - Added `auth_required` and `is_authorized` flags to the `/api/info` response.
   - Implemented a more robust `init()` cycle in `app.js` with clear state transitions.
   - Refined desktop icons: Replaced emojis with modern SVG stroke icons for a more professional glassmorphism look.
   - Enhanced PIN management: Added an "eye" toggle for visibility and a dedicated "Modify PIN" dialog in the desktop app.
   - Optimized web app UI: Centered and resized the PIN input, and added a top-right settings gear with a logout menu.
   - Added `Enter` key support for PIN submission.

## Verification
- [x] Server responds to `/ping` instantly.
- [x] PIN modal appears when unauthorized.
- [x] PIN submission successfully sets session cookie.
- [x] App list loads correctly after login.
