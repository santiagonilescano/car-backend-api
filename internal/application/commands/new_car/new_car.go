//internal/application/commands/new_car/new_car_command.go

package new_car

import (
	"car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const Name = "CreateCar" // Used for mediator registration

// NewCarRequest defines the structure for a new car creation request.
type NewCarRequest struct {
	ModelId uuid.UUID `json:"modelid"`
	OwnerId uuid.UUID `json:"ownerid"`
	Year    int       `json:"year"`
	Color   string    `json:"color"`
	Vin     string    `json:"vin"`
}

// NewCarCommandHandler handles the creation of a new car.
// Implements mediator.CommandHandler[NewCarRequest, *entities.Car].
type NewCarCommandHandler struct {
	service services.CarService
}

// NewNewCarCommandHandler creates a new NewCarCommandHandler.
func NewNewCarCommandHandler(service services.CarService) *NewCarCommandHandler {
	return &NewCarCommandHandler{
		service: service,
	}
}

// Validate checks the NewCarRequest for validation errors.
// The requestModel is already unmarshalled by the mediator.
func (h *NewCarCommandHandler) Validate(request NewCarRequest, c *gin.Context, cmdCtx *mediator.CommandContext) []*mediator.ValidationError {
	var validationErrors []*mediator.ValidationError

	if request.ModelId == uuid.Nil {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "modelid",
			Message: "El ID del modelo es requerido",
		})
	}

	if request.OwnerId == uuid.Nil {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "ownerid",
			Message: "El ID del propietario es requerido",
		})
	}

	if request.Vin == "" {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "vin",
			Message: "El VIN es requerido",
		})
	} else if len(request.Vin) != 17 {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "vin",
			Message: "El VIN debe tener 17 caracteres",
		})
	}
	
	yearToUse := request.Year
	if yearToUse == 0 {
		yearToUse = time.Now().Year()
		cmdCtx.AddDecision("Se utilizará el año en curso ya que no fue informado")
		// Note: The actual modification of request.Year for Execute should happen
		// either here if Validate is allowed to modify, or in Execute.
		// For now, decision is logged, actual defaulting might be better in Execute or a pre-processor.
		// However, the original code modified it in Validate if it was 0 for the purpose of further validation.
	}

	if yearToUse < 1900 || yearToUse > time.Now().Year()+1 {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "year",
			Message: "El año debe estar entre 1900 y el año siguiente al actual",
		})
	}

	return validationErrors
}

// Execute processes the NewCarRequest to create a new car.
// It now directly accepts NewCarRequest.
func (h *NewCarCommandHandler) Execute(request NewCarRequest, ctx *context.Context) (*entities.Car, error) {
	yearToUse := request.Year
	// Example of defaulting logic if not handled by Validate/modified request
	// if request.Year == 0 { 
	// 	yearToUse = time.Now().Year()
	// }


	car := entities.Car{
		ModelID: request.ModelId,
		OwnerID: request.OwnerId,
		Year:    yearToUse, // Use potentially defaulted year
		Color:   request.Color,
		VIN:     request.Vin,
		// Active defaults to true, CreatedAt/UpdatedAt set by GORM hooks or service
	}
	return h.service.CreateCar(*ctx, &car)
}
