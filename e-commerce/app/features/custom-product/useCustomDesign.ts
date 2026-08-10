// app/features/custom-product/useCustomDesign.ts

import { ref, computed, watch, nextTick, reactive } from 'vue'
import type {
  BoardSection, BoardBorder, FloralCrest, CanvasElement,
  CanvasImage, BrushStroke, BrushType, CustomDesignPayloadV3,
  ToolTab, SectionKey, FloralStyle
} from './types'
import {
  FONTS, CORNERS, BORDER_STYLES, SIZES, BRUSH_COLORS,
  BORDER_COLORS, BG_PRESETS, DEFAULT_DRAFT_KEY
} from './constants'
import { normalizeHexColor, calculateDesignChecksum } from './migrate'

export const useCustomDesign = () => {
  /* ─── STATE ───────────────────────────────────────────────────────── */
  const upper = ref<BoardSection>({
    headerText: 'Selamat & Sukses',
    bodyText: 'Atas Pelantikan Saudara/i\nNama Lengkap Anda',
    headerFontSize: 36,
    bodyFontSize: 20,
    headerFont: 'playfair',
    bodyFont: 'inter',
    headerAlign: 'center',
    bodyAlign: 'center',
    bgColor: '#c0392b',
    headerColor: '#ffd700',
    bodyColor: '#ffffff',
    cornerStyle: 'none',
    opacityPercent: 100
  })

  const lower = ref<BoardSection>({
    headerText: '',
    bodyText: 'Nama Pengirim\nNama Instansi / Perusahaan',
    headerFontSize: 26,
    bodyFontSize: 22,
    headerFont: 'bebas',
    bodyFont: 'inter',
    headerAlign: 'center',
    bodyAlign: 'center',
    bgColor: '#1a3a5c',
    headerColor: '#ffffff',
    bodyColor: '#ffffff',
    cornerStyle: 'none',
    opacityPercent: 100
  })

  const heightRatio = ref(0.58)
  const border = ref<BoardBorder>({ style: 'solid', color: '#f5c842', width: 12, center: true })
  const topCrest = ref<FloralCrest>({ enabled: false, style: 'classic', primary: '#e63946', secondary: '#f1faee', size: 40 })
  const bottomCrest = ref<FloralCrest>({ enabled: false, style: 'classic', primary: '#e63946', secondary: '#f1faee', size: 40 })
  const elements = ref<CanvasElement[]>([])

  const activeTab = ref<ToolTab>('text')
  const activeSection = ref<SectionKey>('upper')
  const selectedId = ref<string | null>(null)
  const physicalSize = ref('medium')
  const showReview = ref(false)
  const showFinalizeChoice = ref(false)
  const showThankYou = ref(false)

  const brushType = ref<BrushType>('flower')
  const brushColor = ref('#e85d75')
  const brushSize = ref(48)
  const brushRotation = ref(0)
  const isBrushMode = computed(() => activeTab.value === 'brush')

  const containerRef = ref<HTMLElement | null>(null)
  const boardRef = ref<HTMLElement | null>(null)
  const boardScale = ref(0.75)

  const snapshotDataUrl = ref<string>('')
  const snapshotLoading = ref(false)

  /* ─── CONSTANTS & HELPERS ─────────────────────────────────────────── */
  const boardW = computed(() => 800)
  const boardH = computed(() => {
    if (physicalSize.value === 'small') return 500
    if (physicalSize.value === 'large') return 600
    return 576
  })

  const getFont = (id: string) => FONTS.find(f => f.id === id)?.family ?? "'Inter', sans-serif"
  const sec = computed(() => activeSection.value === 'upper' ? upper.value : lower.value)
  const upperH = computed(() => Math.round(heightRatio.value * boardH.value))
  const lowerH = computed(() => boardH.value - upperH.value)

  const boardBorderStyle = computed((): Record<string, string> => {
    const { style, color, width } = border.value
    if (style === 'none' || width === 0) return {}
    if (style === 'ornate') return { border: `${width}px solid ${color}`, outline: `${Math.max(2, Math.round(width * 0.35))}px solid ${color}`, outlineOffset: '5px' }
    return { border: `${width}px ${style} ${color}` }
  })

  const centerBorderStyle = computed((): Record<string, string> => {
    if (!border.value.center) return {}
    const { style, color, width } = border.value
    if (style === 'none' || width === 0) return {}
    return {
      width: '100%',
      borderTop: `${width}px ${style === 'ornate' ? 'solid' : style} ${color}`,
      ...(style === 'ornate' ? {
        boxShadow: `0 -${Math.max(2, Math.round(width * 0.35))}px 0 ${color}, 0 ${Math.max(2, Math.round(width * 0.35))}px 0 ${color}`
      } : {})
    }
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

  const floralSec = computed(() => activeSection.value === 'upper' ? topCrest.value : bottomCrest.value)

  const selectedEl = computed(() => elements.value.find(e => e.id === selectedId.value) ?? null)
  const selectedImg = computed(() => (selectedEl.value?.type === 'image' ? selectedEl.value : null) as CanvasImage | null)
  const selectedBrush = computed(() => (selectedEl.value?.type === 'brush' ? selectedEl.value : null) as BrushStroke | null)
  const imgElements = computed(() => elements.value.filter(e => e.type === 'image') as CanvasImage[])
  const brushElements = computed(() => elements.value.filter(e => e.type === 'brush') as BrushStroke[])

  /* ─── LIVE PRICE BREAKDOWN FORMULA ───────────────────────────────── */
  const baseSizePrice = computed(() => SIZES.find(s => s.id === physicalSize.value)?.price ?? 200_000)
  const brushFee = computed(() => brushElements.value.length * 2000)

  const uniqueColors = computed(() => {
    const set = new Set<string>()
    const add = (c?: string) => {
      if (c) set.add(normalizeHexColor(c, '#FFFFFF'))
    }
    add(upper.value.bgColor)
    add(lower.value.bgColor)
    add(upper.value.headerColor)
    add(upper.value.bodyColor)
    add(lower.value.headerColor)
    add(lower.value.bodyColor)
    if (border.value.style !== 'none' && border.value.width > 0) {
      add(border.value.color)
    }
    if (topCrest.value.enabled) {
      add(topCrest.value.primary)
      add(topCrest.value.secondary)
    }
    if (bottomCrest.value.enabled) {
      add(bottomCrest.value.primary)
      add(bottomCrest.value.secondary)
    }
    brushElements.value.forEach(b => add(b.color))
    return Array.from(set)
  })

  const colorFee = computed(() => Math.max(0, uniqueColors.value.length - 3) * 10_000)

  const borderFee = computed(() => {
    if (border.value.style !== 'none' && border.value.width > 0) {
      if (['double', 'groove', 'ridge', 'ornate'].includes(border.value.style)) {
        return 15_000
      }
    }
    return 0
  })

  const getCrestFee = (crest: FloralCrest) => {
    if (!crest.enabled) return 0
    if (crest.style === 'grand') return 45_000
    if (crest.style === 'modern') return 30_000
    return 25_000
  }

  const accessoriesFee = computed(() => borderFee.value + getCrestFee(topCrest.value) + getCrestFee(bottomCrest.value))
  const mediaFee = computed(() => imgElements.value.length * 20_000)
  const totalPrice = computed(() => baseSizePrice.value + brushFee.value + colorFee.value + accessoriesFee.value + mediaFee.value)

  /* ─── WATCHERS ────────────────────────────────────────────────────── */
  watch(selectedBrush, (br) => {
    if (br) {
      brushSize.value = br.size
      brushRotation.value = br.rotation
      brushColor.value = br.color
      brushType.value = br.brushType
    }
  })
  watch(brushSize, (v) => { if (selectedBrush.value) selectedBrush.value.size = v })
  watch(brushRotation, (v) => { if (selectedBrush.value) selectedBrush.value.rotation = v })
  watch(brushColor, (v) => { if (selectedBrush.value) selectedBrush.value.color = v })

  /* ─── SCALE & ACTIONS ─────────────────────────────────────────────── */
  const updateScale = () => {
    const el = containerRef.value
    if (!el) return
    const pad = 64
    boardScale.value = Math.max(0.25, Math.min((el.offsetWidth - pad) / boardW.value, (el.offsetHeight - pad) / boardH.value, 1.1))
  }

  const randomizeDesign = () => {
    const rand = <T>(arr: T[]): T => arr[Math.floor(Math.random() * arr.length)]!
    upper.value.bgColor = rand(BG_PRESETS)
    lower.value.bgColor = rand(BG_PRESETS)
    upper.value.headerFont = rand(FONTS).id
    upper.value.bodyFont = rand(FONTS).id
    lower.value.headerFont = rand(FONTS).id
    lower.value.bodyFont = rand(FONTS).id
    border.value.style = rand(BORDER_STYLES).id
    border.value.color = rand(BORDER_COLORS)
    border.value.width = Math.floor(Math.random() * 16) + 4
    border.value.center = Math.random() > 0.5
    upper.value.cornerStyle = rand(CORNERS).id
    lower.value.cornerStyle = rand(CORNERS).id

    topCrest.value.enabled = Math.random() > 0.5
    topCrest.value.style = rand(['classic', 'modern', 'grand'] as FloralStyle[])
    topCrest.value.primary = rand(BG_PRESETS)
    topCrest.value.secondary = rand(BG_PRESETS)

    bottomCrest.value.enabled = Math.random() > 0.5
    bottomCrest.value.style = rand(['classic', 'modern', 'grand'] as FloralStyle[])
    bottomCrest.value.primary = rand(BG_PRESETS)
    bottomCrest.value.secondary = rand(BG_PRESETS)

    boardScale.value = boardScale.value * 0.95
    setTimeout(updateScale, 150)
  }

  watch(physicalSize, () => updateScale())

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

  /* ─── IMAGE & BRUSH HELPERS ───────────────────────────────────────── */
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

  let _suppressBrushPlace = false
  const handleBrushMousedown = (e: MouseEvent, id: string) => {
    selectedId.value = id
    _suppressBrushPlace = true
    startDragEl(e, id)
  }

  const handleBoardClick = (e: MouseEvent) => {
    if (_suppressBrushPlace) { _suppressBrushPlace = false; return }
    if (isBrushMode.value) {
      const r = boardRef.value!.getBoundingClientRect()
      const stroke: BrushStroke = {
        id: 'br-' + Date.now().toString(36) + Math.random().toString(36).slice(2, 5),
        type: 'brush', brushType: brushType.value,
        x: ((e.clientX - r.left) / r.width) * 100,
        y: ((e.clientY - r.top) / r.height) * 100,
        size: brushSize.value, color: brushColor.value, rotation: brushRotation.value,
      }
      elements.value.push(stroke)
      selectedId.value = stroke.id
    } else {
      selectedId.value = null
    }
  }

  /* ─── DRAG LOGIC ─────────────────────────────────────────────────── */
  let _draggingEl = false, _draggingDiv = false, _dragElId = ''
  let _rect = { left: 0, top: 0, width: 0, height: 0 }
  let _dragBX = 0, _dragBY = 0, _dragElX0 = 0, _dragElY0 = 0
  let _divStartY = 0, _divStartR = 0

  const startDragEl = (e: MouseEvent, id: string) => {
    e.stopPropagation(); e.preventDefault()
    bringToFront(id)
    const board = boardRef.value; if (!board) return
    _rect = board.getBoundingClientRect(); _draggingEl = true; _dragElId = id
    _dragBX = (e.clientX - _rect.left) / _rect.width * 100
    _dragBY = (e.clientY - _rect.top) / _rect.height * 100
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
        el.x = Math.max(-5, Math.min(_dragElX0 + ((e.clientX - _rect.left) / _rect.width * 100 - _dragBX), 95))
        el.y = Math.max(-5, Math.min(_dragElY0 + ((e.clientY - _rect.top) / _rect.height * 100 - _dragBY), 95))
      }
    }
    if (_draggingDiv) {
      heightRatio.value = Math.max(0.25, Math.min(_divStartR + (e.clientY - _divStartY) / _rect.height, 0.75))
    }
  }

  const onMouseUp = () => { _draggingEl = false; _draggingDiv = false }

  /* ─── SNAPSHOT & PAYLOAD (v3.0.0) ─────────────────────────────────── */
  const generateBoardSnapshot = (): Promise<string> => {
    return new Promise((resolve) => {
      const W = 800, H = boardH.value
      const uH = Math.round(heightRatio.value * H)
      const lH = H - uH
      const bw = border.value.style !== 'none' ? border.value.width : 0

      const canvas = document.createElement('canvas')
      canvas.width = W
      canvas.height = H
      const ctx = canvas.getContext('2d')!

      // Upper section
      ctx.fillStyle = upper.value.bgColor
      ctx.fillRect(0, 0, W, uH)

      if (upper.value.headerText) {
        const fam = getFont(upper.value.headerFont).replace(/'/g, '')
        const sz = Math.round(upper.value.headerFontSize * 0.9)
        ctx.font = `700 ${sz}px ${fam}, sans-serif`
        ctx.fillStyle = upper.value.headerColor
        ctx.textAlign = upper.value.headerAlign as CanvasTextAlign
        ctx.textBaseline = 'middle'
        const tX = upper.value.headerAlign === 'left' ? 32 : upper.value.headerAlign === 'right' ? W - 32 : W / 2
        ctx.fillText(upper.value.headerText, tX, uH * 0.38, W - 64)
      }
      if (upper.value.bodyText) {
        const fam = getFont(upper.value.bodyFont).replace(/'/g, '')
        const sz = Math.round(upper.value.bodyFontSize * 0.85)
        ctx.font = `${sz}px ${fam}, sans-serif`
        ctx.fillStyle = upper.value.bodyColor
        ctx.textAlign = upper.value.bodyAlign as CanvasTextAlign
        ctx.textBaseline = 'middle'
        const lines = (upper.value.bodyText || '').split('\n')
        const tX = upper.value.bodyAlign === 'left' ? 32 : upper.value.bodyAlign === 'right' ? W - 32 : W / 2
        lines.forEach((line, i) => {
          ctx.fillText(line, tX, uH * 0.68 + i * (sz + 8), W - 64)
        })
      }

      // Lower section
      ctx.fillStyle = lower.value.bgColor
      ctx.fillRect(0, uH, W, lH)

      if (lower.value.headerText) {
        const fam = getFont(lower.value.headerFont).replace(/'/g, '')
        const sz = Math.round(lower.value.headerFontSize * 0.9)
        ctx.font = `700 ${sz}px ${fam}, sans-serif`
        ctx.fillStyle = lower.value.headerColor
        ctx.textAlign = lower.value.headerAlign as CanvasTextAlign
        ctx.textBaseline = 'middle'
        const tX = lower.value.headerAlign === 'left' ? 32 : lower.value.headerAlign === 'right' ? W - 32 : W / 2
        ctx.fillText(lower.value.headerText, tX, uH + lH * 0.32, W - 64)
      }
      if (lower.value.bodyText) {
        const fam = getFont(lower.value.bodyFont).replace(/'/g, '')
        const sz = Math.round(lower.value.bodyFontSize * 0.85)
        ctx.font = `${sz}px ${fam}, sans-serif`
        ctx.fillStyle = lower.value.bodyColor
        ctx.textAlign = lower.value.bodyAlign as CanvasTextAlign
        ctx.textBaseline = 'middle'
        const lines = lower.value.bodyText.split('\n')
        const startY = lower.value.headerText ? uH + lH * 0.58 : uH + lH * 0.4
        const tX = lower.value.bodyAlign === 'left' ? 32 : lower.value.bodyAlign === 'right' ? W - 32 : W / 2
        lines.forEach((line, i) => {
          ctx.fillText(line, tX, startY + i * (sz + 8), W - 64)
        })
      }

      // Divider
      if (border.value.center && border.value.style !== 'none' && bw > 0) {
        ctx.strokeStyle = border.value.color
        ctx.lineWidth = bw
        ctx.beginPath()
        ctx.moveTo(0, uH); ctx.lineTo(W, uH)
        ctx.stroke()
      }

      // Border
      if (border.value.style !== 'none' && bw > 0) {
        ctx.strokeStyle = border.value.color
        ctx.lineWidth = bw
        ctx.strokeRect(bw / 2, bw / 2, W - bw, H - bw)
        if (border.value.style === 'ornate') {
          const gap = Math.max(3, Math.round(bw * 0.4))
          ctx.lineWidth = Math.max(1, Math.round(bw * 0.35))
          ctx.strokeRect(bw + gap, bw + gap, W - (bw + gap) * 2, H - (bw + gap) * 2)
        }
      }

      // Images
      const imgEls = elements.value.filter(e => e.type === 'image') as CanvasImage[]
      const draws = imgEls.map(el => new Promise<void>(res => {
        const img = new Image()
        img.crossOrigin = 'anonymous'
        img.onload = () => {
          const px = (el.x / 100) * W
          const py = (el.y / 100) * H
          const pw = (el.width / 100) * W
          ctx.save()
          ctx.beginPath()
          if (el.frame === 'circle') {
            ctx.arc(px, py, pw / 2, 0, Math.PI * 2)
          } else {
            ctx.rect(px - pw / 2, py - pw / 2, pw, pw)
          }
          ctx.clip()
          ctx.drawImage(img, px - pw / 2, py - pw / 2, pw, pw)
          ctx.restore()
          res()
        }
        img.onerror = () => res()
        img.src = el.src
      }))

      Promise.all(draws).then(() => {
        const brushEls = elements.value.filter(e => e.type === 'brush') as BrushStroke[]
        brushEls.forEach(el => {
          const cx = (el.x / 100) * W
          const cy = (el.y / 100) * H
          const r = el.size / 2
          ctx.save()
          ctx.translate(cx, cy)
          ctx.rotate((el.rotation * Math.PI) / 180)
          ctx.fillStyle = el.color
          if (el.brushType === 'flower') {
            for (let i = 0; i < 5; i++) {
              ctx.save()
              ctx.rotate((i * 72 * Math.PI) / 180)
              ctx.beginPath()
              ctx.ellipse(0, -r * 0.6, r * 0.28, r * 0.5, 0, 0, Math.PI * 2)
              ctx.fill()
              ctx.restore()
            }
            ctx.fillStyle = '#ffd700'
            ctx.beginPath()
            ctx.arc(0, 0, r * 0.28, 0, Math.PI * 2)
            ctx.fill()
          } else {
            ctx.beginPath()
            ctx.arc(0, 0, r * 0.8, 0, Math.PI * 2)
            ctx.fill()
            ctx.fillStyle = 'rgba(0,0,0,0.2)'
            ctx.beginPath()
            ctx.arc(r * 0.15, -r * 0.1, r * 0.5, 0, Math.PI * 2)
            ctx.fill()
            ctx.fillStyle = '#ffe066'
            ctx.globalAlpha = 0.7
            ctx.beginPath()
            ctx.arc(0, r * 0.1, r * 0.18, 0, Math.PI * 2)
            ctx.fill()
            ctx.globalAlpha = 1
          }
          ctx.restore()
        })

        resolve(canvas.toDataURL('image/png', 0.92))
      })
    })
  }

  const buildCustomDesignPayload = (previewBase64: string): CustomDesignPayloadV3 => {
    const payload: CustomDesignPayloadV3 = {
      metadata: {
        version: '3.0.0',
        editorVersion: '3.0.0',
        platform: 'web',
        locale: 'id-ID',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        checksum: '',
        featureFlags: { standaloneModule: true, v3Support: true }
      },
      layout: {
        physicalSizeId: physicalSize.value as any,
        aspectRatioId: 'portrait-3-4',
        upperHeightRatio: heightRatio.value,
        border: {
          style: border.value.style as any,
          colorHex: normalizeHexColor(border.value.color, '#F5C842'),
          widthPx: border.value.width,
          showCenterDivider: border.value.center
        }
      },
      sections: {
        upper: {
          bgColorHex: normalizeHexColor(upper.value.bgColor, '#C0392B'),
          cornerStyle: upper.value.cornerStyle as any,
          opacityPercent: upper.value.opacityPercent ?? 100,
          header: {
            text: upper.value.headerText || null,
            fontId: upper.value.headerFont as any,
            fontSizePx: upper.value.headerFontSize,
            fontColorHex: normalizeHexColor(upper.value.headerColor, '#FFD700'),
            alignment: upper.value.headerAlign as any
          },
          body: {
            text: upper.value.bodyText || null,
            fontId: upper.value.bodyFont as any,
            fontSizePx: upper.value.bodyFontSize,
            fontColorHex: normalizeHexColor(upper.value.bodyColor, '#FFFFFF'),
            alignment: upper.value.bodyAlign as any
          }
        },
        lower: {
          bgColorHex: normalizeHexColor(lower.value.bgColor, '#1A3A5C'),
          cornerStyle: lower.value.cornerStyle as any,
          opacityPercent: lower.value.opacityPercent ?? 100,
          header: {
            text: lower.value.headerText || null,
            fontId: lower.value.headerFont as any,
            fontSizePx: lower.value.headerFontSize,
            fontColorHex: normalizeHexColor(lower.value.headerColor, '#FFFFFF'),
            alignment: lower.value.headerAlign as any
          },
          body: {
            text: lower.value.bodyText || null,
            fontId: lower.value.bodyFont as any,
            fontSizePx: lower.value.bodyFontSize,
            fontColorHex: normalizeHexColor(lower.value.bodyColor, '#FFFFFF'),
            alignment: lower.value.bodyAlign as any
          }
        }
      },
      decorations: {
        topCrest: {
          visible: topCrest.value.enabled,
          variantId: topCrest.value.style as any,
          primaryColorHex: normalizeHexColor(topCrest.value.primary, '#E63946'),
          secondaryColorHex: normalizeHexColor(topCrest.value.secondary, '#F1FAEE'),
          scalePercent: topCrest.value.size
        },
        bottomCrest: {
          visible: bottomCrest.value.enabled,
          variantId: bottomCrest.value.style as any,
          primaryColorHex: normalizeHexColor(bottomCrest.value.primary, '#E63946'),
          secondaryColorHex: normalizeHexColor(bottomCrest.value.secondary, '#F1FAEE'),
          scalePercent: bottomCrest.value.size
        },
        watermark: { enabled: false, text: 'Chia Florist', opacityPercent: 20 }
      },
      elements: elements.value.map((el, index) => {
        if (el.type === 'image') {
          const img = el as CanvasImage
          return {
            id: img.id,
            type: 'image' as const,
            src: img.src,
            frameStyle: img.frame as any,
            crop: { xPercent: img.cropX, yPercent: img.cropY, zoom: img.zoom },
            transform: {
              xPercent: img.x,
              yPercent: img.y,
              scalePercent: img.width,
              rotationDeg: 0,
              zIndex: img.zIndex ?? (index + 10)
            }
          }
        }
        const br = el as BrushStroke
        return {
          id: br.id,
          type: 'brush' as const,
          brushType: br.brushType as any,
          colorHex: normalizeHexColor(br.color, '#E85D75'),
          transform: {
            xPercent: br.x,
            yPercent: br.y,
            scalePercent: br.size,
            rotationDeg: br.rotation,
            zIndex: br.zIndex ?? (index + 10)
          }
        }
      }),
      assets: {
        previewBase64: previewBase64 || null,
        previewAssetId: null,
        previewUrl: null,
        bucketPath: null,
        storageProvider: 'supabase'
      }
    }

    payload.metadata.checksum = calculateDesignChecksum(payload)
    return payload
  }

  /* ─── DRAFT PERSISTENCE ──────────────────────────────────────────── */
  const saveDraft = () => {
    if (!import.meta.client) return
    try {
      const draft = {
        upper: upper.value,
        lower: lower.value,
        border: border.value,
        topCrest: topCrest.value,
        bottomCrest: bottomCrest.value,
        elements: elements.value,
        physicalSize: physicalSize.value,
        heightRatio: heightRatio.value
      }
      localStorage.setItem(DEFAULT_DRAFT_KEY, JSON.stringify(draft))
    } catch (err) {
      console.warn('Failed to auto-save custom board draft:', err)
    }
  }

  const loadDraft = () => {
    if (!import.meta.client) return
    try {
      const raw = localStorage.getItem(DEFAULT_DRAFT_KEY)
      if (!raw) return
      const draft = JSON.parse(raw)
      if (draft.upper) upper.value = draft.upper
      if (draft.lower) lower.value = draft.lower
      if (draft.border) border.value = draft.border
      if (draft.topCrest) topCrest.value = draft.topCrest
      if (draft.bottomCrest) bottomCrest.value = draft.bottomCrest
      if (draft.elements) elements.value = draft.elements
      if (draft.physicalSize) physicalSize.value = draft.physicalSize
      if (draft.heightRatio) heightRatio.value = draft.heightRatio
    } catch (err) {
      console.warn('Failed to load custom board draft:', err)
    }
  }

  const clearDraft = () => {
    if (!import.meta.client) return
    localStorage.removeItem(DEFAULT_DRAFT_KEY)
  }

  watch(
    [upper, lower, border, topCrest, bottomCrest, elements, physicalSize, heightRatio],
    () => { saveDraft() },
    { deep: true }
  )

  const resetDesign = () => {
    clearDraft()
    elements.value = []
    selectedId.value = null
    snapshotDataUrl.value = ''
    upper.value = {
      headerText: 'Selamat & Sukses', bodyText: 'Atas Pelantikan Saudara/i\nNama Lengkap Anda',
      headerFontSize: 36, bodyFontSize: 20, headerFont: 'playfair', bodyFont: 'inter',
      headerAlign: 'center', bodyAlign: 'center',
      bgColor: '#c0392b', headerColor: '#ffd700', bodyColor: '#ffffff', cornerStyle: 'none', opacityPercent: 100
    }
    lower.value = {
      headerText: '', bodyText: 'Nama Pengirim\nNama Instansi / Perusahaan',
      headerFontSize: 26, bodyFontSize: 22, headerFont: 'bebas', bodyFont: 'inter',
      headerAlign: 'center', bodyAlign: 'center',
      bgColor: '#1a3a5c', headerColor: '#ffffff', bodyColor: '#ffffff', cornerStyle: 'none', opacityPercent: 100
    }
    border.value = { style: 'solid', color: '#f5c842', width: 12, center: true }
    topCrest.value = { enabled: false, style: 'classic', primary: '#e63946', secondary: '#f1faee', size: 40 }
    bottomCrest.value = { enabled: false, style: 'classic', primary: '#e63946', secondary: '#f1faee', size: 40 }
    physicalSize.value = 'medium'
    heightRatio.value = 0.58
  }

  // Keyboard shortcut listener
  const onKeyDown = (e: KeyboardEvent) => {
    const tag = (document.activeElement as HTMLElement)?.tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA') return
    if (e.key === 'Delete' || e.key === 'Backspace') deleteSelected()
    if (e.key === 'Escape') selectedId.value = null
  }

  // Review modal snapshot trigger
  watch(showReview, async (open) => {
    if (!open) return
    snapshotLoading.value = true
    snapshotDataUrl.value = ''
    await nextTick()
    try {
      snapshotDataUrl.value = await generateBoardSnapshot()
      const updatedDesign = buildCustomDesignPayload(snapshotDataUrl.value)
      console.log('[Chia Florist] Review Modal Open - Generated Custom Design Payload v3.0:\n', JSON.stringify(updatedDesign, null, 2))
    } catch (e) {
      console.warn('Snapshot generation failed:', e)
    } finally {
      snapshotLoading.value = false
    }
  })

  return reactive({
    upper, lower, heightRatio, border, topCrest, bottomCrest, elements,
    activeTab, activeSection, selectedId, physicalSize, showReview, showFinalizeChoice, showThankYou,
    brushType, brushColor, brushSize, brushRotation, isBrushMode,
    containerRef, boardRef, boardScale, snapshotDataUrl, snapshotLoading,
    boardW, boardH, getFont, sec, upperH, lowerH,
    boardBorderStyle, centerBorderStyle, upperCornerStyle, lowerCornerStyle, floralSec,
    selectedEl, selectedImg, selectedBrush, imgElements, brushElements,
    baseSizePrice, brushFee, uniqueColors, colorFee, borderFee, accessoriesFee, mediaFee, totalPrice,
    updateScale, randomizeDesign, bringToFront, deleteSelected,
    handleDrop, handleFileInput, handleBrushMousedown, handleBoardClick,
    startDragEl, startDragDiv, onMouseMove, onMouseUp, onKeyDown,
    generateBoardSnapshot, buildCustomDesignPayload, loadDraft, saveDraft, clearDraft, resetDesign
  })
}
