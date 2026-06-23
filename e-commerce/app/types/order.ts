// app/types/order.ts

export interface CreateOrderPaymentInput {
  id: string
  is_manual: boolean
}

export interface CreateOrderCourierInput {
  code: string
  service: string
}

export interface CreateOrderItemInput {
  product_id: string
  name: string
  quantity: number
}

export interface CreateOrderShopInput {
  shop_id: string
  name: string
  selected_courier: CreateOrderCourierInput
  items: CreateOrderItemInput[]
}

export interface CreateOrderRequest {
  address_id: string
  selected_payment: CreateOrderPaymentInput
  shops: CreateOrderShopInput[]
}

export interface CreateOrderPaymentAccountResponse {
  account_name: string
  account_number?: string
  phone_number?: string
  qr_string?: string
}

export interface CreateOrderResponse {
  order_id: string
  instruction: string
  payment_account?: CreateOrderPaymentAccountResponse
}
