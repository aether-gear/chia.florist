export interface ProductAvailability {
  slug: string;
  name: string;
  stock: number;
}

export interface ProductBanner {
  thumbnail: string | null;
  preview: string | null;
  detail: string | null;
}

export interface Product {
  id: string;
  sku: string;
  name: string;
  slug: string;
  status: 'active' | 'inactive' | 'archived';
  is_available: boolean;
  price: number;
  stock: number;
  banner: ProductBanner;
  gallery?: ProductBanner[];
  description?: string | null;
  weight?: number | null;
  availability: ProductAvailability[];
}

export interface ProductsResponse {
  limit: number;
  page: number;
  total: number;
  products: Product[];
}

export interface ProductStat {
  id: string;
  sku: string;
  name: string;
  slug: string;
  status: 'active' | 'inactive' | 'archived';
  price: number;
  cost_price: number | null;
  supplier_lead_time_days: number | null;
  gross_margin_pct: number | null;
  view_count: number;
  stock: number;
  sales_velocity_7d: number;
  sales_velocity_30d: number;
  sales_velocity_90d: number;
  conversion_rate: number;
  revenue_contribution_percentage: number;
  return_rate: number | null;
  average_rating: number | null;
  review_count: number | null;
  thumbnail: string | null;
}

export interface ProductStatsResponse {
  page: number;
  limit: number;
  total: number;
  stats: ProductStat[];
}

