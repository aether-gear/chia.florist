<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import {
  useCustomDesign,
  useCustomCart,
  CanvasBoard,
  ToolPanel,
  ReviewModal,
  FinalizeChoiceOverlay,
  ThankYouOverlay,
  TOOL_TABS
} from '~/features/custom-product'
import { useCart } from '~/composables/useCart'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import '~/features/custom-product/custom-product.css'

definePageMeta({ layout: false })
useHead({
  title: 'Chia Florist — Board Designer (v3.0)',
  meta: [{ name: 'description', content: 'Design your custom flower board with our interactive canvas designer.' }]
})

const design = useCustomDesign()
const { isAdding, addCustomDesignToCart } = useCustomCart()
const { formatRupiah } = useCart()
const globalAlert = useGlobalAlert()
const authVm = useAuthViewModel()
const isLoggedIn = useCookie('is_logged_in')

const handleFinalize = () => {
  design.showFinalizeChoice = true
}

const handleOpenReview = () => {
  design.showFinalizeChoice = false
  design.showReview = true
}

const handleAddToCart = async () => {
  // Free customization for all, but adding to cart requires being signed in
  if (isLoggedIn.value !== 'true' && !authVm.isAuthenticated.value) {
    design.saveDraft(true)
    design.showReview = false
    design.showFinalizeChoice = false
    globalAlert.showWarning(
      'Sign In Required',
      'To add your custom flower board to cart, please log in first. Your design has been saved to your drafts!',
      [
        { label: 'Sign In', onClick: () => navigateTo('/login?redirect=/products/custom') },
        { label: 'Keep Editing' }
      ]
    )
    return
  }

  const payload = design.buildCustomDesignPayload(design.snapshotDataUrl)
  await addCustomDesignToCart(
    payload,
    design.snapshotDataUrl,
    design.totalPrice,
    design.physicalSize,
    design.upper.headerText
  )
  design.showReview = false
  globalAlert.showSuccess(
    'Custom Board Added',
    'Your personalized flower board has been added to cart!',
    [
      { label: 'View Cart', onClick: () => navigateTo('/cart') },
      { label: 'Dismiss' }
    ]
  )
}

const handleNewDesign = () => {
  design.resetDesign()
  design.showThankYou = false
}

const pendingTargetUrl = ref('/')
const handleTryLeave = (targetUrl = '/') => {
  if (design.isDirty) {
    pendingTargetUrl.value = targetUrl
    design.showLeaveConfirm = true
  } else {
    navigateTo(targetUrl)
  }
}

const confirmLeaveWithoutSaving = () => {
  design.showLeaveConfirm = false
  design.isDirty = false
  navigateTo(pendingTargetUrl.value)
}

const confirmSaveAndLeave = () => {
  design.saveDraft()
  design.showLeaveConfirm = false
  navigateTo(pendingTargetUrl.value)
}

const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (design.isDirty) {
    e.preventDefault()
    e.returnValue = ''
  }
}

const handleResize = () => {
  design.updateScale()
}

