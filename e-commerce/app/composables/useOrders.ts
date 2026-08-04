// app/composables/useOrders.ts
import { ref, computed } from 'vue'
import { orderService } from '~/services/orderService'
import type { BackendOrder } from '~/types/order'

/**
 * Maps the 4 UI tab keys to the backend status strings.
 *   pending    -> pending
 *   processing -> confirmed | processing
 *   shipping   -> shipped
 *   done       -> delivered | cancelled
 */
export type OrderTab = 'pending' | 'processing' | 'shipping' | 'done'

const TAB_STATUSES: Record<OrderTab, string[]> = {
  pending:    ['pending'],
  processing: ['confirmed', 'processing'],
  shipping:   ['shipped'],
  done:       ['delivered', 'cancelled']
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

      // Client-side guard: 'pending' (Waiting Payment) tab only shows pending payments
      if (tab === 'pending') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status === 'pending')
      }
      // Client-side guard: 'processing' (To Ship) tab only shows non-pending payments
      if (tab === 'processing') {
        fetchedOrders = fetchedOrders.filter(o => o.payment?.status !== 'pending')
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
    formatDate
  }
}
