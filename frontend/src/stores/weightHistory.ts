import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listWeightHistories,
  createWeightHistory as apiCreateWeightHistory,
  updateWeightHistory as apiUpdateWeightHistory,
  deleteWeightHistory as apiDeleteWeightHistory,
  type WeightHistory,
} from "../api";

export const useWeightHistoryStore = defineStore("weightHistory", () => {
  const items = ref<WeightHistory[]>([]);
  const loading = ref(false);

  async function fetchByAnimal(animalId: number) {
    loading.value = true;
    try {
      items.value = await listWeightHistories(animalId);
    } finally {
      loading.value = false;
    }
  }

  async function create(data: Partial<WeightHistory>) {
    await apiCreateWeightHistory(data);
  }

  async function update(id: number, data: Partial<WeightHistory>) {
    await apiUpdateWeightHistory(id, data);
  }

  async function remove(id: number) {
    await apiDeleteWeightHistory(id);
  }

  return { items, loading, fetchByAnimal, create, update, remove };
});
