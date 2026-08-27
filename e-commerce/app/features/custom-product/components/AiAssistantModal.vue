<script setup lang="ts">
import { ref, computed } from 'vue'
import { customDesignAiService } from '~/services/customDesignAiService'
import type { CustomDesignPayloadV3 } from '../types'
import { useGlobalAlert } from '~/composables/useGlobalAlert'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'apply', payload: CustomDesignPayloadV3): void
}>()

const globalAlert = useGlobalAlert()

const promptText = ref('')
const occasion = ref('')
const preferredPalette = ref('')
const recipient = ref('')
const sender = ref('')
const physicalSizeId = ref<'small' | 'medium' | 'large'>('medium')
const showAdvanced = ref(false)

const isGenerating = ref(false)
const generatedPayload = ref<CustomDesignPayloadV3 | null>(null)
const errorMsg = ref<string | null>(null)

const presetPrompts = [
  {
    icon: '💍',
    label: 'Romantic Wedding',
    prompt: 'Papan bunga pernikahan nuansa blush pink & emas elegan untuk Dimas & Sarah dari PT Maju Jaya',
    occasion: 'Pernikahan',
    palette: 'Pastel Gold & Blush Pink'
  },
  {
    icon: '🎉',
    label: 'Grand Opening',
    prompt: 'Papan ucapan Selamat & Sukses atas pembukaan cabang baru Resto Nusantara dari Kopi Kenangan',
    occasion: 'Grand Opening',
    palette: 'Royal Red & Gold'
  },
  {
    icon: '🕊️',
    label: 'Condolences',
    prompt: 'Papan bunga duka cita Turut Berduka Cita atas berpulangnya Bapak Hartono dari Keluarga Besar Wijaya',
    occasion: 'Duka Cita',
    palette: 'Deep Charcoal & Silver'
  },
  {
    icon: '🎓',
    label: 'Graduation',
    prompt: 'Papan bunga wisuda Happy Graduation untuk Annisa Putri, S.Kom dari Sahabat Squad Kampus',
    occasion: 'Wisuda',
    palette: 'Teal & Sunflower Yellow'
  }
]

const selectPreset = (preset: typeof presetPrompts[0]) => {
  promptText.value = preset.prompt
  occasion.value = preset.occasion
  preferredPalette.value = preset.palette
}

const handleGenerate = async () => {
  if (!promptText.value.trim()) {
    errorMsg.value = 'Please enter a prompt or select a preset inspiration below.'
    return
  }

  errorMsg.value = null
  isGenerating.value = true
  generatedPayload.value = null

  try {
    const result = await customDesignAiService.generateDesign({
      prompt: promptText.value.trim(),
      occasion: occasion.value.trim() || undefined,
      preferred_palette: preferredPalette.value.trim() || undefined,
      recipient: recipient.value.trim() || undefined,
      sender: sender.value.trim() || undefined,
      physical_size_id: physicalSizeId.value
    })

    generatedPayload.value = result
  } catch (err: any) {
    console.error('AI Generation error:', err)
    errorMsg.value = err?.data?.message || err?.message || 'Failed to generate design. Please try again.'
  } finally {
    isGenerating.value = false
  }
}

const handleApply = () => {
  if (generatedPayload.value) {
    emit('apply', generatedPayload.value)
    emit('close')
  }
}

const handleReset = () => {
  promptText.value = ''
  occasion.value = ''
  preferredPalette.value = ''
  recipient.value = ''
  sender.value = ''
  generatedPayload.value = null
  errorMsg.value = null
}
</script>

