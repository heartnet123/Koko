<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '~/composables/useAuth'

const auth = useAuth()
const colorMode = useColorMode()

const toggleColorMode = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}

const isMobileMenuOpen = ref(false)

const menuItems = [
  { label: 'Features', to: '#features' },
  { label: 'Why KoKo', to: '#why-koko' },
  { label: 'FAQ', to: '#faq' },
  { label: 'Browse', to: '/browse' },
]
</script>

<template>
  <div class="relative min-h-[100dvh] bg-[var(--ui-bg)] text-[var(--ui-text)] font-sans flex flex-col selection:bg-primary-500/20 selection:text-primary-400 overflow-x-hidden">
    <!-- Ambient glowing mesh background -->
    <div class="ambient-mesh" aria-hidden="true">
      <div class="ambient-mesh-blob-1" />
      <div class="ambient-mesh-blob-2" />
      <div class="ambient-mesh-blob-3" />
    </div>

    <!-- Floating / Sticky Frosted Glass Navbar -->
    <header class="sticky top-0 w-full glass-surface z-50 border-b border-[var(--glass-border-subtle)] transition-all duration-200">
      <div class="max-w-7xl mx-auto px-6 h-16 md:h-20 flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 rounded-xl group">
          <div class="w-2 h-6 bg-gradient-to-b from-primary-400 to-primary-600 rounded-full  transition-transform group-hover:scale-y-110" />
          <span class="text-2xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">KoKo</span>
        </NuxtLink>

        <!-- Desktop Navigation Links -->
        <nav class="hidden md:flex items-center gap-8">
          <NuxtLink
            v-for="item in menuItems"
            :key="item.label"
            :to="item.to"
            class="text-sm font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors cursor-pointer py-1 px-2.5 rounded-lg hover:bg-white/5"
          >
            {{ item.label }}
          </NuxtLink>
        </nav>

        <!-- CTAs / Actions -->
        <div class="hidden md:flex items-center gap-4">
          <!-- Dark mode toggle -->
          <UButton
            :icon="colorMode.value === 'dark' ? 'i-solar-sun-bold-duotone' : 'i-solar-moon-bold-duotone'"
            color="neutral"
            variant="ghost"
            class="rounded-full cursor-pointer hover:bg-white/10"
            aria-label="Toggle theme"
            @click="toggleColorMode"
          />

          <!-- Action Button -->
          <ClientOnly>
            <UButton
              v-if="auth.isAuthenticated.value"
              to="/browse"
              label="Go to Dashboard"
              color="primary"
              icon="i-solar-arrow-right-linear"
              class="rounded-full font-bold shadow-[var(--shadow-diffuse-accent)] px-6 cursor-pointer"
            />
            <div v-else class="flex items-center gap-3">
              <UButton
                to="/login"
                label="Sign In"
                variant="ghost"
                color="neutral"
                class="rounded-full font-semibold cursor-pointer hover:bg-white/10"
              />
              <span v-magnetic class="inline-block transition-transform duration-300 ease-[cubic-bezier(0.16,1,0.3,1)] will-change-transform">
                <UButton
                  to="/login"
                  label="Get Started"
                  color="primary"
                  class="rounded-full font-bold shadow-[var(--shadow-diffuse-accent)] px-6 cursor-pointer"
                />
              </span>
            </div>
          </ClientOnly>
        </div>

        <!-- Mobile Menu Toggle -->
        <div class="flex items-center gap-2 md:hidden">
          <UButton
            :icon="colorMode.value === 'dark' ? 'i-solar-sun-bold-duotone' : 'i-solar-moon-bold-duotone'"
            color="neutral"
            variant="ghost"
            class="rounded-full cursor-pointer"
            aria-label="Toggle theme"
            @click="toggleColorMode"
          />
          <UButton
            icon="i-solar-hamburger-menu-linear"
            color="neutral"
            variant="ghost"
            class="rounded-full cursor-pointer"
            aria-label="Open navigation menu"
            @click="isMobileMenuOpen = !isMobileMenuOpen"
          />
        </div>
      </div>

      <!-- Mobile Dropdown Navigation -->
      <Transition name="fade">
        <div v-if="isMobileMenuOpen" class="md:hidden glass-surface-elevated border-t border-[var(--glass-border)] px-6 py-6 flex flex-col gap-4 ">
          <NuxtLink
            v-for="item in menuItems"
            :key="item.label"
            :to="item.to"
            class="text-base font-semibold text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors py-1.5"
            @click="isMobileMenuOpen = false"
          >
            {{ item.label }}
          </NuxtLink>
          <div class="border-t border-[var(--ui-border-muted)] my-2 pt-4 flex flex-col gap-3">
            <ClientOnly>
              <UButton
                v-if="auth.isAuthenticated.value"
                to="/browse"
                label="Go to Dashboard"
                color="primary"
                block
                class="rounded-full font-bold shadow-[var(--shadow-diffuse-accent)] py-3 text-center justify-center cursor-pointer"
                @click="isMobileMenuOpen = false"
              />
              <div v-else class="flex flex-col gap-2">
                <UButton
                  to="/login"
                  label="Sign In"
                  variant="ghost"
                  color="neutral"
                  block
                  class="rounded-full font-semibold py-2.5 text-center justify-center cursor-pointer"
                  @click="isMobileMenuOpen = false"
                />
                <UButton
                  to="/login"
                  label="Get Started"
                  color="primary"
                  block
                  class="rounded-full font-bold shadow-[var(--shadow-diffuse-accent)] py-3 text-center justify-center cursor-pointer"
                  @click="isMobileMenuOpen = false"
                />
              </div>
            </ClientOnly>
          </div>
        </div>
      </Transition>
    </header>

    <!-- Main Content Slot -->
    <main class="relative flex-1 w-full z-10">
      <slot />
    </main>

    <!-- Glass Footer -->
    <footer class="relative border-t border-[var(--glass-border-subtle)] glass-surface py-16 px-6 mt-16 transition-all duration-200 z-10">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-10">
        <div class="flex flex-col gap-4">
          <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer">
            <div class="w-2 h-6 bg-gradient-to-b from-primary-400 to-primary-600 rounded-full " />
            <span class="text-2xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">KoKo</span>
          </NuxtLink>
          <p class="text-xs text-[var(--ui-text-toned)] max-w-[280px] leading-relaxed font-medium">
            Your cinematic anime catalog, tracker, and media discovery companion.
          </p>
        </div>
        <div>
          <h4 class="text-xs font-bold uppercase tracking-wider text-[var(--ui-text-highlighted)] mb-4">Discover</h4>
          <ul class="flex flex-col gap-2.5">
            <li><NuxtLink to="/browse" class="text-sm font-medium text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors">Browse Library</NuxtLink></li>
            <li><NuxtLink to="/trending" class="text-sm font-medium text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors">Trending Anime</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-bold uppercase tracking-wider text-[var(--ui-text-highlighted)] mb-4">Product</h4>
          <ul class="flex flex-col gap-2.5">
            <li><a href="#features" class="text-sm font-medium text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors">Features</a></li>
            <li><a href="#why-koko" class="text-sm font-medium text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors">Why KoKo</a></li>
            <li><NuxtLink to="/login" class="text-sm font-medium text-[var(--ui-text-toned)] hover:text-[var(--ui-text-highlighted)] transition-colors">Join Now</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-bold uppercase tracking-wider text-[var(--ui-text-highlighted)] mb-4">Legal</h4>
          <p class="text-xs text-[var(--ui-text-toned)] leading-relaxed font-medium">
            KoKo is powered by Jikan API and MyAnimeList. All data and artwork belong to their respective owners.
          </p>
          <p class="text-[11px] text-[var(--ui-text-toned)] mt-4 font-mono">
            &copy; 2026 KoKo. High-speed media tracker.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>