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
    <div v-if="!error && (!response || items.length > 0)" class="flex items-center justify-between px-4 md:px-8">
      <div class="flex items-center gap-3">
        <div class="w-1.5 h-4 bg-primary-500 rounded-full " />
        <h3 class="text-base md:text-lg font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          {{ title }}
        </h3>
        <span v-if="count" class="glass-pill px-2 py-0.5 text-[11px] text-[var(--ui-text-toned)] font-mono rounded-md">
          {{ count.toLocaleString() }}
        </span>
      </div>

      <!-- Arrow Controls in header for desktop -->
      <div v-if="items.length > 0" class="hidden sm:flex items-center gap-1.5">
        <button 
          class="glass-chip w-7 h-7 rounded-lg text-[var(--ui-text-highlighted)] hover:text-primary-400 flex items-center justify-center transition-all cursor-pointer active:scale-95"
          @click="scroll('left')"
          aria-label="Scroll left"
        >
          <UIcon name="i-solar-alt-arrow-left-linear" class="w-3.5 h-3.5" />
        </button>
        <button 
          class="glass-chip w-7 h-7 rounded-lg text-[var(--ui-text-highlighted)] hover:text-primary-400 flex items-center justify-center transition-all cursor-pointer active:scale-95"
          @click="scroll('right')"
          aria-label="Scroll right"
        >
          <UIcon name="i-solar-alt-arrow-right-linear" class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="pending && !response" class="w-full">
      <RailSkeleton />
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center h-44 glass-surface rounded-2xl border border-[var(--glass-border)] gap-3 mx-4 md:px-8 px-4 text-center">
      <UIcon name="i-solar-danger-circle-linear" class="w-6 h-6 text-[var(--ui-error)]" />
      <p class="text-xs text-[var(--ui-text-toned)]">Failed to load {{ title }}</p>
      <UButton size="xs" color="neutral" variant="ghost" label="Retry" class="glass-pill cursor-pointer px-4 font-semibold" @click="refresh" />
    </div>

    <!-- Content (Slider) -->
    <div v-else-if="items.length > 0" class="relative group/rail">
      <!-- Left Edge Fade Mask -->
      <div class="absolute left-0 top-0 bottom-0 w-8 md:w-14 bg-gradient-to-r from-[var(--ui-bg)] to-transparent pointer-events-none z-10" />
      
      <!-- Right Edge Fade Mask -->
      <div class="absolute right-0 top-0 bottom-0 w-8 md:w-14 bg-gradient-to-l from-[var(--ui-bg)] to-transparent pointer-events-none z-10" />

      <!-- Scroller Container -->
      <div
        ref="scrollContainer"
        class="flex gap-4 overflow-x-auto scroll-smooth snap-x hide-scrollbar px-4 md:px-8 pb-3 pt-1"
      >
        <AnimeCard
          v-for="item in items"
          :key="item.mal_id"
          :anime="item"
          class="snap-start"
        />
      </div>

      <!-- Floating Hover Navigation Controls -->
      <button 
        class="absolute left-3 top-1/2 -translate-y-1/2 w-9 h-9 rounded-full glass-chip text-white flex items-center justify-center opacity-0 group-hover/rail:opacity-100 transition-all duration-300 z-20 cursor-pointer shadow-xl hover:scale-110 active:scale-95 hidden sm:flex"
        @click="scroll('left')"
        aria-label="Scroll left"
      >
        <UIcon name="i-solar-alt-arrow-left-linear" class="w-4 h-4" />
      </button>
      <button 
        class="absolute right-3 top-1/2 -translate-y-1/2 w-9 h-9 rounded-full glass-chip text-white flex items-center justify-center opacity-0 group-hover/rail:opacity-100 transition-all duration-300 z-20 cursor-pointer shadow-xl hover:scale-110 active:scale-95 hidden sm:flex"
        @click="scroll('right')"
        aria-label="Scroll right"
      >
        <UIcon name="i-solar-alt-arrow-right-linear" class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>