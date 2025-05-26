package services

import (
	"car-service/internal/domain/entities"
	"context"
	"github.com/google/uuid"
)

// CarService define las operaciones disponibles para los autos
type CarService interface {
	CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetCars(ctx context.Context) ([]*entities.Car, error)
	GetCarByID(ctx context.Context, id uuid.UUID) (*entities.Car, error)
	UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
}
