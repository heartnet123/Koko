<script setup lang="ts">
import { computed, watchEffect } from 'vue'
import { useJikan } from '~/composables/useJikan'
import type { JikanAnime } from '~/types/anime'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Koko - Anime Library',
  ogTitle: 'Koko - Anime Library',
  description: 'Explore the ultimate collection of anime movies and series on Koko.',
  ogDescription: 'Explore the ultimate collection of anime movies and series on Koko.',
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
</script>

<template>
  <div>
    <!-- Authenticated View: Dashboard -->
    <div v-if="auth.isAuthenticated.value" class="max-w-7xl mx-auto px-4 md:px-8 w-full flex flex-col gap-10 mt-4 pb-16">
      <!-- Hero Carousel -->
      <div>
        <HeroSkeleton v-if="heroPending" />
        <HeroCarousel v-else-if="heroSlides.length" :slides="heroSlides" />
        <div v-else class="min-h-[40dvh] flex flex-col items-center justify-center bg-elevated/40 rounded-3xl border border-muted/50 gap-2">
          <p class="text-toned text-sm">Couldn't load current season.</p>
        </div>
      </div>

      <!-- Recommended Rail -->
      <AnimeRail 
        title="Recommended" 
        fetchUrl="/anime?order_by=score&sort=desc&limit=12" 
      />

      <!-- Genre Rails with CSS stagger delay -->
      <div v-if="genresError" class="text-center py-6 text-sm text-toned/80 border border-dashed border-muted/50 rounded-2xl">
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
          <div class="lg:col-span-6 flex flex-col gap-6 md:gap-8 justify-center text-left">
            <h1 class="text-4xl md:text-6xl font-bold tracking-tight text-highlighted leading-[1.1]">
              Anime Tracking,<br>
              <span class="text-primary">Cinematic & Fast.</span>
            </h1>
            <p class="text-sm md:text-base text-toned leading-relaxed max-w-[50ch]">
              Catalog series, track your watchlist, and explore collections in a stunning cinematic interface built for otakus.
            </p>
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 mt-2">
              <UButton
                to="/register"
                label="Join KoKo"
                size="lg"
                color="primary"
                class="rounded-full font-semibold shadow-md shadow-primary/10 px-8 py-3.5 text-center justify-center cursor-pointer"
              />
              <UButton
                to="/browse"
                label="Browse Library"
                size="lg"
                variant="ghost"
                color="neutral"
                class="rounded-full font-medium px-6 py-3.5 text-center justify-center cursor-pointer"
              />
            </div>
          </div>

          <!-- Hero Image with Refraction Glow -->
          <div class="lg:col-span-6 relative w-full flex items-center justify-center">
            <div class="absolute inset-0 bg-primary/10 dark:bg-primary/5 rounded-[40px] blur-3xl -z-10" />
            <div class="w-full relative rounded-3xl overflow-hidden shadow-2xl border border-muted/50 bg-elevated transform transition-all duration-700 hover:scale-[1.01] hover:shadow-primary/5">
              <NuxtImg
                src="/koko_hero.jpg"
                alt="KoKo Cinematic Viewport Preview"
                class="w-full h-full object-cover aspect-[16/10] md:aspect-[16/9] lg:aspect-[4/3]"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- Logo Wall -->
      <section class="py-10 max-w-7xl mx-auto px-6 w-full border-t border-muted/50 flex flex-col md:flex-row items-center justify-between gap-6 opacity-60">
        <span class="text-[10px] font-mono uppercase tracking-wider text-toned">Powered by community standards</span>
        <div class="flex flex-wrap items-center gap-8 justify-center select-none dark:invert dark:brightness-100">
          <img src="https://cdn.simpleicons.org/myanimelist/6B7280" alt="MyAnimeList" class="h-5 opacity-70 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/nuxt/6B7280" alt="Nuxt" class="h-5 opacity-70 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/vuedotjs/6B7280" alt="Vue" class="h-5 opacity-70 hover:opacity-100 transition-opacity" />
          <img src="https://cdn.simpleicons.org/tailwindcss/6B7280" alt="Tailwind CSS" class="h-5 opacity-70 hover:opacity-100 transition-opacity" />
        </div>
      </section>

      <!-- Features Section (Bento Grid) -->
      <section id="features" class="py-20 border-y border-muted/50 bg-elevated/20 transition-all duration-200">
        <div class="max-w-7xl mx-auto px-6 w-full flex flex-col gap-12">
          <div class="max-w-2xl">
            <h2 class="text-3xl font-bold tracking-tight text-highlighted sm:text-4xl leading-tight">
              Designed for the Ultimate View Experience
            </h2>
            <p class="text-sm text-toned mt-3 leading-relaxed">
              Say goodbye to outdated grids and endless clutter. KoKo is built from the ground up for speed, aesthetics, and simplicity.
            </p>
          </div>

          <!-- Bento Grid -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Feature 1: Cinematic UI (col-span-2) -->
            <div class="md:col-span-2 bg-default border border-muted/50 rounded-3xl p-8 flex flex-col justify-between shadow-[0_4px_20px_rgba(0,0,0,0.01)] hover:border-primary/40 hover:shadow-md transition-all duration-300 relative overflow-hidden group min-h-[340px]">
              <div class="absolute -right-10 -bottom-10 w-72 h-44 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-colors" />
              <div class="flex flex-col gap-3 max-w-md z-10">
                <div class="w-10 h-10 rounded-2xl bg-primary/10 flex items-center justify-center text-primary mb-2">
                  <UIcon name="i-solar-gallery-wide-linear" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-highlighted">Cinematic Poster Interface</h3>
                <p class="text-xs text-toned leading-relaxed">
                  Experience anime like never before. High-resolution poster cards, smooth Ken Burns image scaling, and fluid animations mimic premium streaming libraries for a gorgeous distraction-free setup.
                </p>
              </div>
              <!-- Poster Mockup inside the bento cell -->
              <div class="mt-6 flex gap-3 overflow-hidden select-none z-10 translate-y-2 group-hover:translate-y-0 transition-transform duration-500">
                <div v-for="i in 3" :key="i" class="w-24 flex-shrink-0 aspect-[3/4] bg-elevated rounded-xl border border-muted/50 overflow-hidden relative shadow-sm">
                  <div class="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent z-10" />
                  <div class="absolute bottom-2 left-2 right-2 text-[8px] font-medium text-white truncate">Anime Title {{ i }}</div>
                </div>
              </div>
            </div>

            <!-- Feature 2: Personalized Watchlist (col-span-1) -->
            <div class="bg-default border border-muted/50 rounded-3xl p-8 flex flex-col justify-between shadow-[0_4px_20px_rgba(0,0,0,0.01)] hover:border-primary/40 hover:shadow-md transition-all duration-300 relative overflow-hidden group min-h-[340px]">
              <div class="absolute -left-10 -bottom-10 w-52 h-52 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-colors" />
              <div class="flex flex-col gap-3 z-10">
                <div class="w-10 h-10 rounded-2xl bg-primary/10 flex items-center justify-center text-primary mb-2">
                  <UIcon name="i-solar-bookmark-bold-duotone" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-highlighted">Interactive Watchlist</h3>
                <p class="text-xs text-toned leading-relaxed">
                  One-click tracking. Save series, manage your watch history, and plan what to watch next. Everything syncs instantly with your personal profile.
                </p>
              </div>
              <!-- Watchlist visual indicator -->
              <div class="mt-6 bg-elevated border border-muted rounded-2xl p-4 flex flex-col gap-2.5 z-10 transform group-hover:scale-[1.02] transition-transform duration-300 shadow-sm">
                <div class="flex items-center justify-between">
                  <span class="text-[10px] font-semibold text-highlighted">Currently Watching</span>
                  <span class="text-[8px] bg-primary/10 text-primary px-2 py-0.5 rounded-full font-medium">12 Series</span>
                </div>
                <div class="w-full bg-default h-1 rounded-full overflow-hidden">
                  <div class="bg-primary w-2/3 h-full rounded-full" />
                </div>
              </div>
            </div>

            <!-- Feature 3: Instant Discovery (col-span-3) -->
            <div class="md:col-span-3 bg-default border border-muted/50 rounded-3xl p-8 flex flex-col md:flex-row items-center justify-between gap-8 shadow-[0_4px_20px_rgba(0,0,0,0.01)] hover:border-primary/40 hover:shadow-md transition-all duration-300 relative overflow-hidden group">
              <div class="flex flex-col gap-3 max-w-lg z-10">
                <div class="w-10 h-10 rounded-2xl bg-primary/10 flex items-center justify-center text-primary mb-2">
                  <UIcon name="i-solar-magnifer-linear" class="w-5 h-5" />
                </div>
                <h3 class="text-lg font-bold text-highlighted">Fast Browsing & Filtering</h3>
                <p class="text-xs text-toned leading-relaxed">
                  Find what you want, immediately. Highly optimized search querying with instant key-debouncing and genre intersections. Filters anime matching your mood without page refreshes.
                </p>
              </div>
              <!-- Search Mockup -->
              <div class="w-full md:w-80 bg-elevated border border-muted/50 rounded-2xl p-4 flex flex-col gap-3 z-10 shrink-0 transform group-hover:translate-x-1 transition-transform duration-300">
                <div class="flex items-center gap-2 bg-default border border-muted rounded-full px-3 py-1.5 text-[10px] text-toned">
                  <UIcon name="i-solar-magnifer-linear" class="w-3.5 h-3.5 text-primary" />
                  <span>Action, Sci-Fi...</span>
                </div>
                <div class="flex flex-wrap gap-1.5">
                  <span class="text-[9px] bg-primary/10 text-primary border border-primary/20 px-2 py-0.5 rounded-full font-medium">Action</span>
                  <span class="text-[9px] bg-primary/10 text-primary border border-primary/20 px-2 py-0.5 rounded-full font-medium">Sci-Fi</span>
                  <span class="text-[9px] bg-elevated border border-muted text-toned px-2 py-0.5 rounded-full font-medium">+12 more</span>
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
            <h2 class="text-3xl font-bold tracking-tight text-highlighted sm:text-4xl leading-tight">
              Why use KoKo instead of others?
            </h2>
            <p class="text-sm text-toned leading-relaxed">
              Traditional anime trackers are cluttered, slow, and feel like Excel spreadsheets. KoKo is designed for people who appreciate design details, fast loading times, and a focused workspace.
            </p>
          </div>

          <!-- Comparison points layout -->
          <div class="lg:col-span-7 grid grid-cols-1 sm:grid-cols-2 gap-8">
            <div class="flex flex-col gap-3 p-6 bg-elevated/40 border border-muted/30 rounded-3xl hover:border-primary/20 transition-all">
              <div class="text-primary"><UIcon name="i-solar-star-bold" class="w-6 h-6" /></div>
              <h4 class="text-base font-bold text-highlighted">Cinematic Over Metadata</h4>
              <p class="text-xs text-toned leading-relaxed">
                Instead of overwhelming you with wall-to-wall textual statistics, we prioritize large high-resolution posters, visual backdrops, and immersive slides.
              </p>
            </div>

            <div class="flex flex-col gap-3 p-6 bg-elevated/40 border border-muted/30 rounded-3xl hover:border-primary/20 transition-all">
              <div class="text-primary"><UIcon name="i-solar-bolt-bold" class="w-6 h-6" /></div>
              <h4 class="text-base font-bold text-highlighted">Zero Page Reloads</h4>
              <p class="text-xs text-toned leading-relaxed">
                Search and genre selections update instantly in-place. No clicking "Next page" and waiting for database rebuilds.
              </p>
            </div>

            <div class="flex flex-col gap-3 p-6 bg-elevated/40 border border-muted/30 rounded-3xl hover:border-primary/20 transition-all">
              <div class="text-primary"><UIcon name="i-solar-shield-bold" class="w-6 h-6" /></div>
              <h4 class="text-base font-bold text-highlighted">Focused Watchlist</h4>
              <p class="text-xs text-toned leading-relaxed">
                Add items in one click. No complicated forms where you must input episodes watched, status, score, start date, and end date just to save a title.
              </p>
            </div>

            <div class="flex flex-col gap-3 p-6 bg-elevated/40 border border-muted/30 rounded-3xl hover:border-primary/20 transition-all">
              <div class="text-primary"><UIcon name="i-solar-widget-6-bold" class="w-6 h-6" /></div>
              <h4 class="text-base font-bold text-highlighted">Clean and Open Source</h4>
              <p class="text-xs text-toned leading-relaxed">
                No intrusive ads, tracking scripts, or clutter. Completely free and open-source code powered directly by public Jikan APIs.
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- FAQ Section -->
      <section id="faq" class="py-20 border-t border-muted/50 bg-elevated/20 transition-all duration-200">
        <div class="max-w-4xl mx-auto px-6 w-full flex flex-col gap-10">
          <div class="text-center max-w-2xl mx-auto">
            <h2 class="text-3xl font-bold tracking-tight text-highlighted sm:text-4xl">
              Frequently Asked Questions
            </h2>
            <p class="text-sm text-toned mt-3">
              Everything you need to know about the KoKo anime library tracker.
            </p>
          </div>

          <!-- Accordion/Details style list -->
          <div class="flex flex-col gap-4">
            <details class="group bg-default border border-muted/50 rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none">
              <summary class="flex items-center justify-between font-semibold text-highlighted text-sm md:text-base">
                <span>Is KoKo free to use?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-toned">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-toned leading-relaxed">
                Yes! KoKo is completely free. We do not run ads or sell user data. The project is an open source study in designing premium user experiences.
              </p>
            </details>

            <details class="group bg-default border border-muted/50 rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none">
              <summary class="flex items-center justify-between font-semibold text-highlighted text-sm md:text-base">
                <span>Do I need a MyAnimeList account?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-toned">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-toned leading-relaxed">
                No, you do not need an external account. You can sign up directly on KoKo to manage and sync your watchlist locally on our servers.
              </p>
            </details>

            <details class="group bg-default border border-muted/50 rounded-2xl p-6 [&_summary::-webkit-details-marker]:hidden cursor-pointer select-none">
              <summary class="flex items-center justify-between font-semibold text-highlighted text-sm md:text-base">
                <span>Where does the data come from?</span>
                <span class="transition-transform duration-300 group-open:rotate-180 text-toned">
                  <UIcon name="i-solar-alt-arrow-down-linear" class="w-5 h-5" />
                </span>
              </summary>
              <p class="mt-4 text-xs text-toned leading-relaxed">
                All data is powered by the Jikan API, which parses MyAnimeList's extensive public database. This ensures our library is always up-to-date with the latest seasonal airing titles, studios, and genres.
              </p>
            </details>
          </div>
        </div>
      </section>

      <!-- Final CTA Section -->
      <section class="py-24 max-w-7xl mx-auto px-6 w-full">
        <div class="relative bg-primary/5 dark:bg-primary/10 border border-primary/20 rounded-3xl p-12 md:p-16 flex flex-col items-center text-center gap-6 overflow-hidden">
          <div class="absolute -right-20 -bottom-20 w-80 h-80 bg-primary/10 rounded-full blur-3xl" />
          <div class="absolute -left-20 -top-20 w-80 h-80 bg-primary/10 rounded-full blur-3xl" />
          
          <div class="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center text-primary mb-2 animate-bounce">
            <UIcon name="i-solar-heart-bold" class="w-6 h-6" />
          </div>

          <h2 class="text-3xl md:text-4xl font-bold tracking-tight text-highlighted leading-tight max-w-xl">
            Ready to Experience Anime Differently?
          </h2>
          <p class="text-sm text-toned max-w-md">
            Join the KoKo community today. Set up your watchlist, track your favorite series, and find your next watch in seconds.
          </p>
          <div class="flex flex-col sm:flex-row items-center gap-4 mt-2">
            <UButton
              to="/register"
              label="Join KoKo"
              size="lg"
              color="primary"
              class="rounded-full font-semibold shadow-md shadow-primary/10 px-8 py-3.5 cursor-pointer"
            />
            <UButton
              to="/browse"
              label="Browse Library First"
              size="lg"
              variant="ghost"
              color="neutral"
              class="rounded-full font-medium px-6 py-3.5 cursor-pointer"
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
  animation-delay: calc(var(--i) * 100ms);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
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
