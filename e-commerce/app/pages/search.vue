<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productService } from '~/services/productService'
import { useStoreSelection } from '~/composables/useStoreSelection'
import { useCart } from '~/composables/useCart'
import type { CatalogProduct } from '~/types/product'

const route = useRoute()
const router = useRouter()
const { formatRupiah } = useCart()
const storeSelection = useStoreSelection()

// URL query state
const searchQuery = ref((route.query.q as string) || '')
const selectedSort = ref((route.query.sort as string) || 'date:desc')
const selectedShop = ref((route.query.shop as string) || '')

const sortOptions = [
  { value: 'date:desc', label: 'Newest First' },
  { value: 'date:asc', label: 'Oldest First' },
  { value: 'name:asc', label: 'Name (A-Z)' },
  { value: 'name:desc', label: 'Name (Z-A)' },
  { value: 'price:asc', label: 'Price (Low to High)' },
  { value: 'price:desc', label: 'Price (High to Low)' }
]

const popularSearches = [
  'Pernikahan',
  'Duka Cita',
  'Grand Opening',
  'Selamat & Sukses',
  'Wisuda',
  'Custom Board'
]

// Fetch active shops for filter dropdown
const { data: shopsData } = await useAsyncData('search-shops', async () => {
  try {
    return await storeSelection.fetchActiveShops()
  } catch (e) {
    console.error('Failed to load shops for search page:', e)
    return []
  }
})
const shops = computed(() => shopsData.value || [])

// Interactive simulator card representation
const customSimulatorCard: CatalogProduct = {
  id: 'custom',
  name: 'Custom Flower Board Simulator',
  price: 150000,
  rating: 5.0,
  reviews: 89,
  image: '/images/custom-preview.png',
  tag: 'Interactive Game',
  desc: 'Design your own professional flower board in real-time with custom floral arrangements and lettering!',
  isCustomRoute: true,
  isAvailable: true
}

import { filterCatalogProductsByQuery } from '~/utils/searchMatcher'

// Fetch search results on SSR and when query changes
const asyncDataKey = computed(() => `search-${searchQuery.value}-${selectedSort.value}-${selectedShop.value}`)

const { data: searchResponse, status, refresh } = await useAsyncData(
  asyncDataKey.value,
  async () => {
    try {
      const activeShopSlug = selectedShop.value || storeSelection.selectedShop.value?.slug
      const activeShopId = selectedShop.value
        ? shops.value.find(s => s.slug === selectedShop.value)?.id
        : storeSelection.selectedShop.value?.id

      const res = await productService.getPaginatedCatalogProducts({
        sort: selectedSort.value,
        shop_id: activeShopId || undefined,
        shop_slug: activeShopSlug || undefined,
        limit: 50
      })

      const rawProducts = res.products || []
      const filtered = filterCatalogProductsByQuery(
        rawProducts,
        searchQuery.value,
        true,
        customSimulatorCard
      )

      return {
        products: filtered,
        total: filtered.length,
        hasSearched: Boolean(searchQuery.value.trim())
      }
    } catch (err) {
      console.error('Search query failed:', err)
      return {
        products: filterCatalogProductsByQuery([], searchQuery.value, true, customSimulatorCard),
        total: 0,
        hasSearched: Boolean(searchQuery.value.trim())
      }
    }
  },
  {
    watch: [searchQuery, selectedSort, selectedShop]
  }
)

const searchResults = computed(() => searchResponse.value?.products || [])
const totalResults = computed(() => searchResults.value.length)
const isLoading = computed(() => status.value === 'pending')

// Update URL query parameters cleanly
const syncUrlParams = () => {
  const query: Record<string, string> = {}
  if (searchQuery.value.trim()) query.q = searchQuery.value.trim()
  if (selectedSort.value && selectedSort.value !== 'date:desc') query.sort = selectedSort.value
  if (selectedShop.value) query.shop = selectedShop.value

  router.replace({ path: '/search', query })
}

// Form submit or Enter press
const handleSearchSubmit = () => {
  syncUrlParams()
  refresh()
}

const applyQuickSearch = (keyword: string) => {
  searchQuery.value = keyword
  syncUrlParams()
  refresh()
}

const clearSearch = () => {
  searchQuery.value = ''
  syncUrlParams()
  refresh()
}

