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
  <div class="min-h-screen bg-default text-default font-sans flex flex-col selection:bg-primary-100 selection:text-primary-900 overflow-x-hidden">
    <!-- Navbar -->
    <header class="sticky top-0 w-full bg-default/85 backdrop-blur-md z-50 border-b border-muted/50 transition-all duration-200">
      <div class="max-w-7xl mx-auto px-6 h-16 md:h-20 flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary rounded-lg">
          <div class="w-1.5 h-6 bg-primary rounded-full" />
          <span class="text-2xl font-bold tracking-tighter text-highlighted">KoKo</span>
        </NuxtLink>

        <!-- Desktop Navigation Links -->
        <nav class="hidden md:flex items-center gap-8">
          <NuxtLink
            v-for="item in menuItems"
            :key="item.label"
            :to="item.to"
            class="text-sm font-medium text-toned hover:text-default transition-colors cursor-pointer"
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
            class="rounded-full cursor-pointer"
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
              class="rounded-full font-semibold shadow-md shadow-primary/10 cursor-pointer"
            />
            <div v-else class="flex items-center gap-3">
              <UButton
                to="/login"
                label="Sign In"
                variant="ghost"
                color="neutral"
                class="rounded-full font-medium cursor-pointer"
              />
              <UButton
                to="/login"
                label="Get Started"
                color="primary"
                class="rounded-full font-semibold shadow-md shadow-primary/10 cursor-pointer"
              />
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
        <div v-if="isMobileMenuOpen" class="md:hidden border-t border-muted/50 bg-default px-6 py-6 flex flex-col gap-4 shadow-xl">
          <NuxtLink
            v-for="item in menuItems"
            :key="item.label"
            :to="item.to"
            class="text-base font-medium text-toned hover:text-default transition-colors py-1.5"
            @click="isMobileMenuOpen = false"
          >
            {{ item.label }}
          </NuxtLink>
          <div class="border-t border-muted/50 my-2 pt-4 flex flex-col gap-3">
            <ClientOnly>
              <UButton
                v-if="auth.isAuthenticated.value"
                to="/browse"
                label="Go to Dashboard"
                color="primary"
                block
                class="rounded-full font-semibold shadow-md shadow-primary/10 py-3 text-center justify-center cursor-pointer"
                @click="isMobileMenuOpen = false"
              />
              <div v-else class="flex flex-col gap-2">
                <UButton
                  to="/login"
                  label="Sign In"
                  variant="ghost"
                  color="neutral"
                  block
                  class="rounded-full font-medium py-2.5 text-center justify-center cursor-pointer"
                  @click="isMobileMenuOpen = false"
                />
                <UButton
                  to="/login"
                  label="Get Started"
                  color="primary"
                  block
                  class="rounded-full font-semibold shadow-md shadow-primary/10 py-3 text-center justify-center cursor-pointer"
                  @click="isMobileMenuOpen = false"
                />
              </div>
            </ClientOnly>
          </div>
        </div>
      </Transition>
    </header>

    <!-- Main Content Slot -->
    <main class="flex-1 w-full">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="border-t border-muted/50 bg-elevated/30 py-16 px-6 mt-16 transition-all duration-200">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-10">
        <div class="flex flex-col gap-4">
          <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer">
            <div class="w-1.5 h-6 bg-primary rounded-full" />
            <span class="text-2xl font-bold tracking-tighter text-highlighted">KoKo</span>
          </NuxtLink>
          <p class="text-xs text-toned max-w-[250px] leading-relaxed">
            Your cinematic anime catalog, tracker, and social companion. Discover series, build your watchlist, and sync effortlessly.
          </p>
        </div>
        <div>
          <h4 class="text-xs font-semibold uppercase tracking-wider text-highlighted mb-4">Discover</h4>
          <ul class="flex flex-col gap-2">
            <li><NuxtLink to="/browse" class="text-sm text-toned hover:text-default transition-colors">Browse Library</NuxtLink></li>
            <li><NuxtLink to="/trending" class="text-sm text-toned hover:text-default transition-colors">Trending Anime</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-semibold uppercase tracking-wider text-highlighted mb-4">Product</h4>
          <ul class="flex flex-col gap-2">
            <li><a href="#features" class="text-sm text-toned hover:text-default transition-colors">Features</a></li>
            <li><a href="#why-koko" class="text-sm text-toned hover:text-default transition-colors">Why KoKo</a></li>
            <li><NuxtLink to="/login" class="text-sm text-toned hover:text-default transition-colors">Join Now</NuxtLink></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-semibold uppercase tracking-wider text-highlighted mb-4">Legal</h4>
          <p class="text-xs text-toned leading-relaxed">
            KoKo is powered by Jikan API and MyAnimeList. All data and content belong to their respective owners.
          </p>
          <p class="text-[10px] text-toned mt-4 font-mono">
            &copy; 2026 KoKo. Made with love for anime fans.
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
