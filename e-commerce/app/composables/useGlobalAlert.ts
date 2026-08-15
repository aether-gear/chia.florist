// app/composables/useGlobalAlert.ts
import { ref } from 'vue'

export type GlobalAlertStyle = 'error' | 'success' | 'warning' | 'info'

export interface AlertAction {
  label: string
  onClick?: () => void
  variant?: 'primary' | 'secondary'
}

export interface GlobalAlertOptions {
  id: number
  type: GlobalAlertStyle
  header: string
  body: string
  actions?: AlertAction[]
  actionLabel?: string
  onAction?: () => void
  duration?: number
}

export type ActionInput = 
  | string 
  | AlertAction 
  | AlertAction[]

const activeAlert = ref<GlobalAlertOptions | null>(null)
let dismissTimer: ReturnType<typeof setTimeout> | null = null

export const useGlobalAlert = () => {
  const dismissAlert = () => {
    activeAlert.value = null
    if (dismissTimer) {
      clearTimeout(dismissTimer)
      dismissTimer = null
    }
  }

  const showAlert = (
    type: GlobalAlertStyle,
    header: string,
    body: string,
    actionsOrLabel?: ActionInput,
    onActionOrDuration?: (() => void) | number,
    durationMs = 6000
  ) => {
    if (!body) return

    if (dismissTimer) {
      clearTimeout(dismissTimer)
      dismissTimer = null
    }

    let actions: AlertAction[] = []
    let finalDuration = durationMs

    if (Array.isArray(actionsOrLabel)) {
      actions = actionsOrLabel.slice(0, 2)
      if (typeof onActionOrDuration === 'number') {
        finalDuration = onActionOrDuration
      }
    } else if (typeof actionsOrLabel === 'object' && actionsOrLabel !== null) {
      actions = [actionsOrLabel]
      if (typeof onActionOrDuration === 'number') {
        finalDuration = onActionOrDuration
      }
    } else if (typeof actionsOrLabel === 'string') {
      if (typeof onActionOrDuration === 'function') {
        actions = [{ label: actionsOrLabel, onClick: onActionOrDuration }]
      } else {
        actions = [{ label: actionsOrLabel }]
        if (typeof onActionOrDuration === 'number') {
          finalDuration = onActionOrDuration
        }
      }
    } else {
      // Default action based on type
      const defaultLabel = type === 'success' ? 'Confirm' : 'Got it'
      actions = [{ label: defaultLabel }]
      if (typeof actionsOrLabel === 'number') {
        finalDuration = actionsOrLabel
      }
    }

    activeAlert.value = {
      id: Date.now(),
      type,
      header,
      body,
      actions,
      actionLabel: actions[0]?.label,
      onAction: actions[0]?.onClick,
      duration: finalDuration
    }

    if (finalDuration > 0) {
      dismissTimer = setTimeout(() => {
        dismissAlert()
      }, finalDuration)
    }
  }

  const showError = (
    header: string,
    body: string,
    actionsOrLabel?: ActionInput,
    onActionOrDuration?: (() => void) | number,
    duration = 6000
  ) => showAlert('error', header, body, actionsOrLabel ?? 'Got it', onActionOrDuration, duration)

  const showSuccess = (
    header: string,
    body: string,
    actionsOrLabel?: ActionInput,
    onActionOrDuration?: (() => void) | number,
    duration = 6000
  ) => showAlert('success', header, body, actionsOrLabel ?? 'Confirm', onActionOrDuration, duration)

  const showWarning = (
    header: string,
    body: string,
    actionsOrLabel?: ActionInput,
    onActionOrDuration?: (() => void) | number,
    duration = 6000
  ) => showAlert('warning', header, body, actionsOrLabel ?? 'Got it', onActionOrDuration, duration)

  const showInfo = (
    header: string,
    body: string,
    actionsOrLabel?: ActionInput,
    onActionOrDuration?: (() => void) | number,
    duration = 6000
  ) => showAlert('info', header, body, actionsOrLabel ?? 'Got it', onActionOrDuration, duration)

  return {
    activeAlert,
    showAlert,
    showError,
    showSuccess,
    showWarning,
    showInfo,
    dismissAlert
  }
}
