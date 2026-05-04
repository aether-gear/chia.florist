<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const productId = computed(() => route.params.id as string)

// 1. Definisikan Interface untuk data produk agar TypeScript tidak bingung
interface Product {
  id: string
  name: string
  price: number
  rating: number
  reviews: number
  available: boolean
  description: string
  images: string[]
  colors: string[]
  sizes: string[]
}

// 2. Data dummy produk dengan tipe Product[]
const productsData = ref<Product[]>([
  {
    id: 'birthday',
    name: 'Birthday Board',
    price: 192,
    rating: 4.5,
    reviews: 150,
    available: true,
    description: 'This is a beautiful and elegant birthday board, perfect for making your loved ones feel special on their big day!',
    images: [
      '/images/birthday.jpeg', 
      '/images/wedding.jpeg',
      '/images/condolences.jpeg',
      '/images/graduate.jpeg'
    ],
    colors: ['#cbd5e1', '#f43f5e'],
    sizes: ['1.5m', '1.8m', '2m', '2.5m', '3m']
  },
  {
    id: 'wedding',
    name: 'Wedding Flower Board',
    price: 250,
    rating: 4.8,
    reviews: 210,
    available: true,
    description: 'An exquisite flower board crafted for luxurious wedding celebrations.',
    images: [
      '/images/wedding.jpeg',
      '/images/congratulations.jpeg',
      '/images/grandop.jpeg',
      '/images/anniversary.jpeg'
    ],
    colors: ['#ffffff', '#f43f5e'],
    sizes: ['1.5m', '1.8m', '2m']
  }
])

// 3. Mengambil data produk yang aktif sesuai URL secara aman
const product = computed<Product>(() => {
  const found = productsData.value.find((p) => p.id === productId.value)
  if (found) return found
  
  // Fallback object agar TypeScript yakin data tidak akan pernah undefined
  return productsData.value[0] || {
    id: 'fallback',
    name: 'Product Not Found',
    price: 0,
    rating: 0,
    reviews: 0,
    available: false,
    description: '',
    images: [''],
    colors: [''],
    sizes: ['']
  }
})

// State Interaktif Halaman
const activeImage = ref('')
const selectedColor = ref('')
const selectedSize = ref('')
const quantity = ref(1)

// Mengawasi perubahan produk agar state ter-update otomatis tanpa error saat ganti URL
watch(product, (newProduct) => {
  if (newProduct && newProduct.images && newProduct.images.length > 0) {
    activeImage.value = newProduct.images[0] || ''
    selectedColor.value = newProduct.colors[0] || ''
    selectedSize.value = newProduct.sizes[1] || newProduct.sizes[0] || ''
    quantity.value = 1
  }
}, { immediate: true })

// Update gambar utama saat thumbnail diklik
const changeActiveImage = (img: string) => {
  activeImage.value = img
}

