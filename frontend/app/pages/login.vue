<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

definePageMeta({
  alias: ['/register']
})

useSeoMeta({
  title: 'Koko - Account Login & Registration',
  description: 'Sign in or create a Koko account to build your personal watchlist and customize your anime profile.'
})

const auth = useAuth()
const router = useRouter()
const route = useRoute()

const isLoginMode = ref(route.path !== '/register')
const loginEmail = ref('')
const loginPassword = ref('')

const registerUsername = ref('')
const registerEmail = ref('')
const registerPassword = ref('')
const registerConfirmPassword = ref('')

const errorMsg = ref('')
const formLoading = ref(false)

const toggleMode = () => {
  router.push(isLoginMode.value ? '/register' : '/login')
}

watch(() => route.path, (newPath) => {
  isLoginMode.value = newPath !== '/register'
  errorMsg.value = ''
})

const handleLogin = async () => {
  if (!loginEmail.value || !loginPassword.value) {
    errorMsg.value = 'Please enter both your email/username and password.'
    return
  }
  errorMsg.value = ''
  formLoading.value = true

  const res = await auth.login(loginEmail.value, loginPassword.value)
  formLoading.value = false

  if (res.success) {
    router.push('/profile')
  } else {
    errorMsg.value = res.error || 'Login failed.'
  }
}

const handleRegister = async () => {
  if (!registerUsername.value || !registerEmail.value || !registerPassword.value) {
    errorMsg.value = 'All fields are required.'
    return
  }
  if (registerPassword.value.length < 8) {
    errorMsg.value = 'Password must be at least 8 characters long.'
    return
  }
  if (registerPassword.value !== registerConfirmPassword.value) {
    errorMsg.value = 'Passwords do not match.'
    return
  }
  errorMsg.value = ''
  formLoading.value = true

  const res = await auth.register(
    registerUsername.value.trim(),
    registerEmail.value.trim(),
    registerPassword.value
  )
  formLoading.value = false

  if (res.success) {
    router.push('/profile')
  } else {
    errorMsg.value = res.error || 'Registration failed.'
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-80px)] flex items-center justify-center p-4 md:p-8 bg-default">
    <div class="w-full max-w-5xl bg-elevated rounded-3xl overflow-hidden shadow-2xl border border-muted/50 flex lg:flex-row flex-col-reverse min-h-[600px]">
      
      <!-- Form Panel -->
      <div class="flex-1 p-8 md:p-12 flex flex-col justify-center bg-default lg:w-1/2">
        <div class="mb-8">
          <h1 class="text-3xl font-semibold tracking-tighter text-highlighted mb-2">
            {{ isLoginMode ? 'Welcome Back' : 'Create an Account' }}
          </h1>
          <p class="text-toned text-sm">
            {{ isLoginMode ? 'Sign in to access your watchlists and customized recommendations.' : 'Join Koko to keep track of your favorite anime series and movies.' }}
          </p>
        </div>

        <!-- Alert Error Message -->
        <div v-if="errorMsg" class="mb-6 p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs flex items-center gap-3">
          <UIcon name="i-solar-danger-circle-bold-duotone" class="w-5 h-5 flex-shrink-0" />
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Login Form -->
        <form v-if="isLoginMode" @submit.prevent="handleLogin" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <label for="login-email" class="text-xs font-semibold text-highlighted">Email or Username</label>
            <div class="relative group">
              <UIcon name="i-solar-letter-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="login-email"
                v-model="loginEmail"
                type="text"
                required
                placeholder="Enter your email or username"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label for="login-password" class="text-xs font-semibold text-highlighted">Password</label>
            </div>
            <div class="relative group">
              <UIcon name="i-solar-lock-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="login-password"
                v-model="loginPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <UButton
            type="submit"
            color="primary"
            class="w-full py-3 mt-4 rounded-xl flex items-center justify-center font-medium shadow-lg shadow-primary/10"
            :loading="formLoading"
          >
            Sign In
          </UButton>
        </form>

        <!-- Register Form -->
        <form v-else @submit.prevent="handleRegister" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <label for="register-username" class="text-xs font-semibold text-highlighted">Username</label>
            <div class="relative group">
              <UIcon name="i-solar-user-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="register-username"
                v-model="registerUsername"
                type="text"
                required
                placeholder="Choose a username"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-email" class="text-xs font-semibold text-highlighted">Email Address</label>
            <div class="relative group">
              <UIcon name="i-solar-letter-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="register-email"
                v-model="registerEmail"
                type="email"
                required
                placeholder="Enter your email"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-password" class="text-xs font-semibold text-highlighted">Password (min 8 chars)</label>
            <div class="relative group">
              <UIcon name="i-solar-lock-keyhole-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="register-password"
                v-model="registerPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-confirm" class="text-xs font-semibold text-highlighted">Confirm Password</label>
            <div class="relative group">
              <UIcon name="i-solar-lock-keyhole-bold-duotone" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors" />
              <input
                id="register-confirm"
                v-model="registerConfirmPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 bg-elevated rounded-xl text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
              />
            </div>
          </div>

          <UButton
            type="submit"
            color="primary"
            class="w-full py-3 mt-4 rounded-xl flex items-center justify-center font-medium shadow-lg shadow-primary/10"
            :loading="formLoading"
          >
            Create Account
          </UButton>
        </form>

        <div class="mt-8 text-center text-sm text-toned">
          <span>{{ isLoginMode ? "Don't have an account?" : "Already have an account?" }}</span>
          <button @click="toggleMode" class="ml-1 text-primary hover:underline font-semibold cursor-pointer focus:outline-none">
            {{ isLoginMode ? 'Register' : 'Sign In' }}
          </button>
        </div>
      </div>

      <!-- Graphic/Design Panel -->
      <div class="flex-1 relative bg-black/40 min-h-[300px] lg:w-1/2 flex flex-col justify-end p-8 md:p-12 overflow-hidden">
        <NuxtImg
          src="https://images.unsplash.com/photo-1578632767115-351597cf2477?q=80&w=600&auto=format&fit=crop"
          alt="Anime Wallpaper Background"
          class="absolute inset-0 w-full h-full object-cover opacity-60 z-0"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black via-black/45 to-black/20 z-10" />
        <div class="relative z-20 max-w-md">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-1.5 h-6 bg-primary rounded-full" />
            <span class="text-2xl font-bold tracking-tighter text-white">KoKo Anime</span>
          </div>
          <p class="text-white/95 text-xl font-medium tracking-tight mb-2 leading-snug">
            "Your ultimate companion to discover, explore, and track anime movies and series."
          </p>
          <p class="text-white/70 text-xs font-light">
            Keep customized watchlists, view episode releases, rate series, and manage a personalized otakus profile.
          </p>
        </div>
      </div>

    </div>
  </div>
</template>
