// app/services/productService.ts
import type { Product, ApiProductDetail, ApiProductListItem, ApiProductListResponse, CatalogProduct, FindProductsParams, PaginatedCatalogProducts } from '~/types/product'
import { bootstrapConfig } from '~/utils/bootstrap'

export const productService = {
  /**
   * Helper function to map API product detail to UI Product domain model.
   */
  mapApiProduct(apiProduct: ApiProductDetail): Product {
    // Map list of image URLs from banner and gallery
    const images: string[] = []
    if (apiProduct.banner?.preview || apiProduct.banner?.detail || apiProduct.banner?.thumbnail) {
      images.push((apiProduct.banner.preview || apiProduct.banner.detail || apiProduct.banner.thumbnail)!)
    }
    if (Array.isArray(apiProduct.gallery)) {
      apiProduct.gallery.forEach(img => {
        const url = img.preview || img.detail || img.thumbnail
        if (url && !images.includes(url)) {
          images.push(url)
        }
      })
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
      availability: apiProduct.availability,
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
      image: apiProduct.banner?.thumbnail || '',
      tag: apiProduct.sku ? apiProduct.sku.split('-')[1] || 'Collection' : 'Collection',
      desc: `Premium quality ${apiProduct.name} crafted for your special moments.`,
      isAvailable: apiProduct.is_available,
      slug: apiProduct.slug,
      sku: apiProduct.sku,
      status: apiProduct.status,
      stock: apiProduct.stock,
      availability: apiProduct.availability
    }
  },

  /**
   * Fetch paginated catalog products directly from the API (GET /products).
   * Supports shop filtering (query params & X-Shop-ID header).
   */
  async getPaginatedCatalogProducts(params?: FindProductsParams, shopIdHeader?: string): Promise<PaginatedCatalogProducts> {
    const query: Record<string, any> = {}
    if (params?.name) query.name = params.name
    if (params?.id) query.id = params.id
    if (params?.shop_id) query.shop_id = params.shop_id
    if (params?.shop_slug) query.shop_slug = params.shop_slug
    if (params?.page) query.page = params.page
    if (params?.limit) query.limit = params.limit
    if (params?.sort) query.sort = params.sort

    const headers: Record<string, string> = {}
    if (shopIdHeader) {
      headers['X-Shop-ID'] = shopIdHeader
    }

    const response = await bootstrapConfig.fetchApi<ApiProductListResponse>('/products', {
      query,
      headers: Object.keys(headers).length > 0 ? headers : undefined
    })

    if (response && Array.isArray(response.products)) {
      return {
        products: response.products.map(p => this.mapApiCatalogProduct(p)),
        limit: response.limit ?? 10,
        page: response.page ?? 1,
        total: response.total ?? response.products.length
      }
    }

    return {
      products: [],
      limit: params?.limit ?? 10,
      page: params?.page ?? 1,
      total: 0
    }
  },

  /**
   * Fetch catalog products directly from the API.
   * Shortcut returning array of CatalogProduct.
   */
  async getCatalogProducts(params?: FindProductsParams, shopIdHeader?: string): Promise<CatalogProduct[]> {
    const res = await this.getPaginatedCatalogProducts(params, shopIdHeader)
    return res.products
  },

  /**
   * Fetch specific product details by ID or slug directly from the API.
   * No mockup fallbacks.
   */
  async getProductById(idOrSlug: string, shopIdHeader?: string): Promise<Product | null> {
    const headers: Record<string, string> = {}
    if (shopIdHeader) {
      headers['X-Shop-ID'] = shopIdHeader
    }
    const response = await $fetch<ApiProductDetail & { shop_id?: string }>(`/api/products/${idOrSlug}`, {
      headers: Object.keys(headers).length > 0 ? headers : undefined
    })

    if (response && response.id) {
      return this.mapApiProduct(response)
    }
    return null
  },

  /**
   * Fetch shop-specific products and inventory directly (GET /shops/{shopId}/products).
   */
  async getShopProducts(shopId: string): Promise<CatalogProduct[]> {
    try {
      const response = await bootstrapConfig.fetchApi<any>(`/shops/${shopId}/products`)
      if (response && Array.isArray(response.products)) {
        return response.products.map((p: any) => ({
          id: p.id,
          name: p.name,
          price: p.price,
          rating: 4.8,
          reviews: 180,
          image: p.banner?.thumbnail || '',
          tag: p.sku ? p.sku.split('-')[1] || 'Collection' : 'Collection',
          desc: p.description || `Premium quality ${p.name} crafted for your special moments.`,
          isAvailable: p.status === 'active' && (p.inventory?.available > 0 ?? true),
          slug: p.slug,
          sku: p.sku,
          status: p.status,
          stock: p.inventory?.available ?? p.stock,
          inventory: p.inventory,
          shopId: response.shop_id || shopId
        }))
      }
    } catch (e) {
      console.error('Failed to fetch shop products:', e)
    }
    return []
  }
}

