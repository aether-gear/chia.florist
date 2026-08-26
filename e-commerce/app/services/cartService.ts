import { bootstrapConfig } from '~/utils/bootstrap'
import type { BackendCartResponse, ItemOptions } from '~/types/cart'
import type { CheckoutRequest, CheckoutResponse } from '~/types/checkout'
import type { CustomDesignPayload } from '~/composables/useCart'

export interface AddCartItemPayload {
  product_variant_type?: 'standard' | 'custom'
  shop_id: string
  quantity: number
  product_id?: string
  product_name?: string
  physical_size_id?: string
  unit_price?: number
  item_options?: ItemOptions
  size?: string
  jambul?: string
  custom_design?: CustomDesignPayload
}

export const cartService = {
  async getCart(): Promise<BackendCartResponse> {
    return bootstrapConfig.fetchApi<BackendCartResponse>('/carts', {
      method: 'GET'
    })
  },

  async addItem(data: AddCartItemPayload): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>('/carts/items', {
      method: 'POST',
      body: data
    })
  },

  async updateItem(shopId: string, productId: string, quantity: number): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>(`/carts/items/${shopId}/${productId}`, {
      method: 'PUT',
      body: { quantity }
    })
  },

  async removeItem(shopId: string, productId: string): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>(`/carts/items/${shopId}/${productId}`, {
      method: 'DELETE'
    })
  },

  async removeCustomItem(cartItemId: string): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>(`/carts/items/custom/${cartItemId}`, {
      method: 'DELETE'
    })
  },

  async updateCartItemShop(cartItemId: string, shopId: string): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>(`/carts/items/${cartItemId}/shop`, {
      method: 'PATCH',
      body: { shop_id: shopId }
    })
  },

  async checkout(data: CheckoutRequest): Promise<CheckoutResponse> {
    return bootstrapConfig.fetchApi<CheckoutResponse>('/carts/checkout', {
      method: 'POST',
      body: data
    })
  },

  async checkoutCalculate(data: CheckoutRequest): Promise<CheckoutResponse> {
    return bootstrapConfig.fetchApi<CheckoutResponse>('/carts/checkout/calculate', {
      method: 'POST',
      body: data
    })
  }
}

