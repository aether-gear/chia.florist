// app/composables/useAddress.ts
import { ref } from 'vue'
import { addressService } from '~/services/addressService'
import type { UserAddress } from '~/types/address'
import { mapErrorMessage } from '~/utils/errorMessages'

let fetchAddressesPromise: Promise<void> | null = null
let addressWatcherInitialized = false

// Memory caches for location dropdowns
const cachedProvinces = ref<any[] | null>(null)
const cachedCities = new Map<string, any[]>()
const cachedDistricts = new Map<string, any[]>()
const cachedVillages = new Map<string, any[]>()

export const useAddress = () => {
  const addresses = useState<UserAddress[]>('chia-user-addresses-state', () => [])
  const isLoading = useState<boolean>('chia-user-addresses-loading', () => false)
  const error = ref<string | null>(null)
  const isLoggedIn = useCookie('is_logged_in')

  const clearAddresses = () => {
    addresses.value = []
    isLoading.value = false
    error.value = null
    fetchAddressesPromise = null
    if (import.meta.client) {
      try {
        localStorage.removeItem('chia-florist-addresses-cache')
        // Clean any address-related keys from localStorage if any
        const keysToRemove: string[] = []
        for (let i = 0; i < localStorage.length; i++) {
          const key = localStorage.key(i)
          if (key && (key.startsWith('address_') || key.includes('address'))) {
            keysToRemove.push(key)
          }
        }
        keysToRemove.forEach(k => localStorage.removeItem(k))
      } catch (e) {
        console.error('Error removing address cache:', e)
      }
    }
  }

  if (import.meta.client && !addressWatcherInitialized) {
    addressWatcherInitialized = true
    watch(isLoggedIn, (newVal) => {
      if (newVal === 'true') {
        clearAddresses()
        fetchAddresses(true)
      } else {
        clearAddresses()
      }
    }, { immediate: false })
  }

  const fetchAddresses = (force = false): Promise<void> => {
    if (isLoggedIn.value !== 'true') {
      clearAddresses()
      return Promise.resolve()
    }

    if (!force && addresses.value.length > 0) {
      return Promise.resolve()
    }

    if (!force && fetchAddressesPromise) {
      return fetchAddressesPromise
    }

    isLoading.value = true
    error.value = null

    fetchAddressesPromise = (async () => {
      try {
        const response = await addressService.listAddresses()
        addresses.value = response.addresses || []
      } catch (err: any) {
        error.value = mapErrorMessage(err, 'Gagal memuat daftar alamat.')
      } finally {
        isLoading.value = false
        fetchAddressesPromise = null
      }
    })()

    return fetchAddressesPromise
  }

  const saveAddress = async (address: UserAddress) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await addressService.saveAddress(address)
      await fetchAddresses(true)
      return { success: true, message: res.message }
    } catch (err: any) {
      const errMsg = mapErrorMessage(err, 'Gagal menyimpan alamat. Silakan coba lagi.')
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
      await fetchAddresses(true)
      return { success: true, message: res.message }
    } catch (err: any) {
      const errMsg = mapErrorMessage(err, 'Gagal menghapus alamat. Silakan coba lagi.')
      error.value = errMsg
      return { success: false, message: errMsg }
    } finally {
      isLoading.value = false
    }
  }

  // Location loaders with caching
  const loadProvinces = async () => {
    if (cachedProvinces.value && cachedProvinces.value.length > 0) {
      return cachedProvinces.value
    }
    try {
      const res = await addressService.getProvinces()
      cachedProvinces.value = res.provinces || []
      return cachedProvinces.value
    } catch (err) {
      console.error('Error fetching provinces:', err)
      return []
    }
  }

  const loadCities = async (provId: string) => {
    if (!provId) return []
    if (cachedCities.has(provId)) {
      return cachedCities.get(provId)!
    }
    try {
      const res = await addressService.getCities(provId)
      const list = res.cities || []
      cachedCities.set(provId, list)
      return list
    } catch (err) {
      console.error('Error fetching cities:', err)
      return []
    }
  }

  const loadDistricts = async (cityId: string) => {
    if (!cityId) return []
    if (cachedDistricts.has(cityId)) {
      return cachedDistricts.get(cityId)!
    }
    try {
      const res = await addressService.getDistricts(cityId)
      const list = res.districts || []
      cachedDistricts.set(cityId, list)
      return list
    } catch (err) {
      console.error('Error fetching districts:', err)
      return []
    }
  }

  const loadVillages = async (distId: string) => {
    if (!distId) return []
    if (cachedVillages.has(distId)) {
      return cachedVillages.get(distId)!
    }
    try {
      const res = await addressService.getVillages(distId)
      const list = res.villages || []
      cachedVillages.set(distId, list)
      return list
    } catch (err) {
      console.error('Error fetching villages:', err)
      return []
    }
  }

  return {
    addresses,
    isLoading,
    error,
    clearAddresses,
    fetchAddresses,
    saveAddress,
    deleteAddress,
    loadProvinces,
    loadCities,
    loadDistricts,
    loadVillages
  }
}
