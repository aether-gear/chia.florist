<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useCart } from '~/composables/useCart' // 1. Impor Composable Global Cart

// Matikan layout bawaan agar footer/navbar web tidak muncul sama sekali
definePageMeta({
  layout: false
})

useHead({
  title: 'Chia Florist - Board Simulator'
})

// Ambil fungsi penambahan item keranjang dan helper formatter Rupiah dari store
const { addToCart, formatRupiah } = useCart()

// --- STATE NAVIGASI SIMULATOR ---
const isGameStarted = ref(false)
const hasSeenTutorial = ref(false)

// Interface TypeScript
interface SizeOption {
  id: string
  label: string
  price: number
  desc: string
  class: string
}

interface ThemeOption {
  id: string
  label: string
  color: string
  price: number
}

interface FlowerOption {
  id: string
  label: string
  price: number
  icon: string
}

interface FontOption {
  id: string
  label: string
  family: string
  price: number
}

interface CustomSelection {
  size: SizeOption
  theme: ThemeOption
  flower: FlowerOption
  font: FontOption
  text: {
    header: string
    target: string
    sender: string
  }
}

// Data Master Simulator + Rekomendasi (Dikonversi Menjadi Rupiah Murni)
const options = {
  sizes: [
    { id: 'small', label: '1.5m x 2.0m', price: 150000, desc: 'Compact size', class: 'max-w-sm md:max-w-md' },
    { id: 'medium', label: '1.8m x 2.5m', price: 200000, desc: 'Most popular', class: 'max-w-md md:max-w-lg' },
    { id: 'large', label: '2.0m x 3.0m', price: 280000, desc: 'Premium Grand', class: 'max-w-lg md:max-w-xl' }
  ] as SizeOption[],
  themes: [
    { id: 'emerald', label: 'Florist Emerald', color: '#114028', price: 0 },
    { id: 'navy', label: 'Royal Navy', color: '#0c1b30', price: 0 },
    { id: 'maroon', label: 'Wine Maroon', color: '#4c0519', price: 25000 },
    { id: 'black', label: 'Luxury Black', color: '#111111', price: 15000 }
  ] as ThemeOption[],
  flowers: [
    { id: 'basic', label: 'Minimalist (Top)', price: 40000, icon: '🌸' },
    { id: 'standard', label: 'Corner Duo', price: 70000, icon: '🌸🌸' },
    { id: 'luxury', label: 'Full Frame', price: 125000, icon: '👑' }
  ] as FlowerOption[],
  fonts: [
    { id: 'serif', label: 'Classic Luxury', family: 'font-serif', price: 15000 },
    { id: 'modern', label: 'Modern Sans', family: 'font-sans', price: 0 },
    { id: 'script', label: 'Elegant Script', family: 'italic font-serif', price: 20000 }
  ] as FontOption[]
}

// State Default Pilihan User
const currentStep = ref(1)
const selection = ref<CustomSelection>({
  size: options.sizes[0]!,
  theme: options.themes[0]!,
  flower: options.flowers[0]!,
  font: options.fonts[0]!,
  text: {
    header: 'HAPPY WEDDING',
    target: 'Ilham & Belle',
    sender: 'AETHER GEAR'
  }
})

// Kalkulasi Harga Otomatis Secara Numerik
const totalPrice = computed(() => {
  return selection.value.size.price +
         selection.value.theme.price +
         selection.value.flower.price +
         selection.value.font.price
})

// LOGIKA AKHIR SIMULATOR KE KERANJANG GLOBAL
const handleCustomAddToCart = () => {
  addToCart({
    id: 'custom-' + Date.now(), // ID Unik agar tidak saling menimpa
    name: `Custom Board (${selection.value.text.header || 'Design'})`,
    price: totalPrice.value,
    image: '/images/custom-preview.png', // Fallback gambar preview simulator
    size: selection.value.size.label,
    color: selection.value.theme.color,
    shopId: '333f6432-a01c-412f-99f4-0f08ca0d8eb1', // Pass default shop ID for custom products
    isCustom: true
  }, 1)

  // Redirect langsung menuju ke halaman keranjang belanjaan
  navigateTo('/cart')
}

