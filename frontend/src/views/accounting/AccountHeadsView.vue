<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted} from 'vue'
import { DataTable } from '../../components/DataTable'
import RowActions from '../../components/RowActions.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import {
  listAccountHeads,
  createAccountHead,
  updateAccountHead,
  deleteAccountHead,
} from '../../api'
import type { AccountHead } from '../../api'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'
import { getFirstErrorMessage } from "../../utils/error";

const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const items = ref<AccountHead[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  type: 'Income' as 'Income' | 'Expense',
  name: '',
  description: '',
})

const columns = [
  { key: 'type', label: 'Type', component: StatusBadge },
  { key: 'name', label: 'Name', sortable: true },
  { key: 'description', label: 'Description' },
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

const modalTitle = computed(() => (editingId.value ? 'Edit Account Head' : 'Add Account Head'))

async function fetchItems() {
  loading.value = true
  try {
    items.value = await listAccountHeads()
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = { type: 'Income', name: '', description: '' }
}

function openCreate() {
  editingId.value = null
  resetForm()
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find((i) => i.id === id)
  if (!item) return
  editingId.value = id
  form.value = { type: item.type, name: item.name, description: item.description ?? '' }
  showModal.value = true
}

async function handleSave() {
  const payload = { type: form.value.type, name: form.value.name, description: form.value.description || undefined }
  try {
    if (editingId.value) {
      await updateAccountHead(editingId.value, payload)
      success('Updated', 'Account head has been updated')
    } else {
      await createAccountHead(payload)
      success('Created', 'Account head has been created')
    }
    showModal.value = false
    await fetchItems()
  } catch (e: any) {
    showError('Failed', getFirstErrorMessage(e))
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this account head?')) return
  try {
    await deleteAccountHead(id)
    success('Deleted', 'Account head has been deleted')
    await fetchItems()
  } catch (e: any) {
    showError('Failed', getFirstErrorMessage(e))
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Account Heads' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  fetchItems()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Account Heads" subtitle="Manage income and expense categories" />
    <DataTable
      :columns="columns"
      :items="items"
      :loading="loading"
    >
      <template #cell-description="{ item }">
        {{ item.description || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="modalTitle" @close="showModal = false" size="sm">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select v-model="form.type" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option value="Income">Income</option>
            <option value="Expense">Expense</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" type="text" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea v-model="form.description" rows="3" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700">{{ editingId ? 'Update' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
