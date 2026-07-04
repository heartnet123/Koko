<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { Anime } from '~/types/anime'

useSeoMeta({
  title: 'Trending Anime - Koko',
  ogTitle: 'Trending Anime - Koko',
  description: 'Explore the most popular and trending anime on Koko.',
  ogDescription: 'Explore the most popular and trending anime on Koko.',
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
    delete query.page
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
  if (rank === 1) return 'bg-amber-500 text-black border border-amber-300 font-extrabold shadow-[0_2px_8px_rgba(245,158,11,0.3)]'
  if (rank === 2) return 'bg-slate-300 text-black border border-slate-100 font-extrabold shadow-[0_2px_8px_rgba(203,213,225,0.3)]'
  if (rank === 3) return 'bg-amber-700 text-white border border-amber-500 font-extrabold shadow-[0_2px_8px_rgba(180,83,9,0.3)]'
  return 'bg-black/75 backdrop-blur-md text-white border border-white/10 font-bold'
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4 mb-12">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-highlighted">
        Trending Anime
      </h2>
      <p class="text-sm text-toned mt-1">
        Discover the most popular anime series and movies ranked by global popularity.
      </p>
    </div>

    <div class="relative min-h-[300px]">
      <!-- Loading overlay -->
      <div
        v-if="loading"
        class="absolute inset-0 flex items-center justify-center bg-default/40 backdrop-blur-[1px] z-10 text-toned transition-all duration-200"
      >
        <div class="flex items-center gap-2 bg-elevated border border-muted px-4 py-2.5 rounded-full shadow-sm">
          <UIcon name="i-solar-spinner-linear" class="w-5 h-5 animate-spin text-primary" />
          <span class="text-sm font-medium">Updating ranking...</span>
        </div>
      </div>

      <div v-if="animes.length === 0 && !loading" class="flex flex-col items-center justify-center h-64 text-toned border border-dashed border-muted rounded-2xl">
        <UIcon name="i-solar-ghost-bold" class="w-8 h-8 mb-2" />
        No trending anime found.
      </div>

      <div
        v-else
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5 transition-opacity duration-200"
        :class="{ 'opacity-50 pointer-events-none': loading }"
      >
        <NuxtLink
          v-for="(item, index) in animes"
          :key="item.mal_id"
          :to="`/movie/${item.mal_id}`"
          class="flex flex-col group cursor-pointer relative"
        >
          <!-- Rank Badge -->
          <div
            :class="[
              'absolute top-3 left-3 px-2.5 py-1 text-xs rounded-lg shadow-md z-10 select-none tracking-tight flex items-center justify-center min-w-[28px]',
              getRankBadgeClass(getRank(index))
            ]"
          >
            #{{ getRank(index) }}
          </div>

          <div class="relative rounded-2xl overflow-hidden aspect-[3/4] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-muted/50 bg-elevated">
            <NuxtImg
              :src="item.images.jpg.large_image_url || item.images.jpg.image_url"
              :alt="item.title"
              class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
            />
          </div>
          <h4 class="text-sm font-medium text-highlighted tracking-tight truncate">{{ item.title }}</h4>
          <p class="text-xs text-toned mt-0.5">{{ item.type ?? 'Anime' }}</p>
        </NuxtLink>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="animes.length > 0" class="flex flex-col sm:flex-row items-center justify-center gap-3 mt-8">
      <div class="flex items-center gap-2">
        <UButton
          :disabled="isFirstPage"
          variant="ghost"
          icon="i-solar-double-alt-arrow-left-linear"
          @click="firstPage"
        >
          First
        </UButton>
        <UButton
          :disabled="isFirstPage"
          variant="ghost"
          icon="i-solar-alt-arrow-left-linear"
          @click="prevPage"
        >
          Previous
        </UButton>
      </div>

      <p class="text-sm text-toned px-2">
        Page <span class="font-semibold text-highlighted">{{ page }}</span> of
        <span class="font-semibold text-highlighted">{{ pageCount }}</span>
      </p>

      <div class="flex items-center gap-2">
        <UButton
          :disabled="!hasNextPage"
          variant="ghost"
          trailing-icon="i-solar-alt-arrow-right-linear"
          @click="nextPage"
        >
          Next
        </UButton>
        <UButton
          :disabled="isLastPage"
          variant="ghost"
          trailing-icon="i-solar-double-alt-arrow-right-linear"
          @click="lastPage"
        >
          Last
        </UButton>
      </div>
    </div>
  </div>
</template>
