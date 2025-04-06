// internal/application/services/car_service_impl.go

package services

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/errors"
	"car-service/internal/domain/repositories"
	"car-service/internal/domain/services"
	"context"
)

type CarServiceImpl struct {
	repo repositories.CarRepository
}

func NewCarService(repo repositories.CarRepository) services.CarService {
	return &CarServiceImpl{repo: repo}
}

func (s *CarServiceImpl) CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	existingCar, err := s.repo.GetByVIN(car.VIN)
	if err == nil && existingCar != nil {
		return nil, errors.NewBusinessError("DUPLICATE_VIN", "Ya existe un vehículo con este número de VIN")
	}

	return s.repo.Create(ctx, car)
}

func (s *CarServiceImpl) GetCars(ctx context.Context) ([]*entities.Car, error) {
	return s.repo.List()
}
