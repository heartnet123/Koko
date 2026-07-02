<script setup lang="ts">
const slides = [
  {
    id: 1,
    title: 'Romance in\nSunset',
    badge: 'New Episode',
    schedule: 'Every Friday',
    description: 'A heartwarming story of love, fate and second chances.',
    image: 'https://images.unsplash.com/photo-1518173946687-a4c8892bbd9f?q=80&w=1200&auto=format&fit=crop',
  },
]

const current = ref(0)
</script>

<template>
  <section class="relative rounded-3xl overflow-hidden h-80 bg-elevated group">
    <NuxtImg
      :src="slides[current].image"
      alt="Featured anime"
      class="absolute inset-0 w-full h-full object-cover"
    />
    <!-- Gradient overlay -->
    <div class="absolute inset-0 bg-gradient-to-r from-black/60 via-black/30 to-transparent" />

    <!-- Content -->
    <div class="relative h-full flex flex-col justify-center px-12 max-w-xl">
      <h2 class="text-4xl font-semibold tracking-tight text-white leading-[1.1] mb-4 whitespace-pre-line">
        {{ slides[current].title }}
      </h2>
      <div class="flex items-center gap-4 text-sm mb-4">
        <span class="text-primary bg-primary/10 px-2 py-0.5 rounded text-xs font-medium">
          {{ slides[current].badge }}
        </span>
        <span class="text-white/70 font-medium">{{ slides[current].schedule }}</span>
      </div>
      <p class="text-sm text-white/70 leading-relaxed mb-8 max-w-sm">
        {{ slides[current].description }}
      </p>
      <div class="flex items-center gap-3">
        <UButton
          to="/movie/5114"
          icon="i-solar-play-bold"
          label="Play Now"
          color="primary"
          class="rounded-full shadow-lg shadow-primary/20"
        />
        <UButton
          icon="i-solar-add-circle-linear"
          label="Watchlist"
          color="neutral"
          variant="outline"
          class="rounded-full bg-white/10 backdrop-blur border-white/20 text-white hover:bg-white/20"
        />
      </div>
    </div>

    <!-- Arrow -->
    <div class="absolute right-8 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-all translate-x-4 group-hover:translate-x-0">
      <UButton
        icon="i-solar-alt-arrow-right-linear"
        color="neutral"
        variant="ghost"
        class="w-10 h-10 bg-white/80 backdrop-blur-md rounded-full shadow-sm"
        aria-label="Next slide"
      />
    </div>
  </section>
</template>
