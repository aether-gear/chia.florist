<!-- app/pages/checkout.vue -->
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useCart, type CartItem } from '~/composables/useCart'
import { useAddress } from '~/composables/useAddress'
import { cartService } from '~/services/cartService'
import { orderService } from '~/services/orderService'
import { bootstrapConfig } from '~/utils/bootstrap'
import type { CheckoutResponse, CheckoutCourierOption, PaymentMethod, CheckoutShop } from '~/types/checkout'
import type { UserAddress } from '~/types/address'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { triggerAuthAlert } from '~/composables/useSessionState'

useHead({
  title: 'Secure Checkout - Chia Florist'
})

const route = useRoute()
const { cart, orders, loadCart, flushCart, cartSubtotal, cartSubtotalFormatted, checkoutToOrder, formatRupiah } = useCart()
const addressVm = useAddress()
const authVm = useAuthViewModel()

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
const openedCategories = ref<Record<string, boolean>>({})

// Dynamic shop mapping and payment state
const shopsMap = ref<Record<string, string>>({})
const paymentInfoState = useState<any>('last-payment-info', () => null)

const fetchShops = async () => {
  try {
    const res = await bootstrapConfig.fetchApi<{ shops: { id: string; name: string }[] }>('/shops')
    if (res && res.shops) {
      res.shops.forEach(s => {
        shopsMap.value[s.id] = s.name
      })
    }
  } catch (err) {
    console.error('Failed to fetch shops:', err)
  }
}

const paymentMethodTypes = ['qr_code', 'ewallet', 'bank_transfer']

const paymentMethodsGroupedByType = computed(() => {
  const groups: Record<string, PaymentMethod[]> = {}
  paymentMethodTypes.forEach(t => {
    groups[t] = []
  })
  paymentMethods.value.forEach(method => {
    if (groups[method.type]) {
      groups[method.type]!.push(method)
    }
  })
  return groups
})

const formatPaymentType = (type: string) => {
  if (type === 'ewallet') return 'E-Wallet'
  if (type === 'bank_transfer') return 'Bank Transfer'
  if (type === 'qr_code') return 'QR Code / QRIS'
  return type.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}

const toggleCategory = (type: string) => {
  openedCategories.value[type] = !openedCategories.value[type]
}

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
    shopId: (route.query.shopId as string) || '99ef0062-1040-4574-a4be-0123abce5670'
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
      payment_methods: []
    }
    
    // Group ALL items by shopId
    const itemsByShop: Record<string, typeof checkoutItems.value> = {}
    checkoutItems.value.forEach(item => {
      const sId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
      if (!itemsByShop[sId]) {
        itemsByShop[sId] = []
      }
      itemsByShop[sId].push(item)
    })
    
    Object.keys(itemsByShop).forEach(sId => {
      const items = itemsByShop[sId] || []
      const subtotal = items.reduce((acc, item) => acc + (item.price * item.quantity), 0)
      
      const shopEntry: CheckoutShop = {
        shop_id: sId,
        name: shopsMap.value[sId] || 'Chia Florist',
        subtotal: subtotal,
        total: subtotal,
        selected_courier: null,
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
        cost_couriers: []
      }
      merged.shops.push(shopEntry)
    })
    
    merged.subtotal = merged.shops.reduce((acc, s) => acc + s.subtotal, 0)
    merged.total_shipping = 0
    merged.total = merged.subtotal
    
    return merged
  } else {
    merged = JSON.parse(JSON.stringify(res))
  }

  // 1. Gabungkan atribut lokal (size, color, price) untuk produk reguler & kustom
  merged.shops.forEach(shop => {
    if (!shop.name) {
      shop.name = shopsMap.value[shop.shop_id] || 'Chia Florist'
    }
    shop.items.forEach(item => {
      const localItem = checkoutItems.value.find(i =>
        i.id === item.product_id ||
        (i.isCustom && ((item as any).product_variant_type === 'custom' || (item as any).item_type === 'custom'))
      )
      if (localItem) {
        if (localItem.size) {
          (item as any).size = localItem.size
        }
        if (localItem.color) {
          (item as any).color = localItem.color
        }
        if (localItem.customDesign && !(item as any).custom_design) {
          (item as any).custom_design = localItem.customDesign
        }
      }
    })
    
    // Hitung ulang subtotal dan total toko reguler
    shop.subtotal = shop.items.reduce((acc, i) => acc + i.subtotal, 0)
    const fee = shop.selected_courier ? shop.selected_courier.fee : 0
    shop.total = shop.subtotal + fee
  })

  merged.subtotal = merged.shops.reduce((acc, s) => acc + s.subtotal, 0)
  merged.total_shipping = merged.shops.reduce((acc, s) => acc + (s.selected_courier ? s.selected_courier.fee : 0), 0)
  merged.total = merged.subtotal + merged.total_shipping

  return merged
}

