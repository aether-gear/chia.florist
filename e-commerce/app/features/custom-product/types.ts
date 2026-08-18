// app/features/custom-product/types.ts

export type FontId = 'inter' | 'playfair' | 'dancing' | 'bebas' | 'merriweather' | 'pacifico'
export type CornerStyle = 'none' | 'rounded' | 'cut' | 'ornate' | 'floral'
export type FrameStyle = 'none' | 'square' | 'circle'
export type BrushType = 'flower' | 'rose'
export type BorderStyle = 'none' | 'solid' | 'double' | 'dashed' | 'dotted' | 'groove' | 'ridge' | 'ornate'
export type FloralStyle = 'classic' | 'modern' | 'grand'
export type ToolTab = 'text' | 'image' | 'brush' | 'border' | 'corner' | 'floral'
export type SectionKey = 'upper' | 'lower'

export interface BoardSection {
  headerText: string
  bodyText: string
  headerFontSize: number
  bodyFontSize: number
  headerFont: FontId
  bodyFont: FontId
  headerAlign: 'left' | 'center' | 'right'
  bodyAlign: 'left' | 'center' | 'right'
  bgColor: string
  headerColor: string
  bodyColor: string
  cornerStyle: CornerStyle
  opacityPercent?: number
  // Text section borders (decorative horizontal rule under header / above body)
  headerBorder?: boolean
  headerBorderColor?: string
  headerBorderWidth?: number
  bodyBorder?: boolean
  bodyBorderColor?: string
  bodyBorderWidth?: number
}

export interface CanvasImage {
  id: string
  type: 'image'
  src: string
  frame: FrameStyle
  x: number
  y: number
  width: number
  zoom: number
  cropX: number
  cropY: number
  zIndex?: number
}

export interface BrushStroke {
  id: string
  type: 'brush'
  brushType: BrushType
  x: number
  y: number
  size: number
  color: string
  rotation: number
  zIndex?: number
}

export type CanvasElement = CanvasImage | BrushStroke

export interface BoardBorder {
  style: BorderStyle
  color: string
  width: number
  center: boolean
}

export interface FloralCrest {
  enabled: boolean
  style: FloralStyle
  primary: string
  secondary: string
  size: number
}

export interface WatermarkSpec {
  enabled: boolean
  text: string
  opacityPercent: number
}

/* ─── Standardized Custom Design Payload Schema ───────────────── */

export interface TypographySpec {
  text: string | null
  fontId: FontId
  fontSizePx: number
  fontColorHex: string
  alignment: 'left' | 'center' | 'right'
}

export interface BoardSectionSpec {
  bgColorHex: string
  cornerStyle: CornerStyle
  opacityPercent?: number
  header: TypographySpec
  body: TypographySpec
}

export interface BorderSpec {
  style: BorderStyle
  colorHex: string
  widthPx: number
  showCenterDivider: boolean
}

export interface CrestSpec {
  visible: boolean
  variantId: FloralStyle
  primaryColorHex: string
  secondaryColorHex: string
  scalePercent: number
}

export interface BaseElement {
  id: string
  type: 'image' | 'brush'
  transform: {
    xPercent: number
    yPercent: number
    scalePercent: number
    rotationDeg: number
    zIndex: number
  }
}

export interface ImageElement extends BaseElement {
  type: 'image'
  src: string
  frameStyle: FrameStyle
  crop: { xPercent: number; yPercent: number; zoom: number }
}

export interface BrushElement extends BaseElement {
  type: 'brush'
  brushType: BrushType
  colorHex: string
}

export type DesignElement = ImageElement | BrushElement

export interface CustomDesignPayloadV1 {
  metadata: {
    version: '1.0.0'
    editorVersion: string
    platform: string
    locale: string
    createdAt: string
    updatedAt: string
    checksum: string
  }
  layout: {
    physicalSizeId: 'small' | 'medium' | 'large'
    upperHeightRatio: number
    border: BorderSpec
  }
  sections: {
    upper: BoardSectionSpec
    lower: BoardSectionSpec
  }
  decorations: {
    topCrest: CrestSpec
    bottomCrest: CrestSpec
  }
  elements: DesignElement[]
  assets: {
    previewBase64: string | null
    previewAssetId: string | null
    previewUrl: string | null
    bucketPath: string | null
    storageProvider: 'supabase' | 's3' | 'local' | null
  }
}

export interface CustomDesignPayloadV3 {
  metadata: {
    version: '3.0.0'
    editorVersion: string
    platform: string
    locale: string
    createdAt: string
    updatedAt: string
    checksum: string
    featureFlags?: Record<string, boolean>
  }
  layout: {
    physicalSizeId: 'small' | 'medium' | 'large'
    aspectRatioId?: string
    upperHeightRatio: number
    border: BorderSpec
  }
  sections: {
    upper: BoardSectionSpec
    lower: BoardSectionSpec
  }
  decorations: {
    topCrest: CrestSpec
    bottomCrest: CrestSpec
    watermark?: WatermarkSpec
  }
  elements: DesignElement[]
  assets: {
    previewBase64: string | null
    previewAssetId: string | null
    previewUrl: string | null
    bucketPath: string | null
    storageProvider: 'supabase' | 's3' | 'local' | null
  }
}

export type CustomDesignPayload = CustomDesignPayloadV1 | CustomDesignPayloadV3
