<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import { useAuth } from '~/composables/useAuth'
import type { JikanAnime } from '~/types/anime'

const props = defineProps<{
  slides: JikanAnime[]
}>()

const currentSlide = ref(0)
const auth = useAuth()
let resumeTimeoutId: ReturnType<typeof setTimeout> | null = null

const activeSlide = computed(() => {
  if (props.slides.length === 0) return null
  return props.slides[currentSlide.value]
})

const isInWatchlist = computed(() => {
  if (!activeSlide.value) return false
  return auth.inWatchlist(activeSlide.value.mal_id)
})

const nextSlide = () => {
  if (props.slides.length === 0) return
  currentSlide.value = (currentSlide.value + 1) % props.slides.length
}

const { pause, resume } = useIntervalFn(nextSlide, 8000)

const handleManualInteraction = (index: number) => {
  currentSlide.value = index
  pause()
  
  if (resumeTimeoutId) {
    clearTimeout(resumeTimeoutId)
  }
  
  resumeTimeoutId = setTimeout(() => {
    resume()
  }, 10000)
}

const handleAddToWatchlist = async () => {
  if (!activeSlide.value) return
  if (!auth.isAuthenticated.value) {
    navigateTo('/login')
    return
  }
  const slide = activeSlide.value
  const imgUrl = slide.images?.jpg?.large_image_url || slide.images?.jpg?.image_url || ''
  await auth.addToWatchlist(slide.mal_id, slide.title_english || slide.title, imgUrl)
}

onUnmounted(() => {
  if (resumeTimeoutId) {
    clearTimeout(resumeTimeoutId)
  }
})
</script>

<template>
  <div v-if="slides.length > 0 && activeSlide" class="relative w-full min-h-[60dvh] md:min-h-[75dvh] lg:min-h-[80dvh] rounded-3xl overflow-hidden glass-surface-elevated border border-[var(--glass-border)]  group">
    <!-- Slide Transition -->
    <div class="absolute inset-0">
      <Transition name="hero-fade" mode="out-in">
        <div :key="currentSlide" class="absolute inset-0 w-full h-full">
          <!-- Ken Burns image with dynamic blur bleed -->
          <div class="w-full h-full overflow-hidden">
            <NuxtImg
              v-if="activeSlide.images?.jpg?.large_image_url || activeSlide.images?.jpg?.image_url"
              :src="activeSlide.images.jpg.large_image_url || activeSlide.images.jpg.image_url"
              :alt="activeSlide.title"
              class="w-full h-full object-cover animate-kenburns"
            />
          </div>

          <!-- Glass Dark Gradients -->
          <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/95 via-[var(--ui-overlay)]/50 to-transparent" />
          <div class="absolute inset-0 bg-gradient-to-r from-[var(--ui-overlay)]/80 via-transparent to-transparent hidden md:block" />

          <!-- Content overlay -->
          <div class="absolute bottom-8 md:bottom-12 left-6 md:left-12 right-6 max-w-2xl flex flex-col gap-4 z-10 text-white">
            <!-- Glass Metadata Badges -->
            <div class="flex flex-wrap items-center gap-2 text-xs font-semibold">
              <span v-if="activeSlide.score" class="glass-chip text-[var(--rank-gold)] px-2.5 py-1 rounded-lg flex items-center gap-1.5 font-mono shadow-sm">
                <UIcon name="i-solar-star-bold" class="w-3.5 h-3.5 text-yellow-400" />
                {{ Number(activeSlide.score).toFixed(1) }}
              </span>
              <span v-if="activeSlide.type" class="glass-chip text-[var(--ui-text-on-image)] px-2.5 py-1 rounded-lg uppercase tracking-wider text-[11px]">
                {{ activeSlide.type }}
              </span>
              <span v-if="activeSlide.episodes" class="glass-chip text-[var(--ui-text-on-image-muted)] px-2.5 py-1 rounded-lg text-[11px]">
                {{ activeSlide.episodes }} EPS
              </span>
              <span v-if="activeSlide.year" class="glass-chip text-[var(--ui-text-on-image-muted)] px-2.5 py-1 rounded-lg text-[11px] font-mono">
                {{ activeSlide.year }}
              </span>
            </div>

            <!-- Title -->
            <h2 class="text-2xl md:text-4xl lg:text-5xl font-bold tracking-tight text-white leading-tight drop-shadow-[var(--shadow-diffuse)]">
              {{ activeSlide.title_english || activeSlide.title }}
            </h2>
            
            <!-- Synopsis -->
            <p v-if="activeSlide.synopsis" class="text-xs md:text-sm text-[var(--ui-text-on-image-muted)]/90 line-clamp-2 max-w-xl leading-relaxed font-normal drop-shadow">
              {{ activeSlide.synopsis }}
            </p>

            <!-- CTA Actions -->
            <div class="flex items-center gap-3.5 mt-2">
              <UButton
                :to="`/movie/${activeSlide.mal_id}`"
                color="primary"
                size="lg"
                icon="i-solar-play-bold"
                label="View Details"
                class="font-bold px-6 py-3 rounded-2xl shadow-[var(--shadow-diffuse-accent)]  hover:scale-[1.02] active:scale-95 transition-all cursor-pointer"
              />
              <UButton
                v-if="!isInWatchlist"
                color="neutral"
                variant="ghost"
                size="lg"
                icon="i-solar-bookmark-linear"
                label="Watchlist"
                class="glass-chip text-white font-semibold px-5 py-3 rounded-2xl hover:bg-white/20 hover:scale-[1.02] active:scale-95 transition-all cursor-pointer"
                @click="handleAddToWatchlist"
              />
              <UButton
                v-else
                color="neutral"
                variant="ghost"
                size="lg"
                icon="i-solar-bookmark-bold"
                label="In Watchlist"
                class="glass-chip !text-primary-400 font-bold px-5 py-3 rounded-2xl border-primary-400/40 cursor-default"
                disabled
              />
            </div>
          </div>
        </div>
      </Transition>
    </div>

    <!-- Segmented Glass Timeline Indicators (Bottom Right) -->
    <div class="absolute bottom-8 md:bottom-12 right-6 md:right-12 flex items-center gap-2 z-20">
      <button
        v-for="(slide, index) in slides"
        :key="slide.mal_id"
        class="h-1.5 rounded-full transition-all duration-500 cursor-pointer"
        :class="index === currentSlide ? 'w-8 bg-gradient-to-r from-primary-400 to-primary-500 ' : 'w-2 bg-white/30 hover:bg-white/60'"
        @click="handleManualInteraction(index)"
        :aria-label="`Go to slide ${index + 1}`"
      />
    </div>
  </div>
</template>

<style scoped>
.hero-fade-enter-active,
.hero-fade-leave-active {
  transition: opacity 0.7s ease-in-out;
}
.hero-fade-enter-from,
.hero-fade-leave-to {
  opacity: 0;
}
</style>