package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
		INSERT INTO orders
			(customer_id, created_by, status, service_type, items,
			 subtotal, tax_rate, tax_amount, total_price,
			 payment_status, payment_method, notes, due_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, order_number, received_at, updated_at
	`
	return r.db.QueryRow(query,
		o.CustomerID,
		o.CreatedBy,
		o.Status,
		o.ServiceType,
		items,
		o.Subtotal,
		o.TaxRate,
		o.TaxAmount,
		o.TotalPrice,
		o.PaymentStatus,
		database.NullableString(string(o.PaymentMethod)),
		database.NullableString(o.Notes),
		nullableTime(o.DueAt),
	).Scan(&o.ID, &o.OrderNumber, &o.ReceivedAt, &o.UpdatedAt)
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	query := `
		SELECT
			o.id, o.order_number, o.customer_id, o.created_by, o.status, o.service_type,
			o.items, o.subtotal, o.tax_rate, o.tax_amount, o.total_price,
			o.payment_status, o.payment_method,
			o.notes, o.due_at, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone, c.email,
			u.id, u.full_name
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = o.created_by
		WHERE o.id = $1
	`
	o, err := scanOrderRow(r.db.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	return o, err
}

func (r *OrderRepository) List(status string) ([]models.Order, error) {
	base := `
		SELECT
			o.id, o.order_number, o.customer_id, o.created_by, o.status, o.service_type,
			o.items, o.subtotal, o.tax_rate, o.tax_amount, o.total_price,
			o.payment_status, o.payment_method,
			o.notes, o.due_at, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone, c.email,
			u.id, u.full_name
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = o.created_by
	`
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = r.db.Query(base+` WHERE o.status = $1 ORDER BY o.received_at DESC`, status)
	} else {
		rows, err = r.db.Query(base + ` ORDER BY o.received_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()
	return scanOrderRows(rows)
}

