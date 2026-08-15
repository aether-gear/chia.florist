import { ref } from 'vue'
import { useGlobalAlert } from '~/composables/useGlobalAlert'

export const isSessionExpired = ref(false)

export const triggerSessionExpired = () => {
  isSessionExpired.value = true
}

export const clearSessionExpired = () => {
  isSessionExpired.value = false
}

export interface AuthAlert {
  type: 'success' | 'info' | 'warning' | 'error'
  message: string
}

export const authAlert = ref<AuthAlert | null>(null)

let alertTimeout: ReturnType<typeof setTimeout> | null = null

export const triggerAuthAlert = (type: 'success' | 'info' | 'warning' | 'error', message: string, header?: string) => {
  authAlert.value = { type, message }

  const { showAlert } = useGlobalAlert()
  const defaultHeader = type === 'error' ? 'Authentication Error' : type === 'success' ? 'Authentication Success' : type === 'warning' ? 'Authentication Warning' : 'Notice'
  
  showAlert(
    type,
    header || defaultHeader,
    message,
    [
      { label: 'Sign In', onClick: () => navigateTo('/login') },
      { label: 'Dismiss' }
    ]
  )

  if (alertTimeout) {
    clearTimeout(alertTimeout)
  }

  alertTimeout = setTimeout(() => {
    if (authAlert.value && authAlert.value.message === message) {
      authAlert.value = null
    }
  }, 6000)
}

export const clearAuthAlert = () => {
  authAlert.value = null
  const { dismissAlert } = useGlobalAlert()
  dismissAlert()
  if (alertTimeout) {
    clearTimeout(alertTimeout)
    alertTimeout = null
  }
}