const handleShopChange = (slug: string) => {
  selectedShop.value = slug
  syncUrlParams()
  refresh()
}

const handleSortChange = () => {
  syncUrlParams()
  refresh()
}

// Sync route changes if navigated with back/forward buttons
watch(() => route.query, (newQuery) => {
  searchQuery.value = (newQuery.q as string) || ''
  selectedSort.value = (newQuery.sort as string) || 'date:desc'
  selectedShop.value = (newQuery.shop as string) || ''
})

// --- SEO & META TAGS ---
const pageTitle = computed(() => {
  if (searchQuery.value.trim()) {
    return `Search: "${searchQuery.value.trim()}" — Chia Florist`
  }
  return 'Search Flower Collections & Custom Boards — Chia Florist'
})

const pageDescription = computed(() => {
  if (searchQuery.value.trim()) {
    return `Discover flower boards, congratulations bouquets, and custom simulator arrangements matching "${searchQuery.value.trim()}" at Chia Florist. Fast delivery.`
  }
  return 'Find premium flower boards for weddings, inaugurations, graduations, and condolences. Order online or design your own in real-time at Chia Florist.'
})

useHead({
  title: pageTitle,
  meta: [
    { name: 'description', content: pageDescription },
    { property: 'og:title', content: pageTitle },
    { property: 'og:description', content: pageDescription },
    { property: 'og:type', content: 'website' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: pageTitle },
    { name: 'twitter:description', content: pageDescription },
    { name: 'robots', content: 'index, follow' }
  ],
  link: [
    { rel: 'canonical', href: `https://chiaflorist.com/search${searchQuery.value ? `?q=${encodeURIComponent(searchQuery.value)}` : ''}` }
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: computed(() => JSON.stringify({
        '@context': 'https://schema.org',
        '@graph': [
          {
            '@type': 'SearchResultsPage',
            '@id': 'https://chiaflorist.com/search',
            'url': 'https://chiaflorist.com/search',
            'name': pageTitle.value,
            'description': pageDescription.value,
            'potentialAction': {
              '@type': 'SearchAction',
              'target': 'https://chiaflorist.com/search?q={search_term_string}',
              'query-input': 'required name=search_term_string'
            }
          },
          {
            '@type': 'BreadcrumbList',
            'itemListElement': [
              {
                '@type': 'ListItem',
                'position': 1,
                'name': 'Home',
                'item': 'https://chiaflorist.com/'
              },
              {
                '@type': 'ListItem',
                'position': 2,
                'name': 'Search',
                'item': 'https://chiaflorist.com/search'
              },
              ...(searchQuery.value.trim() ? [{
                '@type': 'ListItem',
                'position': 3,
                'name': searchQuery.value.trim(),
                'item': `https://chiaflorist.com/search?q=${encodeURIComponent(searchQuery.value.trim())}`
              }] : [])
            ]
          },
          ...(searchResults.value.length > 0 ? [{
            '@type': 'ItemList',
            'numberOfItems': searchResults.value.length,
            'itemListElement': searchResults.value.map((item, idx) => ({
              '@type': 'ListItem',
              'position': idx + 1,
              'name': item.name,
              'url': item.isCustomRoute || item.id === 'custom'
                ? 'https://chiaflorist.com/products/custom'
                : `https://chiaflorist.com/products/${item.slug || item.id}`
            }))
          }] : [])
        ]
      }))
    }
  ]
})
</script>

