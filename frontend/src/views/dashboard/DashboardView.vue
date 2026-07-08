<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboardStats } from '../../api'

const stats = ref({
  totalAnimals: 0,
  totalSpecies: 0,
  totalInventory: 0,
  transactions: 0,
})
const loading = ref(true)

onMounted(async () => {
  try {
    stats.value = await getDashboardStats()
  } finally {
    loading.value = false
  }
})

const statCards = [
  { label: 'Total Animals', value: 'totalAnimals', color: 'bg-farm-500', icon: '🐄' },
  { label: 'Species', value: 'totalSpecies', color: 'bg-blue-500', icon: '🔬' },
  { label: 'Inventory Items', value: 'totalInventory', color: 'bg-amber-500', icon: '📦' },
  { label: 'Total Revenue', value: 'transactions', color: 'bg-emerald-500', icon: '💰', prefix: '$' },
]
</script>

<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p class="text-sm text-gray-500">Overview of your farm operations</p>
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
      <div v-for="i in 4" :key="i" class="h-28 animate-pulse rounded-xl bg-gray-200" />
    </div>

    <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="card in statCards"
        :key="card.label"
        class="rounded-xl bg-white p-6 shadow-sm ring-1 ring-gray-200"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500">{{ card.label }}</p>
            <p class="mt-1 text-3xl font-bold text-gray-900">
              {{ card.prefix || '' }}{{ stats[card.value as keyof typeof stats] }}
            </p>
          </div>
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-50 text-2xl">
            {{ card.icon }}
          </div>
        </div>
      </div>
    </div>

    <div class="mt-8">
      <h2 class="mb-4 text-lg font-semibold text-gray-900">Quick Actions</h2>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <router-link
          to="/animals"
          class="rounded-xl border border-gray-200 bg-white p-5 transition-colors hover:border-farm-500 hover:shadow-sm"
        >
          <p class="font-medium text-gray-900">Manage Animals</p>
          <p class="mt-1 text-sm text-gray-500">Add, edit, or view livestock records</p>
        </router-link>
        <router-link
          to="/inventory/items"
          class="rounded-xl border border-gray-200 bg-white p-5 transition-colors hover:border-farm-500 hover:shadow-sm"
        >
          <p class="font-medium text-gray-900">Inventory</p>
          <p class="mt-1 text-sm text-gray-500">Track supplies and stock levels</p>
        </router-link>
        <router-link
          to="/accounting/transactions"
          class="rounded-xl border border-gray-200 bg-white p-5 transition-colors hover:border-farm-500 hover:shadow-sm"
        >
          <p class="font-medium text-gray-900">Accounting</p>
          <p class="mt-1 text-sm text-gray-500">Record income and expenses</p>
        </router-link>
      </div>
    </div>
  </div>
</template>