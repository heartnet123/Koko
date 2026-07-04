<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Koko - My Watchlist',
  description: 'Manage and track your personal watchlist of anime series, episodes, and films.'
})

const auth = useAuth()

onMounted(async () => {
  await auth.fetchWatchlist()
})

const handleRemove = async (animeId: number) => {
  await auth.removeFromWatchlist(animeId)
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-6">
    <div class="flex items-baseline justify-between border-b border-muted/50 pb-4">
      <div>
        <h1 class="text-3xl font-semibold tracking-tighter text-highlighted mb-1">My Watchlist</h1>
        <p class="text-toned text-xs">Keep track of the anime you are watching or plan to watch.</p>
      </div>
      <span class="text-xs text-toned font-semibold bg-elevated border border-muted px-3 py-1 rounded-full">
        {{ auth.watchlist.value.length }} Items
      </span>
    </div>

    <!-- Empty State -->
    <div
      v-if="auth.watchlist.value.length === 0"
      class="flex flex-col items-center justify-center text-center py-20 bg-elevated/40 border border-dashed border-muted rounded-3xl p-8"
    >
      <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4 text-primary">
        <UIcon name="i-solar-bookmark-linear" class="w-8 h-8" />
      </div>
      <h2 class="text-lg font-semibold text-highlighted mb-2">Your watchlist is empty</h2>
      <p class="text-toned text-sm max-w-sm mb-6">
        Explore Koko's library, search for anime movies and series, and add them to your watchlist to track your progress.
      </p>
      <UButton
        to="/browse"
        label="Browse Anime"
        color="primary"
        icon="i-solar-compass-linear"
        class="rounded-full shadow-lg shadow-primary/10 px-6 cursor-pointer"
      />
    </div>

    <!-- Watchlist Grid -->
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-6">
      <div
        v-for="item in auth.watchlist.value"
        :key="item.anime_id"
        class="group relative bg-elevated border border-muted/50 rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-all flex flex-col h-full"
      >
        <!-- Card Poster Wrapper -->
        <NuxtLink
          :to="`/movie/${item.anime_id}`"
          class="relative block aspect-[3/4] overflow-hidden group cursor-pointer focus:outline-none"
        >
          <NuxtImg
            :src="item.image_url"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            placeholder
          />
          <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <div class="w-10 h-10 rounded-full bg-white/90 backdrop-blur flex items-center justify-center shadow-lg">
              <UIcon name="i-solar-play-bold" class="w-5 h-5 text-primary ml-0.5" />
            </div>
          </div>
        </NuxtLink>

        <!-- Card Content -->
        <div class="p-3 flex-1 flex flex-col justify-between">
          <NuxtLink
            :to="`/movie/${item.anime_id}`"
            class="text-xs font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2 leading-snug tracking-tight mb-3 cursor-pointer"
          >
            {{ item.title }}
          </NuxtLink>

          <!-- Remove Action Button -->
          <UButton
            icon="i-solar-trash-bin-trash-linear"
            label="Remove"
            variant="subtle"
            color="red"
            size="xs"
            class="w-full py-1.5 rounded-lg flex items-center justify-center text-[10px] font-semibold cursor-pointer"
            @click="handleRemove(item.anime_id)"
          />
        </div>
      </div>
    </div>
  </div>
</template>
