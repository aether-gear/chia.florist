import { deleteCookie } from 'h3'

export default defineEventHandler((event) => {
  deleteCookie(event, 'chast')
  return { message: 'logout success' }
})
