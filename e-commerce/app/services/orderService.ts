// app/services/orderService.ts
import { bootstrapConfig } from '~/utils/bootstrap'
import type { CreateOrderRequest, CreateOrderResponse } from '~/types/order'

export const orderService = {
  async createOrder(data: CreateOrderRequest): Promise<CreateOrderResponse> {
    return bootstrapConfig.fetchApi<CreateOrderResponse>('/order', {
      method: 'POST',
      body: data
    })
  }
}
