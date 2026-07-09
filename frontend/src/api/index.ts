import api from './client'

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface User {
  id: number
  name: string
  email?: string
  phone?: string
  username?: string
  avatar?: string
  role: string
  status: string
  created_at: string
}

export interface Species {
  id: number
  name: string
  breeds?: Breed[]
}

export interface Breed {
  id: number
  species_id: number
  name: string
  species?: Species
}

export interface Animal {
  id: number
  tag_no: string
  species_id: number
  breed_id?: number
  father_id?: number
  mother_id?: number
  gender: string
  birth_date?: string
  purchase_date?: string
  purchase_price: number
  current_weight?: number
  color?: string
  status: string
  remarks?: string
  species?: Species
  breed?: Breed
  father?: Animal
  mother?: Animal
  weight_histories?: WeightHistory[]
  vaccinations?: Vaccination[]
  pregnancies?: Pregnancy[]
  sired_pregnancies?: Pregnancy[]
}

export interface Vaccine {
  id: number
  species_id: number
  name: string
  description?: string
  dose?: string
  minimum_age_value: number
  minimum_age_unit: string
  interval_value: number
  interval_unit: string
  is_repeatable: boolean
  species?: Species
}

export interface InventoryCategory {
  id: number
  name: string
}

export interface InventoryItem {
  id: number
  category_id: number
  name: string
  sku?: string
  unit: string
  purchase_price: number
  selling_price: number
  category?: InventoryCategory
}

export interface InventoryTransaction {
  id: number
  inventory_item_id: number
  transaction_type: string
  quantity: number
  transaction_date: string
  remarks?: string
  inventory_item?: InventoryItem
}

export interface AccountHead {
  id: number
  type: 'Income' | 'Expense'
  name: string
  description?: string
}

export interface AccountTransaction {
  id: number
  account_head_id: number
  transaction_date: string
  amount: number
  payment_method: string
  reference_no?: string
  description?: string
  account_head?: AccountHead
}

export interface Vaccination {
  id: number
  animal_id: number
  vaccine_id: number
  vaccination_date: string
  next_due_date?: string
  doctor_name?: string
  remarks?: string
  animal?: Animal
  vaccine?: Vaccine
}

export interface WeightHistory {
  id: number
  animal_id: number
  weight: number
  record_date: string
  remarks?: string
  animal?: Animal
}

export interface Pregnancy {
  id: number
  animal_id: number
  breeder_id?: number
  mating_date: string
  expected_due_date: string
  actual_birth_date?: string
  status: string
  note?: string
  number_of_children?: number
  number_of_male_children?: number
  number_of_female_children?: number
  number_of_dead_children?: number
  created_at: string
  animal?: Animal
  breeder?: Animal
}

// Auth
export const login = (data: AuthPayload) => api.post<LoginResponse>('/auth/login', data).then(r => r.data)
export const getProfile = () => api.get<User>('/users/profile').then(r => r.data)

// Species
export const listSpecies = () => api.get<{ data: Species[] }>('/species').then(r => r.data.data)
export const listSpeciesPaginated = (params?: Record<string, any>) => api.get<PaginatedResponse<Species>>('/species', { params }).then(r => r.data)
export const getSpecies = (id: number) => api.get<Species>(`/species/${id}`).then(r => r.data)
export const createSpecies = (data: Partial<Species>) => api.post('/species', data)
export const updateSpecies = (id: number, data: Partial<Species>) => api.put(`/species/${id}`, data)
export const deleteSpecies = (id: number) => api.delete(`/species/${id}`)

// Breeds
export const listBreeds = (params?: Record<string, any>) => api.get<{ data: Breed[] }>('/breeds', { params }).then(r => r.data.data)
export const createBreed = (data: Partial<Breed>) => api.post('/breeds', data)
export const updateBreed = (id: number, data: Partial<Breed>) => api.put(`/breeds/${id}`, data)
export const deleteBreed = (id: number) => api.delete(`/breeds/${id}`)

// Animals
export const listAnimals = (params?: Record<string, any>) => api.get<{ data: Animal[] }>('/animals', { params }).then(r => r.data.data)
export const listAnimalsPaginated = (params?: Record<string, any>) => api.get<PaginatedResponse<Animal>>('/animals', { params }).then(r => r.data)
export const getAnimal = (id: number) => api.get<Animal>(`/animals/${id}`).then(r => r.data)
export const getAnimalProfile = (id: number) => api.get<{ animal: Animal; pregnancies: Pregnancy[] }>(`/animals/${id}/profile`).then(r => r.data)
export const createAnimal = (data: Partial<Animal>) => api.post('/animals', data)
export const updateAnimal = (id: number, data: Partial<Animal>) => api.put(`/animals/${id}`, data)
export const deleteAnimal = (id: number) => api.delete(`/animals/${id}`)

