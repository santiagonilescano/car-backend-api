package repositories

import (
	"github.com/google/uuid"
	"car-service/internal/domain/entities"
)

type OwnerRepository interface {
	Create(owner *entities.Owner) error
	GetByID(id uuid.UUID) (*entities.Owner, error)
	GetByEmail(email string) (*entities.Owner, error)
	Update(owner *entities.Owner) error
	Delete(id uuid.UUID) error
	List() ([]*entities.Owner, error)
} 