// app/utils/bootstrap.ts
import { triggerSessionExpired } from '~/composables/useSessionState'

export const bootstrapConfig = {
  /**
   * On the client side, routes through Nuxt's server proxy (/api/core)
   * to avoid CORS preflight issues entirely.
   * On the server side (SSR), calls the backend directly (no CORS).
   */
  getApiBaseUrl(): string {
    if (import.meta.client) {
      // Same-origin proxy path — no cross-origin, no preflight
      return '/api/core'
    }

    // Server-side: call Go backend directly
    try {
      const config = useRuntimeConfig()
      return (config.public.serviceCoreApiUrl as string) || 'http://localhost:7129'
    } catch {
      return 'http://localhost:7129'
    }
  },
  
  /**
   * Helper utility to fetch from the service core API endpoint.
   * Includes credentials so cookies flow through the proxy.
   */
  async fetchApi<T>(endpoint: string, options?: Parameters<typeof $fetch>[1]): Promise<T> {
    const baseUrl = this.getApiBaseUrl().replace(/\/$/, '')
    const cleanEndpoint = endpoint.replace(/^\//, '')
    const url = `${baseUrl}/${cleanEndpoint}`
    
    try {
      return await $fetch<T>(url, {
        credentials: 'include',
        ...options
      })
    } catch (err: any) {
      if (import.meta.client && (err.status === 401 || err.status === 403)) {
        triggerSessionExpired()
      }
      throw err
    }
  }
}
