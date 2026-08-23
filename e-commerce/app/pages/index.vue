<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart'
import { productService } from '~/services/productService'

useHead({
  title: 'Chia Florist — Papan Bunga & Custom Simulator Online',
  meta: [
    {
      name: 'description',
      content: 'Pesan papan bunga ucapan pernikahan, duka cita, peresmian gedung, dan wisuda terbaik atau rancang papan bunga Anda sendiri secara real-time di Chia Florist.'
    },
    { property: 'og:title', content: 'Chia Florist — Papan Bunga & Custom Simulator Online' },
    { property: 'og:description', content: 'Pesan papan bunga ucapan pernikahan, duka cita, peresmian gedung, dan wisuda terbaik atau rancang papan bunga Anda sendiri secara real-time di Chia Florist.' },
    { property: 'og:type', content: 'website' },
    { property: 'og:url', content: 'https://chiaflorist.com/' },
    { property: 'og:image', content: 'https://chiaflorist.com/florist.jpg' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: 'Chia Florist — Papan Bunga & Custom Simulator Online' },
    { name: 'twitter:description', content: 'Pesan papan bunga ucapan pernikahan, duka cita, peresmian gedung, dan wisuda terbaik di Chia Florist.' }
  ],
  link: [
    { rel: 'canonical', href: 'https://chiaflorist.com/' }
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@graph': [
          {
            '@type': 'Florist',
            '@id': 'https://chiaflorist.com/#organization',
            'name': 'Chia Florist',
            'url': 'https://chiaflorist.com/',
            'logo': 'https://chiaflorist.com/images/logo.png',
            'description': 'Penyedia papan bunga ucapan dan buket bunga segar berkualitas tinggi dengan layanan custom simulator online.',
            'telephone': '+62-817-523-4999',
            'address': {
              '@type': 'PostalAddress',
              'streetAddress': 'Jl. Argotirto No 06 RT 04 RW 02 Kp. Air Terjun Kel. Sungai Daeng',
              'addressLocality': 'Kab. Bangka Barat',
              'addressRegion': 'Kep. Bangka Belitung',
              'postalCode': '33311',
              'addressCountry': 'ID'
            }
          },
          {
            '@type': 'WebSite',
            '@id': 'https://chiaflorist.com/#website',
            'url': 'https://chiaflorist.com/',
            'name': 'Chia Florist',
            'potentialAction': {
              '@type': 'SearchAction',
              'target': 'https://chiaflorist.com/search?q={search_term_string}',
              'query-input': 'required name=search_term_string'
            }
          }
        ]
      })
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

export interface HomeOfferingProduct {
  id: string
  name: string
  slug: string
  price: number
  image: string
  isAvailable: boolean
  rating: number
  reviews: number
  status?: string
  isCustomRoute?: boolean
}

// Server-rendered product offerings via useAsyncData for SEO crawlers
const { data: offeringsData, status, error: fetchError } = await useAsyncData<HomeOfferingProduct[]>('home-products', async () => {
  try {
    const list = await productService.getCatalogProducts()
    return list.map((p): HomeOfferingProduct => ({
      id: p.id,
      name: p.name,
      slug: p.slug || '',
      price: p.price,
      image: p.image,
      isAvailable: p.isAvailable,
      rating: p.rating || 4.8,
      reviews: p.reviews || 120,
      status: p.status || (p.isAvailable ? 'active' : 'inactive')
    }))
  } catch (err) {
    console.error('Failed to load homepage offerings during SSR:', err)
    return []
  }
})

const productOfferings = computed(() => offeringsData.value || [])
const isLoading = computed(() => status.value === 'pending')
const hasError = computed(() => Boolean(fetchError.value && productOfferings.value.length === 0))

