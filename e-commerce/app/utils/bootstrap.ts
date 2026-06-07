// app/utils/bootstrap.ts

export const bootstrapConfig = {
  /**
   * Retrieves the core service API base URL from runtime configurations,
   * falling back to http://localhost:8000 if not available.
   */
  getApiBaseUrl(): string {
    try {
      const config = useRuntimeConfig()
      return (config.public.serviceCoreApiUrl as string) || 'http://localhost:8000'
    } catch {
      return 'http://localhost:8000'
    }
  },

  /**
   * Helper utility to fetch from the service core API endpoint.
   */
  async fetchApi<T>(endpoint: string, options?: Parameters<typeof $fetch>[1]): Promise<T> {
    const baseUrl = this.getApiBaseUrl().replace(/\/$/, '')
    const cleanEndpoint = endpoint.replace(/^\//, '')
    const url = `${baseUrl}/${cleanEndpoint}`
    
    return $fetch<T>(url, options)
  }
}
