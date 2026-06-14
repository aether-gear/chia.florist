// app/types/cart.ts

export interface BackendCartItem {
  product_id: string
  shop_id: string
  name: string
  price: number
  subtotal: number
  quantity: number
  images: {
    thumbnail: string
  }
}

export interface BackendCartResponse {
  cart_id: string
  items: BackendCartItem[]
  total: number
}
