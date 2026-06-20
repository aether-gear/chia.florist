export interface Customer {
  id: string;
  name: string;
  username: string;
  phone: string;
  last_login_at: string | null;
}

export interface CustomersResponse {
  page: number;
  limit: number;
  total: number;
  users: Customer[];
}
