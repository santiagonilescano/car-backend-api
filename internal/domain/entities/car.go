package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Car representa un vehículo específico
type Car struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	ModelID   uuid.UUID `gorm:"type:uuid;not null"`
	Model     Model     `gorm:"foreignKey:ModelID"`
	Year      int       `gorm:"not null"` // Año de fabricación del vehículo específico
	Color     string
	VIN       string    `gorm:"unique"` // Vehicle Identification Number
	OwnerID   uuid.UUID `gorm:"type:uuid"`
	Owner     Owner     `gorm:"foreignKey:OwnerID"`
	Active    bool      `gorm:"default:true"` // Indica si el auto está activo en el sistema
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate se ejecuta antes de crear un nuevo registro
func (c *Car) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func NewCar(modelID uuid.UUID, year int, color string, vin string, ownerID uuid.UUID) *Car {
	return &Car{
		ID:        uuid.New(),
		ModelID:   modelID,
		Year:      year,
		Color:     color,
		VIN:       vin,
		OwnerID:   ownerID,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
