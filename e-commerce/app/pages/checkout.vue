<!-- app/pages/checkout.vue -->
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useCart, type CartItem } from '~/composables/useCart'
import { useAddress } from '~/composables/useAddress'
import { cartService } from '~/services/cartService'
import type { CheckoutResponse, CheckoutCourierOption, PaymentMethod } from '~/types/checkout'

useHead({
  title: 'Secure Checkout - Chia Florist'
})

const route = useRoute()
const { cart, cartSubtotal, cartSubtotalFormatted, checkoutToOrder, formatRupiah } = useCart()
const addressVm = useAddress()

// State Management untuk Checkout & Shipping
const checkoutData = ref<CheckoutResponse | null>(null)
const isLoadingCheckout = ref(false)
const isLoadingCalculate = ref(false)
const discount = ref(0)
const selectedAddressId = ref('')
const isProcessing = ref(false)

// State Management untuk Payment Methods
const paymentMethods = ref<PaymentMethod[]>([])
const selectedPaymentMethodId = ref('')

// Map untuk menyimpan opsi kurir per toko
const courierOptionsMap = ref<Record<string, CheckoutCourierOption[]>>({})

// Map untuk melacak kurir terpilih per toko
const selectedCouriers = ref<Record<string, { code: string; service: string }>>({})

// 1. Deteksi Alur Pembelian: Direct Buy Now vs. Dari Cart
const isBuyNow = computed(() => route.query.buyNow === 'true')
const buyNowItem = computed<CartItem | null>(() => {
  if (!isBuyNow.value) return null
  return {
    id: route.query.id as string,
    name: route.query.name as string,
    price: Number(route.query.price || 0),
    image: (route.query.image as string) || '',
    quantity: Number(route.query.qty || 1),
    size: (route.query.size as string) || undefined,
    color: (route.query.color as string) || undefined,
    shopId: (route.query.shopId as string) || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
  }
})

// Sumber item yang akan dicheckout
const checkoutItems = computed(() => {
  if (isBuyNow.value && buyNowItem.value) {
    return [buyNowItem.value]
  }
  return cart.value
})

