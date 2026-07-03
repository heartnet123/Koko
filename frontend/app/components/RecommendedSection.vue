<script setup lang="ts">
import type { Anime } from '~/types/anime'

const { data: response, status } = await useFetch<{ data: Anime[] }>(
  'http://localhost:8080/api/anime?limit=10&order_by=popularity',
  {
    key: 'recommended-anime',
    retry: 1,
    retryDelay: 2000,
    retryStatusCodes: [429]
  }
)

const recommendations = computed(() => {
  const seen = new Set<number>()
  return (response.value?.data ?? []).filter((a: Anime) => {
    if (seen.has(a.mal_id)) return false
    seen.add(a.mal_id)
    return true
  })
})

const loading = computed(() => status.value === 'pending')
const scrollContainer = ref<HTMLElement | null>(null)

const scroll = (direction: 'left' | 'right') => {
  if (!scrollContainer.value) return
  const container = scrollContainer.value
  const scrollAmount = container.clientWidth
  if (direction === 'left') {
    container.scrollLeft -= scrollAmount
  } else {
    container.scrollLeft += scrollAmount
  }
}
</script>

<template>
  <section>
    <div class="flex items-center justify-between mb-5">
      <h3 class="text-lg font-medium tracking-tight text-highlighted">Recommended</h3>
      <div class="flex items-center gap-4">
        <div class="hidden md:flex items-center gap-1.5">
          <UButton
            icon="i-solar-alt-arrow-left-linear"
            color="neutral"
            variant="ghost"
            size="sm"
            class="text-toned hover:text-default cursor-pointer"
            @click="scroll('left')"
          />
          <UButton
            icon="i-solar-alt-arrow-right-linear"
            color="neutral"
            variant="ghost"
            size="sm"
            class="text-toned hover:text-default cursor-pointer"
            @click="scroll('right')"
          />
        </div>
        <UButton
          label="See All"
          to="/browse?order_by=popularity"
          trailing-icon="i-solar-alt-arrow-right-linear"
          color="neutral"
          variant="ghost"
          size="sm"
          class="text-toned hover:text-default cursor-pointer"
        />
      </div>
    </div>

    <div
      v-if="loading"
      class="flex items-center justify-center h-32 text-sm text-toned bg-elevated rounded-2xl border border-muted border-dashed"
    >
      Loading recommendations...
    </div>

    <div v-else ref="scrollContainer" class="flex gap-5 overflow-x-auto scroll-smooth snap-x snap-mandatory scrollbar-none pb-4 pl-4 relative">
      <NuxtLink
        v-for="(item, idx) in recommendations"
        :key="item.mal_id"
        :to="`/movie/${item.mal_id}`"
        class="w-[calc((100%-20px)/2+20px)] sm:w-[calc((100%-40px)/3+20px)] md:w-[calc((100%-60px)/4+20px)] lg:w-[calc((100%-80px)/5+20px)] flex-shrink-0 snap-start flex flex-col group cursor-pointer relative pl-6"
      >
        <!-- Giant rank number behind the poster -->
        <span
          class="absolute left-0 bottom-8 text-7xl md:text-8xl font-bold text-white select-none z-0 leading-none pointer-events-none transition-all duration-300 group-hover:scale-110 group-hover:-translate-y-2"
        >
          {{ idx + 1 }}
        </span>

        <!-- Poster Container -->
        <div class="ml-4 relative z-10 rounded-2xl overflow-hidden aspect-[3/4] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-muted/50 bg-elevated transition-transform duration-500 group-hover:-translate-y-1">
          <NuxtImg
            :src="item.images.jpg.large_image_url || item.images.jpg.image_url"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
          />
          <div class="absolute bottom-3 left-3 w-8 h-8 rounded-full bg-white/90 backdrop-blur shadow-sm flex items-center justify-center opacity-90 group-hover:opacity-100 group-hover:scale-110 transition-all">
            <UIcon name="i-solar-play-bold" class="w-3.5 h-3.5 ml-0.5 text-highlighted" />
          </div>
        </div>
        <h4 class="ml-4 text-sm font-medium text-highlighted tracking-tight truncate">{{ item.title }}</h4>
        <p class="ml-4 text-xs text-toned mt-0.5">{{ item.type ?? 'Anime' }}</p>
      </NuxtLink>
    </div>
  </section>
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
