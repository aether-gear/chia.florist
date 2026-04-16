<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
// Note: composables are auto-imported in nuxt, but we can explicitly import for clarity or let Nuxt handle it
import { useProductViewModel } from '~/composables/viewmodels/useProductViewModel'

const route = useRoute()
const productId = route.params.id as string

const { currentProduct, isLoading, error, fetchProductById } = useProductViewModel()

onMounted(() => {
  fetchProductById(productId)
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-20">
    <div v-if="isLoading" class="text-center py-20">Loading product...</div>
    <div v-else-if="error" class="text-red-500 text-center py-20">{{ error }}</div>
    <div v-else-if="currentProduct" class="grid grid-cols-1 md:grid-cols-2 gap-12">
      <!-- Product Image -->
      <div class="bg-gray-100 rounded-lg overflow-hidden h-[500px]">
        <img :src="currentProduct.imageUrl" :alt="currentProduct.name" class="w-full h-full object-cover" />
      </div>
      
      <!-- Product Info -->
      <div class="space-y-6">
        <h1 class="text-4xl font-bold font-serif">{{ currentProduct.name }}</h1>
        <p class="text-2xl text-[#1b4332] font-semibold">${{ currentProduct.price }}</p>
        <p class="text-gray-600 leading-relaxed">{{ currentProduct.description }}</p>
        
        <div class="pt-8 border-t border-gray-200">
          <CButton variant="solid" size="lg" class="w-full">Add to Cart</CButton>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-20">Product not found.</div>
  </div>
</template>
