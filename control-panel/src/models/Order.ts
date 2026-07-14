export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  shop_id: string;
  shop_name: string;
  courier_code: string;
  courier_service: string;
  shipping_fee: number;
}

export interface OrderPayment {
  id: string;
  status: string;
  provider: string;
  amount: number;
  expires_at: string;
  created_at: string;
}

export interface OrderShipment {
  id: string;
  status: string;
  fulfillment_method: string;
  courier: string;
  service: string;
  tracking_number: string | null;
  cost: number;
  created_at: string;
}

export interface Order {
  id: string;
  number: string;
  user_id: string;
  address_id: string;
  status: string;
  subtotal: number;
  shipping_fee: number;
  total: number;
  created_at: string;
  updated_at: string | null;
  items: OrderItem[];
  payment?: OrderPayment;
  shipment?: OrderShipment;
}

export interface OrdersResponse {
  orders: Order[];
  page: number;
  limit: number;
  total: number;
}
