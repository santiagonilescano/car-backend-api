//internal/application/commands/new_car_command.go

package commands

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
)

type NewCarRequest struct {
}

type NewCarCommand struct {
	service services.CarService
}

func NewNewCarCommand(service services.CarService) *NewCarCommand {
	return &NewCarCommand{
		service: service,
	}
}

func (c *NewCarCommand) Execute(request NewCarRequest) error {
	car := entities.Car{}
	return c.service.CreateCar(context.Background(), &car)
}
