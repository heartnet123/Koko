<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { Anime } from '~/types/anime'

useSeoMeta({
  title: 'Trending Anime — KoKo',
  ogTitle: 'Trending Anime — KoKo',
  description: 'Explore top anime ranked by global community popularity on KoKo.',
  ogDescription: 'Explore top anime ranked by global community popularity on KoKo.',
})

const route = useRoute()
const limit = 25
const page = computed(() => {
  const p = Number(route.query.page)
  return Number.isInteger(p) && p > 0 ? p : 1
})

interface Pagination {
  last_visible_page: number
  has_next_page: boolean
  current_page: number
  items: {
    count: number
    total: number
    per_page: number
  }
}

const { data: response, status } = await useFetch<{ data: Anime[], pagination?: Pagination }>(
  () => `http://localhost:8080/api/anime?limit=${limit}&page=${page.value}&order_by=popularity`,
  {
    retry: 1,
    retryDelay: 2000,
    retryStatusCodes: [429]
  }
)

const animes = computed(() => response.value?.data ?? [])
const pagination = computed(() => response.value?.pagination)
const pageCount = computed(() => pagination.value?.last_visible_page ?? page.value)
const loading = computed(() => status.value === 'pending')
const isFirstPage = computed(() => page.value <= 1)
const isLastPage = computed(() => page.value >= pageCount.value)
const hasNextPage = computed(() => pagination.value?.has_next_page ?? animes.value.length === limit)

