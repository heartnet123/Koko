<script setup lang="ts">
import { computed } from 'vue'
import { useAuth } from '~/composables/useAuth'

const route = useRoute()
const auth = useAuth()

const animeId = computed(() => route.params.id)
const { data: response, status, error } = await useFetch<{ data: any }>(
  () => `http://localhost:8080/api/anime/${animeId.value}`,
  { key: `anime-${animeId.value}` }
)
const anime = computed(() => response.value?.data)
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
  <div v-if="status === 'pending'" class="max-w-7xl mx-auto px-8 w-full mt-4 animate-pulse">
    <div class="h-8 w-24 bg-elevated rounded-lg mb-6" />
    <div class="h-[420px] bg-elevated rounded-[2rem] mb-8" />
    <div class="grid grid-cols-3 gap-8">
      <div class="col-span-2 space-y-6">
        <div class="h-48 bg-elevated rounded-2xl" />
        <div class="h-32 bg-elevated rounded-2xl" />
      </div>
      <div class="space-y-6">
        <div class="h-56 bg-elevated rounded-2xl" />
        <div class="h-64 bg-elevated rounded-2xl" />
      </div>
    </div>
  </div>
  <!-- Error state -->
  <div v-else-if="error" class="max-w-7xl mx-auto px-8 w-full mt-4">
    <NuxtLink to="/" class="inline-flex items-center gap-2 text-sm text-toned hover:text-primary transition-colors mb-6 group">
      <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
      Back to Home
    </NuxtLink>
    <div class="flex flex-col items-center justify-center py-24 text-center">
      <div class="w-16 h-16 rounded-full bg-red-500/10 flex items-center justify-center mb-4">
        <UIcon name="i-solar-danger-triangle-linear" class="w-7 h-7 text-red-400" />
      </div>
      <h2 class="text-xl font-semibold text-highlighted mb-2">Something went wrong</h2>
      <p class="text-sm text-toned max-w-sm mb-6">We couldn't load this anime. Please check your connection and try again.</p>
      <UButton
        label="Try Again"
        icon="i-solar-refresh-linear"
        color="primary"
        class="rounded-full"
        @click="refreshNuxtData(`anime-${animeId}`)"
      />
    </div>
  </div>
  <!-- Main content -->
  <div v-else-if="anime" class="max-w-7xl mx-auto px-8 w-full mt-4 pb-8">
    <!-- Back button -->
    <NuxtLink to="/" class="inline-flex items-center gap-2 text-sm text-toned hover:text-primary transition-colors mb-6 group">
      <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
      Back to Home
    </NuxtLink>
    <!-- Hero -->
    <div class="relative w-full h-[420px] rounded-[2rem] overflow-hidden shadow-lg mb-8 bg-elevated">
      <NuxtImg
        v-if="heroImage"
        :src="heroImage"
        :alt="anime.title"
        class="absolute inset-0 w-full h-full object-cover object-top scale-105 blur-sm opacity-30"
      />
      <!-- Gradient overlays -->
      <div class="absolute inset-0 bg-gradient-to-r from-default via-default/90 to-default/30" />
      <div class="absolute inset-0 bg-gradient-to-t from-default/70 via-transparent to-transparent" />
      <!-- Hero content -->
      <div class="relative h-full flex items-center gap-8 px-10">
        <!-- Poster -->
        <div class="flex-shrink-0 w-56 h-[320px] rounded-2xl overflow-hidden shadow-[0_12px_40px_-8px_rgba(0,0,0,0.4)] border border-white/10">
          <NuxtImg
            :src="anime.images?.jpg?.large_image_url || anime.images?.jpg?.image_url"
            :alt="anime.title"
            class="w-full h-full object-cover"
          />
        </div>
        <!-- Text overlay -->
        <div class="flex flex-col justify-center max-w-xl">
          <div class="flex items-center gap-3 mb-3 flex-wrap">
            <span
              v-if="anime.score != null"
              class="inline-flex items-center gap-1 bg-primary-500/15 text-primary-400 text-xs font-semibold px-2.5 py-1 rounded-lg"
            >
              <UIcon name="i-solar-star-bold" class="w-3.5 h-3.5" />
              {{ formattedScore }}
            </span>
            <span v-if="anime.year != null" class="text-xs text-toned font-medium bg-elevated px-2.5 py-1 rounded-lg">{{ anime.year }}</span>
            <span v-if="anime.type" class="text-xs text-toned font-medium bg-elevated px-2.5 py-1 rounded-lg">{{ anime.type }}</span>
            <span v-if="anime.rating" class="text-xs text-toned font-medium bg-elevated px-2.5 py-1 rounded-lg">{{ anime.rating }}</span>
          </div>
          <h1 class="text-3xl font-bold tracking-tight text-white leading-tight mb-1">
            {{ anime.title_english || anime.title }}
          </h1>
          <p v-if="anime.title_japanese" class="text-sm text-white/50 mb-4">{{ anime.title_japanese }}</p>
          <!-- Genre pills -->
          <div class="flex flex-wrap gap-2 mb-5">
            <span
              v-for="genre in anime.genres"
              :key="genre.mal_id"
              class="text-xs font-medium text-white/80 bg-white/10 backdrop-blur border border-white/10 px-3 py-1 rounded-full"
            >
              {{ genre.name }}
            </span>
          </div>
          <!-- Quick stats row -->
          <div class="flex items-center gap-5 text-sm text-white/70 mb-6">
            <span v-if="anime.episodes != null" class="flex items-center gap-1.5">
              <UIcon name="i-solar-playlist-minimalistic-linear" class="w-4 h-4 text-white/50" />
              {{ anime.episodes }} eps
            </span>
            <span v-if="anime.duration" class="flex items-center gap-1.5">
              <UIcon name="i-solar-stopwatch-linear" class="w-4 h-4 text-white/50" />
              {{ anime.duration }}
            </span>
            <span v-if="anime.status" class="flex items-center gap-1.5">
              <UIcon name="i-solar-clock-circle-linear" class="w-4 h-4 text-white/50" />
              {{ anime.status }}
            </span>
          </div>
          <!-- Action buttons -->
          <div class="flex items-center gap-3">
            <UButton
              icon="i-solar-play-bold"
              label="Play Now"
              color="primary"
              class="rounded-full shadow-lg shadow-primary-500/20 px-6 py-2.5 font-medium cursor-pointer"
            />
            <UButton
              :icon="isInWatchlist ? 'i-solar-minus-circle-linear' : 'i-solar-add-circle-linear'"
              :label="isInWatchlist ? 'Remove from Watchlist' : 'Add to Watchlist'"
              color="neutral"
              variant="outline"
              class="rounded-full bg-white/10 backdrop-blur border-white/20 text-white hover:bg-white/20 px-6 py-2.5 font-medium cursor-pointer"
              @click="toggleWatchlist"
            />
          </div>
        </div>
      </div>
    </div>
    <!-- Two-column body -->
    <div class="grid grid-cols-3 gap-8">
      <!-- Left column (wider) -->
      <div class="col-span-2 space-y-6">
        <!-- Trailer -->
        <section v-if="youtubeEmbedUrl" class="bg-elevated rounded-2xl border border-muted/50 overflow-hidden">
          <div class="px-6 pt-5 pb-4 flex items-center gap-2">
            <UIcon name="i-solar-play-circle-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Trailer</h2>
          </div>
          <div class="aspect-video bg-black">
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
        <section v-if="anime.synopsis" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-document-text-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Synopsis</h2>
          </div>
          <p class="text-sm text-default leading-relaxed whitespace-pre-line">{{ anime.synopsis }}</p>
        </section>
        <!-- Background -->
        <section v-if="anime.background" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-info-circle-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Background</h2>
          </div>
          <p class="text-sm text-default leading-relaxed whitespace-pre-line">{{ anime.background }}</p>
        </section>
        <!-- Relations -->
        <section v-if="anime.relations?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-link-round-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Relations</h2>
          </div>
          <div class="space-y-3">
            <div v-for="(rel, idx) in anime.relations" :key="idx">
              <span class="text-xs font-semibold text-primary-400 uppercase tracking-wider">{{ rel.relation }}</span>
              <div class="flex flex-wrap gap-2 mt-1.5">
                <template v-for="entry in rel.entry" :key="entry.mal_id">
                  <NuxtLink
                    v-if="entry.type === 'anime'"
                    :to="`/movie/${entry.mal_id}`"
                    class="text-sm text-default bg-elevated px-3 py-1.5 rounded-xl border border-muted hover:border-primary-500/30 hover:text-primary-400 transition-colors"
                  >
                    {{ entry.name }}
                    <span class="text-xs text-toned ml-1">({{ entry.type }})</span>
                  </NuxtLink>
                  <span
                    v-else
                    class="text-sm text-toned bg-elevated px-3 py-1.5 rounded-xl border border-muted cursor-default"
                  >
                    {{ entry.name }}
                    <span class="text-xs text-toned/70 ml-1">({{ entry.type }})</span>
                  </span>
                </template>
              </div>
            </div>
          </div>
        </section>
        <!-- Theme Songs -->
        <section v-if="anime.theme?.openings?.length || anime.theme?.endings?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-music-note-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Theme Songs</h2>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div v-if="anime.theme?.openings?.length">
              <h3 class="text-xs font-semibold text-toned uppercase tracking-wider mb-2">Openings</h3>
              <ul class="space-y-1.5">
                <li
                  v-for="(op, idx) in anime.theme.openings"
                  :key="idx"
                  class="text-sm text-default flex items-start gap-2 py-1.5 px-3 rounded-lg hover:bg-white/5 transition-colors"
                >
                  <UIcon name="i-solar-play-linear" class="w-3.5 h-3.5 text-primary-400 mt-0.5 flex-shrink-0" />
                  {{ op }}
                </li>
              </ul>
            </div>
            <div v-if="anime.theme?.endings?.length">
              <h3 class="text-xs font-semibold text-toned uppercase tracking-wider mb-2">Endings</h3>
              <ul class="space-y-1.5">
                <li
                  v-for="(ed, idx) in anime.theme.endings"
                  :key="idx"
                  class="text-sm text-default flex items-start gap-2 py-1.5 px-3 rounded-lg hover:bg-white/5 transition-colors"
                >
                  <UIcon name="i-solar-play-linear" class="w-3.5 h-3.5 text-toned mt-0.5 flex-shrink-0" />
                  {{ ed }}
                </li>
              </ul>
            </div>
          </div>
        </section>
      </div>
      <!-- Right column (sidebar) -->
      <div class="space-y-6">
        <!-- Score card -->
        <div class="bg-elevated rounded-2xl border border-muted/50 p-6 text-center">
          <div class="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-primary-500 to-primary-400 flex items-center justify-center mb-3 shadow-lg shadow-primary-500/20">
            <span class="text-2xl font-bold text-white">{{ formattedScore }}</span>
          </div>
          <p class="text-xs text-toned mb-4">
            {{ anime.scored_by ? `${anime.scored_by.toLocaleString()} votes` : 'No votes yet' }}
          </p>
          <div class="grid grid-cols-2 gap-3">
            <div
              v-for="stat in statCards"
              :key="stat.label"
              class="bg-default/50 rounded-xl p-3 text-center hover:bg-primary-500/10 transition-colors border border-muted/20"
            >
              <UIcon :name="stat.icon" class="w-4.5 h-4.5 text-primary-400 mb-1" />
              <p class="text-sm font-semibold text-highlighted">{{ stat.value }}</p>
              <p class="text-[10px] text-toned uppercase tracking-wider">{{ stat.label }}</p>
            </div>
          </div>
        </div>
        <!-- Info card -->
        <div class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-4">
            <UIcon name="i-solar-info-square-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Information</h2>
          </div>
          <dl class="space-y-3">
            <div
              v-for="row in infoRows"
              :key="row.label"
              class="flex items-start gap-3 text-sm"
            >
              <UIcon :name="row.icon" class="w-4 h-4 text-toned mt-0.5 flex-shrink-0" />
              <div>
                <dt class="text-toned text-xs uppercase tracking-wider">{{ row.label }}</dt>
                <dd class="text-highlighted font-medium">{{ row.value }}</dd>
              </div>
            </div>
          </dl>
        </div>
        <!-- Studios -->
        <div v-if="anime.studios?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-buildings-2-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Studios</h2>
          </div>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="studio in anime.studios"
              :key="studio.mal_id"
              class="text-sm font-medium text-primary-400 bg-primary-500/10 px-3 py-1.5 rounded-xl hover:bg-primary-500/20 transition-colors cursor-default"
            >
              {{ studio.name }}
            </span>
          </div>
        </div>
        <!-- Producers -->
        <div v-if="anime.producers?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-people-nearby-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Producers</h2>
          </div>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="prod in anime.producers"
              :key="prod.mal_id"
              class="text-xs text-toned bg-default/50 px-3 py-1.5 rounded-xl border border-muted"
            >
              {{ prod.name }}
            </span>
          </div>
        </div>
        <!-- Tags (Genres + Themes + Demographics) -->
        <div
          v-if="anime.genres?.length || anime.themes?.length || anime.demographics?.length"
          class="bg-elevated rounded-2xl border border-muted/50 p-6"
        >
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-tag-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Tags</h2>
          </div>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="tag in [...(anime.genres || []), ...(anime.themes || []), ...(anime.demographics || [])]"
              :key="tag.mal_id"
              class="text-xs font-medium text-toned bg-default/50 px-3 py-1.5 rounded-full border border-muted hover:border-primary-500/30 hover:text-primary-400 transition-colors cursor-default"
            >
              {{ tag.name }}
            </span>
          </div>
        </div>
        <!-- Streaming -->
        <div v-if="anime.streaming?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-play-stream-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">Streaming</h2>
          </div>
          <div class="space-y-2">
            <a
              v-for="link in anime.streaming"
              :key="link.name"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center justify-between gap-2 text-sm text-highlighted bg-elevated hover:bg-primary-500/10 px-4 py-2.5 rounded-xl border border-muted hover:border-primary-500/30 transition-all group"
            >
              <span class="font-medium">{{ link.name }}</span>
              <UIcon name="i-solar-arrow-right-up-linear" class="w-4 h-4 text-toned group-hover:text-primary-400 transition-colors" />
            </a>
          </div>
        </div>
        <!-- External links -->
        <div v-if="anime.external?.length" class="bg-elevated rounded-2xl border border-muted/50 p-6">
          <div class="flex items-center gap-2 mb-3">
            <UIcon name="i-solar-global-linear" class="w-5 h-5 text-primary-400" />
            <h2 class="text-base font-semibold text-highlighted tracking-tight">External Links</h2>
          </div>
          <div class="space-y-2">
            <a
              v-for="link in anime.external"
              :key="link.name"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center justify-between gap-2 text-sm text-highlighted bg-elevated hover:bg-primary-500/10 px-4 py-2.5 rounded-xl border border-muted hover:border-primary-500/30 transition-all group"
            >
              <span class="font-medium">{{ link.name }}</span>
              <UIcon name="i-solar-arrow-right-up-linear" class="w-4 h-4 text-toned group-hover:text-primary-400 transition-colors" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>