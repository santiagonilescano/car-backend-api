// cmd/api/controllers/car_controller.go

package controllers

import (
	api "car-service/cmd/api/mediator"
	"car-service/cmd/api/response"
	"car-service/internal/application/commands/new_car"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CarController struct {
	mediator *api.Mediator
}

func NewCarController(mediator *api.Mediator) *CarController {
	return &CarController{mediator: mediator}
}

func (h *CarController) CreateCar(c *gin.Context) {
	h.mediator.HandleGinRequest(c, new_car.Name, new(new_car.NewCarRequest))
}

func (h *CarController) GetCars(c *gin.Context) {
	request := api.QueryRequest[any]{}
	resp, err := h.mediator.SendQuery("GetCars", &request)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, http.StatusInternalServerError, "Error al crear el coche", []string{err.Error()}, nil, nil)
		return
	}
	response.JSON(c, http.StatusOK, http.StatusOK, "Coche obtenido exitosamente", resp, nil, nil)
}
