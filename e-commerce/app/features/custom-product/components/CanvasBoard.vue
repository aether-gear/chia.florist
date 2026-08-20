<script setup lang="ts">
import { computed, onMounted, nextTick } from 'vue'
import type { CanvasImage, BrushStroke } from '../types'
import { useCart } from '~/composables/useCart'

const props = defineProps<{
  design: ReturnType<typeof import('../useCustomDesign').useCustomDesign>
}>()

const { formatRupiah } = useCart()

// 3D transform style for board-scaler
const scalerTransform = computed(() => {
  if (props.design.is3DMode) {
    // No scale() here — board-scaler width/height already encodes boardScale.
    // Pivot is center center (set via CSS .board-scaler--3d), so rotation
    // happens around the visual middle of the board.
    return `perspective(1200px) rotateX(${props.design.rotateX}deg) rotateY(${props.design.rotateY}deg)`
  }
  return undefined
})

// Cursor for chess-bg in 3D mode
const chessCursor = computed(() => props.design.is3DMode ? 'grab' : undefined)

onMounted(() => {
  nextTick(() => props.design.updateScale())
  setTimeout(() => props.design.updateScale(), 50)
  setTimeout(() => props.design.updateScale(), 200)
})
</script>

<template>
  <div class="dr-canvas-area" :ref="(el) => { (design as any).containerRef = el; design.updateScale(); }">
    <div
      class="chess-bg"
      :class="{ 'chess-bg--3d': design.is3DMode }"
      :style="chessCursor ? { cursor: chessCursor } : {}"
      @dragover.prevent
      @mousedown="design.is3DMode ? design.start3DDrag($event) : undefined"
      @touchstart="design.is3DMode ? design.start3DDrag($event) : undefined"
    >
      <!-- 📦 FLOATING ELEMENT COUNT CHIP (Top-Left) -->
      <div class="canvas-info-chip">
        <span v-if="!design.is3DMode">{{ design.elements.length }} element{{ design.elements.length !== 1 ? 's' : '' }}</span>
        <span v-if="design.is3DMode" class="ci-hint" style="color: #7c3aed; font-weight: 700;">⬡ 3D Preview</span>
      </div>

      <!-- ⬡ 3D TOGGLE BUTTON -->
      <button
        class="btn-3d-toggle"
        :class="design.is3DMode ? 'mode-3d' : 'mode-flat'"
        @click="design.toggle3DMode()"
        :title="design.is3DMode ? 'Exit 3D preview' : 'Enter 3D preview'"
      >
        <svg v-if="!design.is3DMode" viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2L2 7l10 5 10-5-10-5z"/><path d="M2 17l10 5 10-5"/><path d="M2 12l10 5 10-5"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
        {{ design.is3DMode ? 'Edit Mode' : '3D Preview' }}
      </button>

      <!-- 🔒 3D MODE BADGE -->
      <div v-if="design.is3DMode" class="badge-3d-mode">
        <span class="hidden sm:inline">🎲 Drag to rotate · Front view only · Click "Edit Mode" to resume editing</span>
        <span class="sm:hidden">🎲 Drag to rotate 3D view</span>
      </div>

      <!-- Board scaler: sized to match scale -->
      <div
        class="board-scaler"
        :class="{ 'board-scaler--3d': design.is3DMode }"
        :style="scalerTransform
          ? { width: Math.round(design.boardW * design.boardScale) + 'px', height: Math.round(design.boardH * design.boardScale) + 'px', transform: scalerTransform }
          : { width: Math.round(design.boardW * design.boardScale) + 'px', height: Math.round(design.boardH * design.boardScale) + 'px' }"
      >
        <svg width="0" height="0" style="position:absolute">
          <defs>
            <g id="fw">
              <path d="M0 -12 C 4 -12 4 -4 12 -12 C 4 -4 12 -4 12 0 C 12 4 4 4 12 12 C 4 4 4 12 0 12 C -4 12 -4 4 -12 12 C -4 4 -12 4 -12 0 C -12 -4 -4 -4 -12 -12 C -4 -4 -4 -12 0 -12" fill="currentColor"/>
              <circle cx="0" cy="0" r="4" fill="#ffd700" opacity="0.9"/>
            </g>
            <g id="lf">
              <path d="M0 0 Q 15 -10 25 0 Q 15 10 0 0 Z" fill="#2d5a27"/>
            </g>
          </defs>
        </svg>

        <!-- 🌸 TOP FLORAL CREST 🌸 -->
        <div v-if="design.topCrest.enabled" class="floral-crest fc-top"
             :style="{ '--p': design.topCrest.primary, '--s': design.topCrest.secondary, width: design.topCrest.size + '%' }">
          <svg v-if="design.topCrest.style === 'classic'" viewBox="0 0 300 150" class="crest-svg">
            <use href="#lf" x="60" y="145" transform="rotate(-30 60 145) scale(2)"/><use href="#lf" x="240" y="145" transform="rotate(210 240 145) scale(2)"/><use href="#lf" x="150" y="55" transform="rotate(-90 150 55) scale(2)"/>
            <g fill="var(--p)">
              <use href="#fw" x="70" y="150"/><use href="#fw" x="74" y="125"/><use href="#fw" x="84" y="103"/><use href="#fw" x="103" y="84"/><use href="#fw" x="125" y="74"/><use href="#fw" x="150" y="70"/><use href="#fw" x="175" y="74"/><use href="#fw" x="197" y="84"/><use href="#fw" x="216" y="103"/><use href="#fw" x="226" y="125"/><use href="#fw" x="230" y="150"/>
            </g>
            <g fill="var(--s)">
              <use href="#fw" x="90" y="150"/><use href="#fw" x="96" y="124"/><use href="#fw" x="115" y="101"/><use href="#fw" x="140" y="91"/><use href="#fw" x="160" y="91"/><use href="#fw" x="185" y="101"/><use href="#fw" x="204" y="124"/><use href="#fw" x="210" y="150"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="110" y="150"/><use href="#fw" x="122" y="122"/><use href="#fw" x="150" y="110"/><use href="#fw" x="178" y="122"/><use href="#fw" x="190" y="150"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="130" y="150"/><use href="#fw" x="150" y="130"/><use href="#fw" x="170" y="150"/>
            </g>
          </svg>
          <svg v-else-if="design.topCrest.style === 'modern'" viewBox="0 0 400 150" class="crest-svg">
            <use href="#lf" x="40" y="145" transform="rotate(-30 40 145) scale(2)"/><use href="#lf" x="360" y="145" transform="rotate(210 360 145) scale(2)"/>
            <g fill="var(--p)">
              <use href="#fw" x="60" y="150"/><use href="#fw" x="85" y="150"/><use href="#fw" x="110" y="150"/><use href="#fw" x="135" y="150"/><use href="#fw" x="160" y="150"/><use href="#fw" x="185" y="150"/><use href="#fw" x="215" y="150"/><use href="#fw" x="240" y="150"/><use href="#fw" x="265" y="150"/><use href="#fw" x="290" y="150"/><use href="#fw" x="315" y="150"/><use href="#fw" x="340" y="150"/>
              <use href="#fw" x="80" y="130"/><use href="#fw" x="100" y="110"/><use href="#fw" x="120" y="90"/><use href="#fw" x="140" y="70"/><use href="#fw" x="320" y="130"/><use href="#fw" x="300" y="110"/><use href="#fw" x="280" y="90"/><use href="#fw" x="260" y="70"/>
            </g>
            <g fill="var(--s)">
              <use href="#fw" x="110" y="130"/><use href="#fw" x="130" y="110"/><use href="#fw" x="150" y="90"/><use href="#fw" x="170" y="70"/><use href="#fw" x="190" y="55"/><use href="#fw" x="210" y="55"/><use href="#fw" x="230" y="70"/><use href="#fw" x="250" y="90"/><use href="#fw" x="270" y="110"/><use href="#fw" x="290" y="130"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="140" y="130"/><use href="#fw" x="160" y="110"/><use href="#fw" x="180" y="90"/><use href="#fw" x="200" y="75"/><use href="#fw" x="220" y="90"/><use href="#fw" x="240" y="110"/><use href="#fw" x="260" y="130"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="170" y="130"/><use href="#fw" x="190" y="110"/><use href="#fw" x="210" y="110"/><use href="#fw" x="230" y="130"/><use href="#fw" x="200" y="130"/>
            </g>
          </svg>
          <svg v-else-if="design.topCrest.style === 'grand'" viewBox="0 0 500 150" class="crest-svg">
            <g fill="var(--s)">
              <use href="#fw" x="40" y="150"/><use href="#fw" x="50" y="120"/><use href="#fw" x="70" y="95"/><use href="#fw" x="100" y="90"/><use href="#fw" x="130" y="95"/><use href="#fw" x="150" y="120"/><use href="#fw" x="160" y="150"/>
              <use href="#fw" x="340" y="150"/><use href="#fw" x="350" y="120"/><use href="#fw" x="370" y="95"/><use href="#fw" x="400" y="90"/><use href="#fw" x="430" y="95"/><use href="#fw" x="450" y="120"/><use href="#fw" x="460" y="150"/>
            </g>
            <g fill="var(--p)">
              <use href="#fw" x="60" y="150"/><use href="#fw" x="80" y="120"/><use href="#fw" x="100" y="110"/><use href="#fw" x="120" y="120"/><use href="#fw" x="140" y="150"/>
              <use href="#fw" x="360" y="150"/><use href="#fw" x="380" y="120"/><use href="#fw" x="400" y="110"/><use href="#fw" x="420" y="120"/><use href="#fw" x="440" y="150"/>
              <use href="#fw" x="170" y="150"/><use href="#fw" x="174" y="125"/><use href="#fw" x="184" y="103"/><use href="#fw" x="203" y="84"/><use href="#fw" x="225" y="74"/><use href="#fw" x="250" y="70"/><use href="#fw" x="275" y="74"/><use href="#fw" x="297" y="84"/><use href="#fw" x="316" y="103"/><use href="#fw" x="326" y="125"/><use href="#fw" x="330" y="150"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="190" y="150"/><use href="#fw" x="196" y="124"/><use href="#fw" x="215" y="101"/><use href="#fw" x="240" y="91"/><use href="#fw" x="260" y="91"/><use href="#fw" x="285" y="101"/><use href="#fw" x="304" y="124"/><use href="#fw" x="310" y="150"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="210" y="150"/><use href="#fw" x="222" y="122"/><use href="#fw" x="250" y="110"/><use href="#fw" x="278" y="122"/><use href="#fw" x="290" y="150"/>
            </g>
          </svg>
        </div>

        <!-- 🌸 BOTTOM FLORAL CREST 🌸 -->
        <div v-if="design.bottomCrest.enabled" class="floral-crest fc-bottom"
             :style="{ '--p': design.bottomCrest.primary, '--s': design.bottomCrest.secondary, width: design.bottomCrest.size + '%' }">
          <svg v-if="design.bottomCrest.style === 'classic'" viewBox="0 0 300 150" class="crest-svg">
            <use href="#lf" x="60" y="145" transform="rotate(-30 60 145) scale(2)"/><use href="#lf" x="240" y="145" transform="rotate(210 240 145) scale(2)"/><use href="#lf" x="150" y="55" transform="rotate(-90 150 55) scale(2)"/>
            <g fill="var(--p)">
              <use href="#fw" x="70" y="150"/><use href="#fw" x="74" y="125"/><use href="#fw" x="84" y="103"/><use href="#fw" x="103" y="84"/><use href="#fw" x="125" y="74"/><use href="#fw" x="150" y="70"/><use href="#fw" x="175" y="74"/><use href="#fw" x="197" y="84"/><use href="#fw" x="216" y="103"/><use href="#fw" x="226" y="125"/><use href="#fw" x="230" y="150"/>
            </g>
            <g fill="var(--s)">
              <use href="#fw" x="90" y="150"/><use href="#fw" x="96" y="124"/><use href="#fw" x="115" y="101"/><use href="#fw" x="140" y="91"/><use href="#fw" x="160" y="91"/><use href="#fw" x="185" y="101"/><use href="#fw" x="204" y="124"/><use href="#fw" x="210" y="150"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="110" y="150"/><use href="#fw" x="122" y="122"/><use href="#fw" x="150" y="110"/><use href="#fw" x="178" y="122"/><use href="#fw" x="190" y="150"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="130" y="150"/><use href="#fw" x="150" y="130"/><use href="#fw" x="170" y="150"/>
            </g>
          </svg>
          <svg v-else-if="design.bottomCrest.style === 'modern'" viewBox="0 0 400 150" class="crest-svg">
            <use href="#lf" x="40" y="145" transform="rotate(-30 40 145) scale(2)"/><use href="#lf" x="360" y="145" transform="rotate(210 360 145) scale(2)"/>
            <g fill="var(--p)">
              <use href="#fw" x="60" y="150"/><use href="#fw" x="85" y="150"/><use href="#fw" x="110" y="150"/><use href="#fw" x="135" y="150"/><use href="#fw" x="160" y="150"/><use href="#fw" x="185" y="150"/><use href="#fw" x="215" y="150"/><use href="#fw" x="240" y="150"/><use href="#fw" x="265" y="150"/><use href="#fw" x="290" y="150"/><use href="#fw" x="315" y="150"/><use href="#fw" x="340" y="150"/>
              <use href="#fw" x="80" y="130"/><use href="#fw" x="100" y="110"/><use href="#fw" x="120" y="90"/><use href="#fw" x="140" y="70"/><use href="#fw" x="320" y="130"/><use href="#fw" x="300" y="110"/><use href="#fw" x="280" y="90"/><use href="#fw" x="260" y="70"/>
            </g>
            <g fill="var(--s)">
              <use href="#fw" x="110" y="130"/><use href="#fw" x="130" y="110"/><use href="#fw" x="150" y="90"/><use href="#fw" x="170" y="70"/><use href="#fw" x="190" y="55"/><use href="#fw" x="210" y="55"/><use href="#fw" x="230" y="70"/><use href="#fw" x="250" y="90"/><use href="#fw" x="270" y="110"/><use href="#fw" x="290" y="130"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="140" y="130"/><use href="#fw" x="160" y="110"/><use href="#fw" x="180" y="90"/><use href="#fw" x="200" y="75"/><use href="#fw" x="220" y="90"/><use href="#fw" x="240" y="110"/><use href="#fw" x="260" y="130"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="170" y="130"/><use href="#fw" x="190" y="110"/><use href="#fw" x="210" y="110"/><use href="#fw" x="230" y="130"/><use href="#fw" x="200" y="130"/>
            </g>
          </svg>
          <svg v-else-if="design.bottomCrest.style === 'grand'" viewBox="0 0 500 150" class="crest-svg">
            <g fill="var(--s)">
              <use href="#fw" x="40" y="150"/><use href="#fw" x="50" y="120"/><use href="#fw" x="70" y="95"/><use href="#fw" x="100" y="90"/><use href="#fw" x="130" y="95"/><use href="#fw" x="150" y="120"/><use href="#fw" x="160" y="150"/>
              <use href="#fw" x="340" y="150"/><use href="#fw" x="350" y="120"/><use href="#fw" x="370" y="95"/><use href="#fw" x="400" y="90"/><use href="#fw" x="430" y="95"/><use href="#fw" x="450" y="120"/><use href="#fw" x="460" y="150"/>
            </g>
            <g fill="var(--p)">
              <use href="#fw" x="60" y="150"/><use href="#fw" x="80" y="120"/><use href="#fw" x="100" y="110"/><use href="#fw" x="120" y="120"/><use href="#fw" x="140" y="150"/>
              <use href="#fw" x="360" y="150"/><use href="#fw" x="380" y="120"/><use href="#fw" x="400" y="110"/><use href="#fw" x="420" y="120"/><use href="#fw" x="440" y="150"/>
              <use href="#fw" x="170" y="150"/><use href="#fw" x="174" y="125"/><use href="#fw" x="184" y="103"/><use href="#fw" x="203" y="84"/><use href="#fw" x="225" y="74"/><use href="#fw" x="250" y="70"/><use href="#fw" x="275" y="74"/><use href="#fw" x="297" y="84"/><use href="#fw" x="316" y="103"/><use href="#fw" x="326" y="125"/><use href="#fw" x="330" y="150"/>
            </g>
            <g fill="#f1faee">
              <use href="#fw" x="190" y="150"/><use href="#fw" x="196" y="124"/><use href="#fw" x="215" y="101"/><use href="#fw" x="240" y="91"/><use href="#fw" x="260" y="91"/><use href="#fw" x="285" y="101"/><use href="#fw" x="304" y="124"/><use href="#fw" x="310" y="150"/>
            </g>
            <g fill="#f5c842">
              <use href="#fw" x="210" y="150"/><use href="#fw" x="222" y="122"/><use href="#fw" x="250" y="110"/><use href="#fw" x="278" y="122"/><use href="#fw" x="290" y="150"/>
            </g>
          </svg>
        </div>

        <!-- ╔══ THE BOARD ══╗ -->
        <div
          :ref="(el) => { (design as any).boardRef = el }"
          class="board-frame"
          :style="{
            width: design.boardW + 'px', height: design.boardH + 'px',
            transform: `scale(${design.boardScale})`, transformOrigin: 'top left',
            cursor: design.is3DMode ? 'inherit' : (design.isBrushMode ? 'crosshair' : 'default'),
            touchAction: 'none',
            ...design.boardBorderStyle,
            ...design.boardCornerStyle,
          }"
          @click="!design.is3DMode && design.handleBoardClick($event)"
          @dragover.prevent
          @drop="!design.is3DMode && design.handleDrop($event)"
        >

          <!-- ▲ UPPER SECTION -->
          <div class="board-section"
            :style="{ top: 0, left: 0, right: 0, height: design.upperH + 'px', backgroundColor: design.upper.bgColor, ...design.upperCornerStyle }">
            <template v-if="design.upper.cornerStyle === 'floral'">
              <span class="floral-corner fc-tl">🌸</span>
              <span class="floral-corner fc-tr">🌸</span>
            </template>
            <div class="sec-inner">
              <div v-if="design.upper.headerText" class="sec-text sec-header"
                :style="{ fontSize: design.upper.headerFontSize + 'px', fontFamily: design.getFont(design.upper.headerFont), color: design.upper.headerColor, textAlign: design.upper.headerAlign }">
                {{ design.upper.headerText }}
              </div>
              <!-- Header border rule -->
              <div v-if="design.upper.headerBorder" class="sec-text-border"
                :style="{ borderTopWidth: (design.upper.headerBorderWidth ?? 2) + 'px', borderTopStyle: 'solid', borderTopColor: design.upper.headerBorderColor ?? design.upper.headerColor }">
              </div>
              <!-- Body border rule -->
              <div v-if="design.upper.bodyBorder" class="sec-text-border"
                :style="{ borderTopWidth: (design.upper.bodyBorderWidth ?? 2) + 'px', borderTopStyle: 'solid', borderTopColor: design.upper.bodyBorderColor ?? design.upper.bodyColor }">
              </div>
              <div class="sec-text sec-body"
                :style="{ fontSize: design.upper.bodyFontSize + 'px', fontFamily: design.getFont(design.upper.bodyFont), color: design.upper.bodyColor, textAlign: design.upper.bodyAlign }">
                {{ design.upper.bodyText }}
              </div>
            </div>
          </div>

          <!-- ── SECTION DIVIDER (drag to resize) ── -->
          <div class="section-divider" :style="{ top: design.upperH + 'px', '--dc': design.border.color || '#ccc', touchAction: 'none' }"
            @mousedown.stop="!design.is3DMode && design.startDragDiv($event)"
            @touchstart.stop="!design.is3DMode && design.startDragDiv($event)"
            @click.stop>
            <div class="div-track" :style="!design.border.center ? { opacity: 0.6 } : { display: 'none' }"/>
            <div v-if="design.border.center" :style="{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center' }">
              <div :style="design.centerBorderStyle"></div>
            </div>
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
            :style="{ top: design.upperH + 'px', left: 0, right: 0, height: design.lowerH + 'px', backgroundColor: design.lower.bgColor, overflow: 'visible', ...design.lowerCornerStyle }">
            <template v-if="design.lower.cornerStyle === 'floral'">
              <span class="floral-corner fc-bl">🌸</span>
              <span class="floral-corner fc-br">🌸</span>
            </template>
            <div class="sec-inner">
              <div v-if="design.lower.headerText" class="sec-text sec-header"
                :style="{ fontSize: design.lower.headerFontSize + 'px', fontFamily: design.getFont(design.lower.headerFont), color: design.lower.headerColor, textAlign: design.lower.headerAlign }">
                {{ design.lower.headerText }}
              </div>
              <!-- Header border rule -->
              <div v-if="design.lower.headerBorder" class="sec-text-border"
                :style="{ borderTopWidth: (design.lower.headerBorderWidth ?? 2) + 'px', borderTopStyle: 'solid', borderTopColor: design.lower.headerBorderColor ?? design.lower.headerColor }">
              </div>
              <!-- Body border rule -->
              <div v-if="design.lower.bodyBorder" class="sec-text-border"
                :style="{ borderTopWidth: (design.lower.bodyBorderWidth ?? 2) + 'px', borderTopStyle: 'solid', borderTopColor: design.lower.bodyBorderColor ?? design.lower.bodyColor }">
              </div>
              <div class="sec-text sec-body"
                :style="{ fontSize: design.lower.bodyFontSize + 'px', fontFamily: design.getFont(design.lower.bodyFont), color: design.lower.bodyColor, textAlign: design.lower.bodyAlign }">
                {{ design.lower.bodyText }}
              </div>
            </div>
          </div>

          <!-- ★ CANVAS ELEMENTS — array order = z-order (last = top) -->
          <template v-for="(el, idx) in design.elements" :key="el.id">

            <!-- Image element -->
            <div v-if="el.type === 'image'" class="canvas-el"
              :class="{ 'el-selected': !design.is3DMode && design.selectedId === el.id }"
              :style="{
                left: (el as CanvasImage).x + '%', top: (el as CanvasImage).y + '%',
                width: (el as CanvasImage).width + '%', aspectRatio: '1/1',
                zIndex: (el as CanvasImage).zIndex ?? (idx + 10), overflow: 'hidden',
                borderRadius: (el as CanvasImage).frame === 'circle' ? '50%' : (el as CanvasImage).frame === 'square' ? '4px' : '0',
                pointerEvents: design.is3DMode ? 'none' : (design.isBrushMode ? 'none' : 'auto'), cursor: 'grab',
                touchAction: 'none'
              }"
              @mousedown.stop="!design.is3DMode && design.startDragEl($event, el.id)"
              @touchstart.stop="!design.is3DMode && design.startDragEl($event, el.id)">
              <img :src="(el as CanvasImage).src" draggable="false"
                :style="{
                  width: '100%', height: '100%', objectFit: 'cover', display: 'block',
                  objectPosition: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                  transform: 'scale(' + (el as CanvasImage).zoom + ')',
                  transformOrigin: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                }"/>
              <div v-if="!design.is3DMode && design.selectedId === el.id && !design.isBrushMode" class="el-del" @click.stop="design.deleteSelected" title="Remove">×</div>
            </div>

            <!-- Brush stroke element -->
            <div v-else-if="el.type === 'brush'" class="canvas-el"
              :class="{ 'el-selected': !design.is3DMode && design.selectedId === el.id }"
              :style="{
                left: (el as BrushStroke).x + '%', top: (el as BrushStroke).y + '%',
                width: (el as BrushStroke).size + 'px', height: (el as BrushStroke).size + 'px',
                transform: `translate(-50%,-50%) rotate(${(el as BrushStroke).rotation}deg)`,
                zIndex: (el as BrushStroke).zIndex ?? (idx + 10), color: (el as BrushStroke).color,
                pointerEvents: design.is3DMode ? 'none' : 'auto', cursor: design.isBrushMode ? 'crosshair' : 'pointer',
                touchAction: 'none'
              }"
              @mousedown.stop="!design.is3DMode && design.handleBrushMousedown($event, el.id)"
              @touchstart.stop="!design.is3DMode && design.handleBrushMousedown($event, el.id)">
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
              <div v-if="!design.is3DMode && design.selectedId === el.id" class="el-del" @click.stop="design.deleteSelected" title="Remove">×</div>
            </div>
          </template>

          <!-- 3D shading overlay — sits on top of all board content -->
          <div class="board-3d-shading" :style="design.board3dShadingStyle"></div>
        </div><!-- /board-frame -->
      </div><!-- /board-scaler -->
    </div><!-- /chess-bg -->
  </div>
</template>
