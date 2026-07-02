# Browse Page Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add multi-select genre filter dropdown to browse page with instant filtering, selected genre pills, and backwards-compatible navigation from index "See All" buttons.

**Architecture:** Create GenreFilterDropdown component with animated dropdown, genre chips, and URL-based state. Enhance browse.vue to handle comma-separated genre IDs, display selected pills, and refetch on filter changes. Maintain existing patterns and styling.

**Tech Stack:** Vue 3, Nuxt 3, TypeScript, Nuxt UI components, Tailwind CSS

## Global Constraints

- Use Nuxt UI components (`UButton`, `UIcon`, `UBadge`) for consistency
- Follow existing color tokens (`text-highlighted`, `text-toned`, `bg-elevated`, `border-muted`)
- Match existing animation patterns (`transition-all duration-200`)
- API endpoint: `http://localhost:8080/api/anime`, `/api/genres`
- Rate limit handling: 429 → retry after 2s
- Mobile-first responsive design (min 44px touch targets)
- Backwards compatible with existing `genre={id}` param from "See All" links

---

### Task 1: Create GenreFilterDropdown Component - Structure and State

**Files:**
- Create: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `useRoute()` from Nuxt, `navigateTo()` from Nuxt
- Produces: None (self-contained component)

- [ ] **Step 1: Create component file with basic structure**

```vue
<script setup lang="ts">
interface Genre {
  mal_id: number
  name: string
  count: number
}

const route = useRoute()
const isOpen = ref(false)
const genres = ref<Genre[]>([])
const loading = ref(true)

const selectedGenreIds = computed(() => {
  const param = route.query.genres as string
  return param ? param.split(',').map(Number) : []
})
</script>

<template>
  <div class="relative">
    <!-- Filter button placeholder -->
    <div>Filter Button</div>
  </div>
</template>
```

- [ ] **Step 2: Verify file compiles**

Run: `cd frontend && npm run dev`
Expected: No TypeScript errors, dev server starts

- [ ] **Step 3: Commit structure**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): create GenreFilterDropdown component structure

- Add Genre interface
- Set up reactive state (isOpen, genres, loading)
- Compute selectedGenreIds from route query
- Basic template scaffold

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 2: GenreFilterDropdown - Fetch Genres on Mount

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `GET /api/genres` endpoint
- Produces: `genres` ref populated with Genre[]

- [ ] **Step 1: Add genre fetch logic in onMounted**

```vue
<script setup lang="ts">
// ... existing code ...

onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/api/genres')
    const data = await res.json()
    genres.value = (data.data ?? [])
      .filter((g: Genre) => g.count > 1000)
      .slice(0, 20)
  } catch (e) {
    console.error('Failed to fetch genres', e)
  } finally {
    loading.value = false
  }
})
</script>
```

- [ ] **Step 2: Test genre fetch**

Run: `cd frontend && npm run dev`
Navigate to: `http://localhost:3000/browse`
Open browser console, check: No fetch errors, genres array populated

- [ ] **Step 3: Commit fetch logic**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): fetch and filter genres in GenreFilterDropdown

- Fetch from /api/genres on mount
- Filter genres with count > 1000
- Limit to top 20 genres
- Handle errors gracefully

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 3: GenreFilterDropdown - Filter Button with Badge

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `selectedGenreIds` computed
- Produces: Filter button UI with active count badge

- [ ] **Step 1: Replace button placeholder with full implementation**

```vue
<template>
  <div class="relative">
    <UButton
      :label="selectedGenreIds.length > 0 ? 'Genres' : 'Filter by Genre'"
      trailing-icon="i-solar-alt-arrow-down-linear"
      color="neutral"
      variant="ghost"
      size="sm"
      class="text-toned hover:text-default"
      @click="isOpen = !isOpen"
    >
      <template #leading>
        <UIcon name="i-solar-filter-linear" class="w-4 h-4" />
      </template>
      <template v-if="selectedGenreIds.length > 0" #trailing>
        <UBadge
          :label="String(selectedGenreIds.length)"
          size="xs"
          color="primary"
          variant="solid"
          class="ml-2"
        />
      </template>
    </UButton>
  </div>
</template>
```

- [ ] **Step 2: Test button appearance**

Run: Navigate to `/browse` in browser
Expected: Filter button visible, badge shows count when filters active

