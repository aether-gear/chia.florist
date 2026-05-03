<script setup lang="ts">
// Footer component
</script>

<template>
  <footer class="bg-[#1b4332] text-white py-12 px-6">
    <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 items-center gap-8">
      
      <div class="flex gap-12 justify-start">
        <div>
          <h4 class="text-base font-medium mb-4 opacity-80">Site</h4>
          <ul class="space-y-3 font-semibold text-sm">
            <li><NuxtLink to="/" class="hover:opacity-70 transition">Home</NuxtLink></li>
            <li><NuxtLink to="/blog" class="hover:opacity-70 transition">Blog</NuxtLink></li>
            <li><NuxtLink to="/products" class="hover:opacity-70 transition">Products</NuxtLink></li>
            <li><NuxtLink to="/delivery" class="hover:opacity-70 transition">Delivery</NuxtLink></li>
            <li><NuxtLink to="/about" class="hover:opacity-70 transition">About Us</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-base font-medium mb-4 opacity-80">Social</h4>
          <ul class="space-y-3 font-semibold text-sm">
            <li><NuxtLink to="/instagram" class="hover:opacity-70 transition">Instagram</NuxtLink></li>
            <li><NuxtLink to="/facebook" class="hover:opacity-70 transition">Facebook</NuxtLink></li>
          </ul>
        </div>
      </div>

      <div class="flex justify-center items-center">
        <img src="/images/footer.png" class="w-48 h-auto object-contain" alt="Logo" />
      </div>

      <div class="flex flex-col items-end text-right">
        <a 
          href="https://wa.me/628175234999" 
          target="_blank" 
          rel="noopener noreferrer"
          class="bg-white text-[#1b4332] px-6 py-2.5 rounded-lg font-bold text-base mb-6 hover:bg-gray-100 transition-all shadow-md flex items-center gap-2.5 hover:scale-102"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12.004 2c-5.518 0-9.996 4.477-9.996 9.996 0 1.764.459 3.42 1.258 4.863L2 22l5.316-1.393c1.393.754 2.972 1.19 4.63 1.19 5.518 0 10.004-4.477 10.004-9.997 0-5.519-4.486-9.996-10.004-9.996zm6.362 14.234c-.267.753-1.332 1.383-1.834 1.433-.467.046-.48.291-2.731-.607-2.616-1.043-4.24-3.723-4.368-3.896-.13-.173-1.045-1.392-1.045-2.656 0-1.264.649-1.884.882-2.133.232-.249.51-.31.68-.31.171 0 .341.002.489.009.155.008.363-.058.568.441.21.511.716 1.743.778 1.874.062.13.104.283.017.457-.086.173-.13.282-.26.432-.13.149-.272.333-.388.446-.13.125-.266.26-.114.521.151.26.669 1.104 1.434 1.785.986.876 1.815 1.147 2.074 1.277.26.13.411.11.564-.065.152-.174.647-.754.822-1.01.175-.259.347-.216.586-.13.238.086 1.516.715 1.776.845.26.13.433.195.498.303.064.108.064.629-.204 1.382z"/>
          </svg>
          <span>Contact us</span>
        </a>
        
        <p class="max-w-[300px] text-xs font-medium leading-relaxed opacity-90">
          Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng, Kota Muntok, Kab. Bangka Barat, Kep. Bangka Belitung, Indonesia 33311
        </p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto mt-12 text-center text-xs font-medium opacity-80">
      Copyright © {{ new Date().getFullYear() }} Chia Florist. All rights reserved.
    </div>
  </footer>
</template>