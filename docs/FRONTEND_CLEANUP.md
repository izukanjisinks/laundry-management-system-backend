# WashPoint Frontend Cleanup Plan

**Goal**: Bring `laungry-management-system-frontend` into exact alignment with the WashPoint.html design.  
**Frontend path**: `C:\development\laungry-management-system-frontend`  
**Design reference**: `docs/WashPoint.html`

---

## Phase 1 — Types & API Foundation
> No UI changes yet. Fix the data layer so everything downstream compiles correctly.

### 1.1 `src/types/index.ts`
- Rename `OrderStatus` `'done'` → `'ready'`
- Add `ServiceType`: `'wash_fold' | 'dry_clean' | 'ironing' | 'wash_iron'`
- Add `PaymentStatus`: `'unpaid' | 'partial' | 'paid'`
- Add `PaymentMethod`: `'cash' | 'card' | 'transfer'`
- Add to `Order`: `order_number`, `service_type`, `subtotal`, `tax_rate`, `tax_amount`, `payment_status`, `payment_method`, `due_at`
- Add to `Customer`: `email`, `total_orders`
- Fix `ReportSummary`: rename `todays_revenue` → `today_revenue`, `total_orders_today` → `total_orders`; add `unpaid_orders`, `daily_orders`
- Add `CatalogItem` interface: `id`, `name`, `slug`, `base_price`, `is_active`, `sort_order`

### 1.2 `src/api/orders.ts`
- Fix `create()` payload to include `service_type`, `due_at` (drop `total_price` — backend calculates it)
- Add `updatePayment(id, paymentStatus, paymentMethod)`
- Add `delete(id)`

### 1.3 `src/api/customers.ts`
- Add `search` param to `list(search?)`
- Add `delete(id)`

### 1.4 `src/api/catalog.ts` *(new file)*
- `list()` → `GET /api/catalog`

### 1.5 `src/api/reports.ts`
- No signature change needed; field names fixed in types

**Exit criteria**: `npm run build` compiles with zero type errors.

---

## Phase 2 — Login Screen
> Match the WashPoint splash: orange brand, droplet logo, "WashPoint" name.

### `src/views/auth/LoginView.vue`
- Replace the "L" letter avatar with the WashPoint droplet SVG (white water-drop on orange circle)
- Change heading from "Laundry MS" → "WashPoint"
- Change subheading to "Staff & Admin Portal"
- Keep the glass-card layout, HexaPrime footer, and existing form logic unchanged

**Exit criteria**: Login page shows WashPoint branding with orange droplet.

---

## Phase 3 — Sidebar & Navigation
> Match the WashPoint nav structure exactly.

### `src/components/layout/AppSidebar.vue`
- Replace `SidebarLogo` text/icon with WashPoint droplet SVG + "WashPoint" wordmark
- Restructure `navMain` to flat primary links (no sub-items needed for MVP):
  - Dashboard (admin only) → `/dashboard`
  - Orders → `/orders`
  - New Order → `/orders/new`
  - Customers → `/customers`
  - Staff (admin only) → `/admin/staff`
- Remove `NavDocuments` / `NavSecondary` (Settings link not in WashPoint design)
- Keep `NavUser` footer (user name + logout)
- Active route highlight uses orange `#F26F21` accent

**Exit criteria**: Sidebar matches WashPoint nav exactly; admin-only items hidden for staff.

---

## Phase 4 — Dashboard
> 4 KPI cards + 7-day bar chart matching the WashPoint dashboard screen.

### `src/views/dashboard/DashboardView.vue`
- Fix field references: `today_revenue` (was `todays_revenue`), `total_orders` (was `total_orders_today`)
- Build 4 stat cards:
  1. Today's Revenue — `K{today_revenue}`
  2. Total Orders — `total_orders`
  3. Unpaid Orders — `unpaid_orders`
  4. Active (received + washing) — computed from `orders_by_status`
- Replace/add 7-day bar chart using `daily_orders` array
  - X axis: day labels (Mon, Tue … or date)
  - Y axis: order count
  - Bar color: orange `#F26F21`
- Remove `SectionCards` generic component; build dashboard-specific layout inline

**Exit criteria**: Dashboard shows 4 cards and animated 7-day chart populated from live API.

---

## Phase 5 — Orders List (Processing Board)
> Tab-filter UI + WP order numbers + service type column.

### `src/views/orders/OrdersView.vue`
- Add status tab bar: **All · Received · Washing · Ready · Pickup** — clicking sets `statusFilter` and refetches
- Replace UUID display with `WP-{order_number}` (zero-pad to 4 digits: `WP-1042`)
- Fix `statusColors`: rename `done` → `ready`; colors: received=blue, washing=amber, ready=green, picked_up=gray
- Add `service_type` column (formatted: "Wash & Fold", "Dry Clean", etc.)
- Make the full row clickable → `/orders/{id}` (not just the ID cell)
- Add payment status badge column

