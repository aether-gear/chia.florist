<script setup lang="ts">
// Default layout
</script>

<template>
  <div class="min-h-screen flex flex-col font-sans">
    <CNavbar />
    <main class="flex-grow">
      <slot />
    </main>
    <CFooter />
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}
</style>
