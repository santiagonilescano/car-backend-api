package update_car

import (
	"car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const Name = "UpdateCar" // Used for mediator registration

// UpdateCarRequest defines the structure for updating an existing car.
// Fields are pointers to allow partial updates (only provided fields are changed).
type UpdateCarRequest struct {
	ID      uuid.UUID  `json:"-"` // ID comes from the URL path, not the request body
	ModelID *uuid.UUID `json:"modelId,omitempty"`
	OwnerID *uuid.UUID `json:"ownerId,omitempty"`
	Year    *int       `json:"year,omitempty"`
	Color   *string    `json:"color,omitempty"`
	VIN     *string    `json:"vin,omitempty"`
	Active  *bool      `json:"active,omitempty"`
}

// UpdateCarCommandHandler handles the update of an existing car.
// Implements mediator.CommandHandler[UpdateCarRequest, *entities.Car].
type UpdateCarCommandHandler struct {
	service services.CarService
}

// NewUpdateCarCommandHandler creates a new UpdateCarCommandHandler.
func NewUpdateCarCommandHandler(service services.CarService) *UpdateCarCommandHandler {
	return &UpdateCarCommandHandler{
		service: service,
	}
}

// Validate checks the UpdateCarRequest for validation errors.
func (h *UpdateCarCommandHandler) Validate(request UpdateCarRequest, c *gin.Context, cmdCtx *mediator.CommandContext) []*mediator.ValidationError {
	var validationErrors []*mediator.ValidationError

	// ID is already validated by the controller (must be a valid UUID).
	// Here, we ensure it's not Nil, though practically the controller's Parse would fail first for an empty ID.
	if request.ID == uuid.Nil {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "id", // Path parameter
			Message: "El ID del auto es requerido en la URL.",
		})
	}

	// Validate VIN if provided
	if request.VIN != nil {
		if *request.VIN == "" {
			validationErrors = append(validationErrors, &mediator.ValidationError{
				Field:   "vin",
				Message: "El VIN no puede estar vacío si se proporciona.",
			})
		} else if len(*request.VIN) != 17 {
			validationErrors = append(validationErrors, &mediator.ValidationError{
				Field:   "vin",
				Message: "El VIN debe tener 17 caracteres.",
			})
		}
	}

	// Validate Year if provided
	if request.Year != nil {
		if *request.Year < 1900 || *request.Year > time.Now().Year()+1 {
			validationErrors = append(validationErrors, &mediator.ValidationError{
				Field:   "year",
				Message: fmt.Sprintf("El año debe estar entre 1900 y %d.", time.Now().Year()+1),
			})
		}
	}
	
	// Check if at least one field to update is provided in the body
	if request.ModelID == nil && request.OwnerID == nil && request.Year == nil &&
		request.Color == nil && request.VIN == nil && request.Active == nil {
		validationErrors = append(validationErrors, &mediator.ValidationError{
			Field:   "requestBody",
			Message: "Al menos un campo debe ser proporcionado para la actualización en el cuerpo de la solicitud.",
		})
	}

	return validationErrors
}

// Execute processes the UpdateCarRequest to update an existing car.
func (h *UpdateCarCommandHandler) Execute(request UpdateCarRequest, ctx *context.Context) (*entities.Car, error) {
	// Fetch the existing car
	carToUpdate, err := h.service.GetCarByID(*ctx, request.ID)
	if err != nil {
		return nil, err // Handles not found or other errors from the service
	}

	// Apply updates from request if fields are provided (not nil)
	updated := false
	if request.ModelID != nil {
		carToUpdate.ModelID = *request.ModelID
		updated = true
	}
	if request.OwnerID != nil {
		carToUpdate.OwnerID = *request.OwnerID
		updated = true
	}
	if request.Year != nil {
		carToUpdate.Year = *request.Year
		updated = true
	}
	if request.Color != nil {
		carToUpdate.Color = *request.Color
		updated = true
	}
	if request.VIN != nil {
		carToUpdate.VIN = *request.VIN
		updated = true
	}
	if request.Active != nil {
		carToUpdate.Active = *request.Active
		updated = true
	}

	if !updated {
		// No actual fields were updated in the request body.
		// Depending on business logic, could return carToUpdate as is, or an error/specific message.
		// For now, we proceed to the service update call, which might handle this (e.g., by just updating UpdatedAt).
		// Or, if no change means no DB call, could return carToUpdate here.
		// Let's assume service.UpdateCar will correctly handle if no actual data fields changed but UpdatedAt needs update.
	}
	
	// The service's UpdateCar method will handle setting UpdatedAt and other business logic.
	return h.service.UpdateCar(*ctx, carToUpdate)
}
