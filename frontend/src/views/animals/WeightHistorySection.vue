<script setup lang="ts">
import { ref, onMounted } from "vue";
import { DataTable } from "../../components/DataTable";
import type { Column } from "../../components/DataTable";
import Drawer from "../../components/Drawer.vue";
import RowActions from "../../components/RowActions.vue";
import DateDisplay from "../../components/DateDisplay.vue";
import { Scale } from "@lucide/vue";
import type { WeightHistory } from "../../api";
import { useToast } from "../../composables/useToast";
import { useWeightHistoryStore } from "../../stores/weightHistory";
import { getFirstErrorMessage } from "../../utils/error";

const props = defineProps<{
  animalId: number;
}>();

const { success, error: showError } = useToast();
const weightHistoryStore = useWeightHistoryStore();

const showWeightForm = ref(false);
const editingWeightId = ref<number | null>(null);
const weightSaving = ref(false);
const weightForm = ref({ weight: 0, record_date: new Date().toISOString().slice(0, 10), remarks: "" });

async function fetchWeights() {
  await weightHistoryStore.fetchByAnimal(props.animalId);
}

async function addWeight() {
  weightSaving.value = true;
  try {
    if (editingWeightId.value) {
      await weightHistoryStore.update(editingWeightId.value, {
        weight: weightForm.value.weight,
        record_date: weightForm.value.record_date,
        remarks: weightForm.value.remarks || undefined,
      });
      success("Updated", "Weight record has been updated");
    } else {
      await weightHistoryStore.create({
        animal_id: props.animalId,
        weight: weightForm.value.weight,
        record_date: weightForm.value.record_date,
        remarks: weightForm.value.remarks || undefined,
      });
      success("Added", "Weight record has been added");
    }
    await fetchWeights();
    weightForm.value = { weight: 0, record_date: new Date().toISOString().slice(0, 10), remarks: "" };
    editingWeightId.value = null;
    showWeightForm.value = false;
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  } finally {
    weightSaving.value = false;
  }
}

async function removeWeight(id: number) {
  if (!confirm("Delete this weight record?")) return;
  try {
    await weightHistoryStore.remove(id);
    await fetchWeights();
    success("Deleted", "Weight record has been removed");
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  }
}

function editWeight(w: WeightHistory) {
  weightForm.value = { weight: w.weight, record_date: w.record_date || "", remarks: w.remarks || "" };
  editingWeightId.value = w.id;
  showWeightForm.value = true;
}

function openAddForm() {
  editingWeightId.value = null;
  weightForm.value = { weight: 0, record_date: new Date().toISOString().slice(0, 10), remarks: "" };
  showWeightForm.value = true;
}

const columns: Column[] = [
  { key: "weight", label: "Weight", sortable: false },
  { key: "record_date", label: "Date", sortable: false },
  { key: "remarks", label: "Remarks", sortable: false },
  {
    key: "actions",
    label: "",
    sortable: false,
    component: RowActions,
    tdPosition: "right",
    componentProps: {
      actions: [
        { label: "Edit", icon: "Pencil", onClick: (item: any) => editWeight(item) },
        { label: "Delete", icon: "Trash2", onClick: (item: any) => removeWeight(item.id), danger: true },
      ],
    },
  },
];

onMounted(() => fetchWeights());
</script>

<template>
  <section class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
    <div class="flex items-center justify-between gap-3 px-6 py-4 border-b border-gray-50 bg-gray-50/50">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center">
          <Scale class="w-4 h-4 text-blue-600" />
        </div>
        <div>
          <h3 class="text-sm font-semibold text-gray-900">Weight History</h3>
          <p class="text-xs text-gray-400">Track weight changes over time</p>
        </div>
      </div>
      <button
        type="button"
        @click="openAddForm"
        class="text-sm font-medium text-brand-600 hover:text-brand-700 transition-colors"
      >
        + Add Weight
      </button>
    </div>

    <div class="p-6 space-y-4">
      <Drawer
        :show="showWeightForm"
        :title="editingWeightId ? 'Edit Weight' : 'Add Weight'"
        @close="showWeightForm = false; editingWeightId = null"
      >
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Weight <span class="text-red-400">*</span></label>
            <input v-model="weightForm.weight" type="number" step="0.01" required placeholder="0.00" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Date <span class="text-red-400">*</span></label>
            <input v-model="weightForm.record_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Remarks</label>
            <textarea v-model="weightForm.remarks" rows="3" placeholder="Optional notes..." class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300" />
          </div>
        </div>
        <template #footer>
          <button type="button" @click="showWeightForm = false" class="px-5 py-2.5 text-sm font-medium text-gray-600 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors">Cancel</button>
          <button type="button" :disabled="weightSaving" @click="addWeight" class="px-6 py-2.5 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50 transition-colors flex items-center gap-2">
            <svg v-if="weightSaving" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
            {{ weightSaving ? "Saving..." : "Save" }}
          </button>
        </template>
      </Drawer>

      <DataTable
        :columns="columns"
        :items="weightHistoryStore.items"
        hide-pagination
      >
        <template #cell-weight="{ item }">
          <span class="text-sm font-medium text-gray-900">{{ item.weight }} kg</span>
        </template>
        <template #cell-record_date="{ item }">
          <DateDisplay :value="item.record_date" />
        </template>
        <template #cell-remarks="{ item }">
          <span class="text-sm text-gray-500">{{ item.remarks || "—" }}</span>
        </template>
      </DataTable>
    </div>
  </section>
</template>
