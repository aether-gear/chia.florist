// app/models/Product.ts
import type { Product } from '~/types/product'

export type { Product }

export interface CartItem {
  product: Product
  quantity: number
}
