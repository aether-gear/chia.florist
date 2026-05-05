<script setup lang="ts">

// Navbar component with Search & Profile Icon

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



    <nav class="hidden md:flex items-center gap-8 font-sans text-sm font-medium">

      <NuxtLink to="/" class="text-gray-800 hover:text-black transition-colors">

        Home

      </NuxtLink>

     

      <a

        href="https://wa.me/628175234999"

        target="_blank"

        rel="noopener noreferrer"

        class="text-gray-800 hover:text-black transition-colors"

      >

        Contact

      </a>



      <NuxtLink to="/cart" class="text-gray-800 hover:text-black transition-colors">

        Cart

      </NuxtLink>

     

      <NuxtLink to="/register" class="text-gray-900 hover:text-black underline underline-offset-4 font-semibold transition-colors">

        Sign Up

      </NuxtLink>

    </nav>



    <div class="flex-1 flex justify-end items-center gap-4">

     

      <div class="relative max-w-[200px] lg:max-w-xs w-full">

        <input

          type="text"

          placeholder="Search flowers..."

          class="w-full bg-[#f5f5f5] text-xs text-gray-700 px-3.5 py-2 pr-8 rounded-full outline-none border border-transparent focus:border-gray-200 focus:bg-white transition-all placeholder:text-gray-400"

        />

        <div class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none">

          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">

            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />

          </svg>

        </div>

      </div>



      <NuxtLink

        to="/login"

        class="p-2 text-gray-700 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all duration-300 flex items-center justify-center group relative"

        title="Account"

      >

        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">

          <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />

        </svg>



        <span class="absolute top-12 right-0 bg-[#1b4332] text-white text-xs px-3 py-1.5 rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-50 pointer-events-none">

          Sign In / Register

        </span>

      </NuxtLink>



    </div>



  </header>

</template>