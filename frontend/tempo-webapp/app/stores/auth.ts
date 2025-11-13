import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const token = useCookie('tempo-token', { maxAge: 60 * 60 * 24 * 3 }); // 3 days

  const isLoggedIn = computed(() => !!token.value);

  async function register(username, email, password) {
    const config = useRuntimeConfig()
    try {
      await $fetch('/users/register', {
        baseURL: config.public.apiBase,
        method: 'POST',
        body: { username, email, password },
      })
    } catch (error) {
      const errorMsg = error.data?.message || 'An error occurred during registration.';
      throw new Error(errorMsg);
    }
  }

  async function login(email, password) {
    const config = useRuntimeConfig()
    try {
      const response = await $fetch<{ token: string }>('/users/login', {
        baseURL: config.public.apiBase,
        method: 'POST',
        body: { email, password },
      })
      token.value = response.token
    } catch (error) {
      const errorMsg = error.data?.message || 'An error occurred during login.';
      throw new Error(errorMsg);
    }
  }

  function logout() {
    token.value = null
  }

  return {
    token,
    isLoggedIn,
    register,
    login,
    logout,
  }
})