<script setup lang="ts">
defineProps<{
  value?: string;
}>();
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium ring-1 bg-gray-50 text-gray-600 ring-gray-200"
  >
    {{ value }}
  </span>
</template>
