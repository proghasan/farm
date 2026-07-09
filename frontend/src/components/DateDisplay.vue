<script setup lang="ts">
defineProps<{
  value?: string | null;
}>();

function format(v?: string | null): string {
  if (!v) return "-";
  const d = new Date(v);
  if (isNaN(d.getTime())) return "-";
  return d.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
</script>

<template>
  <span class="text-sm text-gray-600">{{ format(value) }}</span>
</template>
