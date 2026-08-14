<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'solid' | 'text'
    size?: 'pill' | 'sm' | 'md' | 'lg' | 'auth'
    type?: 'button' | 'submit' | 'reset'
    loading?: boolean
    isLoading?: boolean
    disabled?: boolean
    fullWidth?: boolean
    to?: string
    href?: string
    target?: string
    rel?: string
  }>(),
  {
    variant: 'primary',
    size: 'md',
    type: 'button',
    loading: false,
    isLoading: false,
    disabled: false,
    fullWidth: false
  }
)

const isButtonLoading = computed(() => props.loading || props.isLoading)
const isButtonDisabled = computed(() => props.disabled || isButtonLoading.value)

const baseClasses = 'inline-flex items-center justify-center font-medium transition-all duration-200 ease-in-out cursor-pointer select-none outline-none active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100'

const widthClasses = computed(() => (props.fullWidth ? 'w-full flex' : ''))

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'pill':
      return 'rounded-full px-3.5 py-1.5 text-xs font-semibold'
    case 'sm':
      return 'rounded-xl px-3.5 py-1.5 text-xs font-semibold'
    case 'lg':
    case 'auth':
      return 'rounded-xl px-4 py-3 text-sm font-bold'
    case 'md':
    default:
      return 'rounded-xl px-5 py-2.5 text-sm font-semibold'
  }
})

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'secondary':
      return 'bg-[#1b4332] hover:bg-[#143326] text-white font-bold shadow-xs hover:shadow focus:ring-2 focus:ring-[#1b4332]/30'
    case 'outline':
      return 'border border-gray-200 text-gray-700 hover:bg-gray-50 font-semibold focus:ring-2 focus:ring-gray-200'
    case 'ghost':
    case 'text':
      return 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 font-semibold'
    case 'danger':
      return 'bg-red-500 hover:bg-red-600 text-white font-semibold shadow-xs focus:ring-2 focus:ring-red-500/30'
    case 'solid':
    case 'primary':
    default:
      return 'bg-[#4ade80] hover:bg-[#34d399] text-[#245842] font-bold shadow-xs hover:shadow focus:ring-2 focus:ring-[#4ade80]/30'
  }
})
</script>

<template>
  <!-- Internal Nuxt Link -->
  <NuxtLink
    v-if="to"
    :to="to"
    :class="[baseClasses, sizeClasses, variantClasses, widthClasses, { 'pointer-events-none opacity-50': isButtonDisabled }]"
  >
    <svg
      v-if="isButtonLoading"
      class="animate-spin -ml-1 mr-2 h-4 w-4 text-current"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </NuxtLink>

  <!-- External Link -->
  <a
    v-else-if="href"
    :href="href"
    :target="target"
    :rel="rel"
    :class="[baseClasses, sizeClasses, variantClasses, widthClasses, { 'pointer-events-none opacity-50': isButtonDisabled }]"
  >
    <svg
      v-if="isButtonLoading"
      class="animate-spin -ml-1 mr-2 h-4 w-4 text-current"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </a>

  <!-- Standard Button -->
  <button
    v-else
    :type="type"
    :disabled="isButtonDisabled"
    :class="[baseClasses, sizeClasses, variantClasses, widthClasses]"
  >
    <svg
      v-if="isButtonLoading"
      class="animate-spin -ml-1 mr-2 h-4 w-4 text-current"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </button>
</template>

