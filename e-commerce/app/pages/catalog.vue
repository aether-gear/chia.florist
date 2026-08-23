<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useCart } from '~/composables/useCart'
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'
import { useStoreSelection } from '~/composables/useStoreSelection'
import { productService } from '~/services/productService'

const { formatRupiah } = useCart()
const storeSelection = useStoreSelection()

// Initialize product ViewModel (MVVM Architecture)
const { catalogProducts, isLoading: isVmLoading, error, fetchCatalogProducts, page, limit, total, totalPages } = useProductViewModel()

const searchQuery = ref('')
const selectedSort = ref('date:desc')
const selectedShop = ref('')

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

// Server-side pre-fetch for initial catalog and shops
const { data: initialData, status: initialStatus } = await useAsyncData('catalog-page-init', async () => {
  try {
    const [activeShops, paginatedRes] = await Promise.all([
      storeSelection.fetchActiveShops(),
      productService.getPaginatedCatalogProducts({
        sort: 'date:desc',
        limit: 20
      })
    ])
    return {
      shops: activeShops || [],
      products: paginatedRes?.products || []
    }
  } catch (err) {
    console.error('Failed to load initial catalog on SSR:', err)
    return {
      shops: [],
      products: []
    }
  }
})

const shops = ref<{ id: string; name: string; slug: string }[]>(initialData.value?.shops || [])

const displayProducts = computed(() => {
  return catalogProducts.value?.length > 0
    ? catalogProducts.value
    : (initialData.value?.products || [])
})

const isLoading = computed(() => isVmLoading.value)

const fetchShopsList = async () => {
  try {
    const activeList = await storeSelection.fetchActiveShops()
    if (activeList) {
      shops.value = activeList
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
  if (shops.value.length === 0) {
    fetchShopsList()
  }
})

// Sync local dropdown with global store selection
watch(storeSelection.selectedShop, (newShop) => {
  selectedShop.value = newShop?.slug || ''
  loadProducts()
})

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

useHead({
  title: 'Our Collection — Chia Florist',
  meta: [
    { name: 'description', content: 'Explore our handcrafted selection of flower greeting boards, congratulations bouquets, and custom real-time design simulator.' },
    { property: 'og:title', content: 'Our Collection — Chia Florist' },
    { property: 'og:description', content: 'Explore our handcrafted selection of flower greeting boards, congratulations bouquets, and custom real-time design simulator.' },
    { property: 'og:type', content: 'website' },
    { property: 'og:url', content: 'https://chiaflorist.com/catalog' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: 'Our Collection — Chia Florist' },
    { name: 'twitter:description', content: 'Explore our premium selection of flower greeting boards at Chia Florist.' }
  ],
  link: [
    { rel: 'canonical', href: 'https://chiaflorist.com/catalog' }
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: computed(() => JSON.stringify({
        '@context': 'https://schema.org',
        '@graph': [
          {
            '@type': 'CollectionPage',
            '@id': 'https://chiaflorist.com/catalog',
            'url': 'https://chiaflorist.com/catalog',
            'name': 'Chia Florist Floral Catalog',
            'description': 'Katalog lengkap papan bunga ucapan pernikahan, duka cita, peresmian, dan wisuda.'
          },
          {
            '@type': 'ItemList',
            'numberOfItems': displayProducts.value.length,
            'itemListElement': displayProducts.value.map((item, idx) => ({
              '@type': 'ListItem',
              'position': idx + 1,
              'name': item.name,
              'url': item.isCustomRoute || item.id === 'custom'
                ? 'https://chiaflorist.com/products/custom'
                : `https://chiaflorist.com/products/${(item as any).slug || item.id}`
            }))
          }
        ]
      }))
    }
  ]
})

