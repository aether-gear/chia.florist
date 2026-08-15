<script setup lang="ts">
import { computed } from 'vue'
import { useGlobalAlert, type GlobalAlertStyle, type AlertAction } from '~/composables/useGlobalAlert'

export interface AlertProps {
  /** Alert header text (mandatory when used as component prop) */
  header?: string
  /** Alert body text (mandatory when used as component prop) */
  body?: string
  /** Alert status variant type: 'error' | 'success' | 'warning' | 'info' */
  type?: GlobalAlertStyle
  /** Action button text (e.g. "Got it", "Confirm") */
  actionLabel?: string
  /** Up to 2 action buttons */
  actions?: AlertAction[]
  /** Auto-dismiss duration in milliseconds */
  duration?: number
  /** Visibility override for direct component usage */
  show?: boolean
}

const props = withDefaults(defineProps<AlertProps>(), {
  header: '',
  body: '',
  type: undefined,
  actionLabel: undefined,
  actions: undefined,
  duration: undefined,
  show: undefined
})

const emit = defineEmits<{
  (e: 'action'): void
  (e: 'close'): void
  (e: 'update:show', value: boolean): void
}>()

const { activeAlert, dismissAlert } = useGlobalAlert()

// Visibility: explicit prop override or fallback to active global alert state
const isVisible = computed(() => {
  if (props.show !== undefined) return props.show
  return activeAlert.value !== null
})

// Variant style type
const alertType = computed<GlobalAlertStyle>(() => {
  return props.type || activeAlert.value?.type || 'error'
})

// Header text: explicitly passed via prop or active global alert state
const alertHeader = computed(() => {
  return props.header || activeAlert.value?.header || ''
})

// Body message text: explicitly passed via prop or active global alert state
const alertBody = computed(() => {
  return props.body || activeAlert.value?.body || ''
})

// Resolved actions (maximum 2 buttons)
const resolvedActions = computed<AlertAction[]>(() => {
  if (props.actions && props.actions.length > 0) {
    return props.actions.slice(0, 2)
  }
  if (activeAlert.value?.actions && activeAlert.value.actions.length > 0) {
    return activeAlert.value.actions.slice(0, 2)
  }
  const singleLabel = props.actionLabel || activeAlert.value?.actionLabel || (alertType.value === 'success' ? 'Confirm' : 'Got it')
  const singleOnClick = activeAlert.value?.onAction
  return [{ label: singleLabel, onClick: singleOnClick }]
})

// Countdown timer duration (ms)
const alertDuration = computed(() => {
  if (props.duration !== undefined) return props.duration
  return activeAlert.value?.duration ?? 6000
})

// Animation key to force animation re-trigger on new alert
const alertKey = computed(() => {
  if (props.body || props.header) return `${alertHeader.value}-${alertBody.value}`
  return activeAlert.value?.id || Date.now()
})

const handleActionClick = (action: AlertAction) => {
  emit('action')
  if (action.onClick) {
    action.onClick()
  }
  emit('close')
  emit('update:show', false)
  dismissAlert()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="global-alert">
      <div
        v-if="isVisible && (alertHeader || alertBody)"
        :key="alertKey"
        class="fixed top-5 right-5 sm:top-6 sm:right-6 z-[9999] max-w-sm sm:max-w-md w-[calc(100vw-2.5rem)] font-sans pointer-events-auto overflow-hidden rounded-2xl shadow-xl shadow-black/15"
      >
        <div
          class="relative border border-gray-200/90 rounded-2xl px-4 pt-4 sm:px-4.5 sm:pt-4.5 pb-6 flex items-center justify-between gap-3.5 transition-all duration-300 overflow-hidden"
          style="background: linear-gradient(to left, #ffffff 0%, #ffffff 35%, #e5e7eb 85%, #e5e7eb 100%);"
        >
          <!-- Left Icon (Solid version) -->
          <div class="flex-shrink-0 self-center">
            <!-- Checklist (Success) - Green -->
            <svg
              v-if="alertType === 'success'"
              class="w-8 h-8 text-emerald-500"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                clip-rule="evenodd"
              />
            </svg>

            <!-- Warning sign (Warning / Info) - Yellow -->
            <svg
              v-else-if="alertType === 'warning' || alertType === 'info'"
              class="w-8 h-8 text-amber-500"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z"
                clip-rule="evenodd"
              />
            </svg>

            <!-- X sign (Error) - Red -->
            <svg
              v-else
              class="w-8 h-8 text-red-500"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                clip-rule="evenodd"
              />
            </svg>
          </div>

          <!-- Content Body -->
          <div class="flex-1 min-w-0 space-y-1">
            <!-- Header (Mandatory) -->
            <h4
              v-if="alertHeader"
              class="text-xs sm:text-sm font-bold font-sans tracking-tight leading-snug text-black"
            >
              {{ alertHeader }}
            </h4>
            <!-- Body text (Mandatory) -->
            <p v-if="alertBody" class="text-xs font-semibold text-black/75 leading-snug break-words">
              {{ alertBody }}
            </p>
          </div>

          <!-- Action Buttons (max 2 actions, column flex layout with gap-1) -->
          <div v-if="resolvedActions.length > 0" class="flex-shrink-0 self-center flex flex-col gap-1 min-w-[76px]">
            <button
              v-for="(action, index) in resolvedActions"
              :key="index"
              @click="handleActionClick(action)"
              :class="[
                'px-3 py-1.5 text-xs font-bold rounded-xl transition-all cursor-pointer flex items-center justify-center text-center select-none whitespace-nowrap',
                index === 0
                  ? 'bg-gray-950 hover:bg-gray-900 text-white shadow-xs active:scale-95'
                  : 'bg-white/90 hover:bg-white text-gray-800 border border-gray-200/90 shadow-2xs active:scale-95 text-[11px]'
              ]"
            >
              {{ action.label }}
            </button>
          </div>

          <!-- Timer Progress Line Indicator at the Bottom -->
          <div
            v-if="alertDuration && alertDuration > 0"
            class="absolute bottom-0 left-0 right-0 h-1 overflow-hidden bg-[#4ade80]/25"
          >
            <div
              class="h-full progress-bar-line bg-[#4ade80]"
              :style="{ animationDuration: `${alertDuration}ms` }"
            ></div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* Slide in from right, leave to right */
.global-alert-enter-active {
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.3s ease;
}
.global-alert-leave-active {
  transition: transform 0.35s cubic-bezier(0.7, 0, 0.84, 0), opacity 0.25s ease;
}
.global-alert-enter-from {
  opacity: 0;
  transform: translate3d(calc(100% + 2rem), 0, 0);
}
.global-alert-leave-to {
  opacity: 0;
  transform: translate3d(calc(100% + 2rem), 0, 0);
}

/* Countdown progress bar animation */
.progress-bar-line {
  width: 100%;
  animation-name: timer-countdown;
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@keyframes timer-countdown {
  from {
    width: 100%;
  }
  to {
    width: 0%;
  }
}
</style>
