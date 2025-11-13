<template>
  <div>
    <!-- Loading State -->
    <div v-if="pending && !entryData" class="text-gray-500 dark:text-gray-400">
      Loading entry...
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-4 bg-red-100 text-red-700 rounded-md">
      <p>Could not load the journal entry. It may have been deleted or you don't have permission to view it.</p>
      <NuxtLink to="/journal" class="text-blue-600 hover:underline mt-2 inline-block">Back to Journal</NuxtLink>
    </div>

    <!-- Main Content -->
    <div v-else-if="entryData">
      <div class="flex justify-between items-center mb-6">
        <NuxtLink to="/journal" class="text-blue-600 hover:underline">&larr; Back to Journal</NuxtLink>
        <div>
          <button @click="handleUpdateEntry" class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 mr-2">
            Save Changes
          </button>
          <button @click="confirmDeleteEntry" class="py-2 px-4 border border-red-500 text-red-500 rounded-md shadow-sm text-sm font-medium hover:bg-red-500 hover:text-white">
            Delete Entry
          </button>
        </div>
      </div>

      <!-- Edit Entry Form -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <form class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="md:col-span-2">
              <label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Title</label>
              <input
                id="title"
                type="text"
                v-model="editableTitle"
                class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
             <div>
              <label for="date" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Date</label>
              <input
                id="date"
                type="date"
                v-model="editableEntryDate"
                class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>
          <div>
            <label for="content" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Content</label>
            <textarea
              id="content"
              v-model="editableContent"
              rows="15"
              class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            ></textarea>
          </div>
           <div>
            <label for="mood" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Mood (Optional)</label>
            <input
              id="mood"
              type="text"
              v-model="editableMood"
              class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </form>
      </div>
       <p v-if="successMessage" class="mt-4 text-green-600">{{ successMessage }}</p>
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
const entryId = route.params.entryId;

const { data: entryData, pending, error } = await useApiFetch(`/journal/${entryId}`);

const editableTitle = ref('');
const editableContent = ref('');
const editableMood = ref('');
const editableEntryDate = ref('');
const successMessage = ref('');

// Watch for entryData to become available and then populate the editable refs
watch(entryData, (newData) => {
  if (newData) {
    editableTitle.value = newData.title;
    editableContent.value = newData.content;
    editableMood.value = newData.mood || '';
    // Format the date for the date input, which expects 'YYYY-MM-DD'
    editableEntryDate.value = new Date(newData.entryDate).toISOString().split('T')[0];
  }
}, { immediate: true });


const headers = { Authorization: `Bearer ${authStore.token}` };

// Update the entry
const handleUpdateEntry = async () => {
  if (!editableTitle.value.trim() || !editableEntryDate.value) {
    alert('Title and date cannot be empty.');
    return;
  }
  try {
    await $fetch(`/journal/${entryId}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: {
        title: editableTitle.value,
        content: editableContent.value,
        mood: editableMood.value || null,
        entryDate: new Date(editableEntryDate.value).toISOString(),
      },
      headers,
    });
    successMessage.value = 'Entry updated successfully!';
    setTimeout(() => {
      successMessage.value = '';
    }, 3000);
  } catch (err) {
    alert('Failed to update entry. Please try again.');
  }
};

// Delete the entry
const confirmDeleteEntry = async () => {
  if (confirm(`Are you sure you want to delete this journal entry? This cannot be undone.`)) {
    try {
      await $fetch(`/journal/${entryId}`, {
        baseURL: config.public.apiBase,
        method: 'DELETE',
        headers,
      });
      await router.push('/journal');
    } catch (err) {
      alert('Failed to delete the entry. Please try again.');
    }
  }
}
</script>