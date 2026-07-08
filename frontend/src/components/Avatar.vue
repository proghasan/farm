<script setup lang="ts">
const props = withDefaults(defineProps<{
  name?: string
  size?: 'sm' | 'md' | 'lg'
}>(), { size: 'md' })

const initials = (props.name || '')
  .split(' ')
  .filter(Boolean)
  .slice(0, 2)
  .map(w => w[0].toUpperCase())
  .join('')

const sizeClasses: Record<string, string> = {
  sm: 'w-7 h-7 text-[10px]',
  md: 'w-9 h-9 text-xs',
  lg: 'w-10 h-10 text-sm',
}
</script>

<template>
  <div
    class="rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold shrink-0"
    :class="sizeClasses[size]"
  >
    {{ initials || '?' }}
  </div>
</template>