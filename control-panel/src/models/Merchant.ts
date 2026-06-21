export interface Merchant {
  id: string;
  name: string;
  description: string;
  logo_url: string;
  banner_url: string;
  created_at: string;
}

export interface MerchantsResponse {
  page: number;
  limit: number;
  total: number;
  merchants: Merchant[];
}
