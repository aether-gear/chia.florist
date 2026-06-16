// app/utils/formatter.ts

/**
 * Memformat angka nominal mentah Rupiah dari database Supabase/Golang 
 * menjadi string mata uang Indonesia yang rapi: Rp70.000
 */
export const formatRupiah = (rawPrice: number): string => {
  // Amankan jika input bukan angka atau bernilai null/undefined
  const validPrice = rawPrice || 0

  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(validPrice)
}