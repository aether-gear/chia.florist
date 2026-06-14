import { Client } from 'pg'

export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const config = useRuntimeConfig()
  const backendUrl = `${config.public.serviceCoreApiUrl}/products/${id}`

  // 1. Fetch from Go backend
  let product: any = null
  try {
    product = await $fetch(backendUrl)
  } catch (err: any) {
    throw createError({
      statusCode: err.status || 500,
      statusMessage: err.statusText || 'Backend Error',
      message: err.data?.message || err.message
    })
  }

  // 2. Query database for shop_id
  let shopId = '333f6432-a01c-412f-99f4-0f08ca0d8eb1' // Default fallback (Chia Cipinang)
  const client = new Client({
    user: 'postgres.mqolpawlannysqjokzoq',
    host: 'aws-1-ap-northeast-2.pooler.supabase.com',
    database: 'postgres',
    password: 'Chia.Florist@21',
    port: 6543,
    ssl: { rejectUnauthorized: false }
  })

  try {
    await client.connect()
    const res = await client.query(
      'SELECT shop_id FROM inventory WHERE product_id = $1 ORDER BY stock DESC LIMIT 1;',
      [id]
    )
    if (res.rows.length > 0) {
      shopId = res.rows[0].shop_id
    }
  } catch (dbErr) {
    console.error('Failed to query inventory from DB:', dbErr)
  } finally {
    try {
      await client.end()
    } catch {}
  }

  return {
    ...product,
    shop_id: shopId
  }
})
