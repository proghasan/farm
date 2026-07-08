<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listSpecies, createSpecies, updateSpecies, deleteSpecies } from '../../api'
import type { Species } from '../../api'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import { useToast } from '../../composables/useToast'

const { success, error: showError } = useToast()

const items = ref<(Species & { created_at?: string })[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '' })
const saving = ref(false)

async function fetchData() {
  loading.value = true
  try {
    items.value = await listSpecies() as (Species & { created_at?: string })[]
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { name: '' }
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  editingId.value = id
  form.value = { name: item.name }
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    if (editingId.value) {
      await updateSpecies(editingId.value, form.value)
      success('Updated', 'Species has been updated')
    } else {
      await createSpecies(form.value)
      success('Created', 'Species has been created')
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
  if (!confirm('Are you sure you want to delete this species?')) return
  try {
    await deleteSpecies(id)
    success('Deleted', 'Species has been deleted')
    await fetchData()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Species</h1>
        <p class="text-sm text-gray-500">Manage animal species</p>
      </div>
      <button @click="openCreate" class="inline-flex items-center gap-2 px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-sm font-medium rounded-xl transition-colors">Add New</button>
    </div>

    <DataTable
      :columns="[
        { key: 'name', label: 'Name' },
        { key: 'created_at', label: 'Created At' },
      ]"
      :items="items"
      :loading="loading"
      @edit="openEdit"
      @delete="handleDelete"
    >
      <template #cell-created_at="{ item }">
        {{ item.created_at ? new Date(item.created_at).toLocaleDateString() : '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="editingId ? 'Edit Species' : 'Add Species'" @close="showModal = false" size="sm">
      <form @submit.prevent="save">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" type="text" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Save' }}</button>
      </template>
    </Modal>
  </div>
</template>
