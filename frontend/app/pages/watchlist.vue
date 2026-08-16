<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'My Watchlist — KoKo',
  description: 'Manage and track your personal watchlist of anime series and movies on KoKo.'
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
  <div class="max-w-7xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-8 animate-fade-in-up">
    <!-- Header -->
    <div class="glass-surface p-6 md:p-8 rounded-3xl border border-[var(--glass-border)] shadow-md flex items-center justify-between">
      <div>
        <div class="inline-flex items-center gap-2 mb-2">
          <div class="w-2 h-4 bg-primary-500 rounded-full " />
          <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">Personal Vault</span>
        </div>
        <h1 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          My Watchlist
        </h1>
        <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-1 font-normal">
          Keep track of your active series, scheduled movies, and bookmarked anime.
        </p>
      </div>
      <div class="glass-chip px-3 py-1.5 rounded-xl font-mono text-xs text-primary-300 font-bold border-primary-400/30">
        {{ auth.watchlist.value.length }} {{ auth.watchlist.value.length === 1 ? 'Title' : 'Titles' }}
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="auth.watchlist.value.length === 0"
      class="flex flex-col items-center justify-center text-center py-20 glass-surface border border-dashed border-[var(--glass-border)] rounded-3xl p-8 shadow-sm"
    >
      <div class="w-16 h-16 rounded-2xl glass-pill flex items-center justify-center mb-4 text-primary-400 shadow-inner">
        <UIcon name="i-solar-bookmark-linear" class="w-8 h-8" />
      </div>
      <h2 class="text-lg font-bold text-[var(--ui-text-highlighted)] mb-2">Your watchlist is empty</h2>
      <p class="text-xs md:text-sm text-[var(--ui-text-toned)] max-w-sm mb-6 font-normal">
        Explore KoKo's library, discover trending seasonal releases, and click the bookmark button to build your queue.
      </p>
      <UButton
        to="/browse"
        label="Explore Anime Library"
        color="primary"
        size="lg"
        icon="i-solar-compass-linear"
        class="rounded-2xl font-bold shadow-xl shadow-primary-500/25 px-8 py-3 cursor-pointer hover:scale-[1.02] active:scale-95 transition-all"
      />
    </div>

    <!-- Watchlist Grid -->
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-5 md:gap-6">
      <div
        v-for="item in auth.watchlist.value"
        :key="item.anime_id"
        class="group relative glass-surface border border-[var(--glass-border)] rounded-2xl overflow-hidden shadow-md hover:shadow-[0_12px_30px_rgba(0,220,130,0.25)] hover:border-primary-400/50 transition-all duration-300 flex flex-col h-full hover:-translate-y-1.5"
      >
        <!-- Card Poster Wrapper -->
        <NuxtLink
          :to="`/movie/${item.anime_id}`"
          class="relative block aspect-[2/3] overflow-hidden group cursor-pointer focus:outline-none"
        >
          <NuxtImg
            :src="item.image_url"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-108"
            placeholder
          />
          <div class="absolute inset-0 bg-[var(--ui-overlay)]/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <div class="w-10 h-10 rounded-full glass-chip flex items-center justify-center text-white shadow-xl">
              <UIcon name="i-solar-play-bold" class="w-5 h-5 text-primary-400 ml-0.5" />
            </div>
          </div>
        </NuxtLink>

        <!-- Card Content -->
        <div class="p-3.5 flex-1 flex flex-col justify-between gap-3">
          <NuxtLink
            :to="`/movie/${item.anime_id}`"
            class="text-xs font-semibold text-[var(--ui-text-highlighted)] group-hover:text-primary-400 transition-colors line-clamp-2 leading-snug tracking-tight cursor-pointer min-h-[2rem]"
          >
            {{ item.title }}
          </NuxtLink>

          <!-- Remove Action Button -->
          <button
            type="button"
            class="w-full py-2 rounded-xl flex items-center justify-center gap-1.5 text-[11px] font-semibold text-[var(--ui-error)] hover:text-[var(--ui-error)]/80 glass-pill hover:bg-[var(--ui-error)]/10 hover:border-[var(--ui-error)]/30 transition-all cursor-pointer font-mono"
            @click="handleRemove(item.anime_id)"
          >
            <UIcon name="i-solar-trash-bin-trash-linear" class="w-3.5 h-3.5" />
            <span>Remove</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>