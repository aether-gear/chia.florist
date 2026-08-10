<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart' // <-- INTEGRASI: Ambil useCart untuk kurensi
import { productService } from '~/services/productService'

useHead({
  title: 'Chia Florist - Flower Boards',
  meta: [
    { name: 'description', content: 'We provide customized flower boards for weddings, condolences, graduations, and corporate events.' }
  ]
})

// Ambil helper formatRupiah global agar bisa dipakai langsung di template bawah
const { formatRupiah } = useCart()

// Array gambar background Hero
const backgrounds = [
  '/florist.jpg',          // Gambar 01
  '/flowerist.jpg',          // Gambar 02
  '/flower.jpg'            // Gambar 03
]

// State untuk indeks background yang sedang aktif
const currentIndex = ref(0)
let intervalTimer: any = null

// Fungsi untuk mengganti background saat diklik
const changeBg = (index: number) => {
  currentIndex.value = index
  resetTimer() 
}

// Logic Auto Switch setiap 5 detik
const startTimer = () => {
  intervalTimer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % backgrounds.length
  }, 5000)
}

const resetTimer = () => {
  if (intervalTimer) clearInterval(intervalTimer)
  startTimer()
}

const productOfferings = ref<any[]>([])
const isLoading = ref(false)

onMounted(async () => {
  startTimer()
  isLoading.value = true
  try {
    const list = await productService.getCatalogProducts()
    productOfferings.value = [
      ...list.map(p => ({
        id: p.id,
        name: p.name,
        slug: p.slug || '',
        price: p.price,
        image: p.image,
        isAvailable: p.isAvailable
      })),
      { id: 'custom', name: 'Custom Board Simulator', slug: 'custom', price: 150000, image: '/images/custom-preview.png', isAvailable: true }
    ]
  } catch (err) {
    console.error('Failed to load homepage offerings:', err)
  } finally {
    isLoading.value = false
  }
})

const getProductImageBySlug = (slug: string) => {
  const prod = productOfferings.value.find(p => p.slug === slug)
  return prod?.image || ''
}

onUnmounted(() => {
  if (intervalTimer) clearInterval(intervalTimer)
})
</script>

