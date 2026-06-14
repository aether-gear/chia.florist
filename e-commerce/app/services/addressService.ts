// app/services/addressService.ts
import { bootstrapConfig } from '~/utils/bootstrap'
import type { ListAddressesResponse, SaveAddressPayload, UserAddress } from '~/types/address'

export const addressService = {
  /**
   * List all addresses for the authenticated user
   */
  async listAddresses(): Promise<ListAddressesResponse> {
    return bootstrapConfig.fetchApi<ListAddressesResponse>('/users/me/addresses/', {
      method: 'GET'
    })
  },

  /**
   * Save (create or update) a user address
   */
  async saveAddress(address: UserAddress): Promise<{ message: string }> {
    const payload: SaveAddressPayload = {
      address_id: address.address_id || undefined,
      receiver_name: address.receiver_name,
      phone: address.phone,
      is_default: address.is_default ? 'True' : 'False',
      province_id: address.province_id,
      city_id: address.city_id,
      district_id: address.district_id,
      village_id: address.village_id,
      full_address: address.full_address,
      postal_code: address.postal_code
    }

    return bootstrapConfig.fetchApi<{ message: string }>('/users/me/addresses/', {
      method: 'POST',
      body: payload
    })
  },

  /**
   * Delete a user address
   */
  async deleteAddress(addressId: string): Promise<{ message: string }> {
    return bootstrapConfig.fetchApi<{ message: string }>(`/users/me/addresses/${addressId}`, {
      method: 'DELETE'
    })
  },

  /**
   * Fetch all provinces
   */
  async getProvinces(): Promise<{ provinces: { id: string; name: string }[] }> {
    return bootstrapConfig.fetchApi<{ provinces: { id: string; name: string }[] }>('/provinces/')
  },

  /**
   * Fetch cities in a province
   */
  async getCities(provinceId: string): Promise<{ cities: { id: string; name: string }[] }> {
    return bootstrapConfig.fetchApi<{ cities: { id: string; name: string }[] }>(`/provinces/${provinceId}/cities`)
  },

  /**
   * Fetch districts in a city
   */
  async getDistricts(cityId: string): Promise<{ districts: { id: string; name: string }[] }> {
    return bootstrapConfig.fetchApi<{ districts: { id: string; name: string }[] }>(`/cities/${cityId}/districts`)
  },

  /**
   * Fetch villages in a district
   */
  async getVillages(districtId: string): Promise<{ villages: { id: string; name: string }[] }> {
    return bootstrapConfig.fetchApi<{ villages: { id: string; name: string }[] }>(`/districts/${districtId}/villages`)
  }
}
