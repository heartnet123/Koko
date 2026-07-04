<template>
  <div class="flex flex-col gap-5 px-6 md:px-10">
    <!-- Header Placeholder -->
    <div class="h-5 w-32 animate-shimmer-gradient rounded-md"></div>
    <!-- Cards List -->
    <div class="flex gap-4 overflow-x-hidden pb-4">
      <div 
        v-for="i in 8" 
        :key="i"
        class="flex-shrink-0 aspect-[2/3] w-[160px] md:w-[180px] flex flex-col gap-2"
      >
        <div class="w-full h-full animate-shimmer-gradient rounded-xl"></div>
        <!-- Title placeholder below card -->
        <div class="h-4 w-full animate-shimmer-gradient rounded-md"></div>
        <div class="h-3 w-2/3 animate-shimmer-gradient rounded-md"></div>
      </div>
    </div>
  </div>
</template>
