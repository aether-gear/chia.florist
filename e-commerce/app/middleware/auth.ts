import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  if (import.meta.server) return   // SSR: skip — auth state is client-only

  const authVm = useAuthViewModel()

  if (!authVm.isInitialized.value) {
    try {
      await authVm.fetchCurrentUser()
    } catch (err) {
      console.warn('Auth middleware session check error:', err)
    }
  }

  const isLoggedIn = useCookie('is_logged_in')
  if (!authVm.isAuthenticated.value && isLoggedIn.value !== 'true') {
    return navigateTo('/login')
  }
})
