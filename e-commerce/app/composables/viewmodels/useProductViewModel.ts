// app/composables/viewmodels/useProductViewModel.ts
import { ref, computed } from 'vue'
import { productService } from '~/services/productService'
import type { Product, CatalogProduct, FindProductsParams } from '~/types/product'
import { useStoreSelection } from '~/composables/useStoreSelection'
import { mapErrorMessage } from '~/utils/errorMessages'

export const useProductViewModel = () => {
  const products = ref<Product[]>([])
  const catalogProducts = ref<CatalogProduct[]>([])
  const currentProduct = ref<Product | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const storeSelection = useStoreSelection()

  // Pagination metadata
  const page = ref(1)
  const limit = ref(10)
  const total = ref(0)
  const totalPages = computed(() => (limit.value > 0 ? Math.ceil(total.value / limit.value) : 1))

  /**
   * Fetch all products for the catalog page with filtering and pagination.
   */
  const fetchCatalogProducts = async (params?: FindProductsParams, shopIdHeader?: string) => {
    isLoading.value = true
    error.value = null
    try {
      const activeShop = storeSelection.selectedShop.value
      const queryParams: FindProductsParams = { ...params }
      if (!queryParams.shop_id && !queryParams.shop_slug && activeShop) {
        if (activeShop.id) queryParams.shop_id = activeShop.id
        else if (activeShop.slug) queryParams.shop_slug = activeShop.slug
      }

      const activeShopIdHeader = shopIdHeader || (activeShop?.id ?? undefined)

      const res = await productService.getPaginatedCatalogProducts(queryParams, activeShopIdHeader)
      catalogProducts.value = res.products
      page.value = res.page
      limit.value = res.limit
      total.value = res.total

      if (catalogProducts.value.length === 0) {
        error.value = 'Produk sedang tidak tersedia'
      }
    } catch (e) {
      error.value = mapErrorMessage(e, 'Produk sedang tidak tersedia')
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Fetch a specific product's details by ID or slug.
   */
  const fetchProductById = async (id: string, shopIdHeader?: string) => {
    isLoading.value = true
    error.value = null
    try {
      const activeShopIdHeader = shopIdHeader || (storeSelection.selectedShop.value?.id ?? undefined)
      const result = await productService.getProductById(id, activeShopIdHeader)
      if (result) {
        currentProduct.value = result
      } else {
        error.value = 'Produk tidak ditemukan atau sedang tidak tersedia'
      }
    } catch (e) {
      error.value = mapErrorMessage(e, 'Produk tidak ditemukan atau sedang tidak tersedia')
    } finally {
      isLoading.value = false
    }
  }

  return {
    products: computed(() => products.value),
    catalogProducts: computed(() => catalogProducts.value),
    currentProduct: computed(() => currentProduct.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    page: computed(() => page.value),
    limit: computed(() => limit.value),
    total: computed(() => total.value),
    totalPages,
    fetchCatalogProducts,
    fetchProductById
  }
}

