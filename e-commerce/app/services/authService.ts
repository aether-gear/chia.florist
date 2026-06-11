import type { 
  SignUpRequest, 
  SignUpResponse, 
  VerifyRequest, 
  VerifyResponse, 
  SignInRequest, 
  SignInResponse, 
  GetCurrentUserResponse 
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
    return bootstrapConfig.fetchApi<VerifyResponse>('/auth/verify', {
      method: 'POST',
      body: data
    })
  },
  /**
   * Sign in user via email and password.
   */
  async signIn(data: SignInRequest): Promise<SignInResponse> {
    return bootstrapConfig.fetchApi<SignInResponse>('/auth/signin', {
      method: 'POST',
      body: data
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
  }
}