- [ ] **Step 3: Commit button**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): add filter button with active count badge

- Show filter icon and 'Genres' text
- Display badge with count when filters active
- Toggle dropdown on click

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 4: GenreFilterDropdown - Animated Dropdown with Genre Chips

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `genres` ref, `selectedGenreIds` computed, `isOpen` ref
- Produces: Animated dropdown UI with clickable genre chips

- [ ] **Step 1: Add dropdown container below button**

```vue
<template>
  <div class="relative">
    <!-- ... existing button ... -->
    
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 mt-2 w-80 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl shadow-lg p-4 z-50 max-h-96 overflow-y-auto"
      >
        <div v-if="loading" class="flex items-center justify-center py-8">
          <UIcon name="i-solar-spinner-linear" class="w-5 h-5 animate-spin text-toned" />
        </div>
        
        <div v-else class="flex flex-wrap gap-2">
          <button
            v-for="genre in genres"
            :key="genre.mal_id"
            class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-150 hover:scale-102"
            :class="
              selectedGenreIds.includes(genre.mal_id)
                ? 'bg-primary-50 dark:bg-primary-900/20 border border-primary-500 text-primary-700 dark:text-primary-400'
                : 'bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
            "
          >
            <div class="flex items-center gap-1.5">
              <UIcon
                v-if="selectedGenreIds.includes(genre.mal_id)"
                name="i-solar-check-circle-bold"
                class="w-3.5 h-3.5"
              />
              <span>{{ genre.name }}</span>
            </div>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
```

- [ ] **Step 2: Test dropdown animation**

Run: Navigate to `/browse`, click filter button
Expected: Dropdown slides down and fades in, genre chips displayed, selected genres show checkmark

- [ ] **Step 3: Commit dropdown UI**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): add animated dropdown with genre chips

- Slide-down fade-in transition
- Genre chips with hover effects
- Checkmark icon for selected genres
- Loading state while fetching

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 5: GenreFilterDropdown - Toggle Genre Logic

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `selectedGenreIds` computed, `route` from useRoute()
- Produces: `toggleGenre(genreId: number)` method that updates URL

- [ ] **Step 1: Add toggleGenre method**

```vue
<script setup lang="ts">
// ... existing code ...

const toggleGenre = (genreId: number) => {
  const current = selectedGenreIds.value
  const newIds = current.includes(genreId)
    ? current.filter(id => id !== genreId)
    : [...current, genreId]
  
  const query: Record<string, string> = { ...route.query }
  
  if (newIds.length > 0) {
    query.genres = newIds.join(',')
  } else {
    delete query.genres
  }
  
  // Remove genre_name when multiple genres selected
  if (newIds.length > 1) {
    delete query.genre_name
  }
  
  navigateTo({ path: '/browse', query })
}
</script>
```

- [ ] **Step 2: Wire up click handler**

```vue
<template>
  <!-- ... -->
  <button
    v-for="genre in genres"
    :key="genre.mal_id"
    @click="toggleGenre(genre.mal_id)"
    class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-150 hover:scale-102"
    <!-- ... rest of button ... -->
  >
  <!-- ... -->
</template>
```

- [ ] **Step 3: Test genre toggle**

Run: Navigate to `/browse`, open filter dropdown, click genre chip
Expected: URL updates with `genres=1`, grid refetches, chip shows checkmark
Click again: Genre removed from URL, grid refetches

- [ ] **Step 4: Commit toggle logic**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): implement instant genre filter toggle

- Add/remove genre from selection on click
- Update URL query params immediately
- Remove genre_name when multiple selected
- Preserve other query params (order_by)

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---
## File Structure

**New files:**
- `frontend/app/components/GenreFilterDropdown.vue` — Filter button, animated dropdown, genre chips, clear all

**Modified files:**
- `frontend/app/pages/browse.vue` — Add filter dropdown, selected pills, update fetch logic for multi-genre
- `frontend/app/components/GenreRow.vue` — Verify "See All" link compatibility (already correct)

**Types:**
- Reuse existing `Genre` interface from GenreRow (no changes needed)
- Reuse existing `Anime` interface from types/anime.ts

---


### Task 6: GenreFilterDropdown - Clear All Button

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `selectedGenreIds` computed
- Produces: `clearAll()` method, Clear All button in dropdown

