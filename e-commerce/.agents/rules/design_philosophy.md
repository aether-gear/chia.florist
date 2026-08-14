# Design Philosophy & UI Guide - Chia Florist E-Commerce

This rule establishes the design principles, UI standards, and component guidelines for all AI agents and developers working on the `e-commerce` codebase.

---

## Core UI Standards & Form Geometry

### 1. Form Controls & Container Padding (`py-3`)
- **Uniform Padding**: All form inputs, primary action buttons, secondary back buttons, and alert notifications MUST enforce a uniform vertical padding of `py-3` (with `px-4`).
- **Visual Rhythm**: Ensures balanced element heights across forms and cards.

### 2. Standardized Corner Radius (`rounded-xl` & `rounded-full`)
- **Standard Controls**: Interactive form controls (inputs, major form buttons, dialog controls, cards) MUST use `rounded-xl` (`12px`).
- **Pill Controls**: Header utilities (Navbar "Sign In"), Footer CTAs ("WhatsApp"), tag chips, and micro-triggers MUST use `rounded-full` pill shapes.

### 3. Button Component Hierarchy & Typography
- **Global Component**: All buttons across the application MUST use `<CButton>` (or adhere strictly to its styling).
- **Equal Font Size & Height**: Standard form buttons share `text-sm` font size and `py-3` vertical padding to match adjacent input fields.
- **Font Weight Hierarchy**: Primary and auth confirmation buttons use `font-bold`, while secondary, outline, and ghost buttons use `font-semibold`.

### 4. Color Palette & Contrast
- **Bright Pastel Green Accent**: `#4ade80` (hover `#34d399`) for primary call-to-action buttons and active focus states (`focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20`).
- **Accent Text Color**: Text on bright pastel green background MUST use **Rich Dark Green** (`#245842`), NEVER plain black or white.
- **Brand Palette**:
  - Deep Emerald Green: `#1b4332` (hover `#143326`)
  - Pure White Background: `#ffffff`
  - Neutral Surface: `bg-gray-50`, `border-gray-200`
  - Body Text: `text-gray-900` / `text-gray-600`

---

## Button Component System (`CButton`)

### 5. Standard Variants
- **`primary`**: Bright pastel green fill (`bg-[#4ade80] hover:bg-[#34d399] text-[#245842] font-bold shadow-xs hover:shadow`). Main call-to-action across pages and forms.
- **`secondary`**: Deep emerald fill (`bg-[#1b4332] hover:bg-[#143326] text-white font-bold shadow-xs hover:shadow`) or neutral fill (`bg-gray-100 hover:bg-gray-200 text-gray-800 font-semibold`). Used for secondary solid CTAs, header triggers, and footer brand CTAs.
- **`outline`**: Border standard (`border border-gray-200 hover:bg-gray-50 text-gray-700 font-semibold`) or brand accent outline (`border-2 border-[#4ade80] text-[#245842] hover:bg-[#4ade80]/10 font-bold`). Used for back buttons, secondary options, and neutral triggers.
- **`ghost`**: Transparent background (`hover:bg-gray-100 text-gray-700 font-semibold`). Used for inline actions, icon triggers, and subtle navigation controls.
- **`danger`**: Red alert fill (`bg-red-500 hover:bg-red-600 text-white font-semibold shadow-xs`). Used for destructive actions (Sign Out, Delete, Cancel).

### 6. Sizes & Shapes
- **`pill` / `sm`**: Pill shape (`rounded-full`, `px-3.5 py-1.5`, `text-xs font-semibold`). Used for header navbar triggers (e.g. "Sign In"), footer WhatsApp CTA, filter chips, and inline micro-actions.
- **`md`**: Standard component shape (`rounded-xl`, `px-5 py-2.5`, `text-sm font-semibold`). Used for dialog controls and general card buttons.
- **`lg` / `auth`**: Major CTA & Auth confirmation shape (`rounded-xl`, `py-3 px-4`, `text-sm font-bold`). Enforces `py-3` height matching form input controls for login, registration, and checkout confirmations.

### 7. Interactive States
- **Hover**: Smooth background & color transitions (`transition-all duration-200 ease-in-out`) with elevation shadow boost (`hover:shadow`).
- **Active**: Touch & click micro-animation (`active:scale-[0.99]`).
- **Loading**: Interactive lock (`disabled:opacity-60 disabled:cursor-not-allowed`), animated SVG spinner icon, and loading label state.
- **Disabled**: Reduced opacity (`disabled:opacity-50 disabled:cursor-not-allowed`).

---

## Layout & Content Strategy

### 8. Indonesian Heritage Integration
- **Subtle Cultural Accents**: Minimalist geometric batik/motifs (e.g. Batik Kawung, Parang) as lightweight SVG dividers (`h-12`) on pure white backgrounds.

### 9. Pragmatic Content & Tone of Voice
- **No Fake Placeholders**: Never include placeholder navigation links, categories, or pages that are not implemented in the codebase.
- **Approachability**: Friendly, casual, and helpful copy for customer communication (e.g., WhatsApp assistance).

### 10. Icon & Emoji Usage Policy
- **Clickable Elements Only**: Icons/emojis are permitted ONLY inside interactive elements (buttons, action triggers). Do NOT use decorative icons in static headings or body text.

