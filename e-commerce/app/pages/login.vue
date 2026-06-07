<script setup lang="ts">
import { ref } from 'vue'

useHead({ title: 'Login - Chia Florist' })

const nameOrEmail = ref('')
const password = ref('')
const isSubmitting = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  if (!nameOrEmail.value || !password.value) {
    errorMessage.value = 'Please fill in all fields.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  // =========================================================
  // 🛠️ TRIK BYPASS INSTAN: Langsung lempar ke halaman OTP
  // =========================================================
  localStorage.setItem('register_email', nameOrEmail.value)
  navigateTo('/verify')
  return 
  // =========================================================
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-10 min-h-[80vh] flex items-center">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-16 items-center w-full">
      
      <div class="hidden md:block h-[600px]">
        <img 
          src="/images/florist.jpg"
          alt="Plants Decor" 
          class="w-full h-full object-cover rounded-xl shadow-sm"
        />
      </div>

      <div class="max-w-md w-full mx-auto space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Welcome to Chia Florist</h1>
          <p class="text-gray-600">Enter your details below</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          
          <div v-if="errorMessage" class="bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl">
            ⚠️ {{ errorMessage }}
          </div>

          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              v-model="nameOrEmail"
              placeholder="Email" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="isSubmitting"
            />
          </div>

          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="password" 
              v-model="password"
              placeholder="Password" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="isSubmitting"
            />
          </div>

          <div class="pt-4 space-y-4">
            <button
              type="submit"
              class="w-full bg-[#1b4332] text-white py-4 rounded-md font-medium hover:bg-[#143326] disabled:bg-gray-300 transition-all shadow-md cursor-pointer flex items-center justify-center"
              :disabled="isSubmitting"
            >
              <span>Login</span>
            </button>

            <button 
              type="button" 
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
            >
              <img src="/images/google.png" class="w-5 h-5" alt="Google Icon" />
              Login with Google
            </button>
          </div>
        </form>

        <div class="text-center pt-4 text-gray-600">
          Not registered yet? 
          <NuxtLink to="/register" class="font-medium text-black border-b border-gray-500 pb-0.5 ml-2 hover:text-[#1b4332] transition-colors">
            Sign in
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>