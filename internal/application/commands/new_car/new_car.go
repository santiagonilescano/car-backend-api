//internal/application/commands/new_car/new_car_command.go

package new_car

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"

	"github.com/google/uuid"
)

const Name = "CreateCar"

type NewCarRequest struct {
	ModelId uuid.UUID
	OwnerId uuid.UUID
	Year    int
	Color   string
	Vin     string
}

type NewCarCommand struct {
	service services.CarService
}

func NewNewCarCommand(service services.CarService) *NewCarCommand {
	return &NewCarCommand{
		service: service,
	}
}

func (c *NewCarCommand) Execute(request api.CommandRequest[any], ctx *context.Context) (any, error) {
	carRequest := request.Data.(*NewCarRequest)
	car := entities.Car{
		ModelID: carRequest.ModelId,
		OwnerID: carRequest.OwnerId,
		Year:    carRequest.Year,
		Color:   carRequest.Color,
		VIN:     carRequest.Vin,
	}
	return c.service.CreateCar(*ctx, &car)
}
