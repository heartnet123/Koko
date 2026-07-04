<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '~/composables/useAuth'

useSeoMeta({
  title: 'Koko - Settings',
  description: 'Manage your Koko display name, avatar, biography, and password.'
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
  <div class="max-w-4xl mx-auto px-4 md:px-8 py-8 w-full flex flex-col gap-6">
    <div class="flex items-baseline justify-between border-b border-muted/50 pb-4">
      <div>
        <h1 class="text-3xl font-semibold tracking-tighter text-highlighted mb-1">Settings</h1>
        <p class="text-toned text-xs">Update profile details and account password.</p>
      </div>
      <UButton
        to="/profile"
        label="View Profile"
        icon="i-solar-user-linear"
        variant="ghost"
        color="neutral"
        class="rounded-xl cursor-pointer"
      />
    </div>

    <form v-if="auth.user.value" class="bg-elevated border border-muted/50 rounded-3xl p-6 md:p-8 shadow-lg flex flex-col gap-6" @submit.prevent="handleSave">
      <div class="flex flex-col md:flex-row gap-6 items-start">
        <div class="flex flex-col items-center gap-3 w-full md:w-48">
          <UAvatar
            :src="avatarUrl || auth.user.value.avatar_url || 'https://i.pravatar.cc/150'"
            alt="User avatar preview"
            size="2xl"
            class="ring-4 ring-primary/20"
          />
          <p class="text-xs text-toned text-center">Choose preset or paste image URL.</p>
        </div>

        <div class="flex-1 w-full flex flex-col gap-4">
          <div>
            <label for="display-name" class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Display Name</label>
            <input
              id="display-name"
              v-model="displayName"
              type="text"
              class="w-full px-4 py-2.5 bg-default border border-muted rounded-xl text-sm text-default focus:outline-none focus:border-primary/40 transition-colors"
              placeholder="Display name"
            />
          </div>

          <div>
            <label for="avatar-url" class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Avatar URL</label>
            <input
              id="avatar-url"
              v-model="avatarUrl"
              type="url"
              class="w-full px-4 py-2.5 bg-default border border-muted rounded-xl text-sm text-default focus:outline-none focus:border-primary/40 transition-colors"
              placeholder="https://example.com/avatar.jpg"
            />
          </div>

          <div>
            <span class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-2">Avatar Presets</span>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="preset in avatarPresets"
                :key="preset"
                type="button"
                class="w-12 h-12 rounded-full overflow-hidden ring-2 transition-all focus:outline-none focus-visible:ring-primary"
                :class="avatarUrl === preset ? 'ring-primary' : 'ring-transparent hover:ring-muted'"
                @click="avatarUrl = preset"
              >
                <NuxtImg :src="preset" alt="Avatar preset" class="w-full h-full object-cover" />
              </button>
            </div>
          </div>

          <div>
            <label for="bio" class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">Bio</label>
            <textarea
              id="bio"
              v-model="bio"
              rows="5"
              class="w-full px-4 py-3 bg-default border border-muted rounded-xl text-sm text-default focus:outline-none focus:border-primary/40 transition-colors resize-none"
              placeholder="Tell people about yourself..."
            />
          </div>

          <div>
            <label for="password" class="text-[10px] text-toned uppercase tracking-wider font-semibold block mb-1">New Password</label>
            <input
              id="password"
              v-model="newPassword"
              type="password"
              autocomplete="new-password"
              class="w-full px-4 py-2.5 bg-default border border-muted rounded-xl text-sm text-default focus:outline-none focus:border-primary/40 transition-colors"
              placeholder="Leave blank to keep current password"
            />
          </div>
        </div>
      </div>

      <p v-if="saveError" class="text-sm text-error font-medium">{{ saveError }}</p>
      <p v-if="saveSuccess" class="text-sm text-success font-medium">Settings saved.</p>

      <div class="flex justify-end gap-3 border-t border-muted/50 pt-6">
        <UButton
          to="/profile"
          label="Cancel"
          variant="ghost"
          color="neutral"
          class="rounded-xl cursor-pointer"
        />
        <UButton
          type="submit"
          label="Save Settings"
          icon="i-solar-diskette-linear"
          color="primary"
          :loading="saveLoading"
          class="rounded-xl cursor-pointer"
        />
      </div>
    </form>
  </div>
</template>
