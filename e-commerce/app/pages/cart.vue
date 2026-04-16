<script setup lang="ts">
// View logic handled by a viewmodel (e.g., useCartViewModel)
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-20">
    <h1 class="text-4xl font-bold font-serif mb-8 text-[#1b4332]">Your Cart</h1>
    <div class="bg-white p-8 shadow rounded-lg text-center">
      <p class="text-gray-500 mb-6">Your cart is currently empty.</p>
      <CButton variant="solid" to="/">Continue Shopping</CButton>
    </div>
  </div>
</template>
