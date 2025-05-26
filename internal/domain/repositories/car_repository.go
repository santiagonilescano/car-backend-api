package repositories

import (
	"car-service/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

// CarRepository define las operaciones de persistencia para los autos
type CarRepository interface {
	Create(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Car, error) // Corrected: Added context.Context
	GetByVIN(vin string) (*entities.Car, error) // Assuming this might need context later, but not specified in errors
	Update(ctx context.Context, car *entities.Car) (*entities.Car, error) // Corrected: Added context.Context and changed return type
	Delete(id uuid.UUID) error
	GetByOwnerID(ownerID uuid.UUID) ([]*entities.Car, error) // Assuming this might need context later
	List() ([]*entities.Car, error) // Assuming this might need context later
}
