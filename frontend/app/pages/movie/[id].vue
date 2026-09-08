<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from '~/composables/useAuth'
import type { JikanEpisodesResponse, JikanEpisode } from '~/types/anime'

const route = useRoute()
const auth = useAuth()

const animeId = computed(() => route.params.id as string)
const { data: response, status, error } = await useFetch<{ data: any }>(
  () => `http://localhost:8080/api/anime/${animeId.value}`,
  { key: `anime-${animeId.value}` }
)
const anime = computed(() => response.value?.data)

const episodePage = ref(1)
const { data: episodesResponse, status: episodesStatus } = await useFetch<JikanEpisodesResponse>(
  () => `http://localhost:8080/api/anime/${animeId.value}/episodes?page=${episodePage.value}`,
  {
    key: `anime-episodes-${animeId.value}-${episodePage.value}`,
    watch: [episodePage, animeId],
  }
)
const episodes = computed<JikanEpisode[]>(() => episodesResponse.value?.data || [])
const episodesPagination = computed(() => episodesResponse.value?.pagination)

const formatEpisodeDate = (dateStr?: string | null) => {
  if (!dateStr) return 'TBA'
  try {
    return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return dateStr
  }
}
useHead(() => ({
  title: anime.value
    ? `${anime.value.title_english || anime.value.title} — KoKo`
    : 'Loading… — KoKo',
  meta: [
    { name: 'description', content: anime.value?.synopsis?.slice(0, 160) || 'Anime details on KoKo' },
  ],
}))
const heroImage = computed(() => {
  const imgs = anime.value?.images
  return imgs?.webp?.large_image_url || imgs?.jpg?.large_image_url || imgs?.jpg?.image_url || ''
})
const youtubeEmbedUrl = computed(() => {
  const trailer = anime.value?.trailer
  if (!trailer) return null
  return trailer.embed_url || (trailer.youtube_id ? `https://www.youtube.com/embed/${trailer.youtube_id}` : null)
})
const formattedScore = computed(() => {
  const s = anime.value?.score
  return s != null ? Number(s).toFixed(2) : '—'
})
const infoRows = computed(() => {
  if (!anime.value) return []
  const a = anime.value
  return [
    { label: 'Type', value: a.type, icon: 'i-solar-tv-linear' },
    { label: 'Episodes', value: a.episodes ?? '—', icon: 'i-solar-play-linear' },
    { label: 'Status', value: a.status, icon: 'i-solar-clock-circle-linear' },
    { label: 'Aired', value: a.aired?.string ?? '—', icon: 'i-solar-calendar-linear' },
    { label: 'Duration', value: a.duration, icon: 'i-solar-stopwatch-linear' },
    { label: 'Rating', value: a.rating, icon: 'i-solar-shield-check-linear' },
    { label: 'Source', value: a.source, icon: 'i-solar-book-linear' },
    { label: 'Season', value: a.season ? `${a.season.charAt(0).toUpperCase() + a.season.slice(1)} ${a.year || ''}`.trim() : null, icon: 'i-solar-leaf-linear' },
  ].filter(r => r.value != null && r.value !== '')
})
const statCards = computed(() => {
  if (!anime.value) return []
  const a = anime.value
  return [
    { label: 'Ranked', value: a.rank != null ? `#${a.rank}` : '—', icon: 'i-solar-ranking-linear' },
    { label: 'Popularity', value: a.popularity != null ? `#${a.popularity}` : '—', icon: 'i-solar-fire-linear' },
    { label: 'Members', value: a.members != null ? a.members.toLocaleString() : '—', icon: 'i-solar-users-group-rounded-linear' },
    { label: 'Favorites', value: a.favorites != null ? a.favorites.toLocaleString() : '—', icon: 'i-solar-heart-linear' },
  ]
})

const isInWatchlist = computed(() => {
  return anime.value ? auth.inWatchlist(anime.value.mal_id) : false
})

