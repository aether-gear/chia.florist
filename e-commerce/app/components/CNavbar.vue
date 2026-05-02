<script setup lang="ts">
// Navbar component matching Figma design
</script>

<template>
  <header class="w-full bg-white border-b border-gray-100 px-8 py-3 flex items-center justify-between z-50 sticky top-0">
    
    <div class="flex-1 flex justify-start items-center">
      <NuxtLink to="/" class="flex items-center">
        <img 
          src="/images/logo.png" 
          alt="Chia Florist Logo" 
          class="h-14 w-auto object-contain" 
        />
      </NuxtLink>
    </div>

    <nav class="hidden md:flex items-center gap-10 font-sans text-sm font-medium">
      <NuxtLink to="/" class="text-gray-800 hover:text-black transition-colors">
        Home
      </NuxtLink>
      <NuxtLink to="/contact" class="text-gray-800 hover:text-black transition-colors">
        Contact
      </NuxtLink>
      <NuxtLink to="/about" class="text-gray-800 hover:text-black transition-colors">
        About
      </NuxtLink>
      <NuxtLink to="/register" class="text-gray-900 hover:text-black underline underline-offset-4 font-semibold transition-colors">
        Sign Up
      </NuxtLink>
    </nav>

    <div class="flex-1 flex justify-end">
      <div class="relative max-w-xs w-full">
        <input 
          type="text" 
          placeholder="What are you looking for?" 
          class="w-full bg-[#f5f5f5] text-sm text-gray-700 px-4 py-2 pr-10 rounded-md outline-none border border-transparent focus:border-gray-200 focus:bg-white transition-all placeholder:text-gray-400"
        />
        <div class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>
    </div>

  </header>
</template>