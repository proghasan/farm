<script setup lang="ts">
defineProps<{
  value?: string;
}>();

const colors: Record<string, string> = {
  Active: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  Healthy: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  Sick: "bg-amber-50 text-amber-700 ring-amber-200",
  Deceased: "bg-red-50 text-red-700 ring-red-200",
  Sold: "bg-blue-50 text-blue-700 ring-blue-200",
  Mated: "bg-purple-50 text-purple-700 ring-purple-200",
  Pregnant: "bg-blue-50 text-blue-700 ring-blue-200",
  Delivered: "bg-emerald-50 text-emerald-700 ring-emerald-200",
  Aborted: "bg-red-50 text-red-700 ring-red-200",
  Miscarriage: "bg-rose-50 text-rose-700 ring-rose-200",
  Failed: "bg-gray-50 text-gray-600 ring-gray-200",
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
        'bg-emerald-500':
          value === 'Active' || value === 'Healthy' || value === 'Delivered',
        'bg-red-500': value === 'Deceased' || value === 'Aborted',
        'bg-amber-500': value === 'Sick',
        'bg-blue-500': value === 'Sold' || value === 'Pregnant',
        'bg-purple-500': value === 'Mated',
        'bg-rose-500': value === 'Miscarriage',
        'bg-gray-400': value === 'Failed',
      }"
    />
    {{ value }}
  </span>
</template>
