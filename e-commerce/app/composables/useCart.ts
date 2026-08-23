// app/composables/useCart.ts
import { computed, watch } from 'vue'
import { cartService } from '~/services/cartService'
import { formatRupiah } from '~/utils/formatter' // Import Formatter Rupiah Global
import { logError } from '~/utils/errorMessages'

import { migrateToV3, calculateDesignChecksum, normalizeHexColor } from '~/features/custom-product/migrate'
import type {
  TypographySpec, BoardSectionSpec, BorderSpec, CrestSpec,
  BaseElement, ImageElement, BrushElement, DesignElement,
  CustomDesignPayloadV1, CustomDesignPayloadV3, CustomDesignPayload
} from '~/features/custom-product/types'

export type {
  TypographySpec, BoardSectionSpec, BorderSpec, CrestSpec,
  BaseElement, ImageElement, BrushElement, DesignElement,
  CustomDesignPayloadV1, CustomDesignPayloadV3, CustomDesignPayload
}
export { calculateDesignChecksum, normalizeHexColor }

// Migration utility to seamlessly upgrade raw payloads (v1.0 flat, v1.0.0, or v3.0.0)
export const migrateCustomDesignPayload = (raw: any): CustomDesignPayloadV3 => {
  return migrateToV3(raw)
}

export interface CartItem {
  id: string
  cartItemId?: string
  name: string
  price: number
  subtotal?: number
  image: string
  quantity: number
  size?: string
  color?: string
  isCustom?: boolean
  itemType?: 'standard' | 'custom'
  productVariantType?: 'standard' | 'custom'
  slug?: string
  shopId?: string
  shopName?: string
  shopSlug?: string
  customDesign?: CustomDesignPayload  // only present for custom board orders
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
  const isLoadingCart = useState<boolean>('chia-florist-cart-loading', () => false)

  const clearCart = () => {
    cart.value = []
    isLoadingCart.value = false
    loadCartPromise = null

    // Clear pending debounce timers
    Object.keys(pendingUpdates).forEach(k => {
      if (pendingUpdates[k]?.timeoutId) {
        clearTimeout(pendingUpdates[k].timeoutId!)
      }
      delete pendingUpdates[k]
    })
    Object.keys(pendingAdditions).forEach(k => {
      if (pendingAdditions[k]?.timeoutId) {
        clearTimeout(pendingAdditions[k].timeoutId!)
      }
      delete pendingAdditions[k]
    })

    if (import.meta.client) {
      try {
        localStorage.removeItem('chia-florist-cart-cache')
        const keysToRemove: string[] = []
        for (let i = 0; i < localStorage.length; i++) {
          const key = localStorage.key(i)
          if (key && (key.startsWith('cart_attr_') || key.startsWith('custom_design_'))) {
            keysToRemove.push(key)
          }
        }
        keysToRemove.forEach(k => localStorage.removeItem(k))
      } catch (e) {
        console.error('Error clearing cart local storage:', e)
      }
    }
  }

