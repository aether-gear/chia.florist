<!-- app/pages/profile/[[tab]].vue -->
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useAuthViewModel } from '~/composables/viewmodels/useAuthViewModel'
import { useAddress } from '~/composables/useAddress'
import { useOrders } from '~/composables/useOrders'
import type { OrderTab } from '~/composables/useOrders'
import { useGlobalAlert } from '~/composables/useGlobalAlert'
import { supabaseService } from '~/services/supabaseService'
import { orderService } from '~/services/orderService'
import type { UserAddress } from '~/types/address'
import type { BackendOrder } from '~/types/order'
import { mapErrorMessage } from '~/utils/errorMessages'

definePageMeta({
  middleware: 'auth'
})

useHead({ title: 'My Profile - Chia Florist' })

const route = useRoute()
const authVm = useAuthViewModel()
const addressVm = useAddress()
const globalAlert = useGlobalAlert()

const allowedTabs = ['personal', 'addresses', 'orders']
const activeTab = computed(() => {
  const t = (route.params.tab as string)?.toLowerCase()
  return allowedTabs.includes(t) ? t : 'personal'
})

const switchTab = (tabId: string) => {
  if (tabId === 'personal') {
    navigateTo('/profile/personal')
  } else {
    navigateTo(`/profile/${tabId}`)
  }
}

const initialStatus = ((route.query.status as string)?.toLowerCase() as OrderTab)
const activeOrderStatus = ref<OrderTab>(
  ['all', 'pending', 'processing', 'shipping', 'completed', 'cancelled'].includes(initialStatus)
    ? initialStatus
    : 'all'
)

const setOrderStatus = (statusId: OrderTab) => {
  activeOrderStatus.value = statusId
  navigateTo({
    path: '/profile/orders',
    query: statusId === 'all' ? {} : { status: statusId }
  })
}

const statusLabels: Record<OrderTab, string> = {
  all: 'All Orders',
  pending: 'To Pay',
  processing: 'To Ship',
  shipping: 'To Receive',
  completed: 'Completed',
  cancelled: 'Cancelled / Expired'
}

// Edit address states
const showAddressForm = ref(false)
const editingAddress = ref<UserAddress | null>(null)
const addressFormError = ref<string | null>(null)

const provincesList = ref<any[]>([])
const citiesList = ref<any[]>([])
const districtsList = ref<any[]>([])
const villagesList = ref<any[]>([])

// User basic info state
const user = ref({
  name: 'Loading...',
  username: 'Loading...',
  email: 'Loading...',
  phone: 'Loading...',
  lastLoginAt: '' as string | null
})

// Supabase avatar states
const avatarUrl = ref<string | null>(null)
const publicAvatarUrl = ref<string | null>(null)
const signedAvatarUrl = ref<string | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

// ─── Cropper ────────────────────────────────────────────────────────────────
const showCropper = ref(false)
const cropSrc   = ref('')
const cropFilename = ref('avatar.jpg')
const cropMime     = ref('image/jpeg')
const cropImgEl    = ref<HTMLImageElement | null>(null)

// Container dimensions (responsive to screen width)
const contW = ref(440)
const contH = ref(320)
const MIN_CROP = 36

const updateCropperContainerSize = () => {
  if (typeof window === 'undefined') return
  const maxW = Math.min(window.innerWidth - 48, 440)
  contW.value = Math.max(260, Math.round(maxW))
  contH.value = Math.round(contW.value * 0.72)
}

// Natural image size
const naturalW = ref(0)
const naturalH = ref(0)

// Zoom & image-pan
const cropZoom   = ref(1)
const baseScale  = ref(1)
const imgX = ref(0)  // pan offset from container centre
const imgY = ref(0)

const scale = computed(() => baseScale.value * cropZoom.value)
const maxImgDX = computed(() => Math.max(0, (naturalW.value * scale.value - contW.value) / 2))
const maxImgDY = computed(() => Math.max(0, (naturalH.value * scale.value - contH.value) / 2))

const imgDisplayStyle = computed(() => ({
  width: `${naturalW.value * scale.value}px`,
  height: `${naturalH.value * scale.value}px`,
  maxWidth: 'none',
  maxHeight: 'none',
  position: 'absolute' as const,
  left: '50%',
  top: '50%',
  transform: `translate(-50%, -50%) translate(${imgX.value}px, ${imgY.value}px)`,
  userSelect: 'none' as const,
  pointerEvents: 'none' as const,
  draggable: false
}))

function clampImg() {
  imgX.value = Math.max(-maxImgDX.value, Math.min(maxImgDX.value, imgX.value))
  imgY.value = Math.max(-maxImgDY.value, Math.min(maxImgDY.value, imgY.value))
}

// Crop box (in container pixel coords)
const cbX = ref(40)
const cbY = ref(30)
const cbW = ref(200)
const cbH = ref(200)

// Computed: the four clipping rects that dim outside the crop box
const dimTop    = computed(() => ({ top: 0, left: 0, width: '100%', height: `${cbY.value}px` }))
const dimBottom = computed(() => ({ top: `${cbY.value + cbH.value}px`, left: 0, width: '100%', bottom: 0 }))
const dimLeft   = computed(() => ({ top: `${cbY.value}px`, left: 0, width: `${cbX.value}px`, height: `${cbH.value}px` }))
const dimRight  = computed(() => ({ top: `${cbY.value}px`, left: `${cbX.value + cbW.value}px`, right: 0, height: `${cbH.value}px` }))

// Drag state
type DragMode = 'img' | 'move' | 'nw' | 'n' | 'ne' | 'e' | 'se' | 's' | 'sw' | 'w' | null
const dragMode     = ref<DragMode>(null)
const dragStartX   = ref(0)
const dragStartY   = ref(0)
const initCbX = ref(0); const initCbY = ref(0)
const initCbW = ref(0); const initCbH = ref(0)
const lastPanX = ref(0); const lastPanY = ref(0)

function onImgLoad(e: Event) {
  updateCropperContainerSize()
  const img = e.target as HTMLImageElement
  naturalW.value = img.naturalWidth
  naturalH.value = img.naturalHeight
  baseScale.value = Math.max(contW.value / img.naturalWidth, contH.value / img.naturalHeight)
  cropZoom.value = 1
  imgX.value = 0; imgY.value = 0
  const side = Math.round(Math.min(contW.value, contH.value) * 0.75)
  cbW.value = side; cbH.value = side
  cbX.value = Math.round((contW.value - side) / 2)
  cbY.value = Math.round((contH.value - side) / 2)
}

function startCropHandle(e: MouseEvent | TouchEvent, mode: DragMode) {
  e.stopPropagation()
  if (e.cancelable) e.preventDefault()
  const cx = e instanceof MouseEvent ? e.clientX : (e.touches[0]?.clientX || 0)
  const cy = e instanceof MouseEvent ? e.clientY : (e.touches[0]?.clientY || 0)
  dragMode.value = mode
  dragStartX.value = cx; dragStartY.value = cy
  initCbX.value = cbX.value; initCbY.value = cbY.value
  initCbW.value = cbW.value; initCbH.value = cbH.value
}

function startImgPan(e: MouseEvent | TouchEvent) {
  const cx = e instanceof MouseEvent ? e.clientX : (e.touches[0]?.clientX || 0)
  const cy = e instanceof MouseEvent ? e.clientY : (e.touches[0]?.clientY || 0)
  dragMode.value = 'img'
  lastPanX.value = cx; lastPanY.value = cy
}

function onGlobalMove(cx: number, cy: number) {
  const m = dragMode.value
  if (!m) return

  if (m === 'img') {
    imgX.value += cx - lastPanX.value
    imgY.value += cy - lastPanY.value
    lastPanX.value = cx; lastPanY.value = cy
    clampImg()
    return
  }

  const dx = cx - dragStartX.value
  const dy = cy - dragStartY.value
  let nx = initCbX.value, ny = initCbY.value, nw = initCbW.value, nh = initCbH.value

  if (m === 'move') {
    nx = Math.max(0, Math.min(contW.value - nw, nx + dx))
    ny = Math.max(0, Math.min(contH.value - nh, ny + dy))
  } else {
    let rawSide: number = nw
    if (m === 'e')  rawSide = initCbW.value + dx
    if (m === 'w')  rawSide = initCbW.value - dx
    if (m === 's')  rawSide = initCbH.value + dy
    if (m === 'n')  rawSide = initCbH.value - dy
    if (m === 'se') rawSide = Math.max(initCbW.value + dx, initCbH.value + dy)
    if (m === 'sw') rawSide = Math.max(initCbW.value - dx, initCbH.value + dy)
    if (m === 'ne') rawSide = Math.max(initCbW.value + dx, initCbH.value - dy)
    if (m === 'nw') rawSide = Math.max(initCbW.value - dx, initCbH.value - dy)
    rawSide = Math.max(MIN_CROP, rawSide)

    nw = rawSide; nh = rawSide

    if (m.includes('w')) nx = initCbX.value + initCbW.value - rawSide
    if (m.includes('n')) ny = initCbY.value + initCbH.value - rawSide

    nx = Math.max(0, nx); ny = Math.max(0, ny)
    if (nx + nw > contW.value) { nw = contW.value - nx; nh = nw }
    if (ny + nh > contH.value) { nh = contH.value - ny; nw = nh }
  }

  cbX.value = nx; cbY.value = ny; cbW.value = nw; cbH.value = nh
}

