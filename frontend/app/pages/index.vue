<script setup lang="ts">
import { computed } from 'vue'
import { useJikan } from '~/composables/useJikan'
import type { JikanAnime } from '~/types/anime'

useSeoMeta({
  title: 'Koko - Anime Library',
  ogTitle: 'Koko - Anime Library',
  description: 'Explore the ultimate collection of anime movies and series on Koko.',
  ogDescription: 'Explore the ultimate collection of anime movies and series on Koko.',
})

// Fetch current season for Hero (max 5 slides)
const { data: heroResponse, pending: heroPending } = useJikan<{ data: JikanAnime[] }>('/seasons/now')

const heroSlides = computed(() => {
  return (heroResponse.value?.data ?? []).slice(0, 5)
})

// Fetch genres for rails configuration
const { data: genresResponse, error: genresError } = useJikan<{ data: any[] }>('/genres')

const genres = computed(() => {
  return (genresResponse.value?.data ?? [])
    .filter(g => g.count > 1000)
    .slice(0, 10)
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 md:px-8 w-full flex flex-col gap-10 mt-4 pb-16">
    <!-- Hero Carousel -->
    <div>
      <HeroSkeleton v-if="heroPending" />
      <HeroCarousel v-else-if="heroSlides.length" :slides="heroSlides" />
      <div v-else class="min-h-[40dvh] flex flex-col items-center justify-center bg-elevated/40 rounded-3xl border border-muted/50 gap-2">
        <p class="text-toned text-sm">Couldn't load current season.</p>
      </div>
    </div>

    <!-- Recommended Rail -->
    <AnimeRail 
      title="Recommended" 
      fetchUrl="/anime?order_by=score&sort=desc&limit=12" 
    />

    <!-- Genre Rails with CSS stagger delay -->
    <div v-if="genresError" class="text-center py-6 text-sm text-toned/80 border border-dashed border-muted/50 rounded-2xl">
      Genre rails unavailable
    </div>
    
    <div v-else class="flex flex-col gap-10">
      <div 
        v-for="(genre, idx) in genres" 
        :key="genre.mal_id"
        class="staggered-rail"
        :style="{ '--i': idx }"
      >
        <AnimeRail
          :title="genre.name"
          :count="genre.count"
          :fetchUrl="`/anime?genres=${genre.mal_id}&order_by=score&sort=desc&limit=12`"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.staggered-rail {
  animation: fadeIn 0.6s ease-out both;
  animation-delay: calc(var(--i) * 100ms);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .staggered-rail {
    animation: none !important;
  }
}
</style>
