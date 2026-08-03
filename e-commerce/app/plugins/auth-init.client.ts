import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useCart } from '~/composables/useCart'
import { triggerSessionExpired, triggerAuthAlert } from '~/composables/useSessionState'

export default defineNuxtPlugin(async (nuxtApp) => {
  const authVm = useAuthViewModel()
  const cartVm = useCart()

  if (import.meta.client) {
    const isLoggedIn = useCookie('is_logged_in')
    const isGoogleCallback = sessionStorage.getItem('google_auth_pending') === '1'

    if (isGoogleCallback) {
      // Clear the flag immediately so it doesn't persist across future loads
      sessionStorage.removeItem('google_auth_pending')

      // The backend has set the 'chast' session cookie (HttpOnly — not readable
      // by JS) and redirected back here. Call fetchCurrentUser() which goes
      // through the Nuxt server route and forwards all browser cookies to the
      // backend, so 'chast' is sent and the session is validated server-side.
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

    } else if (isLoggedIn.value === 'true') {
      // Standard session restore for email/password login sessions
      try {
        await authVm.fetchCurrentUser()

        if (authVm.isAuthenticated.value) {
          await cartVm.loadCart(true)
        } else {
          triggerSessionExpired()
          authVm.clearLocalSession()
        }
      } catch (err) {
        console.error('Session initialization failed:', err)
        triggerSessionExpired()
        authVm.clearLocalSession()
      }
    }
  }
})
