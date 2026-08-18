<script setup lang="ts">
import { computed, unref } from 'vue'
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

// Safely unwrap refs regardless of whether design prop fields are Refs or unwrapped
const showReviewModal = computed(() => Boolean(unref(props.design.showReview)))
const snapshotLoading = computed(() => Boolean(unref(props.design.snapshotLoading)))
const snapshotDataUrl = computed(() => unref(props.design.snapshotDataUrl) || '')
const physicalSizeVal = computed(() => unref(props.design.physicalSize) || 'medium')
const upperObj = computed(() => unref(props.design.upper) || {})
const lowerObj = computed(() => unref(props.design.lower) || {})
const borderObj = computed(() => unref(props.design.border) || {})
const topCrestObj = computed(() => unref(props.design.topCrest) || {})
const bottomCrestObj = computed(() => unref(props.design.bottomCrest) || {})
const imgElements = computed(() => unref(props.design.imgElements) || [])
const brushElements = computed(() => unref(props.design.brushElements) || [])
const totalPriceVal = computed(() => unref(props.design.totalPrice) || 200000)

const lowerFirstLine = computed(() => {
  const text = lowerObj.value.bodyText
  if (!text || typeof text !== 'string') return '(No Body)'
  return text.split('\n')[0] || '(No Body)'
})
</script>

<template>
  <div v-if="showReviewModal" class="dr-modal-backdrop" @click.self="emit('close')">
    <div class="dr-modal">
      <div class="dr-modal-head">
        <h2>Design Review &amp; Specs</h2>
        <button class="dr-modal-close" @click="emit('close')">×</button>
      </div>

      <div class="dr-modal-body">
        <div class="rev-grid">
          <!-- Left: generated snapshot -->
          <div class="rev-preview-col">
            <div v-if="snapshotLoading" class="rev-loading">
              <div class="dr-spinner"></div>
              <p>Rendering high-res preview…</p>
            </div>
            <img v-else-if="snapshotDataUrl" :src="snapshotDataUrl" class="rev-snapshot-img" alt="Board Preview"/>
            <div v-else class="rev-no-snap">Preview unavailable</div>
          </div>

          <!-- Right: specifications -->
          <div class="rev-specs-col">
            <h3>Custom Board Specifications</h3>
            <div class="rev-spec-row">
              <span class="lbl">Physical Size:</span>
              <span class="val">{{ SIZES.find(s => s.id === physicalSizeVal)?.label }} ({{ SIZES.find(s => s.id === physicalSizeVal)?.desc }})</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Upper Section:</span>
              <span class="val">{{ upperObj.headerText || '(No Header)' }} · {{ upperObj.bgColor || '#c0392b' }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Lower Section:</span>
              <span class="val">{{ lowerFirstLine }} · {{ lowerObj.bgColor || '#1a3a5c' }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Border:</span>
              <span class="val">{{ borderObj.style || 'solid' }} · {{ borderObj.width || 12 }}px · {{ borderObj.color || '#f5c842' }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Decorations:</span>
              <span class="val">Top: {{ topCrestObj.enabled ? topCrestObj.style : 'None' }} / Bottom: {{ bottomCrestObj.enabled ? bottomCrestObj.style : 'None' }}</span>
            </div>
            <div class="rev-spec-row">
              <span class="lbl">Elements:</span>
              <span class="val">{{ imgElements.length }} image(s), {{ brushElements.length }} stamped flower(s)</span>
            </div>

            <div class="rev-price-box">
              <span>Total Calculated Price:</span>
              <span class="rpb-total">{{ formatRupiah(totalPriceVal) }}</span>
            </div>

            <button class="primary-btn rev-add-btn" :disabled="isAdding" @click="emit('add-to-cart')">
              <span v-if="isAdding" class="dr-spinner sm"></span>
              <span v-else>Confirm &amp; Add to Cart — {{ formatRupiah(totalPriceVal) }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
