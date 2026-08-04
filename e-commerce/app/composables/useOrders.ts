// app/composables/useOrders.ts
import { ref, computed } from 'vue'
import { orderService } from '~/services/orderService'
import type { BackendOrder } from '~/types/order'

/**
 * Maps the 5 UI tab keys to the backend status strings.
 *   pending    -> pending (active unexpired payment)
 *   expired    -> pending | cancelled (payment expired or failed)
 *   processing -> confirmed | processing
 *   shipping   -> shipped
 *   done       -> delivered | cancelled (completed / normal user cancellations)
 */
export type OrderTab = 'pending' | 'expired' | 'processing' | 'shipping' | 'done'

const TAB_STATUSES: Record<OrderTab, string[]> = {
  pending:    ['pending'],
  expired:    ['pending', 'cancelled'],
  processing: ['confirmed', 'processing'],
  shipping:   ['shipped'],
  done:       ['delivered', 'cancelled']
}

/**
 * Checks if an order's payment has expired.
 * An order is NEVER expired if it is paid or in an active fulfillment state.
 */
export const isOrderExpired = (order: BackendOrder): boolean => {
  if (!order) return false

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
      const res = await orderService.listOrders({
        page,
        limit: pageSize.value,
        sort: 'latest:desc',
        status: statuses.join(',')
      })

      let fetchedOrders = res.orders ?? []

      // Client-side guard: 'pending' (Waiting Payment) tab only shows active unexpired pending payments
      if (tab === 'pending') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status === 'pending' && !isOrderExpired(o))
      }
      // Client-side guard: 'expired' (Expired Payment) tab shows orders with expired payments
      if (tab === 'expired') {
        fetchedOrders = fetchedOrders.filter(o => isOrderExpired(o))
      }
      // Client-side guard: 'processing' (To Ship) tab only shows non-pending payments
      if (tab === 'processing') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status !== 'pending' && !isOrderExpired(o))
      }
      // Client-side guard: 'done' tab excludes expired orders if listed under cancelled
      if (tab === 'done') {
        fetchedOrders = fetchedOrders.filter(o => !isOrderExpired(o))
      }

      orders.value      = fetchedOrders
      totalOrders.value = res.total ?? fetchedOrders.length
    } catch (err: any) {
      error.value = err?.data?.message || err?.message || 'Failed to load orders'
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
    formatRupiah,
    formatDate,
    isOrderExpired,
    getTimeRemaining
  }
}
