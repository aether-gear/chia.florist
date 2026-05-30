<script setup lang="ts">
import { ref } from 'vue'

useHead({
  title: 'Our Collection - Chia Florist',
  meta: [
    { name: 'description', content: 'Explore our premium selection of pre-designed flower boards or launch our custom game simulator.' }
  ]
})

// Data dummy katalog disesuaikan dengan ID rute di dalam folder products
const products = ref([
  {
    id: 'birthday',
    name: 'Birthday Flower Board',
    price: 192.00,
    rating: 4.5,
    reviews: 150,
    image: '/images/birthday.jpeg',
    tag: 'Celebration',
    desc: 'Bright and joyful premium flower arrangement to make birthdays and corporate parties unforgettable.'
  },
  {
    id: 'wedding',
    name: 'Wedding Flower Board',
    price: 250.00,
    rating: 4.8,
    reviews: 210,
    image: '/images/wedding.jpeg',
    tag: 'Romantic',
    desc: 'Exquisite, elegant, and grand flower boards tailored beautifully for luxurious wedding celebrations.'
  },
  {
    id: 'custom',
    name: 'Interactive Custom Board Simulator',
    price: 150.00,
    rating: 5.0,
    reviews: 89,
    image: '/images/custom-preview.png',
    tag: 'Interactive Game',
    desc: 'Design your own professional flower board in real-time! Choose your custom layout, foam colors, and fonts.',
    isCustomRoute: true
  }
])

// Logika navigasi masuk ke dalam sub-folder products/
const navigateToProduct = (item: typeof products.value[0]) => {
  if (item.isCustomRoute || item.id === 'custom') {
    navigateTo('/products/custom') // Mengarah ke app/pages/products/custom.vue
  } else {
    navigateTo(`/products/${item.id}`) // Mengarah ke app/pages/products/[id].vue
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <div class="text-center max-w-2xl mx-auto mb-16 space-y-3">
        <span class="text-xs font-black text-emerald-700 uppercase tracking-widest bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
          Chia Florist Collection
        </span>
        <h1 class="text-4xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
          Our Flower Boards
        </h1>
        <p class="text-sm md:text-base text-gray-500 leading-relaxed">
          Select our pre-designed premium flower boards or jump directly into our interactive real-time simulator game to design your custom creation.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div 
          v-for="item in products" 
          :key="item.id"
          @click="navigateToProduct(item)"
          class="group bg-white border border-gray-100 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between cursor-pointer transform hover:-translate-y-1"
        >
          <div>
            <div class="aspect-[4/3] w-full bg-gray-50 relative overflow-hidden border-b border-gray-50">
              <img 
                :src="item.image" 
                :alt="item.name" 
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
              />
              <span class="absolute top-4 left-4 bg-white/90 backdrop-blur-md text-[#1b4332] text-[10px] font-black tracking-widest uppercase px-3 py-1.5 rounded-xl border border-gray-100 shadow-sm">
                {{ item.tag }}
              </span>
            </div>

            <div class="p-6 space-y-3">
              <div class="flex items-center gap-2 text-xs text-yellow-500 font-bold">
                <span>⭐ {{ item.rating.toFixed(1) }}</span>
                <span class="text-gray-300">|</span>
                <span class="text-gray-400 font-medium">({{ item.reviews }} reviews)</span>
              </div>
              
              <h3 class="text-lg font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-snug">
                {{ item.name }}
              </h3>
              
              <p class="text-xs text-gray-400 leading-relaxed line-clamp-2">
                {{ item.desc }}
              </p>
            </div>
          </div>

          <div class="p-6 pt-0 border-t border-gray-50/50 mt-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Starting From</p>
              <p class="text-xl font-extrabold text-gray-900">${{ item.price.toFixed(2) }}</p>
            </div>
            
            <button 
              class="bg-gray-50 group-hover:bg-[#1b4332] text-gray-700 group-hover:text-white border border-gray-200 group-hover:border-[#1b4332] text-xs font-bold px-4 py-2.5 rounded-xl transition-all flex items-center gap-1.5"
            >
              <span>{{ item.isCustomRoute ? 'Launch Game' : 'View Details' }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 transform transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>