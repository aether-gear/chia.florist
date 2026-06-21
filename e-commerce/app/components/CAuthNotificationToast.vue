<script setup lang="ts">
import { authAlert, clearAuthAlert } from '~/composables/useSessionState'

const handleClose = () => {
  clearAuthAlert()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="toast-fade">
      <div 
        v-if="authAlert" 
        class="fixed bottom-6 right-6 z-[9999] max-w-sm w-full bg-white/90 backdrop-blur-md rounded-2xl shadow-2xl border overflow-hidden font-sans transition-all duration-300"
        :class="{
          'border-emerald-100': authAlert.type === 'success',
          'border-blue-100': authAlert.type === 'info',
          'border-amber-100': authAlert.type === 'warning',
          'border-red-100': authAlert.type === 'error'
        }"
      >
        <div class="p-5 flex items-start gap-4">
          <!-- Success Icon -->
          <div 
            v-if="authAlert.type === 'success'" 
            class="flex-shrink-0 w-8 h-8 rounded-full bg-emerald-50 border border-emerald-100 flex items-center justify-center text-emerald-600"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>

          <!-- Info Icon -->
          <div 
            v-else-if="authAlert.type === 'info'" 
            class="flex-shrink-0 w-8 h-8 rounded-full bg-blue-50 border border-blue-100 flex items-center justify-center text-blue-600"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>

          <!-- Warning Icon -->
          <div 
            v-else-if="authAlert.type === 'warning'" 
            class="flex-shrink-0 w-8 h-8 rounded-full bg-amber-50 border border-amber-100 flex items-center justify-center text-amber-600"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>

          <!-- Error Icon -->
          <div 
            v-else 
            class="flex-shrink-0 w-8 h-8 rounded-full bg-red-50 border border-red-100 flex items-center justify-center text-red-600"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>

          <!-- Content -->
          <div class="flex-1 space-y-0.5">
            <h4 
              class="text-xs font-bold uppercase tracking-wider leading-none"
              :class="{
                'text-emerald-800': authAlert.type === 'success',
                'text-blue-800': authAlert.type === 'info',
                'text-amber-800': authAlert.type === 'warning',
                'text-red-800': authAlert.type === 'error'
              }"
            >
              {{ authAlert.type === 'success' ? 'Success' : authAlert.type === 'info' ? 'Notification' : authAlert.type === 'warning' ? 'Alert' : 'Error' }}
            </h4>
            <p class="text-xs text-gray-700 font-semibold leading-normal">
              {{ authAlert.message }}
            </p>
          </div>

          <!-- Close Button -->
          <button 
            @click="handleClose" 
            class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors p-1 rounded-full hover:bg-gray-100/50 cursor-pointer"
            title="Dismiss"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.toast-fade-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-fade-leave-active {
  transition: all 0.25s cubic-bezier(0.7, 0, 0.84, 0);
}
.toast-fade-enter-from {
  transform: translateY(20px) scale(0.95);
  opacity: 0;
}
.toast-fade-leave-to {
  transform: translateY(20px) scale(0.95);
  opacity: 0;
}
</style>
