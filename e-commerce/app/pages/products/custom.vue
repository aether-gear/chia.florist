<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'

definePageMeta({ layout: false })
useHead({
  title: 'Chia Florist — Board Designer',
  meta: [{ name: 'description', content: 'Design your custom flower board with our interactive canvas designer.' }]
})

const { addToCart, formatRupiah } = useCart()

/* ─── TYPES ───────────────────────────────────────────────────────── */
type FontId      = 'inter' | 'playfair' | 'dancing' | 'bebas' | 'merriweather' | 'pacifico'
type CornerStyle = 'none' | 'rounded' | 'cut' | 'ornate' | 'floral'
type FrameStyle  = 'none' | 'square' | 'circle'
type BrushType   = 'flower' | 'rose'
type BorderStyle = 'none' | 'solid' | 'double' | 'dashed' | 'dotted' | 'groove' | 'ridge' | 'ornate'
type ToolTab     = 'text' | 'image' | 'brush' | 'border' | 'corner'
type SectionKey  = 'upper' | 'lower'

interface BoardSection {
  headerText: string; bodyText: string
  headerFontSize: number; bodyFontSize: number
  headerFont: FontId; bodyFont: FontId
  headerAlign: 'left' | 'center' | 'right'
  bodyAlign: 'left' | 'center' | 'right'
  bgColor: string; headerColor: string; bodyColor: string
  cornerStyle: CornerStyle
}
interface CanvasImage {
  id: string; type: 'image'; src: string; frame: FrameStyle
  x: number; y: number; width: number; zoom: number; cropX: number; cropY: number
}
interface BrushStroke {
  id: string; type: 'brush'; brushType: BrushType
  x: number; y: number; size: number; color: string; rotation: number
}
type CanvasElement = CanvasImage | BrushStroke
interface BoardBorder { style: BorderStyle; color: string; width: number }

/* ─── CONSTANTS ───────────────────────────────────────────────────── */
const BOARD_W = 800
const BOARD_H = 534

const FONTS: { id: FontId; label: string; family: string }[] = [
  { id: 'inter',        label: 'Inter',        family: "'Inter', sans-serif" },
  { id: 'playfair',     label: 'Playfair',     family: "'Playfair Display', serif" },
  { id: 'dancing',      label: 'Dancing',      family: "'Dancing Script', cursive" },
  { id: 'bebas',        label: 'Bebas',        family: "'Bebas Neue', sans-serif" },
  { id: 'merriweather', label: 'Merriweather', family: "'Merriweather', serif" },
  { id: 'pacifico',     label: 'Pacifico',     family: "'Pacifico', cursive" },
]

const CORNERS: { id: CornerStyle; label: string }[] = [
  { id: 'none', label: 'None' }, { id: 'rounded', label: 'Rounded' },
  { id: 'cut', label: 'Cut' }, { id: 'ornate', label: 'Ornate' }, { id: 'floral', label: 'Floral' },
]

const BORDER_STYLES: { id: BorderStyle; label: string }[] = [
  { id: 'none', label: 'None' }, { id: 'solid', label: 'Solid' }, { id: 'double', label: 'Double' },
  { id: 'dashed', label: 'Dashed' }, { id: 'dotted', label: 'Dotted' },
  { id: 'groove', label: 'Groove' }, { id: 'ridge', label: 'Ridge' }, { id: 'ornate', label: 'Ornate' },
]

const SIZES = [
  { id: 'small',  label: '1.5 × 2.0m', price: 150_000, desc: 'Compact' },
  { id: 'medium', label: '1.8 × 2.5m', price: 200_000, desc: 'Standard', recommended: true },
  { id: 'large',  label: '2.0 × 3.0m', price: 280_000, desc: 'Grand' },
]

const BRUSH_COLORS  = ['#e85d75','#f4845f','#f9c74f','#90be6d','#4cc9f0','#c77dff','#ffffff','#222222']
const BORDER_COLORS = ['#f5c842','#e63946','#2a9d8f','#264653','#e76f51','#a8dadc','#f1faee','#1d3557']
const BG_PRESETS    = ['#c0392b','#1a3a5c','#145a32','#6c3483','#a04000','#17202a','#f0f0f0','#ffffff']

const TOOL_TABS: { id: ToolTab; label: string }[] = [
  { id: 'text', label: 'Text' }, { id: 'image', label: 'Image' },
  { id: 'brush', label: 'Brush' }, { id: 'border', label: 'Border' }, { id: 'corner', label: 'Corner' },
]

/* ─── STATE ───────────────────────────────────────────────────────── */
const upper = ref<BoardSection>({
  headerText: 'Selamat & Sukses', bodyText: 'Atas Pelantikan Saudara/i\nNama Lengkap Anda',
  headerFontSize: 36, bodyFontSize: 20, headerFont: 'playfair', bodyFont: 'inter',
  headerAlign: 'center', bodyAlign: 'center',
  bgColor: '#c0392b', headerColor: '#ffd700', bodyColor: '#ffffff', cornerStyle: 'none',
})
const lower = ref<BoardSection>({
  headerText: '', bodyText: 'Nama Pengirim\nNama Instansi / Perusahaan',
  headerFontSize: 26, bodyFontSize: 22, headerFont: 'bebas', bodyFont: 'inter',
  headerAlign: 'center', bodyAlign: 'center',
  bgColor: '#1a3a5c', headerColor: '#ffffff', bodyColor: '#ffffff', cornerStyle: 'none',
})

const heightRatio  = ref(0.58)
const border       = ref<BoardBorder>({ style: 'solid', color: '#f5c842', width: 12 })
const elements     = ref<CanvasElement[]>([])

const activeTab     = ref<ToolTab>('text')
const activeSection = ref<SectionKey>('upper')
const selectedId    = ref<string | null>(null)
const physicalSize  = ref('medium')
const showReview    = ref(false)
const showToast     = ref(false)
const isAdding      = ref(false)

const brushType     = ref<BrushType>('flower')
const brushColor    = ref('#e85d75')
const brushSize     = ref(48)
const brushRotation = ref(0)
const isBrushMode   = computed(() => activeTab.value === 'brush')

const containerRef = ref<HTMLElement | null>(null)
const boardRef     = ref<HTMLElement | null>(null)
const boardScale   = ref(0.75)

/* ─── HELPERS ─────────────────────────────────────────────────────── */
const getFont  = (id: FontId) => FONTS.find(f => f.id === id)?.family ?? "'Inter', sans-serif"
const sec      = computed(() => activeSection.value === 'upper' ? upper.value : lower.value)
const upperH   = computed(() => Math.round(heightRatio.value * BOARD_H))
const lowerH   = computed(() => BOARD_H - upperH.value)

const boardBorderStyle = computed((): Record<string, string> => {
  const { style, color, width } = border.value
  if (style === 'none' || width === 0) return {}
  if (style === 'ornate') return { border: `${width}px solid ${color}`, outline: `${Math.max(2, Math.round(width * 0.35))}px solid ${color}`, outlineOffset: '5px' }
  return { border: `${width}px ${style} ${color}` }
})

const upperCornerStyle = computed((): Record<string, string> => {
  const s = upper.value.cornerStyle
  if (s === 'rounded') return { borderTopLeftRadius: '12px', borderTopRightRadius: '12px' }
  if (s === 'cut') return { clipPath: 'polygon(14px 0%,calc(100% - 14px) 0%,100% 14px,100% 100%,0% 100%,0% 14px)' }
  if (s === 'ornate') return { borderTopLeftRadius: '4px 18px', borderTopRightRadius: '4px 18px' }
  return {}
})

const lowerCornerStyle = computed((): Record<string, string> => {
  const s = lower.value.cornerStyle
  if (s === 'rounded') return { borderBottomLeftRadius: '12px', borderBottomRightRadius: '12px' }
  if (s === 'cut') return { clipPath: 'polygon(0% 0%,100% 0%,100% calc(100% - 14px),calc(100% - 14px) 100%,14px 100%,0% calc(100% - 14px))' }
  if (s === 'ornate') return { borderBottomLeftRadius: '4px 18px', borderBottomRightRadius: '4px 18px' }
  return {}
})

const selectedEl   = computed(() => elements.value.find(e => e.id === selectedId.value) ?? null)
const selectedImg  = computed(() => (selectedEl.value?.type === 'image'  ? selectedEl.value : null) as CanvasImage | null)
const imgElements  = computed(() => elements.value.filter(e => e.type === 'image') as CanvasImage[])
const brushElements = computed(() => elements.value.filter(e => e.type === 'brush') as BrushStroke[])
const totalPrice   = computed(() => SIZES.find(s => s.id === physicalSize.value)?.price ?? 200_000)

/* ─── SCALE ───────────────────────────────────────────────────────── */
const updateScale = () => {
  const el = containerRef.value
  if (!el) return
  const pad = 64
  boardScale.value = Math.max(0.25, Math.min((el.offsetWidth - pad) / BOARD_W, (el.offsetHeight - pad) / BOARD_H, 1.1))
}

/* ─── Z-ORDER (Instagram story: last in array = top layer) ────────── */
const bringToFront = (id: string) => {
  const idx = elements.value.findIndex(e => e.id === id)
  if (idx < 0) return
  const [el] = elements.value.splice(idx, 1)
  elements.value.push(el!)
  selectedId.value = id
}

const deleteSelected = () => {
  if (!selectedId.value) return
  elements.value = elements.value.filter(e => e.id !== selectedId.value)
  selectedId.value = null
}

