import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2026-04-16",

  css: [
    "~/assets/css/main.css",
    "~/features/custom-product/custom-product.css"
  ],

  modules: [],

  devServer: {
    host: process.env.HOST || "0.0.0.0",
    port: parseInt(process.env.PORT || "4000"),
  },

  runtimeConfig: {
    public: {
      serviceCoreApiUrl:
        process.env.NUXT_PUBLIC_SERVICE_CORE_API_URL ||
        process.env.SERVICE_CORE_API_URL ||
        "http://127.0.0.1:7129",
      supabaseUrl:
        process.env.SUPABASE_CHIA_URL ||
        "",
      supabaseKey:
        process.env.SUPABASE_CHIA_KEY ||
        process.env.SUPABASE_CHIA_KEY ||
        "",
    },
  },

  vite: {
    plugins: [tailwindcss()],
  },
  app: {
    pageTransition: { name: "page", mode: "out-in" },
  },
});
