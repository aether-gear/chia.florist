<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

useHead({
  title: 'Chia Florist - Flower Boards',
  meta: [
    { name: 'description', content: 'We provide customized flower boards for weddings, condolences, graduations, and corporate events.' }
  ]
})

// 1. Array gambar background Hero (Sesuaikan dengan file yang kamu punya di folder public)
const backgrounds = [
  '/florist.jpg',          // Gambar 01
  '/flowerist.jpg',          // Gambar 02
  '/flower.jpg'            // Gambar 03
]

// State untuk indeks background yang sedang aktif
const currentIndex = ref(0)
let intervalTimer: any = null

// 2. Fungsi untuk mengganti background saat diklik
const changeBg = (index: number) => {
  currentIndex.value = index
  resetTimer() 
}

// 3. Logic Auto Switch setiap 5 detik
const startTimer = () => {
  intervalTimer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % backgrounds.length
  }, 5000)
}

const resetTimer = () => {
  if (intervalTimer) clearInterval(intervalTimer)
  startTimer()
}

onMounted(() => {
  startTimer()
})

onUnmounted(() => {
  if (intervalTimer) clearInterval(intervalTimer)
})

// Data Produk (Diberikan id agar bisa di-redirect secara dinamis)
const productOfferings = ref([
  { id: 'wedding', name: 'Wedding', image: '/images/wedding.jpeg' },
  { id: 'congratulations', name: 'Congratulation', image: '/images/congratulations.jpeg' },
  { id: 'condolences', name: 'Condolences', image: '/images/condolences.jpeg' },
  { id: 'grand-opening', name: 'Grand Opening', image: '/images/grandop.jpeg' },
  { id: 'birthday', name: 'Birthday', image: '/images/birthday.jpeg' },
  { id: 'graduate', name: 'Graduate', image: '/images/graduate.jpeg' },
  { id: 'anniversary', name: 'Anniversary', image: '/images/anniversary.jpeg' },
  { id: 'custom', name: 'Custom', image: '' }
])
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
    Order Now
  </NuxtLink>

  <NuxtLink 
    to="/catalog" 
    class="border-2 border-white hover:bg-white/10 text-white font-bold py-3 px-8 rounded-xl transition text-sm text-center flex items-center justify-center gap-2"
  >
    <span>View Products</span>
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
        
        <div v-for="(item, idx) in productOfferings" :key="idx" class="bg-white-base rounded-xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-md transition">
          <div class="h-64 relative bg-gray-50">
            <div v-if="item.name === 'Custom'" class="w-full h-full bg-accent flex items-center justify-center">
              <span class="text-white-base text-4xl font-bold">?</span>
            </div>
            <img v-else :src="item.image" :alt="item.name" class="w-full h-full object-cover" />
          </div>
          <div class="p-4 flex flex-col gap-3">
            <h3 class="font-bold text-gray-800">{{ item.name }}</h3>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-gray-400 text-xs line-through">$10</span>
                <span class="text-accent font-bold text-base">$8</span>
              </div>
              <NuxtLink :to="`/products/${item.id}`" class="bg-accent text-white-base text-xs px-3 py-1.5 rounded hover-bg-accent-strong transition">
                Buy
              </NuxtLink>
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

    <section class="max-w-7xl mx-auto px-8 py-20">
      <h2 class="text-3xl md:text-4xl font-bold text-center text-accent mb-12">Our Gallery View</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="h-[600px]">
          <img src="/images/graduate.jpeg" alt="Gallery 1" class="w-full h-full object-cover rounded-xl shadow-sm" />
        </div>
        <div class="md:col-span-2 grid grid-cols-2 grid-rows-2 gap-6 h-[600px]">
          <img src="/images/wedding.jpeg" alt="Gallery 2" class="w-full h-full object-cover rounded-xl" />
          <img src="/images/birthday.jpeg" alt="Gallery 3" class="w-full h-full object-cover rounded-xl" />
          <img src="/images/anniversary.jpeg" alt="Gallery 4" class="w-full h-full object-cover rounded-xl" />
          <img src="/images/condolences.jpeg" alt="Gallery 5" class="w-full h-full object-cover rounded-xl" />
        </div>
      </div>
    </section>

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
