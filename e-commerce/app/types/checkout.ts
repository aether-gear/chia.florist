import type { CustomDesignPayload } from '~/composables/useCart'
import type { ItemOptions } from '~/types/cart'

export interface CheckoutAddress {
  id: string
  recipient_name: string
  phone: string | null
  full_address: string
}

export interface CheckoutCourierOption {
  code: string
  name: string
  service: string
  etd: string
  fee: number
}

export interface SelectedCourier {
  code: string
  service: string
  fee: number
}

export interface CheckoutItem {
  product_id?: string
  cart_item_id?: string
  shop_id: string
  name: string
  price: number
  quantity: number
  subtotal: number
  product_variant_type?: 'standard' | 'custom'
  item_type?: 'standard' | 'custom'
  custom_design?: CustomDesignPayload
  item_options?: ItemOptions
  size?: string
  jambul?: string
  color?: string
}

export interface CheckoutShop {
  shop_id: string
  name?: string
  subtotal: number
  total: number
  selected_courier: SelectedCourier | null
  items: CheckoutItem[]
  cost_couriers: CheckoutCourierOption[] | null
}

export interface PaymentMethod {
  id: string
  name: string
  type: string
  description: string
  fee: number
  subtotal: number
  total: number
}

export interface CheckoutResponse {
  address: CheckoutAddress
  shops: CheckoutShop[]
  subtotal: number
  total_shipping: number
  total: number
  payment_methods?: PaymentMethod[]
  selected_payment_method?: PaymentMethod
}

export interface SelectedCourierInput {
  code: string
  service: string
}

export interface CheckoutShopItemInput {
  product_variant_type?: 'standard' | 'custom'
  item_type?: 'standard' | 'custom'
  product_id?: string
  cart_item_id?: string
  product_name?: string
  physical_size_id?: string
  quantity: number
  unit_price?: number
  custom_design?: CustomDesignPayload
  item_options?: ItemOptions
  size?: string
  jambul?: string
}

export type CheckoutItemInput = CheckoutShopItemInput

export interface CheckoutShopInput {
  shop_id: string
  items: CheckoutShopItemInput[]
  courier?: SelectedCourierInput
}

export interface CheckoutRequest {
  address_id?: string
  payment_method_id?: string
  shops: CheckoutShopInput[]
}

