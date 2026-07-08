import { readBody } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/forgot-password/verify`
  const body = await readBody(event)

  try {
    const response = await $fetch.raw(backendUrl, {
      method: 'POST',
      body,
      headers: {
        'content-type': 'application/json'
      }
    })
    return response._data
  } catch (err: any) {
    event.node.res.statusCode = err.status || 500
    return err.data || { message: err.message }
  }
})
