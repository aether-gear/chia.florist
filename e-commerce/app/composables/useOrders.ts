// app/composables/useOrders.ts
import { ref, computed } from 'vue'
import { orderService } from '~/services/orderService'
import type { BackendOrder, GetOrderTrackingTimelineResponse } from '~/types/order'
import { mapErrorMessage } from '~/utils/errorMessages'

export type OrderTab = 'all' | 'pending' | 'processing' | 'shipping' | 'completed' | 'cancelled'


const TAB_STATUSES: Record<OrderTab, string[]> = {
  all:        [],
  pending:    ['pending'],
  processing: ['confirmed', 'processing'],
  shipping:   ['shipped'],
  completed:  ['delivered', 'finished'],
  cancelled:  ['cancelled', 'expired']
}

/**
 * Checks if an order's payment has expired.
 * An order is NEVER expired if it is paid or in an active fulfillment state.
 */
export const isOrderExpired = (order: BackendOrder): boolean => {
  if (!order) return false

  // Explicit backend order status expired
  if (order.status === 'expired') return true

  // Paid orders or active fulfillment states are NEVER expired
  if (order.payment?.status === 'paid') return false
  const nonExpiredStatuses = ['confirmed', 'processing', 'shipped', 'delivered', 'finished']
  if (nonExpiredStatuses.includes(order.status)) return false

  // Explicitly marked expired or failed payment
  if (order.payment?.status === 'expired' || order.payment?.status === 'failed') {
    return true
  }

  // Pending payment past expires_at timestamp
  if (order.payment?.expires_at && (order.status === 'pending' || order.payment?.status === 'pending')) {
    const expiresAt = new Date(order.payment.expires_at).getTime()
    return Date.now() >= expiresAt
  }

  return false
}

/**
 * Generates Shopee-style step-by-step timeline tracking status for an order.
 */
export const getOrderTimelineSteps = (order: BackendOrder) => {
  const isExpired = isOrderExpired(order)
  const isCancelled = order.status === 'cancelled' || order.payment?.status === 'cancelled'

  if (isCancelled || isExpired) {
    return [
      { step: 1, title: 'Order Placed', desc: 'Order submitted to system', done: true, active: false },
      { step: 2, title: isExpired ? 'Payment Expired' : 'Order Cancelled', desc: isExpired ? 'Payment window lapsed' : 'Order has been cancelled', done: true, error: true, active: true }
    ]
  }

  const isPaid = order.payment?.status === 'paid' || ['confirmed', 'processing', 'shipped', 'delivered', 'finished'].includes(order.status)
  const isProcessing = ['processing', 'shipped', 'delivered', 'finished'].includes(order.status)
  const isShipped = ['shipped', 'delivered', 'finished'].includes(order.status)
  const isDelivered = ['delivered', 'finished'].includes(order.status)

  return [
    { step: 1, title: 'Order Placed', desc: 'Awaiting payment verification', done: true, active: order.status === 'pending' && !isPaid },
    { step: 2, title: 'Payment Verified', desc: isPaid ? 'Payment received & confirmed' : 'Awaiting payment transfer', done: isPaid, active: isPaid && order.status === 'confirmed' },
    { step: 3, title: 'Arranging Flowers', desc: isProcessing ? 'Florists preparing bouquet' : 'Waiting for workshop', done: isProcessing, active: order.status === 'processing' },
    { step: 4, title: 'In Transit', desc: isShipped ? 'Package handed over to courier' : 'Awaiting shipment dispatch', done: isShipped, active: order.status === 'shipped' },
    { step: 5, title: 'Order Received', desc: isDelivered ? 'Successfully delivered to recipient' : 'Pending final delivery', done: isDelivered, active: isDelivered }
  ]
}

/**
 * Formats order status into human-readable label and Tailwind badge classes.
 */
export const getOrderStatusBadge = (order: BackendOrder) => {
  if (isOrderExpired(order)) {
    return { label: 'Expired Payment', colorClass: 'bg-rose-50 text-rose-700 border-rose-100' }
  }
  switch (order.status) {
    case 'pending':
      return { label: 'Waiting Payment', colorClass: 'bg-amber-50 text-amber-700 border-amber-100' }
    case 'confirmed':
      return { label: 'Order Confirmed', colorClass: 'bg-indigo-50 text-indigo-700 border-indigo-100' }
    case 'processing':
      return { label: 'Arranging Flowers', colorClass: 'bg-emerald-50 text-emerald-700 border-emerald-100' }
    case 'shipped':
      return { label: 'In Transit', colorClass: 'bg-blue-50 text-blue-700 border-blue-100' }
    case 'delivered':
      return { label: 'Completed', colorClass: 'bg-sky-50 text-sky-700 border-sky-100' }
    case 'finished':
      return { label: 'Completed', colorClass: 'bg-sky-50 text-sky-700 border-sky-100' }
    case 'cancelled':
      return { label: 'Cancelled', colorClass: 'bg-red-50 text-red-700 border-red-100' }
    case 'expired':
      return { label: 'Expired', colorClass: 'bg-rose-50 text-rose-700 border-rose-100' }
    default:
      return { label: order.status || 'Unknown', colorClass: 'bg-gray-100 text-gray-600 border-gray-200' }
  }
}

/**
 * Formats payment status into human-readable label and Tailwind badge classes.
 */
