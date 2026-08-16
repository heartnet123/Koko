<script setup lang="ts">
import { computed } from 'vue'
import type { JikanAnime } from '~/types/anime'
import { useAuth } from '~/composables/useAuth'

const props = defineProps<{
  anime: JikanAnime
}>()

const auth = useAuth()

const isInWatchlist = computed(() => {
  return auth.inWatchlist(props.anime.mal_id)
})

const handleToggleWatchlist = async (e: Event) => {
  e.preventDefault()
  e.stopPropagation()

  if (!auth.isAuthenticated.value) {
    navigateTo('/login')
    return
  }

  if (isInWatchlist.value) {
    await auth.removeFromWatchlist(props.anime.mal_id)
  } else {
    const imgUrl = props.anime.images?.jpg?.large_image_url || props.anime.images?.jpg?.image_url || ''
    await auth.addToWatchlist(props.anime.mal_id, props.anime.title_english || props.anime.title, imgUrl)
  }
}
</script>

<template>
  <NuxtLink
    :to="`/movie/${anime.mal_id}`"
    class="flex-shrink-0 snap-start flex flex-col group cursor-pointer w-[160px] md:w-[185px] transition-transform duration-300 hover:-translate-y-1.5"
  >
    <!-- Poster Container with Glass Border & Spotlight Glow -->
    <div class="relative aspect-[2/3] w-full rounded-2xl overflow-hidden glass-surface border border-[var(--glass-border)] shadow-[var(--shadow-diffuse)] group-hover:shadow-[var(--shadow-diffuse-accent)] group-hover:border-primary-400/50 transition-all duration-300">
      <NuxtImg
        v-if="anime.images?.jpg?.large_image_url || anime.images?.jpg?.image_url"
        :src="anime.images.jpg.large_image_url || anime.images.jpg.image_url"
        :alt="anime.title"
        class="w-full h-full object-cover transition-transform duration-500 ease-out group-hover:scale-108"
        loading="lazy"
      />
      
      <!-- Subtle bottom gradient vignette inside poster -->
      <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none" />

      <!-- Frosted Glass Score Chip (Top-Right) -->
      <div 
        v-if="anime.score" 
        class="absolute top-2.5 right-2.5 glass-chip text-white font-mono text-[10px] px-2 py-0.5 rounded-lg flex items-center gap-1 z-10 shadow-sm"
      >
        <UIcon name="i-solar-star-bold" class="w-3 h-3 text-[var(--rank-gold)]" />
        <span>{{ Number(anime.score).toFixed(1) }}</span>
      </div>

      <!-- Quick Watchlist Button (Top-Left on hover or if in watchlist) -->
      <button
        type="button"
        class="absolute top-2.5 left-2.5 glass-chip p-1.5 rounded-lg z-10 transition-all duration-200 cursor-pointer"
        :class="isInWatchlist ? 'opacity-100 text-primary-400 border-primary-400/60' : 'opacity-0 group-hover:opacity-100 text-white hover:text-primary-300 hover:scale-110'"
        :title="isInWatchlist ? 'Remove from Watchlist' : 'Add to Watchlist'"
        @click="handleToggleWatchlist"
      >
        <UIcon :name="isInWatchlist ? 'i-solar-bookmark-bold' : 'i-solar-bookmark-linear'" class="w-3.5 h-3.5" />
      </button>

      <!-- Episode count badge at bottom left on hover -->
      <div 
        v-if="anime.episodes" 
        class="absolute bottom-2.5 left-2.5 glass-chip text-[10px] text-[var(--ui-text-on-image)] px-2 py-0.5 rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-300 z-10 font-mono"
      >
        {{ anime.episodes }} EPS
      </div>
    </div>
    
    <!-- Info -->
    <div class="mt-2.5 px-1">
      <h4 class="text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] tracking-tight line-clamp-2 min-h-[2.4rem] leading-snug group-hover:text-primary-500 transition-colors duration-200">
        {{ anime.title_english || anime.title }}
      </h4>
      <p class="text-[11px] text-[var(--ui-text-toned)] mt-1 flex items-center justify-between font-mono">
        <span class="truncate max-w-[100px]">{{ anime.type ?? 'Anime' }}</span>
        <span v-if="anime.year">{{ anime.year }}</span>
      </p>
    </div>
  </NuxtLink>
</template>