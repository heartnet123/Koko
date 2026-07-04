import { useAuth } from '~/composables/useAuth'

export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()

  // If client-side and token is present but user not loaded, fetch it
  if (import.meta.client && !auth.user.value && auth.token.value) {
    await auth.fetchUser()
  }

  const protectedRoutes = ['/profile', '/settings', '/watchlist']
  if (protectedRoutes.includes(to.path) && !auth.isAuthenticated.value) {
    return navigateTo('/login')
  }

  // Redirect to profile if trying to access login or register page while authenticated
  if (['/login', '/register'].includes(to.path) && auth.isAuthenticated.value) {
    return navigateTo('/profile')
  }
})