**Exit criteria**: Orders list shows WP numbers, status tabs filter correctly, service type visible.

---

## Phase 6 — New Order Form
> Catalog picker, service type, due date, tax breakdown.

### `src/views/orders/OrderForm.vue`
- On mount, fetch `GET /api/catalog` and `GET /api/customers`
- Replace free-text item input with catalog item rows:
  - Dropdown to pick catalog item (name + base price pre-fills)
  - Quantity stepper (+ / − buttons)
  - Price auto-filled from catalog, editable override allowed
  - Add row / remove row buttons
- Add `service_type` select: Wash & Fold / Dry Clean / Ironing / Wash & Iron
- Add `due_at` date-time picker (optional)
- Replace single "Total" with 3-line breakdown:
  - Subtotal: `K{subtotal}`
  - Tax (7.5%): `K{tax}`
  - **Total: `K{total}`** (bold)
- Remove `total_price` from create payload; send `service_type`, `due_at` instead
- After successful create, redirect to `/orders/{newId}`

**Exit criteria**: New order form uses catalog, shows tax breakdown, posts correct payload.

---

## Phase 7 — Order Detail
> WP order number, full status flow, payment actions, item verification.

### `src/views/orders/OrderDetailView.vue`
- Fix status flow array: `['received', 'washing', 'ready', 'picked_up']` (remove `'done'`)
- Fix `statusColors` map (`ready` key, not `done`)
- Show `WP-{order_number}` as the order title instead of UUID slice
- Display `service_type` (formatted label)
- Replace simple total with subtotal / tax / total breakdown
- Add payment section:
  - Shows current `payment_status` badge
  - If not `paid`: show "Mark as Paid" button → calls `updatePayment(id, 'paid', method)`
  - Payment method select (Cash / Card / Transfer) shown alongside
- Advance status button visible to **all authenticated users** (not admin-only)
- At `picked_up` status: show item-by-item verification checklist (checkboxes per item, read-only display)
- Show `due_at` if set
- Show `received_at` and `picked_up_at` timestamps

**Exit criteria**: Full order lifecycle manageable from detail view; payment recordable.

---

## Phase 8 — Customers
> Search, order count, email, formatted dates.

### `src/views/customers/CustomersView.vue`
- Add search input that calls `customersApi.list(searchTerm)` with debounce (300ms)
- Add `total_orders` column to the table
- Add `email` column (show "—" if empty)

### `src/views/customers/CustomerDetailView.vue`
- Fix `statusColors` (`ready` not `done`)
- Show `total_orders` count badge on the profile card header
- Add `email` row to profile details
- Show "Member since {Month Year}" formatted from `created_at`
- Order history: show `WP-{order_number}` instead of UUID slice
- Add payment status badge column to order history table

### `src/views/customers/CustomerForm.vue`
- Add `email` field (optional)
- Ensure existing form fields match backend: `name`, `phone`, `address`, `notes`, `email`

**Exit criteria**: Customer list is searchable; detail shows full profile with order history.

---

## Phase 9 — Encoding Cleanup
> Fix mojibake in all files (UTF-8 double-encoded characters).

Affected strings throughout all views:
- `â€¦` → `…`
- `â€"` → `—`
- `Loadingâ€¦` → `Loading…`

Files to sweep: `OrdersView.vue`, `OrderForm.vue`, `CustomersView.vue`, `CustomerDetailView.vue`, `CustomerForm.vue`, `StaffManagementView.vue`

**Exit criteria**: No corrupted characters anywhere in the UI.

---

## Phase 10 — Staff Management (Admin)
> Verify existing view matches backend user model.

### `src/views/admin/StaffManagementView.vue`
- Confirm uses `full_name` (not `name`) — matches backend `User` struct
- Add `is_active` toggle / display
- Show `last_login_at` column
- Add delete (deactivate) action

**Exit criteria**: Admin can create, view, and deactivate staff accounts.

---

## Completion Checklist

| Phase | Scope | Status |
|-------|-------|--------|
| 1 | Types & API | ☐ |
| 2 | Login screen | ☐ |
| 3 | Sidebar & nav | ☐ |
| 4 | Dashboard | ☐ |
| 5 | Orders list | ☐ |
| 6 | New order form | ☐ |
| 7 | Order detail | ☐ |
| 8 | Customers | ☐ |
| 9 | Encoding cleanup | ☐ |
| 10 | Staff management | ☐ |

---

## Key Constants (reference while coding)

```
Brand orange:     #F26F21
Order number fmt: WP-{order_number}   (e.g. WP-1042)
Tax rate:         7.5%
Status flow:      received → washing → ready → picked_up
Status colors:    received=blue, washing=amber, ready=green, picked_up=gray
Service types:    wash_fold="Wash & Fold", dry_clean="Dry Clean",
                  ironing="Ironing", wash_iron="Wash & Iron"
Payment methods:  cash="Cash", card="Card", transfer="Transfer"
API base:         http://localhost:8083/api
```
