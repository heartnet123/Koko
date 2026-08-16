<script setup lang="ts">
import { computed, watchEffect } from 'vue'
import { useJikan } from '~/composables/useJikan'
import type { JikanAnime } from '~/types/anime'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Koko — Modern Anime Library & Discovery Engine',
  ogTitle: 'Koko — Modern Anime Library & Discovery Engine',
  description: 'Explore, track, and catalog anime in an ultra-fast, cinematic glassmorphic interface.',
  ogDescription: 'Explore, track, and catalog anime in an ultra-fast, cinematic glassmorphic interface.',
})

const auth = useAuth()

// Dynamic page layout based on authentication state
watchEffect(() => {
  if (import.meta.client) {
    setPageLayout(auth.isAuthenticated.value ? 'default' : 'landing')
  }
})

// Fetch current season for Hero (max 5 slides)
const { data: heroResponse, pending: heroPending } = useJikan<{ data: JikanAnime[] }>('/seasons/now')

const heroSlides = computed(() => {
  return (heroResponse.value?.data ?? []).slice(0, 5)
})

// Fetch genres for rails configuration
const { data: genresResponse, error: genresError } = useJikan<{ data: any[] }>('/genres')

const genres = computed(() => {
  return (genresResponse.value?.data ?? [])
    .filter(g => g.count > 1000)
    .slice(0, 10)
})

const whyItems = [
  { icon: 'i-solar-star-bold', title: 'Cinematic Over Metadata', desc: 'Instead of overwhelming tables, we prioritize high-resolution backdrops, clean trailer playback, and fluid artwork transitions.' },
  { icon: 'i-solar-bolt-bold', title: 'Zero Page Jumps', desc: 'Search queries and genre filters update smoothly in-place with instant client transitions.' },
  { icon: 'i-solar-shield-bold', title: 'One-Click Tracking', desc: 'Bookmark and unbookmark directly from any card or hero banner without tedious multi-field form dialogs.' },
  { icon: 'i-solar-widget-6-bold', title: 'Clean & Ad-Free', desc: 'No intrusive third-party popups, tracking ads, or spam. Just clean, open-source performance.' },
]
</script>

