import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  const authVm = useAuthViewModel()
  const isLoggedIn = useCookie('is_logged_in')
  const rememberMe = useCookie('remember_me')

  if (!authVm.isInitialized.value && (isLoggedIn.value === 'true' || rememberMe.value === 'true')) {
    try {
      await authVm.fetchCurrentUser()
    } catch (err) {
      console.warn('Auth middleware session check error:', err)
    }
  }

  if (!authVm.isAuthenticated.value && isLoggedIn.value !== 'true') {
    return navigateTo('/login')
  }
})


