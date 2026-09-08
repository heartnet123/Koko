# KoKo — Design System

## Brand

| Token | Value | Usage |
|---|---|---|
| Primary | `#16B681` | CTAs, active states, progress (desaturated emerald, ~76% sat) |
| Background | `#F4F6FB` light / `#090B10` dark | App shell, sidebar |
| Text | `#0F172A` light / `#F8FAFC` dark | Headings |

Mapped to Nuxt UI `primary` color via `@theme` in `main.css`.

**Palette rules (design-taste-frontend-v1):**
- One accent only. No purple/indigo/cyan mixes, no gradient text on headlines.
- No neon/outer glow shadows (`0_0_*`). Elevation = `--shadow-diffuse*` (wide, faint, tinted).
- No pure `#000`. Overlays use off-black `--ui-overlay: #090B10`.
- No raw Tailwind palette classes (`text-gray-*`, `text-neutral-*`, `bg-red-500/…`). Semantic tokens only.

## Typography

| Font | Usage |
|---|---|
| Geist (via `@nuxt/fonts`) | All UI text — `--font-sans` |
| Geist Mono | DATA ONLY: scores, counts, indices, meta eyebrows |

- Base weight 400; hierarchy via weight + color, not scale (no oversized H1s).
- Headlines: `tracking-tight`, `text-4xl md:text-5xl` max on landing.

## Layout

```
┌─────────────────────────────────────────────────┐
│  AppSidebar (w-64, sticky, min-h-[100dvh])      │
│   │  main                                       │
│   │  AppHeader (h-20)                           │
│   │  ─────────────                              │
│   │  <slot /> (scrollable)                      │
└─────────────────────────────────────────────────┘
Mobile: floating glass bottom nav + slideover drawer.
```

- `layouts/default.vue` — shell with sidebar + header
- `layouts/landing.vue` — guest marketing layout
- `app.vue` — `UApp > NuxtLayout > NuxtPage transition="page"`

## Color semantics (Nuxt UI)

| Class | Meaning |
|---|---|
| `text-highlighted` | Primary content (replaces `text-gray-900`) |
| `text-default` | Body text |
| `text-toned` | Secondary / muted (replaces `text-gray-500`) |
| `bg-elevated` | Card surfaces (replaces `bg-gray-50`) |
| `bg-default` | White-ish surfaces |
| `border-muted` | Subtle borders (replaces `border-gray-100`) |
| `--ui-text-on-image(-muted)` | Text over hero artwork scrims |
| `--ui-error` / `--ui-success` / `--ui-warn` | Semantic status |
| `--rank-gold/silver/bronze` | Trending podium + score stars |

Never use raw Tailwind palette colors (`text-gray-*`, `bg-gray-*`).

## Shadows

| Token | Use |
|---|---|
| `--shadow-diffuse` | Resting elevation (wide, faint, tinted to surface hue) |
| `--shadow-diffuse-lg` | Large panels, heroes |
| `--shadow-diffuse-accent` | Hover elevation on interactive cards/CTAs |

Glass surfaces carry an inner highlight (`inset 0 1px 1px` refraction) — never stack `shadow-2xl` on glass.

## Full-height rule

Never `h-screen`. Always `min-h-[100dvh]` (mobile URL-bar safe).

## Icons

Collection: `solar` (via `@nuxt/icon` bundled in `@nuxt/ui`).
Format: `i-solar-{name}`. Star icons use `text-[var(--rank-gold)]`.

## Motion (MOTION_INTENSITY 6)

- CSS-first: Ken Burns hero, shimmer skeletons, `animate-breathe` status dots, staggered rail load-ins, `animate-fade-in-up`.
- `v-magnetic` directive (`app/directives/magnetic.ts`): CTAs pull toward cursor via rAF + direct DOM — no reactive state, no re-renders.
- Page transitions: `transition="page"` on `NuxtPage`, 0.3s fade/slide.
- All motion honors `prefers-reduced-motion` (main.css media query).

## Components

### AppSidebar
- Fixed 64px wide, sticky, `min-h-[100dvh]`
- Nav links via `NuxtLink` + `active-class` for active state

### AppHeader
- Sticky, `h-20`, frosted glass (`bg-[var(--ui-bg)]/80 backdrop-blur-md`)
- Search input (native `<input>` with UIcon prefix), weight 400

### HeroCarousel
- Full-bleed image, Ken Burns, 8s autoplay with pause-override, dot indicators
- Play + Watchlist `UButton` calls-to-action