func (r *OrderRepository) ListByCustomer(customerID string) ([]models.Order, error) {
	query := `
		SELECT
			o.id, o.order_number, o.customer_id, o.created_by, o.status, o.service_type,
			o.items, o.subtotal, o.tax_rate, o.tax_amount, o.total_price,
			o.payment_status, o.payment_method,
			o.notes, o.due_at, o.received_at, o.updated_at, o.picked_up_at,
			c.id, c.name, c.phone, c.email,
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
	return scanOrderRows(rows)
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

func (r *OrderRepository) UpdatePayment(id string, paymentStatus models.PaymentStatus, paymentMethod models.PaymentMethod) error {
	query := `
		UPDATE orders
		SET payment_status = $1, payment_method = $2, updated_at = NOW()
		WHERE id = $3
	`
	result, err := r.db.Exec(query,
		paymentStatus,
		database.NullableString(string(paymentMethod)),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
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
		SET service_type = $1, items = $2,
		    subtotal = $3, tax_rate = $4, tax_amount = $5, total_price = $6,
		    notes = $7, due_at = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at
	`
	err = r.db.QueryRow(query,
		o.ServiceType,
		items,
		o.Subtotal,
		o.TaxRate,
		o.TaxAmount,
		o.TotalPrice,
		database.NullableString(o.Notes),
		nullableTime(o.DueAt),
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
	rows, err := r.db.Query(`SELECT status, COUNT(*) FROM orders GROUP BY status`)
	if err != nil {
		return nil, fmt.Errorf("failed to get order counts: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{
		"received":  0,
		"washing":   0,
		"ready":     0,
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
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(total_price), 0)
		FROM orders
		WHERE status = 'picked_up' AND picked_up_at::date = CURRENT_DATE
	`).Scan(&todayRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's revenue: %w", err)
	}

	var totalOrders int
	r.db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&totalOrders)

	var unpaidCount int
	r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE payment_status = 'unpaid' AND status != 'picked_up'`).Scan(&unpaidCount)

	// 7-day daily order counts for dashboard chart
	chartRows, err := r.db.Query(`
		SELECT
			gs.day::date AS day,
			COALESCE(COUNT(o.id), 0) AS count
		FROM generate_series(
			CURRENT_DATE - INTERVAL '6 days',
			CURRENT_DATE,
			'1 day'::interval
		) AS gs(day)
		LEFT JOIN orders o ON o.received_at::date = gs.day::date
		GROUP BY gs.day
		ORDER BY gs.day ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily chart data: %w", err)
	}
	defer chartRows.Close()

	type DayCount struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	}
	var dailyCounts []DayCount
	for chartRows.Next() {
		var dc DayCount
		var day time.Time
		if err := chartRows.Scan(&day, &dc.Count); err != nil {
			return nil, err
		}
		dc.Day = day.Format("2006-01-02")
		dailyCounts = append(dailyCounts, dc)
	}

	return map[string]interface{}{
		"orders_by_status": counts,
		"today_revenue":    todayRevenue.Float64,
		"total_orders":     totalOrders,
		"unpaid_orders":    unpaidCount,
		"daily_orders":     dailyCounts,
	}, nil
}

func scanOrderRow(row *sql.Row) (*models.Order, error) {
	o := &models.Order{}
	c := &models.Customer{}
	creator := &models.User{}
	var notes, paymentMethod, customerEmail sql.NullString
	var pickedUpAt, dueAt sql.NullTime
	var itemsRaw []byte
	var subtotal, taxRate, taxAmount sql.NullFloat64

	err := row.Scan(
		&o.ID, &o.OrderNumber, &o.CustomerID, &o.CreatedBy, &o.Status, &o.ServiceType,
		&itemsRaw, &subtotal, &taxRate, &taxAmount, &o.TotalPrice,
		&o.PaymentStatus, &paymentMethod,
		&notes, &dueAt, &o.ReceivedAt, &o.UpdatedAt, &pickedUpAt,
		&c.ID, &c.Name, &c.Phone, &customerEmail,
		&creator.ID, &creator.FullName,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(itemsRaw, &o.Items); err != nil {
		return nil, fmt.Errorf("failed to parse order items: %w", err)
	}
	o.Subtotal = subtotal.Float64
	o.TaxRate = taxRate.Float64
	o.TaxAmount = taxAmount.Float64
	o.Notes = database.StringOrEmpty(notes)
	o.PaymentMethod = models.PaymentMethod(database.StringOrEmpty(paymentMethod))
	c.Email = database.StringOrEmpty(customerEmail)
	if pickedUpAt.Valid {
		o.PickedUpAt = &pickedUpAt.Time
	}
	if dueAt.Valid {
		o.DueAt = &dueAt.Time
	}
	o.Customer = c
	o.Creator = creator
	return o, nil
}

func scanOrderRows(rows *sql.Rows) ([]models.Order, error) {
	var orders []models.Order
	for rows.Next() {
		o := &models.Order{}
		c := &models.Customer{}
		creator := &models.User{}
		var notes, paymentMethod, customerEmail sql.NullString
		var pickedUpAt, dueAt sql.NullTime
		var itemsRaw []byte
		var subtotal, taxRate, taxAmount sql.NullFloat64

		if err := rows.Scan(
			&o.ID, &o.OrderNumber, &o.CustomerID, &o.CreatedBy, &o.Status, &o.ServiceType,
			&itemsRaw, &subtotal, &taxRate, &taxAmount, &o.TotalPrice,
			&o.PaymentStatus, &paymentMethod,
			&notes, &dueAt, &o.ReceivedAt, &o.UpdatedAt, &pickedUpAt,
			&c.ID, &c.Name, &c.Phone, &customerEmail,
			&creator.ID, &creator.FullName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		if err := json.Unmarshal(itemsRaw, &o.Items); err != nil {
			return nil, fmt.Errorf("failed to parse order items: %w", err)
		}
		o.Subtotal = subtotal.Float64
		o.TaxRate = taxRate.Float64
		o.TaxAmount = taxAmount.Float64
		o.Notes = database.StringOrEmpty(notes)
		o.PaymentMethod = models.PaymentMethod(database.StringOrEmpty(paymentMethod))
		c.Email = database.StringOrEmpty(customerEmail)
		if pickedUpAt.Valid {
			o.PickedUpAt = &pickedUpAt.Time
		}
		if dueAt.Valid {
			o.DueAt = &dueAt.Time
		}
		o.Customer = c
		o.Creator = creator
		orders = append(orders, *o)
	}
	return orders, nil
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
