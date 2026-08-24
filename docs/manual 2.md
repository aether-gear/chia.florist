# Chia Florist E-Commerce Platform — Complete User Manual Guide

> **Document Version:** 3.0.0  
> **Target Application:** Chia Florist (Papan Bunga & Custom Simulator Online)  
> **Scope:** Complete end-to-end user navigation, features, interfaces, and click scenarios in text form.

---

## Table of Contents

1. [Introduction & Platform Overview](#1-introduction--platform-overview)
2. [Global Layout & Navigation Elements](#2-global-layout--navigation-elements)
   - 2.1 [Top Navigation Bar (`CNavbar`)](#21-top-navigation-bar-cnavbar)
   - 2.2 [Global Footer (`CFooter`)](#22-global-footer-cfooter)
3. [User Authentication & Account Management](#3-user-authentication--account-management)
   - 3.1 [Scenario: User Registration (Sign Up)](#31-scenario-user-registration-sign-up)
   - 3.2 [Scenario: Account OTP Verification](#32-scenario-account-otp-verification)
   - 3.3 [Scenario: User Sign In (Login)](#33-scenario-user-sign-in-login)
   - 3.4 [Scenario: Password Reset / Recovery](#34-scenario-password-reset--recovery)
   - 3.5 [Scenario: Profile Management & 2D Avatar Cropper](#35-scenario-profile-management--2d-avatar-cropper)
   - 3.6 [Scenario: Shipping Address Book Management](#36-scenario-shipping-address-book-management)
4. [Product Discovery & Catalog Browsing](#4-product-discovery--catalog-browsing)
   - 4.1 [Scenario: Exploring Homepage & Quick Action Shortcuts](#41-scenario-exploring-homepage--quick-action-shortcuts)
   - 4.2 [Scenario: Full Catalog Browsing with Multi-Store & Sorting Filters](#42-scenario-full-catalog-browsing-with-multi-store--sorting-filters)
   - 4.3 [Scenario: Quick Search Drawer & Dedicated Search Page](#43-scenario-quick-search-drawer--dedicated-search-page)
5. [Standard Product Ordering Flow](#5-standard-product-ordering-flow)
   - 5.1 [Scenario: Product Detail Page & Dynamic Price Calculation](#51-scenario-product-detail-page--dynamic-price-calculation)
   - 5.2 [Scenario: Selecting Size, Color, & Fulfilling Store Branch](#52-scenario-selecting-size-color--fulfilling-store-branch)
   - 5.3 [Scenario: Adding to Cart vs. Instant "Buy Now"](#53-scenario-adding-to-cart-vs-instant-buy-now)
6. [Custom Board Simulator (Interactive 2D Real-Time Designer)](#6-custom-board-simulator-interactive-2d-real-time-designer)
   - 6.1 [Scenario: Launching the Custom Board Workspace & Canvas Overview](#61-scenario-launching-the-custom-board-workspace--canvas-overview)
   - 6.2 [Scenario: Top Navbar Zoom, Scale, & Safe Leave Protection](#62-scenario-top-navbar-zoom-scale--safe-leave-protection)
   - 6.3 [Scenario: Tool Tab 1 — "More" (Draft Save, Randomize, Reset)](#63-scenario-tool-tab-1--more-draft-save-randomize-reset)
   - 6.4 [Scenario: Tool Tab 2 — "Text" (Upper/Lower Typography & Colors)](#64-scenario-tool-tab-2--text-upperlower-typography--colors)
   - 6.5 [Scenario: Tool Tab 3 — "Image" (Logo/Photo Upload, Framing, & Crop)](#65-scenario-tool-tab-3--image-logophoto-upload-framing--crop)
   - 6.6 [Scenario: Tool Tab 4 — "Brush" (Interactive Flower Stamping & Editing)](#66-scenario-tool-tab-4--brush-interactive-flower-stamping--editing)
   - 6.7 [Scenario: Tool Tab 5 — "Border" (Styles, Thickness, & Color)](#67-scenario-tool-tab-5--border-styles-thickness--color)
   - 6.8 [Scenario: Tool Tab 6 — "Corner" (Corner Styles per Section)](#68-scenario-tool-tab-6--corner-corner-styles-per-section)
   - 6.9 [Scenario: Tool Tab 7 — "Floral" (Top Crest & Bottom Base Floral Accents)](#69-scenario-tool-tab-7--floral-top-crest--bottom-base-floral-accents)
   - 6.10 [Scenario: Physical Board Size Selection & Live Cost Breakdown](#610-scenario-physical-board-size-selection--live-cost-breakdown)
   - 6.11 [Scenario: Finalizing, High-Res Specs Review, & Adding to Cart](#611-scenario-finalizing-high-res-specs-review--adding-to-cart)
7. [Shopping Cart Management](#7-shopping-cart-management)
   - 7.1 [Scenario: Reviewing Multi-Branch Grouped Cart Items](#71-scenario-reviewing-multi-branch-grouped-cart-items)
   - 7.2 [Scenario: Transferring Store Branch per Item](#72-scenario-transferring-store-branch-per-item)
   - 7.3 [Scenario: Quantity Adjustments & Throttled Item Removal](#73-scenario-quantity-adjustments--throttled-item-removal)
   - 7.4 [Scenario: Applying Promo Codes & Calculating Total](#74-scenario-applying-promo-codes--calculating-total)
8. [Multi-Branch Checkout & Delivery Courier Selection](#8-multi-branch-checkout--delivery-courier-selection)
   - 8.1 [Scenario: Step 1 — Selecting Destination Address](#81-scenario-step-1--selecting-destination-address)
   - 8.2 [Scenario: Step 2 — Reviewing Items & Choosing Couriers per Store Branch](#82-scenario-step-2--reviewing-items--choosing-couriers-per-store-branch)
   - 8.3 [Scenario: Step 3 — Selecting Payment Gateway Channels](#83-scenario-step-3--selecting-payment-gateway-channels)
   - 8.4 [Scenario: Billing Summary & Confirming Order](#84-scenario-billing-summary--confirming-order)
9. [Payment Gateway & Real-Time Order Tracking](#9-payment-gateway--real-time-order-tracking)
   - 9.1 [Scenario: Payment Gateway Page, 24-Hour Timer, & QRIS/VA Instructions](#91-scenario-payment-gateway-page-24-hour-timer--qrisva-instructions)
   - 9.2 [Scenario: Automated & Manual Payment Verification Checks](#92-scenario-automated--manual-payment-verification-checks)
   - 9.3 [Scenario: Order Lifecycle Statuses](#93-scenario-order-lifecycle-statuses)
   - 9.4 [Scenario: Tracking Live Shipment & Contacting Courier / WhatsApp Support](#94-scenario-tracking-live-shipment--contacting-courier--whatsapp-support)
   - 9.5 [Scenario: Completing Order & Submitting Reviews](#95-scenario-completing-order--submitting-reviews)
10. [Legal Policies & Customer Support](#10-legal-policies--customer-support)
    - 10.1 [Scenario: Accessing Terms of Service & Privacy Policy](#101-scenario-accessing-terms-of-service--privacy-policy)
    - 10.2 [Scenario: Direct WhatsApp Customer Care](#102-scenario-direct-whatsapp-customer-care)

---

# 1. Introduction & Platform Overview

**Chia Florist** is a modern e-commerce web platform specializing in premium handcrafted flower boards (*papan bunga*) for weddings, grand openings, graduations, and condolences. In addition to ready-to-order catalog collections, Chia Florist features a real-time interactive 2D **Custom Board Simulator** that empowers customers to design customized flower boards with custom foam background colors, typography, logos/images, stamped flower arrangements, borders, and decorative crests directly from their web browser.

---

# 2. Global Layout & Navigation Elements

### 2.1 Top Navigation Bar (`CNavbar`)

The Navigation Bar is fixed at the top of every standard page with a frosted glass backdrop (`bg-white/95 backdrop-blur-md`).

#### Layout Components
* **Left Brand Logo:** Clickable Chia Florist logo (`/images/logo.png`) that navigates directly to the Home page (`/`).
* **Center Navigation Links (Desktop):**
  * `Home` (`/`)
  * `Catalog` (`/catalog`)
  * `Custom Board` (`/products/custom`)
* **Right Utility Actions:**
  * **Search Icon Button (🔍):** Opens the slide-out catalog search drawer.
  * **Cart Icon Button (🛒):** Displays a badge indicator showing current total cart items (e.g., `1`, `3`, or `99+`). Clicking opens the Shopping Cart page (`/cart`).
  * **User Profile / Sign In Area:**
    * *Unauthenticated:* Displays the pill button **`Sign In`** linking to `/login`.
    * *Authenticated:* Displays the user's avatar image / icon, name, and online indicator dot. Hovering or clicking reveals the **Account Dropdown Menu** containing:
      * Customer Name & Email Address.
      * **Primary Address preview:** Displays recipient address snippet.
      * Link to **`My Profile`** (`/profile`).
      * Link to **`My Orders`** (`/profile` with orders tab).
      * Button **`Sign Out`** (logs out and clears active auth session).
* **Mobile Hamburger Menu (☰):** Toggles a slide-down mobile panel containing quick search, navigation links, user profile snapshot, address details, and logout action.

---

### 2.2 Global Footer (`CFooter`)

The Footer is positioned at the bottom of standard pages, separated by a traditional Indonesian **Batik Kawung Motif** top border.

#### Layout Components
* **Batik Kawung Top Ribbon:** Aesthetic decorative green batik pattern ribbon.
* **Brand & Branch Info:**
  * Chia Florist Logo.
  * **Main Branch Address:** *Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng, Kab. Bangka Barat, Kep. Bangka Belitung, Indonesia 33311*.
  * **WhatsApp Quick Chat Box:** Message prompt with a **`Tanya via WhatsApp`** pill button opening direct WhatsApp chat (`+62 817-5234-999`).
* **Navigation Columns:**
  * **Navigation:** `Home`, `Catalog`, `Search`.
  * **Legal:** `Terms & Conditions` (`/terms`), `Privacy Policy` (`/privacy`).
  * **Connect:** `Instagram`, `WhatsApp`.
* **Copyright Bar:** `Copyright © [Year] Chia Florist. All rights reserved. Handcrafted with care.`

---

# 3. User Authentication & Account Management

### 3.1 Scenario: User Registration (Sign Up)

#### Flow & Visual Steps

```
[ Navbar: Sign In ] ──> [ Register Page: Step 1 (Initial Method) ]
                               │
            ┌──────────────────┴──────────────────┐
            ▼                                     ▼
 [ Continue with Google ]             [ Enter Email + Submit ]
            │                                     │
            ▼                                     ▼
 [ Google SSO Redirect ]               [ Step 2: Account Form ]
                                                  │
                                                  ▼
                                       [ OTP Verification Modal ]
```

1. **Navigate to Register:**
   * On the Navigation Bar, click **`Sign In`** -> Click the link **`Create an account`** at the bottom, OR navigate directly to `/register`.
2. **Step 1: Choose Registration Method:**
   * **Option A (Google SSO):** Click the green button **`Continue with Google`**. The browser redirects to Google OAuth consent and returns to the home page automatically authenticated.
   * **Option B (Email Form):** Type your email into the email input field and click the green arrow button **`→`** (or press `Enter`).
3. **Step 2: Fill Account Details:**
   * The page smoothly slides into the account details form showing `Registering for: [your-email]` (with a `(Change)` button if you want to switch email).
   * **Full Name:** Type your full name (e.g., `Rayhan Pratama`).
   * **Username:** Type your preferred username (e.g., `rayhan_p`).
   * **Password:** Enter a secure password. Click **`Show`** / **`Hide`** on the right side of the input to toggle password masking.
   * **Phone Number (Optional):** Enter your phone number (e.g., `08123456789`).
   * Click the primary green button **`Create Account`**.

---

### 3.2 Scenario: Account OTP Verification

#### Flow & Visual Steps

1. **Verify Account Prompt:**
   * Upon submitting registration details, the form transitions to the **`Verify Account`** screen showing: `We sent a verification code to: [your-email]`.
2. **Input OTP:**
   * Check your inbox for the 6-digit confirmation code.
   * Type the 6 digits into the large monospaced code field (e.g., `482910`).
3. **Submit Verification:**
   * Click the green button **`Confirm`** (or click **`Back`** to re-edit registration details).
   * On success, a toast confirmation appears and you are automatically logged in and redirected to the Home page.

---

### 3.3 Scenario: User Sign In (Login)

#### Flow & Visual Steps

1. **Navigate to Sign In:**
   * Click **`Sign In`** in the top navigation bar or go to `/login`.
2. **Step 1: Choose Login Method:**
   * **Option A (Google SSO):** Click **`Continue with Google`**.
   * **Option B (Email):** Enter your registered email address and click the green arrow button **`→`**.
3. **Step 2: Enter Password:**
   * The form transitions to the credentials screen with the password input automatically focused and selected.
   * Enter your password. Toggle **`Show`** / **`Hide`** if needed.
   * *(Optional)* Check the **`Remember me`** checkbox to persist your login session.
   * Click the green button **`Sign In`**.
   * On successful validation, you are redirected to the homepage or your previous intended checkout page.

---

### 3.4 Scenario: Password Reset / Recovery

#### Flow & Visual Steps

1. **Trigger Password Reset:**
   * On the login credentials screen (`/login`), click the link **`Forgot password?`** located above the password input.
2. **Request Reset Code:**
   * In the `Reset password` view, type your registered email address.
   * Click **`Send Reset Code`**.
3. **Verify Reset Code:**
   * In the `Verify code` screen, enter the 6-digit OTP received via email.
   * Click **`Verify Code`**.
4. **Set New Password:**
   * In the `Set new password` screen, type your new password.
   * Click **`Reset Password`**.
   * You will receive a success confirmation and return to the login screen ready to sign in with your new credentials.

---

### 3.5 Scenario: Profile Management & 2D Avatar Cropper

#### Flow & Visual Steps

1. **Navigate to Profile:**
   * Click your user avatar on the navbar -> click **`My Profile`** (or browse to `/profile`).
2. **Update Personal Information:**
   * In the **`Personal Information`** tab, you can view your Username and Email (read-only system fields).
   * Edit your **`Full Name`** and **`Phone Number`**.
   * Click **`Save Changes`** (or submit form). A success alert notification confirms the update.
3. **Change Profile Picture (Avatar):**
   * Hover over the circular avatar profile image -> Click **`Change Photo`** (or click **`Upload New Photo`**).
   * Select a picture from your device (`.jpg`, `.png`, `.webp`, max 10MB).
4. **Interactive 2D Cropper Modal:**
   * A dedicated cropper modal opens displaying your uploaded picture.
   * **Zoom Slider:** Move the zoom slider left/right to scale the image.
   * **Pan:** Click and drag the background image to position it inside the crop frame.
   * **Crop Handles:** Drag the corners/edges (`nw`, `ne`, `se`, `sw`, `n`, `s`, `e`, `w`) to resize the square crop box.
   * Click **`Crop & Save Avatar`** (or **`Cancel`** to abort).
   * The cropped image is saved to Supabase storage and your profile picture is instantly updated across the navbar and profile.
5. **Remove Photo:**
   * Click the outline button **`Remove Photo`** and confirm the browser alert to revert to the default avatar icon.

---

### 3.6 Scenario: Shipping Address Book Management

#### Flow & Visual Steps

1. **Navigate to Shipping Addresses:**
   * On the Profile page (`/profile`), click the sidebar/tab **`Shipping Addresses`**.
2. **View Existing Addresses:**
   * Saved addresses are listed with Recipient Name, Phone Number, Full Street Address, Postal Code, and a green **`Default`** badge for the primary destination.
3. **Add New Address:**
   * Click the **`+ Add New Address`** button.
   * In the address modal, fill out the chained location selections:
     1. **Recipient Name:** (e.g., `Ahmad Rizky`).
     2. **Phone Number:** (e.g., `08198765432`).
     3. **Province:** Select from dropdown (e.g., `Kepulauan Bangka Belitung`).
     4. **City / Regency:** Select from dynamically filtered dropdown (e.g., `Kab. Bangka Barat`).
     5. **District:** Select district (e.g., `Muntok`).
     6. **Village / Subdistrict:** Select village (e.g., `Sungai Daeng`).
     7. **Postal Code:** Type postal code (e.g., `33311`).
     8. **Full Street Address:** Type complete street, building, or venue details.
     9. **Default Checkbox:** Check `Set as default shipping address` if desired.
   * Click **`Save Address`**.
4. **Edit / Delete Address:**
   * Click **`Edit`** on an address card to open the form pre-filled with existing data.
   * Click **`Delete`** and confirm to remove an obsolete address.

---

# 4. Product Discovery & Catalog Browsing

### 4.1 Scenario: Exploring Homepage & Quick Action Shortcuts

#### Flow & Visual Steps

```
[ Homepage: Hero Carousel ] ──> Slide 1: "Lihat Katalog" (/catalog)
                            ──> Slide 2: "Coba Simulator" (/products/custom)
                            ──> Slide 3: "Pesan Sekarang" (/catalog)

[ Responsive Shortcuts ] ──> Button 1: "Katalog Papan Bunga"
                         ──> Button 2: "Custom Board Simulator"

[ Featured Products Grid ] ──> Click Card / "Lihat Detail" ──> (/products/[slug])
                           ──> Click Simulator Card ──> (/products/custom)
```

1. **Browse Hero Banner Carousel:**
   * The top hero banner automatically rotates through promotional slides every 5 seconds (pauses when hovering).
   * Use the **`<` (Previous)** and **`>` (Next)** arrow buttons or click the bottom **pagination dots** to switch slides manually.
   * Click the primary CTA button inside the active slide:
     * Slide 1: **`Lihat Katalog`** -> Opens `/catalog`.
     * Slide 2: **`Coba Simulator`** -> Launches `/products/custom`.
     * Slide 3: **`Pesan Sekarang`** -> Opens `/catalog`.
2. **Click Responsive Shortcut Buttons:**
   * Below the banner, click either:
     * **`Katalog Papan Bunga`** (Light gray card): Browse all pre-made floral arrangements.
     * **`Custom Board Simulator`** (Light green pastel card): Jump straight into the 2D custom board creator.
3. **Explore Featured Products Grid:**
   * Scroll down to `Koleksi Produk Pilihan`.
   * Each product card showcases:
     * High-resolution flower board photograph.
     * Rating score (e.g., `⭐ 4.8 | (120)`).
     * Product name and starting price in Indonesian Rupiah (e.g., `Rp 650.000`).
     * Status Badges: `Interactive 2D` (for custom simulator), `Preview Only` (amber badge), or `Sold Out` (red badge).
   * Click any card to navigate directly to its product detail page or custom designer.
   * Click **`Lihat Semua Produk di Katalog →`** at the bottom to view the complete catalog.

---

### 4.2 Scenario: Full Catalog Browsing with Multi-Store & Sorting Filters

#### Flow & Visual Steps

1. **Open Catalog:**
   * Navigate to `/catalog` via the navbar or homepage links.
2. **Search Flower Boards:**
   * Type into the search input box (`Search flower boards...`). Results filter live after a 400ms typing pause.
3. **Filter by Store Branch:**
   * Click the **`Store`** dropdown:
     * Select `All Stores` to view products across all regions.
     * Select a specific branch (e.g., `Chia Florist Muntok` or `Chia Florist Pangkalpinang`) to filter stock by that location.
   * When a specific store is active, a green banner appears: `Showing collection available at [Store Name]`. Click **`View All Stores`** on the banner to reset.
4. **Sort Products:**
   * Click the **`Sort By`** dropdown and pick your preferred ordering:
     * `Newest First` / `Oldest First`
     * `Name (A-Z)` / `Name (Z-A)`
     * `Price (Low to High)` / `Price (High to Low)`
     * `Stock (High to Low)` / `Stock (Low to High)`
     * `Weight (High to Low)` / `Weight (Low to High)`
5. **View Product Information & Stock Badges:**
   * Each catalog card displays:
     * Product category tag (e.g., `Papan Bunga`, `Buket`, `Interactive Game`).
     * Stock level badge (e.g., `📦 5 in stock` or `Sold Out`).
     * Description snippet.
     * Price in Rupiah (`Starting From Rp ...`).
     * Click **`View Details`** or **`Launch Game`** to proceed.

---

### 4.3 Scenario: Quick Search Drawer & Dedicated Search Page

#### Flow & Visual Steps

1. **Option A: Quick Slide-out Search Drawer:**
   * Click the **Search Icon (🔍)** on the top navigation bar.
   * A side drawer slides out from the right with the search field focused.
   * Type search keywords (e.g., `Pernikahan`, `Duka Cita`, `Wisuda`, `Rose`).
   * Matching product cards display with thumbnails, titles, and starting prices.
   * Click any product to jump to its page, OR press **`Enter ↵`** / click **`View all results →`** to open the full search page.
2. **Option B: Dedicated Search Page (`/search`):**
   * Access directly or via query (e.g., `/search?q=pernikahan`).
   * Features breadcrumb navigation (`Home / Search / "Pernikahan"`).
   * Use the top search bar, store branch filter, and sort selector to refine results.
   * If no items match, an empty state displays with suggestions and quick buttons: **`Launch Simulator`** and **`Browse Full Catalog`**.

---

# 5. Standard Product Ordering Flow

### 5.1 Scenario: Product Detail Page & Dynamic Price Calculation

#### Flow & Visual Steps

1. **Open Product Detail:**
   * Click any standard flower board from the catalog (e.g., `/products/papan-bunga-wedding-rose`).
2. **Inspect Photos:**
   * On the left column, click any of the **thumbnail previews** to switch the large high-resolution main photograph.
3. **Review Specifications:**
   * Title, Star Rating, Reviews count.
   * Availability indicator: `📦 In Stock (X available)`, `Sold Out at Selected Store`, or `Preview Only — Not For Sale`.
   * Detailed floral arrangement description, flower types, and build materials.

---

### 5.2 Scenario: Selecting Size, Color, & Fulfilling Store Branch

#### Flow & Visual Steps

1. **Select Color Variant:**
   * In the **`Colours:`** section, click on the colored circular buttons (e.g., Red, Pink, Blue, Gold, White) to highlight your chosen theme.
2. **Select Board Physical Size & Real-Time Price Recalculation:**
   * In the **`Size:`** section, click one of the size chips:
     * **`1.5m`:** Deducts **Rp 20.000** from base price (Compact board).
     * **`1.8m`:** Standard base price (Standard recommended size).
     * **`2m`:** Adds **Rp 30.000** to base price (Jumbo grand size).
   * The displayed price at the top updates instantly in real time.
3. **Select Fulfilling Branch (Mandatory):**
   * Under **`Fulfilling Branch:`**, a list of available in-stock stores appears (e.g., `Chia Florist Muntok (5 in stock)`).
   * Click the store card you wish to fulfill and deliver your order. The selected branch is highlighted with a green border and checkmark (`✓`).
   * *Note:* If no branch is selected, a warning banner alerts: `⚠️ Please select a fulfilling store branch above before adding to cart.`

---

### 5.3 Scenario: Adding to Cart vs. Instant "Buy Now"

#### Flow & Visual Steps

1. **Adjust Quantity:**
   * Click the **`-`** and **`+`** buttons on the quantity stepper (default is `1`).
2. **Action A: Add to Cart:**
   * Click the outline green button **`Add to Cart`**.
   * A green global alert toast pops up: `Added to Cart — [Product Name] (Qty: X) has been added to your shopping cart.`
   * Click **`View Cart`** to proceed to the cart, or click **`Continue`** to keep shopping.
3. **Action B: Buy Now (Direct Checkout):**
   * Click the solid green button **`Buy Now`**.
   * You are routed straight to `/checkout?buyNow=true&...` pre-populated with your chosen size, color, branch, and quantity, skipping the cart screen.

---

# 6. Custom Board Simulator (Interactive 2D Real-Time Designer)

### 6.1 Scenario: Launching the Custom Board Workspace & Canvas Overview

#### Flow & Visual Steps

```
[ Enter Simulator: /products/custom ]
                 │
 ┌───────────────┴────────────────┐
 ▼                                ▼
[ Top Navbar ]                  [ 2D Canvas Workspace ]
- Back to Home (Safe Guard)     - Upper Board (Header)
- Title & Version               - Lower Board (Body / Message)
- Zoom Controls (-, slider, +)  - Live Stamped Flowers / Images
                 │
 ┌───────────────┴────────────────┐
 ▼                                ▼
[ 7 Tool Tabs (Bottom Bar) ]    [ Floating Summary Bar ]
- More | Text | Image | Brush   - Element Count
- Border | Corner | Floral      - Live Total Price
                                - Finalize & Order Button
```

1. **Launch Simulator:**
   * Navigate to `/products/custom` from the navbar, homepage shortcut, or catalog simulator card.
2. **Workspace Layout:**
   * **Canvas Board (Center):** The large 2D interactive flower board divided into:
     * **Upper Section:** Top foam board header area.
     * **Lower Section:** Bottom foam board message/body area.
     * **Frame Border & Crests:** Surrounding wooden stand/frame and optional top/bottom floral crests.
   * **Floating Summary Bar (Bottom Canvas):** Shows total placed elements count, calculated live price (e.g., `Rp 234.000`), and the **`Finalize & Order`** button.
   * **Bottom Tool Tabs Bar:** Contains 7 interactive tool tabs: `More`, `Text`, `Image`, `Brush`, `Border`, `Corner`, `Floral`.
   * **Right Tool Settings Drawer:** Opens with sliders and color pickers when any tab is selected.

---

### 6.2 Scenario: Top Navbar Zoom, Scale, & Safe Leave Protection

#### Flow & Visual Steps

1. **Zoom Controls:**
   * Click **`−`** or press `Ctrl + -` to zoom out.
   * Click **`+`** or press `Ctrl + =` to zoom in.
   * Drag the horizontal **Zoom Slider** (hidden on mobile, visible on desktop).
   * Click the **Scale Chip** (e.g., `100%`) or press `Ctrl + 0` to reset zoom to default.
2. **Safe Leave Protection (Unsaved Drafts):**
   * If you made edits to your board and click the **`← Home`** back button, an alert modal opens: `Unsaved Progress Warning — Are you sure you want to leave? Your unsaved board progress will be lost.`
   * Click **`Save Draft & Leave`** (persists board state in browser storage and exits).
   * Click **`Leave Without Saving`** (discards unsaved modifications and exits).
   * Click **`Cancel`** (returns to editing your board).

---

### 6.3 Scenario: Tool Tab 1 — "More" (Draft Save, Randomize, Reset)

#### Flow & Visual Steps

1. Click the **`More` (•••)** icon button on the left of the bottom toolbar.
2. A popover menu opens with three options:
   * **`Save Progress`:** Saves your current board design, elements, colors, and layout as a draft in local storage. A confirmation toast displays: `Board draft saved!`.
   * **`Randomize Design`:** Instantly generates a surprise layout with randomized background palettes, font pairings, border ornate styles, and floral crest configurations.
   * **`Reset to Default`:** Reverts the entire board canvas back to the default pristine template.

---

### 6.4 Scenario: Tool Tab 2 — "Text" (Upper/Lower Typography & Colors)

#### Flow & Visual Steps

1. Click the **`Text` (T)** tab on the bottom toolbar. The right-side tool panel opens.
2. **Select Active Section:**
   * Click the toggle button **`Upper`** (to edit header) or **`Lower`** (to edit main message body).
3. **Customize Background:**
   * Click the color swatch input or select one of the preset color dots (e.g., Maroon `#c0392b`, Navy `#1a3a5c`, Emerald `#1b4332`, White `#ffffff`).
4. **Edit Header Text (Upper Section):**
   * **Text Input:** Type greeting text (e.g., `HAPPY WEDDING`, `SELAMAT & SUKSES`, `TURUT BERDUKA CITA`).
   * **Font Size Slider:** Adjust slider between `10px` and `96px`.
   * **Font Family Chips:** Choose from `Serif`, `Sans`, `Display`, `Script`, or `Gothic`.
   * **Alignment:** Click `Left`, `Center`, or `Right` alignment icon.
   * **Text Color:** Pick text color using the color picker.
   * **Underline Border:** Check `Underline Border` to enable a decorative underline, and adjust its width (`1px`–`12px`) and color.
5. **Edit Body Message (Lower Section):**
   * Switch toggle to **`Lower`**.
   * **Body Text Input:** Type recipient and sender names (e.g., `Romeo & Juliet\nDari: PT Maju Bersama`).
   * Adjust font size, font family, alignment, text color, and optional `Above-Body Border`.

---

### 6.5 Scenario: Tool Tab 3 — "Image" (Logo/Photo Upload, Framing, & Crop)

#### Flow & Visual Steps

1. Click the **`Image` (🖼️)** tab on the bottom toolbar.
2. **Upload Photo / Logo:**
   * Drag and drop an image file directly onto the canvas, OR click **`Browse Files`** / **`Upload Another Image`**.
3. **Configure Image Styling:**
   * **Frame Style:** Select `None` (raw edges), `Square` (rounded rectangular frame), or `Circle` (circular profile frame).
   * **Width Slider:** Adjust image width scaling on board (`5%` to `80%`).
   * **Zoom & Crop:**
     * **Zoom:** Adjust slider from `1.0×` to `3.0×` inside the frame.
     * **Crop X / Crop Y:** Adjust sliders (`0%` to `100%`) to reposition the image within the frame.
4. **Layer Management:**
   * Under `ALL IMAGES`, click an image in the list to select it and bring it to front.
   * Click the red **`Remove Image`** or list **`×`** button to delete an image layer.

---

### 6.6 Scenario: Tool Tab 4 — "Brush" (Interactive Flower Stamping & Editing)

#### Flow & Visual Steps

1. Click the **`Brush` (🖌️)** tab on the bottom toolbar.
2. **Choose Flower Type:**
   * Click **`Flower`** (5-petal standard daisy) or **`Rose`** (layered rosette blossom).
3. **Set Flower Color & Size:**
   * Select a color from preset swatches or custom color wheel.
   * **Size Slider:** Adjust diameter from `16px` to `120px`.
   * **Angle Slider:** Rotate stamp from `0°` to `360°`.
4. **Stamp on Canvas:**
   * Click anywhere on the board canvas to place a flower stamp.
   * Each stamped flower adds a nominal fee (e.g., **Rp 2.000**) to the live price calculation.
5. **Interactive Live Editing:**
   * Click and drag any placed flower directly on the canvas to move it.
   * Select a placed flower to edit its color, size, or rotation in real time.
   * Click **`×`** in the `PLACED` list or press `Delete` on keyboard to remove selected stamp.

---

### 6.7 Scenario: Tool Tab 5 — "Border" (Styles, Thickness, & Color)

#### Flow & Visual Steps

1. Click the **`Border` (▢)** tab on the bottom toolbar.
2. **Select Border Style:**
   * Choose from styles: `Solid`, `Dashed`, `Dotted`, `Double`, `Groove`, `Ridge`, `Inset`, `Outset`, `Ornate`, or `None`.
3. **Set Color & Width:**
   * Pick border color from palette.
   * Adjust **Width Slider** from `0px` to `32px`.
4. **Center Divider:**
   * Check **`SHOW CENTER BORDER`** to add a divider line between the upper header and lower body sections.

---

### 6.8 Scenario: Tool Tab 6 — "Corner" (Corner Styles per Section)

#### Flow & Visual Steps

1. Click the **`Corner` (◰)** tab on the bottom toolbar.
2. Switch section toggle between **`Upper`** and **`Lower`**.
3. Choose corner cut style:
   * `Square` (Classic 90-degree corners).
   * `Rounded` (Smooth rounded borders).
   * `Ornate` (Elegant scalloped baroque corners).
   * `Cut` (Chamfered angled corners).
   * `Floral` (Corner flower motifs 🌸).

---

### 6.9 Scenario: Tool Tab 7 — "Floral" (Top Crest & Bottom Base Floral Accents)

#### Flow & Visual Steps

1. Click the **`Floral` (🌺)** tab on the bottom toolbar.
2. Switch section toggle between **`Top Crest`** and **`Bottom Base`**.
3. Check **`ENABLE FLORAL DECOR`**.
4. Configure floral styling:
   * **Style:** Choose `Classic`, `Modern`, or `Grand`.
   * **Size Scaling:** Adjust slider (`20%` to `100%`).
   * **Primary Flower Color:** Pick dominant petal color.
   * **Secondary Flower Color:** Pick accent flower color.

---

### 6.10 Scenario: Physical Board Size Selection & Live Cost Breakdown

#### Flow & Visual Steps

1. In the bottom tool panel footer, inspect **`BOARD PHYSICAL SIZE`**:
   * **Compact (1.5m × 1.0m):** `Rp 150.000` — Best for small greetings & intimate venues.
   * **Medium (1.8m × 1.2m):** `Rp 200.000` `[Best]` — Standard full-size banquet board.
   * **Grand (2.0m × 1.5m):** `Rp 280.000` — Large venue inaugurations & luxury weddings.
   * **Royal (2.5m × 1.8m):** `Rp 380.000` — Premium VIP multi-panel presentation.
2. **Inspect Live Price Breakdown Card:**
   * `Base Board Size Price`: (e.g., `Rp 200.000`)
   * `Stamped Flowers`: `+ Rp X.000` (`Count × Rp 2.000`)
   * `Palette Expansion`: `+ Rp X.000` (for extended multi-color palettes)
   * `Decorations & Ornate Borders`: `+ Rp X.000`
   * `Image Components`: `+ Rp X.000` (`Count × Rp 20.000`)
   * **`Estimated Total Price`**: Dynamically summed in Indonesian Rupiah.

---

### 6.11 Scenario: Finalizing, High-Res Specs Review, & Adding to Cart

#### Flow & Visual Steps

```
[ Canvas Summary Bar: "Finalize & Order" ]
                    │
                    ▼
   [ Finalize Choice Overlay Modal ]
   ├── "Keep Customizing" ──> (Returns to Canvas)
   └── "Review Specs & Add to Cart"
                    │
                    ▼
       [ Review Specs Modal ]
       - High-Res Snapshot Preview
       - Technical Specifications Summary
       - Total Calculated Price
                    │
                    ▼ Click "Confirm & Add to Cart"
       [ Thank You Overlay Modal ]
       ├── "Design Another Board"
       └── "Go to Cart & Checkout" ──> (/cart)
```

1. **Click Finalize:**
   * On the floating summary bar above the bottom toolbar, click **`Finalize & Order`**.
2. **Choice Overlay (`FinalizeChoiceOverlay`):**
   * Modal opens with prompt: `What would you like to do next?`
   * Click **`Review Specs & Add to Cart`** (or **`Keep Customizing`** to make more tweaks).
3. **Review Specifications (`ReviewModal`):**
   * The system captures a high-resolution snapshot thumbnail of your 2D canvas.
   * On the right column, review the technical design specs:
     * Physical Size (e.g., `Medium (1.8m x 1.2m)`).
     * Upper Section Header Text & Background color.
     * Lower Section Body Text & Background color.
     * Border Style, Thickness & Color.
     * Top Crest & Bottom Base Decoration configurations.
     * Elements count: `X image(s), Y stamped flower(s)`.
     * Total Calculated Price.
4. **Confirm & Add to Cart:**
   * Click the primary green button **`Confirm & Add to Cart — Rp ...`**.
   * A spinner loads briefly while the custom payload and snapshot thumbnail are registered into your cart.
5. **Success Confirmation (`ThankYouOverlay`):**
   * Modal displays: `🌸 Custom Board Added to Cart!`.
   * Click **`Go to Cart & Checkout`** to proceed to `/cart`, OR click **`Design Another Board`** to start a fresh canvas.

---

# 7. Shopping Cart Management

### 7.1 Scenario: Reviewing Multi-Branch Grouped Cart Items

#### Flow & Visual Steps

1. **Open Shopping Cart:**
   * Click the **Cart Icon (🛒)** on the navbar or navigate to `/cart`.
2. **Empty Cart State:**
   * If no items are in cart, an empty shopping bag illustration displays with message: `Your cart is empty`. Click **`Continue Shopping`** to return to the catalog.
3. **Populated Multi-Branch Cart:**
   * Items in the cart are automatically grouped under store branch cards:
     * Header: `🏪 Fulfilled by [Branch Name]` (e.g., `Chia Florist Muntok Branch`) and item count.
   * Each item row displays:
     * Product thumbnail image (standard flower photo or custom design snapshot).
     * Product Name.
     * Size and Color tags.
     * `✨ Custom Board` badge (for custom simulator designs).
     * Item Subtotal in Rupiah.

---

### 7.2 Scenario: Transferring Store Branch per Item

#### Flow & Visual Steps

1. If multiple store branches are active in the system, each cart item displays a **`Transfer Branch:`** dropdown.
2. Click the dropdown and select another available store branch (e.g., `Chia Florist Pangkalpinang`).
3. An animated spinner appears: `Transferring...`.
4. A green alert confirms: `Branch Transferred — Item transferred to [New Store Branch]!`. The item seamlessly moves into the new store's fulfillment group without needing to re-order.

---

### 7.3 Scenario: Quantity Adjustments & Throttled Item Removal

#### Flow & Visual Steps

1. **Adjust Quantity:**
   * **Step Buttons:** Click **`−`** to decrease (disabled at `1`) or **`+`** to increase (max `80`).
   * **Direct Keyboard Input:** Click directly on the numeric quantity input field, type your desired quantity (e.g., `5`), and press `Enter` or click outside. The subtotal recalculates instantly.
2. **Remove Item:**
   * Click the red text button **`Remove` (🗑️)**.
   * *Anti-spam protection:* All remove buttons are temporarily locked while the selected item is removed with a spinner indicator.
   * A success alert confirms item removal and the cart subtotal updates.

---

### 7.4 Scenario: Applying Promo Codes & Calculating Total

#### Flow & Visual Steps

1. **Order Summary Card (Right Sidebar):**
   * Review `Subtotal`, `Estimated Delivery` (e.g., `Rp 20.000`), and `Total Amount`.
2. **Apply Promo Code:**
   * In the `Do you have a promo code?` input field, type a valid discount coupon (e.g., `CHIAFLORIST`).
   * Click the black button **`Apply`**.
   * A discount line (e.g., `Promo Discount: -Rp 50.000`) is deducted from the total bill and a success alert displays.
3. **Proceed to Checkout:**
   * Click the primary green button **`Proceed to Checkout`** to navigate to `/checkout`.

---

# 8. Multi-Branch Checkout & Delivery Courier Selection

### 8.1 Scenario: Step 1 — Selecting Destination Address

#### Flow & Visual Steps

1. **Access Checkout:**
   * Navigate to `/checkout` from the Cart or via "Buy Now".
   * *Authentication Check:* If you are not signed in, a warning alert redirects you to `/login`. Upon signing in, you return to checkout.
2. **Step 1: Shipping Destination:**
   * All addresses registered in your address book are displayed as selectable radio cards.
   * Each card displays: Recipient Name, `Default` badge, Contact Phone, and Full Address.
   * Click your desired destination address card. It is highlighted with a green border and checkmark radio button.
   * If you need to add a new address, click **`Manage Addresses`** to open your profile settings.

---

### 8.2 Scenario: Step 2 — Reviewing Items & Choosing Couriers per Store Branch

#### Flow & Visual Steps

1. **Review Multi-Branch Items:**
   * Under **`2. Review Items & Courier`**, orders are split by their fulfilling store branches.
   * Review ordered flower boards, sizes, colors, and quantities.
2. **Select Delivery Courier per Store:**
   * Each branch card has its own **`Delivery Courier Service:`** dropdown.
   * Click the dropdown to select a courier option:
     * `JNE (REG) - Rp 25.000`
     * `J&T (Standard) - Rp 22.000`
     * `POS Indonesia (Kilat Khusus) - Rp 20.000`
     * `Chia Florist Internal Fleet Delivery - Rp 15.000`
   * When you select a courier, a subtle spinner indicates `Recalculating rates...` and the shipping cost in the billing sidebar updates in real time.

---

### 8.3 Scenario: Step 3 — Selecting Payment Gateway Channels

#### Flow & Visual Steps

1. Under **`3. Payment Gateway Channel`**, payment options are organized into collapsible accordion categories:
   * **QR Code / QRIS:** QRIS Instant (GoPay, OVO, Dana, LinkAja, BCA QR, ShopeePay).
   * **E-Wallet:** GoPay, OVO, ShopeePay, Dana.
   * **Bank Transfer / Virtual Account:** BCA Virtual Account, Mandiri VA, BNI VA, BRI VA, Permata VA.
2. Click an accordion header to expand its channels.
3. Click the radio card for your preferred payment method. Each channel displays its processing fee (e.g., `Free Fee` or `Fee: Rp 4.000`).

---

### 8.4 Scenario: Billing Summary & Confirming Order

#### Flow & Visual Steps

1. **Review Billing Summary (Right Column):**
   * `Items Subtotal`
   * `Shipping Cost` (summed across all store couriers)
   * `Payment Processing Fee`
   * `Promo Discount` (if applied)
   * **`Total Bill`** (prominently styled in bold green font).
2. **Submit Order:**
   * Click the full-width primary green button **`Confirm & Pay Now`**.
   * The button displays a processing spinner. The order is submitted to the backend orderService, items are converted into an active pending order, and you are automatically redirected to the Payment page (`/payment?orderId=[ORDER_ID]`).

---

# 9. Payment Gateway & Real-Time Order Tracking

### 9.1 Scenario: Payment Gateway Page, 24-Hour Timer, & QRIS/VA Instructions

#### Flow & Visual Steps

```
[ Payment Page: /payment?orderId=... ]
                   │
 ┌─────────────────┴─────────────────┐
 ▼                                   ▼
[ Order Total & 24h Countdown ]    [ Selected Provider Card ]
- Total Bill in Rupiah              - Action URL ("Proceed to Payment Page ↗")
- Red Timer Badge (23:59:59)        - Embedded QRIS QR Code (for Scan & Pay)
                   │
                   ▼
[ Formatted Markdown Instructions ]
- ATM / M-Banking / E-Wallet Step-by-Step Guide
                   │
 ┌─────────────────┴─────────────────┐
 ▼                                   ▼
[ "I Have Paid / Verify Status" ]   [ "Pay Later / View Orders" ]
```

1. **Inspect Payment Header:**
   * **Order Total:** Displays exact invoice total in Rupiah.
   * **Order ID:** Unique order tracking identifier (e.g., `ORD-20260820-0042`).
   * **Payment Time Left:** Active countdown timer starting at 24 hours (`23:59:59`).
2. **Pay via Embedded QR Code (QRIS / E-Wallet):**
   * If QRIS or E-Wallet was selected, an authentic scannable QR Code image is generated directly on the page.
   * Open your banking app or e-wallet (GoPay, OVO, Dana, BCA Mobile) -> Scan the QR Code -> Confirm payment.
3. **Pay via Payment Gateway Portal (Action URL):**
   * If a hosted checkout channel was chosen, click the green button **`Proceed to Payment Page ↗`** to open the secure payment portal.
4. **Follow Step-by-Step Payment Instructions:**
   * The page renders formatted markdown instructions explaining exact transfer steps for ATM, Mobile Banking, Internet Banking, and Minimarket counters.

---

### 9.2 Scenario: Automated & Manual Payment Verification Checks

#### Flow & Visual Steps

1. **Automated Background Polling:**
   * While you are on the payment page, the application automatically polls the payment verification API every 30 seconds.
   * Once your transfer is reconciled by the payment gateway, the page instantly transforms into the **Payment Successful** state without needing a manual refresh.
2. **Manual Instant Verification:**
   * After completing your bank or e-wallet transfer, click the green button **`I Have Paid / Verify Status`**.
   * An animated spinner displays `Verifying Payment...`.
   * On confirmation, a celebratory success toast pops up: `Payment Verified! — Thank you! Your payment has been received and confirmed.`

---

### 9.3 Scenario: Order Lifecycle Statuses

The application handles 6 distinct order states:

| Status Badge | State Description | Available User Actions |
| :--- | :--- | :--- |
| **`Pending`** | Awaiting customer payment within the 24-hour window. | `Pay Now`, `Check Payment Status`, `Cancel Order` |
| **`Processing`** | Payment confirmed; florist is handcrafting the flower board. | `Track Order`, `Chat with Florist on WhatsApp` |
| **`Shipping`** | Courier fleet dispatched with board out for delivery to venue. | `Track Shipment`, `Contact Delivery Driver` |
| **`Completed`** | Flower board successfully delivered and erected at event venue. | `Leave Review`, `Re-Order` |
| **`Expired`** | 24-hour payment window elapsed without transfer. | `Browse Catalog / Re-Order` |
| **`Cancelled`** | Order cancelled by customer or merchant. | `Re-Order Flower Arrangement` |

---

### 9.4 Scenario: Tracking Live Shipment & Contacting Courier / WhatsApp Support

#### Flow & Visual Steps

1. **Navigate to Order Tracking:**
   * Go to `/profile` -> Click the **`Order Tracking`** tab.
   * Filter orders by clicking the sub-tabs: `All Orders`, `To Pay`, `To Ship`, `To Receive`, `Completed`, or `Cancelled / Expired`.
2. **Open Shipment Tracking Modal:**
   * On an order with status `Processing` or `Shipping`, click **`Track Shipment`**.
   * A tracking overlay modal opens featuring:
     * **Courier & AWB:** Courier name and tracking number with a **`Copy` (📋)** button.
     * **Estimated Delivery:** Target delivery arrival time.
     * **Interactive Timeline Milestones:**
       * `✓ Order Placed`
       * `✓ Payment Verified`
       * `✓ Florist Crafting Arrangement`
       * `● Out for Delivery with Courier`
       * `○ Delivered to Venue`
3. **Contact Driver / Florist Support:**
   * Inside the tracking modal, click the green button **`Contact Driver / Support on WhatsApp`**.
   * Opens direct WhatsApp chat with pre-populated order inquiry text: `Hello Chia Florist, I would like to inquire about the delivery status for order [ORDER_ID]`.

---

### 9.5 Scenario: Completing Order & Submitting Reviews

#### Flow & Visual Steps

1. Once your flower board is delivered to the venue, its status transitions to **`Completed`**.
2. Click **`Leave Review`** on the completed order card.
3. Select star rating (`⭐ 1` to `⭐ 5`) and input feedback. A success notification confirms: `Review Submitted — Thank you! Your review has been recorded.`

---

# 10. Legal Policies & Customer Support

### 10.1 Scenario: Accessing Terms of Service & Privacy Policy

#### Flow & Visual Steps

1. **Terms & Conditions:**
   * Click **`Terms of service`** on the login/register forms or footer (`/terms`).
   * Read platform guidelines, user eligibility, merchant fulfillment responsibilities, and return/cancellation policies.
2. **Privacy Policy:**
   * Click **`Privacy Policy`** on the footer (`/privacy`).
   * Review data collection standards, encryption practices, account information security, and cookie policies.

---

### 10.2 Scenario: Direct WhatsApp Customer Care

#### Flow & Visual Steps

1. **Trigger WhatsApp Support:**
   * Click **`Tanya via WhatsApp`** in the website footer or click WhatsApp contact buttons inside order tracking.
2. **Instant Live Chat:**
   * Connects directly with Chia Florist customer service representatives on **`+62 817-5234-999`** for custom bouquet requests, special delivery venue coordination, or order status assistance.

---

*End of User Manual Guide — Chia Florist E-Commerce Platform*