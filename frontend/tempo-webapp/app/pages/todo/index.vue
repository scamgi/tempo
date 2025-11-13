<template>
  <div>
    <h1 class="text-3xl font-bold mb-6 text-gray-800 dark:text-white">To-Do Lists</h1>

    <!-- Add New List Form -->
    <div class="mb-8">
      <form @submit.prevent="addList" class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md flex items-center gap-4">
        <input
          type="text"
          v-model="newListTitle"
          placeholder="Add a new list..."
          class="flex-grow px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
        <button
          type="submit"
          :disabled="!newListTitle.trim()"
          class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed"
        >
          Add List
        </button>
      </form>
      <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="text-gray-500 dark:text-gray-400">
      Loading lists...
    </div>

    <!-- Lists Display -->
    <div v-else-if="lists && lists.length" class="space-y-4">
      <div
        v-for="list in lists"
        :key="list.id"
        class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow"
      >
        <NuxtLink :to="`/todo/${list.id}`" class="font-semibold text-lg text-blue-600 dark:text-blue-400 hover:underline">
          {{ list.title }}
        </NuxtLink>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Created on: {{ new Date(list.createdAt).toLocaleDateString() }}
        </p>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md text-center">
      <p class="text-gray-600 dark:text-gray-300">
        You don't have any to-do lists yet. Create your first one above!
      </p>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  middleware: ['auth']
});

const { data: lists, pending, error, refresh } = await useApiFetch('/lists', {
  lazy: true,
  server: false,
});

const authStore = useAuthStore();
const config = useRuntimeConfig();
const newListTitle = ref('');

const addList = async () => {
  if (!newListTitle.value.trim()) return;

  try {
    // Use $fetch for actions triggered by user interaction
    await $fetch('/lists', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: { title: newListTitle.value },
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    });
    newListTitle.value = '';
    await refresh();
  } catch (err) {
    error.value = 'Failed to create list. Please try again.';
    console.error('Failed to create list:', err);
  }
};
</script>