<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAnimalProfile, deleteWeightHistory } from '../../api'
import type { Animal, Pregnancy } from '../../api'
import PageHeader from '../../components/PageHeader.vue'
import Modal from '../../components/Modal.vue'
import DateDisplay from '../../components/DateDisplay.vue'
import { useToast } from '../../composables/useToast'
import { useHeaderStore } from '../../stores/header'

const route = useRoute()
const router = useRouter()
const headerStore = useHeaderStore()
const { success, error: showError } = useToast()

const animal = ref<Animal | null>(null)
const pregnancies = ref<Pregnancy[]>([])
const loading = ref(true)
const showWeightForm = ref(false)
const weightForm = ref({ weight: 0, record_date: '', remarks: '' })
const weightSaving = ref(false)

const showTab = ref<'history' | 'vaccinations' | 'pregnancies'>('history')

async function load() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const profile = await getAnimalProfile(id)
    animal.value = profile.animal
    pregnancies.value = profile.pregnancies
  } finally {
    loading.value = false
  }
}

async function deleteWeight(id: number) {
  if (!confirm('Delete this weight record?')) return
  try {
    await deleteWeightHistory(id)
    const profile = await getAnimalProfile(animal.value!.id)
    animal.value = profile.animal
    pregnancies.value = profile.pregnancies
    success('Deleted', 'Weight record removed')
  } catch (e: any) {
    showError('Failed', e?.response?.data?.message || 'An error occurred')
  }
}

function statusBadge(status?: string) {
  if (!status) return ''
  const colors: Record<string, string> = {
    Active: 'bg-emerald-50 text-emerald-700 ring-emerald-200',
    Sold: 'bg-blue-50 text-blue-700 ring-blue-200',
    Deceased: 'bg-red-50 text-red-700 ring-red-200',
    Healthy: 'bg-emerald-50 text-emerald-700 ring-emerald-200',
    Sick: 'bg-amber-50 text-amber-700 ring-amber-200',
  }
  return colors[status] || 'bg-gray-50 text-gray-600 ring-gray-200'
}

onMounted(() => {
  headerStore.setBreadcrumb([{ label: 'Dashboard', to: '/dashboard' }, { label: 'Animals', to: '/animals' }, { label: 'Profile' }])
  load()
})
onUnmounted(() => headerStore.clear())
</script>

