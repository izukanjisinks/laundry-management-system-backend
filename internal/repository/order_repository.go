package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"laundry-system/internal/database"
	"laundry-system/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: database.GetDB()}
}

func (r *OrderRepository) Create(o *models.Order) error {
	items, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}
	query := `
		INSERT INTO orders (customer_id, created_by, status, items, total_price, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, received_at, updated_at
	`
	return r.db.QueryRow(query,
		o.CustomerID,
		o.CreatedBy,
		o.Status,
		items,
		o.TotalPrice,
		database.NullableString(o.Notes),
	).Scan(&o.ID, &o.ReceivedAt, &o.UpdatedAt)
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	query := `
		SELECT
			o.id, o.customer_id, o.created_by, o.status, o.items,
			o.total_price, o.notes, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone, c.email,
			u.id, u.full_name
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = o.created_by
		WHERE o.id = $1
	`
	o := &models.Order{}
	c := &models.Customer{}
	creator := &models.User{}
	var notes sql.NullString
	var pickedUpAt sql.NullTime
	var customerEmail sql.NullString
	var itemsRaw []byte

	err := r.db.QueryRow(query, id).Scan(
		&o.ID, &o.CustomerID, &o.CreatedBy, &o.Status, &itemsRaw,
		&o.TotalPrice, &notes, &o.ReceivedAt, &o.UpdatedAt, &pickedUpAt,
		&c.ID, &c.Name, &c.Phone, &customerEmail,
		&creator.ID, &creator.FullName,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := json.Unmarshal(itemsRaw, &o.Items); err != nil {
		return nil, fmt.Errorf("failed to parse order items: %w", err)
	}

	o.Notes = database.StringOrEmpty(notes)
	c.Email = database.StringOrEmpty(customerEmail)
	if pickedUpAt.Valid {
		o.PickedUpAt = &pickedUpAt.Time
	}
	o.Customer = c
	o.Creator = creator
	return o, nil
}

func (r *OrderRepository) List(status string) ([]models.Order, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT
			o.id, o.customer_id, o.created_by, o.status, o.items,
			o.total_price, o.notes, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone,
			u.id, u.full_name
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = o.created_by
	`
	if status != "" {
		rows, err = r.db.Query(baseQuery+` WHERE o.status = $1 ORDER BY o.received_at DESC`, status)
	} else {
		rows, err = r.db.Query(baseQuery + ` ORDER BY o.received_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *OrderRepository) ListByCustomer(customerID string) ([]models.Order, error) {
	query := `
		SELECT
			o.id, o.customer_id, o.created_by, o.status, o.items,
			o.total_price, o.notes, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone,
			u.id, u.full_name
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = o.created_by
		WHERE o.customer_id = $1
		ORDER BY o.received_at DESC
	`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list customer orders: %w", err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *OrderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	var query string
	if status == models.OrderStatusPickedUp {
		query = `UPDATE orders SET status = $1, picked_up_at = NOW(), updated_at = NOW() WHERE id = $2`
	} else {
		query = `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	}
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *OrderRepository) Update(o *models.Order) error {
	items, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}
	query := `
		UPDATE orders
		SET items = $1, total_price = $2, notes = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err = r.db.QueryRow(query,
		items,
		o.TotalPrice,
		database.NullableString(o.Notes),
		o.ID,
	).Scan(&o.UpdatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("order not found")
	}
	return err
}

func (r *OrderRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *OrderRepository) Summary() (map[string]interface{}, error) {
	statusQuery := `
		SELECT status, COUNT(*) FROM orders
		GROUP BY status
	`
	rows, err := r.db.Query(statusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get order counts: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{
		"received":  0,
		"washing":   0,
		"done":      0,
		"picked_up": 0,
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}

	var todayRevenue sql.NullFloat64
	revenueQuery := `
		SELECT COALESCE(SUM(total_price), 0)
		FROM orders
		WHERE status = 'picked_up' AND picked_up_at::date = CURRENT_DATE
	`
	if err := r.db.QueryRow(revenueQuery).Scan(&todayRevenue); err != nil {
		return nil, fmt.Errorf("failed to get today's revenue: %w", err)
	}

	var totalOrders int
	r.db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&totalOrders)

	return map[string]interface{}{
		"orders_by_status": counts,
		"today_revenue":    todayRevenue.Float64,
		"total_orders":     totalOrders,
	}, nil
}

// scanOrders is a shared row scanner for order list queries.
func scanOrders(rows *sql.Rows) ([]models.Order, error) {
	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var c models.Customer
		var creator models.User
		var notes sql.NullString
		var pickedUpAt sql.NullTime
		var itemsRaw []byte

		if err := rows.Scan(
			&o.ID, &o.CustomerID, &o.CreatedBy, &o.Status, &itemsRaw,
			&o.TotalPrice, &notes, &o.ReceivedAt, &o.UpdatedAt, &pickedUpAt,
			&c.ID, &c.Name, &c.Phone,
			&creator.ID, &creator.FullName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal(itemsRaw, &o.Items); err != nil {
			return nil, fmt.Errorf("failed to parse order items: %w", err)
		}

		o.Notes = database.StringOrEmpty(notes)
		if pickedUpAt.Valid {
			o.PickedUpAt = &pickedUpAt.Time
		}
		o.Customer = &c
		o.Creator = &creator
		orders = append(orders, o)
	}
	return orders, nil
}
