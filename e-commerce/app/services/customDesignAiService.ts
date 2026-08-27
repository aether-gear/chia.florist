import { bootstrapConfig } from '~/utils/bootstrap'
import type { CustomDesignPayloadV3 } from '~/features/custom-product/types'

export interface GenerateCustomDesignAIRequest {
  prompt: string
  occasion?: string
  preferred_palette?: string
  recipient?: string
  sender?: string
  physical_size_id?: 'small' | 'medium' | 'large'
}

export const customDesignAiService = {
  async generateDesign(payload: GenerateCustomDesignAIRequest): Promise<CustomDesignPayloadV3> {
    return bootstrapConfig.fetchApi<CustomDesignPayloadV3>('/custom-products/ai/generate', {
      method: 'POST',
      body: payload
    })
  }
}
