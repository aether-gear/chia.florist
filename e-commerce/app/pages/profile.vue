<script setup lang="ts">
import { ref } from 'vue'

useHead({ title: 'My Profile - Chia Florist' })

// State Menu
const activeTab = ref('personal')

// Data Dummy
const user = ref({
  name: 'Rayhan Shidqi',
  email: 'rayhanshidqi07@gmail.com',
  phone: '+62 8953 2620 4046',
  address: 'Jl. Melati No. 12, Jakarta Timur',
  memberSince: '12 Mei 2024'
})

const logout = () => {
  // Logic logout di sini
  navigateTo('/login')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-12 font-sans">
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
              <div class="w-20 h-20 bg-gray-100 rounded-full mb-4 flex items-center justify-center text-[#1b4332] border border-gray-200">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h2 class="font-bold text-gray-900">{{ user.name }}</h2>
              <p class="text-xs text-gray-400 uppercase tracking-widest mt-1">Silver Member</p>
            </div>

            <nav class="space-y-1">
              <button 
                v-for="tab in [{id:'personal', label:'Personal Information'}, {id:'orders', label:'Order History'}, {id:'settings', label:'Account Settings'}]"
                :key="tab.id"
                @click="activeTab = tab.id"
                :class="[
                  'w-full text-left px-4 py-3 rounded-xl text-sm font-medium transition-all',
                  activeTab === tab.id ? 'bg-[#1b4332] text-white shadow-md' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900'
                ]"
              >
                {{ tab.label }}
              </button>
              
              <div class="pt-4 mt-4 border-t border-gray-100">
                <button @click="logout" class="w-full text-left px-4 py-3 rounded-xl text-sm font-medium text-red-500 hover:bg-red-50 transition-all flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
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
                <div v-for="(val, label) in { 'Full Name': user.name, 'Email Address': user.email, 'Phone Number': user.phone }" :key="label" class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">{{ label }}</label>
                  <input type="text" :value="val" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm transition-all" />
                </div>
              </div>

              <div class="space-y-1.5">
                <label class="text-xs font-bold text-gray-500">Shipping Address</label>
                <textarea rows="3" :value="user.address" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-sm transition-all resize-none"></textarea>
              </div>

              <div class="flex justify-end pt-4">
                <button class="bg-[#1b4332] text-white px-8 py-3 rounded-xl text-sm font-bold hover:bg-[#143326] transition-all shadow-sm">
                  Save Changes
                </button>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'personal'" class="grid grid-cols-1 sm:grid-cols-3 gap-6">
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 text-center">
              <span class="text-[10px] font-black uppercase text-gray-400 tracking-widest">Total Orders</span>
              <p class="text-3xl font-bold text-[#1b4332] mt-2">12</p>
            </div>
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 text-center">
              <span class="text-[10px] font-black uppercase text-gray-400 tracking-widest">Rewards Points</span>
              <p class="text-3xl font-bold text-[#1b4332] mt-2">4,250</p>
            </div>
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 text-center">
              <span class="text-[10px] font-black uppercase text-gray-400 tracking-widest">Membership</span>
              <p class="text-3xl font-bold text-[#1b4332] mt-2">Silver</p>
            </div>
          </div>

          <div v-if="activeTab === 'orders'" class="bg-white rounded-2xl p-12 shadow-sm border border-gray-100 text-center animate-fade">
             <div class="w-20 h-20 bg-gray-50 rounded-full mx-auto flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 11V7a4 4 0 00-8 0v4M5 9h12l1 12H4L5 9z" />
                </svg>
             </div>
             <h3 class="font-bold text-gray-900">No Orders Yet</h3>
             <p class="text-sm text-gray-400 mt-2 max-w-xs mx-auto">When you buy flowers, your order history will appear here.</p>
             <NuxtLink to="/products" class="inline-block mt-6 px-6 py-3 bg-[#1b4332] text-white text-sm font-bold rounded-xl">Start Shopping</NuxtLink>
          </div>

        </main>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.4s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>