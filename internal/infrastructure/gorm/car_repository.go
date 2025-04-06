package gorm

import (
	"car-service/internal/domain/entities"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CarRepository implementa la interfaz repositories.CarRepository usando GORM
type CarRepository struct {
	db *gorm.DB
}

// NewCarRepository crea una nueva instancia de CarRepository
func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

// Create guarda un nuevo auto en la base de datos
func (r *CarRepository) Create(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	return car, r.db.WithContext(ctx).Create(car).Error
}

// GetByID obtiene un auto por su ID
func (r *CarRepository) GetByID(id uuid.UUID) (*entities.Car, error) {
	var car entities.Car
	err := r.db.First(&car, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &car, nil
}

// Update actualiza un auto existente
func (r *CarRepository) Update(car *entities.Car) error {
	return r.db.Save(car).Error
}

// Delete elimina un auto por su ID
func (r *CarRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Car{}, "id = ?", id).Error
}

// GetByOwnerID obtiene todos los autos de un propietario
func (r *CarRepository) GetByOwnerID(ownerID uuid.UUID) ([]*entities.Car, error) {
	var cars []*entities.Car
	err := r.db.Where("owner_id = ?", ownerID).Find(&cars).Error
	return cars, err
}

// List obtiene todos los autos
func (r *CarRepository) List() ([]*entities.Car, error) {
	var cars []*entities.Car
	err := r.db.Find(&cars).Error
	return cars, err
}

// GetByVIN obtiene un auto por su número de VIN
func (r *CarRepository) GetByVIN(vin string) (*entities.Car, error) {
	var car entities.Car
	err := r.db.Where("vin = ?", vin).First(&car).Error
	if err != nil {
		return nil, err
	}
	return &car, nil
}
