<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useCart } from '~/composables/useCart'

definePageMeta({ layout: false })

useHead({
  title: 'Chia Florist — Board Designer v2.0',
  meta: [{ name: 'description', content: 'Design your custom flower board in real-time with our interactive board simulator.' }]
})

const { addToCart, formatRupiah } = useCart()

// ─── PHASE MANAGEMENT ─────────────────────────────────────────
type Phase = 'splash' | 'tutorial' | 'builder' | 'review'
const phase = ref<Phase>('splash')
const phaseDir = ref<'forward' | 'back'>('forward')
const currentStep = ref(1)
const TOTAL_STEPS = 4
const stepFlash = ref(false)
const isAddingToCart = ref(false)

const showToast = ref(false)
const toastInfo = ref({
  name: '',
  image: '',
  quantity: 1,
  size: ''
})
let toastTimeout: any = null

// ─── FLOATING PARTICLES (SPLASH) ─────────────────────────────
interface Particle { id: number; x: number; y: number; size: number; opacity: number; dur: number; color: string; delay: number }
const particles = ref<Particle[]>([])
const PCOLORS = ['#10b981','#6ee7b7','#f59e0b','#fcd34d','#d1fae5','#a7f3d0','#ffffff']
const initParticles = () => {
  particles.value = Array.from({ length: 50 }, (_, i) => ({
    id: i,
    x: Math.random() * 100,
    y: Math.random() * 100,
    size: Math.random() * 7 + 2,
    opacity: Math.random() * 0.4 + 0.07,
    dur: Math.random() * 8 + 5,
    color: PCOLORS[Math.floor(Math.random() * PCOLORS.length)]!,
    delay: Math.random() * 6
  }))
}

// ─── INTERFACES ───────────────────────────────────────────────
interface SizeOption { id: string; label: string; price: number; desc: string; ideal: string; class: string; bars: number; recommended?: boolean }
interface ThemeOption { id: string; label: string; color: string; price: number; mood: string; accent: string }
interface FlowerOption { id: string; label: string; price: number; desc: string; emoji: string; recommended?: boolean }
interface FontOption { id: string; label: string; family: string; price: number; preview: string; recommended?: boolean }
interface CustomSelection {
  size: SizeOption; theme: ThemeOption; flower: FlowerOption; font: FontOption
  text: { header: string; target: string; sender: string }
}

// ─── CATALOG DATA ─────────────────────────────────────────────
const options = {
  sizes: [
    { id: 'small',  label: '1.5m × 2.0m', price: 150000, desc: 'Compact',  ideal: 'Indoor events & intimate spaces',     class: 'max-w-xs md:max-w-sm', bars: 1 },
    { id: 'medium', label: '1.8m × 2.5m', price: 200000, desc: 'Standard', ideal: 'Weddings & celebrations',             class: 'max-w-sm md:max-w-md', bars: 2, recommended: true },
    { id: 'large',  label: '2.0m × 3.0m', price: 280000, desc: 'Grand',    ideal: 'Corporate events & grand openings',   class: 'max-w-md md:max-w-lg', bars: 3 }
  ] as SizeOption[],
  themes: [
    { id: 'emerald', label: 'Florist Emerald', color: '#114028', price: 0,     mood: 'Hope & New Beginnings',    accent: '#10b981' },
    { id: 'navy',    label: 'Royal Navy',       color: '#0c1b30', price: 0,     mood: 'Professional Elegance',    accent: '#60a5fa' },
    { id: 'maroon',  label: 'Wine Maroon',       color: '#4c0519', price: 25000, mood: 'Romantic & Luxurious',     accent: '#f43f5e' },
    { id: 'black',   label: 'Luxury Black',      color: '#111111', price: 15000, mood: 'Dignified & Formal',       accent: '#94a3b8' }
  ] as ThemeOption[],
  flowers: [
    { id: 'basic',    label: 'Minimalist Top',  price: 40000,  desc: 'Delicate blooms along the top edge only',        emoji: '🌸' },
    { id: 'standard', label: 'Corner Duo',      price: 70000,  desc: 'Lush clusters at opposing corners',              emoji: '🌺', recommended: true },
    { id: 'luxury',   label: 'Full Frame',       price: 125000, desc: 'Abundant florals enveloping the entire board',   emoji: '👑' }
  ] as FlowerOption[],
  fonts: [
    { id: 'serif',  label: 'Classic Luxury',   family: 'font-serif',        price: 15000, preview: 'Aa — Timeless & stately' },
    { id: 'modern', label: 'Modern Sans',       family: 'font-sans',         price: 0,     preview: 'Aa — Clean & readable',   recommended: true },
    { id: 'script', label: 'Elegant Script',    family: 'italic font-serif', price: 20000, preview: 'Aa — Romantic & personal' }
  ] as FontOption[]
}

const occasionPresets = [
  { label: '💍 Wedding',     header: 'HAPPY WEDDING',           target: 'The Happy Couple' },
  { label: '🏢 Opening',     header: 'GRAND OPENING',           target: 'Your Business Name' },
  { label: '🎂 Birthday',    header: 'HAPPY BIRTHDAY',          target: 'Your Name Here' },
  { label: '🕊️ Condolence',  header: 'DEEPEST CONDOLENCES',     target: 'The Family' },
  { label: '💎 Anniversary', header: 'HAPPY ANNIVERSARY',       target: 'The Couple' },
  { label: '🎓 Graduation',  header: 'CONGRATULATIONS',         target: 'Your Name Here' }
]

// ─── REACTIVE STATE ───────────────────────────────────────────
const selection = ref<CustomSelection>({
  size:   options.sizes[1]!,
  theme:  options.themes[0]!,
  flower: options.flowers[1]!,
  font:   options.fonts[1]!,
  text:   { header: 'HAPPY WEDDING', target: 'The Happy Couple', sender: 'Chia Florist' }
})

const totalPrice = computed(() =>
  selection.value.size.price  +
  selection.value.theme.price +
  selection.value.flower.price +
  selection.value.font.price
)

// Animated price counter
const displayPrice = ref(totalPrice.value)
let priceTimer: ReturnType<typeof setTimeout> | null = null
watch(totalPrice, (newVal, oldVal) => {
  if (priceTimer) clearTimeout(priceTimer)
  const steps = 28, diff = newVal - oldVal
  let step = 0
  const tick = () => {
    step++
    displayPrice.value = Math.round(oldVal + diff * (step / steps))
    if (step < steps) priceTimer = setTimeout(tick, 10)
  }
  tick()
})

const stepLabels = ['Board Scale', 'Foam Theme', 'Decoration', 'Engrave']

const stepTip = computed(() => {
  if (currentStep.value === 1) {
    return selection.value.size.id === 'small'  ? '💡 Best for intimate indoor events and compact venues.'
         : selection.value.size.id === 'medium' ? '💡 Our most popular pick — the perfect balance of presence and value.'
         :                                        '💡 Maximum visual impact for corporate events and grand openings.'
  }
  if (currentStep.value === 2) {
    return selection.value.theme.id === 'emerald' ? '💡 Symbolizes hope, life, and universal celebration.'
         : selection.value.theme.id === 'navy'    ? '💡 Projects authority and professional grace.'
         : selection.value.theme.id === 'maroon'  ? '💡 Deep romance — ideal for anniversaries and weddings.'
         :                                          '💡 Dignified restraint — best for condolences and formal events.'
  }
  if (currentStep.value === 3) return '💡 Select your flower arrangement, then pick a typography style to match.'
  return '💡 Your message will be professionally engraved on the completed board.'
})

// ─── NAVIGATION ───────────────────────────────────────────────
const goToPhase = (p: Phase, dir: 'forward' | 'back' = 'forward') => {
  phaseDir.value = dir
  phase.value = p
}

const next = () => {
  if (currentStep.value < TOTAL_STEPS) {
    stepFlash.value = true
    setTimeout(() => { stepFlash.value = false }, 500)
    currentStep.value++
  } else {
    nextTick(() => { triggerConfetti(); goToPhase('review', 'forward') })
  }
}

const prev = () => {
  if (currentStep.value > 1) currentStep.value--
  else goToPhase('tutorial', 'back')
}

const goToStep = (s: number) => {
  if (s >= 1 && s <= TOTAL_STEPS) currentStep.value = s
}

const applyPreset = (p: typeof occasionPresets[0]) => {
  selection.value.text.header = p.header
  selection.value.text.target = p.target
}

// ─── CONFETTI ─────────────────────────────────────────────────
interface Confetti { id: number; x: number; color: string; delay: number; size: number; rot: number; speed: number }
const confettiPieces = ref<Confetti[]>([])
const triggerConfetti = () => {
  const cols = ['#10b981','#f59e0b','#ffffff','#6ee7b7','#fcd34d','#86efac','#34d399']
  confettiPieces.value = Array.from({ length: 40 }, (_, i) => ({
    id: i,
    x: Math.random() * 100,
    color: cols[Math.floor(Math.random() * cols.length)]!,
    delay: Math.random() * 2.2,
    size: Math.random() * 10 + 5,
    rot: Math.random() * 360,
    speed: Math.random() * 2 + 2
  }))
}

// ─── ADD TO CART ──────────────────────────────────────────────
const closeToastAndGoHome = () => {
  showToast.value = false
  if (toastTimeout) clearTimeout(toastTimeout)
  navigateTo('/')
}

const handleCustomAddToCart = async () => {
  isAddingToCart.value = true
  await new Promise(r => setTimeout(r, 700))
  addToCart({
    id: 'custom-' + Date.now(),
    name: `Custom Board — ${selection.value.text.header || 'My Design'}`,
    price: totalPrice.value,
    image: '/images/custom-preview.png',
    size: selection.value.size.label,
    color: selection.value.theme.color,
    shopId: '99ef0062-1040-4574-a4be-0123abce5670',
    isCustom: true
  }, 1)
  isAddingToCart.value = false
  
  toastInfo.value = {
    name: `Custom Board — ${selection.value.text.header || 'My Design'}`,
    image: '/images/custom-preview.png',
    quantity: 1,
    size: selection.value.size.label
  }
  showToast.value = true
  if (toastTimeout) clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    closeToastAndGoHome()
  }, 2000)
}

