<template>
  <div class="min-h-screen flex flex-col md:flex-row">
    <!-- Mobile Header -->
    <header class="md:hidden bg-white dark:bg-gray-800 shadow-md p-4 flex justify-between items-center">
      <h1 class="text-xl font-bold text-gray-800 dark:text-white">Tempo</h1>
      <button @click="sidebarOpen = !sidebarOpen" class="text-gray-500 focus:outline-none">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path>
        </svg>
      </button>
    </header>

    <!-- Sidebar -->
    <aside
      :class="['bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 w-64 space-y-6 py-7 px-2 absolute inset-y-0 left-0 transform md:relative md:translate-x-0 transition duration-200 ease-in-out', { '-translate-x-full': !sidebarOpen }]">
      <h1 class="text-2xl font-bold text-center text-gray-800 dark:text-white hidden md:block">Tempo</h1>
      <nav>
        <NuxtLink to="/todo"
          class="block py-2.5 px-4 rounded transition duration-200 hover:bg-gray-100 dark:hover:bg-gray-700"
          active-class="bg-blue-500 text-white hover:bg-blue-600 dark:hover:bg-blue-600">
          To-Do Lists
        </NuxtLink>
        <NuxtLink to="/notes"
          class="block py-2.5 px-4 rounded transition duration-200 hover:bg-gray-100 dark:hover:bg-gray-700"
          active-class="bg-blue-500 text-white hover:bg-blue-600 dark:hover:bg-blue-600">
          Notes
        </NuxtLink>
        <NuxtLink to="/journal"
          class="block py-2.5 px-4 rounded transition duration-200 hover:bg-gray-100 dark:hover:bg-gray-700"
          active-class="bg-blue-500 text-white hover:bg-blue-600 dark:hover:bg-blue-600">
          Journal
        </NuxtLink>
      </nav>
      <div class="absolute bottom-4 px-4 w-full left-0">
         <button @click="handleLogout" class="w-full text-left py-2.5 px-4 rounded transition duration-200 hover:bg-red-100 dark:hover:bg-red-700">
          Logout
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-4 sm:p-6 lg:p-8">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth';
const authStore = useAuthStore();
const router = useRouter();
const sidebarOpen = ref(false);

const handleLogout = async () => {
  await authStore.logout();
  router.push('/login');
};
</script>