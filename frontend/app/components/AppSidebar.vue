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

    <!-- Continue Watching widget -->
    <div class="px-6 mt-auto">
      <div class="bg-elevated rounded-2xl p-4 shadow-[0_2px_12px_-4px_rgba(0,0,0,0.05)] border border-muted">
        <h3 class="text-sm font-medium text-highlighted mb-3 tracking-tight">Continue Watching</h3>
        <NuxtLink
          to="/movie/5114"
          class="relative block rounded-xl overflow-hidden mb-3 aspect-video group cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary"
          aria-label="Play Episode 4 of Sunset Romance"
        >
          <NuxtImg
            src="https://images.unsplash.com/photo-1518173946687-a4c8892bbd9f?q=80&w=400&auto=format&fit=crop"
            alt="Episode thumbnail"
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
          <div class="absolute inset-0 bg-black/20 flex items-center justify-center opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-opacity duration-300">
            <div class="w-8 h-8 rounded-full bg-white/90 backdrop-blur flex items-center justify-center shadow-lg">
              <UIcon name="i-solar-play-bold" class="w-4 h-4 text-primary ml-0.5" />
            </div>
          </div>
        </NuxtLink>
        <div class="flex justify-between items-baseline mb-1">
          <p class="text-sm font-medium text-highlighted tracking-tight">Episode 4</p>
          <span class="text-xs text-toned font-medium">60%</span>
        </div>
        <p class="text-xs text-toned mb-2">Drama • Romance</p>
        <div class="w-full bg-elevated h-1.5 rounded-full overflow-hidden">
          <div class="bg-primary w-[60%] h-full rounded-full" />
        </div>
      </div>
    </div>
  </aside>
</template>
