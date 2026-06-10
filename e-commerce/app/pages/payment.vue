<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'

useHead({
  title: 'Secure Payment - Chia Florist'
})

// 1. Ambil fungsi yang benar-benar ada di useCart
const { cartSubtotal, checkoutToOrder } = useCart()

const shippingFee = ref(10)
const selectedMethod = ref('qris')

// 2. totalPayment dideklarasikan SEKALI saja di sini sebagai computed lokal
const totalPayment = computed(() => {
  return (cartSubtotal?.value || 0) + shippingFee.value
})

// --- TIMER LOGIC (Batas Waktu Bayar 24 Jam) ---
const timeLeft = ref(86400) // 24 jam dalam detik
let timerInterval: any = null

const formattedTimer = computed(() => {
  const hours = Math.floor(timeLeft.value / 3600)
  const minutes = Math.floor((timeLeft.value % 3600) / 60)
  const seconds = timeLeft.value % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

onMounted(() => {
  timerInterval = setInterval(() => {
    if (timeLeft.value > 0) timeLeft.value--
  }, 1000)
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

// 3. handleConfirmPayment dideklarasikan SEKALI saja secara bersih
const handleConfirmPayment = () => {
  // Picu pemindahan isi cart ke daftar order tracking di profile
  if (checkoutToOrder) {
    checkoutToOrder(totalPayment.value)
  }
  
  alert('Verifying your payment... Order successfully created!')
  navigateTo('/profile')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-4xl mx-auto px-6">
      
      <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm flex flex-col md:flex-row justify-between items-center gap-6 mb-8">
        <div>
          <span class="text-xs font-bold text-gray-400 uppercase tracking-wider">Order Total</span>
          <p class="text-3xl font-black text-[#1b4332] mt-1">${{ totalPayment.toFixed(2) }}</p>
          <p class="text-xs text-gray-400 mt-2">Order ID: #CHIA-LOCAL</p>
        </div>
        
        <div class="text-center md:text-right bg-red-50 border border-red-100 px-6 py-4 rounded-2xl w-full md:w-auto">
          <span class="text-xs font-bold text-red-600 uppercase tracking-wider">Payment Time Left</span>
          <p class="text-2xl font-mono font-bold text-red-700 mt-1">{{ formattedTimer }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-12 gap-8 items-start">
        
        <div class="md:col-span-7 bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
          <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Select Payment Method</h3>
          
          <div class="space-y-3">
            <div 
              @click="selectedMethod = 'qris'"
              :class="['p-4 rounded-xl border-2 transition-all pointer-events-auto cursor-pointer flex items-center justify-between', selectedMethod === 'qris' ? 'border-[#1b4332] bg-emerald-50/20' : 'border-gray-100 hover:border-gray-200']"
            >
              <div class="flex items-center gap-3">
                <span class="text-2xl">📱</span>
                <span class="text-sm font-bold text-gray-800">QRIS (Automated Verification)</span>
              </div>
              <div class="w-4 h-4 rounded-full border-2 border-gray-300 flex items-center justify-center animate-fade" :class="{'bg-[#1b4332] border-[#1b4332]': selectedMethod === 'qris'}"></div>
            </div>

            <div 
              @click="selectedMethod = 'bank'"
              :class="['p-4 rounded-xl border-2 transition-all pointer-events-auto cursor-pointer flex items-center justify-between', selectedMethod === 'bank' ? 'border-[#1b4332] bg-emerald-50/20' : 'border-gray-100 hover:border-gray-200']"
            >
              <div class="flex items-center gap-3">
                <span class="text-2xl">🏦</span>
                <span class="text-sm font-bold text-gray-800">Bank Transfer (Manual/VA)</span>
              </div>
              <div class="w-4 h-4 rounded-full border-2 border-gray-300 flex items-center justify-center animate-fade" :class="{'bg-[#1b4332] border-[#1b4332]': selectedMethod === 'bank'}"></div>
            </div>
          </div>

          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100">
            <div v-if="selectedMethod === 'qris'" class="text-center space-y-4">
              <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Scan this QR Code to pay</p>
              <div class="w-48 h-48 bg-white border border-gray-200 rounded-2xl mx-auto flex items-center justify-center shadow-sm overflow-hidden">
                <img src="https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=ChiaFlorist" alt="QRIS Code" class="w-40 h-40 object-contain" />
              </div>
              <p class="text-xs text-gray-400 max-w-xs mx-auto">Supports GoPay, OVO, Dana, LinkAja, and all Mobile Banking apps.</p>
            </div>

            <div v-if="selectedMethod === 'bank'" class="space-y-4">
              <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Transfer details</p>
              <div class="bg-white p-4 rounded-xl border border-gray-100 space-y-2">
                <div class="flex justify-between text-sm">
                  <span class="text-gray-400 font-medium">Bank Name:</span>
                  <span class="font-bold text-gray-900">Bank Mandiri</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-gray-400 font-medium">Account Number:</span>
                  <span class="font-mono font-bold text-gray-900 select-all">137-00-123456-7</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-gray-400 font-medium">Account Holder:</span>
                  <span class="font-bold text-gray-900">CHIA FLORIST STUDIO</span>
                </div>
              </div>
              <p class="text-[11px] text-gray-400 leading-relaxed">💡 Please write your client details in the reference note for faster verification.</p>
            </div>
          </div>
        </div>

        <div class="md:col-span-5 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Payment Confirmation</h3>
            <p class="text-xs text-gray-500 leading-relaxed">By clicking the confirmation button below, you declare that you have made a legal transaction according to the billing system total amount.</p>
            <button 
              @click="handleConfirmPayment"
              class="w-full bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-4 rounded-xl transition shadow-md text-center text-sm tracking-wide"
            >
              I Have Paid the Order
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.3s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>