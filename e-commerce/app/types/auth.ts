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
  otp: string
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

export interface UserProfile {
  customer_id: string
  user_id: string
  Name: string
  Username: string
  Phone: string
  AvatarURL: string | null
  LastLoginAt: string | null
  CreatedAt: string
  UpdatedAt: string | null
}

export interface GetProfileResponse {
  profile: UserProfile
}

export interface UpdateProfileRequest {
  name?: string
  phone?: string
  avatar_url?: string
}
