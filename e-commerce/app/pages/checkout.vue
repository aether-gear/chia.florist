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
import { useGlobalAlert } from '~/composables/useGlobalAlert'

useHead({
  title: 'Secure Checkout - Chia Florist',
  meta: [
    { name: 'description', content: 'Complete your flower board order with secure payment and verified courier delivery at Chia Florist.' }
  ]
})

const route = useRoute()
const { cart, orders, loadCart, flushCart, cartSubtotal, checkoutToOrder, formatRupiah } = useCart()
const addressVm = useAddress()
const authVm = useAuthViewModel()
const globalAlert = useGlobalAlert()

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

    const shopsPayloadMap: Record<string, any[]> = {}
    checkoutItems.value.forEach(item => {
      const shopId = item.shopId || '99ef0062-1040-4574-a4be-0123abce5670'
      if (!shopsPayloadMap[shopId]) {
        shopsPayloadMap[shopId] = []
      }
      if (item.isCustom) {
        shopsPayloadMap[shopId].push({
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
        shopsPayloadMap[shopId].push({
          item_type: 'standard',
          product_id: item.id,
          quantity: item.quantity
        })
      }
    })

    const shopsPayload = Object.keys(shopsPayloadMap).map(shopId => ({
      shop_id: shopId,
      items: shopsPayloadMap[shopId]
    }))

    let res: CheckoutResponse | null = null
    if (defaultAddr && shopsPayload.length > 0) {
      const payload: any = { shops: shopsPayload }
      try {
        res = await cartService.checkout(payload)
      } catch (checkoutErr: any) {
        console.warn('Backend returned error during checkout initialization, using fallback client estimation:', checkoutErr)
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
    globalAlert.showError('Checkout Error', 'Unable to proceed with checkout. Some items may be out of stock or unavailable. Redirecting to cart...')
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
    globalAlert.showWarning('Address Required', 'Please select a shipping address before completing your order.')
    return
  }
  if (paymentMethods.value.length === 0 || !selectedPaymentMethodId.value) {
    globalAlert.showWarning('Payment Method Required', 'No payment method available or selected. Please select a payment method before completing your order.')
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

    globalAlert.showSuccess(
      'Order Placed',
      'Order placed successfully! Redirecting to secure payment page...',
      [
        { label: 'Pay Now', onClick: () => navigateTo(`/payment?orderId=${result.order_id}`) },
        { label: 'My Orders', onClick: () => navigateTo('/profile') }
      ]
    )
    navigateTo(`/payment?orderId=${result.order_id}`)
  } catch (err: any) {
    console.error('Checkout processing error:', err)
    globalAlert.showError('Checkout Failed', err.data?.message || err.message || 'Failed to process checkout. Please try again.')
  } finally {
    isProcessing.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-10 sm:py-14 font-sans">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Page Navigation & Title Header -->
      <div class="mb-8 sm:mb-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div class="flex items-center gap-2 mb-2">
            <CButton
              to="/cart"
              variant="ghost"
              size="sm"
              class="-ml-2 text-gray-500 hover:text-gray-900"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
              </svg>
              <span>Back to Cart</span>
            </CButton>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-gray-900 tracking-tight">Secure Checkout</h1>
          <p class="text-sm text-gray-500 mt-1">Confirm your delivery destination, shipping method, and billing totals.</p>
        </div>

        <div class="hidden sm:flex items-center gap-2 bg-emerald-50/80 border border-emerald-100 text-emerald-900 px-3.5 py-2 rounded-xl text-xs font-semibold">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-emerald-700 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
          </svg>
          <span>256-Bit Encrypted Transaction</span>
        </div>
      </div>

      <!-- Loading Checkout State -->
      <div v-if="isLoadingCheckout" class="bg-white border border-gray-100 rounded-3xl p-12 sm:p-16 shadow-xs flex flex-col items-center justify-center min-h-[380px] space-y-4 text-center">
        <div class="relative flex items-center justify-center">
          <div class="w-12 h-12 rounded-full border-3 border-gray-100 border-t-[#1b4332] animate-spin"></div>
        </div>
        <div>
          <h2 class="text-base font-bold text-gray-800">Initializing Checkout Session</h2>
          <p class="text-xs text-gray-400 mt-1 max-w-sm">Fetching verified destination rates and supported payment gateways...</p>
        </div>
      </div>

      <!-- Main Checkout Layout -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-10 items-start">
        
        <!-- Left Column: Checkout Steps & Forms -->
        <div class="lg:col-span-8 space-y-6">
          
          <!-- STEP 1: Alamat Pengiriman -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 sm:p-8 shadow-xs space-y-6">
            <div class="flex items-center justify-between border-b border-gray-100 pb-4">
              <div class="flex items-center gap-3">
                <span class="w-7 h-7 rounded-full bg-[#1b4332] text-white flex items-center justify-center text-xs font-bold shrink-0">1</span>
                <h2 class="font-extrabold text-gray-900 text-lg">Shipping Destination</h2>
              </div>
              <CButton
                to="/profile"
                variant="outline"
                size="sm"
                class="text-xs"
              >
                <span>Manage Addresses</span>
              </CButton>
            </div>
            
            <div v-if="addressVm.isLoading.value" class="flex flex-col items-center justify-center py-8 space-y-3">
              <div class="w-6 h-6 rounded-full border-2 border-gray-200 border-t-[#1b4332] animate-spin"></div>
              <p class="text-gray-400 text-xs font-medium">Loading destination addresses...</p>
            </div>

            <!-- Empty Address State -->
            <div v-else-if="addressVm.addresses.value.length === 0" class="text-center py-8 border-2 border-dashed border-gray-200 rounded-2xl p-6 space-y-3">
              <div class="w-12 h-12 rounded-full bg-emerald-50 text-[#1b4332] flex items-center justify-center mx-auto">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
                </svg>
              </div>
              <p class="text-sm font-semibold text-gray-800">No shipping addresses found</p>
              <p class="text-xs text-gray-500 max-w-sm mx-auto">Please add your event delivery location in your profile settings before continuing.</p>
              <div class="pt-2">
                <CButton
                  to="/profile"
                  variant="primary"
                  size="sm"
                >
                  <span>Add Address in Profile</span>
                </CButton>
              </div>
            </div>

            <!-- Address Radio Card Selection -->
            <div v-else class="space-y-3">
              <label 
                v-for="addr in addressVm.addresses.value" 
                :key="addr.address_id"
                :class="[
                  'border rounded-2xl p-4.5 sm:p-5 flex items-start gap-4 cursor-pointer transition-all duration-200 select-none',
                  selectedAddressId === addr.address_id 
                    ? 'border-[#1b4332] bg-emerald-50/25 ring-1 ring-[#1b4332] shadow-2xs' 
                    : 'border-gray-200 hover:border-gray-300 bg-white hover:bg-gray-50/50',
                  (isLoadingCalculate || isLoadingCheckout) ? 'opacity-60 pointer-events-none cursor-not-allowed' : ''
                ]"
              >
                <!-- Custom Radio Button Indicator -->
                <div class="mt-0.5 relative flex items-center justify-center shrink-0">
                  <input 
                    type="radio" 
                    v-model="selectedAddressId" 
                    :value="addr.address_id" 
                    :disabled="isLoadingCalculate || isLoadingCheckout"
                    class="sr-only" 
                  />
                  <div 
                    class="w-5 h-5 rounded-full border transition-all flex items-center justify-center"
                    :class="selectedAddressId === addr.address_id ? 'border-[#1b4332] bg-[#1b4332]' : 'border-gray-300 bg-white'"
                  >
                    <div v-if="selectedAddressId === addr.address_id" class="w-2 h-2 rounded-full bg-white"></div>
                  </div>
                </div>

                <div class="flex-1 min-w-0 text-xs sm:text-sm">
                  <div class="flex items-center flex-wrap gap-2 mb-1.5">
                    <span class="font-bold text-gray-900 text-sm sm:text-base">{{ addr.receiver_name }}</span>
                    <span v-if="addr.is_default" class="bg-emerald-100 text-emerald-800 font-bold text-[10px] uppercase tracking-wider px-2.5 py-0.5 rounded-full border border-emerald-200">
                      Default
                    </span>
                  </div>

                  <!-- Contact Phone -->
                  <div class="flex items-center gap-1.5 text-gray-600 font-medium mb-1 text-xs">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-gray-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 002.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 01-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 00-1.091-.852H4.5A2.25 2.25 0 002.25 4.5v2.25z" />
                    </svg>
                    <span>{{ addr.phone }}</span>
                  </div>

                  <!-- Full Address -->
                  <div class="flex items-start gap-1.5 text-gray-500 leading-relaxed text-xs">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-gray-400 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
                    </svg>
                    <span>{{ addr.full_address }}<span v-if="addr.postal_code">, {{ addr.postal_code }}</span></span>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <!-- STEP 2: Ringkasan Papan Bunga & Kurir Per Toko -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 sm:p-8 shadow-xs space-y-6">
            <div class="flex justify-between items-center border-b border-gray-100 pb-4">
              <div class="flex items-center gap-3">
                <span class="w-7 h-7 rounded-full bg-[#1b4332] text-white flex items-center justify-center text-xs font-bold shrink-0">2</span>
                <h2 class="font-extrabold text-gray-900 text-lg">Review Items & Courier</h2>
              </div>
              <div v-if="isLoadingCalculate" class="flex items-center gap-1.5 text-xs text-emerald-800 font-semibold bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100">
                <div class="w-3 h-3 rounded-full border-2 border-emerald-700 border-t-transparent animate-spin"></div>
                <span>Recalculating rates...</span>
              </div>
            </div>

            <!-- Shops and items -->
            <div v-if="checkoutData" class="space-y-6">
              <div 
                v-for="shop in checkoutData.shops" 
                :key="shop.shop_id" 
                class="border border-gray-200/80 rounded-2xl p-5 sm:p-6 space-y-5 bg-white shadow-2xs"
              >
                <!-- Header Toko -->
                <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-3.5 border-b border-gray-100 gap-2">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-emerald-800 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 21v-7.5a.75.75 0 01.75-.75h3a.75.75 0 01.75.75V21m-4.5 0H2.25A2.25 2.25 0 010 18.75V10.5m13.5 10.5h7.5a2.25 2.25 0 002.25-2.25V10.5M3 10.5l9-7.5 9 7.5" />
                    </svg>
                    <span class="text-xs font-bold uppercase tracking-wider text-gray-500">Fulfilled by:</span>
                    <span class="text-xs font-bold text-gray-900 bg-gray-100 px-2.5 py-1 rounded-md">{{ shop.name || 'Chia Florist Branch' }}</span>
                  </div>
                  <span class="text-xs font-bold text-emerald-800 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-100 self-start sm:self-auto">
                    Subtotal: {{ formatRupiah(shop.subtotal) }}
                  </span>
                </div>
                
                <!-- Items list in this shop -->
                <div class="divide-y divide-gray-100">
                  <div 
                    v-for="item in shop.items" 
                    :key="item.product_id" 
                    class="flex gap-4 items-center py-4 first:pt-0 last:pb-0"
                  >
                    <div class="w-16 h-16 rounded-xl overflow-hidden bg-gray-50 border border-gray-100 shrink-0 relative">
                      <img 
                        :src="getCartProductImage(item.product_id || '')" 
                        :alt="item.name" 
                        class="w-full h-full object-cover" 
                      />
                    </div>
                    <div class="flex-1 min-w-0">
                      <h3 class="font-bold text-gray-900 text-sm truncate leading-snug">{{ item.name }}</h3>
                      <div class="flex flex-wrap items-center gap-x-2.5 gap-y-1 text-xs text-gray-500 mt-1 font-medium">
                        <span class="bg-gray-100 px-2 py-0.5 rounded text-[11px] font-semibold text-gray-700">Qty: {{ item.quantity }}</span>
                        <span>•</span>
                        <span>{{ formatRupiah(item.price) }}</span>
                        <span v-if="(item as any).size" class="text-gray-400">• Size: {{ (item as any).size }}</span>
                        <div v-if="(item as any).color" class="flex items-center gap-1 text-gray-400">
                          <span>• Color:</span>
                          <span :style="{ backgroundColor: (item as any).color }" class="w-2.5 h-2.5 rounded-full border border-gray-300 inline-block"></span>
                        </div>
                        <span v-if="(item as any).custom_design || item.product_variant_type === 'custom'" class="bg-emerald-50 text-emerald-700 font-bold text-[10px] px-2 py-0.5 rounded-full border border-emerald-200">
                          Custom Design
                        </span>
                      </div>
                    </div>
                    <div class="text-sm sm:text-base font-extrabold text-gray-900 text-right shrink-0">
                      {{ formatRupiah(item.subtotal) }}
                    </div>
                  </div>
                </div>

                <!-- Courier Dropdown Section -->
                <div class="bg-gray-50/70 rounded-2xl p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-3.5 border border-gray-200/80">
                  <div class="flex items-center gap-2 text-xs font-bold text-gray-700">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-[#1b4332] shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 18.75a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 01-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 00-3.213-9.193 2.056 2.056 0 00-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.948c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75" />
                    </svg>
                    <span>Delivery Courier Service:</span>
                  </div>

                  <div class="flex items-center w-full sm:w-auto">
                    <select 
                      :value="getCourierSelectValue(shop.shop_id)"
                      @change="handleCourierChange(shop.shop_id, $event)"
                      :disabled="isLoadingCalculate || isLoadingCheckout"
                      class="w-full sm:w-auto bg-white border border-gray-200 rounded-xl text-xs sm:text-sm px-4 py-3 outline-none focus:border-[#4ade80] focus:ring-2 focus:ring-[#4ade80]/20 transition-all font-bold cursor-pointer text-gray-800 disabled:opacity-50 disabled:cursor-not-allowed shadow-2xs"
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
                <div class="w-14 h-14 rounded-xl overflow-hidden bg-gray-50 border border-gray-100 shrink-0">
                  <img :src="item.image || '/images/custom-preview.png'" :alt="item.name" class="w-full h-full object-cover" />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-bold text-gray-900 text-sm truncate">{{ item.name }}</h3>
                  <p class="text-xs text-gray-400 mt-1 font-semibold">Qty: {{ item.quantity }} | Size: {{ item.size || 'Standard' }}</p>
                </div>
                <div class="text-sm font-extrabold text-gray-900 text-right">
                  {{ formatRupiah(item.price * item.quantity) }}
                </div>
              </div>
            </div>
          </div>

          <!-- STEP 3: Payment Method -->
          <div class="bg-white border border-gray-100 rounded-3xl p-6 sm:p-8 shadow-xs space-y-6">
            <div class="flex items-center gap-3 border-b border-gray-100 pb-4">
              <span class="w-7 h-7 rounded-full bg-[#1b4332] text-white flex items-center justify-center text-xs font-bold shrink-0">3</span>
              <h2 class="font-extrabold text-gray-900 text-lg">Payment Gateway Channel</h2>
            </div>
            
            <div class="space-y-4">
              <div v-if="paymentMethods.length === 0" class="text-center py-8 border-2 border-dashed border-gray-200 rounded-2xl p-4">
                <p class="text-sm text-gray-500 font-medium">No payment channels currently available.</p>
              </div>

              <div v-else class="space-y-3.5">
                <!-- Iterate over payment categories (types) -->
                <div 
                  v-for="type in paymentMethodTypes" 
                  :key="type" 
                  class="border border-gray-200 rounded-2xl overflow-hidden transition-all duration-200"
                  :class="[openedCategories[type] ? 'border-[#1b4332] shadow-2xs' : 'hover:border-gray-300']"
                >
                  <!-- Accordion Header -->
                  <button
                    type="button"
                    @click="toggleCategory(type)"
                    :disabled="isLoadingCalculate || isLoadingCheckout"
                    class="w-full flex items-center justify-between p-4 bg-gray-50/70 hover:bg-gray-100/70 transition-all font-bold text-xs sm:text-sm text-gray-800 text-left outline-none cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    <div class="flex items-center gap-2.5">
                      <!-- Category Icons -->
                      <div class="w-7 h-7 rounded-lg bg-white border border-gray-200 text-[#1b4332] flex items-center justify-center shrink-0">
                        <svg v-if="type === 'ewallet'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3" />
                        </svg>
                        <svg v-else-if="type === 'bank_transfer'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M12 21v-8.25M15.75 21v-8.25M8.25 21v-8.25M3 9l9-6 9 6m-1.5 12V10.5m-15 10.5V10.5M3 21h18M3 9h18" />
                        </svg>
                        <svg v-else-if="type === 'qr_code'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5A1.125 1.125 0 013.75 9.375v-4.5zM3.75 14.625c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5a1.125 1.125 0 01-1.125-1.125v-4.5zM13.5 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5A1.125 1.125 0 0113.5 9.375v-4.5z" />
                        </svg>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />
                        </svg>
                      </div>
                      <span>{{ formatPaymentType(type) }}</span>
                    </div>

                    <svg 
                      xmlns="http://www.w3.org/2000/svg" 
                      class="w-4 h-4 text-gray-400 transition-transform duration-200" 
                      :class="[openedCategories[type] ? 'rotate-180 text-gray-700' : '']"
                      fill="none" 
                      viewBox="0 0 24 24" 
                      stroke="currentColor" 
                      stroke-width="2.5"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                    </svg>
                  </button>

                  <!-- Accordion Content: Vertical Flat Radio List -->
                  <div 
                    v-show="openedCategories[type]" 
                    class="p-4 bg-white border-t border-gray-100 flex flex-col space-y-2.5"
                  >
                    <label 
                      v-for="method in paymentMethodsGroupedByType[type]" 
                      :key="method.id"
                      :class="[
                        'border rounded-xl px-4 py-3 flex items-center justify-between gap-3.5 cursor-pointer transition-all duration-200 select-none',
                        selectedPaymentMethodId === method.id 
                          ? 'border-[#1b4332] bg-emerald-50/25 ring-1 ring-[#1b4332] shadow-2xs' 
                          : 'border-gray-200 hover:border-gray-300 bg-white hover:bg-gray-50/50',
                        (isLoadingCalculate || isLoadingCheckout) ? 'opacity-60 pointer-events-none cursor-not-allowed' : ''
                      ]"
                    >
                      <div class="flex items-center gap-3">
                        <div class="relative flex items-center justify-center shrink-0">
                          <input 
                            type="radio" 
                            v-model="selectedPaymentMethodId" 
                            :value="method.id" 
                            :disabled="isLoadingCalculate || isLoadingCheckout"
                            class="sr-only" 
                          />
                          <div 
                            class="w-4.5 h-4.5 rounded-full border transition-all flex items-center justify-center"
                            :class="selectedPaymentMethodId === method.id ? 'border-[#1b4332] bg-[#1b4332]' : 'border-gray-300 bg-white'"
                          >
                            <div v-if="selectedPaymentMethodId === method.id" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                          </div>
                        </div>
                        <span class="font-bold text-gray-900 text-xs sm:text-sm">{{ method.name }}</span>
                      </div>

                      <div class="text-right">
                        <span v-if="method.fee > 0" class="text-emerald-800 font-bold text-xs bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-100">
                          Fee: {{ formatRupiah(method.fee) }}
                        </span>
                        <span v-else class="text-emerald-800 font-bold text-[11px] bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-100">
                          Free Fee
                        </span>
                      </div>
                    </label>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column: Sidebar Summary Tagihan -->
        <div class="lg:col-span-4 space-y-6 sticky top-24">
          <div class="bg-white border border-gray-100 rounded-3xl p-6 sm:p-7 shadow-xs space-y-6">
            <h2 class="font-extrabold text-gray-900 text-lg border-b border-gray-100 pb-4">Billing Summary</h2>
            
            <div class="space-y-3.5 text-xs sm:text-sm font-medium text-gray-600">
              <div class="flex justify-between items-center">
                <span>Items Subtotal</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(liveSubtotal) }}</span>
              </div>
              
              <div class="flex justify-between items-center">
                <span>Shipping Cost</span>
                <span class="text-gray-900 font-bold">{{ formatRupiah(liveShippingFee) }}</span>
              </div>

              <div class="flex justify-between items-center">
                <span>Payment Processing Fee</span>
                <span class="text-gray-900 font-bold">
                  {{ livePaymentFee > 0 ? formatRupiah(livePaymentFee) : 'Free' }}
                </span>
              </div>
              
              <div class="flex justify-between items-center text-emerald-700" v-if="discount > 0">
                <span>Promo Discount</span>
                <span class="font-bold">-{{ formatRupiah(discount) }}</span>
              </div>
              
              <div class="border-t border-gray-100 pt-4 flex justify-between items-baseline">
                <span class="text-sm font-bold text-gray-900">Total Bill</span>
                <span class="text-2xl font-black text-[#1b4332] tracking-tight">
                  {{ formatRupiah(liveTotalPayment) }}
                </span>
              </div>
            </div>

            <!-- Primary CTA Action Button using CButton -->
            <CButton 
              @click="handlePlaceOrder"
              :disabled="isProcessing || addressVm.addresses.value.length === 0 || isLoadingCalculate || paymentMethods.length === 0 || !selectedPaymentMethodId"
              :loading="isProcessing"
              variant="primary"
              size="lg"
              full-width
              class="shadow-xs hover:shadow"
            >
              <span>Confirm & Pay Now</span>
            </CButton>
          </div>
          
          <!-- Cryptographic Security Card -->
          <div class="bg-emerald-50/60 border border-emerald-100 rounded-2xl p-4 flex gap-3.5 items-start text-emerald-900">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-emerald-700 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
            </svg>
            <p class="text-xs font-medium leading-relaxed">
              Your payment request is cryptographically protected and processed securely with verified payment gateway channels.
            </p>
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