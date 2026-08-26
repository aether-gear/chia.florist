import type { CustomDesignPayload } from '~/composables/useCart'
import type { ItemOptions } from '~/types/cart'

export interface CreateOrderPaymentInput {
  id: string
}

export interface CreateOrderCourierInput {
  code: string
  service: string
}

export interface CreateOrderItemInput {
  product_id?: string
  cart_item_id?: string
  product_variant_type?: 'standard' | 'custom'
  item_type?: 'standard' | 'custom'
  is_custom?: boolean
  name: string
  quantity: number
  item_options?: ItemOptions
  size?: string
  jambul?: string
  custom_design?: CustomDesignPayload
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

export interface PaymentChannelDataResponse {
  channel_type: string
  display_name: string
  action_url?: string
  expires_at?: string
}

export interface CreateOrderResponse {
  order_id: string
  instruction: string
  channel_data?: PaymentChannelDataResponse
}

// ─── Domain Status Definitions ─────────────────────────────────────

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'expired'

export type PaymentStatus =
  | 'pending'
  | 'paid'
  | 'failed'
  | 'expired'
  | 'cancelled'
  | 'refunded'
  | 'refund_pending'
  | 'refund_failed'

// ─── List Orders API ───────────────────────────────────────────────

export interface ListOrdersQuery {
  page?: number
  limit?: number
  sort?: string
  status?: string
}

export interface BackendOrderPayment {
  id: string
  status: PaymentStatus | string
  provider: string
  amount: number
  expires_at: string
  created_at: string
}

export interface BackendOrderItem {
  id: string
  product_id: string | null
  product_name: string
  product_variant_type?: 'standard' | 'custom'
  is_custom?: boolean
  quantity: number
  unit_price: number
  subtotal: number
  shop_id: string
  shop_name: string
  courier_code: string
  courier_service: string
  shipping_fee: number
  item_options?: ItemOptions
}

export interface BackendOrder {
  id: string
  number: string
  customer_id: string
  address_id: string
  status: OrderStatus | string
  subtotal: number
  shipping_fee: number
  total: number
  created_at: string
  items: BackendOrderItem[]
  payment?: BackendOrderPayment
}

export interface ListOrdersResponse {
  limit: number
  orders: BackendOrder[]
  page: number
  total: number
}

export interface GetOrderPaymentDetailsResponse {
  payment_id: string
  status: PaymentStatus | string
  amount: number
  expires_at: string
  channel_type?: string
  display_name?: string
  action_url?: string
  channel_data?: PaymentChannelDataResponse
  account_name?: string
  account_number?: string
  phone_number?: string
  qr_string?: string
  instruction?: string
}

export interface CheckOrderPaymentStatusResponse {
  status: PaymentStatus | string
  synced: boolean
}

// ─── Order Tracking Timeline API ─────────────────────────────────

export interface TimelineEvent {
  status: string
  description: string
  location?: string
  timestamp: string
}

export interface GetOrderTrackingTimelineResponse {
  order_id: string
  shipment_id: string
  courier: string
  tracking_number?: string
  timeline: TimelineEvent[]
}



