<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { mapErrorMessage } from '~/utils/errorMessages'

definePageMeta({
  layout: 'auth'
})

useHead({ 
  title: 'Create Account - Chia Florist',
  meta: [
    { name: 'description', content: 'Create a new account or verify your registration at Chia Florist.' }
  ]
})

const route = useRoute()
const authVm = useAuthViewModel()
const globalAlert = useGlobalAlert()

// 2-Step Registration Flow: 'initial' (Choose Google or Enter Email) -> 'form' (Account Details without Google)
const registerStep = ref<'initial' | 'form'>('initial')
const activePanel = ref<'register' | 'verify'>('register')
const emailInputRef = ref<HTMLInputElement | null>(null)
const nameInputRef = ref<HTMLInputElement | null>(null)

const name = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const phone = ref('')
const otpCode = ref('')
const errorMessage = ref('')

const registrationEmail = computed(() => authVm.registrationEmail.value)

onMounted(() => {
  nextTick(() => {
    emailInputRef.value?.focus()
  })

  if (route.query.error || route.query.error_description || route.query.google_error) {
    globalAlert.showError('Google sign-in failed', "We couldn't sign you in with Google. Please try again.")
  }

  if (route.query.verify === 'true' || route.query.verify === '1') {
    if (authVm.registrationEmail.value) {
      email.value = authVm.registrationEmail.value
    }
    
    setTimeout(() => {
      activePanel.value = 'verify'
    }, 300)
  }
})

const focusNameInput = () => {
  nextTick(() => {
    nameInputRef.value?.focus()
    nameInputRef.value?.select()
  })
  setTimeout(() => {
    nameInputRef.value?.focus()
  }, 150)
}

const handleInitialEmailSubmit = () => {
  if (!email.value) {
    errorMessage.value = 'Please enter your email address.'
    return
  }
  errorMessage.value = ''
  
  // Auto-suggest name/username from email prefix if empty
  if (!name.value) {
    const prefix = email.value.split('@')[0] || ''
    name.value = prefix.charAt(0).toUpperCase() + prefix.slice(1)
    username.value = prefix.toLowerCase().replace(/[^a-z0-9_]/g, '')
  }
  
  registerStep.value = 'form'
  focusNameInput()
}

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
      activePanel.value = 'verify'
    }
  } catch (err: any) {
    errorMessage.value = mapErrorMessage(err, 'Registration failed. Please check your input details.')
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
    errorMessage.value = mapErrorMessage(err, 'Verification failed. Please try again.')
  }
}

const handleGoogleLogin = () => {
  sessionStorage.setItem('google_auth_pending', '1')
  window.location.href = '/api/auth/google'
}

const handleBackToRegister = () => {
  errorMessage.value = ''
  activePanel.value = 'register'
}
</script>

