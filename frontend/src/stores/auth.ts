import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { login as apiLogin, getProfile, type User } from "../api";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem("token"));
  const user = ref<User | null>(null);
  const loading = ref(false);

  const isAuthenticated = computed(() => !!token.value);
  const isAdmin = computed(
    () => user.value?.role === "Owner" || user.value?.role === "Manager",
  );

  async function login(username: string, password: string) {
    loading.value = true;
    try {
      const res = await apiLogin({ login: username, password });
      token.value = res.token;
      user.value = res.user;
      localStorage.setItem("token", res.token);
    } finally {
      loading.value = false;
    }
  }

  async function fetchProfile() {
    try {
      user.value = await getProfile();
    } catch {
      logout();
    }
  }

  function logout() {
    token.value = null;
    user.value = null;
    localStorage.removeItem("token");
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    isAdmin,
    login,
    fetchProfile,
    logout,
  };
});
