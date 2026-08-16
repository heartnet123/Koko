<template>
  <div class="relative w-full min-h-[60dvh] md:min-h-[80dvh] glass-surface rounded-3xl overflow-hidden border border-[var(--glass-border)]">
    <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/80 via-[var(--ui-overlay)]/30 to-transparent" />
    <div class="absolute bottom-10 left-6 md:left-10 right-6 max-w-2xl flex flex-col gap-4 z-10">
      <!-- Title placeholder -->
      <div class="h-8 md:h-12 w-3/4 glass-pill animate-glass-shimmer rounded-xl" />
      <!-- Synopsis placeholder (2 lines) -->
      <div class="space-y-2">
        <div class="h-4 w-full glass-pill animate-glass-shimmer rounded-lg" />
        <div class="h-4 w-5/6 glass-pill animate-glass-shimmer rounded-lg" />
      </div>
      <!-- Meta row placeholder -->
      <div class="h-5 w-1/3 glass-pill animate-glass-shimmer rounded-lg mt-2" />
      <!-- 2 CTA buttons placeholder -->
      <div class="flex items-center gap-4 mt-4">
        <div class="h-10 w-32 glass-pill animate-glass-shimmer rounded-xl" />
        <div class="h-10 w-32 glass-pill animate-glass-shimmer rounded-xl" />
      </div>
    </div>
    <!-- Dot indicator row -->
    <div class="absolute bottom-10 right-6 md:right-10 flex gap-2 z-10">
      <div v-for="i in 5" :key="i" class="w-2.5 h-2.5 rounded-full glass-pill animate-glass-shimmer" />
    </div>
  </div>
</template>