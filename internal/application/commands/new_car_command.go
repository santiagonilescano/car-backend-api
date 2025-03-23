//internal/application/commands/new_car_command.go

package commands

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	api "car-service/internal/infrastructure/api/mediator"
	"context"

	"github.com/google/uuid"
)

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

func (c *NewCarCommand) Execute(request api.CommandRequest[any], ctx context.Context) (any, error) {
	car := entities.Car{
		ModelID: request.Data.(NewCarRequest).ModelId,
		OwnerID: request.Data.(NewCarRequest).OwnerId,
		Year:    request.Data.(NewCarRequest).Year,
		Color:   request.Data.(NewCarRequest).Color,
		VIN:     request.Data.(NewCarRequest).Vin,
	}
	return c.service.CreateCar(context.Background(), &car), nil
}
