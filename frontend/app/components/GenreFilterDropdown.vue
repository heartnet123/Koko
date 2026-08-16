<script setup lang="ts">
const { genreIds: selectedGenreIds } = useGenreQuery()

interface Genre {
  mal_id: number
  name: string
  count: number
}

const route = useRoute()
const isOpen = ref(false)
const { data: genresResponse, status: genresFetchStatus } = await useFetch<{ data: Genre[] }>(
  'http://localhost:8080/api/genres',
  {
    key: 'genres-list',
    retry: 1,
    retryDelay: 2000,
    retryStatusCodes: [429]
  }
)

const genres = computed(() => {
  const list = genresResponse.value?.data ?? []
  return list
    .filter((g: Genre) => g.count > 1000)
    .slice(0, 20)
})

const loading = computed(() => genresFetchStatus.value === 'pending')

onMounted(() => {
  window.addEventListener('click', handleClickOutside)
})

const router = useRouter()

const toggleGenre = (genreId: number) => {
  const currentIds = [...selectedGenreIds.value]
  const idx = currentIds.indexOf(genreId)
  if (idx > -1) {
    currentIds.splice(idx, 1)
  } else {
    currentIds.push(genreId)
  }

  const query = { ...route.query }
  if (currentIds.length > 0) {
    query.genres = currentIds.join(',')
  } else {
    delete query.genres
  }
  
  delete query.genre
  delete query.genre_name
  
  router.push({ query })
}

const clearGenres = () => {
  const query = { ...route.query }
  delete query.genres
  delete query.genre
  delete query.genre_name
  router.push({ query })
}

const dropdownRef = ref<HTMLElement | null>(null)

function handleClickOutside(event: MouseEvent) {
  const target = event.target as Node
  if (isOpen.value && dropdownRef.value && !dropdownRef.value.contains(target) && document.body.contains(target)) {
    isOpen.value = false
  }
}

onBeforeUnmount(() => {
  window.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="relative" ref="dropdownRef">
    <button
      type="button"
      class="glass-surface hover:glass-surface-elevated rounded-xl px-4 py-2.5 min-h-[44px] text-sm font-semibold text-[var(--ui-text-highlighted)] transition-all duration-200 cursor-pointer flex items-center gap-2 border border-[var(--glass-border)] shadow-sm hover:border-primary-500/40 "
      aria-haspopup="listbox"
      :aria-expanded="isOpen"
      @click="isOpen = !isOpen"
    >
      <UIcon name="i-solar-filter-linear" class="w-4 h-4 text-primary-400" />
      <span>{{ selectedGenreIds.length > 0 ? 'Genres' : 'Filter by Genre' }}</span>
      <span
        v-if="selectedGenreIds.length > 0"
        class="bg-primary-500 text-white rounded-full px-2 py-0.5 text-[10px] font-mono font-bold "
      >
        {{ selectedGenreIds.length }}
      </span>
      <UIcon name="i-solar-alt-arrow-down-linear" class="w-3.5 h-3.5 text-[var(--ui-text-toned)] ml-1 transition-transform" :class="{ 'rotate-180': isOpen }" />
    </button>

    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 scale-95 -translate-y-2"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-2"
    >
      <div
        v-show="isOpen"
        class="absolute right-0 mt-2 w-80 md:w-96 rounded-2xl p-4 glass-surface-elevated border border-[var(--glass-border)] z-50 "
      >
        <div class="flex items-center justify-between mb-3 pb-2 border-b border-[var(--glass-border-subtle)]">
          <span class="text-xs font-bold text-[var(--ui-text-highlighted)] uppercase tracking-wider">Select Genres</span>
          <button
            v-if="selectedGenreIds.length > 0"
            type="button"
            class="text-xs font-bold text-primary-400 hover:text-primary-300 transition-colors cursor-pointer"
            @click="clearGenres"
          >
            Clear all
          </button>
        </div>

        <div v-if="loading" class="flex flex-col items-center justify-center py-8 text-[var(--ui-text-toned)]">
          <UIcon name="i-solar-spinner-linear" class="w-6 h-6 animate-spin mb-2 text-primary-400" />
          <span class="text-xs font-semibold">Loading genres...</span>
        </div>

        <div v-else-if="genres.length === 0" class="text-sm text-[var(--ui-text-toned)] text-center py-4">
          No genres found.
        </div>

        <div v-else class="flex flex-wrap gap-2 max-h-64 overflow-y-auto pr-1 scrollbar-none">
          <button
            v-for="genre in genres"
            :key="genre.mal_id"
            type="button"
            :class="[
              'flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-semibold transition-all duration-200 select-none cursor-pointer',
              selectedGenreIds.includes(genre.mal_id)
                ? 'bg-primary-500/20 border border-primary-400 text-primary-300 font-bold '
                : 'glass-pill text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] hover:bg-white/10 hover:border-white/20'
            ]"
            @click="toggleGenre(genre.mal_id)"
          >
            <UIcon
              v-if="selectedGenreIds.includes(genre.mal_id)"
              name="i-solar-check-circle-bold"
              class="w-3.5 h-3.5 text-primary-400"
            />
            <span>{{ genre.name }}</span>
            <span class="text-[10px] opacity-70 font-mono">({{ genre.count }})</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
.scrollbar-none {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>