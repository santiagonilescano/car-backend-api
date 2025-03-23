package repositories

import (
	"github.com/google/uuid"
	"car-service/internal/domain/entities"
)

type ModelRepository interface {
	Create(model *entities.Model) error
	GetByID(id uuid.UUID) (*entities.Model, error)
	GetByBrandID(brandID uuid.UUID) ([]*entities.Model, error)
	GetByNameAndBrand(name string, brandID uuid.UUID) (*entities.Model, error)
	Update(model *entities.Model) error
	Delete(id uuid.UUID) error
	List() ([]*entities.Model, error)
	ListActive() ([]*entities.Model, error)
	ListByCategory(category string) ([]*entities.Model, error)
} 