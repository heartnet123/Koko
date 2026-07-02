# Browse Page Enhancement Design

**Date:** 2026-07-02  
**Status:** Approved

## Overview

Enhance the browse page (`/browse`) with a multi-select genre filter dropdown and ensure seamless navigation from index page "See All" buttons. Add standard browse page features including sorting, filtering, and grid display improvements.

## Current State

- Browse page exists at `/pages/browse.vue`
- Basic grid display with 24 anime items
- Accepts `genre`, `genre_name`, and `order_by` query params
- Index page has GenreRow and RecommendedSection components
- "See All" buttons already navigate to browse page with appropriate filters
- API endpoint: `http://localhost:8080/api/anime`

## Goals

1. Add genre filter UI with multi-select capability
2. Instant filtering on each genre toggle
3. Visual feedback for active filters (chips/pills)
4. Maintain existing "See All" navigation from index
5. Add standard browse features (sorting, grid improvements)

## Architecture

### Component Structure

```
pages/
  browse.vue (enhanced)
components/
  GenreFilterDropdown.vue (new)
  GenreSection.vue (existing)
  GenreRow.vue (existing)
  RecommendedSection.vue (existing)
types/
  anime.ts (existing)
```

### State Management

- **URL as source of truth:** Selected genres stored in query param `genres=1,2,8` (comma-separated IDs)
- **Route reactivity:** `useRoute()` to read, `navigateTo()` to update
- **Local component state:** Dropdown open/closed (ref)
- **Genre list:** Fetched once on mount, cached in ref

### Data Flow

1. User clicks filter button → dropdown animates open
2. User clicks genre chip → toggle selection
3. URL updates immediately via `navigateTo()`
4. `useFetch` reactive key changes → grid refetches
5. Selected genres render as dismissible pills below button

## Components

### GenreFilterDropdown.vue (New)

**Props:** None (reads from route)

**Features:**
- Filter button with badge showing active filter count
- Animated dropdown (slide-down + fade-in)
- Genre chips inside dropdown
  - Checkmark icon when selected
  - Toggle on click (instant apply)
- "Clear all" button when filters active
- Click-outside-to-close behavior
- Loading state while fetching genres

**State:**
```typescript
const isOpen = ref(false)
const genres = ref<Genre[]>([])
const loading = ref(true)
const route = useRoute()
const selectedGenreIds = computed(() => {
  const param = route.query.genres as string
  return param ? param.split(',').map(Number) : []
})
```

**Methods:**
- `toggleGenre(genreId: number)` — add/remove from selection, update URL
- `clearAll()` — remove all genre filters from URL
- `closeDropdown()` — close dropdown (for click-outside)

**Animation:**
- Dropdown: `transition-all duration-200 ease-out`
- Slide from button position + fade opacity 0 → 1
- Scale from 95% → 100%

### Browse.vue (Enhanced)

**Query Params:**
- `genres` (string) — comma-separated genre IDs: `"1,2,8"`
- `order_by` (string) — sort order: `"popularity"`, `"score"`, `"title"`
- `genre_name` (string) — legacy single genre name display (backwards compat)

**State:**
```typescript
const route = useRoute()
const genreIds = computed(() => {
  const param = route.query.genres as string
  return param ? param.split(',') : []
})
const orderBy = computed(() => (route.query.order_by as string) || 'popularity')
```

**Fetch Logic:**
```typescript
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
    key: `browse-${route.fullPath}` // Reactive key for refetch
  }
)
```

**Layout:**
```
[Header: Title + Description]
[GenreFilterDropdown]
[Selected Genre Pills Row] (if filters active)
[Sort Dropdown] (optional: order_by)
[Grid: anime cards]
[Load More / Pagination] (future)
```

## API Integration

### Fetch Genres
- Endpoint: `GET /api/genres`
- Response: `{ data: Genre[] }`
- Filter: `count > 1000` (same as index page)
- Slice: Top 20 genres for dropdown

### Fetch Anime
- Endpoint: `GET /api/anime?genres={ids}&order_by={sort}&limit={n}`
- Multi-genre support: `genres=1,2,8` (comma-separated)
- Response: `{ data: Anime[] }`
- Rate limit handling: 429 → retry after 2s (existing pattern)

