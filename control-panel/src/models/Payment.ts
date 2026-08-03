export type FeeType = 'flat' | 'percentage' | 'mixed';
export type MethodType = 'bank_transfer' | 'ewallet' | 'qr_code';
export type MethodCode = 'gopay' | 'shopeepay' | 'qris' | 'bca_va' | 'mandiri_bill';

export const FEE_TYPES: { value: FeeType; label: string }[] = [
  { value: 'flat', label: 'Flat Fee' },
  { value: 'percentage', label: 'Percentage Fee' },
  { value: 'mixed', label: 'Mixed Fee' },
];

export const METHOD_TYPES: { value: MethodType; label: string }[] = [
  { value: 'bank_transfer', label: 'Bank Transfer' },
  { value: 'ewallet', label: 'E-Wallet' },
  { value: 'qr_code', label: 'QR Code' },
];

export const METHOD_CODES: { value: MethodCode; label: string }[] = [
  { value: 'gopay', label: 'GoPay' },
  { value: 'shopeepay', label: 'ShopeePay' },
  { value: 'qris', label: 'QRIS' },
  { value: 'bca_va', label: 'BCA Virtual Account' },
  { value: 'mandiri_bill', label: 'Mandiri Bill' },
];

export interface PaymentInstruction {
  id: string;
  content: string;
  created_at: string;
  updated_at?: string | null;
}

export interface PaymentMethod {
  id: string;
  name: string;
  code: MethodCode;
  provider: string;
  type: MethodType;
  is_active: boolean;
  description: string;
  fee_type: FeeType;
  fee_fixed: number;
  fee_percentage: number;
  instruction?: PaymentInstruction | null;
}

export interface PaymentMethodsResponse {
  methods: PaymentMethod[];
}
