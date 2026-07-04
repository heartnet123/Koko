<script setup lang="ts">
import { ref } from 'vue'

const isMobileMenuOpen = ref(false)

const menuItems = [
  { icon: 'i-solar-home-smile-linear', label: 'Home', to: '/' },
  { icon: 'i-solar-compass-linear', label: 'Browse', to: '/browse' },
  { icon: 'i-solar-fire-linear', label: 'Trending', to: '/trending' },
  { icon: 'i-solar-bookmark-linear', label: 'Watchlist', to: '/watchlist' },
]

const userItems = [
  { icon: 'i-solar-user-linear', label: 'My Profile', to: '/profile' },
  { icon: 'i-solar-settings-linear', label: 'Settings', to: '/settings' },
]
</script>

<template>
  <div class="min-h-screen bg-default text-default font-sans flex lg:flex-row flex-col selection:bg-primary-100 selection:text-primary-900">
    <!-- Desktop Sidebar -->
    <AppSidebar class="hidden lg:flex" />

    <!-- Mobile Nav Bar (bottom) -->
    <nav class="lg:hidden fixed bottom-0 left-0 right-0 h-16 bg-default/95 backdrop-blur-md border-t border-muted/50 flex items-center justify-around px-4 z-40">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.label"
        :to="item.to"
        class="flex flex-col items-center justify-center gap-1 text-toned hover:text-default transition-colors w-12"
        active-class="text-primary hover:text-primary"
      >
        <UIcon :name="item.icon" class="w-5 h-5" />
        <span class="text-[10px] font-medium">{{ item.label }}</span>
      </NuxtLink>
      <button
        class="flex flex-col items-center justify-center gap-1 text-toned hover:text-default transition-colors w-12 cursor-pointer focus:outline-none"
        @click="isMobileMenuOpen = true"
        aria-label="Open menu"
      >
        <UIcon name="i-solar-hamburger-menu-linear" class="w-5 h-5" />
        <span class="text-[10px] font-medium">Menu</span>
      </button>
    </nav>

    <main class="flex-1 flex flex-col min-w-0 pb-16 lg:pb-0">
      <AppHeader @toggle-menu="isMobileMenuOpen = true" />
      <div class="flex-1 overflow-y-auto pb-12">
        <slot />
      </div>
    </main>

    <!-- Mobile Slideover Drawer -->
    <Transition name="drawer">
      <div v-if="isMobileMenuOpen" class="fixed inset-0 z-50 lg:hidden flex justify-end">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/45 backdrop-blur-xs" @click="isMobileMenuOpen = false" />
        <!-- Content -->
        <div class="relative w-80 max-w-full bg-default border-l border-muted flex flex-col pt-8 pb-6 h-full shadow-2xl">
          <div class="px-8 mb-6 flex items-center justify-between">
            <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary rounded-lg" @click="isMobileMenuOpen = false">
              <div class="w-1.5 h-6 bg-primary rounded-full" />
              <span class="text-2xl font-semibold tracking-tighter text-highlighted">KoKo</span>
            </NuxtLink>
            <button @click="isMobileMenuOpen = false" class="text-toned hover:text-default cursor-pointer focus:outline-none" aria-label="Close menu">
              <UIcon name="i-solar-close-circle-linear" class="w-6 h-6" />
            </button>
          </div>
          <nav class="flex-1 px-4 flex flex-col gap-1 overflow-y-auto">
            <NuxtLink
              v-for="item in userItems"
              :key="item.label"
              :to="item.to"
              class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium text-toned hover:text-default hover:bg-elevated transition-colors"
              active-class="bg-primary/10 text-primary hover:bg-primary/10 hover:text-primary"
              @click="isMobileMenuOpen = false"
            >
              <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0" />
              {{ item.label }}
            </NuxtLink>
          </nav>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.drawer-enter-active, .drawer-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-enter-active .relative, .drawer-leave-active .relative {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.drawer-enter-from {
  opacity: 0;
}
.drawer-enter-from .relative {
  transform: translateX(100%);
}
.drawer-leave-to {
  opacity: 0;
}
.drawer-leave-to .relative {
  transform: translateX(100%);
}
</style>

