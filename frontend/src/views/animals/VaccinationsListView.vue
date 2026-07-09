<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed} from 'vue'
import { listVaccinations, createVaccination, deleteVaccination, listAnimals, listVaccines } from '../../api'
import type { Vaccination, Animal, Vaccine } from '../../api'
import DataTable from '../../components/DataTable.vue'
import Modal from '../../components/Modal.vue'
import PageHeader from '../../components/PageHeader.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'

const searchQuery = computed(() => headerStore.searchQuery)
const { success, error: showError } = useToast()
const headerStore = useHeaderStore()
const items = ref<Vaccination[]>([])
const animalList = ref<Animal[]>([])
const vaccineList = ref<Vaccine[]>([])
const loading = ref(false)
const showModal = ref(false)
const saving = ref(false)

const form = ref({
  animal_id: null as number | null,
  vaccine_id: null as number | null,
  vaccination_date: '',
  next_due_date: '',
  doctor_name: '',
  remarks: '',
})

async function fetchData() {
  loading.value = true
  try {
    const [vaccinations, animals, vaccines] = await Promise.all([
      listVaccinations(),
      listAnimals(),
      listVaccines(),
    ])
    items.value = vaccinations
    animalList.value = animals
    vaccineList.value = vaccines
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = {
    animal_id: null, vaccine_id: null,
    vaccination_date: '', next_due_date: '',
    doctor_name: '', remarks: '',
  }
  showModal.value = true
}

async function save() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      animal_id: form.value.animal_id!,
      vaccine_id: form.value.vaccine_id!,
      vaccination_date: form.value.vaccination_date,
    }
    if (form.value.next_due_date) payload.next_due_date = form.value.next_due_date
    if (form.value.doctor_name) payload.doctor_name = form.value.doctor_name
    if (form.value.remarks) payload.remarks = form.value.remarks

    await createVaccination(payload as any)
    success('Created', 'Vaccination record has been created')
    showModal.value = false
    await fetchData()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('Are you sure you want to delete this vaccination record?')) return
  try {
    await deleteVaccination(id)
    success('Deleted', 'Vaccination record has been deleted')
    await fetchData()
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Vaccinations' }])
  headerStore.setActions([{ label: 'Add New', onClick: openCreate }])
  headerStore.setShowSearch(true)
  fetchData()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div>
    <PageHeader title="Vaccinations" subtitle="Track vaccination history for animals" />
    <DataTable
      :columns="[
        { key: 'animal_display', label: 'Animal' },
        { key: 'vaccine_name', label: 'Vaccine' },
        { key: 'vaccination_date', label: 'Date' },
        { key: 'next_due_date', label: 'Due Date' },
        { key: 'doctor_name', label: 'Doctor' },
      ]"
      :items="items"
      :loading="loading"
      @delete="handleDelete"
    >
      <template #cell-animal_display="{ item }">
        {{ item.animal ? `${item.animal.tag_no}${item.animal.name ? ' - ' + item.animal.name : ''}` : '-' }}
      </template>
      <template #cell-vaccine_name="{ item }">
        {{ item.vaccine?.name || '-' }}
      </template>
      <template #cell-vaccination_date="{ item }">
        {{ item.vaccination_date ? new Date(item.vaccination_date).toLocaleDateString() : '-' }}
      </template>
      <template #cell-next_due_date="{ item }">
        {{ item.next_due_date ? new Date(item.next_due_date).toLocaleDateString() : '-' }}
      </template>
      <template #cell-doctor_name="{ item }">
        {{ item.doctor_name || '-' }}
      </template>
    </DataTable>

    <Modal :show="showModal" title="Add Vaccination" @close="showModal = false" size="md">
      <form @submit.prevent="save">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Animal</label>
            <select v-model="form.animal_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option :value="null" disabled>Select animal</option>
              <option v-for="a in animalList" :key="a.id" :value="a.id">{{ a.tag_no }}{{ a.name ? ' - ' + a.name : '' }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Vaccine</label>
            <select v-model="form.vaccine_id" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full">
              <option :value="null" disabled>Select vaccine</option>
              <option v-for="v in vaccineList" :key="v.id" :value="v.id">{{ v.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Vaccination Date</label>
            <input v-model="form.vaccination_date" type="date" required class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Next Due Date</label>
            <input v-model="form.next_due_date" type="date" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Doctor Name</label>
            <input v-model="form.doctor_name" type="text" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
          <div class="col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Remarks</label>
            <textarea v-model="form.remarks" rows="2" class="rounded-xl border border-gray-200 px-3.5 py-2.5 text-sm focus:border-brand-400 focus:ring-2 focus:ring-brand-500/20 outline-none transition-all w-full" />
          </div>
        </div>
      </form>
      <template #footer>
        <button @click="showModal = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50">Cancel</button>
        <button type="submit" :disabled="saving" class="px-4 py-2 text-sm font-medium text-white bg-brand-600 rounded-xl hover:bg-brand-700 disabled:opacity-50">{{ saving ? 'Saving...' : 'Create' }}</button>
      </template>
    </Modal>
  </div>
</template>
