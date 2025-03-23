// internal/application/services/car_service_impl.go

package services

import (
	"car-service/internal/domain/entities"
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

func (s *CarServiceImpl) CreateCar(ctx context.Context, car *entities.Car) error {
	return s.repo.Create(ctx, car)
}
