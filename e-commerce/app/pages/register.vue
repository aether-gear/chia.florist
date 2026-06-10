<script setup lang="ts">
import { ref } from 'vue'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
useHead({ title: 'Create Account - Chia Florist' })
const authVm = useAuthViewModel()
const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const phone = ref('')
const errorMessage = ref('')
const handleRegister = async () => {
  if (!name.value || !username.value || !email.value || !password.value) {
    errorMessage.value = 'Name, Username, Email, and Password are required.'
    return
  }
  errorMessage.value = ''
  try {
    const response = await authVm.register({
      name: name.value,
      username: username.value,
      email: email.value,
      password: password.value,
      phone: phone.value || undefined
    })
    if (response && response.challenge_id) {
      // Navigate to unified verify view
      navigateTo('/login?verify=true')
    }
  } catch (err: any) {
    errorMessage.value = err.data?.message || 'Registration failed. Please check your input details.'
  }
}
</script>
<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-10 min-h-[80vh] flex items-center">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-16 items-center w-full">
      
      <!-- Left image section -->
      <div class="hidden md:block h-[600px]">
        <img 
          src="/images/florist.jpg"
          alt="Plants Decor" 
          class="w-full h-full object-cover rounded-xl shadow-sm"
        />
      </div>
      <!-- Right registration form -->
      <div class="max-w-md w-full mx-auto space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Create an account</h1>
          <p class="text-gray-600">Enter your details below</p>
        </div>
        <form @submit.prevent="handleRegister" class="space-y-6">
          
          <div v-if="errorMessage" class="bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2">
            <span>⚠️</span>
            <p class="flex-1 leading-normal">{{ errorMessage }}</p>
          </div>
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              v-model="name"
              placeholder="Name" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
              required
            />
          </div>
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              v-model="username"
              placeholder="Username" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
              required
            />
          </div>
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="email" 
              v-model="email"
              placeholder="Email Address" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
              required
            />
          </div>
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="password" 
              v-model="password"
              placeholder="Password" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
              required
            />
          </div>
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="tel" 
              v-model="phone"
              placeholder="Phone Number (Optional)" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
            />
          </div>
          <div class="pt-4 space-y-4">
            <button
              type="submit"
              class="w-full bg-[#1b4332] text-white py-4 rounded-md font-medium hover:bg-[#143326] disabled:bg-gray-300 transition-all shadow-md cursor-pointer flex items-center justify-center"
              :disabled="authVm.isLoading.value"
            >
              <span v-if="authVm.isLoading.value" class="animate-pulse">Creating Account...</span>
              <span v-else>Create Account</span>
            </button>
            <button 
              type="button" 
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
            >
              <img src="/images/google.png" class="w-5 h-5" alt="Google Icon" />
              Sign up with Google
            </button>
          </div>
        </form>
        <div class="text-center pt-4 text-gray-600">
          Already have account? 
          <NuxtLink to="/login" class="font-medium text-black border-b border-gray-500 pb-0.5 ml-2 hover:text-[#1b4332] transition-colors">
            Log in
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>