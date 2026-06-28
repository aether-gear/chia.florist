# Remove All Mock Data — Wire Everything to Real APIs

Replace every static fallback / mock object in the control panel with proper API calls using `fetchApi`. When an API fails, show a real **error state** instead of silently swapping in fake data.

## User Review Required

> [!IMPORTANT]
> The current code uses mock data as a **silent fallback** whenever a backend endpoint fails. Removing mocks means that if the Go backend (`main.go`) does not have a working endpoint, the UI will show an **error state** instead of fabricated data.
> Please confirm that the backend API endpoints listed below are ready (or are ready to be implemented) before executing this plan.

> [!WARNING]
> **Dashboard "Sales Overview" chart** currently uses completely fictional week-revenue data (`mockSales.ts`). There is no existing backend endpoint for sales analytics. Please decide:
> - **Option A**: Wire it to a real `/analytics/sales` or `/orders/summary` backend endpoint (needs backend work).
> - **Option B**: Remove the Sales Overview chart from the dashboard entirely until a real endpoint exists.
> - **Option C**: Keep the chart but clearly label it as "Demo / Coming Soon" with a visual indicator.

## Open Questions

> [!IMPORTANT]
> **Sales data source** — `DashboardPage.tsx` uses `salesData` and `topSellingProducts` from `mockSales.ts` for the bar chart and AI Insights card. What real API should supply this? (e.g., `GET /analytics/sales?period=week`)

> [!IMPORTANT]
> **Merchant Profile** — `useMerchantProfileViewModel.ts` calls `GET /api/core/merchants/profile`. Does this endpoint exist in `main.go`? Should it be merchant-specific (`/merchants/{id}/profile`) instead?

> [!IMPORTANT]
> **WAF Active Rules count** — `getWafSummary()` in `wafData.ts` hardcodes `activeRules: 5`. Should this come from `waf-rules.json` length, or a real `/waf/rules` API response?

## Proposed Changes

### Mock Data Files (to be deleted)

#### [DELETE] [mockSales.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/data/mockSales.ts)
- Contains `salesData` (weekly revenue) and `topSellingProducts` — fully fictional.
- Action depends on answer to the "Sales data source" open question above.

### ViewModels — Remove Inline Mock Objects & Silent Fallbacks

Each viewmodel currently catches errors and swaps in a `mock*` constant. These will be changed to **propagate the error** via `setError(err.message)` so the UI can render a proper error state.

#### [MODIFY] [useProductsViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useProductsViewModel.ts)
- Remove `mockProductsResponse` constant (lines 4–48).
- Switch raw `fetch('/api/core/products')` to `fetchApi('/products')` for consistency.
- On failure → `setError(err.message)`, do **not** fallback to mock.

#### [MODIFY] [useShopViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useShopViewModel.ts)
- Remove `mockShopId`, `mockAddresses`, `mockCouriers`, `mockProducts` (lines 5–75).
- Remove all `catch` blocks that silently set mock data (lines 104–123).
- On any failure → propagate a real `error` state.
- `createAddress` currently uses `shopId || mockShopId`; remove the `mockShopId` fallback.

#### [MODIFY] [usePaymentsViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/usePaymentsViewModel.ts)
- Remove `mockMethods` and `mockAccounts` constants (lines 5–77).
- On failure → `setError(err.message)`.

#### [MODIFY] [useMerchantsViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useMerchantsViewModel.ts)
- Remove `mockMerchantsResponse` constant (lines 5–35).
- On failure → `setError(err.message)`.

#### [MODIFY] [useCustomersViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useCustomersViewModel.ts)
- Remove `mockCustomersResponse` constant (lines 4–38).
- Switch raw `fetch('/api/core/customers')` to `fetchApi('/customers')`.
- On failure → `setError(err.message)`.

#### [MODIFY] [useMerchantProfileViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useMerchantProfileViewModel.ts)
- Remove `mockMerchantProfile` constant (lines 4–48).
- Switch raw `fetch('/api/core/merchants/profile')` to `fetchApi('/merchants/profile')`.
- On failure → `setError(err.message)`.

#### [MODIFY] [useAuthMeViewModel.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/viewmodels/useAuthMeViewModel.ts)
- Remove `mockAuthMeResponse` constant (lines 22–34).
- On failure → `setError(err.message)`, `setData(null)`.
- This one is sensitive: if `/auth/me` fails, the user is unauthenticated. The app's `ProtectedRoute` should handle the redirect.

### Dashboard — Remove Mock Sales Data

#### [MODIFY] [DashboardPage.tsx](file:///d:/__Projects/kage/chia.florist/control-panel/src/pages/dashboard/DashboardPage.tsx)
- Remove `import { salesData, topSellingProducts } from '@/data/mockSales'`.
- **Sales Overview chart**: Pending answer to open question — either wire to real API, remove the chart, or label as "Coming Soon".
- **AI Insights card**: The "Top Selling Product" insight currently reads from `topSellingProducts[0]`. Replace with either real data or generic static text that doesn't reference mock names.
- **WAF Active Rules**: Wire `activeRules` to `waf-rules.json` length (already imported in `wafData.ts`) rather than hardcoding `5`.

#### [MODIFY] [wafData.ts](file:///d:/__Projects/kage/chia.florist/control-panel/src/data/wafData.ts)
- `getWafSummary()` currently hardcodes `activeRules: 5` (line 59).
- Import `wafRules` from `../../waf-rules.json` and return `activeRules: wafRules.length` (or equivalent real count).

## Verification Plan

### Manual Verification
1. Start the Go backend (`go run main.go`) and the Vite dev server (`npm run dev`).
2. Navigate to each page and confirm:
   - **Products page**: shows real products or a clear error state (no mock fallback).
   - **Shop Management page**: shows real addresses/couriers/products or an error.
   - **Payment Accounts page**: shows real accounts or an error.
   - **Merchants / Customers pages**: show real data or errors.
   - **Merchant Profile page**: loads from API or shows error.
   - **Auth**: failed `/auth/me` redirects to login instead of logging in as mock merchant.
   - **Dashboard**: WAF stats from real logs; Sales chart either from real API or removed/labeled.
3. Deliberately shut down the backend and confirm all pages show **error states** (not mock data).
