package repositories

import (
	"car-service/internal/domain/entities"

	"github.com/google/uuid"
)

type OwnerRepository interface {
	Create(owner *entities.Owner) error
	GetByID(id uuid.UUID) (*entities.Owner, error)
	ExistsByID(id uuid.UUID) (bool, error)
	GetByEmail(email string) (*entities.Owner, error)
	Update(owner *entities.Owner) error
	Delete(id uuid.UUID) error
	List() ([]*entities.Owner, error)
}
