<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import {
  listSpeciesPaginated,
  createSpecies,
  updateSpecies,
  deleteSpecies,
} from "../../api";
import type { Species } from "../../api";
import type { Column } from "../../components/DataTable/types";
import { DataTable } from "../../components/DataTable";
import RowActions from "../../components/RowActions.vue";
import Modal from "../../components/Modal.vue";
import PageHeader from "../../components/PageHeader.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";

const { success, error: showError } = useToast();
const headerStore = useHeaderStore();
const items = ref<(Species & { created_at?: string })[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const totalItems = ref(0);
const showModal = ref(false);
const editingId = ref<number | null>(null);
const form = ref({ name: "" });
const saving = ref(false);

async function fetchData() {
  loading.value = true;
  try {
    const result = await listSpeciesPaginated({
      page: page.value,
      per_page: pageSize.value,
    });
    items.value = result.data as (Species & { created_at?: string })[];
    totalItems.value = result.total;
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  editingId.value = null;
  form.value = { name: "" };
  showModal.value = true;
}

function openEdit(id: number) {
  const item = items.value.find((i) => i.id === id);
  if (!item) return;
  editingId.value = id;
  form.value = { name: item.name };
  showModal.value = true;
}

async function save() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateSpecies(editingId.value, form.value);
      success("Updated", "Species has been updated");
    } else {
      await createSpecies(form.value);
      success("Created", "Species has been created");
    }
    showModal.value = false;
    await fetchData();
  } catch (e: any) {
    showError("Failed", e?.response?.data?.message || "An error occurred");
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  if (!confirm("Are you sure you want to delete this species?")) return;
  try {
    await deleteSpecies(id);
    success("Deleted", "Species has been deleted");
    await fetchData();
  } catch (e: any) {
    showError("Failed", e?.response?.data?.message || "An error occurred");
  }
}

const columns: Column[] = [
  { key: "name", label: "Name", sortable: true },
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
    { label: "Species" },
  ]);
  headerStore.setActions([{ label: "Add New", onClick: openCreate }]);
  headerStore.setShowSearch(true);
  fetchData();
});
onUnmounted(() => headerStore.clear());
</script>

<template>
  <div>
    <PageHeader title="Species" subtitle="Manage animal species records" />
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
      <template #cell-created_at="{ item }">
        <DateDisplay :value="item.created_at" />
      </template>
    </DataTable>

    <Modal
      :show="showModal"
      :title="editingId ? 'Edit Species' : 'Add Species'"
      @close="showModal = false"
      size="sm"
    >
      <form @submit.prevent="save">
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
      </form>
      <template #footer>
        <button
          @click="showModal = false"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="saving"
          class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50"
        >
          {{ saving ? "Saving..." : "Save" }}
        </button>
      </template>
    </Modal>
  </div>
</template>
