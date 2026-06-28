/**
 * Chia Florist — Cloudflare Edge Worker
 * 
 * Acts as an edge proxy/cache layer in front of service-core.
 * Deploy to Cloudflare Workers and bind to api.chia.florist.
 *
 * Features:
 * - Routes requests to the Railway/Fly.io service-core origin
 * - Caches GET responses for product/public endpoints
 * - Forwards auth headers transparently for protected endpoints
 * - Returns CORS headers for cross-origin SPA requests
 *
 * Setup:
 *   wrangler deploy infra/cloudflare/worker.js --name chia-florist-api-proxy
 */

// --- Configuration --- #
const ORIGIN = 'https://api.chia.florist';          // service-core origin
const CACHE_TTL_SECONDS = 60;                         // public endpoint cache

// Endpoints that are safe to cache at the edge (GET only, no auth required)
const CACHEABLE_PATH_PREFIXES = [
  '/products',
  '/categories',
  '/shops',
  '/locations',
  '/couriers',
];

// CORS allowed origins
const ALLOWED_ORIGINS = [
  'https://chia.florist',
  'https://app.chia.florist',
  'https://panel.chia.florist',
  'http://localhost:4000',
  'http://localhost:5173',
];

// --- Helper: CORS Headers --- #
function getCorsHeaders(requestOrigin) {
  const origin = ALLOWED_ORIGINS.includes(requestOrigin)
    ? requestOrigin
    : ALLOWED_ORIGINS[0];

  return {
    'Access-Control-Allow-Origin': origin,
    'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With',
    'Access-Control-Max-Age': '86400',
    'Vary': 'Origin',
  };
}

// --- Helper: Is request cacheable? --- #
function isCacheable(request) {
  if (request.method !== 'GET') return false;
  const url = new URL(request.url);
  return CACHEABLE_PATH_PREFIXES.some(prefix => url.pathname.startsWith(prefix));
}

// --- Main Handler --- #
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const requestOrigin = request.headers.get('Origin') || '';
    const corsHeaders = getCorsHeaders(requestOrigin);

    // --- Handle CORS preflight --- #
    if (request.method === 'OPTIONS') {
      return new Response(null, { status: 204, headers: corsHeaders });
    }

    // --- Build origin request --- #
    const originUrl = `${ORIGIN}${url.pathname}${url.search}`;
    const originRequest = new Request(originUrl, {
      method:  request.method,
      headers: request.headers,
      body:    request.method !== 'GET' && request.method !== 'HEAD'
                 ? request.body
                 : undefined,
      redirect: 'follow',
    });

    // --- Cache layer for public GET endpoints --- #
    if (isCacheable(request)) {
      const cache = caches.default;
      const cached = await cache.match(originRequest);
      if (cached) {
        const response = new Response(cached.body, cached);
        response.headers.set('CF-Cache-Status', 'HIT');
        Object.entries(corsHeaders).forEach(([k, v]) => response.headers.set(k, v));
        return response;
      }

      const originResponse = await fetch(originRequest);
      const response = new Response(originResponse.body, originResponse);
      Object.entries(corsHeaders).forEach(([k, v]) => response.headers.set(k, v));

      if (originResponse.ok) {
        response.headers.set('Cache-Control', `public, max-age=${CACHE_TTL_SECONDS}`);
        ctx.waitUntil(cache.put(originRequest, response.clone()));
      }

      return response;
    }

    // --- Pass-through for non-cacheable requests --- #
    const originResponse = await fetch(originRequest);
    const response = new Response(originResponse.body, originResponse);
    Object.entries(corsHeaders).forEach(([k, v]) => response.headers.set(k, v));
    return response;
  },
};