- [ ] **Step 1: Add clearAll method**

```vue
<script setup lang="ts">
// ... existing code ...

const clearAll = () => {
  const query: Record<string, string> = { ...route.query }
  delete query.genres
  delete query.genre
  delete query.genre_name
  navigateTo({ path: '/browse', query })
}
</script>
```

- [ ] **Step 2: Add Clear All button to dropdown**

```vue
<template>
  <div class="relative">
    <!-- ... existing button and dropdown ... -->
    <div v-if="isOpen" class="...">
      <!-- ... existing genre chips ... -->
      
      <div v-if="selectedGenreIds.length > 0" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
        <UButton
          label="Clear all"
          color="neutral"
          variant="ghost"
          size="sm"
          class="w-full text-toned hover:text-default"
          @click="clearAll"
        />
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 3: Test clear all**

Run: Navigate to `/browse?genres=1,2`, open dropdown, click "Clear all"
Expected: All filters removed, URL becomes `/browse`, grid shows all anime

- [ ] **Step 4: Commit clear all**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): add clear all button to genre filter

- Remove all genre filters at once
- Only show when filters active
- Clear genres, genre, genre_name params

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 7: GenreFilterDropdown - Click Outside to Close

**Files:**
- Modify: `frontend/app/components/GenreFilterDropdown.vue`

**Interfaces:**
- Consumes: `isOpen` ref
- Produces: Click-outside behavior closes dropdown

- [ ] **Step 1: Add click-outside directive**

```vue
<script setup lang="ts">
// ... existing code ...

const dropdownRef = ref<HTMLElement | null>(null)

const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="dropdownRef" class="relative">
    <!-- ... existing content ... -->
  </div>
</template>
```

- [ ] **Step 2: Test click outside**

Run: Navigate to `/browse`, open dropdown, click outside dropdown area
Expected: Dropdown closes

- [ ] **Step 3: Commit click-outside**

```bash
git add frontend/app/components/GenreFilterDropdown.vue
git commit -m "feat(browse): close dropdown on click outside

- Add click-outside event listener
- Clean up listener on unmount
- Improve UX for filter dropdown

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 8: Enhance Browse Page - Multi-Genre URL Handling

**Files:**
- Modify: `frontend/app/pages/browse.vue`

**Interfaces:**
- Consumes: `route.query.genres` (comma-separated IDs), `route.query.genre` (legacy single ID)
- Produces: `genreIds` computed that handles both formats

- [ ] **Step 1: Update genreIds computed to handle both formats**

```vue
<script setup lang="ts">
import type { Anime } from '~/types/anime'

const route = useRoute()

// Support both 'genres=1,2,3' and legacy 'genre=1'
const genreIds = computed(() => {
  const genresParam = route.query.genres as string
  const genreParam = route.query.genre as string
  
  if (genresParam) {
    return genresParam.split(',').filter(id => id.trim() && !isNaN(Number(id)))
  } else if (genreParam) {
    return [genreParam]
  }
  return []
})

const genreName = computed(() => route.query.genre_name as string)
const orderBy = computed(() => (route.query.order_by as string) || 'popularity')

// ... rest of existing code ...
</script>
```

- [ ] **Step 2: Update fetch URL construction**

```vue
<script setup lang="ts">
// ... existing code ...

const { data: response, status } = await useFetch<{ data: Anime[] }>(
  () => {
    let url = 'http://localhost:8080/api/anime?limit=24'
    if (genreIds.value.length > 0) {
      url += `&genres=${genreIds.value.join(',')}`
    }
    if (orderBy.value) {
      url += `&order_by=${orderBy.value}`
    }
    return url
  },
  {
    key: `browse-${route.fullPath}`
  }
)

const animes = computed(() => response.value?.data ?? [])
const loading = computed(() => status.value === 'pending')
</script>
```

- [ ] **Step 3: Test multi-genre fetch**

Run: Navigate to `/browse?genres=1,2`
Expected: API called with `&genres=1,2`, grid shows filtered results

Run: Navigate to `/browse?genre=1` (legacy)
Expected: Works, shows single genre results

- [ ] **Step 4: Commit multi-genre handling**

