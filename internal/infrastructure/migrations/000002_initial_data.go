// internal/infrastructure/migrations/000002_initial_data.go

package migrations

import (
	"car-service/internal/domain/entities"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InitialData representa la migración de datos iniciales
type InitialData struct{}

// Up inserta los datos iniciales en la base de datos
func (m *InitialData) Up(db *gorm.DB) error {
	// Verificar si los datos iniciales ya fueron insertados
	var count int64
	db.Model(&SchemaMigration{}).Where("version = ?", "000002_initial_data").Count(&count)
	if count > 0 {
		log.Println("Los datos iniciales ya fueron insertados")
		return nil
	}

	// Crear marcas
	toyotaID := uuid.New()
	hondaID := uuid.New()

	brands := []entities.Brand{
		{ID: toyotaID, Name: "Toyota", Country: "Japan"},
		{ID: hondaID, Name: "Honda", Country: "Japan"},
	}

	if err := db.Create(&brands).Error; err != nil {
		return err
	}

	// Crear modelos
	models := []entities.Model{
		{ID: uuid.New(), Name: "Corolla", StartYear: 2024, BrandID: toyotaID, Category: "Sedan", Active: true},
		{ID: uuid.New(), Name: "Civic", StartYear: 2024, BrandID: hondaID, Category: "Sedan", Active: true},
	}

	if err := db.Create(&models).Error; err != nil {
		return err
	}

	// Crear owner root
	rootOwner := entities.Owner{
		ID:    uuid.New(),
		Name:  "Root Owner",
		Email: "root@example.com",
		Phone: "+1234567890",
	}

	if err := db.Create(&rootOwner).Error; err != nil {
		return err
	}

	// Registrar la migración como ejecutada
	if err := db.Create(&SchemaMigration{
		Version:   "000002_initial_data",
		AppliedAt: time.Now(),
	}).Error; err != nil {
		return err
	}

	log.Println("Datos iniciales insertados correctamente")
	return nil
}

// Down revierte la migración de datos iniciales
func (m *InitialData) Down(db *gorm.DB) error {
	// Eliminar todas las marcas y modelos
	if err := db.Exec("DELETE FROM models").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM brands").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM owners").Error; err != nil {
		return err
	}
	return nil
}