/* ─── IMAGE ───────────────────────────────────────────────────────── */
const makeImage = (src: string, x = 15, y = 15): CanvasImage => ({
  id: 'img-' + Date.now().toString(36) + Math.random().toString(36).slice(2, 5),
  type: 'image', src, frame: 'square', x, y, width: 22, zoom: 1, cropX: 50, cropY: 50,
})

const readFile = (file: File, x: number, y: number) => {
  const reader = new FileReader()
  reader.onload = ev => {
    const img = makeImage(ev.target?.result as string, x, y)
    elements.value.push(img)
    selectedId.value = img.id
    activeTab.value = 'image'
  }
  reader.readAsDataURL(file)
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  const file = e.dataTransfer?.files[0]
  if (!file?.type.startsWith('image/')) return
  const board = boardRef.value; if (!board) return
  const r = board.getBoundingClientRect()
  readFile(file, Math.max(2, Math.min(((e.clientX - r.left) / r.width) * 100 - 11, 76)), Math.max(2, Math.min(((e.clientY - r.top) / r.height) * 100 - 11, 76)))
}

const handleFileInput = (e: Event) => {
  const f = (e.target as HTMLInputElement).files?.[0]; if (!f) return
  readFile(f, 15, 15);
  (e.target as HTMLInputElement).value = ''
}

/* ─── BRUSH PLACEMENT ─────────────────────────────────────────────── */
const handleBoardClick = (e: MouseEvent) => {
  if (isBrushMode.value) {
    const r = boardRef.value!.getBoundingClientRect()
    const stroke: BrushStroke = {
      id: 'br-' + Date.now().toString(36) + Math.random().toString(36).slice(2, 5),
      type: 'brush', brushType: brushType.value,
      x: ((e.clientX - r.left) / r.width) * 100,
      y: ((e.clientY - r.top)  / r.height) * 100,
      size: brushSize.value, color: brushColor.value, rotation: brushRotation.value,
    }
    elements.value.push(stroke)
    selectedId.value = stroke.id
  } else {
    selectedId.value = null
  }
}

/* ─── DRAG ────────────────────────────────────────────────────────── */
let _draggingEl = false, _draggingDiv = false, _dragElId = ''
let _rect = { left: 0, top: 0, width: BOARD_W * 0.75, height: BOARD_H * 0.75 }
let _dragBX = 0, _dragBY = 0, _dragElX0 = 0, _dragElY0 = 0
let _divStartY = 0, _divStartR = 0

const startDragEl = (e: MouseEvent, id: string) => {
  if (isBrushMode.value) return
  e.stopPropagation(); e.preventDefault()
  bringToFront(id)
  const board = boardRef.value; if (!board) return
  _rect = board.getBoundingClientRect(); _draggingEl = true; _dragElId = id
  _dragBX = (e.clientX - _rect.left) / _rect.width  * 100
  _dragBY = (e.clientY - _rect.top)  / _rect.height * 100
  const el = elements.value.find(e => e.id === id)
  if (el) { _dragElX0 = el.x; _dragElY0 = el.y }
}

const startDragDiv = (e: MouseEvent) => {
  e.stopPropagation(); e.preventDefault()
  const board = boardRef.value; if (!board) return
  _rect = board.getBoundingClientRect(); _draggingDiv = true
  _divStartY = e.clientY; _divStartR = heightRatio.value
}

const onMouseMove = (e: MouseEvent) => {
  if (_draggingEl) {
    const el = elements.value.find(ev => ev.id === _dragElId)
    if (el) {
      el.x = Math.max(-5, Math.min(_dragElX0 + ((e.clientX - _rect.left) / _rect.width  * 100 - _dragBX), 95))
      el.y = Math.max(-5, Math.min(_dragElY0 + ((e.clientY - _rect.top)  / _rect.height * 100 - _dragBY), 95))
    }
  }
  if (_draggingDiv) {
    heightRatio.value = Math.max(0.25, Math.min(_divStartR + (e.clientY - _divStartY) / _rect.height, 0.75))
  }
}

const onMouseUp = () => { _draggingEl = false; _draggingDiv = false }

/* ─── CART ────────────────────────────────────────────────────────── */
const addToCartHandler = async () => {
  isAdding.value = true
  await new Promise(r => setTimeout(r, 800))
  addToCart({
    id: 'custom-' + Date.now(),
    name: `Custom Board — ${upper.value.headerText || 'My Design'}`,
    price: totalPrice.value,
    image: '/images/custom-preview.png',
    size: SIZES.find(s => s.id === physicalSize.value)?.label ?? '',
    color: upper.value.bgColor,
    shopId: '99ef0062-1040-4574-a4be-0123abce5670',
    isCustom: true,
  }, 1)
  isAdding.value = false; showReview.value = false; showToast.value = true
  setTimeout(() => { showToast.value = false; navigateTo('/') }, 3200)
}

/* ─── KEYBOARD ────────────────────────────────────────────────────── */
const onKeyDown = (e: KeyboardEvent) => {
  const tag = (document.activeElement as HTMLElement)?.tagName
  if (tag === 'INPUT' || tag === 'TEXTAREA') return
  if (e.key === 'Delete' || e.key === 'Backspace') deleteSelected()
  if (e.key === 'Escape') selectedId.value = null
}