// Helper to merge custom items into CheckoutResponse
const mergeCustomItems = (res: CheckoutResponse | null): CheckoutResponse => {
  const customItems = checkoutItems.value.filter(item => item.isCustom)
  
  let merged: CheckoutResponse
  if (!res) {
    const defaultAddr = addressVm.addresses.value.find(a => a.address_id === selectedAddressId.value) || addressVm.addresses.value.find(a => a.is_default) || addressVm.addresses.value[0]
    merged = {
      address: defaultAddr ? {
        id: defaultAddr.address_id || 'default-addr-id',
        recipient_name: defaultAddr.receiver_name,
        phone: defaultAddr.phone || null,
        full_address: defaultAddr.full_address
      } : {
        id: 'default-addr-id',
        recipient_name: 'No Address Selected',
        phone: '',
        full_address: 'Please add/select a shipping address'
      },
      shops: [],
      subtotal: 0,
      total_shipping: 0,
      total: 0,
      payment_methods: [
        {
          id: "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
          name: "GoPay",
          type: "ewallet",
          description: "GoPay via Midtrans",
          fee: 0,
          subtotal: 0,
          total: 0
        },
        {
          id: "074b02e4-e047-4f60-bdb0-cfeb5481d002",
          name: "DANA",
          type: "ewallet",
          description: "DANA via Midtrans",
          fee: 0,
          subtotal: 0,
          total: 0
        },
        {
          id: "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
          name: "Mandiri",
          type: "bank_transfer",
          description: "Mandiri Bill Payment via Midtrans",
          fee: 0,
          subtotal: 0,
          total: 0
        }
      ]
    }
  } else {
    merged = JSON.parse(JSON.stringify(res))
  }

  // 1. Gabungkan atribut lokal (size, color, price) untuk produk reguler
  merged.shops.forEach(shop => {
    shop.items.forEach(item => {
      const localItem = checkoutItems.value.find(i => i.id === item.product_id)
      if (localItem) {
        if (localItem.size) {
          (item as any).size = localItem.size
        }
        if (localItem.color) {
          (item as any).color = localItem.color
        }
        if (localItem.price) {
          item.price = localItem.price
          item.subtotal = item.price * item.quantity
        }
      }
    })
    
    // Hitung ulang subtotal dan total toko reguler
    shop.subtotal = shop.items.reduce((acc, i) => acc + i.subtotal, 0)
    const fee = shop.selected_courier ? shop.selected_courier.fee : 0
    shop.total = shop.subtotal + fee
  })

  // 2. Tambahkan papan kustom dari simulator
  if (customItems.length === 0) {
    merged.subtotal = merged.shops.reduce((acc, s) => acc + s.subtotal, 0)
    merged.total_shipping = merged.shops.reduce((acc, s) => acc + (s.selected_courier ? s.selected_courier.fee : 0), 0)
    merged.total = merged.subtotal + merged.total_shipping
    return merged
  }

  const customShopsMap: Record<string, typeof customItems> = {}
  customItems.forEach(item => {
    const sId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
    if (!customShopsMap[sId]) {
      customShopsMap[sId] = []
    }
    customShopsMap[sId].push(item)
  })

  Object.keys(customShopsMap).forEach(sId => {
    const items = customShopsMap[sId] || []
    const subtotal = items.reduce((acc, item) => acc + (item.price * item.quantity), 0)
    
    let shopEntry = merged.shops.find(s => s.shop_id === sId)
    if (shopEntry) {
      const entry = shopEntry
      items.forEach(item => {
        if (entry.items && !entry.items.some(i => i.product_id === item.id)) {
          entry.items.push({
            product_id: item.id,
            shop_id: sId,
            name: item.name,
            price: item.price,
            quantity: item.quantity,
            subtotal: item.price * item.quantity,
            size: item.size,
            color: item.color
          } as any)
        }
      })
      entry.subtotal = entry.items ? entry.items.reduce((acc, i) => acc + i.subtotal, 0) : 0
      const fee = entry.selected_courier ? entry.selected_courier.fee : 0
      entry.total = entry.subtotal + fee
    } else {
      const mockOptions = [
        { code: 'jne', name: 'JNE', service: 'REG', etd: '2-3 Days', fee: 20000 },
        { code: 'tiki', name: 'TIKI', service: 'REG', etd: '2-4 Days', fee: 18000 },
        { code: 'pos', name: 'POS', service: 'REG', etd: '3-5 Days', fee: 15000 }
      ]
      
      if (!courierOptionsMap.value[sId]) {
        courierOptionsMap.value[sId] = mockOptions
      }
      
      const selected = selectedCouriers.value[sId] || { code: 'jne', service: 'REG' }
      if (!selectedCouriers.value[sId]) {
        selectedCouriers.value[sId] = selected
      }
      
      const defaultOption = { code: 'jne', name: 'JNE', service: 'REG', etd: '2-3 Days', fee: 20000 }
      const matchedOption = mockOptions.find(o => o.code === selected.code && o.service === selected.service) || mockOptions[0] || defaultOption
      
      shopEntry = {
        shop_id: sId,
        subtotal: subtotal,
        total: subtotal + matchedOption.fee,
        selected_courier: {
          code: matchedOption.code,
          service: matchedOption.service,
          fee: matchedOption.fee
        },
        items: items.map(item => ({
          product_id: item.id,
          shop_id: sId,
          name: item.name,
          price: item.price,
          quantity: item.quantity,
          subtotal: item.price * item.quantity,
          size: item.size,
          color: item.color
        })),
        cost_couriers: mockOptions
      }
      merged.shops.push(shopEntry)
      merged.total_shipping += matchedOption.fee
    }
  })

  merged.subtotal = merged.shops.reduce((acc, s) => acc + s.subtotal, 0)
  merged.total_shipping = merged.shops.reduce((acc, s) => acc + (s.selected_courier ? s.selected_courier.fee : 0), 0)
  merged.total = merged.subtotal + merged.total_shipping

  return merged
}

