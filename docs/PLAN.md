Email: admin@laundry.local
Password: defaultpassword123


# Laundry Management System — Implementation Plan

## Context
Building a drop-off / pickup laundry service management system from scratch. Staff receive customer laundry, process it through stages, and customers return for pickup. The system needs order tracking, customer management, and role-based access for Admin and Staff users.

---

## Tech Stack
- **Backend**: Go (Golang) with Gin router, pgx/v5 for PostgreSQL, JWT auth
- **Frontend**: Vue 3 + TypeScript, Pinia (state), Vue Router, Axios
- **Database**: PostgreSQL 16
- **Dev environment**: Docker Compose (PostgreSQL service)

---

## Project Structure

Follows the same layout as Lodge-Management-System (Go backend only, frontend separate):

```
laundry-management-system/
├── cmd/
│   └── api/
│       └── main.go                      # Entry point, wires router, runs migrations
├── internal/
│   ├── config/
│   │   └── config.go                    # Env vars: DATABASE_URL, JWT_SECRET, PORT
│   ├── database/
│   │   └── postgres.go                  # pgxpool connection
│   ├── middleware/
│   │   ├── auth.go                      # JWT validation, RequireAdmin()
│   │   └── cors.go
│   ├── models/
│   │   ├── user.go
│   │   ├── customer.go
│   │   └── order.go
│   ├── handlers/
│   │   ├── auth_handler.go              # POST /auth/login, POST /auth/refresh
│   │   ├── user_handler.go              # Admin-only CRUD staff accounts
│   │   ├── customer_handler.go
│   │   └── order_handler.go             # Status transitions enforced here
│   ├── services/
│   │   ├── auth_service.go              # bcrypt, JWT sign/verify
│   │   ├── user_service.go
│   │   ├── customer_service.go
│   │   └── order_service.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── customer_repository.go
│   │   └── order_repository.go
│   ├── routes/
│   │   └── routes.go                    # Route definitions, middleware wiring
│   └── utils/
│       └── response.go                  # Standard JSON envelope
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_customers.sql
│   └── 003_create_orders.sql
├── pkg/
│   └── utils/
│       └── (utility functions)
├── scripts/
│   └── seed.sh                          # Database seeding
├── docs/
│   └── PLAN.md                          # This file
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── .env.example
└── .gitignore
```

---

## Database Schema

### `roles`
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
name VARCHAR(50) UNIQUE NOT NULL        -- 'admin', 'staff'
description TEXT
created_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
```

### `permissions`
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
resource VARCHAR(50) NOT NULL           -- 'customers', 'orders', 'users', 'reports'
action VARCHAR(50) NOT NULL             -- 'read', 'create', 'update', 'delete', 'update_status'
description TEXT
created_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
UNIQUE(resource, action)
```

Seed with standard permissions:
- `customers: read, create, update, delete`
- `orders: read, create, update, update_status, delete`
- `users: read, create, update, delete` (admin only)
- `reports: read`

### `role_permissions` (Junction Table)
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE
created_at TIMESTAMPTZ DEFAULT NOW()
UNIQUE(role_id, permission_id)
```

**Admin role** gets all permissions.
**Staff role** gets: `customers:read`, `customers:create`, `customers:update`, `orders:read`, `orders:create`, `orders:update`, `reports:read` (no delete, no update_status, no users management)

### `users`
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
full_name VARCHAR(100) NOT NULL
email VARCHAR(150) UNIQUE NOT NULL
password VARCHAR(255) NOT NULL          -- bcrypt hash
role_id UUID NOT NULL REFERENCES roles(id)
is_active BOOLEAN DEFAULT true
created_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
last_login_at TIMESTAMPTZ
```

Seed with one admin user (password hashed):
- `email: 'admin@laundry.local'`, `role_id: <admin role UUID>`, `is_active: true`

