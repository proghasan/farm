<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import Sidebar from "../components/Sidebar.vue";
import Header from "../components/Header.vue";
import Toast from "../components/Toast.vue";

const sidebarOpen = ref(true);
const route = useRoute();

const pageTitle = computed(() => {
  return (route.meta?.title as string) ?? "Dashboard";
});
</script>

<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <Sidebar :open="sidebarOpen" @toggle="sidebarOpen = !sidebarOpen" />

    <div
      class="flex flex-col flex-1 min-w-0 transition-all duration-300"
      :class="sidebarOpen ? 'ml-64' : 'ml-16'"
    >
      <Header
        :title="pageTitle"
        :sidebar-open="sidebarOpen"
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      />

      <main class="flex-1 overflow-y-auto p-6">
        <RouterView />
      </main>
    </div>

    <Toast />
  </div>
</template>
