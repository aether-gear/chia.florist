<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  variant?: 'solid' | 'outline' | 'text'
  size?: 'sm' | 'md' | 'lg'
  to?: string
}>()

const baseClasses = 'inline-flex items-center justify-center font-medium transition-all duration-300 ease-in-out cursor-pointer'

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return 'px-4 py-2 text-sm'
    case 'lg': return 'px-8 py-4 text-lg'
    default: return 'px-6 py-3 text-base'
  }
})

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'outline': return 'border-2 border-white/80 text-white hover:bg-white hover:text-black'
    case 'text': return 'text-white/80 hover:text-white hover:underline'
    default: return 'bg-accent text-white-base hover-bg-accent-deep shadow-lg hover:shadow-xl'
  }
})
</script>

<template>
  <NuxtLink v-if="to" :to="to" :class="[baseClasses, sizeClasses, variantClasses]">
    <slot />
  </NuxtLink>
  <button v-else :class="[baseClasses, sizeClasses, variantClasses]">
    <slot />
  </button>
</template>
