// internal/infrastructure/migrations/migration.go

package migrations

import (
	"time"

	"gorm.io/gorm"
)

// SchemaMigration registra las migraciones ejecutadas
type SchemaMigration struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}

// Migration define la interfaz para las migraciones
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}
