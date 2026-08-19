<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { mapErrorMessage } from '~/utils/errorMessages'

definePageMeta({
  layout: 'auth'
})

useHead({ 
  title: 'Sign In - Chia Florist',
  meta: [
    { name: 'description', content: 'Sign in to access your Chia Florist profile.' }
  ]
})

const authVm = useAuthViewModel()
const globalAlert = useGlobalAlert()

// 2-Step Login Flow: 'initial' (Choose Google or Enter Email) -> 'credentials' (Password Form without Google)
const loginStep = ref<'initial' | 'credentials'>('initial')
const emailInputRef = ref<HTMLInputElement | null>(null)
const passwordInputRef = ref<HTMLInputElement | null>(null)

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const rememberMe = ref(false)

const errorMessage = ref('')
const successMessage = ref('')

watch(successMessage, (msg) => {
  if (msg) {
    globalAlert.showSuccess('Success', msg)
  }
})
const viewMode = ref<'login' | 'forgot_request' | 'forgot_verify' | 'forgot_reset'>('login')

const forgotEmail = ref('')
const forgotOtp = ref('')
const forgotChallengeId = ref('')
const newPassword = ref('')
const showNewPassword = ref(false)

onMounted(() => {
  nextTick(() => {
    emailInputRef.value?.focus()
  })

  const route = useRoute()
  if (route.query.error || route.query.error_description || route.query.google_error) {
    globalAlert.showError(
      'Google sign-in failed',
      "We couldn't sign you in with Google. Please try again.",
      [
        { label: 'Try Again', onClick: () => emailInputRef.value?.focus() },
        { label: 'Dismiss' }
      ]
    )
  }
})

const focusPasswordInput = () => {
  nextTick(() => {
    passwordInputRef.value?.focus()
    passwordInputRef.value?.select()
  })
  // Secondary backup for transition timing
  setTimeout(() => {
    passwordInputRef.value?.focus()
  }, 150)
}

const handleInitialEmailSubmit = () => {
  if (!email.value) {
    errorMessage.value = 'Please enter your email address.'
    return
  }
  errorMessage.value = ''
  loginStep.value = 'credentials'
  focusPasswordInput()
}

const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = 'Please fill in all required fields.'
    return
  }

  errorMessage.value = ''
  successMessage.value = ''

  try {
    const success = await authVm.login({
      email: email.value,
      password: password.value
    }, rememberMe.value)

    if (success) {
      navigateTo('/')
    }
  } catch (err: any) {
    if (err.status === 403 && (err.data?.message === 'email not verified' || err.data?.message?.includes('not verified'))) {
      if (import.meta.client) {
        localStorage.setItem('register_email', email.value)
      }
      errorMessage.value = 'Email not verified. Redirecting to verification...'
      setTimeout(() => {
        errorMessage.value = ''
        navigateTo('/register?verify=true')
      }, 1200)
    } else {
      errorMessage.value = mapErrorMessage(err, 'Login failed. Please check your credentials.')
    }
  }
}

const handleGoogleLogin = () => {
  const rememberMeCookie = useCookie('remember_me')
  rememberMeCookie.value = rememberMe.value ? 'true' : 'false'

  sessionStorage.setItem('google_auth_pending', '1')
  window.location.href = '/api/auth/google'
}

const handleForgotPasswordRequest = async () => {
  if (!forgotEmail.value) {
    errorMessage.value = 'Please enter your email address.'
    return
  }
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await authVm.requestForgotPassword(forgotEmail.value)
    successMessage.value = res.message || 'Verification code sent to your email.'
    if (res.challenge_id) {
      forgotChallengeId.value = res.challenge_id
    }
    viewMode.value = 'forgot_verify'
  } catch (err: any) {
    errorMessage.value = mapErrorMessage(err, 'Failed to request password reset.')
  }
}