### `customers`
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
name VARCHAR(100) NOT NULL
phone VARCHAR(30) UNIQUE NOT NULL
address TEXT
notes TEXT
created_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
```

### `orders`
```sql
CREATE TYPE order_status AS ENUM ('received', 'washing', 'done', 'picked_up');

id UUID PRIMARY KEY DEFAULT gen_random_uuid()
customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT
created_by UUID NOT NULL REFERENCES users(id)
status order_status NOT NULL DEFAULT 'received'
items JSONB NOT NULL                    -- [{name, qty, price}]
total_price NUMERIC(10,2)
notes TEXT
received_at TIMESTAMPTZ DEFAULT NOW()
updated_at TIMESTAMPTZ DEFAULT NOW()
picked_up_at TIMESTAMPTZ               -- set when status → picked_up
```

Indexes: `orders(customer_id)`, `orders(status)`, `orders(created_by)`

---

## Backend API Endpoints

```
POST   /api/auth/login                  # Public — email + password → token
POST   /api/auth/refresh                # Public — refresh token → new access token

# JWTAuth + RequireAnyRole("admin", "staff")
GET    /api/customers
POST   /api/customers
GET    /api/customers/:id
PUT    /api/customers/:id
GET    /api/customers/:id/orders

GET    /api/orders                      # ?status=washing filter
POST   /api/orders
GET    /api/orders/:id
PUT    /api/orders/:id
PATCH  /api/orders/:id/status           # {status: "washing"}

