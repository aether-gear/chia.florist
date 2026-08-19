<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productService } from '~/services/productService'
import { useStoreSelection } from '~/composables/useStoreSelection'
import { useCart } from '~/composables/useCart'
import type { CatalogProduct } from '~/types/product'
import { filterCatalogProductsByQuery } from '~/utils/searchMatcher'

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
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-8 sm:px-10">

      <!-- Breadcrumb Navigation -->
      <nav aria-label="Breadcrumb" class="flex items-center gap-2 text-xs text-gray-500 font-medium mb-6">
        <NuxtLink to="/" class="hover:text-[#1b4332] transition-colors">Home</NuxtLink>
        <span>/</span>
        <span class="text-gray-900 font-semibold">Search</span>
        <template v-if="searchQuery.trim()">
          <span>/</span>
          <span class="text-[#1b4332] font-bold truncate max-w-[200px]">"{{ searchQuery.trim() }}"</span>
        </template>
      </nav>

      <!-- Clean Header Hero matching catalog.vue & design_philosophy -->
      <div class="text-center max-w-2xl mx-auto mb-12 space-y-3">
        <span class="text-xs font-black text-emerald-700 uppercase tracking-widest bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
          Chia Florist Search
        </span>
        <h1 class="text-4xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
          Search Flower Boards
        </h1>
        <p class="text-sm md:text-base text-gray-500 leading-relaxed">
          Search across our handcrafted flower collections, greeting boards, or launch our real-time custom simulator.
        </p>

        <!-- Search Bar Form matching design_philosophy (py-3, rounded-xl, CButton primary) -->
        <form @submit.prevent="handleSearchSubmit" class="max-w-2xl mx-auto pt-4">
          <div class="flex flex-col sm:flex-row gap-2.5 items-stretch">
            <div class="relative flex-1">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </span>
              <input
                v-model="searchQuery"
                type="search"
                placeholder="Search by occasion (e.g. Wedding, Condolence, Birthday)..."
                class="w-full bg-white border border-gray-200 rounded-xl pl-11 pr-10 py-3 text-sm outline-none focus:bg-white focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all font-medium text-gray-800 shadow-2xs"
              />
              <button
                v-if="searchQuery"
                type="button"
                @click="clearSearch"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-gray-400 hover:text-gray-600 cursor-pointer"
                title="Clear search"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <CButton
              type="submit"
              variant="primary"
              size="lg"
              class="flex-shrink-0"
            >
              <span>Search</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 ml-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
              </svg>
            </CButton>
          </div>
        </form>
      </div>

      <!-- Controls & Filter Bar matching catalog.vue -->
      <div class="mb-10 bg-white border border-gray-100 rounded-xl p-6 shadow-xs flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="flex items-center gap-3">
          <h2 class="text-sm font-bold text-gray-900">
            <span v-if="searchQuery.trim()">Results for <span class="text-[#1b4332] font-black">"{{ searchQuery.trim() }}"</span></span>
            <span v-else>All Collections</span>
          </h2>
          <span class="text-xs bg-emerald-50 text-emerald-800 font-bold px-2.5 py-1 rounded-full border border-emerald-100">
            {{ totalResults }} items found
          </span>
        </div>

        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto justify-end">
          <!-- Store branch selector -->
          <div v-if="shops.length > 0" class="flex items-center gap-2">
            <label for="search-shop-filter" class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Store</label>
            <select
              id="search-shop-filter"
              :value="selectedShop"
              @change="handleShopChange(($event.target as HTMLSelectElement).value)"
              class="bg-gray-50/50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm outline-none focus:bg-white focus:border-emerald-700 transition-all font-semibold text-gray-700 cursor-pointer"
            >
              <option value="">All Stores</option>
              <option v-for="shop in shops" :key="shop.id" :value="shop.slug">
                {{ shop.name }}
              </option>
            </select>
          </div>

          <!-- Sort selector -->
          <div class="flex items-center gap-2">
            <label for="search-sort-filter" class="text-xs font-bold text-gray-500 uppercase tracking-wider whitespace-nowrap">Sort By</label>
            <select
              id="search-sort-filter"
              v-model="selectedSort"
              @change="handleSortChange"
              class="bg-gray-50/50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm outline-none focus:bg-white focus:border-emerald-700 transition-all font-semibold text-gray-700 cursor-pointer"
            >
              <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[300px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-emerald-700"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Searching our collection...</p>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="searchResults.length === 0"
        class="bg-white rounded-xl p-10 sm:p-14 text-center border border-gray-200 shadow-xs max-w-2xl mx-auto space-y-6"
      >
        <div class="w-16 h-16 bg-emerald-50 text-[#1b4332] rounded-full flex items-center justify-center mx-auto border border-emerald-100">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <div class="space-y-2">
          <h3 class="text-xl font-bold text-gray-900">No flower boards match "{{ searchQuery }}"</h3>
          <p class="text-sm text-gray-500 max-w-md mx-auto leading-relaxed">
            We couldn't find any exact matches for your search. Try checking your spelling, using broader keywords, or create a fully customized design in our simulator.
          </p>
        </div>

        <div class="pt-4 flex flex-wrap items-center justify-center gap-3">
          <CButton
            to="/products/custom"
            variant="primary"
            size="md"
          >
            <span>Launch Simulator</span>
          </CButton>
          <CButton
            to="/catalog"
            variant="outline"
            size="md"
          >
            <span>Browse Full Catalog</span>
          </CButton>
        </div>
      </div>

      <!-- Results Product Grid (Cards matching catalog.vue & design_philosophy) -->
      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
      >
        <article
          v-for="product in searchResults"
          :key="product.id"
          class="group bg-white border border-gray-200 rounded-xl overflow-hidden shadow-xs hover:shadow-md transition-all duration-200 flex flex-col justify-between cursor-pointer transform hover:-translate-y-0.5"
          @click="navigateTo(product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`)"
        >
          <div>
            <!-- Product Thumbnail Image -->
            <div class="aspect-[4/3] w-full bg-gray-50 relative overflow-hidden border-b border-gray-100">
              <img
                v-if="product.image"
                :src="product.image"
                :alt="product.name"
                loading="lazy"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              />
              <div v-else class="w-full h-full flex items-center justify-center bg-gray-100 text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>

              <!-- Tag Badge (Pill rounded-full) -->
              <span
                v-if="product.tag"
                class="absolute top-3 left-3 bg-white/95 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1 rounded-full border border-gray-100 shadow-2xs"
              >
                {{ product.tag }}
              </span>

              <!-- Out of stock badge -->
              <span
                v-if="product.isAvailable === false"
                class="absolute top-3 right-3 bg-red-600 text-white text-[10px] font-bold px-2.5 py-1 rounded-full shadow-2xs"
              >
                Out of Stock
              </span>
            </div>

            <!-- Product Body Details -->
            <div class="p-5 space-y-2">
              <div class="flex items-center gap-1.5 text-xs text-amber-500 font-bold">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 fill-current" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
                <span>{{ (product.rating || 5.0).toFixed(1) }}</span>
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

          <!-- Card Footer Price & CTA Button (matching design_philosophy) -->
          <div class="p-5 pt-0">
            <div class="pt-3 border-t border-gray-100 flex items-center justify-between">
              <div>
                <span class="text-[10px] text-gray-400 uppercase font-bold tracking-wider block">Price</span>
                <span class="text-base font-extrabold text-gray-900">
                  {{ formatRupiah(product.price) }}
                </span>
              </div>

              <CButton
                :to="product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`"
                :variant="product.isCustomRoute || product.id === 'custom' ? 'primary' : 'secondary'"
                size="sm"
                @click.stop
              >
                <span>{{ product.isCustomRoute || product.id === 'custom' ? 'Design' : 'View' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                </svg>
              </CButton>
            </div>
          </div>
        </article>
      </div>

    </div>
  </div>
</template>
