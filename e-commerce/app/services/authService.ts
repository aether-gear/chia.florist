import type { 
  SignUpRequest, 
  SignUpResponse, 
  VerifyRequest, 
  VerifyResponse, 
  SignInRequest, 
  SignInResponse, 
  GetCurrentUserResponse,
  GetMeResponse,
  GetProfileResponse,
  UpdateProfileRequest
} from '~/types/auth'
import { bootstrapConfig } from '~/utils/bootstrap'

export const authService = {
  /**
   * Register a new user account.
   */
  async signUp(data: SignUpRequest): Promise<SignUpResponse> {
    return bootstrapConfig.fetchApi<SignUpResponse>('/auth/signup', {
      method: 'POST',
      body: data
    })
  },
  /**
   * Verify account using OTP code and challenge ID.
   */
  async verify(data: VerifyRequest): Promise<VerifyResponse> {
    return $fetch<VerifyResponse>('/api/auth/verify', {
      method: 'POST',
      body: data
    })
  },
  /**
   * Sign in user via email and password.
   */
  async signIn(data: SignInRequest): Promise<SignInResponse> {
    return $fetch<SignInResponse>('/api/auth/signin', {
      method: 'POST',
      body: data
    })
  },
  /**
   * Retrieve session state from local server route.
   */
  async getMe(cookieHeader?: string): Promise<GetMeResponse> {
    const headers: Record<string, string> = {}
    if (cookieHeader) {
      headers['cookie'] = cookieHeader
    }
    return $fetch<GetMeResponse>('/api/auth/me', {
      method: 'GET',
      headers
    })
  },
  /**
   * Sign out the user by hitting local route (which calls backend /auth/logout).
   */
  async signOut(): Promise<void> {
    return $fetch<void>('/api/auth/logout', {
      method: 'POST'
    })
  },
  /**
   * Retrieve profile of the currently logged-in customer.
   */
  async getCurrentUser(cookieHeader?: string): Promise<GetCurrentUserResponse> {
    const headers: Record<string, string> = {}
    
    // If running on server-side (SSR), forward request cookies to the backend API
    if (cookieHeader) {
      headers['cookie'] = cookieHeader
    }
    return bootstrapConfig.fetchApi<GetCurrentUserResponse>('/users/me/', {
      method: 'GET',
      headers
    })
  },
  /**
   * Retrieve profile of the currently authenticated customer.
   */
  async getProfile(cookieHeader?: string): Promise<GetProfileResponse> {
    const headers: Record<string, string> = {}
    if (cookieHeader) {
      headers['cookie'] = cookieHeader
    }
    return bootstrapConfig.fetchApi<GetProfileResponse>('/profile', {
      method: 'GET',
      headers
    })
  },
  /**
   * Update profile of the currently authenticated customer.
   */
  async updateProfile(data: UpdateProfileRequest, cookieHeader?: string): Promise<GetProfileResponse> {
    const headers: Record<string, string> = {}
    if (cookieHeader) {
      headers['cookie'] = cookieHeader
    }
    return bootstrapConfig.fetchApi<GetProfileResponse>('/profile', {
      method: 'PUT',
      headers,
      body: data
    })
  }
}
