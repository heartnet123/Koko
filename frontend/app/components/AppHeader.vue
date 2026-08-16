<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { LocationQueryValue } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

const route = useRoute()
const router = useRouter()
const search = ref('')
const colorMode = useColorMode()
let debounceTimer: ReturnType<typeof setTimeout> | null = null
let pendingNavigations = 0

const auth = useAuth()

const toggleColorMode = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}

const getQueryString = (q: LocationQueryValue | LocationQueryValue[]): string => {
  if (!q) return ''
  return String(Array.isArray(q) ? q[0] : q)
}

// Initialize search from route query
onMounted(() => {
  if (route.query.q) {
    search.value = getQueryString(route.query.q)
  }
})

// Sync from route back to input
watch(() => route.query.q, (newQ) => {
  if (pendingNavigations > 0) return
  const qStr = getQueryString(newQ)
  if (qStr !== search.value) {
    search.value = qStr
  }
})

// Debounce user input and navigate
watch(search, (newValue) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  
  const currentQ = getQueryString(route.query.q)
  if (newValue === currentQ || (!newValue && !currentQ)) return

  const navigate = async () => {
    pendingNavigations++
    const query = { ...route.query }
    if (newValue) {
      query.q = newValue
    } else {
      query.q = undefined
    }
    query.page = undefined
    
    try {
      await router.push({ path: '/browse', query })
    } catch {
      // Ignore navigation abort errors
    } finally {
      pendingNavigations--
    }
  }

  if (route.path !== '/browse' && pendingNavigations === 0) {
    navigate()
  } else {
    debounceTimer = setTimeout(navigate, 450)
  }
})

const handleLogout = async () => {
  await auth.logout()
  router.push('/login')
}

// Dropdown items definition for Nuxt UI
const dropdownItems = computed(() => [
  [
    {
      label: auth.user.value ? `@${auth.user.value.username}` : 'Guest',
      disabled: true
    }
  ],
  [
    {
      label: 'My Profile',
      icon: 'i-solar-user-linear',
      to: '/profile'
    },
    {
      label: 'Watchlist',
      icon: 'i-solar-bookmark-linear',
      to: '/watchlist'
    },
    {
      label: 'Settings',
      icon: 'i-solar-settings-linear',
      to: '/settings'
    }
  ],
  [
    {
      label: 'Logout',
      icon: 'i-solar-logout-linear',
      onSelect: handleLogout
    }
  ]
])

defineEmits<{
  (e: 'toggle-menu'): void
}>()
</script>

<template>
  <header class="h-20 px-4 md:px-8 flex items-center justify-between sticky top-0 glass-surface z-40 border-b border-[var(--glass-border-subtle)] transition-all">
    <!-- Left Menu Trigger & Search -->
    <div class="flex items-center gap-3 w-full max-w-lg">
      <UButton
        icon="i-solar-hamburger-menu-linear"
        color="neutral"
        variant="ghost"
        class="lg:hidden rounded-xl cursor-pointer hover:bg-white/10"
        aria-label="Open menu"
        @click="$emit('toggle-menu')"
      />
      <div class="relative flex-1 group">
        <UIcon
          name="i-solar-magnifer-linear"
          class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--ui-text-toned)] group-focus-within:text-primary-400 transition-colors pointer-events-none"
        />
        <input
          v-model="search"
          type="search"
          placeholder="Search anime, genres, studios..."
          aria-label="Search anime, genres, studios"
          class="w-full pl-11 pr-16 py-2.5 glass-pill rounded-xl text-sm font-semibold text-[var(--ui-text-highlighted)] placeholder:text-[var(--ui-text-toned)] placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-primary-500/40 focus:border-primary-500/40 transition-all"
        />
        <div class="absolute right-3 top-1/2 -translate-y-1/2 hidden sm:flex items-center gap-0.5 text-[10px] font-mono text-[var(--ui-text-toned)] bg-[var(--ui-overlay)]/20 dark:bg-white/10 px-1.5 py-0.5 rounded border border-white/10 pointer-events-none">
          <span class="text-[9px]">⌘</span>K
        </div>
      </div>
    </div>

    <!-- User actions & Theme toggle -->
    <div class="flex items-center gap-2.5 ml-4">
      <!-- Dark mode toggle -->
      <UButton
        :icon="colorMode.value === 'dark' ? 'i-solar-sun-bold-duotone' : 'i-solar-moon-bold-duotone'"
        color="neutral"
        variant="ghost"
        class="rounded-xl cursor-pointer hover:bg-white/10"
        aria-label="Toggle theme"
        @click="toggleColorMode"
      />

      <UButton
        icon="i-solar-bell-linear"
        color="neutral"
        variant="ghost"
        size="md"
        class="rounded-xl hover:bg-white/10"
        aria-label="Notifications"
      />
      
      <ClientOnly>
        <!-- Logged In User Avatar Dropdown -->
        <UDropdownMenu v-if="auth.isAuthenticated.value" :items="dropdownItems" :content="{ align: 'end', side: 'bottom' }">
          <UAvatar
            :src="auth.user.value?.avatar_url || 'https://i.pravatar.cc/40'"
            alt="User avatar"
            size="sm"
            class="cursor-pointer ring-2 ring-primary-500/30 hover:ring-primary-500  transition-all rounded-xl"
          />
        </UDropdownMenu>

        <!-- Logged Out Sign In Button -->
        <UButton
          v-else
          to="/login"
          label="Sign In"
          color="primary"
          size="sm"
          icon="i-solar-login-linear"
          class="rounded-xl shadow-lg shadow-primary-500/20 cursor-pointer font-bold px-4"
        />
      </ClientOnly>
    </div>
  </header>
</template>