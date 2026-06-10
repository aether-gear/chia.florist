export interface MerchantIdentity {
  accountId?: string;
  merchantId?: string;
  name: string;
  slug: string;
  profilePhoto: string;
  coverBanner: string;
  description?: string;
}

export interface MerchantContact {
  email: string;
  phone: string;
  whatsapp?: string;
  customerServiceContact?: string;
  address: string;
  country: string;
  province: string;
  city: string;
  district?: string;
  postalCode?: string;
  fullAddress?: string;
  latitude?: number | null;
  longitude?: number | null;
}

export interface MerchantSettings {
  preferredLanguage: string;
  preferredCurrency: string;
  timezone: string;
}

export interface MerchantOperational {
  openingHours: string;
  closingHours: string;
  businessDays: string[];
  deliveryRadius: number;
  pickupAvailable: boolean;
}

export interface MerchantFinancial {
  bankAccountName: string;
  bankName: string;
  bankAccountNumber: string;
  eWalletInformation?: string;
  taxNumber?: string;
}

export interface MerchantProfile {
  identity: MerchantIdentity;
  contact: MerchantContact;
  settings: MerchantSettings;
  operational: MerchantOperational;
  financial: MerchantFinancial;
}
