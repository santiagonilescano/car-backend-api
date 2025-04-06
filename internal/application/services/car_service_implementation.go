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
	carRepo   repositories.CarRepository
	modelRepo repositories.ModelRepository
	ownerRepo repositories.OwnerRepository
}

func NewCarService(
	carRepo repositories.CarRepository,
	modelRepo repositories.ModelRepository,
	ownerRepo repositories.OwnerRepository,
) services.CarService {
	return &CarServiceImpl{
		carRepo:   carRepo,
		modelRepo: modelRepo,
		ownerRepo: ownerRepo,
	}
}

func (s *CarServiceImpl) CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	existingCar, err := s.carRepo.GetByVIN(car.VIN)
	if err == nil && existingCar != nil {
		return nil, errors.NewBusinessError("DUPLICATE_VIN", "Ya existe un vehículo con este número de VIN")
	}

	modelExists, err := s.modelRepo.ExistsByID(car.ModelID)
	if err != nil {
		return nil, err
	}
	if !modelExists {
		return nil, errors.NewBusinessError("MODEL_NOT_FOUND", "El modelo especificado no existe")
	}

	ownerExists, err := s.ownerRepo.ExistsByID(car.OwnerID)
	if err != nil {
		return nil, err
	}
	if !ownerExists {
		return nil, errors.NewBusinessError("OWNER_NOT_FOUND", "El propietario especificado no existe")
	}

	return s.carRepo.Create(ctx, car)
}

func (s *CarServiceImpl) GetCars(ctx context.Context) ([]*entities.Car, error) {
	return s.carRepo.List()
}
