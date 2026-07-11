// app/services/orderService.ts
import { bootstrapConfig } from '~/utils/bootstrap'
import type { CreateOrderRequest, CreateOrderResponse, ListOrdersQuery, ListOrdersResponse, BackendOrder, GetOrderPaymentDetailsResponse } from '~/types/order'

export const orderService = {
  async createOrder(data: CreateOrderRequest): Promise<CreateOrderResponse> {
    return bootstrapConfig.fetchApi<CreateOrderResponse>('/order', {
      method: 'POST',
      body: data
    })
  },

  async listOrders(query: ListOrdersQuery = {}): Promise<ListOrdersResponse> {
    const params = new URLSearchParams()
    if (query.page)   params.set('page',   String(query.page))
    if (query.limit)  params.set('limit',  String(query.limit))
    if (query.sort)   params.set('sort',   query.sort)
    if (query.status) params.set('status', query.status)
    const qs = params.toString()
    return bootstrapConfig.fetchApi<ListOrdersResponse>(
      `/users/me/orders${qs ? `?${qs}` : ''}`
    )
  },

  async getOrder(orderId: string): Promise<BackendOrder> {
    return bootstrapConfig.fetchApi<BackendOrder>(`/users/me/orders/${orderId}`)
  },

  async getOrderPaymentDetails(orderId: string): Promise<GetOrderPaymentDetailsResponse> {
    return bootstrapConfig.fetchApi<GetOrderPaymentDetailsResponse>(`/users/me/orders/${orderId}/payment`)
  }
}
