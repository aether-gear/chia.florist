<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCart } from '~/composables/useCart'

useHead({ title: 'My Profile - Chia Florist' })

const activeTab = ref('personal')
// State untuk menyaring 4 status pesanan
const activeOrderStatus = ref<'pembayaran' | 'pengemasan' | 'pengiriman' | 'ulasan'>('pembayaran')

// Memisahkan data user secara reaktif yang bersih
const user = ref({
  name: 'Rayhan Shidqi',
  email: 'rayhan.shidqi@email.com',
  phone: '+62 812 3456 7890',
  address: 'Jl. Melati No. 12, Jakarta Timur'
})

// Ambil data orders global dari composable (Aman dari error destructuring)
const { orders } = useCart()

// Menyediakan fungsi logout lokal agar tidak memicu error compiler jika tidak ada di composable
const handleLogout = () => {
  alert('Logging out...')
  navigateTo('/')
}

// Memfilter pesanan yang ditampilkan berdasarkan sub-tab status yang aktif
const filteredOrders = computed(() => {
  if (!orders || !orders.value) return []
  return orders.value.filter(order => order.status === activeOrderStatus.value)
})
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
                v-for="tab in [{id:'personal', label:'Personal Information'}, {id:'orders', label:'Order Tracking'}]"
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
                <label class="text-xs font-bold text-gray-500">Shipping Address</label>
                <textarea rows="3" v-model="user.address" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold resize-none outline-none focus:bg-white focus:border-[#1b4332] transition-all"></textarea>
              </div>
              
              <div class="flex justify-end">
                <button class="bg-[#1b4332] hover:bg-[#143326] text-white px-8 py-3 rounded-xl text-sm font-bold shadow-sm transition">Save Changes</button>
              </div>
            </div>
          </div>

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
              <div v-for="order in filteredOrders" :key="order.orderId" class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm space-y-4">
                <div class="flex justify-between items-center border-b border-gray-50 pb-3 text-xs sm:text-sm font-bold text-gray-500">
                  <span>🗓️ {{ order.date }}</span>
                  <span class="text-mono text-gray-900">{{ order.orderId }}</span>
                </div>

                <div v-for="(item, idx) in order.items" :key="idx" class="flex gap-4 items-center py-2">
                  <div class="w-14 h-14 bg-gray-50 rounded-xl overflow-hidden border border-gray-100 flex-shrink-0">
                    <img :src="item.image || '/images/custom-preview.png'" class="w-full h-full object-cover" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <h4 class="font-bold text-gray-900 text-sm truncate leading-snug">{{ item.name }}</h4>
                    <p class="text-xs text-gray-400 mt-1">Qty: {{ item.quantity }}x <span v-if="item.size">({{ item.size }})</span></p>
                  </div>
                  <div class="text-sm font-bold text-gray-900">${{ (item.price * item.quantity).toFixed(2) }}</div>
                </div>

                <div class="border-t border-gray-50 pt-4 flex justify-between items-center text-sm font-bold">
                  <span class="text-gray-500">Total Paid:</span>
                  <span class="text-xl text-[#1b4332]">${{ order.total.toFixed(2) }}</span>
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
  </div>
</template>

<style scoped>
.animate-fade { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>