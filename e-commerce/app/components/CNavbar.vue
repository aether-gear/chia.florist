<script setup lang="ts">
// Navbar component
</script>

<template>
  <header class="absolute top-0 left-0 w-full z-50 px-8 py-6 flex items-center justify-between">
    <!-- Logo -->
    <div class="flex items-center gap-2 text-white">
      <img 
    src="/images/logo.png" 
    alt="Chia Florist Logo" 
    class="h-20 w-auto" 
  />
    </div>

    <!-- Navigation Links -->
    <nav class="hidden md:flex items-center gap-8 text-white font-medium drop-shadow-sm">
      <NuxtLink to="/" class="hover:text-yellow-300 transition-colors">Home</NuxtLink>
      <NuxtLink to="/about" class="hover:text-yellow-300 transition-colors">About Us</NuxtLink>
      <NuxtLink to="/products" class="hover:text-yellow-300 transition-colors">Products</NuxtLink>
      <NuxtLink to="/contact" class="hover:text-yellow-300 transition-colors">Contact</NuxtLink>
    </nav>

    <!-- Action -->
    <div class="flex items-center gap-4">
      <CButton variant="outline" size="sm" to="/register">Sign In</CButton>
    </div>
  </header>
</template>
