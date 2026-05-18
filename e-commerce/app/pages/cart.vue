<script setup lang="ts">
import { useCart } from '~/composables/useCart'

useHead({
  title: 'Your Shopping Cart - Chia Florist'
})

// Ambil fungsi & data dari global cart composable
const { cart, removeFromCart, updateQuantity, cartSubtotal } = useCart()

// Biaya pengiriman flat (simulasi)
const shippingFee = ref(10)
const promoCode = ref('')
const discount = ref(0)

const totalPayment = computed(() => {
  return cartSubtotal.value + shippingFee.value - discount.value
})

const applyPromo = () => {
  if (promoCode.value.toUpperCase() === 'CHIAFLORIST') {
    discount.value = 15
    alert('Promo code applied successfully! You got a $15 discount.')
  } else {
    alert('Invalid promo code.')
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

      <div v-if="cart.length === 0" class="bg-white border border-gray-100 rounded-3xl p-16 text-center shadow-sm max-w-xl mx-auto mt-10">
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

      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        
        <div class="lg:col-span-8 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
          <div v-for="(item, idx) in cart" :key="idx" class="flex flex-col sm:flex-row gap-6 py-6 border-b border-gray-100 last:border-0 last:pb-0 first:pt-0">
            <div class="w-full sm:w-28 h-28 rounded-2xl overflow-hidden bg-gray-50 border border-gray-100 flex-shrink-0">
              <img :src="item.image || '/images/custom-preview.png'" :alt="item.name" class="w-full h-full object-cover" />
            </div>

            <div class="flex-1 flex flex-col justify-between">
              <div class="flex justify-between items-start gap-4">
                <div>
                  <h3 class="font-bold text-gray-900 text-lg leading-snug">{{ item.name }}</h3>
                  <div class="flex flex-wrap gap-3 mt-2 text-xs text-gray-500 font-medium">
                    <span v-if="item.size" class="bg-gray-100 px-2.5 py-1 rounded-md">Size: {{ item.size }}</span>
                    <div v-if="item.color" class="flex items-center gap-1.5 bg-gray-100 px-2.5 py-1 rounded-md">
                      <span>Color:</span>
                      <span :style="{ backgroundColor: item.color }" class="w-3 h-3 rounded-full border border-gray-300 inline-block"></span>
                    </div>
                    <span v-if="item.isCustom" class="bg-emerald-50 text-emerald-700 px-2.5 py-1 rounded-md font-bold">✨ Custom Board</span>
                  </div>
                </div>
                <div class="text-xl font-bold text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</div>
              </div>

              <div class="flex justify-between items-center mt-6 pt-4 border-t border-dashed border-gray-100">
                <div class="flex border border-gray-200 rounded-lg overflow-hidden bg-gray-50">
                  <button @click="updateQuantity(item.id, item.size, item.color, -1)" class="px-3 py-1.5 hover:bg-gray-200 transition text-gray-600 font-bold">-</button>
                  <span class="px-4 py-1.5 font-bold text-gray-800 flex items-center select-none text-sm">{{ item.quantity }}</span>
                  <button @click="updateQuantity(item.id, item.size, item.color, 1)" class="px-3 py-1.5 hover:bg-gray-200 transition text-gray-600 font-bold">+</button>
                </div>

                <button @click="removeFromCart(item.id, item.size, item.color)" class="text-sm font-semibold text-red-500 hover:text-red-700 flex items-center gap-1.5 transition">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.34 9.149m-8.28 0L5.82 9m1.65-4.361a1.6 1.6 0 0 1 1.76-1.305h4.18c.83 0 1.53.551 1.76 1.305l.14 0.49M8.25 4.5h7.493M4.5 4.5h15M11.06 18h2.23" />
                  </svg>
                  Remove
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="lg:col-span-4 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Order Summary</h3>
            
            <div class="space-y-4 text-sm font-medium text-gray-600">
              <div class="flex justify-between">
                <span>Subtotal</span>
                <span class="text-gray-900 font-bold">${{ cartSubtotal.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Estimated Delivery</span>
                <span class="text-gray-900 font-bold">${{ shippingFee.toFixed(2) }}</span>
              </div>
              <div v-if="discount > 0" class="flex justify-between text-emerald-600">
                <span>Promo Discount</span>
                <span class="font-bold">-${{ discount.toFixed(2) }}</span>
              </div>
              <div class="border-t border-gray-100 pt-4 flex justify-between text-base font-bold text-gray-900">
                <span>Total Amount</span>
                <span class="text-2xl font-black text-[#1b4332]">${{ totalPayment.toFixed(2) }}</span>
              </div>
            </div>

            <div class="pt-4 border-t border-gray-50 space-y-2">
              <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Do you have a promo code?</label>
              <div class="flex gap-2">
                <input v-model="promoCode" type="text" placeholder="e.g. CHIAFLORIST" class="flex-1 bg-gray-50 border border-gray-200 rounded-xl px-4 py-3 text-sm outline-none focus:bg-white focus:border-[#1b4332] transition-all font-semibold" />
                <button @click="applyPromo" class="bg-gray-900 hover:bg-black text-white px-5 py-3 rounded-xl text-xs font-bold transition">Apply</button>
              </div>
            </div>

            <button class="w-full bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-4 rounded-xl transition shadow-md hover:shadow-lg text-center text-sm tracking-wide">
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