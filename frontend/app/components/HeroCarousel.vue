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
let resumeTimeoutId: any = null

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
  <div v-if="slides.length > 0 && activeSlide" class="relative w-full min-h-[60dvh] md:min-h-[80dvh] bg-black rounded-3xl overflow-hidden group">
    <!-- Slide Transition -->
    <div class="absolute inset-0">
      <Transition name="hero-fade" mode="out-in">
        <div :key="currentSlide" class="absolute inset-0 w-full h-full">
          <!-- Ken Burns image container -->
          <div class="w-full h-full overflow-hidden">
            <NuxtImg
              v-if="activeSlide.images?.jpg?.large_image_url || activeSlide.images?.jpg?.image_url"
              :src="activeSlide.images.jpg.large_image_url || activeSlide.images.jpg.image_url"
              :alt="activeSlide.title"
              class="w-full h-full object-cover animate-kenburns"
            />
          </div>

          <!-- Linear overlay -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>

          <!-- Content stack lower-left -->
          <div class="absolute bottom-10 left-6 md:left-10 right-6 max-w-2xl flex flex-col gap-3 md:gap-4 z-10 text-white">
            <h2 class="text-2xl md:text-4xl lg:text-5xl font-bold tracking-tight text-white drop-shadow-md">
              {{ activeSlide.title_english || activeSlide.title }}
            </h2>
            
            <p v-if="activeSlide.synopsis" class="text-sm md:text-base text-neutral-300 line-clamp-2 max-w-xl leading-relaxed drop-shadow">
              {{ activeSlide.synopsis }}
            </p>

            <div class="flex flex-wrap items-center gap-3 text-xs md:text-sm font-mono text-neutral-400">
              <span v-if="activeSlide.type" class="uppercase">{{ activeSlide.type }}</span>
              <span v-if="activeSlide.episodes" class="before:content-['•'] before:mr-2">{{ activeSlide.episodes }} EPS</span>
              <span v-if="activeSlide.score" class="before:content-['•'] before:mr-2 flex items-center gap-1">
                <UIcon name="i-solar-star-bold" class="w-3.5 h-3.5 text-yellow-400" />
                {{ activeSlide.score }}
              </span>
              <span v-if="activeSlide.year" class="before:content-['•'] before:mr-2">{{ activeSlide.year }}</span>
            </div>

            <!-- CTA Actions -->
            <div class="flex items-center gap-4 mt-2">
              <UButton
                :to="`/movie/${activeSlide.mal_id}`"
                color="primary"
                size="md"
                icon="i-solar-play-bold"
                label="View Details"
                class="font-medium px-5 cursor-pointer"
              />
              <UButton
                v-if="!isInWatchlist"
                color="neutral"
                variant="subtle"
                size="md"
                icon="i-solar-bookmark-linear"
                label="Watchlist"
                class="font-medium px-5 cursor-pointer"
                @click="handleAddToWatchlist"
              />
              <UButton
                v-else
                color="neutral"
                variant="subtle"
                size="md"
                icon="i-solar-bookmark-bold"
                label="In Watchlist"
                class="font-medium px-5 text-primary-400 cursor-pointer"
                disabled
              />
            </div>
          </div>
        </div>
      </Transition>
    </div>

    <!-- Dot Indicators bottom-right -->
    <div class="absolute bottom-10 right-6 md:right-10 flex gap-2.5 z-20">
      <button
        v-for="(slide, index) in slides"
        :key="slide.mal_id"
        class="w-2.5 h-2.5 rounded-full transition-all duration-300 cursor-pointer"
        :class="index === currentSlide ? 'bg-primary-500 scale-125' : 'bg-white/40 hover:bg-white/70'"
        @click="handleManualInteraction(index)"
        :aria-label="`Go to slide ${index + 1}`"
      ></button>
    </div>
  </div>
</template>

<style scoped>
.hero-fade-enter-active,
.hero-fade-leave-active {
  transition: opacity 0.8s ease-in-out;
}
.hero-fade-enter-from,
.hero-fade-leave-to {
  opacity: 0;
}
</style>
