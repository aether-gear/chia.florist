<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'

// --- STATE PENCARIAN ---
const isSearchOpen = ref(false)
const searchQuery = ref('')
const searchInput = ref<HTMLInputElement | null>(null)

// --- DATA PRODUK (Source of Truth) ---
const productOfferings = ref([
  { id: 'wedding', name: 'Wedding', image: '/images/wedding.jpeg', price: 8 },
  { id: 'congratulations', name: 'Congratulation', image: '/images/congratulations.jpeg', price: 8 },
  { id: 'condolences', name: 'Condolences', image: '/images/condolences.jpeg', price: 8 },
  { id: 'grand-opening', name: 'Grand Opening', image: '/images/grandop.jpeg', price: 8 },
  { id: 'birthday', name: 'Birthday', image: '/images/birthday.jpeg', price: 8 },
  { id: 'graduate', name: 'Graduate', image: '/images/graduate.jpeg', price: 8 },
  { id: 'anniversary', name: 'Anniversary', image: '/images/anniversary.jpeg', price: 8 },
  { id: 'custom', name: 'Custom Board Simulator', image: '/images/custom-preview.png', price: 10 }
])

// --- LOGIKA PENCARIAN REAL-TIME ---
const filteredResults = computed(() => {
  if (!searchQuery.value) return []
  return productOfferings.value.filter(product => 
    product.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// --- FUNGSI KONTROL ---
const openSearch = async () => {
  isSearchOpen.value = true
  await nextTick()
  searchInput.value?.focus()
}

const closeSearch = () => {
  isSearchOpen.value = false
  searchQuery.value = ''
}
</script>

<template>
  <header class="w-full bg-white border-b border-gray-100 px-8 py-3 flex items-center justify-between z-50 sticky top-0 font-sans">
    
    <!-- SISI KIRI: Logo -->
    <div class="flex-1 flex justify-start items-center">
      <NuxtLink to="/" class="flex items-center">
        <img src="/images/logo.png" alt="Chia Florist Logo" class="h-14 w-auto object-contain" />
      </NuxtLink>
    </div>

    <!-- TENGAH: Menu Navigasi -->
    <nav class="hidden md:flex items-center gap-8 text-sm font-medium">
      <NuxtLink to="/" class="text-gray-800 hover:text-black transition-colors">Home</NuxtLink>
      <a href="https://wa.me/628175234999" target="_blank" class="text-gray-800 hover:text-black transition-colors">Contact</a>
      <NuxtLink to="/cart" class="text-gray-800 hover:text-black transition-colors">Cart</NuxtLink>
      <NuxtLink to="/register" class="text-gray-900 hover:text-black underline underline-offset-4 font-semibold transition-colors">Sign Up</NuxtLink>
    </nav>

    <!-- SISI KANAN: Search & Profile Icons -->
    <div class="flex-1 flex justify-end items-center gap-4">
      <!-- Tombol Pencarian -->
      <button 
        @click="openSearch"
        class="p-2 text-gray-500 hover:text-black hover:bg-gray-50 rounded-full transition-all duration-300"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </button>

      <!-- Akun -->
      <NuxtLink to="/login" class="p-2 text-gray-700 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all group relative">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      </NuxtLink>
    </div>

    <!-- ======================================================== -->
    <!-- SEARCH SIDE DRAWER (Geser dari Kanan) -->
    <!-- ======================================================== -->
    <Teleport to="body">
      <div v-if="isSearchOpen" class="fixed inset-0 z-[100] flex justify-end">
        
        <!-- Backdrop Blur Transparan -->
        <Transition name="fade">
          <div v-if="isSearchOpen" @click="closeSearch" class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        </Transition>

        <!-- Drawer Panel -->
        <Transition name="slide">
          <div v-if="isSearchOpen" class="relative w-full max-w-md bg-white h-full shadow-2xl flex flex-col overflow-hidden">
            
            <!-- Header Drawer -->
            <div class="p-6 flex items-center justify-between border-b border-gray-50">
              <span class="text-xs font-black text-gray-400 uppercase tracking-widest">Search Catalog</span>
              <button @click="closeSearch" class="p-2 text-gray-400 hover:text-black transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Area Input Search -->
            <div class="p-6">
              <div class="relative">
                <input 
                  ref="searchInput"
                  v-model="searchQuery"
                  type="text" 
                  placeholder="Type to search flowers..." 
                  class="w-full bg-[#f5f5f5] text-sm text-gray-800 px-6 py-4 pr-12 rounded-full outline-none border border-transparent focus:border-[#1b4332] focus:bg-white transition-all shadow-inner"
                />
              </div>
            </div>

            <!-- List Hasil Pencarian (Real-time) -->
            <div class="flex-1 overflow-y-auto px-6 pb-6 custom-scrollbar">
              <div v-if="searchQuery && filteredResults.length > 0" class="space-y-4">
                <div class="flex justify-between items-center mb-2">
                  <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Found Products</span>
                  <span class="text-[10px] font-mono text-gray-400 uppercase">{{ filteredResults.length }} results</span>
                </div>
                
                <!-- Hasil Produk -->
                <NuxtLink 
                  v-for="product in filteredResults" 
                  :key="product.id"
                  :to="product.id === 'custom' ? '/products/custom' : `/products/${product.id}`"
                  @click="closeSearch"
                  class="flex gap-4 p-3 rounded-2xl hover:bg-gray-50 transition-all group cursor-pointer border border-transparent hover:border-gray-100"
                >
                  <div class="w-16 h-16 rounded-xl overflow-hidden bg-gray-100 flex-shrink-0">
                    <img :src="product.image" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
                  </div>
                  <div class="flex flex-col justify-center">
                    <h4 class="text-xs font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-tight">
                      {{ product.name }}
                    </h4>
                    <p class="text-[11px] font-semibold text-emerald-600 mt-1">From ${{ product.price }}</p>
                  </div>
                </NuxtLink>
              </div>

              <!-- State: Tidak ada hasil -->
              <div v-else-if="searchQuery && filteredResults.length === 0" class="h-full flex flex-col items-center justify-center opacity-40">
                <p class="text-4xl mb-2">🔍</p>
                <p class="text-xs font-bold">No results found for "{{ searchQuery }}"</p>
              </div>

              <!-- State: Awal Buka -->
              <div v-else class="h-full flex flex-col items-center justify-center opacity-20">
                <p class="text-[10px] font-black uppercase tracking-tighter">Enter a keyword to start...</p>
              </div>
            </div>

          </div>
        </Transition>
      </div>
    </Teleport>
  </header>
</template>

<style scoped>
/* Animasi Slide dari Kanan */
.slide-enter-active, .slide-leave-active {
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from, .slide-leave-to {
  transform: translateX(100%);
}

/* Animasi Fade Backdrop */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.4s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* Custom Scrollbar Mini */
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
</style>