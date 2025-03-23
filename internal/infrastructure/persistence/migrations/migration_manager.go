// internal/infrastructure/persistence/migrations/migration_manager.go

package migrations

import (
	"gorm.io/gorm"
)

// Migrate ejecuta todas las migraciones
func Migrate(db *gorm.DB) error {
	migrations := []Migration{
		&InitialMigration{},
		&InitialData{},
	}

	for _, migration := range migrations {
		if err := migration.Up(db); err != nil {
			return err
		}
	}

	return nil
}

// Rollback revierte la última migración
func Rollback(db *gorm.DB) error {
	migrations := []Migration{
		&InitialData{},
		&InitialMigration{},
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		if err := migrations[i].Down(db); err != nil {
			return err
		}
	}

	return nil
}
