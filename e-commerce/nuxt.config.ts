import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2026-04-16',

  css: ['~/assets/css/main.css'],

  // PASTIKAN BAGIAN MODULES KOSONG ATAU HAPUS NAMA @nuxtjs/tailwindcss
  modules: [], 

  devServer: {
    host: process.env.HOST || '0.0.0.0',
    port: parseInt(process.env.PORT || '4000')
  },

  runtimeConfig: {
    public: {
      serviceCoreApiUrl: process.env.SERVICE_CORE_API_URL || 'http://127.0.0.1:7129'
    }
  },

  // Proxy all /api/core/** requests to the Go backend so the browser
  // never makes a cross-origin call (eliminates CORS preflight entirely).
  routeRules: {
    '/api/core/**': {
      proxy: `${process.env.SERVICE_CORE_API_URL || 'http://127.0.0.1:7129'}/**`
    }
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  app: {
    pageTransition: { name: 'page', mode: 'out-in' }
  }
})