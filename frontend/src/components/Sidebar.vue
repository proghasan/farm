<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'

interface NavItem {
  label: string
  icon: string
  path?: string
  children?: { label: string; path: string; icon: string }[]
}

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ toggle: [] }>()

const route = useRoute()

const navItems: NavItem[] = [
  {
    label: 'Dashboard',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z" /></svg>`,
    path: '/',
  },
  {
    label: 'Livestock',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0l-3-3m3 3l3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" /></svg>`,
    children: [
      {
        label: 'Animals',
        path: '/animals',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 18.75a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 0 1-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 0 0-3.213-9.193 2.056 2.056 0 0 0-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 0 0-10.026 0 1.106 1.106 0 0 0-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12" /></svg>`,
      },
      {
        label: 'Species',
        path: '/species',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 0 0 3 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 0 0 5.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 0 0 9.568 3Z" /><path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6Z" /></svg>`,
      },
      {
        label: 'Breeds',
        path: '/breeds',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" /></svg>`,
      },
      {
        label: 'Vaccines',
        path: '/vaccines',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 0 0-3.7-3.7 48.678 48.678 0 0 0-7.324 0 4.006 4.006 0 0 0-3.7 3.7c-.092 1.21-.139 2.43-.139 3.662s.046 2.453.138 3.662a4.006 4.006 0 0 0 3.7 3.7 48.657 48.657 0 0 0 7.324 0 4.006 4.006 0 0 0 3.7-3.7c.092-1.21.138-2.43.138-3.662Z" /></svg>`,
      },
      {
        label: 'Vaccinations',
        path: '/vaccinations',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>`,
      },
    ],
  },
  {
    label: 'Inventory',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" /></svg>`,
    children: [
      {
        label: 'Categories',
        path: '/inventory/categories',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 0 0 3 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.699 2.45.33a18.095 18.095 0 0 0 5.08-5.13c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 0 0 9.568 3Z" /></svg>`,
      },
      {
        label: 'Items',
        path: '/inventory/items',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 0 1-2.247 2.118H6.622a2.25 2.25 0 0 1-2.247-2.118L3.75 7.5m10 11.25v-6m-6 6v-6m-3.75-6h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125Z" /></svg>`,
      },
      {
        label: 'Transactions',
        path: '/inventory/transactions',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M7.5 21 3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" /></svg>`,
      },
    ],
  },
  {
    label: 'Accounting',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0 1 15.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 0 1 3 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 0 0-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 0 1-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 0 0 3 15h-.75M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm3 0h.008v.008H18V10.5Zm-12 0h.008v.008H6V10.5Z" /></svg>`,
    children: [
      {
        label: 'Account Heads',
        path: '/accounting/heads',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" /></svg>`,
      },
      {
        label: 'Transactions',
        path: '/accounting/transactions',
        icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3 4.5h14.25M3 9h9.75M3 13.5h9.75m4.5-4.5v12m0 0-3.75-3.75M9.75 21 21 21" /></svg>`,
      },
    ],
  },
  {
    label: 'Users',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" /></svg>`,
    path: '/users',
  },
]

const expandedMenus = ref<Set<string>>(new Set())

const activeParent = computed(() =>
  navItems.find(i => i.children?.some(c => c.path === route.path))?.label
)

onMounted(() => {
  const parent = activeParent.value
  if (parent) expandedMenus.value.add(parent)
})

watch(activeParent, (parent) => {
  if (parent) expandedMenus.value.add(parent)
})

function toggleMenu(label: string) {
  if (expandedMenus.value.has(label)) {
    expandedMenus.value.delete(label)
  } else {
    expandedMenus.value.add(label)
  }
}

function isActive(path: string) {
  return route.path === path
}

function isParentActive(item: NavItem) {
  return item.children?.some(c => route.path === c.path)
}
</script>

<template>
  <aside
    class="fixed left-0 top-0 h-full bg-white border-r border-gray-100 z-30 flex flex-col transition-all duration-300 shadow-sm"
    :class="open ? 'w-64' : 'w-16'"
  >
    <!-- Logo -->
    <div class="flex items-center h-16 px-4 border-b border-gray-100 shrink-0">
      <div class="flex items-center gap-3 min-w-0">
        <div class="w-8 h-8 rounded-xl bg-brand-600 flex items-center justify-center shrink-0">
          <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z" />
          </svg>
        </div>
        <Transition name="fade">
          <span v-if="open" class="font-semibold text-gray-900 text-sm tracking-tight whitespace-nowrap">Farm Manager</span>
        </Transition>
      </div>
    </div>

    <!-- Nav -->
    <nav class="flex-1 py-4 overflow-y-auto overflow-x-hidden">
      <ul class="space-y-0.5 px-2">
        <li v-for="item in navItems" :key="item.label">

          <!-- Parent with children -->
          <template v-if="item.children">
            <button
              @click="open && toggleMenu(item.label)"
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 group"
              :class="isParentActive(item)
                ? 'text-brand-600 bg-brand-50'
                : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'"
            >
              <span class="w-5 h-5 shrink-0" v-html="item.icon" />
              <Transition name="fade">
                <span v-if="open" class="flex-1 text-left whitespace-nowrap">{{ item.label }}</span>
              </Transition>
              <Transition name="fade">
                <svg
                  v-if="open"
                  class="w-4 h-4 transition-transform duration-200"
                  :class="expandedMenus.has(item.label) ? 'rotate-90' : ''"
                  fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                </svg>
              </Transition>
            </button>

            <Transition name="expand">
              <ul v-if="open && expandedMenus.has(item.label)" class="mt-0.5 ml-4 space-y-0.5 border-l border-gray-100 pl-3">
                <li v-for="child in item.children" :key="child.label">
                  <RouterLink
                    :to="child.path"
                    class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm transition-all duration-150"
                    :class="isActive(child.path)
                      ? 'text-brand-600 bg-brand-50 font-medium'
                      : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'"
                  >
                    <span class="w-4 h-4 shrink-0" v-html="child.icon" />
                    <span class="whitespace-nowrap">{{ child.label }}</span>
                  </RouterLink>
                </li>
              </ul>
            </Transition>
          </template>

          <!-- Single link -->
          <RouterLink
            v-else
            :to="item.path!"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 group"
            :class="isActive(item.path!)
              ? 'text-brand-600 bg-brand-50'
              : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'"
          >
            <span class="w-5 h-5 shrink-0" v-html="item.icon" />
            <Transition name="fade">
              <span v-if="open" class="whitespace-nowrap">{{ item.label }}</span>
            </Transition>
          </RouterLink>

        </li>
      </ul>
    </nav>

    <!-- Collapse toggle -->
    <button
      @click="emit('toggle')"
      class="absolute -right-3 top-[4.5rem] w-6 h-6 bg-white border border-gray-200 rounded-full flex items-center justify-center shadow-sm hover:shadow-md transition-all text-gray-400 hover:text-gray-700"
    >
      <svg
        class="w-3 h-3 transition-transform duration-300"
        :class="open ? '' : 'rotate-180'"
        fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
      </svg>
    </button>
  </aside>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}
.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}
.expand-enter-to,
.expand-leave-from {
  max-height: 200px;
}
</style>