<script setup lang="ts">
import { useToast } from '../composables/useToast'

const { toasts, remove } = useToast()

const meta = {
  success: {
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" /></svg>`,
    iconClass: 'bg-emerald-100 text-emerald-600',
    barClass: 'bg-emerald-500',
    titleClass: 'text-emerald-900',
  },
  error: {
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" /></svg>`,
    iconClass: 'bg-red-100 text-red-600',
    barClass: 'bg-red-500',
    titleClass: 'text-red-900',
  },
  warning: {
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" /></svg>`,
    iconClass: 'bg-amber-100 text-amber-600',
    barClass: 'bg-amber-500',
    titleClass: 'text-amber-900',
  },
  info: {
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z" /></svg>`,
    iconClass: 'bg-blue-100 text-blue-600',
    barClass: 'bg-blue-500',
    titleClass: 'text-blue-900',
  },
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-5 right-5 z-[100] flex flex-col gap-2.5 w-80 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="pointer-events-auto bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden"
        >
          <div
            class="h-0.5 w-full origin-left"
            :class="meta[toast.type].barClass"
            :style="`animation: shrink ${toast.duration ?? 4000}ms linear forwards`"
          />

          <div class="flex items-start gap-3 px-4 py-3.5">
            <div
              class="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 mt-0.5"
              :class="meta[toast.type].iconClass"
            >
              <span class="w-4 h-4" v-html="meta[toast.type].icon" />
            </div>

            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold" :class="meta[toast.type].titleClass">{{ toast.title }}</p>
              <p v-if="toast.message" class="text-xs text-gray-500 mt-0.5 leading-relaxed">{{ toast.message }}</p>
            </div>

            <button
              @click="remove(toast.id)"
              class="w-6 h-6 flex items-center justify-center rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors shrink-0 -mr-1 mt-0.5"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(24px) scale(0.95);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(24px) scale(0.95);
}
@keyframes shrink {
  from { transform: scaleX(1); }
  to { transform: scaleX(0); }
}
</style>