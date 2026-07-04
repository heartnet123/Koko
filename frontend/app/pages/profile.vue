<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Koko - My Profile',
  description: 'View your Koko profile, account details, and watchlist summary.'
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
  if (!dateStr) return 'Active Member'
  const date = new Date(dateStr)
  return `Joined in ${date.toLocaleString('default', { month: 'long' })} ${date.getFullYear()}`
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-8">
    <div v-if="auth.user.value" class="flex flex-col md:flex-row gap-8 items-start">
      <div class="w-full md:w-80 bg-elevated border border-muted/50 rounded-3xl p-6 flex flex-col items-center text-center shadow-lg">
        <UAvatar
          :src="auth.user.value.avatar_url || 'https://i.pravatar.cc/150'"
          alt="User avatar"
          size="2xl"
          class="ring-4 ring-primary/20 mb-4"
        />
        <h2 class="text-xl font-semibold text-highlighted tracking-tight mb-1">
          {{ auth.user.value.display_name || auth.user.value.username }}
        </h2>
        <p class="text-toned text-xs mb-3">@{{ auth.user.value.username }}</p>
        <span class="text-[10px] text-toned font-medium uppercase tracking-wider bg-default/80 border border-muted px-2.5 py-1 rounded-full mb-6">
          {{ formatDate(auth.user.value.created_at) }}
        </span>

        <p class="text-sm text-default mb-6 line-clamp-4 px-2 italic">
          {{ auth.user.value.bio || 'No biography written yet.' }}
        </p>

        <div class="w-full border-t border-muted/50 pt-6 flex items-center justify-around text-center mb-6">
          <div>
            <span class="text-xl font-bold text-highlighted block">{{ auth.watchlist.value.length }}</span>
            <span class="text-[10px] text-toned uppercase tracking-wider font-semibold">Watchlist</span>
          </div>
        </div>

        <div class="w-full flex flex-col gap-2">
          <UButton
            to="/settings"
            label="Settings"
            icon="i-solar-settings-linear"
            variant="outline"
            color="neutral"
            class="w-full rounded-xl flex justify-center py-2.5 text-xs font-semibold cursor-pointer"
          />
          <UButton
            label="Log Out"
            icon="i-solar-logout-linear"
            variant="subtle"
            color="red"
            class="w-full rounded-xl flex justify-center py-2.5 text-xs font-semibold cursor-pointer"
            @click="handleLogout"
          />
        </div>
      </div>

      <div class="flex-1 w-full bg-elevated border border-muted/50 rounded-3xl p-6 md:p-8 shadow-lg">
        <div class="flex flex-col gap-6">
          <div>
            <h3 class="text-lg font-semibold text-highlighted tracking-tight mb-1 border-b border-muted/50 pb-3">Account Details</h3>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Username</span>
              <span class="text-sm text-highlighted font-medium">@{{ auth.user.value.username }}</span>
            </div>
            <div>
              <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Email Address</span>
              <span class="text-sm text-highlighted font-medium">{{ auth.user.value.email }}</span>
            </div>
            <div>
              <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Display Name</span>
              <span class="text-sm text-highlighted font-medium">{{ auth.user.value.display_name || auth.user.value.username }}</span>
            </div>
            <div>
              <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Custom Avatar</span>
              <span class="text-xs text-toned truncate block max-w-xs">{{ auth.user.value.avatar_url || 'Using default avatar' }}</span>
            </div>
          </div>
          <div class="mt-2">
            <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Bio / About Me</span>
            <div class="p-4 bg-default border border-muted/50 rounded-2xl text-sm text-default leading-relaxed min-h-[100px]">
              {{ auth.user.value.bio || 'No biography written yet.' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
