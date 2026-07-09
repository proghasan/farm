<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import * as LucideIcons from "@lucide/vue";

export interface Action<T = any> {
  label: string;
  onClick: (item: T) => void;
  icon?: string;
  danger?: boolean;
}

const props = defineProps<{
  item: Record<string, any>;
  actions: Action[];
}>();

const open = ref(false);
const menuStyle = ref({ top: "0px", right: "0px" });
const triggerRef = ref<HTMLElement | null>(null);
const containerRef = ref<HTMLElement | null>(null);

function toggle() {
  if (open.value) {
    open.value = false;
    return;
  }
  const btn = triggerRef.value;
  if (!btn) return;
  const rect = btn.getBoundingClientRect();
  menuStyle.value = {
    top: `${rect.bottom + 4}px`,
    right: `${window.innerWidth - rect.right}px`,
  };
  open.value = true;
}

function handleClick(action: Action) {
  action.onClick(props.item);
  open.value = false;
}

function close(e: MouseEvent) {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
    open.value = false;
  }
}

onMounted(() => document.addEventListener("click", close));
onUnmounted(() => document.removeEventListener("click", close));
</script>

<template>
  <div ref="containerRef" class="row-actions relative flex justify-end">
    <button
      ref="triggerRef"
      @click="toggle"
      class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors opacity-0 group-hover:opacity-100"
      :class="open && 'opacity-100 bg-gray-100 text-gray-600'"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="5" r="1.5" />
        <circle cx="12" cy="12" r="1.5" />
        <circle cx="12" cy="19" r="1.5" />
      </svg>
    </button>

    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="open"
          :style="menuStyle"
          class="fixed w-44 bg-white rounded-xl shadow-lg border border-gray-100 overflow-hidden z-50"
        >
          <button
            v-for="action in actions"
            :key="action.label"
            @click="handleClick(action)"
            class="w-full flex items-center gap-2.5 px-3.5 py-2.5 text-sm transition-colors"
            :class="
              action.danger
                ? 'text-red-500 hover:bg-red-50'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            "
          >
            <component
              :is="(LucideIcons as any)[action.icon!]"
              v-if="action.icon"
              class="w-4 h-4"
            />
            {{ action.label }}
          </button>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.dropdown-enter-active {
  transition: all 0.12s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.dropdown-leave-active {
  transition: all 0.08s ease-in;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
