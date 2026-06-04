// app/types/product.ts

export interface ProductImage {
  thumbnail: string
  preview?: string
  detail?: string
}

/**
 * Raw API Response structure for Product Detail (GET /products/{id})
 */
export interface ApiProductDetail {
  id: string
  sku: string
  name: string
  slug: string
  description?: string
  is_available: boolean
  price: number
  weight?: number
  stock: number
  updated_at?: string | null
  images: ProductImage[]
}

/**
 * Raw API Response structure for Product List Item (GET /products/)
 */
export interface ApiProductListItem {
  id: string
  sku: string
  name: string
  slug: string
  is_available: boolean
  price: number
  stock: number
  images: {
    thumbnail: string
  }
}

/**
 * Raw API Response structure for Product Search List (GET /products/)
 */
export interface ApiProductListResponse {
  limit: number
  page: number
  total: number
  products: ApiProductListItem[]
}

/**
 * Domain Model Product type used throughout the Frontend UI (Details page)
 */
export interface Product {
  id: string
  name: string
  price: number
  rating: number
  reviews: number
  available: boolean
  description: string
  images: string[]
  colors: string[]
  sizes: string[]
  sku?: string
  slug?: string
  weight?: number
  stock?: number
}

/**
 * Domain Model Product type used specifically in Product Cards on Catalog list page
 */
export interface CatalogProduct {
  id: string
  name: string
  price: number
  rating: number
  reviews: number
  image: string
  tag: string
  desc: string
  isCustomRoute?: boolean
  isAvailable: boolean
}
