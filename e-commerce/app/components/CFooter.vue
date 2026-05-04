<script setup lang="ts">
// Footer component
</script>

<template>
  <footer class="bg-[#1b4332] text-white py-16 px-8">
    <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 items-start gap-12">
      
      <div class="flex gap-16">
        <div>
          <h4 class="text-lg font-medium mb-6 opacity-80">Site</h4>
          <ul class="space-y-4 font-bold">
            <li><NuxtLink to="/" class="hover:opacity-70 transition">Home</NuxtLink></li>
            <li><NuxtLink to="/blog" class="hover:opacity-70 transition">Blog</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:opacity-70 transition">Products</NuxtLink></li>
            <li><NuxtLink to="/delivery" class="hover:text-white transition">Delivery</NuxtLink></li>
            <li><NuxtLink to="/about" class="hover:text-white transition">About Us</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-lg font-medium mb-6 opacity-80">Social</h4>
          <ul class="space-y-4 font-bold">
            <li><NuxtLink to="/instagram" class="hover:opacity-70 transition">Instagram</NuxtLink></li>
            <li><NuxtLink to="/facebook" class="hover:opacity-70 transition">Facebook</NuxtLink></li>
          </ul>
        </div>
      </div>

      <div class="flex justify-center items-center">
        <img src="/images/footer.png" class="w-90 h-auto" alt="Logo" />
      </div>

      <div class="flex flex-col items-end text-right">
        <NuxtLink 
          to="/contact" 
          class="bg-white text-[#1b4332] px-10 py-3 rounded-lg font-bold text-xl mb-8 hover:bg-gray-100 transition-all shadow-lg"
        >
          Contact us
        </NuxtLink>
        
        <p class="max-w-[300px] text-sm font-medium leading-relaxed opacity-90">
          Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng, Kota Muntok, Kab. Bangka Barat, Kep. Bangka Belitung, Indonesia 33311
        </p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto mt-20 text-center font-medium opacity-80">
      Copyright © {{ new Date().getFullYear() }} Chia Florist. All rights reserved.
    </div>
  </footer>
</template>