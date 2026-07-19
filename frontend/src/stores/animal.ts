import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listAnimalsPaginated,
  getAnimal,
  getAnimalProfile,
  createAnimal as apiCreateAnimal,
  updateAnimal as apiUpdateAnimal,
  deleteAnimal as apiDeleteAnimal,
  type Animal,
  type Pregnancy,
} from "../api";

export const useAnimalStore = defineStore("animal", () => {
  const items = ref<Animal[]>([]);
  const currentAnimal = ref<Animal | null>(null);
  const currentPregnancies = ref<Pregnancy[]>([]);
  const loading = ref(false);
  const page = ref(1);
  const pageSize = ref(20);
  const totalItems = ref(0);
  const search = ref("");

  let lastParams = "";

  async function fetchPaginated(params?: Record<string, any>) {
    const key = JSON.stringify(params || {});
    if (key === lastParams && items.value.length) return;
    lastParams = key;

    loading.value = true;
    try {
      const result = await listAnimalsPaginated(params);
      items.value = result.data;
      totalItems.value = result.total;
      page.value = result.page;
      pageSize.value = result.per_page;
    } finally {
      loading.value = false;
    }
  }

  async function refresh(params?: Record<string, any>) {
    lastParams = "";
    await fetchPaginated(params);
  }

  async function fetchById(id: number) {
    loading.value = true;
    try {
      currentAnimal.value = await getAnimal(id);
    } finally {
      loading.value = false;
    }
  }

  async function fetchProfile(id: number) {
    loading.value = true;
    try {
      const profile = await getAnimalProfile(id);
      currentAnimal.value = profile.animal;
      currentPregnancies.value = profile.pregnancies;
    } finally {
      loading.value = false;
    }
  }

  async function create(data: Partial<Animal>) {
    await apiCreateAnimal(data);
  }

  async function update(id: number, data: Partial<Animal>) {
    await apiUpdateAnimal(id, data);
  }

  async function remove(id: number) {
    await apiDeleteAnimal(id);
  }

  return {
    items,
    currentAnimal,
    currentPregnancies,
    loading,
    page,
    pageSize,
    totalItems,
    search,
    fetchPaginated,
    refresh,
    fetchById,
    fetchProfile,
    create,
    update,
    remove,
  };
});
