<script setup lang="ts">
import { ref } from 'vue'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

useHead({ 
  title: 'Sign In - Chia Florist',
  meta: [
    { name: 'description', content: 'Sign in to access your Chia Florist profile.' }
  ]
})

const authVm = useAuthViewModel()

const email = ref('')
const password = ref('')
const errorMessage = ref('')

const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = 'Please fill in all fields.'
    return
  }

  errorMessage.value = ''
  
  try {
    const success = await authVm.login({
      email: email.value,
      password: password.value
    })
    
    if (success) {
      navigateTo('/')
    }
  } catch (err: any) {
    // Handle 403 unverified case
    if (err.status === 403 && err.data?.message === 'email not verified') {
      // Store current email to display in verify page
      if (import.meta.client) {
        localStorage.setItem('register_email', email.value)
      }
      errorMessage.value = 'Email not verified. Redirecting to verification...'
      setTimeout(() => {
        errorMessage.value = ''
        navigateTo('/register?verify=true')
      }, 1200)
    } else {
      errorMessage.value = err.data?.message || 'Login failed. Please check your credentials.'
    }
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

      <!-- Right auth panels container -->
      <div class="max-w-md w-full mx-auto relative py-4">
        
        <!-- Error Alerts -->
        <div 
          v-if="errorMessage" 
          class="mb-6 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake"
        >
          <span>⚠️</span>
          <p class="flex-1 leading-normal">{{ errorMessage }}</p>
        </div>

        <div class="space-y-8">
          <div>
            <h1 class="text-4xl font-medium tracking-tight mb-2">Welcome to Chia Florist</h1>
            <p class="text-gray-600">Enter your details below</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-6">
            <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
              <input 
                type="email" 
                v-model="email"
                placeholder="Email" 
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

            <div class="pt-4 space-y-4">
              <button
                type="submit"
                class="w-full bg-[#1b4332] text-white py-4 rounded-md font-medium hover:bg-[#143326] disabled:bg-gray-300 transition-all shadow-md cursor-pointer flex items-center justify-center"
                :disabled="authVm.isLoading.value"
              >
                <span v-if="authVm.isLoading.value" class="animate-pulse">Logging in...</span>
                <span v-else>Login</span>
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
              Sign up
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Shake animation for errors */
.animate-shake { animation: shake 0.3s ease-in-out; }
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}
</style>