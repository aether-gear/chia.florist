// app/composables/viewmodels/useProductViewModel.ts
import { ref, computed } from 'vue'
import { productService } from '~/services/productService'
import type { Product, CatalogProduct } from '~/types/product'

export const useProductViewModel = () => {
  const products = ref<Product[]>([])
  const catalogProducts = ref<CatalogProduct[]>([])
  const currentProduct = ref<Product | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Fetch all products for the catalog page and map them.
   */
  const fetchCatalogProducts = async (params?: { name?: string; id?: string; page?: number; limit?: number; sort?: string }) => {
    isLoading.value = true
    error.value = null
    try {
      catalogProducts.value = await productService.getCatalogProducts(params)
      if (catalogProducts.value.length === 0) {
        error.value = 'Produk sedang tidak tersedia'
      }
    } catch (e) {
      error.value = 'Produk sedang tidak tersedia'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Fetch a specific product's details by ID.
   */
  const fetchProductById = async (id: string) => {
    isLoading.value = true
    error.value = null
    try {
      const result = await productService.getProductById(id)
      if (result) {
        currentProduct.value = result
      } else {
        error.value = 'Produk tidak ditemukan atau sedang tidak tersedia'
      }
    } catch (e) {
      error.value = 'Produk tidak ditemukan atau sedang tidak tersedia'
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
    fetchCatalogProducts,
    fetchProductById
  }
}
