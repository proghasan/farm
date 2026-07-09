<script setup lang="ts">
defineProps<{
  placeholder?: string;
}>();

const model = defineModel<string>();
</script>

<template>
  <div
    class="flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-2 w-100 focus-within:ring-2 focus-within:ring-brand-500/20 focus-within:border-brand-400 transition-all"
  >
    <svg
      class="w-4 h-4 text-gray-400 shrink-0"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="2"
      stroke="currentColor"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
      />
    </svg>
    <input
      v-model="model"
      type="text"
      :placeholder="placeholder || 'Search…'"
      class="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none w-full"
    />
  </div>
</template>