useHead({
  title: computed(() => `Chia Florist - ${product.value.name}`),
  meta: [
    { name: 'description', content: computed(() => product.value.description) }
  ]
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-12 mt-10 font-sans">
    
    <nav class="text-sm text-gray-500 mb-12 flex gap-2">
      <NuxtLink to="/" class="hover:text-black transition">Account</NuxtLink>
      <span>/</span>
      <NuxtLink to="/products" class="hover:text-black transition">Boards</NuxtLink>
      <span>/</span>
      <span class="text-black font-medium">{{ product.name }}</span>
    </nav>

    <div class="grid grid-cols-1 md:grid-cols-12 gap-12">
      
      <div class="md:col-span-7 flex flex-col-reverse md:flex-row gap-6">
        
        <div class="flex md:flex-col gap-4 overflow-x-auto md:overflow-visible">
          <button 
            v-for="(img, idx) in product.images" 
            :key="idx"
            @click="changeActiveImage(img)"
            :class="[
              'w-20 h-20 flex-shrink-0 border-2 rounded-lg overflow-hidden transition-all duration-300 transform',
              activeImage === img ? 'border-[#1b4332] scale-105 shadow-sm' : 'border-gray-100 hover:border-gray-300'
            ]"
          >
            <img :src="img" :alt="`${product.name} Thumbnail ${idx + 1}`" class="w-full h-full object-cover" />
          </button>
        </div>

        <div class="flex-1 h-[550px] bg-gray-50 rounded-xl overflow-hidden border border-gray-100">
          <img :src="activeImage" :alt="product.name" class="w-full h-full object-cover transition-opacity duration-300" />
        </div>

      </div>

      <div class="md:col-span-5 space-y-6">
        
        <div class="space-y-2">
          <h1 class="text-3xl font-bold text-gray-900 tracking-tight">{{ product.name }}</h1>
          
          <div class="flex items-center gap-4 text-sm">
            <div class="flex items-center text-yellow-400">
              <svg v-for="i in 5" :key="i" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" :fill="i <= Math.floor(product.rating) ? 'currentColor' : 'none'" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.518 4.674a1 1 0 00.95.69h4.907c.969 0 1.371 1.24.588 1.81l-3.974 2.886a1 1 0 00-.363 1.118l1.518 4.674c.3.921-.755 1.688-1.538 1.118l-3.974-2.886a1 1 0 00-1.176 0l-3.974 2.886c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.974-2.886c-.783-.57-.38-1.81.588-1.81h4.907a1 1 0 00.95-.69l1.518-4.674z" />
              </svg>
              <span class="text-gray-400 ml-2">({{ product.reviews }} Reviews)</span>
            </div>
            <span class="text-gray-300">|</span>
            <span :class="product.available ? 'text-green-600 font-medium' : 'text-red-500 font-medium'">
              {{ product.available ? 'Available' : 'Out of Stock' }}
            </span>
          </div>
        </div>

        <div class="text-3xl font-semibold text-gray-900">
          ${{ product.price.toFixed(2) }}
        </div>

        <p class="text-gray-600 text-sm leading-relaxed border-b border-gray-100 pb-6">
          {{ product.description }}
        </p>

        <div class="space-y-3">
          <label class="text-sm font-semibold text-gray-800">Colours:</label>
          <div class="flex gap-3">
            <button 
              v-for="color in product.colors" 
              :key="color"
              @click="selectedColor = color"
              :style="{ backgroundColor: color }"
              :class="[
                'w-7 h-7 rounded-full border-2 transition transform',
                selectedColor === color ? 'border-black scale-110 shadow' : 'border-transparent'
              ]"
            ></button>
          </div>
        </div>

        <div class="space-y-3">
          <label class="text-sm font-semibold text-gray-800">Size:</label>
          <div class="flex gap-2.5">
            <button 
              v-for="size in product.sizes" 
              :key="size"
              @click="selectedSize = size"
              :class="[
                'min-w-[42px] h-[34px] px-2 rounded border text-xs font-semibold flex items-center justify-center transition',
                selectedSize === size ? 'bg-[#1b4332] text-white border-[#1b4332]' : 'bg-white border-gray-300 text-gray-700 hover:border-gray-500'
              ]"
            >
              {{ size }}
            </button>
          </div>
        </div>

        <div class="flex items-center gap-4 pt-4 border-t border-gray-100">
          <div class="flex border border-gray-300 rounded overflow-hidden">
            <button @click="quantity > 1 ? quantity-- : null" class="px-4 py-2 hover:bg-gray-50 text-gray-600 font-medium">-</button>
            <span class="px-4 py-2 font-semibold text-gray-800 flex items-center select-none">{{ quantity }}</span>
            <button @click="quantity++" class="px-4 py-2 hover:bg-gray-50 text-gray-600 font-medium">+</button>
          </div>

          <button class="flex-grow bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-3 px-6 rounded-lg transition shadow-sm hover:shadow-md">
            Buy Now
          </button>

          <button class="p-3 border border-gray-300 rounded-lg hover:bg-gray-50 hover:border-gray-400 transition text-gray-600">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </svg>
          </button>
        </div>

        <div class="border border-gray-200 rounded-xl divide-y divide-gray-200 mt-8 overflow-hidden bg-white">
          <div class="flex items-center gap-4 p-4">
            <div class="text-gray-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
              </svg>
            </div>
            <div>
              <h4 class="text-sm font-bold text-gray-900">Free Delivery</h4>
              <p class="text-xs text-gray-500 underline cursor-pointer">Enter your postal code for Delivery Availability</p>
            </div>
          </div>
          
          <div class="flex items-center gap-4 p-4">
            <div class="text-gray-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </div>
            <div>
              <h4 class="text-sm font-bold text-gray-900">Return Delivery</h4>
              <p class="text-xs text-gray-500">Free 30 Days Delivery Returns. <span class="underline cursor-pointer">Details</span></p>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>