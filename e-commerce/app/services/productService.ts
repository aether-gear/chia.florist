// app/services/productService.ts
import type { Product, ApiProductDetail, ApiProductListItem, ApiProductListResponse, CatalogProduct } from '~/types/product'
import { bootstrapConfig } from '~/utils/bootstrap'

export const productService = {
  /**
   * Helper function to map API product detail to UI Product domain model.
   */
  mapApiProduct(apiProduct: ApiProductDetail): Product {
    // Map list of image URLs from ProductImage array
    const images = Array.isArray(apiProduct.images)
      ? apiProduct.images.map(img => img.preview || img.detail || img.thumbnail)
      : []

    // Fallback image if empty
    if (images.length === 0) {
      images.push('/images/birthday.jpeg')
    }

    return {
      id: apiProduct.id,
      name: apiProduct.name,
      price: apiProduct.price,
      rating: 4.5,
      reviews: 150,
      available: apiProduct.is_available,
      description: apiProduct.description || 'Premium arrangement from Chia Florist.',
      images: images,
      colors: ['#cbd5e1', '#f43f5e'],
      sizes: ['1.5m', '1.8m', '2m'],
      sku: apiProduct.sku,
      slug: apiProduct.slug,
      weight: apiProduct.weight,
      stock: apiProduct.stock,
      shopId: (apiProduct as any).shop_id
    }
  },

  /**
   * Helper function to map API product list item to CatalogProduct domain model.
   */
  mapApiCatalogProduct(apiProduct: ApiProductListItem): CatalogProduct {
    return {
      id: apiProduct.id,
      name: apiProduct.name,
      price: apiProduct.price,
      rating: 4.8,
      reviews: 180,
      image: apiProduct.images?.thumbnail || '/images/birthday.jpeg',
      tag: apiProduct.sku ? apiProduct.sku.split('-')[1] || 'Collection' : 'Collection',
      desc: `Premium quality ${apiProduct.name} crafted for your special moments.`,
      isAvailable: apiProduct.is_available
    }
  },

  /**
   * Fetch catalog products directly from the API.
   * No mockup fallbacks.
   */
  async getCatalogProducts(params?: { name?: string; page?: number; limit?: number }): Promise<CatalogProduct[]> {
    const query: Record<string, any> = {}
    if (params?.name) query.name = params.name
    if (params?.page) query.page = params.page
    if (params?.limit) query.limit = params.limit

    const response = await bootstrapConfig.fetchApi<ApiProductListResponse>('/products/', { query })
    if (response && Array.isArray(response.products)) {
      return response.products
        .filter(p => p.is_available)
        .map(p => this.mapApiCatalogProduct(p))
    }
    return []
  },

  /**
   * Fetch specific product details by ID directly from the API.
   * No mockup fallbacks.
   */
  async getProductById(id: string): Promise<Product | null> {
    const response = await $fetch<ApiProductDetail & { shop_id?: string }>(`/api/products/${id}`)

    if (response && response.id) {
      return this.mapApiProduct(response)
    }
    return null
  }
}
