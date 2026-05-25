<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCart } from '~/composables/useCart'

useHead({
  title: 'Secure Checkout - Chia Florist'
})

const { cart, cartSubtotal } = useCart()

// State Formulir Checkout
const checkoutForm = ref({
  fullName: '',
  phone: '',
  email: '',
  deliveryDate: '',
  deliveryTime: '08:00 - 12:00',
  venueName: '',
  addressDetails: '',
  notes: ''
})

const shippingFee = ref(10)

const totalPayment = computed(() => {
  return cartSubtotal.value + shippingFee.value
})

const handlePlaceOrder = () => {
  if (!checkoutForm.value.fullName || !checkoutForm.value.venueName || !checkoutForm.value.deliveryDate) {
    alert('Please fill in all the required fields (Name, Venue, and Delivery Date).')
    return
  }
  
  // Nanti di sini kita sambungkan ke Payment Gateway
  alert('Order placed successfully! Proceeding to the payment gateway...')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <nav class="flex text-sm text-gray-400 mb-8 gap-2">
        <NuxtLink to="/cart" class="hover:text-[#1b4332]">Cart</NuxtLink>
        <span>/</span>
        <span class="text-gray-900 font-medium">Checkout</span>
      </nav>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        
        <div class="lg:col-span-7 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-8">
          
          <div>
            <h2 class="text-xl font-bold text-gray-900 mb-4">Contact Information</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Recipient Name *</label>
                <input v-model="checkoutForm.fullName" type="text" placeholder="e.g. John Doe" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
              </div>
              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Phone Number *</label>
                <input v-model="checkoutForm.phone" type="tel" placeholder="e.g. +62812345..." class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
              </div>
            </div>
            <div class="space-y-1.5 mt-4">
              <label class="text-xs font-bold text-gray-500">Email Address</label>
              <input v-model="checkoutForm.email" type="email" placeholder="john@example.com" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
            </div>
          </div>

          <div class="border-t border-gray-100 pt-6">
            <h2 class="text-xl font-bold text-gray-900 mb-4">Delivery Event Details</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Event Date *</label>
                <input v-model="checkoutForm.deliveryDate" type="date" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
              </div>
              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Preferred Time Slot</label>
                <select v-model="checkoutForm.deliveryTime" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all">
                  <option>08:00 - 12:00 (Morning)</option>
                  <option>13:00 - 17:00 (Afternoon)</option>
                  <option>18:00 - 21:00 (Evening)</option>
                </select>
              </div>
            </div>
            <div class="space-y-1.5 mt-4">
              <label class="text-xs font-bold text-gray-500">Venue / Building Name *</label>
              <input v-model="checkoutForm.venueName" type="text" placeholder="e.g. Grand Ballroom Ritz Carlton / Rumah Duka" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
            </div>
            <div class="space-y-1.5 mt-4">
              <label class="text-xs font-bold text-gray-500">Complete Address Specifications</label>
              <textarea v-model="checkoutForm.addressDetails" rows="3" placeholder="Street name, RT/RW, or specific instructions for placement..." class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all resize-none"></textarea>
            </div>
          </div>

          <div class="border-t border-gray-100 pt-6">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-gray-500">Special Notes for Florist / Driver</label>
              <input v-model="checkoutForm.notes" type="text" placeholder="e.g. Please put it near the main entrance gate" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm font-semibold transition-all" />
            </div>
          </div>

        </div>

        <div class="lg:col-span-5 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Review Order Items</h3>
            
            <div class="max-h-60 overflow-y-auto divide-y divide-gray-50 pr-2">
              <div v-for="(item, idx) in cart" :key="idx" class="flex items-center gap-4 py-3 first:pt-0">
                <div class="w-14 h-14 rounded-xl overflow-hidden border border-gray-100 bg-gray-50 flex-shrink-0">
                  <img :src="item.image" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h4 class="font-bold text-gray-900 text-sm truncate">{{ item.name }}</h4>
                  <p class="text-xs text-gray-400 mt-0.5">Qty: {{ item.quantity }}x <span v-if="item.size">({{ item.size }})</span></p>
                </div>
                <div class="text-sm font-bold text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</div>
              </div>
            </div>

            <div class="border-t border-gray-100 pt-4 space-y-3 text-sm font-medium text-gray-600">
              <div class="flex justify-between">
                <span>Subtotal</span>
                <span class="text-gray-900 font-bold">${{ cartSubtotal.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Flat Delivery Fee</span>
                <span class="text-gray-900 font-bold">${{ shippingFee.toFixed(2) }}</span>
              </div>
              <div class="border-t border-gray-100 pt-4 flex justify-between items-center text-gray-900 font-bold">
                <span>Total Payment</span>
                <span class="text-2xl font-black text-[#1b4332]">${{ totalPayment.toFixed(2) }}</span>
              </div>
            </div>

            <button 
              @click="handlePlaceOrder"
              class="w-full bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-4 rounded-xl transition shadow-md text-center text-sm tracking-wide"
            >
              Place Order & Continue to Payment
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>