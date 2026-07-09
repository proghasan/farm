<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";

defineProps<{
  show: boolean;
  title?: string;
}>();

const emit = defineEmits<{ close: [] }>();

function onKeydown(e: KeyboardEvent) {
  if (e.key === "Escape") emit("close");
}

onMounted(() => document.addEventListener("keydown", onKeydown));
onUnmounted(() => document.removeEventListener("keydown", onKeydown));
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer-backdrop">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex justify-end"
        @click.self="emit('close')"
      >
        <div class="absolute inset-0 bg-gray-950/40 backdrop-blur-sm" />

        <Transition name="drawer-panel">
          <div
            v-if="show"
            class="relative bg-white shadow-2xl z-10 w-full max-w-md flex flex-col h-full"
          >
            <div class="flex items-center justify-between px-6 py-5 border-b border-gray-100 shrink-0">
              <h2 class="text-base font-semibold text-gray-900">{{ title }}</h2>
              <button
                @click="emit('close')"
                class="w-8 h-8 flex items-center justify-center rounded-xl hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="px-6 py-5 overflow-y-auto flex-1">
              <slot />
            </div>

            <div
              v-if="$slots.footer"
              class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex items-center justify-end gap-3 shrink-0"
            >
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-backdrop-enter-active,
.drawer-backdrop-leave-active {
  transition: opacity 0.2s ease;
}
.drawer-backdrop-enter-from,
.drawer-backdrop-leave-to {
  opacity: 0;
}
.drawer-panel-enter-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.drawer-panel-leave-active {
  transition: all 0.15s ease-in;
}
.drawer-panel-enter-from {
  transform: translateX(100%);
}
.drawer-panel-leave-to {
  transform: translateX(100%);
}
</style>