// Navigation logic to product details
const navigateToProduct = (item: any) => {
  navigateTo(`/products/${item.slug || item.id}`)
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 sm:px-10">
      
      <div class="text-center max-w-2xl mx-auto mb-10 sm:mb-12 space-y-3">
        <span class="text-xs font-black text-emerald-700 uppercase tracking-widest bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
          Chia Florist Collection
        </span>
        <h1 class="text-3xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
          Our Flower Boards
        </h1>
        <p class="text-xs sm:text-sm text-gray-500 leading-relaxed">
          Pilih koleksi papan bunga berkualitas kami atau buat rancangan kustom Anda sendiri dengan simulator 2D interaktif.
        </p>
      </div>

      <!-- FEATURE SECTION: Custom Board Simulator Showcase Banner -->
      <div class="mb-12 relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#1b4332] via-[#245842] to-[#122e22] text-white p-6 sm:p-8 lg:p-10 shadow-xl border border-emerald-800/40">
        <!-- Decorative Glow Accents -->
        <div class="absolute -right-16 -top-16 w-64 h-64 bg-emerald-400/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -left-16 -bottom-16 w-64 h-64 bg-[#4ade80]/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="relative z-10 grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-10 items-center">
          
          <div class="lg:col-span-8 space-y-4">
            <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/10 backdrop-blur-md border border-white/20 text-emerald-300 text-[10px] sm:text-xs font-bold uppercase tracking-widest">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-300 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42" />
              </svg>
              <span>Custom Board Simulator</span>
            </div>

            <h2 class="text-xl sm:text-3xl font-extrabold tracking-tight text-white leading-snug">
              Buat Desain Papan Bunga Sendiri Sesuai Keinginan
            </h2>

            <p class="text-xs sm:text-sm text-gray-200 leading-relaxed max-w-2xl">
              Gunakan simulator 2D interaktif kami untuk menentukan warna busa, susunan teks ucapan, nama pengirim, serta dekorasi bunga sudut secara fleksibel dengan estimasi harga transparan.
            </p>

            <div class="pt-2 flex flex-wrap items-center gap-3">
              <NuxtLink
                to="/products/custom"
                class="bg-[#4ade80] hover:bg-[#3ec470] text-[#1b4332] font-black text-xs sm:text-sm px-6 py-3 rounded-xl shadow-md hover:shadow-lg transition-all inline-flex items-center gap-2"
              >
                <span>Mulai Rancang Sekarang</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                </svg>
              </NuxtLink>
              <span class="text-xs font-semibold text-emerald-200">Mulai dari Rp 150.000</span>
            </div>
          </div>

          <div class="lg:col-span-4 hidden lg:block">
            <div class="relative group cursor-pointer" @click="navigateTo('/products/custom')">
              <div class="aspect-[4/3] rounded-2xl overflow-hidden border-2 border-white/20 shadow-2xl bg-black/40 relative">
                <img 
                  src="/images/custom-preview.png" 
                  alt="Custom Board Simulator Preview"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent flex flex-col justify-end p-4">
                  <span class="text-[9px] font-black uppercase tracking-widest text-emerald-300">Live 2D Canvas</span>
                  <p class="text-[11px] font-bold text-white mt-0.5">Real-time Flower Board Designer</p>
                </div>
              </div>
            </div>
          </div>

        </div>
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
            @input="handleSearchInput"
            type="text" 
            placeholder="Search flower boards, categories..." 
            class="w-full pl-10 pr-4 py-2.5 bg-gray-50 border border-gray-200 rounded-2xl text-xs font-semibold focus:bg-white focus:border-[#1b4332] outline-none transition"
          />
        </div>

        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
          <!-- Multi-Store Dropdown Filter -->
          <div class="relative flex-1 md:flex-none">
            <select 
              :value="storeSelection.selectedShop.value?.slug || ''"
              @change="handleShopDropdownChange(($event.target as HTMLSelectElement).value)"
              class="w-full md:w-auto pl-4 pr-10 py-2.5 bg-gray-50 border border-gray-200 rounded-2xl text-xs font-bold text-gray-700 outline-none focus:border-[#1b4332] transition appearance-none cursor-pointer"
            >
              <option value="">All Locations / Stores</option>
              <option v-for="shop in shops" :key="shop.id" :value="shop.slug">
                {{ shop.name }}
              </option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </div>
          </div>

          <!-- Sort Filter -->
          <div class="relative flex-1 md:flex-none">
            <select 
              v-model="selectedSort" 
              class="w-full md:w-auto pl-4 pr-10 py-2.5 bg-gray-50 border border-gray-200 rounded-2xl text-xs font-bold text-gray-700 outline-none focus:border-[#1b4332] transition appearance-none cursor-pointer"
            >
              <option value="newest">Sort by: Newest Arrival</option>
              <option value="price_low">Sort by: Price (Low to High)</option>
              <option value="price_high">Sort by: Price (High to Low)</option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-gray-500">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Active Location Notification Banner if store selected -->
      <div v-if="storeSelection.selectedShop.value" class="mb-8 flex items-center justify-between bg-emerald-50 border border-emerald-100 px-5 py-3 rounded-2xl text-xs">
        <div class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-[#1b4332] shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
          </svg>
          <span class="text-gray-600 font-medium">Showing exclusive collection for:</span>
          <span class="font-bold text-[#1b4332]">{{ storeSelection.selectedShop.value.name }}</span>
        </div>
        <button 
          @click="clearSelectedShop"
          class="text-xs font-bold text-emerald-800 hover:text-emerald-950 underline cursor-pointer"
        >
          View All Stores
        </button>
      </div>

      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[300px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-emerald-700"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Loading our collection...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="error && (!catalogProducts || catalogProducts.length === 0)" class="flex flex-col items-center justify-center min-h-[350px] space-y-6 text-center">
        <div class="w-16 h-16 rounded-full bg-emerald-50 text-[#1b4332] flex items-center justify-center mx-auto border border-emerald-100">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L9.568 3z" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6z" />
          </svg>
        </div>
        <div>
          <h3 class="text-2xl font-bold text-gray-900">Produk sedang tidak tersedia</h3>
          <p class="text-gray-500 text-sm mt-2 max-w-md mx-auto">
            Maaf, koleksi produk bunga kami saat ini sedang tidak tersedia. Namun, Anda tetap dapat mendesain papan bunga kustom Anda sendiri menggunakan simulator kami.
          </p>
        </div>
        <NuxtLink 
          to="/products/custom"
          class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold text-xs px-6 py-3 rounded-xl transition shadow-sm"
        >
          Buka Custom Board Simulator
        </NuxtLink>
      </div>

      <!-- Products Grid (Clean Cards without internal buttons) -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div 
          v-for="item in (displayProducts as any[])" 
          :key="item.id"
          @click="navigateToProduct(item)"
          class="group bg-white border border-gray-100 rounded-3xl overflow-hidden shadow-2xs hover:shadow-xl hover:border-emerald-200 transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1"
        >
          <div>
            <!-- Card Image -->
            <div class="aspect-[4/3] w-full bg-gray-100 relative overflow-hidden border-b border-gray-100 flex items-center justify-center">
              <img 
                v-if="item.image"
                :src="item.image" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-xs">
                {{ item.tag || (item.name ? item.name.split(' ')[0] : 'Florist') }}
              </span>
              <span v-if="item.status === 'inactive'" class="absolute top-4 right-4 bg-amber-100 text-amber-900 text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-amber-200 shadow-xs z-20">
                Not Available for Sale
              </span>
              <span v-else-if="!item.isAvailable || item.stock === 0" class="absolute top-4 right-4 bg-red-100 text-red-800 text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-red-200 shadow-xs z-20">
                Sold Out
              </span>
              <span v-else-if="item.stock !== undefined" class="absolute top-4 right-4 bg-emerald-100/90 text-emerald-900 text-[10px] font-extrabold tracking-wider uppercase px-2.5 py-1 rounded-xl border border-emerald-200 shadow-xs z-20 flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-800 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
                </svg>
                <span>{{ item.stock }} in stock</span>
              </span>
            </div>

            <!-- Card Info -->
            <div class="p-6 space-y-3">
              <h3 class="text-lg font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>
              
              <p class="text-xs text-gray-400 leading-relaxed line-clamp-2">
                {{ item.desc || item.description || 'Premium quality flower board crafted elegantly for your special moments.' }}
              </p>
            </div>
          </div>

          <!-- Card Footer (Clean without buttons) -->
          <div class="p-6 pt-0 border-t border-gray-50/50 mt-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Starting From</p>
              <p class="text-xl font-extrabold text-[#1b4332]">{{ formatRupiah(item.price) }}</p>
            </div>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>

<style scoped></style>