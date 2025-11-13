<template>
  <div>
    <!-- Loading State -->
    <div v-if="pending && !noteData" class="text-gray-500 dark:text-gray-400">
      Loading note...
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-4 bg-red-100 text-red-700 rounded-md">
      <p>Could not load the note. It may have been deleted or you don't have permission to view it.</p>
      <NuxtLink to="/notes" class="text-blue-600 hover:underline mt-2 inline-block">Back to all notes</NuxtLink>
    </div>

    <!-- Main Content -->
    <div v-else-if="noteData">
      <div class="flex justify-between items-center mb-6">
        <NuxtLink to="/notes" class="text-blue-600 hover:underline">&larr; Back to Notes</NuxtLink>
        <div>
          <button @click="handleUpdateNote" class="py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 mr-2">
            Save Changes
          </button>
          <button @click="confirmDeleteNote" class="py-2 px-4 border border-red-500 text-red-500 rounded-md shadow-sm text-sm font-medium hover:bg-red-500 hover:text-white">
            Delete Note
          </button>
        </div>
      </div>

      <!-- Edit Note Form -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <form class="space-y-4">
          <div>
            <label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-200">Title</label>
            <input
              id="title"
              type="text"
              v-model="editableTitle"
              class="mt-1 block w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
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
const noteId = route.params.noteId;

const { data: noteData, pending, error } = await useApiFetch(`/notes/${noteId}`);

const editableTitle = ref('');
const editableContent = ref('');
const successMessage = ref('');

// Watch for noteData to become available and then populate the editable refs
watch(noteData, (newData) => {
  if (newData) {
    editableTitle.value = newData.title;
    editableContent.value = newData.content;
  }
}, { immediate: true });


const headers = { Authorization: `Bearer ${authStore.token}` };

// Update the note
const handleUpdateNote = async () => {
  if (!editableTitle.value.trim()) {
    alert('Title cannot be empty.');
    return;
  }
  try {
    await $fetch(`/notes/${noteId}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: {
        title: editableTitle.value,
        content: editableContent.value,
      },
      headers,
    });
    successMessage.value = 'Note updated successfully!';
    // Hide the message after 3 seconds
    setTimeout(() => {
      successMessage.value = '';
    }, 3000);
  } catch (err) {
    alert('Failed to update note. Please try again.');
  }
};

// Delete the note
const confirmDeleteNote = async () => {
  if (confirm(`Are you sure you want to delete this note? This cannot be undone.`)) {
    try {
      await $fetch(`/notes/${noteId}`, {
        baseURL: config.public.apiBase,
        method: 'DELETE',
        headers,
      });
      await router.push('/notes');
    } catch (err) {
      alert('Failed to delete the note. Please try again.');
    }
  }
}
</script>