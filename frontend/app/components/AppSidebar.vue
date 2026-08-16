<script setup lang="ts">
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
  <aside class="hidden lg:flex w-64 flex-shrink-0 glass-surface border-r border-[var(--glass-border-subtle)] flex-col pt-8 pb-6 min-h-[100dvh] sticky top-0 transition-all z-30">
    <!-- Brand -->
    <NuxtLink to="/" class="px-7 mb-10 flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 rounded-xl group">
      <div class="w-2 h-6 bg-gradient-to-b from-primary-400 to-primary-600 rounded-full  transition-transform duration-300 group-hover:scale-y-110" />
      <span class="text-2xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">KoKo</span>
    </NuxtLink>

    <!-- Main Navigation -->
    <nav class="flex-1 px-4 flex flex-col gap-1.5">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.label"
        :to="item.to"
        class="relative flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] hover:bg-white/5 transition-all duration-200 group"
        active-class="glass-pill !text-primary-400 font-bold shadow-sm"
      >
        <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0 transition-transform group-hover:scale-110" />
        <span>{{ item.label }}</span>
      </NuxtLink>

      <div class="my-4 border-t border-[var(--glass-border-subtle)]" />

      <NuxtLink
        v-for="item in userItems"
        :key="item.label"
        :to="item.to"
        class="relative flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] hover:bg-white/5 transition-all duration-200 group"
        active-class="glass-pill !text-primary-400 font-bold shadow-sm"
      >
        <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0 transition-transform group-hover:scale-110" />
        <span>{{ item.label }}</span>
      </NuxtLink>
    </nav>
  </aside>
</template>