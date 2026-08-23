import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useCart } from '~/composables/useCart'
import { triggerSessionExpired } from '~/composables/useSessionState'
import { useGlobalAlert } from '~/composables/useGlobalAlert'

export default defineNuxtPlugin(async (nuxtApp) => {
  const authVm = useAuthViewModel()
  const cartVm = useCart()
  const { showError, showSuccess } = useGlobalAlert()

  if (import.meta.client) {
    const isLoggedIn = useCookie('is_logged_in')
    const rememberMe = useCookie('remember_me')
    const isGoogleCallback = sessionStorage.getItem('google_auth_pending') === '1'
    const wasLoggedIn = isLoggedIn.value === 'true'

    if (isGoogleCallback) {
      // Clear the flag immediately so it doesn't persist across future loads
      sessionStorage.removeItem('google_auth_pending')

      try {
        await authVm.fetchCurrentUser(undefined, true)

        if (authVm.isAuthenticated.value) {
          showSuccess(
            'Signed In Successfully',
            `Welcome, ${authVm.currentUser.value?.name || 'Customer'}!`,
            [
              { label: 'My Profile', onClick: () => navigateTo('/profile/personal') },
              { label: 'Dismiss' }
            ]
          )
          await cartVm.loadCart(true)
        } else {
          authVm.clearLocalSession()
          showError(
            'Google sign-in failed',
            "We couldn't sign you in with Google. Please try again.",
            [
              { label: 'Try Again', onClick: () => navigateTo('/login') },
              { label: 'Dismiss' }
            ]
          )
        }
      } catch (err) {
        console.error('Google OAuth session hydration failed:', err)
        authVm.clearLocalSession()
        showError(
          'Google sign-in failed',
          "We couldn't sign you in with Google. Please try again.",
          [
            { label: 'Try Again', onClick: () => navigateTo('/login') },
            { label: 'Dismiss' }
          ]
        )
      }

    } else if (isLoggedIn.value === 'true' || rememberMe.value === 'true') {
      // Attempt session restore for authenticated or remembered sessions
      try {
        await authVm.fetchCurrentUser(undefined, true)

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
    } else {
      // Unauthenticated visitor: initialize state immediately without calling /auth/me
      await authVm.fetchCurrentUser()
    }
  }
})

