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
</script>

<template>
  <div class="relative">
    <!-- Filter button placeholder -->
    <div>Filter Button</div>
  </div>
</template>
