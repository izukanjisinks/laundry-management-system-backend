package repository

import (
	"database/sql"
	"fmt"
	"time"

	"laundry-system/internal/database"
	"laundry-system/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, full_name, email, password, role_id, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`
	var lastLogin sql.NullTime
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Password,
		&user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}
	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, full_name, email, password, role_id, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	var lastLogin sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Password,
		&user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}
	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		user.FullName, user.Email, user.Password, user.RoleID, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET full_name = $1, email = $2, role_id = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	return r.db.QueryRow(query,
		user.FullName, user.Email, user.RoleID, user.IsActive, user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) UpdatePassword(id, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, id)
	return err
}

func (r *UserRepository) UpdateLastLogin(id string) error {
	now := time.Now()
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, now, id)
	return err
}

func (r *UserRepository) List() ([]models.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.role_id, u.is_active, u.created_at, u.updated_at, u.last_login_at,
		       r.id, r.name, r.description
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY u.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var role models.Role
		var lastLogin sql.NullTime
		if err := rows.Scan(
			&u.ID, &u.FullName, &u.Email, &u.RoleID, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt, &lastLogin,
			&role.ID, &role.Name, &role.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if lastLogin.Valid {
			u.LastLoginAt = &lastLogin.Time
		}
		u.Role = &role
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Delete(id string) error {
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
