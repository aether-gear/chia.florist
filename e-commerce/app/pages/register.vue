<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'

useHead({ 
  title: 'Create Account - Chia Florist',
  meta: [
    { name: 'description', content: 'Create a new account or verify your registration at Chia Florist.' }
  ]
})

const route = useRoute()
const authVm = useAuthViewModel()

const activePanel = ref<'register' | 'verify'>('register')
const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const phone = ref('')
const otpCode = ref('')
const errorMessage = ref('')

// Load values from viewModel if they exist (e.g. from localStorage)
const registrationEmail = computed(() => authVm.registrationEmail.value)

onMounted(() => {
  // If route query verify=true, automatically slide to verify on mount
  if (route.query.verify === 'true' || route.query.verify === '1') {
    if (authVm.registrationEmail.value) {
      email.value = authVm.registrationEmail.value
    }
    
    // Smooth entrance slide
    setTimeout(() => {
      activePanel.value = 'verify'
    }, 400)
  }
})

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
      // Smooth slide transition to the verify panel
      activePanel.value = 'verify'
    }
  } catch (err: any) {
    errorMessage.value = err.data?.message || 'Registration failed. Please check your input details.'
  }
}

const handleVerify = async () => {
  if (!otpCode.value || otpCode.value.length !== 6) {
    errorMessage.value = 'Please enter a valid 6-digit verification code.'
    return
  }

  errorMessage.value = ''

  try {
    const success = await authVm.verifyOtp(otpCode.value)
    if (success) {
      navigateTo('/')
    }
  } catch (err: any) {
    errorMessage.value = err.data?.message || 'Verification failed. Please try again.'
  }
}

const handleGoogleLogin = () => {
  // Signal to auth-init that we're coming back from Google OAuth.
  // sessionStorage survives page navigation in the same tab, so this
  // flag will still be present when the backend redirects us back to /.
  sessionStorage.setItem('google_auth_pending', '1')
  window.location.href = '/api/auth/google'
}

const handleBackToRegister = () => {
  errorMessage.value = ''
  activePanel.value = 'register'
}
</script>

<template>
  <div class="
    max-w-7xl mx-auto px-8 py-20 mt-10 min-h-[80vh] flex flex-col items-center
    gap-10
  ">
    <div class="">
      <img src="/images/logo.png" class="h-20 mx-auto mb-4" alt="Chia Florist Logo" />
    </div>
    <div class="max-w-md w-full mx-auto overflow-hidden relative py-4">
      
      <!-- Error Alerts -->
      <div 
        v-if="errorMessage" 
        class="mb-6 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake"
      >
        <span>⚠️</span>
        <p class="flex-1 leading-normal">{{ errorMessage }}</p>
      </div>

      <div 
        class="flex transition-transform duration-700 ease-[cubic-bezier(0.16,1,0.3,1)]"
        :style="{ transform: activePanel === 'register' ? 'translateX(0%)' : 'translateX(-50%)', width: '200%' }"
      >
        
        <!-- PANEL 1: REGISTER -->
        <div class="w-1/2 pr-6 space-y-8 flex-shrink-0">
          <div>
            <h1 class="text-4xl font-medium tracking-tight mb-2">Create an account</h1>
            <p class="text-gray-600">Enter your details below</p>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-6">
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

            <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors flex items-center justify-between">
              <input 
                :type="showPassword ? 'text' : 'password'" 
                v-model="password"
                placeholder="Password" 
                class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
                :disabled="authVm.isLoading.value"
                required
              />
              <button 
                type="button" 
                @click="showPassword = !showPassword" 
                class="text-sm font-medium text-gray-500 hover:text-black focus:outline-none ml-2"
                tabindex="-1"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
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
                @click="handleGoogleLogin"
                class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
                :disabled="authVm.isLoading.value"
              >
                <img src="/images/google.png" class="w-5 h-5" alt="Google Icon" />
                Sign up with Google
              </button>
            </div>
          </form>

          <div class="text-center pt-4 text-gray-600">
            Already have an account? 
            <NuxtLink to="/login" class="font-medium text-black border-b border-gray-500 pb-0.5 ml-2 hover:text-[#1b4332] transition-colors">
              Log in
            </NuxtLink>
          </div>
        </div>

        <!-- PANEL 2: VERIFICATION -->
        <div class="w-1/2 pl-6 space-y-8 flex-shrink-0">
          <div class="flex flex-col items-center text-center">
            <h2 class="text-xl font-bold text-gray-900 tracking-tight">
              Verify Your Account
            </h2>
            <p class="text-xs text-gray-400 mt-2 max-w-xs leading-relaxed">
              We have sent a secure verification code to <br/>
              <span class="text-gray-700 font-semibold break-all">{{ registrationEmail || email || 'your email' }}</span>
            </p>
          </div>

          <form @submit.prevent="handleVerify" class="space-y-6">
            <div class="space-y-2">
              <label class="text-[11px] font-black uppercase tracking-widest text-gray-400 block text-center">
                Enter Verification Code
              </label>
              <input 
                type="text" 
                v-model="otpCode"
                placeholder="000000"
                maxlength="6"
                class="w-full px-5 py-3.5 bg-gray-50 border border-gray-200 rounded-2xl text-center text-xl font-mono tracking-[0.5em] indent-[0.25em] outline-none focus:bg-white focus:border-[#1b4332] focus:ring-2 focus:ring-[#1b4332]/5 transition-all"
                :disabled="authVm.isLoading.value"
                required
              />
            </div>

            <div class="grid grid-cols-2 gap-4 pt-2">
              <button 
                type="button"
                @click="handleBackToRegister"
                class="w-full py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 text-xs font-bold rounded-xl transition-all outline-none cursor-pointer"
                :disabled="authVm.isLoading.value"
              >
                Back
              </button>
              
              <button 
                type="submit"
                class="w-full py-3 bg-[#1b4332] hover:bg-[#143326] disabled:bg-gray-300 text-white text-xs font-bold rounded-xl shadow-sm transition-all outline-none flex items-center justify-center cursor-pointer"
                :disabled="authVm.isLoading.value"
              >
                <span v-if="authVm.isLoading.value" class="animate-pulse">Verifying...</span>
                <span v-else>Confirm</span>
              </button>
            </div>
          </form>
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