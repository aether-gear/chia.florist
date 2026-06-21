import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { triggerSessionExpired } from '~/composables/useSessionState'

export default defineNuxtPlugin(async (nuxtApp) => {
  const authVm = useAuthViewModel()

  if (import.meta.client) {
    const isLoggedIn = useCookie('is_logged_in')

    if (isLoggedIn.value === 'true') {
      try {
        await authVm.fetchCurrentUser()

        if (!authVm.isAuthenticated.value) {
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
