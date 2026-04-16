<script setup lang="ts">
// Navbar component
</script>

<template>
  <header class="absolute top-0 left-0 w-full z-50 px-8 py-6 flex items-center justify-between">
    <!-- Logo -->
    <div class="flex items-center gap-2 text-white">
      <!-- Simple SVG as placeholder for the flower logo -->
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-yellow-400"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
      <span class="text-2xl font-bold font-serif tracking-wide text-black drop-shadow-md">Chia Florist</span>
    </div>

    <!-- Navigation Links -->
    <nav class="hidden md:flex items-center gap-8 text-white font-medium drop-shadow-sm">
      <NuxtLink to="/" class="hover:text-yellow-300 transition-colors">Home</NuxtLink>
      <NuxtLink to="/about" class="hover:text-yellow-300 transition-colors">About Us</NuxtLink>
      <NuxtLink to="/planters" class="hover:text-yellow-300 transition-colors">Planters</NuxtLink>
      <NuxtLink to="/contact" class="hover:text-yellow-300 transition-colors">Contact</NuxtLink>
    </nav>

    <!-- Action -->
    <div class="flex items-center gap-4">
      <CButton variant="outline" size="sm" to="/signup">Sign In</CButton>
    </div>
  </header>
</template>
