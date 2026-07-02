<script setup lang="ts">
interface Genre {
  mal_id: number
  name: string
  count: number
}

const route = useRoute()
const isOpen = ref(false)
const genres = ref<Genre[]>([])
const loading = ref(true)

// Support both 'genres=1,2,3' and legacy 'genre=1'
const selectedGenreIds = computed(() => {
  const getNumbers = (val: typeof route.query.genres | typeof route.query.genre): number[] => {
    if (!val) return []
    const items = Array.isArray(val) ? val : [val]
    return items
      .filter((item): item is string => typeof item === 'string' && item.trim() !== '')
      .flatMap(item => item.split(','))
      .map(id => id.trim())
      .filter(id => id !== '' && !isNaN(Number(id)))
      .map(Number)
  }

  if (route.query.genres) return getNumbers(route.query.genres)
  if (route.query.genre) return getNumbers(route.query.genre)
  return []
})

onMounted(async () => {
  try {
    let res = await fetch('http://localhost:8080/api/genres')
    
    if (res.status === 429) {
      await new Promise(r => setTimeout(r, 2000))
      res = await fetch('http://localhost:8080/api/genres')
    }
    
    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`)
    }
    
    const d = await res.json()
    genres.value = (d.data ?? [])
      .filter((g: Genre) => g.count > 1000)
      .slice(0, 20)
  } catch (e) {
    console.error('Failed to fetch genres', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="relative">
    <UButton
      icon="i-solar-filter-linear"
      trailing-icon="i-solar-alt-arrow-down-linear"
      color="neutral"
      variant="outline"
      class="rounded-xl px-4 py-2.5 min-h-[44px] text-sm font-medium border-muted/80 bg-elevated hover:bg-default text-highlighted transition-all duration-200 cursor-pointer"
      aria-haspopup="listbox"
      :aria-expanded="isOpen"
      @click="isOpen = !isOpen"
    >
      <span>{{ selectedGenreIds.length > 0 ? 'Genres' : 'Filter by Genre' }}</span>
      <UBadge
        v-if="selectedGenreIds.length > 0"
        size="sm"
        color="primary"
        variant="solid"
        class="rounded-full px-1.5 py-0.5 text-xs font-semibold"
      >
        {{ selectedGenreIds.length }}
      </UBadge>
    </UButton>
  </div>
</template>
