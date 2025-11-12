import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to, from) => {
  const auth = useAuthStore()

  if (!auth.isLoggedIn) {
    if (to.path !== '/login' && to.path !== '/register') {
      return navigateTo('/login', { replace: true })
    }
  } else {
    if (to.path === '/login' || to.path === '/register') {
      return navigateTo('/', { replace: true })
    }
  }
})