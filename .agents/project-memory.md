# Farm Management System — Project Memory

## Overview
A full-stack farm management web application with a Go/Fiber REST API backend and Vue 3 + Tailwind CSS v4 frontend. Enables livestock tracking, inventory management, accounting, and vaccination records for farm operations.

## Tech Stack
- **Language:** Go 1.26
- **Framework:** Fiber v3
- **Database:** MySQL 8+ via GORM (with AutoMigrate)
- **Auth:** JWT (HS256, 72h expiry), bcrypt passwords
- **Validation:** go-playground/validator v10
- **Frontend:** Vue 3 + Pinia + Vue Router + Axios + Tailwind CSS v4 + Vite
- **Dev tooling:** air (hot reload)

## Project Structure

```
farm/
├── main.go              # Entry point: config → DB connect → migrate → seed admin → routes → listen
├── config/config.go     # Env-based config (DB, JWT, port)
├── database/
│   ├── database.go      # MySQL connection via GORM, AutoMigrate for all models
│   └── seed.go          # Seeds species (Cattle/Goat/Chicken), breeds, vaccines
├── middleware/auth.go   # JWT auth middleware, login handler, getUserID/role helpers
├── models/              # GORM model definitions (13 models)
├── handlers/            # HTTP handlers (13 handler files + base utilities)
├── requests/            # Request validation (species create/update with FromContext pattern)
├── routes/routes.go     # Route setup with CORS, JWT-protected API groups
├── frontend/            # Vue 3 SPA (separate package.json)
└── data/                # Postman collection
```

## Database Models (13 tables)

### Core Livestock
1. **User** — id, name, email, phone, username, password, role (Owner/Manager/Veterinarian/Worker/Accountant), status, soft delete
2. **Species** — id, name (unique), created_by → User
3. **Breed** — id, species_id → Species, name
4. **Animal** — id, tag_no (unique), species_id, breed_id, father_id, mother_id, gender, birth_date, purchase_date/price, current_weight, color, status (Healthy/Pregnant/Sick/Sold/Dead), remarks, soft delete. Preloads: Species, Breed, Father, Mother, WeightHistories, Vaccinations
5. **AnimalWeightHistory** — id, animal_id, weight, record_date, remarks. On create: also updates animal.current_weight
6. **AnimalPregnancy** — id, animal_id, breeder_id, mating_date, expected_due_date, actual_birth_date, status (Mated/Pregnant/Delivered/Aborted/Miscarriage/Failed), children counts, note

### Vaccination
7. **Vaccine** — id, species_id, name, description, dose, minimum_age_value/unit, interval_value/unit, is_repeatable. Seeded per species with defaults
8. **AnimalVaccination** — id, animal_id, vaccine_id, vaccination_date, next_due_date, doctor_name, remarks

### Inventory
9. **InventoryCategory** — id, name (unique)
10. **InventoryItem** — id, category_id, name, sku, unit, purchase_price, selling_price
11. **InventoryTransaction** — id, inventory_item_id, transaction_type (Purchase/Sale/Consumption/Adjustment/Return/Damage), quantity, transaction_date, remarks

### Accounting
12. **AccountHead** — id, type (Income/Expense), name, description
13. **AccountTransaction** — id, account_head_id, transaction_date, amount, payment_method (Cash/Bank/Mobile Banking/Other), reference_no, description

## API Routes

All routes under `/api` are JWT-protected except `/api/auth/login`. Paginated endpoints use `?page=1&per_page=20&search=...`.

| Group | Endpoints |
|---|---|
| Health | `GET /health` |
| Auth | `POST /api/auth/login` |
| Users | CRUD + `GET /profile` (self) |
| Species | CRUD (with search, preloads User + Breeds) |
| Breeds | CRUD (filter by species_id) |
| Animals | CRUD + `GET /:id/profile` (full detail + pregnancies) |
| Weight Histories | CRUD (filter by animal_id, auto-updates current_weight) |
| Vaccines | CRUD (filter by species_id) |
| Vaccinations | CRUD (filter by animal_id) |
| Inventory Categories | CRUD |
| Inventory Items | CRUD (filter by category_id) |
| Inventory Transactions | CRUD (filter by item_id, transaction_type) |
| Account Heads | CRUD (filter by type) |
| Account Transactions | CRUD (filter by head_id, type via join) |
| Pregnancies | CRUD (filter by animal_id, status) |

## Common Patterns
- **Handlers** receive `(c fiber.Ctx, db *gorm.DB)` via `wrapH` closure that injects `database.DB`
- **Pagination** via `paginate()` helper — returns `{ data, total, page, per_page, total_pages }`
- **Audit fields**: `created_by` and `updated_by` set from JWT claims via `middleware.GetUserID`
- **Error handling**: `handleError()` maps FK violations → 422, duplicates → 409, others → 400
- **Validation**: `validateBody()` binds + validates with go-playground; species uses custom `FromContext` pattern
- **Seeder**: Seeds default admin (admin/admin123), 3 species, 20 breeds, 21 vaccines

## Auth
- **Roles**: Owner, Manager, Veterinarian, Worker, Accountant
- **Login**: email/phone/username + password → JWT token + user info
- **RBAC**: implicit in code (e.g., only Owner can change roles; Owner/Manager can change status)

## Seeded Data
- Admin user: `admin` / `admin123`
- Species: Cattle, Goat, Chicken with respective breeds and vaccines
- Vaccines include dosage, minimum age, interval schedule, and repeatability info

## Frontend
- Vue 3 SPA with Pinia state management, Vue Router
- Axios for API calls
- Tailwind CSS v4 + Heroicons + Lucide icons
- `vue-tsc` for type checking, Vite for build

## Development
- `air` for hot reload (excludes frontend/, data/, node_modules/)
- `APP_PORT=:8080` default
- JWT secret: `change-me-in-production` (default)
