<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useCart } from '~/composables/useCart'
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'

const route = useRoute()
const productId = computed(() => route.params.id as string)
const { addToCart, formatRupiah } = useCart()

const { currentProduct, isLoading, error, fetchProductById } = useProductViewModel()

watch(productId, (newId) => {
  if (newId) {
    fetchProductById(newId)
  }
}, { immediate: true })

const product = computed(() => currentProduct.value)

const activeImage = ref('')
const selectedColor = ref('')
const selectedSize = ref('1.8m') // Default ukuran tengah standar
const quantity = ref(1)

watch(product, (newProduct) => {
  if (newProduct) {
    activeImage.value = newProduct.images[0] || ''
    selectedColor.value = newProduct.colors[0] || ''
    selectedSize.value = '1.8m'
    quantity.value = 1
  }
}, { immediate: true })

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
  addToCart({
    id: product.value.id,
    name: product.value.name,
    price: displayPrice.value, // Kirim harga ter-kalkulasi baru ke keranjang
    image: activeImage.value,
    size: selectedSize.value,
    color: selectedColor.value,
    isCustom: false
  }, quantity.value)
  navigateTo('/cart')
}

const handleBuyNow = () => {
  if (!product.value) return
  addToCart({
    id: product.value.id,
    name: product.value.name,
    price: displayPrice.value,
    image: activeImage.value,
    size: selectedSize.value,
    color: selectedColor.value,
    isCustom: false
  }, quantity.value)
  navigateTo('/checkout')
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

    <div v-else-if="error" class="flex flex-col items-center justify-center min-h-[400px] space-y-4">
      <span class="text-4xl">⚠️</span>
      <h3 class="text-lg font-bold text-gray-800">Unable to load product</h3>
      <p class="text-gray-500 text-sm max-w-md text-center">{{ error }}</p>
      <button @click="fetchProductById(productId)" class="bg-[#1b4332] text-white px-5 py-2.5 rounded-xl hover:bg-[#143326] transition font-semibold text-xs cursor-pointer">
        Try Again
      </button>
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
            <p class="text-sm text-gray-400 mt-2">
              ({{ product.reviews || 150 }} Reviews) | 
              <span v-if="product.available !== false" class="text-green-600 font-medium">Available</span>
              <span class="text-red-600 font-medium" v-else>Sold Out</span>
            </p>
          </div>

          <div class="text-3xl font-extrabold text-gray-900">
            {{ formatRupiah(displayPrice) }}
          </div>
          <p class="text-gray-600 text-sm leading-relaxed border-b border-gray-100 pb-6">{{ product.description }}</p>

          <div class="space-y-3">
            <label class="text-sm font-semibold text-gray-800">Colours:</label>
            <div class="flex gap-3">
              <button 
                v-for="color in product.colors" 
                :key="color"
                @click="selectedColor = color"
                :style="{ backgroundColor: color }"
                :class="['w-7 h-7 rounded-full border-2 transition', selectedColor === color ? 'border-black scale-110 shadow' : 'border-transparent']"
              ></button>
            </div>
          </div>

          <div class="space-y-3">
            <label class="text-sm font-semibold text-gray-800">Size:</label>
            <div class="flex gap-2.5">
              <button 
                v-for="size in ['1.5m', '1.8m', '2m']" 
                :key="size"
                @click="selectedSize = size"
                :class="['min-w-[42px] h-[34px] px-3 rounded border text-xs font-semibold transition', selectedSize === size ? 'bg-[#1b4332] text-white border-[#1b4332]' : 'bg-white border-gray-300 text-gray-700']"
              >
                {{ size }}
              </button>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 pt-4 border-t border-gray-100">
            <div class="flex border border-gray-300 rounded-xl overflow-hidden bg-gray-50 flex-shrink-0 justify-between items-center w-full sm:w-auto">
              <button @click="quantity > 1 ? quantity-- : null" class="px-4 py-2.5 hover:bg-gray-200 transition font-bold text-gray-600">-</button>
              <span class="px-4 py-2.5 font-semibold text-gray-800 text-sm select-none">{{ quantity }}</span>
              <button @click="quantity++" class="px-4 py-2.5 hover:bg-gray-200 transition font-bold text-gray-600">+</button>
            </div>

            <div class="flex-1 flex gap-3 w-full">
              <button @click="handleAddToCart" class="flex-1 border-2 border-[#1b4332] text-[#1b4332] bg-white hover:bg-emerald-50/50 font-bold py-3 rounded-xl transition text-sm cursor-pointer">
                Add to Cart
              </button>
              <button @click="handleBuyNow" class="flex-1 bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-3 rounded-xl transition shadow-sm text-sm cursor-pointer">
                Buy Now
              </button>
            </div>
          </div>

          <div class="border border-gray-200 rounded-xl divide-y divide-gray-200 bg-white">
            <div class="flex items-center gap-4 p-4">
              <span class="text-2xl">📦</span>
              <div>
                <h4 class="text-sm font-bold text-gray-900">Free Delivery</h4>
                <p class="text-xs text-gray-500 underline cursor-pointer">Enter your postal code for availability</p>
              </div>
            </div>
            <div class="flex items-center gap-4 p-4">
              <span class="text-2xl">🔄</span>
              <div>
                <h4 class="text-sm font-bold text-gray-900">Return Delivery</h4>
                <p class="text-xs text-gray-500">Free 30 Days Delivery Returns.</p>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>