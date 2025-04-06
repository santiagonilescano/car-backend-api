package gorm

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OwnerRepository implementa la interfaz repositories.OwnerRepository usando GORM
type OwnerRepository struct {
	db *gorm.DB
}

// NewOwnerRepository crea una nueva instancia de OwnerRepository
func NewOwnerRepository(db *gorm.DB) repositories.OwnerRepository {
	return &OwnerRepository{
		db: db,
	}
}

// Create guarda un nuevo propietario en la base de datos
func (r *OwnerRepository) Create(owner *entities.Owner) error {
	return r.db.Create(owner).Error
}

// GetByID obtiene un propietario por su ID
func (r *OwnerRepository) GetByID(id uuid.UUID) (*entities.Owner, error) {
	var owner entities.Owner
	err := r.db.First(&owner, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

// ExistsByID verifica si existe un propietario con el ID proporcionado
func (r *OwnerRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Owner{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetByEmail obtiene un propietario por su email
func (r *OwnerRepository) GetByEmail(email string) (*entities.Owner, error) {
	var owner entities.Owner
	err := r.db.Where("email = ?", email).First(&owner).Error
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

// Update actualiza un propietario existente
func (r *OwnerRepository) Update(owner *entities.Owner) error {
	return r.db.Save(owner).Error
}

// Delete elimina un propietario por su ID
func (r *OwnerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Owner{}, "id = ?", id).Error
}

// List obtiene todos los propietarios
func (r *OwnerRepository) List() ([]*entities.Owner, error) {
	var owners []*entities.Owner
	err := r.db.Find(&owners).Error
	return owners, err
}
