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

  // 2. Fetch full product details strictly from backend by slug
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

  // 3. Resolve shop_id using backend shops API based on availability
  let shopId = '333f6432-a01c-412f-99f4-0f08ca0d8eb1' // Default fallback
  try {
    const shopsRes: any = await $fetch(`${backendBaseUrl}/shops`)
    if (shopsRes && Array.isArray(shopsRes.shops) && Array.isArray(product.availability) && product.availability.length > 0) {
      const sortedAvail = [...product.availability].sort((a, b) => b.stock - a.stock)
      const highestStockShopSlug = sortedAvail[0].name
      const matchedShop = shopsRes.shops.find((s: any) => s.slug === highestStockShopSlug)
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
