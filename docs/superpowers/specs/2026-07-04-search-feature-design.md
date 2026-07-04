# Search Feature Design

## Goal
Implement a global search feature that allows users to instantly search for anime from any page.

## Requirements & Constraints
- The search should trigger on every keystroke with a debounce (e.g., 500ms) to prevent API spam.
- If the user types in the search bar on any page, they should immediately be routed to the `/browse` page (if not already there).
- The search input must remain in sync with the `q` URL query parameter.
- The `/browse` page must read the `q` query parameter and pass it to the backend API (`/api/anime?q=...`).

## Architecture & Components

### 1. `AppHeader.vue` (Global Search Input)
- **State**: Keep a local reactive `search` string, initialized from `route.query.q`.
- **Behavior**:
  - Watch the `search` input using a debounced watcher (e.g., via `vueuse` `watchDebounced` or a custom timeout).
  - When the debounced value changes, perform a programmatic navigation `navigateTo({ path: '/browse', query: { ...route.query, q: newSearchValue } })`. (If the value is empty, remove the `q` parameter from the query).
  - Also watch `route.query.q` and update the local `search` string so that navigating back/forward correctly updates the input field.

### 2. `browse.vue` (Search Results Page)
- **State**: Extract `searchQuery = computed(() => route.query.q as string)`.
- **Data Fetching**: Update the `useFetch` URL computed property to include `&q=${searchQuery.value}` if a search query is present.
- **UI Updates**: Update `headerText` and `subtitleText` to mention the search query if it exists (e.g., "Search Results for 'X'" or similar).

## Data Flow
1. User types in `AppHeader.vue`.
2. After 500ms debounce, `navigateTo` is called updating the route to `/browse?q=XYZ`.
3. `browse.vue` detects the route change.
4. `searchQuery` computed property updates.
5. The computed API URL inside `useFetch` updates.
6. `useFetch` automatically triggers a new request to `http://localhost:8080/api/anime?q=XYZ&...`.
7. The anime list updates with the search results.
8. The backend `main.go` already supports the `q` parameter and forwards it to the Jikan API.

## Error Handling
- If the search query returns no results, the existing empty state in `browse.vue` will naturally handle it ("No anime found.").
- The existing debounce will mitigate rate-limiting (HTTP 429) from the backend/Jikan API.

## Testing
- **Manual Verification**:
  - Type in the search box from the home page; verify it routes to `/browse` and loads results.
  - Delete characters; verify it re-fetches and eventually returns to default popular anime when empty.
  - Navigate back in the browser history; verify the search box text updates to match the previous query.
