<script setup lang="ts">
// Footer component
</script>

<template>
  <footer class="bg-[#1b4332] text-white py-12 px-8">
    <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-8">
      <div>
        <h4 class="font-bold mb-4">Site</h4>
        <ul class="space-y-2 text-sm text-gray-300">
          <li><NuxtLink to="/" class="hover:text-white transition-colors">Home</NuxtLink></li>
          <li><NuxtLink to="/blog" class="hover:text-white transition-colors">Blog</NuxtLink></li>
          <li><NuxtLink to="/products" class="hover:text-white transition-colors">Products</NuxtLink></li>
          <li><NuxtLink to="/delivery" class="hover:text-white transition-colors">Delivery</NuxtLink></li>
          <li><NuxtLink to="/about" class="hover:text-white transition-colors">About Us</NuxtLink></li>
        </ul>
      </div>
      <div>
        <h4 class="font-bold mb-4">Social</h4>
        <ul class="space-y-2 text-sm text-gray-300">
          <li><NuxtLink to="/instagram" class="hover:text-white transition-colors">Instagram</NuxtLink></li>
          <li><NuxtLink to="/facebook" class="hover:text-white transition-colors">Facebook</NuxtLink></li>
          <li><NuxtLink to="/whatsapp" class="hover:text-white transition-colors">WhatsApp</NuxtLink></li>
        </ul>
      </div>
      <div>
        <form class="flex" @submit.prevent>
          <button class="bg-yellow-500 text-[#1b4332] px-4 py-2 font-bold hover:bg-yellow-400 transition-colors">Contact Us</button>
        </form>
        <ul class="space-y-2 text-sm text-gray-300">
          <li>chiaflorist@gmail.com</li>
          <li>+62 817-5234-999</li>
          <li>Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng, Kota Muntok, Kab. Bangka Barat, Kep. Bangka Belitung, Indonesia 33311</li>
        </ul>
      </div>
    </div>
    <div class="max-w-7xl mx-auto mt-12 pt-8 border-t border-white/20 text-center text-sm text-gray-400">
      &copy; {{ new Date().getFullYear() }} Copyright © 2026 Chia Florist. All rights reserved.
    </div>
  </footer>
</template>