const toggleWatchlist = async () => {
  if (!auth.isAuthenticated.value) {
    navigateTo('/login')
    return
  }
  if (!anime.value) return
  if (isInWatchlist.value) {
    await auth.removeFromWatchlist(anime.value.mal_id)
  } else {
    const img = anime.value.images?.jpg?.large_image_url || anime.value.images?.jpg?.image_url || ''
    await auth.addToWatchlist(anime.value.mal_id, anime.value.title_english || anime.value.title, img)
  }
}
</script>

<template>
  <!-- Loading skeleton -->
  <div v-if="status === 'pending'" class="max-w-7xl mx-auto px-4 md:px-8 w-full mt-4 animate-fade-in-up">
    <div class="h-8 w-32 glass-pill animate-glass-shimmer rounded-xl mb-6" />
    <div class="h-[440px] glass-surface-elevated rounded-[2.5rem] mb-8 animate-glass-shimmer border border-[var(--glass-border)]" />
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-2 space-y-6">
        <div class="h-64 glass-surface rounded-3xl animate-glass-shimmer border border-[var(--glass-border)]" />
        <div class="h-48 glass-surface rounded-3xl animate-glass-shimmer border border-[var(--glass-border)]" />
      </div>
      <div class="space-y-6">
        <div class="h-64 glass-surface rounded-3xl animate-glass-shimmer border border-[var(--glass-border)]" />
        <div class="h-72 glass-surface rounded-3xl animate-glass-shimmer border border-[var(--glass-border)]" />
      </div>
    </div>
  </div>

  <!-- Error state -->
  <div v-else-if="error" class="max-w-7xl mx-auto px-4 md:px-8 w-full mt-4">
    <NuxtLink to="/" class="inline-flex items-center gap-2 text-xs font-semibold text-[var(--ui-text-toned)] hover:text-primary-400 transition-colors mb-6 group glass-pill px-3 py-1.5 rounded-xl w-fit">
      <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
      Back to Home
    </NuxtLink>
    <div class="flex flex-col items-center justify-center py-20 text-center glass-surface rounded-3xl border border-[var(--glass-border)] px-6">
      <div class="w-16 h-16 rounded-2xl glass-pill flex items-center justify-center mb-4 text-[var(--ui-error)]">
        <UIcon name="i-solar-danger-triangle-linear" class="w-8 h-8" />
      </div>
      <h2 class="text-xl font-bold text-[var(--ui-text-highlighted)] mb-2">Failed to load anime details</h2>
      <p class="text-xs text-[var(--ui-text-toned)] max-w-sm mb-6 font-normal">We encountered an issue communicating with the catalog. Please try again.</p>
      <UButton
        label="Try Again"
        icon="i-solar-refresh-linear"
        color="primary"
        class="rounded-xl px-6 py-2.5 font-bold shadow-[var(--shadow-diffuse-accent)] cursor-pointer"
        @click="refreshNuxtData(`anime-${animeId}`)"
      />
    </div>
  </div>

  <!-- Main content -->
  <div v-else-if="anime" class="max-w-7xl mx-auto px-4 md:px-8 w-full mt-4 pb-12 animate-fade-in-up">
    <!-- Back button -->
    <NuxtLink to="/" class="inline-flex items-center gap-2 text-xs font-semibold text-[var(--ui-text-toned)] hover:text-primary-400 transition-colors mb-6 group glass-pill px-3.5 py-2 rounded-xl w-fit">
      <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
      Back to Library
    </NuxtLink>

    <!-- Glass Hero Banner -->
    <div class="relative w-full min-h-[460px] md:h-[460px] rounded-[2.5rem] overflow-hidden glass-surface-elevated border border-[var(--glass-border)]  mb-8 p-6 md:p-10 flex flex-col justify-end">
      <!-- Background Artwork Blur -->
      <NuxtImg
        v-if="heroImage"
        :src="heroImage"
        :alt="anime.title"
        class="absolute inset-0 w-full h-full object-cover object-top scale-110 blur-md opacity-35"
      />
      <!-- Glass Gradients -->
      <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/95 via-[var(--ui-overlay)]/70 to-transparent" />
      <div class="absolute inset-0 bg-gradient-to-r from-[var(--ui-overlay)]/90 via-[var(--ui-overlay)]/40 to-transparent hidden md:block" />

      <!-- Hero content -->
      <div class="relative flex flex-col md:flex-row items-start md:items-center gap-6 md:gap-8 z-10">
        <!-- Poster with 3D glass edge -->
        <div class="flex-shrink-0 w-44 md:w-56 aspect-[2/3] rounded-2xl overflow-hidden  glass-surface border border-[var(--glass-border)] ring-1 ring-white/20">
          <NuxtImg
            :src="anime.images?.jpg?.large_image_url || anime.images?.jpg?.image_url"
            :alt="anime.title"
            class="w-full h-full object-cover"
          />
        </div>

        <!-- Text details -->
        <div class="flex flex-col justify-center max-w-2xl text-white">
          <div class="flex items-center gap-2 mb-3 flex-wrap font-mono text-xs">
            <span
              v-if="anime.score != null"
              class="glass-chip text-[var(--rank-gold)] px-3 py-1 rounded-xl flex items-center gap-1.5 font-bold shadow-sm"
            >
              <UIcon name="i-solar-star-bold" class="w-3.5 h-3.5 text-[var(--rank-gold)]" />
              {{ formattedScore }}
            </span>
            <span v-if="anime.year != null" class="glass-chip text-[var(--ui-text-on-image)] px-2.5 py-1 rounded-xl">{{ anime.year }}</span>
            <span v-if="anime.type" class="glass-chip text-[var(--ui-text-on-image)] px-2.5 py-1 rounded-xl uppercase">{{ anime.type }}</span>
            <span v-if="anime.rating" class="glass-chip text-[var(--ui-text-on-image-muted)] px-2.5 py-1 rounded-xl">{{ anime.rating }}</span>
          </div>

          <h1 class="text-2xl md:text-4xl lg:text-5xl font-bold tracking-tight text-white leading-tight mb-2 drop-shadow-md">
            {{ anime.title_english || anime.title }}
          </h1>
          <p v-if="anime.title_japanese" class="text-xs md:text-sm text-[var(--ui-text-on-image-muted)] mb-4 font-mono">{{ anime.title_japanese }}</p>

          <!-- Genre pills -->
          <div class="flex flex-wrap gap-2 mb-5">
            <span
              v-for="genre in anime.genres"
              :key="genre.mal_id"
              class="text-xs font-semibold text-[var(--ui-text-on-image)] glass-chip px-3 py-1 rounded-xl hover:border-primary-400/40 transition-colors"
            >
              {{ genre.name }}
            </span>
          </div>

          <!-- Quick stats row -->
          <div class="flex items-center gap-5 text-xs md:text-sm text-[var(--ui-text-on-image-muted)] mb-6 font-mono">
            <span v-if="anime.episodes != null" class="flex items-center gap-1.5">
              <UIcon name="i-solar-playlist-minimalistic-linear" class="w-4 h-4 text-primary-400" />
              {{ anime.episodes }} EPS
            </span>
            <span v-if="anime.duration" class="flex items-center gap-1.5">
              <UIcon name="i-solar-stopwatch-linear" class="w-4 h-4 text-primary-400" />
              {{ anime.duration }}
            </span>
            <span v-if="anime.status" class="flex items-center gap-1.5">
              <UIcon name="i-solar-clock-circle-linear" class="w-4 h-4 text-primary-400" />
              {{ anime.status }}
            </span>
          </div>

          <!-- Action buttons -->
          <div class="flex items-center gap-3.5">
            <UButton
              v-if="youtubeEmbedUrl"
              icon="i-solar-play-bold"
              label="Watch Trailer"
              color="primary"
              size="lg"
              class="rounded-2xl font-bold shadow-[var(--shadow-diffuse-accent)] px-6 py-3 cursor-pointer hover:scale-[1.02] active:scale-95 transition-all"
              @click="() => {
                const el = document.getElementById('trailer-section')
                if (el) el.scrollIntoView({ behavior: 'smooth' })
              }"
            />
            <UButton
              :icon="isInWatchlist ? 'i-solar-bookmark-bold' : 'i-solar-bookmark-linear'"
              :label="isInWatchlist ? 'In Watchlist' : 'Add to Watchlist'"
              color="neutral"
              variant="ghost"
              size="lg"
              class="glass-chip font-bold px-6 py-3 rounded-2xl transition-all cursor-pointer hover:bg-white/20 hover:scale-[1.02] active:scale-95"
              :class="isInWatchlist ? '!text-primary-400 border-primary-400/50' : 'text-white'"
              @click="toggleWatchlist"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Two-column body -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Left column (wider) -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Trailer Section -->
        <section v-if="youtubeEmbedUrl" id="trailer-section" class="glass-surface-elevated rounded-3xl border border-[var(--glass-border)] overflow-hidden shadow-[var(--shadow-diffuse)]">
          <div class="px-6 pt-5 pb-4 flex items-center gap-2 border-b border-[var(--glass-border-subtle)]">
            <UIcon name="i-solar-play-circle-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Official Trailer</h2>
          </div>
          <div class="aspect-video bg-[var(--ui-overlay)]/90">
            <iframe
              :src="youtubeEmbedUrl"
              class="w-full h-full"
              frameborder="0"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowfullscreen
              loading="lazy"
            />
          </div>
        </section>

        <!-- Synopsis -->
        <section v-if="anime.synopsis" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 shadow-sm">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-document-text-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Synopsis</h2>
          </div>
          <p class="text-xs md:text-sm text-[var(--ui-text)] leading-relaxed whitespace-pre-line font-normal">{{ anime.synopsis }}</p>
        </section>

        <!-- Episodes Section -->
        <section id="episodes-section" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <UIcon name="i-solar-clapperboard-play-linear" class="w-5 h-5 text-primary-400" />
              <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Episodes</h2>
            </div>
            <span v-if="episodesPagination" class="text-xs font-mono text-[var(--ui-text-toned)]">
              Page {{ episodePage }} / {{ episodesPagination.last_visible_page }}
            </span>
          </div>

          <!-- Loading State -->
          <div v-if="episodesStatus === 'pending'" class="space-y-2.5">
            <div v-for="i in 5" :key="i" class="h-12 glass-pill animate-glass-shimmer rounded-2xl" />
          </div>

          <!-- Episode List -->
          <div v-else-if="episodes.length" class="space-y-2 max-h-[460px] overflow-y-auto pr-1">
            <NuxtLink
              v-for="ep in episodes"
              :key="ep.mal_id"
              :to="`/watch/${animeId}/${ep.mal_id}`"
              class="flex items-center justify-between p-3 rounded-2xl glass-pill hover:bg-white/10 hover:border-primary-400/40 transition-all group cursor-pointer"
            >
              <div class="flex items-center gap-3 min-w-0">
                <span class="glass-chip text-primary-400 font-mono font-bold text-xs px-2.5 py-1 rounded-xl flex-shrink-0">
                  EP {{ ep.mal_id }}
                </span>
                <div class="min-w-0">
                  <p class="text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] group-hover:text-primary-400 transition-colors truncate">
                    {{ ep.title }}
                  </p>
                  <p v-if="ep.title_japanese" class="text-[10px] text-[var(--ui-text-toned)] truncate font-mono">
                    {{ ep.title_japanese }}
                  </p>
                </div>
                <span v-if="ep.filler" class="text-[10px] font-mono font-bold text-[var(--ui-warn)] glass-chip px-2 py-0.5 rounded-lg flex-shrink-0">
                  Filler
                </span>
                <span v-if="ep.recap" class="text-[10px] font-mono font-bold text-[var(--ui-text-toned)] glass-chip px-2 py-0.5 rounded-lg flex-shrink-0">
                  Recap
                </span>
              </div>

              <div class="flex items-center gap-3 flex-shrink-0 ml-3">
                <span class="text-[11px] font-mono text-[var(--ui-text-toned)] whitespace-nowrap">
                  {{ formatEpisodeDate(ep.aired) }}
                </span>
                <UIcon name="i-solar-play-linear" class="w-4 h-4 text-[var(--ui-text-toned)] group-hover:text-primary-400 group-hover:translate-x-0.5 transition-all" />
              </div>
            </NuxtLink>
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-8">
            <p class="text-xs text-[var(--ui-text-toned)]">No episode data available for this anime.</p>
          </div>

          <!-- Pagination Controls -->
          <div
            v-if="episodesPagination && episodesPagination.last_visible_page > 1"
            class="flex items-center justify-between pt-4 mt-4 border-t border-[var(--glass-border-subtle)]"
          >
            <UButton
              label="Previous"
              icon="i-solar-alt-arrow-left-linear"
              size="xs"
              variant="ghost"
              color="neutral"
              class="glass-pill text-xs font-semibold"
              :disabled="episodePage <= 1"
              @click="episodePage--"
            />
            <span class="text-xs font-mono text-[var(--ui-text-toned)]">
              {{ episodePage }} of {{ episodesPagination.last_visible_page }}
            </span>
            <UButton
              label="Next"
              trailing-icon="i-solar-alt-arrow-right-linear"
              size="xs"
              variant="ghost"
              color="neutral"
              class="glass-pill text-xs font-semibold"
              :disabled="!episodesPagination.has_next_page"
              @click="episodePage++"
            />
          </div>
        </section>

        <!-- Background Information -->
        <section v-if="anime.background" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 shadow-sm">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-info-circle-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Background</h2>
          </div>
          <p class="text-xs md:text-sm text-[var(--ui-text)] leading-relaxed whitespace-pre-line font-normal">{{ anime.background }}</p>
        </section>

        <!-- Relations -->
        <section v-if="anime.relations?.length" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 shadow-sm">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-link-round-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Related Anime & Media</h2>
          </div>
          <div class="space-y-4">
            <div v-for="(rel, idx) in anime.relations" :key="idx" class="glass-pill p-4 rounded-2xl">
              <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">{{ rel.relation }}</span>
              <div class="flex flex-wrap gap-2 mt-2">
                <template v-for="entry in rel.entry" :key="entry.mal_id">
                  <NuxtLink
                    v-if="entry.type === 'anime'"
                    :to="`/movie/${entry.mal_id}`"
                    class="text-xs font-semibold text-[var(--ui-text-highlighted)] glass-chip px-3 py-1.5 rounded-xl hover:border-primary-400/50 hover:text-primary-400 transition-all"
                  >
                    {{ entry.name }}
                    <span class="text-[10px] text-[var(--ui-text-toned)] ml-1 font-mono">({{ entry.type }})</span>
                  </NuxtLink>
                  <span
                    v-else
                    class="text-xs font-medium text-[var(--ui-text-toned)] glass-pill px-3 py-1.5 rounded-xl cursor-default"
                  >
                    {{ entry.name }}
                    <span class="text-[10px] opacity-70 ml-1 font-mono">({{ entry.type }})</span>
                  </span>
                </template>
              </div>
            </div>
          </div>
        </section>

        <!-- Theme Songs -->
        <section v-if="anime.theme?.openings?.length || anime.theme?.endings?.length" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 md:p-8 shadow-sm">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-music-note-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Theme Songs</h2>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div v-if="anime.theme?.openings?.length">
              <h3 class="text-xs font-bold text-[var(--ui-text-toned)] uppercase tracking-wider mb-3 font-mono">Openings</h3>
              <ul class="space-y-2">
                <li
                  v-for="(op, idx) in anime.theme.openings"
                  :key="idx"
                  class="text-xs text-[var(--ui-text-highlighted)] flex items-start gap-2.5 p-2.5 rounded-xl glass-pill font-normal"
                >
                  <UIcon name="i-solar-play-linear" class="w-3.5 h-3.5 text-primary-400 mt-0.5 flex-shrink-0" />
                  <span>{{ op }}</span>
                </li>
              </ul>
            </div>
            <div v-if="anime.theme?.endings?.length">
              <h3 class="text-xs font-bold text-[var(--ui-text-toned)] uppercase tracking-wider mb-3 font-mono">Endings</h3>
              <ul class="space-y-2">
                <li
                  v-for="(ed, idx) in anime.theme.endings"
                  :key="idx"
                  class="text-xs text-[var(--ui-text-highlighted)] flex items-start gap-2.5 p-2.5 rounded-xl glass-pill font-normal"
                >
                  <UIcon name="i-solar-play-linear" class="w-3.5 h-3.5 text-[var(--ui-text-toned)] mt-0.5 flex-shrink-0" />
                  <span>{{ ed }}</span>
                </li>
              </ul>
            </div>
          </div>
        </section>
      </div>

      <!-- Right column (Sidebar) -->
      <div class="space-y-6">
        <!-- Score Card -->
        <div class="glass-surface-elevated rounded-3xl border border-[var(--glass-border)] p-6 text-center shadow-[var(--shadow-diffuse)]">
          <div class="w-20 h-20 mx-auto rounded-2xl glass-surface flex items-center justify-center mb-3 shadow-[var(--shadow-diffuse-lg)] ring-2 ring-primary-400/30">
            <span class="text-3xl font-bold text-primary-400 font-mono">{{ formattedScore }}</span>
          </div>
          <p class="text-xs text-[var(--ui-text-toned)] mb-5 font-mono">
            {{ anime.scored_by ? `${anime.scored_by.toLocaleString()} global votes` : 'No votes recorded' }}
          </p>
          <div class="grid grid-cols-2 gap-3">
            <div
              v-for="stat in statCards"
              :key="stat.label"
              class="glass-pill rounded-2xl p-3.5 text-center hover:border-primary-400/40 transition-colors"
            >
              <UIcon :name="stat.icon" class="w-4 h-4 text-primary-400 mb-1 mx-auto" />
              <p class="text-sm font-bold text-[var(--ui-text-highlighted)] font-mono">{{ stat.value }}</p>
              <p class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-mono font-semibold">{{ stat.label }}</p>
            </div>
          </div>
        </div>

        <!-- Information Card -->
        <div class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 shadow-sm">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-info-square-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Information</h2>
          </div>
          <dl class="space-y-3.5">
            <div
              v-for="row in infoRows"
              :key="row.label"
              class="flex items-start gap-3 text-xs"
            >
              <UIcon :name="row.icon" class="w-4 h-4 text-[var(--ui-text-toned)] mt-0.5 flex-shrink-0" />
              <div class="min-w-0 flex-1">
                <dt class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-mono font-semibold">{{ row.label }}</dt>
                <dd class="text-[var(--ui-text-highlighted)] font-bold truncate">{{ row.value }}</dd>
              </div>
            </div>
          </dl>
        </div>

        <!-- Studios -->
        <div v-if="anime.studios?.length" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 shadow-sm">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-buildings-2-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Studios</h2>
          </div>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="studio in anime.studios"
              :key="studio.mal_id"
              class="text-xs font-bold text-primary-400 glass-chip px-3 py-1.5 rounded-xl cursor-default"
            >
              {{ studio.name }}
            </span>
          </div>
        </div>

        <!-- External & Streaming Links -->
        <div v-if="anime.streaming?.length || anime.external?.length" class="glass-surface rounded-3xl border border-[var(--glass-border)] p-6 shadow-sm">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-global-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">Official Links</h2>
          </div>
          <div class="space-y-2">
            <a
              v-for="link in [...(anime.streaming || []), ...(anime.external || [])].slice(0, 8)"
              :key="link.name"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center justify-between gap-2 text-xs font-semibold text-[var(--ui-text-highlighted)] glass-pill hover:bg-white/10 px-3.5 py-2.5 rounded-xl border border-transparent hover:border-[var(--glass-border)] transition-all group"
            >
              <span class="truncate">{{ link.name }}</span>
              <UIcon name="i-solar-arrow-right-up-linear" class="w-3.5 h-3.5 text-[var(--ui-text-toned)] group-hover:text-primary-400 transition-colors" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>