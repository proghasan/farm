<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import {
  listAccountTransactions,
  createAccountTransaction,
  deleteAccountTransaction,
  listAccountHeads,
} from '../../api'
import type { AccountTransaction, AccountHead } from '../../api'
import { useToast } from '../../composables/useToast'

const { success, error: showError } = useToast()

const items = ref<AccountTransaction[]>([])
const accountHeads = ref<AccountHead[]>([])
const loading = ref(false)
const showModal = ref(false)
const saving = ref(false)

const form = ref({
  account_head_id: 0,
  transaction_date: '',
  amount: 0,
  payment_method: 'Cash',
  reference_no: '',
  description: '',
})

function headName(item: AccountTransaction) {
  return item.account_head?.name ?? '-'
}

async function fetchItems() {
  loading.value = true
  try {
    const [transactions, heads] = await Promise.all([listAccountTransactions(), listAccountHeads()])
    items.value = transactions
    accountHeads.value = heads
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = {
    account_head_id: accountHeads.value[0]?.id ?? 0,
    transaction_date: '',
    amount: 0,
    payment_method: 'Cash',
    reference_no: '',
    description: '',
  }
}

function openCreate() {
  resetForm()
  showModal.value = true
}

async function handleSave() {
  saving.value = true
  try {
    await createAccountTransaction({
      account_head_id: form.value.account_head_id,
      transaction_date: form.value.transaction_date,
      amount: form.value.amount,
      payment_method: form.value.payment_method,
      reference_no: form.value.reference_no || undefined,
      description: form.value.description || undefined,
    })
    success('Created', 'Transaction has been created')
    showModal.value = false
    await fetchItems()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this transaction?')) return
  try {
    await deleteAccountTransaction(id)
    success('Deleted', 'Transaction has been deleted')
    await fetchItems()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(fetchItems)
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Account Transactions</h1>
        <p class="text-sm text-gray-500">Record financial transactions</p>
      </div>
      <button @click="openCreate" class="inline-flex items-center gap-2 px-4 py-2.5 bg-brand-600 hover:bg-brand-700 text-white text-sm font-medium rounded-xl transition-colors">Add New</button>
    </div>

    <DataTable
      :columns="[
        { key: 'transaction_date', label: 'Date' },
        { key: 'account_head_name', label: 'Account Head' },
        { key: 'amount', label: 'Amount ($)' },
        { key: 'payment_method', label: 'Payment Method' },
        { key: 'reference_no', label: 'Reference' },
        { key: 'description', label: 'Description' },
      ]"
      :items="items"
      :loading="loading"
      @delete="handleDelete"
    >
      <template #cell-transaction_date="{ item }">
        {{ item.transaction_date ? new Date(item.transaction_date).toLocaleDateString() : '-' }}
      </template>
      <template #cell-account_head_name="{ item }">
        {{ headName(item) }}
      </template>
      <template #cell-amount="{ item }">
        ${{ Number(item.amount).toFixed(2) }}
      </template>
      <template #cell-reference_no="{ item }">
        {{ item.reference_no || '-' }}
      </template>
      <template #cell-description="{ item }">
        {{ item.description || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" title="Add Transaction" @close="showModal = false" size="md">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Account Head</label>
          <select v-model="form.account_head_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option v-for="h in accountHeads" :key="h.id" :value="h.id">{{ h.name }} ({{ h.type }})</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input v-model="form.transaction_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Amount ($)</label>
          <input v-model="form.amount" type="number" step="0.01" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Payment Method</label>
          <select v-model="form.payment_method" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option value="Cash">Cash</option>
            <option value="Bank">Bank</option>
            <option value="Mobile Banking">Mobile Banking</option>
            <option value="Other">Other</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Reference No</label>
          <input v-model="form.reference_no" type="text" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea v-model="form.description" rows="3" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
