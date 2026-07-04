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
  <aside class="hidden lg:flex w-64 flex-shrink-0 border-r border-muted flex flex-col pt-8 pb-6 bg-default h-screen sticky top-0">
    <!-- Brand -->
    <NuxtLink to="/" class="px-8 mb-10 flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary rounded-lg">
      <div class="w-1.5 h-6 bg-primary rounded-full" />
      <span class="text-2xl font-semibold tracking-tighter text-highlighted">KoKo</span>
    </NuxtLink>

    <!-- Main nav -->
    <nav class="flex-1 px-4 flex flex-col gap-1">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.label"
        :to="item.to"
        class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium text-toned hover:text-default hover:bg-elevated transition-colors"
        active-class="bg-primary/10 text-primary hover:bg-primary/10 hover:text-primary"
      >
        <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0" />
        {{ item.label }}
      </NuxtLink>

      <div class="my-4 border-t border-muted" />

      <NuxtLink
        v-for="item in userItems"
        :key="item.label"
        :to="item.to"
        class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium text-toned hover:text-default hover:bg-elevated transition-colors"
        active-class="bg-primary/10 text-primary hover:bg-primary/10 hover:text-primary"
      >
        <UIcon :name="item.icon" class="w-5 h-5 flex-shrink-0" />
        {{ item.label }}
      </NuxtLink>
    </nav>
  </aside>
</template>
