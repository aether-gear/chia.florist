# Control Panel UI Standardization Audit

This document details the UI audit and gap analysis performed under **Milestone 1** of the shadcn/ui standardization plan. It outlines custom HTML elements, inline styling overrides bypassing theme tokens, dark mode compatibility issues, and the missing shadcn components required for full standardization.

---

## 1. 🔍 Detailed File Audit

The following table lists views and components under `src/pages/` containing native HTML elements or inline style overrides (such as custom Indigo colors or custom card classes) that bypass the centralized shadcn/ui theme variables.

| File / Location | Element / Tag | Description / Styling | Issue / Bypass | Recommendation |
| :--- | :--- | :--- | :--- | :--- |
| **`PlaceholderPage.tsx`** | Icons / Divs | `bg-indigo-100`, `text-indigo-600` | Hardcoded theme colors | Replace with theme variables or standard shadcn classes. |
| **`admin/AddMerchantAccountPage.tsx`** | Icon container / Button | `bg-indigo-100`, `text-indigo-600` (icon); `<Button className="bg-indigo-600 hover:bg-indigo-700">` | Custom primary button coloring; custom icon backgrounds | Refactor button to use standard `variant="default"`. |
| **`admin/AuditLogsPage.tsx`** | Links / Actions | `text-slate-600 hover:text-indigo-600`, `text-slate-400 hover:text-indigo-600` | Custom text state overrides | Map hover behaviors to theme-aligned variables. |
| **`admin/CreateMerchantPage.tsx`** | `<textarea>`, `<Button>` | `textarea` (line 75) custom inline classes; `<Button className="bg-indigo-600 hover:bg-indigo-700">` | Native textarea element; custom button overrides | Replace with `<Textarea>` component; refactor button style. |
| **`admin/payments/CreatePaymentAccountPage.tsx`** | `<select>` | `select` (line 96) custom inline classes inside `<FormControl>` | Native select dropdown | Replace with shadcn `<Select>` component primitives. |
| **`admin/payments/PaymentSettingsPage.tsx`** | Table cell | `text-indigo-600 dark:text-indigo-400` (line 167) | Hardcoded branding colors | Use semantic classes or standard theme variables. |
| **`auth/ForgotPasswordPage.tsx`** | `<Button>` | Multiple occurrences of `bg-indigo-600 hover:bg-indigo-700` | Custom primary button coloring | Refactor to use the standard default primary variant. |
| **`auth/LoginPage.tsx`** | `<input type="checkbox">`, `<Button>`, Error alert | `input` (line 95) native; `<Button className="bg-indigo-600 hover:bg-indigo-700">` (line 111); custom alert div | Native checkbox; hardcoded primary button colors; raw alert box | Replace checkbox with `<Checkbox>`; standardise button variant and error banners. |
| **`dashboard/DashboardPage.tsx`** | `<select>`, Divs | `select` (lines 129, 138) native; custom Suggestion Card `bg-indigo-50 border-indigo-100` | Native select dropdowns; custom card styling | Replace with `<Select>` component; standardise card with `<Card>`. |
| **`dashboard/ProductsPage.tsx`** | Alert | `alert(err.message || ...)` (line 60) | Native browser alert popup | Replace with shadcn `<Toast>` / `use-toast`. |
| **`merchant-profile/components/ProfileContactSection.tsx`** | `<textarea>` | `textarea` (line 79) native | Native textarea element | Replace with shadcn `<Textarea>` component. |
| **`merchant-profile/components/ProfileIdentitySection.tsx`** | `<textarea>` | `textarea` (line 44) native | Native textarea element | Replace with shadcn `<Textarea>` component. |
| **`merchant-profile/components/ProfileOperationalSection.tsx`** | `<input type="checkbox">` | `input` (lines 40, 72) native | Native checkboxes | Replace with shadcn `<Checkbox>` components. |
| **`merchant-profile/components/ProfileSettingsSection.tsx`** | `<select>` | `select` (lines 17, 27, 39) native | Native select dropdowns | Replace with shadcn `<Select>` components. |
| **`security/SecurityPage.tsx`** | `<button>`, `<Button>` | `button` (line 2099) native; custom button colors (`bg-slate-900`, `bg-emerald-600`, `bg-slate-700`) | Native button; custom button color overrides | Replace with standard `<Button>` variants or theme variables. |
| **`shop/ShopManagementPage.tsx`** | `<textarea>`, `<input type="checkbox">`, `<select>`, `<Card>` | Multiple native textareas, checkboxes, selects; raw card overrides (lines 323, 461, 513, etc.) | Large-scale bypass of shadcn elements and theme variables | Full form and container refactoring to use standard primitives. |

