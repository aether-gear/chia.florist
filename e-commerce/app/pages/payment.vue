<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { formatRupiah } from '~/utils/formatter'
import { orderService } from '~/services/orderService'

useHead({
  title: 'Secure Payment - Chia Florist'
})

interface PaymentInfo {
  orderId: string
  instruction: string
  paymentAccount?: {
    account_name: string
    account_number?: string
    phone_number?: string
    qr_string?: string
    action_url?: string
  }
  total: number
  expiresAt?: string
  status: string
}

const paymentInfoState = useState<PaymentInfo | null>('last-payment-info', () => null)
const selectedMethod = ref('qris')
const isLoading = ref(false)
const errorMsg = ref<string | null>(null)

const route = useRoute()

// --- TIMER LOGIC (Batas Waktu Bayar) ---
const timeLeft = ref(86400) // 24 jam default dalam detik
let timerInterval: any = null

const formattedTimer = computed(() => {
  const hours = Math.floor(timeLeft.value / 3600)
  const minutes = Math.floor((timeLeft.value % 3600) / 60)
  const seconds = timeLeft.value % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

onMounted(async () => {
  const orderIdFromQuery = route.query.orderId as string

  if (orderIdFromQuery) {
    isLoading.value = true
    errorMsg.value = null
    try {
      const res = await orderService.getOrderPaymentDetails(orderIdFromQuery)
      paymentInfoState.value = {
        orderId: orderIdFromQuery,
        instruction: res.instruction || '',
        paymentAccount: {
          account_name: res.account_name || res.display_name || 'Chia Florist',
          account_number: res.account_number,
          phone_number: res.phone_number,
          qr_string: res.qr_string || res.action_url
        },
        total: res.amount,
        expiresAt: res.expires_at,
        status: res.status
      }
    } catch (err: any) {
      console.error('Failed to load payment details:', err)
      errorMsg.value = err.data?.message || err.message || 'Failed to load payment details. Please check your login session.'
    } finally {
      isLoading.value = false
    }
  } else if (import.meta.client) {
    const cached = sessionStorage.getItem('chia-last-payment-info')
    if (cached) {
      try {
        paymentInfoState.value = JSON.parse(cached)
      } catch (e) {
        console.error(e)
      }
    }
  }

  // Auto-detect best payment method category
  if (paymentInfoState.value?.paymentAccount?.qr_string) {
    selectedMethod.value = 'qris'
  } else {
    selectedMethod.value = 'bank'
  }

  const updateTimer = () => {
    if (paymentInfoState.value?.expiresAt) {
      const now = new Date().getTime()
      const expiry = new Date(paymentInfoState.value.expiresAt).getTime()
      timeLeft.value = Math.max(0, Math.floor((expiry - now) / 1000))
    } else {
      if (timeLeft.value > 0) timeLeft.value--
    }
  }

  updateTimer()
  timerInterval = setInterval(() => {
    updateTimer()
  }, 1000)
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

const totalPayment = computed(() => {
  return paymentInfoState.value ? paymentInfoState.value.total : 0
})

const orderId = computed(() => {
  return paymentInfoState.value ? paymentInfoState.value.orderId : 'CHIA-LOCAL'
})

const renderMarkdown = (md: string) => {
  if (!md) return ''
  let html = md
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // Headers: # Title
  html = html.replace(/^#\s+(.+)$/gm, '<h2 class="text-xl font-bold text-gray-900 mb-4">$1</h2>')

  // Bold: **text**
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-extrabold text-emerald-800">$1</strong>')

  // Lists: - item
  html = html.replace(/^\s*-\s+(.+)$/gm, '<li class="ml-4 list-disc text-sm text-gray-700 my-1 font-semibold">$1</li>')

  // Paragraphs / Newlines
  html = html.split('\n\n').map(p => {
    if (p.startsWith('<h2') || p.startsWith('<li')) return p
    return `<p class="text-sm text-gray-600 leading-relaxed mb-3">${p.replace(/\n/g, '<br/>')}</p>`
  }).join('')

  return html
}

const instructionHtml = computed(() => {
  return paymentInfoState.value ? renderMarkdown(paymentInfoState.value.instruction) : ''
})

const paymentAccount = computed(() => {
  return paymentInfoState.value?.paymentAccount || null
})

const qrCodeUrl = computed(() => {
  const qrString = paymentAccount.value?.qr_string || 'ChiaFlorist'
  if (qrString.startsWith('data:') || qrString.startsWith('http://') || qrString.startsWith('https://')) {
    return qrString
  }
  return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrString)}`
})

const isChecking = ref(false)
const checkError = ref<string | null>(null)

const handleCheckPayment = async () => {
  const orderIdVal = orderId.value
  if (!orderIdVal || orderIdVal === 'CHIA-LOCAL') return

  isChecking.value = true
  checkError.value = null
  try {
    const res = await orderService.checkOrderPaymentStatus(orderIdVal)
    if (res.status === 'paid') {
      if (paymentInfoState.value) {
        paymentInfoState.value.status = 'paid'
      }
      
      // Update session cache if exists
      if (import.meta.client) {
        const cached = sessionStorage.getItem('chia-last-payment-info')
        if (cached) {
          try {
            const data = JSON.parse(cached)
            data.status = 'paid'
            sessionStorage.setItem('chia-last-payment-info', JSON.stringify(data))
          } catch (e) {
            console.error(e)
          }
        }
      }
      alert('Payment verified successfully! Thank you.')
    } else {
      alert(`Payment status is still pending (status: ${res.status}). If you just paid, please wait a minute and verify again.`)
    }
  } catch (err: any) {
    console.error('Failed to check payment status:', err)
    checkError.value = err.data?.message || err.message || 'Verification failed. Please try again.'
    alert(checkError.value)
  } finally {
    isChecking.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-4xl mx-auto px-6">

      <!-- Loading State -->
      <div v-if="isLoading" class="bg-white border border-gray-100 rounded-3xl p-12 text-center shadow-sm flex flex-col items-center gap-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#1b4332]"></div>
        <p class="text-sm font-bold text-gray-500">Loading payment details...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMsg" class="bg-red-50 border border-red-100 rounded-3xl p-8 text-center shadow-sm space-y-4">
        <div class="text-4xl">⚠️</div>
        <h3 class="font-bold text-red-800 text-lg">Error Loading Payment</h3>
        <p class="text-sm text-red-600 max-w-md mx-auto">{{ errorMsg }}</p>
        <button @click="navigateTo('/profile')" class="mt-2 bg-red-600 hover:bg-red-700 text-white font-bold px-6 py-2.5 rounded-xl transition text-xs cursor-pointer">
          Back to My Orders
        </button>
      </div>

      <!-- Main Payment Details -->
      <div v-else>
        <!-- Paid Success State -->
        <div v-if="paymentInfoState?.status === 'paid'" class="bg-white border border-gray-100 rounded-3xl p-8 md:p-12 text-center shadow-sm space-y-6 animate-fade">
          <div class="w-20 h-20 bg-emerald-50 rounded-full flex items-center justify-center mx-auto ring-8 ring-emerald-50/50">
            <svg class="w-10 h-10 text-emerald-600 animate-bounce-short" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-black text-[#1b4332]">Payment Successful!</h3>
            <p class="text-sm text-gray-500 max-w-md mx-auto">Thank you! Your payment has been received and verified. Our team is now preparing your flower arrangement.</p>
          </div>
          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Order ID:</span>
              <span class="font-mono font-bold text-gray-900 select-all">{{ orderId }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Amount Paid:</span>
              <span class="font-bold text-gray-900">{{ formatRupiah(totalPayment) }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Status:</span>
              <span class="px-2.5 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full border border-emerald-200">Paid &amp; Confirmed</span>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 justify-center pt-4">
            <button @click="navigateTo('/profile')" class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer shadow-sm">
              Track My Order
            </button>
            <button @click="navigateTo('/catalog')" class="border border-gray-200 hover:bg-gray-50 text-gray-700 font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer">
              Continue Shopping
            </button>
          </div>
        </div>

        <div v-else>
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm flex flex-col md:flex-row justify-between items-center gap-6 mb-8">
            <div>
              <span class="text-xs font-bold text-gray-400 uppercase tracking-wider">Order Total</span>
              <p class="text-3xl font-black text-[#1b4332] mt-1">{{ formatRupiah(totalPayment) }}</p>
              <p class="text-xs text-gray-500 font-mono mt-2">Order ID: {{ orderId }}</p>
            </div>

            <div class="text-center md:text-right bg-red-50 border border-red-100 px-6 py-4 rounded-2xl w-full md:w-auto">
              <span class="text-xs font-bold text-red-600 uppercase tracking-wider">Payment Time Left</span>
              <p class="text-2xl font-mono font-bold text-red-700 mt-1">{{ formattedTimer }}</p>
            </div>
          </div>

          <div class="gap-8 items-start">

            <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
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

              <!-- Markdown Payment Instructions from Backend -->
              <div v-if="instructionHtml" class="bg-emerald-50/10 border border-emerald-100/50 rounded-2xl p-6 mb-6 text-left shadow-sm">
                <div v-html="instructionHtml" class="prose max-w-none"></div>
              </div>

              <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100">
                <div v-if="selectedMethod === 'qris'" class="text-center space-y-4">
                  <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Scan this QR Code to pay</p>
                  <div class="w-48 h-48 bg-white border border-gray-200 rounded-2xl mx-auto flex items-center justify-center shadow-sm overflow-hidden">
                    <img :src="qrCodeUrl" alt="QRIS Code" class="w-40 h-40 object-contain" />
                  </div>
                  <p class="text-xs text-gray-400 max-w-xs mx-auto">Supports GoPay, OVO, Dana, LinkAja, and all Mobile Banking apps.</p>
                </div>

                <div v-if="selectedMethod === 'bank'" class="space-y-4">
                  <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Transfer details</p>
                  <div class="bg-white p-4 rounded-xl border border-gray-100 space-y-2">
                    <div class="flex justify-between text-sm" v-if="paymentAccount">
                      <span class="text-gray-400 font-medium">Account Holder:</span>
                      <span class="font-bold text-gray-900">{{ paymentAccount.account_name }}</span>
                    </div>
                    <div class="flex justify-between text-sm" v-if="paymentAccount && paymentAccount.account_number">
                      <span class="text-gray-400 font-medium">Account Number:</span>
                      <span class="font-mono font-bold text-gray-900 select-all">{{ paymentAccount.account_number }}</span>
                    </div>
                    <div class="flex justify-between text-sm" v-if="paymentAccount && paymentAccount.phone_number">
                      <span class="text-gray-400 font-medium">Phone Number (E-Wallet):</span>
                      <span class="font-mono font-bold text-gray-900 select-all">{{ paymentAccount.phone_number }}</span>
                    </div>

                    <template v-if="!paymentAccount">
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
                    </template>
                  </div>
                  <p class="text-[11px] text-gray-400 leading-relaxed">💡 Please write your client details in the reference note for faster verification.</p>
                </div>
              </div>

              <!-- Action Buttons to check/verify payment -->
              <div class="pt-6 border-t border-gray-100 flex flex-col sm:flex-row gap-3">
                <button
                  @click="handleCheckPayment"
                  :disabled="isChecking"
                  class="flex-1 bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-3.5 px-6 rounded-xl transition text-xs cursor-pointer shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  <span v-if="isChecking" class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></span>
                  <span>{{ isChecking ? 'Verifying Payment...' : 'I Have Paid / Verify Status' }}</span>
                </button>
                <button
                  @click="navigateTo('/profile')"
                  class="border border-gray-200 hover:bg-gray-50 text-gray-700 font-bold py-3.5 px-6 rounded-xl transition text-xs cursor-pointer"
                >
                  Pay Later / View Orders
                </button>
              </div>
            </div>

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
@keyframes bounceShort {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
.animate-bounce-short {
  animation: bounceShort 2s infinite ease-in-out;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