export const getPaymentStatusBadge = (status?: string) => {
  switch (status) {
    case 'pending':
      return { label: 'Payment Pending', colorClass: 'bg-amber-50 text-amber-700 border-amber-100' }
    case 'paid':
      return { label: 'Paid & Verified', colorClass: 'bg-emerald-50 text-emerald-700 border-emerald-100' }
    case 'failed':
      return { label: 'Payment Failed', colorClass: 'bg-red-50 text-red-700 border-red-100' }
    case 'expired':
      return { label: 'Payment Expired', colorClass: 'bg-rose-50 text-rose-700 border-rose-100' }
    case 'cancelled':
      return { label: 'Payment Cancelled', colorClass: 'bg-gray-100 text-gray-700 border-gray-200' }
    case 'refunded':
      return { label: 'Refunded', colorClass: 'bg-purple-50 text-purple-700 border-purple-100' }
    case 'refund_pending':
      return { label: 'Refund Pending', colorClass: 'bg-amber-50 text-amber-700 border-amber-100' }
    case 'refund_failed':
      return { label: 'Refund Failed', colorClass: 'bg-red-50 text-red-700 border-red-100' }
    default:
      return { label: status || 'Unspecified', colorClass: 'bg-gray-100 text-gray-600 border-gray-200' }
  }
}

/**
 * Formats time remaining for a pending payment until it expires.
 */
export const getTimeRemaining = (expiresAt?: string): string => {
  if (!expiresAt) return ''
  const diff = new Date(expiresAt).getTime() - Date.now()
  if (diff <= 0) return 'Expired'

  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  if (minutes > 0) {
    return `${minutes}m ${seconds}s`
  }
  return `${seconds}s`
}

export const useOrders = () => {
  const orders      = ref<BackendOrder[]>([])
  const isLoading   = ref(false)
  const error       = ref<string | null>(null)

  // Pagination
  const currentPage  = ref(1)
  const pageSize     = ref(10)
  const totalOrders  = ref(0)
  const totalPages   = computed(() =>
    Math.ceil(totalOrders.value / pageSize.value) || 1
  )

  /**
   * Fetch orders for a given tab using server-side multi-status filtering.
   */
  const fetchOrders = async (tab: OrderTab, page = 1) => {
    isLoading.value = true
    error.value     = null
    currentPage.value = page

    const statuses = TAB_STATUSES[tab]

    try {
      const query: any = {
        page,
        limit: pageSize.value,
        sort: 'latest:desc'
      }

      if (statuses && statuses.length > 0) {
        query.status = statuses.join(',')
      }

      const res = await orderService.listOrders(query)

      let fetchedOrders = res.orders ?? []

      // Client-side guard filters per tab category
      if (tab === 'pending') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status === 'pending' && !isOrderExpired(o))
      } else if (tab === 'cancelled') {
        fetchedOrders = fetchedOrders.filter(o => isOrderExpired(o) || o.status === 'cancelled' || o.payment?.status === 'cancelled')
      } else if (tab === 'processing') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status !== 'pending' && !isOrderExpired(o) && ['confirmed', 'processing'].includes(o.status))
      } else if (tab === 'shipping') {
        fetchedOrders = fetchedOrders.filter(o => !isOrderExpired(o) && o.status === 'shipped')
      } else if (tab === 'completed') {
        fetchedOrders = fetchedOrders.filter(o => !isOrderExpired(o) && ['delivered', 'finished'].includes(o.status))
      }

      orders.value      = fetchedOrders
      totalOrders.value = res.total ?? fetchedOrders.length
    } catch (err: any) {
      error.value = mapErrorMessage(err, 'Gagal memuat pesanan. Silakan coba lagi.')
      orders.value      = []
      totalOrders.value = 0
    } finally {
      isLoading.value = false
    }
  }

  const goToPage = (tab: OrderTab, page: number) => {
    if (page >= 1 && page <= totalPages.value) {
      fetchOrders(tab, page)
    }
  }

  /** Convenience: format currency */
  const formatRupiah = (amount: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount)
  }

  /** Format ISO date to locale string */
  const formatDate = (iso: string): string => {
    if (!iso) return '-'
    return new Date(iso).toLocaleDateString('en-US', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  }

  const trackingData       = ref<GetOrderTrackingTimelineResponse | null>(null)
  const isTrackingLoading  = ref(false)
  const trackingError      = ref<string | null>(null)

  const fetchOrderTracking = async (orderId: string) => {
    if (!orderId) return
    isTrackingLoading.value = true
    trackingError.value     = null
    try {
      const res = await orderService.getOrderTrackingTimeline(orderId)
      trackingData.value = res
    } catch (err: any) {
      trackingError.value = mapErrorMessage(err, 'Gagal memuat data pelacakan pengiriman.')
      trackingData.value = null
    } finally {
      isTrackingLoading.value = false
    }
  }

  return {
    orders,
    isLoading,
    error,
    currentPage,
    pageSize,
    totalOrders,
    totalPages,
    fetchOrders,
    goToPage,
    trackingData,
    isTrackingLoading,
    trackingError,
    fetchOrderTracking,
    formatRupiah,
    formatDate,
    isOrderExpired,
    getOrderTimelineSteps,
    getOrderStatusBadge,
    getPaymentStatusBadge,
    getTimeRemaining
  }
}

