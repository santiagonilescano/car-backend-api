// cmd/api/controllers/car_controller.go

package controllers

import (
	"car-service/cmd/api/mediator" // Single import for mediator
	"car-service/cmd/api/response"
	"car-service/internal/application/commands/new_car"
	"car-service/internal/application/commands/update_car"
	"car-service/internal/application/queries/get_car_by_id"
	"car-service/internal/application/queries/get_cars"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CarController struct {
	mediator *mediator.Mediator
}

func NewCarController(m *mediator.Mediator) *CarController {
	return &CarController{mediator: m}
}

// CreateCar handles POST /cars
func (h *CarController) CreateCar(c *gin.Context) {
	// Pass a pointer to NewCarRequest so mediator can unmarshal the body into it.
	h.mediator.Send(c, mediator.Command, new_car.Name, new(new_car.NewCarRequest))
}

// GetCars handles GET /cars
func (h *CarController) GetCars(c *gin.Context) {
	// For GetCars, the request model is an empty struct.
	// No body to unmarshal, so we pass an instance of the specific request type.
	// The mediator's Send method is designed to handle this; requestModel won't be unmarshalled if bodyBytes is empty.
	h.mediator.Send(c, mediator.Query, get_cars.Name, get_cars.GetCarsRequest{})
}

// GetCarByID handles GET /cars/:id
func (h *CarController) GetCarByID(c *gin.Context) {
	idParam := c.Param("id")
	carID, err := uuid.Parse(idParam)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID de auto inválido", nil, []string{"El ID proporcionado no es un UUID válido."}, nil)
		return
	}

	request := get_car_by_id.GetCarByIDRequest{ID: carID}
	// Pass the request struct directly. Mediator will use it for reflection.
	h.mediator.Send(c, mediator.Query, get_car_by_id.Name, request)
}

// UpdateCar handles PATCH /cars/:id
func (h *CarController) UpdateCar(c *gin.Context) {
	idParam := c.Param("id")
	carID, err := uuid.Parse(idParam)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID de auto inválido", nil, []string{"El ID proporcionado no es un UUID válido."}, nil)
		return
	}

	// Create an instance of UpdateCarRequest, setting the ID from the path.
	// Other fields (pointers) will be populated by the mediator from the JSON body.
	updateRequest := update_car.UpdateCarRequest{ID: carID}

	// Pass a pointer to updateRequest so json.Unmarshal in mediator can populate it.
	h.mediator.Send(c, mediator.Command, update_car.Name, &updateRequest)
}
