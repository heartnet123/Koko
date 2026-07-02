<script setup lang="ts">
interface Genre {
  mal_id: number
  name: string
  count: number
}

const genres = ref<Genre[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/api/genres')
    const data = await res.json()
    genres.value = (data.data ?? []).filter((g: Genre) => g.count > 1000).slice(0, 10)
  } catch (e) {
    console.error('Failed to fetch genres', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="mt-8 mb-8 flex items-center justify-center h-48 text-toned">
    <UIcon name="i-solar-spinner-linear" class="w-6 h-6 animate-spin mr-2" />
    Loading categories...
  </div>

  <div v-else class="flex flex-col gap-12 mt-8 mb-8">
    <GenreRow
      v-for="(genre, idx) in genres"
      :key="genre.mal_id"
      :genre="genre"
      :index="idx"
    />
  </div>
</template>