<template>
  <div class="w-full space-y-6">

    <!-- Error Alert -->
    <div
      v-if="errorMessage"
      class="bg-red-50 border border-red-200 text-red-700 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake shadow-xs"
    >
      <span>⚠️</span>
      <p class="flex-1 leading-normal">{{ errorMessage }}</p>
    </div>

    <!-- REGISTER PANEL -->
    <div v-if="activePanel === 'register'">
      <Transition name="fade-slide" mode="out-in" @after-enter="focusNameInput">

        <!-- STEP 1: INITIAL METHOD SELECTION (Google or Enter Email) -->
        <div v-if="registerStep === 'initial'" key="step-initial" class="space-y-6">

          <!-- Heading -->
          <div class="text-left">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Create an account</h1>
            <p class="text-sm text-gray-500">Sign in or create an account</p>
          </div>

          <!-- Primary Action: Continue with Google (Full Wordmark Logo) -->
          <CButton
            type="button"
            @click="handleGoogleLogin"
            variant="primary"
            size="lg"
            full-width
            :loading="authVm.isLoading.value"
          >
            <span class="flex items-center gap-1">
              <span class="pb-1">Continue with</span>
              <img
                src="/images/google-logo.svg"
                class="h-4.5 w-auto"
                alt="Google"
              />
            </span>
          </CButton>

          <!-- Horizontal Divider: OR -->
          <div class="relative flex py-1 items-center">
            <div class="flex-grow border-t border-gray-200"></div>
            <span class="flex-shrink mx-4 text-xs font-medium text-gray-400 tracking-wide">or</span>
            <div class="flex-grow border-t border-gray-200"></div>
          </div>

          <!-- Initial Email Entry Form -->
          <form @submit.prevent="handleInitialEmailSubmit" class="space-y-4">
            <div class="relative flex items-center">
              <input
                ref="emailInputRef"
                type="email"
                v-model="email"
                placeholder="Email"
                class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400 pr-12"
                :disabled="authVm.isLoading.value"
                required
                autofocus
              />
              <button
                type="submit"
                class="absolute right-2 p-2 bg-[#4ade80] hover:bg-[#34d399] text-[#245842] font-bold rounded-xl transition-all cursor-pointer flex items-center justify-center"
                title="Continue with Email"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
              </button>
            </div>
          </form>

          <!-- Terms & Legal Disclaimer -->
          <p class="text-center text-xs text-gray-400 leading-relaxed pt-1">
            By continuing, you agree to our
            <NuxtLink to="/terms" class="underline hover:text-gray-700 transition-colors">Terms of service</NuxtLink>.
          </p>

          <!-- Account Switch -->
          <div class="text-center pt-2 text-xs text-gray-500">
            Already have an account?
            <NuxtLink to="/login" class="font-bold text-gray-900 hover:text-[#245842] underline ml-1 transition-colors">
              Sign in
            </NuxtLink>
          </div>

        </div>

        <!-- STEP 2: FULL DETAILS FORM (NO GOOGLE BUTTON, AUTO-FOCUSED & SELECTED NAME) -->
        <div v-else-if="registerStep === 'form'" key="step-form" class="space-y-6">

          <!-- Heading -->
          <div class="text-left space-y-1">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900">Create an account</h1>
            <p class="text-sm text-gray-500">
              Registering for <span class="font-medium text-gray-900">{{ email }}</span>
              <button
                type="button"
                @click="registerStep = 'initial'"
                class="text-xs font-semibold text-[#245842] hover:underline ml-1 focus:outline-none"
              >
                (Change)
              </button>
            </p>
          </div>

          <!-- Registration Form -->
          <form @submit.prevent="handleRegister" class="space-y-4">

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div class="space-y-1">
                <label class="text-xs font-medium text-gray-700 block">Full Name</label>
                <input
                  ref="nameInputRef"
                  type="text"
                  v-model="name"
                  placeholder="Name"
                  class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400"
                  :disabled="authVm.isLoading.value"
                  required
                />
              </div>

              <div class="space-y-1">
                <label class="text-xs font-medium text-gray-700 block">Username</label>
                <input
                  type="text"
                  v-model="username"
                  placeholder="Username"
                  class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400"
                  :disabled="authVm.isLoading.value"
                  required
                />
              </div>
            </div>

            <div class="space-y-1">
              <label class="text-xs font-medium text-gray-700 block">Email</label>
              <input
                type="email"
                v-model="email"
                placeholder="Email"
                class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all text-gray-700"
                :disabled="authVm.isLoading.value"
                required
              />
            </div>

            <div class="space-y-1">
              <label class="text-xs font-medium text-gray-700 block">Password</label>
              <div class="relative flex items-center">
                <input
                  :type="showPassword ? 'text' : 'password'"
                  v-model="password"
                  placeholder="Create a password"
                  class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400 pr-14"
                  :disabled="authVm.isLoading.value"
                  required
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute right-4 text-xs font-semibold text-gray-500 hover:text-gray-900 focus:outline-none"
                  tabindex="-1"
                >
                  {{ showPassword ? 'Hide' : 'Show' }}
                </button>
              </div>
            </div>

            <div class="space-y-1">
              <label class="text-xs font-medium text-gray-700 block">Phone Number (Optional)</label>
              <input
                type="tel"
                v-model="phone"
                placeholder="Phone number"
                class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400"
                :disabled="authVm.isLoading.value"
              />
            </div>

            <!-- Action Buttons -->
            <div class="pt-2 space-y-3">
              <CButton
                type="submit"
                variant="primary"
                size="lg"
                full-width
                :loading="authVm.isLoading.value"
              >
                Create Account
              </CButton>

              <CButton
                type="button"
                @click="registerStep = 'initial'"
                variant="outline"
                size="lg"
                full-width
                :disabled="authVm.isLoading.value"
              >
                Back
              </CButton>
            </div>
          </form>

          <!-- Account Switch -->
          <div class="text-center pt-2 text-xs text-gray-500">
            Already have an account?
            <NuxtLink to="/login" class="font-bold text-gray-900 hover:text-[#245842] underline ml-1 transition-colors">
              Sign in
            </NuxtLink>
          </div>

        </div>

      </Transition>
    </div>

    <!-- VERIFICATION PANEL -->
    <div v-else-if="activePanel === 'verify'" class="space-y-6">
      <div class="text-left">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Verify Account</h1>
        <p class="text-sm text-gray-500">
          We sent a verification code to <span class="font-medium text-gray-900">{{ registrationEmail || email || 'your email' }}</span>
        </p>
      </div>

      <form @submit.prevent="handleVerify" class="space-y-4">
        <div class="space-y-1">
          <label class="text-xs font-medium text-gray-700 block text-center">Verification Code</label>
          <input
            type="text"
            v-model="otpCode"
            placeholder="000000"
            maxlength="6"
            class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-center text-xl font-mono tracking-[0.5em] indent-[0.25em] outline-none focus:bg-white focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all"
            :disabled="authVm.isLoading.value"
            required
          />
        </div>

        <div class="grid grid-cols-2 gap-3 pt-2">
          <CButton
            type="button"
            @click="handleBackToRegister"
            variant="outline"
            size="lg"
            full-width
            :disabled="authVm.isLoading.value"
          >
            Back
          </CButton>

          <CButton
            type="submit"
            variant="primary"
            size="lg"
            full-width
            :loading="authVm.isLoading.value"
          >
            Confirm
          </CButton>
        </div>
      </form>
    </div>

  </div>
</template>

<style scoped>
.animate-shake { animation: shake 0.3s ease-in-out; }
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}

/* Step transition animation */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(12px);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
</style>