/* ─── LIFECYCLE ───────────────────────────────────────────────────── */
let ro: ResizeObserver | null = null
onMounted(() => {
  updateScale()
  ro = new ResizeObserver(updateScale)
  if (containerRef.value) ro.observe(containerRef.value)
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
  window.addEventListener('keydown', onKeyDown)
})
onUnmounted(() => {
  ro?.disconnect()
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
  window.removeEventListener('keydown', onKeyDown)
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
        <span class="dr-page-title">Board Designer</span>
      </div>
      <div class="dr-nav-right">
        <span class="dr-scale-chip">{{ Math.round(boardScale * 100) }}%</span>
        <button id="btn-finalize-nav" class="dr-finalize-btn" @click="showReview = true">
          Finalize &amp; Order
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
        </button>
      </div>
    </nav>

    <!-- ═══ BODY ══════════════════════════════════════════════════════ -->
    <div class="dr-body">

      <!-- ── CANVAS AREA (left on md+) ──────────────────────────────── -->
      <div class="dr-canvas-area" ref="containerRef">
        <div class="chess-bg" @dragover.prevent>
          <!-- Board scaler: sized to match scale -->
          <div class="board-scaler"
            :style="{ width: Math.round(BOARD_W * boardScale) + 'px', height: Math.round(BOARD_H * boardScale) + 'px' }">

            <!-- ╔══ THE BOARD ══╗ -->
            <div
              ref="boardRef"
              class="board-frame"
              :style="{
                width: BOARD_W + 'px', height: BOARD_H + 'px',
                transform: `scale(${boardScale})`, transformOrigin: 'top left',
                cursor: isBrushMode ? 'crosshair' : 'default',
                ...boardBorderStyle,
              }"
              @click="handleBoardClick"
              @dragover.prevent
              @drop="handleDrop"
            >
              <!-- ▲ UPPER SECTION -->
              <div class="board-section"
                :style="{ top: 0, left: 0, right: 0, height: upperH + 'px', backgroundColor: upper.bgColor, ...upperCornerStyle }">
                <template v-if="upper.cornerStyle === 'floral'">
                  <span class="floral-corner fc-tl">🌸</span>
                  <span class="floral-corner fc-tr">🌸</span>
                </template>
                <div class="sec-inner">
                  <div v-if="upper.headerText" class="sec-text sec-header"
                    :style="{ fontSize: upper.headerFontSize + 'px', fontFamily: getFont(upper.headerFont), color: upper.headerColor, textAlign: upper.headerAlign }">
                    {{ upper.headerText }}
                  </div>
                  <div class="sec-text sec-body"
                    :style="{ fontSize: upper.bodyFontSize + 'px', fontFamily: getFont(upper.bodyFont), color: upper.bodyColor, textAlign: upper.bodyAlign }">
                    {{ upper.bodyText }}
                  </div>
                </div>
              </div>

              <!-- ── SECTION DIVIDER (drag to resize) ── -->
              <div class="section-divider" :style="{ top: upperH + 'px', '--dc': border.color || '#ccc' }"
                @mousedown.stop="startDragDiv" @click.stop>
                <div class="div-track"/>
                <div class="div-knob">
                  <svg viewBox="0 0 20 6" width="20" height="6">
                    <circle cx="3" cy="3" r="2" fill="currentColor"/>
                    <circle cx="10" cy="3" r="2" fill="currentColor"/>
                    <circle cx="17" cy="3" r="2" fill="currentColor"/>
                  </svg>
                </div>
              </div>

              <!-- ▼ LOWER SECTION -->
              <div class="board-section"
                :style="{ top: upperH + 'px', left: 0, right: 0, height: lowerH + 'px', backgroundColor: lower.bgColor, ...lowerCornerStyle }">
                <template v-if="lower.cornerStyle === 'floral'">
                  <span class="floral-corner fc-bl">🌸</span>
                  <span class="floral-corner fc-br">🌸</span>
                </template>
                <div class="sec-inner">
                  <div v-if="lower.headerText" class="sec-text sec-header"
                    :style="{ fontSize: lower.headerFontSize + 'px', fontFamily: getFont(lower.headerFont), color: lower.headerColor, textAlign: lower.headerAlign }">
                    {{ lower.headerText }}
                  </div>
                  <div class="sec-text sec-body"
                    :style="{ fontSize: lower.bodyFontSize + 'px', fontFamily: getFont(lower.bodyFont), color: lower.bodyColor, textAlign: lower.bodyAlign }">
                    {{ lower.bodyText }}
                  </div>
                </div>
              </div>

              <!-- ★ CANVAS ELEMENTS — array order = z-order (last = top) -->
              <template v-for="(el, idx) in elements" :key="el.id">

                <!-- Image element -->
                <div v-if="el.type === 'image'" class="canvas-el"
                  :class="{ 'el-selected': selectedId === el.id }"
                  :style="{
                    left: (el as CanvasImage).x + '%', top: (el as CanvasImage).y + '%',
                    width: (el as CanvasImage).width + '%', aspectRatio: '1/1',
                    zIndex: idx + 10, overflow: 'hidden',
                    borderRadius: (el as CanvasImage).frame === 'circle' ? '50%' : (el as CanvasImage).frame === 'square' ? '4px' : '0',
                    pointerEvents: isBrushMode ? 'none' : 'auto', cursor: 'grab',
                  }"
                  @mousedown.stop="startDragEl($event, el.id)">
                  <img :src="(el as CanvasImage).src" draggable="false"
                    :style="{
                      width: '100%', height: '100%', objectFit: 'cover', display: 'block',
                      objectPosition: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                      transform: 'scale(' + (el as CanvasImage).zoom + ')',
                      transformOrigin: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                    }"/>
                  <div v-if="selectedId === el.id && !isBrushMode" class="el-del" @click.stop="deleteSelected" title="Remove">×</div>
                </div>

                <!-- Brush stroke element -->
                <div v-else-if="el.type === 'brush'" class="canvas-el"
                  :class="{ 'el-selected': selectedId === el.id }"
                  :style="{
                    left: (el as BrushStroke).x + '%', top: (el as BrushStroke).y + '%',
                    width: (el as BrushStroke).size + 'px', height: (el as BrushStroke).size + 'px',
                    transform: `translate(-50%,-50%) rotate(${(el as BrushStroke).rotation}deg)`,
                    zIndex: idx + 10, color: (el as BrushStroke).color,
                    pointerEvents: isBrushMode ? 'none' : 'auto', cursor: 'pointer',
                  }"
                  @mousedown.stop="!isBrushMode && startDragEl($event, el.id)">
                  <!-- Flower SVG -->
                  <svg v-if="(el as BrushStroke).brushType === 'flower'" viewBox="-20 -20 40 40" width="100%" height="100%">
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" opacity="0.92" transform="rotate(0,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" opacity="0.92" transform="rotate(72,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" opacity="0.92" transform="rotate(144,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" opacity="0.92" transform="rotate(216,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" opacity="0.92" transform="rotate(288,0,0)"/>
                    <circle cx="0" cy="0" r="5.5" fill="#ffe066"/>
                    <circle cx="0" cy="0" r="2.5" fill="#f59e0b"/>
                  </svg>
                  <!-- Rose SVG -->
                  <svg v-else viewBox="-20 -20 40 40" width="100%" height="100%">
                    <path d="M0,-16 C6,-13 14,-7 12,0 C18,-4 20,5 15,10 C10,15 3,16 0,13 C-3,16 -10,15 -15,10 C-20,5 -18,-4 -12,0 C-14,-7 -6,-13 0,-16Z" fill="currentColor"/>
                    <path d="M0,-9 C4,-6 8,-1 6,3 C9,0 11,5 8,8 C5,11 1,11 0,9 C-1,11 -5,11 -8,8 C-11,5 -9,0 -6,3 C-8,-1 -4,-6 0,-9Z" fill="currentColor" opacity="0.55"/>
                    <circle cx="0" cy="2" r="3" fill="#ffe066" opacity="0.75"/>
                  </svg>
                  <div v-if="selectedId === el.id && !isBrushMode" class="el-del" @click.stop="deleteSelected" title="Remove">×</div>
                </div>
              </template>

              <!-- Brush mode overlay -->
              <div v-if="isBrushMode" class="brush-overlay" @click.stop>
                <span>{{ brushType === 'flower' ? '🌸' : '🌹' }} Click anywhere to place</span>
              </div>

            </div><!-- /board-frame -->
          </div><!-- /board-scaler -->
        </div><!-- /chess-bg -->

        <!-- Canvas info bar -->
        <div class="canvas-info">
          <span>{{ elements.length }} element{{ elements.length !== 1 ? 's' : '' }}</span>
          <span v-if="selectedId" class="ci-sel"> · 1 selected</span>
          <button v-if="selectedId" class="ci-desel" @click="selectedId = null">Deselect</button>
          <span v-if="isBrushMode" class="ci-hint">Click canvas to place brush stroke · Del to remove selected</span>
          <span v-else-if="!isBrushMode && selectedId" class="ci-hint">Drag to move · Del to remove</span>
        </div>
      </div><!-- /canvas-area -->

      <!-- ── TOOL PANEL (right on md+) ──────────────────────────────── -->
      <aside class="dr-panel">

        <!-- Tab bar -->
        <div class="tab-bar" role="tablist">
          <button v-for="tab in TOOL_TABS" :key="tab.id"
            class="tab-btn" :class="{ 'tab-active': activeTab === tab.id }"
            @click="activeTab = tab.id" role="tab" :aria-selected="activeTab === tab.id"
            :id="'tab-' + tab.id">
            <!-- Tab icons -->
            <svg v-if="tab.id === 'text'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 7V4h16v3"/><path d="M9 20h6"/><path d="M12 4v16"/>
            </svg>
            <svg v-else-if="tab.id === 'image'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
            </svg>
            <svg v-else-if="tab.id === 'brush'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18.37 2.63 14 7l-1.59-1.59a2 2 0 0 0-2.82 0L8 7l9 9 1.59-1.59a2 2 0 0 0 0-2.82L17 10l4.37-4.37a2.12 2.12 0 1 0-3-3Z"/>
              <path d="M9 8c-2 2.5-2 5-2 5"/>
            </svg>
            <svg v-else-if="tab.id === 'border'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="2" width="20" height="20" rx="2"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
              <path d="M3 9V5a2 2 0 012-2h4"/><path d="M3 15v4a2 2 0 002 2h4"/>
            </svg>
            <span>{{ tab.label }}</span>
          </button>
        </div>

        <!-- Tab content area (scrollable) -->
        <div class="tab-body">

          <!-- ╔══ TEXT TAB ══╗ -->
          <div v-if="activeTab === 'text'" class="tab-pane">
            <div class="sec-toggle">
              <button class="stg-btn" :class="{ active: activeSection === 'upper' }" @click="activeSection = 'upper'">Upper</button>
              <button class="stg-btn" :class="{ active: activeSection === 'lower' }" @click="activeSection = 'lower'">Lower</button>
            </div>

            <!-- Background -->
            <div class="tg">
              <div class="tg-label">BACKGROUND</div>
              <div class="color-row">
                <input id="bg-color-input" type="color" :value="sec.bgColor" @input="sec.bgColor = ($event.target as HTMLInputElement).value" class="csi"/>
                <span class="cval">{{ sec.bgColor }}</span>
              </div>
              <div class="dot-row">
                <button v-for="c in BG_PRESETS" :key="c" class="pdot"
                  :style="{ background: c, outline: sec.bgColor === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#ffffff' ? 'inset 0 0 0 1px #ccc' : 'none' }"
                  @click="sec.bgColor = c"/>
              </div>
            </div>

            <!-- Header -->
            <div class="tg">
              <div class="tg-label">HEADER</div>
              <textarea class="dr-ta" :value="sec.headerText" @input="sec.headerText = ($event.target as HTMLTextAreaElement).value" placeholder="Header text…" rows="2"/>
              <div class="cr">
                <label class="clabel">Size</label>
                <input type="range" min="10" max="96" class="dr-range" :value="sec.headerFontSize" @input="sec.headerFontSize = +($event.target as HTMLInputElement).value"/>
                <span class="cval">{{ sec.headerFontSize }}px</span>
              </div>
              <div class="font-grid">
                <button v-for="f in FONTS" :key="f.id" class="font-chip" :class="{ 'fc-active': sec.headerFont === f.id }"
                  :style="{ fontFamily: f.family }" @click="sec.headerFont = f.id">{{ f.label }}</button>
              </div>
              <div class="align-row">
                <button v-for="a in (['left','center','right'] as const)" :key="a" class="aln-btn"
                  :class="{ 'aln-active': sec.headerAlign === a }" @click="sec.headerAlign = a"
                  :title="a">
                  <svg viewBox="0 0 16 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
                    <line x1="0" y1="2" x2="16" y2="2"/>
                    <line :x1="a==='left'?0:a==='center'?3:6" y1="6" :x2="a==='left'?10:a==='center'?13:16" y2="6"/>
                    <line x1="0" y1="10" x2="16" y2="10"/>
                  </svg>
                </button>
              </div>
              <div class="color-row">
                <label class="clabel">Color</label>
                <input type="color" :value="sec.headerColor" @input="sec.headerColor = ($event.target as HTMLInputElement).value" class="csi"/>
                <span class="cval">{{ sec.headerColor }}</span>
              </div>
            </div>

            <!-- Body -->
            <div class="tg">
              <div class="tg-label">BODY</div>
              <textarea class="dr-ta" :value="sec.bodyText" @input="sec.bodyText = ($event.target as HTMLTextAreaElement).value" placeholder="Body text… (new line = Enter)" rows="3"/>
              <div class="cr">
                <label class="clabel">Size</label>
                <input type="range" min="8" max="72" class="dr-range" :value="sec.bodyFontSize" @input="sec.bodyFontSize = +($event.target as HTMLInputElement).value"/>
                <span class="cval">{{ sec.bodyFontSize }}px</span>
              </div>
              <div class="font-grid">
                <button v-for="f in FONTS" :key="f.id" class="font-chip" :class="{ 'fc-active': sec.bodyFont === f.id }"
                  :style="{ fontFamily: f.family }" @click="sec.bodyFont = f.id">{{ f.label }}</button>
              </div>
              <div class="align-row">
                <button v-for="a in (['left','center','right'] as const)" :key="a" class="aln-btn"
                  :class="{ 'aln-active': sec.bodyAlign === a }" @click="sec.bodyAlign = a" :title="a">
                  <svg viewBox="0 0 16 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
                    <line x1="0" y1="2" x2="16" y2="2"/>
                    <line :x1="a==='left'?0:a==='center'?3:6" y1="6" :x2="a==='left'?10:a==='center'?13:16" y2="6"/>
                    <line x1="0" y1="10" x2="16" y2="10"/>
                  </svg>
                </button>
              </div>
              <div class="color-row">
                <label class="clabel">Color</label>
                <input type="color" :value="sec.bodyColor" @input="sec.bodyColor = ($event.target as HTMLInputElement).value" class="csi"/>
                <span class="cval">{{ sec.bodyColor }}</span>
              </div>
            </div>
          </div>

          <!-- ╔══ IMAGE TAB ══╗ -->
          <div v-else-if="activeTab === 'image'" class="tab-pane">
            <!-- Drop zone (no image selected) -->
            <template v-if="!selectedImg">
              <div class="drop-zone" @dragover.prevent @drop="handleDrop">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="dz-icon">
                  <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
                </svg>
                <p class="dz-title">Drop image here</p>
                <p class="dz-sub">or drag from files onto the board</p>
                <label class="dz-browse" for="dz-file-input">Browse Files</label>
                <input id="dz-file-input" type="file" accept="image/*" @change="handleFileInput" style="display:none"/>
              </div>
              <p class="dz-note">Images are registered as canvas components. Use frame, zoom and crop to style them.</p>
            </template>

            <!-- Image controls (image selected) -->
            <template v-else>
              <div class="img-preview-wrap">
                <img :src="selectedImg.src" class="img-preview-thumb"
                  :style="{ borderRadius: selectedImg.frame === 'circle' ? '50%' : '4px' }"/>
              </div>

              <div class="tg">
                <div class="tg-label">FRAME</div>
                <div class="frame-row">
                  <button v-for="f in (['none','square','circle'] as FrameStyle[])" :key="f"
                    class="frame-btn" :class="{ 'fb-active': selectedImg.frame === f }"
                    @click="selectedImg.frame = f">
                    <svg v-if="f==='none'" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5"><line x1="2" y1="2" x2="14" y2="14"/><line x1="14" y1="2" x2="2" y2="14"/></svg>
                    <svg v-else-if="f==='square'" viewBox="0 0 16 16" fill="currentColor"><rect x="2" y="2" width="12" height="12" rx="2"/></svg>
                    <svg v-else viewBox="0 0 16 16" fill="currentColor"><circle cx="8" cy="8" r="6"/></svg>
                    {{ f === 'none' ? 'None' : f === 'square' ? 'Square' : 'Circle' }}
                  </button>
                </div>
              </div>

              <div class="tg">
                <div class="tg-label">SIZE</div>
                <div class="cr">
                  <label class="clabel">Width</label>
                  <input type="range" min="5" max="80" class="dr-range" v-model.number="selectedImg.width"/>
                  <span class="cval">{{ selectedImg.width }}%</span>
                </div>
              </div>

              <div class="tg">
                <div class="tg-label">ZOOM &amp; CROP</div>
                <div class="cr">
                  <label class="clabel">Zoom</label>
                  <input type="range" min="1" max="3" step="0.05" class="dr-range" v-model.number="selectedImg.zoom"/>
                  <span class="cval">{{ selectedImg.zoom.toFixed(1) }}×</span>
                </div>
                <div class="cr">
                  <label class="clabel">Crop X</label>
                  <input type="range" min="0" max="100" class="dr-range" v-model.number="selectedImg.cropX"/>
                  <span class="cval">{{ selectedImg.cropX }}%</span>
                </div>
                <div class="cr">
                  <label class="clabel">Crop Y</label>
                  <input type="range" min="0" max="100" class="dr-range" v-model.number="selectedImg.cropY"/>
                  <span class="cval">{{ selectedImg.cropY }}%</span>
                </div>
              </div>

              <button class="danger-btn" @click="deleteSelected">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4h6v2"/></svg>
                Remove Image
              </button>
            </template>

            <!-- All images list -->
            <div v-if="imgElements.length" class="el-section">
              <div class="tg-label" style="padding:0.875rem 1rem 0.4rem">ALL IMAGES ({{ imgElements.length }})</div>
              <div class="el-list">
                <div v-for="img in imgElements" :key="img.id" class="el-item"
                  :class="{ 'el-active': selectedId === img.id }" @click="bringToFront(img.id); activeTab='image'">
                  <img :src="img.src" class="el-thumb" :style="{ borderRadius: img.frame === 'circle' ? '50%' : '3px' }"/>
                  <div class="el-meta"><span>{{ img.frame }} frame</span><span>{{ img.width }}% wide</span></div>
                  <button class="el-del-list" @click.stop="selectedId = img.id; deleteSelected()">×</button>
                </div>
              </div>
            </div>
          </div>

          <!-- ╔══ BRUSH TAB ══╗ -->
          <div v-else-if="activeTab === 'brush'" class="tab-pane">
            <div class="brush-info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              <p>Select brush type then <strong>click on the canvas</strong> to place. Drag to reposition.</p>
            </div>

            <div class="tg">
              <div class="tg-label">TYPE</div>
              <div class="brush-types">
                <button class="btype-btn" :class="{ 'btype-active': brushType === 'flower' }" @click="brushType = 'flower'">
                  <svg viewBox="-20 -20 40 40" width="44" height="44" :style="{ color: brushColor }">
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(0,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(72,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(144,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(216,0,0)"/>
                    <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(288,0,0)"/>
                    <circle cx="0" cy="0" r="5" fill="#ffe066"/>
                  </svg>
                  <span>Flower</span>
                </button>
                <button class="btype-btn" :class="{ 'btype-active': brushType === 'rose' }" @click="brushType = 'rose'">
                  <svg viewBox="-20 -20 40 40" width="44" height="44" :style="{ color: brushColor }">
                    <path d="M0,-16 C6,-13 14,-7 12,0 C18,-4 20,5 15,10 C10,15 3,16 0,13 C-3,16 -10,15 -15,10 C-20,5 -18,-4 -12,0 C-14,-7 -6,-13 0,-16Z" fill="currentColor"/>
                    <path d="M0,-9 C4,-6 8,-1 6,3 C9,0 11,5 8,8 C5,11 1,11 0,9 C-1,11 -5,11 -8,8 C-11,5 -9,0 -6,3 C-8,-1 -4,-6 0,-9Z" fill="currentColor" opacity="0.55"/>
                    <circle cx="0" cy="2" r="2.5" fill="#ffe066" opacity="0.8"/>
                  </svg>
                  <span>Rose</span>
                </button>
              </div>
            </div>

            <div class="tg">
              <div class="tg-label">COLOR</div>
              <div class="dot-row">
                <button v-for="c in BRUSH_COLORS" :key="c" class="pdot"
                  :style="{ background: c, outline: brushColor === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#ffffff' ? 'inset 0 0 0 1px #ccc' : 'none' }"
                  @click="brushColor = c"/>
              </div>
              <div class="color-row" style="margin-top:0.35rem">
                <label class="clabel">Custom</label>
                <input type="color" v-model="brushColor" class="csi"/>
                <span class="cval">{{ brushColor }}</span>
              </div>
            </div>

            <div class="tg">
              <div class="tg-label">SIZE &amp; ANGLE</div>
              <div class="cr">
                <label class="clabel">Size</label>
                <input type="range" min="16" max="120" class="dr-range" v-model.number="brushSize"/>
                <span class="cval">{{ brushSize }}px</span>
              </div>
              <div class="cr">
                <label class="clabel">Angle</label>
                <input type="range" min="0" max="360" class="dr-range" v-model.number="brushRotation"/>
                <span class="cval">{{ brushRotation }}°</span>
              </div>
            </div>

            <!-- Live preview -->
            <div class="brush-preview">
              <svg :viewBox="'-20 -20 40 40'" :width="Math.min(brushSize, 80)" :height="Math.min(brushSize, 80)"
                :style="{ color: brushColor, transform: `rotate(${brushRotation}deg)`, display: 'block' }">
                <template v-if="brushType === 'flower'">
                  <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(0,0,0)"/>
                  <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(72,0,0)"/>
                  <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(144,0,0)"/>
                  <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(216,0,0)"/>
                  <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(288,0,0)"/>
                  <circle cx="0" cy="0" r="5" fill="#ffe066"/>
                </template>
                <template v-else>
                  <path d="M0,-16 C6,-13 14,-7 12,0 C18,-4 20,5 15,10 C10,15 3,16 0,13 C-3,16 -10,15 -15,10 C-20,5 -18,-4 -12,0 C-14,-7 -6,-13 0,-16Z" fill="currentColor"/>
                  <path d="M0,-9 C4,-6 8,-1 6,3 C9,0 11,5 8,8 C5,11 1,11 0,9 C-1,11 -5,11 -8,8 C-11,5 -9,0 -6,3 C-8,-1 -4,-6 0,-9Z" fill="currentColor" opacity="0.55"/>
                </template>
              </svg>
              <span class="bp-label">Preview</span>
            </div>

            <!-- Placed brush list -->
            <div v-if="brushElements.length" class="el-section">
              <div class="tg-label" style="padding:0.875rem 1rem 0.4rem">PLACED ({{ brushElements.length }})</div>
              <div class="el-list">
                <div v-for="br in brushElements" :key="br.id" class="el-item"
                  :class="{ 'el-active': selectedId === br.id }" @click="bringToFront(br.id)">
                  <div class="el-brush-icon" :style="{ color: br.color }">
                    <svg viewBox="-20 -20 40 40" width="28" height="28">
                      <template v-if="br.brushType === 'flower'">
                        <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(0,0,0)"/>
                        <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(72,0,0)"/>
                        <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(144,0,0)"/>
                        <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(216,0,0)"/>
                        <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(288,0,0)"/>
                        <circle cx="0" cy="0" r="5" fill="#ffe066"/>
                      </template>
                      <template v-else>
                        <path d="M0,-16 C6,-13 14,-7 12,0 C18,-4 20,5 15,10 C10,15 3,16 0,13 C-3,16 -10,15 -15,10 C-20,5 -18,-4 -12,0 C-14,-7 -6,-13 0,-16Z" fill="currentColor"/>
                      </template>
                    </svg>
                  </div>
                  <span style="flex:1;font-size:0.68rem">{{ br.brushType }} · {{ br.size }}px</span>
                  <button class="el-del-list" @click.stop="selectedId = br.id; deleteSelected()">×</button>
                </div>
              </div>
            </div>
          </div>

          <!-- ╔══ BORDER TAB ══╗ -->
          <div v-else-if="activeTab === 'border'" class="tab-pane">
            <div class="tg">
              <div class="tg-label">STYLE</div>
              <div class="border-grid">
                <button v-for="s in BORDER_STYLES" :key="s.id" class="bsb"
                  :class="{ 'bsb-active': border.style === s.id }" @click="border.style = s.id">
                  <span class="bsb-line" :style="{ borderBottomWidth:'3px', borderBottomStyle: s.id==='ornate'||s.id==='none'?'solid':s.id, borderBottomColor: s.id==='none'?'transparent':'#888' }"/>
                  {{ s.label }}
                </button>
              </div>
            </div>
            <div class="tg">
              <div class="tg-label">COLOR</div>
              <div class="dot-row">
                <button v-for="c in BORDER_COLORS" :key="c" class="pdot"
                  :style="{ background: c, outline: border.color === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#f1faee' ? 'inset 0 0 0 1px #ccc' : 'none' }"
                  @click="border.color = c"/>
              </div>
              <div class="color-row" style="margin-top:0.35rem">
                <label class="clabel">Custom</label>
                <input type="color" v-model="border.color" class="csi"/>
                <span class="cval">{{ border.color }}</span>
              </div>
            </div>
            <div class="tg">
              <div class="tg-label">WIDTH</div>
              <div class="cr">
                <label class="clabel">Size</label>
                <input type="range" min="0" max="32" class="dr-range" v-model.number="border.width"/>
                <span class="cval">{{ border.width }}px</span>
              </div>
            </div>
            <!-- Preview -->
            <div class="tg">
              <div class="tg-label">PREVIEW</div>
              <div class="border-preview"
                :style="{
                  border: border.style !== 'none' && border.width > 0 ? `${border.width}px ${border.style === 'ornate' ? 'solid' : border.style} ${border.color}` : '1px dashed #ccc',
                  outline: border.style === 'ornate' && border.width > 0 ? `${Math.max(1,Math.round(border.width*0.35))}px solid ${border.color}` : 'none',
                  outlineOffset: '3px',
                }">Board Border
              </div>
            </div>
          </div>

          <!-- ╔══ CORNER TAB ══╗ -->
          <div v-else-if="activeTab === 'corner'" class="tab-pane">
            <div class="sec-toggle">
              <button class="stg-btn" :class="{ active: activeSection === 'upper' }" @click="activeSection = 'upper'">Upper</button>
              <button class="stg-btn" :class="{ active: activeSection === 'lower' }" @click="activeSection = 'lower'">Lower</button>
            </div>
            <div class="tg">
              <div class="tg-label">CORNER STYLE — {{ activeSection === 'upper' ? 'UPPER' : 'LOWER' }} SECTION</div>
              <div class="corner-grid">
                <button v-for="c in CORNERS" :key="c.id" class="corner-btn"
                  :class="{ 'cb-active': sec.cornerStyle === c.id }"
                  @click="sec.cornerStyle = c.id">
                  <div class="cb-preview"
                    :style="{
                      borderRadius: c.id==='rounded'?'8px':c.id==='ornate'?'2px 12px':'2px',
                      clipPath: c.id==='cut'?'polygon(8px 0%,calc(100% - 8px) 0%,100% 8px,100% 100%,0% 100%,0% 8px)':'none',
                    }"/>
                  <span>{{ c.label }}</span>
                  <span v-if="c.id === 'floral'" style="font-size:0.75rem">🌸</span>
                </button>
              </div>
            </div>
          </div>

        </div><!-- /tab-body -->

        <!-- ── PANEL FOOTER ────────────────────────────────────────── -->
        <div class="panel-footer">
          <div class="size-sect">
            <div class="tg-label">BOARD SIZE</div>
            <div class="size-opts">
              <button v-for="s in SIZES" :key="s.id" class="size-opt"
                :class="{ 'so-active': physicalSize === s.id }" @click="physicalSize = s.id">
                <span class="so-label">{{ s.label }}</span>
                <span class="so-sub">{{ s.desc }}</span>
                <span class="so-price">{{ formatRupiah(s.price) }}</span>
                <span v-if="s.recommended" class="so-badge">★</span>
              </button>
            </div>
          </div>
          <button id="btn-finalize-footer" class="finalize-btn" @click="showReview = true">
            <span class="fb-left">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 01-8 0"/>
              </svg>
              Finalize &amp; Order
            </span>
            <span class="fb-price">{{ formatRupiah(totalPrice) }}</span>
          </button>
        </div>

      </aside><!-- /panel -->

    </div><!-- /dr-body -->

    <!-- ═══ REVIEW MODAL ══════════════════════════════════════════════ -->
    <Transition name="modal-fade">
      <div v-if="showReview" class="modal-backdrop" @click.self="showReview = false" role="dialog" aria-modal="true">
        <div class="review-modal">
          <button class="rm-close" @click="showReview = false" aria-label="Close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
          <div class="rm-header">
            <p class="rm-eyebrow">REVIEW ORDER</p>
            <h2 class="rm-title">Your Custom Board</h2>
          </div>
          <div class="rm-body">
            <!-- Mini board preview -->
            <div class="rm-board-wrap">
              <div class="rm-board" :style="boardBorderStyle">
                <div :style="{ height: (heightRatio * 100)+'%', backgroundColor: upper.bgColor, display:'flex', alignItems:'center', justifyContent:'center', padding:'8px', overflow:'hidden' }">
                  <span :style="{ fontFamily: getFont(upper.headerFont), color: upper.headerColor, fontSize:'13px', fontWeight:700, textAlign:'center' }">{{ upper.headerText || '(no header)' }}</span>
                </div>
                <div :style="{ height: ((1-heightRatio)*100)+'%', backgroundColor: lower.bgColor, display:'flex', alignItems:'center', justifyContent:'center', padding:'8px', overflow:'hidden' }">
                  <span :style="{ fontFamily: getFont(lower.bodyFont), color: lower.bodyColor, fontSize:'11px', textAlign:'center' }">{{ lower.bodyText.replace(/\n/g,' / ') || '(no text)' }}</span>
                </div>
              </div>
            </div>
            <!-- Specs -->
            <div class="rm-specs">
              <div class="spec-row"><span class="sk">Upper Header</span><span class="sv">{{ upper.headerText || '—' }}</span></div>
              <div class="spec-row"><span class="sk">Upper Body</span><span class="sv">{{ upper.bodyText.replace(/\n/g,' / ') || '—' }}</span></div>
              <div class="spec-row"><span class="sk">Lower Body</span><span class="sv">{{ lower.bodyText.replace(/\n/g,' / ') || '—' }}</span></div>
              <div class="spec-row"><span class="sk">Elements</span><span class="sv">{{ imgElements.length }} image{{ imgElements.length!==1?'s':'' }}, {{ brushElements.length }} brush stroke{{ brushElements.length!==1?'s':'' }}</span></div>
              <div class="spec-row"><span class="sk">Size</span><span class="sv">{{ SIZES.find(s=>s.id===physicalSize)?.label }} — {{ SIZES.find(s=>s.id===physicalSize)?.desc }}</span></div>
              <div class="spec-divider"/>
              <div class="spec-row spec-total"><span class="sk">Total</span><span class="sv spec-price">{{ formatRupiah(totalPrice) }}</span></div>
              <p class="spec-note">* Our team will review your design before production and contact you if adjustments are needed.</p>
            </div>
          </div>
          <div class="rm-footer">
            <button class="rm-back" @click="showReview = false">Back to Designer</button>
            <button id="btn-add-to-cart" class="rm-confirm" :class="{ loading: isAdding }" :disabled="isAdding" @click="addToCartHandler">
              <svg v-if="!isAdding" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <path d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M16.5 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0Zm3 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0Z"/>
              </svg>
              <svg v-else class="spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/></svg>
              {{ isAdding ? 'Adding…' : `Add to Cart — ${formatRupiah(totalPrice)}` }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- ═══ SUCCESS TOAST ══════════════════════════════════════════════ -->
    <Transition name="toast-slide">
      <div v-if="showToast" class="success-toast" role="status">
        <div class="toast-check">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><path d="M5 13l4 4L19 7"/></svg>
        </div>
        <div><p class="toast-title">Added to cart!</p><p class="toast-sub">Redirecting you home…</p></div>
      </div>
    </Transition>

  </div>
</template>

<style>
/* ─── GOOGLE FONTS ─────────────────────────────────────────────────── */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;900&family=Playfair+Display:ital,wght@0,700;0,800;1,600&family=Dancing+Script:wght@600;700&family=Bebas+Neue&family=Merriweather:wght@700;900&family=Pacifico&display=swap');

/* ─── RESET ──────────────────────────────────────────────────────── */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
button { font-family: inherit; }

/* ─── ROOT ───────────────────────────────────────────────────────── */
.dr-root {
  position: fixed; inset: 0;
  display: flex; flex-direction: column;
  font-family: 'Inter', system-ui, sans-serif;
  background: #f4f0eb; color: #1c1813;
  overflow: hidden; color-scheme: light;
}

/* ─── NAVBAR ─────────────────────────────────────────────────────── */
.dr-nav {
  height: 52px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 1.25rem;
  background: #ffffff; border-bottom: 1px solid #e5ddd4;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06); z-index: 30;
}
.dr-back {
  display: flex; align-items: center; gap: 0.35rem;
  color: #6b5d52; text-decoration: none; font-size: 0.78rem; font-weight: 600;
  padding: 0.3rem 0.7rem; border-radius: 6px; transition: background 0.15s, color 0.15s;
}
.dr-back:hover { background: #f4f0eb; color: #1c1813; }
.dr-back svg { width: 0.9rem; height: 0.9rem; }
.dr-nav-center {
  display: flex; align-items: center; gap: 0.5rem;
  position: absolute; left: 50%; transform: translateX(-50%);
}
.dr-brand { font-size: 0.62rem; font-weight: 900; letter-spacing: 0.25em; color: #a8998d; }
.dr-dot { color: #d4c9bc; font-size: 0.65rem; }
.dr-page-title { font-size: 0.84rem; font-weight: 700; color: #1c1813; }
.dr-nav-right { display: flex; align-items: center; gap: 0.7rem; }
.dr-scale-chip { font-size: 0.63rem; font-weight: 700; color: #a8998d; background: #f4f0eb; border: 1px solid #e5ddd4; padding: 0.18rem 0.5rem; border-radius: 4px; letter-spacing: 0.05em; }
.dr-finalize-btn {
  display: flex; align-items: center; gap: 0.45rem;
  background: #c4703e; color: #fff; font-size: 0.77rem; font-weight: 700;
  padding: 0.45rem 0.95rem; border: none; border-radius: 6px; cursor: pointer;
  box-shadow: 0 2px 6px rgba(196,112,62,0.35); transition: background 0.15s, transform 0.1s, box-shadow 0.15s;
}
.dr-finalize-btn:hover { background: #b5622f; transform: translateY(-1px); box-shadow: 0 4px 12px rgba(196,112,62,0.4); }
.dr-finalize-btn svg { width: 0.8rem; height: 0.8rem; }

/* ─── BODY LAYOUT ────────────────────────────────────────────────── */
.dr-body { flex: 1; display: flex; overflow: hidden; }

/* ─── CANVAS AREA ────────────────────────────────────────────────── */
.dr-canvas-area { flex: 1; display: flex; flex-direction: column; min-width: 0; overflow: hidden; }

.chess-bg {
  flex: 1; display: flex; align-items: center; justify-content: center;
  overflow: hidden; position: relative;
  background-image: repeating-conic-gradient(#ede8e0 0% 25%, #f8f5f0 0% 50%);
  background-size: 20px 20px;
}

.board-scaler { position: relative; flex-shrink: 0; }

/* ─── BOARD ──────────────────────────────────────────────────────── */
.board-frame {
  position: absolute; top: 0; left: 0;
  box-shadow: 0 8px 40px rgba(0,0,0,0.2), 0 2px 8px rgba(0,0,0,0.08);
  overflow: hidden;
  transition: border 0.2s, outline 0.2s;
}

/* ─── SECTIONS ───────────────────────────────────────────────────── */
.board-section { position: absolute; left: 0; right: 0; overflow: visible; }
.sec-inner {
  width: 100%; height: 100%; position: relative; z-index: 1;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  padding: 1.25rem 2rem; gap: 0.5rem;
}
.sec-text { line-height: 1.3; word-break: break-word; transition: all 0.25s; }
.sec-header { font-weight: 800; }
.sec-body   { white-space: pre-line; }

/* ─── FLORAL CORNERS ─────────────────────────────────────────────── */
.floral-corner { position: absolute; z-index: 5; pointer-events: none; font-size: 1.5rem; line-height: 1; }
.fc-tl { top: 5px; left: 8px; }
.fc-tr { top: 5px; right: 8px; }
.fc-bl { bottom: 5px; left: 8px; }
.fc-br { bottom: 5px; right: 8px; }

/* ─── SECTION DIVIDER ────────────────────────────────────────────── */
.section-divider {
  position: absolute; left: 0; right: 0; height: 12px; margin-top: -6px;
  display: flex; align-items: center; justify-content: center;
  cursor: row-resize; z-index: 50;
}
.div-track {
  position: absolute; inset: 5px 0;
  background: var(--dc, #f5c842); opacity: 0.6; transition: opacity 0.2s;
}
.div-knob {
  position: relative; z-index: 2; background: #fff; border: 1px solid #e0d8ce;
  border-radius: 4px; padding: 0.15rem 0.45rem;
  opacity: 0; transition: opacity 0.2s, transform 0.2s;
  color: #a8998d; box-shadow: 0 1px 4px rgba(0,0,0,0.1);
}
.section-divider:hover .div-track { opacity: 1; }
.section-divider:hover .div-knob  { opacity: 1; transform: scale(1.05); }

/* ─── CANVAS ELEMENTS ────────────────────────────────────────────── */
.canvas-el { position: absolute; user-select: none; }
.el-selected { outline: 2px dashed rgba(196,112,62,0.85); outline-offset: 2px; }
.el-del {
  position: absolute; top: -10px; right: -10px;
  width: 22px; height: 22px; background: #e63946; color: #fff; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 0.9rem; font-weight: 900; cursor: pointer; z-index: 99;
  box-shadow: 0 2px 6px rgba(230,57,70,0.4); transition: transform 0.15s, background 0.15s;
  border: none; line-height: 1;
}
.el-del:hover { transform: scale(1.15); background: #c1121f; }

/* ─── BRUSH OVERLAY ──────────────────────────────────────────────── */
.brush-overlay {
  position: absolute; inset: 0; pointer-events: none; z-index: 200;
  display: flex; align-items: flex-end; justify-content: center; padding-bottom: 10px;
}
.brush-overlay span {
  background: rgba(0,0,0,0.52); color: #fff; font-size: 0.72rem; font-weight: 600;
  padding: 0.3rem 0.85rem; border-radius: 20px; backdrop-filter: blur(4px); letter-spacing: 0.03em;
}

/* ─── CANVAS INFO BAR ────────────────────────────────────────────── */
.canvas-info {
  flex-shrink: 0; height: 32px; display: flex; align-items: center; gap: 0.45rem;
  padding: 0 1rem; background: #fff; border-top: 1px solid #e5ddd4; font-size: 0.7rem; color: #a8998d;
}
.ci-sel { color: #c4703e; font-weight: 600; }
.ci-desel { font-size: 0.65rem; font-weight: 700; color: #6b5d52; background: #f4f0eb; border: 1px solid #e5ddd4; padding: 0.12rem 0.45rem; border-radius: 4px; cursor: pointer; }
.ci-hint { margin-left: auto; font-size: 0.65rem; font-style: italic; color: #c4703e; }

/* ─── TOOL PANEL ─────────────────────────────────────────────────── */
.dr-panel { width: 340px; flex-shrink: 0; display: flex; flex-direction: column; background: #fff; border-left: 1px solid #e5ddd4; overflow: hidden; }

/* ─── TAB BAR ────────────────────────────────────────────────────── */
.tab-bar { display: flex; flex-shrink: 0; border-bottom: 1px solid #e5ddd4; background: #f9f6f2; }
.tab-btn {
  flex: 1; display: flex; flex-direction: column; align-items: center; gap: 0.22rem;
  padding: 0.55rem 0.2rem; font-size: 0.57rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase;
  color: #a8998d; background: none; border: none; cursor: pointer;
  border-bottom: 2px solid transparent; transition: color 0.15s, border-color 0.15s, background 0.15s;
}
.tab-btn svg { width: 1.1rem; height: 1.1rem; }
.tab-btn:hover { color: #6b5d52; background: rgba(196,112,62,0.04); }
.tab-btn.tab-active { color: #c4703e; border-bottom-color: #c4703e; background: #fff; }

/* ─── TAB BODY ───────────────────────────────────────────────────── */
.tab-body { flex: 1; overflow-y: auto; scrollbar-width: thin; scrollbar-color: #d4c9bc transparent; }
.tab-body::-webkit-scrollbar { width: 4px; }
.tab-body::-webkit-scrollbar-thumb { background: #d4c9bc; border-radius: 2px; }
.tab-pane { display: flex; flex-direction: column; }

/* ─── TOOL GROUP ─────────────────────────────────────────────────── */
.tg { padding: 0.875rem 1rem; border-bottom: 1px solid #f0ebe4; display: flex; flex-direction: column; gap: 0.55rem; }
.tg-label { font-size: 0.56rem; font-weight: 900; letter-spacing: 0.18em; color: #c4703e; text-transform: uppercase; }

/* ─── SECTION TOGGLE ─────────────────────────────────────────────── */
.sec-toggle { display: flex; border-bottom: 1px solid #e5ddd4; background: #f9f6f2; flex-shrink: 0; }
.stg-btn { flex: 1; padding: 0.6rem; font-size: 0.72rem; font-weight: 700; color: #a8998d; background: none; border: none; cursor: pointer; border-bottom: 2px solid transparent; transition: color 0.15s, border-color 0.15s; }
.stg-btn.active { color: #c4703e; border-bottom-color: #c4703e; }

/* ─── INPUTS ─────────────────────────────────────────────────────── */
.dr-ta {
  width: 100%; resize: vertical; background: #f9f6f2; border: 1px solid #e5ddd4;
  border-radius: 6px; color: #1c1813; font-size: 0.8rem; font-family: inherit;
  padding: 0.5rem 0.65rem; outline: none; line-height: 1.5; transition: border-color 0.15s, background 0.15s;
  user-select: text;
}
.dr-ta:focus { border-color: #c4703e; background: #fff; }
.dr-ta::placeholder { color: #c0b0a4; }

.cr { display: flex; align-items: center; gap: 0.45rem; }
.clabel { font-size: 0.6rem; font-weight: 700; color: #a8998d; min-width: 38px; flex-shrink: 0; }
.cval { font-size: 0.63rem; font-weight: 700; color: #6b5d52; min-width: 34px; text-align: right; font-variant-numeric: tabular-nums; }

.dr-range { flex: 1; height: 4px; border-radius: 2px; appearance: none; background: #e5ddd4; outline: none; cursor: pointer; }
.dr-range::-webkit-slider-thumb { appearance: none; width: 14px; height: 14px; background: #c4703e; border-radius: 50%; border: 2px solid #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.2); transition: transform 0.1s; }
.dr-range::-webkit-slider-thumb:hover { transform: scale(1.2); }
.dr-range::-moz-range-thumb { width: 14px; height: 14px; background: #c4703e; border-radius: 50%; border: 2px solid #fff; }

/* ─── FONT GRID ──────────────────────────────────────────────────── */
.font-grid { display: flex; flex-wrap: wrap; gap: 0.3rem; }
.font-chip { padding: 0.28rem 0.55rem; border-radius: 4px; font-size: 0.7rem; font-weight: 600; background: #f4f0eb; border: 1px solid #e5ddd4; color: #6b5d52; cursor: pointer; transition: all 0.15s; }
.font-chip:hover { border-color: #c4703e; color: #c4703e; }
.font-chip.fc-active { background: #c4703e; color: #fff; border-color: #c4703e; }

/* ─── ALIGNMENT ──────────────────────────────────────────────────── */
.align-row { display: flex; gap: 0.3rem; }
.aln-btn { display: flex; align-items: center; padding: 0.28rem 0.45rem; border-radius: 4px; background: #f4f0eb; border: 1px solid #e5ddd4; color: #6b5d52; cursor: pointer; transition: all 0.15s; }
.aln-btn svg { width: 0.85rem; height: 0.65rem; }
.aln-btn:hover { border-color: #c4703e; color: #c4703e; }
.aln-btn.aln-active { background: #c4703e; color: #fff; border-color: #c4703e; }

/* ─── COLORS ─────────────────────────────────────────────────────── */
.color-row { display: flex; align-items: center; gap: 0.45rem; }
.csi { width: 30px; height: 26px; border: 2px solid #e5ddd4; border-radius: 4px; cursor: pointer; padding: 1px; background: none; transition: border-color 0.15s; }
.csi:hover { border-color: #c4703e; }
.dot-row { display: flex; flex-wrap: wrap; gap: 0.35rem; }
.pdot { width: 22px; height: 22px; border-radius: 50%; border: none; cursor: pointer; padding: 0; transition: transform 0.15s; }
.pdot:hover { transform: scale(1.2); }

/* ─── DROP ZONE ──────────────────────────────────────────────────── */
.drop-zone { margin: 1rem; border: 2px dashed #d4c9bc; border-radius: 8px; padding: 1.75rem 1rem; display: flex; flex-direction: column; align-items: center; gap: 0.6rem; text-align: center; transition: border-color 0.2s, background 0.2s; cursor: pointer; }
.drop-zone:hover { border-color: #c4703e; background: rgba(196,112,62,0.04); }
.dz-icon { width: 2.25rem; height: 2.25rem; color: #a8998d; }
.dz-title { font-size: 0.82rem; font-weight: 700; color: #1c1813; }
.dz-sub { font-size: 0.7rem; color: #a8998d; line-height: 1.5; }
.dz-browse { font-size: 0.72rem; font-weight: 700; color: #c4703e; border: 1px solid #c4703e; border-radius: 6px; padding: 0.38rem 0.9rem; cursor: pointer; transition: background 0.15s, color 0.15s; }
.dz-browse:hover { background: #c4703e; color: #fff; }
.dz-note { margin: 0 1rem 1rem; font-size: 0.65rem; color: #a8998d; line-height: 1.55; padding: 0.5rem 0.7rem; background: #f9f6f2; border-radius: 6px; border-left: 2px solid #c4703e; }

/* ─── IMAGE CONTROLS ─────────────────────────────────────────────── */
.img-preview-wrap { padding: 1rem; display: flex; justify-content: center; border-bottom: 1px solid #f0ebe4; }
.img-preview-thumb { width: 96px; height: 96px; object-fit: cover; border-radius: 6px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.frame-row { display: flex; gap: 0.45rem; }
.frame-btn { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 0.3rem; padding: 0.55rem; border: 1px solid #e5ddd4; border-radius: 6px; background: #f4f0eb; color: #6b5d52; cursor: pointer; font-size: 0.62rem; font-weight: 700; transition: all 0.15s; }
.frame-btn svg { width: 0.9rem; height: 0.9rem; }
.frame-btn:hover { border-color: #c4703e; color: #c4703e; }
.frame-btn.fb-active { background: #c4703e; color: #fff; border-color: #c4703e; }
.danger-btn { display: flex; align-items: center; justify-content: center; gap: 0.4rem; margin: 0.875rem 1rem; padding: 0.6rem; border-radius: 6px; font-size: 0.72rem; font-weight: 700; color: #e63946; background: #fff5f5; border: 1px solid rgba(230,57,70,0.25); cursor: pointer; transition: all 0.15s; }
.danger-btn svg { width: 0.85rem; height: 0.85rem; }
.danger-btn:hover { background: #e63946; color: #fff; border-color: #e63946; }

/* ─── ELEMENT LIST ───────────────────────────────────────────────── */
.el-section { border-top: 1px solid #f0ebe4; }
.el-list { display: flex; flex-direction: column; gap: 0.3rem; padding: 0 1rem 0.875rem; }
.el-item { display: flex; align-items: center; gap: 0.5rem; padding: 0.45rem 0.6rem; border-radius: 6px; border: 1px solid #e5ddd4; background: #f9f6f2; cursor: pointer; transition: all 0.15s; font-size: 0.68rem; color: #6b5d52; }
.el-item:hover { border-color: #c4703e; background: rgba(196,112,62,0.04); }
.el-item.el-active { border-color: #c4703e; background: rgba(196,112,62,0.08); color: #c4703e; }
.el-thumb { width: 34px; height: 34px; object-fit: cover; flex-shrink: 0; border-radius: 3px; }
.el-meta { display: flex; flex-direction: column; gap: 0.1rem; flex: 1; }
.el-del-list { width: 20px; height: 20px; border-radius: 50%; background: #f4f0eb; border: 1px solid #e5ddd4; color: #6b5d52; font-size: 0.85rem; cursor: pointer; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all 0.15s; line-height: 1; }
.el-del-list:hover { background: #e63946; color: #fff; border-color: #e63946; }
.el-brush-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }

/* ─── BRUSH PANEL ────────────────────────────────────────────────── */
.brush-info { display: flex; align-items: flex-start; gap: 0.6rem; padding: 0.875rem 1rem; background: rgba(196,112,62,0.06); border-bottom: 1px solid #f0ebe4; }
.brush-info svg { width: 1rem; height: 1rem; flex-shrink: 0; color: #c4703e; margin-top: 1px; }
.brush-info p { font-size: 0.7rem; color: #6b5d52; line-height: 1.5; }
.brush-info strong { color: #c4703e; }
.brush-types { display: flex; gap: 0.5rem; }
.btype-btn { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 0.35rem; padding: 0.75rem 0.5rem; border: 1px solid #e5ddd4; border-radius: 6px; background: #f4f0eb; color: #6b5d52; cursor: pointer; font-size: 0.7rem; font-weight: 700; transition: all 0.15s; }
.btype-btn:hover { border-color: #c4703e; }
.btype-btn.btype-active { background: #fff3ed; border-color: #c4703e; color: #c4703e; }
.brush-preview { display: flex; flex-direction: column; align-items: center; gap: 0.4rem; padding: 1rem; background: #f9f6f2; border-top: 1px solid #f0ebe4; border-bottom: 1px solid #f0ebe4; }
.bp-label { font-size: 0.6rem; font-weight: 700; color: #a8998d; letter-spacing: 0.08em; }

/* ─── BORDER PANEL ───────────────────────────────────────────────── */
.border-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.35rem; }
.bsb { display: flex; flex-direction: column; align-items: center; gap: 0.38rem; padding: 0.5rem 0.3rem; border: 1px solid #e5ddd4; border-radius: 6px; background: #f4f0eb; color: #6b5d52; cursor: pointer; font-size: 0.58rem; font-weight: 700; transition: all 0.15s; }
.bsb-line { width: 100%; display: block; }
.bsb:hover { border-color: #c4703e; }
.bsb.bsb-active { background: #fff3ed; border-color: #c4703e; color: #c4703e; }
.border-preview { height: 52px; background: #f4f0eb; border-radius: 4px; display: flex; align-items: center; justify-content: center; font-size: 0.68rem; font-weight: 700; color: #a8998d; transition: border 0.2s; }

/* ─── CORNER PANEL ───────────────────────────────────────────────── */
.corner-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 0.35rem; }
.corner-btn { display: flex; flex-direction: column; align-items: center; gap: 0.35rem; padding: 0.55rem 0.3rem; border: 1px solid #e5ddd4; border-radius: 6px; background: #f4f0eb; color: #6b5d52; cursor: pointer; font-size: 0.56rem; font-weight: 700; transition: all 0.15s; }
.cb-preview { width: 24px; height: 18px; background: #c0b0a0; }
.corner-btn:hover { border-color: #c4703e; }
.corner-btn.cb-active { background: #fff3ed; border-color: #c4703e; color: #c4703e; }
.corner-btn.cb-active .cb-preview { background: #c4703e; }

/* ─── PANEL FOOTER ───────────────────────────────────────────────── */
.panel-footer { flex-shrink: 0; border-top: 1px solid #e5ddd4; background: #fff; padding: 0.875rem 1rem; display: flex; flex-direction: column; gap: 0.7rem; }
.size-sect { display: flex; flex-direction: column; gap: 0.4rem; }
.size-opts { display: flex; gap: 0.4rem; }
.size-opt { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 0.1rem; padding: 0.5rem 0.25rem; border: 1px solid #e5ddd4; border-radius: 6px; background: #f4f0eb; cursor: pointer; position: relative; transition: all 0.15s; }
.so-label { font-size: 0.65rem; font-weight: 800; color: #1c1813; text-align: center; }
.so-sub { font-size: 0.54rem; color: #a8998d; }
.so-price { font-size: 0.6rem; font-weight: 700; color: #c4703e; }
.so-badge { position: absolute; top: -6px; right: -6px; font-size: 0.55rem; background: #f59e0b; color: #fff; border-radius: 50%; width: 14px; height: 14px; display: flex; align-items: center; justify-content: center; line-height: 1; }
.size-opt:hover { border-color: #c4703e; }
.size-opt.so-active { background: #fff3ed; border-color: #c4703e; border-width: 2px; }
.finalize-btn { display: flex; align-items: center; justify-content: space-between; padding: 0.75rem 1rem; background: linear-gradient(135deg, #c4703e, #a85a2e); color: #fff; border: none; border-radius: 8px; cursor: pointer; box-shadow: 0 3px 12px rgba(196,112,62,0.35); transition: all 0.15s; }
.finalize-btn:hover { background: linear-gradient(135deg, #b5622f, #8f4a22); transform: translateY(-1px); box-shadow: 0 5px 16px rgba(196,112,62,0.45); }
.fb-left { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8rem; font-weight: 700; }
.fb-left svg { width: 0.95rem; height: 0.95rem; }
.fb-price { font-size: 0.85rem; font-weight: 900; font-variant-numeric: tabular-nums; }

/* ─── REVIEW MODAL ───────────────────────────────────────────────── */
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.35); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 200; padding: 1rem; }
.review-modal { background: #fff; border-radius: 12px; width: 100%; max-width: 640px; max-height: 90vh; overflow-y: auto; position: relative; box-shadow: 0 24px 64px rgba(0,0,0,0.2); display: flex; flex-direction: column; }
.rm-close { position: absolute; top: 1rem; right: 1rem; width: 32px; height: 32px; border-radius: 50%; background: #f4f0eb; border: 1px solid #e5ddd4; color: #6b5d52; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all 0.15s; }
.rm-close svg { width: 0.9rem; height: 0.9rem; }
.rm-close:hover { background: #e5ddd4; color: #1c1813; }
.rm-header { padding: 1.5rem 1.5rem 1rem; border-bottom: 1px solid #f0ebe4; }
.rm-eyebrow { font-size: 0.6rem; font-weight: 900; letter-spacing: 0.2em; color: #c4703e; margin-bottom: 0.25rem; }
.rm-title { font-size: 1.4rem; font-weight: 800; color: #1c1813; }
.rm-body { display: flex; gap: 1.5rem; padding: 1.5rem; align-items: flex-start; flex-wrap: wrap; }
.rm-board-wrap { flex: 0 0 auto; width: 180px; }
.rm-board { width: 180px; height: 120px; border-radius: 4px; overflow: hidden; box-shadow: 0 4px 16px rgba(0,0,0,0.12); }
.rm-specs { flex: 1; min-width: 200px; display: flex; flex-direction: column; gap: 0; }
.spec-row { display: flex; justify-content: space-between; align-items: flex-start; padding: 0.45rem 0; border-bottom: 1px solid #f0ebe4; gap: 0.75rem; }
.sk { font-size: 0.68rem; font-weight: 600; color: #a8998d; flex-shrink: 0; }
.sv { font-size: 0.72rem; font-weight: 600; color: #1c1813; text-align: right; word-break: break-word; }
.spec-divider { border-top: 2px solid #e5ddd4; margin: 0.5rem 0; }
.spec-total { margin-top: 0.25rem; }
.spec-price { font-size: 1.1rem; font-weight: 900; color: #c4703e; }
.spec-note { font-size: 0.65rem; color: #a8998d; line-height: 1.5; margin-top: 0.75rem; padding: 0.5rem 0.7rem; background: #f9f6f2; border-radius: 6px; border-left: 2px solid #e5ddd4; }
.rm-footer { display: flex; gap: 0.75rem; padding: 1rem 1.5rem 1.5rem; border-top: 1px solid #f0ebe4; }
.rm-back { flex: 1; padding: 0.75rem; border-radius: 6px; background: #f4f0eb; border: 1px solid #e5ddd4; color: #6b5d52; font-size: 0.78rem; font-weight: 700; cursor: pointer; transition: all 0.15s; }
.rm-back:hover { background: #ede8e0; }
.rm-confirm { flex: 2; display: flex; align-items: center; justify-content: center; gap: 0.45rem; padding: 0.75rem; border-radius: 6px; background: linear-gradient(135deg, #c4703e, #a85a2e); color: #fff; border: none; font-size: 0.8rem; font-weight: 700; cursor: pointer; transition: all 0.15s; box-shadow: 0 2px 8px rgba(196,112,62,0.3); }
.rm-confirm:hover:not(:disabled) { background: linear-gradient(135deg, #b5622f, #8f4a22); }
.rm-confirm:disabled { opacity: 0.7; cursor: not-allowed; }
.rm-confirm svg { width: 1rem; height: 1rem; }
.rm-confirm.loading { opacity: 0.8; }
.spin { animation: spinAnim 0.9s linear infinite; }
@keyframes spinAnim { to { transform: rotate(360deg); } }

/* ─── SUCCESS TOAST ──────────────────────────────────────────────── */
.success-toast { position: fixed; bottom: 1.5rem; left: 50%; transform: translateX(-50%); display: flex; align-items: center; gap: 0.75rem; background: #fff; border: 1px solid #e5ddd4; border-radius: 10px; padding: 0.875rem 1.25rem; box-shadow: 0 8px 24px rgba(0,0,0,0.12); z-index: 300; min-width: 240px; }
.toast-check { width: 32px; height: 32px; border-radius: 50%; background: rgba(107,143,110,0.12); color: #6b8f6e; display: flex; align-items: center; justify-content: center; flex-shrink: 0; border: 1px solid rgba(107,143,110,0.3); }
.toast-check svg { width: 1rem; height: 1rem; }
.toast-title { font-size: 0.82rem; font-weight: 800; color: #1c1813; margin-bottom: 0.1rem; }
.toast-sub { font-size: 0.68rem; color: #a8998d; }

/* ─── TRANSITIONS ────────────────────────────────────────────────── */
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.25s, transform 0.25s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-fade-enter-from .review-modal, .modal-fade-leave-to .review-modal { transform: scale(0.95) translateY(10px); }

.toast-slide-enter-active, .toast-slide-leave-active { transition: opacity 0.3s, transform 0.3s; }
.toast-slide-enter-from, .toast-slide-leave-to { opacity: 0; transform: translateX(-50%) translateY(16px); }

/* ─── RESPONSIVE ─────────────────────────────────────────────────── */
@media (max-width: 768px) {
  .dr-body { flex-direction: column; }
  .dr-canvas-area { flex: none; height: 45vh; }
  .dr-panel { width: 100%; height: 55vh; flex-direction: column; border-left: none; border-top: 1px solid #e5ddd4; }
  .dr-nav-center { display: none; }
  .dr-scale-chip { display: none; }
}
</style>