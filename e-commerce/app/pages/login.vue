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
const showPassword = ref(false)

const errorMessage = ref('')
const successMessage = ref('')
const viewMode = ref<'login' | 'forgot_request' | 'forgot_verify' | 'forgot_reset'>('login')

const forgotEmail = ref('')
const forgotOtp = ref('')
const forgotChallengeId = ref('')
const newPassword = ref('')
const showNewPassword = ref(false)

const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = 'Please fill in all fields.'
    return
  }

  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const success = await authVm.login({
      email: email.value,
      password: password.value
    })
    
    if (success) {
      navigateTo('/')
    }
  } catch (err: any) {
    if (err.status === 403 && err.data?.message === 'email not verified') {
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

const handleGoogleLogin = () => {
  sessionStorage.setItem('google_auth_pending', '1')
  window.location.href = '/api/auth/google'
}

const handleForgotPasswordRequest = async () => {
  if (!forgotEmail.value) {
    errorMessage.value = 'Please enter your email.'
    return
  }
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await authVm.requestForgotPassword(forgotEmail.value)
    successMessage.value = res.message || 'OTP sent to your email.'
    if (res.challenge_id) {
      forgotChallengeId.value = res.challenge_id
    }
    viewMode.value = 'forgot_verify'
  } catch (err: any) {
    errorMessage.value = err.data?.message || err.message || 'Failed to request password reset.'
  }
}

const handleForgotPasswordVerify = async () => {
  if (!forgotOtp.value) {
    errorMessage.value = 'Please enter the OTP.'
    return
  }
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await authVm.verifyForgotPasswordOtp(forgotChallengeId.value, forgotOtp.value)
    successMessage.value = res.message || 'OTP verified.'
    if (res.challenge_id) {
      forgotChallengeId.value = res.challenge_id
    }
    viewMode.value = 'forgot_reset'
  } catch (err: any) {
    errorMessage.value = err.data?.message || err.message || 'OTP verification failed.'
  }
}

const handleForgotPasswordReset = async () => {
  if (!newPassword.value) {
    errorMessage.value = 'Please enter a new password.'
    return
  }
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await authVm.resetPassword(forgotChallengeId.value, newPassword.value)
    successMessage.value = res.message || 'Password reset successful. You can now log in.'
    viewMode.value = 'login'
    email.value = forgotEmail.value
    forgotEmail.value = ''
    forgotOtp.value = ''
    forgotChallengeId.value = ''
    newPassword.value = ''
  } catch (err: any) {
    errorMessage.value = err.data?.message || err.message || 'Failed to reset password.'
  }
}

const switchToLogin = () => {
  viewMode.value = 'login'
  errorMessage.value = ''
  successMessage.value = ''
}

const switchToForgot = () => {
  viewMode.value = 'forgot_request'
  errorMessage.value = ''
  successMessage.value = ''
  forgotEmail.value = email.value
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-10 min-h-[80vh] flex flex-col items-center gap-10">
    <div>
      <img src="/images/logo.png" class="h-20 mx-auto mb-4" alt="Chia Florist Logo" />
    </div>
    
    <div class="max-w-md w-full mx-auto relative py-4">
      
      <!-- Alerts -->
      <div 
        v-if="errorMessage" 
        class="mb-6 bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake"
      >
        <span>⚠️</span>
        <p class="flex-1 leading-normal">{{ errorMessage }}</p>
      </div>

      <div 
        v-if="successMessage" 
        class="mb-6 bg-green-50 border border-green-100 text-green-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2"
      >
        <span>✅</span>
        <p class="flex-1 leading-normal">{{ successMessage }}</p>
      </div>

      <!-- Login View -->
      <div v-if="viewMode === 'login'" class="space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Sign in</h1>
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
          
          <div class="flex justify-end">
            <button type="button" @click="switchToForgot" class="text-sm text-gray-500 hover:text-black transition-colors focus:outline-none">
              Forgot Password?
            </button>
          </div>

          <div class="pt-2 space-y-4">
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
              @click="handleGoogleLogin"
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
              :disabled="authVm.isLoading.value"
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

      <!-- Forgot Password - Request View -->
      <div v-else-if="viewMode === 'forgot_request'" class="space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Reset Password</h1>
          <p class="text-gray-600">Enter your email to receive an OTP</p>
        </div>

        <form @submit.prevent="handleForgotPasswordRequest" class="space-y-6">
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="email" 
              v-model="forgotEmail"
              placeholder="Email" 
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
              <span v-if="authVm.isLoading.value" class="animate-pulse">Sending...</span>
              <span v-else>Send Reset Code</span>
            </button>

            <button 
              type="button"
              @click="switchToLogin"
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
              :disabled="authVm.isLoading.value"
            >
              Back to Login
            </button>
          </div>
        </form>
      </div>

      <!-- Forgot Password - Verify OTP View -->
      <div v-else-if="viewMode === 'forgot_verify'" class="space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Verify OTP</h1>
          <p class="text-gray-600">Enter the OTP sent to your email</p>
        </div>

        <form @submit.prevent="handleForgotPasswordVerify" class="space-y-6">
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              v-model="forgotOtp"
              placeholder="Enter OTP" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400 tracking-widest"
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
              <span v-if="authVm.isLoading.value" class="animate-pulse">Verifying...</span>
              <span v-else>Verify Code</span>
            </button>

            <button 
              type="button"
              @click="switchToLogin"
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all cursor-pointer"
              :disabled="authVm.isLoading.value"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>

      <!-- Forgot Password - Reset Password View -->
      <div v-else-if="viewMode === 'forgot_reset'" class="space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">New Password</h1>
          <p class="text-gray-600">Create a new password for your account</p>
        </div>

        <form @submit.prevent="handleForgotPasswordReset" class="space-y-6">
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors flex items-center justify-between">
            <input 
              :type="showNewPassword ? 'text' : 'password'" 
              v-model="newPassword"
              placeholder="New Password" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
              :disabled="authVm.isLoading.value"
              required
            />
            <button 
              type="button" 
              @click="showNewPassword = !showNewPassword" 
              class="text-sm font-medium text-gray-500 hover:text-black focus:outline-none ml-2"
              tabindex="-1"
            >
              {{ showNewPassword ? 'Hide' : 'Show' }}
            </button>
          </div>

          <div class="pt-4 space-y-4">
            <button
              type="submit"
              class="w-full bg-[#1b4332] text-white py-4 rounded-md font-medium hover:bg-[#143326] disabled:bg-gray-300 transition-all shadow-md cursor-pointer flex items-center justify-center"
              :disabled="authVm.isLoading.value"
            >
              <span v-if="authVm.isLoading.value" class="animate-pulse">Saving...</span>
              <span v-else>Reset Password</span>
            </button>
          </div>
        </form>
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