<template>
  <div v-if="loading" class="text-center py-16 text-gray-400">Loading...</div>
  <div v-else-if="!animal" class="text-center py-16 text-gray-400">Animal not found</div>
  <div v-else class="max-w-4xl">
    <div class="flex items-start justify-between mb-6">
      <div>
        <PageHeader :title="animal.tag_no" :subtitle="`${animal.breed?.species?.name || ''}${animal.breed ? ' / ' + animal.breed.name : ''}`" />
      </div>
      <button @click="router.push('/animals')" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors">Back to List</button>
    </div>

    <!-- Info Cards -->
    <div class="grid grid-cols-3 gap-4 mb-8">
      <div class="bg-white rounded-2xl border border-gray-100 p-5">
        <p class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2">Status</p>
        <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium ring-1" :class="statusBadge(animal.status)">
          <span class="w-1.5 h-1.5 rounded-full" :class="animal.status === 'Active' || animal.status === 'Healthy' ? 'bg-emerald-500' : animal.status === 'Deceased' ? 'bg-red-500' : 'bg-amber-500'" />
          {{ animal.status }}
        </span>
      </div>
      <div class="bg-white rounded-2xl border border-gray-100 p-4">
        <p class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-1">Gender</p>
        <p class="text-sm font-semibold text-gray-900">{{ animal.gender }}</p>
      </div>
      <div class="bg-white rounded-2xl border border-gray-100 p-4">
        <p class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-1">Current Weight</p>
        <p class="text-sm font-semibold text-gray-900">{{ animal.current_weight ? animal.current_weight + ' kg' : 'Not recorded' }}</p>
      </div>
    </div>

    <!-- Details Grid -->
    <div class="bg-white rounded-2xl border border-gray-100 p-6 mb-8">
      <h3 class="text-sm font-semibold text-gray-900 mb-4">Details</h3>
      <div class="grid grid-cols-2 gap-x-8 gap-y-3">
        <div class="flex justify-between"><span class="text-sm text-gray-500">Tag No</span><span class="text-sm font-medium text-gray-900">{{ animal.tag_no }}</span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Species</span><span class="text-sm font-medium text-gray-900">{{ animal.breed?.species?.name || '-' }}</span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Breed</span><span class="text-sm font-medium text-gray-900">{{ animal.breed?.name || '-' }}</span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Birth Date</span><span class="text-sm font-medium text-gray-900"><DateDisplay :value="animal.birth_date" /></span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Purchase Date</span><span class="text-sm font-medium text-gray-900"><DateDisplay :value="animal.purchase_date" /></span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Purchase Price</span><span class="text-sm font-medium text-gray-900">${{ Number(animal.purchase_price).toFixed(2) }}</span></div>
        <div class="flex justify-between"><span class="text-sm text-gray-500">Color</span><span class="text-sm font-medium text-gray-900">{{ animal.color || '-' }}</span></div>
        <div v-if="animal.father" class="flex justify-between col-span-2"><span class="text-sm text-gray-500">Father</span><span class="text-sm font-medium text-gray-900">{{ animal.father.tag_no }}</span></div>
        <div v-if="animal.mother" class="flex justify-between col-span-2"><span class="text-sm text-gray-500">Mother</span><span class="text-sm font-medium text-gray-900">{{ animal.mother.tag_no }}</span></div>
        <div v-if="animal.remarks" class="col-span-2"><span class="text-sm text-gray-500">Remarks</span><p class="text-sm text-gray-900 mt-1">{{ animal.remarks }}</p></div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="border-b border-gray-100 mb-6">
      <nav class="flex gap-6">
        <button @click="showTab = 'history'" class="pb-3 text-sm font-medium border-b-2 transition-colors" :class="showTab === 'history' ? 'text-brand-600 border-brand-600' : 'text-gray-500 border-transparent hover:text-gray-700'">Weight History</button>
        <button @click="showTab = 'vaccinations'" class="pb-3 text-sm font-medium border-b-2 transition-colors" :class="showTab === 'vaccinations' ? 'text-brand-600 border-brand-600' : 'text-gray-500 border-transparent hover:text-gray-700'">Vaccinations</button>
        <button @click="showTab = 'pregnancies'" class="pb-3 text-sm font-medium border-b-2 transition-colors" :class="showTab === 'pregnancies' ? 'text-brand-600 border-brand-600' : 'text-gray-500 border-transparent hover:text-gray-700'">Pregnancies</button>
      </nav>
    </div>

    <!-- Weight History Tab -->
    <div v-if="showTab === 'history'" class="bg-white rounded-2xl border border-gray-100 p-4">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-sm font-semibold text-gray-900">Weight History</h3>
      </div>
      <div class="overflow-x-auto rounded-xl border border-gray-100">
        <table class="w-full text-left text-sm">
          <thead class="border-b bg-gray-50"><tr><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Weight</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Date</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Remarks</th><th class="px-3 py-2 text-right">Action</th></tr></thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="!animal.weight_histories?.length"><td colspan="4" class="px-3 py-8 text-center text-gray-400">No weight records.</td></tr>
            <tr v-for="w in animal.weight_histories" :key="w.id">
              <td class="px-3 py-2">{{ w.weight }} kg</td>
              <td class="px-3 py-2"><DateDisplay :value="w.record_date" /></td>
              <td class="px-3 py-2">{{ w.remarks || '-' }}</td>
              <td class="px-3 py-2 text-right"><button @click="deleteWeight(w.id)" class="text-xs font-medium text-red-500 hover:text-red-600">Delete</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Vaccinations Tab -->
    <div v-if="showTab === 'vaccinations'" class="bg-white rounded-2xl border border-gray-100 p-4">
      <h3 class="text-sm font-semibold text-gray-900 mb-4">Vaccination History</h3>
      <div class="overflow-x-auto rounded-xl border border-gray-100">
        <table class="w-full text-left text-sm">
          <thead class="border-b bg-gray-50"><tr><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Vaccine</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Date</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Due Date</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Doctor</th></tr></thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="!animal.vaccinations?.length"><td colspan="4" class="px-3 py-8 text-center text-gray-400">No vaccination records.</td></tr>
            <tr v-for="v in animal.vaccinations" :key="v.id">
              <td class="px-3 py-2">{{ v.vaccine?.name || '-' }}</td>
              <td class="px-3 py-2"><DateDisplay :value="v.vaccination_date" /></td>
              <td class="px-3 py-2"><DateDisplay :value="v.next_due_date" /></td>
              <td class="px-3 py-2">{{ v.doctor_name || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pregnancies Tab -->
    <div v-if="showTab === 'pregnancies'" class="bg-white rounded-2xl border border-gray-100 p-4">
      <h3 class="text-sm font-semibold text-gray-900 mb-4">Pregnancy History</h3>
      <div class="overflow-x-auto rounded-xl border border-gray-100">
        <table class="w-full text-left text-sm">
          <thead class="border-b bg-gray-50"><tr><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Mating Date</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Due Date</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Status</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Offspring</th><th class="px-3 py-2 font-medium text-gray-500 text-xs uppercase">Breeder</th></tr></thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-if="!pregnancies.length"><td colspan="5" class="px-3 py-8 text-center text-gray-400">No pregnancy records.</td></tr>
            <tr v-for="p in pregnancies" :key="p.id">
              <td class="px-3 py-2"><DateDisplay :value="p.mating_date" /></td>
              <td class="px-3 py-2"><DateDisplay :value="p.expected_due_date" /></td>
              <td class="px-3 py-2"><span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-medium ring-1" :class="p.status === 'Delivered' ? 'bg-emerald-50 text-emerald-700 ring-emerald-200' : p.status === 'Pregnant' ? 'bg-blue-50 text-blue-700 ring-blue-200' : 'bg-gray-50 text-gray-600 ring-gray-200'">{{ p.status }}</span></td>
              <td class="px-3 py-2">{{ p.number_of_children ?? '-' }}</td>
              <td class="px-3 py-2">{{ p.breeder?.tag_no || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>