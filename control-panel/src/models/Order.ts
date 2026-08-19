export interface OrderItem {
  id: string;
  shipment_id?: string | null;
  product_id?: string | null;
  product_variant_type?: string;
  is_custom?: boolean;
  product_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  shop_id: string;
  shop_name: string;
  courier_code?: string | null;
  courier_service?: string | null;
  shipping_fee: number;
}

export interface PaymentChannelData {
  channel_type: string;
  display_name: string;
  action_url?: string | null;
  expires_at?: string | null;
}

export interface OrderPayment {
  id: string;
  status: string;
  provider: string;
  amount: number;
  expires_at?: string | null;
  channel_data?: PaymentChannelData | null;
  created_at: string;
}

export interface ShipmentEvent {
  id: string;
  status: string;
  description: string;
  location: string;
  timestamp: string;
}

export interface OrderShipment {
  id: string;
  order_id?: string;
  status: string;
  fulfillment_method: string;
  courier: string;
  service: string;
  tracking_number: string | null;
  cost: number;
  created_at: string;
  events?: ShipmentEvent[];
  item_ids?: string[];
}

export interface OrderAddress {
  id: string;
  customer_id: string;
  receiver_name: string;
  phone: string;
  is_default: boolean;
  province_id: string;
  city_id: string;
  district_id: string;
  village_id: string;
  full_address: string;
  postal_code: string;
}

export interface Order {
  id: string;
  number: string;
  customer_id?: string;
  user_id?: string;
  address_id: string;
  status: string;
  subtotal: number;
  shipping_fee: number;
  total: number;
  confirmed_at?: string | null;
  handling_expires_at?: string | null;
  created_at: string;
  updated_at: string | null;
  items: OrderItem[];
  payment?: OrderPayment;
  shipment?: OrderShipment;
  shipments?: OrderShipment[];
  address?: OrderAddress;
}

export interface OrdersResponse {
  orders: Order[];
  page: number;
  limit: number;
  total: number;
}
