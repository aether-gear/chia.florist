<script setup lang="ts">
import { ref, onMounted } from 'vue'

useHead({
  title: 'Verify Your Account - Chia Florist',
  meta: [
    { name: 'description', content: 'Confirm your verification code to access your Chia Florist account.' }
  ]
})

const otpCode = ref('')
const isSubmitting = ref(false)
const errorMessage = ref('')
const userEmail = ref('user@email.com') // Default fallback

onMounted(() => {
  // Mengambil email dari localStorage yang disimpan saat user register
  const savedEmail = localStorage.getItem('register_email')
  if (savedEmail) {
    userEmail.value = savedEmail
  }
})

const handleVerify = async () => {
  if (otpCode.value.length < 4) {
    errorMessage.value = 'Please enter a valid verification code.'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    // HIT INTEGRASI KE BACKEND GOLANG TEMANMU
    const response = await $fetch<{ success: boolean; message: string }>('http://localhost:8080/api/v1/auth/verify', {
      method: 'POST',
      body: {
        email: userEmail.value,
        code: otpCode.value
      }
    })

    if (response.success) {
      // Hapus data temporary email register jika sukses
      localStorage.removeItem('register_email')
      
      // Sukses: Sesuai pipeline temanmu, langsung alihkan balik ke aplikasi (homepage)
      navigateTo('/')
    } else {
      // Gagal: User stuck di halaman ini dan menampilkan error dari backend
      errorMessage.value = response.message || 'Verification failed. Please try again.'
      isSubmitting.value = false
    }
  } catch (err: any) {
    // Menangkap error jika OTP salah atau server down
    errorMessage.value = err.data?.message || 'Invalid verification code. You are locked here until verified.'
    isSubmitting.value = false
  }
}

const handleCancel = () => {
  localStorage.removeItem('register_email')
  navigateTo('/login') // Balik ke portal login jika membatalkan
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4 font-sans antialiased">
    
    <div class="w-full max-w-md bg-white border border-gray-200 rounded-3xl p-8 shadow-xl shadow-gray-100/50 transition-all duration-300">
      
      <div class="flex flex-col items-center text-center mb-8">
        <div class="w-14 h-14 bg-[#1b4332] rounded-2xl flex items-center justify-center text-white text-2xl shadow-md mb-5 font-bold tracking-tighter">
          CF
        </div>
        <h2 class="text-xl font-bold text-gray-900 tracking-tight">
          Verify Account to Chia Florist
        </h2>
        <p class="text-xs text-gray-400 mt-2 max-w-xs leading-relaxed">
          We have sent a secure verification code to <br/>
          <span class="text-gray-700 font-semibold break-all">{{ userEmail }}</span>
        </p>
      </div>

      <form @submit.prevent="handleVerify" class="space-y-6">
        
        <div class="space-y-2">
          <label class="text-[11px] font-black uppercase tracking-widest text-gray-400 block text-center">
            Enter Verification Code
          </label>
          <input 
            type="text" 
            v-model="otpCode"
            placeholder="0000"
            maxlength="6"
            class="w-full px-5 py-3.5 bg-gray-50 border border-gray-200 rounded-2xl text-center text-xl font-mono tracking-[0.5em] indent-[0.25em] outline-none focus:bg-white focus:border-[#1b4332] focus:ring-2 focus:ring-[#1b4332]/5 transition-all"
            :disabled="isSubmitting"
          />
        </div>

        <div 
          v-if="errorMessage" 
          class="bg-red-50 border border-red-100 text-red-600 text-xs font-semibold px-4 py-3 rounded-xl flex items-center gap-2 animate-shake"
        >
          <span>⚠️</span>
          <p class="flex-1 leading-normal">{{ errorMessage }}</p>
        </div>

        <div class="grid grid-cols-2 gap-4 pt-2">
          <button 
            type="button"
            @click="handleCancel"
            class="w-full py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 text-xs font-bold rounded-xl transition-all outline-none cursor-pointer"
            :disabled="isSubmitting"
          >
            Cancel
          </button>
          
          <button 
            type="submit"
            class="w-full py-3 bg-[#1b4332] hover:bg-[#143326] disabled:bg-gray-300 text-white text-xs font-bold rounded-xl shadow-sm transition-all outline-none flex items-center justify-center cursor-pointer"
            :disabled="isSubmitting"
          >
            <span v-if="isSubmitting" class="animate-pulse">Verifying...</span>
            <span v-else>Confirm</span>
          </button>
        </div>

      </form>
    </div>

    <div class="mt-8 text-center max-w-xs px-4">
      <p class="text-[10px] text-gray-400 leading-relaxed">
        By confirming this verification process, you apply your account control rules. Codes expire shortly and rate limits apply.
      </p>
    </div>

  </div>
</template>

<style scoped>
/* Animasi micro-shake jika user memasukkan kode OTP yang salah */
.animate-shake { animation: shake 0.3s ease-in-out; }
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}
</style>