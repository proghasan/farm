<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  listInventoryCategories,
  createInventoryCategory,
  updateInventoryCategory,
  deleteInventoryCategory,
} from '../../api'
import type { InventoryCategory } from '../../api'
import { DataTable } from '../../components/DataTable'
import RowActions from '../../components/RowActions.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'

const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const categories = ref<InventoryCategory[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const columns = [
  { key: 'name', label: 'Name', sortable: true },
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

const form = ref({ name: '' })
const saving = ref(false)

async function load() {
  loading.value = true
  try {
    categories.value = await listInventoryCategories()
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
  const cat = categories.value.find(c => c.id === id)
  if (!cat) return
  editingId.value = cat.id
  form.value = { name: cat.name }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editingId.value) {
      await updateInventoryCategory(editingId.value, form.value)
      success('Updated', 'Category has been updated')
    } else {
      await createInventoryCategory(form.value)
      success('Created', 'Category has been created')
    }
    showModal.value = false
    await load()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this category?')) return
  try {
    await deleteInventoryCategory(id)
    success('Deleted', 'Category has been deleted')
    await load()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Inventory Categories' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  load()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Inventory Categories" subtitle="Organize inventory items by category" />
    <DataTable
      :columns="columns"
      :items="categories"
      :loading="loading"
    />

    <Modal :show="showModal" :title="editingId ? 'Edit Category' : 'Add Category'" @close="showModal = false" size="sm">
      <form @submit.prevent="handleSave">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" placeholder="Category name" />
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : editingId ? 'Update' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
