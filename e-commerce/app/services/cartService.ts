// app/services/cartService.ts
import { bootstrapConfig } from '~/utils/bootstrap'
import type { BackendCartResponse } from '~/types/cart'

export const cartService = {
  async getCart(): Promise<BackendCartResponse> {
    return bootstrapConfig.fetchApi<BackendCartResponse>('/carts/', {
      method: 'GET'
    })
  },

  async addItem(data: { product_id: string; shop_id: string; quantity: number }): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>('/carts/items/', {
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
  }
}