// ─── KEYBOARD NAV ─────────────────────────────────────────────
const handleKey = (e: KeyboardEvent) => {
  if (phase.value !== 'builder') return
  if (e.key === 'ArrowRight') next()
  if (e.key === 'ArrowLeft')  prev()
}

onMounted(() => { initParticles(); window.addEventListener('keydown', handleKey) })
onUnmounted(() => { 
  window.removeEventListener('keydown', handleKey)
  if (priceTimer) clearTimeout(priceTimer)
  if (toastTimeout) clearTimeout(toastTimeout)
})
</script>

<template>
  <!-- ━━━━━━━━━━━━━━━━━━ ROOT ━━━━━━━━━━━━━━━━━━ -->
  <div class="sim-root">

    <!-- ── PERSISTENT PARTICLE FIELD ── -->
    <div class="particles-layer" aria-hidden="true">
      <span
        v-for="p in particles" :key="p.id"
        class="particle"
        :style="{
          left: p.x + '%', top: p.y + '%',
          width: p.size + 'px', height: p.size + 'px',
          background: p.color, opacity: p.opacity,
          animationDuration: p.dur + 's',
          animationDelay: p.delay + 's'
        }"
      />
    </div>

    <!-- ── DIAGONAL DECORATION STRIPES ── -->
    <div class="deco-stripes" aria-hidden="true">
      <div class="stripe stripe-1" />
      <div class="stripe stripe-2" />
      <div class="stripe stripe-3" />
    </div>

    <!-- ━━━━━━━━━ PHASE TRANSITIONS ━━━━━━━━━ -->
    <Transition :name="phaseDir === 'forward' ? 'slide-fwd' : 'slide-bck'" mode="out-in">

      <!-- ══════════════════════════════════════════
           PHASE 1 — SPLASH SCREEN
      ══════════════════════════════════════════ -->
      <div v-if="phase === 'splash'" key="splash" class="phase-screen splash-screen">

        <div class="splash-inner">

          <!-- Version badge -->
          <div class="badge-row">
            <span class="badge-patch">◆ NEW PATCH</span>
            <span class="badge-ver">v2.0</span>
          </div>

          <!-- Big title -->
          <div class="splash-title-block">
            <p class="splash-label">CHIA FLORIST</p>
            <h1 class="splash-title">
              <span class="title-outline">BOARD</span>
              <span class="title-solid"> DESIGNER</span>
            </h1>
          </div>

          <!-- Tagline -->
          <p class="splash-tagline">
            Craft your perfect flower board in real-time.<br />
            Choose size, theme, florals &amp; typography — then order with one click.
          </p>

          <!-- Feature chips -->
          <div class="feature-chips">
            <div class="chip"><span>🎨</span><span>Live Preview</span></div>
            <div class="chip"><span>⚡</span><span>4-Step Builder</span></div>
            <div class="chip"><span>🛒</span><span>Direct to Cart</span></div>
          </div>

          <!-- CTA -->
          <button class="btn-start" @click="goToPhase('tutorial', 'forward')">
            <span>START DESIGN</span>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
          </button>

          <!-- Sub note -->
          <p class="splash-note">Press <kbd>→</kbd> / <kbd>←</kbd> to navigate steps during the builder.</p>
        </div>

        <!-- Corner accent graphic -->
        <div class="splash-corner-art" aria-hidden="true">
          <div class="corner-ring ring-1" />
          <div class="corner-ring ring-2" />
          <div class="corner-ring ring-3" />
          <div class="corner-dot" />
        </div>

        <div class="splash-version-watermark">v2.0.0 — Chia Dev</div>
      </div>

      <!-- ══════════════════════════════════════════
           PHASE 2 — TUTORIAL SCREEN
      ══════════════════════════════════════════ -->
      <div v-else-if="phase === 'tutorial'" key="tutorial" class="phase-screen tutorial-screen">

        <div class="tutorial-inner">

          <!-- Header -->
          <div class="tut-header">
            <button class="btn-ghost-back" @click="goToPhase('splash', 'back')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
              <span>Back</span>
            </button>
            <div class="tut-title-block">
              <span class="tut-eyebrow">HOW TO PLAY</span>
              <h2 class="tut-title">Design Your Board</h2>
            </div>
          </div>

          <!-- Step cards -->
          <div class="tut-cards">
            <div class="tut-card" style="animation-delay: 0.05s">
              <div class="tut-card-icon tut-icon-1">📏</div>
              <div class="tut-card-body">
                <div class="tut-step-num">STEP 01</div>
                <h4 class="tut-card-title">Board Scale</h4>
                <p class="tut-card-desc">Choose the physical dimensions of your flower board based on your event's venue and requirements.</p>
              </div>
              <div class="tut-card-accent" />
            </div>

            <div class="tut-card" style="animation-delay: 0.12s">
              <div class="tut-card-icon tut-icon-2">🎨</div>
              <div class="tut-card-body">
                <div class="tut-step-num">STEP 02</div>
                <h4 class="tut-card-title">Foam Theme</h4>
                <p class="tut-card-desc">Select the background foam color of the board — this sets the overall mood and atmosphere of your message.</p>
              </div>
              <div class="tut-card-accent" />
            </div>

            <div class="tut-card" style="animation-delay: 0.19s">
              <div class="tut-card-icon tut-icon-3">💐</div>
              <div class="tut-card-body">
                <div class="tut-step-num">STEP 03</div>
                <h4 class="tut-card-title">Flower &amp; Typography</h4>
                <p class="tut-card-desc">Pick your floral arrangement style and the engraving font. These two choices define the visual personality of the board.</p>
              </div>
              <div class="tut-card-accent" />
            </div>

            <div class="tut-card" style="animation-delay: 0.26s">
              <div class="tut-card-icon tut-icon-4">✍️</div>
              <div class="tut-card-body">
                <div class="tut-step-num">STEP 04</div>
                <h4 class="tut-card-title">Engrave Your Message</h4>
                <p class="tut-card-desc">Type your greeting, the recipient's name, and the sender. Use an occasion preset to get started instantly.</p>
              </div>
              <div class="tut-card-accent" />
            </div>
          </div>

          <!-- CTA -->
          <div class="tut-cta-row">
            <p class="tut-hint">🏆 Every selection updates the live preview in real-time. Design until it's perfect, then add to cart!</p>
            <button class="btn-begin" @click="goToPhase('builder', 'forward')">
              <span>BEGIN DESIGN</span>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
            </button>
          </div>

        </div>
      </div>

      <!-- ══════════════════════════════════════════
           PHASE 3 — BUILDER SCREEN
      ══════════════════════════════════════════ -->
      <div v-else-if="phase === 'builder'" key="builder" class="phase-screen builder-screen">

        <!-- ── TOP HUD BAR ── -->
        <div class="hud-bar">
          <NuxtLink to="/" class="hud-back-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
            <span>Home</span>
          </NuxtLink>

          <div class="hud-center">
            <span class="hud-label">CHIA FLORIST</span>
            <span class="hud-divider">◆</span>
            <span class="hud-title">BOARD DESIGNER</span>
          </div>

          <div class="hud-right">
            <span class="hud-kbd-hint"><kbd>←</kbd><kbd>→</kbd> Navigate</span>
            <span class="hud-ver">v2.0</span>
          </div>
        </div>

        <!-- ── MAIN SPLIT LAYOUT ── -->
        <div class="builder-split">

          <!-- ── LEFT: CANVAS ── -->
          <div class="canvas-panel">

            <!-- Live preview badge -->
            <div class="live-badge">
              <span class="live-dot" />
              <span>LIVE PREVIEW</span>
            </div>

            <!-- Board canvas -->
            <div
              class="board-wrap"
              :class="[selection.size.class, { 'step-flash-anim': stepFlash }]"
            >
              <div
                class="board-frame"
                :style="{ backgroundColor: selection.theme.color }"
                :class="selection.font.family"
              >
                <!-- Flower decorations -->
                <!-- TOP ROW (basic + luxury) -->
                <div v-if="selection.flower.id === 'basic' || selection.flower.id === 'luxury'" class="flower-top-row" aria-hidden="true">
                  <span v-for="n in 7" :key="n" class="flower-top-item">🌸</span>
                </div>
                <!-- CORNER TL (standard + luxury) -->
                <div v-if="selection.flower.id === 'standard' || selection.flower.id === 'luxury'" class="flower-corner flower-tl" aria-hidden="true">🌺</div>
                <!-- CORNER BR (standard + luxury) -->
                <div v-if="selection.flower.id === 'standard' || selection.flower.id === 'luxury'" class="flower-corner flower-br" aria-hidden="true">🌺</div>
                <!-- LEFT SIDE (luxury only) -->
                <div v-if="selection.flower.id === 'luxury'" class="flower-side flower-left" aria-hidden="true">
                  <span>🌸</span><span>🌸</span><span>🌸</span>
                </div>
                <!-- RIGHT SIDE (luxury only) -->
                <div v-if="selection.flower.id === 'luxury'" class="flower-side flower-right" aria-hidden="true">
                  <span>🌸</span><span>🌸</span><span>🌸</span>
                </div>
                <!-- BOTTOM ROW (luxury only) -->
                <div v-if="selection.flower.id === 'luxury'" class="flower-bottom-row" aria-hidden="true">
                  <span v-for="n in 7" :key="n" class="flower-top-item">🌸</span>
                </div>

                <!-- Board inner content -->
                <div class="board-content">
                  <h2 class="board-header">{{ selection.text.header || 'YOUR GREETING' }}</h2>
                  <div class="board-divider" :style="{ borderColor: selection.theme.accent || 'rgba(255,255,255,0.3)' }" />
                  <p class="board-target">{{ selection.text.target || 'Recipient Name' }}</p>
                  <div class="board-sender-row">
                    <span class="board-sender-line" />
                    <p class="board-sender">{{ selection.text.sender || 'Sender Name' }}</p>
                    <span class="board-sender-line" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Price HUD -->
            <div class="price-hud">
              <div class="price-hud-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.407 2.63 1M12 8V7m0 1v1m0 5v1m0-1c-1.11 0-2.08-.407-2.63-1M12 18v1m0-1V17" stroke-linecap="round" stroke-linejoin="round"/></svg>
              </div>
              <div>
                <p class="price-hud-label">ESTIMATED TOTAL</p>
                <p class="price-hud-value">{{ formatRupiah(displayPrice) }}</p>
              </div>
              <div class="price-breakdown-mini">
                <span>Scale: {{ formatRupiah(selection.size.price) }}</span>
                <span v-if="selection.theme.price > 0">Theme: +{{ formatRupiah(selection.theme.price) }}</span>
                <span>Florals: +{{ formatRupiah(selection.flower.price) }}</span>
                <span v-if="selection.font.price > 0">Font: +{{ formatRupiah(selection.font.price) }}</span>
              </div>
            </div>
          </div>

          <!-- ── RIGHT: OPTIONS PANEL ── -->
          <div class="options-panel">

            <!-- Step progress HUD -->
            <div class="step-hud">
              <div class="step-hud-header">
                <span class="step-hud-title">DESIGNER OPTIONS</span>
                <span class="step-counter">{{ currentStep }} / {{ TOTAL_STEPS }}</span>
              </div>
              <div class="step-track">
                <button
                  v-for="i in TOTAL_STEPS" :key="i"
                  class="step-node"
                  :class="{ 'step-active': currentStep === i, 'step-done': currentStep > i }"
                  @click="goToStep(i)"
                  :title="stepLabels[i-1]"
                >
                  <span class="step-node-dot">
                    <svg v-if="currentStep > i" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
                    <span v-else>{{ i }}</span>
                  </span>
                  <span class="step-node-label">{{ stepLabels[i-1] }}</span>
                </button>
                <div class="step-track-line">
                  <div class="step-track-fill" :style="{ width: ((currentStep - 1) / (TOTAL_STEPS - 1) * 100) + '%' }" />
                </div>
              </div>
            </div>

            <!-- Step content area -->
            <div class="options-scroll">
              <Transition name="step-slide" mode="out-in">

                <!-- STEP 1: BOARD SCALE -->
                <div v-if="currentStep === 1" key="step1" class="step-block">
                  <div class="step-heading">
                    <span class="step-eyebrow">STEP 01</span>
                    <h3 class="step-title">Board Scale</h3>
                    <p class="step-desc">Choose the physical dimensions of your flower board.</p>
                  </div>
                  <div class="option-list">
                    <button
                      v-for="(s, idx) in options.sizes" :key="s.id"
                      class="option-card"
                      :class="{ 'option-selected': selection.size.id === s.id }"
                      :style="{ animationDelay: (idx * 0.07) + 's' }"
                      @click="selection.size = s"
                    >
                      <div class="option-card-left">
                        <!-- Visual size bars -->
                        <div class="size-bars">
                          <span v-for="b in 3" :key="b" class="size-bar" :class="{ 'size-bar-filled': b <= s.bars }" />
                        </div>
                        <div>
                          <div class="option-name">
                            {{ s.label }}
                            <span v-if="s.recommended" class="badge-curator">★ Top Pick</span>
                          </div>
                          <div class="option-sub">{{ s.desc }} — {{ s.ideal }}</div>
                        </div>
                      </div>
                      <div class="option-price">{{ formatRupiah(s.price) }}</div>
                      <div class="option-check">
                        <svg v-if="selection.size.id === s.id" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
                      </div>
                      <div class="option-card-accent" />
                    </button>
                  </div>
                  <div class="step-tip">{{ stepTip }}</div>
                </div>

                <!-- STEP 2: FOAM THEME -->
                <div v-else-if="currentStep === 2" key="step2" class="step-block">
                  <div class="step-heading">
                    <span class="step-eyebrow">STEP 02</span>
                    <h3 class="step-title">Foam Theme</h3>
                    <p class="step-desc">Set the background foam color — it defines the board's mood.</p>
                  </div>
                  <div class="option-list">
                    <button
                      v-for="(t, idx) in options.themes" :key="t.id"
                      class="option-card"
                      :class="{ 'option-selected': selection.theme.id === t.id }"
                      :style="{ animationDelay: (idx * 0.07) + 's' }"
                      @click="selection.theme = t"
                    >
                      <div class="option-card-left">
                        <div class="color-swatch" :style="{ background: t.color }">
                          <div class="swatch-accent" :style="{ background: t.accent }" />
                        </div>
                        <div>
                          <div class="option-name">{{ t.label }}</div>
                          <div class="option-sub">{{ t.mood }}</div>
                        </div>
                      </div>
                      <div class="option-price">{{ t.price === 0 ? 'FREE' : '+' + formatRupiah(t.price) }}</div>
                      <div class="option-check">
                        <svg v-if="selection.theme.id === t.id" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
                      </div>
                      <div class="option-card-accent" />
                    </button>
                  </div>
                  <div class="step-tip">{{ stepTip }}</div>
                </div>

                <!-- STEP 3: DECORATION -->
                <div v-else-if="currentStep === 3" key="step3" class="step-block">
                  <div class="step-heading">
                    <span class="step-eyebrow">STEP 03</span>
                    <h3 class="step-title">Flower &amp; Typography</h3>
                    <p class="step-desc">Style your florals then pick an engraving font.</p>
                  </div>

                  <!-- Flower sub-section -->
                  <div class="sub-section">
                    <div class="sub-section-label">◆ FLORAL ARRANGEMENT</div>
                    <div class="option-list">
                      <button
                        v-for="(f, idx) in options.flowers" :key="f.id"
                        class="option-card"
                        :class="{ 'option-selected': selection.flower.id === f.id }"
                        :style="{ animationDelay: (idx * 0.07) + 's' }"
                        @click="selection.flower = f"
                      >
                        <div class="option-card-left">
                          <span class="flower-icon-big">{{ f.emoji }}</span>
                          <div>
                            <div class="option-name">
                              {{ f.label }}
                              <span v-if="f.recommended" class="badge-curator">★ Top Pick</span>
                            </div>
                            <div class="option-sub">{{ f.desc }}</div>
                          </div>
                        </div>
                        <div class="option-price">+{{ formatRupiah(f.price) }}</div>
                        <div class="option-check">
                          <svg v-if="selection.flower.id === f.id" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
                        </div>
                        <div class="option-card-accent" />
                      </button>
                    </div>
                  </div>

                  <!-- Font sub-section -->
                  <div class="sub-section">
                    <div class="sub-section-label">◆ ENGRAVING FONT</div>
                    <div class="option-list">
                      <button
                        v-for="(fn, idx) in options.fonts" :key="fn.id"
                        class="option-card"
                        :class="{ 'option-selected': selection.font.id === fn.id }"
                        :style="{ animationDelay: (idx * 0.07) + 's' }"
                        @click="selection.font = fn"
                      >
                        <div class="option-card-left">
                          <span class="font-preview" :class="fn.family">{{ fn.preview }}</span>
                        </div>
                        <div class="option-right-group">
                          <div class="option-name font-detail-name">
                            {{ fn.label }}
                            <span v-if="fn.recommended" class="badge-curator">★ Top Pick</span>
                          </div>
                          <div class="option-price">{{ fn.price === 0 ? 'FREE' : '+' + formatRupiah(fn.price) }}</div>
                        </div>
                        <div class="option-check">
                          <svg v-if="selection.font.id === fn.id" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
                        </div>
                        <div class="option-card-accent" />
                      </button>
                    </div>
                  </div>

                  <div class="step-tip">{{ stepTip }}</div>
                </div>

                <!-- STEP 4: ENGRAVE MESSAGE -->
                <div v-else-if="currentStep === 4" key="step4" class="step-block">
                  <div class="step-heading">
                    <span class="step-eyebrow">STEP 04</span>
                    <h3 class="step-title">Engrave Your Message</h3>
                    <p class="step-desc">Type your greeting, the recipient's name, and the sender — or tap a preset below.</p>
                  </div>

                  <!-- Occasion presets -->
                  <div class="preset-section">
                    <div class="preset-label">QUICK PRESETS</div>
                    <div class="preset-grid">
                      <button
                        v-for="pr in occasionPresets" :key="pr.label"
                        class="preset-chip"
                        :class="{ 'preset-active': selection.text.header === pr.header }"
                        @click="applyPreset(pr)"
                      >{{ pr.label }}</button>
                    </div>
                  </div>

                  <!-- Text inputs -->
                  <div class="input-group">
                    <div class="input-field">
                      <label class="input-label">MAIN GREETING</label>
                      <input
                        v-model="selection.text.header"
                        type="text"
                        placeholder="e.g. HAPPY WEDDING"
                        class="sim-input"
                        maxlength="40"
                      />
                      <span class="char-count">{{ selection.text.header.length }}/40</span>
                    </div>
                    <div class="input-field">
                      <label class="input-label">RECIPIENT NAME</label>
                      <input
                        v-model="selection.text.target"
                        type="text"
                        placeholder="e.g. John & Jane"
                        class="sim-input"
                        maxlength="40"
                      />
                      <span class="char-count">{{ selection.text.target.length }}/40</span>
                    </div>
                    <div class="input-field">
                      <label class="input-label">SENDER NAME</label>
                      <input
                        v-model="selection.text.sender"
                        type="text"
                        placeholder="e.g. Your Company"
                        class="sim-input"
                        maxlength="40"
                      />
                      <span class="char-count">{{ selection.text.sender.length }}/40</span>
                    </div>
                  </div>

                  <div class="step-tip">{{ stepTip }}</div>
                </div>

              </Transition>
            </div>

            <!-- Bottom nav bar -->
            <div class="nav-bar">
              <button class="btn-nav btn-back-step" @click="prev" :disabled="currentStep === 1 && false">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
                <span>{{ currentStep === 1 ? 'Tutorial' : 'Back' }}</span>
              </button>
              <button class="btn-nav btn-next-step" @click="next">
                <span>{{ currentStep < TOTAL_STEPS ? 'Next Step' : 'Review Design' }}</span>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
              </button>
            </div>

          </div><!-- end options-panel -->
        </div><!-- end builder-split -->

      </div><!-- end builder screen -->

      <!-- ══════════════════════════════════════════
           PHASE 4 — REVIEW SCREEN
      ══════════════════════════════════════════ -->
      <div v-else-if="phase === 'review'" key="review" class="phase-screen review-screen">

        <!-- Confetti -->
        <div class="confetti-layer" aria-hidden="true">
          <span
            v-for="c in confettiPieces" :key="c.id"
            class="confetti-piece"
            :style="{
              left: c.x + '%',
              width: c.size + 'px', height: (c.size * 1.4) + 'px',
              background: c.color,
              animationDelay: c.delay + 's',
              animationDuration: c.speed + 's',
              '--rot': c.rot + 'deg'
            }"
          />
        </div>

        <div class="review-inner">

          <!-- Header -->
          <div class="review-header">
            <button class="btn-ghost-back" @click="goToPhase('builder', 'back'); currentStep = 4">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
              <span>Edit Design</span>
            </button>
            <div class="review-title-block">
              <span class="review-eyebrow">🏆 DESIGN COMPLETE</span>
              <h2 class="review-title">Your Custom Board</h2>
            </div>
          </div>

          <!-- Review content -->
          <div class="review-content">

            <!-- Board mini preview -->
            <div class="review-board-wrap" :class="['max-w-xs']">
              <div
                class="board-frame review-board"
                :style="{ backgroundColor: selection.theme.color }"
                :class="selection.font.family"
              >
                <div v-if="selection.flower.id !== 'basic'" class="flower-corner flower-tl" style="font-size:1.8rem">🌺</div>
                <div v-if="selection.flower.id !== 'basic'" class="flower-corner flower-br" style="font-size:1.8rem">🌺</div>
                <div v-if="selection.flower.id === 'basic' || selection.flower.id === 'luxury'" class="flower-top-row" style="font-size:1rem">
                  <span v-for="n in 5" :key="n">🌸</span>
                </div>
                <div class="board-content">
                  <h2 class="board-header" style="font-size:clamp(0.7rem, 1.4vw, 1rem)">{{ selection.text.header }}</h2>
                  <div class="board-divider" />
                  <p class="board-target" style="font-size:clamp(0.9rem, 2vw, 1.4rem)">{{ selection.text.target }}</p>
                  <div class="board-sender-row">
                    <span class="board-sender-line" />
                    <p class="board-sender" style="font-size:0.65rem">{{ selection.text.sender }}</p>
                    <span class="board-sender-line" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Specs list -->
            <div class="review-specs">
              <div class="specs-title">ORDER SPECIFICATIONS</div>
              <div class="spec-row">
                <span class="spec-label">📏 Board Scale</span>
                <span class="spec-value">{{ selection.size.label }} ({{ selection.size.desc }})</span>
              </div>
              <div class="spec-row">
                <span class="spec-label">🎨 Foam Theme</span>
                <div class="spec-value-row">
                  <span class="spec-color-dot" :style="{ background: selection.theme.color }" />
                  <span class="spec-value">{{ selection.theme.label }}</span>
                </div>
              </div>
              <div class="spec-row">
                <span class="spec-label">💐 Florals</span>
                <span class="spec-value">{{ selection.flower.emoji }} {{ selection.flower.label }}</span>
              </div>
              <div class="spec-row">
                <span class="spec-label">🔤 Font Style</span>
                <span class="spec-value" :class="selection.font.family">{{ selection.font.label }}</span>
              </div>
              <div class="spec-row">
                <span class="spec-label">✍️ Greeting</span>
                <span class="spec-value">{{ selection.text.header }}</span>
              </div>
              <div class="spec-row">
                <span class="spec-label">👤 Recipient</span>
                <span class="spec-value">{{ selection.text.target }}</span>
              </div>
              <div class="spec-row">
                <span class="spec-label">📬 Sender</span>
                <span class="spec-value">{{ selection.text.sender }}</span>
              </div>

              <div class="spec-total-row">
                <span class="spec-total-label">TOTAL PRICE</span>
                <span class="spec-total-value">{{ formatRupiah(totalPrice) }}</span>
              </div>

              <!-- CTA Buttons -->
              <div class="review-cta">
                <button
                  class="btn-add-cart"
                  :class="{ 'btn-loading': isAddingToCart }"
                  :disabled="isAddingToCart"
                  @click="handleCustomAddToCart"
                >
                  <span v-if="!isAddingToCart" class="btn-add-cart-inner">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 0 0-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 0 0-16.536-1.84M7.5 14.25L5.106 5.272M16.5 20.25a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Zm3 0a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z"/></svg>
                    ADD TO CART — {{ formatRupiah(totalPrice) }}
                  </span>
                  <span v-else class="btn-add-cart-inner">
                    <svg class="spin-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3v3m9 9h-3m-6 6v-3m-9-9h3"/></svg>
                    Adding to cart...
                  </span>
                </button>
                <button class="btn-redesign" @click="goToPhase('builder', 'back'); currentStep = 1">
                  ← Start Over
                </button>
              </div>
            </div>

          </div>
        </div>
      </div>

    </Transition>

    <!-- Toast Notification -->
    <Transition name="toast-cust">
      <div v-if="showToast" class="cart-toast-custom" role="alert">
        <div class="toast-body-custom">
          <div class="toast-icon-check-custom">✓</div>
          <div class="toast-img-wrap-custom">
            <img :src="toastInfo.image" class="toast-img-custom" />
          </div>
          <div class="toast-details-custom">
            <h4 class="toast-title-custom">Added to Cart!</h4>
            <p class="toast-name-custom">{{ toastInfo.name }}</p>
            <p class="toast-meta-custom">Qty: {{ toastInfo.quantity }} | {{ toastInfo.size }}</p>
          </div>
        </div>
        <div class="toast-actions-custom">
          <NuxtLink to="/cart" class="btn-toast-view-custom">View Cart</NuxtLink>
          <button @click="closeToastAndGoHome" class="btn-toast-close-custom">×</button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style>