---

## 2. 🔌 Component Gaps Matrix

The following table displays custom UI wrapper components and native HTML elements identified during the audit, mapping them to the official **`shadcn/ui`** primitives that must be installed/composed.

| Current Implementation | Proposed shadcn/ui Component | Status / Location |
| :--- | :--- | :--- |
| **`src/components/Pagination.tsx`** (custom layout and buttons) | **`Pagination`** primitives (`<Pagination>`, `<PaginationContent>`, etc.) | Missing / To be installed |
| **`src/components/StatusBadge.tsx`** (custom slate/emerald hover overrides) | **`Badge`** component variants (`variant="default"`, `variant="secondary"`, etc.) | Installed / To be refactored |
| **`src/components/LoadingState.tsx`** (custom tailwind layout) | **`Skeleton`** component | Installed / To be composed |
| **`src/components/EmptyState.tsx`** (custom tailwind cards) | **`Card`** and **`Button`** primitives | Installed / To be composed |
| **`ProductsPage.tsx`** (native browser `alert()`) | **`Toast`** / **`use-toast`** | Missing / To be installed |
| Native `<input type="checkbox">` tags | **`Checkbox`** component | Installed / To be integrated |
| Native `<select>` tags | **`Select`** component primitives | Installed / To be integrated |
| Native `<textarea>` tags | **`Textarea`** component | Installed / To be integrated |
| Custom action popups or hover triggers | **`Tooltip`** component | Missing / To be installed |

---

## 3. 🌗 Theme Compatibility & Dark Mode Gaps

A critical finding of the audit is that **almost all view pages lack dark mode support**. Because they use hardcoded Tailwind utility colors (like `bg-white`, `bg-slate-50`, and `text-slate-600`) without specifying `dark:` counterparts or utilizing shadcn theme variables (like `bg-background`, `text-foreground`, and `border-border`), they will fail to adapt when switching between light and dark modes.

### ⚠️ High-Risk Files (Missing Dark Mode Support):
1. [PlaceholderPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/PlaceholderPage.tsx)
2. [AuditLogsPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/admin/AuditLogsPage.tsx)
3. [CreateMerchantPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/admin/CreateMerchantPage.tsx)
4. [ForgotPasswordPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/auth/ForgotPasswordPage.tsx)
5. [LoginPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/auth/LoginPage.tsx)
6. [DashboardPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/dashboard/DashboardPage.tsx)
7. [ProductsPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/dashboard/ProductsPage.tsx)
8. [ProfileContactSection.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/merchant-profile/components/ProfileContactSection.tsx)
9. [ProfileFinancialSection.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/merchant-profile/components/ProfileFinancialSection.tsx)
10. [ProfileIdentitySection.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/merchant-profile/components/ProfileIdentitySection.tsx)
11. [ProfileOperationalSection.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/merchant-profile/components/ProfileOperationalSection.tsx)
12. [ProfileSettingsSection.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/merchant-profile/components/ProfileSettingsSection.tsx)
13. [OrdersPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/orders/OrdersPage.tsx)
14. [ShopManagementPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/shop/ShopManagementPage.tsx)

### ✅ Partially Compatible Files (Has Dark Mode Classes):
* [SecurityPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/security/SecurityPage.tsx) - Utilizes some `dark:` class overrides, but still contains custom color buttons that need variables.

---

## 4. 🚀 Recommended Action Plan (Milestone 2 - 4)

1. **Milestone 2**:
   * Run the CLI to install `pagination`, `tooltip`, and `toast`:
     ```bash
     npx shadcn@latest add pagination tooltip toast
     ```
2. **Milestone 3**:
   * Refactor custom pagination, status badges, and loading states.
3. **Milestone 4**:
   * Standardize form controls (inputs, selects, textareas, checkboxes) and buttons inside all page views, while simultaneously resolving all dark-mode and theme variable gaps.
