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

export const useCart = () => {
  // Global reactive Nuxt state
  const cart = useState<CartItem[]>('chia-florist-cart', () => [])
  const orders = useState<Order[]>('chia-florist-orders', () => [])
  
  // Track login status using cookie
  const isLoggedIn = useCookie('is_logged_in')

  // Load cart items from the backend and merge with local custom items
  const loadCart = async () => {
    if (isLoggedIn.value !== 'true') {
      return
    }
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
    }
  }

  // Watch for login/logout changes to sync cart
  if (import.meta.client) {
    watch(isLoggedIn, (newVal) => {
      if (newVal === 'true') {
        loadCart()
      } else {
        // Clear regular items on logout, keeping only custom boards
        cart.value = cart.value.filter(i => i.isCustom)
      }
    }, { immediate: true })
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
      try {
        const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
        await cartService.addItem({
          product_id: item.id,
          shop_id: shopId,
          quantity: qty
        })
        await loadCart()
      } catch (err) {
        console.error('Failed to add item to backend cart:', err)
      }
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
    const newQty = item.quantity + change
    if (newQty <= 0) {
      await removeFromCart(id, size, color)
      return
    }

    if (isLoggedIn.value === 'true') {
      try {
        const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
        await cartService.updateItem(shopId, id, newQty)
        await loadCart()
      } catch (err) {
        console.error('Failed to update quantity in backend cart:', err)
      }
    } else {
      item.quantity = newQty
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
    addToCart,
    removeFromCart,
    updateQuantity,
    cartSubtotal,
    cartCount,
    checkoutToOrder
  }
}