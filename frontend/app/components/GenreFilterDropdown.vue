<script setup lang="ts">
interface Genre {
  mal_id: number;
  name: string;
  count: number;
}

const route = useRoute();
const isOpen = ref(false);
const genres = ref<Genre[]>([
  { mal_id: 1, name: "Action", count: 1254 },
  { mal_id: 5, name: "Adventure", count: 892 },
  { mal_id: 8, name: "Drama", count: 1987 },
  { mal_id: 10, name: "Fantasy", count: 634 },
  { mal_id: 9, name: "Ecchi", count: 298 },
  { mal_id: 12, name: "Adventure", count: 892 },
] as const);
const loading = ref(false);

// Demos selectedGenreIds for local dev from query string; in real app backend supplies active filters via query
const selectedGenreIds = computed(() => {
  const param = route.query.genres as string | undefined;
  return param ? param.split(',').map(Number) : [];
});

function toggleGenre(id: number) {
  // This will be wired to next task: update route and fetch filtered results
}

function clearAll() {
  // This will be wired to next task
}
</script>

<template>
  <div class="relative w-full min-w-[160px]">
    <!-- Filter button -->
    <button
      aria-haspopup="listbox"
      aria-expanded="false"
      @click="isOpen = !isOpen"
      class="flex items-center w-full gap-2 px-4 py-2 text-left bg-white rounded-md shadow-sm ring-1 ring-gray-300 hover:bg-gray-50 disabled:opacity-70"
      :disabled="loading"
    >
      <slot name="button-text">
        <span class="font-medium">
          Genres <span v-if="selectedGenreIds.length" class="ml-1 font-semibold text-blue-600">({{ selectedGenreIds.length }})</span>
        </span>
      </slot>
      <!-- Chevron indicator -->
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="ml-auto h-4 w-4 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fill-rule="evenodd"
          d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
          clip-rule="evenodd"
        />
      </svg>
    </button>

    <!-- Dropdown panel -->
    <div
      v-if="isOpen"
      @click.away="isOpen = false"
      role="listbox"
      class="absolute z-10 mt-1 w-full overflow-hidden rounded-md bg-white shadow-lg ring-1 ring-gray-300"
    >
      <!-- Header controls -->
      <div class="flex items-center justify-between p-3 border-b border-gray-200">
        <button
          class="px-3 py-1 text-sm text-gray-700 rounded-md hover:bg-gray-100"
          :disabled="selectedGenreIds.length === 0"
          @click="clearAll"
        >
          Clear All
        </button>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="p-4 text-sm text-gray-500">Loading genres...</div>

      <!-- Listbox options -->
      <ul v-else class="max-h-60 overflow-auto p-1" role="listbox">
        <li v-for="genre in genres" :key="genre.mal_id" role="option" class="contents">
          <label class="flex items-center gap-2 px-4 py-2 cursor-pointer hover:bg-gray-100 focus-within:bg-gray-100">
            <span class="relative">
              <input
                type="checkbox"
                :checked="selectedGenreIds.includes(genre.mal_id)"
                @change="toggleGenre(genre.mal_id)"
                @keydown.enter.prevent="toggleGenre(genre.mal_id)"
                @keydown.space.prevent="toggleGenre(genre.mal_id)"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                :aria-label="`Toggle ${genre.name} filter`"
                tabindex="0"
              />
            </span>
            <span class="text-sm font-normal">{{ genre.name }}</span>
            <span class="ml-2 text-xs font-mono text-gray-500">{{ genre.count }}</span>
          </label>
        </li>
      </ul>
    </div>
  </div>
</template>
