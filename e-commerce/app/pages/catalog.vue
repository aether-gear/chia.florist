<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'

useHead({
  title: 'Our Collection - Chia Florist',
  meta: [
    { name: 'description', content: 'Explore our premium selection of pre-designed flower boards or launch our custom game simulator.' }
  ]
})

// Initialize product ViewModel (MVVM Architecture)
const { catalogProducts, isLoading, error, fetchCatalogProducts } = useProductViewModel()

onMounted(() => {
  fetchCatalogProducts()
})

// The interactive simulator card runs client-side and should always be present
const customSimulatorCard = {
  id: 'custom',
  name: 'Interactive Custom Board Simulator',
  price: 150.00,
  rating: 5.0,
  reviews: 89,
  image: '/images/custom-preview.png',
  tag: 'Interactive Game',
  desc: 'Design your own professional flower board in real-time! Choose your custom layout, foam colors, and fonts.',
  isCustomRoute: true,
  isAvailable: true
}

// Combine dynamic products with the simulator game card
const displayProducts = computed(() => {
  return [...catalogProducts.value, customSimulatorCard]
})

// Navigation logic to product details or simulator
const navigateToProduct = (item: typeof displayProducts.value[0]) => {
  if (item.isCustomRoute || item.id === 'custom') {
    navigateTo('/products/custom')
  } else {
    navigateTo(`/products/${item.id}`)
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <div class="text-center max-w-2xl mx-auto mb-16 space-y-3">
        <span class="text-xs font-black text-emerald-700 uppercase tracking-widest bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
          Chia Florist Collection
        </span>
        <h1 class="text-4xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
          Our Flower Boards
        </h1>
        <p class="text-sm md:text-base text-gray-500 leading-relaxed">
          Select our pre-designed premium flower boards or jump directly into our interactive real-time simulator game to design your custom creation.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[300px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-emerald-700"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Loading our collection...</p>
      </div>

      <!-- Error / Empty State (Produk Sedang Tidak Tersedia) -->
      <div v-else-if="error && catalogProducts.length === 0" class="flex flex-col items-center justify-center min-h-[350px] space-y-6 text-center">
        <div class="text-5xl">🌸</div>
        <div>
          <h3 class="text-2xl font-bold text-gray-900">Produk sedang tidak tersedia</h3>
          <p class="text-gray-500 text-sm mt-2 max-w-md mx-auto">
            Maaf, koleksi produk bunga kami saat ini sedang tidak tersedia. Namun, Anda tetap dapat mendesain papan bunga kustom Anda sendiri menggunakan simulator kami di bawah ini.
          </p>
        </div>
        
        <div 
          @click="navigateTo('/products/custom')"
          class="group bg-white border border-gray-200 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1 max-w-sm text-left mt-8"
        >
          <div>
            <div class="aspect-[4/3] w-full bg-gray-50 relative overflow-hidden border-b border-gray-50">
              <img 
                :src="customSimulatorCard.image" 
                :alt="customSimulatorCard.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              <span class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-sm">
                {{ customSimulatorCard.tag }}
              </span>
            </div>

            <div class="p-6 space-y-3">
              <div class="flex items-center gap-2 text-xs text-yellow-500 font-bold">
                <span>⭐ {{ customSimulatorCard.rating.toFixed(1) }}</span>
                <span class="text-gray-300">|</span>
                <span class="text-gray-400 font-medium">({{ customSimulatorCard.reviews }} reviews)</span>
              </div>
              
              <h3 class="text-lg font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ customSimulatorCard.name }}
              </h3>
              
              <p class="text-xs text-gray-400 leading-relaxed line-clamp-2">
                {{ customSimulatorCard.desc }}
              </p>
            </div>
          </div>

          <div class="p-6 pt-0 border-t border-gray-50/50 mt-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Starting From</p>
              <p class="text-xl font-extrabold text-gray-900">${{ customSimulatorCard.price.toFixed(2) }}</p>
            </div>
            
            <button 
              class="bg-gray-50 group-hover:bg-[#1b4332] text-gray-700 group-hover:text-white border border-gray-200 group-hover:border-[#1b4332] text-xs font-bold px-4 py-2.5 rounded-xl transition-all flex items-center gap-1.5"
            >
              <span>Launch Game</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 transform transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Content Grid State -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div 
          v-for="item in displayProducts" 
          :key="item.id"
          @click="navigateToProduct(item)"
          class="group bg-white border border-gray-100 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1"
        >
          <div>
            <div class="aspect-[4/3] w-full bg-gray-50 relative overflow-hidden border-b border-gray-50">
              <img 
                :src="item.image" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              <span class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-sm">
                {{ item.tag }}
              </span>
            </div>

            <div class="p-6 space-y-3">
              <div class="flex items-center gap-2 text-xs text-yellow-500 font-bold">
                <span>⭐ {{ item.rating.toFixed(1) }}</span>
                <span class="text-gray-300">|</span>
                <span class="text-gray-400 font-medium">({{ item.reviews }} reviews)</span>
              </div>
              
              <h3 class="text-lg font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>
              
              <p class="text-xs text-gray-400 leading-relaxed line-clamp-2">
                {{ item.desc }}
              </p>
            </div>
          </div>

          <div class="p-6 pt-0 border-t border-gray-50/50 mt-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Starting From</p>
              <p class="text-xl font-extrabold text-gray-900">${{ item.price.toFixed(2) }}</p>
            </div>
            
            <button 
              class="bg-gray-50 group-hover:bg-[#1b4332] text-gray-700 group-hover:text-white border border-gray-200 group-hover:border-[#1b4332] text-xs font-bold px-4 py-2.5 rounded-xl transition-all flex items-center gap-1.5"
            >
              <span>{{ item.isCustomRoute ? 'Launch Game' : 'View Details' }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 transform transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>