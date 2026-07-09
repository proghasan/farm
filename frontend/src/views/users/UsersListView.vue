<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted} from 'vue'
import { DataTable } from '../../components/DataTable'
import RowActions from '../../components/RowActions.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import UserStatusBadge from '../../components/user/UserStatusBadge.vue'
import {
  listUsers,
  createUser,
  updateUser,
  deleteUser,
} from '../../api'
import type { User } from '../../api'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'
import { getErrorMessage } from "../../utils/error";

const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const items = ref<User[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  email: '',
  phone: '',
  username: '',
  password: '',
  role: 'Worker',
  status: 'Active',
})

const columns = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'email', label: 'Email' },
  { key: 'phone', label: 'Phone' },
  { key: 'role', label: 'Role', component: UserStatusBadge },
  { key: 'status', label: 'Status', component: UserStatusBadge },
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

const modalTitle = computed(() => (editingId.value ? 'Edit User' : 'Add User'))

async function fetchItems() {
  loading.value = true
  try {
    items.value = await listUsers()
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = {
    name: '', email: '', phone: '', username: '',
    password: '', role: 'Worker', status: 'Active',
  }
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
  form.value = {
    name: item.name,
    email: item.email ?? '',
    phone: item.phone ?? '',
    username: item.username ?? '',
    password: '',
    role: item.role,
    status: item.status,
  }
  showModal.value = true
}

async function handleSave() {
  try {
    const payload: Record<string, any> = {
      name: form.value.name,
      email: form.value.email || undefined,
      phone: form.value.phone || undefined,
      username: form.value.username || undefined,
      role: form.value.role,
      status: form.value.status,
    }
    if (editingId.value) {
      if (form.value.password) {
        payload.password = form.value.password
      }
      await updateUser(editingId.value, payload)
      success('Updated', 'User has been updated')
    } else {
      payload.password = form.value.password
      await createUser(payload)
      success('Created', 'User has been created')
    }
    showModal.value = false
    await fetchItems()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this user?')) return
  try {
    await deleteUser(id)
    success('Deleted', 'User has been deleted')
    await fetchItems()
  } catch (e: any) {
    showError('Failed', getErrorMessage(e))
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Users' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  fetchItems()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Users" subtitle="Manage system users and permissions" />
    <DataTable
      :columns="columns"
      :items="items"
      :loading="loading"
    >
      <template #cell-email="{ item }">
        {{ item.email || '-' }}
      </template>
      <template #cell-phone="{ item }">
        {{ item.phone || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" :title="modalTitle" @close="showModal = false" size="md">
      <form @submit.prevent="handleSave" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input v-model="form.name" type="text" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="form.email" type="email" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
          <input v-model="form.phone" type="text" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input v-model="form.username" type="text" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Password
            <span v-if="editingId" class="text-xs text-gray-400">(leave empty to keep)</span>
            <span v-else class="text-xs text-red-500">*</span>
          </label>
          <input v-model="form.password" :required="!editingId" type="password" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
          <select v-model="form.role" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option value="Owner">Owner</option>
            <option value="Manager">Manager</option>
            <option value="Veterinarian">Veterinarian</option>
            <option value="Worker">Worker</option>
            <option value="Accountant">Accountant</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select v-model="form.status" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
            <option value="Active">Active</option>
            <option value="Inactive">Inactive</option>
            <option value="Suspended">Suspended</option>
          </select>
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700">{{ editingId ? 'Update' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
