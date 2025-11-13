<template>
  <div>
    <h1 class="text-3xl font-bold mb-6 text-gray-800 dark:text-white">Journal</h1>

    <!-- Add New Entry Form -->
    <div class="mb-8">
      <form @submit.prevent="addEntry" class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="md:col-span-2">
            <label for="new-title" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Title</label>
            <input
              id="new-title"
              type="text"
              v-model="newEntry.title"
              placeholder="A memorable day"
              class="mt-1 w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label for="new-date" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Date</label>
            <input
              id="new-date"
              type="date"
              v-model="newEntry.entryDate"
              class="mt-1 w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>
        <div>
          <label for="new-content" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Content</label>
          <textarea
            id="new-content"
            v-model="newEntry.content"
            placeholder="What's on your mind today?"
            rows="4"
            class="mt-1 w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          ></textarea>
        </div>
        <div>
          <label for="new-mood" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Mood (Optional)</label>
           <input
              id="new-mood"
              type="text"
              v-model="newEntry.mood"
              placeholder="Happy, thoughtful, etc."
              class="mt-1 w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
        </div>
        <div class="text-right">
          <button
            type="submit"
            :disabled="!newEntry.title.trim() || !newEntry.entryDate"
            class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed"
          >
            Add Entry
          </button>
        </div>
      </form>
      <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="text-gray-500 dark:text-gray-400">
      Loading entries...
    </div>

    <!-- Entries Display -->
    <div v-else-if="entries && entries.length" class="space-y-4">
       <div
        v-for="entry in entries"
        :key="entry.id"
        class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow"
      >
        <NuxtLink :to="`/journal/${entry.id}`" class="block">
          <div class="flex justify-between items-start">
            <h2 class="font-semibold text-lg text-blue-600 dark:text-blue-400 hover:underline">{{ entry.title }}</h2>
            <span class="text-sm text-gray-500 dark:text-gray-400">{{ new Date(entry.entryDate).toLocaleDateString() }}</span>
          </div>
          <p v-if="entry.mood" class="text-xs text-gray-400 dark:text-gray-500 mt-1">Mood: {{ entry.mood }}</p>
          <p class="text-sm text-gray-600 dark:text-gray-300 mt-2 truncate">
            {{ entry.content }}
          </p>
        </NuxtLink>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md text-center">
      <p class="text-gray-600 dark:text-gray-300">
        You don't have any journal entries yet. Create your first one above!
      </p>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  middleware: ['auth']
});

const { data: entries, pending, error, refresh } = await useApiFetch('/journal', {
  lazy: true,
  server: false,
});

const authStore = useAuthStore();
const config = useRuntimeConfig();

const newEntry = ref({
  title: '',
  content: '',
  mood: '',
  entryDate: new Date().toISOString().split('T')[0], // Default to today
});

const addEntry = async () => {
  if (!newEntry.value.title.trim() || !newEntry.value.entryDate) return;

  try {
    await $fetch('/journal', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: {
        ...newEntry.value,
        mood: newEntry.value.mood || null, // Send null if mood is empty
        entryDate: new Date(newEntry.value.entryDate).toISOString(), // Ensure correct time format
      },
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    });
    
    // Reset form
    newEntry.value = {
      title: '',
      content: '',
      mood: '',
      entryDate: new Date().toISOString().split('T')[0],
    };

    await refresh();
  } catch (err) {
    error.value = 'Failed to create entry. Please try again.';
    console.error('Failed to create entry:', err);
  }
};
</script>