import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AuthLayout from '../layouts/AuthLayout.vue'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import LoginView from '../views/auth/LoginView.vue'
import Dashboard from '../views/Dashboard.vue'
import AnimalsListView from '../views/animals/AnimalsListView.vue'
import AnimalProfileView from '../views/animals/AnimalProfileView.vue'
import SpeciesListView from '../views/animals/SpeciesListView.vue'
import BreedsListView from '../views/animals/BreedsListView.vue'
import VaccinesListView from '../views/animals/VaccinesListView.vue'
import VaccinationsListView from '../views/animals/VaccinationsListView.vue'
import PregnanciesListView from '../views/animals/PregnanciesListView.vue'
import InventoryCategoriesView from '../views/inventory/InventoryCategoriesView.vue'
import InventoryItemsView from '../views/inventory/InventoryItemsView.vue'
import InventoryTransactionsView from '../views/inventory/InventoryTransactionsView.vue'
import AccountHeadsView from '../views/accounting/AccountHeadsView.vue'
import AccountTransactionsView from '../views/accounting/AccountTransactionsView.vue'
import UsersListView from '../views/users/UsersListView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      component: AuthLayout,
      meta: { guest: true },
      children: [
        { path: '', component: LoginView, meta: { title: 'Login' } },
      ],
    },
    {
      path: '/',
      component: DefaultLayout,
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: Dashboard, meta: { title: 'Dashboard' } },
        { path: 'animals', component: AnimalsListView, meta: { title: 'Animals' } },
        { path: 'animals/:id', component: AnimalProfileView, meta: { title: 'Animal Profile' } },
        { path: 'species', component: SpeciesListView, meta: { title: 'Species' } },
        { path: 'breeds', component: BreedsListView, meta: { title: 'Breeds' } },
        { path: 'vaccines', component: VaccinesListView, meta: { title: 'Vaccines' } },
        { path: 'vaccinations', component: VaccinationsListView, meta: { title: 'Vaccinations' } },
        { path: 'pregnancies', component: PregnanciesListView, meta: { title: 'Pregnancies' } },
        { path: 'inventory/categories', component: InventoryCategoriesView, meta: { title: 'Inventory Categories' } },
        { path: 'inventory/items', component: InventoryItemsView, meta: { title: 'Inventory Items' } },
        { path: 'inventory/transactions', component: InventoryTransactionsView, meta: { title: 'Inventory Transactions' } },
        { path: 'accounting/heads', component: AccountHeadsView, meta: { title: 'Account Heads' } },
        { path: 'accounting/transactions', component: AccountTransactionsView, meta: { title: 'Account Transactions' } },
        { path: 'users', component: UsersListView, meta: { title: 'Users' } },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) next('/login')
  else if (to.meta.guest && auth.isAuthenticated) next('/dashboard')
  else next()
})

export default router