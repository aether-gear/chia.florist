<!-- app/pages/checkout.vue -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { useAddress } from '~/composables/useAddress'

useHead({
  title: 'Secure Checkout - Chia Florist'
})

// INTEGRASI: Ambil data keranjang, fungsi order, dan helper formatRupiah murni dari useCart()
const { cart, cartSubtotal, cartSubtotalFormatted, checkoutToOrder, formatRupiah } = useCart()
const addressVm = useAddress()

// State management internal checkout
const shippingFee = ref(20000) // SINKRONISASI KURS: Diubah ke nominal Rupiah murni (Rp20.000) agar setara dengan database Supabase
const discount = ref(0)
const selectedAddressId = ref('')
const isProcessing = ref(false)

// Jaga alur: Jika keranjang belanja kosong, langsung tendang user kembali ke halaman utama
onMounted(async () => {
  if (cart.value.length === 0) {
    navigateTo('/catalog')
    return
  }
  await addressVm.fetchAddresses()
  
  // Set otomatis ke alamat default jika tersedia di database bapak/ibu kurir
  const defaultAddr = addressVm.addresses.value.find(a => a.is_default)
  if (defaultAddr) {
    // FIX SUCCESS: Menggunakan .value untuk merubah state ref constant
    selectedAddressId.value = defaultAddr.address_id || ''
  }
})

// Menghitung total tagihan akhir secara numerik presisi
const totalPayment = computed(() => {
  return cartSubtotal.value + shippingFee.value - discount.value
})

// Eksekusi checkout memindahkan state item keranjang ke invoice profile order history
const handlePlaceOrder = async () => {
  if (!selectedAddressId.value) {
    alert('Please select a shipping address before completing your order.')
    return
  }

  isProcessing.value = true

  try {
    // Jalankan mutasi transaksi global
    checkoutToOrder(totalPayment.value)
    
    alert('Order placed successfully! Redirecting to secure payment page...')
    navigateTo('/profile') // Alihkan langsung ke tab pesanan saya di profil untuk cek status "pembayaran"
  } catch (err) {
    console.error('Checkout processing error:', err)
    alert('Failed to process checkout. Please try again.')
  } finally {
    isProcessing.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <div class="mb-10">
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Secure Checkout</h1>
        <p class="text-sm text-gray-500 mt-1">Please confirm your shipping metadata and billing totals below.</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        
        <div class="lg:col-span-7 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">1. Shipping Destination</h3>
            
            <div v-if="addressVm.isLoading.value" class="flex flex-col items-center justify-center py-6 space-y-2">
              <div class="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-[#1b4332]"></div>
              <p class="text-gray-400 text-xs">Fetching address cards...</p>
            </div>

            <div class="text-center py-6 border-2 border-dashed border-gray-200 rounded-2xl p-4" v-else-if="addressVm.addresses.value.length === 0">
              <p class="text-sm text-gray-500">No addresses registered to your profile.</p>
              <NuxtLink to="/profile" class="text-xs font-bold text-[#1b4332] underline mt-1 inline-block">Add address in Profile Settings</NuxtLink>
            </div>

            <div class="space-y-3" v-else>
              <label 
                v-for="addr in addressVm.addresses.value" 
                :key="addr.address_id"
                :class="['border rounded-2xl p-4 flex items-start gap-4 cursor-pointer transition-all', selectedAddressId === addr.address_id ? 'border-[#1b4332] bg-emerald-50/5' : 'border-gray-200 hover:border-gray-300']"
              >
                <input 
                  type="radio" 
                  v-model="selectedAddressId" 
                  :value="addr.address_id" 
                  class="mt-1 accent-[#1b4332]" 
                />
                <div class="flex-1 text-xs">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-bold text-gray-900">{{ addr.receiver_name }}</span>
                    <span class="bg-emerald-100 text-emerald-800 font-bold text-[9px] px-2 py-0.2 rounded-full" v-if="addr.is_default">Default</span>
                  </div>
                  <p class="text-gray-600 font-semibold mb-1">📞 {{ addr.phone }}</p>
                  <p class="text-gray-500 leading-normal">{{ addr.full_address }}, {{ addr.postal_code }}</p>
                </div>
              </label>
            </div>
          </div>

          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-4">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">2. Review Ordered Board Items</h3>
            <div class="divide-y divide-gray-100 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
              <div v-for="(item, idx) in cart" :key="idx" class="flex gap-4 items-center py-4 first:pt-0 last:pb-0">
                <div class="w-14 h-14 rounded-xl overflow-hidden bg-gray-50 border border-gray-100 flex-shrink-0">
                  <img :src="item.image || '/images/custom-preview.png'" :alt="item.name" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h4 class="font-bold text-gray-900 text-sm truncate">{{ item.name }}</h4>
                  <p class="text-xs text-gray-400 mt-1 font-semibold">Qty: {{ item.quantity }} | Size: {{ item.size || 'Standard' }}</p>
                </div>
                <div class="text-sm font-extrabold text-gray-900 text-right">
                  {{ formatRupiah(item.price * item.quantity) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="lg:col-span-4 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Billing Summary</h3>
            
            <div class="space-y-4 text-sm font-medium text-gray-600">
              <div class="flex justify-between items-center">
                <span>Subtotal</span>
                <span class="text-gray-900 font-bold">{{ cartSubtotalFormatted }}</span>
              </div>
              
              <div class="flex justify-between items-center">
                <span>Shipping Cost</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(shippingFee) }}</span>
              </div>
              
              <div class="flex justify-between items-center" v-if="discount > 0">
                <span>Promo Discount</span>
                <span class="font-bold">-{{ formatRupiah(discount) }}</span>
              </div>
              
              <div class="border-t border-gray-100 pt-4 flex justify-between items-center text-base font-bold text-gray-900">
                <span>Total Bill</span>
                <span class="text-2xl font-black text-[#1b4332]">
                  {{ formatRupiah(totalPayment) }}
                </span>
              </div>
            </div>

            <button 
              @click="handlePlaceOrder"
              :disabled="isProcessing || addressVm.addresses.value.length === 0"
              class="w-full bg-[#1b4332] hover:bg-[#143326] disabled:bg-gray-300 text-white font-bold py-4 rounded-xl transition shadow-md hover:shadow-lg text-center text-sm tracking-wide cursor-pointer disabled:cursor-not-allowed flex items-center justify-center"
            >
              <span v-if="isProcessing" class="animate-pulse">Processing Order...</span>
              <span v-else>Confirm & Pay Now</span>
            </button>
          </div>
          
          <div class="bg-emerald-50/50 border border-emerald-100 rounded-2xl p-4 flex gap-3 items-center text-emerald-800">
            <span class="text-2xl">🔒</span>
            <p class="text-xs font-medium leading-normal">Your payment request is fully managed under a cryptographically secured end-to-end sandbox module.</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom Scrollbar Mini untuk list item */
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
</style>