// Muat data checkout saat halaman dibuka
onMounted(async () => {
  await authVm.fetchCurrentUser()
  await fetchShops()
  if (!authVm.isAuthenticated.value) {
    triggerAuthAlert('warning', 'Please sign in to proceed with checkout.')
    navigateTo('/login')
    return
  }

  if (checkoutItems.value.length === 0) {
    navigateTo('/catalog')
    return
  }

  isLoadingCheckout.value = true
  try {
    await addressVm.fetchAddresses()

    let defaultAddr = addressVm.addresses.value.find(a => a.is_default)
    if (!defaultAddr && addressVm.addresses.value.length > 0) {
      console.log('Self-healing on mount: Setting first address as default')
      const firstAddr = addressVm.addresses.value[0]
      if (firstAddr) {
        const updated: UserAddress = { ...firstAddr, is_default: true }
        await addressVm.saveAddress(updated)
        await addressVm.fetchAddresses()
        defaultAddr = addressVm.addresses.value.find(a => a.is_default)
      }
    }

    if (defaultAddr) {
      selectedAddressId.value = defaultAddr.address_id || ''
    }

    const shopsMap: Record<string, any[]> = {}
    checkoutItems.value.forEach(item => {
      const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
      if (!shopsMap[shopId]) {
        shopsMap[shopId] = []
      }
      if (item.isCustom) {
        shopsMap[shopId].push({
          cart_item_id: item.cartItemId || item.id,
          product_variant_type: 'custom',
          item_type: 'custom',
          product_name: item.name,
          physical_size_id: item.customDesign?.layout?.physicalSizeId || item.size || 'medium',
          unit_price: item.price,
          quantity: item.quantity,
          custom_design: item.customDesign
        })
      } else {
        shopsMap[shopId].push({
          item_type: 'standard',
          product_id: item.id,
          quantity: item.quantity
        })
      }
    })

    const shopsPayload = Object.keys(shopsMap).map(shopId => ({
      shop_id: shopId,
      items: shopsMap[shopId]
    }))

    let res: CheckoutResponse | null = null
    if (defaultAddr && shopsPayload.length > 0) {
      const payload: any = { shops: shopsPayload }
      try {
        res = await cartService.checkout(payload)
      } catch (checkoutErr: any) {
        console.warn('Backend returned error during checkout initialization, using fallback client estimation:', checkoutErr)
        // Fallback gracefully without throwing so the UI can still load
      }
    }

    const mergedData = mergeCustomItems(res)
    checkoutData.value = mergedData

    if (mergedData) {
      if (mergedData.address?.id) {
        selectedAddressId.value = mergedData.address.id
      }

      if (mergedData.payment_methods && mergedData.payment_methods.length > 0) {
        paymentMethods.value = mergedData.payment_methods
      } else {
        paymentMethods.value = []
      }

      const firstMethod = paymentMethods.value[0]
      if (firstMethod) {
        selectedPaymentMethodId.value = firstMethod.id
        paymentMethodTypes.forEach(type => {
          openedCategories.value[type] = true
        })
      }

      mergedData.shops.forEach(shop => {
        if (shop.cost_couriers) {
          courierOptionsMap.value[shop.shop_id] = shop.cost_couriers
          const firstOption = shop.cost_couriers[0]
          if (!selectedCouriers.value[shop.shop_id] && firstOption) {
            selectedCouriers.value[shop.shop_id] = {
              code: firstOption.code,
              service: firstOption.service
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
      const courier = selectedCouriers.value[shop.shop_id]
      let courierPayload = courier || (shop.selected_courier ? {
        code: shop.selected_courier.code,
        service: shop.selected_courier.service
      } : undefined)

      if (!courierPayload) {
        const options = courierOptionsMap.value[shop.shop_id]
        const firstOption = options && options[0]
        if (firstOption) {
          courierPayload = { code: firstOption.code, service: firstOption.service }
          selectedCouriers.value[shop.shop_id] = { code: firstOption.code, service: firstOption.service }
        } else {
          courierPayload = { code: 'jne', service: 'REG' }
          selectedCouriers.value[shop.shop_id] = { code: 'jne', service: 'REG' }
        }
      }

      const shopItemsPayload = shop.items.map(item => {
        const localItem = checkoutItems.value.find(i => i.id === item.product_id || (i.isCustom && (item.product_variant_type === 'custom' || item.item_type === 'custom')))
        if (localItem?.isCustom || item.product_variant_type === 'custom' || item.item_type === 'custom') {
          const design = localItem?.customDesign || item.custom_design
          const cartItemId = item.cart_item_id || localItem?.cartItemId || localItem?.id || item.product_id
          return {
            cart_item_id: cartItemId,
            product_variant_type: 'custom' as const,
            item_type: 'custom' as const,
            product_name: item.name,
            physical_size_id: design?.layout?.physicalSizeId || 'medium',
            unit_price: item.price,
            quantity: item.quantity,
            custom_design: design
          }
        }
        return {
          product_variant_type: 'standard' as const,
          item_type: 'standard' as const,
          product_id: item.product_id,
          quantity: item.quantity
        }
      })

      backendShopsPayload.push({
        shop_id: shop.shop_id,
        items: shopItemsPayload,
        courier: courierPayload
      })
    })

    let res: CheckoutResponse | null = null
    if (backendShopsPayload.length > 0) {
      const payload = {
        address_id: selectedAddressId.value,
        payment_method_id: selectedPaymentMethodId.value,
        shops: backendShopsPayload
      }
      try {
        res = await cartService.checkoutCalculate(payload)
      } catch (calcErr: any) {
        console.warn('Backend returned error during calculate, using fallback client estimation:', calcErr)
      }
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
  } catch (err: any) {
    console.error('Failed to calculate shipping:', err)
    if (err.status === 409 && err.data?.message && err.data.message.includes('address')) {
      const addresses = addressVm.addresses.value
      const defaultAddr = addresses.find(a => a.is_default)
      if (!defaultAddr && addresses.length > 0) {
        console.log('Self-healing in runCalculate: Setting default address and retrying calculation...')
        try {
          const firstAddr = addresses[0]
          if (firstAddr) {
            const updated: UserAddress = { ...firstAddr, is_default: true }
            await addressVm.saveAddress(updated)
            await addressVm.fetchAddresses()
            setTimeout(() => {
              runCalculate()
            }, 100)
            return
          }
        } catch (saveErr) {
          console.error('Self-healing failed to save default address:', saveErr)
        }
      }
    }
  } finally {
    isLoadingCalculate.value = false
  }
}

let addressTimeout: ReturnType<typeof setTimeout> | null = null

// Pantau perubahan alamat untuk kalkulasi ulang ongkir
watch(selectedAddressId, (newId, oldId) => {
  if (newId && newId !== oldId && !isLoadingCheckout.value) {
    isLoadingCalculate.value = true
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
  if (newId && newId !== oldId && !isLoadingCheckout.value) {
    isLoadingCalculate.value = true
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
      isLoadingCalculate.value = true
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
const livePaymentFee = computed(() => {
  if (checkoutData.value?.selected_payment_method?.id === selectedPaymentMethodId.value) {
    return checkoutData.value.selected_payment_method.fee
  }
  const activeMethod = paymentMethods.value.find(m => m.id === selectedPaymentMethodId.value)
  return activeMethod ? activeMethod.fee : 0
})
const liveTotalPayment = computed(() => {
  if (!checkoutData.value) return cartSubtotal.value - discount.value
  if (checkoutData.value.selected_payment_method?.id === selectedPaymentMethodId.value) {
    return checkoutData.value.total - discount.value
  }
  const sub = checkoutData.value.subtotal
  const ship = checkoutData.value.total_shipping
  const fee = livePaymentFee.value
  return sub + ship + fee - discount.value
})

// Eksekusi checkout memindahkan state item keranjang ke invoice order profile
const handlePlaceOrder = async () => {
  if (!selectedAddressId.value) {
    alert('Please select a shipping address before completing your order.')
    return
  }
  if (paymentMethods.value.length === 0 || !selectedPaymentMethodId.value) {
    alert('No payment method available or selected. Please select a payment method before completing your order.')
    return
  }

  isProcessing.value = true

  try {
    if (!checkoutData.value || !checkoutData.value.shops) {
      throw new Error('Checkout details are not fully loaded.')
    }

    const shopsPayload = checkoutData.value.shops.map(shop => {
      const shopName = shop.name || shopsMap.value[shop.shop_id] || 'Chia Florist'
      const courier = selectedCouriers.value[shop.shop_id]
      if (!courier) {
        throw new Error(`Please select a courier for the shop: ${shopName}`)
      }
      return {
        shop_id: shop.shop_id,
        name: shopName,
        selected_courier: {
          code: courier.code,
          service: courier.service
        },
        items: shop.items.map(item => {
          const localItem = checkoutItems.value.find(i => i.id === item.product_id || (i.isCustom && (item.product_variant_type === 'custom' || item.item_type === 'custom')))
          if (localItem?.isCustom || item.product_variant_type === 'custom' || item.item_type === 'custom') {
            const design = localItem?.customDesign || item.custom_design
            const cartItemId = item.cart_item_id || localItem?.cartItemId || localItem?.id || item.product_id
            return {
              cart_item_id: cartItemId,
              product_variant_type: 'custom' as const,
              item_type: 'custom' as const,
              name: item.name,
              physical_size_id: design?.layout?.physicalSizeId || 'medium',
              quantity: item.quantity,
              unit_price: item.price,
              custom_design: design
            }
          }
          return {
            product_variant_type: 'standard' as const,
            item_type: 'standard' as const,
            product_id: item.product_id,
            name: item.name,
            quantity: item.quantity
          }
        })
      }
    })

    const payload = {
      address_id: selectedAddressId.value,
      selected_payment: {
        id: selectedPaymentMethodId.value
      },
      shops: shopsPayload
    }

    const result = await orderService.createOrder(payload)
    
    const paymentInfo = {
      orderId: result.order_id,
      instruction: result.instruction,
      channelData: result.channel_data,
      total: liveTotalPayment.value,
      status: 'pending'
    }
    paymentInfoState.value = paymentInfo
    if (import.meta.client) {
      sessionStorage.setItem('chia-last-payment-info', JSON.stringify(paymentInfo))
    }

    const orderItems = [...checkoutItems.value]
    const newOrder = {
      orderId: result.order_id,
      date: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }),
      items: orderItems,
      total: liveTotalPayment.value,
      status: 'pending' as const
    }
    orders.value.push(newOrder)

    if (authVm.isAuthenticated.value) {
      for (const item of orderItems) {
        if (import.meta.client) {
          localStorage.removeItem(`cart_attr_${item.id}`)
        }
      }
      await loadCart(true)
    } else {
      for (const item of orderItems) {
        if (import.meta.client) {
          localStorage.removeItem(`cart_attr_${item.id}`)
        }
      }
      cart.value = []
    }

    alert('Order placed successfully! Redirecting to secure payment page...')
    navigateTo(`/payment?orderId=${result.order_id}`)
  } catch (err: any) {
    console.error('Checkout processing error:', err)
    alert(err.data?.message || err.message || 'Failed to process checkout. Please try again.')
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
          <div class="bg-white border border-gray-100 rounded-3xl p-6 md:p-8 shadow-sm space-y-6">
            <h3 class="font-bold text-gray-900 text-lg border-b border-gray-50 pb-4">3. Payment Method</h3>
            
            <!-- Backend Payment Methods Selection (Always Visible) -->
            <div class="space-y-4">
              <div class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2">Select Payment Gateway Channel:</div>
              <div v-if="paymentMethods.length === 0" class="text-center py-6 border-2 border-dashed border-gray-200 rounded-2xl p-4">
                <p class="text-sm text-gray-500">No payment methods available.</p>
              </div>

              <div v-else class="space-y-3">
                <!-- Iterate over payment categories (types) -->
                <div 
                  v-for="type in paymentMethodTypes" 
                  :key="type" 
                  class="border border-gray-100 rounded-2xl overflow-hidden transition-all duration-300"
                  :class="[openedCategories[type] ? 'border-[#1b4332] shadow-sm' : 'hover:border-gray-200']"
                >
                  <!-- Accordion Header -->
                  <button
                    type="button"
                    @click="toggleCategory(type)"
                    :disabled="isLoadingCalculate || isLoadingCheckout"
                    class="w-full flex items-center justify-between p-4 bg-gray-50/50 hover:bg-gray-50/80 transition-all font-bold text-xs text-gray-700 text-left outline-none cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    <div class="flex items-center gap-2">
                      <span class="text-sm" v-if="type === 'ewallet'">📱</span>
                      <span class="text-sm" v-else-if="type === 'bank_transfer'">🏦</span>
                      <span class="text-sm" v-else-if="type === 'qr_code'">🔍</span>
                      <span class="text-sm" v-else>💳</span>
                      <span>{{ formatPaymentType(type) }}</span>
                    </div>
                    <span class="text-xs transition-transform duration-300" :class="[openedCategories[type] ? 'rotate-180' : '']">
                      ▼
                    </span>
                  </button>

                  <!-- Accordion Content: Vertical Flat Radio List -->
                  <div 
                    v-show="openedCategories[type]" 
                    class="p-4 bg-white border-t border-gray-50 flex flex-col space-y-2"
                  >
                    <label 
                      v-for="method in paymentMethodsGroupedByType[type]" 
                      :key="method.id"
                      :class="['border rounded-xl p-3.5 flex items-center justify-between gap-3.5 cursor-pointer transition-all', selectedPaymentMethodId === method.id ? 'border-[#1b4332] bg-emerald-50/5' : 'border-gray-200 hover:border-gray-300', (isLoadingCalculate || isLoadingCheckout) ? 'opacity-50 pointer-events-none cursor-not-allowed' : '']"
                    >
                      <div class="flex items-center gap-3">
                        <input 
                          type="radio" 
                          v-model="selectedPaymentMethodId" 
                          :value="method.id" 
                          :disabled="isLoadingCalculate || isLoadingCheckout"
                          class="accent-[#1b4332]" 
                        />
                        <span class="font-bold text-gray-900 text-xs">{{ method.name }}</span>
                      </div>
                      <div class="text-right text-xs">
                        <span class="text-emerald-700 font-bold" v-if="method.fee > 0">Fee: {{ formatRupiah(method.fee) }}</span>
                        <span class="text-emerald-700 font-bold text-[10px]" v-else>Free Fee</span>
                      </div>
                    </label>
                  </div>
                </div>
              </div>
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

              <div class="flex justify-between items-center">
                <span>Payment Fee</span>
                <span class="text-gray-900 font-bold">
                  {{ livePaymentFee > 0 ? formatRupiah(livePaymentFee) : 'Free' }}
                </span>
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
              :disabled="isProcessing || addressVm.addresses.value.length === 0 || isLoadingCalculate || paymentMethods.length === 0 || !selectedPaymentMethodId"
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