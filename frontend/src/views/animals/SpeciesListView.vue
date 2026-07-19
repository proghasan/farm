<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { watchDebounced } from "@vueuse/core";
import type { Column } from "../../components/DataTable/types";
import { DataTable } from "../../components/DataTable";
import RowActions from "../../components/RowActions.vue";
import Drawer from "../../components/Drawer.vue";
import PageHeader from "../../components/PageHeader.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import Avatar from "../../components/Avatar.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";
import { useSpeciesStore } from "../../stores/species";
import { getFirstErrorMessage } from "../../utils/error";

const { success, error: showError } = useToast();
const headerStore = useHeaderStore();
const speciesStore = useSpeciesStore();
const showModal = ref(false);
const editingId = ref<number | null>(null);
const form = ref({ name: "" });
const saving = ref(false);

watchDebounced(
  () => headerStore.searchQuery,
  () => {
    speciesStore.page = 1;
    speciesStore.fetchPaginated({
      page: speciesStore.page,
      per_page: speciesStore.pageSize,
      search: headerStore.searchQuery || undefined,
    });
  },
  { debounce: 300, maxWait: 1000 },
);

function openCreate() {
  editingId.value = null;
  form.value = { name: "" };
  showModal.value = true;
}

function openEdit(id: number) {
  const item = speciesStore.items.find((i) => i.id === id);
  if (!item) return;
  editingId.value = id;
  form.value = { name: item.name };
  showModal.value = true;
}

async function save() {
  saving.value = true;
  try {
    if (editingId.value) {
      await speciesStore.update(editingId.value, form.value);
      success("Updated", "Species has been updated");
    } else {
      await speciesStore.create(form.value);
      success("Created", "Species has been created");
    }
    showModal.value = false;
    await speciesStore.fetchPaginated({
      page: speciesStore.page,
      per_page: speciesStore.pageSize,
      search: headerStore.searchQuery || undefined,
    });
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  if (!confirm("Are you sure you want to delete this species?")) return;
  try {
    await speciesStore.remove(id);
    success("Deleted", "Species has been deleted");
    await speciesStore.fetchPaginated({
      page: speciesStore.page,
      per_page: speciesStore.pageSize,
      search: headerStore.searchQuery || undefined,
    });
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  }
}

const columns: Column[] = [
  { key: "name", label: "Name", sortable: true },
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
    { label: "Species" },
  ]);
  headerStore.setActions([{ label: "Add New", onClick: openCreate }]);
  headerStore.setShowSearch(true);
  speciesStore.fetchPaginated({
    page: speciesStore.page,
    per_page: speciesStore.pageSize,
    search: headerStore.searchQuery || undefined,
  });
});
onUnmounted(() => headerStore.clear());
</script>

<template>
  <div>
    <PageHeader title="Species" subtitle="Manage animal species records" />
    <DataTable
      :columns="columns"
      :items="speciesStore.items"
      :loading="speciesStore.loading"
      :server-mode="true"
      :total-items="speciesStore.totalItems"
      :current-page="speciesStore.page"
      :page-size="speciesStore.pageSize"
      @update:current-page="
        speciesStore.page = $event;
        speciesStore.fetchPaginated({ page: speciesStore.page, per_page: speciesStore.pageSize, search: headerStore.searchQuery || undefined });
      "
      @update:page-size="
        speciesStore.pageSize = $event;
        speciesStore.page = 1;
        speciesStore.fetchPaginated({ page: speciesStore.page, per_page: speciesStore.pageSize, search: headerStore.searchQuery || undefined });
      "
    >
      <template #cell-user="{ item }">
        <div class="flex items-center gap-2">
          <Avatar :name="item.user?.name" size="sm" />
          <span class="text-sm text-gray-700">{{ item.user?.name || '—' }}</span>
        </div>
      </template>
      <template #cell-created_at="{ item }">
        <DateDisplay :value="item.created_at" />
      </template>
    </DataTable>

    <Drawer
      :show="showModal"
      :title="editingId ? 'Edit Species' : 'Add Species'"
      @close="showModal = false"
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