function onMouseMove(e: MouseEvent)  { onGlobalMove(e.clientX, e.clientY) }
function onMouseUp()                  { dragMode.value = null }
function onTouchMoveCrop(e: TouchEvent) {
  const touch = e.touches[0]
  if (e.touches.length === 1 && touch) {
    if (e.cancelable) e.preventDefault()
    onGlobalMove(touch.clientX, touch.clientY)
  }
}
function onTouchEndCrop()             { dragMode.value = null }

function onCropZoomChange() {
  clampImg()
  if (cbX.value + cbW.value > contW.value) cbW.value = contW.value - cbX.value
  if (cbY.value + cbH.value > contH.value) cbH.value = contH.value - cbY.value
}

function cancelCrop() {
  showCropper.value = false
  cropSrc.value = ''
  if (fileInput.value) fileInput.value.value = ''
}

const loadProfilePicture = async () => {
  const userId = authVm.currentUser.value?.id
  if (!userId) return

  const urls = await supabaseService.getAvatarUrls(userId)
  if (urls) {
    publicAvatarUrl.value = urls.publicUrl
    signedAvatarUrl.value = urls.signedUrl
    avatarUrl.value = urls.signedUrl || urls.publicUrl
    
    if (authVm.currentUser.value) {
      authVm.currentUser.value.avatarUrl = urls.signedUrl || urls.publicUrl
    }
  } else {
    avatarUrl.value = null
    publicAvatarUrl.value = null
    signedAvatarUrl.value = null
    if (authVm.currentUser.value) {
      authVm.currentUser.value.avatarUrl = null
    }
  }
}

const triggerFileSelect = () => {
  fileInput.value?.click()
}

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select a valid image file.'
    return
  }
  if (file.size > 10 * 1024 * 1024) {
    uploadError.value = 'Image size must be less than 10 MB.'
    return
  }
  uploadError.value = null

  cropFilename.value = file.name || 'avatar.jpg'
  cropMime.value = file.type || 'image/jpeg'

  const reader = new FileReader()
  reader.onload = (e) => {
    updateCropperContainerSize()
    cropSrc.value = e.target?.result as string
    showCropper.value = true
  }
  reader.readAsDataURL(file)
}

const performCropAndUpload = async () => {
  const img = cropImgEl.value
  if (!img) return

  const userId = authVm.currentUser.value?.id
  if (!userId) { uploadError.value = 'Not authenticated.'; return }

  const canvas = document.createElement('canvas')
  canvas.width = 800
  canvas.height = 800
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const s = scale.value
  const imgLeft = contW.value / 2 - (naturalW.value * s) / 2 + imgX.value
  const imgTop  = contH.value / 2 - (naturalH.value * s) / 2 + imgY.value

  const srcX = (cbX.value - imgLeft) / s
  const srcY = (cbY.value - imgTop) / s
  const srcW = cbW.value / s
  const srcH = cbH.value / s

  ctx.drawImage(img, srcX, srcY, srcW, srcH, 0, 0, 800, 800)

  canvas.toBlob(async (blob) => {
    if (!blob) { uploadError.value = 'Crop failed.'; return }

    showCropper.value = false
    cropSrc.value = ''
    if (fileInput.value) fileInput.value.value = ''

    isUploading.value = true
    uploadError.value = null
    try {
      const croppedFile = new File([blob], cropFilename.value, { type: cropMime.value })
      const urls = await supabaseService.uploadFile(userId, croppedFile)
      if (urls) {
        publicAvatarUrl.value = urls.publicUrl
        signedAvatarUrl.value = urls.signedUrl
        avatarUrl.value = urls.signedUrl || urls.publicUrl
        if (authVm.currentUser.value) {
          authVm.currentUser.value.avatarUrl = urls.signedUrl || urls.publicUrl
        }
        await authVm.updateUserProfile({ avatar_url: urls.signedUrl || urls.publicUrl })
        globalAlert.showSuccess('Avatar Updated', 'Your profile picture has been updated.')
      } else {
        uploadError.value = 'Upload failed. Please try again.'
      }
    } catch (err: any) {
      uploadError.value = err.message || 'Upload error.'
    } finally {
      isUploading.value = false
    }
  }, cropMime.value)
}

const handleRemovePicture = async () => {
  const userId = authVm.currentUser.value?.id
  if (!userId) return

  if (confirm('Are you sure you want to remove your profile picture?')) {
    isUploading.value = true
    uploadError.value = null
    try {
      const files = await supabaseService.listFiles(userId)
      if (files.length > 0) {
        await supabaseService.deleteFiles(userId, files.map(f => f.name))
      }
      avatarUrl.value = null
      publicAvatarUrl.value = null
      signedAvatarUrl.value = null
      if (authVm.currentUser.value) {
        authVm.currentUser.value.avatarUrl = null
      }

      await authVm.updateUserProfile({
        avatar_url: ''
      })
      globalAlert.showSuccess('Avatar Removed', 'Your profile picture has been removed.')
    } catch (err) {
      console.error('Failed to remove profile picture:', err)
      uploadError.value = 'Failed to remove profile picture.'
    } finally {
      isUploading.value = false
    }
  }
}

const handleUpdateProfile = async () => {
  if (!user.value.name) {
    globalAlert.showWarning('Validation Error', 'Full name is required to update profile.')
    return
  }

  try {
    const res = await authVm.updateUserProfile({
      name: user.value.name,
      username: user.value.username,
      phone: user.value.phone
    })
    if (res.success) {
      globalAlert.showSuccess('Profile Updated', 'Your profile has been updated successfully!')
    } else {
      globalAlert.showError('Update Failed', mapErrorMessage(res.message, 'Failed to update profile.'))
    }
  } catch (err: any) {
    globalAlert.showError('Update Failed', mapErrorMessage(err, 'An error occurred while updating profile.'))
  }
}

// Sync profile data from ViewModel
watch(() => authVm.currentUser.value, (me) => {
  if (me) {
    user.value.name = me.name || ''
    user.value.username = me.username || ''
    user.value.email = me.email || ''
    user.value.phone = me.phone || ''
    user.value.lastLoginAt = me.last_login_at || null
  }
}, { immediate: true })

const loadAddressesData = async () => {
  await addressVm.fetchAddresses()
}

// ─── Orders ────────────────────────────────────────────────────────
const ordersVm = useOrders()
const {
  trackingData,
  isTrackingLoading,
  trackingError,
  fetchOrderTracking
} = ordersVm

// Load orders whenever the orders tab or sub-status tab is activated and scroll smoothly to top
const loadOrders = (tab: OrderTab = activeOrderStatus.value as OrderTab, page = 1) => {
  ordersVm.fetchOrders(tab, page)
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

watch(activeTab, (tab) => {
  if (tab === 'addresses') {
    loadAddressesData()
  } else if (tab === 'orders') {
    loadOrders(activeOrderStatus.value)
  }
}, { immediate: true })

watch(activeOrderStatus, (status) => {
  if (activeTab.value === 'orders') {
    loadOrders(status)
  }
})

watch(() => route.query.status, (newStatus) => {
  if (newStatus && typeof newStatus === 'string') {
    const s = newStatus.toLowerCase() as OrderTab
    if (['all', 'pending', 'processing', 'shipping', 'completed', 'cancelled'].includes(s)) {
      activeOrderStatus.value = s
    }
  }
})

onMounted(async () => {
  updateCropperContainerSize()
  window.addEventListener('resize', updateCropperContainerSize)
  if (!authVm.isInitialized.value) {
    await authVm.fetchCurrentUser()
  }
  if (!authVm.isAuthenticated.value) {
    navigateTo('/login')
  } else {
    await loadProfilePicture()
    if (activeTab.value === 'addresses') {
      loadAddressesData()
    } else if (activeTab.value === 'orders') {
      loadOrders(activeOrderStatus.value)
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCropperContainerSize)
})

// Chained location selections
const initAddressForm = async (addr: UserAddress | null = null) => {
  addressFormError.value = null
  provincesList.value = await addressVm.loadProvinces()
  
  if (addr) {
    editingAddress.value = { ...addr }
    citiesList.value = await addressVm.loadCities(addr.province_id)
    districtsList.value = await addressVm.loadDistricts(addr.city_id)
    villagesList.value = await addressVm.loadVillages(addr.district_id)
  } else {
    editingAddress.value = {
      receiver_name: '',
      phone: '',
      is_default: false,
      province_id: '',
      city_id: '',
      district_id: '',
      village_id: '',
      full_address: '',
      postal_code: ''
    }
    citiesList.value = []
    districtsList.value = []
    villagesList.value = []
  }
  showAddressForm.value = true
}

const onProvinceChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.city_id = ''
    editingAddress.value.district_id = ''
    editingAddress.value.village_id = ''
    citiesList.value = await addressVm.loadCities(editingAddress.value.province_id)
    districtsList.value = []
    villagesList.value = []
  }
}

const onCityChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.district_id = ''
    editingAddress.value.village_id = ''
    districtsList.value = await addressVm.loadDistricts(editingAddress.value.city_id)
    villagesList.value = []
  }
}

const onDistrictChange = async () => {
  if (editingAddress.value) {
    editingAddress.value.village_id = ''
    villagesList.value = await addressVm.loadVillages(editingAddress.value.district_id)
  }
}

const handleSaveAddress = async () => {
  if (!editingAddress.value) return
  const addr = editingAddress.value
  if (!addr.receiver_name || !addr.phone || !addr.province_id || !addr.city_id || !addr.district_id || !addr.village_id || !addr.full_address || !addr.postal_code) {
    addressFormError.value = 'Please fill in all the required fields.'
    return
  }
  addressFormError.value = null
  const result = await addressVm.saveAddress(addr)
  if (result.success) {
    showAddressForm.value = false
    editingAddress.value = null
    globalAlert.showSuccess('Address Saved', 'Shipping address has been saved.')
  } else {
    addressFormError.value = mapErrorMessage(result.message, 'Failed to save address.')
  }
}

const handleDeleteAddress = async (id: string) => {
  if (confirm('Are you sure you want to delete this address?')) {
    const result = await addressVm.deleteAddress(id)
    if (result.success) {
      globalAlert.showSuccess('Address Deleted', 'Shipping address has been removed.')
    } else {
      globalAlert.showError('Delete Failed', mapErrorMessage(result.message, 'Failed to delete address.'))
    }
  }
}

// Logout
const handleLogout = async () => {
  await authVm.logout()
}

const selectedOrder = ref<BackendOrder | null>(null)
const showShippingOverlay = ref(false)
const isCopied = ref(false)

const openOrderDetail = (order: BackendOrder) => {
  selectedOrder.value = order
  fetchOrderTracking(order.id)
}

const openTrackShipment = (order: BackendOrder) => {
  selectedOrder.value = order
  showShippingOverlay.value = true
  fetchOrderTracking(order.id)
}

const closeOrderDetail = () => {
  selectedOrder.value = null
  showShippingOverlay.value = false
}

watch(showShippingOverlay, (show) => {
  if (show && selectedOrder.value?.id) {
    fetchOrderTracking(selectedOrder.value.id)
  }
})

const copyTrackingNumber = (trackingNumber: string) => {
  if (!trackingNumber) return
  navigator.clipboard.writeText(trackingNumber)
  isCopied.value = true
  globalAlert.showSuccess('Copied', 'Tracking number copied to clipboard!')
  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}

const contactDriver = (orderId: string) => {
  window.open(`https://wa.me/628175234999?text=Hello%20Chia%20Florist,%20I%20would%20like%20to%20inquire%20about%20the%20delivery%20status%20for%20order%20${orderId}`, '_blank')
}

const leaveReview = (orderId: string) => {
  globalAlert.showSuccess('Review Submitted', `Thank you! Review form for order ${orderId} has been recorded. ⭐`)
}

const triggerAlert = (message: string) => {
  globalAlert.showInfo('Notice', message)
}

const handleReorder = () => {
  globalAlert.showInfo(
    'Order Expired',
    'This order was cancelled or expired. You can select an arrangement from our catalog to place a new order.',
    [
      { label: 'Browse Catalog', onClick: () => navigateTo('/catalog') },
      { label: 'Got it' }
    ]
  )
}

const isCheckingPayment = ref(false)

