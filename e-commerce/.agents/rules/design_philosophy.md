# Design Philosophy - Chia Florist E-Commerce

This rule establishes the design principles and UI/UX philosophy for all AI agents and developers working on the `e-commerce` codebase. It is derived from the refined design decisions established in components such as `CFooter.vue`.

---

## Core Principles

### 1. Organic Heritage Integration (Modern Cultural Warmth)
- **Subtle Indonesian Motifs**: Integrate traditional Indonesian cultural elements (e.g., minimalist SVG Batik Kawung, Parang patterns) as subtle section dividers or decorative accents.
- **Cultural Authenticity without Clutter**: Keep cultural patterns geometric, lightweight (`h-12`), non-overlapping, and set on pure white (`#ffffff`) backgrounds with brand color strokes (`#1b4332`).

### 2. Light-First & High Contrast Brand Identity
- **White Background Priority**: Always prioritize pure white (`bg-white`) backgrounds for key layouts so the brand logo (`/images/logo.png`) and imagery maintain maximum legibility and contrast.
- **Botanical Palette**:
  - Primary Accent: Deep Emerald Green (`#1b4332` / `text-emerald-900`)
  - Primary Surface: Pure White (`#ffffff`)
  - Secondary Surface: Soft Neutral Light Gray (`bg-gray-50`, `border-gray-100`)
  - Body Text: Charcoal Gray (`text-gray-600` / `text-gray-900`)

### 3. Container Padding Rules
- **No Extra Horizontal Padding on Footer Container**: Do not add horizontal padding (`px-*`) to the main inner footer container. Keep layout margins controlled naturally by `max-w-7xl mx-auto`.

### 4. Pragmatic & Truthful Content Strategy
- **No Non-Existent Placeholders**: Never include placeholder navigation links, categories, or pages that do not yet exist in the codebase (e.g., do not add blogs, collections, or fake shipping pages if the routes are not implemented).
- **Casual & Approachable Tone of Voice**:
  - Use warm, friendly, and approachable copy (in Indonesian / Bahasa) for customer communication.
  - Avoid cold corporate speak or pushy sales jargon. (e.g., *"Ada yang mau ditanyakan atau pengen request buket khusus? Chat kami di WhatsApp ya, kami siap bantu dengan senang hati!"*).
- **Low-Key Functional CTAs**: Present WhatsApp and contact CTAs inside clean, soft container cards with clear solid pill buttons rather than intrusive popups or flashy banners.

### 5. Layout Simplification & DOM Efficiency
- **Flat Containers over Layered Cards**: Avoid unnecessary card containers, heavy borders, or double-nested backgrounds around brand elements (e.g., logo, address) unless functionally required.
- **Clear Typographic Hierarchy**:
  - Section Headers: Crisp uppercase labels with bold weight and dark contrast (`text-xs font-semibold uppercase tracking-wider text-gray-900`).
  - Navigation Links: Medium-weight body text (`text-sm font-medium text-gray-600`).

### 6. Micro-Interactions & Polished Details
- **Animated Underline Hover Effects**:
  - All interactive links should feature a smooth left-to-right underline hover animation:
    `relative inline-block py-0.5 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100`
- **Tactile Button States**: Interactive buttons must include micro-animations (`hover:shadow`, `active:scale-[0.98]`, `transition-all`).
- **Mathematical SVG Precision**: SVG patterns and geometric icons must have calculated, non-overlapping coordinates for sharp display across all resolutions.

### 7. Icon & Emoji Usage Policy
- **No Icons/Emojis in Static Sections**: Do not include decorative icons or emojis inside non-interactive static layout sections, headings, or standard text paragraphs under any circumstances.
- **Actionable Elements Only**: Icons or emojis are strictly permitted ONLY inside clickable/interactive elements (such as buttons, clickable cards, list rows, or action triggers), and the icon/emoji must directly represent the corresponding user action (e.g., WhatsApp logo on a chat trigger button).
