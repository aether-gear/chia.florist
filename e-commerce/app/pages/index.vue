<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { productService } from '~/services/productService'

useHead({
  title: 'Chia Florist - Flower Boards & Custom Simulator',
  meta: [
    {
      name: 'description',
      content: 'Pesan papan bunga ucapan pernikahan, duka cita, peresmian gedung, dan wisuda terbaik atau rancang papan bunga Anda sendiri secara real-time di Chia Florist.'
    }
  ]
})

const { formatRupiah } = useCart()

// Carousel Slides
const slides = [
  {
    image: '/florist.jpg',
    title: 'Papan Bunga Elegan untuk Momen Berharga',
    subtitle: 'Rangkaian bunga segar dan material tahan cuaca untuk pernikahan, duka cita, & peresmian.',
    ctaText: 'Lihat Katalog',
    ctaLink: '/catalog',
    ctaVariant: 'primary' as const
  },
  {
    image: '/flowerist.jpg',
    title: 'Desain Papan Bunga Kustom Real-Time',
    subtitle: 'Kreasikan susunan busa, teks ucapan, dan dekorasi bunga langsung di browser Anda.',
    ctaText: 'Coba Simulator',
    ctaLink: '/products/custom',
    ctaVariant: 'primary' as const
  },
  {
    image: '/flower.jpg',
    title: 'Pengantaran Tepat Waktu Langsung ke Lokasi',
    subtitle: 'Armada kami siap mengantarkan pesanan bunga papan Anda ke seluruh gedung dan venue acara.',
    ctaText: 'Pesan Sekarang',
    ctaLink: '/catalog',
    ctaVariant: 'primary' as const
  }
]

const currentIndex = ref(0)
let intervalTimer: any = null

const nextSlide = () => {
  currentIndex.value = (currentIndex.value + 1) % slides.length
  resetTimer()
}

const prevSlide = () => {
  currentIndex.value = (currentIndex.value - 1 + slides.length) % slides.length
  resetTimer()
}

const goToSlide = (index: number) => {
  currentIndex.value = index
  resetTimer()
}

const stopTimer = () => {
  if (intervalTimer) {
    clearInterval(intervalTimer)
    intervalTimer = null
  }
}

const startTimer = () => {
  stopTimer()
  intervalTimer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % slides.length
  }, 5000)
}

const resetTimer = () => {
  startTimer()
}

const productOfferings = ref<any[]>([])
const isLoading = ref(false)
const hasError = ref(false)

onMounted(async () => {
  startTimer()
  isLoading.value = true
  hasError.value = false

  try {
    const list = await productService.getCatalogProducts()
    productOfferings.value = [
      ...list.map(p => ({
        id: p.id,
        name: p.name,
        slug: p.slug || '',
        price: p.price,
        image: p.image,
        isAvailable: p.isAvailable,
        rating: p.rating || 4.8,
        reviews: p.reviews || 120
      })),
      {
        id: 'custom',
        name: 'Custom Board Simulator',
        slug: 'custom',
        price: 150000,
        image: '/images/custom-preview.png',
        isAvailable: true,
        rating: 5.0,
        reviews: 89,
        isCustomRoute: true
      }
    ]
  } catch (err) {
    console.error('Failed to load homepage offerings:', err)
    hasError.value = true
    // Fallback item so the user can still access custom simulator
    productOfferings.value = [
      {
        id: 'custom',
        name: 'Custom Board Simulator',
        slug: 'custom',
        price: 150000,
        image: '/images/custom-preview.png',
        isAvailable: true,
        rating: 5.0,
        reviews: 89,
        isCustomRoute: true
      }
    ]
  } finally {
    isLoading.value = false
  }
})

onUnmounted(() => {
  stopTimer()
})
</script>

