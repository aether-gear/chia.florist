<!-- app/pages/profile.vue -->
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useCart } from '~/composables/useCart'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useAddress } from '~/composables/useAddress'
import { supabaseService } from '~/services/supabaseService'
import type { UserAddress } from '~/types/address'

useHead({ title: 'My Profile - Chia Florist' })

const authVm = useAuthViewModel()
const addressVm = useAddress()

const activeTab = ref('personal')
const activeOrderStatus = ref<'pending' | 'processing' | 'shipping' | 'done'>('pending')

// Edit address states
const showAddressForm = ref(false)
const editingAddress = ref<UserAddress | null>(null)

const provincesList = ref<any[]>([])
const citiesList = ref<any[]>([])
const districtsList = ref<any[]>([])
const villagesList = ref<any[]>([])

// User basic info state
const user = ref({
  name: 'Loading...',
  username: 'Loading...',
  email: 'Loading...',
  phone: 'Loading...',
  lastLoginAt: '' as string | null
})

// Supabase avatar states
const avatarUrl = ref<string | null>(null)
const publicAvatarUrl = ref<string | null>(null)
const signedAvatarUrl = ref<string | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const loadProfilePicture = async () => {
  const userId = authVm.currentUser.value?.id
  if (!userId) return

  const urls = await supabaseService.getAvatarUrls(userId)
  if (urls) {
    publicAvatarUrl.value = urls.publicUrl
    signedAvatarUrl.value = urls.signedUrl
    avatarUrl.value = urls.signedUrl || urls.publicUrl
  } else {
    avatarUrl.value = null
    publicAvatarUrl.value = null
    signedAvatarUrl.value = null
  }
}

const triggerFileSelect = () => {
  fileInput.value?.click()
}

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    alert('Please select an image file.')
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    alert('Image file size must be less than 5MB.')
    return
  }

  const userId = authVm.currentUser.value?.id
  if (!userId) {
    alert('User not authenticated.')
    return
  }

  isUploading.value = true
  uploadError.value = null
  try {
    const urls = await supabaseService.uploadFile(userId, file)
    if (urls) {
      publicAvatarUrl.value = urls.publicUrl
      signedAvatarUrl.value = urls.signedUrl
      avatarUrl.value = urls.signedUrl || urls.publicUrl
      
      if (authVm.currentUser.value) {
        authVm.currentUser.value.avatarUrl = urls.signedUrl || urls.publicUrl
      }

      // Sync avatar URL with Go backend
      await authVm.updateUserProfile({
        avatar_url: urls.signedUrl || urls.publicUrl
      })
    } else {
      uploadError.value = 'Failed to upload image. Please try again.'
    }
  } catch (err: any) {
    uploadError.value = err.message || 'An error occurred during file upload.'
  } finally {
    isUploading.value = false
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}

const handleRemovePicture = async () => {
  const userId = authVm.currentUser.value?.id
  if (!userId) return

  if (confirm('Are you sure you want to remove your profile picture?')) {
    isUploading.value = true
    try {
      const files = await supabaseService.listFiles(userId)
      if (files.length > 0) {
        await supabaseService.deleteFiles(userId, files.map(f => f.name))
      }
      avatarUrl.value = null
      publicAvatarUrl.value = null
      signedAvatarUrl.value = null
      if (authVm.currentUser.value) {
        authVm.currentUser.value.avatarUrl = null
      }

      // Clear avatar URL in Go backend
      await authVm.updateUserProfile({
        avatar_url: ''
      })
    } catch (err) {
      console.error('Failed to remove profile picture:', err)
      alert('Failed to remove profile picture.')
    } finally {
      isUploading.value = false
    }
  }
}

const handleUpdateProfile = async () => {
  if (!user.value.name) {
    alert('Full name is required.')
    return
  }

  try {
    const res = await authVm.updateUserProfile({
      name: user.value.name,
      username: user.value.username,
      phone: user.value.phone
    })
    if (res.success) {
      alert('Profile updated successfully!')
    } else {
      alert(res.message || 'Failed to update profile.')
    }
  } catch (err: any) {
    alert(err.message || 'An error occurred while updating profile.')
  }
}

// Sync profile data from ViewModel
watch(() => authVm.currentUser.value, (me) => {
  if (me) {
    user.value.name = me.name || ''
    user.value.username = me.username || ''
    user.value.email = me.email || ''
    user.value.phone = me.phone || ''
    user.value.lastLoginAt = me.last_login_at || null
  }
}, { immediate: true })

const loadAddressesData = async () => {
  await addressVm.fetchAddresses()
}

watch(activeTab, (tab) => {
  if (tab === 'addresses') {
    loadAddressesData()
  }
})

onMounted(async () => {
  await authVm.fetchCurrentUser()
  if (!authVm.isAuthenticated.value) {
    navigateTo('/login')
  } else {
    await loadProfilePicture()
  }
})

