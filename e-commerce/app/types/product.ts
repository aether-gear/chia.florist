// app/types/product.ts

export interface ProductImage {
  thumbnail: string
  preview?: string
  detail?: string
}

/**
 * Raw API Response structure for Product Detail (GET /products/{slug})
 */
export interface ApiProductDetail {
  id: string
  sku: string
  name: string
  slug: string
  status: string
  is_available: boolean
  price: number
  stock: number
  description?: string
  weight?: number
  updated_at?: string | null
  banner: {
    thumbnail: string | null
    preview: string | null
    detail: string | null
  }
  gallery: {
    thumbnail: string | null
    preview: string | null
    detail: string | null
  }[]
  availability: {
    slug: string
    name: string
    stock: number
  }[]
}

/**
 * Raw API Response structure for Product List Item (GET /products)
 */
export interface ApiProductListItem {
  id: string
  sku: string
  name: string
  slug: string
  status: string
  is_available: boolean
  price: number
  stock: number
  banner: {
    thumbnail: string | null
    preview: string | null
    detail: string | null
  }
  availability: {
    slug: string
    name: string
    stock: number
  }[]
}

/**
 * Raw API Response structure for Product Search List (GET /products)
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
  shopId?: string
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
  slug?: string
}
