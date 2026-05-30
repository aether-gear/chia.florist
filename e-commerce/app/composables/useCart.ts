// app/composables/useCart.ts
import { computed } from 'vue'

export interface CartItem {
  id: string
  name: string
  price: number
  image: string
  quantity: number
  size?: string
  color?: string
  isCustom?: boolean
}

export interface Order {
  orderId: string
  date: string
  items: CartItem[]
  total: number
  status: 'pembayaran' | 'pengemasan' | 'pengiriman' | 'ulasan'
}

export const useCart = () => {
  // Global reactive state Nuxt
  const cart = useState<CartItem[]>('chia-florist-cart', () => [])
  const orders = useState<Order[]>('chia-florist-orders', () => [])

  // Fungsi Tambah Item ke Keranjang
  const addToCart = (item: Omit<CartItem, 'quantity'>, qty = 1) => {
    const existingItem = cart.value.find(
      i => i.id === item.id && i.size === item.size && i.color === item.color
    )

    if (existingItem) {
      existingItem.quantity += qty
    } else {
      cart.value.push({ ...item, quantity: qty })
    }
  }

  // Fungsi Hapus Item
  const removeFromCart = (id: string, size?: string, color?: string) => {
    const index = cart.value.findIndex(i => i.id === id && i.size === size && i.color === color)
    if (index !== -1) {
      cart.value.splice(index, 1)
    }
  }

  // Perbarui Jumlah Kuantitas (Plus/Minus)
  const updateQuantity = (id: string, size: string | undefined, color: string | undefined, change: number) => {
    const item = cart.value.find(i => i.id === id && i.size === size && i.color === color)
    if (item) {
      item.quantity += change
      if (item.quantity <= 0) {
        removeFromCart(id, size, color)
      }
    }
  }

  // Menghitung Subtotal Belanjaan
  const cartSubtotal = computed(() => {
    return cart.value.reduce((total, item) => total + item.price * item.quantity, 0)
  })

  // Menghitung Total Item di Keranjang
  const cartCount = computed(() => {
    return cart.value.reduce((total, item) => total + item.quantity, 0)
  })

  // Fungsi memindahkan items ke order tracking
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

  // FIX: Kembalikan semua state dan fungsi dalam SATU objek tunggal
  return {
    cart,
    orders,
    addToCart,
    removeFromCart,
    updateQuantity,
    cartSubtotal,
    cartCount,
    checkoutToOrder
  }
}