<script setup lang="ts">
// View logic handled by a viewmodel
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 py-20 mt-20">
    <h1 class="text-4xl font-bold font-serif mb-8 text-accent">My Profile</h1>
    <div class="bg-white-base p-8 shadow rounded-lg grid grid-cols-1 md:grid-cols-3 gap-8">
      <!-- Profile Side -->
      <div class="col-span-1 border-r border-gray-200 pr-8">
        <div class="w-24 h-24 bg-gray-200 rounded-full mx-auto mb-4"></div>
        <h2 class="text-center font-bold text-lg">John Doe</h2>
        <p class="text-center text-gray-500 text-sm mb-6">john.doe@example.com</p>
        <div class="space-y-2">
          <NuxtLink to="/profile" class="block p-2 bg-yellow-100 text-accent rounded font-medium">Personal Info</NuxtLink>
          <NuxtLink to="/orders" class="block p-2 text-gray-600 hover:bg-gray-50 rounded">Order History</NuxtLink>
        </div>
      </div>
      
      <!-- Profile Details -->
      <div class="col-span-2 space-y-6">
        <h3 class="text-2xl font-semibold">Personal Information</h3>
        <p class="text-gray-600">You can manage your account settings here.</p>
        <!-- Dummy form -->
        <div class="grid grid-cols-2 gap-4">
          <input type="text" placeholder="First Name" class="border p-2 rounded outline-none cursor-not-allowed" disabled value="John" />
          <input type="text" placeholder="Last Name" class="border p-2 rounded outline-none cursor-not-allowed" disabled value="Doe" />
        </div>
      </div>
    </div>
  </div>
</template>
