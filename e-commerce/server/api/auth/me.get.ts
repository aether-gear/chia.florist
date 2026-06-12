import { getHeader } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/me`
  const cookies = getHeader(event, 'cookie')

  try {
    const data = await $fetch(backendUrl, {
      method: 'GET',
      headers: cookies ? { cookie: cookies } : {}
    })
    return data
  } catch (err: any) {
    return {
      account_id: "",
      account_type: "",
      is_authenticated: false,
      message: err.data?.message || err.message || 'Unauthorized'
    }
  }
})
