import { sendRedirect } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/auth/google/login`

  // Redirect the browser directly to the backend's Google OAuth endpoint.
  // The browser must navigate here itself so it can interact with Google's
  // OAuth consent screen — a server-side fetch cannot do this.
  return sendRedirect(event, backendUrl, 302)
})
