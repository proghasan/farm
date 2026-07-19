import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listSpecies,
  listSpeciesPaginated,
  createSpecies as apiCreateSpecies,
  updateSpecies as apiUpdateSpecies,
  deleteSpecies as apiDeleteSpecies,
  type Species,
} from "../api";

export const useSpeciesStore = defineStore("species", () => {
  const items = ref<(Species & { created_at?: string })[]>([]);
  const allItems = ref<Species[]>([]);
  const loading = ref(false);
  const page = ref(1);
  const pageSize = ref(20);
  const totalItems = ref(0);

  async function fetchPaginated(params?: Record<string, any>) {
    loading.value = true;
    try {
      const result = await listSpeciesPaginated(params);
      items.value = result.data as (Species & { created_at?: string })[];
      totalItems.value = result.total;
      page.value = result.page;
      pageSize.value = result.per_page;
    } finally {
      loading.value = false;
    }
  }

  async function fetchAll(params?: Record<string, any>) {
    allItems.value = await listSpecies(params);
  }

  async function create(data: Partial<Species>) {
    await apiCreateSpecies(data);
  }

  async function update(id: number, data: Partial<Species>) {
    await apiUpdateSpecies(id, data);
  }

  async function remove(id: number) {
    await apiDeleteSpecies(id);
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
