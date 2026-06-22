package repository

import (
	"database/sql"
	"fmt"

	"laundry-system/internal/database"
	"laundry-system/internal/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{db: database.GetDB()}
}

func (r *CustomerRepository) Create(c *models.Customer) error {
	query := `
		INSERT INTO customers (name, phone, email, address, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		c.Name,
		c.Phone,
		database.NullableString(c.Email),
		database.NullableString(c.Address),
		database.NullableString(c.Notes),
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CustomerRepository) GetByID(id string) (*models.Customer, error) {
	c := &models.Customer{}
	var email, address, notes sql.NullString
	query := `
		SELECT id, name, phone, email, address, notes, created_at, updated_at
		FROM customers WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Phone, &email, &address, &notes,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	c.Email = database.StringOrEmpty(email)
	c.Address = database.StringOrEmpty(address)
	c.Notes = database.StringOrEmpty(notes)
	return c, nil
}

func (r *CustomerRepository) GetByPhone(phone string) (*models.Customer, error) {
	c := &models.Customer{}
	var email, address, notes sql.NullString
	query := `
		SELECT id, name, phone, email, address, notes, created_at, updated_at
		FROM customers WHERE phone = $1
	`
	err := r.db.QueryRow(query, phone).Scan(
		&c.ID, &c.Name, &c.Phone, &email, &address, &notes,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	c.Email = database.StringOrEmpty(email)
	c.Address = database.StringOrEmpty(address)
	c.Notes = database.StringOrEmpty(notes)
	return c, nil
}

func (r *CustomerRepository) List(search string) ([]models.Customer, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		query := `
			SELECT id, name, phone, email, address, notes, created_at, updated_at
			FROM customers
			WHERE name ILIKE $1 OR phone ILIKE $1
			ORDER BY name ASC
		`
		rows, err = r.db.Query(query, "%"+search+"%")
	} else {
		query := `
			SELECT id, name, phone, email, address, notes, created_at, updated_at
			FROM customers
			ORDER BY name ASC
		`
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		var email, address, notes sql.NullString
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &email, &address, &notes,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		c.Email = database.StringOrEmpty(email)
		c.Address = database.StringOrEmpty(address)
		c.Notes = database.StringOrEmpty(notes)
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *CustomerRepository) Update(c *models.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, phone = $2, email = $3, address = $4, notes = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRow(query,
		c.Name,
		c.Phone,
		database.NullableString(c.Email),
		database.NullableString(c.Address),
		database.NullableString(c.Notes),
		c.ID,
	).Scan(&c.UpdatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("customer not found")
	}
	return err
}

func (r *CustomerRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}
