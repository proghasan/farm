import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listBreeds,
  listBreedsPaginated,
  createBreed as apiCreateBreed,
  updateBreed as apiUpdateBreed,
  deleteBreed as apiDeleteBreed,
  type Breed,
} from "../api";

export const useBreedStore = defineStore("breed", () => {
  const items = ref<(Breed & { created_at?: string })[]>([]);
  const allItems = ref<Breed[]>([]);
  const loading = ref(false);
  const page = ref(1);
  const pageSize = ref(20);
  const totalItems = ref(0);
  const search = ref("");

  let lastPaginatedParams = "";
  let lastAllParams = "";

  async function fetchPaginated(params?: Record<string, any>) {
    const key = JSON.stringify(params || {});
    if (key === lastPaginatedParams && items.value.length) return;
    lastPaginatedParams = key;

    loading.value = true;
    try {
      const result = await listBreedsPaginated(params);
      items.value = result.data as (Breed & { created_at?: string })[];
      totalItems.value = result.total;
      page.value = result.page;
      pageSize.value = result.per_page;
    } finally {
      loading.value = false;
    }
  }

  async function fetchAll(params?: Record<string, any>) {
    const key = JSON.stringify(params || {});
    if (key === lastAllParams && allItems.value.length) return;
    lastAllParams = key;

    allItems.value = await listBreeds(params);
  }

  async function create(data: Partial<Breed>) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiCreateBreed(data);
  }

  async function update(id: number, data: Partial<Breed>) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiUpdateBreed(id, data);
  }

  async function remove(id: number) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiDeleteBreed(id);
  }

  return {
    items,
    allItems,
    loading,
    page,
    pageSize,
    totalItems,
    search,
    fetchPaginated,
    fetchAll,
    create,
    update,
    remove,
  };
});
