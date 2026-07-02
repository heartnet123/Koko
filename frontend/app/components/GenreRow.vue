<script setup lang="ts">
import type { Anime } from '~/types/anime'

interface Genre {
  mal_id: number
  name: string
  count: number
}

const props = defineProps<{ genre: Genre; index: number }>()

const GENRE_STYLES = [
  { icon: 'i-solar-heart-bold', color: 'text-pink-500' },
  { icon: 'i-solar-bolt-bold', color: 'text-yellow-500' },
  { icon: 'i-solar-ghost-bold', color: 'text-purple-500' },
  { icon: 'i-solar-sword-bold', color: 'text-red-500' },
  { icon: 'i-solar-planet-bold', color: 'text-blue-500' },
  { icon: 'i-solar-comedy-bold', color: 'text-green-500' },
  { icon: 'i-solar-mask-bold', color: 'text-orange-500' },
  { icon: 'i-solar-star-bold', color: 'text-cyan-500' },
]

const style = computed(() => GENRE_STYLES[props.index % GENRE_STYLES.length])

const animes = ref<Anime[]>([])
const loading = ref(true)
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

onMounted(async () => {
  try {
    const res = await fetch(
      `http://localhost:8080/api/anime?genres=${props.genre.mal_id}&limit=10&order_by=popularity`
    )
    if (res.status === 429) {
      await new Promise(r => setTimeout(r, 2000))
      const retry = await fetch(
        `http://localhost:8080/api/anime?genres=${props.genre.mal_id}&limit=10&order_by=popularity`
      )
      const d = await retry.json()
      animes.value = d.data ?? []
    } else {
      const d = await res.json()
      animes.value = d.data ?? []
    }
  } catch (e) {
    console.error(`Failed to fetch animes for ${props.genre.name}`, e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="flex flex-col">
    <div class="flex items-center justify-between mb-5">
      <div class="flex items-center gap-3">
        <div :class="`flex items-center justify-center w-10 h-10 rounded-xl bg-current/10 border border-current/10 shadow-sm ${style.color}`">
          <UIcon :name="style.icon" class="w-5 h-5" />
        </div>
        <h3 class="text-lg font-medium tracking-tight text-highlighted">Explore {{ genre.name }}</h3>
      </div>
      <div class="flex items-center gap-4">
        <div class="hidden md:flex items-center gap-1.5">
          <UButton
            icon="i-solar-alt-arrow-left-linear"
            color="neutral"
            variant="ghost"
            size="sm"
            class="text-toned hover:text-default"
            @click="scroll('left')"
          />
          <UButton
            icon="i-solar-alt-arrow-right-linear"
            color="neutral"
            variant="ghost"
            size="sm"
            class="text-toned hover:text-default"
            @click="scroll('right')"
          />
        </div>
        <UButton
          :label="`See All ${genre.name}`"
          :to="`/browse?genre=${genre.mal_id}&genre_name=${genre.name}`"
          trailing-icon="i-solar-alt-arrow-right-linear"
          color="neutral"
          variant="ghost"
          size="sm"
          class="text-toned hover:text-default"
        />
      </div>
    </div>

    <div
      v-if="loading"
      class="flex items-center justify-center h-32 text-sm text-toned bg-elevated rounded-2xl border border-muted border-dashed"
    >
      Loading {{ genre.name }} anime...
    </div>

    <div v-else ref="scrollContainer" class="flex gap-5 overflow-x-auto scroll-smooth snap-x snap-mandatory scrollbar-none pb-4">
      <NuxtLink
        v-for="item in animes"
        :key="item.mal_id"
        :to="`/movie/${item.mal_id}`"
        class="w-[calc((100%-20px)/2)] sm:w-[calc((100%-40px)/3)] md:w-[calc((100%-60px)/4)] lg:w-[calc((100%-80px)/5)] flex-shrink-0 snap-start flex flex-col group cursor-pointer"
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
