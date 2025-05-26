package new_car

import "github.com/google/uuid"

type NewCarRequest struct {
	ModelId uuid.UUID `json:"modelid"`
	OwnerId uuid.UUID `json:"ownerid"`
	Year    int       `json:"year"`
	Color   string    `json:"color"`
	Vin     string    `json:"vin"`
}
