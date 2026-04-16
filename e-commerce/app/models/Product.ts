export interface Product {
  id: string
  name: string
  description: string
  price: number
  imageUrl: string
  category: string
}

export interface CartItem {
  product: Product
  quantity: number
}
