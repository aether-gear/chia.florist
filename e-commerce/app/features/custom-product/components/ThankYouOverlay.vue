<script setup lang="ts">
const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'new-design'): void
}>()
</script>

<template>
  <div v-if="show" class="dr-modal-backdrop">
    <div class="dr-modal ty-modal">
      <div class="ty-icon">🌸</div>
      <h2>Custom Board Added to Cart!</h2>
      <p>Your custom board design has been successfully saved and added to your shopping cart.</p>
      <div class="ty-actions">
        <NuxtLink to="/cart" class="primary-btn ty-btn">Go to Cart &amp; Checkout</NuxtLink>
        <button class="secondary-btn ty-btn" @click="emit('new-design')">Design Another Board</button>
      </div>
    </div>
  </div>
</template>
