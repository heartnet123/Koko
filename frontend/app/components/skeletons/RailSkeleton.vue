<template>
  <div class="flex flex-col gap-5 px-6 md:px-10">
    <!-- Header Placeholder -->
    <div class="h-6 w-36 glass-pill animate-glass-shimmer rounded-lg" />
    <!-- Cards List -->
    <div class="flex gap-4 overflow-x-hidden pb-4">
      <div 
        v-for="i in 8" 
        :key="i"
        class="flex-shrink-0 w-[160px] md:w-[180px] flex flex-col gap-2.5"
      >
        <div class="aspect-[2/3] w-full glass-surface rounded-2xl animate-glass-shimmer border border-[var(--glass-border)]" />
        <!-- Title placeholder below card -->
        <div class="h-4 w-full glass-pill animate-glass-shimmer rounded-md" />
        <div class="h-3 w-2/3 glass-pill animate-glass-shimmer rounded-md" />
      </div>
    </div>
  </div>
</template>