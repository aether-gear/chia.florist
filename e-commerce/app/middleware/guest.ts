import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  const authVm = useAuthViewModel()
  const isLoggedIn = useCookie('is_logged_in')
  const rememberMe = useCookie('remember_me')
  const userProfile = useCookie('user_profile')

  // If already authenticated via reactive state or cookies, redirect immediately to home
  if (authVm.isAuthenticated.value || isLoggedIn.value === 'true' || userProfile.value) {
    return navigateTo('/')
  }

  // If remember_me cookie is present and state is not yet initialized, check session
  if (!authVm.isInitialized.value && rememberMe.value === 'true') {
    try {
      await authVm.fetchCurrentUser()
      if (authVm.isAuthenticated.value) {
        return navigateTo('/')
      }
    } catch (err) {
      console.warn('Guest middleware session check error:', err)
    }
  }
})

