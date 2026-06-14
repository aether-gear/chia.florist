// app/composables/useCart.ts
import { computed, watch } from 'vue'
import { cartService } from '~/services/cartService'

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

// Global module-scoped map to track pending debounced product additions
const pendingAdditions: Record<string, PendingEntry> = {}

// Global module-scoped map to track pending debounced product updates
const pendingUpdates: Record<string, PendingEntry> = {}

// Guard: ensure the login watcher is only registered ONCE per client session,
// no matter how many components call useCart().
let cartWatcherInitialized = false

// In-flight loadCart promise — if a fetch is already running, all concurrent
// callers share this same promise instead of firing duplicate requests.
let loadCartPromise: Promise<void> | null = null

export const useCart = () => {
  // Global reactive Nuxt state — useState is a singleton keyed by name
  const cart = useState<CartItem[]>('chia-florist-cart', () => [])
  const orders = useState<Order[]>('chia-florist-orders', () => [])

  // Track login status using cookie
  const isLoggedIn = useCookie('is_logged_in')

  // Load cart items from the backend and merge with local custom items.
  // Deduplicates concurrent calls: if a fetch is already in-flight, reuse it.
  const loadCart = (): Promise<void> => {
    if (isLoggedIn.value !== 'true') return Promise.resolve()

    // Return the existing in-flight promise to avoid a duplicate request
    if (loadCartPromise) return loadCartPromise

    loadCartPromise = (async () => {
      try {
        const response = await cartService.getCart()
        if (response && Array.isArray(response.items)) {
          const backendItems: CartItem[] = response.items.map(item => ({
            id: item.product_id,
            name: item.name,
            price: item.price,
            image: item.images?.thumbnail || '/images/birthday.jpeg',
            quantity: item.quantity,
            shopId: item.shop_id,
            isCustom: false
          }))

          // Preserve any custom items from the current cart
          const customItems = cart.value.filter(i => i.isCustom)
          cart.value = [...backendItems, ...customItems]
        }
      } catch (err) {
        console.error('Failed to load cart from backend:', err)
      } finally {
        // Always clear the promise so future calls can issue a new request
        loadCartPromise = null
      }
    })()

    return loadCartPromise
  }

  // Watch for login/logout changes to sync cart.
  // The `cartWatcherInitialized` flag ensures this block runs only once globally,
  // regardless of how many components call useCart().
  if (import.meta.client && !cartWatcherInitialized) {
    cartWatcherInitialized = true

    watch(isLoggedIn, (newVal) => {
      if (newVal === 'true') {
        loadCart()
      } else {
        // Clear regular items on logout, keeping only custom boards
        cart.value = cart.value.filter(i => i.isCustom)

        // Cancel and clear any pending additions
        for (const id of Object.keys(pendingAdditions)) {
          const pending = pendingAdditions[id]
          if (pending) {
            clearTimeout(pending.timeoutId ?? undefined)
            delete pendingAdditions[id]
          }
        }

        // Cancel and clear any pending updates
        for (const id of Object.keys(pendingUpdates)) {
          const pending = pendingUpdates[id]
          if (pending) {
            clearTimeout(pending.timeoutId ?? undefined)
            delete pendingUpdates[id]
          }
        }
      }
    }, { immediate: true })
  }

  // Flush any pending debounced additions and updates to the backend immediately.
  // Called before remove/checkout to avoid race conditions.
  const flushCart = async () => {
    if (isLoggedIn.value !== 'true') return

    const additionPromises = Object.keys(pendingAdditions).map(async (productId) => {
      const pending = pendingAdditions[productId]
      if (!pending) return

      clearTimeout(pending.timeoutId ?? undefined)
      const qty = pending.quantity
      const shopId = pending.shopId
      delete pendingAdditions[productId]

      try {
        await cartService.addItem({
          product_id: productId,
          shop_id: shopId,
          quantity: qty
        })
      } catch (err) {
        console.error(`Failed to flush addition for product ${productId}:`, err)
      }
    })

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
        console.error(`Failed to flush update for product ${productId}:`, err)
      }
    })

    const allPromises = [...additionPromises, ...updatePromises]

    if (allPromises.length > 0) {
      await Promise.all(allPromises)
      await loadCart()
    }
  }

  // Add Item to Cart
  const addToCart = async (item: Omit<CartItem, 'quantity'>, qty = 1) => {
    if (item.isCustom) {
      const existingItem = cart.value.find(
        i => i.id === item.id && i.size === item.size && i.color === item.color
      )
      if (existingItem) {
        existingItem.quantity += qty
      } else {
        cart.value.push({ ...item, quantity: qty })
      }
      return
    }

    // Regular products
    if (isLoggedIn.value === 'true') {
      const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
      const productId = item.id

      // Ensure entry exists in the pending map
      let pending = pendingAdditions[productId]
      if (!pending) {
        pending = { quantity: 0, shopId, timeoutId: null }
        pendingAdditions[productId] = pending
      }

      // Clear previous debounce timer
      if (pending.timeoutId !== null) {
        clearTimeout(pending.timeoutId)
      }
      pending.quantity += qty

      // Optimistically update frontend state for instant UI response
      const existingItem = cart.value.find(i => i.id === productId)
      if (existingItem) {
        existingItem.quantity += qty
      } else {
        cart.value.push({ ...item, quantity: qty })
      }

      // Schedule the debounced API call
      pending.timeoutId = setTimeout(async () => {
        // Re-read from map — if it was already flushed/cancelled, bail out
        const current = pendingAdditions[productId]
        if (!current) return

        const finalQty = current.quantity
        const targetShopId = current.shopId
        delete pendingAdditions[productId]

        try {
          await cartService.addItem({
            product_id: productId,
            shop_id: targetShopId,
            quantity: finalQty
          })
          await loadCart()
        } catch (err) {
          console.error('Failed to add item to backend cart:', err)
        }
      }, 500) // 500ms debounce window
    } else {
      const existingItem = cart.value.find(i => i.id === item.id)
      if (existingItem) {
        existingItem.quantity += qty
      } else {
        cart.value.push({ ...item, quantity: qty })
      }
    }
  }

  // Remove Item from Cart
  const removeFromCart = async (id: string, size?: string, color?: string) => {
    const item = cart.value.find(i => i.id === id && i.size === size && i.color === color)
    if (!item) return

    if (item.isCustom) {
      const index = cart.value.findIndex(i => i.id === id && i.size === size && i.color === color)
      if (index !== -1) {
        cart.value.splice(index, 1)
      }
      return
    }

    // Regular products
    if (isLoggedIn.value === 'true') {
      await flushCart() // Ensure no race condition with pending additions/updates
      try {
        const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
        await cartService.removeItem(shopId, id)
        await loadCart()
      } catch (err) {
        console.error('Failed to remove item from backend cart:', err)
      }
    } else {
      const index = cart.value.findIndex(i => i.id === id)
      if (index !== -1) {
        cart.value.splice(index, 1)
      }
    }
  }

  // Update Item Quantity
  const updateQuantity = async (id: string, size: string | undefined, color: string | undefined, change: number) => {
    const item = cart.value.find(i => i.id === id && i.size === size && i.color === color)
    if (!item) return

    if (item.isCustom) {
      item.quantity += change
      if (item.quantity <= 0) {
        await removeFromCart(id, size, color)
      }
      return
    }

    // Regular products
    // Use the pending target quantity as the base to accumulate rapid clicks correctly
    const existingUpdate = pendingUpdates[id]
    const currentQty = existingUpdate ? existingUpdate.quantity : item.quantity
    const newQty = currentQty + change

    if (newQty <= 0) {
      if (existingUpdate) {
        clearTimeout(existingUpdate.timeoutId ?? undefined)
        delete pendingUpdates[id]
      }
      await removeFromCart(id, size, color)
      return
    }

    // Optimistically update frontend state for instant UI response
    item.quantity = newQty

    if (isLoggedIn.value === 'true') {
      const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'

      // Ensure entry exists in the pending map
      let pending = existingUpdate
      if (!pending) {
        pending = { quantity: newQty, shopId, timeoutId: null }
        pendingUpdates[id] = pending
      }

      // Clear previous debounce timer
      if (pending.timeoutId !== null) {
        clearTimeout(pending.timeoutId)
      }
      pending.quantity = newQty

      // Schedule the debounced API call
      pending.timeoutId = setTimeout(async () => {
        // Re-read from map — if it was already flushed/cancelled, bail out
        const current = pendingUpdates[id]
        if (!current) return

        const finalQty = current.quantity
        const targetShopId = current.shopId
        delete pendingUpdates[id]

        try {
          await cartService.updateItem(targetShopId, id, finalQty)
          await loadCart()
        } catch (err) {
          console.error('Failed to update quantity in backend cart:', err)
        }
      }, 500) // 500ms debounce window
    }
  }

  // Calculate Subtotal
  const cartSubtotal = computed(() => {
    return cart.value.reduce((total, item) => total + item.price * item.quantity, 0)
  })

  // Calculate Total Items Count
  const cartCount = computed(() => {
    return cart.value.reduce((total, item) => total + item.quantity, 0)
  })

  // Checkout function
  const checkoutToOrder = (totalAmount: number) => {
    if (cart.value.length === 0) return

    const newOrder: Order = {
      orderId: 'CHIA-' + Date.now().toString().slice(-6),
      date: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }),
      items: [...cart.value],
      total: totalAmount,
      status: 'pembayaran'
    }

    orders.value.push(newOrder)
    cart.value = []
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
    cartCount,
    checkoutToOrder
  }
}