<script setup lang="ts">
// Home page
useHead({
  title: 'Chia Florist - Beautiful Flower Boards',
  meta: [
    { name: 'description', content: 'We provide customized flower boards for weddings, condolences, graduations, and corporate events.' }
  ]
})
</script>

<template>
  <div class="relative min-h-screen flex items-center justify-center overflow-hidden">
    <!-- Background Image -->
    <div class="absolute inset-0 z-0">
      <img src="/hero-bg.png" alt="Flower Background" class="w-full h-full object-cover" />
      <!-- Dark overlay for better text readability -->
      <div class="absolute inset-0 bg-black/40"></div>
    </div>

    <!-- Hero Content -->
    <div class="relative z-10 text-center px-4 max-w-4xl mx-auto mt-20">
      <h1 class="text-5xl md:text-7xl font-bold text-white mb-6 leading-tight drop-shadow-lg">
        Beautiful Flower Boards
        <span class="block">for Every Special Moment</span>
      </h1>
      
      <p class="text-lg md:text-xl text-gray-200 mb-10 max-w-2xl mx-auto drop-shadow-md leading-relaxed">
        We provide customized flower boards for weddings, condolences, graduations, and corporate events. Designed with care and delivered on time to make every moment memorable.
      </p>
      
      <div class="flex flex-col sm:flex-row items-center justify-center gap-6">
        <CButton variant="solid" size="lg" to="/products">
          Order Now
        </CButton>
        
        <CButton variant="outline" size="lg" to="/gallery" class="group">
          <svg class="w-6 h-6 mr-2 group-hover:text-black transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          View Products
        </CButton>
      </div>
    </div>

    <!-- Right Side Indicators (as seen in image) -->
    <div class="absolute right-8 top-1/2 transform -translate-y-1/2 flex flex-col gap-6 z-10 hidden lg:flex">
      <span class="text-yellow-400 font-bold text-lg cursor-pointer">01</span>
      <span class="text-white/60 hover:text-white font-medium text-sm cursor-pointer transition-colors">02</span>
      <span class="text-white/60 hover:text-white font-medium text-sm cursor-pointer transition-colors">03</span>
    </div>
  </div>
</template>
