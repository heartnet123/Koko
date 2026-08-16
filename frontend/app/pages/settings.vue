<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Settings — KoKo',
  description: 'Manage your KoKo display name, avatar, biography, and security credentials.'
})

const auth = useAuth()
const router = useRouter()

const displayName = ref('')
const bio = ref('')
const avatarUrl = ref('')
const newPassword = ref('')
const saveSuccess = ref(false)
const saveError = ref('')
const saveLoading = ref(false)

const avatarPresets = [
  'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=150&auto=format&fit=crop',
  'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=150&auto=format&fit=crop',
  'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&auto=format&fit=crop',
  'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=150&auto=format&fit=crop',
  'https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?q=80&w=150&auto=format&fit=crop'
]

const initFields = () => {
  if (!auth.user.value) return
  displayName.value = auth.user.value.display_name || auth.user.value.username
  bio.value = auth.user.value.bio || ''
  avatarUrl.value = auth.user.value.avatar_url || ''
}

onMounted(initFields)

const handleSave = async () => {
  saveSuccess.value = false
  saveError.value = ''

  if (newPassword.value && newPassword.value.length < 8) {
    saveError.value = 'Password must be at least 8 characters long.'
    return
  }

  saveLoading.value = true
  const res = await auth.updateProfile(
    displayName.value.trim(),
    bio.value.trim(),
    avatarUrl.value.trim(),
    newPassword.value || undefined
  )
  saveLoading.value = false

  if (res.success) {
    newPassword.value = ''
    saveSuccess.value = true
    router.push('/profile')
    return
  }

  saveError.value = res.error || 'Failed to save settings.'
}
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-6 animate-fade-in-up">
    <!-- Header -->
    <div class="glass-surface p-6 md:p-8 rounded-3xl border border-[var(--glass-border)] shadow-md flex items-center justify-between">
      <div>
        <div class="inline-flex items-center gap-2 mb-2">
          <div class="w-2 h-4 bg-primary-500 rounded-full " />
          <span class="text-xs font-bold text-primary-400 uppercase tracking-wider font-mono">Preferences</span>
        </div>
        <h1 class="text-2xl md:text-3xl font-bold tracking-tight text-[var(--ui-text-highlighted)]">
          Account Settings
        </h1>
        <p class="text-xs md:text-sm text-[var(--ui-text-toned)] mt-1 font-normal">
          Update your profile avatar, biography, and security credentials.
        </p>
      </div>
      <UButton
        to="/profile"
        label="View Profile"
        icon="i-solar-user-linear"
        variant="ghost"
        color="neutral"
        class="glass-pill rounded-xl px-4 py-2 text-xs font-bold cursor-pointer hover:bg-white/10"
      />
    </div>

    <form v-if="auth.user.value" class="glass-surface-elevated border border-[var(--glass-border)] rounded-3xl p-6 md:p-8 shadow-[var(--shadow-diffuse-lg)] flex flex-col gap-6" @submit.prevent="handleSave">
      <div class="flex flex-col md:flex-row gap-8 items-start">
        <div class="flex flex-col items-center gap-3 w-full md:w-48">
          <UAvatar
            :src="avatarUrl || auth.user.value.avatar_url || 'https://i.pravatar.cc/150'"
            alt="User avatar preview"
            size="2xl"
            class="ring-4 ring-primary-400/30 rounded-2xl shadow-[var(--shadow-diffuse-lg)]"
          />
          <p class="text-xs text-[var(--ui-text-toned)] text-center font-mono">Avatar Preview</p>
        </div>

        <div class="flex-1 w-full flex flex-col gap-5">
          <div>
            <label for="display-name" class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1.5">Display Name</label>
            <input
              id="display-name"
              v-model="displayName"
              type="text"
              class="w-full px-4 py-2.5 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] focus:outline-none focus:ring-2 focus:ring-primary-500/40 transition-all"
              placeholder="Display name"
            />
          </div>

          <div>
            <label for="avatar-url" class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1.5">Avatar Image URL</label>
            <input
              id="avatar-url"
              v-model="avatarUrl"
              type="url"
              class="w-full px-4 py-2.5 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] focus:outline-none focus:ring-2 focus:ring-primary-500/40 transition-all"
              placeholder="https://example.com/avatar.jpg"
            />
          </div>

          <div>
            <span class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-2">Preset Avatars</span>
            <div class="flex flex-wrap gap-2.5">
              <button
                v-for="preset in avatarPresets"
                :key="preset"
                type="button"
                class="w-11 h-11 rounded-xl overflow-hidden ring-2 transition-all cursor-pointer focus:outline-none"
                :class="avatarUrl === preset ? 'ring-primary-400 scale-105 ' : 'ring-transparent hover:ring-white/30'"
                @click="avatarUrl = preset"
              >
                <NuxtImg :src="preset" alt="Avatar preset" class="w-full h-full object-cover" />
              </button>
            </div>
          </div>

          <div>
            <label for="bio" class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1.5">Biography</label>
            <textarea
              id="bio"
              v-model="bio"
              rows="4"
              class="w-full px-4 py-3 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] focus:outline-none focus:ring-2 focus:ring-primary-500/40 transition-all resize-none"
              placeholder="Tell other anime fans about yourself..."
            />
          </div>

          <div>
            <label for="password" class="text-[10px] text-[var(--ui-text-toned)] uppercase tracking-wider font-bold block mb-1.5">New Password (Optional)</label>
            <input
              id="password"
              v-model="newPassword"
              type="password"
              autocomplete="new-password"
              class="w-full px-4 py-2.5 glass-pill rounded-xl text-xs md:text-sm font-semibold text-[var(--ui-text-highlighted)] focus:outline-none focus:ring-2 focus:ring-primary-500/40 transition-all"
              placeholder="Leave blank to retain current password"
            />
          </div>
        </div>
      </div>

      <p v-if="saveError" class="text-xs text-[var(--ui-error)] font-mono font-bold">{{ saveError }}</p>
      <p v-if="saveSuccess" class="text-xs text-[var(--ui-success)] font-mono font-bold">Settings saved successfully.</p>

      <div class="flex justify-end gap-3 border-t border-[var(--glass-border-subtle)] pt-6">
        <UButton
          to="/profile"
          label="Cancel"
          variant="ghost"
          color="neutral"
          class="glass-pill rounded-xl px-5 py-2.5 text-xs font-bold cursor-pointer hover:bg-white/10"
        />
        <UButton
          type="submit"
          label="Save Settings"
          icon="i-solar-diskette-linear"
          color="primary"
          :loading="saveLoading"
          class="rounded-xl px-6 py-2.5 text-xs font-bold shadow-[var(--shadow-diffuse-accent)] cursor-pointer hover:scale-[1.02] active:scale-95 transition-all"
        />
      </div>
    </form>
  </div>
</template>