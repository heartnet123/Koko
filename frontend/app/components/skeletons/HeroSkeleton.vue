<template>
  <div class="relative w-full min-h-[60dvh] md:min-h-[80dvh] bg-elevated rounded-3xl overflow-hidden">
    <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent"></div>
    <div class="absolute bottom-10 left-6 md:left-10 right-6 max-w-2xl flex flex-col gap-4 z-10">
      <!-- Title placeholder -->
      <div class="h-8 md:h-12 w-3/4 animate-shimmer-gradient rounded-lg"></div>
      <!-- Synopsis placeholder (2 lines) -->
      <div class="space-y-2">
        <div class="h-4 w-full animate-shimmer-gradient rounded-lg"></div>
        <div class="h-4 w-5/6 animate-shimmer-gradient rounded-lg"></div>
      </div>
      <!-- Meta row placeholder -->
      <div class="h-4 w-1/3 animate-shimmer-gradient rounded-lg mt-2"></div>
      <!-- 2 CTA buttons placeholder -->
      <div class="flex items-center gap-4 mt-4">
        <div class="h-10 w-32 animate-shimmer-gradient rounded-lg"></div>
        <div class="h-10 w-32 animate-shimmer-gradient rounded-lg"></div>
      </div>
    </div>
    <!-- Dot indicator row -->
    <div class="absolute bottom-10 right-6 md:right-10 flex gap-2 z-10">
      <div class="w-2.5 h-2.5 rounded-full animate-shimmer-gradient" v-for="i in 5" :key="i"></div>
    </div>
  </div>
</template>
