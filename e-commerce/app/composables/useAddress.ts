// app/composables/useAddress.ts
import { ref } from 'vue'
import { addressService } from '~/services/addressService'
import type { UserAddress } from '~/types/address'

export const useAddress = () => {
  const addresses = ref<UserAddress[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const fetchAddresses = async () => {
    isLoading.value = true
    error.value = null
    try {
      const response = await addressService.listAddresses()
      addresses.value = response.addresses || []
    } catch (err: any) {
      error.value = err.data?.message || err.message || 'Failed to load addresses'
    } finally {
      isLoading.value = false
    }
  }

  const saveAddress = async (address: UserAddress) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await addressService.saveAddress(address)
      await fetchAddresses()
      return { success: true, message: res.message }
    } catch (err: any) {
      const errMsg = err.data?.message || err.message || 'Failed to save address'
      error.value = errMsg
      return { success: false, message: errMsg }
    } finally {
      isLoading.value = false
    }
  }

  const deleteAddress = async (addressId: string) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await addressService.deleteAddress(addressId)
      await fetchAddresses()
      return { success: true, message: res.message }
    } catch (err: any) {
      const errMsg = err.data?.message || err.message || 'Failed to delete address'
      error.value = errMsg
      return { success: false, message: errMsg }
    } finally {
      isLoading.value = false
    }
  }

  // Location loaders
  const loadProvinces = async () => {
    try {
      const res = await addressService.getProvinces()
      return res.provinces || []
    } catch (err) {
      console.error('Error fetching provinces:', err)
      return []
    }
  }

  const loadCities = async (provId: string) => {
    if (!provId) return []
    try {
      const res = await addressService.getCities(provId)
      return res.cities || []
    } catch (err) {
      console.error('Error fetching cities:', err)
      return []
    }
  }

  const loadDistricts = async (cityId: string) => {
    if (!cityId) return []
    try {
      const res = await addressService.getDistricts(cityId)
      return res.districts || []
    } catch (err) {
      console.error('Error fetching districts:', err)
      return []
    }
  }

  const loadVillages = async (distId: string) => {
    if (!distId) return []
    try {
      const res = await addressService.getVillages(distId)
      return res.villages || []
    } catch (err) {
      console.error('Error fetching villages:', err)
      return []
    }
  }

  return {
    addresses,
    isLoading,
    error,
    fetchAddresses,
    saveAddress,
    deleteAddress,
    loadProvinces,
    loadCities,
    loadDistricts,
    loadVillages
  }
}
