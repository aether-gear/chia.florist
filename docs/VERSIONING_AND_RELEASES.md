# Monorepo CI Versioning & Release Guidelines

Chia Florist uses **automated semantic versioning and release management** powered by **Google Release Please** and **GitHub Actions**.

Versioning is **independent per module** — modifying one service will only trigger version bumps, changelog entries, and Docker image tags for that specific module.

---

## 1. Commit Trigger Words Cheatsheet

Release Please inspects the prefix and footer of your commit messages on `main` to determine the version bump:

| Bump Type | Version Example | Commit Trigger Prefix | Description & Examples |
| :--- | :--- | :--- | :--- |
| **PATCH** | `0.14.0` → `0.14.1` | `fix:`<br>`perf:` | **Bug fixes & Performance enhancements**<br>• `fix: resolve custom product leak to cart`<br>• `fix(auth): isolate session state on logout`<br>• `perf(db): add index on order_item_custom_designs` |
| **MINOR** | `0.14.0` → `0.15.0` | `feat:` | **New backward-compatible features**<br>• `feat: introduce multi-courier tracking`<br>• `feat(orders): implement automatic refund on cancel`<br>• `feat(ai): add live SKU forecaster` |
| **MAJOR** | `1.0.0` → `2.0.0`<br>*(or `0.x` pre-major)* | `feat!:`<br>`fix!:`<br>`BREAKING CHANGE:` | **Breaking changes / incompatible API changes**<br>• `feat!: restructure order shipment schema`<br>• Or add footer in commit body:<br>`BREAKING CHANGE: domain events now require context` |
| **NO BUMP** | *(No release cut)* | `chore:`<br>`docs:`<br>`style:`<br>`test:`<br>`ci:`<br>`refactor:` *(non-breaking)* | **Internal maintenance, docs, refactors & tests**<br>• `docs: update Staff API documentation`<br>• `test: add unit test for checkout calculation`<br>• `chore: update dependencies`<br>• `ci: update workflow runners` |

> [!NOTE]
> **Pre-1.0 (`0.x.y`) Mode**: Because `"bump-minor-pre-major": true` is enabled, breaking changes during pre-1.0 development bump the **minor** version (`0.14.0` → `0.15.0`) to avoid inflating major version numbers before domain contracts stabilize for `1.0.0`.

---

## 2. How Monorepo Detection & Overlapping Commits Work

> [!IMPORTANT]
> **Release Please tracks modified FILE PATHS, NOT the commit scope name.**
> You do **not** need to write `feat(service-core):` or `fix(e-commerce):`. You should write **domain-first scopes** like `feat(order):`, `fix(auth):`, `feat(custom-product):`, `feat(shipping):`.

### A. Single-Module Commits (Domain Scopes)
Write the commit based on the feature/domain. Release Please checks which directory the changed files live in:

```bash
# Modifies files in service-core/ -> Bumps service-core (0.14.0 -> 0.15.0)
git commit -m "feat(order): implement automatic payment refund on cancellation"
git commit -m "fix(inventory): exclude deleted products from shop inventory"

# Modifies files in e-commerce/ -> Bumps e-commerce (0.8.0 -> 0.8.1)
git commit -m "fix(cart): resolve custom product retention and post-checkout cleanup"
git commit -m "feat(seo): migrate critical data fetching from SPA to Nuxt SSR"

# Modifies files in control-panel/ -> Bumps control-panel (0.5.0 -> 0.5.1)
git commit -m "fix(staff): resolve assigned-shop switcher state on reload"
git commit -m "feat(order-management): add date range filtering and shop isolation"

# Modifies files in intelligence-layer/ -> Bumps intelligence-layer (0.3.0 -> 0.4.0)
git commit -m "feat(demand-forecast): introduce courier SLA and stockout scanner"
```

### B. Cross-Module / Overlapping Commits (Full-Stack Features)
When a feature touches files across multiple directories in one commit or PR, **Release Please automatically detects all affected modules and bumps each one independently**:

```bash
# Example 1: Full-Stack Custom Product feature
# Touches:
#   - service-core/internal/modules/cart/...
#   - e-commerce/app/features/custom-product/...
git commit -m "feat(custom-product): allow customer-selected fulfillment shop"
# -> Result: Bumps BOTH service-core (0.14.0 -> 0.15.0) AND e-commerce (0.8.0 -> 0.9.0)
# -> Adds changelog entry to both CHANGELOG.md files automatically.

# Example 2: Backend Rule + Admin UI Panel
# Touches:
#   - service-core/internal/modules/staff/...
#   - control-panel/src/features/staff/...
git commit -m "feat(staff-rules): enforce shop assignment permissions and status rules"
# -> Result: Bumps BOTH service-core AND control-panel.

# Example 3: AI Inference Integration
# Touches:
#   - intelligence-layer/app/...
#   - service-core/internal/infra/genai/...
git commit -m "feat(ai-designer): integrate chat-to-generate flower board with LLM"
# -> Result: Bumps BOTH intelligence-layer AND service-core.
```

### C. Recommended Domain Scopes Cheatsheet

| Domain Scope | Description / Functional Area | Typical Modules Touched |
| :--- | :--- | :--- |
| `order` / `checkout` | Order workflows, cancellations, refunds, quotes | `service-core`, `e-commerce` |
| `shipping` / `courier` | Courier dispatch, tracking sync, tracking logs | `service-core`, `control-panel` |
| `custom-product` / `designer` | DIY board canvas, v3 schema, flower brushes | `e-commerce`, `service-core` |
| `auth` / `session` | Challenge accounts, login isolation, token refresh | `service-core`, `e-commerce`, `control-panel` |
| `staff` / `permissions` | Shop assignment, staff lifecycle, role rules | `service-core`, `control-panel` |
| `inventory` / `product` | Stockout, deleted_at lifecycle, product pricing | `service-core`, `e-commerce` |
| `ai` / `forecast` | SKU forecaster, SLA anomaly, ML models | `intelligence-layer`, `service-core` |
| `seo` / `ssr` | Nuxt SSR data fetching, view-level SEO metadata | `e-commerce` |
| `infra` / `db` | Database migrations, logging, transactor, WAF | `service-core`, `infra` |


---

## 3. Current Module Baseline Versions

| Module | Directory | Tag Format | Docker Image |
| :--- | :--- | :--- | :--- |
| **`service-core`** | `service-core/` | `service-core-v0.14.0` | `ghcr.io/aether-gear/chia.florist/service-core:0.14.0` |
| **`e-commerce`** | `e-commerce/` | `e-commerce-v0.8.0` | `ghcr.io/aether-gear/chia.florist/e-commerce:0.8.0` |
| **`control-panel`** | `control-panel/` | `control-panel-v0.5.0` | Static Netlify build |
| **`intelligence-layer`** | `intelligence-layer/` | `intelligence-layer-v0.3.0` | `ghcr.io/aether-gear/chia.florist/intelligence-layer:0.3.0` |

---

## 4. Custom Product Contract & Schema Versioning

The DIY & AI custom flower board designer follows strict schema and engine contract versioning:
* **Schema Contract Version**: `3.0.0` (Defined in [custom-product-v3.schema.json](file:///d:/__Projects/kage/chia.florist/docs/specs/custom-product-v3.schema.json))
* **Editor Engine Version**: `3.1.0` (Defined in `e-commerce/app/features/custom-product/constants.ts`)

Every custom design payload exported by the editor or produced by the AI generator carries:
```json
{
  "metadata": {
    "version": "3.0.0",
    "editorVersion": "3.1.0",
    "checksum": "811c9dc5..."
  }
}
```

---

## 5. How the Release Cycle Works

```
1. Developer pushes conventional commits to develop/staging/main
2. Release Please bot inspects commits and opens/updates a "Release PR"
3. Team reviews and merges the Release PR to main
4. GitHub Actions automatically:
   - Cuts the Git Tag (e.g. service-core-v0.15.0)
   - Generates CHANGELOG.md entry
   - Publishes GitHub Release notes
   - Triggers CI workflow to build & push versioned Docker image to GHCR
```
