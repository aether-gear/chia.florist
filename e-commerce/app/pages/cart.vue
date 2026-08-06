<!-- app/pages/cart.vue -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { useToast } from '~/composables/useToast'

useHead({
  title: 'Your Shopping Cart - Chia Florist'
})

const { cart, isLoadingCart, loadCart, removeFromCart, updateQuantity, cartSubtotal, cartSubtotalFormatted, flushCart, formatRupiah } = useCart()

const isPageLoading = ref(true)

onMounted(async () => {
  try {
    await flushCart()
    await loadCart(true)
  } finally {
    isPageLoading.value = false
  }
})

const shippingFee = ref(20000)
const promoCode = ref('')
const discount = ref(0)

const totalPayment = computed(() => {
  return cartSubtotal.value + shippingFee.value - discount.value
})

const applyPromo = () => {
  if (promoCode.value.toUpperCase() === 'CHIAFLORIST') {
    discount.value = 50000
    alert(`Promo code applied successfully! You got a ${formatRupiah(50000)} discount.`)
  } else {
    alert('Invalid promo code.')
  }
}

const handleCheckout = () => {
  navigateTo('/checkout')
}

const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  if (target) {
    target.src = '/images/custom-preview.png'
  }
}

const toast = useToast()

// ── Global remove throttle ────────────────────────────────────────────
// A single boolean locks ALL remove buttons while any one removal is running.
// This stops users from spamming remove across multiple products simultaneously.
const isRemovingAny = ref(false)
const activeRemovingId = ref<string | null>(null)

const handleRemove = async (id: string, size?: string, color?: string) => {
  if (isRemovingAny.value) return   // globally locked — block all remove buttons

  isRemovingAny.value = true
  activeRemovingId.value = id
  try {
    await removeFromCart(id, size, color)
    toast.success('Item removed', 'The product has been removed from your cart.')
  } catch (e) {
    console.error('Remove failed:', e)
    toast.error('Could not remove item', 'Please try again.')
  } finally {
    isRemovingAny.value = false
    activeRemovingId.value = null
  }
}

// ── Bug 2: Direct quantity input (max 80) ─────────────────────────────
const MAX_QTY = 80

// Temporary input string per item — tracks what the user is typing
const qtyInputs = ref<Record<string, string>>({})

const getQtyInput = (id: string, currentQty: number): string => {
  // If user is not actively editing, mirror real quantity
  return qtyInputs.value[id] !== undefined ? qtyInputs.value[id] : String(currentQty)
}

const onQtyFocus = (id: string, currentQty: number) => {
  qtyInputs.value[id] = String(currentQty)
}

const onQtyInput = (id: string, value: string) => {
  // Allow only digits while typing
  qtyInputs.value[id] = value.replace(/[^0-9]/g, '')
}

const onQtyBlur = async (id: string, size: string | undefined, color: string | undefined, currentQty: number) => {
  const raw = qtyInputs.value[id]
  delete qtyInputs.value[id]

  const parsed = parseInt(raw ?? '', 10)
  if (isNaN(parsed) || parsed === currentQty) return

  const clamped = Math.max(0, Math.min(parsed, MAX_QTY))
  const diff = clamped - currentQty
  if (diff !== 0) {
    await updateQuantity(id, size, color, diff)
  }
}