<template>
  <div class="min-h-screen bg-[#fcfbf9] font-brand pb-24">
    <!-- HERO / SEARCH HEADER -->
    <header class="bg-gradient-to-b from-emerald-950 via-[#1b4332] to-[#143225] text-white pt-12 pb-16 px-6 sm:px-10 relative overflow-hidden">
      <div class="absolute inset-0 opacity-10 bg-[radial-gradient(#a7f3d0_1px,transparent_1px)] [background-size:16px_16px]"></div>
      
      <div class="max-w-4xl mx-auto text-center space-y-6 relative z-10">
        <!-- Breadcrumb Navigation for SEO & UX -->
        <nav aria-label="Breadcrumb" class="flex justify-center items-center gap-2 text-xs text-emerald-200/80 font-medium">
          <NuxtLink to="/" class="hover:text-white transition-colors">Home</NuxtLink>
          <span>/</span>
          <span class="text-white font-semibold">Search</span>
          <template v-if="searchQuery.trim()">
            <span>/</span>
            <span class="text-emerald-300 truncate max-w-[200px]">"{{ searchQuery.trim() }}"</span>
          </template>
        </nav>

        <h1 class="text-3xl sm:text-4xl md:text-5xl font-extrabold tracking-tight">
          Find Your Perfect Flower Board
        </h1>
        <p class="text-emerald-100/90 text-sm sm:text-base max-w-xl mx-auto font-light leading-relaxed">
          Search across our handcrafted flower collections, greeting boards, or launch our real-time custom simulator.
        </p>

        <!-- Search Bar Form -->
        <form @submit.prevent="handleSearchSubmit" class="max-w-2xl mx-auto relative mt-8">
          <div class="flex items-center bg-white rounded-2xl shadow-xl p-1.5 border border-emerald-900/30 transition-all focus-within:ring-4 focus-within:ring-emerald-500/20">
            <span class="pl-4 pr-2 text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </span>

            <input
              v-model="searchQuery"
              type="search"
              placeholder="Search by flower type, occasion (e.g. Wedding, Condolence)..."
              class="w-full bg-transparent text-gray-900 text-sm md:text-base font-medium px-2 py-3 outline-none placeholder:text-gray-400"
            />

            <button
              v-if="searchQuery"
              type="button"
              @click="clearSearch"
              class="p-2 text-gray-400 hover:text-gray-600 cursor-pointer mr-1"
              title="Clear search"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <button
              type="submit"
              class="bg-[#1b4332] hover:bg-[#143225] text-white text-xs sm:text-sm font-bold px-5 sm:px-7 py-3 rounded-xl transition-all shadow-md cursor-pointer flex-shrink-0"
            >
              Search
            </button>
          </div>
        </form>

        <!-- Popular Search Chips -->
        <div class="flex flex-wrap items-center justify-center gap-2 pt-2 text-xs">
          <span class="text-emerald-200/80 font-medium mr-1">Popular:</span>
          <button
            v-for="keyword in popularSearches"
            :key="keyword"
            @click="applyQuickSearch(keyword)"
            type="button"
            class="px-3 py-1 rounded-full bg-white/10 hover:bg-white/20 text-white font-medium border border-white/15 transition-all text-[11px] cursor-pointer"
          >
            {{ keyword }}
          </button>
        </div>
      </div>
    </header>

    <!-- MAIN SEARCH CONTENT AREA -->
    <main class="max-w-7xl mx-auto px-6 sm:px-10 -mt-6 relative z-20">
      
      <!-- CONTROLS & FILTER BAR -->
      <section class="bg-white rounded-2xl p-4 sm:p-5 shadow-sm border border-gray-100 mb-8 flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <h2 class="text-sm font-bold text-gray-800">
            <span v-if="searchQuery.trim()">Results for <span class="text-[#1b4332] font-black">"{{ searchQuery.trim() }}"</span></span>
            <span v-else>All Collections</span>
          </h2>
          <span class="text-xs bg-emerald-50 text-emerald-800 font-bold px-2.5 py-1 rounded-full border border-emerald-100">
            {{ totalResults }} items found
          </span>
        </div>

        <div class="flex flex-wrap items-center gap-3 justify-start md:justify-end">
          <!-- Store branch selector -->
          <div v-if="shops.length > 0" class="flex items-center gap-2">
            <label for="search-shop-filter" class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Store:</label>
            <select
              id="search-shop-filter"
              :value="selectedShop"
              @change="handleShopChange(($event.target as HTMLSelectElement).value)"
              class="bg-gray-50 border border-gray-200 rounded-xl px-3.5 py-2 text-xs font-semibold text-gray-700 outline-none focus:border-emerald-700 focus:bg-white transition cursor-pointer"
            >
              <option value="">All Branches</option>
              <option v-for="shop in shops" :key="shop.id" :value="shop.slug">
                {{ shop.name }}
              </option>
            </select>
          </div>

          <!-- Sort selector -->
          <div class="flex items-center gap-2">
            <label for="search-sort-filter" class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Sort:</label>
            <select
              id="search-sort-filter"
              v-model="selectedSort"
              @change="handleSortChange"
              class="bg-gray-50 border border-gray-200 rounded-xl px-3.5 py-2 text-xs font-semibold text-gray-700 outline-none focus:border-emerald-700 focus:bg-white transition cursor-pointer"
            >
              <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>
      </section>

      <!-- LOADING STATE -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[350px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-emerald-700"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Searching our floral arrangements...</p>
      </div>

      <!-- EMPTY STATE -->
      <div
        v-else-if="searchResults.length === 0"
        class="bg-white rounded-3xl p-10 sm:p-14 text-center border border-gray-100 shadow-sm max-w-2xl mx-auto space-y-6"
      >
        <div class="w-20 h-20 bg-emerald-50 text-emerald-700 rounded-full flex items-center justify-center mx-auto text-3xl shadow-inner">
          🔍
        </div>
        <div class="space-y-2">
          <h3 class="text-xl font-bold text-gray-900">No flower boards match "{{ searchQuery }}"</h3>
          <p class="text-xs sm:text-sm text-gray-500 max-w-md mx-auto leading-relaxed">
            We couldn't find any exact matches for your search. Try checking your spelling, using broader keywords, or create a fully customized design in our simulator.
          </p>
        </div>

        <div class="pt-4 flex flex-wrap items-center justify-center gap-3">
          <NuxtLink
            to="/products/custom"
            class="bg-[#1b4332] hover:bg-[#143225] text-white text-xs font-bold px-6 py-3 rounded-xl transition shadow cursor-pointer"
          >
            Launch Custom Simulator
          </NuxtLink>
          <NuxtLink
            to="/catalog"
            class="bg-gray-100 hover:bg-gray-200 text-gray-800 text-xs font-bold px-6 py-3 rounded-xl transition cursor-pointer"
          >
            Browse Full Catalog
          </NuxtLink>
        </div>
      </div>

      <!-- RESULTS PRODUCT GRID -->
      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
      >
        <article
          v-for="product in searchResults"
          :key="product.id"
          class="group bg-white border border-gray-100/80 rounded-3xl overflow-hidden shadow-xs hover:shadow-xl transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1"
          @click="navigateTo(product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`)"
        >
          <div>
            <!-- Product Thumbnail Image -->
            <div class="aspect-[4/3] w-full bg-gray-50 relative overflow-hidden border-b border-gray-50">
              <img
                v-if="product.image"
                :src="product.image"
                :alt="product.name"
                loading="lazy"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              />
              <div v-else class="w-full h-full flex items-center justify-center bg-gray-100 text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>

              <!-- Tag Badge -->
              <span
                v-if="product.tag"
                class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-xs"
              >
                {{ product.tag }}
              </span>

              <!-- Out of stock badge -->
              <span
                v-if="product.isAvailable === false"
                class="absolute top-4 right-4 bg-red-600/90 text-white text-[10px] font-bold px-2.5 py-1 rounded-xl shadow-xs"
              >
                Out of Stock
              </span>
            </div>

            <!-- Product Body Details -->
            <div class="p-5 space-y-2.5">
              <div class="flex items-center gap-2 text-xs text-amber-500 font-bold">
                <span>⭐ {{ (product.rating || 5.0).toFixed(1) }}</span>
                <span class="text-gray-300">|</span>
                <span class="text-gray-400 font-medium">({{ product.reviews || 0 }} reviews)</span>
              </div>

              <h3 class="text-base font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug line-clamp-2">
                {{ product.name }}
              </h3>

              <p class="text-xs text-gray-500 leading-relaxed line-clamp-2">
                {{ product.desc || 'Premium handcrafted floral arrangement from Chia Florist.' }}
              </p>
            </div>
          </div>

          <!-- Card Footer Price & CTA -->
          <div class="p-5 pt-0">
            <div class="pt-3 border-t border-gray-100 flex items-center justify-between">
              <div>
                <span class="text-[10px] text-gray-400 uppercase font-bold tracking-wider block">Price</span>
                <span class="text-sm font-black text-[#1b4332]">
                  {{ formatRupiah(product.price) }}
                </span>
              </div>

              <NuxtLink
                :to="product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`"
                class="px-3.5 py-1.5 rounded-xl text-xs font-bold bg-emerald-50 text-[#1b4332] group-hover:bg-[#1b4332] group-hover:text-white transition-colors"
              >
                {{ product.isCustomRoute || product.id === 'custom' ? 'Design' : 'View' }}
              </NuxtLink>
            </div>
          </div>
        </article>
      </div>
    </main>
  </div>
</template>
