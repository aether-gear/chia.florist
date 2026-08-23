<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useCart } from '~/composables/useCart'
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'
import { useStoreSelection } from '~/composables/useStoreSelection'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { productService } from '~/services/productService'

const route = useRoute()
const productId = computed(() => route.params.id as string)
const { addToCart, formatRupiah } = useCart()
const storeSelection = useStoreSelection()
const globalAlert = useGlobalAlert()

// Redirect if custom
if (productId.value === 'custom') {
  await navigateTo('/products/custom', { replace: true })
}

// Server-side pre-fetching for SEO & fast first paint
const { data: ssrProductData } = await useAsyncData(`product-detail-${productId.value}`, async () => {
  if (!productId.value || productId.value === 'custom') return null
  try {
    const [activeList, prod] = await Promise.all([
      storeSelection.fetchActiveShops(),
      productService.getProductById(productId.value, storeSelection.selectedShop.value?.id).catch(() => null)
    ])
    return {
      shops: activeList || [],
      product: prod || null
    }
  } catch (err) {
    console.error('Failed to load product detail on SSR:', err)
    return null
  }
})

// If product was not found, throw 404 so Nuxt renders the dedicated error page
if (productId.value && productId.value !== 'custom' && !ssrProductData.value?.product) {
  throw createError({
    statusCode: 404,
    statusMessage: 'Produk Tidak Ditemukan',
    fatal: true
  })
}

const { currentProduct: vmProduct, isLoading: isVmLoading, error: vmError, fetchProductById } = useProductViewModel()

// Shops metadata map
const shopsList = ref<{ id: string; name: string; slug: string }[]>(ssrProductData.value?.shops || [])
const selectedShopSlug = ref('')

const fetchShops = async () => {
  try {
    const activeList = await storeSelection.fetchActiveShops()
    if (activeList) {
      shopsList.value = activeList
    }
  } catch (e) {
    console.error('Failed to fetch shops for product detail:', e)
  }
}

// Product computed falls back to SSR preloaded data if ViewModel hasn't fetched yet
const product = computed(() => vmProduct.value || ssrProductData.value?.product || null)
const isLoading = computed(() => isVmLoading.value && !product.value)
const error = computed(() => vmError.value && !product.value ? vmError.value : null)

watch(productId, (newId) => {
  if (newId === 'custom') {
    navigateTo('/products/custom', { replace: true })
    return
  }
  if (newId) {
    fetchShops()
    fetchProductById(newId)
  }
})

const activeImage = ref(product.value?.images?.[0] || '')
const selectedColor = ref(product.value?.colors?.[0] || '')
const selectedSize = ref('1.8m') // Default ukuran tengah standar
const quantity = ref(1)

// Filter only branches with active stock (> 0) for this product
const inStockBranches = computed(() => {
  const avail = (product.value as any)?.availability
  if (Array.isArray(avail)) {
    return avail.filter((a: any) => a.stock > 0)
  }
  return []
})

const selectAvailableBranchForProduct = () => {
  const availableList = inStockBranches.value
  if (!availableList || availableList.length === 0) {
    selectedShopSlug.value = ''
    return
  }

  const globalSlug = storeSelection.selectedShop.value?.slug
  // Only pre-select IF the customer has a global store selected AND that store has stock > 0 for this product
  if (globalSlug) {
    const globalMatch = availableList.find((a: any) => a.name === globalSlug)
    if (globalMatch) {
      selectedShopSlug.value = globalMatch.name
      return
    }
  }

  // If no global store is selected, or if the global store is out of stock for this product:
  // DO NOT auto-select any branch for the user! Keep selection empty.
  selectedShopSlug.value = ''
}

watch(product, (newProduct) => {
  if (newProduct) {
    if (!activeImage.value) activeImage.value = newProduct.images[0] || ''
    if (!selectedColor.value) selectedColor.value = newProduct.colors[0] || ''
    selectedSize.value = '1.8m'
    quantity.value = 1

    selectAvailableBranchForProduct()
  }
}, { immediate: true })

// Watch global store selection and refresh product detail with store context
watch(storeSelection.selectedShop, (newShop) => {
  if (productId.value && productId.value !== 'custom') {
    fetchProductById(productId.value, newShop?.id)
  }
  if (product.value) {
    selectAvailableBranchForProduct()
  }
})

// Dynamic SEO and OpenGraph for product page
const pageTitle = computed(() => product.value ? `${product.value.name} — Chia Florist` : 'Flower Board Details — Chia Florist')
const pageDescription = computed(() => product.value?.description || 'Beli papan bunga ucapan premium berkualitas terbaik di Chia Florist.')
const pageImage = computed(() => product.value?.images?.[0] || 'https://chiaflorist.com/florist.jpg')

