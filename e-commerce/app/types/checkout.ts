// app/types/checkout.ts

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
  product_id: string
  shop_id: string
  name: string
  price: number
  quantity: number
  subtotal: number
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

export interface CheckoutShopInput {
  shop_id: string
  items: {
    product_id: string
    quantity: number
  }[]
  courier?: SelectedCourierInput
}

export interface CheckoutRequest {
  address_id?: string
  payment_method_id?: string
  shops: CheckoutShopInput[]
}
