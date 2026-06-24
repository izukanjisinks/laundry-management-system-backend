package repository

import (
	"database/sql"
	"fmt"

	"laundry-system/internal/database"
	"laundry-system/internal/models"
)

type CatalogRepository struct {
	db *sql.DB
}

func NewCatalogRepository() *CatalogRepository {
	return &CatalogRepository{db: database.GetDB()}
}

func (r *CatalogRepository) List() ([]models.CatalogItem, error) {
	query := `
		SELECT id, name, slug, base_price, is_active, sort_order, created_at, updated_at
		FROM catalog_items
		WHERE is_active = true
		ORDER BY sort_order ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalog items: %w", err)
	}
	defer rows.Close()

	var items []models.CatalogItem
	for rows.Next() {
		var item models.CatalogItem
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Slug, &item.BasePrice,
			&item.IsActive, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan catalog item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
