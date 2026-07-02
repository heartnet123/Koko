<script setup lang="ts">
useSeoMeta({
  title: 'Koko - Anime Library',
  ogTitle: 'Koko - Anime Library',
  description: 'Explore the ultimate collection of anime movies and series on Koko.',
  ogDescription: 'Explore the ultimate collection of anime movies and series on Koko.',
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4">
    <HeroSection />
    <RecommendedSection />
    <GenreSection />
  </div>
</template>
