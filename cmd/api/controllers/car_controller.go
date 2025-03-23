// cmd/api/controllers/car_handler.go

package controllers

import (
	api "car-service/cmd/api/mediator"
	"car-service/internal/application/commands"
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
	var car commands.NewCarRequest
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el mediador para enviar el comando
	request := api.CommandRequest[any]{Data: car}
	if _, err := h.mediator.Send("CreateCar", &request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, car)
}

func (h *CarController) GetCars(c *gin.Context) {
	request := api.QueryRequest[any]{}

	resp, err := h.mediator.SendQuery("GetCars", &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
