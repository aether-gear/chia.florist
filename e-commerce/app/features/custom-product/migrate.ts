// app/features/custom-product/migrate.ts

import type { CustomDesignPayloadV3 } from './types'
import { CUSTOM_PRODUCT_SCHEMA_VERSION, CUSTOM_PRODUCT_ENGINE_VERSION } from './constants'

// Validate 6-digit hex color format (#RRGGBB)
export const normalizeHexColor = (color: string | undefined, fallback = '#FFFFFF'): string => {
  if (!color) return fallback
  const trimmed = color.trim().toUpperCase()
  if (/^#[0-9A-F]{6}$/.test(trimmed)) return trimmed
  if (/^#[0-9A-F]{3}$/.test(trimmed)) {
    return `#${trimmed[1]}${trimmed[1]}${trimmed[2]}${trimmed[2]}${trimmed[3]}${trimmed[3]}`
  }
  return fallback
}

// Calculate SHA-256 / FNV-1a checksum for design deduplication
export const calculateDesignChecksum = (payload: any): string => {
  const copy = { ...payload }
  if (copy.metadata) copy.metadata = { ...copy.metadata, checksum: '' }
  const jsonStr = JSON.stringify(copy)
  let hash = 0x811c9dc5
  for (let i = 0; i < jsonStr.length; i++) {
    hash ^= jsonStr.charCodeAt(i)
    hash += (hash << 1) + (hash << 4) + (hash << 7) + (hash << 8) + (hash << 24)
  }
  return (hash >>> 0).toString(16).padStart(8, '0')
}

// Migration utility to upgrade legacy flat or v1.0.0 payloads to v3.0.0
export const migrateToV3 = (raw: any): CustomDesignPayloadV3 => {
  if (!raw) {
    throw new Error('Empty custom design payload')
  }

  // Already v3.0.0 structured payload
  if (raw.metadata && raw.metadata.version === CUSTOM_PRODUCT_SCHEMA_VERSION && raw.layout && raw.sections) {
    const payload = raw as CustomDesignPayloadV3
    if (!payload.metadata.checksum) {
      payload.metadata.checksum = calculateDesignChecksum(payload)
    }
    return payload
  }

  // Handle v1.0.0 or legacy flat format
  const metadata = raw.metadata || {}
  const layout = raw.layout || {}
  const sections = raw.sections || {}
  const decorations = raw.decorations || {}
  const rawElements = Array.isArray(raw.elements) ? raw.elements : []
  const assets = raw.assets || {}

  // Section defaults
  const upperSec = sections.upper || raw.upper || {}
  const lowerSec = sections.lower || raw.lower || {}

  const upperHeader = upperSec.header || {}
  const upperBody = upperSec.body || {}
  const lowerHeader = lowerSec.header || {}
  const lowerBody = lowerSec.body || {}

  const border = layout.border || raw.border || {}
  const topCrest = decorations.topCrest || raw.topCrest || {}
  const bottomCrest = decorations.bottomCrest || raw.bottomCrest || {}

  const payload: CustomDesignPayloadV3 = {
    metadata: {
      version: CUSTOM_PRODUCT_SCHEMA_VERSION,
      editorVersion: metadata.editorVersion || CUSTOM_PRODUCT_ENGINE_VERSION,
      platform: metadata.platform || 'web',
      locale: metadata.locale || 'id-ID',
      createdAt: metadata.createdAt || raw.generatedAt || new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      checksum: '',
      featureFlags: metadata.featureFlags || { standaloneModule: true, v3Support: true }
    },
    layout: {
      physicalSizeId: (layout.physicalSizeId || raw.physicalSizeId || 'medium') as any,
      aspectRatioId: layout.aspectRatioId || 'portrait-3-4',
      upperHeightRatio: Math.max(0.25, Math.min(Number(layout.upperHeightRatio ?? raw.heightRatio ?? 0.58), 0.75)),
      border: {
        style: (border.style || 'solid') as any,
        colorHex: normalizeHexColor(border.colorHex || border.color, '#F5C842'),
        widthPx: Number(border.widthPx ?? border.width ?? 12),
        showCenterDivider: Boolean(border.showCenterDivider ?? border.center ?? true)
      }
    },
    sections: {
      upper: {
        bgColorHex: normalizeHexColor(upperSec.bgColorHex || upperSec.bgColor, '#C0392B'),
        cornerStyle: (upperSec.cornerStyle || 'none') as any,
        opacityPercent: upperSec.opacityPercent ?? 100,
        header: {
          text: upperHeader.text ?? upperSec.headerText ?? null,
          fontId: (upperHeader.fontId || upperSec.headerFont || 'playfair') as any,
          fontSizePx: Number(upperHeader.fontSizePx ?? upperSec.headerFontSize ?? 36),
          fontColorHex: normalizeHexColor(upperHeader.fontColorHex || upperSec.headerColor, '#FFD700'),
          alignment: (upperHeader.alignment || upperSec.headerAlign || 'center') as any
        },
        body: {
          text: upperBody.text ?? upperSec.bodyText ?? null,
          fontId: (upperBody.fontId || upperSec.bodyFont || 'inter') as any,
          fontSizePx: Number(upperBody.fontSizePx ?? upperSec.bodyFontSize ?? 20),
          fontColorHex: normalizeHexColor(upperBody.fontColorHex || upperSec.bodyColor, '#FFFFFF'),
          alignment: (upperBody.alignment || upperSec.bodyAlign || 'center') as any
        }
      },
      lower: {
        bgColorHex: normalizeHexColor(lowerSec.bgColorHex || lowerSec.bgColor, '#1A3A5C'),
        cornerStyle: (lowerSec.cornerStyle || 'none') as any,
        opacityPercent: lowerSec.opacityPercent ?? 100,
        header: {
          text: lowerHeader.text ?? lowerSec.headerText ?? null,
          fontId: (lowerHeader.fontId || lowerSec.headerFont || 'bebas') as any,
          fontSizePx: Number(lowerHeader.fontSizePx ?? lowerSec.headerFontSize ?? 26),
          fontColorHex: normalizeHexColor(lowerHeader.fontColorHex || lowerSec.headerColor, '#FFFFFF'),
          alignment: (lowerHeader.alignment || lowerSec.headerAlign || 'center') as any
        },
        body: {
          text: lowerBody.text ?? lowerSec.bodyText ?? null,
          fontId: (lowerBody.fontId || lowerSec.bodyFont || 'inter') as any,
          fontSizePx: Number(lowerBody.fontSizePx ?? lowerSec.bodyFontSize ?? 22),
          fontColorHex: normalizeHexColor(lowerBody.fontColorHex || lowerSec.bodyColor, '#FFFFFF'),
          alignment: (lowerBody.alignment || lowerSec.bodyAlign || 'center') as any
        }
      }
    },
    decorations: {
      topCrest: {
        visible: Boolean(topCrest.visible ?? topCrest.enabled),
        variantId: (topCrest.variantId || topCrest.style || 'classic') as any,
        primaryColorHex: normalizeHexColor(topCrest.primaryColorHex || topCrest.primary, '#E63946'),
        secondaryColorHex: normalizeHexColor(topCrest.secondaryColorHex || topCrest.secondary, '#F1FAEE'),
        scalePercent: Number(topCrest.scalePercent ?? topCrest.size ?? 40)
      },
      bottomCrest: {
        visible: Boolean(bottomCrest.visible ?? bottomCrest.enabled),
        variantId: (bottomCrest.variantId || bottomCrest.style || 'classic') as any,
        primaryColorHex: normalizeHexColor(bottomCrest.primaryColorHex || bottomCrest.primary, '#E63946'),
        secondaryColorHex: normalizeHexColor(bottomCrest.secondaryColorHex || bottomCrest.secondary, '#F1FAEE'),
        scalePercent: Number(bottomCrest.scalePercent ?? bottomCrest.size ?? 40)
      },
      watermark: decorations.watermark || { enabled: false, text: 'Chia Florist', opacityPercent: 20 }
    },
    elements: rawElements.map((el: any, index: number) => {
      const transform = el.transform || {}
      const zIndex = transform.zIndex ?? (index + 10)
      if (el.type === 'image') {
        const crop = el.crop || {}
        return {
          id: el.id || `img-${index}-${Date.now()}`,
          type: 'image' as const,
          src: el.src || '',
          frameStyle: (el.frameStyle || el.frame || 'square') as any,
          crop: { xPercent: crop.xPercent ?? el.cropX ?? 50, yPercent: crop.yPercent ?? el.cropY ?? 50, zoom: crop.zoom ?? el.zoom ?? 1 },
          transform: {
            xPercent: transform.xPercent ?? el.x ?? 15,
            yPercent: transform.yPercent ?? el.y ?? 15,
            scalePercent: transform.scalePercent ?? el.width ?? 22,
            rotationDeg: transform.rotationDeg ?? el.rotation ?? 0,
            zIndex
          }
        }
      }
      return {
        id: el.id || `br-${index}-${Date.now()}`,
        type: 'brush' as const,
        brushType: (el.brushType || 'flower') as any,
        colorHex: normalizeHexColor(el.colorHex || el.color, '#E85D75'),
        transform: {
          xPercent: transform.xPercent ?? el.x ?? 50,
          yPercent: transform.yPercent ?? el.y ?? 50,
          scalePercent: transform.scalePercent ?? el.size ?? 48,
          rotationDeg: transform.rotationDeg ?? el.rotation ?? 0,
          zIndex
        }
      }
    }),
    assets: {
      previewBase64: assets.previewBase64 || raw.previewBase64 || null,
      previewAssetId: assets.previewAssetId || raw.previewAssetId || null,
      previewUrl: assets.previewUrl || raw.previewUrl || null,
      bucketPath: assets.bucketPath || raw.bucketPath || null,
      storageProvider: assets.storageProvider || raw.storageProvider || 'supabase'
    }
  }

  payload.metadata.checksum = calculateDesignChecksum(payload)
  return payload
}