// Chained location selections
const initAddressForm = async (addr: UserAddress | null = null) => {
  provincesList.value = await addressVm.loadProvinces()
  
  if (addr) {
    editingAddress.value = { ...addr }
    citiesList.value = await addressVm.loadCities(addr.province_id)
    districtsList.value = await addressVm.loadDistricts(addr.city_id)
    villagesList.value = await addressVm.loadVillages(addr.district_id)
  } else {
    editingAddress.value = {
      receiver_name: '',
      phone: '',
      is_default: false,
      province_id: '',
      city_id: '',
      district_id: '',
      village_id: '',
      full_address: '',
      postal_code: ''
    }
    citiesList.value = []
    districtsList.value = []
    villagesList.value = []
  }
  showAddressForm.value = true
}

const onProvinceChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.city_id = ''
    editingAddress.value.district_id = ''
    editingAddress.value.village_id = ''
    citiesList.value = await addressVm.loadCities(editingAddress.value.province_id)
    districtsList.value = []
    villagesList.value = []
  }
}

const onCityChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.district_id = ''
    editingAddress.value.village_id = ''
    districtsList.value = await addressVm.loadDistricts(editingAddress.value.city_id)
    villagesList.value = []
  }
}

const onDistrictChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.village_id = ''
    villagesList.value = await addressVm.loadVillages(editingAddress.value.district_id)
  }
}

const handleSaveAddress = async () => {
  if (!editingAddress.value) return
  const addr = editingAddress.value
  if (!addr.receiver_name || !addr.phone || !addr.province_id || !addr.city_id || !addr.district_id || !addr.village_id || !addr.full_address || !addr.postal_code) {
    alert('Please fill in all the required fields.')
    return
  }
  const result = await addressVm.saveAddress(addr)
  if (result.success) {
    showAddressForm.value = false
    editingAddress.value = null
  } else {
    alert(result.message)
  }
}

const handleDeleteAddress = async (id: string) => {
  if (confirm('Are you sure you want to delete this address?')) {
    const result = await addressVm.deleteAddress(id)
    if (!result.success) {
      alert(result.message)
    }
  }
}

// INTEGRASI: Ambil fungsi orders dan helper formatRupiah murni dari useCart()
const { orders, formatRupiah } = useCart()

// Logout
const handleLogout = async () => {
  await authVm.logout()
}

// Filter orders
const filteredOrders = computed(() => {
  if (!orders || !orders.value) return []
  return orders.value.filter(order => {
    const status = order.status
    if (activeOrderStatus.value === 'pending') {
      return status === 'pending' || status === 'pembayaran'
    } else if (activeOrderStatus.value === 'processing') {
      return status === 'confirmed' || status === 'processing' || status === 'pengemasan'
    } else if (activeOrderStatus.value === 'shipping') {
      return status === 'shipped' || status === 'delivered' || status === 'pengiriman'
    } else if (activeOrderStatus.value === 'done') {
      return status === 'finished' || status === 'cancelled' || status === 'ulasan'
    }
    return false
  })
})

const selectedOrder = ref<any | null>(null)
const showShippingOverlay = ref(false)

const openOrderDetail = (order: any) => {
  selectedOrder.value = order
}

const closeOrderDetail = () => {
  selectedOrder.value = null
  showShippingOverlay.value = false
}

const toggleShippingOverlay = () => {
  showShippingOverlay.value = !showShippingOverlay.value
}

const contactDriver = (orderId: string) => {
  window.open(`https://wa.me/628175234999?text=Hello%20Chia%20Florist,%20I%20would%20like%20to%20inquire%20about%20the%20delivery%20status%20for%20order%20${orderId}`, '_blank')
}

const leaveReview = (orderId: string) => {
  window.alert(`Thank you! Review form for order ${orderId} will open shortly. Give us your best stars! ⭐`)
}

