<script setup lang="ts">
// CErrorDisplay — Shared branded error display component
// Used by: error.vue (global Nuxt error page) and pages/[...slug].vue (catch-all 404)

interface Props {
  statusCode: 404 | 500 | number
  title?: string
  message?: string
  /** When true, clicking navigation calls clearError() before navigating (used in error.vue) */
  clearOnNavigate?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  clearOnNavigate: false
})

const emit = defineEmits<{
  clear: []
}>()

const router = useRouter()

// --- Derived content ---
const is404 = computed(() => props.statusCode === 404)
const is500 = computed(() => props.statusCode === 500)

const displayTitle = computed(() => {
  if (props.title) return props.title
  if (is404.value) return 'Halaman Tidak Ditemukan'
  if (is500.value) return 'Terjadi Kesalahan Server'
  return 'Terjadi Kesalahan'
})

const displayMessage = computed(() => {
  if (props.message) return props.message
  if (is404.value)
    return 'Maaf, halaman yang Anda cari tidak dapat ditemukan. Mungkin tautan sudah dipindahkan atau tidak lagi tersedia.'
  if (is500.value)
    return 'Kami sedang mengalami gangguan teknis. Tim kami sedang menanganinya. Silakan coba lagi beberapa saat.'
  return 'Terjadi kesalahan yang tidak terduga. Silakan coba lagi.'
})

// --- Navigation helpers ---
const navigate = (path: string) => {
  if (props.clearOnNavigate) {
    emit('clear')
  }
  router.push(path)
}

const goBack = () => {
  if (props.clearOnNavigate) {
    emit('clear')
  }
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}
</script>

<template>
  <div class="min-h-[60vh] flex flex-col items-center justify-center px-6 py-16 text-center font-sans">

    <!-- Illustration -->
    <div class="mb-8 select-none" aria-hidden="true">
      <!-- 404: wilted/searching flower -->
      <template v-if="is404">
        <svg
          width="160"
          height="160"
          viewBox="0 0 160 160"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class="mx-auto"
        >
          <!-- Stem -->
          <path d="M80 130 Q78 100 80 70" stroke="#1b4332" stroke-width="3" stroke-linecap="round" fill="none" />
          <!-- Drooping leaves -->
          <path d="M80 105 Q65 95 60 80" stroke="#4ade80" stroke-width="2.5" stroke-linecap="round" fill="none" />
          <path d="M80 95 Q95 85 98 68" stroke="#4ade80" stroke-width="2.5" stroke-linecap="round" fill="none" />
          <!-- Petals (wilted, drooping) -->
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.6" transform="rotate(0 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.55" transform="rotate(51 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.5" transform="rotate(102 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.55" transform="rotate(153 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.6" transform="rotate(204 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.55" transform="rotate(255 80 70)" />
          <ellipse cx="80" cy="52" rx="9" ry="18" fill="#4ade80" opacity="0.5" transform="rotate(306 80 70)" />
          <!-- Center -->
          <circle cx="80" cy="70" r="11" fill="#1b4332" />
          <!-- Question mark in center -->
          <text x="80" y="75" text-anchor="middle" font-size="13" font-weight="bold" fill="#4ade80" font-family="Inter, sans-serif">?</text>
          <!-- Pot -->
          <rect x="63" y="130" width="34" height="20" rx="4" fill="#1b4332" opacity="0.15" />
          <rect x="60" y="128" width="40" height="6" rx="3" fill="#1b4332" opacity="0.2" />
        </svg>
      </template>

      <!-- 500: sad wilted flower with warning -->
      <template v-else-if="is500">
        <svg
          width="160"
          height="160"
          viewBox="0 0 160 160"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class="mx-auto"
        >
          <!-- Stem bent -->
          <path d="M80 130 Q70 105 72 80" stroke="#1b4332" stroke-width="3" stroke-linecap="round" fill="none" />
          <!-- Leaves -->
          <path d="M76 108 Q60 100 55 85" stroke="#4ade80" stroke-width="2.5" stroke-linecap="round" fill="none" />
          <path d="M75 93 Q90 88 92 74" stroke="#4ade80" stroke-width="2.5" stroke-linecap="round" fill="none" />
          <!-- Petals muted/sad -->
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.45" transform="rotate(0 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.4" transform="rotate(51 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.35" transform="rotate(102 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.4" transform="rotate(153 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.45" transform="rotate(204 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.4" transform="rotate(255 72 68)" />
          <ellipse cx="72" cy="52" rx="9" ry="17" fill="#4ade80" opacity="0.35" transform="rotate(306 72 68)" />
          <!-- Center -->
          <circle cx="72" cy="68" r="11" fill="#1b4332" />
          <!-- Warning mark -->
          <text x="72" y="73" text-anchor="middle" font-size="13" font-weight="bold" fill="#fbbf24" font-family="Inter, sans-serif">!</text>
          <!-- Warning triangle badge -->
          <polygon points="115,32 130,58 100,58" fill="#fbbf24" opacity="0.9" />
          <text x="115" y="53" text-anchor="middle" font-size="14" font-weight="bold" fill="#78350f" font-family="Inter, sans-serif">!</text>
          <!-- Pot -->
          <rect x="55" y="130" width="34" height="20" rx="4" fill="#1b4332" opacity="0.15" />
          <rect x="52" y="128" width="40" height="6" rx="3" fill="#1b4332" opacity="0.2" />
        </svg>
      </template>

      <!-- Generic fallback illustration -->
      <template v-else>
        <svg width="120" height="120" viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg" class="mx-auto">
          <circle cx="60" cy="60" r="50" stroke="#4ade80" stroke-width="3" fill="none" opacity="0.4" />
          <text x="60" y="68" text-anchor="middle" font-size="30" font-weight="bold" fill="#1b4332" font-family="Inter, sans-serif">!</text>
        </svg>
      </template>
    </div>

    <!-- Status Code -->
    <p class="text-8xl font-bold text-[#4ade80] leading-none mb-2 tracking-tight" style="font-family: 'Inter', sans-serif;">
      {{ statusCode }}
    </p>

    <!-- Title -->
    <h1 class="text-2xl md:text-3xl font-bold text-gray-900 mt-4 mb-3" style="font-family: 'Inter', sans-serif;">
      {{ displayTitle }}
    </h1>

    <!-- Message -->
    <p class="text-gray-500 text-sm md:text-base max-w-md mx-auto leading-relaxed mb-10">
      {{ displayMessage }}
    </p>

    <!-- Action Buttons -->
    <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
      <button
        class="inline-flex items-center justify-center font-semibold transition-all duration-200 rounded-xl px-5 py-2.5 text-sm bg-[#4ade80] hover:bg-[#34d399] text-[#245842] shadow-xs hover:shadow focus:ring-2 focus:ring-[#4ade80]/30 cursor-pointer"
        @click="navigate('/')"
      >
        <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        Kembali ke Beranda
      </button>

      <button
        class="inline-flex items-center justify-center font-semibold transition-all duration-200 rounded-xl px-5 py-2.5 text-sm border border-gray-200 text-gray-700 hover:bg-gray-50 focus:ring-2 focus:ring-gray-200 cursor-pointer"
        @click="navigate('/catalog')"
      >
        <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
        </svg>
        Lihat Katalog
      </button>

      <button
        class="inline-flex items-center justify-center font-semibold transition-all duration-200 rounded-xl px-5 py-2.5 text-sm text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-100 cursor-pointer"
        @click="goBack"
      >
        <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Kembali
      </button>
    </div>

  </div>
</template>
