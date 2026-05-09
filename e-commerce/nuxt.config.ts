import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2026-04-16',

  css: ['~/assets/css/main.css'],

  // PASTIKAN BAGIAN MODULES KOSONG ATAU HAPUS NAMA @nuxtjs/tailwindcss
  modules: [], 

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  app: {
    pageTransition: { name: 'page', mode: 'out-in' }
  }
})