// Muat data checkout saat halaman dibuka
onMounted(async () => {
  if (checkoutItems.value.length === 0) {
    navigateTo('/catalog')
    return
  }

  isLoadingCheckout.value = true
  try {
    await addressVm.fetchAddresses()

    const defaultAddr = addressVm.addresses.value.find(a => a.is_default)
    if (defaultAddr) {
      selectedAddressId.value = defaultAddr.address_id || ''
    }

    const shopsMap: Record<string, { product_id: string; quantity: number }[]> = {}
    checkoutItems.value.forEach(item => {
      if (item.isCustom) return
      const shopId = item.shopId || '333f6432-a01c-412f-99f4-0f08ca0d8eb1'
      if (!shopsMap[shopId]) {
        shopsMap[shopId] = []
      }
      shopsMap[shopId].push({
        product_id: item.id,
        quantity: item.quantity
      })
    })

    const shopsPayload = Object.keys(shopsMap).map(shopId => ({
      shop_id: shopId,
      items: shopsMap[shopId]
    }))

    let res: CheckoutResponse | null = null
    if (shopsPayload.length > 0) {
      const payload: any = { shops: shopsPayload }
      res = await cartService.checkout(payload)
    }

    const mergedData = mergeCustomItems(res)
    checkoutData.value = mergedData

    if (mergedData) {
      if (mergedData.address?.id) {
        selectedAddressId.value = mergedData.address.id
      }

      if (mergedData.payment_methods && mergedData.payment_methods.length > 0) {
        paymentMethods.value = mergedData.payment_methods
        selectedPaymentMethodId.value = mergedData.payment_methods[0].id
      }

      mergedData.shops.forEach(shop => {
        if (shop.cost_couriers) {
          courierOptionsMap.value[shop.shop_id] = shop.cost_couriers
          if (!selectedCouriers.value[shop.shop_id] && shop.cost_couriers.length > 0) {
            selectedCouriers.value[shop.shop_id] = {
              code: shop.cost_couriers[0].code,
              service: shop.cost_couriers[0].service
            }
          }
        }
        if (shop.selected_courier) {
          selectedCouriers.value[shop.shop_id] = {
            code: shop.selected_courier.code,
            service: shop.selected_courier.service
          }
        }
      })
    }
  } catch (err) {
    console.error('Failed to initialize checkout:', err)
    alert('Unable to proceed with checkout. Some items may be out of stock or unavailable. Redirecting you back to your cart to review.')
    navigateTo('/cart')
  } finally {
    isLoadingCheckout.value = false
  }
})

