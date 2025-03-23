package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Brand representa una marca de automóvil
type Brand struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"unique;not null"`
	Country   string    // País de origen de la marca
	LogoURL   string    // URL del logo de la marca
	Active    bool      `gorm:"default:true"` // Indica si la marca está activa en el sistema
	Models    []Model   `gorm:"foreignKey:BrandID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewBrand(name, country string, logoURL string) *Brand {
	return &Brand{
		ID:        uuid.New(),
		Name:      name,
		Country:   country,
		LogoURL:   logoURL,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
