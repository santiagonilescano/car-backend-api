//internal/application/commands/new_car/new_car_command.go

package new_car

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const Name = "CreateCar"

type NewCarRequest struct {
	ModelId uuid.UUID `json:"modelid"`
	OwnerId uuid.UUID `json:"ownerid"`
	Year    int       `json:"year"`
	Color   string    `json:"color"`
	Vin     string    `json:"vin"`
}

type NewCarCommand struct {
	service services.CarService
}

func NewNewCarCommand(service services.CarService) *NewCarCommand {
	return &NewCarCommand{
		service: service,
	}
}

func (c *NewCarCommand) Validate(ctx *gin.Context, commandContext *api.CommandContext) []*api.ValidationError {
	var errors []*api.ValidationError

	if ctx.Request.Body == nil || ctx.Request.ContentLength == 0 {
		return []*api.ValidationError{
			{Field: "request", Message: "El cuerpo de la solicitud está vacío"},
		}
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		return []*api.ValidationError{
			{Field: "request", Message: "Error al leer el cuerpo de la solicitud: " + err.Error()},
		}
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	log.Printf("Cuerpo de la solicitud: %s", string(body))

	var request NewCarRequest
	if err := json.Unmarshal(body, &request); err != nil {
		log.Printf("Error al deserializar JSON: %v", err)
		return []*api.ValidationError{
			{Field: "request", Message: "Error al procesar la solicitud: " + err.Error() + ". Los campos deben ser: modelid, ownerid, year, color, vin"},
		}
	}

	if request.ModelId == uuid.Nil {
		errors = append(errors, &api.ValidationError{
			Field:   "modelId",
			Message: "El ID del modelo es requerido",
		})
	}

	if request.OwnerId == uuid.Nil {
		errors = append(errors, &api.ValidationError{
			Field:   "ownerId",
			Message: "El ID del propietario es requerido",
		})
	}

	if request.Vin == "" {
		errors = append(errors, &api.ValidationError{
			Field:   "vin",
			Message: "El VIN es requerido",
		})
	} else if len(request.Vin) != 17 {
		errors = append(errors, &api.ValidationError{
			Field:   "vin",
			Message: "El VIN debe tener 17 caracteres",
		})
	}

	if request.Year == 0 {
		request.Year = time.Now().Year()
		commandContext.AddDecision("Se utilizará el año en curso ya que no fue informado")
	}

	if request.Year < 1900 || request.Year > time.Now().Year()+1 {
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
	return c.service.CreateCar(*ctx, &car)
}