const handleForgotPasswordVerify = async () => {
  if (!forgotOtp.value) {
    errorMessage.value = 'Please enter the verification code.'
    return
  }
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const res = await authVm.verifyForgotPasswordOtp(forgotChallengeId.value, forgotOtp.value)
    successMessage.value = res.message || 'Code verified successfully.'
    if (res.challenge_id) {
      forgotChallengeId.value = res.challenge_id
    }
    viewMode.value = 'forgot_reset'
  } catch (err: any) {
    errorMessage.value = mapErrorMessage(err, 'Verification failed.')
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
    loginStep.value = 'credentials'
    email.value = forgotEmail.value
    forgotEmail.value = ''
    forgotOtp.value = ''
    forgotChallengeId.value = ''
    newPassword.value = ''
    focusPasswordInput()
  } catch (err: any) {
    errorMessage.value = mapErrorMessage(err, 'Failed to reset password.')
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
  <div class="w-full space-y-6">

    <!-- Error & Success Alerts -->
    <div
      v-if="errorMessage"
      class="bg-red-50 border border-red-200 text-red-700 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake shadow-xs"
    >
      <span>⚠️</span>
      <p class="flex-1 leading-normal">{{ errorMessage }}</p>
    </div>

    <div
      v-if="successMessage"
      class="bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 shadow-xs"
    >
      <span>✅</span>
      <p class="flex-1 leading-normal">{{ successMessage }}</p>
    </div>

    <!-- MAIN VIEW: LOGIN -->
    <div v-if="viewMode === 'login'">
      <Transition name="fade-slide" mode="out-in" @after-enter="focusPasswordInput">

        <!-- STEP 1: INITIAL METHOD SELECTION (Google or Enter Email) -->
        <div v-if="loginStep === 'initial'" key="step-initial" class="space-y-6">

          <!-- Heading -->
          <div class="text-left">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Sign in</h1>
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
            Don't have an account?
            <NuxtLink to="/register" class="font-bold text-gray-900 hover:text-[#245842] underline ml-1 transition-colors">
              Create an account
            </NuxtLink>
          </div>

        </div>

        <!-- STEP 2: CREDENTIALS FORM (NO GOOGLE BUTTON, AUTO-FOCUSED & SELECTED PASSWORD) -->
        <div v-else-if="loginStep === 'credentials'" key="step-credentials" class="space-y-6">

          <!-- Heading -->
          <div class="text-left space-y-1">
            <h1 class="text-3xl font-bold tracking-tight text-gray-900">Sign in</h1>
            <p class="text-sm text-gray-500">
              Signing in as <span class="font-medium text-gray-900">{{ email }}</span>
              <button
                type="button"
                @click="loginStep = 'initial'"
                class="text-xs font-semibold text-[#245842] hover:underline ml-1 focus:outline-none"
              >
                (Change)
              </button>
            </p>
          </div>

          <!-- App Password Form -->
          <form @submit.prevent="handleLogin" class="space-y-4">

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
              <div class="flex items-center justify-between">
                <label class="text-xs font-medium text-gray-700 block">Password</label>
                <button
                  type="button"
                  @click="switchToForgot"
                  class="text-xs font-medium text-gray-500 hover:text-gray-900 transition-colors focus:outline-none"
                >
                  Forgot password?
                </button>
              </div>
              <div class="relative flex items-center">
                <input
                  ref="passwordInputRef"
                  :type="showPassword ? 'text' : 'password'"
                  v-model="password"
                  placeholder="Enter your password"
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

            <!-- Remember me Checkbox Option -->
            <div class="pt-1 flex items-center justify-between">
              <label class="flex items-center gap-2 text-xs text-gray-600 hover:text-gray-900 cursor-pointer select-none transition-colors">
                <input
                  type="checkbox"
                  v-model="rememberMe"
                  class="w-4 h-4 rounded border-gray-300 text-[#245842] focus:ring-[#4ade80] accent-[#4ade80]"
                  :disabled="authVm.isLoading.value"
                />
                Remember me
              </label>
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
                Sign In
              </CButton>

              <CButton
                type="button"
                @click="loginStep = 'initial'"
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
            Don't have an account?
            <NuxtLink to="/register" class="font-bold text-gray-900 hover:text-[#245842] underline ml-1 transition-colors">
              Create an account
            </NuxtLink>
          </div>

        </div>

      </Transition>
    </div>

    <!-- FORGOT PASSWORD: REQUEST -->
    <div v-else-if="viewMode === 'forgot_request'" class="space-y-6">
      <div class="text-left">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Reset password</h1>
        <p class="text-sm text-gray-500">Enter your email address to receive a reset code</p>
      </div>

      <form @submit.prevent="handleForgotPasswordRequest" class="space-y-4">
        <div class="space-y-1">
          <label class="text-xs font-medium text-gray-700 block">Email</label>
          <input
            type="email"
            v-model="forgotEmail"
            placeholder="Enter your email"
            class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400"
            :disabled="authVm.isLoading.value"
            required
          />
        </div>

        <div class="pt-2 space-y-3">
          <button
            type="submit"
            class="w-full bg-[#4ade80] hover:bg-[#34d399] text-[#245842] py-3 rounded-xl text-sm font-bold shadow-sm transition-all duration-200 cursor-pointer flex items-center justify-center gap-2"
            :disabled="authVm.isLoading.value"
          >
            <span v-if="authVm.isLoading.value" class="animate-pulse">Sending code...</span>
            <span v-else>Send Reset Code</span>
          </button>

          <button
            type="button"
            @click="switchToLogin"
            class="w-full border border-gray-200 text-gray-700 py-3 rounded-xl text-sm font-semibold hover:bg-gray-50 transition-all cursor-pointer"
            :disabled="authVm.isLoading.value"
          >
            Back to Sign in
          </button>
        </div>
      </form>
    </div>

    <!-- FORGOT PASSWORD: VERIFY OTP -->
    <div v-else-if="viewMode === 'forgot_verify'" class="space-y-6">
      <div class="text-left">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Verify code</h1>
        <p class="text-sm text-gray-500">Enter the verification code sent to your email</p>
      </div>

      <form @submit.prevent="handleForgotPasswordVerify" class="space-y-4">
        <div class="space-y-1">
          <label class="text-xs font-medium text-gray-700 block text-center">Verification Code</label>
          <input
            type="text"
            v-model="forgotOtp"
            placeholder="000000"
            maxlength="6"
            class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-center text-xl font-mono tracking-[0.5em] indent-[0.25em] outline-none focus:bg-white focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all"
            :disabled="authVm.isLoading.value"
            required
          />
        </div>

        <div class="pt-2 space-y-3">
          <button
            type="submit"
            class="w-full bg-[#4ade80] hover:bg-[#34d399] text-[#245842] py-3 rounded-xl text-sm font-bold shadow-sm transition-all duration-200 cursor-pointer flex items-center justify-center gap-2"
            :disabled="authVm.isLoading.value"
          >
            <span v-if="authVm.isLoading.value" class="animate-pulse">Verifying...</span>
            <span v-else>Verify Code</span>
          </button>

          <button
            type="button"
            @click="switchToLogin"
            class="w-full border border-gray-200 text-gray-700 py-3 rounded-xl text-sm font-semibold hover:bg-gray-50 transition-all cursor-pointer"
            :disabled="authVm.isLoading.value"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- FORGOT PASSWORD: RESET -->
    <div v-else-if="viewMode === 'forgot_reset'" class="space-y-6">
      <div class="text-left">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 mb-1">Set new password</h1>
        <p class="text-sm text-gray-500">Create a new password for your account</p>
      </div>

      <form @submit.prevent="handleForgotPasswordReset" class="space-y-4">
        <div class="space-y-1">
          <label class="text-xs font-medium text-gray-700 block">New Password</label>
          <div class="relative flex items-center">
            <input
              :type="showNewPassword ? 'text' : 'password'"
              v-model="newPassword"
              placeholder="Enter new password"
              class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all placeholder:text-gray-400 pr-14"
              :disabled="authVm.isLoading.value"
              required
            />
            <button
              type="button"
              @click="showNewPassword = !showNewPassword"
              class="absolute right-4 text-xs font-semibold text-gray-500 hover:text-gray-900 focus:outline-none"
              tabindex="-1"
            >
              {{ showNewPassword ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>

        <div class="pt-2">
          <button
            type="submit"
            class="w-full bg-[#4ade80] hover:bg-[#34d399] text-[#245842] py-3 rounded-xl text-sm font-bold shadow-sm transition-all duration-200 cursor-pointer flex items-center justify-center gap-2"
            :disabled="authVm.isLoading.value"
          >
            <span v-if="authVm.isLoading.value" class="animate-pulse">Saving new password...</span>
            <span v-else>Reset Password</span>
          </button>
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
