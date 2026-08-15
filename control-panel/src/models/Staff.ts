export interface StaffAccountRole {
  id: string;
  code: string;
  name: string;
}

export interface StaffAccountMember {
  account_id: string;
  user_id: string;
  email: string;
  name: string;
  username: string;
  phone?: string | null;
  avatar_url?: string | null;
  role: StaffAccountRole;
  last_login_at?: string | null;
  created_at: string;
}

export interface StaffAccountsResponse {
  staff_id: string;
  total: number;
  accounts: StaffAccountMember[];
}

export interface Staff {
  id: string;
  user_id: string;
  name: string;
  username: string;
  description?: string | null;
  logo_url?: string | null;
  banner_url?: string | null;
  phone?: string | null;
  avatar_url?: string | null;
  created_at: string;
  updated_at?: string | null;
  last_login_at?: string | null;
  accounts?: StaffAccountMember[];
}

export interface StaffListResponse {
  page: number;
  limit: number;
  total: number;
  staff: Staff[];
}

export interface AddStaffAccountPayload {
  email: string;
  password?: string;
}

export interface CreateStaffPayload {
  name: string;
  username: string;
  description?: string;
  logo_url?: string;
  banner_url?: string;
}

export interface UpdateStaffPayload {
  name: string;
  description?: string;
  logo_url?: string;
  banner_url?: string;
}
