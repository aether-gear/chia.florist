export interface PaymentMethod {
  id: string;
  name: string;
  type: string;
  is_active: boolean;
  description: string;
  fee_type: string;
  fee_fixed: number;
  fee_percentage: number;
}

export interface PaymentAccount {
  id: string;
  method_id: string;
  account_name: string;
  account_number: string | null;
  phone_number: string | null;
  qr_string: string | null;
}

export interface PaymentMethodsResponse {
  methods: PaymentMethod[];
}

export interface PaymentAccountsResponse {
  accounts: PaymentAccount[];
}
