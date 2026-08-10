<script setup lang="ts">
import { useCart } from '~/composables/useCart'
import { SIZES } from '../constants'

const props = defineProps<{
  design: ReturnType<typeof import('../useCustomDesign').useCustomDesign>
  isAdding: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'add-to-cart'): void
}>()

const { formatRupiah } = useCart()
</script>

<template>
  <div v-if="design.showReview" class="dr-modal-backdrop" @click.self="emit('close')">
    <div class="dr-modal">
      <div class="dr-modal-head">
        <h2>Design Review &amp; Specs</h2>
        <button class="dr-modal-close" @click="emit('close')">×</button>
      </div>

      <div class="dr-modal-body">
        <div class="rev-grid">
          <!-- Left: generated snapshot -->
          <div class="rev-preview-col">
            <div v-if="design.snapshotLoading" class="rev-loading">
              <div class="dr-spinner"></div>
              <p>Rendering high-res preview…</p>
            </div>
            <img v-else-if="design.snapshotDataUrl" :src="design.snapshotDataUrl" class="rev-snapshot-img" alt="Board Preview"/>
            <div v-else class="rev-no-snap">Preview unavailable</div>
          </div>

          <!-- Right: specifications -->
          <div class="rev-specs-col">
            <h3>Custom Board Specifications</h3>
            <div class="rev-spec-row">
              <span class="lbl">Physical Size:</span>
              <span class="val">{{ SIZES.find(s => s.id === design.physicalSize)?.label }} ({{ SIZES.find(s => s.id === design.physicalSize)?.desc }})</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Upper Section:</span>
              <span class="val">{{ design.upper.headerText || '(No Header)' }} · {{ design.upper.bgColor }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Lower Section:</span>
              <span class="val">{{ design.lower.bodyText.split('\n')[0] || '(No Body)' }} · {{ design.lower.bgColor }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Border:</span>
              <span class="val">{{ design.border.style }} · {{ design.border.width }}px · {{ design.border.color }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Decorations:</span>
              <span class="val">Top: {{ design.topCrest.enabled ? design.topCrest.style : 'None' }} / Bottom: {{ design.bottomCrest.enabled ? design.bottomCrest.style : 'None' }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Elements:</span>
              <span class="val">{{ design.imgElements.length }} image(s), {{ design.brushElements.length }} stamped flower(s)</span>
            </div>

            <div class="rev-price-box">
              <span>Total Calculated Price:</span>
              <span class="rpb-total">{{ formatRupiah(design.totalPrice) }}</span>
            </div>

            <button class="primary-btn rev-add-btn" :disabled="isAdding" @click="emit('add-to-cart')">
              <span v-if="isAdding" class="dr-spinner sm"></span>
              <span v-else>Confirm &amp; Add to Cart — {{ formatRupiah(design.totalPrice) }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