## URL Strategy

### Multi-Genre Filtering
- Query param: `genres=1,2,8` (comma-separated IDs)
- Backwards compatible with single `genre={id}` from existing "See All" links
- When navigating from index "See All":
  - Single genre: `/browse?genre=1&genre_name=Action`
  - Browse page normalizes to: `/browse?genres=1`

### Preserving Other Params
- Always maintain `order_by` when toggling genres
- Drop `genre_name` when multiple genres selected (no single name to display)

### Example URLs
```
/browse?genres=1,2,8&order_by=popularity
/browse?genres=1
/browse?order_by=score
/browse
```

## UI/UX Details

### Filter Button
- Icon: `i-solar-filter-linear` or `i-solar-settings-linear`
- Text: "Genres" or "Filter"
- Badge: Small circle with count (e.g., "3") when filters active
- Styling: `variant="ghost"`, `color="neutral"`

### Dropdown
- Width: `w-80` (320px)
- Max height: `max-h-96` with scroll
- Position: Absolute, below button, right-aligned
- Backdrop: Click outside to close
- Padding: `p-4`
- Border: `border border-muted`, `rounded-xl`
- Background: `bg-surface`, `backdrop-blur`

### Genre Chips (Inside Dropdown)
- Layout: Flexbox wrap
- Styling when unselected: `border border-muted`, `bg-elevated`
- Styling when selected: `border border-primary`, `bg-primary/10`, checkmark icon
- Hover: Scale 102%, shadow increase
- Transition: `transition-all duration-150`

### Selected Pills (Outside Dropdown)
- Displayed below filter button when filters active
- Dismissible (X icon on right)
- Match styling: `bg-primary/10`, `border border-primary`
- Genre name displayed (fetch genre list to map ID → name)

### Clear All Button
- Only visible when filters active
- Text: "Clear all"
- Styling: `variant="ghost"`, `color="neutral"`, `size="sm"`
- Position: Bottom of dropdown

### Loading States
- Genre dropdown: Skeleton chips or spinner
- Anime grid: Existing spinner pattern
- Empty state: "No anime found" with ghost icon (already exists)

### Mobile Considerations
- Dropdown full-width on small screens
- Genre chips wrap naturally
- Filter button remains accessible in header
- Touch-friendly hit areas (min 44px)

## Error Handling

### API Failures
- Genres fetch fails: Show error toast, hide dropdown
- Anime fetch fails: Show error state in grid
- Rate limiting (429): Retry after 2s (existing pattern)

### Invalid Query Params
- Non-numeric genre IDs: Filter out, proceed with valid ones
- Empty `genres` param: Treat as no filter
- Missing genres: Gracefully handle, show all anime

## Testing Considerations

- Multi-genre filtering combines correctly (AND logic)
- URL updates trigger refetch
- Back/forward navigation preserves filter state
- "See All" links from index work correctly
- Dropdown closes on outside click
- Selected pills remove genres correctly
- Clear all resets to unfiltered state

## Future Enhancements (Out of Scope)

- Pagination / infinite scroll
- Additional filters (year, rating, type)
- Sort by multiple fields
- Save filter presets
- Filter by tags/themes beyond genres
- Advanced search with text input

## Implementation Notes

- Use existing `Genre` interface from GenreRow component
- Reuse anime card styling from existing grid
- Match animation patterns from other components
- Maintain consistent spacing/sizing with current design system
- Use Nuxt UI components (`UButton`, `UIcon`, `UBadge`)
- Follow existing color token patterns (`text-highlighted`, `text-toned`, `bg-elevated`)

## Acceptance Criteria

✓ Genre filter dropdown button visible in browse page header  
✓ Clicking button animates dropdown open/closed  
✓ Genres displayed as clickable chips in dropdown  
✓ Clicking genre chip toggles selection (instant filtering)  
✓ Selected genres shown as pills below button  
✓ URL updates with `genres` query param  
✓ Grid refetches and displays filtered results  
✓ "Clear all" removes all genre filters  
✓ "See All" navigation from index page works  
✓ Backwards compatible with existing single-genre links  
✓ Loading and empty states display correctly  
✓ Mobile-responsive layout
