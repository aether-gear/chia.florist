// app/utils/errorMessages.ts
// Centralized error message mapping for Issue #286.
// Translates raw backend error messages / HTTP status codes into
// consistent, user-friendly strings — no technical details exposed to users.

/**
 * Extract the raw backend error message string from a fetch error,
 * falling back through the chain: err.data?.message → err.message → ''
 */
const getRawMessage = (err: unknown): string => {
  if (!err || typeof err !== 'object') return ''
  const e = err as Record<string, any>
  return (e.data?.message ?? e.data?.error ?? e.message ?? '').toString().toLowerCase().trim()
}

/**
 * Extract the HTTP status code from a fetch error.
 */
const getStatusCode = (err: unknown): number => {
  if (!err || typeof err !== 'object') return 0
  const e = err as Record<string, any>
  return Number(e.status ?? e.statusCode ?? e.response?.status ?? 0)
}

/**
 * Maps raw backend error message strings (case-insensitive substring match)
 * to user-friendly strings.
 */
const MESSAGE_MAP: Array<{ match: string[]; friendly: string }> = [
  // Authentication errors
  {
    match: ['invalid credentials', 'wrong password', 'incorrect password', 'invalid password', 'wrong email'],
    friendly: 'Email atau kata sandi tidak valid. Silakan periksa kembali.'
  },
  {
    match: ['email not verified', 'not verified', 'account not verified'],
    friendly: 'Email Anda belum diverifikasi. Silakan periksa kotak masuk email Anda.'
  },
  {
    match: ['email already registered', 'email already exists', 'already registered', 'duplicate email', 'email taken', 'already in use', 'already exist'],
    friendly: 'Email ini sudah terdaftar. Silakan gunakan email lain atau masuk ke akun Anda.'
  },
  {
    match: ['username already', 'username taken', 'username exist'],
    friendly: 'Nama pengguna ini sudah digunakan. Silakan pilih nama pengguna lain.'
  },
  {
    match: ['unauthorized', 'unauthenticated', 'session expired', 'token expired', 'invalid token', 'no active session'],
    friendly: 'Sesi Anda telah berakhir. Silakan masuk kembali.'
  },
  {
    match: ['forbidden', 'access denied', 'insufficient permission', 'not allowed'],
    friendly: 'Anda tidak memiliki akses untuk melakukan tindakan ini.'
  },
  {
    match: ['otp invalid', 'otp expired', 'invalid otp', 'wrong otp', 'invalid code', 'wrong code', 'code expired', 'verification code'],
    friendly: 'Kode verifikasi tidak valid atau sudah kadaluarsa. Silakan coba lagi.'
  },
  // Resource errors
  {
    match: ['record not found', 'not found', 'no record', 'does not exist', 'product not found', 'item not found'],
    friendly: 'Data yang Anda cari tidak dapat ditemukan.'
  },
  {
    match: ['address not found', 'no address'],
    friendly: 'Alamat pengiriman tidak ditemukan. Silakan tambahkan alamat terlebih dahulu.'
  },
  // Stock / inventory
  {
    match: ['out of stock', 'no stock', 'stock is 0', 'unavailable'],
    friendly: 'Produk ini sedang tidak tersedia. Silakan cek kembali nanti.'
  },
  {
    match: ['insufficient stock', 'not enough stock', 'exceeds available stock', 'stock not enough'],
    friendly: 'Stok tidak mencukupi untuk jumlah yang diminta. Silakan kurangi jumlah pembelian.'
  },
  // Validation
  {
    match: ['validation failed', 'validation error', 'invalid input', 'invalid request body', 'invalid data', 'bad request'],
    friendly: 'Periksa kembali data yang Anda masukkan dan coba lagi.'
  },
  // Conflicts
  {
    match: ['conflict', 'duplicate key', 'already processed', 'already paid'],
    friendly: 'Terjadi konflik data. Silakan muat ulang halaman dan coba lagi.'
  },
  // Rate limiting
  {
    match: ['too many requests', 'rate limit', 'too many attempts'],
    friendly: 'Terlalu banyak percobaan. Silakan tunggu beberapa saat sebelum mencoba lagi.'
  },
  // Server errors
  {
    match: ['internal server error', 'server error', 'unexpected error', 'something went wrong', 'system error'],
    friendly: 'Terjadi kesalahan pada server kami. Silakan coba lagi dalam beberapa saat.'
  },
  // Checkout / payment
  {
    match: ['payment failed', 'payment error'],
    friendly: 'Pembayaran gagal diproses. Silakan coba metode pembayaran lain.'
  },
  {
    match: ['order not found', 'order expired'],
    friendly: 'Pesanan tidak ditemukan atau sudah kadaluarsa.'
  },
]