const setPage = (targetPage: number) => {
  const nextPage = Math.min(Math.max(targetPage, 1), pageCount.value)
  if (nextPage === page.value) return

  const query = { ...route.query }
  if (nextPage === 1) {
    query.page = undefined
  } else {
    query.page = String(nextPage)
  }

  navigateTo({ path: route.path, query })
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const firstPage = () => setPage(1)
const prevPage = () => setPage(page.value - 1)
const nextPage = () => {
  if (hasNextPage.value) setPage(page.value + 1)
}
const lastPage = () => setPage(pageCount.value)

const getRank = (index: number) => {
  return index + 1 + (page.value - 1) * limit
}

const getRankBadgeClass = (rank: number) => {
  if (rank === 1) return 'bg-[var(--rank-gold)] text-[var(--ui-overlay)] border border-[var(--rank-gold)]/60 font-extrabold '
  if (rank === 2) return 'bg-[var(--rank-silver)] text-[var(--ui-overlay)] border border-slate-100 font-extrabold '
  if (rank === 3) return 'bg-[var(--rank-bronze)] text-[var(--ui-text-on-image)] border border-[var(--rank-gold)]/60 font-extrabold '
  return 'glass-chip text-white border border-white/15 font-bold'
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 md:px-8 w-full flex flex-col gap-8 mt-4 mb-16 animate-fade-in-up">
    <!-- Header banner -->
    <div class="glass-surface p-6 md:p-8 rounded-3xl border border-[var(--glass-border)] shadow-md flex items-center justify-between">
      <div>
        <div class="inline-flex items-center gap-2 mb-2">
          <div class="w-2 h-4 bg-[var(--rank-gold)] rounded-full " />
          <span class="text-xs font-bold text-[var(--rank-gold)] uppercase tracking-wider font-mono">Global Charts</span>
        </div>
        <h2 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          Trending Anime
        </h2>
        <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-1 font-normal">
          Discover top titles ranked by community popularity and seasonal engagement.
        </p>
      </div>
      <div class="hidden sm:flex items-center gap-2 glass-pill px-3 py-1.5 rounded-xl font-mono text-xs text-[var(--ui-text-toned)]">
        <UIcon name="i-solar-fire-bold" class="w-4 h-4 text-[var(--rank-gold)]" />
        <span>Live Ranking</span>
      </div>
    </div>

    <div class="relative min-h-[360px]">
      <!-- Loading Glass Overlay -->
      <div
        v-if="loading"
        class="absolute inset-0 flex items-center justify-center bg-[var(--ui-overlay)]/20 backdrop-blur-xs z-10 text-[var(--ui-text-toned)] transition-all duration-200 rounded-3xl"
      >
        <div class="flex items-center gap-3 glass-surface-elevated border border-[var(--glass-border)] px-5 py-3 rounded-2xl shadow-[var(--shadow-diffuse-lg)]">
          <UIcon name="i-solar-spinner-linear" class="w-5 h-5 animate-spin text-primary-400" />
          <span class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Updating leaderboard...</span>
        </div>
      </div>

      <div v-if="animes.length === 0 && !loading" class="flex flex-col items-center justify-center h-72 text-[var(--ui-text-toned)] glass-surface border border-dashed border-[var(--glass-border)] rounded-3xl p-8 text-center">
        <UIcon name="i-solar-ghost-bold" class="w-10 h-10 mb-3 text-primary-400" />
        <h3 class="text-sm font-bold text-[var(--ui-text-highlighted)]">No trending anime found</h3>
        <p class="text-xs text-[var(--ui-text-toned)] mt-1 font-normal">Please check back in a moment.</p>
      </div>

      <div
        v-else
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5 md:gap-6 transition-opacity duration-200"
        :class="{ 'opacity-50 pointer-events-none': loading }"
      >
        <div
          v-for="(item, index) in animes"
          :key="item.mal_id"
          class="relative flex flex-col group cursor-pointer"
        >
          <!-- Rank Badge -->
          <div
            :class="[
              'absolute top-2.5 left-2.5 px-2.5 py-1 text-[11px] rounded-xl shadow-[var(--shadow-diffuse)] z-20 select-none tracking-tight flex items-center justify-center font-mono',
              getRankBadgeClass(getRank(index))
            ]"
          >
            #{{ getRank(index) }}
          </div>

          <NuxtLink
            :to="`/movie/${item.mal_id}`"
            class="flex flex-col w-full transition-transform duration-300 hover:-translate-y-1.5"
          >
            <div class="relative aspect-[2/3] w-full rounded-2xl overflow-hidden glass-surface border border-[var(--glass-border)] shadow-md group-hover:shadow-[0_12px_30px_rgba(0,220,130,0.25)] group-hover:border-primary-400/50 transition-all duration-300">
              <NuxtImg
                :src="item.images.jpg.large_image_url || item.images.jpg.image_url"
                :alt="item.title"
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-108"
                loading="lazy"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none" />
              
              <!-- Score Chip -->
              <div 
                v-if="item.score" 
                class="absolute top-2.5 right-2.5 glass-chip text-white font-mono text-[10px] px-2 py-0.5 rounded-lg flex items-center gap-1 z-10"
              >
                <UIcon name="i-solar-star-bold" class="w-3 h-3 text-yellow-400" />
                <span>{{ Number(item.score).toFixed(1) }}</span>
              </div>
            </div>

            <div class="mt-2.5 px-1">
              <h4 class="text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] tracking-tight line-clamp-2 min-h-[2.4rem] leading-snug group-hover:text-primary-500 transition-colors duration-200">
                {{ item.title }}
              </h4>
              <p class="text-[11px] text-[var(--ui-text-toned)] mt-1 flex items-center justify-between font-mono">
                <span class="truncate max-w-[100px]">{{ item.type ?? 'Anime' }}</span>
                <span v-if="item.year">{{ item.year }}</span>
              </p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Glass Pagination Controls -->
    <div v-if="animes.length > 0" class="flex flex-col sm:flex-row items-center justify-center gap-3 mt-4 glass-surface p-4 rounded-2xl border border-[var(--glass-border)] w-fit mx-auto shadow-sm">
      <div class="flex items-center gap-1.5">
        <button
          :disabled="isFirstPage"
          class="glass-pill px-3 py-1.5 rounded-xl text-xs font-semibold text-[var(--ui-text-highlighted)] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-white/10 transition-all flex items-center gap-1 cursor-pointer"
          @click="firstPage"
        >
          <UIcon name="i-solar-double-alt-arrow-left-linear" class="w-3.5 h-3.5" />
          <span>First</span>
        </button>
        <button
          :disabled="isFirstPage"
          class="glass-pill px-3 py-1.5 rounded-xl text-xs font-semibold text-[var(--ui-text-highlighted)] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-white/10 transition-all flex items-center gap-1 cursor-pointer"
          @click="prevPage"
        >
          <UIcon name="i-solar-alt-arrow-left-linear" class="w-3.5 h-3.5" />
          <span>Prev</span>
        </button>
      </div>

      <div class="text-xs text-[var(--ui-text-toned)] px-3 font-mono font-semibold">
        Page <span class="text-primary-400 font-bold">{{ page }}</span> of
        <span class="text-[var(--ui-text-highlighted)] font-bold">{{ pageCount }}</span>
      </div>

      <div class="flex items-center gap-1.5">
        <button
          :disabled="!hasNextPage"
          class="glass-pill px-3 py-1.5 rounded-xl text-xs font-semibold text-[var(--ui-text-highlighted)] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-white/10 transition-all flex items-center gap-1 cursor-pointer"
          @click="nextPage"
        >
          <span>Next</span>
          <UIcon name="i-solar-alt-arrow-right-linear" class="w-3.5 h-3.5" />
        </button>
        <button
          :disabled="isLastPage"
          class="glass-pill px-3 py-1.5 rounded-xl text-xs font-semibold text-[var(--ui-text-highlighted)] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-white/10 transition-all flex items-center gap-1 cursor-pointer"
          @click="lastPage"
        >
          <span>Last</span>
          <UIcon name="i-solar-double-alt-arrow-right-linear" class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>