<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed} from 'vue'
import { listVaccines, createVaccine, updateVaccine, deleteVaccine, listSpecies } from '../../api'
import type { Vaccine, Species } from '../../api'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'

const searchQuery = computed(() => headerStore.searchQuery)
const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const items = ref<Vaccine[]>([])
const speciesList = ref<Species[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = ref({
  species_id: null as number | null,
  name: '',
  description: '',
  dose: '',
  minimum_age_value: 0,
  minimum_age_unit: 'Day',
  interval_value: 0,
  interval_unit: 'Day',
  is_repeatable: false,
})

const ageUnits = ['Day', 'Week', 'Month', 'Year']

async function fetchData() {
  loading.value = true
  try {
    const [vaccines, species] = await Promise.all([listVaccines(), listSpecies()])
    items.value = vaccines
    speciesList.value = species
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = {
    species_id: null, name: '', description: '', dose: '',
    minimum_age_value: 0, minimum_age_unit: 'Day',
    interval_value: 0, interval_unit: 'Day', is_repeatable: false,
  }
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  editingId.value = id
  form.value = {
    species_id: item.species_id,
    name: item.name,
    description: item.description || '',
    dose: item.dose || '',
    minimum_age_value: item.minimum_age_value,
    minimum_age_unit: item.minimum_age_unit,
    interval_value: item.interval_value,
    interval_unit: item.interval_unit,
    is_repeatable: item.is_repeatable,
  }
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      species_id: form.value.species_id!,
      name: form.value.name,
      minimum_age_value: form.value.minimum_age_value,
      minimum_age_unit: form.value.minimum_age_unit,
      interval_value: form.value.interval_value,
      interval_unit: form.value.interval_unit,
      is_repeatable: form.value.is_repeatable,
    }
    if (form.value.description) payload.description = form.value.description
    if (form.value.dose) payload.dose = form.value.dose

    if (editingId.value) {
      await updateVaccine(editingId.value, payload as any)
      success('Updated', 'Vaccine has been updated')
    } else {
      await createVaccine(payload as any)
      success('Created', 'Vaccine has been created')
    }
    showModal.value = false
    await fetchData()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this vaccine?')) return
  try {
    await deleteVaccine(id)
    success('Deleted', 'Vaccine has been deleted')
    await fetchData()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Vaccines' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  fetchData()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Vaccines" subtitle="Manage vaccine and immunization records" />
    <DataTable
      :columns="[
        { key: 'name', label: 'Name' },
        { key: 'species_name', label: 'Species' },
        { key: 'dose', label: 'Dose' },
        { key: 'minimum_age', label: 'Min Age' },
        { key: 'interval', label: 'Interval' },
        { key: 'is_repeatable', label: 'Repeatable' },
      ]"
      :items="items"
      :loading="loading"
      @edit="openEdit"
      @delete="handleDelete"
    >
      <template #cell-species_name="{ item }">
        {{ item.species?.name || '-' }}
      </template>
      <template #cell-dose="{ item }">
        {{ item.dose || '-' }}
      </template>
      <template #cell-minimum_age="{ item }">
        {{ item.minimum_age_value }} {{ item.minimum_age_unit }}{{ item.minimum_age_value > 1 ? 's' : '' }}
      </template>
      <template #cell-interval="{ item }">
        {{ item.interval_value }} {{ item.interval_unit }}{{ item.interval_value > 1 ? 's' : '' }}
      </template>
      <template #cell-is_repeatable="{ item }">
        {{ item.is_repeatable ? 'Yes' : 'No' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="editingId ? 'Edit Vaccine' : 'Add Vaccine'" @close="showModal = false" size="md">
      <form @submit.prevent="save">
        <div class="grid grid-cols-2 gap-4">
          <div class="col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input v-model="form.name" type="text" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Species</label>
            <select v-model="form.species_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option :value="null" disabled>Select species</option>
              <option v-for="s in speciesList" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Dose</label>
            <input v-model="form.dose" type="text" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Minimum Age Value</label>
            <input v-model="form.minimum_age_value" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Minimum Age Unit</label>
            <select v-model="form.minimum_age_unit" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option v-for="u in ageUnits" :key="u" :value="u">{{ u }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Interval Value</label>
            <input v-model="form.interval_value" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Interval Unit</label>
            <select v-model="form.interval_unit" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option v-for="u in ageUnits" :key="u" :value="u">{{ u }}</option>
            </select>
          </div>
          <div class="col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea v-model="form.description" rows="2" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div class="col-span-2">
            <label class="flex items-center gap-2 text-sm font-medium text-gray-700">
              <input v-model="form.is_repeatable" type="checkbox" class="rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
              Is Repeatable
            </label>
          </div>
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Save' }}</button>
      </template>
    </Modal>
  </div>
</template>
