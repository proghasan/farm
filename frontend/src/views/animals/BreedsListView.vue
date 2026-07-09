<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { listBreeds, createBreed, updateBreed, deleteBreed, listSpecies } from '../../api'
import type { Breed, Species } from '../../api'
import { DataTable } from '../../components/DataTable'
import RowActions from '../../components/RowActions.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'
import { getErrorMessage } from "../../utils/error";

const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const items = ref<Breed[]>([])
const speciesList = ref<Species[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', species_id: null as number | null })
const saving = ref(false)

async function fetchData() {
  loading.value = true
  try {
    items.value = await listBreeds()
    speciesList.value = await listSpecies()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { name: '', species_id: null }
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  editingId.value = id
  form.value = { name: item.name, species_id: item.species_id }
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    const payload = { name: form.value.name, species_id: form.value.species_id! }
    if (editingId.value) {
      await updateBreed(editingId.value, payload)
      success('Updated', 'Breed has been updated')
    } else {
      await createBreed(payload)
      success('Created', 'Breed has been created')
    }
    showModal.value = false
    await fetchData()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this breed?')) return
  try {
    await deleteBreed(id)
    success('Deleted', 'Breed has been deleted')
    await fetchData()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  }
}

const columns = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'species_name', label: 'Species' },
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

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Breeds' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  fetchData()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Breeds" subtitle="Manage animal breed records" />
<DataTable
      :columns="columns"
      :items="items"
      :loading="loading"
    >
      <template #cell-species_name="{ item }">
        {{ item.species?.name || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="editingId ? 'Edit Breed' : 'Add Breed'" @close="showModal = false" size="sm">
      <form @submit.prevent="save" class="space-y-4">
        <div>
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
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Save' }}</button>
      </template>
    </Modal>
  </div>
</template>
