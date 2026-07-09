import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { BreadcrumbItem } from '../components/Breadcrumb.vue'

export interface HeaderAction {
  label: string
  onClick: () => void
}

export const useHeaderStore = defineStore('header', () => {
  const searchQuery = ref('')
  const showSearch = ref(false)
  const actions = ref<HeaderAction[]>([])
  const breadcrumb = ref<BreadcrumbItem[]>([])

  function setActions(a: HeaderAction[]) {
    actions.value = a
  }

  function setShowSearch(v: boolean) {
    showSearch.value = v
  }

  function setBreadcrumb(items: BreadcrumbItem[]) {
    breadcrumb.value = items
  }

  function clear() {
    searchQuery.value = ''
    showSearch.value = false
    actions.value = []
    breadcrumb.value = []
  }

  return { searchQuery, showSearch, actions, breadcrumb, setActions, setShowSearch, setBreadcrumb, clear }
})