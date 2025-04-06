// cmd/api/controllers/car_controller.go

package controllers

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/application/commands/new_car"
	"car-service/internal/application/queries/get_cars"

	"github.com/gin-gonic/gin"
)

type CarController struct {
	mediator *api.Mediator
}

func NewCarController(mediator *api.Mediator) *CarController {
	return &CarController{mediator: mediator}
}

func (h *CarController) CreateCar(c *gin.Context) {
	h.mediator.Send(c, api.Command, new_car.Name, new(new_car.NewCarRequest))
}

func (h *CarController) GetCars(c *gin.Context) {
	h.mediator.Send(c, api.Query, get_cars.Name, api.QueryRequest[any]{})
}
