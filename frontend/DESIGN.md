# KoKo — Design System

## Brand

| Token | Value | Usage |
|---|---|---|
| Primary | `#635BFF` | CTAs, active states, progress |
| Background | `#FCFCFD` | App shell, sidebar |
| Text | `#1A1A2E` | Headings |

Mapped to Nuxt UI `primary` color via `@theme` in `main.css`.

---

## Layout

```
┌─────────────────────────────────────────────────┐
│  AppSidebar (w-64, sticky)  │  main             │
│                              │  AppHeader (h-20) │
│  Brand logo                  │  ───────────────  │
│  Nav items                   │  <slot />         │
│  User items                  │  (scrollable)     │
└─────────────────────────────────────────────────┘
```

- `layouts/default.vue` — shell with sidebar + header
- `app.vue` — `UApp > NuxtLayout > NuxtPage`

---

## Color semantics (Nuxt UI)

| Class | Meaning |
|---|---|
| `text-highlighted` | Primary content (replaces `text-gray-900`) |
| `text-default` | Body text |
| `text-toned` | Secondary / muted (replaces `text-gray-500`) |
| `bg-elevated` | Card surfaces (replaces `bg-gray-50`) |
| `bg-default` | White-ish surfaces |
| `border-muted` | Subtle borders (replaces `border-gray-100`) |

Never use raw Tailwind palette colors (`text-gray-*`, `bg-gray-*`).

---

## Components

### AppSidebar
- Fixed 64px wide, sticky, `h-screen`
- Nav links via `NuxtLink` + `active-class` for active state

### AppHeader
- Sticky, `h-20`, frosted glass (`bg-[#FCFCFD]/80 backdrop-blur-md`)
- Search input (native `<input>` with UIcon prefix)
- Bell button + avatar

### HeroSection
- Full-bleed image with gradient overlay
- Play + Watchlist `UButton` calls-to-action
- Slides array — swap `src` to drive featured content

### RecommendedSection
- Fetches `api.jikan.moe/v4/top/anime`
- 5-col grid, `aspect-[3/4]` cards with play overlay

### GenreSection
- Fetches genre list, filters `count > 1000`, top 10
- Renders `GenreRow` per genre (staggered with index)

### GenreRow
- Per-genre 5-col grid
- Rate-limit aware: retries after 2s on HTTP 429
- Color/icon rotates through 8 preset styles via `index % 8`

---

## Icons

Collection: `solar` (via `@nuxt/icon` bundled in `@nuxt/ui`).
Format: `i-solar-{name}`.

---

## API

All data from [Jikan v4](https://api.jikan.moe/v4) (unofficial MyAnimeList REST).
No auth required. Rate limit: ~3 req/s.

---

## File tree

```
app/
├── app.vue                          # UApp shell
├── assets/css/main.css              # Tailwind + Nuxt UI + @theme brand + custom animations
├── layouts/
│   └── default.vue                  # AppSidebar + AppHeader shell
├── pages/
│   └── index.vue                    # Composes HeroCarousel + AnimeRail list
├── composables/
│   ├── useAuth.ts                   # Auth + watchlist management
│   ├── useJikan.ts                  # Fetch wrapper for Jikan proxy with retry logic
│   └── useGenreQuery.ts             # Genre routing query parser
└── components/
    ├── AppSidebar.vue
    ├── AppHeader.vue
    ├── AnimeCard.vue                # Portrait aspect-[2/3] card
    ├── AnimeRail.vue                # Lazy horizontal anime scroller
    ├── HeroCarousel.vue             # Autoplay Ken Burns carousel
    └── skeletons/
        ├── HeroSkeleton.vue         # Shimmer skeleton for Hero
        └── RailSkeleton.vue         # Shimmer skeleton for AnimeRail
```

---

## Minimal Redesign

### Decision Log
1. **Skeletons Isolation**: Skeletons are modularized inside `skeletons/` to avoid component clutter.
2. **Hero Datasource**: `/seasons/now` endpoint is mapped in backend to load current airing titles for the hero carousel.
3. **useJikan Composable**: Crafted a clean utility using native Vue refs and `$fetch` supporting automatic retry logic.
4. **Ken Burns Motion**: Defined pure CSS scale/translate keyframes for smooth hardware-accelerated movements.
5. **Shimmer Gradients**: Created a linear-gradient keyframe animation on background positioning for skeletons.
6. **Horizontal Scrolling**: Implemented touch-smooth horizontal scroll container utilizing Tailwind's `snap-x` and browser scrollbar hiding.
7. **Lazy Fetching via intersection**: Utilized VueUse's `useIntersectionObserver` to trigger network request only when scroller enters viewport (rootMargin: '200px').
8. **Item De-duplication**: Embedded set-based filtering on `mal_id`s within rails to clean up API results.
9. **Retry Mechanisms**: Provided direct click callbacks on error states to let users reload individual failed categories.
10. **Autoplay Timing**: Programmed 8-second rotation cycles for Hero slide intervals.
11. **Autoplay Pause Override**: Implemented timeout overrides that pause auto-rotation for 10 seconds on dot interactions.
12. **Hero Transitions**: Wrapped slides in Vue `<Transition name="hero-fade">` to smoothly blend layouts.
13. **Watchlist State Matching**: Wired bookmarks directly to the `inWatchlist(id)` getter inside `useAuth()`.
14. **Watchlist Flow Protection**: Programmed automatic redirection to `/login` inside card bookmarks when unauthenticated.
15. **Staggered Page Load**: Applied stagger delays to the categories list to introduce rails with visual flow.
16. **Prefers Reduced Motion Support**: Applied a CSS media query that overrides both `@keyframes` animations when OS preference is active.

### Architecture

```mermaid
graph TD
  Index[index.vue Orchestrator] -->|Loads immediate| HeroSkel[HeroSkeleton]
  Index -->|Loads immediate| Hero[HeroCarousel]
  Index -->|Loads staggered| Rails[AnimeRails]
  
  Hero -->|Watchlist state| Auth[useAuth Composable]
  Hero -->|Renders slides| Card[AnimeCard]
  Rails -->|Viewport Intersect| useJikan[useJikan Fetcher]
  Rails -->|Renders items| Card
  Card -->|Route Link| Details[/movie/:id Details Page]
```

### Component Responsibilities
- `HeroCarousel.vue`: Handles image scale animation, automatic 8s interval timer, dot buttons, and watchlist sync.
- `AnimeRail.vue`: Orchestrates intersection observer, lazy loads details, displays slider chevrons, and lists card nodes.
- `AnimeCard.vue`: Renders score chip, title info, type/year metadata, and provides the redirect route.

### Manual Smoke Test Steps
1. Load home page: Verify shimmer skeleton placeholder is displayed.
2. Hero slide loading: Ensure the first slide loads and begins Ken Burns panning.
3. Rotation verification: Confirm the carousel rotates to slide 2 after 8s.
4. Pause check: Click slide indicator 4, confirm transition, and ensure the rotation halts for 10s.
5. Lazy intersection test: Scroll down. Confirm in network logs that rails query endpoints only when visible.
6. Watchlist redirect check: Log out. Click "Watchlist" on Hero. Confirm it redirects to `/login`.
7. Card click verification: Click any anime poster to ensure redirection to `/movie/[id]`.
8. Reduced motion: Enable reduced-motion in system settings. Confirm Ken Burns and skeleton shimmer animations freeze.
9. Error retry test: Trigger a mock 429 rate limit. Verify "Retry" button renders and recovers successfully on click.
