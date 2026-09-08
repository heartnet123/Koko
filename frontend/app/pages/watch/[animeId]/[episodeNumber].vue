<script setup lang="ts">
import { computed } from 'vue'
import type { JikanEpisodesResponse, JikanEpisode } from '~/types/anime'

const route = useRoute()

const animeId = computed(() => route.params.animeId as string)
const episodeNumber = computed(() => Number(route.params.episodeNumber))

// Anime details
const { data: animeResponse, status: animeStatus } = await useFetch<{ data: any }>(
  () => `http://localhost:8080/api/anime/${animeId.value}`,
  { key: `anime-${animeId.value}` }
)
const anime = computed(() => animeResponse.value?.data)

// Episodes to lookup title/metadata
const pageForEpisode = computed(() => Math.max(1, Math.ceil(episodeNumber.value / 100)))
const { data: episodesResponse } = await useFetch<JikanEpisodesResponse>(
  () => `http://localhost:8080/api/anime/${animeId.value}/episodes?page=${pageForEpisode.value}`,
  { key: `anime-episodes-${animeId.value}-${pageForEpisode.value}` }
)
const currentEpisode = computed<JikanEpisode | undefined>(() => {
  return episodesResponse.value?.data?.find(ep => ep.mal_id === episodeNumber.value)
})

const animeTitle = computed(() => anime.value?.title_english || anime.value?.title || `Anime #${animeId.value}`)
const episodeTitle = computed(() => currentEpisode.value?.title || `Episode ${episodeNumber.value}`)
const totalEpisodes = computed<number | null>(() => anime.value?.episodes ?? null)

useHead(() => ({
  title: `Watch ${animeTitle.value} Ep ${episodeNumber.value} — KoKo`,
  meta: [
    { name: 'description', content: `Watch ${animeTitle.value} Episode ${episodeNumber.value} on KoKo.` },
  ],
}))
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 md:px-8 w-full mt-4 pb-12 animate-fade-in-up">
    <!-- Back to Details -->
    <NuxtLink
      :to="`/movie/${animeId}`"
      class="inline-flex items-center gap-2 text-xs font-semibold text-[var(--ui-text-toned)] hover:text-primary-400 transition-colors mb-6 group glass-pill px-3.5 py-2 rounded-xl w-fit cursor-pointer"
    >
      <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
      Back to {{ animeTitle }}
    </NuxtLink>

    <!-- Player Container -->
    <div class="glass-surface-elevated rounded-3xl border border-[var(--glass-border)] overflow-hidden shadow-[var(--shadow-diffuse)] mb-6">
      <div class="aspect-video bg-[var(--ui-overlay)]/95 relative flex flex-col items-center justify-center p-6 text-center">
        <!-- Player Placeholder Surface -->
        <div class="w-20 h-20 rounded-full glass-surface flex items-center justify-center mb-4 ring-2 ring-primary-400/40 text-primary-400 shadow-[var(--shadow-diffuse-accent)] cursor-pointer hover:scale-105 active:scale-95 transition-all">
          <UIcon name="i-solar-play-bold" class="w-8 h-8 translate-x-0.5" />
        </div>
        <p class="text-xs font-mono text-[var(--ui-text-toned)] uppercase tracking-wider mb-1">
          Stream Source
        </p>
        <h2 class="text-sm md:text-base font-bold text-[var(--ui-text-highlighted)]">
          Episode {{ episodeNumber }}: {{ episodeTitle }}
        </h2>
      </div>
    </div>

    <!-- Episode Info & Navigation Controls -->
    <div class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
      <div class="min-w-0">
        <div class="flex items-center gap-2 mb-2">
          <span class="glass-chip text-primary-400 font-mono font-bold text-xs px-2.5 py-1 rounded-xl">
            EP {{ episodeNumber }}
          </span>
          <span v-if="anime?.type" class="glass-chip text-[var(--ui-text-toned)] font-mono text-xs px-2.5 py-1 rounded-xl uppercase">
            {{ anime.type }}
          </span>
        </div>
        <h1 class="text-xl md:text-2xl font-bold text-[var(--ui-text-highlighted)] tracking-tight truncate">
          {{ episodeTitle }}
        </h1>
        <p class="text-xs text-[var(--ui-text-toned)] mt-1">
          {{ animeTitle }}
        </p>
      </div>

      <!-- Episode Switcher -->
      <div class="flex items-center gap-3 w-full md:w-auto justify-end">
        <UButton
          v-if="episodeNumber > 1"
          :to="`/watch/${animeId}/${episodeNumber - 1}`"
          icon="i-solar-alt-arrow-left-linear"
          label="Prev Ep"
          variant="ghost"
          color="neutral"
          class="glass-pill text-xs font-semibold px-4 py-2 cursor-pointer"
        />
        <UButton
          v-if="!totalEpisodes || episodeNumber < totalEpisodes"
          :to="`/watch/${animeId}/${episodeNumber + 1}`"
          trailing-icon="i-solar-alt-arrow-right-linear"
          label="Next Ep"
          color="primary"
          class="rounded-xl text-xs font-semibold px-5 py-2 cursor-pointer shadow-[var(--shadow-diffuse-accent)]"
        />
      </div>
    </div>
  </div>
</template>
