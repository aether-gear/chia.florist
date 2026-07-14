# Antigravity Project Rules - control-panel

This file contains project-scoped rules and behavioral constraints for the `control-panel` workspace.

## 1. Directory Modification Restrictions
* Do not modify, create, or delete any files in other modules outside the `control-panel` workspace.
* You may only read, search, or look up code/documentation in other modules to understand their structure and behavior.

## 2. Unimplemented Features in Other Modules
* If a feature needs to be added to the `control-panel` but depends on functionality not yet present in other modules:
  * Consult the parent API documentation located at [docs/api](file:///D:/__Projects/kage/chia.florist/docs/api) or check the respective module's source code.
  * Do **not** attempt to modify or add the feature to other modules yourself.
  * Instead, suggest that the user ask the responsible programmer via Discord or WhatsApp, or create a GitHub issue and mention them.

## 3. UI Design Philosophy & Theme Constraints
*   **Creative Minimalist Layout**: Avoid traditional card containers or heavy borders. Prefer airy, borderless, content-driven layouts with high-contrast visual focus points (like retro 1-bit or dithered illustrations representing the botanical theme of `chia.florist`).
*   **Green Accent Overrides**: The primary brand accent is organic green, not black or slate. Ensure all interactive components (buttons, links, active rings, checked checkboxes) utilize the custom HSL design tokens configured in `index.css`:
    *   **Light Mode**: Deep forest green (`--primary` = `142.4 71.8% 29.2%`) and emerald green (`--ring` = `142.4 71.8% 39.2%`).
    *   **Dark Mode**: Soft sage/mint green (`--primary` = `142.4 70% 85%`) and light green (`--ring` = `142.4 70.6% 45.3%`).
*   **Typography Hierarchy**:
    *   Use **Outfit** (`font-display`) for page headers, hero statements, and title components.
    *   Use **Inter** (`font-sans`) for form labels, descriptions, and general text copy.
*   **Responsive Auth Design**:
    *   Auth and settings templates must use a split-screen design (`lg:grid-cols-12`) where one side handles inputs and forms and the other handles dithered illustration art.
    *   Harness responsive visibility (e.g., `hidden lg:block` for the art panel) to conceal heavy visual elements on smaller displays (mobile/tablet) to optimize loading performance.
*   **Checkbox Class Variable Usage**: Avoid hardcoding neutral slate borders and backgrounds. Use variables like `bg-primary`, `border-input`, and `text-primary-foreground` to support active accent adjustments dynamically.
