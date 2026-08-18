<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useCart } from '~/composables/useCart' // Import useCart untuk mengambil formatRupiah global
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'
import { bootstrapConfig } from '~/utils/bootstrap'

useHead({
  title: 'Our Collection - Chia Florist',
  meta: [
    { name: 'description', content: 'Explore our premium selection of pre-designed flower boards or launch our custom game simulator.' }
  ]
})

// Ambil helper formatRupiah dari global useCart composable
const { formatRupiah } = useCart()

// Initialize product ViewModel (MVVM Architecture)
const { catalogProducts, isLoading, error, fetchCatalogProducts, page, limit, total, totalPages } = useProductViewModel()

const searchQuery = ref('')
const selectedSort = ref('date:desc')
const selectedShop = ref('')
const shops = ref<{ id: string; name: string; slug: string }[]>([])

const sortOptions = [
  { value: 'date:desc', label: 'Newest First' },
  { value: 'date:asc', label: 'Oldest First' },
  { value: 'name:asc', label: 'Name (A-Z)' },
  { value: 'name:desc', label: 'Name (Z-A)' },
  { value: 'price:asc', label: 'Price (Low to High)' },
  { value: 'price:desc', label: 'Price (High to Low)' },
  { value: 'stock:desc', label: 'Stock (High to Low)' },
  { value: 'stock:asc', label: 'Stock (Low to High)' },
  { value: 'weight:desc', label: 'Weight (High to Low)' },
  { value: 'weight:asc', label: 'Weight (Low to High)' }
]

import { useStoreSelection } from '~/composables/useStoreSelection'

const storeSelection = useStoreSelection()

const fetchShopsList = async () => {
  try {
    const res = await bootstrapConfig.fetchApi<{ shops: { id: string; name: string; slug: string }[] }>('/shops?active=true')
    if (res && Array.isArray(res.shops)) {
      shops.value = res.shops
    }
  } catch (e) {
    console.error('Failed to fetch shops list:', e)
  }
}

const loadProducts = () => {
  const activeShop = storeSelection.selectedShop.value
  fetchCatalogProducts({
    name: searchQuery.value || undefined,
    sort: selectedSort.value,
    shop_id: activeShop?.id || undefined,
    shop_slug: activeShop?.slug || undefined
  })
}

const handleShopDropdownChange = (slug: string) => {
  if (!slug) {
    storeSelection.selectShop(null)
  } else {
    const matched = shops.value.find(s => s.slug === slug)
    if (matched) {
      storeSelection.selectShop(matched)
    }
  }
}

onMounted(() => {
  fetchShopsList()
  loadProducts()
})

// Sync local dropdown with global store selection
watch(storeSelection.selectedShop, (newShop) => {
  selectedShop.value = newShop?.slug || ''
  loadProducts()
}, { immediate: true })

// Reload products when sort selection changes
watch(selectedSort, () => {
  loadProducts()
})

// Debounced search watcher
let searchTimeout: any = null
watch(searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadProducts()
  }, 400)
})

// Interactive simulator card runs client-side and should always be present
const customSimulatorCard = {
  id: 'custom',
  name: 'Custom Board Simulator',
  price: 150000, // Nominal Rupiah murni agar tidak ter-render Rp150
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
  const products = catalogProducts.value || []
  return [...products, customSimulatorCard]
})

