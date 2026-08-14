<script setup lang="ts">
// Dedicated layout for Auth pages (Login, Register, Password Reset)
</script>

<template>
  <div class="min-h-screen bg-white flex flex-col justify-between items-center px-4 py-8 md:py-12 font-sans text-gray-900 selection:bg-[#4ade80] selection:text-[#245842] antialiased">
    <!-- Header: Brand Logo -->
    <header class="w-full flex justify-center pt-2 pb-4">
      <NuxtLink to="/" class="inline-block transition-transform hover:scale-105 duration-200">
        <img src="/images/logo.png" class="h-8 md:h-10 w-auto object-contain" alt="Chia Florist Logo" />
      </NuxtLink>
    </header>

    <!-- Main Form Container -->
    <main class="w-full max-w-sm sm:max-w-md mx-auto flex-1 flex flex-col justify-center py-4">
      <slot />
    </main>

    <!-- Bottom Footer Links -->
    <footer class="w-full text-center pt-8 pb-2 text-xs text-gray-500 space-x-3">
      <NuxtLink to="/privacy" class="hover:text-gray-900 hover:underline transition-colors">Privacy policy</NuxtLink>
      <span class="text-gray-300">•</span>
      <NuxtLink to="/terms" class="hover:text-gray-900 hover:underline transition-colors">Terms of service</NuxtLink>
    </footer>
  </div>
</template>

<style scoped>
body {
  margin: 0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}
</style>
