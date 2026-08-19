// app/utils/searchMatcher.ts
import type { CatalogProduct } from '~/types/product'

/**
 * Keyword synonym groups to bridge English / Indonesian floral terms and occasions.
 */
const SYNONYM_GROUPS: string[][] = [
  ['custom', 'simulator', 'kustom', 'desain', 'design', 'interactive', 'game', 'papan kustom', 'custom board'],
  ['wedding', 'pernikahan', 'nikah', 'kawin', 'wed', 'pengantin', 'mempelai', 'happy wedding'],
  ['duka', 'duka cita', 'berduka', 'condolence', 'rip', 'belasungkawa', 'wafat', 'turut berduka cita'],
  ['wisuda', 'graduation', 'grad', 'lulus', 'sarjana', 'kelulusan', 'happy graduation'],
  ['birthday', 'ulang tahun', 'ultah', 'bday', 'hbd', 'milad', 'happy birthday'],
  ['opening', 'grand opening', 'peresmian', 'selamat & sukses', 'selamat dan sukses', 'congrats', 'congratulations', 'sukses', 'inauguration', 'toko baru', 'kantor baru', 'pelantikan']
]

/**
 * Evaluates whether a product matches a user search query with keyword & synonym matching.
 */
export function matchesProductSearch(product: CatalogProduct, query: string): boolean {
  const cleanQuery = query.trim().toLowerCase()
  if (!cleanQuery) return true

  const isCustomProduct = product.id === 'custom' || product.isCustomRoute || product.name.toLowerCase().includes('simulator')
  const cleanName = product.name.toLowerCase()
  const cleanTag = (product.tag || '').toLowerCase()
  const cleanSlug = (product.slug || '').toLowerCase()

  // 1. Single character search (e.g. 'a' through 'z'):
  // Matches directly against visible product name, tag, or slug (ignoring internal SKU prefixes)
  if (cleanQuery.length === 1) {
    const singleCharTarget = `${cleanName} ${cleanTag} ${cleanSlug}`
    return singleCharTarget.includes(cleanQuery)
  }

  // 2. Multi-character search:
  const isCustomQuery = ['custom', 'simulator', 'kustom', 'desain'].some(k => cleanQuery.includes(k))

  // If query is specifically targeting custom board/simulator, reject standard products
  if (isCustomQuery && !isCustomProduct && !cleanName.includes('custom')) {
    return false
  }

  // If product is the custom board simulator:
  // If query is specifically an occasion (e.g. wedding, duka, wisuda, birthday), custom should NOT appear.
  if (isCustomProduct) {
    const isOtherOccasionQuery = SYNONYM_GROUPS.some((group, idx) => idx > 0 && group.some(s => s === cleanQuery || cleanQuery.includes(s)))
    if (isOtherOccasionQuery) {
      return false
    }
  }

  // Search text for multi-character queries
  const coreProductText = `${cleanName} ${cleanTag} ${cleanSlug}`

  // Direct substring match on core identity
  if (coreProductText.includes(cleanQuery)) {
    return true
  }

  // Tokenized match (every word/token in query must match directly or via synonym)
  const queryTokens = cleanQuery.split(/\s+/).filter(Boolean)
  if (queryTokens.length === 0) return false

  const allTokensMatch = queryTokens.every(token => {
    // Direct match in core text
    if (coreProductText.includes(token)) return true

    // Match via synonym group
    for (const group of SYNONYM_GROUPS) {
      const tokenInGroup = group.some(s => s === token || s.includes(token) || token.includes(s))
      if (tokenInGroup) {
        const hasMatchingSynonymInProduct = group.some(s => coreProductText.includes(s))
        if (hasMatchingSynonymInProduct) return true
      }
    }

    return false
  })

  return allTokensMatch
}

/**
 * Filter a list of catalog products based on query relevance.
 */
export function filterCatalogProductsByQuery(
  products: CatalogProduct[],
  query: string,
  includeCustomSimulator = true,
  customCard?: CatalogProduct
): CatalogProduct[] {
  const cleanQuery = query.trim()
  
  let pool = [...products]

  // Add custom simulator card if provided and not already in list
  if (includeCustomSimulator && customCard && !pool.some(p => p.id === 'custom')) {
    pool = [customCard, ...pool]
  }

  if (!cleanQuery) {
    return pool
  }

  return pool.filter(p => matchesProductSearch(p, cleanQuery))
}