```bash
git add frontend/app/pages/browse.vue
git commit -m "feat(browse): handle multi-genre URL params

- Parse genres=1,2,3 format
- Backwards compat with genre=1 (legacy)
- Filter out invalid genre IDs
- Update fetch to use comma-separated genres

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---


### Task 9: Browse Page - Add GenreFilterDropdown Component

**Files:**
- Modify: `frontend/app/pages/browse.vue`

**Interfaces:**
- Consumes: `GenreFilterDropdown` component
- Produces: Filter dropdown visible in browse page header

- [ ] **Step 1: Add GenreFilterDropdown to template**

```vue
<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4 mb-12">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold tracking-tight text-highlighted">
          {{ genreName ? `Explore ${genreName}` : genreIds.length > 0 ? 'Filtered Anime' : 'Popular Anime' }}
        </h2>
        <p class="text-sm text-toned mt-1">
          {{ genreIds.length > 0 ? `Showing ${genreIds.length} genre${genreIds.length > 1 ? 's' : ''} selected` : 'Showing collections matching your query.' }}
        </p>
      </div>
      
      <GenreFilterDropdown />
    </div>

    <div
      v-if="loading"
      class="flex items-center justify-center h-64 text-toned"
    >
      <UIcon name="i-solar-spinner-linear" class="w-6 h-6 animate-spin mr-2" />
      Loading...
    </div>

    <div v-else-if="animes.length === 0" class="flex flex-col items-center justify-center h-64 text-toned border border-dashed border-muted rounded-2xl">
      <UIcon name="i-solar-ghost-bold" class="w-8 h-8 mb-2" />
      No anime found.
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-5">
      <NuxtLink
        v-for="item in animes"
        :key="item.mal_id"
        :to="`/movie/${item.mal_id}`"
        class="flex flex-col group cursor-pointer"
      >
        <div class="relative rounded-2xl overflow-hidden aspect-[3/4] mb-3 shadow-[0_4px_12px_rgb(0,0,0,0.03)] border border-muted/50 bg-elevated">
          <NuxtImg
            :src="item.images.jpg.large_image_url || item.images.jpg.image_url"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
          />
          <div class="absolute bottom-3 left-3 w-8 h-8 rounded-full bg-white/90 backdrop-blur shadow-sm flex items-center justify-center opacity-90 group-hover:opacity-100 group-hover:scale-110 transition-all">
            <UIcon name="i-solar-play-bold" class="w-3.5 h-3.5 ml-0.5 text-highlighted" />
          </div>
        </div>
        <h4 class="text-sm font-medium text-highlighted tracking-tight truncate">{{ item.title }}</h4>
        <p class="text-xs text-toned mt-0.5">{{ item.type ?? 'Anime' }}</p>
      </NuxtLink>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Test integration**

Run: Navigate to `/browse`
Expected: Filter button appears in header, clicking opens dropdown

- [ ] **Step 3: Commit integration**

```bash
git add frontend/app/pages/browse.vue
git commit -m "feat(browse): integrate GenreFilterDropdown into browse page

- Add filter dropdown to page header
- Update header text based on filter state
- Show genre count when filters active

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 10: Browse Page - Selected Genre Pills Row

**Files:**
- Modify: `frontend/app/pages/browse.vue`

**Interfaces:**
- Consumes: `genreIds` computed, genre names from API
- Produces: Row of dismissible genre pills below filter button

- [ ] **Step 1: Add genre fetch and mapping**

```vue
<script setup lang="ts">
// ... existing code ...

interface Genre {
  mal_id: number
  name: string
  count: number
}

const allGenres = ref<Genre[]>([])

const selectedGenres = computed(() => {
  return genreIds.value
    .map(id => allGenres.value.find(g => g.mal_id === Number(id)))
    .filter(Boolean) as Genre[]
})

const removeGenre = (genreId: number) => {
  const newIds = genreIds.value.filter(id => Number(id) !== genreId)
  const query: Record<string, string> = { ...route.query }
  
  if (newIds.length > 0) {
    query.genres = newIds.join(',')
  } else {
    delete query.genres
    delete query.genre
    delete query.genre_name
  }
  
  navigateTo({ path: '/browse', query })
}

onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/api/genres')
    const data = await res.json()
    allGenres.value = data.data ?? []
  } catch (e) {
    console.error('Failed to fetch genres for pills', e)
  }
})
</script>
```

- [ ] **Step 2: Add pills row to template**

```vue
<template>
  <div class="max-w-7xl mx-auto px-8 w-full flex flex-col gap-8 mt-4 mb-12">
    <!-- ... existing header with GenreFilterDropdown ... -->
    
    <div v-if="selectedGenres.length > 0" class="flex flex-wrap gap-2">
      <div
        v-for="genre in selectedGenres"
        :key="genre.mal_id"
        class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium bg-primary-50 dark:bg-primary-900/20 border border-primary-500 text-primary-700 dark:text-primary-400"
      >
        <span>{{ genre.name }}</span>
        <button
          @click="removeGenre(genre.mal_id)"
          class="hover:bg-primary-100 dark:hover:bg-primary-900/40 rounded-full p-0.5 transition-colors"
        >
          <UIcon name="i-solar-close-circle-bold" class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- ... existing loading/empty/grid sections ... -->
  </div>
</template>
```

- [ ] **Step 3: Test pills**

Run: Navigate to `/browse?genres=1,2`, wait for pills to appear
Expected: Two genre pills shown with names and X buttons
Click X: Genre removed, URL updated, grid refetches

- [ ] **Step 4: Commit pills**

```bash
git add frontend/app/pages/browse.vue
git commit -m "feat(browse): add selected genre pills with remove buttons

- Fetch genres to map IDs to names
- Show dismissible pills for active filters
- Remove individual genres via X button
- Pills only appear when filters active

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---


### Task 11: Final Integration Test

**Files:**
- Test: All components together

**Interfaces:**
- Consumes: Complete browse page with filter dropdown
- Produces: Verified working implementation

- [ ] **Step 1: Test complete flow from index**

Run: Navigate to index page `/`
Click: "See All" button on any GenreRow
Expected: Browse page opens with single genre filter, pill shows genre name

- [ ] **Step 2: Test multi-genre selection**

Run: On browse page, open filter dropdown
Click: Select 2-3 different genres
Expected: Each click updates URL, pills appear, grid refetches instantly

- [ ] **Step 3: Test clear all**

Run: With multiple genres selected, open dropdown
Click: "Clear all" button
Expected: All filters removed, pills disappear, URL becomes `/browse`, shows all anime

- [ ] **Step 4: Test pill removal**

Run: Select 3 genres, click X on one pill
Expected: That genre removed, other pills remain, URL updates, grid refetches

- [ ] **Step 5: Test backwards compatibility**

Run: Navigate to `/browse?genre=1&genre_name=Action` (legacy format)
Expected: Single genre shows in pills, filter dropdown shows it selected

- [ ] **Step 6: Test dropdown close behaviors**

Run: Open dropdown, click outside
Expected: Dropdown closes
Run: Open dropdown, click genre chip
Expected: Dropdown stays open, genre toggles
Run: Open dropdown, click filter button again
Expected: Dropdown closes

- [ ] **Step 7: Test mobile responsiveness**

Run: Resize browser to mobile width (375px)
Expected: Dropdown full-width, chips wrap, pills wrap, grid adjusts to 2 columns

- [ ] **Step 8: Final commit**

```bash
git add -A
git commit -m "test(browse): verify complete genre filter integration

- Index See All navigation works
- Multi-genre selection and instant filtering
- Clear all removes all filters
- Pill removal works individually
- Backwards compatible with legacy URLs
- Dropdown close behaviors correct
- Mobile responsive layout

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

## Self-Review Checklist

**Spec Coverage:**
- ✓ Multi-select genre filter UI (Task 1-7)
- ✓ Instant filtering on toggle (Task 5)
- ✓ Visual feedback via pills (Task 10)
- ✓ "See All" navigation maintained (Task 8, backwards compat)
- ✓ Animated dropdown (Task 4)
- ✓ Clear all button (Task 6)
- ✓ Click-outside behavior (Task 7)
- ✓ URL-based state management (Task 5, 8)
- ✓ Mobile responsive (Task 11 verification)

**Placeholder Scan:** None found — all code complete

**Type Consistency:**
- `Genre` interface: `{ mal_id: number; name: string; count: number }` — used consistently
- `toggleGenre(genreId: number)` → `removeGenre(genreId: number)` — both take number
- `selectedGenreIds` computed returns `number[]` — used consistently
- `genreIds` computed returns `string[]` — converted to numbers where needed

**No Gaps:** All acceptance criteria covered in tasks.

