// app/composables/useCart.ts
import { computed, watch } from 'vue'
import { cartService } from '~/services/cartService'
import { formatRupiah } from '~/utils/formatter' // Import Formatter Rupiah Global

/* ─── Standardized Custom Design Payload Schema (v1.0.0) ───────────────── */

export interface TypographySpec {
  text: string | null
  fontId: 'inter' | 'playfair' | 'dancing' | 'bebas' | 'merriweather' | 'pacifico'
  fontSizePx: number
  fontColorHex: string // 6-digit hex format: #RRGGBB
  alignment: 'left' | 'center' | 'right'
}

export interface BoardSectionSpec {
  bgColorHex: string
  cornerStyle: 'none' | 'rounded' | 'cut' | 'ornate' | 'floral'
  header: TypographySpec
  body: TypographySpec
}

export interface BorderSpec {
  style: 'none' | 'solid' | 'double' | 'dashed' | 'dotted' | 'groove' | 'ridge' | 'ornate'
  colorHex: string
  widthPx: number
  showCenterDivider: boolean
}

export interface CrestSpec {
  visible: boolean
  variantId: 'classic' | 'modern' | 'grand'
  primaryColorHex: string
  secondaryColorHex: string
  scalePercent: number
}

export interface BaseElement {
  id: string
  type: 'image' | 'brush'
  transform: {
    xPercent: number
    yPercent: number
    scalePercent: number
    rotationDeg: number
  }
}

export interface ImageElement extends BaseElement {
  type: 'image'
  src: string
  frameStyle: 'none' | 'square' | 'circle'
  crop: { xPercent: number; yPercent: number; zoom: number }
}

export interface BrushElement extends BaseElement {
  type: 'brush'
  brushType: 'flower' | 'rose'
  colorHex: string
}

export type DesignElement = ImageElement | BrushElement

export interface CustomDesignPayloadV1 {
  metadata: {
    version: '1.0.0'
    editorVersion: string
    platform: string
    locale: string
    createdAt: string
    updatedAt: string
    checksum: string
  }
  layout: {
    physicalSizeId: 'small' | 'medium' | 'large'
    upperHeightRatio: number
    border: BorderSpec
  }
  sections: {
    upper: BoardSectionSpec
    lower: BoardSectionSpec
  }
  decorations: {
    topCrest: CrestSpec
    bottomCrest: CrestSpec
  }
  elements: DesignElement[]
  assets: {
    previewBase64: string | null
    previewAssetId: string | null
    previewUrl: string | null
    bucketPath: string | null
    storageProvider: 'supabase' | 's3' | 'local' | null
  }
}

export type CustomDesignPayload = CustomDesignPayloadV1

// Calculate SHA-256 / FNV-1a checksum for design deduplication
export const calculateDesignChecksum = (payload: any): string => {
  const copy = { ...payload }
  if (copy.metadata) copy.metadata = { ...copy.metadata, checksum: '' }
  const jsonStr = JSON.stringify(copy)
  let hash = 0x811c9dc5
  for (let i = 0; i < jsonStr.length; i++) {
    hash ^= jsonStr.charCodeAt(i)
    hash += (hash << 1) + (hash << 4) + (hash << 7) + (hash << 8) + (hash << 24)
  }
  return (hash >>> 0).toString(16).padStart(8, '0')
}

