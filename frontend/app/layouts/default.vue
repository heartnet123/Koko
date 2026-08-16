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
  <div class="relative min-h-[100dvh] bg-[var(--ui-bg)] text-[var(--ui-text)] font-sans flex lg:flex-row flex-col selection:bg-primary-500/20 selection:text-primary-400">
    <!-- Ambient glowing mesh background -->
    <div class="ambient-mesh" aria-hidden="true">
      <div class="ambient-mesh-blob-1" />
      <div class="ambient-mesh-blob-2" />
      <div class="ambient-mesh-blob-3" />
    </div>

    <!-- Desktop Sidebar -->
    <AppSidebar class="hidden lg:flex z-30" />

    <!-- Mobile Floating Glass Navbar (Bottom) -->
    <div class="lg:hidden fixed bottom-3 left-0 right-0 px-4 z-40">
      <nav class="max-w-md mx-auto glass-surface rounded-2xl px-3 py-2 flex items-center justify-around ">
        <NuxtLink
          v-for="item in menuItems"
          :key="item.label"
          :to="item.to"
          class="flex flex-col items-center justify-center gap-1 text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-all duration-200 py-1 px-2 rounded-xl group"
          active-class="!text-primary-500 bg-primary-500/10 shadow-inner"
        >
          <UIcon :name="item.icon" class="w-5 h-5 transition-transform group-hover:scale-110 group-active:scale-95" />
          <span class="text-[10px] font-semibold tracking-tight">{{ item.label }}</span>
        </NuxtLink>
        <button
          class="flex flex-col items-center justify-center gap-1 text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-all duration-200 py-1 px-2 rounded-xl cursor-pointer focus:outline-none group"
          @click="isMobileMenuOpen = true"
          aria-label="Open menu"
        >
          <UIcon name="i-solar-hamburger-menu-linear" class="w-5 h-5 transition-transform group-hover:scale-110 group-active:scale-95" />
          <span class="text-[10px] font-semibold tracking-tight">Menu</span>
        </button>
      </nav>
    </div>

    <!-- Main Content Area -->
    <main class="relative flex-1 flex flex-col min-w-0 pb-20 lg:pb-0 z-10">
      <AppHeader @toggle-menu="isMobileMenuOpen = true" />
      <div class="flex-1 overflow-y-auto pb-12">
        <slot />
      </div>
    </main>

    <!-- Mobile Slideover Drawer (Glassmorphic) -->
    <Transition name="drawer">
      <div v-if="isMobileMenuOpen" class="fixed inset-0 z-50 lg:hidden flex justify-end">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-[var(--ui-overlay)]/60 backdrop-blur-sm" @click="isMobileMenuOpen = false" />
        <!-- Content -->
        <div class="relative w-80 max-w-full glass-surface-elevated border-l border-[var(--glass-border)] flex flex-col pt-8 pb-6 h-full ">
          <div class="px-6 mb-6 flex items-center justify-between">
            <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 rounded-lg" @click="isMobileMenuOpen = false">
              <div class="w-2 h-6 bg-gradient-to-b from-primary-400 to-primary-600 rounded-full " />
              <span class="text-2xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">KoKo</span>
            </NuxtLink>
            <button @click="isMobileMenuOpen = false" class="text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] cursor-pointer focus:outline-none p-1 rounded-lg hover:bg-white/10 transition-colors" aria-label="Close menu">
              <UIcon name="i-solar-close-circle-linear" class="w-6 h-6" />
            </button>
          </div>
          <nav class="flex-1 px-4 flex flex-col gap-1.5 overflow-y-auto">
            <NuxtLink
              v-for="item in menuItems"
              :key="item.label"
              :to="item.to"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] hover:bg-white/5 transition-all"
              active-class="glass-pill !text-primary-400 font-bold shadow-sm"
              @click="isMobileMenuOpen = false"
            >
              <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0" />
              {{ item.label }}
            </NuxtLink>

            <div class="my-3 border-t border-[var(--ui-border-muted)]" />

            <NuxtLink
              v-for="item in userItems"
              :key="item.label"
              :to="item.to"
              class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] hover:bg-white/5 transition-all"
              active-class="glass-pill !text-primary-400 font-bold shadow-sm"
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