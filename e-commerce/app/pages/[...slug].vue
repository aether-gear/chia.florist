<script setup lang="ts">
// pages/[...slug].vue — Catch-all page for unknown client-side routes
// Renders inside the default layout (navbar + footer) and shows a 404 error display.
// This handles the case where a user navigates to a non-existent route on the client side.

definePageMeta({
  layout: 'default'
})

useHead({
  title: '404 — Halaman Tidak Ditemukan | Chia Florist',
  meta: [
    {
      name: 'robots',
      content: 'noindex'
    }
  ]
})

// Set the HTTP response status code to 404 for SSR requests hitting this catch-all
if (import.meta.server) {
  const event = useRequestEvent()
  if (event) {
    setResponseStatus(event, 404)
  }
}
</script>

<template>
  <CErrorDisplay :status-code="404" />
</template>
