// app/types/address.ts

export interface UserAddress {
  address_id?: string
  user_id?: string
  receiver_name: string
  phone: string
  is_default: boolean
  province_id: string
  city_id: string
  district_id: string
  village_id: string
  full_address: string
  postal_code: string
  created_at?: string
}

export interface SaveAddressPayload {
  address_id?: string
  receiver_name: string
  phone: string
  is_default: string // "True" or "False"
  province_id: string
  city_id: string
  district_id: string
  village_id: string
  full_address: string
  postal_code: string
}

export interface ListAddressesResponse {
  addresses: UserAddress[]
}

export interface LocationItem {
  id: string
  name: string
}
