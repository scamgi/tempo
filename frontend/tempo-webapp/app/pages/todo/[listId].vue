<template>
  <div>
    <!-- Loading State -->
    <div v-if="pending && !listData" class="text-gray-500 dark:text-gray-400">
      Loading list...
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-4 bg-red-100 text-red-700 rounded-md">
      <p>Could not load the to-do list. It may have been deleted or you don't have permission to view it.</p>
      <NuxtLink to="/todo" class="text-blue-600 hover:underline mt-2 inline-block">Back to all lists</NuxtLink>
    </div>

    <!-- Main Content -->
    <div v-else-if="listData">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold text-gray-800 dark:text-white">{{ listData.title }}</h1>
        <button @click="confirmDeleteList" class="py-2 px-4 border border-red-500 text-red-500 rounded-md shadow-sm text-sm font-medium hover:bg-red-500 hover:text-white">
          Delete List
        </button>
      </div>

      <!-- Add New Item Form -->
      <div class="mb-6">
        <form @submit.prevent="addItem" class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md flex items-center gap-4">
          <input
            type="text"
            v-model="newItemTask"
            placeholder="Add a new task..."
            class="flex-grow px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
          <button
            type="submit"
            :disabled="!newItemTask.trim()"
            class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed"
          >
            Add Task
          </button>
        </form>
      </div>

      <!-- Items List -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <ul v-if="listData.items && listData.items.length" class="space-y-3">
          <li v-for="item in listData.items" :key="item.id" class="flex items-center justify-between p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700">
            <div class="flex items-center">
              <input
                type="checkbox"
                :checked="item.isCompleted"
                @change="toggleItemCompletion(item)"
                class="h-5 w-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <span :class="['ml-3', { 'line-through text-gray-500': item.isCompleted }]">
                {{ item.task }}
              </span>
            </div>
            <button @click="deleteItem(item.id)" class="text-gray-400 hover:text-red-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </button>
          </li>
        </ul>
        <p v-else class="text-gray-500 dark:text-gray-400 text-center">No tasks in this list yet. Add one above!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  middleware: ['auth']
});

const route = useRoute();
const router = useRouter();
const config = useRuntimeConfig();
const authStore = useAuthStore();
const listId = route.params.listId;

const { data: listData, pending, error, refresh } = await useApiFetch(`/lists/${listId}`);

const newItemTask = ref('');
const headers = { Authorization: `Bearer ${authStore.token}` };

// Add a new item
const addItem = async () => {
  if (!newItemTask.value.trim()) return;
  await $fetch(`/lists/${listId}/items`, {
    baseURL: config.public.apiBase,
    method: 'POST',
    body: { task: newItemTask.value },
    headers,
  });
  newItemTask.value = '';
  await refresh();
};

// Toggle item completion
const toggleItemCompletion = async (item) => {
  await $fetch(`/items/${item.id}`, {
    baseURL: config.public.apiBase,
    method: 'PUT',
    body: { isCompleted: !item.isCompleted },
    headers,
  });
  await refresh();
};

// Delete an item
const deleteItem = async (itemId) => {
  await $fetch(`/items/${itemId}`, {
    baseURL: config.public.apiBase,
    method: 'DELETE',
    headers,
  });
  await refresh();
};

// Delete the entire list
const confirmDeleteList = async () => {
  if (confirm(`Are you sure you want to delete the list "${listData.value.title}"? This cannot be undone.`)) {
    try {
      await $fetch(`/lists/${listId}`, {
        baseURL: config.public.apiBase,
        method: 'DELETE',
        headers,
      });
      await router.push('/todo');
    } catch (err) {
      alert('Failed to delete the list. Please try again.');
    }
  }
}
</script>