// internal/application/services/car_service_impl.go

package services

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/errors"
	"car-service/internal/domain/repositories"
	"car-service/internal/domain/services"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// It's okay if GetByVIN returns gorm.ErrRecordNotFound, means VIN is unique

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

func (s *CarServiceImpl) GetCarByID(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
	car, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		// Handle gorm.ErrRecordNotFound specifically if you want to return a custom domain error
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError("CAR_NOT_FOUND", "Auto no encontrado") 
		}
		return nil, err
	}
	return car, nil
}

func (s *CarServiceImpl) UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	// The car object passed here is assumed to be the one fetched and then modified by the command handler.
	// Validations for ModelID and OwnerID ensure they exist if they are being changed to a new value.
	// If ModelID/OwnerID are uuid.Nil, it means they are not being updated or are being cleared (if business logic allows).
	// For this implementation, we assume non-nil means it's an intended update to that ID.

	if car.ModelID != uuid.Nil { // Check if ModelID is being set/updated
		modelExists, err := s.modelRepo.ExistsByID(car.ModelID)
		if err != nil {
			return nil, err
		}
		if !modelExists {
			return nil, errors.NewBusinessError("MODEL_NOT_FOUND", "El modelo especificado no existe")
		}
	}

	if car.OwnerID != uuid.Nil { // Check if OwnerID is being set/updated
		ownerExists, err := s.ownerRepo.ExistsByID(car.OwnerID)
		if err != nil {
			return nil, err
		}
		if !ownerExists {
			return nil, errors.NewBusinessError("OWNER_NOT_FOUND", "El propietario especificado no existe")
		}
	}

	// Validate VIN uniqueness if it's being changed
	if car.VIN != "" { // VIN is part of the update
		existingCarWithVin, err := s.carRepo.GetByVIN(car.VIN)
		if err == nil && existingCarWithVin != nil && existingCarWithVin.ID != car.ID {
			return nil, errors.NewBusinessError("DUPLICATE_VIN", "Ya existe otro vehículo con este número de VIN")
		}
		if err != nil && err != gorm.ErrRecordNotFound { // An actual DB error occurred
			return nil, err
		}
		// If gorm.ErrRecordNotFound, VIN is unique, which is good.
	}
	
	car.UpdatedAt = time.Now()
	return s.carRepo.Update(ctx, car)
}
