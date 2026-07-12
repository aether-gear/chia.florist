import { getHeader, setCookie, parseCookies } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/me`
  const cookies = getHeader(event, 'cookie')
  
  const requestCookies = parseCookies(event)
  const rememberMe = requestCookies['remember_me'] === 'true'

  try {
    const response = await $fetch.raw(backendUrl, {
      method: 'GET',
      headers: cookies ? { cookie: cookies } : {}
    })

    const data = response._data
    const setCookieHeaders = response.headers.getSetCookie 
      ? response.headers.getSetCookie() 
      : [response.headers.get('set-cookie')].filter(Boolean) as string[]

    const backendSetCookies = new Set<string>()

    for (const setCookieHeader of setCookieHeaders) {
      if (setCookieHeader) {
        const parts = setCookieHeader.split(';')
        const [cookieNameValue, ...options] = parts
        const eqIdx = cookieNameValue.indexOf('=')
        const name = cookieNameValue.substring(0, eqIdx).trim()
        const value = cookieNameValue.substring(eqIdx + 1).trim()

        backendSetCookies.add(name)

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

        setCookie(event, name, value, {
          path,
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite,
          expires: rememberMe ? expires : undefined
        })
      }
    }

    if (!rememberMe) {
      const cookiesToSanitize = ['chast', 'malkist', 'hotpot', 'ladle']
      for (const cookieName of cookiesToSanitize) {
        if (requestCookies[cookieName] && !backendSetCookies.has(cookieName)) {
          setCookie(event, cookieName, requestCookies[cookieName], {
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            sameSite: 'lax'
          })
        }
      }
    }

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
