<script setup lang="ts">
import type { Anime } from '~/types/anime'

interface Genre {
  mal_id: number
  name: string
  count: number
}

const { genreIds } = useGenreQuery()

const route = useRoute()
const genreName = computed(() => route.query.genre_name as string)
const orderBy = computed(() => (route.query.order_by as string) || 'popularity')
const searchQuery = computed(() => route.query.q as string)

const { data: genresResponse } = await useFetch<{ data: Genre[] }>(
  'http://localhost:8080/api/genres',
  {
    key: 'genres-list',
    retry: 1,
    retryDelay: 2000,
    retryStatusCodes: [429]
  }
)
const genresList = computed(() => genresResponse.value?.data ?? [])

const selectedGenres = computed(() => {
  return genreIds.value.map(id => {
    const found = genresList.value.find(g => g.mal_id === id)
    let name = found?.name
    if (!name && id === Number(route.query.genre) && route.query.genre_name) {
      name = route.query.genre_name as string
    }
    return {
      mal_id: id,
      name: name || `Genre ${id}`
    }
  })
})

const removeGenre = (genreId: number) => {
  const currentIds = genreIds.value.filter(id => id !== genreId)
  const query = { ...route.query }
  
  if (currentIds.length > 0) {
    query.genres = currentIds.join(',')
  } else {
    query.genres = undefined
  }
  
  query.genre = undefined
  query.genre_name = undefined
  query.page = undefined
  
  navigateTo({
    path: route.path,
    query
  })
}

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
  () => {
    let url = `http://localhost:8080/api/anime?limit=${limit}&page=${page.value}`
    if (genreIds.value.length > 0) {
      url += `&genres=${genreIds.value.join(',')}`
    }
    if (searchQuery.value) {
      url += `&q=${encodeURIComponent(searchQuery.value)}`
    }
    if (orderBy.value) {
      url += `&order_by=${orderBy.value}`
    }
    return url
  },
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

const headerText = computed(() => {
  if (searchQuery.value) return `Search Results for "${searchQuery.value}"`
  const count = genreIds.value.length
  if (count === 0) return 'Explore Library'
  if (count === 1 && genreName.value) return `Explore ${genreName.value}`
  return 'Filtered Collections'
})

const subtitleText = computed(() => {
  if (searchQuery.value) return 'Matching titles from our high-performance anime catalog.'
  const count = genreIds.value.length
  if (count === 0) return 'Browse top anime ranked by popularity and community scores.'
  if (count === 1) return 'Displaying collections matching 1 active filter tag.'
  return `Displaying collections matching ${count} active filter tags.`
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 md:px-8 w-full flex flex-col gap-8 mt-4 mb-16 animate-fade-in-up">
    <!-- Header with Filter Controls -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 glass-surface p-6 md:p-8 rounded-3xl border border-[var(--glass-border)] shadow-md">
      <div>
        <div class="inline-flex items-center gap-2 mb-2">
          <div class="w-2 h-4 bg-primary-500 rounded-full " />
          <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">Anime Catalog</span>
        </div>
        <h2 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          {{ headerText }}
        </h2>
        <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-1 font-normal">
          {{ subtitleText }}
        </p>
      </div>
      <GenreFilterDropdown />
    </div>

    <!-- Active Genre Filter Tags -->
    <TransitionGroup
      v-if="selectedGenres.length > 0"
      name="list"
      tag="div"
      class="flex flex-wrap gap-2 items-center -mt-2 relative"
    >
      <div
        v-for="genre in selectedGenres"
        :key="genre.mal_id"
        class="glass-chip rounded-xl flex items-center gap-2 pl-3 pr-2 py-1.5 text-xs font-semibold text-primary-300 border-primary-400/40 shadow-sm"
      >
        <span>{{ genre.name }}</span>
        <button
          type="button"
          class="text-[var(--ui-text-toned)] hover:text-white transition-colors cursor-pointer p-0.5 rounded hover:bg-white/10"
          @click="removeGenre(genre.mal_id)"
          :aria-label="`Remove ${genre.name} filter`"
        >
          <UIcon name="i-solar-close-circle-linear" class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>

    <!-- Grid Results Area -->
    <div class="relative min-h-[360px]">
      <!-- Loading Glass Overlay -->
      <div
        v-if="loading"
        class="absolute inset-0 flex items-center justify-center bg-[var(--ui-overlay)]/20 backdrop-blur-xs z-10 text-[var(--ui-text-toned)] transition-all duration-200 rounded-3xl"
      >
        <div class="flex items-center gap-3 glass-surface-elevated border border-[var(--glass-border)] px-5 py-3 rounded-2xl shadow-xl">
          <UIcon name="i-solar-spinner-linear" class="w-5 h-5 animate-spin text-primary-400" />
          <span class="text-xs font-bold text-[var(--ui-text-highlighted)] font-mono">Fetching catalog...</span>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="animes.length === 0 && !loading" class="flex flex-col items-center justify-center h-72 text-[var(--ui-text-toned)] glass-surface border border-dashed border-[var(--glass-border)] rounded-3xl p-8 text-center">
        <UIcon name="i-solar-ghost-bold" class="w-10 h-10 mb-3 text-primary-400" />
        <h3 class="text-sm font-bold text-[var(--ui-text-highlighted)]">No anime found</h3>
        <p class="text-xs text-[var(--ui-text-toned)] mt-1 max-w-xs font-normal">Try adjusting your search terms or clearing some genre filters.</p>
      </div>

      <!-- Anime Card Grid -->
      <div
        v-else
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5 md:gap-6 transition-opacity duration-200"
        :class="{ 'opacity-50 pointer-events-none': loading }"
      >
        <AnimeCard
          v-for="item in animes"
          :key="item.mal_id"
          :anime="item"
          class="w-full"
        />
      </div>
    </div>

    <!-- Glass Pagination Controls -->
    <div class="flex flex-col sm:flex-row items-center justify-center gap-3 mt-4 glass-surface p-4 rounded-2xl border border-[var(--glass-border)] w-fit mx-auto shadow-sm">
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

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.2s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
.list-move {
  transition: transform 0.2s ease;
}
.list-leave-active {
  position: absolute;
}
</style>