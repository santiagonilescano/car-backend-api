//internal/application/commands/new_car/new_car_command.go

package new_car

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
	"time"

	"github.com/google/uuid"
)

const Name = "CreateCar"

type NewCarCommand struct {
	service services.CarService
}

func CreateNewCarCommand(service services.CarService) *NewCarCommand {
	return &NewCarCommand{
		service: service,
	}
}

func (c *NewCarCommand) Validate(request api.CommandRequest[any], commandContext *api.CommandContext) []*api.ValidationError {
	var errors []*api.ValidationError
	carRequest := request.Data.(*NewCarRequest)
	if carRequest.ModelId == uuid.Nil {
		errors = append(errors, &api.ValidationError{
			Field:   "modelId",
			Message: "El ID del modelo es requerido",
		})
	}

	if carRequest.OwnerId == uuid.Nil {
		errors = append(errors, &api.ValidationError{
			Field:   "ownerId",
			Message: "El ID del propietario es requerido",
		})
	}

	if carRequest.Vin == "" {
		errors = append(errors, &api.ValidationError{
			Field:   "vin",
			Message: "El VIN es requerido",
		})
	} else if len(carRequest.Vin) != 17 {
		errors = append(errors, &api.ValidationError{
			Field:   "vin",
			Message: "El VIN debe tener 17 caracteres",
		})
	}

	if carRequest.Year < 1900 || carRequest.Year > time.Now().Year()+1 {
		errors = append(errors, &api.ValidationError{
			Field:   "year",
			Message: "El año debe estar entre 1900 y el año siguiente al actual",
		})
	}
	return errors
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
	carResult, err := c.service.CreateCar(*ctx, &car)
	if err != nil {
		return nil, err
	}
	return CreateCarResponse(carResult), nil
}
