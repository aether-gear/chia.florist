<script setup lang="ts">
const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'review'): void
}>()
</script>

<template>
  <div v-if="show" class="dr-modal-backdrop" @click.self="emit('close')">
    <div class="dr-modal choice-modal">
      <div class="dr-modal-head">
        <h2>What would you like to do next?</h2>
        <button class="dr-modal-close" @click="emit('close')">×</button>
      </div>
      <div class="choice-body">
        <p class="choice-subtitle">Your board design is ready. You can proceed to review your design specifications or continue customizing.</p>
        <div class="choice-actions">
          <button class="choice-btn primary" @click="emit('review')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M5 12l5 5L20 7"/></svg>
            Review Specs &amp; Add to Cart
          </button>
          <button class="choice-btn secondary" @click="emit('close')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
            Keep Customizing
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
