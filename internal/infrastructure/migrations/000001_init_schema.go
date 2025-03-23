// internal/infrastructure/migrations/000001_init_schema.go

package migrations

import (
	"car-service/internal/domain/entities"

	"gorm.io/gorm"
)

// InitialMigration representa la migración inicial
type InitialMigration struct{}

// Up aplica la migración inicial
func (m *InitialMigration) Up(db *gorm.DB) error {
	// Crear tabla de migraciones si no existe
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return err
	}

	// Crear las tablas en orden según las dependencias
	return db.AutoMigrate(
		&entities.Brand{},
		&entities.Model{},
		&entities.Owner{},
		&entities.Car{},
	)
}

// Down revierte la migración inicial
func (m *InitialMigration) Down(db *gorm.DB) error {
	// Eliminar tablas en orden inverso
	db.Migrator().DropTable(&entities.Car{})
	db.Migrator().DropTable(&entities.Owner{})
	db.Migrator().DropTable(&entities.Model{})
	db.Migrator().DropTable(&entities.Brand{})
	return nil
}
