# --- Stage 1: Build --- #
FROM oven/bun:1.2-alpine AS builder

WORKDIR /app

# Install dependencies (frozen lock file for reproducibility)
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source and build Nuxt SSR output
COPY . .
RUN bun run build

# --- Stage 2: Runtime --- #
FROM node:22-alpine AS runtime

RUN apk add --no-cache wget

WORKDIR /app

# Only copy the compiled Nuxt output
COPY --from=builder /app/.output ./.output

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=4000

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 \
  CMD wget -qO- http://localhost:4000/ || exit 1

CMD ["node", ".output/server/index.mjs"]