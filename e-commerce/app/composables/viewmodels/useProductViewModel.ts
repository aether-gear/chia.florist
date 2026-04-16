import { ref, computed } from 'vue'
import { productService } from '~/services/productService'
import type { Product } from '~/models/Product'

export const useProductViewModel = () => {
  const products = ref<Product[]>([])
  const currentProduct = ref<Product | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const fetchProducts = async () => {
    isLoading.value = true
    error.value = null
    try {
      products.value = await productService.getProducts()
    } catch (e) {
      error.value = 'Failed to load products'
    } finally {
      isLoading.value = false
    }
  }

  const fetchProductById = async (id: string) => {
    isLoading.value = true
    error.value = null
    try {
      currentProduct.value = await productService.getProductById(id)
    } catch (e) {
      error.value = 'Failed to load product details'
    } finally {
      isLoading.value = false
    }
  }

  return {
    products: computed(() => products.value),
    currentProduct: computed(() => currentProduct.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    fetchProducts,
    fetchProductById
  }
}
