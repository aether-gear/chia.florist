import { deleteCookie, getHeader } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/logout`
  const cookies = getHeader(event, 'cookie')

  try {
    await $fetch(backendUrl, {
      method: 'POST',
      headers: {
        'X-Account-Type': 'customer',
        ...(cookies ? { cookie: cookies } : {})
      }
    })
  } catch (err: any) {
    console.error('Backend logout failed:', err.data?.message || err.message)
  }

  // Clear customer cookies locally
  deleteCookie(event, 'chast')
  deleteCookie(event, 'malkist')
  return { message: 'logout success' }
})
