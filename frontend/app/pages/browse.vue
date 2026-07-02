<script setup lang="ts">
import type { Anime } from '~/types/anime'

const route = useRoute()
const genreId = computed(() => route.query.genre as string)
const genreName = computed(() => route.query.genre_name as string)
const orderBy = computed(() => (route.query.order_by as string) || 'popularity')

const { data: response, status } = await useFetch<{ data: Anime[] }>(
  () => {
    let url = 'http://localhost:8080/api/anime?limit=24'
    if (genreId.value) {
      url += `&genres=${genreId.value}`
    }
    if (orderBy.value) {
      url += `&order_by=${orderBy.value}`
    }
    return url
  },
  {
    key: `browse-${route.fullPath}`
  }
)
const animes = computed(() => response.value?.data ?? [])
const loading = computed(() => status.value === 'pending')
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4 mb-12">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-highlighted">
        {{ genreName ? `Explore ${genreName}` : 'Popular Anime' }}
      </h2>
      <p class="text-sm text-toned mt-1">
        Showing collections matching your query.
      </p>
    </div>

    <div
      v-if="loading"
      class="flex items-center justify-center h-64 text-toned"
    >
      <UIcon name="i-solar-spinner-linear" class="w-6 h-6 animate-spin mr-2" />
      Loading...
    </div>

    <div v-else-if="animes.length === 0" class="flex flex-col items-center justify-center h-64 text-toned border border-dashed border-muted rounded-2xl">
      <UIcon name="i-solar-ghost-bold" class="w-8 h-8 mb-2" />
      No anime found.
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5">
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
</template>
