package get_cars

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/domain/services"
	"context"
)

const Name = "GetCars"

type GetCarsQuery struct {
	service services.CarService
}

func NewGetCarsQuery(service services.CarService) *GetCarsQuery {
	return &GetCarsQuery{service: service}
}

func (q *GetCarsQuery) Execute(request api.QueryRequest[any], ctx context.Context) (any, error) {
	return q.service.GetCars(ctx)
}