const triggerAlert = (message: string) => {
  window.alert(message)
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-12 font-sans">
    <div class="max-w-7xl mx-auto px-8">
      
      <nav class="flex text-sm text-gray-400 mb-8 gap-2">
        <NuxtLink to="/" class="hover:text-[#1b4332]">Home</NuxtLink>
        <span>/</span>
        <span class="text-gray-900 font-medium">My Profile</span>
      </nav>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        
        <aside class="lg:col-span-1 space-y-6">
          <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
            <div class="flex flex-col items-center text-center mb-8">
              <!-- Interactive Avatar Container -->
              <div 
                class="relative group cursor-pointer w-20 h-20 mb-4 rounded-full overflow-hidden border border-gray-200 shadow-sm flex items-center justify-center bg-gray-100" 
                @click="triggerFileSelect"
                title="Change Profile Picture"
              >
                <!-- Avatar Image -->
                <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105" alt="Avatar" />
                <!-- Default Icon -->
                <div v-else class="text-3xl text-[#1b4332]">👤</div>

                <!-- Hover Overlay -->
                <div class="absolute inset-0 bg-black/40 flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  <span class="text-white text-xs font-bold">Change</span>
                  <span class="text-white text-[10px]">Photo</span>
                </div>

                <!-- Uploading Spinner -->
                <div v-if="isUploading" class="absolute inset-0 bg-black/60 flex items-center justify-center">
                  <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                </div>
              </div>

              <!-- Hidden File Input -->
              <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileUpload" />

              <h2 class="font-bold text-gray-900 leading-tight">{{ user.name }}</h2>
              <p class="text-[10px] font-black text-emerald-600 uppercase tracking-widest mt-1">Silver Member</p>
            </div>
            <nav class="space-y-1">
              <button 
                v-for="tab in [{id:'personal', label:'Personal Information'}, {id:'addresses', label:'Shipping Addresses'}, {id:'orders', label:'Order Tracking'}]"
                :key="tab.id"
                @click="activeTab = tab.id"
                :class="['w-full text-left px-4 py-3 rounded-xl text-sm font-semibold transition-all', activeTab === tab.id ? 'bg-[#1b4332] text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50']"
              >
                {{ tab.label }}
              </button>
              
              <div class="pt-4 mt-4 border-t border-gray-100">
                <button @click="handleLogout" class="w-full text-left px-4 py-3 rounded-xl text-sm font-semibold text-red-500 hover:bg-red-50 transition-all flex items-center gap-2">
                  Logout
                </button>
              </div>
            </nav>
          </div>
        </aside>

        <main class="lg:col-span-3 space-y-6">
          
          <!-- Personal tab -->
          <div v-if="activeTab === 'personal'" class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden animate-fade">
            <div class="px-8 py-6 border-b border-gray-50">
              <h3 class="font-bold text-gray-900 text-lg">Edit Profile</h3>
              <p class="text-xs text-gray-400 mt-1">Manage your personal details and contact information.</p>
            </div>
            
            <div class="p-8 space-y-6">
              <!-- PROFILE PICTURE SECTION -->
              <div class="border-b border-gray-100 pb-6 space-y-4">
                <h4 class="text-xs font-bold text-gray-700 uppercase tracking-wider">Profile Picture</h4>
                <div class="flex flex-col sm:flex-row items-start sm:items-center gap-6">
                  <!-- Image Preview -->
                  <div class="w-24 h-24 rounded-full overflow-hidden border border-gray-200 bg-gray-50 flex items-center justify-center text-3xl text-gray-400 relative">
                    <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="Profile picture" />
                    <span v-else>👤</span>
                    <div v-if="isUploading" class="absolute inset-0 bg-black/50 flex items-center justify-center">
                      <div class="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
                    </div>
                  </div>

                  <div class="space-y-3 flex-1">
                    <div class="flex flex-wrap gap-2">
                      <button @click="triggerFileSelect" :disabled="isUploading" class="bg-[#1b4332] hover:bg-[#143326] text-white px-4 py-2 rounded-xl text-xs font-bold shadow-sm transition disabled:opacity-50 cursor-pointer">
                        Upload New Photo
                      </button>
                      <button v-if="avatarUrl" @click="handleRemovePicture" :disabled="isUploading" class="border border-red-200 text-red-600 hover:bg-red-50 px-4 py-2 rounded-xl text-xs font-bold transition disabled:opacity-50 cursor-pointer">
                        Remove Photo
                      </button>
                    </div>
                    <p class="text-[10px] text-gray-400 leading-normal">Allowed JPG, PNG or WEBP. Max size 5MB.</p>
                  </div>
                </div>


                
                <!-- Upload Error -->
                <div v-if="uploadError" class="bg-red-50 border border-red-100 rounded-xl p-3 text-red-600 text-xs font-semibold">
                  ⚠️ {{ uploadError }}
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Full Name</label>
                  <input type="text" v-model="user.name" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Username</label>
                  <input type="text" v-model="user.username" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Phone Number</label>
                  <input type="text" v-model="user.phone" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Email Address</label>
                  <input type="email" :value="user.email" disabled class="w-full px-4 py-3 bg-gray-100 border border-gray-200 rounded-xl text-sm font-semibold outline-none text-gray-500 cursor-not-allowed" />
                </div>
              </div>

              <div v-if="user.lastLoginAt" class="flex justify-end mt-2">
                <p class="text-[10px] text-gray-400 font-medium bg-gray-50 px-2 py-1 rounded-md border border-gray-100">Last Login: {{ new Date(user.lastLoginAt).toLocaleString() }}</p>
              </div>
              
              <div class="flex justify-end">
                <button @click="handleUpdateProfile" :disabled="authVm.isLoading.value" class="bg-[#1b4332] hover:bg-[#143326] text-white px-8 py-3 rounded-xl text-sm font-bold shadow-sm transition disabled:opacity-50 flex items-center gap-2 cursor-pointer">
                  <span v-if="authVm.isLoading.value" class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></span>
                  <span>Save Changes</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Shipping Addresses tab -->
          <div v-if="activeTab === 'addresses'" class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden animate-fade">
            <div class="px-8 py-6 border-b border-gray-50 flex justify-between items-center">
              <div>
                <h3 class="font-bold text-gray-900 text-lg">Shipping Addresses</h3>
                <p class="text-xs text-gray-400 mt-1">Manage your delivery addresses for secure checkout.</p>
              </div>
              <button @click="initAddressForm()" class="bg-[#1b4332] hover:bg-[#143326] text-white px-4 py-2.5 rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
                + Add New Address
              </button>
            </div>

            <div class="p-8 space-y-6">
              <!-- Loading -->
              <div v-if="addressVm.isLoading.value" class="flex flex-col items-center justify-center py-12 space-y-4">
                <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-[#1b4332]"></div>
                <p class="text-gray-500 text-xs font-medium">Loading addresses...</p>
              </div>

              <!-- Empty state -->
              <div v-else-if="addressVm.addresses.value.length === 0" class="text-center py-16 space-y-4">
                <div class="text-5xl">📍</div>
                <h4 class="font-bold text-gray-900 text-base">No Address Registered</h4>
                <p class="text-xs text-gray-400 max-w-sm mx-auto">Please add a shipping address to place orders and calculate courier fees.</p>
              </div>

              <!-- Address List -->
              <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div 
                  v-for="addr in addressVm.addresses.value" 
                  :key="addr.address_id" 
                  :class="['border rounded-2xl p-6 relative flex flex-col justify-between transition-all duration-300', addr.is_default ? 'border-emerald-500 bg-emerald-50/10' : 'border-gray-200 hover:border-gray-300 bg-white']"
                >
                  <div>
                    <div class="flex items-center gap-2 mb-3">
                      <h4 class="font-bold text-gray-900 text-sm">{{ addr.receiver_name }}</h4>
                      <span v-if="addr.is_default" class="bg-emerald-100 text-emerald-800 text-[10px] font-bold px-2 py-0.5 rounded-full border border-emerald-200">Default</span>
                    </div>
                    <p class="text-xs text-gray-600 font-semibold mb-2">📞 {{ addr.phone }}</p>
                    <p class="text-xs text-gray-500 leading-relaxed">{{ addr.full_address }}</p>
                    <p class="text-[11px] text-gray-400 mt-2 font-medium">Postal Code: {{ addr.postal_code }}</p>
                  </div>

                  <div class="mt-6 pt-4 border-t border-gray-100 flex justify-end gap-3">
                    <button @click="initAddressForm(addr)" class="text-xs font-bold text-gray-600 hover:text-black transition cursor-pointer">
                      Edit
                    </button>
                    <button 
                      v-if="!addr.is_default" 
                      @click="handleDeleteAddress(addr.address_id!)" 
                      class="text-xs font-bold text-red-500 hover:text-red-600 transition cursor-pointer"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Orders tab -->
          <div v-if="activeTab === 'orders'" class="space-y-6 animate-fade">
            <div class="bg-white border border-gray-100 p-2 rounded-2xl shadow-sm grid grid-cols-4 gap-1 text-center font-medium">
              <button 
                v-for="status in [{id:'pending', label:'Pending'}, {id:'processing', label:'Processing'}, {id:'shipping', label:'Shipping'}, {id:'done', label:'Done'}]"
                :key="status.id"
                @click="activeOrderStatus = status.id as any"
                :class="['py-3 text-xs sm:text-sm rounded-xl transition-all font-bold', activeOrderStatus === status.id ? 'bg-[#1b4332] text-white shadow-sm' : 'text-gray-400 hover:text-gray-900']"
              >
                {{ status.label }}
              </button>
            </div>

            <div v-if="filteredOrders.length > 0" class="space-y-4">
              <div v-for="order in filteredOrders" :key="order.orderId" class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-6 hover:shadow-md transition duration-300">
                <div class="space-y-3 flex-1">
                  <div class="flex flex-wrap items-center gap-3">
                    <span class="text-xs font-bold text-gray-400">🗓️ {{ order.date }}</span>
                    <span class="text-xs font-mono font-bold text-gray-900 bg-gray-50 px-2.5 py-1 rounded-md border border-gray-100">{{ order.orderId }}</span>
                    <div>
                      <span v-if="order.status === 'pending' || order.status === 'pembayaran'" class="px-2.5 py-0.5 bg-amber-50 text-amber-700 text-xs font-bold rounded-full border border-amber-100">
                        Pending Payment
                      </span>
                      <span v-else-if="order.status === 'confirmed'" class="px-2.5 py-0.5 bg-indigo-50 text-indigo-700 text-xs font-bold rounded-full border border-indigo-100">
                        Confirmed
                      </span>
                      <span v-else-if="order.status === 'processing' || order.status === 'pengemasan'" class="px-2.5 py-0.5 bg-emerald-50 text-emerald-700 text-xs font-bold rounded-full border border-emerald-100">
                        Arranging Flowers
                      </span>
                      <span v-else-if="order.status === 'shipped'" class="px-2.5 py-0.5 bg-blue-50 text-blue-700 text-xs font-bold rounded-full border border-blue-100">
                        Shipped
                      </span>
                      <span v-else-if="order.status === 'delivered' || order.status === 'pengiriman'" class="px-2.5 py-0.5 bg-sky-50 text-sky-700 text-xs font-bold rounded-full border border-sky-100">
                        Delivered
                      </span>
                      <span v-else-if="order.status === 'finished' || order.status === 'ulasan'" class="px-2.5 py-0.5 bg-gray-100 text-gray-700 text-xs font-bold rounded-full border border-gray-200">
                        Finished
                      </span>
                      <span v-else-if="order.status === 'cancelled'" class="px-2.5 py-0.5 bg-red-50 text-red-700 text-xs font-bold rounded-full border border-red-100">
                        Cancelled
                      </span>
                    </div>
                  </div>

                  <div class="flex items-center gap-4 py-1">
                    <div class="flex -space-x-3 overflow-hidden">
                      <img 
                        v-for="(item, idx) in order.items.slice(0, 3)" 
                        :key="idx"
                        :src="item.image || '/images/custom-preview.png'" 
                        class="inline-block h-10 w-10 rounded-lg ring-2 ring-white object-cover border border-gray-100 bg-gray-50"
                      />
                      <div v-if="order.items.length > 3" class="inline-flex items-center justify-center h-10 w-10 rounded-lg ring-2 ring-white bg-gray-100 border border-gray-100 text-xs font-bold text-gray-500">
                        +{{ order.items.length - 3 }}
                      </div>
                    </div>
                    <div class="min-w-0">
                      <p class="text-xs font-bold text-gray-900 truncate max-w-md">
                        {{ order.items[0]?.name }}{{ order.items.length > 1 ? ` and ${order.items.length - 1} other item(s)` : '' }}
                      </p>
                      <p class="text-[10px] text-gray-400 mt-0.5">Total Quantity: {{ order.items.reduce((acc, item) => acc + item.quantity, 0) }}</p>
                    </div>
                  </div>
                </div>

                <div class="flex items-center justify-between md:justify-end gap-6 border-t md:border-t-0 border-gray-50 pt-4 md:pt-0">
                  <div class="text-left md:text-right">
                    <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Total Bill</p>
                    <p class="text-base font-extrabold text-[#1b4332] mt-0.5">{{ formatRupiah(order.total) }}</p>
                  </div>
                  <button 
                    @click="openOrderDetail(order)" 
                    class="bg-[#1b4332] hover:bg-[#143326] text-white px-5 py-2.5 rounded-xl text-xs font-bold shadow-sm transition cursor-pointer"
                  >
                    View Details
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="bg-white border border-gray-100 rounded-2xl p-16 text-center shadow-sm">
              <div class="text-5xl mb-4">📑</div>
              <h4 class="font-bold text-gray-900 text-lg">No Orders Found</h4>
              <p class="text-sm text-gray-400 mt-1">There are no transactions in the "{{ activeOrderStatus }}" status tab yet.</p>
            </div>
          </div>

        </main>
      </div>
    </div>

    <!-- Modal Address Form -->
    <div v-if="showAddressForm" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl max-w-xl w-full max-h-[90vh] overflow-y-auto shadow-2xl border border-gray-100 flex flex-col justify-between">
        <div class="px-8 py-6 border-b border-gray-100 flex justify-between items-center">
          <h3 class="font-extrabold text-gray-900 text-base">
            {{ editingAddress?.address_id ? 'Edit Address' : 'Add New Address' }}
          </h3>
          <button @click="showAddressForm = false" class="text-gray-400 hover:text-black font-bold text-lg cursor-pointer">✕</button>
        </div>

        <div class="p-8 space-y-4 flex-1">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Receiver Name *</label>
              <input v-model="editingAddress!.receiver_name" type="text" placeholder="e.g. John Doe" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Phone Number *</label>
              <input v-model="editingAddress!.phone" type="text" placeholder="e.g. 0812345678" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Province *</label>
              <select v-model="editingAddress!.province_id" @change="onProvinceChange" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition">
                <option value="">Select Province</option>
                <option v-for="p in provincesList" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">City *</label>
              <select v-model="editingAddress!.city_id" @change="onCityChange" :disabled="!editingAddress!.province_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select City</option>
                <option v-for="c in citiesList" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">District *</label>
              <select v-model="editingAddress!.district_id" @change="onDistrictChange" :disabled="!editingAddress!.city_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select District</option>
                <option v-for="d in districtsList" :key="d.id" :value="d.id">{{ d.name }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Village *</label>
              <select v-model="editingAddress!.village_id" :disabled="!editingAddress!.district_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select Village</option>
                <option v-for="v in villagesList" :key="v.id" :value="v.id">{{ v.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Postal Code *</label>
              <input v-model="editingAddress!.postal_code" type="text" placeholder="e.g. 17131" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
            <div class="flex items-center pt-6 gap-2">
              <input v-model="editingAddress!.is_default" type="checkbox" id="default-check" class="h-4 w-4 accent-[#1b4332] rounded focus:ring-0 cursor-pointer" />
              <label for="default-check" class="text-xs font-bold text-gray-700 cursor-pointer">Set as default address</label>
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-xs font-bold text-gray-500">Complete Address *</label>
            <textarea v-model="editingAddress!.full_address" rows="3" placeholder="Street name, house number, details..." class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition resize-none"></textarea>
          </div>
        </div>

        <div class="px-8 py-6 border-t border-gray-100 flex justify-end gap-3">
          <button @click="showAddressForm = false" class="px-5 py-2.5 border border-gray-200 rounded-xl text-xs font-bold text-gray-600 hover:bg-gray-50 transition cursor-pointer">
            Cancel
          </button>
          <button @click="handleSaveAddress" class="px-5 py-2.5 bg-[#1b4332] hover:bg-[#143326] text-white rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
            Save Address
          </button>
        </div>
      </div>
    </div>

    <!-- Independent Order Detail Modal Overlay -->
    <div v-if="selectedOrder" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl max-w-2xl w-full max-h-[85vh] overflow-hidden shadow-2xl border border-gray-100 flex flex-col justify-between relative animate-fade">
        
        <!-- Header -->
        <div class="px-8 py-5 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
          <div>
            <div class="flex items-center gap-3">
              <h3 class="font-extrabold text-gray-900 text-base">Order Details</h3>
              <span class="text-xs font-mono font-bold text-gray-500 bg-gray-100 px-2 py-0.5 rounded border border-gray-200">{{ selectedOrder.orderId }}</span>
            </div>
            <p class="text-[11px] text-gray-400 mt-0.5">Placed on {{ selectedOrder.date }}</p>
          </div>
          <button @click="closeOrderDetail" class="text-gray-400 hover:text-black font-bold text-lg cursor-pointer">✕</button>
        </div>

        <!-- Body -->
        <div class="p-8 space-y-6 flex-1 overflow-y-auto custom-scrollbar">
          <!-- Status block -->
          <div class="flex items-center justify-between bg-emerald-50/10 border border-emerald-100 p-4 rounded-2xl">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Current Status</p>
              <div class="mt-1 flex items-center gap-2">
                <span v-if="selectedOrder.status === 'pending' || selectedOrder.status === 'pembayaran'" class="px-2.5 py-0.5 bg-amber-50 text-amber-700 text-xs font-bold rounded-full border border-amber-100">
                  Pending Payment
                </span>
                <span v-else-if="selectedOrder.status === 'confirmed'" class="px-2.5 py-0.5 bg-indigo-50 text-indigo-700 text-xs font-bold rounded-full border border-indigo-100">
                  Confirmed
                </span>
                <span v-else-if="selectedOrder.status === 'processing' || selectedOrder.status === 'pengemasan'" class="px-2.5 py-0.5 bg-emerald-50 text-emerald-700 text-xs font-bold rounded-full border border-emerald-100">
                  Arranging Flowers
                </span>
                <span v-else-if="selectedOrder.status === 'shipped'" class="px-2.5 py-0.5 bg-blue-50 text-blue-700 text-xs font-bold rounded-full border border-blue-100">
                  Shipped
                </span>
                <span v-else-if="selectedOrder.status === 'delivered' || selectedOrder.status === 'pengiriman'" class="px-2.5 py-0.5 bg-sky-50 text-sky-700 text-xs font-bold rounded-full border border-sky-100">
                  Delivered
                </span>
                <span v-else-if="selectedOrder.status === 'finished' || selectedOrder.status === 'ulasan'" class="px-2.5 py-0.5 bg-gray-100 text-gray-700 text-xs font-bold rounded-full border border-gray-200">
                  Finished
                </span>
                <span v-else-if="selectedOrder.status === 'cancelled'" class="px-2.5 py-0.5 bg-red-50 text-red-700 text-xs font-bold rounded-full border border-red-100">
                  Cancelled
                </span>
              </div>
            </div>

            <!-- Shipping Detail Button -->
            <button 
              @click="toggleShippingOverlay" 
              class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-xl text-xs font-bold transition flex items-center gap-1.5 cursor-pointer"
            >
              <span>🚚</span>
              <span>Track Shipping</span>
            </button>
          </div>

          <!-- Items list -->
          <div class="space-y-4">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider">Ordered Items</h4>
            <div class="divide-y divide-gray-100 border border-gray-100 rounded-2xl overflow-hidden bg-gray-50/20">
              <div v-for="(item, idx) in selectedOrder.items" :key="idx" class="flex gap-4 items-center p-4 bg-white">
                <div class="w-14 h-14 bg-gray-50 rounded-xl overflow-hidden border border-gray-100 flex-shrink-0">
                  <img :src="item.image || '/images/custom-preview.png'" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h5 class="font-bold text-gray-900 text-xs truncate leading-snug">{{ item.name }}</h5>
                  <div class="flex gap-2 mt-1.5 text-[10px] text-gray-500 font-semibold">
                    <span v-if="item.size" class="bg-gray-100 px-2 py-0.5 rounded">Size: {{ item.size }}</span>
                    <span v-if="item.color" class="flex items-center gap-1">
                      Color:
                      <span :style="{ backgroundColor: item.color }" class="w-2.5 h-2.5 rounded-full border border-gray-300 inline-block"></span>
                    </span>
                    <span v-if="item.isCustom" class="bg-emerald-50 text-emerald-700 px-2 py-0.5 rounded font-black">Custom Design</span>
                  </div>
                </div>
                <div class="text-right text-xs">
                  <p class="font-bold text-gray-900">{{ formatRupiah(item.price) }}</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Qty: {{ item.quantity }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Billing Info -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider">Billing Summary</h4>
            <div class="border border-gray-100 rounded-2xl p-4 bg-gray-50/30 text-xs space-y-3 font-semibold text-gray-600">
              <div class="flex justify-between">
                <span>Items Subtotal</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(selectedOrder.total - (selectedOrder.shipping?.cost || 20000)) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Shipping Cost</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(selectedOrder.shipping?.cost || 20000) }}</span>
              </div>
              <div class="border-t border-gray-100 pt-3 flex justify-between text-sm font-bold text-gray-900">
                <span>Total Amount Paid</span>
                <span class="text-base text-[#1b4332] font-black">{{ formatRupiah(selectedOrder.total) }}</span>
              </div>
            </div>
          </div>

        </div>

        <!-- Footer Actions -->
        <div class="px-8 py-5 border-t border-gray-100 flex justify-between items-center bg-gray-50/50">
          <button @click="closeOrderDetail" class="px-5 py-2.5 border border-gray-200 rounded-xl text-xs font-bold text-gray-600 hover:bg-gray-50 transition cursor-pointer">
            Close
          </button>

          <div>
            <div v-if="selectedOrder.status === 'pending' || selectedOrder.status === 'pembayaran'">
              <button @click="triggerAlert('Our team is reviewing your transaction recipe. Please wait max 15 minutes.')" class="px-5 py-2.5 bg-gray-900 hover:bg-black text-white text-xs font-bold rounded-xl transition cursor-pointer">
                Check Bank Receipt
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'confirmed' || selectedOrder.status === 'processing' || selectedOrder.status === 'pengemasan'">
              <button @click="triggerAlert('Your flower board is being carefully arranged by our expert florists.')" class="px-5 py-2.5 bg-[#1b4332] text-white hover:bg-[#143326] text-xs font-bold rounded-xl transition cursor-pointer">
                View Workshop Progress
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'shipped' || selectedOrder.status === 'delivered' || selectedOrder.status === 'pengiriman'">
              <button @click="contactDriver(selectedOrder.orderId)" class="px-5 py-2.5 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-xl transition flex items-center gap-2 cursor-pointer">
                Contact Courier (WhatsApp)
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'finished' || selectedOrder.status === 'ulasan'">
              <button @click="leaveReview(selectedOrder.orderId)" class="px-5 py-2.5 bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold rounded-xl transition cursor-pointer">
                Write Product Review
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'cancelled'">
              <button @click="triggerAlert('This order was cancelled. You can rebuild your cart to re-order.')" class="px-5 py-2.5 bg-red-50 text-red-700 border border-red-100 hover:bg-red-100 text-xs font-bold rounded-xl transition cursor-pointer">
                Order Again
              </button>
            </div>
          </div>
        </div>

        <!-- Shipping Detail Slide-out Overlay Panel (Nested) -->
        <div 
          v-if="showShippingOverlay" 
          class="absolute inset-y-0 right-0 w-full sm:w-[450px] bg-white shadow-2xl border-l border-gray-100 flex flex-col justify-between z-50 transform transition-transform duration-300 ease-in-out animate-slide-in"
        >
          <!-- Overlay Header -->
          <div class="px-6 py-5 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
            <div class="flex items-center gap-2">
              <span class="text-base">🚚</span>
              <h3 class="font-extrabold text-gray-900 text-sm">Shipping Information</h3>
            </div>
            <button @click="showShippingOverlay = false" class="text-gray-400 hover:text-black font-bold text-base cursor-pointer">✕</button>
          </div>

          <!-- Overlay Body -->
          <div class="p-6 space-y-6 flex-1 overflow-y-auto custom-scrollbar text-xs">
            <!-- Courier Summary -->
            <div class="border border-gray-100 rounded-2xl p-4 bg-gray-50/30 space-y-2">
              <div class="flex justify-between font-semibold">
                <span class="text-gray-400">Courier Partner</span>
                <span class="text-gray-900 font-bold">{{ selectedOrder.shipping?.courier || 'JNE Express' }} ({{ selectedOrder.shipping?.service || 'REG' }})</span>
              </div>
              <div class="flex justify-between font-semibold">
                <span class="text-gray-400">Tracking Code</span>
                <span class="text-gray-900 font-mono font-bold">{{ selectedOrder.shipping?.trackingNumber || 'AWB-8839214751' }}</span>
              </div>
            </div>

            <!-- Destination Address -->
            <div class="space-y-2">
              <h4 class="text-[10px] font-bold text-gray-500 uppercase tracking-wider">Shipping Destination</h4>
              <div class="border border-gray-100 rounded-2xl p-4 space-y-2 font-semibold">
                <div class="flex justify-between">
                  <span class="text-gray-400">Recipient</span>
                  <span class="text-gray-900">{{ selectedOrder.shipping?.recipientName || 'Jane Doe' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-400">Contact Number</span>
                  <span class="text-gray-900">{{ selectedOrder.shipping?.phone || '081234567890' }}</span>
                </div>
                <div class="border-t border-gray-50 pt-2 mt-2">
                  <p class="text-gray-400 mb-1">Full Address</p>
                  <p class="text-gray-700 leading-normal font-medium">{{ selectedOrder.shipping?.address || 'Jl. Merdeka No. 45, Kebayoran Baru, Jakarta Selatan, 12110' }}</p>
                </div>
              </div>
            </div>

            <!-- Shipping Timeline -->
            <div class="space-y-4">
              <h4 class="text-[10px] font-bold text-gray-500 uppercase tracking-wider">Tracking Timeline</h4>
              
              <div class="relative pl-6 space-y-6 before:content-[''] before:absolute before:left-2 before:top-2 before:bottom-2 before:w-[2px] before:bg-gray-100">
                <!-- Milestone: Placed -->
                <div class="relative">
                  <span 
                    class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white flex items-center justify-center"
                    :class="[selectedOrder.status !== 'cancelled' ? 'bg-emerald-500 ring-4 ring-emerald-100' : 'bg-red-500 ring-4 ring-red-100']"
                  ></span>
                  <p class="font-bold text-gray-900">Order Placed Successfully</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Mock Date: 2026-06-26 10:00 AM</p>
                </div>

                <!-- Milestone: Payment (for non-pending) -->
                <div v-if="selectedOrder.status !== 'pending' && selectedOrder.status !== 'pembayaran' && selectedOrder.status !== 'cancelled'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Payment Verified by Admin</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Mock Date: 2026-06-26 10:15 AM</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Payment Verification</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Awaiting bank receipt upload / review</p>
                </div>

                <!-- Milestone: Arranging Flowers (for processing or later) -->
                <div v-if="['processing', 'shipped', 'delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Flower Arrangement Complete</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Expert florist team finished production</p>
                </div>
                <div v-else-if="selectedOrder.status === 'confirmed'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-blue-500 ring-4 ring-blue-100"></span>
                  <p class="font-bold text-gray-900">Queued for Arrangement</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Waiting for florist slot allocation</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Flower Board Production</p>
                </div>

                <!-- Milestone: Shipped (for shipped or later) -->
                <div v-if="['shipped', 'delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Handed Over to Courier</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Parcel picked up by {{ selectedOrder.shipping?.courier || 'JNE' }} partner</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Hand Over to Courier Partner</p>
                </div>

                <!-- Milestone: Delivered (for delivered or later) -->
                <div v-if="['delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Delivered to Destination</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Signed by recipient at destination address</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Out for Delivery / Completed Destination</p>
                </div>

                <!-- Milestone: Completed / Finished -->
                <div v-if="selectedOrder.status === 'finished'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Order Completed</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Transaction finalized by customer</p>
                </div>

                <!-- Milestone: Cancelled -->
                <div v-if="selectedOrder.status === 'cancelled'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-red-500 ring-4 ring-red-100"></span>
                  <p class="font-bold text-red-600">Order Cancelled</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">This transaction has been terminated</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Overlay Footer -->
          <div class="px-6 py-5 border-t border-gray-100 flex justify-end bg-gray-50/50">
            <button @click="showShippingOverlay = false" class="px-5 py-2.5 bg-[#1b4332] text-white hover:bg-[#143326] rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
              Go Back to Details
            </button>
          </div>
        </div>



      </div>
    </div>

  </div>
</template>

<style scoped>
.animate-fade { animation: fadeIn 0.4s ease-out; }
.animate-slide-in { animation: slideIn 0.3s ease-out forwards; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes slideIn { from { transform: translateX(100%); } to { transform: translateX(0); } }
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
</style>