package get_cars

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services"
	"context"
)

const Name = "GetCars"

// GetCarsRequest can be an empty struct if no parameters are needed for this query.
type GetCarsRequest struct {
	// Add any request parameters here if needed in the future
}

type GetCarsQueryHandler struct {
	service services.CarService
}

func NewGetCarsQueryHandler(service services.CarService) *GetCarsQueryHandler {
	return &GetCarsQueryHandler{service: service}
}

// Execute handles the GetCarsRequest and returns a list of cars.
// It now matches the QueryHandler interface: QueryHandler[GetCarsRequest, []*entities.Car]
func (q *GetCarsQueryHandler) Execute(request GetCarsRequest, ctx context.Context) ([]*entities.Car, error) {
	// request parameter is available if needed, e.g., request.SomeFilter
	return q.service.GetCars(ctx)
}