# JWTAuth + RequireRole("admin")
GET    /api/users                       # List all staff users
POST   /api/users                       # Create new staff account
PUT    /api/users/:id                   # Edit staff account
DELETE /api/users/:id                   # Deactivate staff account
GET    /api/reports/summary             # Orders by status + today's revenue
```

### Middleware Chain Example (in `internal/routes/routes.go`)
```
protectedGroup := router.Group("/api").Use(JWTAuth())
{
  customerGroup := protectedGroup.Group("/customers").Use(RequireAnyRole("admin", "staff"))
  {
    customerGroup.GET("", CustomerList)
    customerGroup.POST("", CustomerCreate)
    ...
  }

  adminGroup := protectedGroup.Group("/users").Use(RequireRole("admin"))
  {
    adminGroup.GET("", UserList)
    adminGroup.POST("", UserCreate)
    ...
  }
}
```

### Status Transition Rules (enforced in `internal/services/order_service.go`)
- `received → washing → done → picked_up` only (no skipping, no reversing)
- Invalid transition returns `422 Unprocessable Entity`

### Standard Response Envelope
```json
{ "success": true, "data": { ... }, "error": null }
```

### Key Go Dependencies
- `github.com/gin-gonic/gin` — router
- `github.com/golang-jwt/jwt/v5` — JWT
- `github.com/jackc/pgx/v5` — PostgreSQL driver (pgxpool)
- `golang.org/x/crypto/bcrypt` — password hashing
- `github.com/joho/godotenv` — env loading
- `github.com/google/uuid` — UUID generation
- `github.com/golang-migrate/migrate/v4` — migrations

---

## Auth & RBAC Flow

### Authentication (JWT)
1. `POST /api/auth/login` — bcrypt compare password, issue access token (15 min)
2. Token payload includes: `userID`, `email`, `roleID`, `orgID`, `branchID`
3. All protected requests: `Authorization: Bearer <accessToken>`

### Middleware Chain
- `internal/middleware/auth.go` — `JWTAuth()` validates token, fetches full User from DB (with Role), injects into context
- Sets context keys: `userID`, `email`, `user` (full User object), `role` (Role name), `orgID`, `branchID`
- User account must be `isActive`; inactive users get 401

### RBAC (Permission-Based Access Control)
Permission model: `resource:action` (e.g., `customers:read`, `orders:update_status`)

**Middleware functions** (`internal/middleware/rbac.go`):
- `RequirePermission(resource, action)` — fine-grained check: user must have `resource:action`
- `RequireAnyPermission(permissions...)` — user must have at least one of the given permissions

**Permission lookup flow**:
1. Extract user from context (added by `JWTAuth()`)
2. Fetch user's role from DB
3. Query role's permissions from `role_permissions` junction table
4. Check if `resource:action` exists in permission set
5. Return 403 if not found

### Standard Permissions
**Customers**:
- `customers:read` — view all customers
- `customers:create` — create new customer
- `customers:update` — edit customer details
- `customers:delete` — delete customer

**Orders**:
- `orders:read` — view orders
- `orders:create` — create order
- `orders:update` — edit order details
- `orders:update_status` — change order status (received → washing → done → picked_up)
- `orders:delete` — delete order

**Users** (staff management):
- `users:read` — view staff accounts
- `users:create` — create staff account
- `users:update` — edit staff account
- `users:delete` — deactivate staff account

**Reports**:
- `reports:read` — view dashboard and reports

### Laundry System Roles & Permissions
- **admin**: has all permissions (auto-assigned on role creation)
- **staff**: has `customers:read`, `customers:create`, `customers:update`, `orders:read`, `orders:create`, `orders:update`, `reports:read`
  - Does NOT have: `customers:delete`, `orders:delete`, `orders:update_status`, `users:*`

### Token Refresh
- On 401, Axios interceptor calls `POST /api/auth/refresh` with refresh token
- Backend validates refresh token, issues new access token, retries original request
- On second 401, logout + redirect to `/login`

---

## Frontend Auth & Routing

### Pinia `auth.ts` state
- `token`, `refreshToken`, `user: User | null` — persisted to `localStorage`
- Computed: `isAdmin`, `isAuthenticated`

### Route Guard (`router/index.ts` — `beforeEach`)
1. No token + auth-required route → redirect `/login`
2. Role=staff + admin-required route → redirect `/orders`
3. Already authenticated + `/login` → redirect `/orders`

---

## Implementation Order

### Phase 1 — Infrastructure
1. `.gitignore`, `docker-compose.yml` (PostgreSQL 16), `Makefile`
2. Go module init (`go mod init laundry-system`)
3. `internal/config/config.go` — load env vars (DATABASE_URL, JWT_SECRET, PORT)
4. `internal/database/postgres.go` — pgxpool connection with migration runner
5. Migration SQL files in `migrations/` (001, 002, 003)
6. `cmd/api/main.go` — HTTP server startup, run migrations on boot

### Phase 2 — Core Models & Database
7. `internal/models/role.go` — Role struct, predefined roles (admin, staff)
8. `internal/models/user.go` — User struct with Role embed, `HasPermission(permission string) bool` method
9. `internal/models/customer.go`, `order.go`
10. Verify migrations execute successfully, seed admin role and admin user

### Phase 3 — Auth Backend
9. `internal/services/auth_service.go` — bcrypt hash/compare, JWT sign/verify
10. `internal/repository/user_repository.go` — `GetUserByID`, `GetUserByEmail`, `Create`, `Update`
11. `internal/repository/role_repository.go` — `GetByName`, `List`
12. `internal/handlers/auth_handler.go` — `POST /api/auth/login` (email+password → token)
13. `internal/middleware/auth.go` — `JWTAuth()` validates token, fetches full User from DB, injects to context
14. `internal/middleware/rbac.go` — `RequireRole(role)`, `RequireAnyRole(roles...)`, `RequirePermission(permission)` middleware
15. Implement `User.HasPermission(permission)` in models/user.go (role-action matrix)
16. `scripts/seed.sh` — seed admin role + admin user (default email: `admin@laundry.local`, password: prompt or env var)
17. Test login with curl/Postman, verify JWT is valid

### Phase 4 — Customer & Order Backend
15. `internal/repository/customer_repository.go` — CRUD methods
16. `internal/services/customer_service.go`
17. `internal/handlers/customer_handler.go` — all customer endpoints
18. `internal/repository/order_repository.go` — CRUD + status transition methods
19. `internal/services/order_service.go` — business logic (status validation)
20. `internal/handlers/order_handler.go` — all order endpoints + `GET /api/reports/summary`
21. `internal/routes/routes.go` — wire all handlers, apply middleware, define routes
22. Test all endpoints with Postman/curl

### Phase 5 — User (Staff) Management Backend
23. `internal/repository/user_repository.go` — extend with admin CRUD operations
24. `internal/handlers/user_handler.go` — admin-only endpoints (Create, Update, Delete staff)
25. Test user endpoints with admin JWT

### Phase 6 — Frontend Scaffold
26. Initialize Vue 3 project in separate `laundry-management-frontend` directory (or use existing frontend pattern)
27. `src/types/index.ts`, `src/api/client.ts` (Axios + interceptors with token refresh)
28. `src/stores/auth.ts` (Pinia), `src/router/index.ts` with route guard

### Phase 7 — Frontend Auth & Layout
29. `LoginView.vue`, `AppSidebar.vue`, `AppHeader.vue`
30. Verify login flow → JWT token → authenticated state persists

### Phase 8 — Frontend Core Views
31. `CustomersView.vue` + `CustomerForm.vue` (modal or page)
32. `CustomerDetailView.vue` (profile + order history)
33. `OrdersView.vue` (status-tab filter, list, "New Order" button)
34. `OrderForm.vue` (customer picker dropdown, add items dynamically, auto-total)
35. `OrderDetailView.vue` (full details, progressive status-change buttons)
36. `StatusBadge.vue` (reusable status display component)

### Phase 9 — Frontend Admin Views
37. `DashboardView.vue` (summary cards: orders by status, today's revenue)
38. `StaffManagementView.vue` (list, create/edit/delete staff accounts)
39. Verify route guards prevent staff from accessing admin views

### Phase 10 — Polish & Hardening
40. Backend: form validation, return 400 with field-level errors
41. Frontend: loading states, error toasts, empty-state messages
42. Confirmation dialogs for destructive actions (delete customer/order/staff)
43. Complete token refresh flow test (let access token expire, verify retry)

---

## Verification

### Backend (curl/Postman)
| # | Test | Expected |
|---|------|----------|
| 1 | Login valid admin | 200 + tokens |
| 2 | Login wrong password | 401 |
| 3 | `GET /customers` no token | 401 |
| 4 | `GET /users` with staff JWT | 403 |
| 5 | Create customer (staff JWT) | 201 |
| 6 | Create order | 201, status=received |
| 7 | PATCH status: received→washing | 200 |
| 8 | PATCH status: washing→picked_up (skip done) | 422 |
| 9 | PATCH status: done→picked_up | 200, picked_up_at set |
| 10 | Expired access + valid refresh | silent retry succeeds |

### Frontend Walkthrough
1. Open app → redirected to `/login`
2. Login as admin → Dashboard + Staff links visible in sidebar
3. Create customer → appears in list
4. Create order → appears with "Received" badge
5. Advance order through all 4 statuses
6. Check Customer Detail → order in history
7. Login as staff → Dashboard/Staff links hidden; direct nav to `/dashboard` redirects to `/orders`

### Critical Files (get these right first)
**Backend:**
- `internal/models/user.go` → `User.HasPermission(permission)` — role-action matrix is the source of truth for authorization
- `internal/middleware/auth.go` → `JWTAuth()` — fetches full User from DB; if buggy, role is wrong
- `internal/middleware/rbac.go` → `RequireRole()`, `RequireAnyRole()` — guards endpoints; misconfigured middleware = auth bypass
- `internal/routes/routes.go` — must apply correct middleware chain to each route group; a missing middleware = endpoint unprotected
- `internal/services/order_service.go` — status transition validation (core business logic)
- `migrations/001_create_roles.sql` and migration runner in `internal/database/postgres.go` — roles table must exist before users can be seeded

**Frontend:**
- `src/router/index.ts` — route-level auth guard prevents unauthorized nav
- `src/api/client.ts` — token refresh interceptor; bugs break all authenticated API calls
