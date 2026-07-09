<script setup lang="ts">
import { ref, computed } from "vue";
import type { Column } from "./types";

const props = withDefaults(
  defineProps<{
    columns: Column[];
    items: Record<string, any>[];
    loading?: boolean;
    pageSize?: number;
    serverMode?: boolean;
    totalItems?: number;
    currentPage?: number;
  }>(),
  {
    loading: false,
    pageSize: 10,
    serverMode: false,
    totalItems: 0,
    currentPage: 1,
  },
);

const emit = defineEmits<{
  "update:currentPage": [page: number];
  "update:pageSize": [size: number];
  sort: [payload: { key: string; direction: "asc" | "desc" }];
}>();

const sortKey = ref("");
const sortDir = ref<"asc" | "desc">("asc");
const localPage = ref(1);
const localPageSize = ref(props.pageSize);

const effectivePage = computed(() =>
  props.serverMode ? props.currentPage : localPage.value,
);
const effectivePageSize = computed(() =>
  props.serverMode ? props.pageSize : localPageSize.value,
);

function setPage(p: number | "...") {
  if (typeof p !== "number") return;
  if (props.serverMode) {
    emit("update:currentPage", p);
  } else {
    localPage.value = p;
  }
}

function setPageSize(size: number) {
  if (props.serverMode) {
    emit("update:pageSize", size);
  } else {
    localPageSize.value = size;
  }
  localPage.value = 1;
}

function toggleSort(key: string) {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === "asc" ? "desc" : "asc";
  } else {
    sortKey.value = key;
    sortDir.value = "asc";
  }
  if (props.serverMode) {
    emit("sort", { key: sortKey.value, direction: sortDir.value });
  } else {
    localPage.value = 1;
  }
}

const sorted = computed(() => {
  if (props.serverMode) return props.items;
  let data = props.items;
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

const totalItems_ = computed(() =>
  props.serverMode ? props.totalItems : sorted.value.length,
);

const totalPages = computed(() =>
  Math.max(1, Math.ceil(totalItems_.value / effectivePageSize.value)),
);

const paginated = computed(() => {
  if (props.serverMode) return props.items;
  const start = (effectivePage.value - 1) * effectivePageSize.value;
  return sorted.value.slice(start, start + effectivePageSize.value);
});

const visiblePages = computed(() => {
  const pages: (number | "...")[] = [];
  const total = totalPages.value;
  const cur = effectivePage.value;
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  pages.push(1);
  if (cur > 3) pages.push("...");
  for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++)
    pages.push(i);
  if (cur < total - 2) pages.push("...");
  pages.push(total);
  return pages;
});

const displayFrom = computed(
  () => (effectivePage.value - 1) * effectivePageSize.value + 1,
);

const displayTo = computed(() =>
  Math.min(effectivePage.value * effectivePageSize.value, totalItems_.value),
);
</script>

<template>
  <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
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
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-50">
          <tr v-if="loading">
            <td
              :colspan="columns.length"
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
              <template v-if="col.component">
                <component
                  :is="col.component"
                  :item="item"
                  :value="item[col.key]"
                  v-bind="col.componentProps || {}"
                />
              </template>
              <slot v-else :name="`cell-${col.key}`" :item="item" :index="idx">
                <span class="text-sm text-gray-600">{{
                  item[col.key] ?? "-"
                }}</span>
              </slot>
            </td>
          </tr>
          <tr v-if="!loading && paginated.length === 0">
            <td
              :colspan="columns.length"
              class="px-4 py-16 text-center"
            >
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
                <p class="text-xs text-gray-400">
                  Try adjusting your search
                </p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      class="flex flex-col sm:flex-row items-center justify-between gap-3 px-5 py-4 border-t border-gray-100"
    >
      <div class="flex items-center gap-3 text-sm text-gray-500">
        <span>
          Showing
          <span class="font-medium text-gray-700">{{ displayFrom }}</span>
          –
          <span class="font-medium text-gray-700">{{ displayTo }}</span>
          of
          <span class="font-medium text-gray-700">{{ totalItems_ }}</span>
        </span>
        <select
          :value="effectivePageSize"
          @change="setPageSize(Number(($event.target as HTMLSelectElement).value))"
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
          @click="setPage(effectivePage - 1)"
          :disabled="effectivePage === 1"
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
              : p === effectivePage
                ? 'bg-brand-600 text-white shadow-sm'
                : 'border border-gray-200 text-gray-600 hover:bg-gray-50 hover:text-gray-900'
          "
        >
          {{ p }}
        </button>
        <button
          @click="setPage(effectivePage + 1)"
          :disabled="effectivePage === totalPages"
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