### AnimeRail
- Per-genre lazy horizontal scroller (intersection observer, 200px rootMargin)
- Rate-limit aware: retries after 2s on HTTP 429

### AnimeCard
- Portrait `aspect-[2/3]` card, hover = accent diffusion shadow + border tint

## API

All data from [Jikan v4](https://api.jikan.moe/v4) (unofficial MyAnimeList REST).
No auth required. Rate limit: ~3 req/s.

## File tree

```
app/
├── app.vue                          # UApp shell + page transitions
├── assets/css/main.css              # Tailwind + Nuxt UI + @theme brand + tokens
├── directives/magnetic.ts           # v-magnetic CTA micro-physics
├── layouts/
│   ├── default.vue                  # Sidebar + header shell (+ mobile bottom nav)
│   └── landing.vue                  # Guest marketing shell
├── pages/
│   ├── index.vue                    # Landing (guest) + dashboard (authed)
│   ├── browse.vue / trending.vue / watchlist.vue
│   ├── login.vue / profile.vue / settings.vue
│   └── movie/[id].vue               # Detail page
├── composables/
│   ├── useAuth.ts                   # Auth + watchlist management
│   ├── useJikan.ts                  # Fetch wrapper for Jikan proxy with retry logic
│   └── useGenreQuery.ts             # Genre routing query parser
└── components/
    ├── AppSidebar.vue / AppHeader.vue / AnimeCard.vue / AnimeRail.vue
    ├── HeroCarousel.vue / GenreFilterDropdown.vue
    └── skeletons/ (HeroSkeleton, RailSkeleton)
```

---

## Redesign Log (design-taste-frontend-v1 pass)

### Decision Log
1-16: prior entries — see git history (`DESIGN.md` pre-2026 entries).
17. **Fonts**: Cascadia Mono (everywhere) → Geist + Geist Mono. Mono reserved for data/labels; UI sans. Dropped CDN `<link>`, served via `@nuxt/fonts`.
18. **Accent desaturation**: `#00DC82` (100% sat) → `#16B681` (~76%). Full ramp in `@theme`.
19. **No-neon policy**: all `shadow-[0_0_*rgba(0,220,130,…)]` outer glows removed; replaced with `--shadow-diffuse*` + inner glass highlight.
20. **Off-black overlays**: `bg-black/*` → `bg-[var(--ui-overlay)]/*` (`#090B10`); no pure black anywhere.
21. **Gradient headline** (indigo/cyan mix) → single-hue accent span with weight contrast; H1 scale capped at `md:text-5xl`.
22. **Card-wall removal**: "Why KoKo" 2×2 equal cards → `divide-y` border rows with mono indices (no boxes).
23. **Hero asymmetry**: landing hero 6/6 → 5/7 split with vertically offset visual (`lg:mt-12`).
24. **dvh hardening**: `h-screen` / `min-h-screen` / `calc(100vh-*)` → `100dvh` variants.
25. **Magnetic CTAs**: `v-magnetic` directive (rAF, direct DOM, no reactive state) on hero / Get Started / Join KoKo.
26. **Page transitions**: `NuxtPage transition="page"` (0.3s fade/slide) with reduced-motion override.
27. **Breathing status dot**: `animate-pulse` → `animate-breathe` (spring-like scale/opacity loop).
28. **Form fields**: `font-mono` stripped from inputs/labels; weight 400; mono kept only for data.
29. **Semantic tokens**: `--ui-error/-success/-warn`, `--rank-{gold,silver,bronze}` replace raw `red-/amber-/yellow-*` classes.

### Manual Smoke Test Steps (post-redesign)
1. Load home page: shimmer skeleton, then hero first slide starts Ken Burns.
2. Hero rotation: slide 2 after 8s; clicking indicator pauses rotation 10s.
3. Lazy rails: network log shows per-genre queries only when scroller enters viewport.
4. Magnetic CTAs: move cursor over "Get Started Free" — button translates toward cursor, springs back on leave.
5. Page transitions: navigate Home ↔ Browse — fade/slide, no flash.
6. Watchlist redirect: logged out, click bookmark → `/login`.
7. Card click: poster → `/movie/[id]`; detail hero renders with off-black scrims.
8. Fonts: inspect computed font-family — Geist loaded locally (no CDN request).
9. Shadows: hover cards — diffusion tint, no green glow.
10. Reduced motion: enable OS reduce — Ken Burns, shimmer, breathe, page transitions freeze.