// 2. Fungsi Hitung Ulang Biaya Pengiriman (Saat Kurir / Alamat Berubah)
const runCalculate = async () => {
  if (!checkoutData.value) return
  if (!selectedAddressId.value || !selectedPaymentMethodId.value) return
  isLoadingCalculate.value = true

  try {
    const backendShopsPayload: any[] = []

    checkoutData.value.shops.forEach(shop => {
      const nonCustomItems = shop.items.filter(item => {
        const localItem = checkoutItems.value.find(i => i.id === item.product_id)
        return !localItem?.isCustom
      })

      if (nonCustomItems.length > 0) {
        const courier = selectedCouriers.value[shop.shop_id]
        let courierPayload = courier || (shop.selected_courier ? {
          code: shop.selected_courier.code,
          service: shop.selected_courier.service
        } : undefined)

        if (!courierPayload) {
          const options = courierOptionsMap.value[shop.shop_id]
          if (options && options.length > 0) {
            courierPayload = { code: options[0].code, service: options[0].service }
            selectedCouriers.value[shop.shop_id] = { code: options[0].code, service: options[0].service }
          } else {
            courierPayload = { code: 'jne', service: 'REG' }
            selectedCouriers.value[shop.shop_id] = { code: 'jne', service: 'REG' }
          }
        }

        backendShopsPayload.push({
          shop_id: shop.shop_id,
          items: nonCustomItems.map(item => ({
            product_id: item.product_id,
            quantity: item.quantity
          })),
          courier: courierPayload
        })
      }
    })

    let res: CheckoutResponse | null = null
    if (backendShopsPayload.length > 0) {
      const payload = {
        address_id: selectedAddressId.value,
        payment_method_id: selectedPaymentMethodId.value,
        shops: backendShopsPayload
      }
      res = await cartService.checkoutCalculate(payload)
    }

    let newCheckoutData: CheckoutResponse
    if (res) {
      newCheckoutData = JSON.parse(JSON.stringify(res))
    } else {
      newCheckoutData = JSON.parse(JSON.stringify(checkoutData.value))
      const selectedAddr = addressVm.addresses.value.find(a => a.address_id === selectedAddressId.value)
      if (selectedAddr) {
        newCheckoutData.address = {
          id: selectedAddr.address_id || 'default-addr-id',
          recipient_name: selectedAddr.receiver_name,
          phone: selectedAddr.phone || null,
          full_address: selectedAddr.full_address
        }
      }
    }

    const mergedData = mergeCustomItems(res ? newCheckoutData : null)
    checkoutData.value = mergedData

    if (res) {
      res.shops.forEach(shop => {
        if (shop.selected_courier) {
          selectedCouriers.value[shop.shop_id] = {
            code: shop.selected_courier.code,
            service: shop.selected_courier.service
          }
        }
      })
    }
  } catch (err) {
    console.error('Failed to calculate shipping:', err)
  } finally {
    isLoadingCalculate.value = false
  }
}

let addressTimeout: ReturnType<typeof setTimeout> | null = null

// Pantau perubahan alamat untuk kalkulasi ulang ongkir
watch(selectedAddressId, (newId, oldId) => {
  if (newId && oldId && newId !== oldId) {
    if (addressTimeout) {
      clearTimeout(addressTimeout)
    }
    addressTimeout = setTimeout(() => {
      runCalculate()
    }, 300)
  }
})

let paymentMethodTimeout: ReturnType<typeof setTimeout> | null = null

// Pantau perubahan metode pembayaran untuk kalkulasi ulang ongkir
watch(selectedPaymentMethodId, (newId, oldId) => {
  if (newId && oldId && newId !== oldId) {
    if (paymentMethodTimeout) {
      clearTimeout(paymentMethodTimeout)
    }
    paymentMethodTimeout = setTimeout(() => {
      runCalculate()
    }, 300)
  }
})

let courierTimeout: ReturnType<typeof setTimeout> | null = null

onUnmounted(() => {
  if (paymentMethodTimeout) {
    clearTimeout(paymentMethodTimeout)
  }
  if (addressTimeout) {
    clearTimeout(addressTimeout)
  }
  if (courierTimeout) {
    clearTimeout(courierTimeout)
  }
})

// Helper pencari opsi kurir
const getCourierOptions = (shopId: string): CheckoutCourierOption[] => {
  return courierOptionsMap.value[shopId] || []
}

// Helper nilai value select kurir
const getCourierSelectValue = (shopId: string): string => {
  const selected = selectedCouriers.value[shopId]
  return selected ? `${selected.code}|${selected.service}` : ''
}

// Penanganan perubahan pilihan kurir oleh user
const handleCourierChange = (shopId: string, event: Event) => {
  const select = event.target as HTMLSelectElement
  const val = select.value
  if (val) {
    const [code, service] = val.split('|')
    if (code && service) {
      selectedCouriers.value[shopId] = { code, service }
      if (courierTimeout) {
        clearTimeout(courierTimeout)
      }
      courierTimeout = setTimeout(() => {
        runCalculate()
      }, 300)
    }
  }
}

// Helper untuk mencocokkan URL gambar lokal
const getCartProductImage = (productId: string): string => {
  const localItem = checkoutItems.value.find(i => i.id === productId)
  return localItem?.image || '/images/custom-preview.png'
}

