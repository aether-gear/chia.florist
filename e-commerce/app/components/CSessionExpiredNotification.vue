<script setup lang="ts">
import { isSessionExpired, clearSessionExpired } from '~/composables/useSessionState'

const handleSignIn = () => {
  clearSessionExpired()
  navigateTo('/login')
}

const handleDismiss = () => {
  clearSessionExpired()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="slide-fade">
      <div 
        v-if="isSessionExpired" 
        class="fixed top-6 right-6 z-[9999] max-w-sm w-full bg-white/90 backdrop-blur-md rounded-2xl shadow-2xl border border-red-100 overflow-hidden font-sans"
      >
        <div class="p-6">
          <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-50 flex items-center justify-center border border-red-100 animate-pulse">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>

            <!-- Content -->
            <div class="flex-1 space-y-1">
              <h4 class="text-sm font-bold text-gray-900 leading-tight">Session Expired</h4>
              <p class="text-xs text-gray-500 font-semibold leading-relaxed">
                Your secure session has expired. Please sign in again to access your profile and check out your cart.
              </p>
            </div>

            <!-- Close icon -->
            <button 
              @click="handleDismiss" 
              class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors p-1 rounded-full hover:bg-gray-100/50"
              title="Dismiss"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Actions -->
          <div class="mt-4 flex items-center justify-end gap-3">
            <button 
              @click="handleDismiss" 
              class="px-3.5 py-1.5 text-xs font-semibold text-gray-600 hover:text-gray-800 transition-colors cursor-pointer"
            >
              Dismiss
            </button>
            <button 
              @click="handleSignIn" 
              class="px-4 py-2 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-lg transition-all shadow-md hover:shadow-lg cursor-pointer"
            >
              Sign In
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.slide-fade-enter-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.7, 0, 0.84, 0);
}
.slide-fade-enter-from {
  transform: translateY(-20px) scale(0.95);
  opacity: 0;
}
.slide-fade-leave-to {
  transform: translateY(-20px) scale(0.95);
  opacity: 0;
}
</style>
