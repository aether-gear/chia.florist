// Memory cache for backend shops to avoid repeating GET /shops on every product request
let cachedShops: { data: any[]; expiresAt: number } | null = null
const uuidToSlugCache = new Map<string, { slug: string; expiresAt: number }>()
const CACHE_TTL_MS = 5 * 60 * 1000 // 5 minutes

export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const config = useRuntimeConfig()
  const backendBaseUrl = config.public.serviceCoreApiUrl || 'http://127.0.0.1:7129'

  let slug = id || ''

  // Virtual handler for client-side interactive custom simulator
  if (slug === 'custom') {
    return {
      id: 'custom',
      sku: 'CUSTOM-BOARD-V3',
      name: 'Custom Board Simulator',
      slug: 'custom',
      status: 'active',
      is_available: true,
      price: 150000,
      stock: 999,
      description: 'Interactive Custom Flower Board Simulator',
      banner: { thumbnail: '/images/custom-preview.png', preview: '/images/custom-preview.png', detail: '/images/custom-preview.png' },
      gallery: [],
      availability: [],
      shop_id: '99ef0062-1040-4574-a4be-0123abce5670'
    }
  }

  // 2. If parameter is a UUID string, resolve it to slug via GET /products?id={uuid} with caching
  const isUuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(slug)
  if (isUuid) {
    const cached = uuidToSlugCache.get(slug)
    if (cached && Date.now() < cached.expiresAt) {
      slug = cached.slug
    } else {
      try {
        const listRes: any = await $fetch(`${backendBaseUrl}/products?id=${slug}`)
        if (listRes && Array.isArray(listRes.products) && listRes.products.length > 0 && listRes.products[0].slug) {
          slug = listRes.products[0].slug
          uuidToSlugCache.set(id!, { slug, expiresAt: Date.now() + CACHE_TTL_MS })
        }
      } catch (lookupErr) {
        console.error(`Failed to resolve product slug for UUID ${slug}:`, lookupErr)
      }
    }
  }

  // 3. Fetch full product details strictly from backend by slug
  let product: any = null
  try {
    product = await $fetch(`${backendBaseUrl}/products/${slug}`)
  } catch (err: any) {
    throw createError({
      statusCode: err.status || 500,
      statusMessage: err.statusText || 'Backend Error',
      message: err.data?.message || err.message
    })
  }

  // 4. Resolve shop_id using backend shops API based on availability with caching
  let shopId = '333f6432-a01c-412f-99f4-0f08ca0d8eb1' // Default fallback
  try {
    let shopsList: any[] = []
    if (cachedShops && Date.now() < cachedShops.expiresAt) {
      shopsList = cachedShops.data
    } else {
      const shopsRes: any = await $fetch(`${backendBaseUrl}/shops`)
      if (shopsRes && Array.isArray(shopsRes.shops)) {
        shopsList = shopsRes.shops
        cachedShops = { data: shopsList, expiresAt: Date.now() + CACHE_TTL_MS }
      }
    }

    if (shopsList.length > 0 && Array.isArray(product.availability) && product.availability.length > 0) {
      const sortedAvail = [...product.availability].sort((a, b) => b.stock - a.stock)
      const highestStockShopSlug = sortedAvail[0].name
      const matchedShop = shopsList.find((s: any) => s.slug === highestStockShopSlug)
      if (matchedShop) {
        shopId = matchedShop.id
      }
    }
  } catch (shopErr) {
    console.error('Failed to resolve shop ID dynamically from shops API:', shopErr)
  }

  return {
    ...product,
    shop_id: shopId
  }
})
