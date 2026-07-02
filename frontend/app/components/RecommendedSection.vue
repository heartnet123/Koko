<script setup lang="ts">
interface Anime {
  mal_id: number
  title: string
  images: { jpg: { image_url: string; large_image_url: string } }
  score: number | null
  type: string | null
}

const recommendations = ref<Anime[]>([])
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
    const res = await fetch('http://localhost:8080/api/anime?limit=10&order_by=popularity')
    const data = await res.json()
    const seen = new Set<number>()
    recommendations.value = (data.data ?? []).filter((a: Anime) => {
      if (seen.has(a.mal_id)) return false
      seen.add(a.mal_id)
      return true
    })
  } catch (e) {
    console.error('Failed to fetch recommendations', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section>
    <div class="flex items-center justify-between mb-5">
      <h3 class="text-lg font-medium tracking-tight text-highlighted">Recommended</h3>
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-1.5">
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
          label="See All"
          to="/browse?order_by=popularity"
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
      Loading recommendations...
    </div>

    <div v-else ref="scrollContainer" class="flex gap-5 overflow-x-auto scroll-smooth snap-x snap-mandatory scrollbar-none pb-4">
      <NuxtLink
        v-for="item in recommendations"
        :key="item.mal_id"
        :to="`/movie/${item.mal_id}`"
        class="w-[calc((100%-80px)/5)] flex-shrink-0 snap-start flex flex-col group cursor-pointer"
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
