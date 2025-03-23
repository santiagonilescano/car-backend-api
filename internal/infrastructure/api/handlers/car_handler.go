// internal/infrastructure/api/handlers/car_handler.go

package handlers

import (
	"car-service/internal/application/commands"
	api "car-service/internal/infrastructure/api/mediator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	mediator *api.Mediator
}

func NewCarHandler(mediator *api.Mediator) *CarHandler {
	return &CarHandler{mediator: mediator}
}

func (h *CarHandler) CreateCar(c *gin.Context) {
	var car commands.NewCarRequest
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar el mediador para enviar el comando
	if err := h.mediator.Send("CreateCar", &car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, car)
}
