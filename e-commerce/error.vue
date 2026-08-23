<script setup lang="ts">
// error.vue — Nuxt global error page
// Rendered by Nuxt when a fatal/unhandled error occurs (404, 500, createError, etc.)
// NuxtLayout is NOT available here, so the branded shell is built inline.
//
// Docs: https://nuxt.com/docs/guide/directory-structure/error

import type { NuxtError } from '#app'

const props = defineProps<{
  error: NuxtError
}>()

// Resolve the status code — Nuxt may surface it as number or string
const statusCode = computed<404 | 500 | number>(() => {
  const code = Number(props.error?.statusCode ?? 500)
  return code
})

// Clear the error and navigate
const handleClear = async (path = '/') => {
  await clearError({ redirect: path })
}
</script>

<template>
  <div class="min-h-screen flex flex-col font-sans bg-white text-gray-900 antialiased" style="font-family: 'Inter', system-ui, -apple-system, sans-serif;">

    <!-- Minimal branded header (NuxtLayout unavailable in error.vue) -->
    <header class="w-full border-b border-gray-100 bg-white">
      <div class="container-standard flex items-center h-16">
        <a href="/" class="inline-block transition-transform hover:scale-105 duration-200">
          <img src="/images/logo.png" class="h-8 w-auto object-contain" alt="Chia Florist" />
        </a>
      </div>
    </header>

    <!-- Error body -->
    <main class="flex-grow">
      <CErrorDisplay
        :status-code="statusCode"
        :clear-on-navigate="true"
        @clear="handleClear()"
      />
    </main>

    <!-- Minimal footer -->
    <footer class="w-full border-t border-gray-100 py-6 text-center text-xs text-gray-400" style="font-family: 'Inter', system-ui, sans-serif;">
      <p>© {{ new Date().getFullYear() }} Chia Florist. All rights reserved.</p>
    </footer>

  </div>
</template>

<style scoped>
.container-standard {
  max-width: 80rem;
  margin-left: auto;
  margin-right: auto;
  padding-left: 2rem;
  padding-right: 2rem;
}

@media (min-width: 640px) {
  .container-standard {
    padding-left: 2.5rem;
    padding-right: 2.5rem;
  }
}
</style>
