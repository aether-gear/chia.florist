<script setup lang="ts">
import type { FrameStyle, FloralStyle } from '../types'
import { FONTS, CORNERS, BORDER_STYLES, SIZES, BRUSH_COLORS, BORDER_COLORS, BG_PRESETS, TOOL_TABS } from '../constants'
import { useCart } from '~/composables/useCart'

const props = defineProps<{
  design: ReturnType<typeof import('../useCustomDesign').useCustomDesign>
}>()

const emit = defineEmits<{
  (e: 'finalize'): void
}>()

const { formatRupiah } = useCart()
</script>

<template>
  <aside class="dr-panel">

    <!-- Tab bar -->
    <div class="tab-bar" role="tablist">
      <button v-for="tab in TOOL_TABS" :key="tab.id"
        class="tab-btn" :class="{ 'tab-active': design.activeTab === tab.id }"
        @click="design.activeTab = tab.id" role="tab" :aria-selected="design.activeTab === tab.id"
        :id="'tab-' + tab.id">
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
        <svg v-else-if="tab.id === 'corner'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <path d="M3 9V5a2 2 0 012-2h4"/><path d="M3 15v4a2 2 0 002 2h4"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 22C12 22 20 18 20 12C20 6 12 2 12 2C12 2 4 6 4 12C4 18 12 22 12 22Z"/><circle cx="12" cy="12" r="3"/>
        </svg>
        <span>{{ tab.label }}</span>
      </button>
    </div>

    <!-- Tab content area -->
    <div class="tab-body">

      <!-- ╔══ TEXT TAB ══╗ -->
      <div v-if="design.activeTab === 'text'" class="tab-pane">
        <div class="sec-toggle">
          <button class="stg-btn" :class="{ active: design.activeSection === 'upper' }" @click="design.activeSection = 'upper'">Upper</button>
          <button class="stg-btn" :class="{ active: design.activeSection === 'lower' }" @click="design.activeSection = 'lower'">Lower</button>
        </div>

        <!-- Background -->
        <div class="tg">
          <div class="tg-label">BACKGROUND</div>
          <div class="color-row">
            <input id="bg-color-input" type="color" v-model="design.sec.bgColor" class="csi"/>
            <span class="cval">{{ design.sec.bgColor }}</span>
          </div>
          <div class="dot-row">
            <button v-for="c in BG_PRESETS" :key="c" class="pdot"
              :style="{ background: c, outline: design.sec.bgColor === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#ffffff' ? 'inset 0 0 0 1px #ccc' : 'none' }"
              @click="design.sec.bgColor = c"/>
          </div>
        </div>

        <!-- Header -->
        <div class="tg">
          <div class="tg-label">HEADER</div>
          <textarea class="dr-ta" v-model="design.sec.headerText" placeholder="Header text…" rows="2"/>
          <div class="cr">
            <label class="clabel">Size</label>
            <input type="range" min="10" max="96" class="dr-range" v-model.number="design.sec.headerFontSize"/>
            <span class="cval">{{ design.sec.headerFontSize }}px</span>
          </div>
          <div class="font-grid">
            <button v-for="f in FONTS" :key="f.id" class="font-chip" :class="{ 'fc-active': design.sec.headerFont === f.id }"
              :style="{ fontFamily: f.family }" @click="design.sec.headerFont = f.id">{{ f.label }}</button>
          </div>
          <div class="align-row">
            <button v-for="a in (['left','center','right'] as const)" :key="a" class="aln-btn"
              :class="{ 'aln-active': design.sec.headerAlign === a }" @click="design.sec.headerAlign = a"
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
            <input type="color" v-model="design.sec.headerColor" class="csi"/>
            <span class="cval">{{ design.sec.headerColor }}</span>
          </div>
        </div>

        <!-- Body -->
        <div class="tg">
          <div class="tg-label">BODY</div>
          <textarea class="dr-ta" v-model="design.sec.bodyText" placeholder="Body text… (new line = Enter)" rows="3"/>
          <div class="cr">
            <label class="clabel">Size</label>
            <input type="range" min="8" max="72" class="dr-range" v-model.number="design.sec.bodyFontSize"/>
            <span class="cval">{{ design.sec.bodyFontSize }}px</span>
          </div>
          <div class="font-grid">
            <button v-for="f in FONTS" :key="f.id" class="font-chip" :class="{ 'fc-active': design.sec.bodyFont === f.id }"
              :style="{ fontFamily: f.family }" @click="design.sec.bodyFont = f.id">{{ f.label }}</button>
          </div>
          <div class="align-row">
            <button v-for="a in (['left','center','right'] as const)" :key="a" class="aln-btn"
              :class="{ 'aln-active': design.sec.bodyAlign === a }" @click="design.sec.bodyAlign = a" :title="a">
              <svg viewBox="0 0 16 12" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
                <line x1="0" y1="2" x2="16" y2="2"/>
                <line :x1="a==='left'?0:a==='center'?3:6" y1="6" :x2="a==='left'?10:a==='center'?13:16" y2="6"/>
                <line x1="0" y1="10" x2="16" y2="10"/>
              </svg>
            </button>
          </div>
          <div class="color-row">
            <label class="clabel">Color</label>
            <input type="color" v-model="design.sec.bodyColor" class="csi"/>
            <span class="cval">{{ design.sec.bodyColor }}</span>
          </div>
        </div>
      </div>

      <!-- ╔══ IMAGE TAB ══╗ -->
      <div v-else-if="design.activeTab === 'image'" class="tab-pane">
        <template v-if="design.imgElements.length === 0">
          <div class="drop-zone" @dragover.prevent @drop="design.handleDrop">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="dz-icon">
              <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
            </svg>
            <p class="dz-title">Drop image here</p>
            <p class="dz-sub">or drag from files onto the board</p>
            <label class="dz-browse" for="dz-file-input">Browse Files</label>
            <input id="dz-file-input" type="file" accept="image/*" @change="design.handleFileInput" style="display:none"/>
          </div>
          <p class="dz-note">Images are registered as canvas components. Use frame, zoom and crop to style them.</p>
        </template>

        <template v-else>
          <div style="padding: 1rem 1rem 0;">
            <label class="primary-btn" style="display:flex; justify-content:center; cursor:pointer; align-items:center; background:#4a4a4a; color:#fff; padding:0.6rem; border-radius:4px; font-weight:500; font-size:0.85rem">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" style="width:16px;height:16px;margin-right:6px;"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              Upload Another Image
              <input type="file" accept="image/*" @change="design.handleFileInput" style="display:none"/>
            </label>
          </div>

          <template v-if="design.selectedImg">
            <div class="img-preview-wrap">
              <img :src="design.selectedImg.src" class="img-preview-thumb"
                :style="{ borderRadius: design.selectedImg.frame === 'circle' ? '50%' : '4px' }"/>
            </div>

            <div class="tg">
              <div class="tg-label">FRAME</div>
              <div class="frame-row">
                <button v-for="f in (['none','square','circle'] as FrameStyle[])" :key="f"
                  class="frame-btn" :class="{ 'fb-active': design.selectedImg.frame === f }"
                  @click="design.selectedImg.frame = f">
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
                <input type="range" min="5" max="80" class="dr-range" v-model.number="design.selectedImg.width"/>
                <span class="cval">{{ design.selectedImg.width }}%</span>
              </div>
            </div>

            <div class="tg">
              <div class="tg-label">ZOOM &amp; CROP</div>
              <div class="cr">
                <label class="clabel">Zoom</label>
                <input type="range" min="1" max="3" step="0.05" class="dr-range" v-model.number="design.selectedImg.zoom"/>
                <span class="cval">{{ design.selectedImg.zoom.toFixed(1) }}×</span>
              </div>
              <div class="cr">
                <label class="clabel">Crop X</label>
                <input type="range" min="0" max="100" class="dr-range" v-model.number="design.selectedImg.cropX"/>
                <span class="cval">{{ design.selectedImg.cropX }}%</span>
              </div>
              <div class="cr">
                <label class="clabel">Crop Y</label>
                <input type="range" min="0" max="100" class="dr-range" v-model.number="design.selectedImg.cropY"/>
                <span class="cval">{{ design.selectedImg.cropY }}%</span>
              </div>
            </div>

            <button class="danger-btn" @click="design.deleteSelected">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4h6v2"/></svg>
              Remove Image
            </button>
          </template>
        </template>

        <div v-if="design.imgElements.length" class="el-section">
          <div class="tg-label" style="padding:0.875rem 1rem 0.4rem">ALL IMAGES ({{ design.imgElements.length }})</div>
          <div class="el-list">
            <div v-for="img in design.imgElements" :key="img.id" class="el-item"
              :class="{ 'el-active': design.selectedId === img.id }" @click="design.bringToFront(img.id); design.activeTab='image'">
              <img :src="img.src" class="el-thumb" :style="{ borderRadius: img.frame === 'circle' ? '50%' : '3px' }"/>
              <div class="el-meta"><span>{{ img.frame }} frame</span><span>{{ img.width }}% wide</span></div>
              <button class="el-del-list" @click.stop="design.selectedId = img.id; design.deleteSelected()">×</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ╔══ BRUSH TAB ══╗ -->
      <div v-else-if="design.activeTab === 'brush'" class="tab-pane">
        <div class="brush-info">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <p>Click empty canvas to stamp flowers. Drag or click any placed flower to move, rotate, &amp; edit it live!</p>
        </div>

        <div class="tg">
          <div class="tg-label">TYPE</div>
          <div class="brush-types">
            <button class="btype-btn" :class="{ 'btype-active': design.brushType === 'flower' }" @click="design.brushType = 'flower'">
              <svg viewBox="-20 -20 40 40" width="44" height="44" :style="{ color: design.brushColor }">
                <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(0,0,0)"/>
                <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(72,0,0)"/>
                <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(144,0,0)"/>
                <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(216,0,0)"/>
                <ellipse cx="0" cy="-10" rx="5" ry="9" fill="currentColor" transform="rotate(288,0,0)"/>
                <circle cx="0" cy="0" r="5" fill="#ffe066"/>
              </svg>
              <span>Flower</span>
            </button>
            <button class="btype-btn" :class="{ 'btype-active': design.brushType === 'rose' }" @click="design.brushType = 'rose'">
              <svg viewBox="-20 -20 40 40" width="44" height="44" :style="{ color: design.brushColor }">
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
              :style="{ background: c, outline: design.brushColor === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#ffffff' ? 'inset 0 0 0 1px #ccc' : 'none' }"
              @click="design.brushColor = c"/>
          </div>
          <div class="color-row" style="margin-top:0.35rem">
            <label class="clabel">Custom</label>
            <input type="color" v-model="design.brushColor" class="csi"/>
            <span class="cval">{{ design.brushColor }}</span>
          </div>
        </div>

        <div class="tg">
          <div class="tg-label">SIZE &amp; ANGLE</div>
          <div class="cr">
            <label class="clabel">Size</label>
            <input type="range" min="16" max="120" class="dr-range" v-model.number="design.brushSize"/>
            <span class="cval">{{ design.brushSize }}px</span>
          </div>
          <div class="cr">
            <label class="clabel">Angle</label>
            <input type="range" min="0" max="360" class="dr-range" v-model.number="design.brushRotation"/>
            <span class="cval">{{ design.brushRotation }}°</span>
          </div>
        </div>

        <div class="brush-preview">
          <svg :viewBox="'-20 -20 40 40'" :width="Math.min(design.brushSize, 80)" :height="Math.min(design.brushSize, 80)"
            :style="{ color: design.brushColor, transform: `rotate(${design.brushRotation}deg)`, display: 'block' }">
            <template v-if="design.brushType === 'flower'">
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
          <span class="bp-label">{{ design.selectedBrush ? 'Editing selected' : 'Preview (next stamp)' }}</span>
        </div>

        <div v-if="design.brushElements.length" class="el-section">
          <div class="tg-label" style="padding:0.875rem 1rem 0.4rem">PLACED ({{ design.brushElements.length }})</div>
          <div class="el-list">
            <div v-for="br in design.brushElements" :key="br.id" class="el-item"
              :class="{ 'el-active': design.selectedId === br.id }"
              @click="design.selectedId = br.id; design.bringToFront(br.id)">
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
              <div style="flex:1;min-width:0">
                <span style="display:block;font-size:0.68rem">{{ br.brushType }} · {{ br.size }}px · {{ br.rotation }}°</span>
                <span style="display:block;font-size:0.65rem;color:#999">{{ design.selectedId === br.id ? '✏️ editing' : 'click to edit' }}</span>
              </div>
              <button class="el-del-list" @click.stop="design.selectedId = br.id; design.deleteSelected()">×</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ╔══ BORDER TAB ══╗ -->
      <div v-else-if="design.activeTab === 'border'" class="tab-pane">
        <div class="tg">
          <div class="tg-label">STYLE</div>
          <div class="border-grid">
            <button v-for="s in BORDER_STYLES" :key="s.id" class="bsb"
              :class="{ 'bsb-active': design.border.style === s.id }" @click="design.border.style = s.id">
              <span class="bsb-line" :style="{ borderBottomWidth:'3px', borderBottomStyle: s.id==='ornate'||s.id==='none'?'solid':s.id, borderBottomColor: s.id==='none'?'transparent':'#888' }"/>
              {{ s.label }}
            </button>
          </div>
        </div>
        <div class="tg">
          <div class="tg-label">COLOR</div>
          <div class="dot-row">
            <button v-for="c in BORDER_COLORS" :key="c" class="pdot"
              :style="{ background: c, outline: design.border.color === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px', boxShadow: c === '#f1faee' ? 'inset 0 0 0 1px #ccc' : 'none' }"
              @click="design.border.color = c"/>
          </div>
          <div class="color-row" style="margin-top:0.35rem">
            <label class="clabel">Custom</label>
            <input type="color" v-model="design.border.color" class="csi"/>
            <span class="cval">{{ design.border.color }}</span>
          </div>
        </div>
        <div class="tg">
          <div class="tg-label">WIDTH</div>
          <div class="cr">
            <label class="clabel">Size</label>
            <input type="range" min="0" max="32" class="dr-range" v-model.number="design.border.width"/>
            <span class="cval">{{ design.border.width }}px</span>
          </div>
        </div>

        <div class="tg" style="margin-top: 1rem;">
          <label style="display:flex; align-items:center; gap:0.6rem; font-size:0.75rem; font-weight:700; color:#555; cursor:pointer; letter-spacing: 0.05em;">
            <input type="checkbox" v-model="design.border.center" style="width:16px;height:16px;accent-color:#c4703e;"/>
            SHOW CENTER BORDER
          </label>
        </div>

        <div class="tg">
          <div class="tg-label">PREVIEW</div>
          <div class="border-preview"
            :style="{
              border: design.border.style !== 'none' && design.border.width > 0 ? `${design.border.width}px ${design.border.style === 'ornate' ? 'solid' : design.border.style} ${design.border.color}` : '1px dashed #ccc',
              outline: design.border.style === 'ornate' && design.border.width > 0 ? `${Math.max(1,Math.round(design.border.width*0.35))}px solid ${design.border.color}` : 'none',
              outlineOffset: '3px',
            }">Board Border
          </div>
        </div>
      </div>

      <!-- ╔══ CORNER TAB ══╗ -->
      <div v-else-if="design.activeTab === 'corner'" class="tab-pane">
        <div class="sec-toggle">
          <button class="stg-btn" :class="{ active: design.activeSection === 'upper' }" @click="design.activeSection = 'upper'">Upper</button>
          <button class="stg-btn" :class="{ active: design.activeSection === 'lower' }" @click="design.activeSection = 'lower'">Lower</button>
        </div>
        <div class="tg">
          <div class="tg-label">CORNER STYLE — {{ design.activeSection === 'upper' ? 'UPPER' : 'LOWER' }} SECTION</div>
          <div class="corner-grid">
            <button v-for="c in CORNERS" :key="c.id" class="corner-btn"
              :class="{ 'cb-active': design.sec.cornerStyle === c.id }"
              @click="design.sec.cornerStyle = c.id">
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

      <!-- ╔══ FLORAL TAB ══╗ -->
      <div v-else-if="design.activeTab === 'floral'" class="tab-pane">
        <div class="sec-toggle">
          <button class="stg-btn" :class="{ active: design.activeSection === 'upper' }" @click="design.activeSection = 'upper'">Top Crest</button>
          <button class="stg-btn" :class="{ active: design.activeSection === 'lower' }" @click="design.activeSection = 'lower'">Bottom Base</button>
        </div>

        <div class="tg" style="margin-top: 1rem;">
          <label style="display:flex; align-items:center; gap:0.6rem; font-size:0.75rem; font-weight:700; color:#555; cursor:pointer; letter-spacing:0.05em;">
            <input type="checkbox" v-model="design.floralSec.enabled" style="width:16px;height:16px;accent-color:#c4703e;"/>
            ENABLE FLORAL DECOR
          </label>
        </div>

        <template v-if="design.floralSec.enabled">
          <div class="tg">
            <div class="tg-label">STYLE</div>
            <div class="frame-row">
              <button v-for="s in (['classic','modern','grand'] as FloralStyle[])" :key="s"
                class="frame-btn" :class="{ 'fb-active': design.floralSec.style === s }"
                @click="design.floralSec.style = s" style="text-transform: capitalize;">
                {{ s }}
              </button>
            </div>
          </div>

          <div class="tg">
            <div class="tg-label">SIZE SCALING</div>
            <div class="cr">
              <input type="range" min="20" max="100" class="dr-range" v-model.number="design.floralSec.size"/>
              <span class="cval">{{ design.floralSec.size }}%</span>
            </div>
          </div>

          <div class="tg">
            <div class="tg-label">PRIMARY FLOWER COLOR</div>
            <div class="dot-row">
              <button v-for="c in BG_PRESETS" :key="c" class="pdot"
                :style="{ background: c, outline: design.floralSec.primary === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px' }"
                @click="design.floralSec.primary = c"/>
            </div>
            <div class="color-row" style="margin-top:0.35rem">
              <label class="clabel">Custom</label>
              <input type="color" v-model="design.floralSec.primary" class="csi"/>
              <span class="cval">{{ design.floralSec.primary }}</span>
            </div>
          </div>

          <div class="tg">
            <div class="tg-label">SECONDARY FLOWER COLOR</div>
            <div class="dot-row">
              <button v-for="c in BG_PRESETS" :key="c" class="pdot"
                :style="{ background: c, outline: design.floralSec.secondary === c ? '2px solid #c4703e' : '2px solid transparent', outlineOffset: '2px' }"
                @click="design.floralSec.secondary = c"/>
            </div>
            <div class="color-row" style="margin-top:0.35rem">
              <label class="clabel">Custom</label>
              <input type="color" v-model="design.floralSec.secondary" class="csi"/>
              <span class="cval">{{ design.floralSec.secondary }}</span>
            </div>
          </div>
        </template>
      </div>
    </div><!-- /tab-body -->

    <!-- ╔══ BOTTOM PANEL: SIZE & PRICE BREAKDOWN ══╗ -->
    <div class="panel-foot">
      <div class="size-group">
        <label class="size-group-label">BOARD PHYSICAL SIZE</label>
        <div class="size-row">
          <button v-for="sz in SIZES" :key="sz.id"
            class="sz-card" :class="{ 'sz-active': design.physicalSize === sz.id }"
            @click="design.physicalSize = sz.id">
            <span class="sz-name">{{ sz.label }}</span>
            <span class="sz-desc">{{ sz.desc }}</span>
            <span class="sz-price">{{ formatRupiah(sz.price) }}</span>
            <span v-if="sz.recommended" class="sz-badge">Best</span>
          </button>
        </div>
      </div>

      <!-- Live Dynamic Price Breakdown -->
      <div class="price-breakdown-card">
        <div class="pbc-header">
          <span class="pbc-title">PRICE ESTIMATE</span>
          <span class="pbc-total">{{ formatRupiah(design.totalPrice) }}</span>
        </div>
        <div class="pbc-rows">
          <div class="pbc-row">
            <span>Base Board Size ({{ SIZES.find(s=>s.id===design.physicalSize)?.label }})</span>
            <span>{{ formatRupiah(design.baseSizePrice) }}</span>
          </div>
          <div v-if="design.brushFee > 0" class="pbc-row">
            <span>Stamped Flowers ({{ design.brushElements.length }} × Rp 2.000)</span>
            <span>+ {{ formatRupiah(design.brushFee) }}</span>
          </div>
          <div v-if="design.colorFee > 0" class="pbc-row">
            <span>Palette Expansion ({{ design.uniqueColors.length }} colors)</span>
            <span>+ {{ formatRupiah(design.colorFee) }}</span>
          </div>
          <div v-if="design.accessoriesFee > 0" class="pbc-row">
            <span>Decorations &amp; Ornate Borders</span>
            <span>+ {{ formatRupiah(design.accessoriesFee) }}</span>
          </div>
          <div v-if="design.mediaFee > 0" class="pbc-row">
            <span>Image Components ({{ design.imgElements.length }} × Rp 20.000)</span>
            <span>+ {{ formatRupiah(design.mediaFee) }}</span>
          </div>
        </div>
      </div>

      <button class="primary-btn finalize-btn" @click="emit('finalize')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M5 12l5 5L20 7"/></svg>
        Finalize &amp; Order
      </button>
    </div>

  </aside>
</template>
