<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { LocationQueryValue } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

const route = useRoute()
const router = useRouter()
const search = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null
let pendingNavigations = 0

const auth = useAuth()

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

// Sync from route back to input (e.g. user clicks back button)
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
      delete query.q
    }
    delete query.page
    
    try {
      await router.push({ path: '/browse', query })
    } catch (err) {
      // Ignore navigation abort errors
    } finally {
      pendingNavigations--
    }
  }

  if (route.path !== '/browse' && pendingNavigations === 0) {
    navigate()
  } else {
    debounceTimer = setTimeout(navigate, 500)
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
  <header class="h-20 px-4 md:px-8 flex items-center justify-between sticky top-0 bg-default/80 backdrop-blur-md z-40 border-b border-muted/50">
    <!-- Left Menu Trigger & Search -->
    <div class="flex items-center gap-3 w-full max-w-lg">
      <UButton
        icon="i-solar-hamburger-menu-linear"
        color="neutral"
        variant="ghost"
        class="lg:hidden rounded-full cursor-pointer focus:outline-none"
        aria-label="Open menu"
        @click="$emit('toggle-menu')"
      />
      <div class="relative flex-1 group">
        <UIcon
          name="i-solar-magnifer-linear"
          class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors"
        />
        <input
          v-model="search"
          type="search"
          placeholder="Search anime, genres, studios..."
          aria-label="Search anime, genres, studios"
          class="w-full pl-10 pr-4 py-2.5 bg-elevated rounded-full text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
        />
      </div>
    </div>

    <!-- User actions -->
    <div class="flex items-center gap-3 ml-6">
      <UButton
        icon="i-solar-bell-linear"
        color="neutral"
        variant="ghost"
        size="md"
        class="rounded-full"
        aria-label="Notifications"
      />
      
      <ClientOnly>
        <!-- Logged In User Avatar Dropdown -->
        <UDropdownMenu v-if="auth.isAuthenticated.value" :items="dropdownItems" :content="{ align: 'end', side: 'bottom' }">
          <UAvatar
            :src="auth.user.value?.avatar_url || 'https://i.pravatar.cc/40'"
            alt="User avatar"
            size="sm"
            class="cursor-pointer ring-2 ring-primary/20 hover:ring-primary/50 transition-all"
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
          class="rounded-full shadow-md shadow-primary/10 cursor-pointer font-medium"
        />
      </ClientOnly>
    </div>
  </header>
</template>
