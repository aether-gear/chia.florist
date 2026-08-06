export interface AnalyticsFilterState {
  from?: string;
  to?: string;
  granularity?: 'daily' | 'weekly' | 'monthly';
  shop_id?: string;
  top_n?: number;
}

// Order Metrics
export interface OrderMetricsSummary {
  total_orders: number;
  total_gmv: number;
  total_revenue: number;
  total_shipping_fee: number;
  aov: number;
  cancellation_rate: number;
  pending_count: number;
  confirmed_count: number;
  processing_count: number;
  shipped_count: number;
  delivered_count: number;
  cancelled_count: number;
}

export interface OrderTimeSeriesPoint {
  date: string;
  order_count: number;
  gmv: number;
  aov: number;
}

export interface TopProduct {
  product_id: string;
  product_name: string;
  quantity: number;
  revenue: number;
}

export interface TopShop {
  shop_id: string;
  shop_name: string;
  revenue: number;
  orders: number;
}

export interface OrderMetricsResponse {
  summary: OrderMetricsSummary;
  time_series: OrderTimeSeriesPoint[];
  top_products: TopProduct[];
  top_shops: TopShop[];
}

// Payment Metrics
export interface PaymentMetricsSummary {
  total_paid: number;
  total_pending: number;
  total_expired: number;
  total_refunded: number;
  payment_success_rate: number;
  avg_time_to_pay: number;
}

export interface PaymentMethodBreakdown {
  method_id: string;
  method_name: string;
  method_type: string;
  count: number;
  amount: number;
  success_rate: number;
}

export interface PaymentMetricsResponse {
  summary: PaymentMetricsSummary;
  breakdown: PaymentMethodBreakdown[];
}

// Shipment Metrics
export interface ShipmentMetricsSummary {
  total: number;
  delivered: number;
  failed: number;
  returned: number;
  cancelled: number;
  delivery_rate: number;
  avg_fulfillment_sec: number;
}

export interface CourierBreakdown {
  courier: string;
  service: string;
  count: number;
  delivery_rate: number;
  avg_cost: number;
}

export interface ShipmentMetricsResponse {
  summary: ShipmentMetricsSummary;
  couriers: CourierBreakdown[];
}

// Inventory Metrics
export interface InventoryMetricsResponse {
  total_products: number;
  total_stock: number;
  total_reserved: number;
  total_available: number;
  stockout_count: number;
  low_stock_count: number;
}

// Product Metrics
export interface ProductRevenueRank {
  product_id: string;
  product_name: string;
  revenue: number;
  units_sold: number;
  conversion_rate?: number;
  return_rate?: number | null;
  gross_margin_pct: number;
  sales_velocity_7d: number;
  sales_velocity_30d: number;
}

export interface ProductVolumeRank {
  product_id: string;
  product_name: string;
  revenue: number;
  units_sold: number;
  conversion_rate?: number;
  return_rate?: number | null;
  gross_margin_pct: number;
  sales_velocity_7d: number;
  sales_velocity_30d: number;
}

export interface ProductMetricsResponse {
  top_by_revenue: ProductRevenueRank[];
  top_by_volume: ProductVolumeRank[];
  avg_conversion: number;
  avg_return_rate: number;
  invoice_void_rate: number;
}