/**
 * Maps an HTTP status code to a friendly fallback message.
 */
const STATUS_MAP: Record<number, string> = {
  400: 'Permintaan tidak valid. Periksa kembali data yang Anda masukkan.',
  401: 'Sesi Anda telah berakhir. Silakan masuk kembali.',
  403: 'Anda tidak memiliki akses untuk melakukan tindakan ini.',
  404: 'Data yang Anda cari tidak dapat ditemukan.',
  408: 'Permintaan habis waktu. Silakan coba lagi.',
  409: 'Terjadi konflik data. Silakan muat ulang halaman dan coba lagi.',
  410: 'Data ini sudah tidak tersedia.',
  422: 'Data yang dikirim tidak valid. Silakan periksa kembali.',
  429: 'Terlalu banyak permintaan. Silakan tunggu sejenak sebelum mencoba lagi.',
  500: 'Terjadi kesalahan pada server kami. Silakan coba lagi dalam beberapa saat.',
  502: 'Server sedang tidak dapat dijangkau. Silakan coba lagi nanti.',
  503: 'Layanan sedang dalam pemeliharaan. Silakan coba lagi nanti.',
  504: 'Server tidak merespons dalam waktu yang cukup. Silakan coba lagi.',
}

/**
 * Translates any caught error into a user-friendly message.
 *
 * Priority order:
 *  1. Raw backend message string matched against MESSAGE_MAP
 *  2. HTTP status code matched against STATUS_MAP
 *  3. Generic fallback
 *
 * @param err - The caught error (FetchError, Error, or unknown)
 * @param fallback - Optional custom fallback if no mapping found
 */
export const mapErrorMessage = (err: unknown, fallback = 'Terjadi kesalahan. Silakan coba lagi.'): string => {
  const raw = getRawMessage(err)
  const status = getStatusCode(err)

  // 1. Try raw message mapping first (most specific)
  if (raw) {
    for (const entry of MESSAGE_MAP) {
      if (entry.match.some(keyword => raw.includes(keyword))) {
        return entry.friendly
      }
    }
  }

  // 2. Fall back to HTTP status code mapping
  if (status > 0) {
    if (STATUS_MAP[status]) return STATUS_MAP[status]
    // Generic 5xx
    if (status >= 500) return STATUS_MAP[500]
    // Generic 4xx
    if (status >= 400) return STATUS_MAP[400]
  }

  // 3. Final fallback
  return fallback
}

/**
 * Convenience: determine if an error is a network connectivity issue.
 */
export const isNetworkError = (err: unknown): boolean => {
  if (!err || typeof err !== 'object') return false
  const e = err as Record<string, any>
  const msg = (e.message ?? '').toString().toLowerCase()
  return msg.includes('network') || msg.includes('fetch') || msg.includes('econnrefused') || e.status === 0
}

/**
 * Convenience: determine if an error is a 404 Not Found.
 */
export const isNotFoundError = (err: unknown): boolean => {
  return getStatusCode(err) === 404
}

/**
 * Log the original error for debugging while suppressing it from users.
 * Use in catch blocks where only logging is needed.
 */
export const logError = (context: string, err: unknown): void => {
  if (import.meta.dev || import.meta.server) {
    console.error(`[${context}]`, err)
  }
}
