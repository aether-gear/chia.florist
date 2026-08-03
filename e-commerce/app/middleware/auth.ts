// app/middleware/auth.ts
// Redirects unauthenticated users to the login page.
// Pages that require authentication should declare:
//   definePageMeta({ middleware: ['auth'] })
export default defineNuxtRouteMiddleware((_to, _from) => {
  if (import.meta.server) return   // SSR: skip — auth state is client-only
  const isLoggedIn = useCookie('is_logged_in')
  if (isLoggedIn.value !== 'true') {
    return navigateTo('/login')
  }
})
