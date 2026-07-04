<script setup lang="ts">
import type { JikanAnime } from '~/types/anime'

defineProps<{
  anime: JikanAnime
}>()
</script>

<template>
  <NuxtLink
    :to="`/movie/${anime.mal_id}`"
    class="flex-shrink-0 snap-start flex flex-col group cursor-pointer w-[160px] md:w-[180px]"
  >
    <!-- Poster Container -->
    <div class="relative aspect-[2/3] w-full rounded-xl overflow-hidden bg-elevated ring-1 ring-neutral-200/40 dark:ring-neutral-800/40 transition-all duration-300 group-hover:ring-primary-500/40 group-hover:shadow-md">
      <NuxtImg
        v-if="anime.images?.jpg?.large_image_url || anime.images?.jpg?.image_url"
        :src="anime.images.jpg.large_image_url || anime.images.jpg.image_url"
        :alt="anime.title"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        loading="lazy"
      />
      <!-- Score Chip -->
      <div 
        v-if="anime.score" 
        class="absolute top-2 right-2 bg-black/60 backdrop-blur-md text-white font-mono text-[10px] px-1.5 py-0.5 rounded flex items-center gap-1 z-10"
      >
        <UIcon name="i-solar-star-bold" class="w-3 h-3 text-yellow-400" />
        {{ anime.score.toFixed(1) }}
      </div>
    </div>
    
    <!-- Info -->
    <div class="mt-2 px-1">
      <h4 class="text-sm font-medium text-highlighted tracking-tight line-clamp-2 min-h-[2.5rem] leading-snug group-hover:text-primary-500 transition-colors duration-200">
        {{ anime.title }}
      </h4>
      <p class="text-xs text-toned mt-1 flex items-center justify-between">
        <span>{{ anime.type ?? 'Anime' }}</span>
        <span v-if="anime.year" class="font-mono">{{ anime.year }}</span>
      </p>
    </div>
  </NuxtLink>
</template>
