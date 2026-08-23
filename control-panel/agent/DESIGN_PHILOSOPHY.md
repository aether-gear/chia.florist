# chia.florist - Control Panel UI Design Philosophy

This document outlines the visual guidelines, typography, layout structures, and coloring systems established during the authentication UI refactoring. Future development and modifications must adhere to these design rules to maintain a cohesive, premium brand experience.

---

## 🌿 1. Core Visual Aesthetic: "Organic Minimalist"

The design language of `chia.florist` moves away from dark, heavy, industrial slate/black layouts. Instead, it embraces an **organic, minimalist design** inspired by nature, greenhouses, and high-end botanical shops. 
*   **Breathing Room**: Keep interfaces airy, with generous padding (`p-6` to `p-16`).
*   **Borderless Forms**: Avoid embedding forms inside heavy card borders and floating shadows. Center forms directly inside clean backgrounds.
*   **Artistic Accents**: Use high-contrast retro artwork (e.g., 1-bit or dithered pixel-art illustrations of botanical gardens) to create a striking visual center-point.

---

## 🎨 2. Theme & Color Tokens (HSL)

Accents are strictly **organic green**, not black. All interactive states (buttons, focus rings, custom inputs, checkbox bgs) must resolve to the following global variables defined in `src/index.css`:

### Light Mode (`:root`)
*   `--primary: 142.4 71.8% 29.2%` (Elegant Forest Green — HSL representation of `#155e37`)
*   `--primary-foreground: 0 0% 100%` (Contrast White Text)
*   `--ring: 142.4 71.8% 39.2%` (Focus Accent Emerald — HSL representation of `#047857`)
*   `--radius: 0.625rem` (Softer, modern border cornering)

### Dark Mode (`.dark`)
*   `--primary: 142.4 70% 85%` (Soft Sage/Mint Green)
*   `--primary-foreground: 142.4 71.8% 12%` (Deep Green text for contrast)
*   `--ring: 142.4 70.6% 45.3%` (Active Mint border ring)

---

## ✍️ 3. Typography Hierarchy

Imported via `index.html` and configured in `tailwind.config.js`:
*   **Display Font**: **Outfit** (`font-display`). Use this font for main headers, page titles, hero text, and statements to elevate the brand's visual identity.
*   **Body Font**: **Inter** (`font-sans`). Use this font for labels, inputs, descriptive text, tables, and standard copy to ensure optimal readability.

---

## 📐 4. Layout Structure

### Split-Screen Auth & Settings Templates
For landing pages, authentication steps, and visual settings, employ a **split-screen grid** layout:
*   **Left Column (Form Panel)**: `lg:col-span-5 xl:col-span-4`. Full height on mobile. Contains brand logo, forms, and footer links centered vertically with breathing room.
*   **Right Column (Art Panel)**: `lg:col-span-7 xl:col-span-8 bg-zinc-950`. Hidden on mobile viewports (`hidden lg:block`). Displays the custom dithered illustration cover.

### Artwork Rendering
To display retro/dithered illustration files correctly without browser blurring, apply crisp image rendering filters:
```css
image-rendering: pixelated;
filter: brightness(1.02) contrast(1.25);
```

---

## 🎛️ 6. Overlay Right Panel Drawer Rule

*   **Mandatory Panel Standard**: ALWAYS use the global **Overlay Right Panel** (`Sheet` drawer component with `w-full sm:max-w-none md:w-[48vw] md:min-w-[480px]`) for any data input, entity creation, editing, account binding, or permission assignment workflows (e.g. creating staff, editing staff, binding user accounts, assigning shop permissions).
*   **No Center Modals/Popups for Inputs**: Avoid using centered modal dialogs for multi-field forms or editing. All input, editing, and configuration forms must slide in from the right edge as a full-height overlay panel with sticky header and action footer.

