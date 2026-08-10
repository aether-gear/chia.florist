<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import {
  useCustomDesign,
  useCustomCart,
  CanvasBoard,
  ToolPanel,
  ReviewModal,
  FinalizeChoiceOverlay,
  ThankYouOverlay
} from '~/features/custom-product'
import '~/features/custom-product/custom-product.css'

definePageMeta({ layout: false })
useHead({
  title: 'Chia Florist — Board Designer (v3.0)',
  meta: [{ name: 'description', content: 'Design your custom flower board with our interactive canvas designer.' }]
})

const design = useCustomDesign()
const { isAdding, addCustomDesignToCart } = useCustomCart()

const handleFinalize = () => {
  design.showFinalizeChoice.value = true
}

const handleOpenReview = () => {
  design.showFinalizeChoice.value = false
  design.showReview.value = true
}

const handleAddToCart = async () => {
  const payload = design.buildCustomDesignPayload(design.snapshotDataUrl.value)
  await addCustomDesignToCart(
    payload,
    design.snapshotDataUrl.value,
    design.totalPrice.value,
    design.physicalSize.value,
    design.upper.value.headerText
  )
  design.showReview.value = false
  design.showThankYou.value = true
}

const handleNewDesign = () => {
  design.resetDesign()
  design.showThankYou.value = false
}

onMounted(() => {
  design.loadDraft()
  design.updateScale()
  window.addEventListener('mousemove', design.onMouseMove)
  window.addEventListener('mouseup', design.onMouseUp)
  window.addEventListener('keydown', design.onKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', design.onMouseMove)
  window.removeEventListener('mouseup', design.onMouseUp)
  window.removeEventListener('keydown', design.onKeyDown)
})
</script>

<template>
  <div class="dr-root">
    <!-- ═══ NAVBAR ═══════════════════════════════════════════════════ -->
    <nav class="dr-nav">
      <NuxtLink to="/" class="dr-back">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        Home
      </NuxtLink>
      <div class="dr-nav-center">
        <span class="dr-brand">CHIA FLORIST</span>
        <span class="dr-dot">◆</span>
        <span class="dr-page-title">Board Designer v3.0</span>
      </div>
      <div class="dr-nav-right">
        <span class="dr-scale-chip">{{ Math.round(design.boardScale.value * 100) }}%</span>
        <button class="dr-random-btn" @click="design.randomizeDesign" title="Randomize Design">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="16 3 21 3 21 8"/><line x1="4" y1="20" x2="21" y2="3"/><polyline points="21 16 21 21 16 21"/><line x1="15" y1="15" x2="21" y2="21"/><line x1="4" y1="4" x2="9" y2="9"/>
          </svg>
          Randomize
        </button>
      </div>
    </nav>

    <!-- ═══ BODY (Client Only for Interactive Board Designer) ════════ -->
    <ClientOnly>
      <div class="dr-body">
        <CanvasBoard :design="design" />
        <ToolPanel :design="design" @finalize="handleFinalize" />
      </div>

      <!-- ═══ MODALS & OVERLAYS ═════════════════════════════════════════ -->
      <FinalizeChoiceOverlay
        :show="design.showFinalizeChoice.value"
        @close="design.showFinalizeChoice.value = false"
        @review="handleOpenReview"
      />
      <ReviewModal
        :design="design"
        :is-adding="isAdding"
        @close="design.showReview.value = false"
        @add-to-cart="handleAddToCart"
      />
      <ThankYouOverlay
        :show="design.showThankYou.value"
        @new-design="handleNewDesign"
      />

      <template #fallback>
        <div class="dr-body flex items-center justify-center min-h-[400px]">
          <div class="flex flex-col items-center gap-3">
            <div class="dr-spinner"></div>
            <p class="text-xs text-gray-500 font-semibold">Loading Board Designer v3.0…</p>
          </div>
        </div>
      </template>
    </ClientOnly>
  </div>
</template>
