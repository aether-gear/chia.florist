<script setup lang="ts">
// Logic untuk form bisa ditambahkan di sini nanti
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-10 min-h-[80vh] flex items-center">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-16 items-center w-full">
      
      <div class="hidden md:block h-[600px]">
        <img 
          src="/images/flower.jpg" 
          alt="Plants Decor" 
          class="w-full h-full object-cover rounded-xl shadow-sm"
        />
      </div>

      <div class="max-w-md w-full mx-auto space-y-8">
        <div>
          <h1 class="text-4xl font-medium tracking-tight mb-2">Create an account</h1>
          <p class="text-gray-600">Enter your details below</p>
        </div>

        <form @submit.prevent class="space-y-6">
          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              placeholder="Name" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
            />
          </div>

          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="text" 
              placeholder="Email or Phone Number" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
            />
          </div>

          <div class="border-b border-gray-300 py-2 focus-within:border-black transition-colors">
            <input 
              type="password" 
              placeholder="Password" 
              class="w-full outline-none bg-transparent text-lg placeholder:text-gray-400"
            />
          </div>

          <div class="pt-4 space-y-4">
            <button
              type="submit"
              class="w-full bg-accent text-white-base py-4 rounded-md font-medium hover-bg-accent-strong transition-all shadow-md"
            >
              Create Account
            </button>

            <button 
              type="button" 
              class="w-full border border-gray-300 py-4 rounded-md font-medium flex items-center justify-center gap-3 hover:bg-gray-50 transition-all"
            >
              <img src="/images/google.png" class="w-5 h-5" alt="Google Icon" />
              Sign up with Google
            </button>
          </div>
        </form>

        <div class="text-center pt-4 text-gray-600">
          Already have account? 
          <NuxtLink to="/login" class="font-medium text-black border-b border-gray-500 pb-0.5 ml-2 hover-text-accent transition-colors">
            Log in
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Blok style dikosongkan untuk menghindari error @apply di Tailwind v4 */
</style>