// Vaccines
export const listVaccines = (params?: Record<string, any>) => api.get<{ data: Vaccine[] }>('/vaccines', { params }).then(r => r.data.data)
export const createVaccine = (data: Partial<Vaccine>) => api.post('/vaccines', data)
export const updateVaccine = (id: number, data: Partial<Vaccine>) => api.put(`/vaccines/${id}`, data)
export const deleteVaccine = (id: number) => api.delete(`/vaccines/${id}`)

// Vaccinations
export const listVaccinations = () => api.get<{ data: Vaccination[] }>('/vaccinations').then(r => r.data.data)
export const createVaccination = (data: Partial<Vaccination>) => api.post('/vaccinations', data)
export const deleteVaccination = (id: number) => api.delete(`/vaccinations/${id}`)

// Weight Histories
export const listWeightHistories = () => api.get<{ data: WeightHistory[] }>('/weight-histories').then(r => r.data.data)
export const createWeightHistory = (data: Partial<WeightHistory>) => api.post('/weight-histories', data)
export const deleteWeightHistory = (id: number) => api.delete(`/weight-histories/${id}`)

// Inventory Categories
export const listInventoryCategories = () => api.get<{ data: InventoryCategory[] }>('/inventory-categories').then(r => r.data.data)
export const createInventoryCategory = (data: Partial<InventoryCategory>) => api.post('/inventory-categories', data)
export const updateInventoryCategory = (id: number, data: Partial<InventoryCategory>) => api.put(`/inventory-categories/${id}`, data)
export const deleteInventoryCategory = (id: number) => api.delete(`/inventory-categories/${id}`)

// Inventory Items
export const listInventoryItems = () => api.get<{ data: InventoryItem[] }>('/inventory-items').then(r => r.data.data)
export const createInventoryItem = (data: Partial<InventoryItem>) => api.post('/inventory-items', data)
export const updateInventoryItem = (id: number, data: Partial<InventoryItem>) => api.put(`/inventory-items/${id}`, data)
export const deleteInventoryItem = (id: number) => api.delete(`/inventory-items/${id}`)

// Inventory Transactions
export const listInventoryTransactions = () => api.get<{ data: InventoryTransaction[] }>('/inventory-transactions').then(r => r.data.data)
export const createInventoryTransaction = (data: Partial<InventoryTransaction>) => api.post('/inventory-transactions', data)
export const deleteInventoryTransaction = (id: number) => api.delete(`/inventory-transactions/${id}`)

// Account Heads
export const listAccountHeads = () => api.get<{ data: AccountHead[] }>('/account-heads').then(r => r.data.data)
export const createAccountHead = (data: Partial<AccountHead>) => api.post('/account-heads', data)
export const updateAccountHead = (id: number, data: Partial<AccountHead>) => api.put(`/account-heads/${id}`, data)
export const deleteAccountHead = (id: number) => api.delete(`/account-heads/${id}`)

// Account Transactions
export const listAccountTransactions = () => api.get<{ data: AccountTransaction[] }>('/account-transactions').then(r => r.data.data)
export const createAccountTransaction = (data: Partial<AccountTransaction>) => api.post('/account-transactions', data)
export const deleteAccountTransaction = (id: number) => api.delete(`/account-transactions/${id}`)

// Users
export const listUsers = () => api.get<{ data: User[] }>('/users').then(r => r.data.data)
export const createUser = (data: Partial<User>) => api.post('/users', data)
export const updateUser = (id: number, data: Partial<User>) => api.put(`/users/${id}`, data)
export const deleteUser = (id: number) => api.delete(`/users/${id}`)

// Pregnancies
export const listPregnancies = () => api.get<{ data: Pregnancy[] }>('/pregnancies').then(r => r.data.data)
export const getPregnancy = (id: number) => api.get<Pregnancy>(`/pregnancies/${id}`).then(r => r.data)
export const createPregnancy = (data: Partial<Pregnancy>) => api.post('/pregnancies', data)
export const updatePregnancy = (id: number, data: Partial<Pregnancy>) => api.put(`/pregnancies/${id}`, data)
export const deletePregnancy = (id: number) => api.delete(`/pregnancies/${id}`)

// Dashboard stats
export const getDashboardStats = () =>
  Promise.all([
    listAnimals(),
    listSpecies(),
    listInventoryItems(),
    listAccountTransactions(),
  ]).then(([animals, species, items, txns]) => ({
    totalAnimals: animals.length,
    totalSpecies: species.length,
    totalInventory: items.length,
    transactions: txns.reduce((s, t) => s + t.amount, 0),
  }))

type AuthPayload = { login: string; password: string }