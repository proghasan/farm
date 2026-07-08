<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboardStats } from '../api'

const stats = ref({
  totalAnimals: 0,
  totalSpecies: 0,
  totalInventory: 0,
  totalUsers: 0,
})
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await getDashboardStats()
    stats.value = data
  } finally {
    loading.value = false
  }
})

const statCards = [
  {
    label: 'Total Animals',
    value: 'totalAnimals',
    change: '+2',
    up: true,
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0l-3-3m3 3l3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" /></svg>`,
    color: 'bg-emerald-50 text-emerald-600',
  },
  {
    label: 'Species',
    value: 'totalSpecies',
    change: null,
    up: true,
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 0 0 3 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 0 0 5.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 0 0 9.568 3Z" /></svg>`,
    color: 'bg-blue-50 text-blue-600',
  },
  {
    label: 'Inventory Items',
    value: 'totalInventory',
    change: null,
    up: true,
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" /></svg>`,
    color: 'bg-amber-50 text-amber-600',
  },
  {
    label: 'Users',
    value: 'totalUsers',
    change: null,
    up: true,
    icon: `<svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" /></svg>`,
    color: 'bg-violet-50 text-violet-600',
  },
]
</script>

<template>
  <div class="space-y-6">
    <!-- Stats grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
      <div
        v-for="card in statCards"
        :key="card.label"
        class="bg-white rounded-2xl p-5 border border-gray-100 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" :class="card.color">
            <span class="w-5 h-5" v-html="card.icon" />
          </div>
          <span
            v-if="card.change"
            class="inline-flex items-center gap-0.5 text-xs font-medium px-2 py-1 rounded-lg"
            :class="card.up ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-500'"
          >
            <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 10.5 12 3m0 0 7.5 7.5M12 3v18" />
            </svg>
            {{ card.change }}
          </span>
        </div>
        <p class="text-2xl font-bold text-gray-900">{{ loading ? '...' : stats[card.value as keyof typeof stats] }}</p>
        <p class="text-sm text-gray-400 mt-0.5">{{ card.label }}</p>
      </div>
    </div>

    <!-- Quick actions -->
    <div class="bg-brand-600 rounded-2xl p-6 text-white">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <h3 class="font-semibold text-lg">Farm Operations</h3>
          <p class="text-brand-200 text-sm mt-0.5">Manage livestock, inventory, and finances from here.</p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <RouterLink
            to="/animals"
            class="px-4 py-2 bg-white/20 hover:bg-white/30 text-white text-sm font-medium rounded-xl transition-colors backdrop-blur-sm border border-white/20"
          >
            Add Animal
          </RouterLink>
          <RouterLink
            to="/inventory/items"
            class="px-4 py-2 bg-white text-brand-700 text-sm font-medium rounded-xl hover:bg-brand-50 transition-colors"
          >
            Inventory
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Quick links -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <RouterLink
        to="/vaccinations"
        class="bg-white rounded-2xl p-5 border border-gray-100 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center text-blue-600">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-900">Upcoming Vaccinations</p>
            <p class="text-xs text-gray-400">View scheduled vaccinations</p>
          </div>
        </div>
      </RouterLink>

      <RouterLink
        to="/accounting/transactions"
        class="bg-white rounded-2xl p-5 border border-gray-100 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl bg-emerald-50 flex items-center justify-center text-emerald-600">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0 1 15.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 0 1 3 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 0 0-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 0 1-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 0 0 3 15h-.75M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-900">Recent Transactions</p>
            <p class="text-xs text-gray-400">View financial records</p>
          </div>
        </div>
      </RouterLink>

      <RouterLink
        to="/users"
        class="bg-white rounded-2xl p-5 border border-gray-100 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl bg-violet-50 flex items-center justify-center text-violet-600">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-900">Team Members</p>
            <p class="text-xs text-gray-400">Manage farm staff</p>
          </div>
        </div>
      </RouterLink>
    </div>
  </div>
</template>