import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useCart } from '~/composables/useCart'
import { triggerSessionExpired, triggerAuthAlert } from '~/composables/useSessionState'

export default defineNuxtPlugin(async (nuxtApp) => {
  const authVm = useAuthViewModel()
  const cartVm = useCart()

  if (import.meta.client) {
    const isLoggedIn = useCookie('is_logged_in')
    const rememberMe = useCookie('remember_me')
    const isGoogleCallback = sessionStorage.getItem('google_auth_pending') === '1'
    const wasLoggedIn = isLoggedIn.value === 'true'

    if (isGoogleCallback) {
      // Clear the flag immediately so it doesn't persist across future loads
      sessionStorage.removeItem('google_auth_pending')

      try {
        await authVm.fetchCurrentUser()

        if (authVm.isAuthenticated.value) {
          triggerAuthAlert('success', `Signed in successfully. Welcome, ${authVm.currentUser.value?.name || 'Customer'}!`)
          await cartVm.loadCart(true)
        } else {
          authVm.clearLocalSession()
          triggerAuthAlert('error', 'Google sign-in failed. Please try again.')
        }
      } catch (err) {
        console.error('Google OAuth session hydration failed:', err)
        authVm.clearLocalSession()
      }

    } else if (isLoggedIn.value === 'true' || rememberMe.value === 'true' || !isLoggedIn.value) {
      // Attempt session restore for email/password or persistent sessions
      try {
        await authVm.fetchCurrentUser()

        if (authVm.isAuthenticated.value) {
          await cartVm.loadCart(true)
        } else {
          if (wasLoggedIn) {
            triggerSessionExpired()
          }
          authVm.clearLocalSession()
        }
      } catch (err) {
        console.error('Session initialization failed:', err)
        if (wasLoggedIn) {
          triggerSessionExpired()
        }
        authVm.clearLocalSession()
      }
    }
  }
})