<template>
  <div v-if="show" class="dr-modal-backdrop" @click.self="emit('close')">
    <div class="dr-modal ai-modal max-w-xl w-full">
      <!-- Modal Header -->
      <div class="dr-modal-head border-b border-gray-100 pb-3 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-amber-500 via-rose-500 to-indigo-600 flex items-center justify-center text-white text-sm shadow-sm">
            ✨
          </div>
          <div>
            <h2 class="text-base font-bold text-gray-900 leading-tight">AI Flower Board Designer</h2>
            <p class="text-xs text-gray-500">Describe what you want and AI will arrange the entire board</p>
          </div>
        </div>
        <button class="dr-modal-close" @click="emit('close')" aria-label="Close modal">×</button>
      </div>

      <!-- Modal Body -->
      <div class="choice-body space-y-4 py-4 max-h-[75vh] overflow-y-auto px-1">
        <!-- Preset Chips -->
        <div>
          <label class="text-xs font-semibold text-gray-600 mb-1.5 block">Quick Inspiration Presets:</label>
          <div class="flex flex-wrap gap-1.5">
            <button
              v-for="p in presetPrompts"
              :key="p.label"
              class="text-xs px-2.5 py-1.5 rounded-full border border-gray-200 hover:border-emerald-500 hover:bg-emerald-50 text-gray-700 font-medium transition-all flex items-center gap-1.5 shadow-2xs"
              @click="selectPreset(p)"
            >
              <span>{{ p.icon }}</span>
              <span>{{ p.label }}</span>
            </button>
          </div>
        </div>

        <!-- Prompt Textarea -->
        <div>
          <div class="flex justify-between items-center mb-1">
            <label class="text-xs font-semibold text-gray-700">Custom Prompt / Occasion Details:</label>
            <span class="text-[11px] text-gray-400">{{ promptText.length }}/1000</span>
          </div>
          <textarea
            v-model="promptText"
            rows="3"
            class="w-full p-2.5 text-xs rounded-lg border border-gray-300 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-hidden resize-none transition-all placeholder:text-gray-400"
            placeholder="e.g., Buatkan papan bunga pernikahan tema rustic pastel untuk Dimas & Sarah dari PT Maju Bersama, ada ucapan Selamat Menempuh Hidup Baru..."
            maxlength="1000"
            :disabled="isGenerating"
          ></textarea>
        </div>

        <!-- Advanced Options Toggle -->
        <div>
          <button
            type="button"
            class="text-xs font-medium text-emerald-600 hover:text-emerald-700 flex items-center gap-1 cursor-pointer"
            @click="showAdvanced = !showAdvanced"
          >
            <span>{{ showAdvanced ? '▲ Hide Advanced Options' : '▼ Add Specific Details (Optional)' }}</span>
          </button>

          <div v-if="showAdvanced" class="mt-2.5 p-3 bg-gray-50 rounded-lg border border-gray-200 grid grid-cols-1 sm:grid-cols-2 gap-2.5 text-xs">
            <div>
              <label class="font-medium text-gray-600 block mb-1">Occasion / Event</label>
              <input
                v-model="occasion"
                type="text"
                placeholder="e.g., Wedding, Grand Opening"
                class="w-full p-2 rounded border border-gray-300 bg-white"
              />
            </div>
            <div>
              <label class="font-medium text-gray-600 block mb-1">Color Palette</label>
              <input
                v-model="preferredPalette"
                type="text"
                placeholder="e.g., Pastel Gold, Navy Silver"
                class="w-full p-2 rounded border border-gray-300 bg-white"
              />
            </div>
            <div>
              <label class="font-medium text-gray-600 block mb-1">Recipient Name</label>
              <input
                v-model="recipient"
                type="text"
                placeholder="e.g., Dimas & Sarah"
                class="w-full p-2 rounded border border-gray-300 bg-white"
              />
            </div>
            <div>
              <label class="font-medium text-gray-600 block mb-1">Sender Name</label>
              <input
                v-model="sender"
                type="text"
                placeholder="e.g., PT Maju Bersama"
                class="w-full p-2 rounded border border-gray-300 bg-white"
              />
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="errorMsg" class="p-2.5 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-xs flex items-center gap-2">
          <span>⚠️</span>
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Loading Animation State -->
        <div v-if="isGenerating" class="py-6 flex flex-col items-center justify-center gap-3 bg-gray-50 rounded-xl border border-gray-200">
          <div class="relative w-10 h-10">
            <div class="w-10 h-10 border-3 border-emerald-200 border-t-emerald-600 rounded-full animate-spin"></div>
          </div>
          <div class="text-center">
            <p class="text-xs font-bold text-gray-800 animate-pulse">Generating Flower Board Layout…</p>
            <p class="text-[11px] text-gray-500">Choosing typography, harmony palettes, and floral crests</p>
          </div>
        </div>

        <!-- Generated Preview Card -->
        <div v-if="generatedPayload && !isGenerating" class="p-3.5 bg-emerald-50/80 border border-emerald-200 rounded-xl space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-emerald-800 flex items-center gap-1.5">
              <span>🎉</span> Design Generated Successfully!
            </span>
            <span class="text-[11px] text-emerald-700 font-mono font-semibold uppercase px-2 py-0.5 bg-emerald-100 rounded-full">
              Size: {{ generatedPayload.layout.physicalSizeId }}
            </span>
          </div>

          <div class="bg-white p-3 rounded-lg border border-emerald-100 text-xs space-y-1.5 shadow-2xs">
            <div class="flex items-start gap-2">
              <span class="text-gray-400 font-bold">Upper:</span>
              <div>
                <p class="font-bold text-gray-900">{{ generatedPayload.sections.upper.header.text }}</p>
                <p class="text-gray-600 text-[11px]">{{ generatedPayload.sections.upper.body.text }}</p>
              </div>
            </div>
            <div class="border-t border-gray-100 pt-1.5 flex items-start gap-2">
              <span class="text-gray-400 font-bold">Lower:</span>
              <div>
                <p class="font-bold text-gray-900">{{ generatedPayload.sections.lower.header.text }}</p>
                <p class="text-gray-600 text-[11px]">{{ generatedPayload.sections.lower.body.text }}</p>
              </div>
            </div>
            <div class="border-t border-gray-100 pt-1.5 flex items-center gap-2 text-[11px] text-gray-500">
              <span>Theme:</span>
              <div class="flex items-center gap-1">
                <span class="w-3.5 h-3.5 rounded-full border border-gray-200" :style="{ backgroundColor: generatedPayload.sections.upper.bgColorHex }"></span>
                <span class="w-3.5 h-3.5 rounded-full border border-gray-200" :style="{ backgroundColor: generatedPayload.sections.lower.bgColorHex }"></span>
                <span class="w-3.5 h-3.5 rounded-full border border-gray-200" :style="{ backgroundColor: generatedPayload.layout.border.colorHex }"></span>
              </div>
              <span class="ml-auto font-medium">Border: {{ generatedPayload.layout.border.style }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="choice-actions pt-3 border-t border-gray-100 flex items-center justify-end gap-2">
        <button
          v-if="generatedPayload"
          class="choice-btn secondary"
          @click="handleReset"
          :disabled="isGenerating"
        >
          Try Another
        </button>
        <button
          v-if="!generatedPayload"
          class="choice-btn secondary"
          @click="emit('close')"
          :disabled="isGenerating"
        >
          Cancel
        </button>
        <button
          v-if="!generatedPayload"
          class="choice-btn primary flex items-center gap-1.5"
          @click="handleGenerate"
          :disabled="isGenerating"
        >
          <span>✨ Generate Board</span>
        </button>
        <button
          v-else
          class="choice-btn primary bg-emerald-600 hover:bg-emerald-700 text-white flex items-center gap-1.5 font-bold"
          @click="handleApply"
        >
          <span>🎨 Apply to Canvas</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ai-modal {
  background: #ffffff;
  border-radius: 1rem;
  padding: 1.25rem;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}
</style>
