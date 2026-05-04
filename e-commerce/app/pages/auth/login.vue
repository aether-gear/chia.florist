<script setup lang="ts">
// Auth ViewModel would handle the login function here

const handleSubmit = () => {
  // Logic here
}
</script>

<template>
  <div class="flex min-h-screen">
    <!-- Left Image -->
    <div class="w-1/2 hidden md:block relative">
      <img
        src="https://images.unsplash.com/photo-1563241598-be9846e4c767?q=80&w=1500&auto=format&fit=crop"
        class="w-full h-full object-cover"
        alt="Login bg"
      />
    </div>

    <!-- Form -->
    <div class="flex-1 flex items-center justify-center p-8 bg-white">
      <div class="w-full max-w-md space-y-6">
        <div>
          <h1 class="text-4xl font-serif font-bold text-[#1b4332] mb-2">Welcome Back</h1>
          <p class="text-gray-500">Please enter your details to log in</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <input type="email" placeholder="Email Address"
            class="w-full border p-3 rounded-lg outline-none focus:ring-2 focus:ring-[#1b4332]" />

          <input type="password" placeholder="Password"
            class="w-full border p-3 rounded-lg outline-none focus:ring-2 focus:ring-[#1b4332]" />

          <CButton variant="solid" class="w-full py-4 text-lg">Log in</CButton>
        </form>

        <p class="text-center text-sm text-gray-600 pt-4">
          Don't have an account?
          <NuxtLink to="/auth/register" class="text-yellow-600 font-bold hover:underline ml-1">
            Sign Up
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