useHead({
  title: pageTitle,
  meta: [
    { name: 'description', content: pageDescription },
    { property: 'og:title', content: pageTitle },
    { property: 'og:description', content: pageDescription },
    { property: 'og:type', content: 'product' },
    { property: 'og:image', content: pageImage },
    { property: 'og:url', content: computed(() => `https://chiaflorist.com/products/${productId.value}`) },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: pageTitle },
    { name: 'twitter:description', content: pageDescription },
    { name: 'twitter:image', content: pageImage },
    { name: 'robots', content: 'index, follow' }
  ],
  link: [
    { rel: 'canonical', href: computed(() => `https://chiaflorist.com/products/${productId.value}`) }
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: computed(() => {
        if (!product.value) return '{}'
        return JSON.stringify({
          '@context': 'https://schema.org',
          '@graph': [
            {
              '@type': 'Product',
              '@id': `https://chiaflorist.com/products/${productId.value}#product`,
              'name': product.value.name,
              'description': product.value.description,
              'image': product.value.images || [pageImage.value],
              'sku': product.value.sku || `CHIA-${product.value.id}`,
              'offers': {
                '@type': 'Offer',
                'url': `https://chiaflorist.com/products/${productId.value}`,
                'priceCurrency': 'IDR',
                'price': product.value.price,
                'itemCondition': 'https://schema.org/NewCondition',
                'availability': product.value.available ? 'https://schema.org/InStock' : 'https://schema.org/OutOfStock',
                'seller': {
                  '@type': 'Organization',
                  'name': 'Chia Florist'
                }
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
                  'name': 'Catalog',
                  'item': 'https://chiaflorist.com/catalog'
                },
                {
                  '@type': 'ListItem',
                  'position': 3,
                  'name': product.value.name,
                  'item': `https://chiaflorist.com/products/${productId.value}`
                }
              ]
            }
          ]
        })
      })
    }
  ]
})

const branchWarning = ref<string | null>(null)

const handleSelectFulfillingBranch = (avail: { name: string; stock: number }) => {
  // Do not allow selecting empty/out-of-stock shops
  if (!avail || avail.stock <= 0) return

  branchWarning.value = null
  selectedShopSlug.value = avail.name
  const matched = shopsList.value.find(s => s.slug === avail.name)
  if (matched && storeSelection.selectedShop.value?.id !== matched.id) {
    storeSelection.selectShop(matched)
  }
}

// Helper to resolve chosen shop_id
const resolvedShopId = computed(() => {
  if (selectedShopSlug.value) {
    const match = shopsList.value.find(s => s.slug === selectedShopSlug.value)
    if (match) return match.id
  }
  return ''
})

// Helper to compute active stock level for selected branch
const selectedBranchStock = computed(() => {
  if (selectedShopSlug.value && (product.value as any)?.availability) {
    const match = (product.value as any).availability.find((a: any) => a.name === selectedShopSlug.value)
    if (match) return match.stock
  }
  return 0
})

// FIX DINAMIS: Kalkulasi perubahan harga berdasarkan modifikasi ukuran (Size)
const displayPrice = computed(() => {
  if (!product.value) return 0
  const basePrice = Number(product.value.price)
  
  if (selectedSize.value === '1.5m') {
    return basePrice - 20000 // Jika ukuran kecil, potong Rp20.000
  } else if (selectedSize.value === '2m') {
    return basePrice + 30000 // Jika ukuran jumbo, tambah Rp30.000
  }
  return basePrice
})

const handleAddToCart = () => {
  if (!product.value) return
  if (!selectedShopSlug.value || !resolvedShopId.value) {
    branchWarning.value = 'Please select a fulfilling store branch above before adding to cart.'
    return
  }
  branchWarning.value = null

  addToCart({
    id: product.value.id,
    name: product.value.name,
    price: displayPrice.value, // Kirim harga ter-kalkulasi baru ke keranjang
    image: activeImage.value,
    size: selectedSize.value,
    color: selectedColor.value,
    shopId: resolvedShopId.value,
    isCustom: false
  }, quantity.value)
  
  globalAlert.showSuccess(
    'Added to Cart',
    `${product.value.name} (Qty: ${quantity.value}) has been added to your shopping cart.`,
    [
      { label: 'View Cart', onClick: () => navigateTo('/cart') },
      { label: 'Continue' }
    ]
  )
}

const handleBuyNow = () => {
  if (!product.value) return
  if (!selectedShopSlug.value || !resolvedShopId.value) {
    branchWarning.value = 'Please select a fulfilling store branch above before proceeding.'
    return
  }
  branchWarning.value = null

  navigateTo({
    path: '/checkout',
    query: {
      buyNow: 'true',
      id: product.value.id,
      name: product.value.name,
      price: displayPrice.value.toString(),
      image: activeImage.value,
      size: selectedSize.value,
      color: selectedColor.value,
      qty: quantity.value.toString(),
      shopId: resolvedShopId.value
    }
  })
}

