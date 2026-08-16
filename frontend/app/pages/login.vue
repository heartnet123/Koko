<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

definePageMeta({
  alias: ['/register']
})

useSeoMeta({
  title: 'Koko — Account Authentication',
  description: 'Sign in or create a KoKo account to build your personal anime watchlist and profile.'
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
  <div class="min-h-[calc(100dvh-80px)] flex items-center justify-center p-4 md:p-8 animate-fade-in-up">
    <div class="w-full max-w-5xl glass-surface-elevated rounded-3xl overflow-hidden  border border-[var(--glass-border)] flex lg:flex-row flex-col-reverse min-h-[600px]">
      
      <!-- Form Panel -->
      <div class="flex-1 p-8 md:p-12 flex flex-col justify-center lg:w-1/2">
        <div class="mb-8">
          <div class="inline-flex items-center gap-2 mb-2">
            <div class="w-2 h-4 bg-primary-500 rounded-full " />
            <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">Authentication</span>
          </div>
          <h1 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)] mb-2">
            {{ isLoginMode ? 'Welcome Back' : 'Create an Account' }}
          </h1>
          <p class="text-[var(--ui-text-toned)] text-xs md:text-sm font-normal">
            {{ isLoginMode ? 'Sign in to access your watchlists and custom recommendations.' : 'Join KoKo to catalog your favorite anime series and sync across sessions.' }}
          </p>
        </div>

        <!-- Alert Error Message -->
        <div v-if="errorMsg" class="mb-6 p-4 rounded-2xl bg-[var(--ui-error)]/10 border border-[var(--ui-error)]/30 text-[var(--ui-error)] text-xs flex items-center gap-3 font-mono">
          <UIcon name="i-solar-danger-circle-bold-duotone" class="w-5 h-5 flex-shrink-0" />
          <span>{{ errorMsg }}</span>
        </div>

        <!-- Login Form -->
        <form v-if="isLoginMode" @submit.prevent="handleLogin" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <label for="login-email" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Email or Username</label>
            <div class="relative group">
              <UIcon name="i-solar-letter-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="login-email"
                v-model="loginEmail"
                type="text"
                required
                placeholder="Enter your email or username"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="login-password" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Password</label>
            <div class="relative group">
              <UIcon name="i-solar-lock-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="login-password"
                v-model="loginPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <UButton
            type="submit"
            color="primary"
            class="w-full py-3.5 mt-4 rounded-2xl flex items-center justify-center font-bold shadow-xl shadow-primary-500/25 cursor-pointer hover:scale-[1.01] active:scale-95 transition-all"
            :loading="formLoading"
          >
            Sign In
          </UButton>
        </form>

        <!-- Register Form -->
        <form v-else @submit.prevent="handleRegister" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <label for="register-username" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Username</label>
            <div class="relative group">
              <UIcon name="i-solar-user-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="register-username"
                v-model="registerUsername"
                type="text"
                required
                placeholder="Choose a username"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-email" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Email Address</label>
            <div class="relative group">
              <UIcon name="i-solar-letter-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="register-email"
                v-model="registerEmail"
                type="email"
                required
                placeholder="Enter your email"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-password" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Password (min 8 chars)</label>
            <div class="relative group">
              <UIcon name="i-solar-lock-keyhole-linear" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="register-password"
                v-model="registerPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="register-confirm" class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Confirm Password</label>
            <div class="relative group">
              <UIcon name="i-solar-lock-keyhole-bold-duotone" class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none" />
              <input
                id="register-confirm"
                v-model="registerConfirmPassword"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-11 pr-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all font-mono"
              />
            </div>
          </div>

          <UButton
            type="submit"
            color="primary"
            class="w-full py-3.5 mt-4 rounded-2xl flex items-center justify-center font-bold shadow-xl shadow-primary-500/25 cursor-pointer hover:scale-[1.01] active:scale-95 transition-all"
            :loading="formLoading"
          >
            Create Account
          </UButton>
        </form>

        <div class="mt-8 text-center text-xs md:text-sm text-[var(--ui-text-toned)] font-semibold">
          <span>{{ isLoginMode ? "Don't have an account?" : "Already have an account?" }}</span>
          <button @click="toggleMode" class="ml-1.5 text-primary-400 hover:text-primary-300 font-bold cursor-pointer focus:outline-none underline">
            {{ isLoginMode ? 'Register Now' : 'Sign In' }}
          </button>
        </div>
      </div>

      <!-- Graphic Artwork Panel -->
      <div class="flex-1 relative bg-[var(--ui-overlay)]/60 min-h-[300px] lg:w-1/2 flex flex-col justify-end p-8 md:p-12 overflow-hidden border-b lg:border-b-0 lg:border-l border-[var(--glass-border)]">
        <NuxtImg
          src="https://images.unsplash.com/photo-1578632767115-351597cf2477?q=80&w=800&auto=format&fit=crop"
          alt="Anime Wallpaper Background"
          class="absolute inset-0 w-full h-full object-cover opacity-50 z-0 scale-105"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-[#090B10] via-[#090B10]/50 to-transparent z-10" />
        <div class="relative z-20 max-w-md">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-2 h-6 bg-gradient-to-b from-primary-400 to-primary-600 rounded-full " />
            <span class="text-2xl font-bold tracking-tight text-white">KoKo Anime</span>
          </div>
          <p class="text-white text-lg md:text-xl font-bold tracking-tight mb-2 leading-snug">
            "Your ultimate companion to discover, explore, and track anime."
          </p>
          <p class="text-[var(--ui-text-on-image-muted)] text-xs leading-relaxed font-normal">
            Keep customized watchlists, view episode releases, rate series, and manage a personalized otaku vault.
          </p>
        </div>
      </div>

    </div>
  </div>
</template>