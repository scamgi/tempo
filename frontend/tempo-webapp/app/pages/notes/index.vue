<template>
  <div>
    <h1 class="text-3xl font-bold mb-6 text-gray-800 dark:text-white">Notes</h1>

    <!-- Add New Note Form -->
    <div class="mb-8">
      <form @submit.prevent="addNote" class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md space-y-4">
        <div>
          <input
            type="text"
            v-model="newNoteTitle"
            placeholder="Note title..."
            class="w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
        <div>
          <textarea
            v-model="newNoteContent"
            placeholder="Start writing your note here..."
            rows="3"
            class="w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          ></textarea>
        </div>
        <div class="text-right">
          <button
            type="submit"
            :disabled="!newNoteTitle.trim()"
            class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed"
          >
            Add Note
          </button>
        </div>
      </form>
      <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="text-gray-500 dark:text-gray-400">
      Loading notes...
    </div>

    <!-- Notes Display -->
    <div v-else-if="notes && notes.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="note in notes"
        :key="note.id"
        class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow flex flex-col"
      >
        <NuxtLink :to="`/notes/${note.id}`" class="font-semibold text-lg text-blue-600 dark:text-blue-400 hover:underline mb-2">
          {{ note.title }}
        </NuxtLink>
        <p class="text-sm text-gray-600 dark:text-gray-400 flex-grow">
          {{ note.content.substring(0, 100) }}{{ note.content.length > 100 ? '...' : '' }}
        </p>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-3 pt-2 border-t border-gray-200 dark:border-gray-700">
          Updated: {{ new Date(note.updatedAt).toLocaleDateString() }}
        </p>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md text-center">
      <p class="text-gray-600 dark:text-gray-300">
        You don't have any notes yet. Create your first one above!
      </p>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  middleware: ['auth']
});

const { data: notes, pending, error, refresh } = await useApiFetch('/notes', {
  lazy: true,
  server: false,
});

const authStore = useAuthStore();
const config = useRuntimeConfig();
const newNoteTitle = ref('');
const newNoteContent = ref('');

const addNote = async () => {
  if (!newNoteTitle.value.trim()) return;

  try {
    await $fetch('/notes', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: { 
        title: newNoteTitle.value,
        content: newNoteContent.value
      },
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    });
    newNoteTitle.value = '';
    newNoteContent.value = '';
    await refresh();
  } catch (err) {
    error.value = 'Failed to create note. Please try again.';
    console.error('Failed to create note:', err);
  }
};
</script>
