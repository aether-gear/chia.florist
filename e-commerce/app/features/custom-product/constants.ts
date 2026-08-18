// app/features/custom-product/constants.ts

import type { FontId, CornerStyle, BorderStyle, ToolTab } from './types'

export const FONTS: { id: FontId; label: string; family: string }[] = [
  { id: 'inter', label: 'Inter', family: "'Inter', sans-serif" },
  { id: 'playfair', label: 'Playfair', family: "'Playfair Display', serif" },
  { id: 'dancing', label: 'Dancing', family: "'Dancing Script', cursive" },
  { id: 'bebas', label: 'Bebas', family: "'Bebas Neue', sans-serif" },
  { id: 'merriweather', label: 'Merriweather', family: "'Merriweather', serif" },
  { id: 'pacifico', label: 'Pacifico', family: "'Pacifico', cursive" },
]

export const CORNERS: { id: CornerStyle; label: string }[] = [
  { id: 'none', label: 'None' },
  { id: 'rounded', label: 'Rounded' },
  { id: 'cut', label: 'Cut' },
  { id: 'ornate', label: 'Ornate' },
  { id: 'floral', label: 'Floral' },
]

export const BORDER_STYLES: { id: BorderStyle; label: string }[] = [
  { id: 'none', label: 'None' },
  { id: 'solid', label: 'Solid' },
  { id: 'double', label: 'Double' },
  { id: 'dashed', label: 'Dashed' },
  { id: 'dotted', label: 'Dotted' },
  { id: 'groove', label: 'Groove' },
  { id: 'ridge', label: 'Ridge' },
  { id: 'ornate', label: 'Ornate' },
]

export const SIZES = [
  { id: 'small', label: '1.5 × 2.0m', price: 150_000, desc: 'Compact' },
  { id: 'medium', label: '1.8 × 2.5m', price: 200_000, desc: 'Standard', recommended: true },
  { id: 'large', label: '2.0 × 3.0m', price: 280_000, desc: 'Grand' },
]

export const BRUSH_COLORS = ['#e85d75', '#f4845f', '#f9c74f', '#90be6d', '#4cc9f0', '#c77dff', '#ffffff', '#222222']
export const BORDER_COLORS = ['#f5c842', '#e63946', '#2a9d8f', '#264653', '#e76f51', '#a8dadc', '#f1faee', '#1d3557']
export const BG_PRESETS = ['#c0392b', '#1a3a5c', '#145a32', '#6c3483', '#a04000', '#17202a', '#f0f0f0', '#ffffff']

export const TOOL_TABS: { id: ToolTab; label: string }[] = [
  { id: 'text', label: 'Text' },
  { id: 'image', label: 'Image' },
  { id: 'brush', label: 'Brush' },
  { id: 'border', label: 'Border' },
  { id: 'corner', label: 'Corner' },
  { id: 'floral', label: 'Floral' },
]

export const DEFAULT_DRAFT_KEY = 'chia-florist-custom-draft'