  const loadCart = (force = false): Promise<void> => {
    if (isLoggedIn.value !== 'true') {
      clearCart()
      return Promise.resolve()
    }

    if (!force && cart.value.length > 0) {
      return Promise.resolve()
    }

    if (!force && loadCartPromise) return loadCartPromise

    isLoadingCart.value = true
    loadCartPromise = (async () => {
      try {
        const response = await cartService.getCart()
        if (response && Array.isArray(response.items)) {
          const backendItems: CartItem[] = response.items.map((item: any) => {
            let size = '1.8m'
            let color = '#1b4332'
            let price = Number(item.price ?? item.unit_price ?? 0)
            let subtotal = Number(item.subtotal ?? (price * Number(item.quantity || 1)))
            const cartItemId = item.cart_item_id || item.id || `item-${Date.now()}`

            if (item.product_variant_type === 'custom' || item.item_type === 'custom' || item.custom_design) {
              const migratedDesign = item.custom_design ? migrateCustomDesignPayload(item.custom_design) : undefined
              return {
                id: cartItemId,
                cartItemId: cartItemId,
                name: item.name || item.product_name || 'Custom Board',
                price: price,
                subtotal: subtotal,
                image: migratedDesign?.assets?.previewUrl || migratedDesign?.assets?.previewBase64 || item.images?.thumbnail || '/images/custom-preview.png',
                quantity: Number(item.quantity),
                shopId: item.shop_id,
                shopName: item.shop_name,
                shopSlug: item.shop_slug,
                isCustom: true,
                itemType: 'custom',
                productVariantType: 'custom',
                customDesign: migratedDesign,
                size: migratedDesign?.layout?.physicalSizeId || size,
                color: migratedDesign?.sections?.upper?.bgColorHex || color
              }
            }

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
              id: item.product_id || cartItemId,
              cartItemId: cartItemId,
              name: item.name,
              price: price,
              subtotal: subtotal,
              image: item.images?.thumbnail || '',
              quantity: Number(item.quantity),
              slug: item.slug || item.product_slug,
              shopId: item.shop_id,
              shopName: item.shop_name,
              shopSlug: item.shop_slug,
              isCustom: false,
              itemType: 'standard',
              productVariantType: 'standard',
              size: size, 
              color: color
            }
          })

          cart.value = backendItems

          if (import.meta.client) {
            localStorage.setItem('chia-florist-cart-cache', JSON.stringify(backendItems))
          }
        } else {
          cart.value = []
        }
      } catch (err) {
        logError('useCart', err)
      } finally {
        isLoadingCart.value = false
        loadCartPromise = null
      }
    })()