// Logic Auto Switch Background Hero
const currentIndex = ref(0)
let intervalTimer: any = null

const backgrounds = [
  '/florist.jpg',
  '/flowerist.jpg',
  '/flower.jpg'
]

const changeBg = (index: number) => {
  currentIndex.value = index
  resetTimer() 
}

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

// Data Rekomendasi Tips Langkah
const sizeTip = computed(() => selection.value.size.id === 'small' ? 'Sangat cocok untuk ucapan di ruangan terbatas (indoor).' : selection.value.size.id === 'medium' ? 'Rekomendasi paling pas untuk acara pernikahan/wedding.' : 'Sangat megah untuk acara Grand Opening perusahaan.')
const themeTip = computed(() => selection.value.theme.id === 'emerald' ? 'Melambangkan harapan dan kebahagiaan universal.' : selection.value.theme.id === 'navy' ? 'Sangat elegan dan profesional untuk acara korporat.' : selection.value.theme.id === 'maroon' ? 'Mewah dan romantis untuk ulang tahun pernikahan.' : 'Sangat tepat untuk duka cita/condolences.')
const flowerTip = computed(() => selection.value.flower.id === 'basic' ? 'Simpel namun tetap sopan untuk acara santai.' : selection.value.flower.id === 'standard' ? 'Pilihan seimbang dan estetik di dua sudut papan.' : 'Tampak sangat mewah dengan bunga yang mengelilingi papan.')
const fontTip = computed(() => selection.value.font.id === 'serif' ? 'Font resmi yang memberikan kesan megah.' : selection.value.font.id === 'modern' ? 'Terbaca jelas dan cocok untuk segala jenis ucapan.' : 'Gaya tulisan tangan yang romantis dan personal.')

const startGame = () => { isGameStarted.value = true }
const completeTutorial = () => { hasSeenTutorial.value = true }
const goToStep = (step: number) => { currentStep.value = step }
const next = () => { if (currentStep.value < 4) currentStep.value++ }
const prev = () => { if (currentStep.value > 1) currentStep.value-- }
</script>

