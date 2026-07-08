<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  listInventoryItems,
  createInventoryItem,
  updateInventoryItem,
  deleteInventoryItem,
  listInventoryCategories,
} from '../../api'
import type { InventoryItem, InventoryCategory } from '../../api'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import { useToast } from '../../composables/useToast'

const { success, error: showError } = useToast()

const items = ref<InventoryItem[]>([])
const categories = ref<InventoryCategory[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const form = ref({
  category_id: 0,
  name: '',
  sku: '',
  unit: '',
  purchase_price: 0,
  selling_price: 0,
})

async function load() {
  loading.value = true
  try {
    const [itms, cats] = await Promise.all([
      listInventoryItems(),
      listInventoryCategories(),
    ])
    items.value = itms
    categories.value = cats
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { category_id: 0, name: '', sku: '', unit: '', purchase_price: 0, selling_price: 0 }
  showModal.value = true
}

function openEdit(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  editingId.value = item.id
  form.value = {
    category_id: item.category_id,
    name: item.name,
    sku: item.sku ?? '',
    unit: item.unit,
    purchase_price: item.purchase_price,
    selling_price: item.selling_price,
  }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editingId.value) {
      await updateInventoryItem(editingId.value, form.value)
      success('Updated', 'Item has been updated')
    } else {
      await createInventoryItem(form.value)
      success('Created', 'Item has been created')
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
  if (!confirm('Are you sure you want to delete this item?')) return
  try {
    await deleteInventoryItem(id)
    success('Deleted', 'Item has been deleted')
    await load()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Inventory Items</h1>
        <p class="text-sm text-gray-500">Manage inventory items</p>
      </div>
      <button @click="openCreate" class="inline-flex items-center gap-2 px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-sm font-medium rounded-xl transition-colors">Add New</button>
    </div>

    <DataTable
      :columns="[
        { key: 'name', label: 'Name' },
        { key: 'category_name', label: 'Category' },
        { key: 'sku', label: 'SKU' },
        { key: 'unit', label: 'Unit' },
        { key: 'purchase_price', label: 'Purchase Price' },
        { key: 'selling_price', label: 'Selling Price' },
      ]"
      :items="items"
      :loading="loading"
      @edit="openEdit"
      @delete="handleDelete"
    >
      <template #cell-category_name="{ item }">
        {{ item.category?.name }}
      </template>
      <template #cell-sku="{ item }">
        {{ item.sku || '-' }}
      </template>
      <template #cell-purchase_price="{ item }">
        ${{ Number(item.purchase_price).toFixed(2) }}
      </template>
      <template #cell-selling_price="{ item }">
        ${{ Number(item.selling_price).toFixed(2) }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="editingId ? 'Edit Item' : 'Add Item'" @close="showModal = false" size="md">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select v-model.number="form.category_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option :value="0" disabled>Select category</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" placeholder="Item name" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">SKU</label>
          <input v-model="form.sku" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" placeholder="SKU (optional)" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Unit</label>
          <input v-model="form.unit" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" placeholder="e.g. kg, pcs, liters" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Purchase Price</label>
            <input v-model.number="form.purchase_price" type="number" step="0.01" min="0" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Selling Price</label>
            <input v-model.number="form.selling_price" type="number" step="0.01" min="0" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : editingId ? 'Update' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
