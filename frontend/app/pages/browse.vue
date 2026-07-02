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
    delete query.genres
  }
  
  delete query.genre
  delete query.genre_name
  
  navigateTo({
    path: route.path,
    query
  })
}

const { data: response, status } = await useFetch<{ data: Anime[] }>(
  () => {
    let url = 'http://localhost:8080/api/anime?limit=24'
    if (genreIds.value.length > 0) {
      url += `&genres=${genreIds.value.join(',')}`
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
const loading = computed(() => status.value === 'pending')

const headerText = computed(() => {
  const count = genreIds.value.length
  if (count === 0) return 'Popular Anime'
  if (count === 1 && genreName.value) return `Explore ${genreName.value}`
  return 'Explore Anime'
})

const subtitleText = computed(() => {
  const count = genreIds.value.length
  if (count === 0) return 'Showing collections matching your query.'
  if (count === 1) return 'Showing collections matching 1 active filter.'
  return `Showing collections matching ${count} active filters.`
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4 mb-12">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-highlighted">
          {{ headerText }}
        </h2>
        <p class="text-sm text-toned mt-1">
          {{ subtitleText }}
        </p>
      </div>
      <GenreFilterDropdown />
    </div>

    <!-- Selected Genre Pills Row -->
    <TransitionGroup
      v-if="selectedGenres.length > 0"
      name="list"
      tag="div"
      class="flex flex-wrap gap-2 items-center -mt-4 relative"
    >
      <UBadge
        v-for="genre in selectedGenres"
        :key="genre.mal_id"
        size="md"
        variant="subtle"
        class="rounded-full flex items-center gap-1.5 pl-3 pr-1 py-1 text-primary bg-primary/10 border border-primary/50 font-medium"
      >
        <span>{{ genre.name }}</span>
        <UButton
          icon="i-solar-close-circle-bold"
          variant="ghost"
          color="primary"
          class="w-11 h-11 -my-3 -mr-1 rounded-full flex items-center justify-center cursor-pointer text-primary/70 hover:text-primary hover:bg-transparent p-0"
          @click="removeGenre(genre.mal_id)"
          :aria-label="`Remove ${genre.name} filter`"
        />
      </UBadge>
    </TransitionGroup>

    <div class="relative min-h-[300px]">
      <!-- Loading overlay -->
      <div
        v-if="loading"
        class="absolute inset-0 flex items-center justify-center bg-default/40 backdrop-blur-[1px] z-10 text-toned transition-all duration-200"
      >
        <div class="flex items-center gap-2 bg-elevated border border-muted px-4 py-2.5 rounded-full shadow-sm">
          <UIcon name="i-solar-spinner-linear" class="w-5 h-5 animate-spin text-primary" />
          <span class="text-sm font-medium">Updating results...</span>
        </div>
      </div>

      <div v-if="animes.length === 0 && !loading" class="flex flex-col items-center justify-center h-64 text-toned border border-dashed border-muted rounded-2xl">
        <UIcon name="i-solar-ghost-bold" class="w-8 h-8 mb-2" />
        No anime found.
      </div>

      <div
        v-else
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5 transition-opacity duration-200"
        :class="{ 'opacity-50 pointer-events-none': loading }"
      >
        <NuxtLink
          v-for="item in animes"
          :key="item.mal_id"
          :to="`/movie/${item.mal_id}`"
          class="flex flex-col group cursor-pointer"
        >
          <div class="relative rounded-2xl overflow-hidden aspect-[3/4] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-muted/50 bg-elevated">
            <NuxtImg
              :src="item.images.jpg.large_image_url || item.images.jpg.image_url"
              :alt="item.title"
              class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
            />
            <div class="absolute bottom-3 left-3 w-8 h-8 rounded-full bg-white/90 backdrop-blur shadow-sm flex items-center justify-center opacity-90 group-hover:opacity-100 group-hover:scale-110 transition-all">
              <UIcon name="i-solar-play-bold" class="w-3.5 h-3.5 ml-0.5 text-highlighted" />
            </div>
          </div>
          <h4 class="text-sm font-medium text-highlighted tracking-tight truncate">{{ item.title }}</h4>
          <p class="text-xs text-toned mt-0.5">{{ item.type ?? 'Anime' }}</p>
        </NuxtLink>
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
