<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { watchDebounced } from "@vueuse/core";
import type { Column } from "../../components/DataTable/types";
import { DataTable } from "../../components/DataTable";
import RowActions from "../../components/RowActions.vue";
import PageHeader from "../../components/PageHeader.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import Avatar from "../../components/Avatar.vue";
import AnimalStatusBadge from "../../components/animal/AnimalStatusBadge.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";
import { useAnimalStore } from "../../stores/animal";
import { getFirstErrorMessage } from "../../utils/error";

const router = useRouter();
const headerStore = useHeaderStore();
const animalStore = useAnimalStore();
const { success, error: showError } = useToast();

watchDebounced(
  () => headerStore.searchQuery,
  () => {
    animalStore.page = 1;
    animalStore.fetchPaginated({
      page: animalStore.page,
      per_page: animalStore.pageSize,
      search: headerStore.searchQuery || undefined,
    });
  },
  { debounce: 300, maxWait: 1000 },
);

const columns: Column[] = [
  { key: "tag_no", label: "Tag No", sortable: true },
  { key: "species_name", label: "Species" },
  { key: "breed_name", label: "Breed" },
  { key: "gender", label: "Gender" },
  { key: "status", label: "Status", component: AnimalStatusBadge },
  { key: "user", label: "Created By" },
  { key: "birth_date", label: "Birth Date" },
  {
    key: "actions",
    label: "Action",
    component: RowActions,
    tdPosition: "left",
    componentProps: {
      actions: [
        { label: "View", icon: "Eye", onClick: (item: any) => router.push("/animals/" + item.id) },
        { label: "Edit", icon: "Pencil", onClick: (item: any) => router.push("/animals/" + item.id + "/edit") },
        { label: "Delete", icon: "Trash2", onClick: (item: any) => handleDelete(item.id), danger: true },
      ],
    },
  },
];

async function handleDelete(id: number) {
  if (!confirm("Are you sure you want to delete this animal?")) return;
  try {
    await animalStore.remove(id);
    success("Deleted", "Animal record has been deleted");
    await animalStore.fetchPaginated({
      page: animalStore.page,
      per_page: animalStore.pageSize,
      search: headerStore.searchQuery || undefined,
    });
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([
    { label: "Dashboard", to: "/dashboard" },
    { label: "Animals" },
  ]);
  headerStore.setActions([{ label: "Add New", onClick: () => router.push("/animals/new") }]);
  headerStore.setShowSearch(true);
  animalStore.fetchPaginated({
    page: animalStore.page,
    per_page: animalStore.pageSize,
    search: headerStore.searchQuery || undefined,
  });
});
onUnmounted(() => headerStore.clear());
</script>

<template>
  <div>
    <PageHeader
      title="Animals"
      subtitle="Manage livestock and animal records"
    />
    <DataTable
      :columns="columns"
      :items="animalStore.items"
      :loading="animalStore.loading"
      :server-mode="true"
      :total-items="animalStore.totalItems"
      :current-page="animalStore.page"
      :page-size="animalStore.pageSize"
      @update:current-page="animalStore.page = $event; animalStore.fetchPaginated({ page: animalStore.page, per_page: animalStore.pageSize, search: headerStore.searchQuery || undefined })"
      @update:page-size="animalStore.pageSize = $event; animalStore.page = 1; animalStore.fetchPaginated({ page: animalStore.page, per_page: animalStore.pageSize, search: headerStore.searchQuery || undefined })"
    >
      <template #cell-species_name="{ item }">
        {{ item.breed?.species?.name || "-" }}
      </template>
      <template #cell-breed_name="{ item }">
        {{ item.breed?.name || "-" }}
      </template>
      <template #cell-user="{ item }">
        <div class="flex items-center gap-2">
          <Avatar :name="item.user?.name" size="sm" />
          <span class="text-sm text-gray-700">{{ item.user?.name || "—" }}</span>
        </div>
      </template>
      <template #cell-birth_date="{ item }">
        <DateDisplay :value="item.birth_date" />
      </template>
    </DataTable>
  </div>
</template>
