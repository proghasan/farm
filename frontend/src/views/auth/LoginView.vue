<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../../stores/auth";

const auth = useAuthStore();
const router = useRouter();
const username = ref("");
const password = ref("");
const error = ref("");

async function handleLogin() {
  error.value = "";
  try {
    await auth.login(username.value, password.value);
    await auth.fetchProfile();
    router.push("/dashboard");
  } catch (e: any) {
    error.value = e.response?.data?.error || "Login failed";
  }
}
</script>

<template>
  <div class="w-full max-w-sm">
    <div
      class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden"
    >
      <div class="h-1 bg-brand-600" />
      <div class="p-8">
        <div class="mb-6 text-center">
          <div
            class="mx-auto mb-3 w-12 h-12 rounded-xl bg-brand-600 flex items-center justify-center"
          >
            <svg
              class="w-6 h-6 text-white"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="2"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z"
              />
            </svg>
          </div>
          <h1 class="text-xl font-bold text-gray-900">Farm Manager</h1>
          <p class="text-sm text-gray-400 mt-1">Sign in to your account</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div
            v-if="error"
            class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2"
          >
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Username</label
            >
            <input
              v-model="username"
              type="text"
              required
              class="w-full rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all"
              placeholder="Enter username"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Password</label
            >
            <input
              v-model="password"
              type="password"
              required
              class="w-full rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all"
              placeholder="Enter password"
            />
          </div>

          <button
            type="submit"
            :disabled="auth.loading"
            class="w-full rounded-xl bg-brand-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-brand-700 transition-colors disabled:opacity-50"
          >
            {{ auth.loading ? "Signing in..." : "Sign in" }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