// Validate 6-digit hex color format (#RRGGBB)
export const normalizeHexColor = (color: string | undefined, fallback = '#FFFFFF'): string => {
  if (!color) return fallback
  const trimmed = color.trim().toUpperCase()
  if (/^#[0-9A-F]{6}$/.test(trimmed)) return trimmed
  if (/^#[0-9A-F]{3}$/.test(trimmed)) {
    return `#${trimmed[1]}${trimmed[1]}${trimmed[2]}${trimmed[2]}${trimmed[3]}${trimmed[3]}`
  }
  return fallback
}

// Migration utility to seamlessly upgrade legacy v1.0 flat payloads to v1.0.0
export const migrateCustomDesignPayload = (raw: any): CustomDesignPayloadV1 => {
  if (!raw) {
    throw new Error('Empty custom design payload')
  }

  // Already v1.0.0 structured payload
  if (raw.metadata && raw.metadata.version === '1.0.0' && raw.layout && raw.sections) {
    const payload = raw as CustomDesignPayloadV1
    if (!payload.metadata.checksum) {
      payload.metadata.checksum = calculateDesignChecksum(payload)
    }
    return payload
  }

  // Legacy v1.0 flat format migration
  const upper = raw.upper || {}
  const lower = raw.lower || {}
  const border = raw.border || {}
  const topCrest = raw.topCrest || {}
  const bottomCrest = raw.bottomCrest || {}

  const payload: CustomDesignPayloadV1 = {
    metadata: {
      version: '1.0.0',
      editorVersion: '1.0.0',
      platform: 'web',
      locale: 'id-ID',
      createdAt: raw.generatedAt || new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      checksum: ''
    },
    layout: {
      physicalSizeId: (raw.physicalSizeId || 'medium') as any,
      upperHeightRatio: Math.max(0.25, Math.min(Number(raw.heightRatio || 0.58), 0.75)),
      border: {
        style: (border.style || 'solid') as any,
        colorHex: normalizeHexColor(border.color, '#F5C842'),
        widthPx: Number(border.width || 12),
        showCenterDivider: Boolean(border.center ?? true)
      }
    },
    sections: {
      upper: {
        bgColorHex: normalizeHexColor(upper.bgColor, '#C0392B'),
        cornerStyle: (upper.cornerStyle || 'none') as any,
        header: {
          text: upper.headerText || null,
          fontId: (upper.headerFont || 'playfair') as any,
          fontSizePx: Number(upper.headerFontSize || 36),
          fontColorHex: normalizeHexColor(upper.headerColor, '#FFD700'),
          alignment: (upper.headerAlign || 'center') as any
        },
        body: {
          text: upper.bodyText || null,
          fontId: (upper.bodyFont || 'inter') as any,
          fontSizePx: Number(upper.bodyFontSize || 20),
          fontColorHex: normalizeHexColor(upper.bodyColor, '#FFFFFF'),
          alignment: (upper.bodyAlign || 'center') as any
        }
      },
      lower: {
        bgColorHex: normalizeHexColor(lower.bgColor, '#1A3A5C'),
        cornerStyle: (lower.cornerStyle || 'none') as any,
        header: {
          text: lower.headerText || null,
          fontId: (lower.headerFont || 'bebas') as any,
          fontSizePx: Number(lower.headerFontSize || 26),
          fontColorHex: normalizeHexColor(lower.headerColor, '#FFFFFF'),
          alignment: (lower.headerAlign || 'center') as any
        },
        body: {
          text: lower.bodyText || null,
          fontId: (lower.bodyFont || 'inter') as any,
          fontSizePx: Number(lower.bodyFontSize || 22),
          fontColorHex: normalizeHexColor(lower.bodyColor, '#FFFFFF'),
          alignment: (lower.bodyAlign || 'center') as any
        }
      }
    },
    decorations: {
      topCrest: {
        visible: Boolean(topCrest.enabled),
        variantId: (topCrest.style || 'classic') as any,
        primaryColorHex: normalizeHexColor(topCrest.primary, '#E63946'),
        secondaryColorHex: normalizeHexColor(topCrest.secondary, '#F1FAEE'),
        scalePercent: Number(topCrest.size || 40)
      },
      bottomCrest: {
        visible: Boolean(bottomCrest.enabled),
        variantId: (bottomCrest.style || 'classic') as any,
        primaryColorHex: normalizeHexColor(bottomCrest.primary, '#E63946'),
        secondaryColorHex: normalizeHexColor(bottomCrest.secondary, '#F1FAEE'),
        scalePercent: Number(bottomCrest.size || 40)
      }
    },
    elements: Array.isArray(raw.elements) ? raw.elements.map((el: any) => {
      if (el.type === 'image') {
        return {
          id: el.id || `elem-${Math.random()}`,
          type: 'image' as const,
          src: el.src || '',
          frameStyle: (el.frame || 'square') as any,
          crop: { xPercent: el.cropX ?? 50, yPercent: el.cropY ?? 50, zoom: el.zoom ?? 1 },
          transform: {
            xPercent: el.x ?? 15,
            yPercent: el.y ?? 15,
            scalePercent: el.width ?? 22,
            rotationDeg: el.rotation ?? 0
          }
        }
      }
      return {
        id: el.id || `elem-${Math.random()}`,
        type: 'brush' as const,
        brushType: (el.brushType || 'flower') as any,
        colorHex: normalizeHexColor(el.color, '#E85D75'),
        transform: {
          xPercent: el.x ?? 50,
          yPercent: el.y ?? 50,
          scalePercent: el.size ?? 48,
          rotationDeg: el.rotation ?? 0
        }
      }
    }) : [],
    assets: {
      previewBase64: raw.previewBase64 || raw.assets?.previewBase64 || null,
      previewAssetId: raw.previewAssetId || raw.assets?.previewAssetId || null,
      previewUrl: raw.previewUrl || raw.assets?.previewUrl || null,
      bucketPath: raw.bucketPath || raw.assets?.bucketPath || null,
      storageProvider: raw.storageProvider || raw.assets?.storageProvider || 'supabase'
    }
  }

  payload.metadata.checksum = calculateDesignChecksum(payload)
  return payload
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
  shopId?: string
  customDesign?: CustomDesignPayloadV1  // only present for custom board orders
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

  const loadCart = (force = false): Promise<void> => {
    if (isLoggedIn.value !== 'true') {
      isLoadingCart.value = false
      return Promise.resolve()
    }

    if (loadCartPromise) return loadCartPromise

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

            if (item.product_variant_type === 'custom' || item.item_type === 'custom' || item.custom_design) {
              const migratedDesign = item.custom_design ? migrateCustomDesignPayload(item.custom_design) : undefined
              const cartItemId = item.cart_item_id || item.id || `custom-${Date.now()}`
              return {
                id: cartItemId,
                cartItemId: cartItemId,
                name: item.product_name || item.name || 'Custom Board',
                price: price,
                subtotal: subtotal,
                image: migratedDesign?.assets?.previewUrl || migratedDesign?.assets?.previewBase64 || item.images?.thumbnail || '/images/custom-preview.png',
                quantity: Number(item.quantity),
                shopId: item.shop_id,
                isCustom: true,
                itemType: 'custom',
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
              id: item.product_id,
              name: item.name,
              price: price,
              image: item.images?.thumbnail || '',
              quantity: Number(item.quantity),
              shopId: item.shop_id,
              isCustom: false,
              itemType: 'standard',
              size: size, 
              color: color
            }
          })

          if (isLoggedIn.value === 'true') {
            cart.value = backendItems
          } else {
            const customItems = cart.value.filter(localItem => 
              localItem.isCustom && 
              !backendItems.some(b => b.id === localItem.id || (b.isCustom && b.name === localItem.name && b.size === localItem.size))
            )
            cart.value = [...backendItems, ...customItems]
          }

          if (import.meta.client) {
            localStorage.setItem('chia-florist-cart-cache', JSON.stringify(backendItems))
          }
        }
      } catch (err) {
        console.error('Failed to load cart from backend:', err)
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
        loadCart(true)
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
          await cartService.removeCustomItem(id)
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
        if (!item.isCustom) {
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

  return {
    cart,
    orders,
    isLoadingCart,
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