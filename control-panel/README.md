# Chia Florist — Control Panel

Internal staff and administration dashboard for Chia Florist, providing order fulfillment workspaces, shop management, staff permissions, courier assignments, and operational analytics.

## Setup & Running Locally

```bash
# Install dependencies
npm ci

# Start development server
npm run dev

# Build for production
npm run build

# Run type checks and linter
npm run lint
npx tsc --noEmit
```

## Versioning & Commit Guidelines

`control-panel` releases are automated via GitHub Actions (Current baseline: `control-panel-v0.6.0`).
Production builds are deployed to Netlify on release.

| Commit Type | Version Bump | Example |
| :--- | :--- | :--- |
| `fix:` / `perf:` | **Patch** (`0.6.0` → `0.6.1`) | `fix(staff): resolve assigned shop switcher state` |
| `feat:` | **Minor** (`0.6.0` → `0.7.0`) | `feat(orders): add date range filter in global view` |
| `feat!:` / `BREAKING CHANGE:` | **Major** (or minor pre-1.0) | `feat!: refactor staff assignment workflow UI` |
| `docs:` / `test:` / `chore:` | *No release* | `chore: fix lint rule violations in orders page` |

See root [docs/VERSIONING_AND_RELEASES.md](../docs/VERSIONING_AND_RELEASES.md) for full monorepo guidelines.
