<script setup lang="ts">
import { ref, watch, onUnmounted, nextTick } from 'vue'
import { productService } from '~/services/productService'
import { formatRupiah } from '~/utils/formatter'
import type { CatalogProduct } from '~/types/product'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useStoreSelection, type Shop } from '~/composables/useStoreSelection'

const authVm = useAuthViewModel()
const { selectedShop, activeShops, isLoadingShops, fetchActiveShops, selectShop } = useStoreSelection()

// --- STATE STORE PICKER ---
const isStoreModalOpen = ref(false)

const openStoreModal = async () => {
  isStoreModalOpen.value = true
  await fetchActiveShops()
}

const closeStoreModal = () => {
  isStoreModalOpen.value = false
}

const handleSelectShop = (shop: Shop | null) => {
  selectShop(shop)
  closeStoreModal()
}

// --- STATE PENCARIAN ---
const isSearchOpen = ref(false)
const searchQuery = ref('')
const searchInput = ref<HTMLInputElement | null>(null)
const searchResults = ref<CatalogProduct[]>([])
const isLoadingSearch = ref(false)

// Simulator card data (Frontend Interactive component)
const customSimulatorCard: CatalogProduct = {
  id: 'custom',
  name: 'Custom Flower Board Simulator',
  price: 150000,
  rating: 5.0,
  reviews: 89,
  image: '/images/custom-preview.png',
  tag: 'Interactive Game',
  desc: 'Design your own professional flower board in real-time!',
  isCustomRoute: true,
  isAvailable: true
}

let debounceTimeout: ReturnType<typeof setTimeout> | null = null

// --- LOGIKA PENCARIAN REAL-TIME DENGAN DEBOUNCE ---
watch(searchQuery, (newQuery) => {
  if (debounceTimeout) {
    clearTimeout(debounceTimeout)
  }

  const query = newQuery.trim()
  if (!query) {
    searchResults.value = []
    return
  }

  isLoadingSearch.value = true

  debounceTimeout = setTimeout(async () => {
    try {
      // Fetch matching products from Go API backend
      const apiResults = await productService.getCatalogProducts({
        name: query,
        shop_id: selectedShop.value?.id
      })

      // Check if user search matches our custom board simulator
      const matchesCustom = customSimulatorCard.name.toLowerCase().includes(query.toLowerCase()) ||
                            'custom'.includes(query.toLowerCase()) ||
                            'simulator'.includes(query.toLowerCase())

      if (matchesCustom) {
        searchResults.value = [customSimulatorCard, ...apiResults]
      } else {
        searchResults.value = apiResults
      }
    } catch (err) {
      console.error('Search failed:', err)
      searchResults.value = []
    } finally {
      isLoadingSearch.value = false
    }
  }, 300) // 300ms debounce
})

