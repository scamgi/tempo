<template>
  <form @submit.prevent="handleLogin" class="space-y-4">
    <div v-if="errorMsg" class="p-3 bg-red-100 text-red-700 rounded-md">
      {{ errorMsg }}
    </div>
    <div>
      <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Email</label>
      <input type="email" id="email" v-model="email" required
        class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" />
    </div>
    <div>
      <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Password</label>
      <input type="password" id="password" v-model="password" required
        class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" />
    </div>
    <button type="submit" :disabled="isLoading"
      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-blue-300">
      {{ isLoading ? 'Logging in...' : 'Login' }}
    </button>
    <div class="text-center">
      <NuxtLink to="/register" class="text-sm text-blue-600 hover:underline">
        Don't have an account? Register
      </NuxtLink>
    </div>
  </form>
</template>

<script setup>
definePageMeta({
  layout: 'auth'
});

import { useAuthStore } from '~/stores/auth';

const authStore = useAuthStore();
const router = useRouter();

const email = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  isLoading.value = true;
  errorMsg.value = '';
  try {
    await authStore.login(email.value, password.value);
    router.push('/');
  } catch (error) {
    errorMsg.value = error.message || 'Failed to login. Please check your credentials.';
  } finally {
    isLoading.value = false;
  }
};
</script>