<template>
  <div class="w-full bg-white text-gray-900 font-brand">

    <!-- HERO BANNER CAROUSEL (Tokopedia-Style Contained Width) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4 sm:pt-6">
      <div 
        class="relative w-full rounded-2xl sm:rounded-3xl overflow-hidden shadow-xs border border-gray-100 bg-gray-900 aspect-[16/8] sm:aspect-[21/9] min-h-[220px] max-h-[420px]"
        @mouseenter="stopTimer"
        @mouseleave="startTimer"
      >
        <!-- Slides -->
        <div 
          v-for="(slide, index) in slides" 
          :key="index"
          :class="[
            'absolute inset-0 w-full h-full transition-opacity duration-700 ease-in-out',
            currentIndex === index ? 'opacity-100 z-10' : 'opacity-0 z-0 pointer-events-none'
          ]"
        >
          <!-- Background Image -->
          <img 
            :src="slide.image" 
            :alt="slide.title" 
            class="absolute inset-0 w-full h-full object-cover"
          />
          <!-- Dark Contrast Gradient Overlay -->
          <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-black/50 to-transparent z-10"></div>

          <!-- Slide Content (Without Pill Badge above Header) -->
          <div class="relative z-20 h-full flex flex-col justify-center max-w-xl p-6 sm:p-10 lg:p-12 text-white">
            <h1 class="text-xl sm:text-3xl lg:text-4xl font-extrabold leading-tight mb-1.5 sm:mb-3 drop-shadow-sm line-clamp-2">
              {{ slide.title }}
            </h1>

            <p class="text-xs sm:text-sm text-gray-200 mb-4 sm:mb-6 line-clamp-2 max-w-md hidden sm:block">
              {{ slide.subtitle }}
            </p>

            <div>
              <CButton 
                :to="slide.ctaLink" 
                variant="primary" 
                size="md"
                class="shadow-md"
              >
                <span>{{ slide.ctaText }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 ml-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                </svg>
              </CButton>
            </div>
          </div>
        </div>

        <!-- Left Arrow Button -->
        <button 
          @click="prevSlide" 
          aria-label="Previous Slide"
          class="absolute left-3 top-1/2 -translate-y-1/2 z-20 w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-white/80 hover:bg-white text-gray-800 flex items-center justify-center shadow-md backdrop-blur-sm transition-all hover:scale-105 active:scale-95 cursor-pointer"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 sm:w-5 sm:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
        </button>

        <!-- Right Arrow Button -->
        <button 
          @click="nextSlide" 
          aria-label="Next Slide"
          class="absolute right-3 top-1/2 -translate-y-1/2 z-20 w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-white/80 hover:bg-white text-gray-800 flex items-center justify-center shadow-md backdrop-blur-sm transition-all hover:scale-105 active:scale-95 cursor-pointer"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 sm:w-5 sm:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
          </svg>
        </button>

        <!-- Pagination Dots -->
        <div class="absolute bottom-3 sm:bottom-4 right-4 sm:right-6 z-20 flex items-center gap-1.5 sm:gap-2">
          <button 
            v-for="(slide, index) in slides" 
            :key="index"
            @click="goToSlide(index)"
            :aria-label="`Go to slide ${index + 1}`"
            :class="[
              'h-2 rounded-full transition-all duration-300 cursor-pointer',
              currentIndex === index 
                ? 'w-6 bg-[#4ade80]' 
                : 'w-2 bg-white/60 hover:bg-white'
            ]"
          />
        </div>
      </div>
    </section>

    <!-- HORIZONTALLY SCROLLABLE SHORTCUT BUTTONS -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-5">
      <div class="flex items-center gap-3 sm:gap-4 overflow-x-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden py-1">
        
        <!-- Button 1: Katalog Papan Bunga (Light Gray Button) -->
        <NuxtLink
          to="/catalog"
          class="flex items-center gap-3.5 px-5 py-3.5 bg-gray-100 hover:bg-gray-200/80 active:scale-[0.99] border border-gray-200/80 rounded-2xl shadow-2xs hover:shadow-xs transition-all shrink-0 min-w-[250px] sm:min-w-[280px] group cursor-pointer text-left"
        >
          <div class="w-11 h-11 rounded-xl bg-white text-[#1b4332] flex items-center justify-center shrink-0 border border-gray-200 shadow-2xs group-hover:scale-105 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="text-sm sm:text-base font-bold text-gray-800 group-hover:text-gray-900 transition-colors">
              Katalog Papan Bunga
            </h2>
            <p class="text-xs text-gray-500 truncate mt-0.5">
              Pilihan bunga papan siap pesan
            </p>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-gray-400 group-hover:text-gray-700 group-hover:translate-x-0.5 transition-all shrink-0 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
          </svg>
        </NuxtLink>

        <!-- Button 2: Custom Board Simulator (Light Green Pastel Opacity 50 Button) -->
        <NuxtLink
          to="/products/custom"
          class="flex items-center gap-3.5 px-5 py-3.5 bg-[#4ade80]/50 hover:bg-[#4ade80]/65 active:scale-[0.99] border border-[#4ade80]/70 rounded-2xl shadow-2xs hover:shadow-xs transition-all shrink-0 min-w-[250px] sm:min-w-[280px] group cursor-pointer text-left"
        >
          <div class="w-11 h-11 rounded-xl bg-white text-[#245842] flex items-center justify-center shrink-0 border border-[#4ade80]/40 shadow-2xs group-hover:scale-105 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42" />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="text-sm sm:text-base font-bold text-[#245842] group-hover:text-[#1b4332] transition-colors">
              Custom Board Simulator
            </h2>
            <p class="text-xs text-[#245842]/80 truncate mt-0.5">
              Desain kustom interaktif 2D
            </p>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-[#245842]/70 group-hover:text-[#245842] group-hover:translate-x-0.5 transition-all shrink-0 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
          </svg>
        </NuxtLink>

      </div>
    </section>

    <!-- FEATURED PRODUCTS GRID (With pb-40 / 10rem bottom padding) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-8 pb-40">
      
      <!-- Section Header -->
      <div class="flex flex-col sm:flex-row sm:items-end justify-between mb-8 pb-4 border-b border-gray-100 gap-4">
        <div>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-gray-900 tracking-tight">
            Koleksi Produk Pilihan
          </h2>
          <p class="text-xs sm:text-sm text-gray-500 mt-1">
            Papan bunga berkualitas tinggi siap rangkai dan kirim cepat ke lokasi tujuan.
          </p>
        </div>
        <div>
          <NuxtLink 
            to="/catalog" 
            class="text-xs sm:text-sm font-bold text-[#1b4332] hover:text-[#143326] flex items-center gap-1 transition-colors"
          >
            <span>Lihat Semua Produk</span>
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
          </NuxtLink>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6">
        <div v-for="n in 4" :key="n" class="bg-gray-50 border border-gray-100 rounded-2xl p-4 animate-pulse space-y-4">
          <div class="aspect-[4/3] bg-gray-200 rounded-xl w-full"></div>
          <div class="h-4 bg-gray-200 rounded w-3/4"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2"></div>
          <div class="h-8 bg-gray-200 rounded-xl w-full"></div>
        </div>
      </div>

      <!-- Products Grid -->
      <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6">
        <div 
          v-for="item in productOfferings" 
          :key="item.id"
          class="group bg-white border border-gray-100 rounded-2xl overflow-hidden shadow-2xs hover:shadow-md transition-all duration-300 flex flex-col justify-between cursor-pointer"
          @click="navigateTo(item.id === 'custom' || item.isCustomRoute ? '/products/custom' : `/products/${item.slug || item.id}`)"
        >
          <div>
            <!-- Product Image -->
            <div class="aspect-[4/3] relative bg-gray-50 overflow-hidden">
              <img 
                :src="item.image || '/images/custom-preview.png'" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              
              <!-- Badges -->
              <span 
                v-if="item.id === 'custom' || item.isCustomRoute" 
                class="absolute top-2.5 left-2.5 bg-[#1b4332] text-white text-[9px] sm:text-[10px] font-bold tracking-wider uppercase px-2 py-0.5 rounded-lg border border-[#143326] shadow-2xs z-10"
              >
                Interactive 2D
              </span>

              <span 
                v-if="item.status === 'inactive'" 
                class="absolute top-2.5 right-2.5 bg-amber-100 text-amber-900 text-[9px] sm:text-[10px] font-bold tracking-wider uppercase px-2 py-0.5 rounded-lg border border-amber-200 shadow-2xs z-10"
              >
                Preview Only
              </span>
              <span 
                v-else-if="!item.isAvailable" 
                class="absolute top-2.5 right-2.5 bg-red-100 text-red-800 text-[9px] sm:text-[10px] font-bold tracking-wider uppercase px-2 py-0.5 rounded-lg border border-red-200 shadow-2xs z-10"
              >
                Sold Out
              </span>
            </div>

            <!-- Product Info -->
            <div class="p-3.5 sm:p-4 space-y-2">
              <!-- Rating & Reviews -->
              <div class="flex items-center gap-1 text-[11px] text-yellow-500 font-bold">
                <span>⭐ {{ item.rating ? item.rating.toFixed(1) : '4.8' }}</span>
                <span class="text-gray-300">|</span>
                <span class="text-gray-400 font-normal">({{ item.reviews || 80 }})</span>
              </div>

              <!-- Product Name -->
              <h3 class="font-bold text-gray-900 text-xs sm:text-sm line-clamp-2 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>

              <!-- Price Display -->
              <div>
                <p class="text-[10px] text-gray-400 font-medium">Harga mulai dari</p>
                <p class="text-sm sm:text-base font-extrabold text-[#1b4332]">
                  {{ formatRupiah(item.price) }}
                </p>
              </div>
            </div>
          </div>

          <!-- Action Button -->
          <div class="p-3.5 sm:p-4 pt-0">
            <CButton
              v-if="item.isAvailable && item.status !== 'inactive'"
              :to="item.id === 'custom' || item.isCustomRoute ? '/products/custom' : `/products/${item.slug || item.id}`"
              :variant="item.id === 'custom' || item.isCustomRoute ? 'primary' : 'secondary'"
              size="sm"
              class="w-full"
              @click.stop
            >
              {{ item.id === 'custom' || item.isCustomRoute ? 'Desain Sekarang' : 'Lihat Detail' }}
            </CButton>

            <CButton
              v-else-if="item.status === 'inactive'"
              :to="`/products/${item.slug || item.id}`"
              variant="outline"
              size="sm"
              class="w-full text-amber-800 border-amber-300 hover:bg-amber-50"
              @click.stop
            >
              Lihat Preview
            </CButton>

            <CButton
              v-else
              disabled
              variant="outline"
              size="sm"
              class="w-full opacity-60 cursor-not-allowed"
            >
              Habis Terjual
            </CButton>
          </div>
        </div>
      </div>

      <!-- Bottom View All Button -->
      <div class="mt-12 text-center">
        <CButton 
          to="/catalog" 
          variant="outline" 
          size="md"
        >
          <span>Lihat Semua Produk di Katalog</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 ml-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
          </svg>
        </CButton>
      </div>

    </section>

  </div>
</template>

<style scoped></style>