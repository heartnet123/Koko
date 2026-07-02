<script setup lang="ts">
const { genreIds: selectedGenreIds } = useGenreQuery()

interface Genre {
  mal_id: number
  name: string
  count: number
}

const route = useRoute()
const isOpen = ref(false)
const genres = ref<Genre[]>([])
const loading = ref(true)

onMounted(async () => {
  window.addEventListener('click', handleClickOutside)
  try {
    let res = await fetch('http://localhost:8080/api/genres')
    
    if (res.status === 429) {
      await new Promise(r => setTimeout(r, 2000))
      res = await fetch('http://localhost:8080/api/genres')
    }
    
    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`)
    }
    
    const d = await res.json()
    genres.value = (d.data ?? [])
      .filter((g: Genre) => g.count > 1000)
      .slice(0, 20)
  } catch (e) {
    console.error('Failed to fetch genres', e)
  } finally {
    loading.value = false
  }
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
    <UButton
      icon="i-solar-filter-linear"
      trailing-icon="i-solar-alt-arrow-down-linear"
      color="neutral"
      variant="outline"
      class="rounded-xl px-4 py-2.5 min-h-[44px] text-sm font-medium border-muted/80 bg-elevated hover:bg-default text-highlighted transition-all duration-200 cursor-pointer"
      aria-haspopup="listbox"
      :aria-expanded="isOpen"
      @click="isOpen = !isOpen"
    >
      <span>{{ selectedGenreIds.length > 0 ? 'Genres' : 'Filter by Genre' }}</span>
      <UBadge
        v-if="selectedGenreIds.length > 0"
        size="sm"
        color="primary"
        variant="solid"
        class="rounded-full px-1.5 py-0.5 text-xs font-semibold"
      >
        {{ selectedGenreIds.length }}
      </UBadge>
    </UButton>

    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 scale-95 -translate-y-2"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-2"
    >
      <div
        v-show="isOpen"
        class="absolute right-0 mt-2 w-80 md:w-96 rounded-2xl p-4 shadow-lg bg-elevated border border-muted z-50"
      >
        <div class="flex items-center justify-between mb-3 pb-2 border-b border-muted">
          <span class="text-xs font-semibold text-toned uppercase tracking-wider">Genres</span>
          <button
            v-if="selectedGenreIds.length > 0"
            type="button"
            class="text-xs font-semibold text-primary hover:text-primary/80 transition-colors cursor-pointer"
            @click="clearGenres"
          >
            Clear all
          </button>
        </div>

        <div v-if="loading" class="flex flex-col items-center justify-center py-8 text-toned">
          <UIcon name="i-solar-spinner-linear" class="w-6 h-6 animate-spin mb-2" />
          <span class="text-xs">Loading genres...</span>
        </div>

        <div v-else-if="genres.length === 0" class="text-sm text-toned text-center py-4">
          No genres found.
        </div>

        <div v-else class="flex flex-wrap gap-2 max-h-60 overflow-y-auto pr-1 scrollbar-none">
          <button
            v-for="genre in genres"
            :key="genre.mal_id"
            type="button"
            :class="[
              'flex items-center gap-2 px-3 py-2 rounded-xl text-xs font-medium transition-all duration-200 select-none cursor-pointer min-h-[44px]',
              selectedGenreIds.includes(genre.mal_id)
                ? 'bg-primary/10 border border-primary/50 text-primary font-semibold'
                : 'bg-default border border-muted text-toned hover:bg-elevated hover:text-highlighted'
            ]"
            @click="toggleGenre(genre.mal_id)"
          >
            <UIcon
              v-if="selectedGenreIds.includes(genre.mal_id)"
              name="i-solar-check-circle-bold"
              class="w-4 h-4 text-primary"
            />
            <span>{{ genre.name }}</span>
            <span class="text-[10px] opacity-60">({{ genre.count }})</span>
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