useHead({
  title: computed(() => product.value ? `Chia Florist - ${product.value.name}` : 'Chia Florist - Loading...')
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-12 mt-10 font-sans">
    
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[400px] space-y-4">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#1b4332]"></div>
      <p class="text-gray-500 font-medium animate-pulse text-sm">Loading product details...</p>
    </div>

    <div v-else-if="product && product.status === 'archived'" class="flex flex-col items-center justify-center min-h-[400px] space-y-4 text-center animate-fade-in">
      <div class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0l-3-3m3 3l3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
        </svg>
      </div>
      <h3 class="text-lg font-bold text-gray-800">Product No Longer Available</h3>
      <p class="text-gray-500 text-sm max-w-md">This item has been archived and is no longer listed in our store catalog.</p>
      <NuxtLink to="/catalog" class="bg-[#1b4332] text-white px-5 py-2.5 rounded-xl hover:bg-[#143326] transition font-semibold text-xs inline-block">
        Back to Catalog
      </NuxtLink>
    </div>

    <div v-else-if="product" class="animate-fade-in">
      <nav class="text-sm text-gray-500 mb-12 flex gap-2">
        <NuxtLink to="/" class="hover:text-black transition">Home</NuxtLink>
        <span>/</span>
        <NuxtLink to="/catalog" class="hover:text-black transition">Catalog</NuxtLink>
        <span>/</span>
        <span class="text-black font-medium">{{ product.name }}</span>
      </nav>

      <div class="grid grid-cols-1 md:grid-cols-12 gap-12">
        
        <div class="md:col-span-7 flex flex-col-reverse md:flex-row gap-6">
          <div class="flex md:flex-col gap-4">
            <button 
              v-for="(img, idx) in product.images" 
              :key="idx"
              @click="activeImage = img"
              :class="['w-20 h-20 border-2 rounded-lg overflow-hidden transition-all', activeImage === img ? 'border-[#1b4332] scale-105 shadow-sm' : 'border-gray-100']"
            >
              <img :src="img" class="w-full h-full object-cover" />
            </button>
          </div>
          <div class="flex-1 h-[500px] bg-gray-50 rounded-xl overflow-hidden border border-gray-100">
            <img :src="activeImage" class="w-full h-full object-cover" />
          </div>
        </div>

        <div class="md:col-span-5 space-y-6">
          <div>
            <h1 class="text-3xl font-bold text-gray-900 tracking-tight">{{ product.name }}</h1>
            <div class="flex items-center gap-3 mt-2 flex-wrap">
              <span v-if="product.status === 'inactive'" class="text-amber-800 font-extrabold text-xs bg-amber-50 px-3 py-1 rounded-full border border-amber-200">
                Preview Only — Not For Sale
              </span>
              <span v-else-if="selectedShopSlug && selectedBranchStock > 0" class="text-emerald-700 font-extrabold text-xs bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200 flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
                </svg>
                In Stock ({{ selectedBranchStock }} available)
              </span>
              <span v-else-if="selectedShopSlug && selectedBranchStock <= 0" class="text-red-700 font-extrabold text-xs bg-red-50 px-3 py-1 rounded-full border border-red-200">
                Sold Out at Selected Store
              </span>
              <span v-else class="text-amber-700 font-extrabold text-xs bg-amber-50 px-3 py-1 rounded-full border border-amber-200 flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-amber-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
                </svg>
                Select Fulfilling Branch Below
              </span>
            </div>
          </div>

          <div class="text-3xl font-extrabold text-gray-900">
            {{ formatRupiah(displayPrice) }}
          </div>
          <p class="text-gray-600 text-sm leading-relaxed border-b border-gray-100 pb-6">{{ product.description }}</p>

          <!-- Inactive Product Notice Banner -->
          <div v-if="product.status === 'inactive'" class="p-4 bg-amber-50 border border-amber-200 rounded-2xl text-xs font-semibold text-amber-800 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-amber-700 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
            </svg>
            <span>This product is currently available for preview only and cannot be ordered online at this time.</span>
          </div>

          <!-- Store / Branch Availability Selection (In-Stock Branches Only) -->
          <div v-if="inStockBranches.length > 0 && product.status !== 'inactive'" class="space-y-3 pt-2">
            <label class="text-sm font-semibold text-gray-800 flex items-center justify-between">
              <span>Fulfilling Branch:</span>
              <span class="text-xs font-normal text-gray-500">Available in-stock stores</span>
            </label>
            <div class="space-y-2">
              <div 
                v-for="avail in inStockBranches" 
                :key="avail.name"
                @click="handleSelectFulfillingBranch(avail)"
                :class="[
                  selectedShopSlug === avail.name ? 'border-[#1b4332] bg-emerald-50/60 ring-1 ring-[#1b4332] cursor-pointer' : 
                  'border-gray-200 bg-white hover:border-emerald-300 cursor-pointer'
                ]"
                class="p-3 rounded-xl border flex items-center justify-between transition-all group"
              >
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-gray-500 group-hover:text-emerald-700 transition-colors shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 21v-7.5a.75.75 0 01.75-.75h3a.75.75 0 01.75.75V21m-4.5 0H2.25A2.25 2.25 0 010 18.75V10.5m13.5 10.5h7.5a2.25 2.25 0 002.25-2.25V10.5M3 10.5l9-7.5 9 7.5" />
                  </svg>
                  <span class="text-xs font-bold text-gray-800 group-hover:text-emerald-800 transition-colors">{{ avail.slug }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs font-bold text-emerald-700">
                    {{ avail.stock }} in stock
                  </span>
                  <span v-if="selectedShopSlug === avail.name" class="text-xs font-bold text-[#1b4332]">✓</span>
                </div>
              </div>
            </div>
            <!-- Inline Warning Below Branch List -->
            <div 
              v-if="branchWarning || !selectedShopSlug" 
              class="mt-2.5 p-3 rounded-xl flex items-center gap-2 text-xs font-semibold transition-all"
              :class="branchWarning ? 'bg-red-50 text-red-700 border border-red-200 shadow-2xs' : 'bg-amber-50 text-amber-800 border border-amber-200/80'"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" :class="branchWarning ? 'text-red-600' : 'text-amber-600'" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
              </svg>
              <span>{{ branchWarning || 'Please select a fulfilling store branch above to check stock and purchase.' }}</span>
            </div>
          </div>

          <!-- Out of Stock Message if no branches have inventory -->
          <div v-else-if="product && product.status !== 'inactive' && (product.available === false || inStockBranches.length === 0)" class="p-4 bg-red-50 border border-red-200 rounded-2xl text-xs font-bold text-red-700 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-red-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <span>This product is currently out of stock across all store branches.</span>
          </div>

          <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 pt-4 border-t border-gray-100">
            <div class="flex border border-gray-300 rounded-xl overflow-hidden bg-gray-50 flex-shrink-0 justify-between items-center w-full sm:w-auto">
              <button @click="quantity > 1 ? quantity-- : null" :disabled="product.status === 'inactive'" class="px-4 py-2.5 hover:bg-gray-200 transition font-bold text-gray-600 disabled:opacity-50">-</button>
              <span class="px-4 py-2.5 font-semibold text-gray-800 text-sm select-none">{{ quantity }}</span>
              <button @click="quantity++" :disabled="product.status === 'inactive'" class="px-4 py-2.5 hover:bg-gray-200 transition font-bold text-gray-600 disabled:opacity-50">+</button>
            </div>

            <div class="flex-1 flex gap-3 w-full">
              <button 
                :disabled="!product.available || product.status === 'inactive' || (!!selectedShopSlug && selectedBranchStock <= 0)"
                @click="handleAddToCart()" 
                :class="[(!product.available || product.status === 'inactive' || (!!selectedShopSlug && selectedBranchStock <= 0)) ? 'opacity-50 cursor-not-allowed border-gray-300 text-gray-400 bg-gray-50' : 'border-2 border-[#1b4332] text-[#1b4332] bg-white hover:bg-emerald-50/50 cursor-pointer']"
                class="flex-1 font-bold py-3 rounded-xl transition text-sm"
              >
                {{ product.status === 'inactive' ? 'Not For Sale' : (selectedShopSlug && selectedBranchStock <= 0) ? 'Out of Stock' : 'Add to Cart' }}
              </button>
              <button 
                :disabled="!product.available || product.status === 'inactive' || (!!selectedShopSlug && selectedBranchStock <= 0)"
                @click="handleBuyNow()" 
                :class="[(!product.available || product.status === 'inactive' || (!!selectedShopSlug && selectedBranchStock <= 0)) ? 'opacity-50 cursor-not-allowed bg-gray-300 text-gray-500' : 'bg-[#1b4332] hover:bg-[#143326] text-white cursor-pointer']"
                class="flex-1 font-bold py-3 rounded-xl transition shadow-sm text-sm"
              >
                {{ product.status === 'inactive' ? 'Preview Only' : 'Buy Now' }}
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>

    <!-- Fallback Not Found state if product is null -->
    <div v-else class="min-h-[60vh] flex items-center justify-center">
      <CErrorDisplay :status-code="404" />
    </div>
  </div>
</template>