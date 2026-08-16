<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'My Profile — KoKo',
  description: 'View your KoKo account details, preferences, and watchlist statistics.'
})

const auth = useAuth()
const router = useRouter()

onMounted(async () => {
  await auth.fetchWatchlist()
})

const handleLogout = async () => {
  await auth.logout()
  router.push('/login')
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return 'Active Otaku'
  const date = new Date(dateStr)
  return `Member since ${date.toLocaleString('default', { month: 'short' })} ${date.getFullYear()}`
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-8 animate-fade-in-up">
    <!-- Header -->
    <div class="glass-surface p-6 md:p-8 rounded-3xl border border-[var(--glass-border)] shadow-md flex items-center justify-between">
      <div>
        <div class="inline-flex items-center gap-2 mb-2">
          <div class="w-2 h-4 bg-primary-500 rounded-full " />
          <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">User Space</span>
        </div>
        <h1 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          Account Profile
        </h1>
        <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-1 font-normal">
          Manage your identity, viewing history summary, and preferences.
        </p>
      </div>
      <UButton
        to="/settings"
        label="Edit Profile"
        icon="i-solar-pen-linear"
        variant="ghost"
        color="neutral"
        class="glass-pill rounded-xl px-4 py-2 text-xs font-bold cursor-pointer hover:bg-white/10"
      />
    </div>

    <div v-if="auth.user.value" class="flex flex-col md:flex-row gap-8 items-start">
      <!-- Left Card (Identity & Stats) -->
      <div class="w-full md:w-80 glass-surface-elevated border border-[var(--glass-border)] rounded-3xl p-6 flex flex-col items-center text-center shadow-[var(--shadow-diffuse-lg)]">
        <div class="relative mb-4">
          <UAvatar
            :src="auth.user.value.avatar_url || 'https://i.pravatar.cc/150'"
            alt="User avatar"
            size="2xl"
            class="ring-4 ring-primary-400/30 rounded-2xl shadow-[var(--shadow-diffuse-lg)]"
          />
          <div class="absolute -bottom-1 -right-1 w-5 h-5 rounded-full bg-green-500 ring-2 ring-[var(--ui-bg)]" title="Online" />
        </div>

        <h2 class="text-xl font-bold text-[var(--ui-text-highlighted)] tracking-tight mb-1">
          {{ auth.user.value.display_name || auth.user.value.username }}
        </h2>
        <p class="text-primary-400 text-xs font-mono font-semibold mb-3">@{{ auth.user.value.username }}</p>
        <span class="text-[11px] text-[var(--ui-text-toned)] font-mono uppercase tracking-wider glass-pill px-3 py-1 rounded-full mb-6">
          {{ formatDate(auth.user.value.created_at) }}
        </span>

        <p class="text-xs text-[var(--ui-text)] mb-6 line-clamp-4 px-2 italic font-normal">
          "{{ auth.user.value.bio || 'No biography written yet.' }}"
        </p>

        <div class="w-full border-t border-[var(--glass-border-subtle)] pt-6 flex items-center justify-around text-center mb-6">
          <div class="glass-pill px-6 py-3 rounded-2xl w-full">
            <span class="text-2xl font-bold text-primary-400 font-mono block">{{ auth.watchlist.value.length }}</span>
            <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-mono font-bold">Watchlist Items</span>
          </div>
        </div>

        <div class="w-full flex flex-col gap-2.5">
          <UButton
            to="/settings"
            label="Settings"
            icon="i-solar-settings-linear"
            variant="ghost"
            color="neutral"
            class="w-full glass-pill rounded-xl flex justify-center py-2.5 text-xs font-bold cursor-pointer hover:bg-white/10"
          />
          <button
            type="button"
            class="w-full py-2.5 rounded-xl flex items-center justify-center gap-1.5 text-xs font-bold text-[var(--ui-error)] hover:text-[var(--ui-error)]/80 glass-pill hover:bg-[var(--ui-error)]/10 hover:border-[var(--ui-error)]/30 transition-all cursor-pointer font-mono"
            @click="handleLogout"
          >
            <UIcon name="i-solar-logout-linear" class="w-4 h-4" />
            <span>Log Out</span>
          </button>
        </div>
      </div>

      <!-- Right Card (Account Details) -->
      <div class="flex-1 w-full glass-surface border border-[var(--glass-border)] rounded-3xl p-6 md:p-8 shadow-[var(--shadow-diffuse-lg)]">
        <div class="flex flex-col gap-6">
          <div>
            <h3 class="text-base font-bold text-[var(--ui-text-highlighted)] tracking-tight mb-1 border-b border-[var(--glass-border-subtle)] pb-3">Account Details</h3>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div class="glass-pill p-4 rounded-2xl">
              <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1">Username</span>
              <span class="text-xs font-bold text-[var(--ui-text-highlighted)]">@{{ auth.user.value.username }}</span>
            </div>
            <div class="glass-pill p-4 rounded-2xl">
              <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1">Email Address</span>
              <span class="text-xs font-bold text-[var(--ui-text-highlighted)] truncate block">{{ auth.user.value.email }}</span>
            </div>
            <div class="glass-pill p-4 rounded-2xl">
              <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1">Display Name</span>
              <span class="text-xs font-bold text-[var(--ui-text-highlighted)]">{{ auth.user.value.display_name || auth.user.value.username }}</span>
            </div>
            <div class="glass-pill p-4 rounded-2xl">
              <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1">Avatar Preset</span>
              <span class="text-[11px] text-[var(--ui-text-toned)] truncate block max-w-xs">{{ auth.user.value.avatar_url ? 'Custom avatar URL' : 'Default gravatar' }}</span>
            </div>
          </div>
          <div class="mt-2">
            <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-mono font-bold block mb-2">About / Biography</span>
            <div class="p-5 glass-pill rounded-2xl text-xs md:text-sm text-[var(--ui-text)] leading-relaxed min-h-[110px] font-normal">
              {{ auth.user.value.bio || 'No biography written yet.' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>