<template>
  <div class="w-full bg-white-base text-brand font-brand">
    
    <div class="relative min-h-screen flex items-center justify-center overflow-hidden">
      <div class="absolute inset-0 z-0">
        <img 
          v-for="(bg, index) in backgrounds" 
          :key="index"
          :src="bg" 
          alt="Florist Background" 
          :class="[
            'absolute inset-0 w-full h-full object-cover transition-opacity duration-1000 ease-in-out',
            currentIndex === index ? 'opacity-100 z-10' : 'opacity-0 z-0'
          ]" 
        />
        <div class="absolute inset-0 bg-black/40 z-20"></div>
      </div>

      <div class="relative z-30 text-center px-4 max-w-4xl mx-auto mt-20">
        <h1 class="text-5xl md:text-7xl font-bold text-white mb-6 leading-tight drop-shadow-lg">
          Beautiful Flower Boards
          <span class="block">for Every Special Moment</span>
        </h1>
        
        <p class="text-lg md:text-xl text-gray-200 mb-10 max-w-2xl mx-auto drop-shadow-md leading-relaxed">
          We provide customized flower boards for weddings, condolences, graduations, and corporate events. Designed with care and delivered on time to make every moment memorable.
        </p>
        
        <div class="mt-10 flex flex-wrap justify-center gap-4">
          <NuxtLink 
            to="/catalog" 
            class="bg-[#1b4332] hover:bg-[#143326] text-white font-bold py-3 px-8 rounded-xl transition shadow-md text-sm text-center"
          >
            View Products
          </NuxtLink>

          <NuxtLink 
            to="/products/custom" 
            class="border-2 border-white hover:bg-white/10 text-white font-bold py-3 px-8 rounded-xl transition text-sm text-center flex items-center justify-center gap-2"
          >
            <span>Design Custom Board 🌸</span>
          </NuxtLink>
        </div>
      </div>

      <div class="absolute right-8 top-1/2 transform -translate-y-1/2 flex flex-col gap-6 z-30 hidden lg:flex">
        <span 
          v-for="(bg, index) in backgrounds" 
          :key="index"
          @click="changeBg(index)"
          :class="[
            'font-bold cursor-pointer transition-all duration-300 transform',
            currentIndex === index 
              ? 'text-yellow-400 text-lg scale-110' 
              : 'text-white/60 hover:text-white text-sm hover:scale-105'
          ]"
        >
          0{{ index + 1 }}
        </span>
      </div>
    </div>

    <section class="max-w-7xl mx-auto px-8 py-24 grid grid-cols-1 md:grid-cols-2 gap-12 items-start">
      <div>
        <h2 class="text-3xl md:text-4xl font-bold text-accent leading-tight max-w-sm">
          We Help Choose The Perfect Flower Board
        </h2>
      </div>
      <div>
        <p class="text-gray-600 text-sm md:text-base leading-relaxed">
          Our flower boards are designed to suit various occasions, from joyful celebrations to heartfelt condolences. Each arrangement is crafted to deliver your message with elegance and meaning.
        </p>
      </div>
    </section>

    <section class="max-w-7xl mx-auto px-8 pb-24 grid grid-cols-1 md:grid-cols-3 gap-8">
      <div class="bg-white-base border border-gray-100 p-8 rounded-lg shadow-sm hover:shadow-md transition">
        <div class="text-accent mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-4 text-gray-800">Wedding Flower Board</h3>
        <p class="text-gray-500 text-sm leading-relaxed">Bring the beauty of nature to your outdoor spaces with our wide selection of outdoor plants</p>
      </div>

      <div class="bg-accent text-white-base p-8 rounded-lg shadow-md transition transform hover:scale-105">
        <div class="text-white-base mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-4">Condolence Flower Board</h3>
        <p class="text-gray-200 text-sm leading-relaxed">Bring a touch of greenery to your living spaces with our collection of indoor plants, perfect for purifying the air and adding a natural touch to your home.</p>
      </div>

      <div class="bg-white-base border border-gray-100 p-8 rounded-lg shadow-sm hover:shadow-md transition">
        <div class="text-accent mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 7v10M20 7v10M12 4v16m0-4h4m-4-8h4m-4 4H8m0-4h4" />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-4 text-gray-800">Celebration Flower Board</h3>
        <p class="text-gray-500 text-sm leading-relaxed">Add a touch of style to your indoor or outdoor spaces with our collection of pots plants, available in a variety of sizes and designs to fit any decor</p>
      </div>
    </section>

    <section class="max-w-7xl mx-auto px-8 py-20">
      <h2 class="text-3xl md:text-4xl font-bold text-center text-accent mb-12">What We Offer To You</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        
        <div v-for="(item, idx) in productOfferings" :key="idx" class="bg-white-base rounded-xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-md transition flex flex-col justify-between cursor-pointer" @click="navigateTo(item.id === 'custom' || item.slug === 'custom' ? '/products/custom' : `/products/${item.slug || item.id}`)">
          <div class="h-64 relative bg-gray-50">
            <img :src="item.image || '/images/custom-preview.png'" :alt="item.name" class="w-full h-full object-cover" />
            <span v-if="item.id === 'custom'" class="absolute top-3 left-3 bg-[#1b4332] text-white text-[10px] font-black tracking-widest uppercase px-2.5 py-1 rounded-lg border border-[#143326] shadow-sm z-20">
              Interactive
            </span>
            <span v-if="!item.isAvailable" class="absolute top-3 right-3 bg-red-100 text-red-800 text-[10px] font-black tracking-widest uppercase px-2.5 py-1 rounded-lg border border-red-200 shadow-sm z-20">
              Sold Out
            </span>
          </div>
          <div class="p-4 flex flex-col gap-3">
            <h3 class="font-bold text-gray-800">{{ item.name }}</h3>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-gray-400 text-xs line-through" v-if="item.id !== 'custom'">{{ formatRupiah(item.price + 35000) }}</span>
                <span class="text-accent font-bold text-base">{{ formatRupiah(item.price) }}</span>
              </div>
              <NuxtLink 
                v-if="item.isAvailable" 
                :to="item.id === 'custom' || item.slug === 'custom' ? '/products/custom' : `/products/${item.slug || item.id}`" 
                class="bg-accent text-white-base text-xs px-3 py-1.5 rounded hover:bg-accent/90 transition font-bold"
                @click.stop
              >
                {{ item.id === 'custom' || item.slug === 'custom' ? 'Design Now' : 'Buy' }}
              </NuxtLink>
              <button 
                v-else 
                disabled 
                class="bg-gray-100 text-gray-400 text-xs px-3 py-1.5 rounded cursor-not-allowed border border-gray-200"
              >
                Sold Out
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <section class="max-w-7xl mx-auto px-8 py-20 grid grid-cols-1 md:grid-cols-2 gap-16 items-center">
      <div class="h-[450px]">
        <img src="/images/florist.jpg" alt="Gallery preview" class="w-full h-full object-cover rounded-xl shadow-sm" />
      </div>
      <div class="grid grid-cols-2 gap-8">
        
        <div class="space-y-2">
          <div class="text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-800">Quality Product</h4>
          <p class="text-xs text-gray-500 leading-relaxed">Our flowers are of the highest quality, carefully selected and sourced from reputable</p>
        </div>

        <div class="space-y-2">
          <div class="text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-800">Always Fresh</h4>
          <p class="text-xs text-gray-500 leading-relaxed">Our flowers are always fresh, handpicked and delivered promptly for maximum longevity and enjoyment.</p>
        </div>

        <div class="space-y-2">
          <div class="text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-800">Work Smart</h4>
          <p class="text-xs text-gray-500 leading-relaxed">We work smart, using innovative techniques and technology to streamline our processes</p>
        </div>

        <div class="space-y-2">
          <div class="text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.121 14.121L19 19m-4.879-4.879l-1.414-1.414M12 10.586l-1.414-1.414M9.879 9.879L5 5m4.879 4.879l1.414 1.414M12 10.586l1.414 1.414M9.879 9.879L5 19" />
            </svg>
          </div>
          <h4 class="font-bold text-gray-800">Excellent Service</h4>
          <p class="text-xs text-gray-500 leading-relaxed">We pride ourselves on providing excellent service, going above and beyond to meet our customers' needs</p>
        </div>

      </div>
    </section>

    <!-- <section class="max-w-7xl mx-auto px-8 py-20">
      <h2 class="text-3xl md:text-4xl font-bold text-center text-accent mb-12">Our Gallery View</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

        <div class="h-[600px] rounded-xl overflow-hidden shadow-sm relative group bg-gray-200 flex items-center justify-center">
          <img 
            v-if="getProductImageBySlug('graduate')" 
            :src="getProductImageBySlug('graduate')" 
            alt="Graduate Gallery" 
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
          />
          <div v-else class="text-center text-gray-400 p-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span class="text-xs font-black uppercase tracking-wider block">Graduate</span>
          </div>
        </div>

        <div class="md:col-span-2 grid grid-cols-2 grid-rows-2 gap-6 h-[600px]">
          <div class="rounded-xl overflow-hidden shadow-sm relative group bg-gray-200 flex items-center justify-center">
            <img 
              v-if="getProductImageBySlug('wedding')" 
              :src="getProductImageBySlug('wedding')" 
              alt="Wedding Gallery" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
            />
            <div v-else class="text-center text-gray-400 p-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="text-[10px] font-black uppercase tracking-wider block">Wedding</span>
            </div>
          </div>

          <div class="rounded-xl overflow-hidden shadow-sm relative group bg-gray-200 flex items-center justify-center">
            <img 
              v-if="getProductImageBySlug('birthday')" 
              :src="getProductImageBySlug('birthday')" 
              alt="Birthday Gallery" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
            />
            <div v-else class="text-center text-gray-400 p-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="text-[10px] font-black uppercase tracking-wider block">Birthday</span>
            </div>
          </div>

          <div class="rounded-xl overflow-hidden shadow-sm relative group bg-gray-200 flex items-center justify-center">
            <img 
              v-if="getProductImageBySlug('anniversary')" 
              :src="getProductImageBySlug('anniversary')" 
              alt="Anniversary Gallery" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
            />
            <div v-else class="text-center text-gray-400 p-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="text-[10px] font-black uppercase tracking-wider block">Anniversary</span>
            </div>
          </div>

          <div class="rounded-xl overflow-hidden shadow-sm relative group bg-gray-200 flex items-center justify-center">
            <img 
              v-if="getProductImageBySlug('condolences')" 
              :src="getProductImageBySlug('condolences')" 
              alt="Condolences Gallery" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
            />
            <div v-else class="text-center text-gray-400 p-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span class="text-[10px] font-black uppercase tracking-wider block">Condolences</span>
            </div>
          </div>
        </div>
      </div>
    </section> -->

    <section class="max-w-7xl mx-auto px-8 py-20">
      <h2 class="text-3xl md:text-4xl font-bold text-center text-accent mb-12">What Do They Say About Us</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        
        <div class="bg-white-base p-8 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-4 mb-6">
            <img src="/images/ilham.jpeg" alt="Ilham" class="w-12 h-12 rounded-full object-cover" />
            <h4 class="font-bold text-gray-800">Ilham Priambodo</h4>
          </div>
          <p class="text-gray-500 text-sm leading-relaxed">"Highly recommend this website for quality flowers and plants. Great prices, timely delivery and excellent customer service."</p>
        </div>

        <div class="bg-white-base p-8 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-4 mb-6">
            <img src="/images/rafata.jpeg" alt="Rafata" class="w-12 h-12 rounded-full object-cover" />
            <h4 class="font-bold text-gray-800">Rafata Alfatih</h4>
          </div>
          <p class="text-gray-500 text-sm leading-relaxed">"Great service, beautiful flowers, timely delivery. Highly recommend."</p>
        </div>

        <div class="bg-white-base p-8 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-4 mb-6">
            <img src="/images/rayhan.jpeg" alt="Rayhan" class="w-12 h-12 rounded-full object-cover" />
            <h4 class="font-bold text-gray-800">Rayhan Shidqi</h4>
          </div>
          <p class="text-gray-500 text-sm leading-relaxed">"I am very happy with my purchase from this website, the plants were healthy and arrived on time."</p>
        </div>

      </div>
    </section>

  </div>
</template>

<style scoped></style>