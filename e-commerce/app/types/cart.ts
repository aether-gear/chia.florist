// app/types/cart.ts

export interface BackendCartItem {
  cart_item_id: string
  product_variant_type: 'standard' | 'custom'
  product_id: string | null
  shop_id: string
  shop_name: string
  shop_slug: string
  name: string
  price: number
  quantity: number
  subtotal: number
  images: {
    thumbnail?: string | null
  }
  custom_design?: Record<string, any>
}

export interface BackendCartResponse {
  cart_id: string
  items: BackendCartItem[]
  total: number
}

