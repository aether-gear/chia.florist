// app/composables/useCart.ts
import { computed, watch } from 'vue'
import { cartService } from '~/services/cartService'
import { formatRupiah } from '~/utils/formatter' // Import Formatter Rupiah Global

export interface CartItem {
  id: string
  name: string
  price: number
  image: string
  quantity: number
  size?: string
  color?: string
  isCustom?: boolean
  shopId?: string
}

export interface Order {
  orderId: string
  date: string
  items: CartItem[]
  total: number
  status: 'pembayaran' | 'pengemasan' | 'pengiriman' | 'ulasan'
}

interface PendingEntry {
  quantity: number
  shopId: string
  timeoutId: ReturnType<typeof setTimeout> | null
}

const pendingAdditions: Record<string, PendingEntry> = {}
const pendingUpdates: Record<string, PendingEntry> = {}
let cartWatcherInitialized = false
let loadCartPromise: Promise<void> | null = null

export const useCart = () => {
  const cart = useState<CartItem[]>('chia-florist-cart', () => [])
  const orders = useState<Order[]>('chia-florist-orders', () => [])
  const isLoggedIn = useCookie('is_logged_in')

  const loadCart = (): Promise<void> => {
    if (isLoggedIn.value !== 'true') return Promise.resolve()
    if (loadCartPromise) return loadCartPromise

    loadCartPromise = (async () => {
      try {
        const response = await cartService.getCart()
        if (response && Array.isArray(response.items)) {
          const backendItems: CartItem[] = response.items.map((item: any) => {
            let size = '1.8m'
            let color = '#1b4332'
            let price = Number(item.price)

            if (import.meta.client) {
              const savedAttr = localStorage.getItem(`cart_attr_${item.product_id}`)
              if (savedAttr) {
                try {
                  const parsed = JSON.parse(savedAttr)
                  if (parsed.size) size = parsed.size
                  if (parsed.color) color = parsed.color
                  if (parsed.price) price = Number(parsed.price)
                } catch (e) {
                  console.error('Failed to parse saved cart attributes:', e)
                }
              }
            }

            return {
              id: item.product_id,
              name: item.name,
              price: price,
              image: item.images?.thumbnail || '/images/birthday.jpeg',
              quantity: Number(item.quantity),
              shopId: item.shop_id,
              isCustom: false,
              size: size, 
              color: color
            }
          })

          const customItems = cart.value.filter(i => i.isCustom)
          cart.value = [...backendItems, ...customItems]
        }
      } catch (err) {
        console.error('Failed to load cart from backend:', err)
      } finally {
        loadCartPromise = null
      }
    })()

    return loadCartPromise
  }

  if (import.meta.client && !cartWatcherInitialized) {
    cartWatcherInitialized = true
    watch(isLoggedIn, (newVal) => {
      if (newVal === 'true') {
        loadCart()
      } else {
        cart.value = cart.value.filter(i => i.isCustom)
      }
    }, { immediate: true })
  }

  const flushCart = async () => {
    if (isLoggedIn.value !== 'true') return
    const updatePromises = Object.keys(pendingUpdates).map(async (productId) => {
      const pending = pendingUpdates[productId]
      if (!pending) return
      clearTimeout(pending.timeoutId ?? undefined)
      const qty = pending.quantity
      const shopId = pending.shopId
      delete pendingUpdates[productId]
      try {
        await cartService.updateItem(shopId, productId, qty)
      } catch (err) {
        console.error(err)
      }
    })
    await Promise.all(updatePromises)
    await loadCart()
  }

  const addToCart = async (item: Omit<CartItem, 'quantity'>, qty = 1) => {
    if (import.meta.client && !item.isCustom) {
      localStorage.setItem(`cart_attr_${item.id}`, JSON.stringify({
        size: item.size,
        color: item.color,
        price: item.price
      }))
    }

    if (item.isCustom) {
      const existingItem = cart.value.find(
        i => i.id === item.id || (i.isCustom && i.name === item.name && i.size === item.size && i.color === item.color)
      )
      if (existingItem) {
        existingItem.quantity += qty
      } else {
        cart.value.push({ ...item, quantity: qty })
      }
      return
    }

    if (isLoggedIn.value === 'true') {
      try {
        const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
        
        // FIX SUCCESS: Gunakan type-casting 'as any' saat melempar payload ke addItem 
        // agar tidak memicu error structural contract pada model parameter backend kalian
        await cartService.addItem({ 
          product_id: item.id, 
          shop_id: shopId, 
          quantity: qty,
          size: item.size || '1.8m',
          color: item.color || '#1b4332'
        } as any)
        
        await loadCart()
      } catch (err) {
        console.error(err)
      }
    } else {
      const existingItem = cart.value.find(i => i.id === item.id)
      if (existingItem) existingItem.quantity += qty
      else cart.value.push({ ...item, quantity: qty })
    }
  }

  const removeFromCart = async (id: string, size?: string, color?: string) => {
    if (import.meta.client) {
      localStorage.removeItem(`cart_attr_${id}`)
    }
    const item = cart.value.find(i => i.id === id)
    if (!item) return

    if (item.isCustom) {
      cart.value = cart.value.filter(i => i.id !== id)
      return
    }

    if (isLoggedIn.value === 'true') {
      try {
        const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
        await cartService.removeItem(shopId, id)
        await loadCart()
      } catch (err) {
        console.error(err)
      }
    } else {
      cart.value = cart.value.filter(i => i.id !== id)
    }
  }

  const updateQuantity = async (id: string, size: string | undefined, color: string | undefined, change: number) => {
    const item = cart.value.find(i => i.id === id)
    if (!item) return

    const newQty = item.quantity + change
    if (newQty <= 0) {
      await removeFromCart(id, size, color)
      return
    }

    item.quantity = newQty

    if (item.isCustom) return

    if (isLoggedIn.value === 'true') {
      const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'

      if (pendingUpdates[id]?.timeoutId) {
        clearTimeout(pendingUpdates[id].timeoutId!)
      }

      pendingUpdates[id] = {
        quantity: newQty,
        shopId,
        timeoutId: setTimeout(async () => {
          try {
            await cartService.updateItem(shopId, id, newQty)
          } catch (err) {
            console.error('Backend sync failed:', err)
          } finally {
            delete pendingUpdates[id]
          }
        }, 400) as any
      }
    }
  }

  const cartSubtotal = computed(() => {
    return cart.value.reduce((total, item) => total + (Number(item.price) * Number(item.quantity)), 0)
  })

  const cartSubtotalFormatted = computed(() => formatRupiah(cartSubtotal.value))
  const cartCount = computed(() => cart.value.reduce((total, item) => total + Number(item.quantity), 0))

  const checkoutToOrder = async (totalAmount: number, itemsToOrder?: CartItem[]) => {
    const orderItems = itemsToOrder || [...cart.value]
    if (orderItems.length === 0) return

    // Clear backend cart for checked out items
    if (isLoggedIn.value === 'true') {
      await flushCart() // Ensure no race condition with pending additions/updates
      for (const item of orderItems) {
        if (import.meta.client) {
          localStorage.removeItem(`cart_attr_${item.id}`)
        }
        if (!item.isCustom) {
          try {
            const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
            await cartService.removeItem(shopId, item.id)
          } catch (err) {
            console.error(`Failed to remove item ${item.id} from backend cart on checkout:`, err)
          }
        }
      }
    } else {
      for (const item of orderItems) {
        if (import.meta.client) {
          localStorage.removeItem(`cart_attr_${item.id}`)
        }
      }
    }

    const newOrder: Order = {
      orderId: 'CHIA-' + Date.now().toString().slice(-6),
      date: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }),
      items: orderItems,
      total: totalAmount,
      status: 'pembayaran'
    }

    orders.value.push(newOrder)

    if (!itemsToOrder) {
      cart.value = []
    } else {
      // Filter out the items that were ordered
      cart.value = cart.value.filter(cartItem => 
        !orderItems.some(ordered => 
          ordered.id === cartItem.id && 
          ordered.size === cartItem.size && 
          ordered.color === cartItem.color
        )
      )
    }
  }

  return {
    cart,
    orders,
    loadCart,
    flushCart,
    addToCart,
    removeFromCart,
    updateQuantity,
    cartSubtotal,
    cartSubtotalFormatted,
    cartCount,
    checkoutToOrder,
    formatRupiah
  }
}