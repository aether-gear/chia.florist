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
