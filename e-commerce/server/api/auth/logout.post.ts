import { deleteCookie, getHeader } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/logout`
  const cookies = getHeader(event, 'cookie')

  try {
    await $fetch(backendUrl, {
      method: 'POST',
      headers: cookies ? { cookie: cookies } : {}
    })
  } catch (err: any) {
    console.error('Backend logout failed:', err.data?.message || err.message)
  }

  // Clear cookie locally
  deleteCookie(event, 'chast')
  return { message: 'logout success' }
})
