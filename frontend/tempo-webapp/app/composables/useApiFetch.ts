import { useAuthStore } from '~/stores/auth';

export const useApiFetch: typeof useFetch = (request, opts?) => {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();
  const defaults: typeof opts = {
    baseURL: config.public.apiBase,
    // cache request
    key: request as string,
    // set user token if connected
    headers: authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {},

    onResponseError: ({ response }) => {
      if (response.status === 401) {
        authStore.logout();
        navigateTo('/login');
      }
    }
  };
  // for nice deep defaults, please use unjs/defu or similar
  const params = { ...defaults, ...opts };
  return useFetch(request, params);
};
