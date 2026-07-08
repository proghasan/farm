<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import Avatar from "./Avatar.vue";

defineProps<{
  title: string;
  sidebarOpen: boolean;
}>();

defineEmits<{ toggleSidebar: [] }>();

const auth = useAuthStore();
const router = useRouter();

const avatarOpen = ref(false);

const profile = computed(() => auth.user);

function closeAll(e: MouseEvent) {
  const target = e.target as HTMLElement;
  if (!target.closest("#avatar-menu")) avatarOpen.value = false;
}

function logout() {
  auth.logout();
  router.push("/login");
}

onMounted(() => document.addEventListener("click", closeAll));
onUnmounted(() => document.removeEventListener("click", closeAll));
</script>

<template>
  <header
    class="h-16 bg-white border-b border-gray-100 flex items-center justify-between px-6 shrink-0 z-20"
  >
    <div class="flex items-center gap-3">
      <h1 class="text-lg font-semibold text-gray-900">{{ title }}</h1>
      <nav class="hidden sm:flex items-center gap-1 text-sm text-gray-400">
        <span>Farm</span>
        <svg
          class="w-3.5 h-3.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="2"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="m8.25 4.5 7.5 7.5-7.5 7.5"
          />
        </svg>
        <span class="text-gray-700 font-medium">{{ title }}</span>
      </nav>
    </div>

    <div class="flex items-center gap-2">
      <!-- Search -->
      <div
        class="hidden md:flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-2 w-56 focus-within:ring-2 focus-within:ring-brand-500/20 focus-within:border-brand-400 transition-all"
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
          type="text"
          placeholder="Search…"
          class="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none w-full"
        />
        <kbd
          class="text-xs text-gray-400 bg-white border border-gray-200 rounded px-1.5 py-0.5 font-mono"
          >⌘K</kbd
        >
      </div>

      <!-- Settings -->
      <button
        class="w-9 h-9 flex items-center justify-center rounded-xl hover:bg-gray-100 text-gray-500 hover:text-gray-700 transition-colors"
      >
        <svg
          class="w-5 h-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
          />
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
          />
        </svg>
      </button>

      <!-- Avatar + Dropdown -->
      <div id="avatar-menu" class="relative">
        <button
          @click="avatarOpen = !avatarOpen"
          class="cursor-pointer ring-2 rounded-full transition-all"
          :class="
            avatarOpen
              ? 'ring-brand-400'
              : 'ring-transparent hover:ring-brand-300'
          "
        >
          <Avatar v-if="profile" :name="profile.name" />
        </button>

        <Transition name="dropdown">
          <div
            v-if="avatarOpen"
            class="absolute right-0 top-12 w-56 bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden z-50"
          >
            <!-- User info -->
            <div v-if="profile" class="px-4 py-3.5 border-b border-gray-100">
              <div class="flex items-center gap-3">
                <Avatar :name="profile.name" />
                <div class="min-w-0">
                  <p class="text-sm font-semibold text-gray-900 truncate">
                    {{ profile.name }}
                  </p>
                  <p class="text-xs text-gray-400 truncate">
                    {{ profile?.email }}
                  </p>
                </div>
              </div>
              <span
                class="mt-2.5 inline-flex items-center gap-1.5 text-xs font-medium text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-md ring-1 ring-emerald-200"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                {{ profile.role }}
              </span>
            </div>

            <!-- Menu items -->
            <div class="py-1.5">
              <button
                class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-colors"
              >
                <svg
                  class="w-4 h-4 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
                  />
                </svg>
                View Profile
              </button>
              <button
                class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-colors"
              >
                <svg
                  class="w-4 h-4 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                  />
                </svg>
                Settings
              </button>
            </div>

            <!-- Divider + Logout -->
            <div class="border-t border-gray-100 py-1.5">
              <button
                @click="logout"
                class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-500 hover:bg-red-50 transition-colors"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9"
                  />
                </svg>
                Sign Out
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </header>
</template>

<style scoped>
.dropdown-enter-active {
  transition: all 0.15s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.dropdown-leave-active {
  transition: all 0.1s ease-in;
}
.dropdown-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(-6px);
}
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.97) translateY(-4px);
}
</style>
