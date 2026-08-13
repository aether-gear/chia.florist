import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  if (import.meta.server) return   // SSR: skip — auth state is client-only

  const authVm = useAuthViewModel()
  const isLoggedIn = useCookie('is_logged_in')

  if (!authVm.isInitialized.value && isLoggedIn.value === 'true') {
    try {
      await authVm.fetchCurrentUser()
    } catch (err) {
      console.warn('Auth middleware session check error:', err)
    }
  }

  if (!authVm.isAuthenticated.value) {
    return navigateTo('/login')
  }
})
