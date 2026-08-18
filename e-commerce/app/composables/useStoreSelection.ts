// app/composables/useStoreSelection.ts
import { ref, computed, watch } from 'vue'
import { bootstrapConfig } from '~/utils/bootstrap'

export interface Shop {
  id: string
  name: string
  slug: string
  description?: string
  is_active?: boolean
}

const STORAGE_KEY = 'chia-selected-shop'
let isInitialized = false

export const useStoreSelection = () => {
  const selectedShop = useState<Shop | null>('chia-selected-shop-state', () => null)
  const activeShops = useState<Shop[]>('chia-active-shops-state', () => [])
  const isLoadingShops = useState<boolean>('chia-loading-shops-state', () => false)

  // Initialize from localStorage on client side
  if (import.meta.client && !isInitialized) {
    isInitialized = true
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        selectedShop.value = JSON.parse(saved)
      }
    } catch (e) {
      console.error('Failed to parse saved shop from localStorage:', e)
    }

    watch(selectedShop, (newShop) => {
      try {
        if (newShop) {
          localStorage.setItem(STORAGE_KEY, JSON.stringify(newShop))
        } else {
          localStorage.removeItem(STORAGE_KEY)
        }
      } catch (e) {
        console.error('Failed to persist shop to localStorage:', e)
      }
    }, { deep: true })
  }

  const fetchActiveShops = async (): Promise<Shop[]> => {
    if (activeShops.value.length > 0) return activeShops.value
    isLoadingShops.value = true
    try {
      const res = await bootstrapConfig.fetchApi<{ shops: Shop[] }>('/shops?active=true')
      if (res && Array.isArray(res.shops)) {
        activeShops.value = res.shops
        return res.shops
      }
    } catch (err) {
      console.error('Failed to fetch active shops:', err)
    } finally {
      isLoadingShops.value = false
    }
    return []
  }

  const selectShop = (shop: Shop | null) => {
    selectedShop.value = shop
  }

  const clearSelectedShop = () => {
    selectedShop.value = null
  }

  return {
    selectedShop: computed(() => selectedShop.value),
    activeShops: computed(() => activeShops.value),
    isLoadingShops: computed(() => isLoadingShops.value),
    fetchActiveShops,
    selectShop,
    clearSelectedShop
  }
}
