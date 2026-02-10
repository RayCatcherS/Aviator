# Aviator Design System (Glassmorphism)

All UI elements in Aviator (Desktop and Web) must follow this design language to ensure a premium, unified experience.

## Theme Colors
- **Core Background**: Dark Slate (`#0f172a / bg-slate-900`) with animated blobs.
- **Glass Base**: `rgba(255, 255, 255, 0.03)` with `backdrop-filter: blur(16px)`.
- **Primary Color**: Cyan 400 (`#22d3ee`) - Used for active states and primary buttons.
- **Danger Color**: Red 500/600 - Used for destructive actions (Delete, Logout).

## Typography
- **Font**: Inter (Sans-Serif).
- **Style**: Modern, clean, with wide letter-spacing for monospaced elements.

## Component: Glass Button
Buttons should rarely use solid backgrounds. Instead, use the **Tinted Glass** style:

| State | Primary (Cyan) | Danger (Red) | Ghost (Standard) |
| :--- | :--- | :--- | :--- |
| **Base** | Transparent Cyan tint + border | Transparent Red tint + border | Transparent |
| **Hover** | Glow + opaque tint | Glow + opaque tint | Slight highlight |
| **Active** | `scale-96` | `scale-96` | `scale-96` |

### CSS Classes Reference (Web)
- `.glass-button`: Base layout and blur.
- `.glass-button.primary`: Cyan variant.
- `.glass-button.danger`: Red variant.
- `.glass-button.ghost`: No borders, low opacity.

## Iconography
- **Style**: SVG Stroke Icons (e.g., Lucide/Feather style).
- **Unification**: Always use the SVGs defined in the Desktop application for unified components (Settings, Launch, Offline, Power).
- **Stroke Width**: `2.0` or `2.5` for better visibility on glass backgrounds.
- **Color**: Matches the button's variant or Slate 400 by default.

## Interaction Rules
1. **Interactive Cards**: Should glow and lift (`translateY(-5px)`) on hover.
2. **Settings Gear**: Highlights on hover, but only turns cyan when the menu is active.
3. **No Placeholders**: Icons or gradients should always occupy space to prevent layout shifts.
