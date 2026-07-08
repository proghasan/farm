<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

defineProps<{
  show: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg'
}>()

const emit = defineEmits<{ close: [] }>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-backdrop">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="emit('close')"
      >
        <div class="absolute inset-0 bg-gray-950/40 backdrop-blur-sm" @click="emit('close')" />

        <Transition name="modal-panel">
          <div
            v-if="show"
            class="relative bg-white rounded-2xl shadow-2xl w-full z-10 overflow-hidden"
            :class="{
              'max-w-sm': size === 'sm',
              'max-w-md': !size || size === 'md',
              'max-w-2xl': size === 'lg',
            }"
          >
            <div class="h-1 bg-brand-600" />

            <div class="flex items-center justify-between px-6 py-5 border-b border-gray-100">
              <div>
                <h2 class="text-base font-semibold text-gray-900">{{ title }}</h2>
                <slot name="subtitle">
                  <p class="text-sm text-gray-400 mt-0.5">Fill in the details below</p>
                </slot>
              </div>
              <button
                @click="emit('close')"
                class="w-8 h-8 flex items-center justify-center rounded-xl hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="px-6 py-5 max-h-[60vh] overflow-y-auto">
              <slot />
            </div>

            <div v-if="$slots.footer" class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex items-center justify-end gap-3">
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-backdrop-enter-active,
.modal-backdrop-leave-active {
  transition: opacity 0.2s ease;
}
.modal-backdrop-enter-from,
.modal-backdrop-leave-to {
  opacity: 0;
}
.modal-panel-enter-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.modal-panel-leave-active {
  transition: all 0.15s ease-in;
}
.modal-panel-enter-from {
  opacity: 0;
  transform: scale(0.92) translateY(8px);
}
.modal-panel-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(4px);
}
</style>