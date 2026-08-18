import { setCookie, readBody } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/signin`
  const body = await readBody(event)
  const { rememberMe, ...backendBody } = body || {}

  try {
    const response = await $fetch.raw(backendUrl, {
      method: 'POST',
      body: backendBody,
      headers: {
        'content-type': 'application/json'
      }
    })

    const data = response._data
    const setCookieHeaders = response.headers.getSetCookie ? response.headers.getSetCookie() : [response.headers.get('set-cookie')].filter(Boolean) as string[]
    
    for (const setCookieHeader of setCookieHeaders) {
      if (setCookieHeader) {
        const parts = setCookieHeader.split(';')
        const [cookieNameValue, ...options] = parts
        const eqIdx = cookieNameValue.indexOf('=')
        const name = cookieNameValue.substring(0, eqIdx).trim()
        const value = cookieNameValue.substring(eqIdx + 1).trim()

        let expires: Date | undefined
        let path = '/'
        let sameSite: 'lax' | 'strict' | 'none' | undefined = 'lax'

        for (const option of options) {
          const [optName, optVal] = option.split('=').map(s => s.trim())
          if (optName.toLowerCase() === 'expires') {
            expires = new Date(optVal)
          } else if (optName.toLowerCase() === 'path') {
            path = optVal
          } else if (optName.toLowerCase() === 'samesite') {
            sameSite = optVal.toLowerCase() as any
          }
        }

        // Enforce 30-day expiration when rememberMe is enabled
        const thirtyDays = new Date(Date.now() + 30 * 24 * 3600 * 1000)
        const cookieExpires = rememberMe
          ? (expires && expires.getTime() > thirtyDays.getTime() ? expires : thirtyDays)
          : undefined

        // Strip Secure flag in local development
        setCookie(event, name, value, {
          path,
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite,
          expires: cookieExpires
        })
      }
    }

    // Set remember_me cookie for the client
    setCookie(event, 'remember_me', rememberMe ? 'true' : 'false', {
      path: '/',
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      expires: rememberMe ? new Date(Date.now() + 30 * 24 * 3600 * 1000) : undefined
    })

    return data
  } catch (err: any) {
    event.node.res.statusCode = err.status || 500
    return err.data || { message: err.message }
  }
})