<template>
  <div class="h-screen w-full flex flex-col bg-[#0b2518] text-white overflow-hidden select-none font-sans">
    
    <div v-if="!isGameStarted" class="flex-1 flex flex-col items-center justify-center p-6 text-center bg-gradient-to-b from-[#0b2518] to-[#04120c] relative animate-in">
      <div class="max-w-2xl px-4 space-y-6">
        <div class="inline-block bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-black px-4 py-1.5 rounded-full tracking-widest uppercase mb-2 animate-pulse">
          🎮 Interactive Experience
        </div>
        <h1 class="text-4xl md:text-6xl font-black tracking-tighter text-white leading-none drop-shadow-lg">
          CHIA FLORIST <span class="text-yellow-400 block mt-2">BOARD SIMULATOR</span>
        </h1>
        <p class="text-sm md:text-base text-slate-300 font-medium leading-relaxed max-w-xl mx-auto">
          Selamat datang di simulator kustomisasi papan bunga interaktif. Di sini kamu bisa merancang, memilih tema warna, menata bunga hias, dan memasukkan tulisan kustom sesuai keinginanmu secara real-time.
        </p>

        <div class="pt-6">
          <button 
            @click="startGame"
            class="group bg-emerald-500 hover:bg-emerald-400 text-black text-base font-black px-10 py-4 rounded-2xl transition-all duration-300 shadow-[0_0_30px_rgba(16,185,129,0.3)] hover:shadow-[0_0_40px_rgba(16,185,129,0.5)] flex items-center gap-3 mx-auto hover:scale-105 active:scale-95 cursor-pointer"
          >
            <span>START GAME</span>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 transform transition group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M13 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>
      <div class="absolute bottom-6 right-6 text-[10px] font-bold opacity-40">v1.0.6 - Chia Dev</div>
    </div>

    <div v-else-if="!hasSeenTutorial" class="flex-1 flex flex-col items-center justify-center p-6 text-center bg-gradient-to-b from-[#0b2518] to-[#04120c] relative animate-in">
      <div class="max-w-2xl px-4 space-y-8">
        <div>
          <h2 class="text-3xl font-black tracking-tight text-white mb-2">💡 CARA BERMAIN</h2>
          <p class="text-xs text-emerald-300 uppercase tracking-widest font-bold">Langkah mudah untuk merakit papan bungamu</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-left max-w-xl mx-auto select-none">
          <div class="bg-white/5 border border-white/10 p-4 rounded-xl flex gap-3 items-start">
            <span class="text-xl">📏</span>
            <div>
              <h4 class="font-bold text-sm text-yellow-400">1. Ukuran Papan</h4>
              <p class="text-xs text-slate-300">Pilih skala dimensi papan yang sesuai dengan lokasi acaramu.</p>
            </div>
          </div>
          <div class="bg-white/5 border border-white/10 p-4 rounded-xl flex gap-3 items-start">
            <span class="text-xl">🎨</span>
            <div>
              <h4 class="font-bold text-sm text-yellow-400">2. Tema Warna</h4>
              <p class="text-xs text-slate-300">Pilih warna kain styrofoam latar belakang papan.</p>
            </div>
          </div>
          <div class="bg-white/5 border border-white/10 p-4 rounded-xl flex gap-3 items-start">
            <span class="text-xl">💐</span>
            <div>
              <h4 class="font-bold text-sm text-yellow-400">3. Bunga & Gaya Font</h4>
              <p class="text-xs text-slate-300">Pasang rangkaian bunga dan gaya huruf yang paling cocok.</p>
            </div>
          </div>
          <div class="bg-white/5 border border-white/10 p-4 rounded-xl flex gap-3 items-start">
            <span class="text-xl">✍️</span>
            <div>
              <h4 class="font-bold text-sm text-yellow-400">4. Grafir Pesan</h4>
              <p class="text-xs text-slate-300">Ketik pesan ucapan lengkap, nama penerima, dan pengirim.</p>
            </div>
          </div>
        </div>

        <div class="pt-4">
          <button @click="completeTutorial" class="bg-emerald-500 hover:bg-emerald-400 text-black text-sm font-black px-8 py-3.5 rounded-xl transition shadow-lg hover:scale-105 cursor-pointer">
            MULAI SIMULASI 🚀
          </button>
        </div>
      </div>
    </div>

    <div v-else class="flex-1 flex flex-col overflow-hidden animate-in">
      
      <div class="h-16 bg-[#041a10] border-b border-white/10 flex items-center justify-between px-6 flex-shrink-0 z-50">
        <NuxtLink to="/" class="flex items-center gap-2 bg-white/5 hover:bg-white/10 px-4 py-2 rounded-xl border border-white/10 backdrop-blur-md transition-all text-xs font-bold text-white hover:scale-105">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to Home
        </NuxtLink>
        <h1 class="text-sm font-black text-emerald-400 tracking-wider hidden md:block">CHIA FLORIST - BOARD BUILDER</h1>
        <div class="text-xs font-bold text-white/40 bg-white/5 px-3 py-1.5 rounded-full border border-white/5">v1.0.6</div>
      </div>

      <div class="flex-1 flex flex-col lg:flex-row overflow-hidden bg-[#072215]">
        
        <div class="flex-1 relative flex items-center justify-center p-6 lg:p-12 bg-[#051a0f] border-b lg:border-b-0 lg:border-r border-white/5 min-h-[350px]">
          
          <div class="absolute top-6 left-6 flex items-center gap-2 bg-black/40 px-3 py-1.5 rounded-full border border-white/10 backdrop-blur-md z-30">
            <span class="flex h-2 w-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span class="text-[10px] font-bold uppercase tracking-wider text-emerald-300">Live Simulation View</span>
          </div>

          <div 
            id="board-canvas"
            :style="{ backgroundColor: selection.theme.color }"
            :class="[
              'relative aspect-[4/3] w-full flex flex-col justify-between p-6 md:p-10 text-center border-[12px] border-[#8b5a2b] transition-all duration-500 shadow-[0_20px_50px_rgba(0,0,0,0.6)] rounded-sm select-none overflow-hidden max-h-[55vh]',
              selection.size.class,
              selection.font.family
            ]"
          >
            <div v-if="selection.flower.id === 'basic' || selection.flower.id === 'luxury'" class="absolute -top-6 left-0 right-0 flex justify-around text-3xl md:text-4xl pointer-events-none z-10 select-none">
              🌸🌸🌸🌸🌸🌸🌸🌸
            </div>
            <div v-if="selection.flower.id === 'standard' || selection.flower.id === 'luxury'" class="absolute -top-3 -left-3 text-5xl md:text-6xl pointer-events-none z-10 select-none">🌺</div>
            <div v-if="selection.flower.id === 'standard' || selection.flower.id === 'luxury'" class="absolute -bottom-3 -right-3 text-5xl md:text-6xl pointer-events-none z-10 select-none">🌺</div>
            <div v-if="selection.flower.id === 'luxury'" class="absolute inset-y-0 -left-6 flex flex-col justify-around text-xl md:text-2xl pointer-events-none z-10 select-none">🌸<br>🌸<br>🌸</div>
            <div v-if="selection.flower.id === 'luxury'" class="absolute inset-y-0 -right-6 flex flex-col justify-around text-xl md:text-2xl pointer-events-none z-10 select-none">🌸<br>🌸<br>🌸</div>

            <div class="h-full w-full flex flex-col justify-around text-center break-words pointer-events-none px-2 z-10">
              <h2 class="text-xl md:text-2xl font-black uppercase tracking-widest text-white drop-shadow-md leading-tight">
                {{ selection.text.header || 'Header' }}
              </h2>
              <div class="w-full">
                <p class="text-2xl md:text-4xl font-extrabold text-yellow-300 drop-shadow-md italic leading-snug">
                  {{ selection.text.target || 'Nama Penerima' }}
                </p>
              </div>
              <div class="border-t border-white/20 pt-2 px-4 w-full">
                <p class="text-sm md:text-base font-medium text-slate-100">
                  {{ selection.text.sender || 'Nama Pengirim' }}
                </p>
              </div>
            </div>
          </div>

          <div class="absolute bottom-6 right-6 flex items-center gap-4 bg-black/80 p-4 rounded-xl border border-white/10 backdrop-blur-md">
            <div class="w-10 h-10 bg-emerald-500 rounded-full flex items-center justify-center text-black shadow">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.407 2.63 1m-2.63-1V7m0 1v1m0 5v1m0-1c-1.11 0-2.08-.407-2.63-1m2.63 1v1m0-1V14" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-tighter">Estimasi Harga</p>
              <p class="text-2xl font-black text-emerald-400">{{ formatRupiah(totalPrice) }}</p>
            </div>
          </div>
        </div>

        <div class="w-full lg:w-[450px] bg-[#03140c] border-l border-white/5 flex flex-col justify-between h-full flex-shrink-0">
          
          <div class="p-6 bg-[#061e11] border-b border-white/5">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xs font-black tracking-widest text-slate-400 uppercase">Simulator Options</h2>
              <span class="text-xs font-bold text-emerald-400 bg-emerald-500/10 px-3 py-1 rounded-full border border-emerald-500/20">Step {{ currentStep }} / 4</span>
            </div>
            <div class="flex gap-1">
              <button 
                v-for="i in 4" :key="i"
                @click="goToStep(i)"
                :class="['h-1.5 flex-1 rounded-full transition-all duration-300 cursor-pointer', currentStep === i ? 'bg-emerald-400 scale-102' : currentStep > i ? 'bg-slate-700' : 'bg-slate-800']"
              ></button>
            </div>
          </div>

          <div class="flex-1 overflow-y-auto p-6 space-y-6">
            
            <div v-if="currentStep === 1" class="space-y-4 animate-in">
              <div>
                <h3 class="text-base font-bold text-slate-100">1. Board Size</h3>
                <p class="text-xs text-slate-400">Pilih dimensi papan sesuai kebutuhan acaramu.</p>
              </div>
              <div class="grid grid-cols-1 gap-3">
                <button 
                  v-for="s in options.sizes" :key="s.id"
                  @click="selection.size = s"
                  :class="['flex flex-col gap-1 p-4 rounded-xl border-2 transition text-left cursor-pointer', selection.size.id === s.id ? 'border-emerald-500 bg-emerald-500/10' : 'border-white/5 bg-white/5 hover:border-white/20']"
                >
                  <div class="flex justify-between items-center w-full">
                    <span class="font-bold text-slate-100 text-sm">{{ s.label }} ({{ s.desc }})</span>
                    <span :class="['font-black text-sm', selection.size.id === s.id ? 'text-emerald-400' : 'text-slate-300']">{{ formatRupiah(s.price) }}</span>
                  </div>
                  <p class="text-[11px] text-emerald-300 mt-1 leading-normal font-medium">💡 {{ sizeTip }}</p>
                </button>
              </div>
            </div>

            <div v-if="currentStep === 2" class="space-y-4 animate-in">
              <div>
                <h3 class="text-base font-bold text-slate-100">2. Foam Theme</h3>
                <p class="text-xs text-slate-400">Tentukan warna dasar latar belakang untuk nuansa papan bunga.</p>
              </div>
              <div class="grid grid-cols-1 gap-3">
                <button 
                  v-for="c in options.themes" :key="c.id"
                  @click="selection.theme = c"
                  :class="['p-4 rounded-xl border-2 transition-all flex flex-col gap-2 text-left cursor-pointer', selection.theme.id === c.id ? 'border-emerald-500 bg-emerald-500/10' : 'border-white/5 bg-white/5 hover:border-white/20']"
                >
                  <div class="flex justify-between items-center w-full">
                    <div class="flex items-center gap-3">
                      <div :style="{ backgroundColor: c.color }" class="w-8 h-6 rounded border border-white/10 shadow-md"></div>
                      <span class="text-xs font-bold text-slate-100">{{ c.label }}</span>
                    </div>
                    <span class="text-xs font-bold text-emerald-400">{{ c.price === 0 ? 'FREE' : `+${formatRupiah(c.price)}` }}</span>
                  </div>
                  <p class="text-[11px] text-emerald-300 mt-1 leading-normal font-medium">💡 {{ themeTip }}</p>
                </button>
              </div>
            </div>

            <div v-if="currentStep === 3" class="space-y-6 animate-in">
              <div class="space-y-3">
                <div>
                  <h3 class="text-sm font-bold text-slate-100">3a. Flowers Layout</h3>
                  <p class="text-xs text-slate-400">Pilih penempatan bunga hias.</p>
                </div>
                <div class="grid grid-cols-1 gap-2">
                  <button 
                    v-for="f in options.flowers" :key="f.id"
                    @click="selection.flower = f"
                    :class="['flex flex-col gap-1 p-3.5 rounded-xl border-2 transition text-left cursor-pointer', selection.flower.id === f.id ? 'border-emerald-500 bg-emerald-500/10 text-emerald-300' : 'border-white/5 bg-white/5 hover:border-white/20']"
                  >
                    <div class="flex justify-between items-center w-full">
                      <span class="text-xs font-bold text-slate-100">{{ f.icon }} {{ f.label }}</span>
                      <span class="text-xs font-bold text-emerald-400">+{{ formatRupiah(f.price) }}</span>
                    </div>
                    <p class="text-[11px] text-emerald-200 mt-1 leading-normal font-medium">💡 {{ flowerTip }}</p>
                  </button>
                </div>
              </div>

              <div class="space-y-3">
                <div>
                  <h3 class="text-sm font-bold text-slate-100">3b. Font Style</h3>
                  <p class="text-xs text-slate-400">Tentukan bentuk huruf teks pesan.</p>
                </div>
                <div class="grid grid-cols-1 gap-2">
                  <button 
                    v-for="font in options.fonts" :key="font.id"
                    @click="selection.font = font"
                    :class="['flex flex-col gap-1 p-3.5 rounded-xl border-2 transition text-left cursor-pointer', selection.font.id === font.id ? 'border-emerald-500 bg-emerald-500/10 text-emerald-300' : 'border-white/5 bg-white/5 hover:border-white/20']"
                  >
                    <div class="flex justify-between items-center w-full">
                      <span :class="['text-xs italic tracking-wider text-slate-100', font.family]">Aa {{ font.label }}</span>
                      <span class="text-xs font-bold text-emerald-400">+{{ formatRupiah(font.price) }}</span>
                    </div>
                    <p class="text-[11px] text-emerald-200 mt-1 leading-normal font-medium">💡 {{ fontTip }}</p>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="currentStep === 4" class="space-y-5 animate-in">
              <div>
                <h3 class="text-base font-bold text-slate-100">4. Tulis Grafir Pesan</h3>
                <p class="text-xs text-slate-400">Masukkan kalimat ucapan untuk dicetak pada papan.</p>
              </div>
              <div class="space-y-4">
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-emerald-400 uppercase tracking-widest">Ucapan Utama</label>
                  <input v-model="selection.text.header" type="text" class="w-full bg-white/5 border border-white/10 p-3 rounded-xl focus:border-emerald-500 outline-none transition text-sm text-white" placeholder="Contoh: HAPPY WEDDING" />
                </div>
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-emerald-400 uppercase tracking-widest">Nama Penerima</label>
                  <input v-model="selection.text.target" type="text" class="w-full bg-white/5 border border-white/10 p-3 rounded-xl focus:border-emerald-500 outline-none transition text-sm text-white" placeholder="Contoh: Sita & Panca" />
                </div>
                <div class="space-y-1">
                  <label class="text-[10px] font-black text-emerald-400 uppercase tracking-widest">Nama Pengirim</label>
                  <input v-model="selection.text.sender" type="text" class="w-full bg-white/5 border border-white/10 p-3 rounded-xl focus:border-emerald-500 outline-none transition text-sm text-white" placeholder="Contoh: PT. Maju Bersama" />
                </div>
              </div>
            </div>

          </div>

          <div class="p-6 bg-black/40 border-t border-white/10 grid grid-cols-2 gap-4 flex-shrink-0">
            <button 
              @click="prev" 
              :disabled="currentStep === 1"
              class="py-4 rounded-xl font-bold text-sm bg-white/5 hover:bg-white/10 disabled:opacity-30 transition border border-white/5 cursor-pointer"
            >
              BACK
            </button>

            <button 
              v-if="currentStep < 4"
              @click="next"
              class="py-4 rounded-xl font-bold text-sm bg-emerald-500 text-black hover:bg-emerald-400 transition shadow-[0_0_20px_rgba(16,185,129,0.3)] cursor-pointer"
            >
              NEXT STEP
            </button>

            <button 
              v-else
              @click="handleCustomAddToCart"
              class="py-4 rounded-xl font-bold text-sm bg-yellow-500 text-black hover:bg-yellow-400 transition shadow-[0_0_20px_rgba(234,179,8,0.3)] cursor-pointer"
            >
              ADD TO CART
            </button>
          </div>

        </div>

      </div>
    </div>

  </div>
</template>

<style>
@import "tailwindcss";

.animate-in {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateX(12px); }
  to { opacity: 1; transform: translateX(0); }
}

::-webkit-scrollbar {
  width: 5px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #234b35;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #2f6145;
}
</style>