<script setup>
// This page just redirects to the main 'todo' page on load.
definePageMeta({
  middleware: ['auth']
})
onMounted(() => {
  navigateTo('/todo')
})
</script>