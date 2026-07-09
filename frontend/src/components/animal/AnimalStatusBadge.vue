<script setup lang="ts">
defineProps<{
  value?: string;
}>();

const colors: Record<string, string> = {
  Healthy: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  Pregnant: "bg-blue-50 text-blue-700 ring-blue-200",
  Sick: "bg-amber-50 text-amber-700 ring-amber-200",
  Sold: "bg-gray-50 text-gray-600 ring-gray-200",
  Dead: "bg-red-50 text-red-700 ring-red-200",
};
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium ring-1"
    :class="colors[value ?? ''] || 'bg-gray-50 text-gray-600 ring-gray-200'"
  >
    <span
      class="w-1.5 h-1.5 rounded-full"
      :class="{
        'bg-emerald-500': value === 'Healthy',
        'bg-blue-500': value === 'Pregnant',
        'bg-amber-500': value === 'Sick',
        'bg-gray-400': value === 'Sold',
        'bg-red-500': value === 'Dead',
      }"
    />
    {{ value }}
  </span>
</template>
