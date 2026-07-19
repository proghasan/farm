<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { FileText, ShoppingCart, Upload } from "@lucide/vue";
import PageHeader from "../../components/PageHeader.vue";
import WeightHistorySection from "./WeightHistorySection.vue";
import { useToast } from "../../composables/useToast";
import { useHeaderStore } from "../../stores/header";
import { useAnimalStore } from "../../stores/animal";
import { useSpeciesStore } from "../../stores/species";
import { useBreedStore } from "../../stores/breed";
import { getFirstErrorMessage } from "../../utils/error";

const route = useRoute();
const router = useRouter();
const { success, error: showError } = useToast();
const headerStore = useHeaderStore();
const animalStore = useAnimalStore();
const speciesStore = useSpeciesStore();
const breedStore = useBreedStore();

const isEdit = computed(() => !!route.params.id);
const animalId = computed(() => Number(route.params.id));

const saving = ref(false);
const loading = ref(false);

const form = ref({
  tag_no: "",
  species_id: null as number | null,
  breed_id: null as number | null,
  gender: "Male",
  birth_date: new Date().toISOString().slice(0, 10),
  purchase_date: new Date().toISOString().slice(0, 10),
  purchase_price: 0,
  current_weight: null as number | null,
  last_vaccine: new Date().toISOString().slice(0, 10),
  color: "",
  status: "Healthy",
  remarks: "",
});

const filteredBreeds = computed(() =>
  breedStore.allItems.filter((b) => b.species_id === form.value.species_id),
);

async function loadSpeciesAndBreeds() {
  await Promise.all([
    speciesStore.fetchAll({ all: "true" }),
    breedStore.fetchAll({ all: "true" }),
  ]);
}

async function loadAnimal() {
  if (!isEdit.value) return;
  loading.value = true;
  try {
    await animalStore.fetchById(animalId.value);
    const animal = animalStore.currentAnimal!;
    form.value = {
      tag_no: animal.tag_no,
      species_id: animal.species_id,
      breed_id: animal.breed_id ?? null,
      gender: animal.gender,
      birth_date: animal.birth_date || "",
      purchase_date: animal.purchase_date || "",
      purchase_price: animal.purchase_price || 0,
      current_weight: animal.current_weight ?? null,
      last_vaccine: animal.last_vaccine || "",
      color: animal.color || "",
      status: animal.status,
      remarks: animal.remarks || "",
    };
  } finally {
    loading.value = false;
  }
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
      breed_id: form.value.breed_id!,
    };
    if (form.value.birth_date) payload.birth_date = form.value.birth_date;
    if (form.value.purchase_date)
      payload.purchase_date = form.value.purchase_date;
    if (form.value.color) payload.color = form.value.color;
    if (form.value.remarks) payload.remarks = form.value.remarks;
    if (form.value.current_weight != null)
      payload.current_weight = form.value.current_weight;
    if (form.value.last_vaccine) payload.last_vaccine = form.value.last_vaccine;

    if (isEdit.value) {
      await animalStore.update(animalId.value, payload as any);
      success("Updated", "Animal record has been updated");
    } else {
      await animalStore.create(payload as any);
      success("Created", "Animal record has been created");
    }
    router.push("/animals");
  } catch (e: any) {
    showError("Failed", getFirstErrorMessage(e));
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  const label = isEdit.value ? "Edit Animal" : "Add Animal";
  headerStore.setBreadcrumb([
    { label: "Dashboard", to: "/dashboard" },
    { label: "Animals", to: "/animals" },
    { label },
  ]);
  headerStore.setShowSearch(false);
  loadSpeciesAndBreeds();
  loadAnimal();
});
onUnmounted(() => headerStore.clear());
</script>

