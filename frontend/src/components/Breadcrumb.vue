<script setup lang="ts">
export interface BreadcrumbItem {
  label: string
  to?: string
}

defineProps<{
  items: BreadcrumbItem[]
}>()
</script>

<template>
  <nav class="flex items-center gap-1 text-sm text-gray-400">
    <template v-for="(item, i) in items" :key="i">
      <svg v-if="i > 0" class="w-3.5 h-3.5 text-gray-300 shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
      </svg>
      <router-link v-if="item.to" :to="item.to" class="hover:text-gray-600 transition-colors">{{ item.label }}</router-link>
      <span v-else class="text-gray-700 font-medium">{{ item.label }}</span>
    </template>
  </nav>
</template>