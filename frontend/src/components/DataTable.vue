<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
}

const props = withDefaults(
  defineProps<{
    columns: Column[];
    items: Record<string, any>[];
    loading?: boolean;
    searchable?: boolean;
    pageSize?: number;
  }>(),
  {
    loading: false,
    searchable: true,
    pageSize: 10,
  },
);

const emit = defineEmits<{
  edit: [id: number];
  delete: [id: number];
}>();

const search = ref("");
const sortKey = ref("");
const sortDir = ref<"asc" | "desc">("asc");
const currentPage = ref(1);
const pageSize = ref(props.pageSize);
const selected = ref<Set<number>>(new Set());
const openMenuId = ref<number | null>(null);

function toggleSort(key: string) {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === "asc" ? "desc" : "asc";
  } else {
    sortKey.value = key;
    sortDir.value = "asc";
  }
  currentPage.value = 1;
}

function toggleMenu(id: number) {
  openMenuId.value = openMenuId.value === id ? null : id;
}

function closeMenus(e: MouseEvent) {
  if (!(e.target as HTMLElement).closest(".row-menu")) {
    openMenuId.value = null;
  }
}

onMounted(() => document.addEventListener("click", closeMenus));
onUnmounted(() => document.removeEventListener("click", closeMenus));

const filtered = computed(() => {
  let data = props.items;
  if (search.value) {
    const q = search.value.toLowerCase();
    data = data.filter((item) =>
      props.columns.some((col) =>
        String(item[col.key] ?? "")
          .toLowerCase()
          .includes(q),
      ),
    );
  }
  if (sortKey.value) {
    data = [...data].sort((a, b) => {
      const av = String(a[sortKey.value] ?? "");
      const bv = String(b[sortKey.value] ?? "");
      return sortDir.value === "asc"
        ? av.localeCompare(bv)
        : bv.localeCompare(av);
    });
  }
  return data;
});

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filtered.value.length / pageSize.value)),
);

const paginated = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return filtered.value.slice(start, start + pageSize.value);
});

const visiblePages = computed(() => {
  const pages: (number | "...")[] = [];
  const total = totalPages.value;
  const cur = currentPage.value;
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  pages.push(1);
  if (cur > 3) pages.push("...");
  for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++)
    pages.push(i);
  if (cur < total - 2) pages.push("...");
  pages.push(total);
  return pages;
});

function setPage(p: number | "...") {
  if (typeof p === "number") currentPage.value = p;
}

function clearFilters() {
  search.value = "";
  currentPage.value = 1;
}
</script>

