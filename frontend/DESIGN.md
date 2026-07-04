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
├── assets/css/main.css              # Tailwind + Nuxt UI + @theme brand
├── layouts/
│   └── default.vue                  # AppSidebar + AppHeader shell
├── pages/
│   └── index.vue                    # Composes Hero + Recommended + Genre
└── components/
    ├── AppSidebar.vue
    ├── AppHeader.vue
    ├── HeroSection.vue
    ├── RecommendedSection.vue
    ├── GenreSection.vue
    └── GenreRow.vue
```
