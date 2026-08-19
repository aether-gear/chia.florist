<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { formatRupiah } from '~/utils/formatter'
import { orderService } from '~/services/orderService'
import { mapErrorMessage } from '~/utils/errorMessages'

useHead({
  title: 'Secure Payment - Chia Florist'
})

interface PaymentInfo {
  orderId: string
  instruction: string
  channelData?: {
    channel_type: string
    display_name: string
    action_url?: string
    expires_at?: string
  }
  total: number
  expiresAt?: string
  status: string
}

const paymentInfoState = useState<PaymentInfo | null>('last-payment-info', () => null)
const isLoading = ref(false)
const errorMsg = ref<string | null>(null)
const globalAlert = useGlobalAlert()

const route = useRoute()

// --- TIMER & POLLING LOGIC ---
const timeLeft = ref(86400) // 24 hours default in seconds
let timerInterval: any = null
let pollInterval: any = null

const formattedTimer = computed(() => {
  const hours = Math.floor(timeLeft.value / 3600)
  const minutes = Math.floor((timeLeft.value % 3600) / 60)
  const seconds = timeLeft.value % 60
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

const stopPolling = () => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

const startPolling = () => {
  stopPolling()
  if (!import.meta.client) return

  // Poll payment status silently every 30 seconds while pending
  pollInterval = setInterval(async () => {
    if (paymentInfoState.value?.status !== 'pending') {
      stopPolling()
      return
    }
    const orderIdVal = orderId.value
    if (!orderIdVal || orderIdVal === 'CHIA-LOCAL') return

    try {
      const res = await orderService.checkOrderPaymentStatus(orderIdVal)
      if (res.status !== 'pending' && paymentInfoState.value) {
        paymentInfoState.value.status = res.status
        stopPolling()
        if (res.status === 'paid') {
          globalAlert.showSuccess(
            'Payment Verified!',
            'Your payment has been received and confirmed.',
            [
              { label: 'View Order', onClick: () => navigateTo('/profile') },
              { label: 'Got it' }
            ]
          )
        } else if (res.status === 'expired' || res.status === 'cancelled') {
          globalAlert.showWarning('Order ' + res.status, 'The payment status for this order is now ' + res.status + '.')
        }
      }
    } catch (e) {
      // Ignore background poll errors
    }
  }, 30000)
}

const handleTimerZero = async () => {
  stopPolling()
  const orderIdVal = orderId.value
  if (orderIdVal && orderIdVal !== 'CHIA-LOCAL') {
    try {
      const res = await orderService.checkOrderPaymentStatus(orderIdVal)
      if (paymentInfoState.value) {
        paymentInfoState.value.status = res.status === 'pending' ? 'expired' : res.status
      }
    } catch (e) {
      if (paymentInfoState.value) {
        paymentInfoState.value.status = 'expired'
      }
    }
  } else if (paymentInfoState.value) {
    paymentInfoState.value.status = 'expired'
  }
  globalAlert.showWarning(
    'Payment Window Expired',
    'The time allocated for completing your payment has lapsed.',
    [
      { label: 'Browse Catalog', onClick: () => navigateTo('/catalog') },
      { label: 'My Orders', onClick: () => navigateTo('/profile') }
    ]
  )
}

onMounted(async () => {
  const orderIdFromQuery = route.query.orderId as string

  if (orderIdFromQuery) {
    isLoading.value = true
    errorMsg.value = null
    try {
      const res = await orderService.getOrderPaymentDetails(orderIdFromQuery)
      
      // Compute status upfront so we don't render instructions if already expired
      let initialStatus = res.status || 'pending'
      if (initialStatus !== 'paid' && res.expires_at) {
        const isPast = new Date().getTime() >= new Date(res.expires_at).getTime()
        if (isPast) {
          initialStatus = 'expired'
        }
      }

      paymentInfoState.value = {
        orderId: orderIdFromQuery,
        instruction: res.instruction || '',
        channelData: res.channel_data ? {
          channel_type: res.channel_data.channel_type,
          display_name: res.channel_data.display_name || 'Payment Gateway',
          action_url: res.channel_data.action_url
        } : (res.channel_type ? {
          channel_type: res.channel_type,
          display_name: res.display_name || 'Payment Gateway',
          action_url: res.action_url
        } : undefined),
        total: res.amount,
        expiresAt: res.expires_at,
        status: initialStatus
      }
    } catch (err: any) {
      console.error('Failed to load payment details:', err)
      errorMsg.value = mapErrorMessage(err, 'Failed to load payment details. Please check your login session.')
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

  const updateTimer = () => {
    if (paymentInfoState.value?.expiresAt) {
      const now = new Date().getTime()
      const expiry = new Date(paymentInfoState.value.expiresAt).getTime()
      const diff = Math.max(0, Math.floor((expiry - now) / 1000))
      timeLeft.value = diff

      if (diff <= 0 && paymentInfoState.value.status === 'pending') {
        handleTimerZero()
      }
    } else {
      if (timeLeft.value > 0) {
        timeLeft.value--
        if (timeLeft.value <= 0 && paymentInfoState.value?.status === 'pending') {
          handleTimerZero()
        }
      }
    }
  }

  updateTimer()
  timerInterval = setInterval(() => {
    updateTimer()
  }, 1000)

  if (paymentInfoState.value?.status === 'pending') {
    startPolling()
  }
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
  stopPolling()
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

const qrCodeUrl = computed(() => {
  const actionUrl = paymentInfoState.value?.channelData?.action_url
  if (!actionUrl) return null
  return `https://api.qrserver.com/v1/create-qr-code/?size=250x250&data=${encodeURIComponent(actionUrl)}`
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
      stopPolling()

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
      globalAlert.showSuccess(
        'Payment Verified!',
        'Thank you! Your payment has been received and confirmed.',
        [
          { label: 'View Order', onClick: () => navigateTo('/profile') },
          { label: 'Got it' }
        ]
      )
    } else if (res.status === 'expired' || res.status === 'cancelled') {
      if (paymentInfoState.value) {
        paymentInfoState.value.status = res.status
      }
      stopPolling()
      globalAlert.showWarning('Order ' + res.status, `Your payment status is currently ${res.status}.`)
    } else {
      globalAlert.showInfo('Payment Pending', 'Payment status is still pending. If you just transferred, please allow a moment for confirmation.')
    }
  } catch (err: any) {
    console.error('Failed to check payment status:', err)
    checkError.value = mapErrorMessage(err, 'Verification failed. Please try again.')
    globalAlert.showError('Verification Error', checkError.value)
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

        <!-- Expired State -->
        <div v-else-if="paymentInfoState?.status === 'expired'" class="bg-white border border-gray-100 rounded-3xl p-8 md:p-12 text-center shadow-sm space-y-6 animate-fade">
          <div class="w-20 h-20 bg-amber-50 rounded-full flex items-center justify-center mx-auto ring-8 ring-amber-50/50">
            <span class="text-4xl">⏰</span>
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-black text-amber-900">Payment Window Expired</h3>
            <p class="text-sm text-gray-500 max-w-md mx-auto">The payment time limit for this order has lapsed and reserved stock has been released.</p>
          </div>
          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Order ID:</span>
              <span class="font-mono font-bold text-gray-900 select-all">{{ orderId }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Amount:</span>
              <span class="font-bold text-gray-900">{{ formatRupiah(totalPayment) }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Status:</span>
              <span class="px-2.5 py-0.5 bg-amber-100 text-amber-800 text-[10px] font-bold rounded-full border border-amber-200">Expired</span>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 justify-center pt-4">
            <button @click="navigateTo('/catalog')" class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer shadow-sm">
              Browse Catalog / Re-Order
            </button>
            <button @click="navigateTo('/profile')" class="border border-gray-200 hover:bg-gray-50 text-gray-700 font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer">
              Back to My Orders
            </button>
          </div>
        </div>

        <!-- Cancelled State -->
        <div v-else-if="paymentInfoState?.status === 'cancelled'" class="bg-white border border-gray-100 rounded-3xl p-8 md:p-12 text-center shadow-sm space-y-6 animate-fade">
          <div class="w-20 h-20 bg-red-50 rounded-full flex items-center justify-center mx-auto ring-8 ring-red-50/50">
            <span class="text-4xl">🚫</span>
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-black text-red-900">Order Cancelled</h3>
            <p class="text-sm text-gray-500 max-w-md mx-auto">This order has been cancelled. You can rebuild your arrangement or place a new order from our catalog.</p>
          </div>
          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Order ID:</span>
              <span class="font-mono font-bold text-gray-900 select-all">{{ orderId }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Status:</span>
              <span class="px-2.5 py-0.5 bg-red-100 text-red-800 text-[10px] font-bold rounded-full border border-red-200">Cancelled</span>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 justify-center pt-4">
            <button @click="navigateTo('/catalog')" class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer shadow-sm">
              Re-Order Flower Arrangement
            </button>
            <button @click="navigateTo('/profile')" class="border border-gray-200 hover:bg-gray-50 text-gray-700 font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer">
              View Order History
            </button>
          </div>
        </div>

        <!-- Failed State -->
        <div v-else-if="paymentInfoState?.status === 'failed'" class="bg-white border border-gray-100 rounded-3xl p-8 md:p-12 text-center shadow-sm space-y-6 animate-fade">
          <div class="w-20 h-20 bg-red-50 rounded-full flex items-center justify-center mx-auto ring-8 ring-red-50/50">
            <span class="text-4xl">❌</span>
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-black text-red-900">Payment Failed</h3>
            <p class="text-sm text-gray-500 max-w-md mx-auto">Your payment transaction was declined or failed to process. Please try checking out again with another payment method.</p>
          </div>
          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Order ID:</span>
              <span class="font-mono font-bold text-gray-900 select-all">{{ orderId }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Status:</span>
              <span class="px-2.5 py-0.5 bg-red-100 text-red-800 text-[10px] font-bold rounded-full border border-red-200">Failed</span>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 justify-center pt-4">
            <button @click="navigateTo('/checkout')" class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer shadow-sm">
              Try Checkout Again
            </button>
            <button @click="navigateTo('/profile')" class="border border-gray-200 hover:bg-gray-50 text-gray-700 font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer">
              Back to My Orders
            </button>
          </div>
        </div>

        <!-- Refunded / Refund Pending State -->
        <div v-else-if="paymentInfoState?.status === 'refunded' || paymentInfoState?.status === 'refund_pending'" class="bg-white border border-gray-100 rounded-3xl p-8 md:p-12 text-center shadow-sm space-y-6 animate-fade">
          <div class="w-20 h-20 bg-purple-50 rounded-full flex items-center justify-center mx-auto ring-8 ring-purple-50/50">
            <span class="text-4xl">💸</span>
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-black text-purple-900">{{ paymentInfoState?.status === 'refunded' ? 'Payment Refunded' : 'Refund Pending' }}</h3>
            <p class="text-sm text-gray-500 max-w-md mx-auto">Your payment for this order is being processed for refund. Please check your bank/e-wallet account statement.</p>
          </div>
          <div class="bg-gray-50 rounded-2xl p-6 border border-gray-100 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Order ID:</span>
              <span class="font-mono font-bold text-gray-900 select-all">{{ orderId }}</span>
            </div>
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-gray-400 font-medium">Status:</span>
              <span class="px-2.5 py-0.5 bg-purple-100 text-purple-800 text-[10px] font-bold rounded-full border border-purple-200">{{ paymentInfoState?.status === 'refunded' ? 'Refunded' : 'Refund Pending' }}</span>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row gap-3 justify-center pt-4">
            <button @click="navigateTo('/profile')" class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold px-6 py-3 rounded-xl transition text-xs cursor-pointer shadow-sm">
              View My Orders
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
              <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Payment Method</h3>

              <!-- Payment Channel Info Card -->
              <div v-if="paymentInfoState?.channelData" class="flex flex-col justify-between p-6 rounded-2xl border border-gray-100 bg-gray-50/30 gap-6">
                <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 w-full">
                  <div class="flex items-center gap-3">
                    <div class="w-12 h-12 bg-[#1b4332]/5 rounded-xl flex items-center justify-center text-xl">
                      <span v-if="paymentInfoState.channelData.channel_type === 'ewallet'">📱</span>
                      <span v-else-if="paymentInfoState.channelData.channel_type === 'bank_transfer'">🏦</span>
                      <span v-else-if="paymentInfoState.channelData.channel_type === 'qr_code'">🔍</span>
                      <span v-else>💳</span>
                    </div>
                    <div>
                      <span class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Selected Provider</span>
                      <h4 class="font-bold text-gray-900 text-sm mt-0.5">{{ paymentInfoState.channelData.display_name }}</h4>
                    </div>
                  </div>
                  <!-- Action URL Redirection button if action_url is present -->
                  <div v-if="paymentInfoState.channelData.action_url" class="w-full md:w-auto">
                    <a
                      :href="paymentInfoState.channelData.action_url"
                      target="_blank"
                      class="w-full md:w-auto inline-flex items-center justify-center gap-2 bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-2.5 px-5 rounded-xl transition text-xs shadow-sm"
                    >
                      Proceed to Payment Page ↗
                    </a>
                  </div>
                </div>

                <!-- Direct QR Code Display for QRIS / GoPay / ShopeePay -->
                <div 
                  v-if="qrCodeUrl && (paymentInfoState.channelData.channel_type === 'qr_code' || paymentInfoState.channelData.channel_type === 'ewallet' || paymentInfoState.channelData.display_name.toLowerCase().includes('gopay') || paymentInfoState.channelData.display_name.toLowerCase().includes('qris'))" 
                  class="border-t border-gray-100 pt-6 flex flex-col items-center text-center space-y-4"
                >
                  <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Scan this QR Code with your phone to pay</p>
                  <div class="w-56 h-56 bg-white border border-gray-100 rounded-3xl flex items-center justify-center shadow-sm overflow-hidden p-4">
                    <img :src="qrCodeUrl" alt="Payment QR Code" class="w-48 h-48 object-contain" />
                  </div>
                  <p class="text-[10px] text-gray-400 max-w-xs leading-relaxed">
                    Supports GoPay, OVO, Dana, LinkAja, ShopeePay, and all QRIS-compliant Mobile Banking apps.
                  </p>
                </div>
              </div>

              <!-- Markdown Payment Instructions from Backend -->
              <div v-if="instructionHtml" class="bg-emerald-50/10 border border-emerald-100/50 rounded-2xl p-6 mb-6 text-left shadow-sm">
                <div v-html="instructionHtml" class="prose max-w-none"></div>
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
