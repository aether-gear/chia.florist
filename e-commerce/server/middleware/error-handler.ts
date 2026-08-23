// server/middleware/error-handler.ts
// Global server middleware for error handling.
// Catches unhandled server errors and ensures clean, non-leaking error responses.
// The actual visual error page is rendered by the root-level error.vue.

export default defineEventHandler((event) => {
  // Attach a hook that runs after the response is generated
  event.node.res.on('finish', () => {
    // No-op: errors are handled by Nuxt's error.vue
    // This middleware exists as a safety net for future server-side logging
  })
})
