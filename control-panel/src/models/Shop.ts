export interface ShopAddress {
  id: string;
  label: string;
  phone: string;
  is_active: boolean;
  province_id: string;
  city_id: string;
  district_id: string;
  village_id: string;
  full_address: string;
  postal_code: string;
  created_at: string;
  updated_at: string | null;
}

export interface ShopCourier {
  code: string;
  active: boolean;
}

export interface ShopProductInventory {
  total_stock: number;
  reserved_stock: number;
  available: number;
}

export interface ShopProduct {
  id: string;
  sku: string;
  name: string;
  slug: string;
  description: string;
  status: 'active' | 'inactive' | 'archived';
  price: number;
  weight: number;
  inventory: ShopProductInventory;
  created_at: string;
  updated_at: string | null;
}

export interface ShopAddressesResponse {
  shop_id: string;
  addresses: ShopAddress[];
}

export interface ShopCouriersResponse {
  shop_id: string;
  couriers: ShopCourier[];
}

export interface ShopProductsResponse {
  shop_id: string;
  products: ShopProduct[];
}
