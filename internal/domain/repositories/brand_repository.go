package repositories

import (
	"github.com/google/uuid"
	"car-service/internal/domain/entities"
)

type BrandRepository interface {
	Create(brand *entities.Brand) error
	GetByID(id uuid.UUID) (*entities.Brand, error)
	GetByName(name string) (*entities.Brand, error)
	Update(brand *entities.Brand) error
	Delete(id uuid.UUID) error
	List() ([]*entities.Brand, error)
	ListActive() ([]*entities.Brand, error)
} 