const onQtyKeydown = (e: KeyboardEvent, id: string, size: string | undefined, color: string | undefined, currentQty: number) => {
  if (e.key === 'Enter') {
    (e.target as HTMLInputElement).blur()
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <div class="mb-10">
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Shopping Cart</h1>
        <p class="text-sm text-gray-500 mt-1">Review your selections before proceeding to secure checkout.</p>
      </div>

      <!-- 1. Loading Skeleton -->
      <div v-if="isPageLoading || isLoadingCart" class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        <div class="lg:col-span-8 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
          <div v-for="n in 2" :key="n" class="flex flex-col sm:flex-row gap-6 py-6 border-b border-gray-100 last:border-0 animate-pulse">
            <div class="w-full sm:w-28 h-28 rounded-2xl bg-gray-200 flex-shrink-0"></div>
            <div class="flex-1 flex flex-col justify-between space-y-4">
              <div class="flex justify-between items-start gap-4">
                <div class="space-y-2 flex-1">
                  <div class="h-5 bg-gray-200 rounded-md w-3/4"></div>
                  <div class="h-4 bg-gray-200 rounded-md w-1/2"></div>
                </div>
                <div class="h-6 bg-gray-200 rounded-md w-24"></div>
              </div>
              <div class="flex justify-between items-center pt-4 border-t border-dashed border-gray-100">
                <div class="h-8 bg-gray-200 rounded-lg w-28"></div>
                <div class="h-6 bg-gray-200 rounded-md w-20"></div>
              </div>
            </div>
          </div>
        </div>
        <div class="lg:col-span-4 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6 animate-pulse">
          <div class="h-6 bg-gray-200 rounded-md w-1/2"></div>
          <div class="space-y-4">
            <div class="flex justify-between"><div class="h-4 bg-gray-200 rounded-md w-20"></div><div class="h-4 bg-gray-200 rounded-md w-24"></div></div>
            <div class="flex justify-between"><div class="h-4 bg-gray-200 rounded-md w-28"></div><div class="h-4 bg-gray-200 rounded-md w-20"></div></div>
            <div class="border-t border-gray-100 pt-4 flex justify-between"><div class="h-6 bg-gray-200 rounded-md w-28"></div><div class="h-6 bg-gray-200 rounded-md w-32"></div></div>
          </div>
          <div class="h-12 bg-gray-200 rounded-xl w-full"></div>
        </div>
      </div>

      <!-- 2. Empty Cart -->
      <div v-else-if="cart.length === 0" class="bg-white border border-gray-100 rounded-3xl p-16 text-center shadow-sm max-w-xl mx-auto mt-10">
        <div class="w-24 h-24 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-6 text-[#1b4332]">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5V6a3.75 3.75 0 1 0-7.5 0v4.5m11.356-1.993 1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 0 1-1.12-1.243l1.264-12A1.125 1.125 0 0 1 5.513 7.5h12.974c.576 0 1.059.435 1.119 1.007ZM8.625 10.5a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm7.5 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-gray-900">Your cart is empty</h2>
        <p class="text-sm text-gray-500 mt-2 max-w-sm mx-auto">Looks like you haven't added any beautiful arrangements to your cart yet.</p>
        <NuxtLink to="/" class="inline-block mt-8 bg-[#1b4332] hover:bg-[#143326] text-white text-sm font-bold px-8 py-3.5 rounded-xl transition shadow-sm">
          Continue Shopping
        </NuxtLink>
      </div>

      <!-- 3. Populated Cart -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        
        <div class="lg:col-span-8 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
          <div v-for="(item, idx) in cart" :key="item.id || idx" class="flex flex-col sm:flex-row gap-6 py-6 border-b border-gray-100 last:border-0 last:pb-0 first:pt-0">
            <div class="w-full sm:w-28 h-28 rounded-2xl overflow-hidden bg-gray-50 border border-gray-100 flex-shrink-0 relative">
              <img 
                :src="item.image || '/images/custom-preview.png'" 
                :alt="item.name" 
                class="w-full h-full object-cover"
                @error="handleImageError"
              />
            </div>

            <div class="flex-1 flex flex-col justify-between min-w-0">
              <div class="flex flex-col sm:flex-row justify-between items-start gap-4">
                <div class="min-w-0 flex-1">
                  <h3 class="font-bold text-gray-900 text-lg leading-snug break-words">{{ item.name }}</h3>
                  <div class="flex flex-wrap gap-2.5 mt-2 text-xs text-gray-500 font-medium">
                    <span v-if="item.size" class="bg-gray-100 px-2.5 py-1 rounded-md">Size: {{ item.size }}</span>
                    <div v-if="item.color" class="flex items-center gap-1.5 bg-gray-100 px-2.5 py-1 rounded-md">
                      <span>Color:</span>
                      <span :style="{ backgroundColor: item.color }" class="w-3 h-3 rounded-full border border-gray-300 inline-block"></span>
                    </div>
                    <span v-if="item.isCustom" class="bg-emerald-50 text-emerald-700 px-2.5 py-1 rounded-md font-bold">✨ Custom Board</span>
                  </div>
                </div>
                <div class="text-xl font-extrabold text-gray-900 flex-shrink-0">
                  {{ formatRupiah(item.subtotal ?? (item.price * item.quantity)) }}
                </div>
              </div>

              <div class="flex justify-between items-center mt-6 pt-4 border-t border-dashed border-gray-100">

                <!-- ── Quantity controls with direct input ── -->
                <div class="flex items-center border border-gray-200 rounded-lg overflow-hidden bg-gray-50">
                  <button
                    @click="updateQuantity(item.id, item.size, item.color, -1)"
                    class="px-3 py-1.5 hover:bg-gray-200 transition text-gray-600 font-bold cursor-pointer select-none"
                    :disabled="item.quantity <= 1"
                    title="Decrease quantity"
                  >−</button>

                  <!-- Direct quantity input -->
                  <input
                    type="text"
                    inputmode="numeric"
                    :value="getQtyInput(item.id, item.quantity)"
                    @focus="onQtyFocus(item.id, item.quantity)"
                    @input="onQtyInput(item.id, ($event.target as HTMLInputElement).value)"
                    @blur="onQtyBlur(item.id, item.size, item.color, item.quantity)"
                    @keydown="onQtyKeydown($event, item.id, item.size, item.color, item.quantity)"
                    class="w-12 py-1.5 text-center font-bold text-gray-800 text-sm bg-white border-x border-gray-200 outline-none focus:ring-1 focus:ring-inset focus:ring-[#1b4332]"
                    :title="`Max quantity: ${MAX_QTY}`"
                  />

                  <button
                    @click="updateQuantity(item.id, item.size, item.color, 1)"
                    class="px-3 py-1.5 hover:bg-gray-200 transition text-gray-600 font-bold cursor-pointer select-none"
                    :disabled="item.quantity >= MAX_QTY"
                    title="Increase quantity"
                  >+</button>
                </div>

                <!-- ── Remove button: ALL buttons lock while any removal runs ── -->
                <button
                  @click="handleRemove(item.id, item.size, item.color)"
                  :disabled="isRemovingAny"
                  class="text-sm font-semibold flex items-center gap-1.5 transition cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
                  :class="isRemovingAny ? 'text-gray-400' : 'text-red-500 hover:text-red-700'"
                >
                  <!-- Spinner only on the item being removed -->
                  <svg v-if="activeRemovingId === item.id" class="h-4 w-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"/>
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.34 9.149m-8.28 0L5.82 9m1.65-4.361a1.6 1.6 0 0 1 1.76-1.305h4.18c.83 0 1.53.551 1.76 1.305l.14 0.49M8.25 4.5h7.493M4.5 4.5h15M11.06 18h2.23" />
                  </svg>
                  {{ activeRemovingId === item.id ? 'Removing…' : 'Remove' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Order Summary -->
        <div class="lg:col-span-4 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Order Summary</h3>
            
            <div class="space-y-4 text-sm font-medium text-gray-600">
              <div class="flex justify-between items-center">
                <span>Subtotal</span>
                <span class="text-gray-900 font-bold">{{ cartSubtotalFormatted }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span>Estimated Delivery</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(shippingFee) }}</span>
              </div>
              <div class="flex justify-between items-center text-emerald-600" v-if="discount > 0">
                <span>Promo Discount</span>
                <span class="font-bold">-{{ formatRupiah(discount) }}</span>
              </div>
              <div class="border-t border-gray-100 pt-4 flex justify-between items-center text-base font-bold text-gray-900">
                <span>Total Amount</span>
                <span class="text-2xl font-extrabold text-[#1b4332]">{{ formatRupiah(totalPayment) }}</span>
              </div>
            </div>

            <div class="pt-4 border-t border-gray-50 space-y-2">
              <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Do you have a promo code?</label>
              <div class="flex gap-2">
                <input v-model="promoCode" type="text" placeholder="e.g. CHIAFLORIST" class="flex-1 bg-gray-50 border border-gray-200 rounded-xl px-4 py-3 text-sm outline-none focus:bg-white focus:border-[#1b4332] transition-all font-semibold" />
                <button @click="applyPromo" class="bg-gray-900 hover:bg-black text-white px-5 py-3 rounded-xl text-xs font-bold transition cursor-pointer">Apply</button>
              </div>
            </div>

            <button 
              @click="handleCheckout"
              class="w-full bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-4 rounded-xl transition shadow-md hover:shadow-lg text-center text-sm tracking-wide cursor-pointer"
            >
              Proceed to Checkout
            </button>
          </div>
          
          <div class="bg-emerald-50/50 border border-emerald-100 rounded-2xl p-4 flex gap-3 items-center text-emerald-800">
            <span class="text-2xl">🔒</span>
            <p class="text-xs font-medium leading-normal">Secure Checkout Guaranteed. Your data is encrypted and completely safe with us.</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-enter-active, .slide-leave-active {
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from, .slide-leave-to {
  transform: translateX(100%);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.4s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
</style>