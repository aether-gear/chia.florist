// app/composables/useToast.ts
// Lightweight global toast queue. Shows one toast at a time.
// Design mirrors the success-toast used in custom.vue (cart added notification).
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  type: ToastType
  title: string
  subtitle?: string
}

const toasts = ref<Toast[]>([])
let _counter = 0

export const useToast = () => {
  const show = (type: ToastType, title: string, subtitle?: string, durationMs = 3200) => {
    const id = ++_counter
    toasts.value.push({ id, type, title, subtitle })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, durationMs)
  }

  const success = (title: string, subtitle?: string) => show('success', title, subtitle)
  const error   = (title: string, subtitle?: string) => show('error',   title, subtitle, 4000)
  const info    = (title: string, subtitle?: string) => show('info',    title, subtitle)

  const dismiss = (id: number) => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  return { toasts, show, success, error, info, dismiss }
}
