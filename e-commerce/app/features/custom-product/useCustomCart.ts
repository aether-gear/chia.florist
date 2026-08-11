// app/features/custom-product/useCustomCart.ts

import { ref } from 'vue'
import { useCart } from '~/composables/useCart'
import { supabaseService } from '~/services/supabaseService'
import type { CustomDesignPayloadV3 } from './types'
import { SIZES } from './constants'

export const useCustomCart = () => {
  const { addToCart } = useCart()
  const isAdding = ref(false)

  const addCustomDesignToCart = async (
    design: CustomDesignPayloadV3,
    snapshotUrl: string,
    totalPrice: number,
    physicalSize: string,
    headerText?: string
  ) => {
    isAdding.value = true

    if (snapshotUrl && snapshotUrl.startsWith('data:image/')) {
      try {
        const response = await fetch(snapshotUrl)
        const blob = await response.blob()
        const uploadRes = await supabaseService.uploadCustomPreview(blob)
        if (uploadRes) {
          design.assets.previewUrl = uploadRes.publicUrl
          design.assets.bucketPath = uploadRes.bucketPath
          design.assets.storageProvider = 'supabase'
        }
      } catch (err) {
        console.warn('[Chia Florist] Could not upload custom preview to Supabase:', err)
      }
    }

    const previewUrl = design.assets.previewUrl || snapshotUrl || '/images/custom-preview.png'
    console.log('[Chia Florist] Finalized Custom Design Payload v3.0 for Cart/Order:\n', JSON.stringify(design, null, 2))

    const itemId = 'custom-' + Date.now()
    await addToCart({
      id: itemId,
      name: `Custom Board — ${headerText || 'My Design'}`,
      price: totalPrice,
      image: previewUrl,
      size: SIZES.find(s => s.id === physicalSize)?.label ?? '',
      color: design.sections.upper.bgColorHex,
      shopId: '99ef0062-1040-4574-a4be-0123abce5670',
      isCustom: true,
      itemType: 'custom',
      customDesign: design
    }, 1)

    isAdding.value = false
    return itemId
  }

  return {
    isAdding,
    addCustomDesignToCart
  }
}
