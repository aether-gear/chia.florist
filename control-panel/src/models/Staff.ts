export interface Staff {
  id: string;
  user_id: string;
  name: string;
  username: string;
  phone?: string | null;
  avatar_url?: string | null;
  created_at: string;
  updated_at?: string | null;
  last_login_at?: string | null;
}

export interface StaffListResponse {
  page: number;
  limit: number;
  total: number;
  staff: Staff[];
}

export interface AddStaffAccountPayload {
  email: string;
  name: string;
  username: string;
  password?: string;
  phone?: string;
}

export interface CreateStaffPayload {
  name: string;
  description?: string;
  logo_url?: string;
  banner_url?: string;
}
