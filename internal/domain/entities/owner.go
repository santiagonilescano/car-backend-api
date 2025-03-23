package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Owner struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Phone     string
	Address   string
	Cars      []Car `gorm:"foreignKey:OwnerID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewOwner(name, email, phone, address string) *Owner {
	return &Owner{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
