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
  const search = ref("");

  let lastPaginatedParams = "";
  let lastAllParams = "";

  async function fetchPaginated(params?: Record<string, any>) {
    const key = JSON.stringify(params || {});
    if (key === lastPaginatedParams && items.value.length) return;
    lastPaginatedParams = key;

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
    const key = JSON.stringify(params || {});
    if (key === lastAllParams && allItems.value.length) return;
    lastAllParams = key;

    allItems.value = await listSpecies(params);
  }

  async function create(data: Partial<Species>) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiCreateSpecies(data);
  }

  async function update(id: number, data: Partial<Species>) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiUpdateSpecies(id, data);
  }

  async function remove(id: number) {
    lastPaginatedParams = "";
    lastAllParams = "";
    await apiDeleteSpecies(id);
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
