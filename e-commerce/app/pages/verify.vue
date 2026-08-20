<script setup lang="ts">
import { onMounted } from 'vue'
definePageMeta({
  middleware: 'guest'
})

useHead({ title: 'Redirecting to Verification - Chia Florist' })
onMounted(() => {
  navigateTo('/register?verify=true')
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center font-sans antialiased">
    <div class="text-center space-y-4">
      <div class="w-12 h-12 border-4 border-[#1b4332] border-t-transparent rounded-full animate-spin mx-auto"></div>
      <p class="text-gray-500 text-sm font-semibold">Redirecting you to secure verification page...</p>
    </div>
  </div>
</template>

<style scoped>
/* Animasi micro-shake jika user memasukkan kode OTP yang salah */
.animate-shake { animation: shake 0.3s ease-in-out; }
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}
</style>