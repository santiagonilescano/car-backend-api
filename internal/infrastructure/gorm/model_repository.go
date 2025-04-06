package gorm

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelRepository implementa la interfaz repositories.ModelRepository usando GORM
type ModelRepository struct {
	db *gorm.DB
}

// NewModelRepository crea una nueva instancia de ModelRepository
func NewModelRepository(db *gorm.DB) repositories.ModelRepository {
	return &ModelRepository{
		db: db,
	}
}

// Create guarda un nuevo modelo en la base de datos
func (r *ModelRepository) Create(model *entities.Model) error {
	return r.db.Create(model).Error
}

// GetByID obtiene un modelo por su ID
func (r *ModelRepository) GetByID(id uuid.UUID) (*entities.Model, error) {
	var model entities.Model
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// ExistsByID verifica si existe un modelo con el ID proporcionado
func (r *ModelRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Model{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetByBrandID obtiene todos los modelos de una marca
func (r *ModelRepository) GetByBrandID(brandID uuid.UUID) ([]*entities.Model, error) {
	var models []*entities.Model
	err := r.db.Where("brand_id = ?", brandID).Find(&models).Error
	return models, err
}

// GetByNameAndBrand obtiene un modelo por su nombre y marca
func (r *ModelRepository) GetByNameAndBrand(name string, brandID uuid.UUID) (*entities.Model, error) {
	var model entities.Model
	err := r.db.Where("name = ? AND brand_id = ?", name, brandID).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// Update actualiza un modelo existente
func (r *ModelRepository) Update(model *entities.Model) error {
	return r.db.Save(model).Error
}

// Delete elimina un modelo por su ID
func (r *ModelRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Model{}, "id = ?", id).Error
}

// List obtiene todos los modelos
func (r *ModelRepository) List() ([]*entities.Model, error) {
	var models []*entities.Model
	err := r.db.Find(&models).Error
	return models, err
}

// ListActive obtiene todos los modelos activos
func (r *ModelRepository) ListActive() ([]*entities.Model, error) {
	var models []*entities.Model
	err := r.db.Where("active = ?", true).Find(&models).Error
	return models, err
}

// ListByCategory obtiene todos los modelos de una categoría
func (r *ModelRepository) ListByCategory(category string) ([]*entities.Model, error) {
	var models []*entities.Model
	err := r.db.Where("category = ?", category).Find(&models).Error
	return models, err
}