onMounted(() => {
  startTimer()
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

    <!-- RESPONSIVE SHORTCUT BUTTONS -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-5">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
        
        <!-- Button 1: Katalog Papan Bunga (Light Gray Button) -->
        <NuxtLink
          to="/catalog"
          class="flex items-center gap-3.5 px-4 sm:px-5 py-3.5 bg-gray-100 hover:bg-gray-200/80 active:scale-[0.99] border border-gray-200/80 rounded-2xl shadow-2xs hover:shadow-xs transition-all group cursor-pointer text-left w-full"
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
          class="flex items-center gap-3.5 px-4 sm:px-5 py-3.5 bg-[#4ade80]/50 hover:bg-[#4ade80]/65 active:scale-[0.99] border border-[#4ade80]/70 rounded-2xl shadow-2xs hover:shadow-xs transition-all group cursor-pointer text-left w-full"
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

    <!-- FEATURED SECTION: Custom Board Simulator Highlight -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#1b4332] via-[#245842] to-[#122e22] text-white p-6 sm:p-10 lg:p-12 shadow-xl border border-emerald-800/40">
        <!-- Decorative Glow Accents -->
        <div class="absolute -right-20 -top-20 w-80 h-80 bg-emerald-400/10 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -left-20 -bottom-20 w-80 h-80 bg-[#4ade80]/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="relative z-10 grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">
          
          <!-- Text & Highlights -->
          <div class="lg:col-span-7 space-y-5">
            <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-white/10 backdrop-blur-md border border-white/20 text-emerald-300 text-xs font-bold uppercase tracking-widest">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-300 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42" />
              </svg>
              <span>Interactive 2D Board Simulator</span>
            </div>

            <h2 class="text-2xl sm:text-4xl font-extrabold tracking-tight text-white leading-tight sm:leading-snug">
              Desain Papan Bunga Kustom Sendiri Secara Real-Time
            </h2>

            <p class="text-xs sm:text-sm text-gray-200 leading-relaxed max-w-xl">
              Ingin tampilan papan bunga yang benar-benar unik? Kreasikan warna busa dasar, atur teks ucapan, nama pengirim, dan ornamen bunga sudut secara langsung melalui simulator interaktif kami.
            </p>

            <!-- Features Bullet Points -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
              <div class="flex items-center gap-2.5 text-xs sm:text-sm font-medium text-gray-200">
                <div class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-300 flex items-center justify-center shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </div>
                <span>Visualisasi 2D langsung di layar</span>
              </div>
              <div class="flex items-center gap-2.5 text-xs sm:text-sm font-medium text-gray-200">
                <div class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-300 flex items-center justify-center shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </div>
                <span>Pilihan warna & font lengkap</span>
              </div>
              <div class="flex items-center gap-2.5 text-xs sm:text-sm font-medium text-gray-200">
                <div class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-300 flex items-center justify-center shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </div>
                <span>Simpan draft & review desain</span>
              </div>
              <div class="flex items-center gap-2.5 text-xs sm:text-sm font-medium text-gray-200">
                <div class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-300 flex items-center justify-center shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </div>
                <span>Estimasi harga transparan</span>
              </div>
            </div>

            <div class="pt-4 flex flex-wrap items-center gap-3.5">
              <CButton 
                to="/products/custom" 
                variant="primary" 
                size="lg"
                class="bg-[#4ade80] hover:bg-[#3ec470] text-[#1b4332] font-black shadow-lg hover:shadow-xl transition-all"
              >
                <span>Coba Simulator Interaktif</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 ml-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                </svg>
              </CButton>
              
              <NuxtLink 
                to="/catalog"
                class="px-5 py-3 rounded-xl border border-white/20 text-white hover:bg-white/10 text-xs font-bold transition flex items-center gap-1.5"
              >
                <span>Lihat Koleksi Jadi</span>
              </NuxtLink>
            </div>
          </div>

          <!-- Preview Visual Mockup Showcase -->
          <div class="lg:col-span-5">
            <div class="relative group cursor-pointer" @click="navigateTo('/products/custom')">
              <div class="aspect-[4/3] rounded-2xl overflow-hidden border-2 border-white/20 shadow-2xl bg-black/40 relative">
                <img 
                  src="/images/custom-preview.png" 
                  alt="Custom Board Simulator Preview"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent flex flex-col justify-end p-5">
                  <span class="text-[10px] font-black uppercase tracking-widest text-emerald-300">Live 2D Canvas</span>
                  <p class="text-xs font-bold text-white mt-0.5">Rancang papan ucapan pernikahan, duka cita & wisuda</p>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </section>

    <!-- FEATURED PRODUCTS GRID (With pb-40 / 10rem bottom padding) -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-6 pb-40">
      
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
        </div>
      </div>

      <!-- Products Grid (Clean Cards without internal buttons) -->
      <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6">
        <div 
          v-for="item in productOfferings" 
          :key="item.id"
          class="group bg-white border border-gray-100 rounded-2xl overflow-hidden shadow-2xs hover:shadow-md hover:border-emerald-200 transition-all duration-300 flex flex-col justify-between cursor-pointer"
          @click="navigateTo(`/products/${item.slug || item.id}`)"
        >
          <div>
            <!-- Product Image -->
            <div class="aspect-[4/3] relative bg-gray-50 overflow-hidden">
              <img 
                :src="item.image || '/images/custom-preview.png'" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              
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
              <!-- Product Name -->
              <h3 class="font-bold text-gray-900 text-xs sm:text-sm line-clamp-2 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>

              <!-- Price Display -->
              <div class="pt-1 border-t border-gray-50">
                <p class="text-[10px] text-gray-400 font-medium">Harga mulai dari</p>
                <p class="text-sm sm:text-base font-extrabold text-[#1b4332]">
                  {{ formatRupiah(item.price) }}
                </p>
              </div>
            </div>
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