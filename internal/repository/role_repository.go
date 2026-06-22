package repository

import (
	"database/sql"
	"fmt"

	"laundry-system/internal/database"
	"laundry-system/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.GetDB()}
}

func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	role := &models.Role{}
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) GetByID(id string) (*models.Role, error) {
	role := &models.Role{}
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) GetPermissionsByRoleID(roleID string) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.resource, p.action, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

func (r *RoleRepository) List() ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}
