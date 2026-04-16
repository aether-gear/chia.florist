<script setup lang="ts">
const name = ref('')
const email = ref('')
const password = ref('')

const handleSubmit = () => {
  console.log({
    name: name.value,
    email: email.value,
    password: password.value
  })
}
</script>

<template>
  <div class="flex">

    <!-- Left Image -->
    <div class="w-1/2 hidden md:block bg-gray-300 relative">
      <img
        src="https://placehold.co/900x700"
        class="w-full h-full object-cover"
      />
    </div>

    <!-- Form -->
    <div class="flex-1 flex items-center justify-center">
      <div class="w-full max-w-md space-y-6">

        <div>
          <h1 class="text-3xl font-semibold">Create an account</h1>
          <p class="text-gray-500 text-sm">Enter your details below</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">

          <input v-model="name" type="text" placeholder="Name"
            class="w-full border-b py-2 outline-none" />

          <input v-model="email" type="text" placeholder="Email or Phone Number"
            class="w-full border-b py-2 outline-none" />

          <input v-model="password" type="password" placeholder="Password"
            class="w-full border-b py-2 outline-none" />

          <button type="submit"
            class="w-full bg-teal-900 text-white py-3 rounded">
            Create Account
          </button>
        </form>

        <button class="w-full border py-3">
          Sign up with Google
        </button>

        <p class="text-center text-sm">
          Already have account?
          <NuxtLink to="/login" class="underline ml-1">
            Log in
          </NuxtLink>
        </p>

      </div>
    </div>

  </div>
</template>