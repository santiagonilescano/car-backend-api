package get_car_by_id

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"

	"github.com/google/uuid"
)

const Name = "GetCarByID" // Used for mediator registration

// GetCarByIDRequest defines the structure for requesting a car by its ID.
type GetCarByIDRequest struct {
	ID uuid.UUID
}

// GetCarByIDQueryHandler handles the GetCarByIDRequest.
// Implements mediator.QueryHandler[GetCarByIDRequest, *entities.Car].
type GetCarByIDQueryHandler struct {
	service services.CarService
}

// NewGetCarByIDQueryHandler creates a new GetCarByIDQueryHandler.
func NewGetCarByIDQueryHandler(service services.CarService) *GetCarByIDQueryHandler {
	return &GetCarByIDQueryHandler{service: service}
}

// Execute processes the GetCarByIDRequest and returns the car if found.
func (h *GetCarByIDQueryHandler) Execute(request GetCarByIDRequest, ctx context.Context) (*entities.Car, error) {
	return h.service.GetCarByID(ctx, request.ID)
}
