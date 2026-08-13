<script setup lang="ts">
// CFooter Component
</script>

<template>
  <footer class="bg-white text-gray-900 border-t border-gray-100 p-0 w-full font-brand">

    <!-- Batik Top Border Divider (h-12, White Background, Clean Kawung Motif) -->
    <div class="w-full h-12 bg-white overflow-hidden flex items-center justify-center border-y border-y-[#1b4332]/30">
      <svg class="w-full h-12 text-[#1b4332] opacity-75" xmlns="http://www.w3.org/2000/svg" width="100%" height="48">
        <defs>
          <pattern id="batik-kawung" x="0" y="0" width="48" height="48" patternUnits="userSpaceOnUse">
            <rect width="48" height="48" fill="#ffffff" />

            <!-- Central Floral Motif (Center at 24, 24) -->
            <ellipse cx="24" cy="16" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="24" cy="32" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="16" cy="24" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="32" cy="24" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />
            <circle cx="24" cy="24" r="2" fill="currentColor" />

            <!-- Interlocking Corner Motifs -->
            <ellipse cx="0" cy="16" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="48" cy="16" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="0" cy="32" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="48" cy="32" rx="6" ry="7.5" fill="none" stroke="currentColor" stroke-width="1.2" />

            <ellipse cx="16" cy="0" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="16" cy="48" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="32" cy="0" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />
            <ellipse cx="32" cy="48" rx="7.5" ry="6" fill="none" stroke="currentColor" stroke-width="1.2" />

            <!-- Corner Centers -->
            <circle cx="0" cy="0" r="2" fill="currentColor" />
            <circle cx="48" cy="0" r="2" fill="currentColor" />
            <circle cx="0" cy="48" r="2" fill="currentColor" />
            <circle cx="48" cy="48" r="2" fill="currentColor" />
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#batik-kawung)" />
      </svg>
    </div>

    <div class="max-w-7xl mx-auto px-8 sm:px-10 space-y-12 flex flex-col pt-12 pb-8">

      <div class="space-y-3">
        <img src="/images/logo.png" class="h-14 w-auto object-contain" alt="Chia Florist" />
      </div>

      <!-- Main Footer Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">

        <!-- Main Branch Address & WhatsApp Container -->
        <div class="lg:col-span-5 space-y-6">
          <div class="space-y-6">
            <p class="text-xs font-bold uppercase tracking-wider text-emerald-800">Main Branch</p>
            <p class="text-xs text-gray-500 leading-relaxed max-w-sm">
              Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng, Kab. Bangka Barat, Kep. Bangka Belitung, Indonesia 33311
            </p>
          </div>

          <!-- WhatsApp Container with Flat White Background & Button CTA -->
          <div class="p-4 bg-gray-50 border border-gray-100 rounded-2xl space-y-3 max-w-sm">
            <p class="text-xs text-gray-600 leading-relaxed">
              Ada yang mau ditanyakan atau pengen request buket khusus? Chat kami di WhatsApp ya, kami siap bantu dengan senang hati!
            </p>
            <div>
              <a
                href="https://wa.me/628175234999"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 bg-[#1b4332] hover:bg-[#143326] text-white px-4 py-2 rounded-full text-xs font-semibold transition-all duration-200 shadow-sm hover:shadow active:scale-[0.98]"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M12.004 2c-5.518 0-9.996 4.477-9.996 9.996 0 1.764.459 3.42 1.258 4.863L2 22l5.316-1.393c1.393.754 2.972 1.19 4.63 1.19 5.518 0 10.004-4.477 10.004-9.997 0-5.519-4.486-9.996-10.004-9.996zm6.362 14.234c-.267.753-1.332 1.383-1.834 1.433-.467.046-.48.291-2.731-.607-2.616-1.043-4.24-3.723-4.368-3.896-.13-.173-1.045-1.392-1.045-2.656 0-1.264.649-1.884.882-2.133.232-.249.51-.31.68-.31.171 0 .341.002.489.009.155.008.363-.058.568.441.21.511.716 1.743.778 1.874.062.13.104.283.017.457-.086.173-.13.282-.26.432-.13.149-.272.333-.388.446-.13.125-.266.26-.114.521.151.26.669 1.104 1.434 1.785.986.876 1.815 1.147 2.074 1.277.26.13.411.11.564-.065.152-.174.647-.754.822-1.01.175-.259.347-.216.586-.13.238.086 1.516.715 1.776.845.26.13.433.195.498.303.064.108.064.629-.204 1.382z"/>
                </svg>
                <span>Tanya via WhatsApp</span>
              </a>
            </div>
          </div>
        </div>

        <!-- Navigation Columns with Swapped Colors -->
        <div class="lg:col-span-7 grid grid-cols-2 sm:grid-cols-3 gap-8">

          <!-- Navigation Column (Home, Catalog, Search) -->
          <div class="space-y-5.5">
            <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-900">Navigation</h4>
            <ul class="space-y-1">
              <li>
                <NuxtLink to="/" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Home</NuxtLink>
              </li>
              <li>
                <NuxtLink to="/catalog" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Catalog</NuxtLink>
              </li>
              <li>
                <NuxtLink to="/catalog" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Search</NuxtLink>
              </li>
            </ul>
          </div>

          <!-- Legal Column -->
          <div class="space-y-5.5">
            <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-900">Legal</h4>
            <ul class="space-y-1">
              <li>
                <NuxtLink to="/terms" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Terms & Conditions</NuxtLink>
              </li>
              <li>
                <NuxtLink to="/privacy" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Privacy Policy</NuxtLink>
              </li>
            </ul>
          </div>

          <!-- Connect Column -->
          <div class="space-y-5.5">
            <h4 class="text-xs font-semibold uppercase tracking-wider text-gray-900">Connect</h4>
            <ul class="space-y-1">
              <li>
                <a href="https://instagram.com" target="_blank" rel="noopener noreferrer" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">Instagram</a>
              </li>
              <li>
                <a href="https://wa.me/628175234999" target="_blank" rel="noopener noreferrer" class="relative inline-block py-0.5 text-sm md:text-md font-medium text-gray-600 hover:text-[#1b4332] transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100">WhatsApp</a>
              </li>
            </ul>
          </div>

        </div>

    </div>

      <!-- Copyright Bottom Bar -->
      <div class="pt-6 border-t border-gray-100 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs font-medium text-gray-400">
        <p>Copyright © {{ new Date().getFullYear() }} Chia Florist. All rights reserved.</p>
        <p>Handcrafted with care</p>
      </div>

    </div>
  </footer>
</template>
