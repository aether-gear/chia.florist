<script setup lang="ts">
import { ref, computed, watch, onUnmounted, nextTick } from 'vue'
import { productService } from '~/services/productService'
import { formatRupiah } from '~/utils/formatter'
import type { CatalogProduct } from '~/types/product'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useCart } from '~/composables/useCart'
import { useAddress } from '~/composables/useAddress'
import { filterCatalogProductsByQuery } from '~/utils/searchMatcher'

const authVm = useAuthViewModel()
const cart = useCart()
const route = useRoute()
const { addresses, fetchAddresses } = useAddress()

const ensureAddressesLoaded = () => {
  if (authVm.isAuthenticated.value && addresses.value.length === 0) {
    fetchAddresses().catch((err) => console.warn('Failed to load user addresses for navbar:', err))
  }
}

// Get primary/default user address
const defaultAddress = computed(() => {
  if (!addresses.value || addresses.value.length === 0) return null
  return addresses.value.find(a => a.is_default) || addresses.value[0]
})

// --- NAVIGATION LINKS ---
const navLinks = [
  { name: 'Home', to: '/' },
  { name: 'Catalog', to: '/catalog' },
  { name: 'Custom Board', to: '/products/custom' }
]

// --- MOBILE MENU STATE ---
const isMobileMenuOpen = ref(false)
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}
const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

// --- PROFILE DROPDOWN STATE ---
const isProfileOpen = ref(false)
let profileCloseTimeout: ReturnType<typeof setTimeout> | null = null

const openProfileDropdown = () => {
  if (profileCloseTimeout) clearTimeout(profileCloseTimeout)
  isProfileOpen.value = true
  ensureAddressesLoaded()
}

const closeProfileDropdown = () => {
  profileCloseTimeout = setTimeout(() => {
    isProfileOpen.value = false
  }, 200)
}

const toggleProfileDropdown = () => {
  isProfileOpen.value = !isProfileOpen.value
  if (isProfileOpen.value) {
    ensureAddressesLoaded()
  }
}

// --- SEARCH STATE & LOGIC ---
const isSearchOpen = ref(false)
const searchQuery = ref('')
const searchInput = ref<HTMLInputElement | null>(null)
const searchResults = ref<CatalogProduct[]>([])
const isLoadingSearch = ref(false)

const customSimulatorCard: CatalogProduct = {
  id: 'custom',
  name: 'Custom Flower Board Simulator',
  price: 150000,
  rating: 5.0,
  reviews: 89,
  image: '/images/custom-preview.png',
  tag: 'Interactive Game',
  desc: 'Design your own professional flower board in real-time!',
  isCustomRoute: true,
  isAvailable: true
}

let debounceTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, (newQuery) => {
  if (debounceTimeout) clearTimeout(debounceTimeout)

  const query = newQuery.trim()
  if (!query) {
    searchResults.value = []
    return
  }

  isLoadingSearch.value = true

  debounceTimeout = setTimeout(async () => {
    try {
      const apiResults = await productService.getCatalogProducts()
      searchResults.value = filterCatalogProductsByQuery(
        apiResults,
        query,
        true,
        customSimulatorCard
      )
    } catch (err) {
      console.error('Search failed:', err)
      searchResults.value = []
    } finally {
      isLoadingSearch.value = false
    }
  }, 300)
})

onUnmounted(() => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  if (profileCloseTimeout) clearTimeout(profileCloseTimeout)
})

const openSearch = async () => {
  isSearchOpen.value = true
  await nextTick()
  searchInput.value?.focus()
}

const closeSearch = () => {
  isSearchOpen.value = false
  searchQuery.value = ''
  searchResults.value = []
}

const goToDedicatedSearch = () => {
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim()
    closeSearch()
    navigateTo(`/search?q=${encodeURIComponent(q)}`)
  }
}

// Close mobile menu on route change
watch(() => route.path, () => {
  closeMobileMenu()
  isProfileOpen.value = false
})
</script>

