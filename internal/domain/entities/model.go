package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model representa un modelo específico de automóvil
type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"not null"`
	BrandID   uuid.UUID `gorm:"type:uuid;not null"`
	Brand     Brand     `gorm:"foreignKey:BrandID"`
	StartYear int       // Año en que comenzó la producción
	EndYear   int       // Año en que terminó la producción (0 si sigue en producción)
	Category  string    // Categoría del vehículo (SUV, Sedán, etc.)
	Active    bool      `gorm:"default:true"` // Indica si el modelo está activo en el sistema
	Cars      []Car     `gorm:"foreignKey:ModelID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewModel(name string, brandID uuid.UUID, startYear int, category string) *Model {
	return &Model{
		ID:        uuid.New(),
		Name:      name,
		BrandID:   brandID,
		StartYear: startYear,
		Category:  category,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
