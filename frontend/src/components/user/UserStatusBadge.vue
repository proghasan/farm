<script setup lang="ts">
defineProps<{
  value?: string;
}>();

const colors: Record<string, string> = {
  Active: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  Inactive: "bg-gray-50 text-gray-600 ring-gray-200",
  Suspended: "bg-amber-50 text-amber-700 ring-amber-200",
  Owner: "bg-purple-50 text-purple-700 ring-purple-200",
  Manager: "bg-blue-50 text-blue-700 ring-blue-200",
  Veterinarian: "bg-cyan-50 text-cyan-700 ring-cyan-200",
  Worker: "bg-gray-50 text-gray-600 ring-gray-200",
  Accountant: "bg-emerald-50 text-emerald-700 ring-emerald-200",
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
        'bg-emerald-500': value === 'Active' || value === 'Accountant',
        'bg-gray-400': value === 'Inactive' || value === 'Worker',
        'bg-amber-500': value === 'Suspended',
        'bg-purple-500': value === 'Owner',
        'bg-blue-500': value === 'Manager',
        'bg-cyan-500': value === 'Veterinarian',
      }"
    />
    {{ value }}
  </span>
</template>
