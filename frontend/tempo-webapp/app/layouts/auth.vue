<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
    <div class="max-w-md w-full p-6 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
      <h2 class="text-2xl font-bold text-center text-gray-800 dark:text-white mb-6">{{ pageTitle }}</h2>
      <slot />
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const pageTitle = computed(() => {
  if (route.path.includes('login')) return 'Welcome Back!';
  if (route.path.includes('register')) return 'Create an Account';
  return 'Tempo';
});
</script>