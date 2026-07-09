<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import {
  listAnimalsPaginated,
  createAnimal,
  updateAnimal,
  deleteAnimal,
  listSpecies,
  listBreeds,
  listWeightHistories,
  createWeightHistory,
  deleteWeightHistory,
} from "../../api";
import type { Animal, Species, Breed, WeightHistory } from "../../api";
import { DataTable } from "../../components/DataTable";
import RowActions from "../../components/RowActions.vue";
import Modal from "../../components/Modal.vue";
import PageHeader from "../../components/PageHeader.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import AnimalStatusBadge from "../../components/animal/AnimalStatusBadge.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";

const router = useRouter();
const headerStore = useHeaderStore();
const { success, error: showError } = useToast();
const items = ref<Animal[]>([]);
const speciesList = ref<Species[]>([]);
const breedsList = ref<Breed[]>([]);
const weightHistories = ref<WeightHistory[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const totalItems = ref(0);
const showModal = ref(false);
const editingId = ref<number | null>(null);
const saving = ref(false);

const form = ref({
  tag_no: "",
  species_id: null as number | null,
  breed_id: null as number | null,
  gender: "Male",
  birth_date: "",
  purchase_date: "",
  purchase_price: 0,
  color: "",
  status: "Healthy",
  remarks: "",
});

const weightForm = ref({ weight: 0, record_date: "", remarks: "" });
const showWeightForm = ref(false);
const weightSaving = ref(false);

const filteredBreeds = computed(() =>
  breedsList.value.filter((b) => b.species_id === form.value.species_id),
);

const weightHistoriesForAnimal = computed(() =>
  weightHistories.value.filter((w) => w.animal_id === editingId.value),
);

async function fetchData() {
  loading.value = true;
  try {
    const [result, species, breeds, weights] = await Promise.all([
      listAnimalsPaginated({ page: page.value, per_page: pageSize.value }),
      listSpecies(),
      listBreeds(),
      listWeightHistories(),
    ]);
    items.value = result.data;
    totalItems.value = result.total;
    speciesList.value = species;
    breedsList.value = breeds;
    weightHistories.value = weights;
  } finally {
    loading.value = false;
  }
}

const columns = [
  { key: "tag_no", label: "Tag No", sortable: true },
  { key: "species_name", label: "Species" },
  { key: "gender", label: "Gender" },
  { key: "status", label: "Status", component: AnimalStatusBadge },
  { key: "birth_date", label: "Birth Date" },
  {
    key: "actions",
    label: "Action",
    component: RowActions,
    componentProps: {
      actions: [
        { label: "View", icon: "Eye", onClick: (item: any) => router.push("/animals/" + item.id) },
        { label: "Edit", icon: "Pencil", onClick: (item: any) => openEdit(item.id) },
        { label: "Delete", icon: "Trash2", onClick: (item: any) => handleDelete(item.id), danger: true },
      ],
    },
  },
];

function openCreate() {
  editingId.value = null;
  form.value = {
    tag_no: "",
    species_id: null,
    breed_id: null,
    gender: "Male",
    birth_date: "",
    purchase_date: "",
    purchase_price: 0,
    color: "",
    status: "Healthy",
    remarks: "",
  };
  showWeightForm.value = false;
  showModal.value = true;
}

function openEdit(id: number) {
  const item = items.value.find((i) => i.id === id);
  if (!item) return;
  editingId.value = id;
  form.value = {
    tag_no: item.tag_no,
    species_id: item.species_id,
    breed_id: item.breed_id ?? null,
    gender: item.gender,
    birth_date: item.birth_date || "",
    purchase_date: item.purchase_date || "",
    purchase_price: item.purchase_price || 0,
    color: item.color || "",
    status: item.status,
    remarks: item.remarks || "",
  };
  showWeightForm.value = false;
  showModal.value = true;
}

async function save() {
  saving.value = true;
  try {
    const payload: Record<string, any> = {
      tag_no: form.value.tag_no,
      species_id: form.value.species_id!,
      gender: form.value.gender,
      status: form.value.status,
      purchase_price: form.value.purchase_price || 0,
    };
    if (form.value.breed_id) payload.breed_id = form.value.breed_id;
    if (form.value.birth_date) payload.birth_date = form.value.birth_date;
    if (form.value.purchase_date)
      payload.purchase_date = form.value.purchase_date;
    if (form.value.color) payload.color = form.value.color;
    if (form.value.remarks) payload.remarks = form.value.remarks;

    if (editingId.value) {
      await updateAnimal(editingId.value, payload as any);
      success("Updated", "Animal record has been updated");
    } else {
      await createAnimal(payload as any);
      success("Created", "Animal record has been created");
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
  if (!confirm("Are you sure you want to delete this animal?")) return;
  try {
    await deleteAnimal(id);
    success("Deleted", "Animal record has been deleted");
    await fetchData();
  } catch (e: any) {
    showError("Failed", e?.response?.data?.message || "An error occurred");
  }
}

async function addWeight() {
  if (!editingId.value) return;
  weightSaving.value = true;
  try {
    await createWeightHistory({
      animal_id: editingId.value,
      weight: weightForm.value.weight,
      record_date: weightForm.value.record_date,
      remarks: weightForm.value.remarks || undefined,
    });
    weightHistories.value = await listWeightHistories();
    weightForm.value = { weight: 0, record_date: "", remarks: "" };
    showWeightForm.value = false;
    success("Added", "Weight record has been added");
  } catch (e: any) {
    showError("Failed", e?.response?.data?.message || "An error occurred");
  } finally {
    weightSaving.value = false;
  }
}

async function removeWeight(id: number) {
  if (!confirm("Delete this weight record?")) return;
  try {
    await deleteWeightHistory(id);
    weightHistories.value = await listWeightHistories();
    success("Deleted", "Weight record has been removed");
  } catch (e: any) {
    showError("Failed", e?.response?.data?.message || "An error occurred");
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([
    { label: "Dashboard", to: "/dashboard" },
    { label: "Animals" },
  ]);
  headerStore.setActions([{ label: "Add New", onClick: openCreate }]);
  headerStore.setShowSearch(true);
  fetchData();
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
      :items="items"
      :loading="loading"
      :server-mode="true"
      :total-items="totalItems"
      :current-page="page"
      :page-size="pageSize"
      @update:current-page="page = $event; fetchData()"
      @update:page-size="pageSize = $event; page = 1; fetchData()"
    >
      <template #cell-species_name="{ item }">
        {{ item.species?.name || "-" }}
      </template>
      <template #cell-birth_date="{ item }">
        <DateDisplay :value="item.birth_date" />
      </template>
    </DataTable>

    <Modal
      :show="showModal"
      :title="editingId ? 'Edit Animal' : 'Add Animal'"
      @close="showModal = false"
      size="lg"
    >
      <form @submit.prevent="save">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Tag No</label
            >
            <input
              v-model="form.tag_no"
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
              @change="form.breed_id = null"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            >
              <option :value="null" disabled>Select species</option>
              <option v-for="s in speciesList" :key="s.id" :value="s.id">
                {{ s.name }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Breed</label
            >
            <select
              v-model="form.breed_id"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            >
              <option :value="null">Select breed</option>
              <option v-for="b in filteredBreeds" :key="b.id" :value="b.id">
                {{ b.name }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Gender</label
            >
            <select
              v-model="form.gender"
              required
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            >
              <option value="Male">Male</option>
              <option value="Female">Female</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Status</label
            >
            <select
              v-model="form.status"
              required
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            >
              <option value="Healthy">Healthy</option>
              <option value="Pregnant">Pregnant</option>
              <option value="Sick">Sick</option>
              <option value="Sold">Sold</option>
              <option value="Dead">Dead</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Birth Date</label
            >
            <input
              v-model="form.birth_date"
              type="date"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Purchase Date</label
            >
            <input
              v-model="form.purchase_date"
              type="date"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Purchase Price</label
            >
            <input
              v-model="form.purchase_price"
              type="number"
              min="0"
              step="0.01"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Color</label
            >
            <input
              v-model="form.color"
              type="text"
              class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
            />
          </div>
        </div>
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Remarks</label
          >
          <textarea
            v-model="form.remarks"
            rows="2"
            class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
          />
        </div>

        <div v-if="editingId" class="border-t pt-4 mt-4">
          <div class="mb-3 flex items-center justify-between">
            <h4 class="text-sm font-semibold text-gray-900">Weight History</h4>
            <button
              type="button"
              @click="showWeightForm = !showWeightForm"
              class="text-sm font-medium text-brand-600 hover:text-brand-700"
            >
              {{ showWeightForm ? "Cancel" : "Add Weight" }}
            </button>
          </div>

          <div v-if="showWeightForm" class="mb-3 flex items-end gap-3">
            <div class="flex-1">
              <label class="block text-xs font-medium text-gray-700 mb-1"
                >Weight (kg)</label
              >
              <input
                v-model="weightForm.weight"
                type="number"
                step="0.01"
                required
                class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
              />
            </div>
            <div class="flex-1">
              <label class="block text-xs font-medium text-gray-700 mb-1"
                >Date</label
              >
              <input
                v-model="weightForm.record_date"
                type="date"
                required
                class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
              />
            </div>
            <button
              type="button"
              :disabled="weightSaving"
              @click="addWeight"
              class="px-4 py-2.5 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50"
            >
              {{ weightSaving ? "..." : "Save" }}
            </button>
          </div>

          <div class="overflow-x-auto rounded-xl border border-gray-100">
            <table class="w-full text-left text-sm">
              <thead class="border-b bg-gray-50">
                <tr>
                  <th
                    class="px-3 py-2 font-medium text-gray-500 text-xs uppercase"
                  >
                    Weight
                  </th>
                  <th
                    class="px-3 py-2 font-medium text-gray-500 text-xs uppercase"
                  >
                    Date
                  </th>
                  <th
                    class="px-3 py-2 font-medium text-gray-500 text-xs uppercase"
                  >
                    Remarks
                  </th>
                  <th class="px-3 py-2 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-50">
                <tr v-if="weightHistoriesForAnimal.length === 0">
                  <td colspan="4" class="px-3 py-4 text-center text-gray-400">
                    No weight records.
                  </td>
                </tr>
                <tr v-for="w in weightHistoriesForAnimal" :key="w.id">
                  <td class="px-3 py-2">{{ w.weight }} kg</td>
                  <td class="px-3 py-2">
<DateDisplay :value="w.record_date" />
                  </td>
                  <td class="px-3 py-2">{{ w.remarks || "-" }}</td>
                  <td class="px-3 py-2 text-right">
                    <button
                      type="button"
                      @click="removeWeight(w.id)"
                      class="text-xs font-medium text-red-500 hover:text-red-600"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
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
