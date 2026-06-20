// app/composables/viewmodels/useAuthViewModel.ts
import { ref, computed } from 'vue'
import { authService } from '~/services/authService'
import type { UserMe, SignUpRequest, VerifyRequest, SignInRequest } from '~/types/auth'

// Shared global state so all components see the same auth state
const currentUser = ref<UserMe | null>(null)
const challengeId = ref<string | null>(null)
const registrationEmail = ref<string | null>(null)
const isInitialized = ref(false)

export const useAuthViewModel = () => {
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const isAuthenticated = computed(() => currentUser.value !== null)

  /**
   * Safe method to fetch the current logged-in user.
   */
  const fetchCurrentUser = async (cookieHeader?: string) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authService.getMe(cookieHeader)
      if (response && response.is_authenticated) {
        const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
        const isLoggedIn = useCookie('is_logged_in')
        
        isLoggedIn.value = 'true'
        
        if (!userProfile.value) {
          userProfile.value = {
            id: response.account_id,
            name: 'Customer',
            username: 'customer',
            phone: '',
            last_login_at: new Date().toISOString()
          }
        } else {
          userProfile.value.id = response.account_id
        }
        
        currentUser.value = userProfile.value as UserMe
      } else {
        currentUser.value = null
        const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
        const isLoggedIn = useCookie('is_logged_in')
        userProfile.value = null
        isLoggedIn.value = null
      }
    } catch (err: any) {
      currentUser.value = null
      const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
      const isLoggedIn = useCookie('is_logged_in')
      userProfile.value = null
      isLoggedIn.value = null
      console.warn('Failed to fetch user state:', err)
    } finally {
      isLoading.value = false
      isInitialized.value = true
    }
  }

  /**
   * Authenticate user.
   */
  const login = async (credentials: SignInRequest) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await authService.signIn(credentials)
      if (response.message === 'login success') {
        if (import.meta.client) {
          localStorage.removeItem('chia-florist-cart-cache')
        }
        // Fetch profile to populate global user state
        await fetchCurrentUser()
        const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
        userProfile.value = {
          id: 'temp-id',
          name: credentials.email.split('@')[0],
          username: credentials.email.split('@')[0],
          email: credentials.email,
          phone: '',
          last_login_at: new Date().toISOString()
        }
        
        const isLoggedIn = useCookie('is_logged_in')
        isLoggedIn.value = 'true'
        
        currentUser.value = userProfile.value as UserMe
        return true
      }
      return false
    } catch (err: any) {
      error.value = err.data?.message || 'Login failed. Please check your credentials.'
      currentUser.value = null
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Register a new user account.
   */
  const register = async (signUpData: SignUpRequest) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authService.signUp(signUpData)
      if (response && response.challenge_id) {
        challengeId.value = response.challenge_id
        registrationEmail.value = signUpData.email
        
        // Save to localStorage as a fallback/state transfer mechanism
        if (import.meta.client) {
          localStorage.setItem('auth_challenge_id', response.challenge_id)
          localStorage.setItem('register_email', signUpData.email)
          localStorage.setItem('register_name', signUpData.name)
          localStorage.setItem('register_username', signUpData.username)
          localStorage.setItem('register_phone', signUpData.phone || '')
        }
        return response
      }
      throw new Error('Registration did not return a verification challenge.')
    } catch (err: any) {
      error.value = err.data?.message || 'Registration failed. Please check the inputs.'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Verify account with OTP code.
   */
  const verifyOtp = async (otpCode: number) => {
    isLoading.value = true
    error.value = null
    
    // Restore from localStorage if state was lost on page reload
    if (!challengeId.value && import.meta.client) {
      challengeId.value = localStorage.getItem('auth_challenge_id')
    }
    if (!challengeId.value) {
      const errMsg = 'No active verification challenge found. Please sign up again.'
      error.value = errMsg
      isLoading.value = false
      throw new Error(errMsg)
    }
    try {
      const reqData: VerifyRequest = {
        challenge_id: challengeId.value,
        otp: otpCode
      }
      
      const response = await authService.verify(reqData)
      if (response.message === 'verify success') {
        let name = 'Verified User'
        let username = 'user'
        let email = ''
        let phone = ''
        
        if (import.meta.client) {
          name = localStorage.getItem('register_name') || 'Verified User'
          username = localStorage.getItem('register_username') || 'user'
          email = localStorage.getItem('register_email') || ''
          phone = localStorage.getItem('register_phone') || ''
          
          localStorage.removeItem('auth_challenge_id')
          localStorage.removeItem('register_email')
          localStorage.removeItem('register_name')
          localStorage.removeItem('register_username')
          localStorage.removeItem('register_phone')
          localStorage.removeItem('chia-florist-cart-cache')
        }
        challengeId.value = null
        registrationEmail.value = null
        const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
        userProfile.value = {
          id: 'temp-id',
          name,
          username,
          email,
          phone,
          last_login_at: new Date().toISOString()
        }
        const isLoggedIn = useCookie('is_logged_in')
        isLoggedIn.value = 'true'
        currentUser.value = userProfile.value as UserMe
        return true
      }
      return false
    } catch (err: any) {
      error.value = err.data?.message || 'OTP Verification failed. Please check the code.'
      throw err
    } finally {
      isLoading.value = false
    }
  }


  /**
   * Log out the current user.
   */
  const logout = async () => {
    isLoading.value = true
    try {
      await authService.signOut()
    } catch (err) {
      console.error('Error hitting logout route:', err)
    } finally {
      currentUser.value = null
      challengeId.value = null
      registrationEmail.value = null
      
      const isLoggedIn = useCookie('is_logged_in')
      isLoggedIn.value = null
      const userProfile = useCookie<Partial<UserMe> | null>('user_profile')
      userProfile.value = null
      if (import.meta.client) {
        localStorage.removeItem('auth_challenge_id')
        localStorage.removeItem('register_email')
        localStorage.removeItem('register_name')
        localStorage.removeItem('register_username')
        localStorage.removeItem('register_phone')
        localStorage.removeItem('chia-florist-cart-cache')
      }
      isLoading.value = false
      navigateTo('/login')
    }
  }

  // Restore active challenge details on instantiation on the client
  if (import.meta.client && !registrationEmail.value) {
    registrationEmail.value = localStorage.getItem('register_email')
    challengeId.value = localStorage.getItem('auth_challenge_id')
  }

  return {
    currentUser: computed(() => currentUser.value),
    challengeId: computed(() => challengeId.value),
    registrationEmail: computed(() => registrationEmail.value),
    isInitialized: computed(() => isInitialized.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    isAuthenticated,
    fetchCurrentUser,
    login,
    register,
    verifyOtp,
    logout
  }
}
