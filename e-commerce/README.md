# Chia Florist — E-Commerce (Nuxt 4)

Customer-facing e-commerce storefront for Chia Florist, featuring product browsing, DIY Custom Flower Board Designer, SEO SSR optimization, and checkout flows.

## Setup

Make sure to install dependencies:

```bash
# bun (recommended)
bun install

# or npm
npm install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
bun run dev
```

## Production Build

Build the application for production:

```bash
bun run build
bun run preview
```

## Versioning & Commit Guidelines

`e-commerce` releases are automated via GitHub Actions (Current baseline: `e-commerce-v0.9.0`).
Docker images are automatically tagged and published to GHCR upon release.

| Commit Type | Version Bump | Example |
| :--- | :--- | :--- |
| `fix:` / `perf:` | **Patch** (`0.9.0` → `0.9.1`) | `fix(cart): resolve custom product retention in cart` |
| `feat:` | **Minor** (`0.9.0` → `0.10.0`) | `feat(custom-product): add size selection for regular product` |
| `feat!:` / `BREAKING CHANGE:` | **Major** (or minor pre-1.0) | `feat!: migrate data fetching from SPA to Nuxt SSR` |
| `docs:` / `test:` / `chore:` | *No release* | `docs: add status for cart endpoints` |

### Custom Product Versioning
* **Payload Schema Version**: `3.0.0` ([docs/specs/custom-product-v3.schema.json](../docs/specs/custom-product-v3.schema.json))
* **Editor Engine Version**: `3.1.0` (`app/features/custom-product/constants.ts`)

See root [docs/VERSIONING_AND_RELEASES.md](../docs/VERSIONING_AND_RELEASES.md) for full monorepo guidelines.