<template>
  <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-100">
            <th
              v-for="col in columns"
              :key="col.key"
              @click="col.sortable !== false && toggleSort(col.key)"
              class="px-4 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide whitespace-nowrap"
              :class="
                col.sortable !== false
                  ? 'cursor-pointer hover:text-gray-700 select-none'
                  : ''
              "
            >
              <div class="flex items-center gap-1.5">
                {{ col.label }}
                <span
                  v-if="col.sortable !== false"
                  class="flex flex-col gap-0.5"
                >
                  <svg
                    class="w-2.5 h-2.5 transition-colors"
                    :class="
                      sortKey === col.key && sortDir === 'asc'
                        ? 'text-brand-500'
                        : 'text-gray-300'
                    "
                    fill="currentColor"
                    viewBox="0 0 10 6"
                  >
                    <path d="M5 0l5 6H0z" />
                  </svg>
                  <svg
                    class="w-2.5 h-2.5 transition-colors"
                    :class="
                      sortKey === col.key && sortDir === 'desc'
                        ? 'text-brand-500'
                        : 'text-gray-300'
                    "
                    fill="currentColor"
                    viewBox="0 0 10 6"
                  >
                    <path d="M5 6l-5-6h10z" />
                  </svg>
                </span>
              </div>
            </th>
            <th
              class="px-4 py-3.5 text-right text-xs font-semibold text-gray-500 uppercase tracking-wide"
            >
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-50">
          <tr v-if="loading">
            <td
              :colspan="columns.length + 1"
              class="px-4 py-16 text-center text-sm text-gray-400"
            >
              Loading...
            </td>
          </tr>
          <tr
            v-else
            v-for="(item, idx) in paginated"
            :key="item.id"
            class="hover:bg-gray-50 transition-colors group"
          >
            <td v-for="col in columns" :key="col.key" class="px-4 py-3.5">
              <slot :name="`cell-${col.key}`" :item="item" :index="idx">
                <span class="text-sm text-gray-600">{{
                  item[col.key] ?? "-"
                }}</span>
              </slot>
            </td>
            <td class="px-4 py-3.5">
              <div class="row-menu relative flex justify-end">
                <button
                  @click="toggleMenu(item.id)"
                  class="w-7 h-7 flex items-center justify-center rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors opacity-0 group-hover:opacity-100"
                  :class="
                    openMenuId === item.id &&
                    'opacity-100 bg-gray-100 text-gray-600'
                  "
                >
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                  </svg>
                </button>

                <Transition name="dropdown">
                  <div
                    v-if="openMenuId === item.id"
                    class="absolute right-0 top-8 w-40 bg-white rounded-xl shadow-lg border border-gray-100 overflow-hidden z-20"
                  >
                    <button
                      @click="
                        emit('edit', item.id);
                        openMenuId = null;
                      "
                      class="w-full flex items-center gap-2.5 px-3.5 py-2.5 text-sm text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-colors"
                    >
                      <svg
                        class="w-4 h-4 text-gray-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125"
                        />
                      </svg>
                      Edit
                    </button>
                    <div class="border-t border-gray-100 mx-2" />
                    <button
                      @click="
                        emit('delete', item.id);
                        openMenuId = null;
                      "
                      class="w-full flex items-center gap-2.5 px-3.5 py-2.5 text-sm text-red-500 hover:bg-red-50 transition-colors"
                    >
                      <svg
                        class="w-4 h-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                        />
                      </svg>
                      Delete
                    </button>
                  </div>
                </Transition>
              </div>
            </td>
          </tr>
          <tr v-if="!loading && paginated.length === 0">
            <td :colspan="columns.length + 1" class="px-4 py-16 text-center">
              <div class="flex flex-col items-center gap-3">
                <div
                  class="w-12 h-12 bg-gray-100 rounded-xl flex items-center justify-center"
                >
                  <svg
                    class="w-6 h-6 text-gray-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
                    />
                  </svg>
                </div>
                <p class="text-sm font-medium text-gray-700">
                  No results found
                </p>
                <p class="text-xs text-gray-400">Try adjusting your search</p>
                <button
                  v-if="search"
                  @click="clearFilters"
                  class="text-xs text-brand-600 font-medium hover:text-brand-700 mt-1"
                >
                  Clear search
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div
      class="flex flex-col sm:flex-row items-center justify-between gap-3 px-5 py-4 border-t border-gray-100"
    >
      <div class="flex items-center gap-3 text-sm text-gray-500">
        <span>
          Showing
          <span class="font-medium text-gray-700">{{
            Math.min((currentPage - 1) * pageSize + 1, filtered.length)
          }}</span>
          –
          <span class="font-medium text-gray-700">{{
            Math.min(currentPage * pageSize, filtered.length)
          }}</span>
          of
          <span class="font-medium text-gray-700">{{ filtered.length }}</span>
        </span>
        <select
          v-model="pageSize"
          @change="currentPage = 1"
          class="border border-gray-200 rounded-lg px-2 py-1 text-xs text-gray-600 bg-gray-50 outline-none focus:ring-1 focus:ring-brand-300"
        >
          <option :value="10">10</option>
          <option :value="20">20</option>
          <option :value="50">50</option>
          <option :value="100">100</option>
        </select>
      </div>

      <div class="flex items-center gap-1">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="w-8 h-8 flex items-center justify-center rounded-lg border border-gray-200 text-gray-500 hover:bg-gray-50 hover:text-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all text-sm"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15.75 19.5 8.25 12l7.5-7.5"
            />
          </svg>
        </button>
        <button
          v-for="p in visiblePages"
          :key="String(p)"
          @click="setPage(p)"
          class="min-w-8 h-8 px-1.5 flex items-center justify-center rounded-lg text-sm font-medium transition-all"
          :class="
            p === '...'
              ? 'cursor-default text-gray-400'
              : p === currentPage
                ? 'bg-brand-600 text-white shadow-sm'
                : 'border border-gray-200 text-gray-600 hover:bg-gray-50 hover:text-gray-900'
          "
        >
          {{ p }}
        </button>
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="w-8 h-8 flex items-center justify-center rounded-lg border border-gray-200 text-gray-500 hover:bg-gray-50 hover:text-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m8.25 4.5 7.5 7.5-7.5 7.5"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dropdown-enter-active {
  transition: all 0.12s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.dropdown-leave-active {
  transition: all 0.08s ease-in;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
