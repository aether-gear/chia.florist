import { setCookie, readBody } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/verify`
  const body = await readBody(event)

  try {
    const response = await $fetch.raw(backendUrl, {
      method: 'POST',
      body,
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

        // Strip Secure flag in local development
        setCookie(event, name, value, {
          path,
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite,
          expires
        })
      }
    }

    return data
  } catch (err: any) {
    event.node.res.statusCode = err.status || 500
    return err.data || { message: err.message }
  }
})
