<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useCart } from '~/composables/useCart'
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'

const route = useRoute()
const productId = computed(() => route.params.id as string)
const { addToCart, formatRupiah } = useCart()

const showToast = ref(false)
const toastInfo = ref({
  name: '',
  image: '',
  quantity: 1,
  size: ''
})
let toastTimeout: ReturnType<typeof setTimeout> | null = null

onUnmounted(() => {
  if (toastTimeout) clearTimeout(toastTimeout)
})

const { currentProduct, isLoading, error, fetchProductById } = useProductViewModel()

watch(productId, (newId) => {
  if (newId === 'custom') {
    navigateTo('/products/custom', { replace: true })
    return
  }
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
    shopId: product.value.shopId || '99ef0062-1040-4574-a4be-0123abce5670',
    isCustom: false
  }, quantity.value)
  
  toastInfo.value = {
    name: product.value.name,
    image: activeImage.value || '/images/birthday.jpeg',
    quantity: quantity.value,
    size: selectedSize.value
  }
  showToast.value = true
  if (toastTimeout) clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    showToast.value = false
  }, 4000)
}

const handleBuyNow = () => {
  if (!product.value) return
  const shopId = product.value.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
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
      shopId: shopId
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
              <button 
                :disabled="!product.available"
                @click="product.available && handleAddToCart()" 
                :class="[!product.available ? 'opacity-50 cursor-not-allowed border-gray-300 text-gray-400 bg-gray-50' : 'border-2 border-[#1b4332] text-[#1b4332] bg-white hover:bg-emerald-50/50 cursor-pointer']"
                class="flex-1 font-bold py-3 rounded-xl transition text-sm"
              >
                Add to Cart
              </button>
              <button 
                :disabled="!product.available"
                @click="product.available && handleBuyNow()" 
                :class="[!product.available ? 'opacity-50 cursor-not-allowed bg-gray-300 text-gray-500' : 'bg-[#1b4332] hover:bg-[#143326] text-white cursor-pointer']"
                class="flex-1 font-bold py-3 rounded-xl transition shadow-sm text-sm"
              >
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
    
    <!-- Toast Notification -->
    <Transition name="toast">
      <div v-if="showToast" class="cart-toast" role="alert">
        <div class="toast-body">
          <div class="toast-icon-check">✓</div>
          <div class="toast-img-wrap">
            <img :src="toastInfo.image" class="toast-img" />
          </div>
          <div class="toast-details">
            <h4 class="toast-title">Added to Cart!</h4>
            <p class="toast-name">{{ toastInfo.name }}</p>
            <p class="toast-meta">Qty: {{ toastInfo.quantity }} | Size: {{ toastInfo.size }}</p>
          </div>
        </div>
        <div class="toast-actions">
          <NuxtLink to="/cart" class="btn-toast-view">View Cart</NuxtLink>
          <button @click="showToast = false" class="btn-toast-close">×</button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
/* Custom Toast Notification */
.cart-toast {
  position: fixed;
  top: 90px;
  right: 24px;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 360px;
  max-width: calc(100vw - 48px);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(27, 67, 50, 0.15);
  border-radius: 20px;
  padding: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
  font-family: 'Inter', system-ui, sans-serif;
}

.toast-body {
  display: flex;
  align-items: center;
  gap: 14px;
}

.toast-icon-check {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(27, 67, 50, 0.1);
  color: #1b4332;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
  flex-shrink: 0;
}

.toast-img-wrap {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  overflow: hidden;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  flex-shrink: 0;
}

.toast-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.toast-details {
  flex-grow: 1;
  min-width: 0;
}

.toast-title {
  font-size: 14px;
  font-weight: 800;
  color: #1b4332;
  margin: 0;
}

.toast-name {
  font-size: 13px;
  font-weight: 600;
  color: #111827;
  margin: 2px 0 0 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.toast-meta {
  font-size: 11px;
  color: #6b7280;
  margin: 2px 0 0 0;
  font-weight: 500;
}

.toast-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px dashed rgba(27, 67, 50, 0.1);
  padding-top: 12px;
}

.btn-toast-view {
  font-size: 12px;
  font-weight: 700;
  color: #1b4332;
  text-decoration: none;
  background: rgba(27, 67, 50, 0.05);
  padding: 6px 14px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.btn-toast-view:hover {
  background: #1b4332;
  color: #ffffff;
}

.btn-toast-close {
  background: transparent;
  border: none;
  font-size: 16px;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
  line-height: 1;
}

.btn-toast-close:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #374151;
}

/* Toast Transition */
.toast-enter-active, .toast-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from {
  transform: translateX(120%) scale(0.9);
  opacity: 0;
}
.toast-leave-to {
  transform: translateX(120%) scale(0.9);
  opacity: 0;
}
</style>