<script setup lang="ts">
const search = ref('')
defineEmits<{
  (e: 'toggle-menu'): void
}>()
</script>

<template>
  <header class="h-20 px-4 md:px-8 flex items-center justify-between sticky top-0 bg-default/80 backdrop-blur-md z-40 border-b border-muted/50">
    <!-- Left Menu Trigger & Search -->
    <div class="flex items-center gap-3 w-full max-w-lg">
      <UButton
        icon="i-solar-menu-hamburger-linear"
        color="neutral"
        variant="ghost"
        class="lg:hidden rounded-full cursor-pointer focus:outline-none"
        aria-label="Open menu"
        @click="$emit('toggle-menu')"
      />
      <div class="relative flex-1 group">
        <UIcon
          name="i-solar-magnifer-linear"
          class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-toned group-focus-within:text-primary transition-colors"
        />
        <input
          v-model="search"
          type="search"
          placeholder="Search anime, genres, studios..."
          aria-label="Search anime, genres, studios"
          class="w-full pl-10 pr-4 py-2.5 bg-elevated rounded-full text-sm text-default placeholder:text-toned border border-muted focus:outline-none focus:border-primary/40 focus:bg-default transition-all"
        />
      </div>
    </div>

    <!-- User actions -->
    <div class="flex items-center gap-3 ml-6">
      <UButton
        icon="i-solar-bell-linear"
        color="neutral"
        variant="ghost"
        size="md"
        class="rounded-full"
        aria-label="Notifications"
      />
      <UAvatar
        src="https://i.pravatar.cc/40"
        alt="User avatar"
        size="sm"
        class="cursor-pointer ring-2 ring-primary/20 hover:ring-primary/50 transition-all"
      />
    </div>
  </header>
</template>