<template>
  <header class="w-full bg-white/95 backdrop-blur-md border-b border-gray-100 sticky top-0 z-50 font-brand transition-all">
    <div class="max-w-7xl mx-auto px-8 sm:px-10 py-3.5 flex items-center justify-between gap-4">

      <!-- BRAND LOGO -->
      <div class="flex items-center">
        <NuxtLink to="/" class="flex items-center gap-2 group">
          <img src="/images/logo.png" alt="Chia Florist" class="h-7 sm:h-8 w-auto object-contain transition-transform group-hover:scale-105" />
        </NuxtLink>
      </div>

      <!-- DESKTOP NAVIGATION LINKS (Visible on lg+) -->
      <nav class="hidden lg:flex items-center gap-8">
        <NuxtLink
          v-for="item in navLinks"
          :key="item.to"
          :to="item.to"
          class="relative py-1 text-sm font-medium transition-colors after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-[2px] after:bottom-0 after:left-0 after:bg-[#1b4332] after:origin-bottom-left after:transition-transform after:duration-300 hover:after:scale-x-100"
          :class="[route.path === item.to ? 'text-[#1b4332] font-semibold after:scale-x-100' : 'text-gray-600 hover:text-[#1b4332]']"
        >
          {{ item.name }}
        </NuxtLink>
      </nav>

      <!-- RIGHT UTILITIES & AUTH -->
      <div class="flex items-center gap-2 sm:gap-4">

        <!-- SEARCH BUTTON (Hidden on screens < lg, moved to mobile panel) -->
        <button
          @click="openSearch"
          class="hidden lg:inline-flex p-2 text-gray-500 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all duration-200 cursor-pointer"
          title="Search catalog"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5.5 w-5.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>

        <!-- CART BUTTON WITH BADGE (Visible on all screens including mobile) -->
        <NuxtLink
          to="/cart"
          class="relative p-2 text-gray-500 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-all duration-200 cursor-pointer"
          title="Shopping Cart"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5.5 w-5.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 0 0-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 0 0-16.536-1.84M7.5 14.25L5.106 5.25M16.5 20.25a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Zm3 0a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z" />
          </svg>

          <!-- Dynamic Badge -->
          <span
            v-if="cart.cartCount.value > 0"
            class="absolute top-0.5 right-0.5 min-w-[18px] h-[18px] bg-[#1b4332] text-white text-[10px] font-bold rounded-full px-1 flex items-center justify-center border-2 border-white shadow-xs"
          >
            {{ cart.cartCount.value > 99 ? '99+' : cart.cartCount.value }}
          </span>
        </NuxtLink>

        <!-- PROFILE & AUTH SECTION (Hidden on screens < lg, moved to mobile panel) -->
        <div class="hidden lg:block relative border-l border-gray-100 pl-2 sm:pl-4">

          <!-- AUTHENTICATED USER DROPDOWN -->
          <template v-if="authVm.isAuthenticated.value">
            <div
              class="relative"
              @mouseenter="openProfileDropdown"
              @mouseleave="closeProfileDropdown"
            >
              <button
                @click="toggleProfileDropdown"
                class="flex items-center gap-2 py-1 px-2 hover:bg-gray-50 rounded-full transition-all duration-200 cursor-pointer group"
                title="Account Menu"
              >
                <div class="relative w-8 h-8 rounded-full overflow-hidden bg-emerald-50 border border-emerald-100 flex items-center justify-center text-[#1b4332] group-hover:bg-[#1b4332] group-hover:text-white transition-all duration-200 shadow-2xs">
                  <img
                    v-if="authVm.currentUser.value?.avatarUrl"
                    :src="authVm.currentUser.value.avatarUrl"
                    class="w-full h-full object-cover"
                    alt="User Avatar"
                  />
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span class="absolute bottom-0 right-0 w-2 h-2 bg-emerald-500 border-2 border-white rounded-full"></span>
                </div>

                <span class="hidden sm:inline text-xs font-bold text-gray-700 group-hover:text-[#1b4332] max-w-[90px] truncate">
                  {{ authVm.currentUser.value?.name || 'Account' }}
                </span>

                <svg xmlns="http://www.w3.org/2000/svg" class="hidden sm:block h-3.5 w-3.5 text-gray-400 group-hover:text-[#1b4332] transition-transform duration-200" :class="[isProfileOpen ? 'rotate-180' : '']" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <!-- DROPDOWN MENU -->
              <Transition name="dropdown">
                <div
                  v-if="isProfileOpen"
                  class="absolute right-0 mt-1.5 w-64 bg-white rounded-2xl shadow-xl border border-gray-100 py-2 z-50 overflow-hidden"
                >
                  <!-- USER DETAILS HEADER -->
                  <div class="px-4 py-3 border-b border-gray-100 bg-gray-50/50 space-y-2">
                    <div>
                      <p class="text-xs font-bold text-gray-900 truncate">
                        {{ authVm.currentUser.value?.name || 'Customer' }}
                      </p>
                      <p class="text-[11px] text-gray-500 truncate mt-0.5">
                        {{ authVm.currentUser.value?.email || 'user@chiaflorist.com' }}
                      </p>
                    </div>

                    <!-- USER PRIMARY ADDRESS -->
                    <div class="pt-2 border-t border-gray-200/60">
                      <div class="flex items-center gap-1 text-[10px] font-bold text-emerald-800 uppercase tracking-wider">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-emerald-700 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                          <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        <span>Primary Address</span>
                      </div>
                      <p v-if="defaultAddress" class="text-[11px] text-gray-600 line-clamp-2 mt-0.5 leading-tight">
                        {{ defaultAddress.full_address }}
                      </p>
                      <p v-else class="text-[11px] text-gray-400 italic mt-0.5">
                        No address added yet
                      </p>
                    </div>
                  </div>

                  <!-- MENU OPTIONS -->
                  <div class="py-1">
                    <NuxtLink
                      to="/profile/personal"
                      @click="isProfileOpen = false"
                      class="flex items-center gap-2.5 px-4 py-2.5 text-xs font-medium text-gray-700 hover:bg-emerald-50/60 hover:text-[#1b4332] transition-colors"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                      </svg>
                      <span>My Profile</span>
                    </NuxtLink>

                    <NuxtLink
                      to="/profile/addresses"
                      @click="isProfileOpen = false"
                      class="flex items-center gap-2.5 px-4 py-2.5 text-xs font-medium text-gray-700 hover:bg-emerald-50/60 hover:text-[#1b4332] transition-colors"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                      <span>Shipping Addresses</span>
                    </NuxtLink>

                    <NuxtLink
                      to="/profile/orders"
                      @click="isProfileOpen = false"
                      class="flex items-center gap-2.5 px-4 py-2.5 text-xs font-medium text-gray-700 hover:bg-emerald-50/60 hover:text-[#1b4332] transition-colors"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
                      </svg>
                      <span>My Orders</span>
                    </NuxtLink>
                  </div>

                  <div class="border-t border-gray-100 pt-1">
                    <button
                      @click="() => { isProfileOpen = false; authVm.logout() }"
                      class="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs font-medium text-red-600 hover:bg-red-50 transition-colors cursor-pointer text-left"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                      </svg>
                      <span>Sign Out</span>
                    </button>
                  </div>
                </div>
              </Transition>
            </div>
          </template>

          <!-- UNAUTHENTICATED SIGN IN -->
          <template v-else>
            <CButton
              to="/login"
              variant="secondary"
              size="pill"
              title="Sign In"
            >
              Sign In
            </CButton>
          </template>
        </div>

        <!-- MOBILE MENU HAMBURGER TOGGLE (Visible on screens < lg) -->
        <button
          @click="toggleMobileMenu"
          class="lg:hidden p-2 text-gray-600 hover:text-[#1b4332] hover:bg-gray-50 rounded-full transition-colors cursor-pointer"
          title="Toggle menu"
        >
          <svg v-if="!isMobileMenuOpen" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

      </div>

    </div>

    <!-- MOBILE MENU DRAWER PANEL (Visible on screens < lg) -->
    <Transition name="mobile-drawer">
      <div v-if="isMobileMenuOpen" class="lg:hidden border-t border-gray-100 bg-white px-4 pt-4 pb-6 space-y-4">

        <!-- SEARCH TRIGGER IN MOBILE PANEL -->
        <div>
          <button
            @click="() => { closeMobileMenu(); openSearch() }"
            class="w-full flex items-center gap-2.5 px-4 py-2.5 bg-gray-50 hover:bg-emerald-50/60 rounded-2xl text-xs font-medium text-gray-500 hover:text-[#1b4332] transition-colors cursor-pointer border border-gray-100"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <span>Search flowers, bouquets & boards...</span>
          </button>
        </div>

        <!-- NAVIGATION LINKS -->
        <nav class="space-y-1">
          <NuxtLink
            v-for="item in navLinks"
            :key="item.to"
            :to="item.to"
            @click="closeMobileMenu"
            class="block px-3 py-2.5 rounded-xl text-sm font-medium transition-colors"
            :class="[route.path === item.to ? 'bg-emerald-50 text-[#1b4332] font-bold' : 'text-gray-700 hover:bg-gray-50 hover:text-[#1b4332]']"
          >
            {{ item.name }}
          </NuxtLink>
        </nav>

        <!-- PROFILE & AUTH IN MOBILE PANEL -->
        <div class="pt-3 border-t border-gray-100 space-y-2">
          <template v-if="authVm.isAuthenticated.value">
            <div class="px-3 py-2.5 bg-gray-50 rounded-xl space-y-2 border border-gray-100">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full bg-emerald-50 border border-emerald-100 flex items-center justify-center text-[#1b4332] text-xs font-bold flex-shrink-0">
                  <img v-if="authVm.currentUser.value?.avatarUrl" :src="authVm.currentUser.value.avatarUrl" class="w-full h-full object-cover rounded-full" />
                  <span v-else>{{ authVm.currentUser.value?.name?.charAt(0).toUpperCase() || 'U' }}</span>
                </div>
                <div class="truncate">
                  <p class="text-xs font-bold text-gray-900 truncate">{{ authVm.currentUser.value?.name || 'Customer' }}</p>
                  <p class="text-[10px] text-gray-500 truncate">{{ authVm.currentUser.value?.email || 'user@chiaflorist.com' }}</p>
                </div>
              </div>

              <!-- USER PRIMARY ADDRESS -->
              <div class="pt-2 border-t border-gray-200/60">
                <div class="flex items-center gap-1 text-[10px] font-bold text-emerald-800 uppercase tracking-wider">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-emerald-700 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span>Primary Address</span>
                </div>
                <p v-if="defaultAddress" class="text-[11px] text-gray-600 line-clamp-2 mt-0.5 leading-tight">
                  {{ defaultAddress.full_address }}
                </p>
                <p v-else class="text-[11px] text-gray-400 italic mt-0.5">
                  No address added yet
                </p>
              </div>
            </div>

            <NuxtLink
              to="/profile/personal"
              @click="closeMobileMenu"
              class="flex items-center gap-2 px-3 py-2 text-xs font-medium text-gray-700 hover:text-[#1b4332]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              <span>My Profile</span>
            </NuxtLink>

            <NuxtLink
              to="/profile/addresses"
              @click="closeMobileMenu"
              class="flex items-center gap-2 px-3 py-2 text-xs font-medium text-gray-700 hover:text-[#1b4332]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>Shipping Addresses</span>
            </NuxtLink>

            <NuxtLink
              to="/profile/orders"
              @click="closeMobileMenu"
              class="flex items-center gap-2 px-3 py-2 text-xs font-medium text-gray-700 hover:text-[#1b4332]"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-700" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
              </svg>
              <span>My Orders</span>
            </NuxtLink>

            <button
              @click="() => { closeMobileMenu(); authVm.logout() }"
              class="w-full flex items-center gap-2 text-left px-3 py-2 text-xs font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              <span>Sign Out</span>
            </button>
          </template>

          <template v-else>
            <CButton
              to="/login"
              @click="closeMobileMenu"
              variant="secondary"
              size="lg"
              full-width
            >
              Sign In
            </CButton>
          </template>
        </div>
      </div>
    </Transition>

    <!-- SEARCH MODAL TELEPORT -->
    <Teleport to="body">
      <div v-if="isSearchOpen" class="fixed inset-0 z-[100] flex justify-end">
        <Transition name="fade">
          <div v-if="isSearchOpen" @click="closeSearch" class="absolute inset-0 bg-black/40 backdrop-blur-xs"></div>
        </Transition>

        <Transition name="slide">
          <div v-if="isSearchOpen" class="relative w-full max-w-md bg-white h-full shadow-2xl flex flex-col overflow-hidden font-brand">

            <div class="p-6 flex items-center justify-between border-b border-gray-100">
              <span class="text-xs font-bold text-gray-400 uppercase tracking-widest">Search Catalog</span>
              <button @click="closeSearch" class="p-2 text-gray-400 hover:text-gray-800 transition-colors cursor-pointer">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="p-6 border-b border-gray-50">
              <div class="relative">
                <input
                  ref="searchInput"
                  v-model="searchQuery"
                  @keydown.enter="goToDedicatedSearch"
                  type="text"
                  placeholder="Search flower boards, hand bouquets..."
                  class="w-full bg-surface text-sm text-gray-800 px-5 py-3.5 pr-10 rounded-2xl outline-none border border-transparent focus:border-[#1b4332] focus:bg-white transition-all shadow-inner"
                />
                <span v-if="searchQuery" @click="searchQuery = ''" class="absolute right-3.5 top-1/2 -translate-y-1/2 text-xs text-gray-400 hover:text-gray-600 cursor-pointer">
                  ✕
                </span>
              </div>
              <p class="text-[10px] text-gray-400 mt-2 flex items-center justify-between">
                <span>Press <strong>Enter ↵</strong> for full search page</span>
                <NuxtLink
                  v-if="searchQuery"
                  :to="`/search?q=${encodeURIComponent(searchQuery)}`"
                  @click="closeSearch"
                  class="text-emerald-700 hover:underline font-bold"
                >
                  View all results →
                </NuxtLink>
              </p>
            </div>

            <div class="flex-1 overflow-y-auto px-6 py-4 custom-scrollbar">
              <div v-if="isLoadingSearch" class="flex flex-col items-center justify-center py-12 space-y-3">
                <div class="animate-spin rounded-full h-7 w-7 border-t-2 border-b-2 border-[#1b4332]"></div>
                <p class="text-xs text-gray-400 font-medium animate-pulse">Searching catalog...</p>
              </div>

              <div v-else-if="searchQuery && searchResults.length > 0" class="space-y-3">
                <div class="flex justify-between items-center mb-1">
                  <span class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Results</span>
                  <span class="text-[10px] font-mono text-gray-400">{{ searchResults.length }} items</span>
                </div>

                <NuxtLink
                  v-for="product in searchResults"
                  :key="product.id"
                  :to="product.isCustomRoute || product.id === 'custom' ? '/products/custom' : `/products/${product.slug || product.id}`"
                  @click="closeSearch"
                  class="flex gap-3.5 p-3 rounded-2xl hover:bg-emerald-50/50 transition-all group cursor-pointer border border-transparent hover:border-emerald-100"
                >
                  <div class="w-14 h-14 rounded-xl overflow-hidden bg-gray-100 flex-shrink-0 flex items-center justify-center">
                    <img v-if="product.image" :src="product.image" :alt="product.name" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div class="flex flex-col justify-center">
                    <h4 class="text-xs font-bold text-gray-900 group-hover:text-[#1b4332] transition-colors leading-tight">
                      {{ product.name }}
                    </h4>
                    <p class="text-[11px] font-semibold text-[#1b4332] mt-1">Starting From {{ formatRupiah(product.price) }}</p>
                  </div>
                </NuxtLink>

                <div class="pt-2">
                  <NuxtLink
                    :to="`/search?q=${encodeURIComponent(searchQuery)}`"
                    @click="closeSearch"
                    class="block text-center py-2.5 px-4 rounded-xl bg-emerald-50 hover:bg-emerald-100 text-[#1b4332] text-xs font-bold transition"
                  >
                    Open dedicated search page for "{{ searchQuery }}" →
                  </NuxtLink>
                </div>
              </div>

              <div v-else-if="searchQuery && searchResults.length === 0" class="h-full flex flex-col items-center justify-center py-12 text-gray-400 space-y-3">
                <p class="text-3xl">🔍</p>
                <p class="text-xs font-bold">No products match "{{ searchQuery }}"</p>
                <NuxtLink
                  :to="`/search?q=${encodeURIComponent(searchQuery)}`"
                  @click="closeSearch"
                  class="text-xs text-emerald-700 underline font-bold"
                >
                  Search on dedicated page
                </NuxtLink>
              </div>

              <div v-else class="h-full flex flex-col items-center justify-center py-12 text-gray-300">
                <p class="text-[11px] font-semibold uppercase tracking-wider">Type to start searching...</p>
              </div>
            </div>

          </div>
        </Transition>
      </div>
    </Teleport>
  </header>
</template>

<style scoped>
/* Profile Dropdown Transition */
.dropdown-enter-active, .dropdown-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* Mobile Drawer Transition */
.mobile-drawer-enter-active, .mobile-drawer-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.mobile-drawer-enter-from, .mobile-drawer-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Slide Transition for Search Drawer */
.slide-enter-active, .slide-leave-active {
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from, .slide-leave-to {
  transform: translateX(100%);
}

/* Fade Transition for Backdrop */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* Custom Scrollbar Mini */
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--color-border-soft);
  border-radius: 10px;
}
</style>
