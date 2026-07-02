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
  const genresVal = route.query.genres
  const genreVal = route.query.genre
  
  if (genresVal) {
    const genresStr = Array.isArray(genresVal) ? genresVal[0] : genresVal
    if (typeof genresStr === 'string') {
      return genresStr.split(',')
        .map(id => id.trim())
        .filter(id => id && !isNaN(Number(id)))
        .map(Number)
    }
  } else if (genreVal) {
    const genreStr = Array.isArray(genreVal) ? genreVal[0] : genreVal
    const num = Number(genreStr)
    return isNaN(num) ? [] : [num]
  }
  return []
})
</script>

<template>
  <div class="relative">
    <!-- Filter button placeholder -->
    <div>Filter Button</div>
  </div>
</template>
