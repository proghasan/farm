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

  async function fetchPaginated(params?: Record<string, any>) {
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
    allItems.value = await listBreeds(params);
  }

  async function create(data: Partial<Breed>) {
    await apiCreateBreed(data);
  }

  async function update(id: number, data: Partial<Breed>) {
    await apiUpdateBreed(id, data);
  }

  async function remove(id: number) {
    await apiDeleteBreed(id);
  }

  return {
    items,
    allItems,
    loading,
    page,
    pageSize,
    totalItems,
    fetchPaginated,
    fetchAll,
    create,
    update,
    remove,
  };
});