    return loadCartPromise
  }


  if (import.meta.client && !cartWatcherInitialized) {
    cartWatcherInitialized = true
    watch(isLoggedIn, (newVal) => {
      if (newVal === 'true') {
        clearCart()
        loadCart(true)
      } else {
        clearCart()
      }
    }, { immediate: false })
  }

  const flushCart = async () => {
    if (isLoggedIn.value !== 'true') return
    const pendingKeys = Object.keys(pendingUpdates)
    if (pendingKeys.length === 0) return

    const updatePromises = pendingKeys.map(async (productId) => {
      const pending = pendingUpdates[productId]
      if (!pending) return
      clearTimeout(pending.timeoutId ?? undefined)
      const qty = pending.quantity
      const shopId = pending.shopId
      delete pendingUpdates[productId]
      try {
        await cartService.updateItem(shopId, productId, qty)
      } catch (err) {
        logError('useCart', err)
      }
    })
    await Promise.all(updatePromises)
    await loadCart(true)
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
      const designPayload = item.customDesign ? migrateCustomDesignPayload(item.customDesign) : undefined

      if (isLoggedIn.value === 'true') {
        try {
          const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
          await cartService.addItem({
            product_variant_type: 'custom',
            item_type: 'custom',
            shop_id: shopId,
            quantity: qty,
            product_name: item.name,
            physical_size_id: designPayload?.layout?.physicalSizeId || item.size || 'medium',
            unit_price: item.price,
            custom_design: designPayload
          })
          await loadCart(true)
        } catch (err) {
          console.error('Failed to add custom item to backend cart:', err)
        }
      } else {
        const existingItem = cart.value.find(
          i => i.id === item.id || (i.isCustom && i.name === item.name && i.size === item.size && i.color === item.color)
        )
        if (existingItem) {
          existingItem.quantity += qty
          if (designPayload) existingItem.customDesign = designPayload
        } else {
          cart.value.push({ ...item, customDesign: designPayload, quantity: qty, itemType: 'custom' })
        }
      }
      
      if (import.meta.client && designPayload) {
        try {
          localStorage.setItem(`custom_design_${item.id}`, JSON.stringify(designPayload))
          console.info('[Chia Florist] Custom Design Payload v1.0.0 (persisted):', designPayload)
        } catch (e) {
          console.warn('Could not persist custom design to localStorage:', e)
        }
      }
      return
    }

    if (isLoggedIn.value === 'true') {
      try {
        const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
        
        await cartService.addItem({ 
          product_variant_type: 'standard',
          item_type: 'standard',
          product_id: item.id, 
          shop_id: shopId, 
          quantity: qty
        })
        
        await loadCart(true)
      } catch (err) {
        console.error(err)
      }
    } else {
      const existingItem = cart.value.find(i => i.id === item.id)
      if (existingItem) existingItem.quantity += qty
      else cart.value.push({ ...item, quantity: qty, itemType: 'standard' })
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
      if (isLoggedIn.value === 'true') {
        try {
          const cartItemId = item.cartItemId || id
          await cartService.removeCustomItem(cartItemId)
          await loadCart(true)
        } catch (err) {
          console.error('Failed to remove custom item from backend:', err)
        }
      }
      return
    }

    if (isLoggedIn.value === 'true') {
      try {
        const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
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

    const MAX_CART_QTY = 80
    const newQty = Math.min(item.quantity + change, MAX_CART_QTY)
    if (newQty <= 0) {
      await removeFromCart(id, size, color)
      return
    }
    if (newQty === item.quantity) return  // already at cap — nothing to do

    item.quantity = newQty

    if (isLoggedIn.value === 'true') {
      const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'

      if (pendingUpdates[id]?.timeoutId) {
        clearTimeout(pendingUpdates[id].timeoutId!)
      }

      pendingUpdates[id] = {
        quantity: newQty,
        shopId,
        timeoutId: setTimeout(async () => {
          try {
            if (item.isCustom) {
              await cartService.removeCustomItem(id)
              const designPayload = item.customDesign ? migrateCustomDesignPayload(item.customDesign) : undefined
              await cartService.addItem({
                product_variant_type: 'custom',
                item_type: 'custom',
                shop_id: shopId,
                quantity: newQty,
                product_name: item.name,
                physical_size_id: designPayload?.layout?.physicalSizeId || item.size || 'medium',
                unit_price: item.price,
                custom_design: designPayload
              })
            } else {
              await cartService.updateItem(shopId, id, newQty)
            }
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
        if (item.isCustom) {
          try {
            const cartItemId = (item as any).cartItemId || item.id
            await cartService.removeCustomItem(cartItemId)
          } catch (err) {
            console.error(`Failed to remove custom item ${item.id} from backend cart on checkout:`, err)
          }
        } else {
          try {
            const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
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

  const changeCartItemShop = async (cartItemId: string, shopId: string) => {
    const item = cart.value.find(i => (i as any).cartItemId === cartItemId || i.id === cartItemId)
    if (!item) return

    const storeSelection = useStoreSelection()
    const targetShop = storeSelection.activeShops.value.find(s => s.id === shopId)

    if (isLoggedIn.value === 'true') {
      try {
        const idToTransfer = (item as any).cartItemId || item.id
        await cartService.updateCartItemShop(idToTransfer, shopId)
        item.shopId = shopId
        if (targetShop) {
          item.shopName = targetShop.name
          item.shopSlug = targetShop.slug
        }
        await loadCart(true)
      } catch (err) {
        console.error('Failed to update cart item shop:', err)
        throw err
      }
    } else {
      item.shopId = shopId
      if (targetShop) {
        item.shopName = targetShop.name
        item.shopSlug = targetShop.slug
      }
    }
  }

  return {
    cart,
    orders,
    isLoadingCart,
    loadCart,
    clearCart,
    flushCart,
    addToCart,
    removeFromCart,
    updateQuantity,
    changeCartItemShop,
    cartSubtotal,
    cartSubtotalFormatted,
    cartCount,
    checkoutToOrder,
    formatRupiah
  }
}