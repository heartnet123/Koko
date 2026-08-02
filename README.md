<div align="center">

<img src="./frontend/public/koko_hero.jpg" alt="KoKo Banner" width="720" style="border-radius: 12px; margin-bottom: 16px;" />

# KoKo

**A high-performance anime library, discovery engine, and personal media tracker.**

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Nuxt](https://img.shields.io/badge/Nuxt-v4.4-00DC82?style=flat-square&logo=nuxt.js&logoColor=white)](https://nuxt.com/)
[![Vue](https://img.shields.io/badge/Vue.js-v3.5-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.10-008080?style=flat-square&logo=go&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://supabase.com/)

[Overview](#overview) • [Key Features](#key-features) • [Tech Stack](#tech-stack) • [Getting Started](#getting-started) • [API Reference](#api-reference) • [Architecture](#architecture)

</div>

---

## Overview

**KoKo** is a full-stack media platform engineered for anime discovery, detailed title inspection, and user watchlist management. It pairs a sleek, hardware-accelerated **Nuxt 4** single-page application with a concurrent **Go + Gin** backend microservice. 

To overcome third-party rate limits and maintain sub-millisecond page transitions, KoKo features an in-memory caching proxy that interfaces with [Jikan v4](https://api.jikan.moe/v4) (the open-source MyAnimeList API), backed by a PostgreSQL database for user identity and persistent watchlists.

> [!NOTE]
> KoKo operates seamlessly out of the box. Public anime discovery endpoints rely on cached proxy routes, requiring no external API keys.

---

## Key Features

- **Dynamic Hero Carousel**: Featured seasonal titles with interactive slides, Ken Burns motion effects, and automatic rotation intervals.
- **Lazy-Loaded Anime Rails**: Horizontal scrollers backed by VueUse `useIntersectionObserver` to defer network fetching until rails scroll into view.
- **Catalog & Search Filter**: Browse anime titles filtered by genre, search query, popularity, score, and pagination.
- **Detailed Title Inspector**: Comprehensive detail pages featuring trailer embeds, rating breakdowns, genre tags, and recommendation rails.
- **Secure Authentication**: JWT-based user authentication using `bcrypt` password hashing, supported via HTTP-Only cookies and Bearer tokens.
- **Personal Watchlist**: Custom collection management allowing users to bookmark titles and sync items across sessions.
- **Thread-Safe Proxy Cache**: In-memory `sync.RWMutex` caching layer in Go, reducing downstream latency and protecting against upstream HTTP 429 rate limits.
- **Responsive & Accessible UI**: Modern design system built with Nuxt UI, custom semantic tokens, light/dark accessibility, and system reduced-motion detection.

---

## Tech Stack

### Frontend
- **Framework**: [Nuxt 4](https://nuxt.com/) (Vue 3.5+, TypeScript)
- **UI & Components**: [Nuxt UI v4](https://ui.nuxt.com/), Solar Linear Icons (`@nuxt/icon`)
- **Styling**: Tailwind CSS v4, custom theme variables, hardware-accelerated CSS keyframe animations
- **Utilities**: `@vueuse/core` for intersection observers and reactive state management

### Backend
- **Language**: [Go 1.25](https://go.dev/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: PostgreSQL (Supabase compatible) via `github.com/lib/pq`
- **Security**: JWT (`golang.org/x/crypto/bcrypt`, `github.com/golang-jwt/jwt/v5`)
- **Cache**: Concurrent in-memory key-value cache with TTL expiration

---

## Architecture

```
                                 ┌─────────────────────────┐
                                 │   Jikan REST API v4     │
                                 └────────────▲────────────┘
                                              │ HTTP (Cached 5m)
┌─────────────────────────┐      ┌────────────┴────────────┐      ┌─────────────────────────┐
│   Nuxt 4 Web Client     │<---->│     Go / Gin Server     │<---->│  PostgreSQL / Supabase  │
│  (Port 3000 / SPA)      │ JSON │    (Port 8080 / API)    │ SQL  │       (Users/List)      │
└─────────────────────────┘      └─────────────────────────┘      └─────────────────────────┘
```

> [!TIP]
> The backend server automatically runs schema migrations on startup, ensuring PostgreSQL tables (`users` and `watchlist`) are provisioned automatically.

---

## Project Structure

```text
koko/
├── backend/                  # Go backend microservice
│   ├── main.go               # Router setup, auth handlers, proxy & cache engine
│   ├── db.go                 # Database connection & schema initializations
│   ├── go.mod                # Go module specification
│   └── .env                  # Backend environment variables
├── frontend/                 # Nuxt 4 web application
│   ├── app/
│   │   ├── components/       # HeroCarousel, AnimeRail, AnimeCard, AppSidebar, etc.
│   │   ├── composables/      # useAuth, useJikan, useGenreQuery
│   │   ├── pages/            # Nuxt file-based routes (index, browse, movie, profile)
│   │   └── assets/css/       # Custom design system tokens and animation keyframes
│   ├── public/               # Static assets & brand imagery
│   ├── nuxt.config.ts        # Nuxt application config
│   └── DESIGN.md             # Frontend design tokens and decision log
└── docs/                     # Documentation & database scripts
    ├── supabase_schema.sql   # Standalone PostgreSQL SQL schema
    └── design/               # Architectural and UI planning specs
```

---

## Getting Started

### Prerequisites

Ensure you have the following installed on your local machine:
- **Go**: `v1.25` or higher
- **Node.js**: `v20.x` or higher (or **Bun** `v1.x`)
- **PostgreSQL**: Local PostgreSQL instance or a [Supabase](https://supabase.com/) project

### 1. Database Setup

Create a PostgreSQL database and set up your connection string. You can manually run `./docs/supabase_schema.sql` or allow the backend auto-migration to initialize tables on first startup.

### 2. Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Create a `.env` file in `backend/`:
   ```env
   DATABASE_URL="postgres://user:password@localhost:5432/koko?sslmode=disable"
   JWT_SECRET="your-super-secret-jwt-key"
   ```

3. Run the Go server:
   ```bash
   go run main.go db.go
   ```
   The backend API will start on `http://localhost:8080`.

> [!IMPORTANT]
> To execute internal health checks and cache self-tests, run:
> ```bash
> RUN_SELF_TEST=true go run main.go db.go
> ```

### 3. Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   bun install
   # or
   npm install
   ```

3. Start the Nuxt development server:
   ```bash
   bun run dev
   # or
   npm run dev
   ```

4. Open `http://localhost:3000` in your web browser.

---

## Environment Variables

| Variable | Location | Description | Default / Example |
|---|---|---|---|
| `DATABASE_URL` | `backend/.env` | PostgreSQL connection string | `postgres://user:pass@host:5432/db` |
| `JWT_SECRET` | `backend/.env` | Secret key used to sign JWT tokens | `your-secure-jwt-secret-key` |
| `RUN_SELF_TEST` | Runtime flag | Triggers automated backend self-tests on boot | `true` / `false` |

---

## API Reference

### Authentication & User Endpoints

| Method | Endpoint | Protection | Description |
|---|---|---|---|
| `POST` | `/api/users` | Public | Register a new user account |
| `POST` | `/api/auth/session` | Public | Authenticate user & receive JWT cookie |
| `DELETE` | `/api/auth/session` | Public | Log out user & clear auth session |
| `GET` | `/api/users/me` | Protected | Fetch current user profile |
| `PATCH` | `/api/users/me` | Protected | Update display name, bio, avatar, or password |
| `GET` | `/api/users/me/watchlist` | Protected | Retrieve saved watchlist items |
| `POST` | `/api/users/me/watchlist` | Protected | Add title to watchlist |
| `DELETE` | `/api/users/me/watchlist/:anime_id` | Protected | Remove title from watchlist |

### Public Anime Endpoints (Cached Proxy)

| Method | Endpoint | Query Parameters | Description |
|---|---|---|---|
| `GET` | `/api/recommendations` | None | Fetch top recommended anime titles |
| `GET` | `/api/genres` | None | Fetch available anime genres |
| `GET` | `/api/seasons/now` | None | Fetch currently airing seasonal anime |
| `GET` | `/api/anime` | `q`, `genres`, `page`, `limit`, `order_by`, `sort` | Search & filter anime catalog |
| `GET` | `/api/anime/:id` | None | Fetch full metadata for a specific anime |
