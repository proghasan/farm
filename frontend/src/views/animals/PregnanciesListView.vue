<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  listPregnancies, createPregnancy, updatePregnancy, deletePregnancy,
  listAnimals,
} from '../../api'
import type { Pregnancy, Animal } from '../../api'
import { DataTable } from '../../components/DataTable'
import RowActions from '../../components/RowActions.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import DateDisplay from '../../components/DateDisplay.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'
import { getErrorMessage } from "../../utils/error";

const headerStore = useHeaderStore()
const { success, error: showError } = useToast()

const items = ref<Pregnancy[]>([])
const animalList = ref<Animal[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = ref({
  animal_id: 0,
  breeder_id: null as number | null,
  mating_date: '',
  expected_due_date: '',
  actual_birth_date: '',
  status: 'Mated',
  note: '',
  number_of_children: null as number | null,
  number_of_male_children: null as number | null,
  number_of_female_children: null as number | null,
  number_of_dead_children: null as number | null,
})

const columns = [
  { key: 'animal_display', label: 'Animal' },
  { key: 'mating_date', label: 'Mating Date' },
  { key: 'expected_due_date', label: 'Due Date' },
  { key: 'status', label: 'Status', component: StatusBadge },
  { key: 'number_of_children', label: 'Offspring' },
  {
    key: 'actions',
    label: '',
    component: RowActions,
    componentProps: {
      actions: [
        { label: 'Edit', icon: 'Pencil', onClick: (item: any) => openEdit(item.id) },
        { label: 'Delete', icon: 'Trash2', onClick: (item: any) => handleDelete(item.id), danger: true },
      ],
    },
  },
]

const statuses = ['Mated', 'Pregnant', 'Delivered', 'Aborted', 'Miscarriage', 'Failed']

function today() {
  return new Date().toISOString().slice(0, 10)
}

function resetForm() {
  form.value = {
    animal_id: 0,
    breeder_id: null,
    mating_date: today(),
    expected_due_date: '',
    actual_birth_date: '',
    status: 'Mated',
    note: '',
    number_of_children: null,
    number_of_male_children: null,
    number_of_female_children: null,
    number_of_dead_children: null,
  }
}

async function load() {
  loading.value = true
  try {
    const [pregs, animals] = await Promise.all([listPregnancies(), listAnimals()])
    items.value = pregs
    animalList.value = animals
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  resetForm()
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  editingId.value = id
  form.value = {
    animal_id: item.animal_id,
    breeder_id: item.breeder_id ?? null,
    mating_date: item.mating_date,
    expected_due_date: item.expected_due_date,
    actual_birth_date: item.actual_birth_date ?? '',
    status: item.status,
    note: item.note ?? '',
    number_of_children: item.number_of_children ?? null,
    number_of_male_children: item.number_of_male_children ?? null,
    number_of_female_children: item.number_of_female_children ?? null,
    number_of_dead_children: item.number_of_dead_children ?? null,
  }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      animal_id: form.value.animal_id,
      mating_date: form.value.mating_date,
      expected_due_date: form.value.expected_due_date,
      status: form.value.status,
    }
    if (form.value.breeder_id) payload.breeder_id = form.value.breeder_id
    if (form.value.actual_birth_date) payload.actual_birth_date = form.value.actual_birth_date
    if (form.value.note) payload.note = form.value.note
    if (form.value.number_of_children !== null) payload.number_of_children = form.value.number_of_children
    if (form.value.number_of_male_children !== null) payload.number_of_male_children = form.value.number_of_male_children
    if (form.value.number_of_female_children !== null) payload.number_of_female_children = form.value.number_of_female_children
    if (form.value.number_of_dead_children !== null) payload.number_of_dead_children = form.value.number_of_dead_children

    if (editingId.value) {
      await updatePregnancy(editingId.value, payload)
      success('Updated', 'Pregnancy record has been updated')
    } else {
      await createPregnancy(payload)
      success('Created', 'Pregnancy record has been created')
    }
    showModal.value = false
    await load()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this pregnancy record?')) return
  try {
    await deletePregnancy(id)
    success('Deleted', 'Pregnancy record has been deleted')
    await load()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Pregnancies' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  load()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Pregnancies" subtitle="Track pregnancy and breeding records" />
    <DataTable
      :columns="columns"
      :items="items"
      :loading="loading"
    >
      <template #cell-animal_display="{ item }">
        {{ item.animal?.tag_no || '-' }}
      </template>
      <template #cell-mating_date="{ item }">
        <DateDisplay :value="item.mating_date" />
      </template>
      <template #cell-expected_due_date="{ item }">
        <DateDisplay :value="item.expected_due_date" />
      </template>
      <template #cell-number_of_children="{ item }">
        {{ item.number_of_children ?? '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="editingId ? 'Edit Pregnancy' : 'Add Pregnancy'" @close="showModal = false" size="lg">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Animal</label>
            <select v-model.number="form.animal_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option :value="0" disabled>Select animal</option>
              <option v-for="a in animalList.filter(a => a.gender === 'Female')" :key="a.id" :value="a.id">{{ a.tag_no }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Breeder (optional)</label>
            <select v-model.number="form.breeder_id" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option :value="null">Not specified</option>
              <option v-for="a in animalList.filter(a => a.gender === 'Male')" :key="a.id" :value="a.id">{{ a.tag_no }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Mating Date</label>
            <input v-model="form.mating_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Expected Due Date</label>
            <input v-model="form.expected_due_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Actual Birth Date</label>
            <input v-model="form.actual_birth_date" type="date" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select v-model="form.status" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option v-for="s in statuses" :key="s" :value="s">{{ s }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Total Offspring</label>
            <input v-model.number="form.number_of_children" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Male Offspring</label>
            <input v-model.number="form.number_of_male_children" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Female Offspring</label>
            <input v-model.number="form.number_of_female_children" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Dead Offspring</label>
            <input v-model.number="form.number_of_dead_children" type="number" min="0" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Notes</label>
          <textarea v-model="form.note" rows="2" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
      </form>
      <template #footer>
        <button type="button" @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : editingId ? 'Update' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>