onMounted(() => {
  design.loadDraft()
  design.updateScale()
  nextTick(() => design.updateScale())
  setTimeout(() => design.updateScale(), 100)
  setTimeout(() => design.updateScale(), 400)
  window.addEventListener('mousemove', design.onMouseMove)
  window.addEventListener('mouseup', design.onMouseUp)
  window.addEventListener('touchmove', design.onTouchMove, { passive: false })
  window.addEventListener('touchend', design.onTouchEnd)
  window.addEventListener('touchcancel', design.onTouchEnd)
  window.addEventListener('resize', handleResize)
  window.addEventListener('orientationchange', handleResize)
  window.addEventListener('keydown', design.onKeyDown)
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', design.onMouseMove)
  window.removeEventListener('mouseup', design.onMouseUp)
  window.removeEventListener('touchmove', design.onTouchMove)
  window.removeEventListener('touchend', design.onTouchEnd)
  window.removeEventListener('touchcancel', design.onTouchEnd)
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('orientationchange', handleResize)
  window.removeEventListener('keydown', design.onKeyDown)
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<template>
  <div class="dr-root">
    <!-- ═══ NAVBAR ═══════════════════════════════════════════════════ -->
    <nav class="dr-nav">
      <button class="dr-back" @click="handleTryLeave('/')" title="Back to Home">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        <span class="hidden xs:inline">Home</span>
      </button>
      <div class="dr-nav-center">
        <span class="dr-brand hidden sm:inline">CHIA FLORIST</span>
        <span class="dr-dot hidden sm:inline">◆</span>
        <span class="dr-page-title">Board Designer</span>
      </div>
      <div class="dr-nav-right">
        <!-- Zoom controls -->
        <div class="dr-zoom-controls">
          <button class="dr-zoom-btn" @click="design.zoomOut()" title="Zoom out (Ctrl+-)">
            <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <input type="range" min="20" max="150" step="5"
            :value="Math.round(design.boardScale * 100)"
            :style="{ '--v': Math.round(design.boardScale * 100) }"
            @input="(e) => design.setZoom(Number((e.target as HTMLInputElement).value))"
            @change="(e) => design.setZoom(Number((e.target as HTMLInputElement).value))"
            class="dr-zoom-slider hidden md:block" title="Zoom"/>
          <button class="dr-zoom-btn" @click="design.zoomIn()" title="Zoom in (Ctrl+=)">
            <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          </button>
          <button class="dr-scale-chip" @click="design.resetZoom()" title="Reset zoom (Ctrl+0)">{{ Math.round(design.boardScale * 100) }}%</button>
        </div>
      </div>
    </nav>

    <!-- ═══ BODY (Client Only for Interactive Board Designer) ════════ -->
    <ClientOnly>
      <div class="dr-body">
        <!-- Canvas area -->
        <div class="dr-canvas-wrapper">
          <CanvasBoard :design="design" />

          <!-- ✦ Canvas Summary Bar — elements + price + finalize (floating above toolbar) -->
          <div class="canvas-summary-bar">
            <span class="csb-count">
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
              {{ design.elements.length }} element{{ design.elements.length !== 1 ? 's' : '' }}
            </span>
            <span class="csb-divider">·</span>
            <span class="csb-price">{{ formatRupiah(design.totalPrice) }}</span>
            <button class="csb-finalize" @click="handleFinalize">
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M5 12l5 5L20 7"/></svg>
              Finalize &amp; Order
            </button>
          </div>
        </div>

        <!-- Right panel: width 0 when no tab active, 340px when open -->
        <ToolPanel
          :design="design"
          :class="{ 'panel-open': design.activeTab }"
        />
      </div>

      <!-- ⬡ ZZZ Bottom Toolbar — All 7 Tool & Action Buttons -->
      <nav class="zzz-bottom-bar" role="tablist" aria-label="Tool selection">
        
        <!-- ✦ More Action Button (First Item in Bar) -->
        <div class="relative flex-1 md:flex-initial flex justify-center">
          <button
            class="zzz-tab-btn zzz-more-btn w-full"
            :class="{ 'tab-active': design.showMoreMenu }"
            @click="design.showMoreMenu = !design.showMoreMenu"
            title="More Actions (Save, Reset, Randomize)"
            role="button"
            aria-label="More Options"
          >
            <span v-if="design.showMoreMenu" class="zzz-indicator"></span>
            <div class="zzz-tab-icon">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                <circle cx="5" cy="5" r="2.2"/><circle cx="12" cy="5" r="2.2"/><circle cx="19" cy="5" r="2.2"/>
                <circle cx="5" cy="12" r="2.2"/><circle cx="12" cy="12" r="2.2"/><circle cx="19" cy="12" r="2.2"/>
                <circle cx="5" cy="19" r="2.2"/><circle cx="12" cy="19" r="2.2"/><circle cx="19" cy="19" r="2.2"/>
              </svg>
            </div>
            <span class="zzz-tab-label">More</span>
          </button>

          <!-- MORE QUICK ACTIONS MENU POPOVER -->
          <div v-if="design.showMoreMenu" class="more-menu-popover" @click.stop>
            <button class="more-menu-item" @click="design.saveDraft(); design.showMoreMenu = false">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/>
              </svg>
              <span>Save Progress</span>
            </button>
            <button class="more-menu-item" @click="design.randomizeDesign(); design.showMoreMenu = false">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8"/><line x1="4" y1="20" x2="21" y2="3"/><polyline points="21 16 21 21 16 21"/><line x1="15" y1="15" x2="21" y2="21"/><line x1="4" y1="4" x2="9" y2="9"/>
              </svg>
              <span>Randomize Design</span>
            </button>
            <button class="more-menu-item danger" @click="design.resetDesign(); design.showMoreMenu = false">
              <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/>
              </svg>
              <span>Reset to Default</span>
            </button>
          </div>
        </div>

        <!-- ✦ 6 Tool Tabs: Text, Image, Brush, Border, Corner, Floral -->
        <button
          v-for="tab in TOOL_TABS"
          :key="tab.id"
          class="zzz-tab-btn flex-1 md:flex-initial"
          :class="{ 'tab-active': design.activeTab === tab.id }"
          @click="design.showMoreMenu = false; (design.activeTab === tab.id ? (design.activeTab = null) : (design.activeTab = tab.id))"
          role="tab"
          :aria-selected="design.activeTab === tab.id"
        >
          <span v-if="design.activeTab === tab.id" class="zzz-indicator"></span>
          <div class="zzz-tab-icon">
            <svg v-if="tab.id === 'text'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 7V4h16v3"/><path d="M9 20h6"/><path d="M12 4v16"/></svg>
            <svg v-else-if="tab.id === 'image'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
            <svg v-else-if="tab.id === 'brush'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M18.37 2.63 14 7l-1.59-1.59a2 2 0 0 0-2.82 0L8 7l9 9 1.59-1.59a2 2 0 0 0 0-2.82L17 10l4.37-4.37a2.12 2.12 0 1 0-3-3Z"/><path d="M9 8c-2 2.5-2 5-2 5"/></svg>
            <svg v-else-if="tab.id === 'border'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="2"/></svg>
            <svg v-else-if="tab.id === 'corner'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9V5a2 2 0 012-2h4"/><path d="M3 15v4a2 2 0 002 2h4"/></svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22C12 22 20 18 20 12C20 6 12 2 12 2C12 2 4 6 4 12C4 18 12 22 12 22Z"/><circle cx="12" cy="12" r="3"/></svg>
          </div>
          <span class="zzz-tab-label">{{ tab.label }}</span>
        </button>
      </nav>

      <!-- ✦ TOAST POPUP NOTIFICATION -->
      <Transition name="toast-fade">
        <div v-if="design.saveToastNotice" class="dr-toast-popup">
          {{ design.saveToastMessage }}
        </div>
      </Transition>

      <!-- ✦ UNSAVED LEAVE WARNING CONFIRMATION MODAL -->
      <div v-if="design.showLeaveConfirm" class="dr-modal-backdrop" @click.self="design.showLeaveConfirm = false">
        <div class="dr-modal choice-modal">
          <div class="dr-modal-head">
            <h2>Unsaved Progress Warning</h2>
            <button class="dr-modal-close" @click="design.showLeaveConfirm = false">×</button>
          </div>
          <div class="choice-body">
            <p class="choice-subtitle">
              Are you sure you want to leave? Your unsaved board progress will be lost.
            </p>
            <div class="choice-actions">
              <button class="choice-btn primary" @click="confirmSaveAndLeave">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                Save Draft &amp; Leave
              </button>
              <button class="choice-btn secondary" @click="confirmLeaveWithoutSaving">
                Leave Without Saving
              </button>
              <button class="choice-btn secondary" @click="design.showLeaveConfirm = false">
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ═══ MODALS & OVERLAYS ═════════════════════════════════════════ -->
      <FinalizeChoiceOverlay
        :show="design.showFinalizeChoice"
        @close="design.showFinalizeChoice = false"
        @review="handleOpenReview"
      />
      <ReviewModal
        :design="design"
        :is-adding="isAdding"
        @close="design.showReview = false"
        @add-to-cart="handleAddToCart"
      />
      <ThankYouOverlay
        :show="design.showThankYou"
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
