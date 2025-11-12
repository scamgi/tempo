<template>
  <div>
    <h1 class="text-3xl font-bold mb-6 text-gray-800 dark:text-white">To-Do Lists</h1>
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
      <p class="text-gray-600 dark:text-gray-300">
        Welcome to your To-Do list. Manage your tasks, organize by project, and stay productive.
        This is a placeholder page. Functionality to fetch, create, and manage tasks will be implemented here.
      </p>
      <!-- Example UI elements -->
      <div class="mt-6">
        <input type="text" placeholder="Add a new list..."
          class="w-full md:w-1/2 px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500">
        <button
          class="mt-2 md:mt-0 md:ml-2 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700">
          Add List
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