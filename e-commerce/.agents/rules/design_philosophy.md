# Design Philosophy & UI Guide - Chia Florist E-Commerce

This rule establishes the design principles, UI standards, and component guidelines for all AI agents and developers working on the `e-commerce` codebase.

---

## Core UI Standards & Form Geometry

### 1. Form Controls & Container Padding (`py-3`)
- **Uniform Padding**: All form inputs, primary action buttons, secondary back buttons, and alert notifications MUST enforce a uniform vertical padding of `py-3` (with `px-4`).
- **Visual Rhythm**: Ensures balanced element heights across forms and cards.

### 2. Standardized Corner Radius (`rounded-xl`)
- **Consistent Radius**: All interactive controls (inputs, primary buttons, secondary outline buttons, alert boxes, cards) MUST use `rounded-xl` (`12px`).
- Do NOT mix `rounded-lg`, `rounded-2xl`, or heavy pill shapes on standard form controls.

### 3. Button Typography & Font Consistency (`text-sm`)
- **Equal Font Size**: Primary filled buttons and secondary outline buttons (e.g. "Back", "Cancel") MUST share the exact same font size (`text-sm`) and line-height.
- **Font Weight Hierarchy**: Primary buttons use `font-bold`, while secondary outline buttons use `font-semibold`.
- Buttons with identical container shapes MUST have matching font sizes (`text-sm`) regardless of primary/outline variant.

### 4. Color Palette & Contrast
- **Bright Pastel Green Accent**: `#4ade80` (hover `#34d399`) for primary call-to-action buttons and active focus states (`focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20`).
- **Accent Text Color**: Text on bright pastel green background MUST use **Rich Dark Green** (`#245842`), NEVER plain black or white.
- **Brand Palette**:
  - Deep Emerald Green: `#1b4332`
  - Pure White Background: `#ffffff`
  - Neutral Surface: `bg-gray-50`, `border-gray-200`
  - Body Text: `text-gray-900` / `text-gray-600`

---

## Layout & Content Strategy

### 5. Indonesian Heritage Integration
- **Subtle Cultural Accents**: Minimalist geometric batik/motifs (e.g. Batik Kawung, Parang) as lightweight SVG dividers (`h-12`) on pure white backgrounds.

### 6. Pragmatic Content & Tone of Voice
- **No Fake Placeholders**: Never include placeholder navigation links, categories, or pages that are not implemented in the codebase.
- **Approachability**: Friendly, casual, and helpful copy for customer communication (e.g., WhatsApp assistance).

### 7. Icon & Emoji Usage Policy
- **Clickable Elements Only**: Icons/emojis are permitted ONLY inside interactive elements (buttons, action triggers). Do NOT use decorative icons in static headings or body text.