<template>
  <div>
    <PageHeader
      :title="isEdit ? 'Edit Animal' : 'Add Animal'"
      :subtitle="isEdit ? 'Update animal record' : 'Register a new animal'"
    />

    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="text-sm text-gray-400">Loading...</div>
    </div>

    <form v-else @submit.prevent="save">
      <div class="grid grid-cols-2 gap-6 pb-32">
        <div class="space-y-6">
          <section
            class="bg-white rounded-2xl border border-gray-100 overflow-hidden"
          >
            <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-50 bg-gray-50/50">
              <div class="w-8 h-8 rounded-lg bg-brand-50 flex items-center justify-center">
                <FileText class="w-4 h-4 text-brand-600" />
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-900">Basic Information</h3>
                <p class="text-xs text-gray-400">
                  Tag number, species, breed and status
                </p>
              </div>
            </div>

            <div class="p-6 grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Tag No <span class="text-red-400">*</span></label
                >
                <input
                  v-model="form.tag_no"
                  type="text"
                  placeholder="e.g. AN-001"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Species <span class="text-red-400">*</span></label
                >
                <select
                  v-model="form.species_id"
                  @change="form.breed_id = null"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
                >
                  <option :value="null" disabled>Select species</option>
                  <option v-for="s in speciesStore.allItems" :key="s.id" :value="s.id">
                    {{ s.name }}
                  </option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Breed <span class="text-red-400">*</span></label
                >
                <select
                  v-model="form.breed_id"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
                >
                  <option :value="null" disabled>Select breed</option>
                  <option v-for="b in filteredBreeds" :key="b.id" :value="b.id">
                    {{ b.name }}
                  </option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Gender <span class="text-red-400">*</span></label
                >
                <select
                  v-model="form.gender"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
                >
                  <option value="Male">Male</option>
                  <option value="Female">Female</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Status <span class="text-red-400">*</span></label
                >
                <select
                  v-model="form.status"
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
                  >Color</label
                >
                <input
                  v-model="form.color"
                  type="text"
                  placeholder="e.g. Brown, White"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300"
                />
              </div>
              <div class="col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Remarks</label
                >
                <textarea
                  v-model="form.remarks"
                  rows="2"
                  placeholder="Any additional notes..."
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300"
                />
              </div>
            </div>
          </section>

          <section
            class="bg-white rounded-2xl border border-gray-100 overflow-hidden"
          >
            <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-50 bg-gray-50/50">
              <div class="w-8 h-8 rounded-lg bg-amber-50 flex items-center justify-center">
                <ShoppingCart class="w-4 h-4 text-amber-600" />
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-900">Purchase Information</h3>
                <p class="text-xs text-gray-400">
                  Optional — fill in when buying an animal
                </p>
              </div>
            </div>

            <div class="p-6 grid grid-cols-2 gap-4">
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
                <div class="relative">
                  <span
                    class="absolute left-3.5 top-1/2 -translate-y-1/2 text-sm text-gray-400"
                    >$</span
                  >
                  <input
                    v-model="form.purchase_price"
                    type="number"
                    min="0"
                    step="0.01"
                    placeholder="0.00"
                    class="rounded-xl border border-gray-200 pl-8 pr-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300"
                  />
                </div>
              </div>
            </div>
          </section>
        </div>

        <div class="space-y-6">
          <section
            class="bg-white rounded-2xl border border-gray-100 overflow-hidden"
          >
            <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-50 bg-gray-50/50">
              <div class="w-8 h-8 rounded-lg bg-emerald-50 flex items-center justify-center">
                <Upload class="w-4 h-4 text-emerald-600" />
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-900">Birth Information</h3>
                <p class="text-xs text-gray-400">
                  Date of birth and current weight
                </p>
              </div>
            </div>

            <div class="p-6 grid grid-cols-2 gap-4">
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
                  >Current Weight
                  <span class="text-gray-400 font-normal">(kg)</span></label
                >
                <input
                  v-model="form.current_weight"
                  type="number"
                  step="0.01"
                  min="0"
                  placeholder="0.00"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full placeholder:text-gray-300"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Last Vaccine</label
                >
                <input
                  v-model="form.last_vaccine"
                  type="date"
                  class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full"
                />
              </div>
            </div>
          </section>

          <WeightHistorySection v-if="isEdit" :animal-id="animalId" />
        </div>
      </div>

      <div
        class="fixed bottom-0 left-0 right-0 bg-white/80 backdrop-blur-lg border-t border-gray-100 px-8 py-4 z-10"
      >
        <div class="flex items-center justify-end gap-3">
          <button
            type="button"
            @click="router.push('/animals')"
            class="px-5 py-2.5 text-sm font-medium text-gray-600 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 hover:text-gray-700 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-2.5 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50 transition-colors flex items-center gap-2"
          >
            <svg
              v-if="saving"
              class="w-4 h-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              />
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
              />
            </svg>
            {{ saving ? "Saving..." : "Save" }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>