onUnmounted(() => {
  if (debounceTimeout) {
    clearTimeout(debounceTimeout)
  }
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
  searchResults.value = []
}

const route = useRoute()
</script>

<template>
  <header class="
    w-full bg-white border-b border-gray-100 px-6
    py-3 flex items-center justify-between z-50 sticky top-0 font-sans
    md:px-8
  ">

    <div class="flex-1 flex justify-start items-center gap-4">
      <NuxtLink to="/" class="flex">
        <img src="/images/logo.png" alt="Chia Florist" class="h-6 md:h-6.75 lg:h-7.25 xl:h-8" />
      </NuxtLink>

      <!-- Store Picker Pill Button (Hidden on /cart page) -->
      <button
        v-if="route.path !== '/cart'"
        @click="openStoreModal"
        class="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-50/80 hover:bg-emerald-100 border border-emerald-200/80 text-emerald-800 text-xs font-bold transition-all duration-300 shadow-2xs group cursor-pointer"
        title="Change Store Location"
      >
        <span class="text-sm">📍</span>
        <span class="max-w-[140px] truncate">
          {{ selectedShop ? selectedShop.name : 'All Stores' }}
        </span>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-emerald-600 group-hover:translate-y-0.5 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
        </svg>
      </button>
    </div>

    <div class="hidden lg:block flex-1"></div>

    <div class="flex-1 flex justify-end items-center gap-2 sm:gap-4">

      <!-- Mobile Store Picker Button (Hidden on /cart page) -->
      <button
        v-if="route.path !== '/cart'"
        @click="openStoreModal"
        class="sm:hidden p-2 text-emerald-700 hover:bg-emerald-50 rounded-full transition-all duration-300"
        title="Change Store"
      >
        <span class="text-base">📍</span>
      </button>

      <button
        @click="openSearch"
        class="p-2 text-gray-500 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all duration-300"
        title="Search"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </button>

      <NuxtLink
        to="/cart"
        class="p-2 text-gray-500 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all duration-300"
        title="Cart"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 0 0-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 0 0-16.536-1.84M7.5 14.25L5.106 5.25M16.5 20.25a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Zm3 0a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z" />
        </svg>
      </NuxtLink>

      <!-- Profile & Auth status indicator -->
      <div class="flex items-center gap-1 sm:gap-2 border-l border-gray-100 pl-2 sm:pl-4">
        <template v-if="authVm.isAuthenticated.value">
          <NuxtLink
            to="/profile"
            class="flex items-center gap-2 py-1 px-2.5 hover:bg-gray-50 rounded-full transition-all duration-300 group"
            title="Profile Settings"
          >
            <div class="relative w-8 h-8 rounded-full overflow-hidden bg-emerald-50 border border-emerald-100 flex items-center justify-center text-[#1b4332] group-hover:bg-[#1b4332] group-hover:text-white transition-all duration-300">
              <img v-if="authVm.currentUser.value?.avatarUrl" :src="authVm.currentUser.value.avatarUrl" class="w-full h-full object-cover" />
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              <!-- Green online dot overlay -->
              <span class="absolute bottom-0 right-0 w-2.5 h-2.5 bg-emerald-500 border-2 border-white rounded-full"></span>
            </div>
            <span class="hidden sm:inline text-xs font-bold text-gray-700 group-hover:text-[#1b4332] transition-colors max-w-[90px] truncate">
              {{ authVm.currentUser.value?.name || 'Customer' }}
            </span>
          </NuxtLink>

          <button
            @click="authVm.logout"
            class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-full transition-all duration-300 cursor-pointer"
            title="Sign Out"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
            </svg>
          </button>
        </template>

        <template v-else>
          <NuxtLink
            to="/login"
            class="flex items-center gap-1.5 py-1.5 px-4 bg-[#1b4332] hover:bg-[#143326] text-white rounded-full text-xs font-bold transition-all duration-300 shadow-sm hover:shadow cursor-pointer"
            title="Sign In"
          >
            <span>Sign In</span>
          </NuxtLink>
        </template>
      </div>
    </div>

    <!-- Store Picker Modal Teleport -->
    <Teleport to="body">
      <div v-if="isStoreModalOpen" class="fixed inset-0 z-[110] flex items-center justify-center p-4">
        <Transition name="fade">
          <div v-if="isStoreModalOpen" @click="closeStoreModal" class="absolute inset-0 bg-black/50 backdrop-blur-xs"></div>
        </Transition>

        <Transition name="slide">
          <div v-if="isStoreModalOpen" class="relative w-full max-w-md bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100 z-10">
            <div class="p-6 border-b border-gray-100 flex justify-between items-center bg-emerald-50/50">
              <div>
                <h3 class="text-lg font-extrabold text-gray-900 flex items-center gap-2">
                  <span>📍</span> Select Store Location
                </h3>
                <p class="text-xs text-gray-500 mt-0.5">Browse flower boards available at your nearest store.</p>
              </div>
              <button @click="closeStoreModal" class="p-2 text-gray-400 hover:text-gray-700 rounded-full transition-colors cursor-pointer">
                ✕
              </button>
            </div>

            <div class="p-6 max-h-[60vh] overflow-y-auto space-y-3 custom-scrollbar">
              <!-- Option 1: All Stores -->
              <div
                @click="handleSelectShop(null)"
                :class="[!selectedShop ? 'border-emerald-600 bg-emerald-50/60 ring-2 ring-emerald-600/20' : 'border-gray-200 hover:border-emerald-300 bg-white']"
                class="p-4 rounded-2xl border transition-all cursor-pointer flex items-center justify-between group"
              >
                <div>
                  <h4 class="font-bold text-sm text-gray-900 group-hover:text-emerald-700 transition-colors">All Stores</h4>
                  <p class="text-xs text-gray-500">Show catalog products across all active stores</p>
                </div>
                <span v-if="!selectedShop" class="text-emerald-600 font-bold text-sm">✓ Selected</span>
              </div>

              <div v-if="isLoadingShops" class="py-8 text-center text-xs text-gray-400 font-medium animate-pulse">
                Loading active stores...
              </div>

              <!-- Option 2+: Active Stores -->
              <div
                v-else
                v-for="shop in activeShops"
                :key="shop.id"
                @click="handleSelectShop(shop)"
                :class="[selectedShop?.id === shop.id ? 'border-emerald-600 bg-emerald-50/60 ring-2 ring-emerald-600/20' : 'border-gray-200 hover:border-emerald-300 bg-white']"
                class="p-4 rounded-2xl border transition-all cursor-pointer flex items-center justify-between group"
              >
                <div>
                  <h4 class="font-bold text-sm text-gray-900 group-hover:text-emerald-700 transition-colors">{{ shop.name }}</h4>
                  <p class="text-xs text-gray-500 mt-0.5">{{ shop.description || `Branch: ${shop.slug}` }}</p>
                </div>
                <span v-if="selectedShop?.id === shop.id" class="text-emerald-600 font-bold text-sm">✓ Selected</span>
              </div>
            </div>

            <div class="p-4 bg-gray-50 border-t border-gray-100 text-center">
              <button
                @click="closeStoreModal"
                class="w-full py-3 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-xl transition cursor-pointer"
              >
                Done
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Teleport>

    <!-- Search Teleport -->
    <Teleport to="body">
      <div v-if="isSearchOpen" class="fixed inset-0 z-[100] flex justify-end">

        <Transition name="fade">
          <div v-if="isSearchOpen" @click="closeSearch" class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        </Transition>

        <Transition name="slide">
          <div v-if="isSearchOpen" class="relative w-full max-w-md bg-white h-full shadow-2xl flex flex-col overflow-hidden">

            <div class="p-6 flex items-center justify-between border-b border-gray-50">
              <span class="text-xs font-black text-gray-400 uppercase tracking-widest">Search Catalog</span>
              <button @click="closeSearch" class="p-2 text-gray-400 hover:text-black transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="p-6">
              <div class="relative">
                <input
                  ref="searchInput"
                  v-model="searchQuery"
                  type="text"
                  placeholder="Type to search flowers..."
                  class="w-full bg-surface text-sm text-gray-800 px-6 py-4 pr-12 rounded-full outline-none border border-transparent focus:border-accent focus:bg-white transition-all shadow-inner"
                />
              </div>
            </div>

            <div class="flex-1 overflow-y-auto px-6 pb-6 custom-scrollbar">
              <div v-if="isLoadingSearch" class="flex flex-col items-center justify-center py-12 space-y-3">
                <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-emerald-700"></div>
                <p class="text-xs text-gray-500 font-medium animate-pulse">Searching...</p>
              </div>

              <div v-else-if="searchQuery && searchResults.length > 0" class="space-y-4">
                <div class="flex justify-between items-center mb-2">
                  <span class="text-[10px] font-black text-gray-400 uppercase tracking-widest">Found Products</span>
                  <span class="text-[10px] font-mono text-gray-400 uppercase">{{ searchResults.length }} results</span>
                </div>

                <NuxtLink
                  v-for="product in searchResults"
                  :key="product.id"
                  :to="product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`"
                  @click="closeSearch"
                  class="flex gap-4 p-3 rounded-2xl hover:bg-gray-50 transition-all group cursor-pointer border border-transparent hover:border-gray-100"
                >
                  <div class="w-16 h-16 rounded-xl overflow-hidden bg-gray-200 flex-shrink-0 flex items-center justify-center">
                    <img v-if="product.image" :src="product.image" :alt="product.name" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" />
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div class="flex flex-col justify-center">
                    <h4 class="text-xs font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-tight">
                      {{ product.name }}
                    </h4>
                    <p class="text-[11px] font-semibold text-emerald-600 mt-1">Starting From {{ formatRupiah(product.price) }}</p>
                  </div>
                </NuxtLink>
              </div>

              <div v-else-if="searchQuery && searchResults.length === 0" class="h-full flex flex-col items-center justify-center opacity-40">
                <p class="text-4xl mb-2">🔍</p>
                <p class="text-xs font-bold">No results found for "{{ searchQuery }}"</p>
              </div>

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
  background: var(--color-border-soft);
  border-radius: 10px;
}
</style>