// Navigation logic to product details or simulator
const navigateToProduct = (item: any) => {
  if (item.isCustomRoute || item.id === 'custom') {
    navigateTo('/products/custom')
  } else {
    // Navigate using slug if available, fallback to id
    navigateTo(`/products/${item.slug || item.id}`)
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-8 sm:px-10">
      
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

      <!-- Search and Sort Filters Section -->
      <div class="mb-10 bg-white border border-gray-100 rounded-3xl p-6 shadow-sm flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="relative w-full md:w-96">
          <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Search flower boards..." 
            class="w-full bg-gray-50/50 border border-gray-200 rounded-2xl pl-11 pr-4 py-3 text-sm outline-none focus:bg-white focus:border-emerald-700 transition-all font-medium text-gray-800"
          />
        </div>

        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto justify-end">
          <div v-if="shops.length > 0" class="flex items-center gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Store</label>
            <select 
              v-model="selectedShop" 
              @change="handleShopDropdownChange(($event.target as HTMLSelectElement).value)"
              class="bg-gray-50/50 border border-gray-200 rounded-2xl px-4 py-3 text-sm outline-none focus:bg-white focus:border-emerald-700 transition-all font-semibold text-gray-700 cursor-pointer"
            >
              <option value="">All Stores</option>
              <option v-for="shop in shops" :key="shop.id" :value="shop.slug">
                {{ shop.name }}
              </option>
            </select>
          </div>

          <div class="flex items-center gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Sort By</label>
            <select 
              v-model="selectedSort" 
              class="bg-gray-50/50 border border-gray-200 rounded-2xl px-4 py-3 text-sm outline-none focus:bg-white focus:border-emerald-700 transition-all font-semibold text-gray-700 cursor-pointer"
            >
              <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Active Store Pill Indicator Banner -->
      <div v-if="storeSelection.selectedShop.value" class="mb-8 p-4 rounded-2xl bg-emerald-50/80 border border-emerald-200/80 flex items-center justify-between shadow-2xs">
        <div class="flex items-center gap-2.5 text-emerald-900 text-xs md:text-sm font-bold">
          <span class="text-base">📍</span>
          <span>Showing collection available at <span class="underline decoration-emerald-400 font-extrabold">{{ storeSelection.selectedShop.value.name }}</span></span>
        </div>
        <button 
          @click="storeSelection.selectShop(null)" 
          class="text-xs font-extrabold text-emerald-700 hover:text-emerald-900 bg-white px-3 py-1.5 rounded-xl border border-emerald-200 shadow-2xs transition cursor-pointer"
        >
          View All Stores
        </button>
      </div>

      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[300px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-emerald-700"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Loading our collection...</p>
      </div>

      <div v-else-if="error && (!catalogProducts || catalogProducts.length === 0)" class="flex flex-col items-center justify-center min-h-[350px] space-y-6 text-center">
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
              <p class="text-xl font-extrabold text-gray-900">{{ formatRupiah(customSimulatorCard.price) }}</p>
            </div>
            
            <button 
              class="bg-gray-50 group-hover:bg-[#1b4332] text-gray-700 group-hover:text-white border border-gray-200 group-hover:border-[#1b4332] text-xs font-bold px-4 py-2.5 rounded-xl transition-all flex items-center gap-1.5 cursor-pointer"
            >
              <span>Launch Game</span>
              <svg xmlns="http://www.w3.org/2000/xl" class="h-3.5 w-3.5 transform transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div 
          v-for="item in (displayProducts as any[])" 
          :key="item.id"
          @click="navigateToProduct(item)"
          class="group bg-white border border-gray-100 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1"
        >
          <div>
            <div class="aspect-[4/3] w-full bg-gray-200 relative overflow-hidden border-b border-gray-100 flex items-center justify-center">
              <img 
                v-if="item.image"
                :src="item.image" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-sm">
                {{ item.tag || (item.name ? item.name.split(' ')[0] : 'Florist') }}
              </span>
              <span v-if="item.status === 'inactive'" class="absolute top-4 right-4 bg-amber-100 text-amber-900 text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-amber-200 shadow-sm z-20">
                Not Available for Sale
              </span>
              <span v-else-if="!item.isAvailable || item.stock === 0" class="absolute top-4 right-4 bg-red-100 text-red-800 text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-red-200 shadow-sm z-20">
                Sold Out
              </span>
              <span v-else-if="item.stock !== undefined" class="absolute top-4 right-4 bg-emerald-100/90 text-emerald-900 text-[10px] font-extrabold tracking-wider uppercase px-2.5 py-1 rounded-xl border border-emerald-200 shadow-sm z-20 flex items-center gap-1">
                <span>📦</span> {{ item.stock }} in stock
              </span>
            </div>

            <div class="p-6 space-y-3">
              <div class="flex items-center justify-between text-xs font-bold">
                <div class="flex items-center gap-1.5 text-yellow-500">
                  <span>⭐ {{ item.rating ? item.rating.toFixed(1) : '4.8' }}</span>
                  <span class="text-gray-300">|</span>
                  <span class="text-gray-400 font-medium">({{ item.reviews || 150 }} reviews)</span>
                </div>
                <span v-if="item.stock !== undefined" class="text-[11px] font-semibold text-emerald-700">
                  {{ item.stock > 0 ? `${item.stock} left` : 'Out of stock' }}
                </span>
              </div>
              
              <h3 class="text-lg font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>
              
              <p class="text-xs text-gray-400 leading-relaxed line-clamp-2">
                {{ item.desc || item.description || 'Premium quality flower board crafted elegantly for your special moments.' }}
              </p>
            </div>
          </div>

          <div class="p-6 pt-0 border-t border-gray-50/50 mt-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Starting From</p>
              <p class="text-xl font-extrabold text-gray-900">{{ formatRupiah(item.price) }}</p>
            </div>
            
            <button 
              class="bg-gray-50 group-hover:bg-[#1b4332] text-gray-700 group-hover:text-white border border-gray-200 group-hover:border-[#1b4332] text-xs font-bold px-4 py-2.5 rounded-xl transition-all flex items-center gap-1.5 cursor-pointer"
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

<style scoped></style>