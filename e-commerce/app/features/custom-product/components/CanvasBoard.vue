<script setup lang="ts">
import type { CanvasImage, BrushStroke } from '../types'
import { useCart } from '~/composables/useCart'

const props = defineProps<{
  design: ReturnType<typeof import('../useCustomDesign').useCustomDesign>
}>()

const { formatRupiah } = useCart()
</script>

<template>
  <div class="dr-canvas-area" ref="design.containerRef">
    <div class="chess-bg" @dragover.prevent>
      <!-- Board scaler: sized to match scale -->
      <div class="board-scaler"
        :style="{ width: Math.round(design.boardW * design.boardScale) + 'px', height: Math.round(design.boardH * design.boardScale) + 'px' }">
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
          ref="design.boardRef"
          class="board-frame"
          :style="{
            width: design.boardW + 'px', height: design.boardH + 'px',
            transform: `scale(${design.boardScale})`, transformOrigin: 'top left',
            cursor: design.isBrushMode ? 'crosshair' : 'default',
            ...design.boardBorderStyle,
          }"
          @click="design.handleBoardClick"
          @dragover.prevent
          @drop="design.handleDrop"
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
              <div class="sec-text sec-body"
                :style="{ fontSize: design.upper.bodyFontSize + 'px', fontFamily: design.getFont(design.upper.bodyFont), color: design.upper.bodyColor, textAlign: design.upper.bodyAlign }">
                {{ design.upper.bodyText }}
              </div>
            </div>
          </div>

          <!-- ── SECTION DIVIDER (drag to resize) ── -->
          <div class="section-divider" :style="{ top: design.upperH + 'px', '--dc': design.border.color || '#ccc' }"
            @mousedown.stop="design.startDragDiv" @click.stop>
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
              :class="{ 'el-selected': design.selectedId === el.id }"
              :style="{
                left: (el as CanvasImage).x + '%', top: (el as CanvasImage).y + '%',
                width: (el as CanvasImage).width + '%', aspectRatio: '1/1',
                zIndex: (el as CanvasImage).zIndex ?? (idx + 10), overflow: 'hidden',
                borderRadius: (el as CanvasImage).frame === 'circle' ? '50%' : (el as CanvasImage).frame === 'square' ? '4px' : '0',
                pointerEvents: design.isBrushMode ? 'none' : 'auto', cursor: 'grab',
              }"
              @mousedown.stop="design.startDragEl($event, el.id)">
              <img :src="(el as CanvasImage).src" draggable="false"
                :style="{
                  width: '100%', height: '100%', objectFit: 'cover', display: 'block',
                  objectPosition: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                  transform: 'scale(' + (el as CanvasImage).zoom + ')',
                  transformOrigin: (el as CanvasImage).cropX + '% ' + (el as CanvasImage).cropY + '%',
                }"/>
              <div v-if="design.selectedId === el.id && !design.isBrushMode" class="el-del" @click.stop="design.deleteSelected" title="Remove">×</div>
            </div>

            <!-- Brush stroke element -->
            <div v-else-if="el.type === 'brush'" class="canvas-el"
              :class="{ 'el-selected': design.selectedId === el.id }"
              :style="{
                left: (el as BrushStroke).x + '%', top: (el as BrushStroke).y + '%',
                width: (el as BrushStroke).size + 'px', height: (el as BrushStroke).size + 'px',
                transform: `translate(-50%,-50%) rotate(${(el as BrushStroke).rotation}deg)`,
                zIndex: (el as BrushStroke).zIndex ?? (idx + 10), color: (el as BrushStroke).color,
                pointerEvents: 'auto', cursor: design.isBrushMode ? 'crosshair' : 'pointer',
              }"
              @mousedown.stop="design.handleBrushMousedown($event, el.id)">
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
              <div v-if="design.selectedId === el.id" class="el-del" @click.stop="design.deleteSelected" title="Remove">×</div>
            </div>
          </template>
        </div><!-- /board-frame -->
      </div><!-- /board-scaler -->
    </div><!-- /chess-bg -->

    <!-- Canvas info bar -->
    <div class="canvas-info">
      <span>{{ design.elements.length }} element{{ design.elements.length !== 1 ? 's' : '' }}</span>
      <span v-if="design.selectedId" class="ci-sel"> · 1 selected</span>
      <button v-if="design.selectedId" class="ci-desel" @click="design.selectedId = null">Deselect</button>
      <span v-if="design.isBrushMode" class="ci-hint">Click canvas to place brush stroke · Del to remove selected</span>
      <span v-else-if="!design.isBrushMode && design.selectedId" class="ci-hint">Drag to move · Del to remove</span>
    </div>
  </div>
</template>
