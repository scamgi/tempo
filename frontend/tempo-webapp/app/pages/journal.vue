<template>
  <div>
    <h1 class="text-3xl font-bold mb-6 text-gray-800 dark:text-white">Journal</h1>
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
      <p class="text-gray-600 dark:text-gray-300">
        Reflect on your day, track your mood, and cherish your memories.
        This is a placeholder page. Functionality to fetch, create, and manage journal entries will be implemented here.
      </p>
       <div class="mt-6">
        <button
          class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700">
          + New Entry
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ['auth']
});
</script>