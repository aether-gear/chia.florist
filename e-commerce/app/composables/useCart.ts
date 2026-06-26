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
  status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'finished' | 'cancelled' | 'pembayaran' | 'pengemasan' | 'pengiriman' | 'ulasan'
  shipping?: {
    courier: string
    service: string
    trackingNumber: string
    recipientName: string
    phone: string
    address: string
    cost: number
  }
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
  const orders = useState<Order[]>('chia-florist-orders', () => [
    {
      orderId: 'CHIA-982145',
      date: '26 Jun 2026',
      items: [
        {
          id: 'prod-001',
          name: 'Papan Bunga Congratulation Grand Opening Premium',
          price: 850000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '2.0m x 1.5m',
          color: '#1b4332'
        }
      ],
      total: 850000,
      status: 'pending',
      shipping: {
        courier: 'JNE Express',
        service: 'REG',
        trackingNumber: 'AWB-8839214751',
        recipientName: 'Jane Doe',
        phone: '081234567890',
        address: 'Jl. Merdeka No. 45, Kebayoran Baru, Jakarta Selatan, 12110',
        cost: 20000
      }
    },
    {
      orderId: 'CHIA-881234',
      date: '25 Jun 2026',
      items: [
        {
          id: 'prod-002',
          name: 'Standing Flower Congratulation Special',
          price: 1200000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '1.8m',
          color: '#c1121f'
        }
      ],
      total: 1200000,
      status: 'confirmed',
      shipping: {
        courier: 'J&T Express',
        service: 'EZ',
        trackingNumber: 'JT-9948214002',
        recipientName: 'John Smith',
        phone: '081987654321',
        address: 'Sudirman Central Business District (SCBD), Tower B Lt. 12, Jakarta Selatan, 12190',
        cost: 15000
      }
    },
    {
      orderId: 'CHIA-775621',
      date: '25 Jun 2026',
      items: [
        {
          id: 'prod-003',
          name: 'Hand Bouquet Rose & Lily Sweet Pink',
          price: 550000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: 'Standard',
          color: '#ffb3c1'
        }
      ],
      total: 550000,
      status: 'processing',
      shipping: {
        courier: 'SiCepat',
        service: 'REG',
        trackingNumber: 'SI-7738291044',
        recipientName: 'Alice Johnson',
        phone: '087711223344',
        address: 'Kuningan Place, Block C-3, Jakarta Selatan, 12950',
        cost: 18000
      }
    },
    {
      orderId: 'CHIA-664321',
      date: '24 Jun 2026',
      items: [
        {
          id: 'prod-004',
          name: 'Papan Bunga Pernikahan Elegant Double Board',
          price: 950000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '2.5m x 1.5m',
          color: '#1b4332'
        }
      ],
      total: 950000,
      status: 'shipped',
      shipping: {
        courier: 'JNE Express',
        service: 'YES',
        trackingNumber: 'AWB-8839214755',
        recipientName: 'Michael Brown',
        phone: '081299887766',
        address: 'Jl. Kemang Raya No. 10, Kemang, Jakarta Selatan, 12730',
        cost: 30000
      }
    },
    {
      orderId: 'CHIA-553199',
      date: '24 Jun 2026',
      items: [
        {
          id: 'prod-005',
          name: 'Papan Bunga Wisuda Congratulations Modern',
          price: 650000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '1.8m',
          color: '#0077b6'
        }
      ],
      total: 650000,
      status: 'delivered',
      shipping: {
        courier: 'Grab Express',
        service: 'Instant',
        trackingNumber: 'GRAB-229981442',
        recipientName: 'David Lee',
        phone: '081344556677',
        address: 'Apartemen Sudirman Hill, Unit 15-A, Jakarta Pusat, 10210',
        cost: 25000
      }
    },
    {
      orderId: 'CHIA-442188',
      date: '20 Jun 2026',
      items: [
        {
          id: 'prod-001',
          name: 'Papan Bunga Congratulation Grand Opening Premium',
          price: 850000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '2.0m x 1.5m',
          color: '#1b4332'
        }
      ],
      total: 850000,
      status: 'finished',
      shipping: {
        courier: 'JNE Express',
        service: 'REG',
        trackingNumber: 'AWB-8839214760',
        recipientName: 'Sarah Connor',
        phone: '085522334455',
        address: 'Menteng Residence, Block D No. 5, Jakarta Pusat, 10310',
        cost: 20000
      }
    },
    {
      orderId: 'CHIA-331077',
      date: '18 Jun 2026',
      items: [
        {
          id: 'prod-002',
          name: 'Standing Flower Congratulation Special',
          price: 1200000,
          image: '/images/custom-preview.png',
          quantity: 1,
          size: '1.8m',
          color: '#c1121f'
        }
      ],
      total: 1200000,
      status: 'cancelled',
      shipping: {
        courier: 'JNE Express',
        service: 'REG',
        trackingNumber: 'AWB-8839214770',
        recipientName: 'Sarah Connor',
        phone: '085522334455',
        address: 'Menteng Residence, Block D No. 5, Jakarta Pusat, 10310',
        cost: 20000
      }
    }
  ])
  const isLoggedIn = useCookie('is_logged_in')

  const loadCart = (force = false): Promise<void> => {
    if (isLoggedIn.value !== 'true') return Promise.resolve()

    if (!force && import.meta.client) {
      const cached = localStorage.getItem('chia-florist-cart-cache')
      if (cached) {
        try {
          const parsed = JSON.parse(cached)
          const customItems = cart.value.filter(i => i.isCustom)
          cart.value = [...parsed, ...customItems]
          return Promise.resolve()
        } catch (e) {
          console.error('Failed to parse cart cache:', e)
        }
      }
    }

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
              image: item.images?.thumbnail || '',
              quantity: Number(item.quantity),
              shopId: item.shop_id,
              isCustom: false,
              size: size, 
              color: color
            }
          })

          const customItems = cart.value.filter(i => i.isCustom)
          cart.value = [...backendItems, ...customItems]

          if (import.meta.client) {
            localStorage.setItem('chia-florist-cart-cache', JSON.stringify(backendItems))
          }
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
    watch(isLoggedIn, (newVal, oldVal) => {
      if (newVal === 'true') {
        const shouldForce = oldVal !== undefined && oldVal !== 'true'
        loadCart(shouldForce)
      } else {
        cart.value = cart.value.filter(i => i.isCustom)
        localStorage.removeItem('chia-florist-cart-cache')
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
        
        await loadCart(true)
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
        await loadCart(true)
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
            await loadCart(true)
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
      await loadCart(true)
    } else {
      for (const item of orderItems) {
        if (import.meta.client) {
          localStorage.removeItem(`cart_attr_${item.id}`)
        }
      }
    }

    const newOrder: Order = {
      orderId: 'CHIA-' + Date.now().toString().slice(-6),
      date: new Date().toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' }),
      items: orderItems,
      total: totalAmount,
      status: 'pending',
      shipping: {
        courier: 'JNE Express',
        service: 'REG',
        trackingNumber: 'AWB-' + Math.floor(1000000000 + Math.random() * 9000000000),
        recipientName: 'Jane Doe',
        phone: '081234567890',
        address: 'Jl. Merdeka No. 45, Kebayoran Baru, Jakarta Selatan, 12110',
        cost: 20000
      }
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