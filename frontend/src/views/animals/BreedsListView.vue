<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { watchDebounced } from "@vueuse/core";
import {
  listBreedsPaginated,
  createBreed,
  updateBreed,
  deleteBreed,
  listSpecies,
} from "../../api";
import type { Breed, Species } from "../../api";
import type { Column } from "../../components/DataTable/types";
import { DataTable } from "../../components/DataTable";
import RowActions from "../../components/RowActions.vue";
import Drawer from "../../components/Drawer.vue";
import PageHeader from "../../components/PageHeader.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import Avatar from "../../components/Avatar.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";
import { getFirstErrorMessage } from "../../utils/error";

const { success, error: showError } = useToast();
const headerStore = useHeaderStore();
const items = ref<(Breed & { created_at?: string })[]>([]);
const speciesList = ref<Species[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const totalItems = ref(0);
const showModal = ref(false);
const editingId = ref<number | null>(null);
const form = ref({ name: "", species_id: null as number | null });
const saving = ref(false);

async function fetchData() {
  loading.value = true;
  try {
    const result = await listBreedsPaginated({
      page: page.value,
      per_page: pageSize.value,
      search: headerStore.searchQuery || undefined,
    });
    items.value = result.data as (Breed & { created_at?: string })[];
    totalItems.value = result.total;
    speciesList.value = await listSpecies({ all: "true" });
  } finally {
    loading.value = false;
  }
}

watchDebounced(
  () => headerStore.searchQuery,
  () => {
    page.value = 1;
    fetchData();
  },
  { debounce: 300, maxWait: 1000 },
);

function openCreate() {
  editingId.value = null;
  form.value = { name: "", species_id: null };
  showModal.value = true;
}

function openEdit(id: number) {
  const item = items.value.find((i) => i.id === id);
  if (!item) return;
  editingId.value = id;
  form.value = { name: item.name, species_id: item.species_id };
  showModal.value = true;
}

async function save() {
  saving.value = true;
  try {
    const payload = {
      name: form.value.name,
      species_id: form.value.species_id!,
    };
    if (editingId.value) {
      await updateBreed(editingId.value, payload);
      success("Updated", "Breed has been updated");
    } else {
      await createBreed(payload);
      success("Created", "Breed has been created");
    }
    showModal.value = false;
    await fetchData();
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  if (!confirm("Are you sure you want to delete this breed?")) return;
  try {
    await deleteBreed(id);
    success("Deleted", "Breed has been deleted");
    await fetchData();
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  }
}

const columns: Column[] = [
  { key: "name", label: "Name", sortable: true },
  { key: "species_name", label: "Species" },
  { key: "user", label: "Created By" },
  { key: "created_at", label: "Created At" },
  {
    key: "actions",
    label: "Action",
    component: RowActions,
    tdPosition: "left",
    componentProps: {
      actions: [
        {
          label: "Edit",
          icon: "Pencil",
          onClick: (item: any) => openEdit(item.id),
        },
        {
          label: "Delete",
          icon: "Trash2",
          onClick: (item: any) => handleDelete(item.id),
          danger: true,
        },
      ],
    },
  },
];

onMounted(() => {
  headerStore.setBreadcrumb([
    { label: "Dashboard", to: "/dashboard" },
    { label: "Breeds" },
  ]);
  headerStore.setActions([{ label: "Add New", onClick: openCreate }]);
  headerStore.setShowSearch(true);
  fetchData();
});
onUnmounted(() => headerStore.clear());
</script>

<template>
  <div>
    <PageHeader title="Breeds" subtitle="Manage animal breed records" />
    <DataTable
      :columns="columns"
      :items="items"
      :loading="loading"
      :server-mode="true"
      :total-items="totalItems"
      :current-page="page"
      :page-size="pageSize"
      @update:current-page="
        page = $event;
        fetchData();
      "
      @update:page-size="
        pageSize = $event;
        page = 1;
        fetchData();
      "
    >
      <template #cell-species_name="{ item }">
        {{ item.species?.name || "-" }}
      </template>
      <template #cell-user="{ item }">
        <div class="flex items-center gap-2">
          <Avatar :name="item.user?.name" size="sm" />
          <span class="text-sm text-gray-700">{{
            item.user?.name || "—"
          }}</span>
        </div>
      </template>
      <template #cell-created_at="{ item }">
        <DateDisplay :value="item.created_at" />
      </template>
    </DataTable>

    <Drawer
      :show="showModal"
      :title="editingId ? 'Edit Breed' : 'Add Breed'"
      @close="showModal = false"
    >
      <form @submit.prevent="save" class="space-y-2">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Name</label
          >
          <input
            v-model="form.name"
            type="text"
            required
            class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Species</label
          >
          <select
            v-model="form.species_id"
            required
            class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
          >
            <option :value="null" disabled>Select species</option>
            <option v-for="s in speciesList" :key="s.id" :value="s.id">
              {{ s.name }}
            </option>
          </select>
        </div>
      </form>
      <template #footer>
        <button
          @click="save"
          :disabled="saving"
          class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50"
        >
          {{ saving ? "Saving..." : "Save" }}
        </button>
        <button
          @click="showModal = false"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50"
        >
          Cancel
        </button>
      </template>
    </Drawer>
  </div>
</template>
