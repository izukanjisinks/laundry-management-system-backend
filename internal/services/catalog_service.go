package services

import (
	"laundry-system/internal/models"
	"laundry-system/internal/repository"
)

type CatalogService struct {
	catalogRepo *repository.CatalogRepository
}

func NewCatalogService(catalogRepo *repository.CatalogRepository) *CatalogService {
	return &CatalogService{catalogRepo: catalogRepo}
}

func (s *CatalogService) List() ([]models.CatalogItem, error) {
	return s.catalogRepo.List()
}
