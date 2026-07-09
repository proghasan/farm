<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed} from 'vue'
import {
  listInventoryTransactions,
  createInventoryTransaction,
  deleteInventoryTransaction,
  listInventoryItems,
} from '../../api'
import type { InventoryTransaction, InventoryItem } from '../../api'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'

const searchQuery = computed(() => headerStore.searchQuery)
const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const transactions = ref<InventoryTransaction[]>([])
const items = ref<InventoryItem[]>([])
const loading = ref(false)
const showModal = ref(false)
const saving = ref(false)

const form = ref({
  inventory_item_id: 0,
  transaction_type: 'Purchase',
  quantity: 0,
  transaction_date: '',
  remarks: '',
})

const transactionTypes = ['Purchase', 'Sale', 'Consumption', 'Adjustment', 'Return', 'Damage']

function formatDate(d: string) {
  if (!d) return ''
  return new Date(d).toLocaleDateString()
}

function today() {
  return new Date().toISOString().slice(0, 10)
}

async function load() {
  loading.value = true
  try {
    const [txns, invItems] = await Promise.all([
      listInventoryTransactions(),
      listInventoryItems(),
    ])
    transactions.value = txns
    items.value = invItems
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = {
    inventory_item_id: 0,
    transaction_type: 'Purchase',
    quantity: 0,
    transaction_date: today(),
    remarks: '',
  }
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    await createInventoryTransaction(form.value)
    success('Created', 'Transaction has been created')
    showModal.value = false
    await load()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this transaction?')) return
  try {
    await deleteInventoryTransaction(id)
    success('Deleted', 'Transaction has been deleted')
    await load()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Inventory Transactions' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  load()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Inventory Transactions" subtitle="Track inventory movements and adjustments" />
    <DataTable
      :columns="[
        { key: 'item_name', label: 'Item' },
        { key: 'transaction_type', label: 'Type' },
        { key: 'quantity', label: 'Quantity' },
        { key: 'transaction_date', label: 'Date' },
        { key: 'remarks', label: 'Remarks' },
      ]"
      :items="transactions"
      :loading="loading"
      @delete="handleDelete"
    >
      <template #cell-item_name="{ item }">
        {{ item.inventory_item?.name }}
      </template>
      <template #cell-transaction_date="{ item }">
        {{ formatDate(item.transaction_date) }}
      </template>
      <template #cell-remarks="{ item }">
        {{ item.remarks || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" title="Add Transaction" @close="showModal = false" size="md">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Item</label>
          <select v-model.number="form.inventory_item_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option :value="0" disabled>Select item</option>
            <option v-for="item in items" :key="item.id" :value="item.id">{{ item.name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select v-model="form.transaction_type" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option v-for="t in transactionTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
            <input v-model.number="form.quantity" type="number" step="0.01" min="0" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
            <input v-model="form.transaction_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Remarks</label>
          <input v-model="form.remarks" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" placeholder="Optional remarks" />
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
