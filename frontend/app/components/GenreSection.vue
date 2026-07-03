<script setup lang="ts">
interface Genre {
  mal_id: number
  name: string
  count: number
}

const { data: genresResponse, status } = await useFetch<{ data: Genre[] }>(
  'http://localhost:8080/api/genres',
  {
    key: 'genres-top-list',
    retry: 1,
    retryDelay: 2000,
    retryStatusCodes: [429]
  }
)

const genres = computed(() => {
  return (genresResponse.value?.data ?? [])
    .filter((g: Genre) => g.count > 1000)
    .slice(0, 10)
})

const loading = computed(() => status.value === 'pending')
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
