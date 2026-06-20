import { ref } from 'vue'

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

export const triggerAuthAlert = (type: 'success' | 'info' | 'warning' | 'error', message: string) => {
  authAlert.value = { type, message }

  if (alertTimeout) {
    clearTimeout(alertTimeout)
  }

  alertTimeout = setTimeout(() => {
    if (authAlert.value && authAlert.value.message === message) {
      authAlert.value = null
    }
  }, 4000)
}

export const clearAuthAlert = () => {
  authAlert.value = null
  if (alertTimeout) {
    clearTimeout(alertTimeout)
    alertTimeout = null
  }
}
