// app/composables/useOrders.ts
import { ref, computed } from 'vue'
import { orderService } from '~/services/orderService'
import type { BackendOrder } from '~/types/order'

/**
 * Maps the 4 UI tab keys to the backend status strings.
 *   pending    -> pending
 *   processing -> confirmed | processing
 *   shipping   -> shipped | delivered
 *   done       -> finished | cancelled
 */
export type OrderTab = 'pending' | 'processing' | 'shipping' | 'done'

const TAB_STATUSES: Record<OrderTab, string[]> = {
  pending:    ['pending'],
  processing: ['confirmed', 'processing'],
  shipping:   ['shipped', 'delivered'],
  done:       ['finished', 'cancelled']
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
   * Fetch orders for a given tab. The backend `status` query param accepts a
   * single status value, so we make one request per supported status and merge
   * the results when a tab maps to multiple statuses.
   */
  const fetchOrders = async (tab: OrderTab, page = 1) => {
    isLoading.value = true
    error.value     = null
    currentPage.value = page

    const statuses = TAB_STATUSES[tab]

    try {
      if (statuses.length === 1) {
        // Single status — use server-side filtering
        const res = await orderService.listOrders({
          page,
          limit: pageSize.value,
          sort: 'latest:desc',
          status: statuses[0]
        })
        console.log(`[useOrders] Tab: ${tab}, Status: ${statuses[0]}, Count: ${res.orders?.length || 0}`, res.orders?.map(o => ({ number: o.number, orderStatus: o.status, paymentStatus: o.payment?.status })))
        
        let fetchedOrders = res.orders ?? []
        // Client-side guard: 'pending' (Waiting Payment) tab only shows pending payments
        if (tab === 'pending') {
          fetchedOrders = fetchedOrders.filter(o => o.payment?.status === 'pending')
        }
        
        orders.value      = fetchedOrders
        totalOrders.value = fetchedOrders.length
      } else {
        // Multiple statuses — fetch each and merge (no server-side multi-filter)
        const results = await Promise.all(
          statuses.map(s =>
            orderService.listOrders({
              page: 1,
              limit: 50, // generous upper bound per status
              sort: 'latest:desc',
              status: s
            })
          )
        )
        let merged = results.flatMap(r => r.orders ?? [])
        console.log(`[useOrders] Tab: ${tab}, Statuses: ${statuses.join(',')}, Count: ${merged.length}`, merged.map(o => ({ number: o.number, orderStatus: o.status, paymentStatus: o.payment?.status })))
        
        // Deduplicate orders by ID to prevent duplicates if the backend returns the same order in multiple queries
        const seen = new Set<string>()
        merged = merged.filter(o => {
          if (seen.has(o.id)) return false
          seen.add(o.id)
          return true
        })

        // Client-side guard: 'processing' (To Ship) tab only shows non-pending payments
        if (tab === 'processing') {
          merged = merged.filter(o => o.payment?.status !== 'pending')
        }
        
        // Sort merged list newest first
        merged.sort((a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        )
        // Manual pagination on client side
        const start = (page - 1) * pageSize.value
        orders.value      = merged.slice(start, start + pageSize.value)
        totalOrders.value = merged.length
      }
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
