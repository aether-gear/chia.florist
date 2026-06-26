// app/types/auth.ts

export interface SignUpRequest {
  name: string
  username: string
  email: string
  password: string
  phone?: string
}

export interface SignUpResponse {
  message: string
  challenge_id: string
}

export interface VerifyRequest {
  user_agent?: string
  ip_address?: string
  challenge_id: string
  otp: number
}

export interface VerifyResponse {
  message: string
}

export interface SignInRequest {
  email: string
  password: string
  user_agent?: string
  ip_address?: string
}

export interface SignInResponse {
  message: string
}

export interface UserMe {
  id: string
  name: string
  username: string
  email?: string
  phone: string
  last_login_at: string | null
  avatarUrl?: string | null
}

export interface GetCurrentUserResponse {
  me: UserMe
}

export interface GetMeResponse {
  account_id: string
  account_type: string
  is_authenticated: boolean
}
