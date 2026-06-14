<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useCart } from '~/composables/useCart'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useAddress } from '~/composables/useAddress'
import type { UserAddress } from '~/types/address'

useHead({ title: 'My Profile - Chia Florist' })

const authVm = useAuthViewModel()
const addressVm = useAddress()

const activeTab = ref('personal')
const activeOrderStatus = ref<'pembayaran' | 'pengemasan' | 'pengiriman' | 'ulasan'>('pembayaran')

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
  email: 'Loading...',
  phone: 'Loading...',
  address: 'Jl. Melati No. 12, Jakarta Timur'
})

// Sync profile data from ViewModel
watch(() => authVm.currentUser.value, (me) => {
  if (me) {
    user.value.name = me.name || ''
    user.value.email = me.email || ''
    user.value.phone = me.phone || ''
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

// Global orders from useCart
const { orders } = useCart()

// Logout
const handleLogout = async () => {
  await authVm.logout()
}

// Filter orders
const filteredOrders = computed(() => {
  if (!orders || !orders.value) return []
  return orders.value.filter(order => order.status === activeOrderStatus.value)
})

const contactDriver = (orderId: string) => {
  window.open(`https://wa.me/628175234999?text=Halo%20Chia%20Florist,%20saya%20ingin%20menanyakan%20status%20kurir%20untuk%20pesanan%20${orderId}`, '_blank')
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
              <div class="w-20 h-20 bg-gray-100 rounded-full mb-4 flex items-center justify-center text-[#1b4332] border border-gray-200 text-2xl">👤</div>
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
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Full Name</label>
                  <input type="text" v-model="user.name" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Email Address</label>
                  <input type="email" v-model="user.email" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Phone Number</label>
                  <input type="text" v-model="user.phone" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
              </div>
              
              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Shipping Address (Legacy)</label>
                <textarea rows="3" v-model="user.address" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold resize-none outline-none focus:bg-white focus:border-[#1b4332] transition-all"></textarea>
              </div>
              
              <div class="flex justify-end">
                <button class="bg-[#1b4332] hover:bg-[#143326] text-white px-8 py-3 rounded-xl text-sm font-bold shadow-sm transition">Save Changes</button>
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
                v-for="status in [{id:'pembayaran', label:'Pembayaran'}, {id:'pengemasan', label:'Pengemasan'}, {id:'pengiriman', label:'Pengiriman'}, {id:'ulasan', label:'Ulasan'}]"
                :key="status.id"
                @click="activeOrderStatus = status.id as any"
                :class="['py-3 text-xs sm:text-sm rounded-xl transition-all font-bold', activeOrderStatus === status.id ? 'bg-[#1b4332] text-white shadow-sm' : 'text-gray-400 hover:text-gray-900']"
              >
                {{ status.label }}
              </button>
            </div>

            <div v-if="filteredOrders.length > 0" class="space-y-4">
              <div v-for="order in filteredOrders" :key="order.orderId" class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm space-y-5">
                
                <div class="flex flex-col sm:flex-row sm:items-center justify-between border-b border-gray-50 pb-4 gap-2">
                  <div class="flex items-center gap-4 text-xs sm:text-sm font-bold text-gray-500">
                    <span>🗓️ {{ order.date }}</span>
                    <span class="text-mono text-gray-900 bg-gray-50 px-2.5 py-1 rounded-md border border-gray-100">{{ order.orderId }}</span>
                  </div>
                  
                  <div>
                    <span v-if="order.status === 'pembayaran'" class="px-3 py-1 bg-amber-50 text-amber-700 text-xs font-bold rounded-full border border-amber-100">
                      Menunggu Verifikasi Pembayaran
                    </span>
                    <span v-else-if="order.status === 'pengemasan'" class="px-3 py-1 bg-emerald-50 text-emerald-700 text-xs font-bold rounded-full border border-emerald-100">
                      Papan Bunga Sedang Dirangkai
                    </span>
                    <span v-else-if="order.status === 'pengiriman'" class="px-3 py-1 bg-blue-50 text-blue-700 text-xs font-bold rounded-full border border-blue-100">
                      Dalam Perjalanan Menuju Gedung
                    </span>
                    <span v-else-if="order.status === 'ulasan'" class="px-3 py-1 bg-gray-100 text-gray-700 text-xs font-bold rounded-full border border-gray-200">
                      Pesanan Selesai Diterima
                    </span>
                  </div>
                </div>

                <div v-for="(item, idx) in order.items" :key="idx" class="flex gap-4 items-center py-2">
                  <div class="w-16 h-16 bg-gray-50 rounded-xl overflow-hidden border border-gray-100 flex-shrink-0">
                    <img :src="item.image || '/images/custom-preview.png'" class="w-full h-full object-cover" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <h4 class="font-bold text-gray-900 text-sm truncate leading-snug">{{ item.name }}</h4>
                    <div class="flex gap-2 mt-1.5 text-[11px] text-gray-500 font-semibold">
                      <span v-if="item.size" class="bg-gray-100 px-2 py-0.5 rounded">Size: {{ item.size }}</span>
                      <span v-if="item.isCustom" class="bg-emerald-50 text-emerald-700 px-2 py-0.5 rounded font-black">Custom UI Design</span>
                    </div>
                  </div>
                  <div class="text-sm font-bold text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</div>
                </div>

                <div class="border-t border-gray-50 pt-4 flex flex-col sm:flex-row justify-between items-stretch sm:items-center gap-4">
                  <div class="flex items-center gap-2 text-sm font-bold">
                    <span class="text-gray-400">Total Paid:</span>
                    <span class="text-xl text-[#1b4332]">${{ order.total.toFixed(2) }}</span>
                  </div>

                  <div class="flex justify-end">
                    <div v-if="order.status === 'pembayaran'">
                      <button @click="triggerAlert('Our team is reviewing your transaction recipe. Please wait max 15 minutes.')" class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-xs font-bold rounded-xl transition">
                        Cek Bukti Transfer
                      </button>
                    </div>
                    <div v-else-if="order.status === 'pengemasan'">
                      <button @click="triggerAlert('Papan bunga kamu sedang dikerjakan secara teliti oleh florist ahli kami.')" class="px-4 py-2 bg-emerald-50 text-emerald-800 border border-emerald-100 text-xs font-bold rounded-xl transition">
                        Lihat Progres Workshop
                      </button>
                    </div>
                    <div v-else-if="order.status === 'pengiriman'">
                      <button @click="contactDriver(order.orderId)" class="px-4 py-2 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-xl transition flex items-center gap-2">
                        Hubungi Kurir (WA)
                      </button>
                    </div>
                    <div v-else-if="order.status === 'ulasan'">
                      <button @click="leaveReview(order.orderId)" class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold rounded-xl transition">
                        Tulis Testimoni Produk
                      </button>
                    </div>
                  </div>
                </div>

              </div>
            </div>

            <div v-else class="bg-white border border-gray-100 rounded-2xl p-16 text-center shadow-sm">
              <div class="text-5xl mb-4">📑</div>
              <h4 class="font-bold text-gray-900 text-lg">Tidak Ada Pesanan</h4>
              <p class="text-sm text-gray-400 mt-1">Belum ada transaksi di tab status "{{ activeOrderStatus }}" ini.</p>
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

  </div>
</template>

<style scoped>
.animate-fade { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>