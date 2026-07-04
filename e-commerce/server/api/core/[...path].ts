import { defineEventHandler, proxyRequest, setCookie, removeResponseHeader } from "h3";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const targetBase =
    process.env.SERVICE_CORE_API_URL ||
    config.public.serviceCoreApiUrl ||
    "http://127.0.0.1:7129";
  const subPath = event.context.params?.path || "";

  const targetUrl = `${targetBase.replace(/\/$/, "")}/${subPath.replace(/^\//, "")}`;
  console.log(
    `[PROXY] Proxying request from ${event.node.req.url} to ${targetUrl}`,
  );

  try {
    return await proxyRequest(event, targetUrl, {
      onResponse(event, response) {
        const setCookieHeaders = response.headers.getSetCookie 
          ? response.headers.getSetCookie() 
          : [response.headers.get('set-cookie')].filter(Boolean) as string[]

        if (setCookieHeaders.length > 0) {
          removeResponseHeader(event, 'set-cookie')

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

              setCookie(event, name, value, {
                path,
                httpOnly: true,
                secure: process.env.NODE_ENV === 'production',
                sameSite,
                expires
              })
            }
          }
        }
      }
    });
  } catch (err: any) {
    console.error(`[PROXY ERROR] Failed to proxy to ${targetUrl}:`, err);
    throw err;
  }
});
