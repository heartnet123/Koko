import { ref, computed } from 'vue'

export interface User {
  id: number
  username: string
  email: string
  display_name: string
  avatar_url: string
  bio: string
  created_at: string
  updated_at: string
}

export interface WatchlistItem {
  anime_id: number
  title: string
  image_url: string
  added_at: string
}

export function useAuth() {
  const user = useState<User | null>('auth-user', () => null)
  const token = useState<string | null>('auth-token', () => null)
  const watchlist = useState<WatchlistItem[]>('auth-watchlist', () => [])
  const loading = ref(false)

  // Sync token from localStorage on client-side
  if (import.meta.client) {
    const savedToken = localStorage.getItem('koko_auth_token')
    if (savedToken) {
      token.value = savedToken
    }
  }

  const isAuthenticated = computed(() => !!user.value)

  // Helper for API fetch configuration
  const apiFetch = async <T = any>(path: string, options: any = {}) => {
    const headers = { ...options.headers }
    if (token.value) {
      headers['Authorization'] = `Bearer ${token.value}`
    }

    const config = {
      baseURL: 'http://localhost:8080',
      credentials: 'include' as const,
      ...options,
      headers
    }

    return await $fetch<T>(path, config)
  }

  // Fetch current user profile
  const fetchUser = async () => {
    if (!token.value) {
      user.value = null
      return null
    }

    loading.value = true
    try {
      const response = await apiFetch<{ user: User }>('/api/users/me')
      user.value = response.user
      return response.user
    } catch (err) {
      // Clear invalid token
      logoutLocal()
      return null
    } finally {
      loading.value = false
    }
  }

  // Login
  const login = async (email: string, password: string) => {
    loading.value = true
    try {
      const response = await apiFetch<{ token: string; user: User }>('/api/auth/session', {
        method: 'POST',
        body: { email, password }
      })

      token.value = response.token
      user.value = response.user
      if (import.meta.client) {
        localStorage.setItem('koko_auth_token', response.token)
      }
      // Load watchlist on login
      await fetchWatchlist()
      return { success: true }
    } catch (err: any) {
      return {
        success: false,
        error: err.data?.message || 'Login failed. Please try again.'
      }
    } finally {
      loading.value = false
    }
  }

  // Register
  const register = async (username: string, email: string, password: string) => {
    loading.value = true
    try {
      const response = await apiFetch<{ token: string; user: User }>('/api/users', {
        method: 'POST',
        body: { username, email, password }
      })

      token.value = response.token
      user.value = response.user
      if (import.meta.client) {
        localStorage.setItem('koko_auth_token', response.token)
      }
      watchlist.value = []
      return { success: true }
    } catch (err: any) {
      return {
        success: false,
        error: err.data?.message || 'Registration failed. Please try again.'
      }
    } finally {
      loading.value = false
    }
  }

  // Local logout cleanup
  const logoutLocal = () => {
    token.value = null
    user.value = null
    watchlist.value = []
    if (import.meta.client) {
      localStorage.removeItem('koko_auth_token')
    }
  }

  // Logout
  const logout = async () => {
    loading.value = true
    try {
      await apiFetch('/api/auth/session', {
        method: 'DELETE'
      })
    } catch (err) {
      // Ignore network errors on logout, proceed with local cleanup
    } finally {
      logoutLocal()
      loading.value = false
    }
  }

  // Update profile
  const updateProfile = async (displayName: string, bio: string, avatarUrl: string, password?: string) => {
    loading.value = true
    try {
      const response = await apiFetch<{ user: User }>('/api/users/me', {
        method: 'PATCH',
        body: {
          display_name: displayName,
          bio,
          avatar_url: avatarUrl,
          password
        }
      })
      user.value = response.user
      return { success: true }
    } catch (err: any) {
      return {
        success: false,
        error: err.data?.message || 'Failed to update profile.'
      }
    } finally {
      loading.value = false
    }
  }

  // Watchlist: Fetch
  const fetchWatchlist = async () => {
    if (!token.value) return
    try {
      const response = await apiFetch<WatchlistItem[]>('/api/users/me/watchlist')
      watchlist.value = response
    } catch (err) {
      watchlist.value = []
    }
  }

  // Watchlist: Add
  const addToWatchlist = async (animeId: number, title: string, imageUrl: string) => {
    if (!token.value) return false
    try {
      await apiFetch('/api/users/me/watchlist', {
        method: 'POST',
        body: {
          anime_id: animeId,
          title,
          image_url: imageUrl
        }
      })
      // Refresh list
      await fetchWatchlist()
      return true
    } catch (err) {
      return false
    }
  }

  // Watchlist: Remove
  const removeFromWatchlist = async (animeId: number) => {
    if (!token.value) return false
    try {
      await apiFetch(`/api/users/me/watchlist/${animeId}`, {
        method: 'DELETE'
      })
      // Refresh list
      await fetchWatchlist()
      return true
    } catch (err) {
      return false
    }
  }

  // Check if anime is in watchlist
  const inWatchlist = (animeId: number) => {
    return watchlist.value.some(item => item.anime_id === animeId)
  }

  return {
    user,
    token,
    watchlist,
    loading,
    isAuthenticated,
    fetchUser,
    login,
    register,
    logout,
    updateProfile,
    fetchWatchlist,
    addToWatchlist,
    removeFromWatchlist,
    inWatchlist
  }
}