const handleCheckPaymentStatus = async (orderId: string) => {
  isCheckingPayment.value = true
  try {
    const res = await orderService.checkOrderPaymentStatus(orderId)
    if (res.status === 'paid') {
      globalAlert.showSuccess(
        'Payment Verified',
        'Payment verified! Your order is now being processed.',
        [
          { label: 'View Orders', onClick: () => loadOrders('all') },
          { label: 'Got it' }
        ]
      )
      closeOrderDetail()
      loadOrders('pending')
    } else {
      globalAlert.showInfo('Payment Pending', `Payment status is still pending (status: ${res.status}). If you have already transferred, please allow up to 15 minutes for automated reconciliation.`)
    }
  } catch (err: any) {
    console.error('Failed to check payment status:', err)
    globalAlert.showError('Verification Failed', mapErrorMessage(err, 'Failed to check payment status. Please try again.'))
  } finally {
    isCheckingPayment.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50/50 py-6 sm:py-12 font-sans">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <nav class="flex text-xs sm:text-sm text-gray-400 mb-6 sm:mb-8 gap-2">
        <NuxtLink to="/" class="hover:text-[#1b4332]">Home</NuxtLink>
        <span>/</span>
        <span class="text-gray-900 font-medium">My Profile</span>
        <span v-if="activeTab !== 'personal'">/</span>
        <span v-if="activeTab === 'addresses'" class="text-[#1b4332] font-semibold">Shipping Addresses</span>
        <span v-if="activeTab === 'orders'" class="text-[#1b4332] font-semibold">Order Tracking</span>
      </nav>

      <!-- Mobile Horizontal Tab Switcher -->
      <div class="lg:hidden mb-6 bg-white p-1.5 rounded-2xl shadow-xs border border-gray-100 grid grid-cols-3 gap-1">
        <button 
          v-for="tab in [{id:'personal', label:'Profile'}, {id:'addresses', label:'Addresses'}, {id:'orders', label:'Orders'}]"
          :key="tab.id"
          @click="switchTab(tab.id)"
          :class="['py-2.5 px-2 rounded-xl text-xs font-bold text-center transition-all cursor-pointer', activeTab === tab.id ? 'bg-[#1b4332] text-white shadow-xs' : 'text-gray-500 hover:bg-gray-50']"
        >
          {{ tab.label }}
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6 sm:gap-8">
        
        <aside class="lg:col-span-1 space-y-6">
          <div class="bg-white rounded-2xl p-5 sm:p-6 shadow-sm border border-gray-100">
            <div class="flex flex-col items-center text-center mb-6 lg:mb-8">
              <!-- Interactive Avatar Container -->
              <div 
                class="relative group cursor-pointer w-20 h-20 mb-4 rounded-full overflow-hidden border border-gray-200 shadow-sm flex items-center justify-center bg-gray-100" 
                @click="triggerFileSelect"
                title="Change Profile Picture"
              >
                <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105" alt="Avatar" />
                <div v-else class="text-3xl text-[#1b4332]">👤</div>

                <div class="absolute inset-0 bg-black/40 flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  <span class="text-white text-xs font-bold">Change</span>
                  <span class="text-white text-[10px]">Photo</span>
                </div>

                <div v-if="isUploading" class="absolute inset-0 bg-black/60 flex items-center justify-center">
                  <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                </div>
              </div>

              <!-- Hidden File Input -->
              <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileUpload" />

              <h2 class="font-bold text-gray-900 leading-tight">{{ user.name }}</h2>
              <p class="text-[10px] font-black text-emerald-600 uppercase tracking-widest mt-1">Silver Member</p>
            </div>

            <!-- Desktop Sidebar Tabs with sub-route links -->
            <nav class="hidden lg:block space-y-1">
              <button 
                v-for="tab in [{id:'personal', label:'Personal Information'}, {id:'addresses', label:'Shipping Addresses'}, {id:'orders', label:'Order Tracking'}]"
                :key="tab.id"
                @click="switchTab(tab.id)"
                :class="['w-full text-left px-4 py-3 rounded-xl text-sm font-semibold transition-all cursor-pointer', activeTab === tab.id ? 'bg-[#1b4332] text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50']"
              >
                {{ tab.label }}
              </button>
              
              <div class="pt-4 mt-4 border-t border-gray-100">
                <button @click="handleLogout" class="w-full text-left px-4 py-3 rounded-xl text-sm font-semibold text-red-500 hover:bg-red-50 transition-all flex items-center gap-2 cursor-pointer">
                  Logout
                </button>
              </div>
            </nav>

            <!-- Mobile Logout Quick Button -->
            <div class="lg:hidden pt-4 border-t border-gray-100">
              <button @click="handleLogout" class="w-full text-center py-2 px-3 rounded-xl text-xs font-bold text-red-500 hover:bg-red-50 transition-all cursor-pointer">
                Logout
              </button>
            </div>
          </div>
        </aside>

        <main class="lg:col-span-3 space-y-6">
          
          <!-- Personal tab -->
          <div v-if="activeTab === 'personal'" class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden animate-fade">
            <div class="px-5 sm:px-8 py-5 sm:py-6 border-b border-gray-50">
              <h3 class="font-bold text-gray-900 text-base sm:text-lg">Edit Profile</h3>
              <p class="text-xs text-gray-400 mt-0.5">Manage your personal details and contact information.</p>
            </div>
            
            <div class="p-5 sm:p-8 space-y-6">
              <!-- PROFILE PICTURE SECTION -->
              <div class="border-b border-gray-100 pb-6 space-y-4">
                <h4 class="text-xs font-bold text-gray-700 uppercase tracking-wider">Profile Picture</h4>
                <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6">
                  <div class="w-20 h-20 sm:w-24 sm:h-24 rounded-full overflow-hidden border border-gray-200 bg-gray-50 flex items-center justify-center text-3xl text-gray-400 relative flex-shrink-0">
                    <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="Profile picture" />
                    <span v-else>👤</span>
                    <div v-if="isUploading" class="absolute inset-0 bg-black/50 flex items-center justify-center">
                      <div class="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
                    </div>
                  </div>

                  <div class="space-y-3 flex-1">
                    <div class="flex flex-wrap gap-2">
                      <button @click="triggerFileSelect" :disabled="isUploading" class="bg-[#1b4332] hover:bg-[#143326] text-white px-4 py-2.5 rounded-xl text-xs font-bold shadow-sm transition disabled:opacity-50 cursor-pointer">
                        Upload New Photo
                      </button>
                      <button v-if="avatarUrl" @click="handleRemovePicture" :disabled="isUploading" class="border border-red-200 text-red-600 hover:bg-red-50 px-4 py-2.5 rounded-xl text-xs font-bold transition disabled:opacity-50 cursor-pointer">
                        Remove Photo
                      </button>
                    </div>
                    <p class="text-[10px] text-gray-400 leading-normal">Allowed JPG, PNG or WEBP. Max size 5MB.</p>
                  </div>
                </div>

                <div v-if="uploadError" class="bg-red-50 border border-red-100 rounded-xl p-3 text-red-600 text-xs font-semibold">
                  ⚠️ {{ uploadError }}
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Full Name</label>
                  <input type="text" v-model="user.name" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Username</label>
                  <input type="text" :value="user.username" disabled class="w-full px-4 py-3 bg-gray-100 border border-gray-200 rounded-xl text-sm font-semibold outline-none text-gray-500 cursor-not-allowed" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Phone Number</label>
                  <input type="text" v-model="user.phone" class="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm font-semibold outline-none focus:bg-white focus:border-[#1b4332] transition-all" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-bold text-gray-500">Email Address</label>
                  <input type="email" :value="user.email" disabled class="w-full px-4 py-3 bg-gray-100 border border-gray-200 rounded-xl text-sm font-semibold outline-none text-gray-500 cursor-not-allowed" />
                </div>
              </div>

              <div v-if="user.lastLoginAt" class="flex justify-end mt-2">
                <p class="text-[10px] text-gray-400 font-medium bg-gray-50 px-2 py-1 rounded-md border border-gray-100">Last Login: {{ new Date(user.lastLoginAt).toLocaleString() }}</p>
              </div>
              
              <div class="flex justify-end">
                <button @click="handleUpdateProfile" :disabled="authVm.isLoading.value" class="w-full sm:w-auto bg-[#1b4332] hover:bg-[#143326] text-white px-8 py-3 rounded-xl text-sm font-bold shadow-sm transition disabled:opacity-50 flex items-center justify-center gap-2 cursor-pointer">
                  <span v-if="authVm.isLoading.value" class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></span>
                  <span>Save Changes</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Shipping Addresses tab -->
          <div v-if="activeTab === 'addresses'" class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden animate-fade">
            <div class="px-5 sm:px-8 py-5 sm:py-6 border-b border-gray-50 flex flex-wrap gap-3 justify-between items-center">
              <div>
                <h3 class="font-bold text-gray-900 text-base sm:text-lg">Shipping Addresses</h3>
                <p class="text-xs text-gray-400 mt-0.5">Manage your delivery addresses for secure checkout.</p>
              </div>
              <button @click="initAddressForm()" class="bg-[#1b4332] hover:bg-[#143326] text-white px-4 py-2.5 rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
                + Add New Address
              </button>
            </div>

            <div class="p-5 sm:p-8 space-y-6">
              <!-- Loading -->
              <div v-if="addressVm.isLoading.value" class="flex flex-col items-center justify-center py-12 space-y-4">
                <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-[#1b4332]"></div>
                <p class="text-gray-500 text-xs font-medium">Loading addresses...</p>
              </div>

              <!-- Empty state -->
              <div v-else-if="addressVm.addresses.value.length === 0" class="text-center py-16 space-y-4">
                <div class="text-5xl">📍</div>
                <h4 class="font-bold text-gray-900 text-base">No Address Registered</h4>
                <p class="text-xs text-gray-400 max-w-sm mx-auto">Please add a shipping address to place orders and calculate courier fees.</p>
              </div>

              <!-- Address List -->
              <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div 
                  v-for="addr in addressVm.addresses.value" 
                  :key="addr.address_id" 
                  :class="['border rounded-2xl p-5 sm:p-6 relative flex flex-col justify-between transition-all duration-300', addr.is_default ? 'border-emerald-500 bg-emerald-50/10' : 'border-gray-200 hover:border-gray-300 bg-white']"
                >
                  <div>
                    <div class="flex items-center gap-2 mb-3">
                      <h4 class="font-bold text-gray-900 text-sm">{{ addr.receiver_name }}</h4>
                      <span v-if="addr.is_default" class="bg-emerald-100 text-emerald-800 text-[10px] font-bold px-2 py-0.5 rounded-full border border-emerald-200">Default</span>
                    </div>
                    <p class="text-xs text-gray-600 font-semibold mb-2">📞 {{ addr.phone }}</p>
                    <p class="text-xs text-gray-500 leading-relaxed">{{ addr.full_address }}</p>
                    <p class="text-[11px] text-gray-400 mt-2 font-medium">Postal Code: {{ addr.postal_code }}</p>
                  </div>

                  <div class="mt-6 pt-4 border-t border-gray-100 flex justify-end gap-3">
                    <button @click="initAddressForm(addr)" class="text-xs font-bold text-gray-600 hover:text-black transition cursor-pointer">
                      Edit
                    </button>
                    <button 
                      v-if="!addr.is_default" 
                      @click="handleDeleteAddress(addr.address_id!)" 
                      class="text-xs font-bold text-red-500 hover:text-red-600 transition cursor-pointer"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Orders tab -->
          <div v-if="activeTab === 'orders'" class="space-y-6 animate-fade">
            <!-- Order Category Tabs Bar -->
            <div class="bg-white border border-gray-100 p-1.5 rounded-2xl shadow-xs flex overflow-x-auto sm:grid sm:grid-cols-6 gap-1 text-center font-medium scrollbar-none">
              <button 
                v-for="status in [
                  { id: 'all', label: 'All' },
                  { id: 'pending', label: 'To Pay' },
                  { id: 'processing', label: 'To Ship' },
                  { id: 'shipping', label: 'To Receive' },
                  { id: 'completed', label: 'Completed' },
                  { id: 'cancelled', label: 'Cancelled' }
                ]"
                :key="status.id"
                @click="setOrderStatus(status.id as any)"
                :class="['py-2 px-3 sm:px-1 text-xs rounded-xl transition-all font-bold cursor-pointer whitespace-nowrap flex-shrink-0 sm:flex-shrink', activeOrderStatus === status.id ? 'bg-[#1b4332] text-white shadow-xs' : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50']"
              >
                {{ status.label }}
              </button>
            </div>

            <!-- Loading State -->
            <div v-if="ordersVm.isLoading.value" class="bg-white border border-gray-100 rounded-2xl p-8 flex flex-col items-center gap-3 shadow-xs">
              <div class="animate-spin rounded-full h-8 w-8 border-3 border-[#1b4332] border-t-transparent"></div>
              <p class="text-xs text-gray-400 font-semibold">Loading your orders...</p>
            </div>

            <!-- Error State -->
            <div v-else-if="ordersVm.error.value" class="bg-red-50 border border-red-100 rounded-2xl p-6 text-center shadow-xs">
              <div class="text-3xl mb-2">⚠️</div>
              <h4 class="font-bold text-red-700 text-xs">Failed to load orders</h4>
              <p class="text-[11px] text-red-500 mt-0.5">{{ ordersVm.error.value }}</p>
              <button @click="loadOrders()" class="mt-3 px-4 py-1.5 bg-red-600 text-white text-xs font-bold rounded-xl hover:bg-red-700 transition cursor-pointer">Retry</button>
            </div>

            <!-- Orders list (Clear Sleek Cards - Clickable card without View Details button) -->
            <template v-else>
              <div v-if="ordersVm.orders.value.length > 0" class="space-y-4">
                <div 
                  v-for="order in ordersVm.orders.value" 
                  :key="order.id" 
                  @click="openOrderDetail(order)"
                  class="group bg-white border border-gray-100 rounded-2xl p-4 sm:p-5 shadow-xs space-y-4 hover:shadow-md hover:border-emerald-200 transition duration-300 cursor-pointer"
                >
                  <!-- Store Header & Order Info Bar -->
                  <div class="flex flex-wrap items-center justify-between gap-2.5 pb-3 border-b border-gray-100">
                    <div class="flex items-center gap-2">
                      <span class="text-sm">🏬</span>
                      <h4 class="font-bold text-gray-900 text-xs group-hover:text-[#1b4332] transition-colors">Chia Florist Workshop</h4>
                      <span class="hidden xs:inline-block text-[9px] bg-emerald-50 text-emerald-700 font-bold px-2 py-0.5 rounded-md border border-emerald-100">Official Store</span>
                      <span class="text-xs text-gray-300">|</span>
                      <span class="text-[11px] font-mono text-gray-500 font-bold">#{{ order.number }}</span>
                    </div>

                    <!-- Right side: Status Badges -->
                    <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
                      <span :class="['px-2.5 py-0.5 text-[11px] sm:text-xs font-bold rounded-full border', ordersVm.getOrderStatusBadge(order).colorClass]">
                        {{ ordersVm.getOrderStatusBadge(order).label }}
                      </span>
                      <span v-if="order.payment?.status && order.payment.status !== 'pending' && order.payment.status !== 'paid'" :class="['px-2.5 py-0.5 text-[11px] sm:text-xs font-bold rounded-full border', ordersVm.getPaymentStatusBadge(order.payment.status).colorClass]">
                        {{ ordersVm.getPaymentStatusBadge(order.payment.status).label }}
                      </span>
                      <span v-if="!ordersVm.isOrderExpired(order) && order.status === 'pending' && order.payment?.expires_at" class="text-[10px] sm:text-[11px] font-bold text-amber-600 bg-amber-50 px-2 py-0.5 rounded-md border border-amber-100">
                        ⏳ Pay in {{ ordersVm.getTimeRemaining(order.payment.expires_at) }}
                      </span>
                    </div>
                  </div>

                  <!-- Product Preview Item -->
                  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-1">
                    <div class="flex items-center gap-3 min-w-0">
                      <div class="w-12 h-12 rounded-xl bg-emerald-50 border border-emerald-100 flex-shrink-0 flex items-center justify-center text-xl">
                        🌸
                      </div>
                      <div class="min-w-0">
                        <h5 class="font-bold text-gray-900 text-xs truncate leading-snug">
                          {{ order.items[0]?.product_name }}{{ order.items.length > 1 ? ` and ${order.items.length - 1} other item(s)` : '' }}
                        </h5>
                        <p class="text-[10px] text-gray-400 mt-0.5 font-medium truncate">
                          Total Qty: {{ order.items.reduce((acc, item) => acc + item.quantity, 0) }} • Courier: {{ order.items[0]?.courier_code ? `${order.items[0].courier_code.toUpperCase()} ${order.items[0].courier_service}` : 'Standard Delivery' }}
                        </p>
                      </div>
                    </div>

                    <div class="text-left sm:text-right flex-shrink-0 pl-15 sm:pl-0">
                      <p class="text-[10px] text-gray-400 font-medium">Order Total</p>
                      <p class="font-black text-[#1b4332] text-sm mt-0.5">{{ ordersVm.formatRupiah(order.total) }}</p>
                    </div>
                  </div>

                  <!-- Card Footer: Contextual Direct Action Buttons (without redundant View Details button) -->
                  <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 pt-3 border-t border-gray-50">
                    <span class="text-[11px] text-gray-400 font-medium">Placed on {{ ordersVm.formatDate(order.created_at) }}</span>

                    <div class="flex flex-wrap items-center justify-end gap-2">
                      <button 
                        v-if="!ordersVm.isOrderExpired(order) && order.status === 'pending'"
                        @click.stop="navigateTo(`/payment?orderId=${order.id}`)"
                        class="bg-amber-500 hover:bg-amber-600 text-white px-3.5 py-2 sm:py-1.5 rounded-xl text-xs font-bold shadow-xs transition cursor-pointer"
                      >
                        Pay Now
                      </button>
                      <button 
                        v-if="['confirmed', 'processing', 'shipped', 'delivered', 'finished'].includes(order.status)"
                        @click.stop="openTrackShipment(order)"
                        class="bg-blue-600 hover:bg-blue-700 text-white px-3.5 py-2 sm:py-1.5 rounded-xl text-xs font-bold shadow-xs transition flex items-center gap-1 cursor-pointer"
                      >
                        <span>🚚</span>
                        <span>Track Package</span>
                      </button>
                      <button 
                        v-if="order.status === 'shipped' || order.status === 'delivered'"
                        @click.stop="contactDriver(order.id)"
                        class="bg-[#1b4332] hover:bg-[#143326] text-white px-3 py-2 sm:py-1.5 rounded-xl text-xs font-bold shadow-xs transition flex items-center gap-1 cursor-pointer"
                      >
                        <span>💬</span>
                        <span>Contact Courier</span>
                      </button>
                    </div>

                  </div>
                </div>
              </div>

              <div v-else class="bg-white border border-gray-100 rounded-2xl p-8 sm:p-12 text-center shadow-xs space-y-3">
                <div class="text-4xl">📑</div>
                <h4 class="font-bold text-gray-900 text-sm">No Orders Found</h4>
                <p class="text-xs text-gray-400 max-w-xs mx-auto">There are no orders in the "{{ statusLabels[activeOrderStatus] }}" category.</p>
                <button @click="navigateTo('/catalog')" class="mt-1 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold px-5 py-2.5 rounded-xl transition cursor-pointer">
                  Browse Catalog
                </button>
              </div>

              <!-- Pagination with Smooth Scroll to Top -->
              <div v-if="ordersVm.totalPages.value > 1" class="flex items-center justify-center gap-2 pt-2">
                <button
                  @click="loadOrders(activeOrderStatus as OrderTab, ordersVm.currentPage.value - 1)"
                  :disabled="ordersVm.currentPage.value <= 1"
                  class="px-3 py-1.5 text-xs font-bold rounded-xl border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition cursor-pointer"
                >← Prev</button>
                <span class="text-xs text-gray-500 font-semibold px-2">
                  Page {{ ordersVm.currentPage.value }} of {{ ordersVm.totalPages.value }}
                </span>
                <button
                  @click="loadOrders(activeOrderStatus as OrderTab, ordersVm.currentPage.value + 1)"
                  :disabled="ordersVm.currentPage.value >= ordersVm.totalPages.value"
                  class="px-3 py-1.5 text-xs font-bold rounded-xl border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition cursor-pointer"
                >Next →</button>
              </div>
            </template>
          </div>

        </main>
      </div>
    </div>

    <!-- Modal Address Form -->
    <div v-if="showAddressForm" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-3 sm:p-4">
      <div class="bg-white rounded-3xl max-w-xl w-full max-h-[90dvh] overflow-y-auto shadow-2xl border border-gray-100 flex flex-col justify-between mx-auto">
        <div class="px-5 sm:px-8 py-4 sm:py-6 border-b border-gray-100 flex justify-between items-center">
          <h3 class="font-extrabold text-gray-900 text-base">
            {{ editingAddress?.address_id ? 'Edit Address' : 'Add New Address' }}
          </h3>
          <button @click="showAddressForm = false" class="text-gray-400 hover:text-black font-bold text-lg p-1 cursor-pointer">✕</button>
        </div>

        <div class="p-5 sm:p-8 space-y-4 flex-1">
          <div v-if="addressFormError" class="bg-red-50 border border-red-100 rounded-xl p-3 text-red-600 text-xs font-semibold">
            ⚠️ {{ addressFormError }}
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Receiver Name *</label>
              <input v-model="editingAddress!.receiver_name" type="text" placeholder="e.g. John Doe" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Phone Number *</label>
              <input v-model="editingAddress!.phone" type="text" placeholder="e.g. 0812345678" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Province *</label>
              <select v-model="editingAddress!.province_id" @change="onProvinceChange" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition">
                <option value="">Select Province</option>
                <option v-for="p in provincesList" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">City *</label>
              <select v-model="editingAddress!.city_id" @change="onCityChange" :disabled="!editingAddress!.province_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select City</option>
                <option v-for="c in citiesList" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">District *</label>
              <select v-model="editingAddress!.district_id" @change="onDistrictChange" :disabled="!editingAddress!.city_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select District</option>
                <option v-for="d in districtsList" :key="d.id" :value="d.id">{{ d.name }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Village *</label>
              <select v-model="editingAddress!.village_id" :disabled="!editingAddress!.district_id" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition disabled:opacity-50">
                <option value="">Select Village</option>
                <option v-for="v in villagesList" :key="v.id" :value="v.id">{{ v.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-xs font-bold text-gray-500">Postal Code *</label>
              <input v-model="editingAddress!.postal_code" type="text" placeholder="e.g. 17131" class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition" />
            </div>
            <div class="flex items-center pt-2 sm:pt-6 gap-2">
              <input v-model="editingAddress!.is_default" type="checkbox" id="default-check" class="h-4 w-4 accent-[#1b4332] rounded focus:ring-0 cursor-pointer" />
              <label for="default-check" class="text-xs font-bold text-gray-700 cursor-pointer">Set as default address</label>
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-xs font-bold text-gray-500">Complete Address *</label>
            <textarea v-model="editingAddress!.full_address" rows="3" placeholder="Street name, house number, details..." class="w-full px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-xl focus:bg-white focus:border-[#1b4332] outline-none text-xs font-semibold transition resize-none"></textarea>
          </div>
        </div>

        <div class="px-5 sm:px-8 py-4 sm:py-6 border-t border-gray-100 flex justify-end gap-3">
          <button @click="showAddressForm = false" class="px-5 py-2.5 border border-gray-200 rounded-xl text-xs font-bold text-gray-600 hover:bg-gray-50 transition cursor-pointer">
            Cancel
          </button>
          <button @click="handleSaveAddress" class="px-5 py-2.5 bg-[#1b4332] hover:bg-[#143326] text-white rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
            Save Address
          </button>
        </div>
      </div>
    </div>

    <!-- Independent Order Detail Modal Overlay -->
    <div v-if="selectedOrder" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-3 sm:p-4">
      <div class="bg-white rounded-3xl max-w-4xl w-full max-h-[90dvh] overflow-hidden shadow-2xl border border-gray-100 flex flex-col justify-between relative animate-fade mx-auto">
        
        <!-- Header -->
        <div class="px-5 sm:px-8 py-4 sm:py-5 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
          <div>
            <div class="flex flex-wrap items-center gap-2 sm:gap-3">
              <h3 class="font-extrabold text-gray-900 text-sm sm:text-base">Order Details</h3>
              <span class="text-xs font-mono font-bold text-gray-500 bg-gray-100 px-2 py-0.5 rounded border border-gray-200">#{{ selectedOrder.number }}</span>
            </div>
            <p class="text-[11px] text-gray-400 mt-0.5">Placed on {{ ordersVm.formatDate(selectedOrder.created_at) }}</p>
          </div>
          <button @click="closeOrderDetail" class="text-gray-400 hover:text-black font-bold text-lg p-1 cursor-pointer">✕</button>
        </div>

        <!-- Body -->
        <div class="p-4 sm:p-8 space-y-5 sm:space-y-6 flex-1 overflow-y-auto custom-scrollbar">
          <!-- Current Status Banner with Live Tracking Action -->
          <div class="flex flex-wrap items-center justify-between bg-emerald-50/20 border border-emerald-100 p-3.5 sm:p-4 rounded-2xl gap-3">
            <div>
              <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wider">Current Order Status</p>
              <div class="mt-1 flex flex-wrap items-center gap-2">
                <span :class="['px-3 py-1 text-xs font-bold rounded-full border shadow-2xs', ordersVm.getOrderStatusBadge(selectedOrder).colorClass]">
                  Order: {{ ordersVm.getOrderStatusBadge(selectedOrder).label }}
                </span>
                <span v-if="selectedOrder.payment?.status" :class="['px-3 py-1 text-xs font-bold rounded-full border shadow-2xs', ordersVm.getPaymentStatusBadge(selectedOrder.payment.status).colorClass]">
                  Payment: {{ ordersVm.getPaymentStatusBadge(selectedOrder.payment.status).label }}
                </span>
              </div>
            </div>
            <button 
              v-if="['confirmed', 'processing', 'shipped', 'delivered', 'finished'].includes(selectedOrder.status)"
              @click="showShippingOverlay = true"
              class="px-4 py-2 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-xl shadow-xs transition flex items-center gap-1.5 cursor-pointer"
            >
              <span>🚚</span>
              <span>Open Tracking Overlay</span>
            </button>
          </div>

          <!-- Shipping Summary Section & Tracking Overlay Launcher -->
          <div class="bg-blue-50/30 border border-blue-100 rounded-2xl p-4 sm:p-5 space-y-4">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h4 class="text-xs font-extrabold text-blue-900 uppercase tracking-wider flex items-center gap-2">
                  <span>🚚</span>
                  <span>Shipping Information</span>
                </h4>
                <p class="text-[11px] text-gray-500 mt-1">
                  Courier: <span class="font-bold text-gray-900 uppercase">{{ trackingData?.courier || selectedOrder.items[0]?.courier_code || 'Standard Courier' }}</span>
                  <span v-if="selectedOrder.items[0]?.courier_service"> ({{ selectedOrder.items[0].courier_service }})</span>
                </p>
              </div>

              <button 
                @click="showShippingOverlay = true"
                class="px-4 sm:px-5 py-2 sm:py-2.5 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-xs transition flex items-center gap-2 cursor-pointer"
              >
                <span>📦</span>
                <span>Open Tracking Overlay</span>
              </button>
            </div>

            <!-- Resi number pill if available -->
            <div v-if="trackingData?.tracking_number" class="flex flex-wrap items-center justify-between bg-white p-3 rounded-xl border border-blue-100 text-xs gap-2">
              <div class="flex items-center gap-2">
                <span class="text-gray-400 font-semibold">Waybill / Resi:</span>
                <span class="font-mono font-bold text-gray-900">{{ trackingData.tracking_number }}</span>
              </div>
              <button @click="copyTrackingNumber(trackingData.tracking_number)" class="text-[10px] text-blue-700 bg-blue-50 hover:bg-blue-100 px-2.5 py-1 rounded-lg border border-blue-200 font-bold transition cursor-pointer">
                {{ isCopied ? 'Copied!' : 'Copy Resi' }}
              </button>
            </div>
          </div>

          <!-- Payment Information Details Box -->
          <div class="bg-amber-50/20 border border-amber-100/80 rounded-2xl p-4 sm:p-5 space-y-3">
            <h4 class="text-xs font-extrabold text-amber-900 uppercase tracking-wider flex items-center gap-2">
              <span>💳</span>
              <span>Payment Details</span>
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4 text-xs">
              <div>
                <p class="text-gray-400 font-semibold text-[10px]">PAYMENT METHOD</p>
                <p class="font-bold text-gray-800 mt-0.5">{{ selectedOrder.payment?.provider ? selectedOrder.payment.provider.toUpperCase() : 'Online Payment Gateway' }}</p>
              </div>
              <div>
                <p class="text-gray-400 font-semibold text-[10px]">PAYMENT STATUS</p>
                <p class="font-bold text-gray-800 mt-0.5">{{ ordersVm.getPaymentStatusBadge(selectedOrder.payment?.status).label }}</p>
              </div>
              <div v-if="selectedOrder.payment?.expires_at && selectedOrder.status === 'pending'">
                <p class="text-gray-400 font-semibold text-[10px]">EXPIRATION COUNTDOWN</p>
                <p class="font-bold text-amber-700 mt-0.5">⏳ {{ ordersVm.getTimeRemaining(selectedOrder.payment.expires_at) }}</p>
              </div>
              <div>
                <p class="text-gray-400 font-semibold text-[10px]">ORDER ID</p>
                <p class="font-mono font-bold text-gray-800 mt-0.5 select-all truncate">{{ selectedOrder.id }}</p>
              </div>
            </div>
          </div>

          <!-- Items list -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider">Ordered Products</h4>
            <div class="divide-y divide-gray-100 border border-gray-100 rounded-2xl overflow-hidden bg-gray-50/20">
              <div v-for="(item, idx) in selectedOrder.items" :key="idx" class="flex gap-3 sm:gap-4 items-center p-3 sm:p-4 bg-white">
                <div class="w-12 h-12 sm:w-14 sm:h-14 bg-emerald-50 rounded-xl border border-emerald-100 flex-shrink-0 flex items-center justify-center text-xl sm:text-2xl">
                  🌸
                </div>
                <div class="flex-1 min-w-0">
                  <h5 class="font-bold text-gray-900 text-xs truncate leading-snug">{{ item.product_name }}</h5>
                  <div class="flex flex-wrap gap-1.5 sm:gap-2 mt-1 text-[10px] text-gray-500 font-semibold">
                    <span class="bg-gray-100 px-2 py-0.5 rounded">Shop: {{ item.shop_name }}</span>
                    <span v-if="item.courier_code" class="bg-blue-50 text-blue-700 px-2 py-0.5 rounded">🚚 {{ item.courier_code.toUpperCase() }} {{ item.courier_service }}</span>
                  </div>
                </div>
                <div class="text-right text-xs flex-shrink-0">
                  <p class="font-bold text-gray-900">{{ ordersVm.formatRupiah(item.unit_price) }}</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Qty: {{ item.quantity }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Billing Info -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider">Billing Summary</h4>
            <div class="border border-gray-100 rounded-2xl p-4 bg-gray-50/30 text-xs space-y-3 font-semibold text-gray-600">
              <div class="flex justify-between">
                <span>Items Subtotal</span>
                <span class="text-gray-900 font-bold">{{ ordersVm.formatRupiah(selectedOrder.subtotal) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Shipping Cost</span>
                <span class="text-gray-900 font-bold">{{ ordersVm.formatRupiah(selectedOrder.shipping_fee) }}</span>
              </div>
              <div class="border-t border-gray-100 pt-3 flex justify-between text-sm font-bold text-gray-900">
                <span>Total Amount Paid</span>
                <span class="text-base text-[#1b4332] font-black">{{ ordersVm.formatRupiah(selectedOrder.total) }}</span>
              </div>
            </div>
          </div>

        </div>

        <!-- Footer Actions -->
        <div class="px-5 sm:px-8 py-4 sm:py-5 border-t border-gray-100 flex flex-wrap justify-between items-center bg-gray-50/50 gap-3">
          <button @click="closeOrderDetail" class="px-5 py-2.5 border border-gray-200 rounded-xl text-xs font-bold text-gray-600 hover:bg-gray-50 transition cursor-pointer">
            Close
          </button>

          <div class="flex flex-wrap items-center gap-2">
            <div v-if="selectedOrder.status === 'pending'" class="flex flex-wrap gap-2">
              <button 
                @click="navigateTo(`/payment?orderId=${selectedOrder.id}`)" 
                class="px-5 py-2.5 bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold rounded-xl transition cursor-pointer"
              >
                Pay Now
              </button>
              <button 
                @click="handleCheckPaymentStatus(selectedOrder.id)" 
                :disabled="isCheckingPayment"
                class="px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 text-xs font-bold rounded-xl transition cursor-pointer flex items-center gap-1.5 disabled:opacity-50"
              >
                <span v-if="isCheckingPayment" class="animate-spin rounded-full h-3.5 w-3.5 border-2 border-gray-500 border-t-transparent"></span>
                <span>Check Payment Status</span>
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'confirmed' || selectedOrder.status === 'processing'">
              <button @click="triggerAlert('Your flower board is being carefully arranged by our expert florists.')" class="px-5 py-2.5 bg-[#1b4332] text-white hover:bg-[#143326] text-xs font-bold rounded-xl transition cursor-pointer">
                View Workshop Progress
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'shipped' || selectedOrder.status === 'delivered'">
              <button @click="contactDriver(selectedOrder.id)" class="px-5 py-2.5 bg-[#1b4332] hover:bg-[#143326] text-white text-xs font-bold rounded-xl transition flex items-center gap-2 cursor-pointer">
                Contact Courier (WhatsApp)
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'finished'">
              <button @click="leaveReview(selectedOrder.id)" class="px-5 py-2.5 bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold rounded-xl transition cursor-pointer">
                Write Product Review
              </button>
            </div>
            <div v-else-if="selectedOrder.status === 'cancelled' || selectedOrder.status === 'expired' || ordersVm.isOrderExpired(selectedOrder)">
              <button @click="handleReorder" class="px-5 py-2.5 bg-rose-50 text-rose-700 border border-rose-100 hover:bg-rose-100 text-xs font-bold rounded-xl transition cursor-pointer">
                Order Again
              </button>
            </div>
          </div>
        </div>

        <!-- Shipping Detail Slide-out Overlay Panel (Nested) -->
        <div 
          v-if="showShippingOverlay" 
          class="absolute inset-y-0 right-0 w-full sm:w-[450px] bg-white shadow-2xl border-l border-gray-100 flex flex-col justify-between z-50 transform transition-transform duration-300 ease-in-out animate-slide-in"
        >
          <!-- Overlay Header -->
          <div class="px-5 sm:px-6 py-4 sm:py-5 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
            <div class="flex items-center gap-2">
              <span class="text-base">🚚</span>
              <h3 class="font-extrabold text-gray-900 text-sm">Shipping Information</h3>
            </div>
            <button @click="showShippingOverlay = false" class="text-gray-400 hover:text-black font-bold text-base p-1 cursor-pointer">✕</button>
          </div>

          <!-- Overlay Body -->
          <div class="p-5 sm:p-6 space-y-6 flex-1 overflow-y-auto custom-scrollbar text-xs">
            <div class="border border-gray-100 rounded-2xl p-4 bg-gray-50/30 space-y-2">
              <div class="flex justify-between font-semibold">
                <span class="text-gray-400">Courier Partner</span>
                <span class="text-gray-900 font-bold uppercase">
                  {{ trackingData?.courier || selectedOrder.items[0]?.courier_code || '—' }}
                  <span v-if="selectedOrder.items[0]?.courier_service" class="normal-case"> ({{ selectedOrder.items[0].courier_service }})</span>
                </span>
              </div>
              <div v-if="trackingData?.tracking_number" class="flex justify-between items-center font-semibold pt-1 border-t border-gray-100/60">
                <span class="text-gray-400">Tracking Number</span>
                <div class="flex items-center gap-1.5 font-mono font-bold text-gray-900">
                  <span>{{ trackingData.tracking_number }}</span>
                  <button @click="copyTrackingNumber(trackingData.tracking_number)" class="text-[10px] text-emerald-700 bg-emerald-50 hover:bg-emerald-100 px-2 py-0.5 rounded border border-emerald-200 transition cursor-pointer">
                    {{ isCopied ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
              </div>
              <div class="flex justify-between font-semibold">
                <span class="text-gray-400">Shipping Fee</span>
                <span class="text-gray-900 font-mono font-bold">{{ ordersVm.formatRupiah(selectedOrder.shipping_fee) }}</span>
              </div>
            </div>

            <!-- Payment Info -->
            <div class="space-y-2" v-if="selectedOrder.payment">
              <h4 class="text-[10px] font-bold text-gray-500 uppercase tracking-wider">Payment Information</h4>
              <div class="border border-gray-100 rounded-2xl p-4 space-y-2 font-semibold">
                <div class="flex justify-between">
                  <span class="text-gray-400">Provider</span>
                  <span class="text-gray-900 uppercase">{{ selectedOrder.payment.provider }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-400">Amount</span>
                  <span class="text-gray-900">{{ ordersVm.formatRupiah(selectedOrder.payment.amount) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-400">Payment Status</span>
                  <span :class="['font-bold', selectedOrder.payment.status === 'paid' ? 'text-emerald-600' : 'text-amber-600']">{{ selectedOrder.payment.status }}</span>
                </div>
              </div>
            </div>

            <!-- Shipping Timeline -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h4 class="text-[10px] font-bold text-gray-500 uppercase tracking-wider">Order &amp; Tracking Timeline</h4>
                <button 
                  v-if="selectedOrder.id" 
                  @click="fetchOrderTracking(selectedOrder.id)" 
                  class="text-[10px] text-emerald-700 hover:underline font-semibold flex items-center gap-1 cursor-pointer"
                >
                  <span>🔄</span> Refresh
                </button>
              </div>

              <!-- Loading Skeleton -->
              <div v-if="isTrackingLoading" class="p-6 text-center space-y-3">
                <div class="inline-block animate-spin rounded-full h-5 w-5 border-2 border-emerald-600 border-t-transparent"></div>
                <p class="text-gray-400 text-[11px]">Fetching latest tracking updates from courier...</p>
              </div>
              
              <!-- Dynamic Real-time Courier Timeline -->
              <div 
                v-else-if="trackingData?.timeline && trackingData.timeline.length > 0"
                class="relative pl-6 space-y-6 before:content-[''] before:absolute before:left-2 before:top-2 before:bottom-2 before:w-[2px] before:bg-emerald-100"
              >
                <div 
                  v-for="(event, idx) in trackingData.timeline" 
                  :key="idx" 
                  class="relative"
                >
                  <span 
                    class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white"
                    :class="[idx === 0 ? 'bg-emerald-500 ring-4 ring-emerald-100 animate-pulse' : 'bg-emerald-400']"
                  ></span>
                  <div class="flex items-center justify-between">
                    <p class="font-bold text-gray-900 capitalize">{{ event.status.replace(/_/g, ' ') }}</p>
                    <span class="text-[10px] font-mono text-gray-400">{{ ordersVm.formatDate(event.timestamp) }}</span>
                  </div>
                  <p class="text-[11px] text-gray-600 mt-0.5">{{ event.description }}</p>
                  <p v-if="event.location" class="text-[10px] text-gray-400 mt-0.5 flex items-center gap-1">
                    <span>📍</span> {{ event.location }}
                  </p>
                </div>
              </div>

              <!-- Fallback Status Milestones -->
              <div v-else class="relative pl-6 space-y-6 before:content-[''] before:absolute before:left-2 before:top-2 before:bottom-2 before:w-[2px] before:bg-gray-100">
                <div class="relative">
                  <span 
                    class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white flex items-center justify-center"
                    :class="[selectedOrder.status !== 'cancelled' ? 'bg-emerald-500 ring-4 ring-emerald-100' : 'bg-red-500 ring-4 ring-red-100']"
                  ></span>
                  <p class="font-bold text-gray-900">Order Placed Successfully</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">{{ ordersVm.formatDate(selectedOrder.created_at) }}</p>
                </div>

                <div v-if="selectedOrder.status !== 'pending' && selectedOrder.status !== 'cancelled'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Payment Verified by Admin</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Payment confirmed &amp; order accepted</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Payment Verification</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Awaiting payment confirmation</p>
                </div>

                <div v-if="['processing', 'shipped', 'delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Flower Arrangement Complete</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Expert florist team finished production</p>
                </div>
                <div v-else-if="selectedOrder.status === 'confirmed'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-blue-500 ring-4 ring-blue-100"></span>
                  <p class="font-bold text-gray-900">Queued for Arrangement</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Waiting for florist slot allocation</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Flower Board Production</p>
                </div>

                <div v-if="['shipped', 'delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Handed Over to Courier</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Parcel picked up by {{ selectedOrder.items[0]?.courier_code || 'courier' }} partner</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Hand Over to Courier Partner</p>
                </div>

                <div v-if="['delivered', 'finished'].includes(selectedOrder.status)" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Delivered to Destination</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Signed by recipient at destination address</p>
                </div>
                <div v-else-if="selectedOrder.status !== 'cancelled'" class="relative opacity-40">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-gray-200"></span>
                  <p class="font-bold text-gray-500">Out for Delivery / Completed Destination</p>
                </div>

                <div v-if="selectedOrder.status === 'finished'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-emerald-500 ring-4 ring-emerald-100"></span>
                  <p class="font-bold text-gray-900">Order Completed</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Transaction finalized by customer</p>
                </div>

                <div v-if="selectedOrder.status === 'cancelled'" class="relative">
                  <span class="absolute -left-[23px] top-0 w-3 h-3 rounded-full border-2 border-white bg-red-500 ring-4 ring-red-100"></span>
                  <p class="font-bold text-red-600">Order Cancelled</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">This transaction has been terminated</p>
                </div>
              </div>
            </div>

          </div>

          <!-- Overlay Footer -->
          <div class="px-5 sm:px-6 py-4 sm:py-5 border-t border-gray-100 flex justify-end bg-gray-50/50">
            <button @click="showShippingOverlay = false" class="px-5 py-2.5 bg-[#1b4332] text-white hover:bg-[#143326] rounded-xl text-xs font-bold shadow-sm transition cursor-pointer">
              Go Back to Details
            </button>
          </div>
        </div>

      </div>
    </div>

    <!-- ═══ Image Crop Modal ═══════════════════════════════════════════════ -->
    <Teleport to="body">
      <Transition name="cropper-fade">
        <div
          v-if="showCropper"
          class="fixed inset-0 z-[9999] flex items-center justify-center p-3 sm:p-4"
          style="background: rgba(0,0,0,0.35); backdrop-filter: blur(4px)"
          @mousemove="onMouseMove"
          @mouseup="onMouseUp"
          @touchmove="onTouchMoveCrop"
          @touchend="onTouchEndCrop"
          @touchcancel="onTouchEndCrop"
        >
          <div class="bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col animate-fade max-w-[96vw]" style="width: 520px">

            <!-- Header -->
            <div class="px-4 sm:px-6 pt-4 sm:pt-5 pb-3 flex items-center justify-between border-b border-gray-100">
              <div>
                <h3 class="font-extrabold text-gray-900 text-sm">Crop Profile Picture</h3>
                <p class="text-[10px] sm:text-[11px] text-gray-400 mt-0.5">Drag corners or edges to resize · Drag crop box to move · Pan outside</p>
              </div>
              <button @click="cancelCrop" class="text-gray-400 hover:text-gray-600 transition text-lg leading-none p-1 cursor-pointer">✕</button>
            </div>

            <!-- Canvas editing area -->
            <div class="relative bg-gray-900 overflow-hidden select-none"
              :style="{ width: `${contW}px`, height: `${contH}px`, margin: '0 auto', touchAction: 'none' }"
              @mousedown="startImgPan"
              @touchstart="startImgPan"
            >
              <img
                v-if="cropSrc"
                :src="cropSrc"
                :style="imgDisplayStyle"
                @load="onImgLoad"
                ref="cropImgEl"
                alt=""
              />

              <!-- Dim overlays -->
              <div class="absolute pointer-events-none" :style="{ ...dimTop,    background: 'rgba(0,0,0,0.5)' }"></div>
              <div class="absolute pointer-events-none" :style="{ ...dimBottom, background: 'rgba(0,0,0,0.5)' }"></div>
              <div class="absolute pointer-events-none" :style="{ ...dimLeft,   background: 'rgba(0,0,0,0.5)' }"></div>
              <div class="absolute pointer-events-none" :style="{ ...dimRight,  background: 'rgba(0,0,0,0.5)' }"></div>

              <!-- Crop selection box -->
              <div
                class="absolute cursor-move"
                :style="{
                  left: `${cbX}px`, top: `${cbY}px`,
                  width: `${cbW}px`, height: `${cbH}px`,
                  border: '2px solid #fff',
                  boxSizing: 'border-box',
                  touchAction: 'none'
                }"
                @mousedown.stop="startCropHandle($event, 'move')"
                @touchstart.stop="startCropHandle($event, 'move')"
              >
                <!-- Rule of thirds grid -->
                <div class="absolute inset-0 pointer-events-none" style="
                  background-image: linear-gradient(rgba(255,255,255,0.2) 1px, transparent 1px),
                                    linear-gradient(90deg, rgba(255,255,255,0.2) 1px, transparent 1px);
                  background-size: 33.33% 33.33%;
                "></div>

                <!-- 8 resize handles -->
                <div class="crop-handle crop-handle-corner" style="top:-5px;left:-5px;cursor:nw-resize"
                  @mousedown.stop="startCropHandle($event,'nw')" @touchstart.stop="startCropHandle($event,'nw')"></div>
                <div class="crop-handle crop-handle-corner" style="top:-5px;right:-5px;cursor:ne-resize"
                  @mousedown.stop="startCropHandle($event,'ne')" @touchstart.stop="startCropHandle($event,'ne')"></div>
                <div class="crop-handle crop-handle-corner" style="bottom:-5px;left:-5px;cursor:sw-resize"
                  @mousedown.stop="startCropHandle($event,'sw')" @touchstart.stop="startCropHandle($event,'sw')"></div>
                <div class="crop-handle crop-handle-corner" style="bottom:-5px;right:-5px;cursor:se-resize"
                  @mousedown.stop="startCropHandle($event,'se')" @touchstart.stop="startCropHandle($event,'se')"></div>
                
                <div class="crop-handle crop-handle-edge" style="top:-4px;left:calc(50% - 12px);cursor:n-resize"
                  @mousedown.stop="startCropHandle($event,'n')" @touchstart.stop="startCropHandle($event,'n')"></div>
                <div class="crop-handle crop-handle-edge" style="bottom:-4px;left:calc(50% - 12px);cursor:s-resize"
                  @mousedown.stop="startCropHandle($event,'s')" @touchstart.stop="startCropHandle($event,'s')"></div>
                <div class="crop-handle crop-handle-edge" style="left:-4px;top:calc(50% - 12px);cursor:w-resize"
                  @mousedown.stop="startCropHandle($event,'w')" @touchstart.stop="startCropHandle($event,'w')"></div>
                <div class="crop-handle crop-handle-edge" style="right:-4px;top:calc(50% - 12px);cursor:e-resize"
                  @mousedown.stop="startCropHandle($event,'e')" @touchstart.stop="startCropHandle($event,'e')"></div>
              </div>
            </div>

            <!-- Controls -->
            <div class="px-4 sm:px-6 py-3 sm:py-4 space-y-3">
              <div class="flex items-center gap-3">
                <span class="text-xs font-bold text-gray-400 w-10">Zoom</span>
                <input
                  type="range"
                  v-model.number="cropZoom"
                  min="1" max="4" step="0.01"
                  class="flex-1 h-1.5 rounded-full appearance-none cursor-pointer accent-[#1b4332] bg-gray-200"
                  @input="onCropZoomChange"
                />
                <span class="text-xs font-bold text-gray-500 w-10 text-right">{{ Math.round(cropZoom * 100) }}%</span>
              </div>

              <div class="flex flex-wrap items-center justify-between gap-2">
                <span class="text-[11px] text-gray-400">
                  ⬜ {{ cbW }} × {{ cbW }} px (square)
                </span>
                <div class="flex items-center gap-2">
                  <button
                    @click="cancelCrop"
                    class="px-4 py-2 border border-gray-200 rounded-xl text-xs font-bold text-gray-600 hover:bg-gray-50 transition cursor-pointer"
                  >Cancel</button>
                  <button
                    @click="performCropAndUpload"
                    :disabled="isUploading"
                    class="px-4 sm:px-5 py-2 bg-[#1b4332] hover:bg-[#143326] disabled:opacity-60 text-white text-xs font-bold rounded-xl shadow-sm transition flex items-center gap-2 cursor-pointer"
                  >
                    <span v-if="isUploading" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                    <span>{{ isUploading ? 'Uploading…' : 'Crop & Save' }}</span>
                  </button>
                </div>
              </div>
            </div>

          </div>
        </div>
      </Transition>
    </Teleport>

  </div>
</template>

<style scoped>
.animate-fade { animation: fadeIn 0.4s ease-out; }
.animate-slide-in { animation: slideIn 0.3s ease-out forwards; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes slideIn { from { transform: translateX(100%); } to { transform: translateX(0); } }
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
.scrollbar-none {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

/* Cropper modal transition */
.cropper-fade-enter-active,
.cropper-fade-leave-active { transition: opacity 0.2s ease; }
.cropper-fade-enter-from,
.cropper-fade-leave-to { opacity: 0; }

/* Crop handles with enlarged touch hit target */
.crop-handle {
  position: absolute;
  background: #fff;
  border: 2px solid #1b4332;
  border-radius: 2px;
  touch-action: none;
}
.crop-handle::before {
  content: '';
  position: absolute;
  top: -14px;
  left: -14px;
  right: -14px;
  bottom: -14px;
}
.crop-handle-corner {
  width: 10px;
  height: 10px;
}
.crop-handle-edge {
  width: 24px;
  height: 8px;
}
.crop-handle-edge[style*="left:-4px"],
.crop-handle-edge[style*="right:-4px"] {
  width: 8px;
  height: 24px;
}
</style>
