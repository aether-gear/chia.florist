import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  if (import.meta.server) return   // SSR: skip — auth state is client-only

  const authVm = useAuthViewModel()
  const isLoggedIn = useCookie('is_logged_in')
  const rememberMe = useCookie('remember_me')

  // Check if session might exist and needs initialization
  if (!authVm.isInitialized.value && (isLoggedIn.value === 'true' || rememberMe.value === 'true')) {
    try {
      await authVm.fetchCurrentUser()
    } catch (err) {
      console.warn('Guest middleware session check error:', err)
    }
  }

  // If already authenticated, redirect to home
  if (authVm.isAuthenticated.value) {
    return navigateTo('/')
  }
})
