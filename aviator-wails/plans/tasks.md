# Tasks Checklist

## Debugging & Infrastructure
- [x] Add `/ping` endpoint for basic health check.
- [x] Implement trace logging in `server.go` to find hang points.
- [x] Fix lock contention in `ConfigManager`.

## Authentication Flow
- [x] Priority routing for `/api/info`.
- [x] Dynamic check for PIN authorization status.
- [x] PIN modal auto-focus and Enter key submit handler.
- [x] Secure `HttpOnly` cookie session management.
- [x] Dedicated "Modify PIN" dialog with visibility toggle.
- [x] Refined desktop icons (SVG instead of emojis).
- [x] Refined web app PIN entry (centered and resized).
- [x] Web app settings gear + dropdown menu with logout.

## Cleanup
- [x] Remove verbose debug logs from Go backend.
- [x] Remove diagnostic console logs from Javascript frontend.
- [x] Bump app version to `2.4` for cache busting.
- [x] Document fixes in the `plans/` directory.
