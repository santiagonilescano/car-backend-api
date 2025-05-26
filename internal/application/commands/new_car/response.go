package new_car

import (
	"car-service/internal/domain/entities"
	"time"
)

type NewCarResponse struct {
	ID        string    `json:"id"`
	ModelID   string    `json:"modelId"`
	OwnerID   string    `json:"ownerId"`
	Year      int       `json:"year"`
	Color     string    `json:"color"`
	VIN       string    `json:"vin"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateCarResponse(car *entities.Car) *NewCarResponse {
	return &NewCarResponse{
		ID:        car.ID.String(),
		ModelID:   car.ModelID.String(),
		OwnerID:   car.OwnerID.String(),
		Year:      car.Year,
		Color:     car.Color,
		VIN:       car.VIN,
		CreatedAt: car.CreatedAt,
	}
}