// 3. Live Values untuk Ringkasan Pembayaran
const liveSubtotal = computed(() => {
  return checkoutData.value ? checkoutData.value.subtotal : cartSubtotal.value
})
const liveShippingFee = computed(() => {
  return checkoutData.value ? checkoutData.value.total_shipping : 0
})
const liveTotalPayment = computed(() => {
  return checkoutData.value ? checkoutData.value.total : (cartSubtotal.value - discount.value)
})

// Eksekusi checkout memindahkan state item keranjang ke invoice order profile
const handlePlaceOrder = async () => {
  if (!selectedAddressId.value) {
    alert('Please select a shipping address before completing your order.')
    return
  }

  isProcessing.value = true

  try {
    // Selesaikan pemesanan dan sinkronkan keranjang backend
    await checkoutToOrder(liveTotalPayment.value, checkoutItems.value)
    
    alert('Order placed successfully! Redirecting to secure payment page...')
    navigateTo('/profile')
  } catch (err) {
    console.error('Checkout processing error:', err)
    alert('Failed to process checkout. Please try again.')
  } finally {
    isProcessing.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-16 font-sans">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
      
      <div class="mb-10">
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Secure Checkout</h1>
        <p class="text-sm text-gray-500 mt-1">Please confirm your shipping metadata and billing totals below.</p>
      </div>

      <div v-if="isLoadingCheckout" class="flex flex-col items-center justify-center min-h-[400px] space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#1b4332]"></div>
        <p class="text-gray-500 font-medium animate-pulse text-sm">Initializing secure checkout session...</p>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
        
        <div class="lg:col-span-7 space-y-6">
          
          <!-- 1. Alamat Pengiriman -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">1. Shipping Destination</h3>
            
            <div v-if="addressVm.isLoading.value" class="flex flex-col items-center justify-center py-6 space-y-2">
              <div class="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-[#1b4332]"></div>
              <p class="text-gray-400 text-xs">Fetching address cards...</p>
            </div>

            <div class="text-center py-6 border-2 border-dashed border-gray-200 rounded-2xl p-4" v-else-if="addressVm.addresses.value.length === 0">
              <p class="text-sm text-gray-500">No addresses registered to your profile.</p>
              <NuxtLink to="/profile" class="text-xs font-bold text-[#1b4332] underline mt-1 inline-block">Add address in Profile Settings</NuxtLink>
            </div>

            <div class="space-y-3" v-else>
              <label 
                v-for="addr in addressVm.addresses.value" 
                :key="addr.address_id"
                :class="['border rounded-2xl p-4 flex items-start gap-4 cursor-pointer transition-all', selectedAddressId === addr.address_id ? 'border-[#1b4332] bg-emerald-50/5' : 'border-gray-200 hover:border-gray-300', (isLoadingCalculate || isLoadingCheckout) ? 'opacity-50 pointer-events-none cursor-not-allowed' : '']"
              >
                <input 
                  type="radio" 
                  v-model="selectedAddressId" 
                  :value="addr.address_id" 
                  :disabled="isLoadingCalculate || isLoadingCheckout"
                  class="mt-1 accent-[#1b4332]" 
                />
                <div class="flex-1 text-xs">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-bold text-gray-900">{{ addr.receiver_name }}</span>
                    <span class="bg-emerald-100 text-emerald-800 font-bold text-[9px] px-2 py-0.2 rounded-full" v-if="addr.is_default">Default</span>
                  </div>
                  <p class="text-gray-600 font-semibold mb-1">📞 {{ addr.phone }}</p>
                  <p class="text-gray-500 leading-normal">{{ addr.full_address }}, {{ addr.postal_code }}</p>
                </div>
              </label>
            </div>
          </div>

          <!-- 2. Ringkasan Papan Bunga & Kurir Per Toko -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <div class="flex justify-between items-center border-b border-gray-50 pb-4">
              <h3 class="font-bold text-gray-900 text-lg">2. Review Ordered Flower Boards</h3>
              <div v-if="isLoadingCalculate" class="flex items-center gap-1.5 text-xs text-emerald-700 font-semibold">
                <div class="animate-spin rounded-full h-3.5 w-3.5 border-t-2 border-b-2 border-emerald-700"></div>
                <span>Recalculating...</span>
              </div>
            </div>

            <div v-if="checkoutData" class="space-y-6">
              <div v-for="shop in checkoutData.shops" :key="shop.shop_id" class="border border-gray-100 rounded-2xl p-5 space-y-4 shadow-sm">
                <!-- Header Toko -->
                <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-bold text-gray-500 uppercase tracking-wider">Seller Shop ID:</span>
                    <span class="text-xs font-mono font-bold text-gray-700 bg-gray-100 px-2 py-0.5 rounded">{{ shop.shop_id.slice(0, 8) }}...</span>
                  </div>
                  <span class="text-xs font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full">Subtotal: {{ formatRupiah(shop.subtotal) }}</span>
                </div>
                
                <!-- Items list in this shop -->
                <div class="divide-y divide-gray-50">
                  <div v-for="item in shop.items" :key="item.product_id" class="flex gap-4 items-center py-3 first:pt-0 last:pb-0">
                    <div class="w-14 h-14 rounded-xl overflow-hidden bg-gray-50 border border-gray-100 flex-shrink-0">
                      <img :src="getCartProductImage(item.product_id)" :alt="item.name" class="w-full h-full object-cover" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <h4 class="font-bold text-gray-900 text-sm truncate">{{ item.name }}</h4>
                      <div class="flex flex-wrap items-center gap-x-2 gap-y-0.5 text-xs text-gray-400 mt-1 font-semibold">
                        <span>Qty: {{ item.quantity }}</span>
                        <span>|</span>
                        <span>Price: {{ formatRupiah(item.price) }}</span>
                        <span v-if="(item as any).size">| Size: {{ (item as any).size }}</span>
                        <span v-if="(item as any).color" class="flex items-center gap-1">
                          | Color:
                          <span :style="{ backgroundColor: (item as any).color }" class="w-2.5 h-2.5 rounded-full border border-gray-300 inline-block"></span>
                        </span>
                      </div>
                    </div>
                    <div class="text-sm font-extrabold text-gray-900 text-right">
                      {{ formatRupiah(item.subtotal) }}
                    </div>
                  </div>
                </div>

                <!-- Kurir untuk Toko ini -->
                <div class="bg-gray-50/50 rounded-2xl p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border border-gray-100">
                  <div class="text-xs font-bold text-gray-500">
                    🚚 Choose Delivery Courier:
                  </div>
                  <div class="flex items-center">
                    <select 
                      :value="getCourierSelectValue(shop.shop_id)"
                      @change="handleCourierChange(shop.shop_id, $event)"
                      :disabled="isLoadingCalculate || isLoadingCheckout"
                      class="bg-white border border-gray-200 rounded-xl text-xs p-2.5 outline-none focus:border-emerald-700 transition-all font-bold cursor-pointer text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <option 
                        v-for="courier in getCourierOptions(shop.shop_id)" 
                        :key="courier.code + '|' + courier.service"
                        :value="courier.code + '|' + courier.service"
                      >
                        {{ courier.name || courier.code.toUpperCase() }} ({{ courier.service }}) - {{ formatRupiah(courier.fee) }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
            </div>

            <!-- Fallback UI apabila data backend checkout gagal termuat -->
            <div v-else class="divide-y divide-gray-100 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
              <div v-for="(item, idx) in checkoutItems" :key="idx" class="flex gap-4 items-center py-4 first:pt-0 last:pb-0">
                <div class="w-14 h-14 rounded-xl overflow-hidden bg-gray-50 border border-gray-100 flex-shrink-0">
                  <img :src="item.image || '/images/custom-preview.png'" :alt="item.name" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h4 class="font-bold text-gray-900 text-sm truncate">{{ item.name }}</h4>
                  <p class="text-xs text-gray-400 mt-1 font-semibold">Qty: {{ item.quantity }} | Size: {{ item.size || 'Standard' }}</p>
                </div>
                <div class="text-sm font-extrabold text-gray-900 text-right">
                  {{ formatRupiah(item.price * item.quantity) }}
                </div>
              </div>
            </div>
          </div>

          <!-- 3. Payment Method -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">3. Payment Method</h3>
            
            <div v-if="paymentMethods.length === 0" class="text-center py-6 border-2 border-dashed border-gray-200 rounded-2xl p-4">
              <p class="text-sm text-gray-500">No payment methods available.</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4" v-else>
              <label 
                v-for="method in paymentMethods" 
                :key="method.id"
                :class="['border rounded-2xl p-4 flex items-start gap-4 cursor-pointer transition-all', selectedPaymentMethodId === method.id ? 'border-[#1b4332] bg-emerald-50/5' : 'border-gray-200 hover:border-gray-300', (isLoadingCalculate || isLoadingCheckout) ? 'opacity-50 pointer-events-none cursor-not-allowed' : '']"
              >
                <input 
                  type="radio" 
                  v-model="selectedPaymentMethodId" 
                  :value="method.id" 
                  :disabled="isLoadingCalculate || isLoadingCheckout"
                  class="mt-1 accent-[#1b4332]" 
                />
                <div class="flex-1 text-xs">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-bold text-gray-900">{{ method.name }}</span>
                    <span class="bg-emerald-100 text-emerald-800 font-bold text-[9px] px-2 py-0.2 rounded-full uppercase">{{ method.type.replace('_', ' ') }}</span>
                  </div>
                  <p class="text-gray-500 leading-normal mb-1 font-semibold">{{ method.description }}</p>
                  <p class="text-emerald-700 font-bold" v-if="method.fee > 0">Fee: {{ formatRupiah(method.fee) }}</p>
                  <p class="text-emerald-700 font-bold text-[10px]" v-else>Free Transaction Fee</p>
                </div>
              </label>
            </div>
          </div>
        </div>

        <!-- Sidebar Summary Tagihan -->
        <div class="lg:col-span-4 space-y-6">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">Billing Summary</h3>
            
            <div class="space-y-4 text-sm font-medium text-gray-600">
              <div class="flex justify-between items-center">
                <span>Subtotal</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(liveSubtotal) }}</span>
              </div>
              
              <div class="flex justify-between items-center">
                <span>Shipping Cost</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(liveShippingFee) }}</span>
              </div>
              
              <div class="flex justify-between items-center" v-if="discount > 0">
                <span>Promo Discount</span>
                <span class="font-bold">-{{ formatRupiah(discount) }}</span>
              </div>
              
              <div class="border-t border-gray-100 pt-4 flex justify-between items-center text-base font-bold text-gray-900">
                <span>Total Bill</span>
                <span class="text-2xl font-black text-[#1b4332]">
                  {{ formatRupiah(liveTotalPayment) }}
                </span>
              </div>
            </div>

            <button 
              @click="handlePlaceOrder"
              :disabled="isProcessing || addressVm.addresses.value.length === 0 || isLoadingCalculate || !selectedPaymentMethodId"
              class="w-full bg-[#1b4332] hover:bg-[#143326] disabled:bg-gray-300 text-white font-bold py-4 rounded-xl transition shadow-md hover:shadow-lg text-center text-sm tracking-wide cursor-pointer disabled:cursor-not-allowed flex items-center justify-center"
            >
              <span v-if="isProcessing">Processing Order...</span>
              <span v-else>Confirm & Pay Now</span>
            </button>
          </div>
          
          <div class="bg-emerald-50/50 border border-emerald-100 rounded-2xl p-4 flex gap-3 items-center text-emerald-800">
            <span class="text-2xl">🔒</span>
            <p class="text-xs font-medium leading-normal">Your payment request is fully managed under a cryptographically secured end-to-end sandbox module.</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
</style>