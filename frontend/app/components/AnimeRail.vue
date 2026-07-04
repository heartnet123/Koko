<script setup lang="ts">
import { ref, computed } from 'vue'
import { useIntersectionObserver } from '@vueuse/core'
import { useJikan } from '~/composables/useJikan'
import type { JikanAnime } from '~/types/anime'

const props = defineProps<{
  title: string
  count?: number
  fetchUrl: string
}>()

const railEl = ref<HTMLElement | null>(null)
const scrollContainer = ref<HTMLElement | null>(null)
const hasLoaded = ref(false)

const { data: response, pending, error, refresh } = useJikan<{ data: JikanAnime[] }>(
  () => props.fetchUrl,
  { immediate: false }
)

const items = computed(() => {
  const list = response.value?.data ?? []
  const seen = new Set<number>()
  return list.filter(item => {
    if (!item.mal_id || seen.has(item.mal_id)) return false
    seen.add(item.mal_id)
    return true
  })
})

const { stop } = useIntersectionObserver(
  railEl,
  ([{ isIntersecting }]) => {
    if (isIntersecting && !hasLoaded.value) {
      hasLoaded.value = true
      refresh()
      stop()
    }
  },
  { rootMargin: '200px' }
)

const scroll = (direction: 'left' | 'right') => {
  if (!scrollContainer.value) return
  const container = scrollContainer.value
  const scrollAmount = container.clientWidth * 0.8
  if (direction === 'left') {
    container.scrollLeft -= scrollAmount
  } else {
    container.scrollLeft += scrollAmount
  }
}
</script>

<template>
  <div ref="railEl" class="relative flex flex-col gap-4">
    <!-- Header -->
    <div v-if="!error && (!response || items.length > 0)" class="flex items-center justify-between px-6 md:px-10">
      <div class="flex items-baseline gap-3">
        <h3 class="text-base md:text-lg font-medium tracking-tight text-highlighted">
          {{ title }}
        </h3>
        <span v-if="count" class="text-xs text-toned font-mono">({{ count }})</span>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="pending && !response" class="w-full animate-pulse">
      <RailSkeleton />
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center h-48 bg-elevated/40 rounded-2xl border border-muted/50 gap-2 mx-6 md:mx-10 px-4">
      <p class="text-sm text-toned">Failed to load {{ title }}</p>
      <UButton size="xs" color="neutral" variant="subtle" label="Retry" @click="refresh" />
    </div>

    <!-- Content (Slider) -->
    <div v-else-if="items.length > 0" class="relative group/rail">
      <!-- Left Edge Fade Mask -->
      <div class="absolute left-0 top-0 bottom-0 w-12 md:w-16 bg-gradient-to-r from-[var(--ui-bg)] to-transparent pointer-events-none z-10"></div>
      
      <!-- Right Edge Fade Mask -->
      <div class="absolute right-0 top-0 bottom-0 w-12 md:w-16 bg-gradient-to-l from-[var(--ui-bg)] to-transparent pointer-events-none z-10"></div>

      <!-- Scroller Container -->
      <div
        ref="scrollContainer"
        class="flex gap-4 overflow-x-auto scroll-smooth snap-x hide-scrollbar px-6 md:px-10 pb-4"
      >
        <AnimeCard
          v-for="item in items"
          :key="item.mal_id"
          :anime="item"
          class="snap-start"
        />
      </div>

      <!-- Optional hover navigation controls -->
      <button 
        class="absolute left-4 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/60 text-white flex items-center justify-center opacity-0 group-hover/rail:opacity-100 transition-opacity duration-300 z-20 cursor-pointer backdrop-blur-sm shadow hover:bg-black/80"
        @click="scroll('left')"
        aria-label="Scroll left"
      >
        <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4" />
      </button>
      <button 
        class="absolute right-4 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/60 text-white flex items-center justify-center opacity-0 group-hover/rail:opacity-100 transition-opacity duration-300 z-20 cursor-pointer backdrop-blur-sm shadow hover:bg-black/80"
        @click="scroll('right')"
        aria-label="Scroll right"
      >
        <UIcon name="i-solar-alt-arrow-right-linear" class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
