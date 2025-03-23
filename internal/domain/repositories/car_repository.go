package repositories

import (
	"car-service/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

// CarRepository define las operaciones de persistencia para los autos
type CarRepository interface {
	Create(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetByID(id uuid.UUID) (*entities.Car, error)
	Update(car *entities.Car) error
	Delete(id uuid.UUID) error
	GetByOwnerID(ownerID uuid.UUID) ([]*entities.Car, error)
	List() ([]*entities.Car, error)
}
