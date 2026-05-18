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

export const useCart = () => {
  // Global reactive state yang tersimpan di memori Nuxt
  const cart = useState<CartItem[]>('chia-florist-cart', () => [])

  // Fungsi Tambah Item ke Keranjang
  const addToCart = (item: Omit<CartItem, 'quantity'>, qty = 1) => {
    // Cek apakah produk dengan size & warna yang sama sudah ada di keranjang
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

  // Mengitung Subtotal Belanjaan
  const cartSubtotal = computed(() => {
    return cart.value.reduce((total, item) => total + item.price * item.quantity, 0)
  })

  // Menghitung Total Item di Keranjang (untuk badge di navbar)
  const cartCount = computed(() => {
    return cart.value.reduce((total, item) => total + item.quantity, 0)
  })

  return {
    cart,
    addToCart,
    removeFromCart,
    updateQuantity,
    cartSubtotal,
    cartCount
  }
}