/* ─── GOOGLE FONTS ────────────────────────────────────────────── */
@import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=Inter:wght@400;600;700;900&display=swap');

/* ─── ROOT ────────────────────────────────────────────────────── */
.sim-root {
  position: fixed; inset: 0;
  background: #030d06;
  color: #fff;
  overflow: hidden;
  font-family: 'Inter', system-ui, sans-serif;
  user-select: none;
}

/* ─── PARTICLES ───────────────────────────────────────────────── */
.particles-layer {
  position: absolute; inset: 0; pointer-events: none; z-index: 0; overflow: hidden;
}
.particle {
  position: absolute; border-radius: 50%;
  animation: particleRise linear infinite;
}
@keyframes particleRise {
  0%   { transform: translateY(0) scale(1); opacity: var(--p-op, 0.2); }
  100% { transform: translateY(-160px) scale(0.3); opacity: 0; }
}

/* ─── DIAGONAL STRIPES ────────────────────────────────────────── */
.deco-stripes {
  position: absolute; inset: 0; pointer-events: none; z-index: 0; overflow: hidden;
}
.stripe {
  position: absolute; height: 2px; opacity: 0.035;
  transform-origin: left center;
}
.stripe-1 { width: 200%; top: 25%; transform: rotate(-12deg); background: linear-gradient(90deg, transparent, #10b981, transparent); }
.stripe-2 { width: 200%; top: 55%; transform: rotate(-12deg); background: linear-gradient(90deg, transparent, #f59e0b, transparent); opacity: 0.025; }
.stripe-3 { width: 200%; top: 78%; transform: rotate(-12deg); background: linear-gradient(90deg, transparent, #10b981, transparent); }

/* ─── PHASE SCREEN ────────────────────────────────────────────── */
.phase-screen {
  position: absolute; inset: 0; z-index: 10; overflow: hidden;
}

/* ─── PHASE TRANSITIONS ───────────────────────────────────────── */
.slide-fwd-enter-active { animation: phaseIn  0.42s cubic-bezier(0.22, 1, 0.36, 1) forwards; }
.slide-fwd-leave-active { animation: phaseOut 0.28s cubic-bezier(0.55, 0, 1, 0.45) forwards; }
.slide-bck-enter-active { animation: phaseInBack  0.42s cubic-bezier(0.22, 1, 0.36, 1) forwards; }
.slide-bck-leave-active { animation: phaseOutBack 0.28s cubic-bezier(0.55, 0, 1, 0.45) forwards; }

@keyframes phaseIn     { from { opacity: 0; transform: translateX(50px) skewX(-1.5deg); } to { opacity: 1; transform: none; } }
@keyframes phaseOut    { from { opacity: 1; transform: none; } to { opacity: 0; transform: translateX(-50px) skewX(1.5deg); } }
@keyframes phaseInBack  { from { opacity: 0; transform: translateX(-50px) skewX(1.5deg); } to { opacity: 1; transform: none; } }
@keyframes phaseOutBack { from { opacity: 1; transform: none; } to { opacity: 0; transform: translateX(50px) skewX(-1.5deg); } }

/* ─── STEP CONTENT TRANSITIONS ────────────────────────────────── */
.step-slide-enter-active { animation: stepIn  0.28s cubic-bezier(0.22, 1, 0.36, 1); }
.step-slide-leave-active { animation: stepOut 0.18s ease-in; }
@keyframes stepIn  { from { opacity: 0; transform: translateX(22px); } to { opacity: 1; transform: none; } }
@keyframes stepOut { from { opacity: 1; transform: none; } to { opacity: 0; transform: translateX(-18px); } }

/* ─── STEP FLASH ANIMATION ────────────────────────────────────── */
.step-flash-anim {
  animation: stepFlash 0.5s ease;
}
@keyframes stepFlash {
  0%   { box-shadow: 0 0 0 0 rgba(16,185,129,0); }
  30%  { box-shadow: 0 0 0 8px rgba(16,185,129,0.3), 0 0 50px rgba(16,185,129,0.15); }
  100% { box-shadow: 0 0 0 0 rgba(16,185,129,0); }
}

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   SPLASH SCREEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
.splash-screen {
  display: flex; align-items: center; justify-content: center;
  background: linear-gradient(145deg, #030d06 0%, #050f07 40%, #061209 100%);
}
.splash-inner {
  position: relative; z-index: 5;
  max-width: 680px; width: 100%; padding: 2rem;
  text-align: center;
  animation: splashEnter 0.7s cubic-bezier(0.22, 1, 0.36, 1);
}
@keyframes splashEnter {
  from { opacity: 0; transform: translateY(30px); }
  to   { opacity: 1; transform: none; }
}

.badge-row {
  display: inline-flex; gap: 0.5rem; align-items: center;
  margin-bottom: 1.5rem;
}
.badge-patch {
  background: rgba(16,185,129,0.15); border: 1px solid rgba(16,185,129,0.4);
  color: #10b981; font-size: 0.6rem; font-weight: 900; letter-spacing: 0.2em;
  padding: 0.3rem 0.8rem; border-radius: 2px;
  text-transform: uppercase;
  animation: pulseBadge 2s ease infinite;
}
@keyframes pulseBadge {
  0%,100% { box-shadow: 0 0 8px rgba(16,185,129,0.2); }
  50%      { box-shadow: 0 0 20px rgba(16,185,129,0.5); }
}
.badge-ver {
  background: rgba(245,158,11,0.15); border: 1px solid rgba(245,158,11,0.3);
  color: #f59e0b; font-size: 0.6rem; font-weight: 900; letter-spacing: 0.15em;
  padding: 0.3rem 0.7rem; border-radius: 2px;
}

.splash-title-block { margin-bottom: 1.25rem; }
.splash-label {
  font-family: 'Bebas Neue', 'Inter', sans-serif;
  font-size: clamp(0.9rem, 2vw, 1.1rem);
  letter-spacing: 0.5em; color: rgba(255,255,255,0.45); margin-bottom: 0.3rem;
}
.splash-title {
  font-family: 'Bebas Neue', 'Inter', sans-serif;
  font-size: clamp(3.5rem, 9vw, 6.5rem);
  line-height: 0.9; letter-spacing: -0.01em;
  display: flex; align-items: baseline; justify-content: center; gap: 0.2em;
  flex-wrap: wrap;
}
.title-outline {
  -webkit-text-stroke: 2px #10b981;
  color: transparent;
  filter: drop-shadow(0 0 20px rgba(16,185,129,0.4));
}
.title-solid {
  color: #fff;
  background: linear-gradient(90deg, #fff 0%, #d1fae5 50%, #fff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shimmer 3s linear infinite;
}
@keyframes shimmer {
  from { background-position: 200% center; }
  to   { background-position: -200% center; }
}

.splash-tagline {
  color: rgba(255,255,255,0.55); font-size: 0.875rem; line-height: 1.7;
  max-width: 500px; margin: 0 auto 1.75rem; font-weight: 500;
}

.feature-chips {
  display: flex; gap: 0.75rem; justify-content: center; flex-wrap: wrap;
  margin-bottom: 2rem;
}
.chip {
  display: flex; align-items: center; gap: 0.4rem;
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1);
  padding: 0.4rem 1rem; border-radius: 2px;
  font-size: 0.7rem; font-weight: 700; letter-spacing: 0.05em;
  color: rgba(255,255,255,0.6);
  transition: all 0.2s;
}
.chip:hover { background: rgba(16,185,129,0.1); border-color: rgba(16,185,129,0.3); color: #10b981; }

.btn-start {
  display: inline-flex; align-items: center; gap: 0.75rem;
  background: #10b981; color: #000;
  font-size: 0.9rem; font-weight: 900; letter-spacing: 0.15em;
  padding: 1rem 2.5rem; border-radius: 2px; border: none; cursor: pointer;
  box-shadow: 0 0 30px rgba(16,185,129,0.35);
  transition: all 0.2s cubic-bezier(0.22, 1, 0.36, 1);
  position: relative; overflow: hidden;
}
.btn-start svg { width: 1.2rem; height: 1.2rem; transition: transform 0.2s; }
.btn-start:hover { background: #34d399; box-shadow: 0 0 50px rgba(16,185,129,0.55); transform: translateY(-2px) scale(1.02); }
.btn-start:hover svg { transform: translateX(4px); }
.btn-start:active { transform: scale(0.97); }

.splash-note {
  margin-top: 1rem; color: rgba(255,255,255,0.2);
  font-size: 0.65rem; letter-spacing: 0.05em;
}
.splash-note kbd {
  background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.12);
  padding: 0.1rem 0.4rem; border-radius: 2px; font-size: 0.6rem;
}

.splash-corner-art {
  position: absolute; bottom: -80px; right: -80px;
  width: 380px; height: 380px; pointer-events: none;
}
.corner-ring {
  position: absolute; inset: 0; border-radius: 50%;
  border: 1px solid rgba(16,185,129,0.12);
  animation: rotateRing linear infinite;
}
.ring-1 { animation-duration: 20s; }
.ring-2 { inset: 30px; border-color: rgba(245,158,11,0.08); animation-duration: 30s; animation-direction: reverse; }
.ring-3 { inset: 70px; border-color: rgba(16,185,129,0.06); animation-duration: 45s; }
.corner-dot {
  position: absolute; width: 8px; height: 8px; background: #10b981;
  border-radius: 50%; top: 50%; left: 50%; transform: translate(-50%,-50%);
  box-shadow: 0 0 20px #10b981;
}
@keyframes rotateRing { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.splash-version-watermark {
  position: absolute; bottom: 1rem; right: 1.5rem;
  font-size: 0.6rem; font-weight: 700; color: rgba(255,255,255,0.12);
  letter-spacing: 0.1em;
}

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   TUTORIAL SCREEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
.tutorial-screen {
  display: flex; align-items: center; justify-content: center;
  background: linear-gradient(160deg, #030d06 0%, #050f07 100%);
  overflow-y: auto; padding: 2rem 1rem;
}
.tutorial-inner {
  max-width: 860px; width: 100%;
  display: flex; flex-direction: column; gap: 2rem;
}

.tut-header {
  display: flex; align-items: center; gap: 1.5rem;
}
.tut-eyebrow {
  display: block; font-size: 0.6rem; font-weight: 900;
  letter-spacing: 0.25em; color: #10b981; margin-bottom: 0.25rem;
}
.tut-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(2rem, 5vw, 3rem);
  letter-spacing: 0.05em; line-height: 1;
  color: #fff;
}

.tut-cards {
  display: grid; grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: 1rem;
}
.tut-card {
  position: relative; overflow: hidden;
  background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.07);
  border-radius: 2px; padding: 1.25rem 1.25rem 1.25rem 1rem;
  display: flex; gap: 1rem; align-items: flex-start;
  animation: cardFlyIn 0.5s cubic-bezier(0.22, 1, 0.36, 1) both;
  transition: background 0.2s, border-color 0.2s;
}
.tut-card:hover { background: rgba(16,185,129,0.05); border-color: rgba(16,185,129,0.2); }
@keyframes cardFlyIn {
  from { opacity: 0; transform: translateX(30px) skewX(-2deg); }
  to   { opacity: 1; transform: none; }
}
.tut-card-icon {
  font-size: 1.5rem; flex-shrink: 0;
  width: 3rem; height: 3rem; border-radius: 2px;
  display: flex; align-items: center; justify-content: center;
}
.tut-icon-1 { background: rgba(16,185,129,0.1); }
.tut-icon-2 { background: rgba(245,158,11,0.1); }
.tut-icon-3 { background: rgba(251,146,60,0.1); }
.tut-icon-4 { background: rgba(167,139,250,0.1); }

.tut-step-num {
  font-size: 0.55rem; font-weight: 900; letter-spacing: 0.2em;
  color: rgba(255,255,255,0.3); margin-bottom: 0.25rem;
}
.tut-card-title { font-size: 0.9rem; font-weight: 800; color: #f59e0b; margin-bottom: 0.35rem; }
.tut-card-desc { font-size: 0.75rem; color: rgba(255,255,255,0.5); line-height: 1.6; }

.tut-card-accent {
  position: absolute; left: 0; top: 0; bottom: 0;
  width: 3px; background: #10b981;
  transform: scaleY(0); transform-origin: top;
  transition: transform 0.3s;
}
.tut-card:hover .tut-card-accent { transform: scaleY(1); }

.tut-cta-row {
  display: flex; flex-direction: column; align-items: flex-start; gap: 1rem;
}
.tut-hint {
  font-size: 0.75rem; color: rgba(255,255,255,0.35);
  background: rgba(16,185,129,0.05); border: 1px solid rgba(16,185,129,0.12);
  padding: 0.75rem 1rem; border-radius: 2px; max-width: 600px; line-height: 1.6;
}
.btn-begin {
  display: inline-flex; align-items: center; gap: 0.6rem;
  background: #f59e0b; color: #000;
  font-size: 0.85rem; font-weight: 900; letter-spacing: 0.15em;
  padding: 0.875rem 2rem; border-radius: 2px; border: none; cursor: pointer;
  box-shadow: 0 0 25px rgba(245,158,11,0.3);
  transition: all 0.2s;
}
.btn-begin svg { width: 1rem; height: 1rem; transition: transform 0.2s; }
.btn-begin:hover { background: #fbbf24; transform: translateY(-2px); box-shadow: 0 0 40px rgba(245,158,11,0.5); }
.btn-begin:hover svg { transform: translateX(4px); }
.btn-begin:active { transform: scale(0.97); }

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   SHARED: GHOST BACK BUTTON
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
.btn-ghost-back {
  display: flex; align-items: center; gap: 0.4rem;
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.6); font-size: 0.72rem; font-weight: 700;
  letter-spacing: 0.05em; padding: 0.5rem 1rem; border-radius: 2px;
  cursor: pointer; transition: all 0.2s; flex-shrink: 0;
}
.btn-ghost-back svg { width: 0.9rem; height: 0.9rem; }
.btn-ghost-back:hover {
  background: rgba(255,255,255,0.08); border-color: rgba(255,255,255,0.2);
  color: #fff; transform: translateX(-2px);
}

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   BUILDER SCREEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
.builder-screen {
  display: flex; flex-direction: column;
  background: #050e07;
}

/* HUD BAR */
.hud-bar {
  height: 3.75rem; flex-shrink: 0;
  background: rgba(3,9,5,0.95); backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255,255,255,0.06);
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 1.5rem; z-index: 20;
}
.hud-back-link {
  display: flex; align-items: center; gap: 0.4rem;
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.08);
  color: rgba(255,255,255,0.5); font-size: 0.7rem; font-weight: 700;
  letter-spacing: 0.08em; padding: 0.4rem 0.9rem; border-radius: 2px;
  text-decoration: none; transition: all 0.2s;
}
.hud-back-link svg { width: 0.85rem; height: 0.85rem; }
.hud-back-link:hover { background: rgba(255,255,255,0.08); color: #fff; }

.hud-center {
  display: flex; align-items: center; gap: 0.6rem;
  position: absolute; left: 50%; transform: translateX(-50%);
}
.hud-label {
  font-size: 0.6rem; font-weight: 900; letter-spacing: 0.25em;
  color: rgba(255,255,255,0.3);
}
.hud-divider { color: #10b981; font-size: 0.5rem; }
.hud-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 1rem; letter-spacing: 0.15em; color: #10b981;
}

.hud-right {
  display: flex; align-items: center; gap: 0.75rem;
}
.hud-kbd-hint {
  display: flex; gap: 0.25rem; align-items: center;
  font-size: 0.55rem; color: rgba(255,255,255,0.2);
}
.hud-kbd-hint kbd {
  background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.1);
  padding: 0.1rem 0.35rem; border-radius: 2px; font-size: 0.5rem;
}
.hud-ver {
  font-size: 0.6rem; font-weight: 900; color: rgba(255,255,255,0.18);
  background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.06);
  padding: 0.25rem 0.6rem; border-radius: 2px; letter-spacing: 0.1em;
}

/* SPLIT LAYOUT */
.builder-split {
  flex: 1; display: flex; overflow: hidden;
}

/* CANVAS PANEL */
.canvas-panel {
  flex: 1; position: relative; display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  background: #040d06; border-right: 1px solid rgba(255,255,255,0.04);
  padding: 2rem 1.5rem; gap: 1.5rem; min-height: 0;
}

.live-badge {
  position: absolute; top: 1.25rem; left: 1.25rem;
  display: flex; align-items: center; gap: 0.4rem;
  background: rgba(0,0,0,0.5); backdrop-filter: blur(8px);
  border: 1px solid rgba(255,255,255,0.08);
  padding: 0.35rem 0.8rem; border-radius: 2px; z-index: 5;
}
.live-badge span:last-child {
  font-size: 0.6rem; font-weight: 800; letter-spacing: 0.15em; color: #10b981;
}
.live-dot {
  width: 6px; height: 6px; border-radius: 50%; background: #10b981;
  animation: pulseDot 1.5s ease infinite;
}
@keyframes pulseDot {
  0%,100% { box-shadow: 0 0 0 0 rgba(16,185,129,0.4); }
  50%      { box-shadow: 0 0 0 5px rgba(16,185,129,0); }
}

/* BOARD WRAP */
.board-wrap {
  width: 100%; position: relative; max-height: 60vh;
  display: flex; align-items: center; justify-content: center;
  transition: max-width 0.5s cubic-bezier(0.22, 1, 0.36, 1);
}

/* BOARD FRAME - premium wood border via gradient trick */
.board-frame {
  position: relative; width: 100%; aspect-ratio: 4/3;
  max-height: 58vh;
  border-radius: 3px; overflow: visible;
  /* Wood-like border using gradient + padding */
  padding: 14px;
  background-clip: padding-box;
  box-shadow:
    0 0 0 14px transparent,
    0 24px 60px rgba(0,0,0,0.7),
    0 0 80px rgba(0,0,0,0.3);
  outline: 14px solid;
  outline-color: transparent;
  /* The trick: inner bg + gradient border */
  background-image: none; /* color set via :style */
  border: 14px solid transparent;
  background-origin: border-box;
  transition: background-color 0.6s ease, border-color 0.4s;
  /* wood gradient border */
  box-shadow:
    0 24px 80px rgba(0,0,0,0.8),
    inset 0 0 60px rgba(0,0,0,0.2);
}
.board-frame::before {
  content: '';
  position: absolute;
  inset: -14px;
  border-radius: 3px;
  background: linear-gradient(135deg,
    #a0652a 0%, #c49a6c 20%, #8b5a2b 40%,
    #d4a574 55%, #8b5a2b 70%, #c49a6c 85%, #6b3d1a 100%);
  z-index: -1;
  box-shadow: 0 0 0 1px rgba(0,0,0,0.5), inset 0 0 8px rgba(0,0,0,0.3);
}

/* Board content positioning */
.board-content {
  position: absolute; inset: 0;
  display: flex; flex-direction: column; align-items: center; justify-content: space-around;
  padding: 1rem; z-index: 5; text-align: center;
}
.board-header {
  font-size: clamp(0.75rem, 2vw, 1.3rem);
  font-weight: 900; letter-spacing: 0.15em;
  text-transform: uppercase; color: #fff;
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
  line-height: 1.2; word-break: break-word;
  transition: all 0.4s;
}
.board-divider {
  width: 60%; border-top: 1px solid rgba(255,255,255,0.3);
  margin: 0 auto; flex-shrink: 0;
  transition: border-color 0.4s;
}
.board-target {
  font-size: clamp(1rem, 3.5vw, 2.2rem);
  font-weight: 800; color: #fde68a;
  text-shadow: 0 2px 12px rgba(0,0,0,0.6);
  font-style: italic; line-height: 1.2;
  word-break: break-word;
  transition: all 0.4s;
}
.board-sender-row {
  display: flex; align-items: center; gap: 0.5rem;
  width: 100%; justify-content: center; flex-shrink: 0;
}
.board-sender-line {
  flex: 1; max-width: 60px; height: 1px; background: rgba(255,255,255,0.2);
}
.board-sender {
  font-size: clamp(0.55rem, 1.2vw, 0.8rem);
  font-weight: 600; color: rgba(255,255,255,0.75);
  letter-spacing: 0.08em;
  transition: all 0.4s;
}

/* Flower decorations */
.flower-top-row {
  position: absolute; top: -12px; left: 0; right: 0;
  display: flex; justify-content: space-around; z-index: 10; pointer-events: none;
  font-size: clamp(1rem, 2.5vw, 2rem);
}
.flower-top-item { display: inline-block; animation: floatFlower 3s ease-in-out infinite; }
.flower-top-item:nth-child(even) { animation-delay: 0.5s; animation-direction: alternate-reverse; }
@keyframes floatFlower {
  0%,100% { transform: translateY(0) rotate(-5deg); }
  50%      { transform: translateY(-4px) rotate(5deg); }
}
.flower-bottom-row {
  position: absolute; bottom: -12px; left: 0; right: 0;
  display: flex; justify-content: space-around; z-index: 10; pointer-events: none;
  font-size: clamp(0.75rem, 1.8vw, 1.4rem);
}
.flower-corner {
  position: absolute; z-index: 10; pointer-events: none;
  font-size: clamp(1.5rem, 4vw, 3rem);
  transition: all 0.4s;
  animation: flowerCornerAnim 4s ease-in-out infinite alternate;
}
@keyframes flowerCornerAnim {
  from { transform: rotate(-5deg) scale(1); }
  to   { transform: rotate(5deg) scale(1.08); }
}
.flower-tl { top: -12px; left: -10px; }
.flower-br { bottom: -12px; right: -10px; animation-direction: alternate-reverse; }
.flower-side {
  position: absolute; top: 0; bottom: 0;
  display: flex; flex-direction: column; justify-content: space-around;
  font-size: clamp(0.7rem, 1.5vw, 1.2rem); z-index: 10; pointer-events: none;
}
.flower-left { left: -14px; }
.flower-right { right: -14px; }

/* PRICE HUD */
.price-hud {
  display: flex; align-items: center; gap: 0.85rem;
  background: rgba(0,0,0,0.7); backdrop-filter: blur(12px);
  border: 1px solid rgba(16,185,129,0.2); border-radius: 2px;
  padding: 0.75rem 1.25rem; position: absolute; bottom: 1.25rem; right: 1.25rem;
  box-shadow: 0 0 20px rgba(16,185,129,0.1);
}
.price-hud-icon {
  width: 2.25rem; height: 2.25rem; border-radius: 2px;
  background: rgba(16,185,129,0.15); border: 1px solid rgba(16,185,129,0.25);
  display: flex; align-items: center; justify-content: center;
  color: #10b981; flex-shrink: 0;
}
.price-hud-icon svg { width: 1.1rem; height: 1.1rem; }
.price-hud-label {
  font-size: 0.55rem; font-weight: 800; letter-spacing: 0.15em;
  color: rgba(255,255,255,0.35); margin-bottom: 0.15rem;
}
.price-hud-value {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 1.6rem; color: #10b981; line-height: 1;
  letter-spacing: 0.05em;
  transition: color 0.2s;
}
.price-breakdown-mini {
  border-left: 1px solid rgba(255,255,255,0.08);
  padding-left: 0.85rem;
  display: flex; flex-direction: column; gap: 0.1rem;
}
.price-breakdown-mini span {
  font-size: 0.6rem; color: rgba(255,255,255,0.3); font-weight: 600;
}

/* OPTIONS PANEL */
.options-panel {
  width: clamp(320px, 30vw, 440px); flex-shrink: 0;
  background: #03100a; border-left: 1px solid rgba(255,255,255,0.05);
  display: flex; flex-direction: column; overflow: hidden;
}

/* STEP HUD */
.step-hud {
  flex-shrink: 0; padding: 1rem 1.25rem 0.875rem;
  background: #040f08; border-bottom: 1px solid rgba(255,255,255,0.05);
}
.step-hud-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 0.875rem;
}
.step-hud-title {
  font-size: 0.6rem; font-weight: 900; letter-spacing: 0.2em;
  color: rgba(255,255,255,0.3);
}
.step-counter {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 0.9rem; letter-spacing: 0.05em;
  color: #10b981;
}

.step-track {
  position: relative; display: flex; gap: 0;
  align-items: flex-start; justify-content: space-between;
}
.step-track-line {
  position: absolute; top: 0.75rem; left: 0.75rem; right: 0.75rem; height: 2px;
  background: rgba(255,255,255,0.06); z-index: 0;
}
.step-track-fill {
  height: 100%; background: #10b981;
  transition: width 0.4s cubic-bezier(0.22, 1, 0.36, 1);
  box-shadow: 0 0 8px rgba(16,185,129,0.5);
}
.step-node {
  display: flex; flex-direction: column; align-items: center; gap: 0.3rem;
  background: none; border: none; cursor: pointer; padding: 0;
  flex: 1; z-index: 1; min-width: 0;
  transition: opacity 0.2s;
}
.step-node:hover { opacity: 0.8; }
.step-node-dot {
  width: 1.5rem; height: 1.5rem; border-radius: 50%;
  border: 2px solid rgba(255,255,255,0.15);
  background: #03100a;
  display: flex; align-items: center; justify-content: center;
  font-size: 0.6rem; font-weight: 800; color: rgba(255,255,255,0.3);
  transition: all 0.3s;
}
.step-node-dot svg { width: 0.65rem; height: 0.65rem; }
.step-node.step-active .step-node-dot {
  border-color: #10b981; background: #10b981; color: #000;
  box-shadow: 0 0 14px rgba(16,185,129,0.5);
}
.step-node.step-done .step-node-dot {
  border-color: rgba(16,185,129,0.5); background: rgba(16,185,129,0.12); color: #10b981;
}
.step-node-label {
  font-size: 0.5rem; font-weight: 700; letter-spacing: 0.04em;
  color: rgba(255,255,255,0.2); text-align: center;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 100%;
  padding: 0 0.1rem;
}
.step-node.step-active .step-node-label { color: #10b981; }
.step-node.step-done .step-node-label { color: rgba(16,185,129,0.6); }

/* STEP SCROLL AREA */
.options-scroll {
  flex: 1; overflow-y: auto; padding: 1.25rem;
  scrollbar-width: thin; scrollbar-color: #10b981 transparent;
}
.options-scroll::-webkit-scrollbar { width: 4px; }
.options-scroll::-webkit-scrollbar-thumb { background: rgba(16,185,129,0.35); border-radius: 2px; }

.step-heading { margin-bottom: 1.1rem; }
.step-eyebrow {
  font-size: 0.55rem; font-weight: 900; letter-spacing: 0.25em;
  color: rgba(255,255,255,0.25); display: block; margin-bottom: 0.2rem;
}
.step-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 1.5rem; letter-spacing: 0.06em; color: #fff; line-height: 1; margin-bottom: 0.3rem;
}
.step-desc { font-size: 0.72rem; color: rgba(255,255,255,0.4); line-height: 1.5; }
.step-tip {
  font-size: 0.7rem; color: #6ee7b7;
  background: rgba(16,185,129,0.06); border-left: 2px solid #10b981;
  padding: 0.6rem 0.75rem; margin-top: 1rem;
  border-radius: 0 2px 2px 0; line-height: 1.5;
}

/* OPTION CARDS */
.option-list { display: flex; flex-direction: column; gap: 0.6rem; }
.option-card {
  position: relative; overflow: hidden;
  display: flex; align-items: center; gap: 0.75rem;
  background: rgba(255,255,255,0.03); border: 1px solid rgba(255,255,255,0.07);
  border-left: 3px solid transparent;
  border-radius: 2px; padding: 0.8rem 0.85rem 0.8rem 0.75rem;
  text-align: left; cursor: pointer; width: 100%;
  animation: cardFlyIn 0.35s cubic-bezier(0.22, 1, 0.36, 1) both;
  transition: background 0.2s, border-color 0.2s, transform 0.15s;
}
.option-card:hover {
  background: rgba(255,255,255,0.05); border-color: rgba(255,255,255,0.15);
  border-left-color: rgba(16,185,129,0.5); transform: translateX(2px);
}
.option-card.option-selected {
  background: rgba(16,185,129,0.08); border-color: rgba(16,185,129,0.3);
  border-left-color: #10b981;
}
.option-card.option-selected::after {
  content: '';
  position: absolute; inset: 0;
  background: repeating-linear-gradient(
    -45deg,
    transparent, transparent 4px,
    rgba(16,185,129,0.02) 4px, rgba(16,185,129,0.02) 8px
  );
  pointer-events: none;
}
.option-card-left {
  flex: 1; display: flex; align-items: center; gap: 0.65rem; min-width: 0;
}
.option-name {
  font-size: 0.78rem; font-weight: 800; color: rgba(255,255,255,0.88);
  display: flex; align-items: center; gap: 0.4rem; flex-wrap: wrap;
  margin-bottom: 0.15rem;
}
.option-sub { font-size: 0.63rem; color: rgba(255,255,255,0.35); line-height: 1.4; }
.option-price {
  font-size: 0.7rem; font-weight: 900; color: #10b981;
  white-space: nowrap; flex-shrink: 0; font-family: 'Bebas Neue', sans-serif;
  font-size: 0.85rem; letter-spacing: 0.03em;
}
.option-check {
  width: 1.25rem; height: 1.25rem; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  color: #10b981;
}
.option-check svg { width: 0.9rem; height: 0.9rem; }
.option-card-accent {
  position: absolute; top: 0; right: 0; bottom: 0; width: 2px;
  background: linear-gradient(to bottom, transparent, rgba(16,185,129,0.3), transparent);
  transform: scaleY(0); transform-origin: center;
  transition: transform 0.25s;
}
.option-card.option-selected .option-card-accent { transform: scaleY(1); }

/* SIZE BARS */
.size-bars { display: flex; gap: 2px; flex-shrink: 0; }
.size-bar {
  width: 4px; height: 20px; border-radius: 1px;
  background: rgba(255,255,255,0.1); transition: background 0.2s;
}
.size-bar-filled { background: #10b981; box-shadow: 0 0 6px rgba(16,185,129,0.5); }

/* CURATOR BADGE */
.badge-curator {
  font-size: 0.5rem; font-weight: 900; letter-spacing: 0.05em;
  background: rgba(245,158,11,0.15); border: 1px solid rgba(245,158,11,0.3);
  color: #f59e0b; padding: 0.1rem 0.4rem; border-radius: 2px;
}

/* COLOR SWATCH */
.color-swatch {
  width: 2.5rem; height: 2rem; border-radius: 2px; flex-shrink: 0;
  border: 1px solid rgba(255,255,255,0.1); position: relative; overflow: hidden;
}
.swatch-accent {
  position: absolute; bottom: 0; left: 0; right: 0; height: 4px;
}

/* FLOWER & FONT STEP SPECIFICS */
.flower-icon-big { font-size: 1.5rem; flex-shrink: 0; }
.font-preview {
  font-size: 0.85rem; color: rgba(255,255,255,0.7); flex: 1;
  display: block; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.option-right-group { display: flex; flex-direction: column; gap: 0.15rem; align-items: flex-end; flex-shrink: 0; }
.font-detail-name { font-size: 0.7rem; font-weight: 800; color: rgba(255,255,255,0.7); white-space: nowrap; }

.sub-section { margin-bottom: 1.1rem; }
.sub-section-label {
  font-size: 0.55rem; font-weight: 900; letter-spacing: 0.2em;
  color: #f59e0b; margin-bottom: 0.6rem; display: block;
}

/* ENGRAVE STEP - PRESETS */
.preset-section { margin-bottom: 1rem; }
.preset-label {
  font-size: 0.55rem; font-weight: 900; letter-spacing: 0.2em;
  color: rgba(255,255,255,0.3); margin-bottom: 0.5rem; display: block;
}
.preset-grid {
  display: flex; flex-wrap: wrap; gap: 0.4rem;
}
.preset-chip {
  font-size: 0.65rem; font-weight: 700; cursor: pointer;
  background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.09);
  color: rgba(255,255,255,0.5); padding: 0.3rem 0.65rem; border-radius: 2px;
  transition: all 0.2s;
}
.preset-chip:hover { background: rgba(245,158,11,0.08); border-color: rgba(245,158,11,0.25); color: #f59e0b; }
.preset-chip.preset-active { background: rgba(245,158,11,0.12); border-color: rgba(245,158,11,0.4); color: #f59e0b; }

/* INPUTS */
.input-group { display: flex; flex-direction: column; gap: 0.8rem; }
.input-field { position: relative; }
.input-label {
  display: block; font-size: 0.55rem; font-weight: 900; letter-spacing: 0.2em;
  color: rgba(255,255,255,0.3); margin-bottom: 0.3rem;
}
.sim-input {
  width: 100%; background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.09);
  border-bottom: 2px solid rgba(16,185,129,0.2); border-radius: 2px 2px 0 0;
  color: #fff; font-size: 0.8rem; font-weight: 700; font-family: inherit;
  padding: 0.65rem 2.5rem 0.65rem 0.75rem; outline: none;
  transition: border-color 0.2s, background 0.2s;
  user-select: text;
}
.sim-input:focus { background: rgba(255,255,255,0.06); border-bottom-color: #10b981; }
.sim-input::placeholder { color: rgba(255,255,255,0.2); font-weight: 500; }
.char-count {
  position: absolute; right: 0.6rem; top: 50%; transform: translateY(25%);
  font-size: 0.55rem; color: rgba(255,255,255,0.18); font-weight: 600;
  pointer-events: none;
}

/* NAV BAR */
.nav-bar {
  flex-shrink: 0; display: grid; grid-template-columns: 1fr 1.6fr;
  gap: 0.75rem; padding: 1rem 1.25rem;
  background: rgba(3,9,5,0.85); border-top: 1px solid rgba(255,255,255,0.05);
  backdrop-filter: blur(8px);
}
.btn-nav {
  display: flex; align-items: center; justify-content: center; gap: 0.4rem;
  font-size: 0.75rem; font-weight: 900; letter-spacing: 0.12em;
  padding: 0.875rem; border-radius: 2px; cursor: pointer; border: none;
  transition: all 0.2s;
}
.btn-nav svg { width: 0.9rem; height: 0.9rem; flex-shrink: 0; }
.btn-back-step {
  background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.5);
}
.btn-back-step:hover { background: rgba(255,255,255,0.08); color: #fff; border-color: rgba(255,255,255,0.2); }
.btn-next-step {
  background: #10b981; color: #000;
  box-shadow: 0 0 20px rgba(16,185,129,0.3);
}
.btn-next-step:hover { background: #34d399; box-shadow: 0 0 35px rgba(16,185,129,0.5); transform: translateY(-1px); }
.btn-next-step:active { transform: scale(0.97); }

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   REVIEW SCREEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
.review-screen {
  display: flex; align-items: flex-start; justify-content: center;
  background: linear-gradient(160deg, #030d06 0%, #050f07 100%);
  overflow-y: auto; padding: 2rem 1rem;
}
.review-inner {
  max-width: 900px; width: 100%;
  display: flex; flex-direction: column; gap: 1.75rem;
  animation: splashEnter 0.5s cubic-bezier(0.22, 1, 0.36, 1);
}

.review-header {
  display: flex; align-items: center; gap: 1.5rem;
}
.review-eyebrow {
  display: block; font-size: 0.6rem; font-weight: 900; letter-spacing: 0.2em;
  color: #f59e0b; margin-bottom: 0.25rem;
}
.review-title {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(2rem, 5vw, 3rem); letter-spacing: 0.06em; line-height: 1; color: #fff;
}

.review-content {
  display: flex; gap: 2rem; align-items: flex-start;
  flex-wrap: wrap;
}
.review-board-wrap {
  flex: 0 0 auto; width: clamp(200px, 40%, 300px);
  animation: bounceIn 0.7s cubic-bezier(0.22, 1, 0.36, 1) both;
}
@keyframes bounceIn {
  0%   { transform: scale(0.85) translateY(20px); opacity: 0; }
  60%  { transform: scale(1.03) translateY(-3px); opacity: 1; }
  100% { transform: none; }
}
.review-board { min-height: 200px; }

.review-specs {
  flex: 1; min-width: 260px;
  background: rgba(255,255,255,0.03); border: 1px solid rgba(255,255,255,0.07);
  border-radius: 2px; padding: 1.25rem; display: flex; flex-direction: column; gap: 0;
}
.specs-title {
  font-size: 0.6rem; font-weight: 900; letter-spacing: 0.2em;
  color: rgba(255,255,255,0.3); margin-bottom: 1rem;
  padding-bottom: 0.75rem; border-bottom: 1px solid rgba(255,255,255,0.06);
}
.spec-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 0.55rem 0; border-bottom: 1px solid rgba(255,255,255,0.04);
  gap: 1rem;
}
.spec-label {
  font-size: 0.68rem; font-weight: 700; color: rgba(255,255,255,0.4);
  flex-shrink: 0;
}
.spec-value {
  font-size: 0.72rem; font-weight: 700; color: rgba(255,255,255,0.8);
  text-align: right; word-break: break-word;
}
.spec-value-row { display: flex; align-items: center; gap: 0.4rem; }
.spec-color-dot { width: 12px; height: 12px; border-radius: 2px; flex-shrink: 0; border: 1px solid rgba(255,255,255,0.15); }

.spec-total-row {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: 0.75rem; padding-top: 0.75rem;
  border-top: 2px solid rgba(16,185,129,0.3);
}
.spec-total-label {
  font-size: 0.65rem; font-weight: 900; letter-spacing: 0.15em;
  color: rgba(255,255,255,0.5);
}
.spec-total-value {
  font-family: 'Bebas Neue', sans-serif;
  font-size: 1.75rem; color: #10b981; letter-spacing: 0.03em;
  text-shadow: 0 0 20px rgba(16,185,129,0.4);
}

.review-cta {
  display: flex; flex-direction: column; gap: 0.6rem; margin-top: 1rem;
}
.btn-add-cart {
  width: 100%; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #000; font-size: 0.8rem; font-weight: 900; letter-spacing: 0.1em;
  padding: 1rem; border-radius: 2px; border: none; cursor: pointer;
  box-shadow: 0 0 30px rgba(245,158,11,0.35);
  transition: all 0.25s;
}
.btn-add-cart-inner { display: flex; align-items: center; gap: 0.5rem; }
.btn-add-cart-inner svg { width: 1rem; height: 1rem; }
.btn-add-cart:not(.btn-loading):hover {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  box-shadow: 0 0 50px rgba(245,158,11,0.55); transform: translateY(-2px);
}
.btn-add-cart.btn-loading { opacity: 0.7; cursor: not-allowed; }
.btn-redesign {
  width: 100%; background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.08);
  color: rgba(255,255,255,0.4); font-size: 0.72rem; font-weight: 700;
  padding: 0.7rem; border-radius: 2px; cursor: pointer;
  transition: all 0.2s;
}
.btn-redesign:hover { background: rgba(255,255,255,0.07); color: rgba(255,255,255,0.7); }

.spin-icon { animation: spinAnim 1s linear infinite; }
@keyframes spinAnim { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

/* ─── CONFETTI ────────────────────────────────────────────────── */
.confetti-layer { position: fixed; inset: 0; pointer-events: none; z-index: 5; }
.confetti-piece {
  position: absolute; top: -20px;
  border-radius: 2px;
  animation: confettiFall linear both;
}
@keyframes confettiFall {
  0%   { transform: translateY(0) rotate(var(--rot, 45deg)); opacity: 1; }
  80%  { opacity: 1; }
  100% { transform: translateY(110vh) rotate(calc(var(--rot, 45deg) + 540deg)); opacity: 0; }
}

/* ─── RESPONSIVE ──────────────────────────────────────────────── */
@media (max-width: 900px) {
  .builder-split { flex-direction: column; }
  .canvas-panel { min-height: 45vh; border-right: none; border-bottom: 1px solid rgba(255,255,255,0.05); padding: 1.5rem 1rem; }
  .options-panel { width: 100%; height: 60vh; }
  .board-wrap { max-height: 40vh; }
  .board-frame { min-height: 180px; }
  .price-hud { bottom: 0.75rem; right: 0.75rem; padding: 0.5rem 0.85rem; }
  .price-breakdown-mini { display: none; }
  .hud-center { position: static; transform: none; }
  .hud-bar { gap: 0.5rem; flex-wrap: wrap; height: auto; padding: 0.5rem 1rem; }
  .hud-kbd-hint { display: none; }
}

@media (max-width: 640px) {
  .tut-cards { grid-template-columns: 1fr; }
  .splash-title { font-size: clamp(3rem, 14vw, 5rem); }
  .review-content { flex-direction: column; align-items: center; }
  .review-board-wrap { width: 100%; max-width: 280px; }
  .hud-center { display: none; }
}

/* Custom Dark-Themed Toast Notification for Simulator */
.cart-toast-custom {
  position: fixed;
  top: 90px;
  right: 24px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 360px;
  max-width: calc(100vw - 48px);
  background: rgba(16, 24, 18, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5);
  font-family: 'Inter', system-ui, sans-serif;
  color: #ffffff;
}

.toast-body-custom {
  display: flex;
  align-items: center;
  gap: 14px;
}

.toast-icon-check-custom {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
  flex-shrink: 0;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.toast-img-wrap-custom {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  overflow: hidden;
  background: #111;
  border: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.toast-img-custom {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.toast-details-custom {
  flex-grow: 1;
  min-width: 0;
}

.toast-title-custom {
  font-size: 14px;
  font-weight: 800;
  color: #10b981;
  margin: 0;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.toast-name-custom {
  font-size: 13px;
  font-weight: 600;
  color: #ffffff;
  margin: 4px 0 0 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.toast-meta-custom {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  margin: 2px 0 0 0;
  font-weight: 500;
}

.toast-actions-custom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px dashed rgba(16, 185, 129, 0.2);
  padding-top: 12px;
}

.btn-toast-view-custom {
  font-size: 11px;
  font-weight: 700;
  color: #000000;
  text-decoration: none;
  background: #10b981;
  padding: 6px 14px;
  border-radius: 4px;
  transition: all 0.2s ease;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.btn-toast-view-custom:hover {
  background: #34d399;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.4);
}

.btn-toast-close-custom {
  background: transparent;
  border: none;
  font-size: 18px;
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.2s ease;
  line-height: 1;
}

.btn-toast-close-custom:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
}

/* Toast Transition for Custom */
.toast-cust-enter-active, .toast-cust-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-cust-enter-from {
  transform: translateX(120%) scale(0.9);
  opacity: 0;
}
.toast-cust-leave-to {
  transform: translateX(120%) scale(0.9);
  opacity: 0;
}
</style>