<script setup lang="ts">
// Footer component
</script>

<template>
  <footer class="bg-[#1b4332] text-white py-12 px-8">
    <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-8">
      <div>
        <h3 class="text-xl font-serif font-bold text-yellow-500 mb-4">Chia Florist</h3>
        <p class="text-sm text-gray-300">Beautiful flower boards for every special moment. Designed with care, delivered on time.</p>
      </div>
      <div>
        <h4 class="font-bold mb-4">Quick Links</h4>
        <ul class="space-y-2 text-sm text-gray-300">
          <li><NuxtLink to="/" class="hover:text-white transition-colors">Home</NuxtLink></li>
          <li><NuxtLink to="/about" class="hover:text-white transition-colors">About Us</NuxtLink></li>
          <li><NuxtLink to="/planters" class="hover:text-white transition-colors">Planters</NuxtLink></li>
        </ul>
      </div>
      <div>
        <h4 class="font-bold mb-4">Contact Info</h4>
        <ul class="space-y-2 text-sm text-gray-300">
          <li>hello@chia.florist</li>
          <li>+1 234 567 890</li>
          <li>123 Floral Street, NY</li>
        </ul>
      </div>
      <div>
        <h4 class="font-bold mb-4">Newsletter</h4>
        <form class="flex" @submit.prevent>
          <input type="email" placeholder="Your email" class="px-4 py-2 w-full text-black outline-none" />
          <button class="bg-yellow-500 text-[#1b4332] px-4 py-2 font-bold hover:bg-yellow-400 transition-colors">Subscribe</button>
        </form>
      </div>
    </div>
    <div class="max-w-7xl mx-auto mt-12 pt-8 border-t border-white/20 text-center text-sm text-gray-400">
      &copy; {{ new Date().getFullYear() }} Chia Florist. All rights reserved.
    </div>
  </footer>
</template>