<template>
  <div>
    <!-- Authenticated View: Dashboard -->
    <div v-if="auth.isAuthenticated.value" class="max-w-7xl mx-auto px-4 md:px-8 w-full flex flex-col gap-10 mt-4 pb-16">
      <!-- Hero Carousel -->
      <div>
        <HeroSkeleton v-if="heroPending" />
        <HeroCarousel v-else-if="heroSlides.length" :slides="heroSlides" />
        <div v-else class="min-h-[40dvh] flex flex-col items-center justify-center glass-surface rounded-3xl border border-[var(--glass-border)] gap-2">
          <p class="text-[var(--ui-text-toned)] text-xs font-mono">Couldn't load current season.</p>
        </div>
      </div>

      <!-- Recommended Rail -->
      <AnimeRail 
        title="Recommended For You" 
        fetchUrl="/anime?order_by=score&sort=desc&limit=12" 
      />

      <!-- Genre Rails with CSS stagger delay -->
      <div v-if="genresError" class="text-center py-6 text-xs text-[var(--ui-text-toned)] border border-dashed border-[var(--ui-border-muted)] rounded-2xl glass-pill">
        Genre rails unavailable
      </div>
      
      <div v-else class="flex flex-col gap-10">
        <div 
          v-for="(genre, idx) in genres" 
          :key="genre.mal_id"
          class="staggered-rail"
          :style="{ '--i': idx }"
        >
          <AnimeRail
            :title="genre.name"
            :count="genre.count"
            :fetchUrl="`/anime?genres=${genre.mal_id}&order_by=score&sort=desc&limit=12`"
          />
        </div>
      </div>
    </div>

    <!-- Public Landing Page for Guest Visitors -->
    <div v-else class="w-full flex flex-col gap-16 md:gap-24 animate-fade-in-up">
      <!-- Hero Section -->
      <section class="max-w-7xl mx-auto px-6 w-full min-h-[80dvh] md:min-h-[85dvh] flex items-center pt-8 md:pt-16 pb-12">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 w-full items-center">
          <!-- Copy & CTAs -->
          <div class="lg:col-span-5 flex flex-col gap-6 md:gap-8 justify-center text-left">
            <div class="inline-flex items-center gap-2 glass-pill px-3 py-1 rounded-full w-fit">
              <span class="w-2 h-2 rounded-full bg-primary-400 animate-pulse" />
              <span class="text-xs font-bold text-primary-400 tracking-wider uppercase font-mono">v4.0 Glass Edition</span>
            </div>

            <h1 class="text-4xl md:text-5xl font-semibold tracking-tight text-[var(--ui-text-highlighted)] leading-[1.08]">
              Anime Tracking,<br>
              <span class="font-bold text-primary-600">
                Cinematic & Fast.
              </span>
            </h1>
            <p class="text-sm md:text-base text-[var(--ui-text-toned)] leading-relaxed max-w-[48ch] font-normal">
              Catalog series, track your watchlist, and explore collections in a stunning glassmorphic interface engineered for anime enthusiasts.
            </p>
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 mt-2">
              <UButton
                to="/login"
                label="Get Started Free"
                size="lg"
                color="primary"
                icon="i-solar-rocket-bold"
                class="rounded-2xl font-bold shadow-[var(--shadow-diffuse-accent)] px-8 py-3.5 text-center justify-center cursor-pointer hover:scale-[1.02] active:scale-95 transition-all"
              />
              <UButton
                to="/browse"
                label="Browse Library"
                size="lg"
                variant="ghost"
                color="neutral"
                icon="i-solar-compass-linear"
                class="glass-pill rounded-2xl font-semibold px-6 py-3.5 text-center justify-center cursor-pointer hover:bg-white/10 hover:scale-[1.02] active:scale-95 transition-all"
              />
            </div>
          </div>

          <!-- Hero Image with Glass Refraction Glow -->
          <div class="lg:col-span-7 relative w-full flex items-center justify-center lg:mt-12">
            <div class="absolute inset-0 bg-primary-500/20 rounded-[40px] blur-3xl -z-10" />
            <div class="w-full relative rounded-3xl overflow-hidden glass-surface border border-[var(--glass-border)]  transform transition-all duration-700 hover:scale-[1.01] group">
              <NuxtImg
                src="/koko_hero.jpg"
                alt="KoKo Cinematic Viewport Preview"
                class="w-full h-full object-cover aspect-[16/10] md:aspect-[16/9] lg:aspect-[4/3] transition-transform duration-700 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/60 via-transparent to-transparent pointer-events-none" />
              <div class="absolute bottom-4 left-4 right-4 glass-chip px-4 py-2.5 rounded-xl flex items-center justify-between text-white text-xs font-mono">
                <span class="flex items-center gap-2">
                  <UIcon name="i-solar-play-circle-bold" class="w-4 h-4 text-primary-400" />
                  Instant Trailer Previews
                </span>
                <span class="text-[var(--ui-text-on-image-muted)] text-[10px]">Zero Latency</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Logo Wall -->
      <section class="py-10 max-w-7xl mx-auto px-6 w-full border-t border-[var(--glass-border-subtle)] flex flex-col md:flex-row items-center justify-between gap-6 opacity-70">
        <span class="text-[11px] font-mono uppercase tracking-wider text-[var(--ui-text-toned)]">Powered by community standards</span>
        <div class="flex flex-wrap items-center gap-8 justify-center select-none dark:invert dark:brightness-100">
          <img src="https://cdn.simpleicons.org/myanimelist/6B7280" alt="MyAnimeList" class="h-5 opacity-80 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/nuxt/6B7280" alt="Nuxt" class="h-5 opacity-80 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/vuedotjs/6B7280" alt="Vue" class="h-5 opacity-80 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/tailwindcss/6B7280" alt="Tailwind CSS" class="h-5 opacity-80 hover:opacity-100 transition-opacity" />
        </div>
      </section>

      <!-- Features Section (Glass Bento Grid) -->
      <section id="features" class="py-20 border-y border-[var(--glass-border-subtle)] relative">
        <div class="max-w-7xl mx-auto px-6 w-full flex flex-col gap-12">
          <div class="max-w-2xl">
            <h2 class="text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)] sm:text-4xl leading-tight">
              Crafted for the Ultimate Viewing Experience
            </h2>
            <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-3 leading-relaxed font-normal">
              Say goodbye to outdated tables and endless clutter. KoKo is built with frosted glass ergonomics, instant proxy caching, and responsive media rails.
            </p>
          </div>

          <!-- Bento Grid with Frosted Glass Panels -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Feature 1: Cinematic UI (col-span-2) -->
            <div class="md:col-span-2 glass-surface-elevated rounded-3xl p-8 flex flex-col justify-between border border-[var(--glass-border)] shadow-[var(--shadow-diffuse-lg)] hover:border-primary-400/50 hover:shadow-[var(--shadow-diffuse-accent)] transition-all duration-300 relative overflow-hidden group min-h-[340px]">
              <div class="absolute -right-10 -bottom-10 w-72 h-44 bg-primary-500/10 rounded-full blur-3xl group-hover:bg-primary-500/20 transition-colors" />
              <div class="flex flex-col gap-3 max-w-md z-10">
                <div class="w-10 h-10 rounded-2xl glass-pill flex items-center justify-center text-primary-400 mb-2 shadow-inner">
                  <UIcon name="i-solar-gallery-wide-linear" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-[var(--ui-text-highlighted)]">Cinematic Poster Interface</h3>
                <p class="text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                  Experience anime with high-resolution poster cards, smooth Ken Burns scaling, and fluid micro-animations mirroring premium streaming platforms.
                </p>
              </div>
              <!-- Poster Mockup inside the bento cell -->
              <div class="mt-6 flex gap-3 overflow-hidden select-none z-10 translate-y-2 group-hover:translate-y-0 transition-transform duration-500">
                <div v-for="i in 3" :key="i" class="w-24 flex-shrink-0 aspect-[3/4] glass-surface rounded-xl border border-[var(--glass-border)] overflow-hidden relative shadow-md">
                  <div class="absolute inset-0 bg-gradient-to-t from-[var(--ui-overlay)]/80 to-transparent z-10" />
                  <div class="absolute bottom-2 left-2 right-2 text-[9px] font-semibold text-white truncate font-mono">Season {{ i }}</div>
                </div>
              </div>
            </div>

            <!-- Feature 2: Personalized Watchlist (col-span-1) -->
            <div class="glass-surface-elevated rounded-3xl p-8 flex flex-col justify-between border border-[var(--glass-border)] shadow-[var(--shadow-diffuse-lg)] hover:border-primary-400/50 hover:shadow-[var(--shadow-diffuse-accent)] transition-all duration-300 relative overflow-hidden group min-h-[340px]">
              <div class="absolute -left-10 -bottom-10 w-52 h-52 bg-primary-500/10 rounded-full blur-3xl group-hover:bg-primary-500/20 transition-colors" />
              <div class="flex flex-col gap-3 z-10">
                <div class="w-10 h-10 rounded-2xl glass-pill flex items-center justify-center text-primary-400 mb-2 shadow-inner">
                  <UIcon name="i-solar-bookmark-bold-duotone" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-[var(--ui-text-highlighted)]">Instant Watchlist</h3>
                <p class="text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                  One-click tracking. Save series, manage your collection, and keep titles organized across all your devices.
                </p>
              </div>
              <!-- Watchlist visual indicator -->
              <div class="mt-6 glass-surface border border-[var(--glass-border)] rounded-2xl p-4 flex flex-col gap-2.5 z-10 transform group-hover:scale-[1.02] transition-transform duration-300 shadow-md">
                <div class="flex items-center justify-between">
                  <span class="text-[11px] font-bold text-[var(--ui-text-highlighted)] font-mono">Currently Watching</span>
                  <span class="text-[9px] bg-primary-500/20 text-primary-300 px-2 py-0.5 rounded-full font-mono font-bold">12 Series</span>
                </div>
                <div class="w-full bg-[var(--ui-overlay)]/20 dark:bg-white/10 h-1.5 rounded-full overflow-hidden">
                  <div class="bg-gradient-to-r from-primary-400 to-primary-500 w-2/3 h-full rounded-full " />
                </div>
              </div>
            </div>

            <!-- Feature 3: Instant Discovery (col-span-3) -->
            <div class="md:col-span-3 glass-surface-elevated rounded-3xl p-8 flex flex-col md:flex-row items-center justify-between gap-8 border border-[var(--glass-border)] shadow-[var(--shadow-diffuse-lg)] hover:border-primary-400/50 hover:shadow-[var(--shadow-diffuse-accent)] transition-all duration-300 relative overflow-hidden group">
              <div class="flex flex-col gap-3 max-w-lg z-10">
                <div class="w-10 h-10 rounded-2xl glass-pill flex items-center justify-center text-primary-400 mb-2 shadow-inner">
                  <UIcon name="i-solar-magnifer-linear" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-[var(--ui-text-highlighted)]">Fast Browsing & Filtering</h3>
                <p class="text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                  Find what you want immediately. Optimized debounce querying with instant genre intersections and sub-second page transitions.
                </p>
              </div>
              <!-- Search Mockup -->
              <div class="w-full md:w-80 glass-surface border border-[var(--glass-border)] rounded-2xl p-4 flex flex-col gap-3 z-10 shrink-0 transform group-hover:translate-x-1 transition-transform duration-300 shadow-md">
                <div class="flex items-center gap-2 glass-pill rounded-full px-3 py-1.5 text-[11px] text-[var(--ui-text-toned)] font-mono">
                  <UIcon name="i-solar-magnifer-linear" class="w-3.5 h-3.5 text-primary-400" />
                  <span>Action, Cyberpunk...</span>
                </div>
                <div class="flex flex-wrap gap-1.5 font-mono">
                  <span class="text-[10px] bg-primary-500/20 text-primary-300 border border-primary-400/40 px-2 py-0.5 rounded-lg font-semibold">Action</span>
                  <span class="text-[10px] bg-primary-500/20 text-primary-300 border border-primary-400/40 px-2 py-0.5 rounded-lg font-semibold">Sci-Fi</span>
                  <span class="text-[10px] glass-pill text-[var(--ui-text-toned)] px-2 py-0.5 rounded-lg">+14 genres</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Why KoKo Section -->
      <section id="why-koko" class="py-20 transition-all duration-200">
        <div class="max-w-7xl mx-auto px-6 w-full grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          <div class="lg:col-span-5 flex flex-col gap-4">
            <h2 class="text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)] sm:text-4xl leading-tight">
              Why use KoKo?
            </h2>
            <p class="text-xs md:text-sm text-[var(--ui-text-toned)] leading-relaxed font-normal">
              Traditional anime portals are cluttered, slow, and feel like complex spreadsheets. KoKo is built for people who appreciate design details, speed, and immersive visuals.
            </p>
          </div>

          <!-- Comparison Rows — divided lines, no card boxes -->
          <div class="lg:col-span-7 divide-y divide-[var(--ui-border-muted)] border-y border-[var(--ui-border-muted)]">
            <div
              v-for="(item, i) in whyItems"
              :key="item.title"
              class="group flex items-start gap-5 py-7 px-1 hover:translate-x-1 transition-transform duration-300"
            >
              <div class="w-11 h-11 flex-shrink-0 rounded-2xl glass-pill flex items-center justify-center text-primary-500">
                <UIcon :name="item.icon" class="w-5 h-5" />
              </div>
              <div class="flex-1 min-w-0">
                <h4 class="text-sm md:text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight">{{ item.title }}</h4>
                <p class="text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal mt-1">{{ item.desc }}</p>
              </div>
              <span class="hidden sm:block text-[11px] font-mono text-[var(--ui-text-toned)]/70 pt-1">{{ `0${i + 1}` }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- FAQ Section (Glass Accordion) -->
      <section id="faq" class="py-20 border-t border-[var(--glass-border-subtle)] relative">
        <div class="max-w-4xl mx-auto px-6 w-full flex flex-col gap-10">
          <div class="max-w-2xl">
            <h2 class="text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)] sm:text-4xl">
              Frequently Asked Questions
            </h2>
            <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-3 font-normal">
              Everything you need to know about the KoKo platform.
            </p>
          </div>

          <!-- Accordion Details with Glass Styling -->
          <div class="flex flex-col gap-4">
            <details class="group glass-surface border border-[var(--glass-border)] rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none shadow-sm transition-all hover:border-primary-400/30">
              <summary class="flex items-center justify-between font-bold text-[var(--ui-text-highlighted)] text-sm md:text-base">
                <span>Is KoKo free to use?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-[var(--ui-text-toned)]">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                Yes! KoKo is completely free. We do not run ads or sell user data. The project is an open source study in designing premium user experiences.
              </p>
            </details>

            <details class="group glass-surface border border-[var(--glass-border)] rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none shadow-sm transition-all hover:border-primary-400/30">
              <summary class="flex items-center justify-between font-bold text-[var(--ui-text-highlighted)] text-sm md:text-base">
                <span>Do I need an external account?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-[var(--ui-text-toned)]">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                No external account is required. You can sign up directly on KoKo to store and sync your watchlist locally across sessions.
              </p>
            </details>

            <details class="group glass-surface border border-[var(--glass-border)] rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none shadow-sm transition-all hover:border-primary-400/30">
              <summary class="flex items-center justify-between font-bold text-[var(--ui-text-highlighted)] text-sm md:text-base">
                <span>Where does anime metadata come from?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-[var(--ui-text-toned)]">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-[var(--ui-text-toned)] leading-relaxed font-normal">
                All data is fetched through our high-speed Go proxy cache interfacing with the open Jikan REST v4 API.
              </p>
            </details>
          </div>
        </div>
      </section>

      <!-- Final CTA Section -->
      <section class="py-24 max-w-7xl mx-auto px-6 w-full">
        <div class="relative glass-surface-elevated border border-[var(--glass-border)] rounded-3xl p-12 md:p-16 flex flex-col items-center text-center gap-6 overflow-hidden ">
          <div class="absolute -right-20 -bottom-20 w-80 h-80 bg-primary-500/15 rounded-full blur-3xl" />
          <div class="absolute -left-20 -top-20 w-80 h-80 bg-cyan-500/15 rounded-full blur-3xl" />
          
          <div class="w-12 h-12 rounded-2xl glass-pill flex items-center justify-center text-primary-400 mb-2 shadow-inner">
            <UIcon name="i-solar-heart-bold" class="w-6 h-6" />
          </div>

          <h2 class="text-3xl md:text-4xl font-bold tracking-tight text-[var(--ui-text-highlighted)] leading-tight max-w-xl">
            Ready to Experience Anime Differently?
          </h2>
          <p class="text-xs md:text-sm text-[var(--ui-text-toned)] max-w-md font-normal">
            Join KoKo today. Organize your collection, discover new seasons, and enjoy a clean cinematic interface.
          </p>
          <div class="flex flex-col sm:flex-row items-center gap-4 mt-2">
            <UButton
              to="/login"
              label="Join KoKo"
              size="lg"
              color="primary"
              class="rounded-2xl font-bold shadow-[var(--shadow-diffuse-accent)] px-8 py-3.5 cursor-pointer hover:scale-[1.02] active:scale-95 transition-all"
            />
            <UButton
              to="/browse"
              label="Browse Library First"
              size="lg"
              variant="ghost"
              color="neutral"
              class="glass-pill rounded-2xl font-semibold px-6 py-3.5 cursor-pointer hover:bg-white/10 hover:scale-[1.02] active:scale-95 transition-all"
            />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.staggered-rail {
  animation: fadeIn 0.6s ease-out both;
  animation-delay: calc(var(--i) * 80ms);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .staggered-rail {
    animation: none !important;
  }
}
</style>