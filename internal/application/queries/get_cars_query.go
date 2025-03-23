package queries

import (
	"car-service/internal/domain/services"
	"context"

	api "car-service/internal/infrastructure/api/mediator"
)

type GetCarsQuery struct {
	service services.CarService
}

func NewGetCarsQuery(service services.CarService) *GetCarsQuery {
	return &GetCarsQuery{service: service}
}

func (q *GetCarsQuery) Execute(request api.QueryRequest[any], ctx context.Context) (any, error) {
	return q.service.GetCars(ctx)
}
